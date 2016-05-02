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

	node := NewNode()
	node.Setup(
		ID("new-node-instance-id-from-API"),
		Provider{
			ID:     DigitalOcean,
			APIKey: "some-key",
		},
		NetworkInterface{
			ID: ID("eth0-from-API"),
			IP: net.ParseIP("192.100.10.1"),
		},
		NetworkInterface{
			ID: ID("eth1-from-API"),
			IP: net.ParseIP("192.100.10.2"),
		},
	)

	// Add new node
	asg.AddNode(node)

	return nil
}

func (lc *Terminate) Execute(asg *AutoScalingGroup) error {

	// API call

	// retry here
	asg.RemoveNode(lc.NodeID)

	return nil
}

func (lc *Relaunch) Execute(asg *AutoScalingGroup) error {

	// API calls here

	node := NewNode()
	node.Setup(
		ID("new-node-instance-id-from-API"),
		Provider{
			ID:     DigitalOcean,
			APIKey: "some-key",
		},
		NetworkInterface{
			ID: ID("eth0-from-API"),
			IP: net.ParseIP("192.100.10.1"),
		},
		NetworkInterface{
			ID: ID("eth1-from-API"),
			IP: net.ParseIP("192.100.10.2"),
		},
	)

	// Add new node
	asg.AddNode(node)

	// oOnly when new is added, remove old
	asg.RemoveNode(lc.NodeID)

	return nil
}
