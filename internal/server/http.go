package server

import (
	"fmt"
	"time"

	"gateway/api/gateway"
	"gateway/internal/conf"
	"gateway/internal/middlewares"
	"gateway/internal/pkg/metrics"
	"gateway/internal/service"

	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	kratosHTTP "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new HTTP server.
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
	router := gin.New()
	router.Use(
		gin.LoggerWithConfig(
			gin.LoggerConfig{
				Formatter: LogFormatter,
				Output:    gin.DefaultWriter,
				SkipPaths: []string{
					`/swagger/`,
					`/swagger/swagger-ui.css`,
					`/swagger/swagger-ui.css.map`,
					`/swagger/swagger-ui-bundle.js`,
					`/swagger/swagger-ui-bundle.js.map`,
					`/swagger/swagger.yaml`,
					`/swagger/favicon.ico`,
					`/auth/`,
					`/auth/favicon.ico`,
					`/auth/lib/spectre.min.css`,
					`/auth/lib/spectre-exp.min.css`,
					`/auth/lib/spectre-icons.min.css`,
					`/auth/lib/auth.css`,
					`/auth/lib/auth.js`,
					`/auth/lib/vue.global.prod.js`,
				},
			},
		),
		gin.Recovery(),
		kgin.Middlewares(
			middlewares.Duration(metric, logger),
			tracing.Server(),
		),
	)
	router.Static(`/auth`, `./static/auth`)
	router.Static(`/swagger`, `./static/swagger`)

	router.GET(`/api/swagger`, gw.GetSwagger)

	gateway.RegisterHandlersWithOptions(
		router,
		gw,
		gateway.GinServerOptions{
			BaseURL:     ``, // c.BaseUrl,
			Middlewares: []gateway.MiddlewareFunc{
				// TODO CHECK MIDDLEWARES FOR KRATOS BELOW
			},
		},
	)

	srv.HandlePrefix(`/`, router)

	return srv
}

func LogFormatter(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf(
		"ACCESS ts=%v status=%d latency=%v client.ip=%s method=%s path=%-7s error=%#v\n",
		param.TimeStamp.Format(time.RFC3339),
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}
