/*
Copyright 2019 The Crossplane Authors.

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

package container

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/source"

	runtimev1alpha1 "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
	"github.com/crossplaneio/crossplane-runtime/pkg/reconciler/claimbinding"
	"github.com/crossplaneio/crossplane-runtime/pkg/reconciler/claimdefaulting"
	"github.com/crossplaneio/crossplane-runtime/pkg/reconciler/claimscheduling"
	"github.com/crossplaneio/crossplane-runtime/pkg/resource"
	computev1alpha1 "github.com/crossplaneio/crossplane/apis/compute/v1alpha1"

	"github.com/crossplaneio/stack-existing-cluster/apis/container/v1beta1"
)

// A ExistingClusterClaimSchedulingController reconciles KubernetesCluster claims
// that include a class selector but omit their class and resource references by
// picking a random matching GKEClusterClass, if any.
type ExistingClusterClaimSchedulingController struct{}

// SetupWithManager sets up the ExistingClusterClaimSchedulingController using the
// supplied manager.
func (c *ExistingClusterClaimSchedulingController) SetupWithManager(mgr ctrl.Manager) error {
	name := strings.ToLower(fmt.Sprintf("scheduler.%s.%s.%s",
		computev1alpha1.KubernetesClusterKind,
		v1beta1.ExistingClusterKind,
		v1beta1.Group))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&computev1alpha1.KubernetesCluster{}).
		WithEventFilter(resource.NewPredicates(resource.AllOf(
			resource.HasClassSelector(),
			resource.HasNoClassReference(),
			resource.HasNoManagedResourceReference(),
		))).
		Complete(claimscheduling.NewReconciler(mgr,
			resource.ClaimKind(computev1alpha1.KubernetesClusterGroupVersionKind),
			resource.ClassKind(v1beta1.ExistingClusterClassGroupVersionKind),
		))
}

// A ExistingClusterClaimDefaultingController reconciles KubernetesCluster claims
// that omit their resource ref, class ref, and class selector by choosing a
// default GKEClusterClass if one exists.
type ExistingClusterClaimDefaultingController struct{}

// SetupWithManager sets up the ExistingClusterClaimDefaultingController using the
// supplied manager.
func (c *ExistingClusterClaimDefaultingController) SetupWithManager(mgr ctrl.Manager) error {
	name := strings.ToLower(fmt.Sprintf("defaulter.%s.%s.%s",
		computev1alpha1.KubernetesClusterKind,
		v1beta1.ExistingClusterKind,
		v1beta1.Group))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&computev1alpha1.KubernetesCluster{}).
		WithEventFilter(resource.NewPredicates(resource.AllOf(
			resource.HasNoClassSelector(),
			resource.HasNoClassReference(),
			resource.HasNoManagedResourceReference(),
		))).
		Complete(claimdefaulting.NewReconciler(mgr,
			resource.ClaimKind(computev1alpha1.KubernetesClusterGroupVersionKind),
			resource.ClassKind(v1beta1.ExistingClusterClassGroupVersionKind),
		))
}

// A ExistingClusterClaimController reconciles KubernetesCluster claims with
// GKEClusters, dynamically provisioning them if needed.
type ExistingClusterClaimController struct{}

// SetupWithManager adds a controller that reconciles KubernetesCluster resource claims.
func (c *ExistingClusterClaimController) SetupWithManager(mgr ctrl.Manager) error {
	name := strings.ToLower(fmt.Sprintf("%s.%s.%s",
		computev1alpha1.KubernetesClusterKind,
		v1beta1.ExistingClusterClassKind,
		v1beta1.Group))

	p := resource.NewPredicates(resource.AnyOf(
		resource.HasClassReferenceKind(resource.ClassKind(v1beta1.ExistingClusterClassGroupVersionKind)),
		resource.HasManagedResourceReferenceKind(resource.ManagedKind(v1beta1.ExistingClusterGroupVersionKind)),
		resource.IsManagedKind(resource.ManagedKind(v1beta1.ExistingClusterGroupVersionKind), mgr.GetScheme()),
	))

	r := claimbinding.NewReconciler(mgr,
		resource.ClaimKind(computev1alpha1.KubernetesClusterGroupVersionKind),
		resource.ClassKind(v1beta1.ExistingClusterClassGroupVersionKind),
		resource.ManagedKind(v1beta1.ExistingClusterGroupVersionKind),
		claimbinding.WithManagedConfigurators(
			claimbinding.ManagedConfiguratorFn(ConfigureExistingCluster),
			claimbinding.ManagedConfiguratorFn(claimbinding.ConfigureReclaimPolicy),
			claimbinding.ManagedConfiguratorFn(claimbinding.ConfigureNames),
		))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		Watches(&source.Kind{Type: &v1beta1.ExistingCluster{}}, &resource.EnqueueRequestForClaim{}).
		For(&computev1alpha1.KubernetesCluster{}).
		WithEventFilter(p).
		Complete(r)
}

// ConfigureExistingCluster configures the supplied resource (presumed to be a
// ExistingCluster) using the supplied resource claim (presumed to be a
// KubernetesCluster) and resource class.
func ConfigureExistingCluster(_ context.Context, cm resource.Claim, cs resource.Class, mg resource.Managed) error {
	if _, cmok := cm.(*computev1alpha1.KubernetesCluster); !cmok {
		return errors.Errorf("expected resource claim %s to be %s", cm.GetName(), computev1alpha1.KubernetesClusterGroupVersionKind)
	}

	rs, csok := cs.(*v1beta1.ExistingClusterClass)
	if !csok {
		return errors.Errorf("expected resource class %s to be %s", cs.GetName(), v1beta1.ExistingClusterClassGroupVersionKind)
	}

	i, mgok := mg.(*v1beta1.ExistingCluster)
	if !mgok {
		return errors.Errorf("expected managed resource %s to be %s", mg.GetName(), v1beta1.ExistingClusterGroupVersionKind)
	}

	spec := &v1beta1.ExistingClusterSpec{
		ResourceSpec: runtimev1alpha1.ResourceSpec{
			ReclaimPolicy: v1beta1.DefaultReclaimPolicy,
		},
		ForProvider: rs.SpecTemplate.ForProvider,
	}

	spec.WriteConnectionSecretToReference = &runtimev1alpha1.SecretReference{
		Namespace: rs.SpecTemplate.WriteConnectionSecretsToNamespace,
		Name:      string(cm.GetUID()),
	}
	spec.ProviderReference = rs.SpecTemplate.ProviderReference
	spec.ReclaimPolicy = rs.SpecTemplate.ReclaimPolicy

	i.Spec = *spec

	return nil
}
