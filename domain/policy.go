package domain

import (
	"time"

	"github.com/juju/errors"
)

type (
	Policies []Policy

	Policy interface {
		// Evaluate takes in AutoScalingGroup and calculates what commands should be added or removed
		Evaluate(*AutoScalingGroup) error
	}

	// DesiredNodeAmountPerProviderPolicy evaluates current state and creates Commands per provider
	DesiredHealthyNodeAmountPerProviderPolicy struct {
		Min, Max, Desired, Current int
		HealthyThreshold           float64
		CheckInterval              time.Duration
		Provider                   string
	}
)

// NewDesiredNodeAmountPerProviderPolicy constructor
func NewDesiredNodeAmountPerProviderPolicy(min, max, desired int, healthyThreshold float64, checkInterval time.Duration, providerID string) (Policy, error) {

	if desired > max {
		return nil, errors.Errorf("Desired %d can not be more than max %d", desired, max)
	}

	if desired < min {
		return nil, errors.Errorf("Desired %d can not be less than min %d", desired, min)
	}

	if min > max {
		return nil, errors.Errorf("Min %d can not be more than max %d", min, max)
	}

	return &DesiredHealthyNodeAmountPerProviderPolicy{
		Min:              min,
		Max:              max,
		Desired:          desired,
		Current:          0,
		HealthyThreshold: healthyThreshold,
		CheckInterval:    checkInterval,
		Provider:         providerID,
	}, nil
}

func (dsp *DesiredHealthyNodeAmountPerProviderPolicy) Evaluate(asg *AutoScalingGroup) error {

	dsp.countCurrent(asg.Nodes)

	if dsp.Current == dsp.Desired {
		return nil
	}

	if dsp.Current < dsp.Desired {
		amt := dsp.Desired - dsp.Current
		for i := 0; i < amt; i++ {
			asg.Commands[Order(i)] = &Launch{}
		}
	}

	if dsp.Current > dsp.Desired {
		amt := dsp.Current - dsp.Desired
		for i := 0; i < amt; i++ {
			asg.Commands[Order(i)] = &Terminate{}
		}
	}

	return nil
}

func (dsp *DesiredHealthyNodeAmountPerProviderPolicy) countCurrent(nodes NodeSet) error {
	for _, node := range nodes {
		if node.Provider.ID != dsp.Provider {
			continue
		}
		now := time.Now()
		before := now.Add(dsp.CheckInterval)
		val := node.CalculateMetricValue(HealthMetricType, before, now)

		if val >= dsp.HealthyThreshold {
			dsp.Current = dsp.Current + 1
		}
	}

	return nil
}
