package linode

import (
	"context"
	"fmt"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/events"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/linode/linodego/v2"

	"github.com/dweomer/linode-controller-manager/api/v1alpha1"
)

// entityMapping pairs lookup/create functions with a channel for a managed resource type.
type entityMapping struct {
	channel chan event.GenericEvent
	lookup  func(ctx context.Context, c client.Client, id int64) (client.Object, error)
	stub    func(ns string, id int64, label string) client.Object
}

// mapEntity builds an entityMapping using generics.
func mapEntity[T any, PT interface {
	*T
	client.Object
}, L any, PL interface {
	*L
	client.ObjectList
}](
	ch chan event.GenericEvent,
	items func(PL) []T,
	stub func(ns string, id int64, label string) PT,
) entityMapping {
	return entityMapping{
		channel: ch,
		lookup: func(ctx context.Context, c client.Client, id int64) (client.Object, error) {
			var list L
			pl := PL(&list)
			if err := c.List(ctx, pl, client.MatchingFields{".spec.id": strconv.FormatInt(id, 10)}); err != nil {
				return nil, err
			}
			if results := items(pl); len(results) > 0 {
				return PT(&results[0]), nil
			}
			return nil, nil
		},
		stub: func(ns string, id int64, label string) client.Object {
			return stub(ns, id, label)
		},
	}
}

// EventPoller polls Linode account events and fans them out as:
//   - native K8s Events on managed resources (observability)
//   - GenericEvents into per-type channels (reconciliation triggers)
//   - stub CRs for entities discovered via Linode events with no existing CR
type EventPoller struct {
	ctl client.Client
	api *AtomicClient
	rec events.EventRecorder
	dur time.Duration

	ns string

	entities map[linodego.EntityType]entityMapping
}

// NewEventPoller creates an EventPoller with entity-to-resource mappings for all
// managed types. Reconcilers subscribe to the returned channels via source.Channel.
func NewEventPoller(ctlClient client.Client, apiClient *AtomicClient, recorder events.EventRecorder, interval time.Duration, namespace string) *EventPoller {
	p := &EventPoller{
		ctl:      ctlClient,
		api:      apiClient,
		rec:      recorder,
		dur:      interval,
		ns:       namespace,
		entities: make(map[linodego.EntityType]entityMapping),
	}
	p.entities[linodego.EntityLinode] = mapEntity(
		make(chan event.GenericEvent, 64),
		func(l *v1alpha1.InstanceList) []v1alpha1.Instance { return l.Items },
		newInstanceStub,
	)
	p.entities[linodego.EntityFirewall] = mapEntity(
		make(chan event.GenericEvent, 64),
		func(l *v1alpha1.FirewallList) []v1alpha1.Firewall { return l.Items },
		newFirewallStub,
	)
	p.entities[linodego.EntityVPC] = mapEntity(
		make(chan event.GenericEvent, 64),
		func(l *v1alpha1.VPCList) []v1alpha1.VPC { return l.Items },
		newVPCStub,
	)
	return p
}

func newInstanceStub(ns string, id int64, label string) *v1alpha1.Instance {
	return &v1alpha1.Instance{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "instance-",
			Namespace:    ns,
		},
		Spec: v1alpha1.InstanceSpec{
			ID:    id,
			Label: label,
		},
	}
}

func newFirewallStub(ns string, id int64, label string) *v1alpha1.Firewall {
	return &v1alpha1.Firewall{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "firewall-",
			Namespace:    ns,
		},
		Spec: v1alpha1.FirewallSpec{
			ID:    id,
			Label: label,
		},
	}
}

func newVPCStub(ns string, id int64, label string) *v1alpha1.VPC {
	return &v1alpha1.VPC{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "vpc-",
			Namespace:    ns,
		},
		Spec: v1alpha1.VPCSpec{
			ID:    id,
			Label: label,
		},
	}
}

// Channel returns the GenericEvent channel for a Linode entity type.
func (p *EventPoller) Channel(entityType linodego.EntityType) <-chan event.GenericEvent {
	if m, ok := p.entities[entityType]; ok {
		return m.channel
	}
	return nil
}

