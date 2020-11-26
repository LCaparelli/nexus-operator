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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/m88i/nexus-operator/pkg/admission"
	"github.com/m88i/nexus-operator/pkg/logger"
)

const logName = "nexus-resource"

// log is for logging in this package.
var log = logger.GetLogger(logName)

func (n *Nexus) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(n).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-apps-m88i-io-m88i-io-v1alpha1-validatable,mutating=true,failurePolicy=fail,groups=apps.m88i.io.m88i.io,resources=validatable,verbs=create;update,versions=v1alpha1,name=mnexus.kb.io

var _ webhook.Defaulter = &Nexus{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (n *Nexus) Default() {
	log = logger.GetLoggerWithResource(logName, n)
	defer func() { log = logger.GetLogger(logName) }()
	log.Info("Setting defaults")
	admission.SetDefaults(n)
}

// if validation upon deletion ever becomes necessary, insert "delete" in the verbs
// +kubebuilder:webhook:verbs=create;update,path=/validate-apps-m88i-io-m88i-io-v1alpha1-validatable,mutating=false,failurePolicy=fail,groups=apps.m88i.io.m88i.io,resources=validatable,versions=v1alpha1,name=vnexus.kb.io

var _ webhook.Validator = &Nexus{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (n *Nexus) ValidateCreate() error {
	log = logger.GetLoggerWithResource(logName, n)
	defer func() { log = logger.GetLogger(logName) }()
	log.Info("Validating create request")
	return admission.Validate(n)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (n *Nexus) ValidateUpdate(old runtime.Object) error {
	log = logger.GetLoggerWithResource(logName, n)
	defer func() { log = logger.GetLogger(logName) }()
	log.Info("Validating update request")
	// we don't really care about the old validatable for validation
	return admission.Validate(n)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (*Nexus) ValidateDelete() error {
	return nil
}
