package endpoints

import (
	"net/http"

	"github.com/nildev/lib/utils"
)

type (
	// SetupASGRequest type
	AddNodeRequest struct {
		Node
		ASGID string
	}

	AddNodeResponse struct{}
)

func AddNodeHandler(rw http.ResponseWriter, r *http.Request) {

	returnCode := http.StatusOK
	utils.Respond(rw, "", returnCode)
}
