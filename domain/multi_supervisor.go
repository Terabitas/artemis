package domain

import "sync"

var wg sync.WaitGroup

type (
	// MultiSupervisor manages multiple ASG
	MultiSupervisor struct {
		sync.RWMutex
		autoScalingGroups AutoScalingGroupSet
		stop              chan bool
	}
)

// MakeMultiSupervisor creates new supervisor that supports multiple ASG
func MakeMultiSupervisor() *MultiSupervisor {
	return &MultiSupervisor{
		autoScalingGroups: NewAutoScalingGroupSet(),
		stop:              make(chan bool),
	}
}

// Run supervisor
func (s *MultiSupervisor) Run() {
	// for main routine
	wg.Add(1)

	// wait for everything to exit
	wg.Wait()
}

// Add new ASG to run
func (s *MultiSupervisor) Add(asg *AutoScalingGroup) {
	s.runASG(asg)
}

// Get ASG
func (s *MultiSupervisor) Get(id ID) *AutoScalingGroup {
	if !s.exists(id) {
		return nil
	}

	return s.get(id)
}

// Remove ASG
func (s *MultiSupervisor) Remove(id ID) {
	if s.exists(id) {
		s.get(id).Stop()
		delete(s.autoScalingGroups, id)
	}
}

// Private stuff

func (s *MultiSupervisor) runASG(asg *AutoScalingGroup) {
	did := asg.ID

	if s.exists(asg.ID) {
		// already running
		return
	}

	// add to map
	s.add(asg)

	// Starting actual ASG routine
	// `asg.Run()` will run until ASG is stopped or removed
	go func(asg *AutoScalingGroup) {
		asg.Run()
		wg.Done()
	}(s.get(did))
}

func (s *MultiSupervisor) exists(id ID) bool {
	_, ok := s.autoScalingGroups[id]
	return ok
}

func (s *MultiSupervisor) get(id ID) *AutoScalingGroup {
	if s.exists(id) {
		return s.autoScalingGroups[id]
	}
	return nil
}

func (s *MultiSupervisor) add(asg *AutoScalingGroup) {
	wg.Add(1)

	s.Lock()
	s.autoScalingGroups[asg.ID] = asg
	s.Unlock()
}
