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

	"github.com/crossplaneio/crossplane-runtime/pkg/meta"
	v1beta12 "github.com/crossplaneio/stack-existing-cluster/apis/v1beta1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
	runtimev1alpha1 "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
	"github.com/crossplaneio/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplaneio/crossplane-runtime/pkg/resource"

	"github.com/crossplaneio/stack-existing-cluster/apis/container/v1beta1"
)

// Error strings.
const (
	errGetProvider       = "cannot get Provider"
	errGetProviderSecret = "cannot get Provider Secret"
	errNotCluster        = "managed resource is not a ExistingCluster"
)

// ExistingClusterController is responsible for adding the ExistingCluster
// controller and its corresponding reconciler to the manager with any runtime configuration.
type ExistingClusterController struct{}

// SetupWithManager creates a new Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func (c *ExistingClusterController) SetupWithManager(mgr ctrl.Manager) error {
	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1beta1.ExistingClusterGroupVersionKind),
		managed.WithExternalConnecter(&clusterConnector{kube: mgr.GetClient()}))

	name := strings.ToLower(fmt.Sprintf("%s.%s", v1beta1.ExistingClusterKindAPIVersion, v1beta1.Group))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1beta1.ExistingCluster{}).
		Complete(r)
}

type clusterConnector struct {
	kube client.Client
}

func (c *clusterConnector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	i, ok := mg.(*v1beta1.ExistingCluster)
	if !ok {
		return nil, errors.New(errNotCluster)
	}

	p := &v1beta12.Provider{}
	if err := c.kube.Get(ctx, meta.NamespacedNameOf(i.Spec.ProviderReference), p); err != nil {
		return nil, errors.Wrap(err, errGetProvider)
	}

	s := &corev1.Secret{}
	n := types.NamespacedName{Namespace: p.Spec.CredentialsSecretRef.Namespace, Name: p.Spec.CredentialsSecretRef.Name}
	if err := c.kube.Get(ctx, n, s); err != nil {
		return nil, errors.Wrap(err, errGetProviderSecret)
	}

	return &clusterExternal{kube: c.kube, configData: s.Data[runtimev1alpha1.ResourceCredentialsSecretKubeconfigKey]}, nil
}

type clusterExternal struct {
	kube       client.Client
	configData []byte
}

func (e *clusterExternal) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) { // nolint:gocyclo
	cr, ok := mg.(*v1beta1.ExistingCluster)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotCluster)
	}

	cr.Status.AtProvider.Status = v1beta1.ClusterStateRunning
	cr.Status.SetConditions(v1alpha1.Available())
	resource.SetBindable(cr)

	return managed.ExternalObservation{
		ResourceExists:    true,
		ResourceUpToDate:  true,
		ConnectionDetails: connectionDetails(e.configData),
	}, nil
}

func (e *clusterExternal) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1beta1.ExistingCluster)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotCluster)
	}
	cr.SetConditions(v1alpha1.Creating())

	return managed.ExternalCreation{}, nil
}

func (e *clusterExternal) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	_, ok := mg.(*v1beta1.ExistingCluster)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotCluster)
	}
	return managed.ExternalUpdate{}, nil
}

func (e *clusterExternal) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1beta1.ExistingCluster)
	if !ok {
		return errors.New(errNotCluster)
	}
	cr.SetConditions(runtimev1alpha1.Deleting())
	return nil
}

// connectionSecret return secret object for cluster instance
func connectionDetails(rawConfig []byte) managed.ConnectionDetails {
	cd := managed.ConnectionDetails{
		runtimev1alpha1.ResourceCredentialsSecretKubeconfigKey: rawConfig,
	}
	return cd
}
