package clients

import (
	"context"
	"errors"
	"time"

	v1 "gateway/api/auth/v1"
	"gateway/internal/pkg/logger"
	"gateway/internal/pkg/metrics"
	"gateway/internal/pkg/strings"

	"github.com/AlekSi/pointer"
	"github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	metricPrefix = `clients.auth`
)

var (
	ErrSessionExpiredOrNotFound = errors.New(`session expired or not found`)
)

type AuthClient interface {
	Check(ctx context.Context, token string) (*CheckResult, error)
	Login(ctx context.Context, username, password string, remember bool) (string, *time.Time, error)
	LoginByCode(ctx context.Context, username, code string, remember bool) (string, *time.Time, error)
	ResetPassword(ctx context.Context, username string) error
	NewPassword(ctx context.Context, username, password, passwordResetHash string) error
	ChangePassword(ctx context.Context, username, oldPassword, newPassword string) error
	GenerateCode(ctx context.Context, username string) error
}

type Auth struct {
	client v1.AuthClient
	metric metrics.Metrics
	logger *log.Helper
}

func NewAuth(
	ctx context.Context,
	endpoint string,
	timeout time.Duration,
	metric metrics.Metrics,
	logs log.Logger,
) (*Auth, error) {
	client, err := Default(ctx, endpoint, timeout)
	if err != nil {
		return nil, err
	}
	return &Auth{
		client: v1.NewAuthClient(client),
		metric: metric,
		logger: logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, `clients/auth`),
	}, nil
}

type User struct {
	Type        string  `json:"type"`
	DisplayName string  `json:"displayName"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
}

type Session struct {
	Until     time.Time `json:"until"`
	IP        *string   `json:"IP,omitempty"`
	UserAgent *string   `json:"userAgent,omitempty"`
	DeviceID  *string   `json:"deviceId,omitempty"`
}

type CheckResult struct {
	User    *User    `json:"user"`
	Session *Session `json:"session"`
}

func (a *Auth) postProcess(ctx context.Context, method string, err error) {
	if err != nil {
		a.logger.WithContext(ctx).Errorf(`client auth method %s failed: %v`, method, err)
		a.metric.Increment(strings.Metric(metricPrefix, method, `failure`))
	} else {
		a.metric.Increment(strings.Metric(metricPrefix, method, `success`))
	}
}

func StatsFromContext(ctx context.Context) *v1.Stats {
	g, ok := gin.FromGinContext(ctx)
	if !ok {
		return nil
	}
	stats := &v1.Stats{
		Ip:        g.ClientIP(),
		UserAgent: g.GetHeader("User-Agent"),
	}
	deviceID := g.GetHeader("DeviceId")
	if deviceID != `` {
		stats.DeviceId = &deviceID
	}
	return stats
}

func (a *Auth) Check(ctx context.Context, token string) (*CheckResult, error) {
	method := `check`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	res, err := a.client.Check(ctx, &v1.CheckRequest{Token: token})
	if err != nil {
		if statusErr, ok := status.FromError(err); ok {
			if statusErr.Code() == codes.NotFound {
				return nil, ErrSessionExpiredOrNotFound
			}
		}
		return nil, err
	}
	return &CheckResult{
		User: &User{
			Type:        res.User.Type,
			DisplayName: res.User.DisplayName,
			Email:       res.User.Email,
			Phone:       res.User.Phone,
		},
		Session: &Session{
			Until:     res.Session.Until.AsTime(),
			IP:        res.Session.Ip,
			UserAgent: res.Session.UserAgent,
			DeviceID:  res.Session.DeviceId,
		},
	}, nil
}

func (a *Auth) Login(ctx context.Context, username, password string, remember bool) (string, *time.Time, error) {
	method := `login`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	in := &v1.LoginRequest{
		Username: username,
		Password: password,
		Remember: remember,
		Stats:    StatsFromContext(ctx),
	}
	out, err := a.client.Login(ctx, in)
	if err != nil {
		return "", nil, err
	}
	return out.Token, pointer.ToTime(out.Until.AsTime()), nil
}

func (a *Auth) LoginByCode(ctx context.Context, username, code string, remember bool) (string, *time.Time, error) {
	method := `loginByCode`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	in := &v1.LoginByCodeRequest{
		Username: username,
		Code:     code,
		Remember: remember,
		Stats:    StatsFromContext(ctx),
	}
	out, err := a.client.LoginByCode(ctx, in)
	if err != nil {
		return "", nil, err
	}
	return out.Token, pointer.ToTime(out.Until.AsTime()), nil
}

func (a *Auth) ResetPassword(ctx context.Context, username string) error {
	method := `resetPassword`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	in := &v1.ResetPasswordRequest{
		Username: username,
		Stats:    StatsFromContext(ctx),
	}
	_, err = a.client.ResetPassword(ctx, in)
	return err
}

func (a *Auth) NewPassword(ctx context.Context, username, password, passwordResetHash string) error {
	method := `newPassword`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	in := &v1.NewPasswordRequest{
		Username:          username,
		Password:          password,
		PasswordResetHash: passwordResetHash,
		Stats:             StatsFromContext(ctx),
	}
	_, err = a.client.NewPassword(ctx, in)
	return err
}

func (a *Auth) ChangePassword(ctx context.Context, username, oldPassword, newPassword string) error {
	method := `changePassword`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	in := &v1.ChangePasswordRequest{
		Username:    username,
		OldPassword: oldPassword,
		NewPassword: newPassword,
		Stats:       StatsFromContext(ctx),
	}
	_, err = a.client.ChangePassword(ctx, in)
	return err
}

func (a *Auth) GenerateCode(ctx context.Context, username string) error {
	method := `generateCode`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	in := &v1.GenerateCodeRequest{
		Username: username,
		Stats:    StatsFromContext(ctx),
	}
	_, err = a.client.GenerateCode(ctx, in)
	return err
}

// TODO Implement History(context.Context, *HistoryRequest) (*HistoryResponse, error)
