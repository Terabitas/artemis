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

func prepareNode(fail int, id ID, period int) *Node {
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
	healthMetricSeries := prepareMetrics(fail, period)
	node.AddMetrics(healthMetricSeries)

	return node
}

func prepareNodes(fail int) NodeSet {
	nodeSet := NewNodeSet(prepareNode(fail, ID("node1"), 5))
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

func (s *ASGSuite) TestIfASGEvaluatesThatOneOfThreeDifferentNodeShouldBeLaunchedAndBeforeThatSameNodeShouldBeTerminates(c *C) {
	plc, err := NewDesiredNodeAmountPerProviderPolicy(ID("policy-1"), 1, 6, 3, 3, 0.7, time.Duration(-5*time.Second), Provider{
		ID:     DigitalOcean,
		APIKey: "some-key",
	})
	c.Assert(err, IsNil)
	policies := NewPolicySet(plc)

	node1 := prepareNode(2, ID("node1"), 5)
	node2 := prepareNode(0, ID("node2"), 5)
	node3 := prepareNode(0, ID("node3"), 5)

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

func (s *ASGSuite) TestIfNodeIsRemovedSuccessfully(c *C) {
	policies := NewPolicySet()

	node1 := prepareNode(2, ID("node1"), 5)
	node2 := prepareNode(0, ID("node2"), 5)
	node3 := prepareNode(0, ID("node3"), 5)

	nodes := NewNodeSet(node1, node2, node3)

	asg := NewAutoScalingGroup(ID("asg-1"))
	err := asg.Setup(nodes, policies)
	c.Assert(err, IsNil)

	c.Assert(len(asg.Nodes), Equals, 3)

	err = asg.RemoveNode(ID("node2"))
	c.Assert(err, IsNil)

	c.Assert(len(asg.Nodes), Equals, 2)
}

func (s *ASGSuite) TestIfNodeIsAddedSuccessfully(c *C) {
	policies := NewPolicySet()

	node1 := prepareNode(2, ID("node1"), 5)
	node2 := prepareNode(0, ID("node2"), 5)
	node3 := prepareNode(0, ID("node3"), 5)

	nodes := NewNodeSet(node1, node2)

	asg := NewAutoScalingGroup(ID("asg-1"))
	err := asg.Setup(nodes, policies)
	c.Assert(err, IsNil)

	c.Assert(len(asg.Nodes), Equals, 2)

	err = asg.AddNode(node3)
	c.Assert(err, IsNil)

	c.Assert(len(asg.Nodes), Equals, 3)
	c.Assert(asg.Nodes.GetByID(ID("node3")), DeepEquals, node3)
}

func (s *ASGSuite) TestIfASGRemovesFailingNodeAndReplacesItWithNewOne(c *C) {
	plc, err := NewDesiredNodeAmountPerProviderPolicy(ID("policy-1"), 1, 1, 1, 1, 1, time.Duration(-5*time.Second), Provider{
		ID:     DigitalOcean,
		APIKey: "some-key",
	})

	c.Assert(err, IsNil)
	policies := NewPolicySet(plc)

	node1 := prepareNode(0, ID("node1"), 0)

	nodes := NewNodeSet(node1)

	asg := NewAutoScalingGroup(ID("asg-1"))
	err = asg.Setup(nodes, policies)
	c.Assert(err, IsNil)

	// At this point we have prepared ASG with one node which has no metrics yet
	// Let's generate metrics for last 5 sec and register for our node1
	healthyMetrics := prepareMetrics(0, 5)
	err = asg.AddMetrics(ID("node1"), healthyMetrics)
	c.Assert(err, IsNil)

	// We evaluate ASG, everything should be fine
	err = asg.Evaluate()
	c.Assert(err, IsNil)

	// We wait 5 sec
	time.Sleep(time.Second * 5)

	// Now we add generate metrics which indicates that node is not healthy
	// this should result in command to relaunch
	unhealthyMetrics := prepareMetrics(5, 5)
	err = asg.AddMetrics(ID("node1"), unhealthyMetrics)
	c.Assert(err, IsNil)

	// We evaluate ASG, it should evaluate to relaunch cmd
	err = asg.Evaluate()
	c.Assert(err, IsNil)

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

	err = asg.Execute()
	c.Assert(err, IsNil)
	c.Assert(len(asg.Commands), Equals, 0)
	c.Assert(len(asg.Nodes), Equals, 1)
	c.Assert(asg.Nodes.GetByID(ID("node1")), IsNil)
}

func (s *ASGSuite) TestIfASGAddsUpNewNodesWhenDesiredIsIncreasedAfterInitialSetup(c *C) {
	plc, err := NewDesiredNodeAmountPerProviderPolicy(ID("policy-1"), 1, 1, 1, 1, 1, time.Duration(-5*time.Second), Provider{
		ID:     DigitalOcean,
		APIKey: "some-key",
	})

	c.Assert(err, IsNil)
	policies := NewPolicySet(plc)

	node1 := prepareNode(0, ID("node1"), 0)

	nodes := NewNodeSet(node1)

	asg := NewAutoScalingGroup(ID("asg-1"))
	err = asg.Setup(nodes, policies)
	c.Assert(err, IsNil)

	// At this point we have prepared ASG with one node which has no metrics yet
	// Let's generate metrics for last 5 sec and register for our node1
	healthyMetrics := prepareMetrics(0, 5)
	err = asg.AddMetrics(ID("node1"), healthyMetrics)
	c.Assert(err, IsNil)

	// We evaluate ASG, everything should be fine
	err = asg.Evaluate()
	c.Assert(err, IsNil)

	// We wait 5 sec
	time.Sleep(time.Second * 5)

	// Lets add metrics to simulate that other node is fine
	healthyMetrics = prepareMetrics(0, 5)
	err = asg.AddMetrics(ID("node1"), healthyMetrics)
	c.Assert(err, IsNil)

	// We change policy now and do increase desired and max by 1
	plc, err = NewDesiredNodeAmountPerProviderPolicy(ID("policy-1"), 1, 2, 2, 1, 1, time.Duration(-5*time.Second), Provider{
		ID:     DigitalOcean,
		APIKey: "some-key",
	})
	c.Assert(err, IsNil)

	err = asg.ChangePolicy(plc)
	c.Assert(err, IsNil)

	// We evaluate ASG, it should result in one new Launch command
	err = asg.Evaluate()
	c.Assert(err, IsNil)

	c.Assert(len(asg.Commands), Equals, 1)
	c.Assert(asg.Commands[Order(1)], DeepEquals, &Launch{
		BaseCommand: BaseCommand{
			Provider: Provider{
				ID:     DigitalOcean,
				APIKey: "some-key",
			},
		},
	})

	// After execution we should have two nodes
	err = asg.Execute()
	c.Assert(err, IsNil)
	c.Assert(len(asg.Commands), Equals, 0)
	c.Assert(len(asg.Nodes), Equals, 2)
}
