/*
Copyright 2026.

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
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/qdrant/go-client/qdrant"
	qdrantv1alpha1 "github.com/yagurcloud/yagur-controllers/api/v1alpha1"
)

// CollectionReconciler reconciles a Collection object
type CollectionReconciler struct {
	client.Client
	Scheme       *runtime.Scheme
	QdrantClient *qdrant.Client
	RequeueAfter time.Duration
}

// +kubebuilder:rbac:groups=qdrant.yagur.io,resources=collections,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=qdrant.yagur.io,resources=collections/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=qdrant.yagur.io,resources=collections/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Collection object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.23.3/pkg/reconcile
func (r *CollectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// 1. Fetch the Collection CR
	collection := &qdrantv1alpha1.Collection{}
	if err := r.Get(ctx, req.NamespacedName, collection); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Ask Qdrant - does this collection exist
	exists, err := r.QdrantClient.CollectionExists(ctx, collection.Name)
	if err != nil {
		log.Error(err, "Failed to check if collection exists")
		return ctrl.Result{}, err
	}

	// 3. Create if missing
	if !exists {
		log.Info("Creating collection", "name", collection.Name)
		req, err := collection.Spec.ToQdrantParams(collection.Name)
		if err != nil {
			return ctrl.Result{}, err
		}
		err = r.QdrantClient.CreateCollection(ctx, req)
		if err != nil {
			log.Error(err, "Failed to create collection")
			return ctrl.Result{}, err
		}
	}

	// 4. Update status with observed state
	actual, err := r.QdrantClient.GetCollectionInfo(ctx, collection.Name)
	if err != nil {
		return ctrl.Result{}, err
	}
	collection.Status.Status = actual.Status.String()
	if err := r.Status().Update(ctx, collection); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: r.RequeueAfter}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CollectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&qdrantv1alpha1.Collection{}).
		Named("collection").
		Complete(r)
}
