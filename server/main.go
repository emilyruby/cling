package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"github.com/emilyruby/cling/api"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"github.com/dgrijalva/jwt-go"
)

const (
	port = ":50051"
)

var jwtKey = []byte("AllYourBase") // TODO: read this from env

type UserClaims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

type server struct{}

func checkTokenValidity(ctx context.Context, token string) (context.Context, error) {
	// Initialize a new instance of `Claims`
	claims := &UserClaims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if !tkn.Valid {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	ctx = context.WithValue(ctx, "user", claims.Username)
	return ctx, nil
}

func newSessionTokenForUser(username string) (string, error) {
	// Create the Claims
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := UserClaims{
		username,
		jwt.StandardClaims{
			// TODO: add expiriration
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "cling",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func checkLoginCredentials(username string, password string) error {
	return nil
}

func (s *server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	if fullMethodName == "/Cling/Login" {
		return ctx, nil
	} else {
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "no token")
		}
		return checkTokenValidity(ctx, token)
	}
}

func (s *server) Login(ctx context.Context, in *api.LoginRequest) (*api.LoginReply, error) {
	err := checkLoginCredentials(in.Username, in.Password)
	if err != nil {
		return nil, err
	}

	token, err := newSessionTokenForUser(in.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	md := metadata.Pairs("authorization", "Bearer " + token)
	err = grpc.SetHeader(ctx, md)
	if err != nil {
		return nil, err
	}

	return &api.LoginReply{}, nil
}

func (s *server) NewPost(ctx context.Context, in *api.Post) (*api.PostConfirmation, error) {
	return &api.PostConfirmation{PostID: "123"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer( grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
        grpc_auth.UnaryServerInterceptor(nil),
    )))
	api.RegisterClingServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}