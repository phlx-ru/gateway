// Package gateway provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package gateway

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	externalRef0 "gateway/api/gateway/auth"
	externalRef1 "gateway/api/gateway/common"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// PostChangePasswordJSONBody defines parameters for PostChangePassword.
type PostChangePasswordJSONBody = externalRef0.ChangePasswordRequestBody

// GetAuthCheckParams defines parameters for GetAuthCheck.
type GetAuthCheckParams struct {
	// Auth token from Query
	AuthToken *externalRef1.AuthTokenInQuery `form:"authToken,omitempty" json:"authToken,omitempty"`
}

// PostGenerateCodeJSONBody defines parameters for PostGenerateCode.
type PostGenerateCodeJSONBody = externalRef0.GenerateCodeRequestBody

// PostAuthLoginJSONBody defines parameters for PostAuthLogin.
type PostAuthLoginJSONBody = externalRef0.LoginRequestBody

// PostAuthLoginByCodeJSONBody defines parameters for PostAuthLoginByCode.
type PostAuthLoginByCodeJSONBody = externalRef0.LoginByCodeRequestBody

// PostNewPasswordJSONBody defines parameters for PostNewPassword.
type PostNewPasswordJSONBody = externalRef0.NewPasswordRequestBody

// PostResetPasswordJSONBody defines parameters for PostResetPassword.
type PostResetPasswordJSONBody = externalRef0.ResetPasswordRequestBody

// PostChangePasswordJSONRequestBody defines body for PostChangePassword for application/json ContentType.
type PostChangePasswordJSONRequestBody = PostChangePasswordJSONBody

// PostGenerateCodeJSONRequestBody defines body for PostGenerateCode for application/json ContentType.
type PostGenerateCodeJSONRequestBody = PostGenerateCodeJSONBody

// PostAuthLoginJSONRequestBody defines body for PostAuthLogin for application/json ContentType.
type PostAuthLoginJSONRequestBody = PostAuthLoginJSONBody

// PostAuthLoginByCodeJSONRequestBody defines body for PostAuthLoginByCode for application/json ContentType.
type PostAuthLoginByCodeJSONRequestBody = PostAuthLoginByCodeJSONBody

// PostNewPasswordJSONRequestBody defines body for PostNewPassword for application/json ContentType.
type PostNewPasswordJSONRequestBody = PostNewPasswordJSONBody

// PostResetPasswordJSONRequestBody defines body for PostResetPassword for application/json ContentType.
type PostResetPasswordJSONRequestBody = PostResetPasswordJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /api/1/auth/changePassword)
	PostChangePassword(c *gin.Context)

	// (GET /api/1/auth/check)
	GetAuthCheck(c *gin.Context, params GetAuthCheckParams)

	// (POST /api/1/auth/generateCode)
	PostGenerateCode(c *gin.Context)

	// (POST /api/1/auth/login)
	PostAuthLogin(c *gin.Context)

	// (POST /api/1/auth/loginByCode)
	PostAuthLoginByCode(c *gin.Context)

	// (POST /api/1/auth/newPassword)
	PostNewPassword(c *gin.Context)

	// (POST /api/1/auth/resetPassword)
	PostResetPassword(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
}

type MiddlewareFunc func(c *gin.Context)

// PostChangePassword operation middleware
func (siw *ServerInterfaceWrapper) PostChangePassword(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostChangePassword(c)
}

