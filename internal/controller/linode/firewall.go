package linode

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/dweomer/linode-controller-manager/api/v1alpha1"
)

// FirewallReconciler reconciles a Firewall object
type FirewallReconciler struct {
	client.Client
	Events <-chan event.GenericEvent
}

// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=firewalls,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=firewalls/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=firewalls/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *FirewallReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)
	_ = r.Get(nil, req.NamespacedName, nil)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FirewallReconciler) SetupWithManager(mgr ctrl.Manager) error {
	bld := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Firewall{}).
		Named("firewall")
	if r.Events != nil {
		bld.WatchesRawSource(source.Channel(r.Events, &handler.EnqueueRequestForObject{}))
	}
	return bld.Complete(r)
}
