package endpoints

import (
	"net/http"

	"github.com/nildev/lib/utils"
)

func AddNodeHandler(rw http.ResponseWriter, r *http.Request) {

	returnCode := http.StatusOK
	utils.Respond(rw, "", returnCode)
}