// GetAuthCheck operation middleware
func (siw *ServerInterfaceWrapper) GetAuthCheck(c *gin.Context) {

	var err error

	c.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAuthCheckParams

	// ------------- Optional query parameter "authToken" -------------
	if paramValue := c.Query("authToken"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "authToken", c.Request.URL.Query(), &params.AuthToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter authToken: %s", err)})
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetAuthCheck(c, params)
}

// PostGenerateCode operation middleware
func (siw *ServerInterfaceWrapper) PostGenerateCode(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostGenerateCode(c)
}

// PostAuthLogin operation middleware
func (siw *ServerInterfaceWrapper) PostAuthLogin(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostAuthLogin(c)
}

// PostAuthLoginByCode operation middleware
func (siw *ServerInterfaceWrapper) PostAuthLoginByCode(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostAuthLoginByCode(c)
}

// PostNewPassword operation middleware
func (siw *ServerInterfaceWrapper) PostNewPassword(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostNewPassword(c)
}

// PostResetPassword operation middleware
func (siw *ServerInterfaceWrapper) PostResetPassword(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostResetPassword(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL     string
	Middlewares []MiddlewareFunc
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router *gin.Engine, si ServerInterface) *gin.Engine {
	return RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router *gin.Engine, si ServerInterface, options GinServerOptions) *gin.Engine {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
	}

	router.POST(options.BaseURL+"/api/1/auth/changePassword", wrapper.PostChangePassword)

	router.GET(options.BaseURL+"/api/1/auth/check", wrapper.GetAuthCheck)

	router.POST(options.BaseURL+"/api/1/auth/generateCode", wrapper.PostGenerateCode)

	router.POST(options.BaseURL+"/api/1/auth/login", wrapper.PostAuthLogin)

	router.POST(options.BaseURL+"/api/1/auth/loginByCode", wrapper.PostAuthLoginByCode)

	router.POST(options.BaseURL+"/api/1/auth/newPassword", wrapper.PostNewPassword)

	router.POST(options.BaseURL+"/api/1/auth/resetPassword", wrapper.PostResetPassword)

	return router
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+yZ3W4bRRTHX2U0cNFKy65TigR710SiCh9taYu4KLmY2hN7qb1jZtaJTGXJH5RSpR9Q",
	"xC0V4gVMiKlxEvsVznkjdGbWyTq206gNppFyt7Mzc+bMmf9vzuzsfZ5XlaqKZZwYHt7nRuZrOkrqt/Il",
	"WZH21bIUWuortaREpbu29LHSFZHwkH/y1W3ucWNb8zCt5R5P6lUql5KkyhuNhsejeF1R/4I0eR1Vk0jF",
	"PORXbqyyqyKRm6LORFxgVF4hh0xELZiReiPKS7auNFsRuqjItNQVc339lqtKBwmDwGyKYlFqP1KBbRJQ",
	"2ygpU5OiG4N7fENq48Ze8pf8Jd7wuKrKWFQjHvL3/Zyf4x6viqRk5x6IahQsBaKWlIJ8ScRFeUMYs6l0",
	"YXou8AKb0Id96MMedKGHbQa7MIS/6CW7AEPYZ/A39FjNSB2LivTY4TvYhyHsQQ+bDNvQg13o4fdUDV0G",
	"fdiFPoMu7GATethi+IQawADb2LSNqPc/DEYwxIfYxq2LHsMWtqGLTdyimr4bYduVRlQBQzLrM3jOsAW7",
	"2MGH5DbDDrZgBD18QGNvk0+wjz/TdEZUh2032BDbtrbtfx1TzJRJKCiqKrWgkKwWeMhvKJOsTAbO41p+",
	"W5MmWVaFOvXIqziRse0sqtVylLfdg28MhfW+E5egp3e1XOch9916uPd+XVTK7wSHKk7fmyPrdTMzZsMq",
	"UktTVbFxIr+Uu3zyEQ56HjxdUyvpJBoev5zLTdjKq0pFxSewJrVWelmMfSVbH7yRrdU4IamVOc03EUXD",
	"wzucpsbXGh43tUpF6DpJ9w+nFqcRGNDCu0eS7zAjGHzGsEWqoKod0uyB0Ga0pRlMIiTz92aQ8xyG8BK2",
	"sQldfDRmZ8f6s49bpMm2lXsHH0HPjdKFbRKiRe6ldbWbcpA6QQ48Hlc5pvCZBSFreLa5Lv4A/QxY2CLs",
	"sAV9C8yvxIsjcoBt6MP2GMCTm6EA97Dt0dSGszFzDma8daAV5QzOrsqEdugVG+ApbedeR9t2sW6mxeuf",
	"nrayL+eW3sDWlzHNQunoO1mw1i599AbWbiv1uYjrqXfmlMEb51Qe3pnMpnc4X2uszSCzKrSoyERqY/uc",
	"1InDbnaNb6t7Ml6Nv6hJXec0Tgb4F5YP2vJ7NmNNINbHBxndUWE0n6mnzJFwUuF7jJ4dP5SvWjaV7cA+",
	"PiM+s0h1j+4fRRmT6OWKKsi3PAH7DH6xsSWu+9jEjnNoYDe5NgxsNm5CjwZNAaekSpGwcey6OLvX1Hpn",
	"dpT7Nv6UjEe21zbtc26sdG+jKUzO6WBGY387c3fMVyX3q9klWVRqz+rgPLEfl9h/y+jCpvUpgY0zppUY",
	"tZgpsqMkllUxit9+BE9+tJ27gVkA7Zl3QKC+igfa2T+zwVkUDHYpXknBax0BUtP/1RFgQQz8NL20lGls",
	"TjsQJ22B/ezh9elMyS/Xz0DucRnhmEzy/5KRxnChfLgxzyk5HUpmJZE9qrPqws5RdGK5eUbuTez03GkL",
	"W/AnNcMW2cp81S7gyuRaJmCL4iSzSOdnqtO8LBlZiQ7xET7Gpwwf4BP88fC45URKX0e4BbswsIU5yjvK",
	"lZZGJmeErAVQc3MiHIviZmIRzsk5jpzfx7qehmQ48aFiLw0cKF12YSYjnSMU9udTeJFhhxTkLiytUO2d",
	"B+zNvdRIL270xvgOZvq/hYqZSUQxiotMxhuRVnGFVsbjNV1O/0qYMAjSXw9+Xuii8u++p6Up+brGG94c",
	"o2WVF2UWxetazDWW/hcxvm1cUnYVj7VXkHdrxQl7YRAc9A4/zOVyPHMXNbWVzDwecI/TtsFDt9qNtca/",
	"AQAA//8Jz6i4UBoAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	pathPrefix := path.Dir(pathToFile)

	for rawPath, rawFunc := range externalRef0.PathToRawSpec(path.Join(pathPrefix, "./auth/schema.yaml")) {
		if _, ok := res[rawPath]; ok {
			// it is not possible to compare functions in golang, so always overwrite the old value
		}
		res[rawPath] = rawFunc
	}
	for rawPath, rawFunc := range externalRef1.PathToRawSpec(path.Join(pathPrefix, "./common/schema.yaml")) {
		if _, ok := res[rawPath]; ok {
			// it is not possible to compare functions in golang, so always overwrite the old value
		}
		res[rawPath] = rawFunc
	}
	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
