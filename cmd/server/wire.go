//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"gateway/internal/biz"
	"gateway/internal/clients"
	"gateway/internal/conf"
	"gateway/internal/pkg/metrics"
	"gateway/internal/server"
	"gateway/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(
	context.Context,
	*conf.Server,
	clients.AuthClient,
	metrics.Metrics,
	log.Logger,
) (
	*kratos.App,
	error,
) {
	panic(wire.Build(server.ProviderSet, service.ProviderSet, biz.ProviderSet, newApp))
}
