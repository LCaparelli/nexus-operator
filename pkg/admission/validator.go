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

	"github.com/m88i/nexus-operator/pkg/cluster/discovery"
	"github.com/m88i/nexus-operator/pkg/framework/kind"
	"github.com/m88i/nexus-operator/pkg/logger"
)

// Validate returns an error if the given Validatable is invalid
func Validate(validatable Validatable) error {
	v := &validator{
		Validatable: validatable,
		log:         logger.GetLoggerWithResource("validator", validatable),
	}

	routeAvailable, err := discovery.IsRouteAvailable()
	if err != nil {
		v.log.Error(err, "Discovery failure", "kind", kind.RouteKind)
	}

	ingressAvailable, err := discovery.AnyIngressAvailable()
	if err != nil {
		v.log.Error(err, "Discovery failure", "kind", kind.IngressKind)
	}

	// if there were errors, these are false. Safe to use
	v.routeAvailable = routeAvailable
	v.ingressAvailable = ingressAvailable
	return v.validate()
}

type validator struct {
	Validatable
	log              logger.Logger
	routeAvailable   bool
	ingressAvailable bool
}

func (v *validator) validate() error {
	return v.validateNetworking()
}

func (v *validator) validateNetworking() error {
	if !v.ShouldExpose() {
		v.log.Debug("'spec.networking.expose' set to 'false', ignoring networking configuration")
		return nil
	}

	if !v.ingressAvailable && v.IsExposingAs(ingressExposeType) {
		v.log.Warn("Ingresses are not available on your cluster. Make sure to be running Kubernetes > 1.14 or if you're running Openshift set ", "spec.networking.exposeAs", routeExposeType, "Also try", nodePortExposeType)
		return fmt.Errorf("ingress expose required, but unavailable")
	}

	if !v.routeAvailable && v.IsExposingAs(routeExposeType) {
		v.log.Warn("Routes are not available on your cluster. If you're running Kubernetes 1.14 or higher try setting ", "'spec.networking.exposeAs'", ingressExposeType, "Also try", nodePortExposeType)
		return fmt.Errorf("route expose required, but unavailable")
	}

	if v.IsExposingAs(nodePortExposeType) && v.NodePort() == 0 {
		v.log.Warn("NodePort networking requires a port. Check the Nexus resource 'spec.networking.nodePort' parameter")
		return fmt.Errorf("nodeport expose required, but no port informed")
	}

	if v.IsExposingAs(ingressExposeType) && len(v.Host()) == 0 {
		v.log.Warn("Ingress networking requires a host. Check the Nexus resource 'spec.networking.host' parameter")
		return fmt.Errorf("ingress expose required, but no host informed")
	}

	if len(v.SecretName()) > 0 && !v.IsExposingAs(ingressExposeType) {
		v.log.Warn("'spec.networking.tls.secretName' is only available when using an Ingress. Try setting ", "spec.networking.exposeAs'", ingressExposeType)
		return fmt.Errorf("tls secret name informed, but using route")
	}

	if v.IsTLSMandatory() && !v.IsExposingAs(routeExposeType) {
		v.log.Warn("'spec.networking.tls.mandatory' is only available when using a Route. Try setting ", "spec.networking.exposeAs'", routeExposeType)
		return fmt.Errorf("tls set to mandatory, but using ingress")
	}

	return nil
}
