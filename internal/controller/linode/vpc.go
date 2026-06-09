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

// VPCReconciler reconciles a VPC object.
type VPCReconciler struct {
	Events <-chan event.GenericEvent
}

// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=vpcs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=vpcs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=vpcs/finalizers,verbs=update

func (r *VPCReconciler) Reconcile(ctx context.Context, obj *v1alpha1.VPC) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	_ = log
	_ = FromContext(ctx)
	return ctrl.Result{}, nil
}

func (r *VPCReconciler) SetupWithManager(mgr ctrl.Manager, api *AtomicClient) error {
	bld := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.VPC{}).
		Named("vpc")
	if r.Events != nil {
		bld.WatchesRawSource(source.Channel(r.Events, &handler.EnqueueRequestForObject{}))
	}
	return bld.Complete(&Reconciler[v1alpha1.VPC, *v1alpha1.VPC]{
		ctl: mgr.GetClient(),
		api: api,
		obj: r,
	})
}
