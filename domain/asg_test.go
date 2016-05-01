package domain

import (
	"net"
	"time"

	. "gopkg.in/check.v1"
)

type ASGSuite struct{}

var _ = Suite(&ASGSuite{})

func prepareMetrics(fail int, period int) MetricSeries {
	// Add metrics for lats 60 seconds
	now := time.Now()
	healthMetricSeries := NewMetricSeries()
	for i := 0; i < period-fail; i++ {
		b := now.Add(time.Duration(time.Second * time.Duration(-1*i)))
		healthMetricSeries[b] = NewHealthMetric(1, b)
	}

	// Simulate fails
	if fail > 0 {
		for i := period - fail; i < period; i++ {
			b := now.Add(time.Duration(time.Second * time.Duration(-1*i)))
			healthMetricSeries[b] = NewHealthMetric(0, b)
		}
	}

	return healthMetricSeries
}

func prepareNode(fail int, id ID) *Node {
	node := NewNode()
	node.Setup(
		id,
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
	healthMetricSeries := prepareMetrics(fail, 5)
	node.AddMetrics(healthMetricSeries)

	return node
}

func prepareNodes(fail int) NodeSet {
	nodeSet := NewNodeSet(prepareNode(fail, ID("node1")))
	return nodeSet
}

func (s *ASGSuite) TestIfASGEvaluatesThatOneNodeShouldBeLaunchedAndBeforeThatSameNodeShouldBeTerminated(c *C) {
	plc, err := NewDesiredNodeAmountPerProviderPolicy(ID("policy-1"), 1, 1, 1, 3, 0.7, time.Duration(-5*time.Second), Provider{
		ID:     DigitalOcean,
		APIKey: "some-key",
	})
	c.Assert(err, IsNil)
	policies := NewPolicySet(plc)

	// Two health check metrics will fail
	nodes := prepareNodes(2)
	asg := NewAutoScalingGroup(ID("asg-1"))
	err = asg.Setup(nodes, policies)
	c.Assert(err, IsNil)

	// Check should fail but as we set that 3 consecutive checks should fail, 2 more should fail
	err = asg.Evaluate()
	c.Assert(err, IsNil)

	// Simulate wait period, we do checks every 5 second
	time.Sleep(time.Second * 5)

	// we will check in 5 sec interval, so generate metrics for last 5 sec
	// fail once, should pass policy as we tolerate one failed data point (0.7 < 4/5)
	// this will reset consecutive checks to 0 again, so after next 3 fails it should fail
	m := prepareMetrics(1, 5)
	asg.AddMetrics(ID("node1"), m)
	err = asg.Evaluate()
	c.Assert(err, IsNil)

	// Simulate wait period, we do checks every 5 second
	time.Sleep(time.Second * 5)

	// Now we fail three times in a row
	m = prepareMetrics(3, 5)
	asg.AddMetrics(ID("node1"), m)
	err = asg.Evaluate()
	c.Assert(err, IsNil)

	// Simulate wait period, we do checks every 5 second
	time.Sleep(time.Second * 5)

	m = prepareMetrics(5, 5)
	asg.AddMetrics(ID("node1"), m)
	err = asg.Evaluate()
	c.Assert(err, IsNil)

	// Simulate wait period, we do checks every 5 second
	time.Sleep(time.Second * 5)

	m = prepareMetrics(3, 5)
	asg.AddMetrics(ID("node1"), m)
	err = asg.Evaluate()
	c.Assert(err, IsNil)

	c.Assert(len(asg.Commands), Equals, 2)
	c.Assert(asg.Commands[Order(1)], DeepEquals, &Terminate{
		BaseCommand: BaseCommand{
			Provider: Provider{
				ID:     DigitalOcean,
				APIKey: "some-key",
			},
		},
		NodeID: ID("node1"),
	})
	c.Assert(asg.Commands[Order(2)], DeepEquals, &Launch{
		BaseCommand: BaseCommand{
			Provider: Provider{
				ID:     DigitalOcean,
				APIKey: "some-key",
			},
		},
	})
}

func (s *ASGSuite) TestIfASGEvaluatesThatOneOfThreeDifferentNodeShouldBeLaunchedAndBeforeThatSameNodeShouldBeTerminates(c *C) {
	plc, err := NewDesiredNodeAmountPerProviderPolicy(ID("policy-1"), 1, 6, 3, 3, 0.7, time.Duration(-5*time.Second), Provider{
		ID:     DigitalOcean,
		APIKey: "some-key",
	})
	c.Assert(err, IsNil)
	policies := NewPolicySet(plc)

	node1 := prepareNode(2, ID("node1"))
	node2 := prepareNode(0, ID("node2"))
	node3 := prepareNode(0, ID("node3"))

	nodes := NewNodeSet(node1, node2, node3)

	asg := NewAutoScalingGroup(ID("asg-1"))
	err = asg.Setup(nodes, policies)
	c.Assert(err, IsNil)

	// Check should fail but as we set that 3 consecutive checks should fail, 2 more should fail
	err = asg.Evaluate()
	c.Assert(err, IsNil)

	// Simulate wait period, we do checks every 5 second
	time.Sleep(time.Second * 5)

	// we will check in 5 sec interval, so generate metrics for last 5 sec
	// fail once, should pass policy as we tolerate one failed data point (0.7 < 4/5)
	// this will reset consecutive checks to 0 again, so after next 3 fails it should fail
	m := prepareMetrics(1, 5)

	// Simulate not failing nodes
	healthyMetrics := prepareMetrics(0, 5)

	asg.AddMetrics(ID("node1"), m)
	asg.AddMetrics(ID("node2"), healthyMetrics)
	asg.AddMetrics(ID("node3"), healthyMetrics)
	err = asg.Evaluate()
	c.Assert(err, IsNil)

	// Simulate wait period, we do checks every 5 second
	time.Sleep(time.Second * 5)

	// Now we fail three times in a row
	m = prepareMetrics(3, 5)
	healthyMetrics = prepareMetrics(0, 5)
	asg.AddMetrics(ID("node1"), m)
	asg.AddMetrics(ID("node2"), healthyMetrics)
	asg.AddMetrics(ID("node3"), healthyMetrics)
	err = asg.Evaluate()
	c.Assert(err, IsNil)

	// Simulate wait period, we do checks every 5 second
	time.Sleep(time.Second * 5)

	m = prepareMetrics(5, 5)
	healthyMetrics = prepareMetrics(0, 5)
	asg.AddMetrics(ID("node1"), m)
	asg.AddMetrics(ID("node2"), healthyMetrics)
	asg.AddMetrics(ID("node3"), healthyMetrics)
	err = asg.Evaluate()
	c.Assert(err, IsNil)

	// Simulate wait period, we do checks every 5 second
	time.Sleep(time.Second * 5)

	m = prepareMetrics(3, 5)
	healthyMetrics = prepareMetrics(0, 5)
	asg.AddMetrics(ID("node1"), m)
	asg.AddMetrics(ID("node2"), healthyMetrics)
	asg.AddMetrics(ID("node3"), healthyMetrics)
	err = asg.Evaluate()
	c.Assert(err, IsNil)

	c.Assert(len(asg.Commands), Equals, 2)
	c.Assert(asg.Commands[Order(1)], DeepEquals, &Terminate{
		BaseCommand: BaseCommand{
			Provider: Provider{
				ID:     DigitalOcean,
				APIKey: "some-key",
			},
		},
		NodeID: ID("node1"),
	})
	c.Assert(asg.Commands[Order(2)], DeepEquals, &Launch{
		BaseCommand: BaseCommand{
			Provider: Provider{
				ID:     DigitalOcean,
				APIKey: "some-key",
			},
		},
	})
}
