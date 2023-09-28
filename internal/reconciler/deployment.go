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

type deploymentReconciler struct {
	client       client.Client
	scheme       *runtime.Scheme
	reconcilable resource.Reconcilable
}

func NewDeploymentReconciler(client client.Client, scheme *runtime.Scheme, reconcilable resource.Reconcilable) Reconciler {
	return &deploymentReconciler{client: client, scheme: scheme, reconcilable: reconcilable}
}

func (dr deploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request, carrier *spacecraftv1.Spacecraft) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	newDeployment := dr.reconcilable.MakeResource(carrier)
	oldDeployment := dr.reconcilable.MakeDefaultResource()
	if err := dr.client.Get(ctx, req.NamespacedName, oldDeployment); err != nil {
		if errors.IsNotFound(err) {
			if err := dr.client.Create(ctx, newDeployment); err != nil {
				logger.Error(err, "Create deployment failed", "deploymentName", newDeployment.GetName())
				return ctrl.Result{}, err
			}
		}
	} else if oldDeployment != newDeployment {
		if err := dr.client.Update(ctx, newDeployment); err != nil {
			logger.Error(err, "Update deployment failed", "deploymentName", newDeployment.GetName())
			return ctrl.Result{}, err
		}
	}

	// log.Log.Info("Get oldDeployment from cluster", "spaceCraft", oldDeployment)

	return ctrl.Result{}, nil
}
