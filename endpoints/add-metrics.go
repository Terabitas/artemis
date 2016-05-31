package endpoints

import (
	"net/http"

	"encoding/json"
	"io/ioutil"
	"time"

	"bitbucket.org/nildev/lib/Godeps/_workspace/src/github.com/juju/errors"
	log "github.com/Sirupsen/logrus"
	"github.com/nildev/artemis/domain"
	"github.com/nildev/artemis/version"
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
	ctxLog := log.WithField("version", version.Version).WithField("git-hash", version.GitHash).WithField("build-time", version.BuiltTimestamp)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ctxLog.Error(err)
		utils.Respond(rw, err.Error(), http.StatusBadRequest)
		return
	}

	req := &AddMetricsRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		ctxLog.Error(err)
		utils.Respond(rw, err.Error(), http.StatusBadRequest)
		return
	}

	asg := ASGSupervisor.Get(domain.ID(req.ID))
	if asg == nil {
		err := errors.Errorf("ASG with ID [%s], could not be found! Have you created it with /setup endpoint?", domain.ID(req.ID))
		ctxLog.Error(err)
		utils.Respond(rw, err.Error(), http.StatusBadRequest)
		return
	}

	healthMetricSeries := domain.NewMetricSeries()

	for _, m := range req.Metrics {
		ctxLog.Infof("Adding: [%v] [%v] \n", m.Value, m.Time)
		healthMetricSeries[m.Time] = domain.NewHealthMetric(m.Value, m.Time)
	}

	err = asg.AddMetrics(domain.ID(req.NodeID), healthMetricSeries)
	if err := json.Unmarshal(body, req); err != nil {
		ctxLog.Error(err)
		utils.Respond(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	returnCode := http.StatusOK
	utils.Respond(rw, nil, returnCode)
}
