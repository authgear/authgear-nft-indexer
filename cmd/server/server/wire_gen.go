// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"github.com/authgear/authgear-nft-indexer/pkg/handler"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	"github.com/authgear/authgear-nft-indexer/pkg/web3"
	"github.com/authgear/authgear-server/pkg/util/httputil"
	"net/http"
)

// Injectors from wire.go:

func NewHealthCheckAPIHandler(p *handler.RequestProvider) http.Handler {
	factory := p.LogFactory
	jsonResponseWriterLogger := httputil.NewJSONResponseWriterLogger(factory)
	jsonResponseWriter := &httputil.JSONResponseWriter{
		Logger: jsonResponseWriterLogger,
	}
	healthCheckHandlerLogger := handler.NewHealthCheckHandlerLogger(factory)
	config := p.Config
	db := p.Database
	request := p.Request
	context := handler.ProvideRequestContext(request)
	healthCheckAPIHandler := &handler.HealthCheckAPIHandler{
		JSON:     jsonResponseWriter,
		Logger:   healthCheckHandlerLogger,
		Config:   config,
		Database: db,
		Context:  context,
	}
	return healthCheckAPIHandler
}

func NewListCollectionAPIHandler(p *handler.RequestProvider) http.Handler {
	factory := p.LogFactory
	jsonResponseWriterLogger := httputil.NewJSONResponseWriterLogger(factory)
	jsonResponseWriter := &httputil.JSONResponseWriter{
		Logger: jsonResponseWriterLogger,
	}
	listCollectionHandlerLogger := handler.NewListCollectionHandlerLogger(factory)
	request := p.Request
	context := handler.ProvideRequestContext(request)
	db := p.Database
	nftCollectionQuery := query.NFTCollectionQuery{
		Ctx:     context,
		Session: db,
	}
	listCollectionAPIHandler := &handler.ListCollectionAPIHandler{
		JSON:               jsonResponseWriter,
		Logger:             listCollectionHandlerLogger,
		NFTCollectionQuery: nftCollectionQuery,
	}
	return listCollectionAPIHandler
}

func NewListOwnerNFTAPIHandler(p *handler.RequestProvider) http.Handler {
	factory := p.LogFactory
	jsonResponseWriterLogger := httputil.NewJSONResponseWriterLogger(factory)
	jsonResponseWriter := &httputil.JSONResponseWriter{
		Logger: jsonResponseWriterLogger,
	}
	listOwnerNFTHandlerLogger := handler.NewListOwnerNFTHandlerLogger(factory)
	config := p.Config
	request := p.Request
	context := handler.ProvideRequestContext(request)
	db := p.Database
	nftOwnerQuery := query.NFTOwnerQuery{
		Ctx:     context,
		Session: db,
	}
	nftCollectionQuery := query.NFTCollectionQuery{
		Ctx:     context,
		Session: db,
	}
	listOwnerNFTAPIHandler := &handler.ListOwnerNFTAPIHandler{
		JSON:               jsonResponseWriter,
		Logger:             listOwnerNFTHandlerLogger,
		Config:             config,
		NFTOwnerQuery:      nftOwnerQuery,
		NFTCollectionQuery: nftCollectionQuery,
	}
	return listOwnerNFTAPIHandler
}

func NewGetCollectionAPIHandler(p *handler.RequestProvider) http.Handler {
	factory := p.LogFactory
	jsonResponseWriterLogger := httputil.NewJSONResponseWriterLogger(factory)
	jsonResponseWriter := &httputil.JSONResponseWriter{
		Logger: jsonResponseWriterLogger,
	}
	listCollectionHandlerLogger := handler.NewListCollectionHandlerLogger(factory)
	request := p.Request
	context := handler.ProvideRequestContext(request)
	db := p.Database
	nftCollectionQuery := query.NFTCollectionQuery{
		Ctx:     context,
		Session: db,
	}
	getCollectionAPIHandler := &handler.GetCollectionAPIHandler{
		JSON:               jsonResponseWriter,
		Logger:             listCollectionHandlerLogger,
		NFTCollectionQuery: nftCollectionQuery,
	}
	return getCollectionAPIHandler
}

func NewGetCollectionMetadataAPIHandler(p *handler.RequestProvider) http.Handler {
	factory := p.LogFactory
	jsonResponseWriterLogger := httputil.NewJSONResponseWriterLogger(factory)
	jsonResponseWriter := &httputil.JSONResponseWriter{
		Logger: jsonResponseWriterLogger,
	}
	getCollectionMetadataHandlerLogger := handler.NewGetCollectionMetadataHandlerLogger(factory)
	config := p.Config
	alchemyAPI := &web3.AlchemyAPI{
		Config: config,
	}
	limiter := p.RateLimiter
	getCollectionMetadataAPIHandler := &handler.GetCollectionMetadataAPIHandler{
		JSON:        jsonResponseWriter,
		Logger:      getCollectionMetadataHandlerLogger,
		AlchemyAPI:  alchemyAPI,
		RateLimiter: limiter,
	}
	return getCollectionMetadataAPIHandler
}

func NewProbeCollectionAPIHandler(p *handler.RequestProvider) http.Handler {
	factory := p.LogFactory
	jsonResponseWriterLogger := httputil.NewJSONResponseWriterLogger(factory)
	jsonResponseWriter := &httputil.JSONResponseWriter{
		Logger: jsonResponseWriterLogger,
	}
	probeCollectionHandlerLogger := handler.NewProbeCollectionHandlerLogger(factory)
	config := p.Config
	alchemyAPI := &web3.AlchemyAPI{
		Config: config,
	}
	limiter := p.RateLimiter
	probeCollectionAPIHandler := &handler.ProbeCollectionAPIHandler{
		JSON:        jsonResponseWriter,
		Logger:      probeCollectionHandlerLogger,
		AlchemyAPI:  alchemyAPI,
		RateLimiter: limiter,
	}
	return probeCollectionAPIHandler
}
