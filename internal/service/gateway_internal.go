package service

import (
	"context"

	"gateway/api/gateway"
	"gateway/internal/clients"
	"github.com/go-kratos/kratos/v2/errors"
)

func (s *GatewayService) getAuthCheck(ctx context.Context, token string) (*gateway.CheckResponseOK, error) {
	if token == "" {
		return nil, errors.BadRequest(`token_not_found`, `not found authToken in query or X-Auth-Token in headers`)
	}
	res, err := s.authClient.Check(ctx, token)
	if err == clients.ErrSessionExpiredOrNotFound {
		return nil, errors.Unauthorized(`invalid_auth_token`, `auth session is expired or not found`)
	}
	if err != nil {
		return nil, errors.InternalServer(`auth_client_error`, err.Error())
	}

	response := &gateway.CheckResponseOK{}

	if res.User != nil {
		response.User.DisplayName = res.User.DisplayName
		response.User.Type = res.User.Type
		response.User.Phone = res.User.Phone
		response.User.Email = res.User.Email
	}

	if res.Session != nil {
		response.Session.Until = res.Session.Until
		response.Session.Ip = res.Session.IP
		response.Session.UserAgent = res.Session.UserAgent
	}

	return response, nil
}
