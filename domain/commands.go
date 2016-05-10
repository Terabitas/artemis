package domain

import (
	"fmt"
	"net"
	"time"

	"strconv"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

type (
	TokenSource struct {
		AccessToken string
	}

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

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func (lc *Launch) Execute(asg *AutoScalingGroup) error {

	dropletName := "auto-" + strconv.Itoa(time.Now().Nanosecond())
	tokenSource := &TokenSource{
		AccessToken: lc.Provider.APIKey,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	createRequest := &godo.DropletCreateRequest{
		Name:   dropletName,
		Region: lc.Provider.Region,
		Size:   lc.Provider.Size,
		Image: godo.DropletCreateImage{
			Slug: lc.Provider.Image,
		},
		PrivateNetworking: true,
		SSHKeys: []godo.DropletCreateSSHKey{
			godo.DropletCreateSSHKey{
				Fingerprint: lc.Provider.SSHKey,
			},
		},
	}

	newDroplet, _, err := client.Droplets.Create(createRequest)

	if err != nil {
		fmt.Printf("Could not launch droplet : %s\n\n", err)
		return err
	}

	status := newDroplet.Status
	for {
		fmt.Printf("Droplet [%d] status [%s] : \n\n", newDroplet.ID, status)
		if status == "active" {
			break
		}

		d, _, err := client.Droplets.Get(newDroplet.ID)
		if err != nil {
			fmt.Printf("Could not get status for droplet : %s\n\n", err)
			return err
		}

		status = d.Status
		// timeout needed
		time.Sleep(time.Second * 5)
	}

	publicIP, err := newDroplet.PublicIPv4()
	if err != nil {
		fmt.Printf("Could not get public IP : %s\n\n", err)
		return err
	}

	privateIP, err := newDroplet.PrivateIPv4()
	if err != nil {
		fmt.Printf("Could not get private IP : %s\n\n", err)
		return err
	}

	fmt.Printf("Setting up new Node for [%s] \n", asg.ID)
	// API call
	node := NewNode()
	node.Setup(
		ID(strconv.Itoa(newDroplet.ID)),
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

	// Only when health metrics are received then return

	time.Sleep(time.Second * 3)

	return nil
}

func (lc *Terminate) Execute(asg *AutoScalingGroup) error {

	tokenSource := &TokenSource{
		AccessToken: lc.Provider.APIKey,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	nid, err := strconv.Atoi(string(lc.NodeID))
	if err != nil {
		fmt.Printf("Could not convert to int [%s]: %s\n\n", lc.NodeID, err)
		return err
	}

	// Remove bad one
	_, err = client.Droplets.Delete(nid)
	asg.RemoveNode(lc.NodeID)

	if err != nil {
		fmt.Printf("Could not delete node [%s]: %s\n\n", lc.NodeID, err)
		return err
	}

	fmt.Printf("Execute Terminate [%s] \n", asg.ID)

	return nil
}

func (lc *Relaunch) Execute(asg *AutoScalingGroup) error {

	dropletName := "auto-" + strconv.Itoa(time.Now().Nanosecond())
	tokenSource := &TokenSource{
		AccessToken: lc.Provider.APIKey,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	// Launch new

	createRequest := &godo.DropletCreateRequest{
		Name:   dropletName,
		Region: lc.Provider.Region,
		Size:   lc.Provider.Size,
		Image: godo.DropletCreateImage{
			Slug: lc.Provider.Image,
		},
		PrivateNetworking: true,
		SSHKeys: []godo.DropletCreateSSHKey{
			godo.DropletCreateSSHKey{
				Fingerprint: lc.Provider.SSHKey,
			},
		},
	}

	newDroplet, _, err := client.Droplets.Create(createRequest)

	if err != nil {
		fmt.Printf("Could not launch droplet : %s\n\n", err)
		return err
	}

	status := newDroplet.Status
	for {
		fmt.Printf("Droplet [%d] status [%s] : \n\n", newDroplet.ID, status)
		if status == "active" {
			break
		}

		d, _, err := client.Droplets.Get(newDroplet.ID)
		if err != nil {
			fmt.Printf("Could not get status for droplet : %s\n\n", err)
			return err
		}

		status = d.Status
		// timeout needed
		time.Sleep(time.Second * 5)
	}

	publicIP, err := newDroplet.PublicIPv4()
	if err != nil {
		fmt.Printf("Could not get public IP : %s\n\n", err)
		return err
	}

	privateIP, err := newDroplet.PrivateIPv4()
	if err != nil {
		fmt.Printf("Could not get private IP : %s\n\n", err)
		return err
	}

	fmt.Printf("Setting up new Node for [%s] \n", asg.ID)
	// API call
	node := NewNode()
	node.Setup(
		ID(strconv.Itoa(newDroplet.ID)),
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

	// Only when health metrics are received then return

	time.Sleep(time.Second * 3)

	nid, err := strconv.Atoi(string(lc.NodeID))
	if err != nil {
		fmt.Printf("Could not convert to int [%s]: %s\n\n", lc.NodeID, err)
		return err
	}

	// Remove bad one
	_, err = client.Droplets.Delete(nid)
	asg.RemoveNode(lc.NodeID)
	if err != nil {
		fmt.Printf("Could not delete node [%s]: %s\n\n", lc.NodeID, err)
		return err
	}

	fmt.Printf("Execute Relaunch [%s] \n", asg.ID)

	return nil
}
