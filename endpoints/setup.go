package endpoints

import (
	"net/http"

	"encoding/json"
	"io/ioutil"

	"github.com/nildev/artemis/domain"
	"github.com/nildev/lib/utils"
)

type (
	// SetupASGRequest type
	SetupASGRequest struct {
		ID string
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

	asg.Setup(domain.NewNodeSet(), domain.NewPolicySet())

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
