package endpoints

import (
	"net/http"

	"encoding/json"
	"io/ioutil"

	"github.com/nildev/artemis/domain"
	"github.com/nildev/lib/utils"
)

type (
	// RemoveASGRequest type
	RemoveASGRequest struct {
		ID string
	}

	RemoveASGResponse struct{}
)

// RemoveASGHandler API handler
func RemoveASGHandler(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Respond(rw, err.Error(), http.StatusBadRequest)
		return
	}

	req := &RemoveASGRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		utils.Respond(rw, err.Error(), http.StatusBadRequest)
		return
	}

	asg := domain.NewAutoScalingGroup(domain.ID(req.ID))

	asg.Setup(domain.NewNodeSet(), domain.NewPolicySet())

	// Remove ASG routine
	ASGSupervisor.Remove(domain.ID(req.ID))

	outResp := &RemoveASGResponse{}
	out, err := json.Marshal(outResp)
	if err != nil {
		utils.Respond(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.Respond(rw, string(out), http.StatusCreated)
}
