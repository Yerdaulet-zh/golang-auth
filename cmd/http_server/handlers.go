package httpserver

import (
	"net/http"

	httpAdapter "github.com/golang-auth/internal/adapters/handlers/http"
	"github.com/golang-auth/internal/adapters/repository/postgre"
)

// HealthHandler defines the dependencies for health checks
type HealthHandler struct {
	dbClient *postgre.Client
}

func NewHealthHandler(db *postgre.Client) *HealthHandler {
	return &HealthHandler{dbClient: db}
}

func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	httpAdapter.OpsProcessed.Inc()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	if err := h.dbClient.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Service Unavailable: Database unreachable"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}
