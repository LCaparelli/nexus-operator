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
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"

	"github.com/m88i/nexus-operator/pkg/cluster/discovery"
	"github.com/m88i/nexus-operator/pkg/framework"
	"github.com/m88i/nexus-operator/pkg/framework/kind"
	"github.com/m88i/nexus-operator/pkg/logger"
)

// SetDefaults destructively sets defaults for a given Defaultable
func SetDefaults(defaultable Defaultable) {
	d := &defaulter{
		Defaultable: defaultable,
		log:         logger.GetLoggerWithResource("defaulter", defaultable),
	}

	routeAvailable, err := discovery.IsRouteAvailable()
	if err != nil {
		d.log.Error(err, "Discovery failure", "kind", kind.RouteKind)
	}

	ingressAvailable, err := discovery.AnyIngressAvailable()
	if err != nil {
		d.log.Error(err, "Discovery failure", "kind", kind.IngressKind)
	}

	// if there were errors, these are false. Safe to use
	d.routeAvailable = routeAvailable
	d.ingressAvailable = ingressAvailable
	d.setDefaults()
}

type defaulter struct {
	Defaultable
	log logger.Logger

	routeAvailable, ingressAvailable bool
}

func (d *defaulter) logChange(field string, value interface{}) {
	valueStr := fmt.Sprintf("%+v", value)
	d.log.Debug("Setting default for "+field, "value", valueStr)
}

func (d *defaulter) setDefaults() {
	d.setDeploymentDefaults()
	d.setUpdateDefaults()
	d.setNetworkingDefaults()
	d.setPersistenceDefaults()
	d.setSecurityDefaults()
}

func (d *defaulter) setDeploymentDefaults() {
	d.setResourcesDefaults()
	d.setImageDefaults()
	d.setProbeDefaults()
}

func (d *defaulter) setResourcesDefaults() {
	resources := d.Resources()
	if resources.Requests == nil && resources.Limits == nil {
		d.logChange("spec.resources", DefaultResources)
		d.SetResources(DefaultResources)
	}
}

func (d *defaulter) setImageDefaults() {
	if d.ShouldUseRedHatImage() {
		if len(d.Image()) > 0 {
			d.log.Warn("Nexus CR configured to the use Red Hat Certified Image, ignoring 'spec.image' field.")
		}
		d.logChange("spec.image", NexusCertifiedImage)
		d.SetImage(NexusCertifiedImage)
	} else if len(d.Image()) == 0 {
		d.logChange("spec.image", NexusCommunityImage)
		d.SetImage(NexusCommunityImage)
	}

	if pullPolicy := d.ImagePullPolicy(); len(pullPolicy) > 0 &&
		pullPolicy != corev1.PullAlways &&
		pullPolicy != corev1.PullIfNotPresent &&
		pullPolicy != corev1.PullNever {

		d.log.Warn("Invalid 'spec.imagePullPolicy', unsetting the value. The pull policy will be determined by the image tag. Valid values are: " + string(corev1.PullAlways) + ", " + string(corev1.PullIfNotPresent) + " or " + string(corev1.PullNever))
		d.logChange("spec.imagePullPolicy", "")
		d.SetPullPolicy("")
	}
}

func (d *defaulter) setProbeDefaults() {
	if d.IsProbeSet(kind.LivenessProbe) {
		d.SetProbeValue(kind.LivenessProbe, kind.FailureThreshold, d.ensureMinimum(
			d.ProbeValue(kind.LivenessProbe, kind.FailureThreshold), 1, "spec.livenessProbe.failureThreshold"))

		d.SetProbeValue(kind.LivenessProbe, kind.InitialDelaySeconds, d.ensureMinimum(
			d.ProbeValue(kind.LivenessProbe, kind.InitialDelaySeconds), 1, "spec.livenessProbe.initialDelaySeconds"))

		d.SetProbeValue(kind.LivenessProbe, kind.PeriodSeconds, d.ensureMinimum(
			d.ProbeValue(kind.LivenessProbe, kind.PeriodSeconds), 1, "spec.livenessProbe.periodSeconds"))

		d.SetProbeValue(kind.LivenessProbe, kind.TimeoutSeconds, d.ensureMinimum(
			d.ProbeValue(kind.LivenessProbe, kind.TimeoutSeconds), 1, "spec.livenessProbe.timeoutSeconds"))
	} else {
		d.logChange("spec.livenessProbe", d.DefaultLivenessProbe())
		d.UseDefaultLivenessProbe()
	}

	// SuccessThreshold for Liveness Probes must be 1
	if d.ProbeValue(kind.LivenessProbe, kind.SuccessThreshold) != 1 {
		d.logChange("spec.readinessProbe.successThreshold", 1)
		d.SetProbeValue(kind.LivenessProbe, kind.SuccessThreshold, 1)
	}

	if d.IsProbeSet(kind.ReadinessProbe) {
		d.SetProbeValue(kind.ReadinessProbe, kind.FailureThreshold, d.ensureMinimum(
			d.ProbeValue(kind.ReadinessProbe, kind.FailureThreshold), 1, "spec.livenessProbe.failureThreshold"))

		d.SetProbeValue(kind.ReadinessProbe, kind.InitialDelaySeconds, d.ensureMinimum(
			d.ProbeValue(kind.ReadinessProbe, kind.InitialDelaySeconds), 1, "spec.livenessProbe.initialDelaySeconds"))

		d.SetProbeValue(kind.ReadinessProbe, kind.PeriodSeconds, d.ensureMinimum(
			d.ProbeValue(kind.ReadinessProbe, kind.PeriodSeconds), 1, "spec.livenessProbe.periodSeconds"))

		d.SetProbeValue(kind.ReadinessProbe, kind.TimeoutSeconds, d.ensureMinimum(
			d.ProbeValue(kind.ReadinessProbe, kind.TimeoutSeconds), 1, "spec.livenessProbe.timeoutSeconds"))

		d.SetProbeValue(kind.ReadinessProbe, kind.SuccessThreshold, d.ensureMinimum(
			d.ProbeValue(kind.ReadinessProbe, kind.SuccessThreshold), 1, "spec.livenessProbe.successThreshold"))
	} else {
		d.logChange("spec.readinessProbe", d.DefaultReadinessProbe())
		d.UseDefaultReadinessProbe()
	}
}

