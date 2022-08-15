package handler

import (
	"net/http"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/uptrace/bun"
)

type JSONResponseWriter interface {
	WriteResponse(rw http.ResponseWriter, resp *authgearapi.Response)
}

type RequestProvider struct {
	Config         config.Config
	Database       *bun.DB
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

func RouteHandler(factory func(*RequestProvider) http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := &RequestProvider{
			Request:        r,
			ResponseWriter: w,
		}

		router := factory(p)
		router.ServeHTTP(w, r)
	})
}
