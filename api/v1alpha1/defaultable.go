package v1alpha1

import (
	"errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/m88i/nexus-operator/pkg/admission"
	"github.com/m88i/nexus-operator/pkg/framework/kind"
)

const (
	probeDefaultInitialDelaySeconds = int32(240)
	probeDefaultTimeoutSeconds      = int32(15)
	probeDefaultPeriodSeconds       = int32(10)
	probeDefaultSuccessThreshold    = int32(1)
	probeDefaultFailureThreshold    = int32(3)
)

var (
	DefaultProbe = &NexusProbe{
		InitialDelaySeconds: probeDefaultInitialDelaySeconds,
		TimeoutSeconds:      probeDefaultTimeoutSeconds,
		PeriodSeconds:       probeDefaultPeriodSeconds,
		SuccessThreshold:    probeDefaultSuccessThreshold,
		FailureThreshold:    probeDefaultFailureThreshold,
	}

	AllDefaultsCommunityNexus = Nexus{
		ObjectMeta: v1.ObjectMeta{Name: "default-community-nexus", Namespace: "default"},
		Spec: NexusSpec{
			Replicas:                    0,
			Image:                       admission.NexusCommunityImage,
			ImagePullPolicy:             "",
			Resources:                   admission.DefaultResources,
			UseRedHatImage:              false,
			GenerateRandomAdminPassword: false,
			ServiceAccountName:          "default-community-nexus",
			LivenessProbe:               DefaultProbe.DeepCopy(),
			ReadinessProbe:              DefaultProbe.DeepCopy(),
		},
	}
)

func (n *Nexus) Resources() corev1.ResourceRequirements {
	return n.Spec.Resources
}

func (n *Nexus) SetResources(requirements corev1.ResourceRequirements) {
	n.Spec.Resources = requirements
}

func (n *Nexus) ShouldUseRedHatImage() bool {
	return n.Spec.UseRedHatImage
}

func (n *Nexus) Image() string {
	return n.Spec.Image
}

func (n *Nexus) SetImage(img string) {
	n.Spec.Image = img
}

func (n *Nexus) ImagePullPolicy() corev1.PullPolicy {
	return n.Spec.ImagePullPolicy
}

func (n *Nexus) SetPullPolicy(policy corev1.PullPolicy) {
	n.Spec.ImagePullPolicy = policy
}

func (n *Nexus) IsProbeSet(probeType kind.ProbeType) bool {
	if probeType == kind.ReadinessProbe {
		return n.Spec.ReadinessProbe != nil
	}
	return n.Spec.LivenessProbe != nil
}

func (n *Nexus) ProbeValue(probeType kind.ProbeType, field kind.ProbeField) int32 {
	if probeType == kind.LivenessProbe {
		return n.livenessProbeValue(field)
	}
	return n.readinessProbeValue(field)
}

func (n *Nexus) livenessProbeValue(field kind.ProbeField) int32 {
	switch field {
	case kind.SuccessThreshold:
		return n.Spec.LivenessProbe.SuccessThreshold
	case kind.TimeoutSeconds:
		return n.Spec.LivenessProbe.TimeoutSeconds
	case kind.PeriodSeconds:
		return n.Spec.LivenessProbe.PeriodSeconds
	case kind.InitialDelaySeconds:
		return n.Spec.LivenessProbe.InitialDelaySeconds
	case kind.FailureThreshold:
		return n.Spec.LivenessProbe.FailureThreshold
	default:
		// something is VERY wrong, there is definitely a bug
		log.Error(errors.New("unknown liveness probe field requested"), "BUG, please report this")
		return 0
	}
}

func (n *Nexus) readinessProbeValue(field kind.ProbeField) int32 {
	switch field {
	case kind.SuccessThreshold:
		return n.Spec.ReadinessProbe.SuccessThreshold
	case kind.TimeoutSeconds:
		return n.Spec.ReadinessProbe.TimeoutSeconds
	case kind.PeriodSeconds:
		return n.Spec.ReadinessProbe.PeriodSeconds
	case kind.InitialDelaySeconds:
		return n.Spec.ReadinessProbe.InitialDelaySeconds
	case kind.FailureThreshold:
		return n.Spec.ReadinessProbe.FailureThreshold
	default:
		// something is VERY wrong, there is definitely a bug
		log.Error(errors.New("unknown readiness probe field requested"), "BUG, please report this")
		return 0
	}
}

