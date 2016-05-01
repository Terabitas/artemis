package domain

import (
	"net"
	"time"

	. "gopkg.in/check.v1"
)

type NodeSuite struct{}

var _ = Suite(&NodeSuite{})

func (s *NodeSuite) TestIfNodeMetricsAreAddedCorrectly(c *C) {
	node := NewNode()
	node.Setup(
		ID("node1"),
		Provider{
			ID:     DigitalOcean,
			APIKey: "some-key",
		},
		NetworkInterface{
			ID: ID("eth0"),
			IP: net.ParseIP("192.100.10.1"),
		},
		NetworkInterface{
			ID: ID("eth0"),
			IP: net.ParseIP("192.100.10.2"),
		},
	)
	// Add metrics for lats 60 seconds
	now := time.Now()
	healthMetricSeries := MetricSeries{}
	for i := 0; i < 60-5; i++ {
		b := now.Add(time.Duration(time.Second * time.Duration(-1*i)))
		healthMetricSeries[b] = NewHealthMetric(1, b)
	}

	for i := 55; i < 60; i++ {
		b := now.Add(time.Duration(time.Second * time.Duration(-1*i)))
		healthMetricSeries[b] = NewHealthMetric(0, b)
	}

	node.AddMetrics(healthMetricSeries)

	c.Assert(len(node.Metrics), Equals, 60)
}

func (s *NodeSuite) TestIfNodeMetricsAreCleared(c *C) {
	node := NewNode()
	node.Setup(
		ID("node1"),
		Provider{
			ID:     DigitalOcean,
			APIKey: "some-key",
		},
		NetworkInterface{
			ID: ID("eth0"),
			IP: net.ParseIP("192.100.10.1"),
		},
		NetworkInterface{
			ID: ID("eth0"),
			IP: net.ParseIP("192.100.10.2"),
		},
	)
	// Add metrics for last 60 seconds
	now := time.Now()
	healthMetricSeries := MetricSeries{}
	for i := 0; i < 60; i++ {
		b := now.Add(time.Duration(time.Second * time.Duration(-1*i)))
		healthMetricSeries[b] = NewHealthMetric(1, b)
	}

	node.AddMetrics(healthMetricSeries)

	time.Sleep(time.Second * 3)

	// Add metrics for lats 60 seconds
	now = time.Now()
	required := now.Add(-1 * time.Minute)
	healthMetricSeries = MetricSeries{}
	for i := 0; i < 60; i++ {
		b := now.Add(time.Duration(time.Second * time.Duration(-1*i)))
		healthMetricSeries[b] = NewHealthMetric(1, b)
	}

	node.AddMetrics(healthMetricSeries)

	for _, m := range node.Metrics {
		c.Assert(m.GetTimestamp().After(required), Equals, true)
	}
}
