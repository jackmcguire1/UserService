package healthcheck

import (
	"github.com/jackmcguire1/UserService/pkg/utils"
	"net/http"
)

type HealthCheckHandler struct {
	LogVerbosity string
}

type HealthCheckResp struct {
	Alive        bool
	LogVerbosity string
}

func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	w.Write([]byte(utils.ToJSON(&HealthCheckResp{
		Alive:        true,
		LogVerbosity: h.LogVerbosity,
	})))
	w.WriteHeader(http.StatusOK)

	return
}
