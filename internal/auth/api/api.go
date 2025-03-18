package api

//go:generate mockgen -destination=api_mock.go -source=api.go -package=api
import (
	"context"
	"errors"

	"github.com/GlebRadaev/password-manager/internal/auth/models"
	"github.com/GlebRadaev/password-manager/internal/auth/service"

	"github.com/GlebRadaev/password-manager/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	Register(ctx context.Context, username, password, email string) (string, error)
	Login(ctx context.Context, username, password string) (string, int64, error)
	ValidateToken(ctx context.Context, token string) (bool, string, error)
	GenerateOTP(ctx context.Context, userID, deviceID string) (string, error)
	ValidateOTP(ctx context.Context, userID, otpCode, deviceID string) (bool, error)
	ListSessions(ctx context.Context, userID string) ([]models.Session, error)
	TerminateSession(ctx context.Context, sessionID string) error
}

type Api struct {
	auth.UnimplementedAuthServiceServer
	srv Service
}

func New(srv Service) *Api {
	return &Api{srv: srv}
}

func (a *Api) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if err := ValidateRegisterRequest(req); err != nil {
		return nil, err
	}

	userID, err := a.srv.Register(ctx, req.Username, req.Password, req.Email)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to register: %v", err)
	}

	return &auth.RegisterResponse{UserId: userID, Message: "User registered successfully"}, nil
}

func (a *Api) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if err := ValidateLoginRequest(req); err != nil {
		return nil, err
	}

	token, expiresIn, err := a.srv.Login(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Errorf(codes.Internal, "login failed: %v", err)
	}

	return &auth.LoginResponse{Token: token, ExpiresIn: expiresIn}, nil
}

func (a *Api) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	if err := ValidateTokenRequest(req); err != nil {
		return nil, err
	}

	valid, userID, err := a.srv.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token validation failed: %v", err)
	}

	return &auth.ValidateTokenResponse{Valid: valid, UserId: userID}, nil
}

func (a *Api) GenerateOTP(ctx context.Context, req *auth.GenerateOTPRequest) (*auth.GenerateOTPResponse, error) {
	if err := ValidateGenerateOTPRequest(req); err != nil {
		return nil, err
	}

	otpCode, err := a.srv.GenerateOTP(ctx, req.UserId, req.DeviceId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate OTP: %v", err)
	}

	return &auth.GenerateOTPResponse{OtpCode: otpCode}, nil
}

func (a *Api) ValidateOTP(ctx context.Context, req *auth.ValidateOTPRequest) (*auth.ValidateOTPResponse, error) {
	if err := ValidateOTPRequest(req); err != nil {
		return nil, err
	}

	valid, err := a.srv.ValidateOTP(ctx, req.UserId, req.OtpCode, req.DeviceId)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOTP) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid OTP")
		}
		return nil, status.Errorf(codes.Internal, "failed to validate OTP: %v", err)
	}

	return &auth.ValidateOTPResponse{Valid: valid}, nil
}

// Новые методы для управления сессиями
func (a *Api) ListSessions(ctx context.Context, req *auth.ListSessionsRequest) (*auth.ListSessionsResponse, error) {
	sessions, err := a.srv.ListSessions(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list sessions: %v", err)
	}

	var protoSessions []*auth.Session
	for _, session := range sessions {
		protoSessions = append(protoSessions, &auth.Session{
			SessionId:  session.SessionID,
			DeviceInfo: session.DeviceInfo,
			CreatedAt:  session.CreatedAt.Unix(),
			ExpiresAt:  session.ExpiresAt.Unix(),
		})
	}

	return &auth.ListSessionsResponse{Sessions: protoSessions}, nil
}

func (a *Api) TerminateSession(ctx context.Context, req *auth.TerminateSessionRequest) (*auth.TerminateSessionResponse, error) {
	if err := a.srv.TerminateSession(ctx, req.SessionId); err != nil {
		if errors.Is(err, service.ErrSessionNotFound) {
			return nil, status.Errorf(codes.NotFound, "session not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to terminate session: %v", err)
	}

	return &auth.TerminateSessionResponse{Message: "Session terminated successfully"}, nil
}
