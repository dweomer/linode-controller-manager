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

// InstanceReconciler reconciles an Instance object.
// The linodego client is available via linode.FromContext(ctx).
type InstanceReconciler struct {
	Events <-chan event.GenericEvent
}

// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=instances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=instances/finalizers,verbs=update

func (r *InstanceReconciler) Reconcile(ctx context.Context, obj *v1alpha1.Instance) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	_ = log
	_ = FromContext(ctx) // linodego client
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager, api *AtomicClient) error {
	bld := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Instance{}).
		Named("instance")
	if r.Events != nil {
		bld.WatchesRawSource(source.Channel(r.Events, &handler.EnqueueRequestForObject{}))
	}
	return bld.Complete(&Reconciler[v1alpha1.Instance, *v1alpha1.Instance]{
		ctl: mgr.GetClient(),
		api: api,
		obj: r,
	})
}
