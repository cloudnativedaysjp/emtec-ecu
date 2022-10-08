// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: pkg/ws-proxy/schema/track.proto

package schema

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TrackServiceClient is the client API for TrackService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TrackServiceClient interface {
	GetTrack(ctx context.Context, in *GetTrackRequest, opts ...grpc.CallOption) (*Track, error)
	ListTrack(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListTrackResponse, error)
	EnableAutomation(ctx context.Context, in *SwitchAutomationRequest, opts ...grpc.CallOption) (*Track, error)
	DisableAutomation(ctx context.Context, in *SwitchAutomationRequest, opts ...grpc.CallOption) (*Track, error)
}

type trackServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTrackServiceClient(cc grpc.ClientConnInterface) TrackServiceClient {
	return &trackServiceClient{cc}
}

func (c *trackServiceClient) GetTrack(ctx context.Context, in *GetTrackRequest, opts ...grpc.CallOption) (*Track, error) {
	out := new(Track)
	err := c.cc.Invoke(ctx, "/schema.TrackService/GetTrack", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) ListTrack(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListTrackResponse, error) {
	out := new(ListTrackResponse)
	err := c.cc.Invoke(ctx, "/schema.TrackService/ListTrack", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) EnableAutomation(ctx context.Context, in *SwitchAutomationRequest, opts ...grpc.CallOption) (*Track, error) {
	out := new(Track)
	err := c.cc.Invoke(ctx, "/schema.TrackService/EnableAutomation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) DisableAutomation(ctx context.Context, in *SwitchAutomationRequest, opts ...grpc.CallOption) (*Track, error) {
	out := new(Track)
	err := c.cc.Invoke(ctx, "/schema.TrackService/DisableAutomation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TrackServiceServer is the server API for TrackService service.
// All implementations must embed UnimplementedTrackServiceServer
// for forward compatibility
type TrackServiceServer interface {
	GetTrack(context.Context, *GetTrackRequest) (*Track, error)
	ListTrack(context.Context, *emptypb.Empty) (*ListTrackResponse, error)
	EnableAutomation(context.Context, *SwitchAutomationRequest) (*Track, error)
	DisableAutomation(context.Context, *SwitchAutomationRequest) (*Track, error)
	mustEmbedUnimplementedTrackServiceServer()
}

// UnimplementedTrackServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTrackServiceServer struct {
}

func (UnimplementedTrackServiceServer) GetTrack(context.Context, *GetTrackRequest) (*Track, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTrack not implemented")
}
func (UnimplementedTrackServiceServer) ListTrack(context.Context, *emptypb.Empty) (*ListTrackResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTrack not implemented")
}
func (UnimplementedTrackServiceServer) EnableAutomation(context.Context, *SwitchAutomationRequest) (*Track, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnableAutomation not implemented")
}
func (UnimplementedTrackServiceServer) DisableAutomation(context.Context, *SwitchAutomationRequest) (*Track, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisableAutomation not implemented")
}
func (UnimplementedTrackServiceServer) mustEmbedUnimplementedTrackServiceServer() {}

// UnsafeTrackServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TrackServiceServer will
// result in compilation errors.
type UnsafeTrackServiceServer interface {
	mustEmbedUnimplementedTrackServiceServer()
}

func RegisterTrackServiceServer(s grpc.ServiceRegistrar, srv TrackServiceServer) {
	s.RegisterService(&TrackService_ServiceDesc, srv)
}

func _TrackService_GetTrack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTrackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).GetTrack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schema.TrackService/GetTrack",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).GetTrack(ctx, req.(*GetTrackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_ListTrack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).ListTrack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schema.TrackService/ListTrack",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).ListTrack(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_EnableAutomation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SwitchAutomationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).EnableAutomation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schema.TrackService/EnableAutomation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).EnableAutomation(ctx, req.(*SwitchAutomationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_DisableAutomation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SwitchAutomationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).DisableAutomation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schema.TrackService/DisableAutomation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).DisableAutomation(ctx, req.(*SwitchAutomationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TrackService_ServiceDesc is the grpc.ServiceDesc for TrackService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TrackService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "schema.TrackService",
	HandlerType: (*TrackServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTrack",
			Handler:    _TrackService_GetTrack_Handler,
		},
		{
			MethodName: "ListTrack",
			Handler:    _TrackService_ListTrack_Handler,
		},
		{
			MethodName: "EnableAutomation",
			Handler:    _TrackService_EnableAutomation_Handler,
		},
		{
			MethodName: "DisableAutomation",
			Handler:    _TrackService_DisableAutomation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/ws-proxy/schema/track.proto",
}
