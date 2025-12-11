package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/DangeL187/erax"

	"device/internal/shared/config"
	"device/internal/shared/tokens"
)

var (
	ErrRetry         = errors.New("retry")
	ErrLoginRequired = errors.New("login required")
	ErrWrongCreds    = errors.New("wrong credentials")
)

type AuthService struct {
	cfg    *config.Config
	tokens *tokens.Tokens
	client *http.Client

	auth   chan struct{}
	authMu sync.RWMutex
}

func (as *AuthService) Login(ctx context.Context) bool {
	zap.L().Info("Logining...")
	url := as.cfg.AuthServerURL + "/devices/login"

	for {
		select {
		case <-ctx.Done():
			return false
		default:
		}

		accessToken, refreshToken, err := as.login(url, as.cfg.DeviceID, as.cfg.DevicePassword)
		if err == nil {
			as.tokens.SetAccess(accessToken)
			as.tokens.SetRefresh(refreshToken)

			as.authMu.Lock()
			select {
			case <-as.auth:
			default:
				close(as.auth)
			}
			as.authMu.Unlock()

			return true
		}

		err = erax.Wrap(err, "failed to login")
		zap.S().Errorf("\n%f", err)
		switch {
		case errors.Is(err, ErrRetry):
			zap.S().Infof("Retrying in %ds", int(as.cfg.ConnectRetryInterval.Seconds()))
			time.Sleep(as.cfg.ConnectRetryInterval)
		case errors.Is(err, ErrWrongCreds):
			return false
		default:
			return false
		}
	}
}

func (as *AuthService) Refresh(ctx context.Context) {
	zap.L().Info("Refreshing...")
	url := as.cfg.AuthServerURL + "/devices/refresh"

	for {
		select {
		case <-ctx.Done():
			return
		default:
			accessToken, err := as.refresh(url, as.tokens.GetRefresh())
			if err == nil {
				as.tokens.SetAccess(accessToken)

				select {
				case <-as.auth:
				default:
					close(as.auth)
				}
				return
			}

			err = erax.Wrap(err, "failed to refresh access token")
			switch {
			case errors.Is(err, ErrRetry):
				zap.S().Errorf("\n%f", err)
				zap.S().Infof("Retrying in %ds", int(as.cfg.ConnectRetryInterval.Seconds()))
				time.Sleep(as.cfg.ConnectRetryInterval)
			case errors.Is(err, ErrLoginRequired):
				zap.L().Info("Access token expired")
				as.Login(ctx)
				return
			default:
				zap.S().Errorf("\n%f", err)
				return
			}
		}
	}
}

func (as *AuthService) ResetAuth() {
	as.authMu.Lock()
	as.auth = make(chan struct{})
	as.authMu.Unlock()
}

func (as *AuthService) WaitForAuth(ctx context.Context) bool {
	as.authMu.RLock()
	ch := as.auth
	as.authMu.RUnlock()

	select {
	case <-ctx.Done():
		return false
	case <-ch:
	}
	return true
}

func (as *AuthService) login(url, deviceID, password string) (string, string, error) {
	body := map[string]string{"device_id": deviceID, "password": password}
	data, _ := json.Marshal(body)

	resp, err := as.client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return "", "", erax.WrapWithError(err, ErrRetry, "failed to send request")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		return "", "", erax.WithMeta(ErrWrongCreds, "response_data", string(respData))
	}

	var accessToken, refreshToken string
	for _, c := range resp.Cookies() {
		switch c.Name {
		case "access_token":
			accessToken = c.Value
		case "refresh_token":
			refreshToken = c.Value
		}
	}

	if accessToken == "" || refreshToken == "" {
		return "", "", errors.New("tokens not found in cookies")
	}

	zap.L().Debug("logged in", zap.String("access", accessToken), zap.String("refresh", refreshToken))
	return accessToken, refreshToken, nil
}

func (as *AuthService) refresh(url, refreshToken string) (string, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", erax.WrapWithError(err, ErrRetry, "failed to create request")
	}

	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken})

	resp, err := as.client.Do(req)
	if err != nil {
		return "", erax.WrapWithError(err, ErrRetry, "failed to send request")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		return "", erax.WithMeta(ErrWrongCreds, "response_data", string(respData))
	}

	var accessToken string
	for _, c := range resp.Cookies() {
		switch c.Name {
		case "access_token":
			accessToken = c.Value
		}
	}

	if accessToken == "" {
		return "", errors.New("tokens not found in cookies")
	}

	zap.L().Debug("token refreshed", zap.String("access", accessToken))
	return accessToken, nil
}

func NewAuthService(cfg *config.Config, tokens *tokens.Tokens) *AuthService {
	client := &http.Client{Timeout: time.Minute}

	return &AuthService{
		auth:   make(chan struct{}),
		cfg:    cfg,
		client: client,
		tokens: tokens,
	}
}
