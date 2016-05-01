package domain

import (
	"math"
	"time"
)

type (
	// Node type
	Node struct {
		ID            ID
		Provider      Provider
		PrivateIface  NetworkInterface
		PublicIface   NetworkInterface
		State         NodeState
		Metrics       MetricSeries
		KeepMetricFor time.Duration
	}

	// NodeSet set
	NodeSet map[ID]*Node
)

// NewNodeSet constructor
func NewNodeSet(nodes ...*Node) NodeSet {
	ns := NodeSet{}
	for _, n := range nodes {
		ns[n.ID] = n
	}
	return ns
}

// NewNode constructor
func NewNode() *Node {
	return &Node{
		State: NodeStateNew,
	}
}

// Create new node
func (n *Node) Setup(id ID, provider Provider, prIface, puIface NetworkInterface) error {
	n.State = NodeStateUnhealthy
	n.ID = id
	n.Provider = provider
	n.PrivateIface = prIface
	n.PublicIface = puIface
	n.KeepMetricFor = time.Minute * -1
	n.Metrics = NewMetricSeries()

	return nil
}

// ChangeProvider ...
func (n *Node) ChangeProvider(provider Provider) error {
	n.Provider = provider
	return nil
}

// ChangeState ...
func (n *Node) ChangeState(State NodeState) error {
	n.State = State
	return nil
}

// ChangeNetworkInterfaces ...
func (n *Node) ChangeNetworkInterfaces(prIface, puIface *NetworkInterface) error {
	if prIface != nil {
		n.PrivateIface = *prIface
	}

	if puIface != nil {
		n.PublicIface = *puIface
	}

	return nil
}

// Remove ...
func (n *Node) Remove() error {
	n.State = NodeStateDeleted
	return nil
}

// AddMetrics ...
func (n *Node) AddMetrics(metrics MetricSeries) error {
	n.clearMetrics()
	for t, m := range metrics {
		n.Metrics[t] = m
	}

	return nil
}

// CalculateMetricValue calculates avg of requested metric
func (n *Node) CalculateMetricValue(metricType MetricType, from, to time.Time) float64 {
	rez := 0.0
	value := 0.0
	dataPoints := 0
	for t, m := range n.Metrics {
		if isRequiredMetric(m, metricType) {
			if t.After(from) && t.Before(to) {
				value = value + m.GetValue()
				dataPoints = dataPoints + 1
			}
		}
	}

	if value > 0 {
		rez = value / float64(dataPoints)
	}

	return round(rez, .5, 2)
}

// clearMetrics removes metrics that are older than n.KeepMetricFor
func (n *Node) clearMetrics() error {
	requiredTime := time.Now().Add(n.KeepMetricFor)
	for t, _ := range n.Metrics {
		if t.Before(requiredTime) {
			delete(n.Metrics, t)
		}
	}

	return nil
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
