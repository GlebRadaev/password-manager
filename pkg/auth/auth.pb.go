// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.12.4
// source: auth.proto

package auth

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RegisterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Email    string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	mi := &file_auth_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *RegisterRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *RegisterRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type RegisterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *RegisterResponse) Reset() {
	*x = RegisterResponse{}
	mi := &file_auth_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterResponse) ProtoMessage() {}

func (x *RegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterResponse.ProtoReflect.Descriptor instead.
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterResponse) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *RegisterResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type LoginRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	mi := &file_auth_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{2}
}

func (x *LoginRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *LoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type LoginResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token     string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	ExpiresIn int64  `protobuf:"varint,2,opt,name=expires_in,json=expiresIn,proto3" json:"expires_in,omitempty"`
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	mi := &file_auth_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{3}
}

func (x *LoginResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *LoginResponse) GetExpiresIn() int64 {
	if x != nil {
		return x.ExpiresIn
	}
	return 0
}

type ValidateTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *ValidateTokenRequest) Reset() {
	*x = ValidateTokenRequest{}
	mi := &file_auth_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidateTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateTokenRequest) ProtoMessage() {}

func (x *ValidateTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateTokenRequest.ProtoReflect.Descriptor instead.
func (*ValidateTokenRequest) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{4}
}

func (x *ValidateTokenRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type ValidateTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Valid  bool   `protobuf:"varint,1,opt,name=valid,proto3" json:"valid,omitempty"`
	UserId string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *ValidateTokenResponse) Reset() {
	*x = ValidateTokenResponse{}
	mi := &file_auth_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidateTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateTokenResponse) ProtoMessage() {}

func (x *ValidateTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateTokenResponse.ProtoReflect.Descriptor instead.
func (*ValidateTokenResponse) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{5}
}

func (x *ValidateTokenResponse) GetValid() bool {
	if x != nil {
		return x.Valid
	}
	return false
}

func (x *ValidateTokenResponse) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type GenerateOTPRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *GenerateOTPRequest) Reset() {
	*x = GenerateOTPRequest{}
	mi := &file_auth_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GenerateOTPRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateOTPRequest) ProtoMessage() {}

func (x *GenerateOTPRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateOTPRequest.ProtoReflect.Descriptor instead.
func (*GenerateOTPRequest) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{6}
}

func (x *GenerateOTPRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type GenerateOTPResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OtpCode string `protobuf:"bytes,1,opt,name=otp_code,json=otpCode,proto3" json:"otp_code,omitempty"`
}

func (x *GenerateOTPResponse) Reset() {
	*x = GenerateOTPResponse{}
	mi := &file_auth_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GenerateOTPResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateOTPResponse) ProtoMessage() {}

func (x *GenerateOTPResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateOTPResponse.ProtoReflect.Descriptor instead.
func (*GenerateOTPResponse) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{7}
}

func (x *GenerateOTPResponse) GetOtpCode() string {
	if x != nil {
		return x.OtpCode
	}
	return ""
}

type ValidateOTPRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	OtpCode string `protobuf:"bytes,2,opt,name=otp_code,json=otpCode,proto3" json:"otp_code,omitempty"`
}

func (x *ValidateOTPRequest) Reset() {
	*x = ValidateOTPRequest{}
	mi := &file_auth_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidateOTPRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateOTPRequest) ProtoMessage() {}

func (x *ValidateOTPRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateOTPRequest.ProtoReflect.Descriptor instead.
func (*ValidateOTPRequest) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{8}
}

func (x *ValidateOTPRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ValidateOTPRequest) GetOtpCode() string {
	if x != nil {
		return x.OtpCode
	}
	return ""
}

type ValidateOTPResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Valid bool `protobuf:"varint,1,opt,name=valid,proto3" json:"valid,omitempty"`
}

func (x *ValidateOTPResponse) Reset() {
	*x = ValidateOTPResponse{}
	mi := &file_auth_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidateOTPResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateOTPResponse) ProtoMessage() {}

func (x *ValidateOTPResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateOTPResponse.ProtoReflect.Descriptor instead.
func (*ValidateOTPResponse) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{9}
}

func (x *ValidateOTPResponse) GetValid() bool {
	if x != nil {
		return x.Valid
	}
	return false
}

var File_auth_proto protoreflect.FileDescriptor

