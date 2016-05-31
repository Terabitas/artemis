package domain

import (
	"fmt"
	"net"
)

type (
	LaunchLocal struct {
		BaseCommand
	}

	TerminateLocal struct {
		BaseCommand

		NodeID ID
	}

	RelaunchLocal struct {
		BaseCommand

		NodeID ID
	}
)

func (lc *LaunchLocal) Execute(asg *AutoScalingGroup) error {
	publicIP := "1.1.1.2"
	privateIP := "1.1.1.1"

	fmt.Printf("Setting up new Node for [%s] \n", asg.ID)
	// API call
	node := NewNode()
	node.Setup(
		ID("droplet-1"),
		lc.Provider,
		NetworkInterface{
			IP: net.ParseIP(privateIP),
		},
		NetworkInterface{
			IP: net.ParseIP(publicIP),
		},
	)

	// Add new node
	asg.AddNode(node)

	return nil
}

func (lc *TerminateLocal) Execute(asg *AutoScalingGroup) error {

	asg.RemoveNode("drople-1")

	fmt.Printf("Execute Terminate [%s] \n", asg.ID)

	return nil
}

func (lc *RelaunchLocal) Execute(asg *AutoScalingGroup) error {

	publicIP := "1.1.1.2"
	privateIP := "1.1.1.1"

	fmt.Printf("Setting up new Node for [%s] \n", asg.ID)
	// API call
	node := NewNode()
	node.Setup(
		ID("droplet-1"),
		lc.Provider,
		NetworkInterface{
			IP: net.ParseIP(privateIP),
		},
		NetworkInterface{
			IP: net.ParseIP(publicIP),
		},
	)

	// Add new node
	asg.AddNode(node)

	asg.RemoveNode("droplet-1")

	fmt.Printf("Execute Relaunch [%s] \n", asg.ID)

	return nil
}
