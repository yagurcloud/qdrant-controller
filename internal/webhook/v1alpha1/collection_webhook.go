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

package v1alpha1

import (
	"context"

	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/qdrant/go-client/qdrant"
	qdrantv1alpha1 "github.com/yagurcloud/qdrant-controller/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// nolint:unused
// log is for logging in this package.
var collectionlog = logf.Log.WithName("collection-resource")

// SetupCollectionWebhookWithManager registers the webhook for Collection in the manager.
func SetupCollectionWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &qdrantv1alpha1.Collection{}).
		WithValidator(&CollectionCustomValidator{}).
		WithDefaulter(&CollectionCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-qdrant-yagur-io-v1alpha1-collection,mutating=true,failurePolicy=fail,sideEffects=None,groups=qdrant.yagur.io,resources=collections,verbs=create;update,versions=v1alpha1,name=mcollection-v1alpha1.kb.io,admissionReviewVersions=v1

// CollectionCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Collection when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type CollectionCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Collection.
func (d *CollectionCustomDefaulter) Default(_ context.Context, obj *qdrantv1alpha1.Collection) error {
	collectionlog.Info("Defaulting for Collection", "name", obj.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-qdrant-yagur-io-v1alpha1-collection,mutating=false,failurePolicy=fail,sideEffects=None,groups=qdrant.yagur.io,resources=collections,verbs=create;update,versions=v1alpha1,name=vcollection-v1alpha1.kb.io,admissionReviewVersions=v1

// CollectionCustomValidator struct is responsible for validating the Collection resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type CollectionCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

func (v *CollectionCustomValidator) ValidateCreate(_ context.Context, obj *qdrantv1alpha1.Collection) (admission.Warnings, error) {
	collectionlog.Info("Validation for Collection upon creation", "name", obj.GetName())

	var allErrs field.ErrorList

	if obj.Spec.VectorSize <= 0 {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec").Child("vectorSize"),
			obj.Spec.VectorSize,
			"must be greater than 0",
		))
	}

	_, ok := qdrant.Distance_value[obj.Spec.Distance]
	if !ok {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec").Child("distance"),
			obj.Spec.Distance,
			"must be one of: Cosine, Euclid, Dot, Manhattan",
		))
	}

	if len(allErrs) > 0 {
		return nil,
			apierrors.NewInvalid(
				obj.GroupVersionKind().GroupKind(),
				obj.Name,
				allErrs,
			)
	}

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Collection.
func (v *CollectionCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj *qdrantv1alpha1.Collection) (admission.Warnings, error) {
	collectionlog.Info("Validation for Collection upon update", "name", newObj.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Collection.
func (v *CollectionCustomValidator) ValidateDelete(_ context.Context, obj *qdrantv1alpha1.Collection) (admission.Warnings, error) {
	collectionlog.Info("Validation for Collection upon deletion", "name", obj.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
