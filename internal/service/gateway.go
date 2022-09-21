package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"

	"gateway/api/gateway"
	"gateway/internal/biz"
	"gateway/internal/pkg/logger"
	"gateway/internal/pkg/metrics"
)

const (
	metricPrefix = `service.gateway`
)

type GatewayService struct {
	usecase *biz.GatewayUsecase
	metric  metrics.Metrics
	logger  *log.Helper
}

func NewGatewayService(
	usecase *biz.GatewayUsecase,
	metric metrics.Metrics,
	logs log.Logger,
) *GatewayService {
	return &GatewayService{
		usecase: usecase,
		metric:  metric,
		logger:  logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "service/gateway"),
	}
}

func (g *GatewayService) GetSwagger(c *gin.Context) {
	swagger, err := gateway.GetSwagger()
	if err != nil {
		g.responseError(c, errors.InternalServer(`internal_error`, err.Error()))
		return
	}
	swagger.InternalizeRefs(c.Request.Context(), nil) // TODO Bug with nested refs and properties includes
	g.responseOK(c, swagger)
}
