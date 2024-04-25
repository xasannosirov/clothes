// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: media_service.proto

package media_service

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("media_service.proto", fileDescriptor_cad0f3d72e604dd5) }

var fileDescriptor_cad0f3d72e604dd5 = []byte{
	// 191 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xce, 0x4d, 0x4d, 0xc9,
	0x4c, 0x8c, 0x2f, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0xe2, 0x45, 0x11, 0x94, 0x12, 0x84, 0x70, 0x73, 0xf3, 0x53, 0x52, 0x73, 0x20, 0x2a, 0x8c, 0x1e,
	0x33, 0x72, 0xf1, 0xf8, 0x82, 0x44, 0x83, 0x21, 0x6a, 0x84, 0xac, 0xb8, 0xd8, 0x9c, 0x8b, 0x52,
	0x13, 0x4b, 0x52, 0x85, 0x44, 0xf4, 0x50, 0x8d, 0x04, 0x2b, 0x93, 0x92, 0xc2, 0x26, 0x1a, 0x9e,
	0x59, 0x92, 0xe1, 0xe9, 0x22, 0xe4, 0xc6, 0xc5, 0xec, 0x9e, 0x5a, 0x22, 0xa4, 0x88, 0x4b, 0x49,
	0x40, 0x51, 0x7e, 0x4a, 0x69, 0x72, 0x89, 0xa7, 0x8b, 0x94, 0x0c, 0x9a, 0x12, 0x98, 0x4c, 0x6e,
	0x62, 0x7a, 0x6a, 0xb1, 0x90, 0x3f, 0x17, 0x9b, 0x4b, 0x6a, 0x4e, 0x6a, 0x49, 0x2a, 0x31, 0x46,
	0x29, 0xa1, 0x29, 0x81, 0xe8, 0x04, 0x2b, 0x0c, 0x4a, 0x2d, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0x75,
	0xd2, 0x38, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x67, 0x3c,
	0x96, 0x63, 0x88, 0x12, 0x4b, 0x4f, 0xcd, 0x03, 0x87, 0x80, 0x3e, 0x8a, 0xee, 0x24, 0x36, 0xb0,
	0xa0, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xe0, 0xb3, 0x7e, 0xc8, 0x4f, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MediaServiceClient is the client API for MediaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MediaServiceClient interface {
	Create(ctx context.Context, in *Media, opts ...grpc.CallOption) (*MediaWithID, error)
	Get(ctx context.Context, in *MediaWithProductID, opts ...grpc.CallOption) (*ProductImages, error)
	Delete(ctx context.Context, in *MediaWithProductID, opts ...grpc.CallOption) (*DeleteMediaResponse, error)
}

type mediaServiceClient struct {
	cc *grpc.ClientConn
}

func NewMediaServiceClient(cc *grpc.ClientConn) MediaServiceClient {
	return &mediaServiceClient{cc}
}

func (c *mediaServiceClient) Create(ctx context.Context, in *Media, opts ...grpc.CallOption) (*MediaWithID, error) {
	out := new(MediaWithID)
	err := c.cc.Invoke(ctx, "/media_service.MediaService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mediaServiceClient) Get(ctx context.Context, in *MediaWithProductID, opts ...grpc.CallOption) (*ProductImages, error) {
	out := new(ProductImages)
	err := c.cc.Invoke(ctx, "/media_service.MediaService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mediaServiceClient) Delete(ctx context.Context, in *MediaWithProductID, opts ...grpc.CallOption) (*DeleteMediaResponse, error) {
	out := new(DeleteMediaResponse)
	err := c.cc.Invoke(ctx, "/media_service.MediaService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MediaServiceServer is the server API for MediaService service.
type MediaServiceServer interface {
	Create(context.Context, *Media) (*MediaWithID, error)
	Get(context.Context, *MediaWithProductID) (*ProductImages, error)
	Delete(context.Context, *MediaWithProductID) (*DeleteMediaResponse, error)
}

// UnimplementedMediaServiceServer can be embedded to have forward compatible implementations.
type UnimplementedMediaServiceServer struct {
}

func (*UnimplementedMediaServiceServer) Create(ctx context.Context, req *Media) (*MediaWithID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedMediaServiceServer) Get(ctx context.Context, req *MediaWithProductID) (*ProductImages, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedMediaServiceServer) Delete(ctx context.Context, req *MediaWithProductID) (*DeleteMediaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

func RegisterMediaServiceServer(s *grpc.Server, srv MediaServiceServer) {
	s.RegisterService(&_MediaService_serviceDesc, srv)
}

func _MediaService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Media)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MediaServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/media_service.MediaService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MediaServiceServer).Create(ctx, req.(*Media))
	}
	return interceptor(ctx, in, info, handler)
}

func _MediaService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MediaWithProductID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MediaServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/media_service.MediaService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MediaServiceServer).Get(ctx, req.(*MediaWithProductID))
	}
	return interceptor(ctx, in, info, handler)
}

func _MediaService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MediaWithProductID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MediaServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/media_service.MediaService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MediaServiceServer).Delete(ctx, req.(*MediaWithProductID))
	}
	return interceptor(ctx, in, info, handler)
}

var _MediaService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "media_service.MediaService",
	HandlerType: (*MediaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _MediaService_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _MediaService_Get_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _MediaService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "media_service.proto",
}
