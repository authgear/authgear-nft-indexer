package server

import (
	"net/http"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/handler"
	"github.com/authgear/authgear-server/pkg/util/httproute"
)

func NewRouter(config config.Config) *httproute.Router {
	router := httproute.NewRouter()

	// db := database.GetDatabase(config.Database)

	router.Add(httproute.Route{
		Methods:     []string{"GET"},
		PathPattern: "/ping",
	}, http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("pong"))
	}))

	route := httproute.Route{}
	router.Add(handler.ConfigureRegisterCollectionRoute(route), handler.RouteHandler(NewRegisterCollectionAPIHandler))
	router.Add(handler.ConfigureDeregisterCollectionRoute(route), handler.RouteHandler(NewDeregisterCollectionAPIHandler))
	router.Add(handler.ConfigureListCollectionRoute(route), handler.RouteHandler(NewListCollectionAPIHandler))
	router.Add(handler.ConfigureListCollectionOwnerRoute(route), handler.RouteHandler(NewListCollectionOwnerAPIHandler))
	router.Add(handler.ConfigureListOwnerNFTRoute(route), handler.RouteHandler(NewListOwnerNFTAPIHandler))

	return router
}
