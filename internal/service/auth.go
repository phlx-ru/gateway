package service

import (
	"github.com/gin-gonic/gin"

	"gateway/api/gateway"
	pkgStrings "gateway/internal/pkg/strings"
)

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
