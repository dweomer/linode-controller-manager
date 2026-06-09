package linode

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/linode/linodego/v2"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	SecretName = "linode-token"
	SecretKey  = "token"
)

type Managed interface {
	SetupWithManager(ctrl.Manager) error
}

// linodeReconciler watches the Linode token Secret and maintains a validated
// linodego client reference that other reconcilers consume via api.
type linodeReconciler struct {
	ctl client.Client
	api *AtomicClient
}

func NewReconciler(ctl client.Client, api *AtomicClient) Managed {
	return &linodeReconciler{
		ctl: ctl, api: api,
	}
}

// +kubebuilder:rbac:groups="",namespace=linode-system,resources=secrets,verbs=get;list;watch

func (r *linodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	var secret corev1.Secret
	if err := r.ctl.Get(ctx, req.NamespacedName, &secret); err != nil {
		if client.IgnoreNotFound(err) == nil {
			log.Info("Token secret deleted, clearing Linode client")
			r.api.Store(nil)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	token, ok := secret.Data[SecretKey]
	if !ok {
		log.Error(fmt.Errorf("missing key %q", SecretKey), "Invalid token secret")
		r.api.Store(nil)
		return ctrl.Result{}, nil
	}

	lc, err := linodego.NewClient(&http.Client{})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("creating linodego client: %w", err)
	}
	lc.SetToken(strings.TrimSpace(string(token)))

	// Validate connectivity.
	if _, err := lc.GetAccount(ctx); err != nil {
		log.Error(err, "Linode API connectivity check failed")
		r.api.Store(nil)
		return ctrl.Result{}, err
	}

	log.Info("Linode API client validated")
	r.api.Store(&lc)

	return ctrl.Result{}, nil
}

func (r *linodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Secret{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				return e.Object.GetName() == SecretName
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				return e.ObjectNew.GetName() == SecretName
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				return e.Object.GetName() == SecretName
			},
			GenericFunc: func(e event.GenericEvent) bool {
				return e.Object.GetName() == SecretName
			},
		}).
		Named(SecretName).
		Complete(r)
}
