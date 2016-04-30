package domain

type (
	AutoScalingGroup struct {
		ID       ID
		State    State
		Nodes    NodeSet
		Policies Policies
		Commands CommandSet
	}
)

func NewAutoScalingGroup(id ID) *AutoScalingGroup {
	return &AutoScalingGroup{
		ID: id,
	}
}

func (asg *AutoScalingGroup) Create(nodes NodeSet, policies Policies) error {

	asg.State = ASGStateActive
	asg.Nodes = nodes
	asg.Policies = policies
	asg.Commands = NewCommandSet()

	return nil
}

func (asg *AutoScalingGroup) AddMetrics(node ID, metrics MetricSeries) error {

	return nil
}

func (asg *AutoScalingGroup) Execute() error {

	asg.State = ASGStateDeleted

	for _, policy := range asg.Policies {
		policy.Evaluate(asg)
	}

	for _, cmd := range asg.Commands {
		cmd.Execute()
	}

	return nil
}

func (asg *AutoScalingGroup) Remove() error {

	asg.State = ASGStateDeleted

	return nil
}
