package kind

type ProbeType int
type ProbeField int

const (
	LivenessProbe ProbeType = iota
	ReadinessProbe

	FailureThreshold ProbeField = iota
	InitialDelaySeconds
	PeriodSeconds
	TimeoutSeconds
	SuccessThreshold
)
