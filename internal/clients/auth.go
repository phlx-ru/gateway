package clients

import (
	"context"
	"errors"
	"time"

	v1 "gateway/api/auth/v1"
	"gateway/internal/pkg/logger"
	"gateway/internal/pkg/metrics"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	metricAuthCheckTimings = `clients.auth.check.timings`
	metricAuthCheckSuccess = `clients.auth.check.success`
	metricAuthCheckFailure = `clients.auth.check.failure`
)

var (
	ErrSessionExpiredOrNotFound = errors.New(`session expired or not found`)
)

type AuthClient interface {
	Check(ctx context.Context, token string) (*CheckResult, error)
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

func (a *Auth) Check(ctx context.Context, token string) (*CheckResult, error) {
	defer a.metric.NewTiming().Send(metricAuthCheckTimings)
	var err error
	defer func() {
		if err != nil {
			a.logger.WithContext(ctx).Errorf(`failed to check auth: %v`, err)
			a.metric.Increment(metricAuthCheckFailure)
		} else {
			a.metric.Increment(metricAuthCheckSuccess)
		}
	}()
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
