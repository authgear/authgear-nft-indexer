package handler

import (
	"context"
	"net/http"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/google/wire"
	"github.com/uptrace/bun"
)

func ProvideRequestContext(r *http.Request) context.Context { return r.Context() }
func ProvideConfig(r *RequestProvider) config.Config        { return r.Config }
func ProvideDatabase(r *RequestProvider) *bun.DB            { return r.Database }

var DependencySet = wire.NewSet(
	ProvideRequestContext,

	wire.FieldsOf(new(*RequestProvider),
		"Config",
		"Database",
		"Request",
	),
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