func (n *Nexus) SetProbeValue(probeType kind.ProbeType, field kind.ProbeField, value int32) {
	if probeType == kind.LivenessProbe {
		n.setLivenessProbeValue(field, value)
	}
	n.setReadinessProbeValue(field, value)
}

func (n *Nexus) setLivenessProbeValue(field kind.ProbeField, value int32) {
	switch field {
	case kind.SuccessThreshold:
		n.Spec.LivenessProbe.SuccessThreshold = value
	case kind.TimeoutSeconds:
		n.Spec.LivenessProbe.TimeoutSeconds = value
	case kind.PeriodSeconds:
		n.Spec.LivenessProbe.PeriodSeconds = value
	case kind.InitialDelaySeconds:
		n.Spec.LivenessProbe.InitialDelaySeconds = value
	case kind.FailureThreshold:
		n.Spec.LivenessProbe.FailureThreshold = value
	default:
		// something is VERY wrong, there is definitely a bug
		log.Error(errors.New("attempted to set unknown liveness probe field"), "BUG, please report this")
	}
}

func (n *Nexus) setReadinessProbeValue(field kind.ProbeField, value int32) {
	switch field {
	case kind.SuccessThreshold:
		n.Spec.ReadinessProbe.SuccessThreshold = value
	case kind.TimeoutSeconds:
		n.Spec.ReadinessProbe.TimeoutSeconds = value
	case kind.PeriodSeconds:
		n.Spec.ReadinessProbe.PeriodSeconds = value
	case kind.InitialDelaySeconds:
		n.Spec.ReadinessProbe.InitialDelaySeconds = value
	case kind.FailureThreshold:
		n.Spec.ReadinessProbe.FailureThreshold = value
	default:
		// something is VERY wrong, there is definitely a bug
		log.Error(errors.New("attempted to set unknown readiness probe field"), "BUG, please report this")
	}
}

func (n *Nexus) DefaultLivenessProbe() interface{} {
	return *DefaultProbe
}

func (n *Nexus) UseDefaultLivenessProbe() {
	n.Spec.LivenessProbe = DefaultProbe.DeepCopy()
}

func (n *Nexus) DefaultReadinessProbe() interface{} {
	return *DefaultProbe
}

func (n *Nexus) UseDefaultReadinessProbe() {
	n.Spec.ReadinessProbe = DefaultProbe.DeepCopy()
}

func (n *Nexus) IsUpdateDisabled() bool {
	return n.Spec.AutomaticUpdate.Disabled
}

func (n *Nexus) DisableUpdate() {
	n.Spec.AutomaticUpdate.Disabled = true
}

func (n *Nexus) IsMinorVersionSet() bool {
	return n.Spec.AutomaticUpdate.MinorVersion != nil
}

func (n *Nexus) SetMinorVersion(minor int) {
	n.Spec.AutomaticUpdate.MinorVersion = &minor
}

func (n *Nexus) MinorVersion() int {
	return *n.Spec.AutomaticUpdate.MinorVersion
}

func (n *Nexus) IsPersistent() bool {
	return n.Spec.Persistence.Persistent
}

func (n *Nexus) VolumeSize() string {
	return n.Spec.Persistence.VolumeSize
}

func (n *Nexus) SetVolumeSize(volumeSize string) {
	n.Spec.Persistence.VolumeSize = volumeSize
}

func (n *Nexus) ServiceAccount() string {
	return n.Spec.ServiceAccountName
}

func (n *Nexus) SetServiceAccount(svcAccnt string) {
	n.Spec.ServiceAccountName = svcAccnt
}

func (n *Nexus) IsExposeSet() bool {
	return n.Spec.Networking.Expose != nil
}

func (n *Nexus) ShouldExpose() bool {
	if n.Spec.Networking.Expose == nil {
		return false
	}
	return *n.Spec.Networking.Expose
}

func (n *Nexus) ExposeAs() string {
	return string(n.Spec.Networking.ExposeAs)
}

func (n *Nexus) SetExpose(expose bool) {
	n.Spec.Networking.Expose = &expose
}

func (n *Nexus) SetExposeAs(exposeAs string) {
	n.Spec.Networking.ExposeAs = NexusNetworkingExposeType(exposeAs)
}
