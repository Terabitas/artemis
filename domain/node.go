package domain

import (
	"math"
	"time"
)

type (
	Node struct {
		ID            ID
		Provider      Provider
		PrivateIface  NetworkInterface
		PublicIface   NetworkInterface
		State         NodeState
		Metrics       MetricSeries
		KeepMetricFor time.Duration
	}

	NodeSet map[ID]*Node
)

func NewNode() *Node {
	return &Node{
		State: NodeStateNew,
	}
}

func (n *Node) Create(id ID, provider Provider, prIface, puIface NetworkInterface) error {

	n.State = NodeStateNew
	n.ID = id
	n.Provider = provider
	n.PrivateIface = prIface
	n.PublicIface = puIface
	n.KeepMetricFor = time.Second * 60
	n.Metrics = NewMetricSeries()

	return nil
}

func (n *Node) ChangeProvider(provider Provider) error {

	n.Provider = provider

	return nil
}

func (n *Node) ChangeNetworkInterfaces(prIface, puIface *NetworkInterface) error {

	if prIface != nil {
		n.PrivateIface = *prIface
	}

	if puIface != nil {
		n.PublicIface = *puIface
	}

	return nil
}

func (n *Node) Remove() error {

	n.State = NodeStateDeleted

	return nil
}

func (n *Node) AddMetrics(metrics MetricSeries) error {
	n.clearMetrics()
	for t, m := range metrics {
		n.Metrics[t] = m
	}

	return nil
}

func (n *Node) clearMetrics() error {
	requiredTime := time.Now().Truncate(n.KeepMetricFor)
	for t, _ := range n.Metrics {
		if t.Before(requiredTime) {
			delete(n.Metrics, t)
		}
	}

	return nil
}

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

	return Round(rez, .5, 2)
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
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
