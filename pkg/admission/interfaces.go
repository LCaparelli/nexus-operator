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
	"github.com/RHsyseng/operator-utils/pkg/resource"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/m88i/nexus-operator/pkg/framework/kind"
	"github.com/m88i/nexus-operator/pkg/logger"
)

const logName = "admission"

// log is for logging in this package.
var log = logger.GetLogger(logName)

type Defaultable interface {
	resource.KubernetesResource
	runtime.Object
	Resources() corev1.ResourceRequirements
	SetResources(corev1.ResourceRequirements)
	ShouldUseRedHatImage() bool
	Image() string
	SetImage(string)
	ImagePullPolicy() corev1.PullPolicy
	SetPullPolicy(corev1.PullPolicy)
	IsProbeSet(kind.ProbeType) bool
	ProbeValue(kind.ProbeType, kind.ProbeField) int32
	SetProbeValue(kind.ProbeType, kind.ProbeField, int32)
	DefaultLivenessProbe() interface{}
	UseDefaultLivenessProbe()
	DefaultReadinessProbe() interface{}
	UseDefaultReadinessProbe()
	IsUpdateDisabled() bool
	DisableUpdate()
	IsMinorVersionSet() bool
	SetMinorVersion(int)
	MinorVersion() int
	IsPersistent() bool
	VolumeSize() string
	SetVolumeSize(string)
	ServiceAccount() string
	SetServiceAccount(string)
	IsExposeSet() bool
	ShouldExpose() bool
	ExposeAs() string
	SetExpose(bool)
	SetExposeAs(string)
}

type Validatable interface {
	resource.KubernetesResource
	ShouldExpose() bool
	IsExposingAs(string) bool
	NodePort() int32
	Host() string
	SecretName() string
	IsTLSMandatory() bool
}
