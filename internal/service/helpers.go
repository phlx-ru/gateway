package service

import (
	"context"
	"net/http"

	pkgStrings "gateway/internal/pkg/strings"
	"gateway/internal/pkg/validate"

	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/errors"
)

func (g *GatewayService) postProcess(ctx context.Context, method string, err error) {
	if err != nil {
		g.logger.WithContext(ctx).Errorf(`method %s failed: %v`, method, err)
		g.metric.Increment(pkgStrings.Metric(method, `failure`))
	} else {
		g.metric.Increment(pkgStrings.Metric(method, `success`))
	}
}

func (g *GatewayService) validate(s any) error {
	validationErrors, err := validate.Struct(s)
	if err != nil {
		return errors.InternalServer(`validator_fails`, err.Error())
	}
	if validationErrors != nil {
		metadata := validate.AsCustomValidationTranslations(validationErrors)
		return errors.BadRequest(`validation_error`, `incorrect input`).WithMetadata(metadata)
	}
	return nil
}

func (g *GatewayService) responseError(c *gin.Context, err error) {
	c.Header(`Content-Type`, `application/json`) // Bugfix for error render
	kgin.Error(c, err)
}

func (g *GatewayService) responseValidationError(c *gin.Context, err error) {
	validationErr := errors.BadRequest(`validation_error`, err.Error())
	g.responseError(c, validationErr)
}

func (g *GatewayService) responseOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func (g *GatewayService) responseNoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, ``)
}
