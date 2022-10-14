package handler

import (
	"net/http"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/lib/infra/redis/appredis"
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
		p := &RequestProvider{
			Config:         rh.Config,
			Database:       rh.Database,
			LogFactory:     rh.LogFactory,
			Request:        r,
			ResponseWriter: w,
		}

		router := factory(p)
		router.ServeHTTP(w, r)
	})
}
