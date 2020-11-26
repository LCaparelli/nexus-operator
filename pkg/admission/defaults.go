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
	corev1 "k8s.io/api/core/v1"
	k8sres "k8s.io/apimachinery/pkg/api/resource"
)

const (
	NexusCommunityImage = "docker.io/sonatype/nexus3"
	NexusCertifiedImage = "registry.connect.redhat.com/sonatype/validatable-repository-manager"

	DefaultVolumeSize = "10Gi"
)

var (
	DefaultResources = corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    k8sres.MustParse("2"),
			corev1.ResourceMemory: k8sres.MustParse("2Gi"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    k8sres.MustParse("1"),
			corev1.ResourceMemory: k8sres.MustParse("2Gi"),
		},
	}
)
