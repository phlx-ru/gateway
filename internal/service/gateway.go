package service

import (
	"gateway/api/gateway"
	"gateway/internal/biz"
	"gateway/internal/pkg/logger"
	"gateway/internal/pkg/metrics"
	pkgStrings "gateway/internal/pkg/strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
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

func (g *GatewayService) GetAuthCheck(c *gin.Context, params gateway.GetAuthCheckParams) {
	method := `getAuthCheck`
	defer g.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { g.postProcess(c.Request.Context(), method, err) }()

	res, err := g.usecase.GetAuthCheck(c.Request.Context(), params)
	if err != nil {
		g.responseError(c, err)
		return
	}

	g.responseOK(c, res)
}

func (g *GatewayService) PostAuthLogin(c *gin.Context) {
	method := `postAuthLogin`
	defer g.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { g.postProcess(c.Request.Context(), method, err) }()

	var in *gateway.PostAuthLoginJSONRequestBody
	if err = c.Bind(&in); err != nil {
		g.responseValidationError(c, err)
		return
	}

	if err = g.validate(in); err != nil {
		g.responseError(c, err)
		return
	}

	out, err := g.usecase.PostAuthLogin(c.Request.Context(), in)
	if err != nil {
		g.responseError(c, err)
		return
	}

	g.responseOK(c, out)
}

func (g *GatewayService) PostGenerateCode(c *gin.Context) {
	method := `postGenerateCode`
	defer g.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { g.postProcess(c.Request.Context(), method, err) }()

	var in *gateway.PostGenerateCodeJSONRequestBody
	if err = c.Bind(&in); err != nil {
		g.responseError(c, err)
		return
	}

	if err = g.validate(in); err != nil {
		g.responseError(c, err)
		return
	}

	err = g.usecase.PostGenerateCode(c.Request.Context(), in)
	if err != nil {
		g.responseError(c, err)
		return
	}

	g.responseNoContent(c)
}

func (g *GatewayService) PostAuthLoginByCode(c *gin.Context) {
	method := `postAuthLoginByCode`
	defer g.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { g.postProcess(c.Request.Context(), method, err) }()

	var in *gateway.PostAuthLoginByCodeJSONRequestBody
	if err = c.Bind(&in); err != nil {
		g.responseError(c, err)
		return
	}

	if err = g.validate(in); err != nil {
		g.responseError(c, err)
		return
	}

	out, err := g.usecase.PostAuthLoginByCode(c.Request.Context(), in)
	if err != nil {
		g.responseError(c, err)
		return
	}

	g.responseOK(c, out)
}

func (g *GatewayService) PostResetPassword(c *gin.Context) {
	method := `postResetPassword`
	defer g.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { g.postProcess(c.Request.Context(), method, err) }()

	var in *gateway.PostResetPasswordJSONRequestBody
	if err = c.Bind(&in); err != nil {
		g.responseError(c, err)
		return
	}

	if err = g.validate(in); err != nil {
		g.responseError(c, err)
		return
	}

	err = g.usecase.PostResetPassword(c.Request.Context(), in)
	if err != nil {
		g.responseError(c, err)
		return
	}

	g.responseNoContent(c)
}

func (g *GatewayService) PostChangePassword(c *gin.Context) {
	method := `postChangePassword`
	defer g.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { g.postProcess(c.Request.Context(), method, err) }()

	var in *gateway.PostChangePasswordJSONRequestBody
	if err = c.Bind(&in); err != nil {
		g.responseError(c, err)
		return
	}

	if err = g.validate(in); err != nil {
		g.responseError(c, err)
		return
	}

	err = g.usecase.PostChangePassword(c.Request.Context(), in)
	if err != nil {
		g.responseError(c, err)
		return
	}

	g.responseNoContent(c)
}

func (g *GatewayService) PostNewPassword(c *gin.Context) {
	method := `postNewPassword`
	defer g.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { g.postProcess(c.Request.Context(), method, err) }()

	var in *gateway.PostNewPasswordJSONRequestBody
	if err = c.Bind(&in); err != nil {
		g.responseError(c, err)
		return
	}

	if err = g.validate(in); err != nil {
		g.responseError(c, err)
		return
	}

	err = g.usecase.PostNewPassword(c.Request.Context(), in)
	if err != nil {
		g.responseError(c, err)
		return
	}

	g.responseNoContent(c)
}
