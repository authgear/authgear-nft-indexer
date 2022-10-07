package server

import (
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/handler"
	"github.com/authgear/authgear-server/pkg/lib/infra/redis/appredis"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
	"github.com/uptrace/bun"
)

func NewRouter(config config.Config, session *bun.DB, redis *appredis.Handle, lf *log.Factory) *httproute.Router {
	router := httproute.NewRouter()

	routeHandler := handler.RouteHandler{
		Config:     config,
		Database:   session,
		Redis:      redis,
		LogFactory: lf,
	}
	route := httproute.Route{}
	router.Add(handler.ConfigureHealthCheckRoute(route), routeHandler.Handle(NewHealthCheckAPIHandler))
	router.Add(handler.ConfigureListCollectionRoute(route), routeHandler.Handle(NewListCollectionAPIHandler))
	router.Add(handler.ConfigureListOwnerNFTRoute(route), routeHandler.Handle(NewListOwnerNFTAPIHandler))
	router.Add(handler.ConfigureGetCollectionRoute(route), routeHandler.Handle(NewGetCollectionAPIHandler))
	router.Add(handler.ConfigureGetCollectionMetadataRoute(route), routeHandler.Handle(NewGetCollectionMetadataAPIHandler))
	router.Add(handler.ConfigureProbeCollectionRoute(route), routeHandler.Handle(NewProbeCollectionAPIHandler))
	return router
}
