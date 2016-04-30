package domain

import "net"

const (
	DigitalOcean = "digitalocean"
	Linoid       = "linoid"
	Vultr        = "vultr"

	NodeStateNew        = NodeState(0)
	NodeStateUnhealthy  = NodeState(1)
	NodeStateActive     = NodeState(2)
	NodeStateTerminated = NodeState(4)
	NodeStateDeleted    = NodeState(8)

	ASGStateNew     = State(0)
	ASGStateActive  = State(1)
	ASGStateDeleted = State(2)

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

	NodeState int

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