var file_auth_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x70,
	0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8f, 0x01,
	0x0a, 0x0f, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x36, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x1a, 0xfa, 0x42, 0x17, 0x72, 0x15, 0x10, 0x03, 0x18, 0x14, 0x32, 0x0f,
	0x5e, 0x5b, 0x61, 0x2d, 0x7a, 0x41, 0x2d, 0x5a, 0x30, 0x2d, 0x39, 0x5f, 0x5d, 0x2b, 0x24, 0x52,
	0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x08, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xfa, 0x42, 0x06,
	0x72, 0x04, 0x10, 0x08, 0x18, 0x32, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x12, 0x1d, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x60, 0x01, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22,
	0x45, 0x0a, 0x10, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x5c, 0x0a, 0x0c, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x72, 0x04, 0x10,
	0x03, 0x18, 0x14, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x09, 0xfa, 0x42, 0x06, 0x72, 0x04, 0x10, 0x08, 0x18, 0x32, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x22, 0x44, 0x0a, 0x0d, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x65,
	0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x49, 0x6e, 0x22, 0x35, 0x0a, 0x14, 0x56, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1d, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x0a, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x22, 0x46, 0x0a, 0x15, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x37, 0x0a, 0x12, 0x47, 0x65, 0x6e,
	0x65, 0x72, 0x61, 0x74, 0x65, 0x4f, 0x54, 0x50, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x21, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x22, 0x30, 0x0a, 0x13, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x4f, 0x54,
	0x50, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x74, 0x70,
	0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x74, 0x70,
	0x43, 0x6f, 0x64, 0x65, 0x22, 0x5c, 0x0a, 0x12, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x4f, 0x54, 0x50, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05,
	0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x23, 0x0a,
	0x08, 0x6f, 0x74, 0x70, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x98, 0x01, 0x06, 0x52, 0x07, 0x6f, 0x74, 0x70, 0x43, 0x6f,
	0x64, 0x65, 0x22, 0x2b, 0x0a, 0x13, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x4f, 0x54,
	0x50, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x32,
	0x92, 0x04, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x5f, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x19, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74,
	0x68, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x3a, 0x01, 0x2a, 0x22, 0x11, 0x2f,
	0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x12, 0x53, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x16, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x61, 0x75, 0x74, 0x68, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x4c, 0x6f, 0x67,
	0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x13, 0x3a, 0x01, 0x2a, 0x22, 0x0e, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f,
	0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x71, 0x0a, 0x0d, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74,
	0x68, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74,
	0x68, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x12,
	0x17, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2d, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x6c, 0x0a, 0x0b, 0x47, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x65, 0x4f, 0x54, 0x50, 0x12, 0x1c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75,
	0x74, 0x68, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x4f, 0x54, 0x50, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x4f, 0x54, 0x50, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x3a, 0x01, 0x2a, 0x22,
	0x15, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61,
	0x74, 0x65, 0x2d, 0x6f, 0x74, 0x70, 0x12, 0x6c, 0x0a, 0x0b, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x4f, 0x54, 0x50, 0x12, 0x1c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x4f, 0x54, 0x50, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x56,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x4f, 0x54, 0x50, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x3a, 0x01, 0x2a, 0x22, 0x15, 0x2f,
	0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2d, 0x6f, 0x74, 0x70, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_auth_proto_rawDescOnce sync.Once
	file_auth_proto_rawDescData = file_auth_proto_rawDesc
)

func file_auth_proto_rawDescGZIP() []byte {
	file_auth_proto_rawDescOnce.Do(func() {
		file_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_auth_proto_rawDescData)
	})
	return file_auth_proto_rawDescData
}

var file_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_auth_proto_goTypes = []any{
	(*RegisterRequest)(nil),       // 0: api.auth.RegisterRequest
	(*RegisterResponse)(nil),      // 1: api.auth.RegisterResponse
	(*LoginRequest)(nil),          // 2: api.auth.LoginRequest
	(*LoginResponse)(nil),         // 3: api.auth.LoginResponse
	(*ValidateTokenRequest)(nil),  // 4: api.auth.ValidateTokenRequest
	(*ValidateTokenResponse)(nil), // 5: api.auth.ValidateTokenResponse
	(*GenerateOTPRequest)(nil),    // 6: api.auth.GenerateOTPRequest
	(*GenerateOTPResponse)(nil),   // 7: api.auth.GenerateOTPResponse
	(*ValidateOTPRequest)(nil),    // 8: api.auth.ValidateOTPRequest
	(*ValidateOTPResponse)(nil),   // 9: api.auth.ValidateOTPResponse
}
var file_auth_proto_depIdxs = []int32{
	0, // 0: api.auth.AuthService.Register:input_type -> api.auth.RegisterRequest
	2, // 1: api.auth.AuthService.Login:input_type -> api.auth.LoginRequest
	4, // 2: api.auth.AuthService.ValidateToken:input_type -> api.auth.ValidateTokenRequest
	6, // 3: api.auth.AuthService.GenerateOTP:input_type -> api.auth.GenerateOTPRequest
	8, // 4: api.auth.AuthService.ValidateOTP:input_type -> api.auth.ValidateOTPRequest
	1, // 5: api.auth.AuthService.Register:output_type -> api.auth.RegisterResponse
	3, // 6: api.auth.AuthService.Login:output_type -> api.auth.LoginResponse
	5, // 7: api.auth.AuthService.ValidateToken:output_type -> api.auth.ValidateTokenResponse
	7, // 8: api.auth.AuthService.GenerateOTP:output_type -> api.auth.GenerateOTPResponse
	9, // 9: api.auth.AuthService.ValidateOTP:output_type -> api.auth.ValidateOTPResponse
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_auth_proto_init() }
func file_auth_proto_init() {
	if File_auth_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_auth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_auth_proto_goTypes,
		DependencyIndexes: file_auth_proto_depIdxs,
		MessageInfos:      file_auth_proto_msgTypes,
	}.Build()
	File_auth_proto = out.File
	file_auth_proto_rawDesc = nil
	file_auth_proto_goTypes = nil
	file_auth_proto_depIdxs = nil
}
