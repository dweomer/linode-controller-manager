package linode

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/linode/linodego/v2"

	"github.com/dweomer/linode-controller-manager/api/v1alpha1"
)

var (
	_ linodego.Instance
)

// InstanceReconciler reconciles an Instance object
type InstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=instances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=linode.com,namespace=linode-system,resources=instances/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)
	_ = r.Get(nil, req.NamespacedName, nil)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Instance{}).
		Named("instance").
		Complete(r)
}
