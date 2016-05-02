package domain

import (
	"net"
	"time"
)

type (
	CommandSet map[Order]Command

	Command interface {
		Execute(*AutoScalingGroup) error
	}

	CMDError struct {
		Code    int
		Message string
	}

	BaseCommand struct {
		Provider Provider
		State    CommandState
		Error    *CMDError
		Timeout  time.Duration
	}

	BaseCommands []BaseCommand

	Launch struct {
		BaseCommand
	}

	Terminate struct {
		BaseCommand

		NodeID ID
	}

	Relaunch struct {
		BaseCommand

		NodeID ID
	}
)

func (lc *Launch) Execute(asg *AutoScalingGroup) error {

	// API call

	return nil
}

func (lc *Terminate) Execute(asg *AutoScalingGroup) error {

	// API call

	// retry here
	asg.RemoveNode(lc.NodeID)

	return nil
}

func (lc *Relaunch) Execute(asg *AutoScalingGroup) error {

	// API call

	// retry here

	node := NewNode()

	node.Setup(
		ID("new-node"),
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

	asg.AddNode(node)

	// only when new is added, remove old
	asg.RemoveNode(lc.NodeID)

	return nil
}
