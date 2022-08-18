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
		"Request",
	),
	wire.Struct(new(HealthCheckAPIHandler), "*"),
	NewHealthCheckHandlerLogger,
	wire.Struct(new(RegisterCollectionAPIHandler), "*"),
	NewRegisterCollectionHandlerLogger,
	wire.Struct(new(DeregisterCollectionAPIHandler), "*"),
	NewDeregisterCollectionHandlerLogger,
	wire.Struct(new(ListCollectionAPIHandler), "*"),
	NewListCollectionHandlerLogger,
	wire.Struct(new(ListCollectionOwnersAPIHandler), "*"),
	NewListCollectionOwnerHandlerLogger,
	wire.Struct(new(ListOwnerNFTAPIHandler), "*"),
	NewListOwnerNFTHandlerLogger,
)
