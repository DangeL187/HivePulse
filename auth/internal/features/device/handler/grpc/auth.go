package grpc

import (
	"context"
	"log"

	"github.com/DangeL187/erax"

	"auth/internal/app"
	pb "auth/internal/infra/grpc/proto/auth"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	app *app.App
}

func (a *AuthHandler) AuthDevice(_ context.Context, req *pb.AuthDeviceRequest) (*pb.AuthDeviceResponse, error) {
	deviceID, tokenType, err := a.app.JWTManager.ParseToken(req.Token)
	if err != nil {
		err = erax.Wrap(err, "failed to parse token")
		log.Printf("\n%f\n", err)

		return &pb.AuthDeviceResponse{
			DeviceId: 0,
			Error:    "unauthorized",
		}, nil
	}
	if tokenType != "access" {
		log.Printf("wrong token type (%s)\n", tokenType)

		return &pb.AuthDeviceResponse{
			DeviceId: 0,
			Error:    "unauthorized",
		}, nil
	}

	return &pb.AuthDeviceResponse{
		DeviceId: uint64(deviceID),
		Error:    "",
	}, nil
}

func (a *AuthHandler) GetPublicKey(_ context.Context, _ *pb.GetPublicKeyRequest) (*pb.GetPublicKeyResponse, error) {
	publicKey, err := a.app.JWTManager.GetPublicKey()
	if err != nil {
		return nil, erax.Wrap(err, "failed to get public key")
	}

	return &pb.GetPublicKeyResponse{
		PublicKey: publicKey,
	}, nil
}

func NewAuthHandler(app *app.App) *AuthHandler {
	return &AuthHandler{app: app}
}
