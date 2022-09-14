// Package auth provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package auth

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	externalRef0 "gateway/api/gateway/common"

	"github.com/getkin/kin-openapi/openapi3"
)

// CheckResponse defines model for checkResponse.
type CheckResponse struct {
	// Данные авторизационной сессии
	Session struct {
		// Возможный DeviceID если запрос пришёл от мобильного устройства
		DeviceId *PropertyDeviceId `json:"deviceId,omitempty"`

		// IP-адрес
		Ip *PropertyIP `json:"ip,omitempty"`

		// Время по UTC, до которого сессия активна
		Until PropertySessionUntil `json:"until"`

		// User-Agent
		UserAgent *PropertyUserAgent `json:"userAgent,omitempty"`
	} `json:"session"`

	// Данные об авторизованном пользователе
	User struct {
		// Отображаемое имя пользователя
		DisplayName PropertyUserDisplayName `json:"displayName"`

		// Электронная почта
		Email *PropertyUserEmail `json:"email,omitempty"`

		// Российский номер мобильного телефона в формате 9009009090
		Phone *PropertyUserPhoneTyped `json:"phone,omitempty"`

		// Тип пользователя из набора (admin|dispatcher|driver)
		Type PropertyUserType `json:"type"`
	} `json:"user"`
}

// LoginRequestBody defines model for loginRequestBody.
type LoginRequestBody struct {
	// Пароль пользователя — от 8 до 255 символов
	Password PropertyPassword `json:"password"`

	// Номер телефона (в произвольной форме) или адрес электронной почты
	Username PropertyUsername `json:"username"`
}

// LoginResponse defines model for loginResponse.
type LoginResponse struct {
	// Авторизационный токен для пользовательской сессии
	Token externalRef0.PropertyAuthToken `json:"token"`

	// Время по UTC, до которого сессия активна
	Until PropertySessionUntil `json:"until"`
}

// Возможный DeviceID если запрос пришёл от мобильного устройства
type PropertyDeviceId = string

// IP-адрес
type PropertyIP = string

// Пароль пользователя — от 8 до 255 символов
type PropertyPassword = string

// Время по UTC, до которого сессия активна
type PropertySessionUntil = time.Time

// User-Agent
type PropertyUserAgent = string

// Отображаемое имя пользователя
type PropertyUserDisplayName = string

// Электронная почта
type PropertyUserEmail = string

// Российский номер мобильного телефона в формате 9009009090
type PropertyUserPhoneTyped = string

// Тип пользователя из набора (admin|dispatcher|driver)
type PropertyUserType = string

// Номер телефона (в произвольной форме) или адрес электронной почты
type PropertyUsername = string

// Возможный DeviceID если запрос пришёл от мобильного устройства
type DeviceId = PropertyDeviceId

// IP-адрес
type Ip = PropertyIP

// User-Agent
type UserAgent = PropertyUserAgent

// CheckResponseOK defines model for checkResponseOK.
type CheckResponseOK = CheckResponse

// LoginResponseOK defines model for loginResponseOK.
type LoginResponseOK = LoginResponse

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"/4yEWoF3EAAA",
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

	for rawPath, rawFunc := range externalRef0.PathToRawSpec(path.Join(pathPrefix, "../common/schema.yaml")) {
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
