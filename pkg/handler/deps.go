package handler

import (
	"context"
	"net/http"

	"github.com/google/wire"
)

func ProvideRequestContext(r *http.Request) context.Context { return r.Context() }

var DependencySet = wire.NewSet(
	ProvideRequestContext,

	wire.FieldsOf(new(*RequestProvider),
		"Config",
		"LogFactory",
		"Database",
		"RateLimiter",
		"Request",
	),
	wire.Struct(new(HealthCheckAPIHandler), "*"),
	NewHealthCheckHandlerLogger,
	wire.Struct(new(WatchCollectionAPIHandler), "*"),
	NewWatchCollectionHandlerLogger,
	wire.Struct(new(ListCollectionAPIHandler), "*"),
	NewListCollectionHandlerLogger,
	wire.Struct(new(ListOwnerNFTAPIHandler), "*"),
	NewListOwnerNFTHandlerLogger,
	wire.Struct(new(GetCollectionAPIHandler), "*"),
	NewGetCollectionHandlerLogger,
)
