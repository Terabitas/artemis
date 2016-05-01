package domain

import "github.com/juju/errors"

type (
	// AutoScalingGroup ...
	AutoScalingGroup struct {
		ID       ID
		State    State
		Nodes    NodeSet
		Policies PolicySet
		Commands CommandSet
	}
)

// NewAutoScalingGroup constructor
func NewAutoScalingGroup(id ID) *AutoScalingGroup {
	return &AutoScalingGroup{
		ID:    id,
		State: ASGStateNew,
	}
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
	err := asg.Policies.Replace(policy)

	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

// AddMetrics ...
func (asg *AutoScalingGroup) AddMetrics(node ID, metrics MetricSeries) error {

	if _, ok := asg.Nodes[node]; !ok {
		errors.Errorf("Node by ID %s was not found", node)
	}

	asg.Nodes[node].AddMetrics(metrics)
	return nil
}

// Evaluate goes through the policies and generates commands that needs
// to be executed
func (asg *AutoScalingGroup) Evaluate() error {
	asg.State = ASGStateDeleted

	for _, policy := range asg.Policies {
		err := policy.Evaluate(asg)
		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

// Remove ...
func (asg *AutoScalingGroup) Remove() error {
	asg.State = ASGStateDeleted
	return nil
}
