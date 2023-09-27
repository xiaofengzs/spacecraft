package reconciler

import (
	"context"
	spacecraftv1 "github.com/xiaofengzs/spacecraft/api/v1"
	"github.com/xiaofengzs/spacecraft/internal/resource"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type serviceReconciler struct {
	client       client.Client
	scheme       *runtime.Scheme
	reconcilable resource.Reconcilable
}

func NewServiceReconciler(client client.Client, scheme *runtime.Scheme, reconcilable resource.Reconcilable) Reconciler {
	return &serviceReconciler{client: client, scheme: scheme, reconcilable: reconcilable}
}

func (sr serviceReconciler) Reconcile(ctx context.Context, req ctrl.Request, carrier *spacecraftv1.Spacecraft) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	if err := sr.reconcilable.Validate(carrier); err != nil {
		return ctrl.Result{}, err
	}

	newService := sr.reconcilable.MakeResource(carrier)
	foundService := sr.reconcilable.MakeDefaultResource()

	// enable service, service not found , create
	// enable service, service found, update
	// disable service, service found, delete
	// disable service, service not found, do nothing

	if err := sr.client.Get(ctx, req.NamespacedName, foundService); err != nil {
		if errors.IsNotFound(err) && carrier.Spec.EnableService {
			if err := sr.client.Create(ctx, newService); err != nil {
				logger.Error(err, "Create servcie failed", "service", newService)
				return ctrl.Result{}, err
			}
		}

		if errors.IsNotFound(err) && !carrier.Spec.EnableService {
			return ctrl.Result{}, nil
		}

		if !errors.IsNotFound(err) && carrier.Spec.EnableService {
			logger.Error(err, "Get service failed", "err", err)
			return ctrl.Result{}, err
		}
	}

	if newService != foundService && carrier.Spec.EnableService {
		if err := sr.client.Update(ctx, newService); err != nil {
			logger.Error(err, "Update service failed", "err", err)
			return ctrl.Result{}, err
		}
	}

	if !carrier.Spec.EnableService {
		if err := sr.client.Delete(ctx, foundService); err != nil {
			logger.Error(err, "Delete service failed", "err", err)
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}
