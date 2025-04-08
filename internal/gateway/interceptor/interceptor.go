// Package interceptor provides gRPC interceptors for authentication and authorization.
// It includes middleware for validating JWT tokens in incoming requests.
package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/GlebRadaev/password-manager/pkg/auth"
)

type contextKey string

const (
	// UserIDKey is the context key
	UserIDKey contextKey = "user_id"
	// AuthorizationHeader is the HTTP header key for authorization
	AuthorizationHeader = "authorization"
	// BearerPrefix is the prefix for Bearer tokens
	BearerPrefix = "Bearer "
)

// AuthInterceptor creates a gRPC unary interceptor that validates JWT tokens.
// It skips validation for public methods (login, register, etc.) and injects
// the user ID into the context for authenticated requests.
func AuthInterceptor(authClient auth.AuthServiceClient) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if isPublicMethod(method) {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		var token string

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if authHeader := md.Get(AuthorizationHeader); len(authHeader) > 0 {
				if strings.HasPrefix(authHeader[0], BearerPrefix) {
					token = strings.TrimPrefix(authHeader[0], BearerPrefix)
				}
			}
		}

		if token == "" {
			if header, ok := metadata.FromOutgoingContext(ctx); ok {
				if authHeader := header.Get(AuthorizationHeader); len(authHeader) > 0 {
					if strings.HasPrefix(authHeader[0], BearerPrefix) {
						token = strings.TrimPrefix(authHeader[0], BearerPrefix)
					}
				}
			}
		}

		if token == "" {
			return status.Error(codes.Unauthenticated, "authorization token missing")
		}

		resp, err := authClient.ValidateToken(ctx, &auth.ValidateTokenRequest{Token: token})
		if err != nil || !resp.Valid {
			return status.Error(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, UserIDKey, resp.UserId)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// isPublicMethod checks if the given gRPC method should bypass authentication.
func isPublicMethod(method string) bool {
	publicMethods := map[string]bool{
		"/api.auth.AuthService/Register":      true,
		"/api.auth.AuthService/Login":         true,
		"/api.auth.AuthService/ValidateToken": true,
		"/api.auth.AuthService/RefreshToken":  true,
	}
	return publicMethods[method]
}
