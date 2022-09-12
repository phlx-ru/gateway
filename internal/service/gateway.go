package service

import (
	"net/http"

	"gateway/api/gateway"
	"gateway/internal/clients"
	"gateway/internal/pkg/logger"
	"gateway/internal/pkg/metrics"

	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

const (
	metricGetAuthCheckTimings = `service.gateway.getAuthCheck.timings`
	metricGetAuthCheckSuccess = `service.gateway.getAuthCheck.success`
	metricGetAuthCheckFailure = `service.gateway.getAuthCheck.failure`
)

type GatewayService struct {
	authClient clients.AuthClient
	metric     metrics.Metrics
	logger     *log.Helper
}

func NewGatewayService(
	authClient clients.AuthClient,
	metric metrics.Metrics,
	logs log.Logger,
) *GatewayService {
	return &GatewayService{
		authClient: authClient,
		metric:     metric,
		logger:     logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "service/gateway"),
	}
}

func (s *GatewayService) GetAuthCheck(c *gin.Context, params gateway.GetAuthCheckParams) {
	defer s.metric.NewTiming().Send(metricGetAuthCheckTimings)
	var err error
	defer func() {
		if err != nil {
			s.logger.WithContext(c.Request.Context()).Errorf(`failed GetAuthCheck: %v`, err)
			s.metric.Increment(metricGetAuthCheckFailure)
		} else {
			s.metric.Increment(metricGetAuthCheckSuccess)
		}
	}()

	token := ""
	if params.AuthToken != nil {
		token = *params.AuthToken
	}
	if params.XAuthToken != nil {
		token = *params.XAuthToken
	}
	if token == "" {
		err = errors.BadRequest(`token_not_found`, `not found authToken in query or X-Auth-Token in headers`)
		kgin.Error(c, err)
		return
	}

	res, err := s.getAuthCheck(c.Request.Context(), token)
	if err != nil {
		kgin.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
