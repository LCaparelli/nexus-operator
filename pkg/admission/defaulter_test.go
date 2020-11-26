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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"

	"github.com/m88i/nexus-operator/api/v1alpha1"
	"github.com/m88i/nexus-operator/pkg/cluster/discovery"
	"github.com/m88i/nexus-operator/pkg/framework"
	"github.com/m88i/nexus-operator/pkg/test"
)

var shouldExpose = true

func TestSetDefaults_deployment(t *testing.T) {
	discovery.SetClient(test.NewFakeClientBuilder().Build())
	minimumDefaultProbe := &v1alpha1.NexusProbe{
		InitialDelaySeconds: 0,
		TimeoutSeconds:      1,
		PeriodSeconds:       1,
		SuccessThreshold:    1,
		FailureThreshold:    1,
	}

	tests := []struct {
		name  string
		input *v1alpha1.Nexus
		want  *v1alpha1.Nexus
	}{
		{
			"'spec.resources' left blank",
			func() *v1alpha1.Nexus {
				nexus := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				nexus.Spec.Resources = v1.ResourceRequirements{}
				return nexus
			}(),
			&v1alpha1.AllDefaultsCommunityNexus,
		},
		{
			"'spec.useRedHatImage' set to true and 'spec.image' not left blank",
			func() *v1alpha1.Nexus {
				nexus := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				nexus.Spec.UseRedHatImage = true
				nexus.Spec.Image = "some-image"
				return nexus
			}(),
			func() *v1alpha1.Nexus {
				n := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				n.Spec.UseRedHatImage = true
				n.Spec.Image = NexusCertifiedImage
				return n
			}(),
		},
		{
			"'spec.useRedHatImage' set to false and 'spec.image' left blank",
			func() *v1alpha1.Nexus {
				nexus := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				nexus.Spec.Image = ""
				return nexus
			}(),
			&v1alpha1.AllDefaultsCommunityNexus,
		},
		{
			"'spec.livenessProbe.successThreshold' not equal to 1",
			func() *v1alpha1.Nexus {
				nexus := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				nexus.Spec.LivenessProbe.SuccessThreshold = 2
				return nexus
			}(),
			&v1alpha1.AllDefaultsCommunityNexus,
		},
		{
			"'spec.livenessProbe.*' and 'spec.readinessProbe.*' don't meet minimum values",
			func() *v1alpha1.Nexus {
				nexus := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				nexus.Spec.LivenessProbe = &v1alpha1.NexusProbe{
					InitialDelaySeconds: -1,
					TimeoutSeconds:      -1,
					PeriodSeconds:       -1,
					SuccessThreshold:    -1,
					FailureThreshold:    -1,
				}
				nexus.Spec.ReadinessProbe = nexus.Spec.LivenessProbe.DeepCopy()
				return nexus
			}(),
			func() *v1alpha1.Nexus {
				nexus := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				nexus.Spec.LivenessProbe = minimumDefaultProbe.DeepCopy()
				nexus.Spec.ReadinessProbe = minimumDefaultProbe.DeepCopy()
				return nexus
			}(),
		},
		{
			"Unset 'spec.livenessProbe' and 'spec.readinessProbe'",
			func() *v1alpha1.Nexus {
				nexus := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				nexus.Spec.LivenessProbe = nil
				nexus.Spec.ReadinessProbe = nil
				return nexus
			}(),
			v1alpha1.AllDefaultsCommunityNexus.DeepCopy(),
		},
		{
			"Invalid 'spec.imagePullPolicy'",
			func() *v1alpha1.Nexus {
				nexus := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				nexus.Spec.ImagePullPolicy = "invalid"
				return nexus
			}(),
			v1alpha1.AllDefaultsCommunityNexus.DeepCopy(),
		},
	}

	for _, tt := range tests {
		SetDefaults(tt.input)
		if !reflect.DeepEqual(tt.input, tt.want) {
			t.Errorf("%s\nWant: %+v\nGot: %+v", tt.name, tt.want, tt.input)
		}
	}
}

