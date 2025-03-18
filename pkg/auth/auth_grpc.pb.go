// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: auth.proto

package auth

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	AuthService_Register_FullMethodName         = "/api.auth.AuthService/Register"
	AuthService_Login_FullMethodName            = "/api.auth.AuthService/Login"
	AuthService_ValidateToken_FullMethodName    = "/api.auth.AuthService/ValidateToken"
	AuthService_GenerateOTP_FullMethodName      = "/api.auth.AuthService/GenerateOTP"
	AuthService_ValidateOTP_FullMethodName      = "/api.auth.AuthService/ValidateOTP"
	AuthService_ListSessions_FullMethodName     = "/api.auth.AuthService/ListSessions"
	AuthService_TerminateSession_FullMethodName = "/api.auth.AuthService/TerminateSession"
)

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthServiceClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	ValidateToken(ctx context.Context, in *ValidateTokenRequest, opts ...grpc.CallOption) (*ValidateTokenResponse, error)
	GenerateOTP(ctx context.Context, in *GenerateOTPRequest, opts ...grpc.CallOption) (*GenerateOTPResponse, error)
	ValidateOTP(ctx context.Context, in *ValidateOTPRequest, opts ...grpc.CallOption) (*ValidateOTPResponse, error)
	ListSessions(ctx context.Context, in *ListSessionsRequest, opts ...grpc.CallOption) (*ListSessionsResponse, error)
	TerminateSession(ctx context.Context, in *TerminateSessionRequest, opts ...grpc.CallOption) (*TerminateSessionResponse, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, AuthService_Register_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, AuthService_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ValidateToken(ctx context.Context, in *ValidateTokenRequest, opts ...grpc.CallOption) (*ValidateTokenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ValidateTokenResponse)
	err := c.cc.Invoke(ctx, AuthService_ValidateToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) GenerateOTP(ctx context.Context, in *GenerateOTPRequest, opts ...grpc.CallOption) (*GenerateOTPResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GenerateOTPResponse)
	err := c.cc.Invoke(ctx, AuthService_GenerateOTP_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ValidateOTP(ctx context.Context, in *ValidateOTPRequest, opts ...grpc.CallOption) (*ValidateOTPResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ValidateOTPResponse)
	err := c.cc.Invoke(ctx, AuthService_ValidateOTP_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ListSessions(ctx context.Context, in *ListSessionsRequest, opts ...grpc.CallOption) (*ListSessionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListSessionsResponse)
	err := c.cc.Invoke(ctx, AuthService_ListSessions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) TerminateSession(ctx context.Context, in *TerminateSessionRequest, opts ...grpc.CallOption) (*TerminateSessionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TerminateSessionResponse)
	err := c.cc.Invoke(ctx, AuthService_TerminateSession_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
// All implementations must embed UnimplementedAuthServiceServer
// for forward compatibility.
type AuthServiceServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	ValidateToken(context.Context, *ValidateTokenRequest) (*ValidateTokenResponse, error)
	GenerateOTP(context.Context, *GenerateOTPRequest) (*GenerateOTPResponse, error)
	ValidateOTP(context.Context, *ValidateOTPRequest) (*ValidateOTPResponse, error)
	ListSessions(context.Context, *ListSessionsRequest) (*ListSessionsResponse, error)
	TerminateSession(context.Context, *TerminateSessionRequest) (*TerminateSessionResponse, error)
	mustEmbedUnimplementedAuthServiceServer()
}

// UnimplementedAuthServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAuthServiceServer struct{}

func (UnimplementedAuthServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedAuthServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAuthServiceServer) ValidateToken(context.Context, *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateToken not implemented")
}
func (UnimplementedAuthServiceServer) GenerateOTP(context.Context, *GenerateOTPRequest) (*GenerateOTPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateOTP not implemented")
}
func (UnimplementedAuthServiceServer) ValidateOTP(context.Context, *ValidateOTPRequest) (*ValidateOTPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateOTP not implemented")
}
func (UnimplementedAuthServiceServer) ListSessions(context.Context, *ListSessionsRequest) (*ListSessionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSessions not implemented")
}
func (UnimplementedAuthServiceServer) TerminateSession(context.Context, *TerminateSessionRequest) (*TerminateSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TerminateSession not implemented")
}
func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}
func (UnimplementedAuthServiceServer) testEmbeddedByValue()                     {}

// UnsafeAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServiceServer will
// result in compilation errors.
type UnsafeAuthServiceServer interface {
	mustEmbedUnimplementedAuthServiceServer()
}

func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer) {
	// If the following call pancis, it indicates UnimplementedAuthServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

func _AuthService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ValidateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ValidateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_ValidateToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ValidateToken(ctx, req.(*ValidateTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_GenerateOTP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateOTPRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).GenerateOTP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_GenerateOTP_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).GenerateOTP(ctx, req.(*GenerateOTPRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ValidateOTP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateOTPRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ValidateOTP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_ValidateOTP_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ValidateOTP(ctx, req.(*ValidateOTPRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ListSessions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSessionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ListSessions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_ListSessions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ListSessions(ctx, req.(*ListSessionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_TerminateSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TerminateSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).TerminateSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_TerminateSession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).TerminateSession(ctx, req.(*TerminateSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthService_ServiceDesc is the grpc.ServiceDesc for AuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.auth.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _AuthService_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _AuthService_Login_Handler,
		},
		{
			MethodName: "ValidateToken",
			Handler:    _AuthService_ValidateToken_Handler,
		},
		{
			MethodName: "GenerateOTP",
			Handler:    _AuthService_GenerateOTP_Handler,
		},
		{
			MethodName: "ValidateOTP",
			Handler:    _AuthService_ValidateOTP_Handler,
		},
		{
			MethodName: "ListSessions",
			Handler:    _AuthService_ListSessions_Handler,
		},
		{
			MethodName: "TerminateSession",
			Handler:    _AuthService_TerminateSession_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}
