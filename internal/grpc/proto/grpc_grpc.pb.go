// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: proto/grpc.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// URLsClient is the client API for URLs service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type URLsClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	GetList(ctx context.Context, in *GetListRequest, opts ...grpc.CallOption) (*GetListResponse, error)
	Save(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveResponse, error)
	SaveList(ctx context.Context, in *SaveListRequest, opts ...grpc.CallOption) (*SaveListResponse, error)
	DelList(ctx context.Context, in *DelListRequest, opts ...grpc.CallOption) (*DelListResponse, error)
}

type uRLsClient struct {
	cc grpc.ClientConnInterface
}

func NewURLsClient(cc grpc.ClientConnInterface) URLsClient {
	return &uRLsClient{cc}
}

func (c *uRLsClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/grpc.URLs/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLsClient) GetList(ctx context.Context, in *GetListRequest, opts ...grpc.CallOption) (*GetListResponse, error) {
	out := new(GetListResponse)
	err := c.cc.Invoke(ctx, "/grpc.URLs/GetList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLsClient) Save(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveResponse, error) {
	out := new(SaveResponse)
	err := c.cc.Invoke(ctx, "/grpc.URLs/Save", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLsClient) SaveList(ctx context.Context, in *SaveListRequest, opts ...grpc.CallOption) (*SaveListResponse, error) {
	out := new(SaveListResponse)
	err := c.cc.Invoke(ctx, "/grpc.URLs/SaveList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLsClient) DelList(ctx context.Context, in *DelListRequest, opts ...grpc.CallOption) (*DelListResponse, error) {
	out := new(DelListResponse)
	err := c.cc.Invoke(ctx, "/grpc.URLs/DelList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// URLsServer is the server API for URLs service.
// All implementations must embed UnimplementedURLsServer
// for forward compatibility
type URLsServer interface {
	Get(context.Context, *GetRequest) (*GetResponse, error)
	GetList(context.Context, *GetListRequest) (*GetListResponse, error)
	Save(context.Context, *SaveRequest) (*SaveResponse, error)
	SaveList(context.Context, *SaveListRequest) (*SaveListResponse, error)
	DelList(context.Context, *DelListRequest) (*DelListResponse, error)
	mustEmbedUnimplementedURLsServer()
}

// UnimplementedURLsServer must be embedded to have forward compatible implementations.
type UnimplementedURLsServer struct {
}

func (UnimplementedURLsServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedURLsServer) GetList(context.Context, *GetListRequest) (*GetListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetList not implemented")
}
func (UnimplementedURLsServer) Save(context.Context, *SaveRequest) (*SaveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Save not implemented")
}
func (UnimplementedURLsServer) SaveList(context.Context, *SaveListRequest) (*SaveListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveList not implemented")
}
func (UnimplementedURLsServer) DelList(context.Context, *DelListRequest) (*DelListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelList not implemented")
}
func (UnimplementedURLsServer) mustEmbedUnimplementedURLsServer() {}

// UnsafeURLsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to URLsServer will
// result in compilation errors.
type UnsafeURLsServer interface {
	mustEmbedUnimplementedURLsServer()
}

func RegisterURLsServer(s grpc.ServiceRegistrar, srv URLsServer) {
	s.RegisterService(&URLs_ServiceDesc, srv)
}

func _URLs_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLsServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.URLs/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLsServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLs_GetList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLsServer).GetList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.URLs/GetList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLsServer).GetList(ctx, req.(*GetListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLs_Save_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLsServer).Save(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.URLs/Save",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLsServer).Save(ctx, req.(*SaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLs_SaveList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLsServer).SaveList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.URLs/SaveList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLsServer).SaveList(ctx, req.(*SaveListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLs_DelList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLsServer).DelList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.URLs/DelList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLsServer).DelList(ctx, req.(*DelListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// URLs_ServiceDesc is the grpc.ServiceDesc for URLs service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var URLs_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.URLs",
	HandlerType: (*URLsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _URLs_Get_Handler,
		},
		{
			MethodName: "GetList",
			Handler:    _URLs_GetList_Handler,
		},
		{
			MethodName: "Save",
			Handler:    _URLs_Save_Handler,
		},
		{
			MethodName: "SaveList",
			Handler:    _URLs_SaveList_Handler,
		},
		{
			MethodName: "DelList",
			Handler:    _URLs_DelList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/grpc.proto",
}