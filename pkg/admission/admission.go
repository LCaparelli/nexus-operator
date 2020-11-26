package admission

const (
	// importing these from v1alpha1 gets an import cycle, so we mirror them here
	ingressExposeType  = "Ingress"
	routeExposeType    = "Route"
	nodePortExposeType = "NodePort"
)

var allDefaultsCommunityNexus =