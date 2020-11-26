package v1alpha1

func (n *Nexus) IsExposingAs(exposeAs string) bool {
	return n.Spec.Networking.ExposeAs == NexusNetworkingExposeType(exposeAs)
}

func (n *Nexus) NodePort() int32 {
	return n.Spec.Networking.NodePort
}

func (n *Nexus) Host() string {
	return n.Spec.Networking.Host
}

func (n *Nexus) SecretName() string {
	return n.Spec.Networking.TLS.SecretName
}

func (n *Nexus) IsTLSMandatory() bool {
	return n.Spec.Networking.TLS.Mandatory
}
