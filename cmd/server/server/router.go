package server

import (
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/handler"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
	"github.com/uptrace/bun"
)

func NewRouter(config config.Config, session *bun.DB, lf *log.Factory) *httproute.Router {
	router := httproute.NewRouter()

	routeHandler := handler.RouteHandler{
		Config:     config,
		Database:   session,
		LogFactory: lf,
	}
	route := httproute.Route{}
	router.Add(handler.ConfigureHealthCheckRoute(route), routeHandler.Handle(NewHealthCheckAPIHandler))
	router.Add(handler.ConfigureWatchCollectionRoute(route), routeHandler.Handle(NewWatchCollectionAPIHandler))
	router.Add(handler.ConfigureListCollectionRoute(route), routeHandler.Handle(NewListCollectionAPIHandler))
	router.Add(handler.ConfigureListOwnerNFTRoute(route), routeHandler.Handle(NewListOwnerNFTAPIHandler))

	return router
}
