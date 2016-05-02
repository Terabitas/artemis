package domain

import (
	"time"

	"github.com/juju/errors"
)

type (
	// Policies list
	PolicySet map[ID]Policy

	// Policy type
	Policy interface {
		// Evaluate takes in AutoScalingGroup and calculates what commands should be added or removed
		Evaluate(*AutoScalingGroup) error
		Update(Policy) error
		GetID() ID
	}

	// DesiredNodeAmountPerProviderPolicy evaluates current state and creates Commands per provider
	DesiredHealthyNodeAmountPerProviderPolicy struct {
		ID                         ID
		Min, Max, Desired, Current int
		HealthyThreshold           float64
		CheckInterval              time.Duration
		Provider                   Provider
		ConsecutiveChecks          int
		ConsecutiveChecksNum       map[ID]int
	}
)

// NewPolicySet constructor
func NewPolicySet(plc ...Policy) PolicySet {
	policies := PolicySet{}
	for _, p := range plc {
		policies[p.GetID()] = p
	}
	return policies
}

// NewDesiredNodeAmountPerProviderPolicy constructor
func NewDesiredNodeAmountPerProviderPolicy(id ID, min, max, desired, consecutiveChecks int, healthyThreshold float64, checkInterval time.Duration, provider Provider) (Policy, error) {

	if desired > max {
		return nil, errors.Errorf("Desired %d can not be more than max %d", desired, max)
	}

	if desired < min {
		return nil, errors.Errorf("Desired %d can not be less than min %d", desired, min)
	}

	if min > max {
		return nil, errors.Errorf("Min %d can not be more than max %d", min, max)
	}

	if consecutiveChecks <= 0 {
		return nil, errors.Errorf("ConsecutiveChecks %d can not be less or equal to 0", consecutiveChecks)
	}

	return &DesiredHealthyNodeAmountPerProviderPolicy{
		ID:                   id,
		Min:                  min,
		Max:                  max,
		Desired:              desired,
		Current:              0,
		ConsecutiveChecksNum: map[ID]int{},
		HealthyThreshold:     healthyThreshold,
		CheckInterval:        checkInterval,
		Provider:             provider,
		ConsecutiveChecks:    consecutiveChecks,
	}, nil
}

func (dsp *DesiredHealthyNodeAmountPerProviderPolicy) GetID() ID {
	return dsp.ID
}

// Update will reset checks state
func (dsp *DesiredHealthyNodeAmountPerProviderPolicy) Update(plc Policy) error {
	v, ok := plc.(*DesiredHealthyNodeAmountPerProviderPolicy)
	if !ok {
		return errors.Errorf("Given policy is not *DesiredHealthyNodeAmountPerProviderPolicy")
	}

	dsp.Desired = v.Desired
	dsp.Max = v.Max
	dsp.Min = v.Min
	dsp.HealthyThreshold = v.HealthyThreshold
	dsp.CheckInterval = v.CheckInterval
	dsp.Provider = v.Provider
	dsp.ConsecutiveChecks = v.ConsecutiveChecks
	dsp.ConsecutiveChecksNum = map[ID]int{}

	return nil
}

// Evaluate what commands should be executed by given ASG
func (dsp *DesiredHealthyNodeAmountPerProviderPolicy) Evaluate(asg *AutoScalingGroup) error {

	dsp.Current = 0
	dsp.countCurrent(asg.Nodes)

	if dsp.Current == dsp.Desired {
		return nil
	}

	commandOrder := len(asg.Commands)
	if dsp.Current < dsp.Desired {
		amt := dsp.Desired - dsp.Current

		// Relaunch nodes
		handled := 0
		for nodeID, v := range dsp.ConsecutiveChecksNum {
			// Relaunch those nodes which has failed checks
			if v == dsp.ConsecutiveChecks {
				commandOrder++
				asg.Commands[Order(commandOrder)] = &Relaunch{
					BaseCommand: BaseCommand{
						Provider: dsp.Provider,
					},
					NodeID: nodeID,
				}

				handled++
			}
		}

		// Now of desired has been increased so even though all nodes
		// are healthy we need to launch new ones
		if amt > handled {
			for i := 0; i < amt-handled; i++ {
				commandOrder++
				asg.Commands[Order(commandOrder)] = &Launch{
					BaseCommand: BaseCommand{
						Provider: dsp.Provider,
					},
				}
			}
		}
	}

	// If desired has been minimized, terminate the difference
	if dsp.Current > dsp.Desired {
		amt := dsp.Current - dsp.Desired

		handled := 0
		for nodeID := range asg.Nodes {
			if handled == amt {
				break
			}
			commandOrder++
			asg.Commands[Order(commandOrder)] = &Terminate{
				BaseCommand: BaseCommand{
					Provider: dsp.Provider,
				},
				NodeID: nodeID,
			}
			handled++
		}
	}

	return nil
}

func (dsp *DesiredHealthyNodeAmountPerProviderPolicy) countCurrent(nodes NodeSet) error {
	for _, node := range nodes {
		if _, ok := dsp.ConsecutiveChecksNum[node.ID]; !ok {
			dsp.ConsecutiveChecksNum[node.ID] = 0
		}

		if node.Provider.ID != dsp.Provider.ID {
			continue
		}

		now := time.Now()
		before := now.Add(dsp.CheckInterval)
		val := node.CalculateMetricValue(HealthMetricType, before, now)

		if val >= dsp.HealthyThreshold {
			node.ChangeState(NodeStateActive)
			dsp.Current++
			// reset
			dsp.ConsecutiveChecksNum[node.ID] = 0
		} else {
			dsp.ConsecutiveChecksNum[node.ID]++
			node.ChangeState(NodeStateUnhealthy)
			if dsp.ConsecutiveChecksNum[node.ID] < dsp.ConsecutiveChecks {
				dsp.Current++
				node.ChangeState(NodeStateActive)
			}
		}
	}

	return nil
}

// Update policy
func (p PolicySet) Update(policy Policy) error {
	if _, ok := p[policy.GetID()]; !ok {
		return errors.Errorf("Policy by id %s was not found", policy.GetID())
	}

	return p[policy.GetID()].Update(policy)
}
