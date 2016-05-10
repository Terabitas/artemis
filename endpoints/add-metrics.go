package endpoints

import (
	"net/http"

	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/nildev/artemis/domain"
	"github.com/nildev/lib/utils"
)

type (
	Metric struct {
		Value float64
		Time  time.Time
	}

	AddMetricsRequest struct {
		ID      string
		NodeID  string
		Metrics []Metric
	}
)

func AddMetricsHandler(rw http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Respond(rw, err.Error(), http.StatusBadRequest)
		return
	}

	req := &AddMetricsRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		utils.Respond(rw, err.Error(), http.StatusBadRequest)
		return
	}

	asg := ASGSupervisor.Get(domain.ID(req.ID))
	if asg == nil {
		utils.Respond(rw, err.Error(), http.StatusBadRequest)
		return
	}

	healthMetricSeries := domain.NewMetricSeries()

	for _, m := range req.Metrics {
		//fmt.Printf("Adding: [%v] [%v] \n", m.Value, m.Time)
		healthMetricSeries[m.Time] = domain.NewHealthMetric(m.Value, m.Time)
	}

	err = asg.AddMetrics(domain.ID(req.NodeID), healthMetricSeries)
	if err := json.Unmarshal(body, req); err != nil {
		utils.Respond(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	returnCode := http.StatusOK
	utils.Respond(rw, nil, returnCode)
}
