package biz

import (
	"context"

	"gateway/api/gateway"
	"gateway/api/gateway/auth"
	"gateway/internal/clients"
	"gateway/internal/pkg/logger"
	"gateway/internal/pkg/metrics"
	"gateway/internal/pkg/request"
	"github.com/AlekSi/pointer"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type GatewayUsecase struct {
	authClient clients.AuthClient
	metric     metrics.Metrics
	logger     *log.Helper
}

func NewGatewayUsecase(
	authClient clients.AuthClient,
	metric metrics.Metrics,
	logs log.Logger,
) *GatewayUsecase {
	return &GatewayUsecase{
		authClient: authClient,
		metric:     metric,
		logger:     logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "service/gateway"),
	}
}

func (u *GatewayUsecase) GetAuthCheck(ctx context.Context, params gateway.GetAuthCheckParams) (
	*auth.CheckResponseOK,
	error,
) {
	token := request.AuthTokenFromContext(ctx)
	if params.AuthToken != nil {
		token = *params.AuthToken
	}
	if token == "" {
		return nil, errors.BadRequest(`token_not_found`, `not found authToken in query or in authorization header`)
	}
	res, err := u.authClient.Check(ctx, token)
	if err != nil {
		if errors.Is(err, clients.ErrSessionExpiredOrNotFound) {
			return nil, errors.Unauthorized(`invalid_auth_token`, `auth session is expired or not found`)
		}
		return nil, errors.InternalServer(`auth_client_error`, err.Error())
	}

	response := &auth.CheckResponseOK{}

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

func (u *GatewayUsecase) PostAuthLogin(
	ctx context.Context,
	in *gateway.PostAuthLoginJSONRequestBody,
) (*auth.LoginResponseOK, error) {
	authToken, until, err := u.authClient.Login(ctx, in.Username, in.Password, pointer.GetBool(in.Remember))
	if err != nil {
		return nil, err
	}
	return &auth.LoginResponseOK{
		Token: authToken,
		Until: pointer.GetTime(until),
	}, nil
}

func (u *GatewayUsecase) PostGenerateCode(
	ctx context.Context,
	in *gateway.PostGenerateCodeJSONRequestBody,
) error {
	return u.authClient.GenerateCode(ctx, in.Username)
}

func (u *GatewayUsecase) PostAuthLoginByCode(
	ctx context.Context,
	in *gateway.PostAuthLoginByCodeJSONRequestBody,
) (*auth.LoginResponseOK, error) {
	authToken, until, err := u.authClient.LoginByCode(ctx, in.Username, in.Code, pointer.GetBool(in.Remember))
	if err != nil {
		return nil, err
	}
	return &auth.LoginResponseOK{
		Token: authToken,
		Until: pointer.GetTime(until),
	}, nil
}

func (u *GatewayUsecase) PostResetPassword(ctx context.Context, in *gateway.PostResetPasswordJSONRequestBody) error {
	return u.authClient.ResetPassword(ctx, in.Username)
}

func (u *GatewayUsecase) PostChangePassword(ctx context.Context, in *gateway.PostChangePasswordJSONRequestBody) error {
	return u.authClient.ChangePassword(ctx, in.Username, in.OldPassword, in.NewPassword)
}

func (u *GatewayUsecase) PostNewPassword(ctx context.Context, in *gateway.PostNewPasswordJSONRequestBody) error {
	return u.authClient.NewPassword(ctx, in.Username, in.Password, in.PasswordResetHash)
}
