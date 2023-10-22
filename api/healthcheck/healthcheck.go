package healthcheck

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/jackmcguire1/UserService/pkg/utils"
)

type HealthCheckHandler struct {
	LogVerbosity string
	StartTime    time.Time
	Logger       *slog.Logger
}

type HealthCheckResp struct {
	LogVerbosity string `json:"logVerbosity"`
	UpTime       string `json:"upTime"`
}

func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "OPTIONS,GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Requested-With,Origin,Accept")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	h.Logger.Info("fetching healthcheck")

	data := &HealthCheckResp{
		LogVerbosity: h.LogVerbosity,
		UpTime:       time.Since(h.StartTime).String(),
	}
	w.Write(utils.ToRAWJSON(data))
	w.WriteHeader(http.StatusOK)

	h.Logger.
		With("users", utils.ToJSON(data)).
		Debug("returning healthcheck")

	return
}
