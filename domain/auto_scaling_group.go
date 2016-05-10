package domain

import (
	"sort"

	"strings"

	"time"

	"fmt"

	"github.com/juju/errors"
)

type (
	// AutoScalingGroup ...
	AutoScalingGroup struct {
		ID       ID
		State    State
		Nodes    NodeSet
		Policies PolicySet
		Commands CommandSet

		stop bool
	}

	// AutoScalingGroupSet type
	AutoScalingGroupSet map[ID]*AutoScalingGroup
)

// NewAutoScalingGroup constructor
func NewAutoScalingGroup(id ID) *AutoScalingGroup {
	return &AutoScalingGroup{
		ID:    id,
		State: ASGStateNew,
	}
}

// NewAutoScalingGroupSet constructor
func NewAutoScalingGroupSet(asgs ...*AutoScalingGroup) AutoScalingGroupSet {
	asgsg := AutoScalingGroupSet{}
	for _, a := range asgs {
		asgsg[a.ID] = a
	}

	return asgsg
}

// Setup ...
func (asg *AutoScalingGroup) Setup(nodes NodeSet, policies PolicySet) error {
	asg.State = ASGStateActive
	asg.Nodes = nodes
	asg.Policies = policies
	asg.Commands = NewCommandSet()

	return nil
}

// ChangePolicy changes the policy, currently it overrides it and all state is lost
func (asg *AutoScalingGroup) ChangePolicy(policy Policy) error {
	if asg.State == ASGStateNew {
		return errors.Errorf("ASG is in ASGStateNew state, use Setup() first!")
	}

	err := asg.Policies.Update(policy)

	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

// AddMetrics ...
func (asg *AutoScalingGroup) AddMetrics(node ID, metrics MetricSeries) error {
	if asg.State == ASGStateNew {
		return errors.Errorf("ASG is in ASGStateNew state, use Setup() first!")
	}

	if _, ok := asg.Nodes[node]; !ok {
		errors.Errorf("Node by ID %s was not found", node)
	}

	asg.Nodes[node].AddMetrics(metrics)
	return nil
}

// RemoveNode ...
func (asg *AutoScalingGroup) RemoveNode(node ID) error {
	if asg.State == ASGStateNew {
		return errors.Errorf("ASG is in ASGStateNew state, use Setup() first!")
	}

	if _, ok := asg.Nodes[node]; !ok {
		errors.Errorf("Node by ID %s was not found", node)
	}

	fmt.Printf("Remove node [%s] \n", node)
	delete(asg.Nodes, node)
	return nil
}

// AddNode ...
func (asg *AutoScalingGroup) AddNode(node *Node) error {
	if asg.State == ASGStateNew {
		return errors.Errorf("ASG is in ASGStateNew state, use Setup() first!")
	}

	if _, ok := asg.Nodes[node.ID]; !ok {
		errors.Errorf("Node with ID %s already exists", node.ID)
	}

	fmt.Printf("Add node [%s] \n", node.ID)
	asg.Nodes[node.ID] = node
	return nil
}

// Evaluate goes through the policies and generates commands that needs
// to be executed
func (asg *AutoScalingGroup) Evaluate() error {
	if asg.State == ASGStateExecuting {
		return errors.Errorf("ASG is in ASGStateExecuting state, try later")
	}

	if asg.State == ASGStateNew {
		return errors.Errorf("ASG is in ASGStateNew state, use Setup() first!")
	}

	fmt.Printf("Nodes: %d \n", len(asg.Nodes))

	for _, policy := range asg.Policies {
		err := policy.Evaluate(asg)
		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

// Execute required commands created by policies
func (asg *AutoScalingGroup) Execute() error {
	if asg.State == ASGStateNew {
		return errors.Errorf("ASG is in ASGStateNew state, use Setup() first!")
	}

	fmt.Printf("Commands: %d \n", len(asg.Commands))
	asg.State = ASGStateExecuting
	var keys []int
	keys = []int{}
	for k := range asg.Commands {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	errs := []string{}
	for _, k := range keys {
		fmt.Printf("Command: %+v \n", asg.Commands[Order(k)])
		// If case of error we add it to slice of errors
		// and we do move on
		// Commands are atomic and if one fails it should not influence others
		err := asg.Commands[Order(k)].Execute(asg)

		// how to deal with this ?
		// we just return what has failed
		if err != nil {
			errs = append(errs, err.Error())
		}

		delete(asg.Commands, Order(k))
	}

	if len(errs) > 0 {
		return errors.Errorf("Execution finished with these errors - %s", strings.Join(errs, ":"))
	}

	asg.State = ASGStateActive
	return nil
}

// Run ASG and start monitoring nodes
func (asg *AutoScalingGroup) Run() error {
	for {
		// Stop
		if asg.stop {
			return nil
		}

		err := asg.Evaluate()
		if err != nil {
			return err
		}
		err = asg.Execute()

		if err != nil {
			return err
		}
		time.Sleep(time.Second * 5)
		fmt.Printf("[%s] OK \n", asg.ID)
	}
}

// Stop ASG
func (asg *AutoScalingGroup) Stop() error {
	asg.stop = true
	return nil
}

// Remove ...
func (asg *AutoScalingGroup) Remove() error {
	if asg.State == ASGStateNew {
		return errors.Errorf("ASG is in ASGStateNew state, use Setup() first!")
	}

	asg.State = ASGStateDeleted
	return nil
}
