package healthcheck

import (
	"net/http"
	"time"

	"github.com/jackmcguire1/UserService/pkg/utils"
)

type HealthCheckHandler struct {
	LogVerbosity string
	StartTime    time.Time
}

type HealthCheckResp struct {
	LogVerbosity string `json:"logVerbosity"`
	UpTime       string `json:"upTime"`
}

func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	w.Write([]byte(utils.ToJSON(&HealthCheckResp{
		LogVerbosity: h.LogVerbosity,
		UpTime:       time.Since(h.StartTime).String(),
	})))
	w.WriteHeader(http.StatusOK)

	return
}
