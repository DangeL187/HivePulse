package infra

import (
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"github.com/DangeL187/erax"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "ingress/internal/infra/grpc/proto/auth"
)

type GRPCAuthenticator struct {
	grpcAuthClient pb.AuthServiceClient
	grpcClientConn *grpc.ClientConn

	publicKey ed25519.PublicKey
}

func (a *GRPCAuthenticator) Auth(ctx context.Context, deviceToken string) error {
	if a.publicKey != nil {
		if err := a.verifyToken(deviceToken); err == nil {
			return nil
		}

		if err := a.updatePublicJWT(ctx); err != nil {
			return erax.Wrap(err, "failed to update public JWT")
		}

		if err := a.verifyToken(deviceToken); err == nil {
			return nil
		}

		return errors.New("unauthorized after refreshing public JWT")
	}

	if err := a.updatePublicJWT(ctx); err != nil {
		return erax.Wrap(err, "failed to update public JWT")
	}

	if err := a.verifyToken(deviceToken); err != nil {
		return erax.Wrap(err, "failed to verify device token")
	}

	return nil
}

func (a *GRPCAuthenticator) Close() error {
	err := a.grpcClientConn.Close()
	if err != nil {
		return erax.Wrap(err, "failed to close gRPC client connection")
	}
	return nil
}

func (a *GRPCAuthenticator) updatePublicJWT(ctx context.Context) error {
	resp, err := a.grpcAuthClient.GetPublicKey(ctx, &pb.GetPublicKeyRequest{})
	if err != nil {
		return erax.Wrap(err, "failed to get public key")
	}

	pubDER, err := base64.StdEncoding.DecodeString(resp.PublicKey)
	if err != nil {
		return errors.New("failed to decode base64 public key: " + err.Error())
	}

	pubIfc, err := x509.ParsePKIXPublicKey(pubDER)
	if err != nil {
		return errors.New("failed to parse DER public key: " + err.Error())
	}

	pubKey, ok := pubIfc.(ed25519.PublicKey)
	if !ok {
		return errors.New("not an ed25519 public key")
	}

	a.publicKey = pubKey

	return nil
}

func (a *GRPCAuthenticator) verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodEdDSA {
			return nil, errors.New("invalid signing method")
		}
		return a.publicKey, nil
	})

	if err != nil {
		return erax.Wrap(err, "invalid token")
	}
	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}

func NewGRPCAuthenticator(grpcAddr string) (*GRPCAuthenticator, error) {
	a := &GRPCAuthenticator{}

	var err error
	a.grpcClientConn, err = grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, erax.Wrap(err, "failed to create grpc client")
	}

	a.grpcAuthClient = pb.NewAuthServiceClient(a.grpcClientConn)

	return a, nil
}