func TestSetDefaults_automaticUpdate(t *testing.T) {
	discovery.SetClient(test.NewFakeClientBuilder().Build())
	nexus := &v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{AutomaticUpdate: v1alpha1.NexusAutomaticUpdate{}}}
	nexus.Spec.Image = NexusCommunityImage

	SetDefaults(nexus)
	latestMinor, err := framework.GetLatestMinor()
	if err != nil {
		// If we couldn't fetch the tags updates should be disabled
		assert.True(t, nexus.Spec.AutomaticUpdate.Disabled)
	} else {
		assert.Equal(t, latestMinor, *nexus.Spec.AutomaticUpdate.MinorVersion)
	}

	// Now an invalid image
	nexus = &v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{AutomaticUpdate: v1alpha1.NexusAutomaticUpdate{}}}
	nexus.Spec.Image = "some-image"
	SetDefaults(nexus)
	assert.True(t, nexus.Spec.AutomaticUpdate.Disabled)

	// Informed a minor which does not exist
	nexus = &v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{AutomaticUpdate: v1alpha1.NexusAutomaticUpdate{}}}
	nexus.Spec.Image = NexusCommunityImage
	bogusMinor := -1
	nexus.Spec.AutomaticUpdate.MinorVersion = &bogusMinor
	SetDefaults(nexus)
	latestMinor, err = framework.GetLatestMinor()
	if err != nil {
		// If we couldn't fetch the tags updates should be disabled
		assert.True(t, nexus.Spec.AutomaticUpdate.Disabled)
	} else {
		assert.Equal(t, latestMinor, *nexus.Spec.AutomaticUpdate.MinorVersion)
	}
}

func TestSetDefaults_networking(t *testing.T) {
	tests := []struct {
		name       string
		fakeClient *test.FakeClient
		input      *v1alpha1.Nexus
		want       *v1alpha1.Nexus
	}{
		{
			"'spec.networking.exposeAs' left blank on OCP",
			test.NewFakeClientBuilder().OnOpenshift().Build(),
			func() *v1alpha1.Nexus {
				n := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				n.Spec.Networking.Expose = &shouldExpose
				return n
			}(),
			func() *v1alpha1.Nexus {
				n := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				n.Spec.Networking.Expose = &shouldExpose
				n.Spec.Networking.ExposeAs = v1alpha1.RouteExposeType
				return n
			}(),
		},
		{
			"'spec.networking.exposeAs' left blank on K8s",
			test.NewFakeClientBuilder().WithIngress().Build(),
			func() *v1alpha1.Nexus {
				n := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				n.Spec.Networking.Expose = &shouldExpose
				return n
			}(),
			func() *v1alpha1.Nexus {
				n := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				n.Spec.Networking.Expose = &shouldExpose
				n.Spec.Networking.ExposeAs = v1alpha1.IngressExposeType
				return n
			}(),
		},
		{
			"'spec.networking.exposeAs' left blank on K8s, but Ingress unavailable",
			test.NewFakeClientBuilder().Build(),
			func() *v1alpha1.Nexus {
				n := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				n.Spec.Networking.Expose = &shouldExpose
				return n
			}(),
			func() *v1alpha1.Nexus {
				n := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				n.Spec.Networking.Expose = &shouldExpose
				n.Spec.Networking.ExposeAs = v1alpha1.NodePortExposeType
				return n
			}(),
		},
	}

	for _, tt := range tests {
		discovery.SetClient(tt.fakeClient)
		SetDefaults(tt.input)
		if !reflect.DeepEqual(tt.input, tt.want) {
			t.Errorf("%s\nWant: %v\nGot: %v", tt.name, tt.want, tt.input)
		}
	}
}

func TestSetDefaults_persistence(t *testing.T) {
	discovery.SetClient(test.NewFakeClientBuilder().Build())
	tests := []struct {
		name  string
		input *v1alpha1.Nexus
		want  *v1alpha1.Nexus
	}{
		{
			"'spec.persistence.volumeSize' left blank",
			func() *v1alpha1.Nexus {
				n := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				n.Spec.Persistence.Persistent = true
				n.Spec.Persistence.VolumeSize = ""
				return n
			}(),
			func() *v1alpha1.Nexus {
				n := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				n.Spec.Persistence.Persistent = true
				n.Spec.Persistence.VolumeSize = DefaultVolumeSize
				return n
			}(),
		},
	}
	for _, tt := range tests {
		SetDefaults(tt.input)
		if !reflect.DeepEqual(tt.input, tt.want) {
			t.Errorf("%s\nWant: %+v\nGot: %+v", tt.name, tt.want, tt.input)
		}
	}
}

func TestSetDefaults_security(t *testing.T) {
	discovery.SetClient(test.NewFakeClientBuilder().Build())
	tests := []struct {
		name  string
		input *v1alpha1.Nexus
		want  *v1alpha1.Nexus
	}{
		{
			"'spec.serviceAccountName' left blank",
			func() *v1alpha1.Nexus {
				nexus := v1alpha1.AllDefaultsCommunityNexus.DeepCopy()
				nexus.Spec.ServiceAccountName = ""
				return nexus
			}(),
			&v1alpha1.AllDefaultsCommunityNexus,
		},
	}
	for _, tt := range tests {
		SetDefaults(tt.input)
		if !reflect.DeepEqual(tt.input, tt.want) {
			t.Errorf("%s\nWant: %+v\nGot: %+v", tt.name, tt.want, tt.input)
		}
	}
}