// must be called only after image defaults have been set
func (d *defaulter) setUpdateDefaults() {
	if d.IsUpdateDisabled() {
		return
	}

	image := strings.Split(d.Image(), ":")[0]
	if image != NexusCommunityImage {
		d.log.Warn("Automatic Updates are enabled, but 'spec.image' is not using the community image. Disabling automatic updates", "Community Image", NexusCommunityImage)
		d.logChange("spec.automaticUpdate.disabled", true)
		d.DisableUpdate()
		return
	}

	if !d.IsMinorVersionSet() {
		d.log.Debug("Automatic Updates are enabled, but no minor was informed. Fetching the most recent...")
		minor, err := framework.GetLatestMinor()
		if err != nil {
			d.handleTagFetchError(err)
			return
		}
		d.logChange("spec.automaticUpdate.minorVersion", minor)
		d.SetMinorVersion(minor)
	}

	d.log.Debug("Fetching the latest micro from minor", "MinorVersion", d.MinorVersion())
	tag, ok := framework.GetLatestMicro(d.MinorVersion())
	if !ok {
		// the informed minor doesn't exist, let's try the latest minor
		d.log.Warn("Latest tag for minor version not found. Trying the latest minor instead", "Informed tag", d.MinorVersion())
		minor, err := framework.GetLatestMinor()
		if err != nil {
			d.handleTagFetchError(err)
			return
		}
		d.log.Info("Setting 'spec.automaticUpdate.minorVersion to", "MinorTag", minor)
		d.logChange("spec.automaticUpdate.minorVersion", minor)
		d.SetMinorVersion(minor)
		// no need to check for the tag existence here,
		// we would have gotten an error from GetLatestMinor() if it didn't
		tag, _ = framework.GetLatestMicro(minor)
	}
	newImage := fmt.Sprintf("%s:%s", image, tag)
	if newImage != d.Image() {
		d.logChange("spec.image", newImage)
		d.SetImage(newImage)
	}
}

func (d *defaulter) setNetworkingDefaults() {
	if !d.IsExposeSet() && len(d.ExposeAs()) == 0 {
		d.logChange("spec.networking.expose", false)
		d.SetExpose(false)
		return
	}

	if d.ShouldExpose() && len(d.ExposeAs()) == 0 {
		// expose is true, but exposeAs is blank
		// let's figure out the best way to expose
		if d.routeAvailable {
			d.logChange("spec.networking.exposeAs", routeExposeType)
			d.SetExposeAs(routeExposeType)
		} else if d.ingressAvailable {
			d.logChange("spec.networking.exposeAs", ingressExposeType)
			d.SetExposeAs(ingressExposeType)
		} else {
			// we're on kubernetes < 1.14
			// try setting nodePort, validation will catch it if impossible
			d.log.Info("On Kubernetes, but Ingresses are not available")
			d.logChange("spec.networking.exposeAs", nodePortExposeType)
			d.SetExposeAs(nodePortExposeType)
		}
	} else if !d.IsExposeSet() {
		// expose is unset but exposeAs is not blank
		// let's set expose to true
		d.logChange("spec.networking.expose", true)
		d.SetExpose(true)
	}
}

func (d *defaulter) setPersistenceDefaults() {
	if d.IsPersistent() && len(d.VolumeSize()) == 0 {
		d.logChange("spec.persistence.volumeSize", DefaultVolumeSize)
		d.SetVolumeSize(DefaultVolumeSize)
	}
}

func (d *defaulter) setSecurityDefaults() {
	if len(d.ServiceAccount()) == 0 {
		d.logChange("spec.ServiceAccountName", d.GetName())
		d.SetServiceAccount(d.GetName())
	}
}

func (d *defaulter) ensureMinimum(value, minimum int32, field string) int32 {
	if value < minimum {
		d.log.Warn(field + " below minimum.")
		d.logChange(field, minimum)
		return minimum
	}
	return value
}

func (d *defaulter) handleTagFetchError(err error) {
	d.log.Error(err, "Unable to fetch the most recent minor. Disabling automatic updates.")
	d.logChange("spec.automaticUpdate.disabled", true)
	d.DisableUpdate()
	createChangedNexusEvent(d.Defaultable, "spec.automaticUpdate.disabled")
}
