// Code generated by protoc-gen-go.
// source: api.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	api.proto

It has these top-level messages:
	StatusRequest
	StatusResponse
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gengo/grpc-gateway/third_party/googleapis/google/api"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type StatusRequest struct {
}

func (m *StatusRequest) Reset()                    { *m = StatusRequest{} }
func (m *StatusRequest) String() string            { return proto.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()               {}
func (*StatusRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type StatusResponse struct {
	Version   string `protobuf:"bytes,1,opt,name=version" json:"version,omitempty"`
	BuildTime int64  `protobuf:"varint,2,opt,name=buildTime" json:"buildTime,omitempty"`
	BuildHash string `protobuf:"bytes,3,opt,name=buildHash" json:"buildHash,omitempty"`
	AppStart  int64  `protobuf:"varint,4,opt,name=appStart" json:"appStart,omitempty"`
}

func (m *StatusResponse) Reset()                    { *m = StatusResponse{} }
func (m *StatusResponse) String() string            { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()               {}
func (*StatusResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*StatusRequest)(nil), "pb.StatusRequest")
	proto.RegisterType((*StatusResponse)(nil), "pb.StatusResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion1

// Client API for StatusService service

type StatusServiceClient interface {
	GetStatus(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
}

type statusServiceClient struct {
	cc *grpc.ClientConn
}

func NewStatusServiceClient(cc *grpc.ClientConn) StatusServiceClient {
	return &statusServiceClient{cc}
}

func (c *statusServiceClient) GetStatus(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := grpc.Invoke(ctx, "/pb.StatusService/GetStatus", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for StatusService service

type StatusServiceServer interface {
	GetStatus(context.Context, *StatusRequest) (*StatusResponse, error)
}

func RegisterStatusServiceServer(s *grpc.Server, srv StatusServiceServer) {
	s.RegisterService(&_StatusService_serviceDesc, srv)
}

func _StatusService_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(StatusServiceServer).GetStatus(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _StatusService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.StatusService",
	HandlerType: (*StatusServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStatus",
			Handler:    _StatusService_GetStatus_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 218 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x4c, 0x2c, 0xc8, 0xd4,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x92, 0x92, 0x49, 0xcf, 0xcf, 0x4f, 0xcf,
	0x49, 0xd5, 0x07, 0x8a, 0xea, 0x27, 0xe6, 0xe5, 0xe5, 0x97, 0x24, 0x96, 0x64, 0xe6, 0xe7, 0x15,
	0x43, 0x54, 0x28, 0xf1, 0x73, 0xf1, 0x06, 0x03, 0x45, 0x4a, 0x8b, 0x83, 0x52, 0x0b, 0x4b, 0x53,
	0x8b, 0x4b, 0x94, 0x9a, 0x18, 0xb9, 0xf8, 0x60, 0x22, 0xc5, 0x05, 0x40, 0x85, 0xa9, 0x42, 0x12,
	0x5c, 0xec, 0x65, 0xa9, 0x45, 0xc5, 0x40, 0x5d, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x30,
	0xae, 0x90, 0x0c, 0x17, 0x67, 0x52, 0x69, 0x66, 0x4e, 0x4a, 0x48, 0x66, 0x6e, 0xaa, 0x04, 0x13,
	0x50, 0x8e, 0x39, 0x08, 0x21, 0x00, 0x97, 0xf5, 0x48, 0x2c, 0xce, 0x90, 0x60, 0x06, 0xeb, 0x44,
	0x08, 0x08, 0x49, 0x71, 0x71, 0x24, 0x16, 0x14, 0x00, 0xad, 0x2a, 0x2a, 0x91, 0x60, 0x01, 0x6b,
	0x85, 0xf3, 0x8d, 0xa2, 0x61, 0xae, 0x0a, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0x15, 0xf2, 0xe2,
	0xe2, 0x74, 0x4f, 0x2d, 0x81, 0x88, 0x09, 0x09, 0xea, 0x15, 0x24, 0xe9, 0xa1, 0xb8, 0x5a, 0x4a,
	0x08, 0x59, 0x08, 0xe2, 0x6c, 0x25, 0xb1, 0xa6, 0xcb, 0x4f, 0x26, 0x33, 0x09, 0x28, 0x71, 0x83,
	0xbd, 0x5e, 0x0c, 0x96, 0xb4, 0x62, 0xd4, 0x4a, 0x62, 0x03, 0xfb, 0xdc, 0x18, 0x10, 0x00, 0x00,
	0xff, 0xff, 0x7f, 0xe5, 0xc8, 0x7e, 0x28, 0x01, 0x00, 0x00,
}