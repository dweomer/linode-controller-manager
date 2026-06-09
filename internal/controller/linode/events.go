package linode

import (
	"context"
	"fmt"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/linode/linodego/v2"

	"github.com/dweomer/linode-controller-manager/api/v1alpha1"
)

// entityMapping pairs a lookup function with a channel for a managed resource type.
type entityMapping struct {
	channel chan event.GenericEvent
	lookup  func(ctx context.Context, c client.Client, id int64) (client.Object, error)
}

// mapEntity builds an entityMapping using generics. Go infers the pointer type
// parameters (PT, PL) via constraint type inference from the items function.
func mapEntity[T any, PT interface {
	*T
	client.Object
}, L any, PL interface {
	*L
	client.ObjectList
}](
	ch chan event.GenericEvent,
	items func(PL) []T,
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
	}
}

// EventPoller polls Linode account events and fans them out as:
//   - native K8s Events on managed resources (observability)
//   - GenericEvents into per-type channels (reconciliation triggers)
type EventPoller struct {
	client.Client
	Linode   *linodego.Client
	Recorder record.EventRecorder
	Interval time.Duration

	entities map[linodego.EntityType]entityMapping
}

// NewEventPoller creates an EventPoller with entity-to-resource mappings for all
// managed types. Reconcilers subscribe to the returned channels via source.Channel.
func NewEventPoller(c client.Client, linode *linodego.Client, recorder record.EventRecorder, interval time.Duration) *EventPoller {
	p := &EventPoller{
		Client:   c,
		Linode:   linode,
		Recorder: recorder,
		Interval: interval,
		entities: make(map[linodego.EntityType]entityMapping),
	}
	p.entities[linodego.EntityLinode] = mapEntity(
		make(chan event.GenericEvent, 64),
		func(l *v1alpha1.InstanceList) []v1alpha1.Instance { return l.Items },
	)
	p.entities[linodego.EntityFirewall] = mapEntity(
		make(chan event.GenericEvent, 64),
		func(l *v1alpha1.FirewallList) []v1alpha1.Firewall { return l.Items },
	)
	p.entities[linodego.EntityVPC] = mapEntity(
		make(chan event.GenericEvent, 64),
		func(l *v1alpha1.VPCList) []v1alpha1.VPC { return l.Items },
	)
	return p
}

// Channel returns the GenericEvent channel for a Linode entity type.
// Reconcilers use this to subscribe via source.Channel.
func (p *EventPoller) Channel(entityType linodego.EntityType) <-chan event.GenericEvent {
	if m, ok := p.entities[entityType]; ok {
		return m.channel
	}
	return nil
}

// +kubebuilder:rbac:groups="",namespace=linode-system,resources=events,verbs=create;patch

// Start implements manager.Runnable. It blocks until the context is cancelled.
func (p *EventPoller) Start(ctx context.Context) error {
	log := logf.FromContext(ctx).WithName("event-poller")
	log.Info("Starting Linode event poller", "interval", p.Interval)

	// Start from now — informer-based reconciliation handles pre-existing state.
	lastEventID := 0

	ticker := time.NewTicker(p.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping Linode event poller")
			return nil
		case <-ticker.C:
			newHighWater, err := p.poll(ctx, lastEventID)
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

func (p *EventPoller) poll(ctx context.Context, lastEventID int) (int, error) {
	log := logf.FromContext(ctx).WithName("event-poller")

	filter := fmt.Sprintf(`{"id": {"+gt": %d}}`, lastEventID)
	opts := linodego.ListOptions{Filter: filter}

	events, err := p.Linode.ListEvents(ctx, &opts)
	if err != nil {
		return lastEventID, err
	}

	highWater := lastEventID
	// Events come newest-first; process oldest-first.
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

		obj, err := mapping.lookup(ctx, p.Client, int64(ev.Entity.ID))
		if err != nil {
			log.Error(err, "Failed to look up entity", "entityType", ev.Entity.Type, "entityID", ev.Entity.ID)
			continue
		}
		if obj == nil {
			continue
		}

		eventType := corev1.EventTypeNormal
		if ev.Status == linodego.EventFailed {
			eventType = corev1.EventTypeWarning
		}
		p.Recorder.Event(obj, eventType, string(ev.Action), ev.Message)

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

// NeedsLeaderElection implements manager.LeaderElectionRunnable.
func (p *EventPoller) NeedsLeaderElection() bool {
	return true
}

// SetupWithManager adds the event poller as a manager runnable.
func (p *EventPoller) SetupWithManager(mgr ctrl.Manager) error {
	return mgr.Add(p)
}
