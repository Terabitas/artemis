package endpoints

import (
	"net/http"

	"encoding/json"
	"io/ioutil"

	"net"

	"time"

	"github.com/nildev/artemis/domain"
	"github.com/nildev/lib/utils"
)

type (

	// SetupASGRequest type
	SetupASGRequest struct {
		ID string

		Nodes        []Node
		HealthPolicy HealthPolicy
	}

	SetupASGResponse struct{}
)

// SetupHandler API handler
func SetupHandler(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Respond(rw, err.Error(), http.StatusBadRequest)
		return
	}

	req := &SetupASGRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		utils.Respond(rw, err.Error(), http.StatusBadRequest)
		return
	}

	asg := domain.NewAutoScalingGroup(domain.ID(req.ID))

	plc, err := domain.NewDesiredNodeAmountPerProviderPolicy(
		domain.ID(req.HealthPolicy.ID),
		req.HealthPolicy.Min,
		req.HealthPolicy.Max,
		req.HealthPolicy.Desired,
		req.HealthPolicy.ConsecutiveChecks,
		req.HealthPolicy.HealthyThreshold,
		time.Duration(req.HealthPolicy.CheckInterval)*time.Second,
		domain.Provider{
			ID:     req.HealthPolicy.Provider.ID,
			APIKey: req.HealthPolicy.Provider.APIKey,
			Region: req.HealthPolicy.Provider.Region,
			Size:   req.HealthPolicy.Provider.Size,
			Image:  req.HealthPolicy.Provider.Image,
			SSHKey: req.HealthPolicy.Provider.SSHKey,
		},
	)
	policySet := domain.NewPolicySet(plc)

	nodeSet := domain.NewNodeSet()
	for _, n := range req.Nodes {
		node := domain.NewNode()
		node.Setup(
			domain.ID(n.ID),
			domain.Provider{
				ID:     n.Provider.ID,
				APIKey: n.Provider.APIKey,
			},
			domain.NetworkInterface{
				ID: domain.ID(n.PublicIFace.ID),
				IP: net.ParseIP(n.PublicIFace.IP),
			},
			domain.NetworkInterface{
				ID: domain.ID(n.PrivateIFace.ID),
				IP: net.ParseIP(n.PrivateIFace.IP),
			},
		)

		nodeSet[node.ID] = node
	}

	asg.Setup(nodeSet, policySet)

	// Start ASG routine
	ASGSupervisor.Add(asg)

	outResp := &SetupASGResponse{}
	out, err := json.Marshal(outResp)
	if err != nil {
		utils.Respond(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.Respond(rw, string(out), http.StatusCreated)
}
