package server

import (
	"gateway/api/gateway"
	"gateway/internal/conf"
	"gateway/internal/middlewares"
	"gateway/internal/pkg/metrics"
	"gateway/internal/service"

	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	kratosHTTP "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(
	c *conf.Server,
	gw *service.GatewayService,
	metric metrics.Metrics,
	logger log.Logger,
) *kratosHTTP.Server {
	var opts = []kratosHTTP.ServerOption{
		kratosHTTP.Timeout(c.Http.Timeout.AsDuration()),
	}
	if c.Http.Network != "" {
		opts = append(opts, kratosHTTP.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, kratosHTTP.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, kratosHTTP.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := kratosHTTP.NewServer(opts...)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(
		kgin.Middlewares(
			middlewares.Duration(metric, logger),
			tracing.Server(),
		),
	)
	router.Static(`/auth`, `./static/auth`)

	router.GET(`/api/swagger`, gw.GetSwagger)

	router.GET(
		// TODO REMOVE
		"/helloworld/:name", func(c *gin.Context) {
			name := c.Param("name")
			if name == "error" {
				// 返回kratos error
				kgin.Error(c, errors.Unauthorized("auth_error", "no authentication"))
			} else {
				c.JSON(200, map[string]string{"welcome": name})
			}
		},
	)

	gateway.RegisterHandlersWithOptions(
		router,
		gw,
		gateway.GinServerOptions{
			BaseURL:     ``, //c.BaseUrl,
			Middlewares: []gateway.MiddlewareFunc{
				// TODO CHECK MIDDLEWARES FOR KRATOS BELOW
			},
		},
	)

	srv.HandlePrefix(`/`, router)

	return srv
}
