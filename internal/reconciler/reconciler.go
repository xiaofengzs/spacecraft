package reconciler

import (
	"context"
	spacecraftv1 "github.com/xiaofengzs/spacecraft/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Reconciler interface {
	Reconcile(ctx context.Context, req ctrl.Request, carrier *spacecraftv1.Spacecraft) (ctrl.Result, error)
}