// +kubebuilder:rbac:groups="",namespace=linode-system,resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=instances;firewalls;vpcs,verbs=create

// Start implements manager.Runnable. It blocks until the context is cancelled.
func (p *EventPoller) Start(ctx context.Context) error {
	log := logf.FromContext(ctx).WithName("event-poller")
	log.Info("Starting Linode event poller", "interval", p.dur)

	lastEventID := 0

	ticker := time.NewTicker(p.dur)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping Linode event poller")
			return nil
		case <-ticker.C:
			lc := p.api.Load()
			if lc == nil {
				log.V(1).Info("Linode client not ready, skipping poll")
				continue
			}
			newHighWater, err := p.poll(ctx, lc, lastEventID)
			if err != nil {
				log.Error(err, "Failed to poll Linode events")
				continue
			}
			if newHighWater > lastEventID {
				lastEventID = newHighWater
			}
		}
	}
}

func (p *EventPoller) poll(ctx context.Context, lc *linodego.Client, lastEventID int) (int, error) {
	log := logf.FromContext(ctx).WithName("event-poller")

	filter := fmt.Sprintf(`{"id": {"+gt": %d}}`, lastEventID)
	opts := linodego.ListOptions{Filter: filter}

	events, err := lc.ListEvents(ctx, &opts)
	if err != nil {
		return lastEventID, err
	}

	highWater := lastEventID
	for i := len(events) - 1; i >= 0; i-- {
		ev := events[i]
		if ev.ID > highWater {
			highWater = ev.ID
		}
		if ev.Entity == nil {
			continue
		}

		mapping, ok := p.entities[ev.Entity.Type]
		if !ok {
			continue
		}

		obj, err := mapping.lookup(ctx, p.ctl, int64(ev.Entity.ID))
		if err != nil {
			log.Error(err, "Failed to look up entity", "entityType", ev.Entity.Type, "entityID", ev.Entity.ID)
			continue
		}

		if obj == nil {
			obj = mapping.stub(p.ns, int64(ev.Entity.ID), ev.Entity.Label)
			setEventAnnotations(obj, ev)
			if err := p.ctl.Create(ctx, obj); err != nil {
				log.Error(err, "Failed to create stub for discovered entity",
					"entityType", ev.Entity.Type, "entityID", ev.Entity.ID, "label", ev.Entity.Label)
				continue
			}
			log.Info("Created stub for discovered entity",
				"entityType", ev.Entity.Type, "entityID", ev.Entity.ID, "object", client.ObjectKeyFromObject(obj))
		} else {
			setEventAnnotations(obj, ev)
			if err := p.ctl.Update(ctx, obj); err != nil {
				log.Error(err, "Failed to annotate entity with event",
					"entityType", ev.Entity.Type, "object", client.ObjectKeyFromObject(obj))
			}
		}

		eventType := corev1.EventTypeNormal
		if ev.Status == linodego.EventFailed {
			eventType = corev1.EventTypeWarning
		}
		p.rec.Eventf(obj, nil, eventType,
			string(ev.Action),
			fmt.Sprintf("LinodeEvent/%d", ev.ID),
			"%s %s", ev.Action, ev.Status)

		select {
		case mapping.channel <- event.GenericEvent{Object: obj}:
		default:
			log.V(1).Info("Channel full, dropping reconcile trigger",
				"entityType", ev.Entity.Type,
				"object", types.NamespacedName{Namespace: obj.GetNamespace(), Name: obj.GetName()})
		}
	}

	return highWater, nil
}

func setEventAnnotations(obj client.Object, ev linodego.Event) {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations[AnnotationEventAction] = string(ev.Action)
	annotations[AnnotationEventID] = strconv.Itoa(ev.ID)
	obj.SetAnnotations(annotations)
}

// NeedsLeaderElection implements manager.LeaderElectionRunnable.
func (p *EventPoller) NeedsLeaderElection() bool {
	return true
}

// SetupWithManager adds the event poller as a manager runnable.
func (p *EventPoller) SetupWithManager(mgr ctrl.Manager) error {
	return mgr.Add(p)
}
