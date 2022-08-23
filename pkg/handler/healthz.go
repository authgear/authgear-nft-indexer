package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
	"github.com/uptrace/bun"
)

func ConfigureHealthCheckRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("GET").
		WithPathPattern("/healthz")
}

type HealthCheckHandlerLogger struct{ *log.Logger }

func NewHealthCheckHandlerLogger(lf *log.Factory) HealthCheckHandlerLogger {
	return HealthCheckHandlerLogger{lf.New("healthz")}
}

type HealthCheckAPIHandler struct {
	JSON     JSONResponseWriter
	Logger   HealthCheckHandlerLogger
	Config   config.Config
	Database *bun.DB
	Context  context.Context
}

func (h *HealthCheckAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	err := h.CheckHealth()
	if err != nil {
		h.Logger.WithError(err).Errorf("health check failed")
		http.Error(resp, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	body := []byte("OK")
	resp.Header().Set("Content-Type", "text/plain")
	resp.Header().Set("Content-Length", strconv.Itoa(len(body)))
	resp.WriteHeader(http.StatusOK)
	_, _ = resp.Write(body)
}

func (h *HealthCheckAPIHandler) CheckHealth() error {
	_, err := h.Database.QueryContext(h.Context, "SELECT 42")
	if err != nil {
		return err
	}
	h.Logger.Debugf("global database connection healthz passed")
	return nil
}
