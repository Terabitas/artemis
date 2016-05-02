package domain

import "net"

const (
	DigitalOcean = "digitalocean"
	Linoid       = "linoid"
	Vultr        = "vultr"

	NodeStateNew        = NodeState(0)
	NodeStateActive     = NodeState(1)
	NodeStateUnhealthy  = NodeState(2)
	NodeStateTerminated = NodeState(4)
	NodeStateDeleted    = NodeState(8)

	ASGStateNew       = State(0)
	ASGStateActive    = State(1)
	ASGStateExecuting = State(2)
	ASGStateDeleted   = State(4)

	CMDStateNew        = CommandState(0)
	CMDStateInProgress = CommandState(1)
	CMDStateDone       = CommandState(2)
	CMDStateFailed     = CommandState(4)

	HealthMetricType MetricType = "health"
)

type (
	ID string

	MetricType string

	Provider struct {
		ID     string
		APIKey string
	}

	State int

	Order int

	NodeState    int
	CommandState int

	NetworkInterface struct {
		ID ID
		IP net.IP
	}
)

func NewCommandSet(cmd ...Command) CommandSet {
	return CommandSet{}
}

func (cs CommandSet) Merge(newSet CommandSet) {

}

func isRequiredMetric(metric Metric, typ MetricType) bool {
	ok := false

	switch typ {
	case HealthMetricType:
		_, ok = metric.(HealthMetric)
		return ok
	}

	return ok
}
