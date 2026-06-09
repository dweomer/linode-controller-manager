package linode

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/linode/linodego/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	AnnotationEventAction = "linode.com/event-action"
	AnnotationEventID     = "linode.com/event-id"
)

type AtomicClient = atomic.Pointer[linodego.Client]

// Reconciler wraps a reconcile.ObjectReconciler into a standard reconcile.Reconciler.
// It injects the linodego.Client into the context, fetches the typed object,
// and delegates to the inner obj reconciler.
type Reconciler[T any, P interface {
	*T
	client.Object
}] struct {
	ctl client.Client
	api *AtomicClient
	obj reconcile.ObjectReconciler[P]
}

func (r *Reconciler[T, P]) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	api := r.api.Load()
	if api == nil {
		log.V(1).Info("linode client not ready")
		return ctrl.Result{RequeueAfter: 37 * time.Second}, nil
	}
	ctx = NewContext(ctx, api)
	var (
		obj T
		ptr = P(&obj)
	)
	if err := r.ctl.Get(ctx, req.NamespacedName, ptr); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	return r.obj.Reconcile(ctx, ptr)
}
