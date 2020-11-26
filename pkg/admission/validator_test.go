// Copyright 2020 Nexus Operator and/or its authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package admission

import (
	"testing"

	"github.com/m88i/nexus-operator/api/v1alpha1"
	"github.com/m88i/nexus-operator/pkg/cluster/discovery"
	"github.com/m88i/nexus-operator/pkg/test"
)

func TestValidator_validate_networking(t *testing.T) {
	tests := []struct {
		name       string
		fakeClient *test.FakeClient
		input      *v1alpha1.Nexus
		wantError  bool
	}{
		{
			"'spec.networking.expose' left blank or set to false",
			test.NewFakeClientBuilder().Build(),
			&v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: false}}},
			false,
		},
		{
			"Valid Nexus with Ingress on K8s",
			test.NewFakeClientBuilder().WithIngress().Build(),
			&v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.IngressExposeType, Host: "example.com"}}},
			false,
		},
		{
			"Valid Nexus with Ingress and TLS secret on K8s",
			test.NewFakeClientBuilder().WithIngress().Build(),
			&v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.IngressExposeType, Host: "example.com", TLS: v1alpha1.NexusNetworkingTLS{SecretName: "tt-secret"}}}},
			false,
		},
		{
			"Valid Nexus with Ingress on K8s, but Ingress unavailable (Kubernetes < 1.14)",
			test.NewFakeClientBuilder().Build(),
			&v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.IngressExposeType, Host: "example.com"}}},
			true,
		},
		{
			"Invalid Nexus with Ingress on K8s and no 'spec.networking.host'",
			test.NewFakeClientBuilder().WithIngress().Build(),
			&v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.IngressExposeType}}},
			true,
		},
		{
			"Invalid Nexus with Ingress on K8s and 'spec.networking.mandatory' set to 'true'",
			test.NewFakeClientBuilder().WithIngress().Build(),
			&v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.IngressExposeType, Host: "example.com", TLS: v1alpha1.NexusNetworkingTLS{Mandatory: true}}}},
			true,
		},
		{
			"Invalid Nexus with Route on K8s",
			test.NewFakeClientBuilder().WithIngress().Build(),
			&v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.RouteExposeType}}},
			true,
		},
		{
			"Valid Nexus with Route on OCP",
			test.NewFakeClientBuilder().OnOpenshift().Build(),
			&v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.RouteExposeType}}},
			false,
		},
		{
			"Valid Nexus with Route on OCP with mandatory TLS",
			test.NewFakeClientBuilder().OnOpenshift().Build(),
			&v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.RouteExposeType, TLS: v1alpha1.NexusNetworkingTLS{Mandatory: true}}}},
			false,
		},
		{
			"Invalid Nexus with Route on OCP, but using secret name",
			test.NewFakeClientBuilder().OnOpenshift().Build(),
			&v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.RouteExposeType, TLS: v1alpha1.NexusNetworkingTLS{SecretName: "test-secret"}}}},
			true,
		},
		{
			"Invalid Nexus with Ingress on OCP",
			test.NewFakeClientBuilder().OnOpenshift().Build(),
			&v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.IngressExposeType, Host: "example.com"}}},
			true,
		},
		{
			"Valid Nexus with Node Port",
			test.NewFakeClientBuilder().Build(),
			&v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.NodePortExposeType, NodePort: 8080}}},
			false,
		},
		{
			"Invalid Nexus with Node Port and no port informed",
			test.NewFakeClientBuilder().Build(),
			&v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.NodePortExposeType}}},
			true,
		},
	}

	for _, tt := range tests {
		discovery.SetClient(tt.fakeClient)
		if err := Validate(tt.input); (err != nil) != tt.wantError {
			t.Errorf("%s\nWantError: %v\tError: %v", tt.name, tt.wantError, err)
		}
	}
}
