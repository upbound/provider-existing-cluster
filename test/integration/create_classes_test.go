// +build integration

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

package integration

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/crossplaneio/crossplane-runtime/pkg/test/integration"
	crossplaneapis "github.com/crossplaneio/crossplane/apis"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/turkenh/provider-existing-cluster/apis"
	containerv1beta1 "github.com/turkenh/provider-existing-cluster/apis/container/v1beta1"
	"github.com/turkenh/provider-existing-cluster/pkg/controller"
)

func TestCreateAllClasses(t *testing.T) {
	cases := map[string]struct {
		reason string
		test   func(c client.Client) error
	}{
		"CreateV1Beta1ExistingClusterClass": {
			reason: "A v1beta1 ExistingClusterClass should be created without error.",
			test: func(c client.Client) error {
				dat, err := ioutil.ReadFile("../../examples/container/kubernetescluster/resource-class.yaml")
				if err != nil {
					return err
				}
				s := &containerv1beta1.ExistingClusterClass{}
				if err := yaml.Unmarshal(dat, s); err != nil {
					return err
				}

				defer func() {
					if err := c.Delete(context.Background(), s); err != nil {
						t.Error(err)
					}
				}()

				return c.Create(context.Background(), s)
			},
		},
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", "../../kubeconfig.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../../sa.json"); err != nil {
		t.Fatal(err)
	}

	i, err := integration.New(cfg,
		integration.WithCRDPaths("../../config/crd"),
		integration.WithCleaners(
			integration.NewCRDCleaner(),
			integration.NewCRDDirCleaner()),
	)

	if err != nil {
		t.Fatal(err)
	}

	if err := apis.AddToScheme(i.GetScheme()); err != nil {
		t.Fatal(err)
	}

	if err := crossplaneapis.AddToScheme(i.GetScheme()); err != nil {
		t.Fatal(err)
	}

	if err := (&controller.Controllers{}).SetupWithManager(i); err != nil {
		t.Fatal(err)
	}

	i.Run()

	defer func() {
		if err := i.Cleanup(); err != nil {
			t.Fatal(err)
		}
	}()

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := tc.test(i.GetClient())
			if err != nil {
				t.Error(err)
			}
		})
	}
}
