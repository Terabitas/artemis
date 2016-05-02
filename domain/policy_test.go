package domain

import (
	"net"
	"time"

	. "gopkg.in/check.v1"
)

type PolicySuite struct{}

var _ = Suite(&PolicySuite{})

func prepareAsg(fail int) *AutoScalingGroup {
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
	for i := 0; i < 60-fail; i++ {
		b := now.Add(time.Duration(time.Second * time.Duration(-1*i)))
		healthMetricSeries[b] = NewHealthMetric(1, b)
	}

	// Simulate fails
	if fail > 0 {
		for i := 60 - fail; i <= 60; i++ {
			b := now.Add(time.Duration(time.Second * time.Duration(-1*i)))
			healthMetricSeries[b] = NewHealthMetric(0, b)
		}
	}

	node.AddMetrics(healthMetricSeries)

	nodeSet := NodeSet{
		ID("node1"): node,
	}

	asg := NewAutoScalingGroup(ID("test"))
	asg.Setup(nodeSet, PolicySet{})

	return asg
}

func (s *PolicySuite) TestIfDesiredNodeAmountPerProviderPolicyEvaluatesCorrectly(c *C) {
	asg := prepareAsg(0)

	plc, err := NewDesiredNodeAmountPerProviderPolicy(ID("policy-1"), 1, 1, 1, 3, 0.9, time.Duration(-60*time.Second), Provider{
		ID:     DigitalOcean,
		APIKey: "some-key",
	})
	c.Assert(err, IsNil)
	plc.Evaluate(asg)

	c.Assert(len(asg.Commands), Equals, 0)

	// 55 of 60 - should not fail
	asg = prepareAsg(5)

	plc, err = NewDesiredNodeAmountPerProviderPolicy(ID("policy-1"), 1, 1, 1, 3, 0.9, time.Duration(-60*time.Second), Provider{
		ID:     DigitalOcean,
		APIKey: "some-key",
	})
	c.Assert(err, IsNil)
	plc.Evaluate(asg)

	c.Assert(len(asg.Commands), Equals, 0)

	// 50 of 60 should be less than 0.9 and because we allow only one check to fail it should launch new instance
	// terminate old
	asg = prepareAsg(10)

	plc, err = NewDesiredNodeAmountPerProviderPolicy(ID("policy-1"), 1, 1, 1, 1, 0.9, time.Duration(-60*time.Second), Provider{
		ID:     DigitalOcean,
		APIKey: "some-key",
	})
	c.Assert(err, IsNil)
	plc.Evaluate(asg)

	c.Assert(len(asg.Commands), Equals, 1)
	c.Assert(asg.Commands[Order(1)], DeepEquals, &Relaunch{
		BaseCommand: BaseCommand{
			Provider: Provider{
				ID:     DigitalOcean,
				APIKey: "some-key",
			},
		},
		NodeID: ID("node1"),
	})
}
