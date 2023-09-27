/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	xiaofengv1 "github.com/xiaofengzs/spacecraft/api/v1"
	"github.com/xiaofengzs/spacecraft/internal/reconciler"
	"github.com/xiaofengzs/spacecraft/internal/resource"
)

// SpacecraftReconciler reconciles a Spacecraft object
type SpacecraftReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	reconcilers []reconciler.Reconciler
}

func NewSpacecraftReconciler(c client.Client, s *runtime.Scheme) *SpacecraftReconciler {
	var subReconcilers []reconciler.Reconciler
	dr := reconciler.NewDeploymentReconciler(c, s, resource.NewDeployment())
	sr := reconciler.NewServiceReconciler(c, s, resource.NewService())
	subReconcilers = append(subReconcilers, dr, sr)
	return &SpacecraftReconciler{Client: c, Scheme: s, reconcilers: subReconcilers}
}

//+kubebuilder:rbac:groups=xiaofeng.cloud,resources=spacecrafts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=xiaofeng.cloud,resources=spacecrafts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=xiaofeng.cloud,resources=spacecrafts/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps/v1,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Spacecraft object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *SpacecraftReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	log.Log.Info("receive request", "key", req.NamespacedName)

	spaceCraft := &xiaofengv1.Spacecraft{}

	if err := r.Client.Get(context.TODO(), req.NamespacedName, spaceCraft); err != nil {
		if errors.IsNotFound(err) {
			log.Log.Info("This space craft spec has beed delete", "space craft", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		log.Log.Error(err, "Get spaceCraft error")
		return ctrl.Result{}, err
	}

	for _, reconciler := range r.reconcilers {
		if result, err := reconciler.Reconcile(context.TODO(), req, spaceCraft); err != nil {
			log.Log.Error(err, "reconcile sub resource failed")
			return ctrl.Result{}, err
		} else {
			log.Log.Info("receive request", "key", result)
		}
	}

	log.Log.Info("Get space craft from cluster", "spaceCraft", spaceCraft)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpacecraftReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&xiaofengv1.Spacecraft{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
