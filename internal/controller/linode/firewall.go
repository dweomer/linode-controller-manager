package linode

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/dweomer/linode-controller-manager/api/v1alpha1"
)

// FirewallReconciler reconciles a Firewall object.
type FirewallReconciler struct {
	Events <-chan event.GenericEvent
}

// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=firewalls,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=firewalls/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=firewalls/finalizers,verbs=update

func (r *FirewallReconciler) Reconcile(ctx context.Context, obj *v1alpha1.Firewall) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	_ = log
	_ = FromContext(ctx)
	return ctrl.Result{}, nil
}

func (r *FirewallReconciler) SetupWithManager(mgr ctrl.Manager, api *AtomicClient) error {
	bld := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Firewall{}).
		Named("firewall")
	if r.Events != nil {
		bld.WatchesRawSource(source.Channel(r.Events, &handler.EnqueueRequestForObject{}))
	}
	return bld.Complete(&Reconciler[v1alpha1.Firewall, *v1alpha1.Firewall]{
		ctl: mgr.GetClient(),
		api: api,
		obj: r,
	})
}
