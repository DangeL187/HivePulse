package runtime

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"sync"
	"time"

	"github.com/DangeL187/erax"
)

type authResponse struct {
	DeviceID string
	ErrorMsg string
}

type authenticator interface {
	Auth(ctx context.Context, deviceToken string) error
	Close() error
}

type publisher interface {
	Publish(topic string, payload any) error
}

type AuthService struct {
	errRespChan chan authResponse
	errRespWg   sync.WaitGroup

	authenticator authenticator
	publisher     publisher
}

func (as *AuthService) Run(ctx context.Context) {
	as.runErrRespWorkers(ctx, 5)
}

func (as *AuthService) Stop() {
	err := as.authenticator.Close()
	if err != nil {
		zap.L().Error("failed to close authenticator", zap.Error(err))
	}

	close(as.errRespChan)

	as.errRespWg.Wait()
}

func (as *AuthService) Auth(deviceID, deviceToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := as.authenticator.Auth(ctx, deviceToken)
	if err != nil {
		as.errRespChan <- authResponse{
			DeviceID: deviceID,
			ErrorMsg: err.Error(),
		}
		return erax.Wrap(err, "failed to auth device")
	}

	return nil
}

func (as *AuthService) runErrRespWorkers(ctx context.Context, workerCount int) {
	as.errRespWg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func() {
			defer as.errRespWg.Done()
			for {
				select {
				case job, ok := <-as.errRespChan:
					if !ok {
						return
					}
					payload, err := json.Marshal(map[string]string{"error": job.ErrorMsg})
					if err != nil {
						zap.L().Error("failed to marshal auth response", zap.Error(err))
						continue
					}
					topic := "devices/" + job.DeviceID + "/auth_response"
					err = as.publisher.Publish(topic, payload)
					if err != nil {
						zap.L().Error("failed to publish auth response", zap.Error(err))
						continue
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

func NewAuthService(authenticator authenticator, publisher publisher) (*AuthService, error) {
	s := &AuthService{
		errRespChan:   make(chan authResponse, 1024),
		authenticator: authenticator,
		publisher:     publisher,
	}

	return s, nil
}
