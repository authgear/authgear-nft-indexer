package handler

import (
	"net/http"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/lib/infra/redis/appredis"
	agratelimit "github.com/authgear/authgear-server/pkg/lib/ratelimit"
	"github.com/authgear/authgear-server/pkg/util/log"
	"github.com/uptrace/bun"
)

type JSONResponseWriter interface {
	WriteResponse(rw http.ResponseWriter, resp *authgearapi.Response)
}

type RequestProvider struct {
	Config         config.Config
	Database       *bun.DB
	Request        *http.Request
	LogFactory     *log.Factory
	RateLimiter    *agratelimit.Limiter
	ResponseWriter http.ResponseWriter
}

type RouteHandler struct {
	Config     config.Config
	Database   *bun.DB
	Redis      *appredis.Handle
	LogFactory *log.Factory
}

func (rh *RouteHandler) Handle(factory func(*RequestProvider) http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rlf := NewRateLimiterFactory(rh.LogFactory, rh.Redis)

		query := r.URL.Query()
		appID := query.Get("app_id")

		rl := new(agratelimit.Limiter)
		if appID != "" {
			rl = rlf.New(appID)
		}

		p := &RequestProvider{
			Config:         rh.Config,
			Database:       rh.Database,
			LogFactory:     rh.LogFactory,
			RateLimiter:    rl,
			Request:        r,
			ResponseWriter: w,
		}

		router := factory(p)
		router.ServeHTTP(w, r)
	})
}
