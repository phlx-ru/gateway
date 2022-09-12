package main

import (
	"context"
	"flag"
	"os"
	"path"

	"gateway/internal/clients"
	"gateway/internal/conf"
	pkgConfig "gateway/internal/pkg/config"
	"gateway/internal/pkg/logger"
	"gateway/internal/pkg/metrics"
	"gateway/internal/pkg/runtime"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/joho/godotenv"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = `gateway_server`
	// Version is the version of the compiled software.
	Version = `1.1.1`
	// flagconf is the config flag.
	flagconf string
	// dotenv is loaded from config path .env file
	dotenv string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&dotenv, "dotenv", ".env", ".env file, eg: -dotenv .env")
}

func newApp(ctx context.Context, logger log.Logger, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Context(ctx),
		kratos.Server(
			hs,
		),
	)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flag.Parse()

	var err error

	ctx := context.Background()

	envPath := path.Join(flagconf, dotenv)
	err = godotenv.Overload(envPath)
	if err != nil {
		return err
	}

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(pkgConfig.EnvReplaceDecoder),
	)
	defer func() {
		_ = c.Close()
	}()

	if err = c.Load(); err != nil {
		return err
	}

	var bc conf.Bootstrap
	if err = c.Scan(&bc); err != nil {
		return err
	}

	logs := logger.New(id, Name, Version, bc.Log.Level)
	logHelper := logger.NewHelper(logs, "scope", "server")

	metric, err := metrics.New(bc.Metrics.Address, Name, bc.Metrics.Mute)
	if err != nil {
		return err
	}
	defer metric.Close()
	metric.Increment("starts.count")

	go runtime.CollectGoMetrics(ctx, metric)

	a := bc.Client.Grpc.Auth
	authClient, err := clients.NewAuth(ctx, a.Endpoint, a.Timeout.AsDuration(), metric, logs)
	if err != nil {
		return err
	}

	app, err := wireApp(ctx, bc.Server, authClient, metric, logs)
	if err != nil {
		return err
	}

	// start and wait for stop signal
	if err = app.Run(); err != nil {
		return err
	}

	logHelper.Info("app terminates")

	return nil
}
