// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: p2p.proto

package p2p_pb

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

const (
	PeerService_GetPeers_FullMethodName = "/p2p_proto.PeerService/GetPeers"
)

// PeerServiceClient is the client API for PeerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PeerServiceClient interface {
	GetPeers(ctx context.Context, in *GetPeersRequest, opts ...grpc.CallOption) (*GetPeersResponse, error)
}

type peerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPeerServiceClient(cc grpc.ClientConnInterface) PeerServiceClient {
	return &peerServiceClient{cc}
}

func (c *peerServiceClient) GetPeers(ctx context.Context, in *GetPeersRequest, opts ...grpc.CallOption) (*GetPeersResponse, error) {
	out := new(GetPeersResponse)
	err := c.cc.Invoke(ctx, PeerService_GetPeers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PeerServiceServer is the server API for PeerService service.
// All implementations must embed UnimplementedPeerServiceServer
// for forward compatibility
type PeerServiceServer interface {
	GetPeers(context.Context, *GetPeersRequest) (*GetPeersResponse, error)
	mustEmbedUnimplementedPeerServiceServer()
}

// UnimplementedPeerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPeerServiceServer struct {
}

func (UnimplementedPeerServiceServer) GetPeers(context.Context, *GetPeersRequest) (*GetPeersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPeers not implemented")
}
func (UnimplementedPeerServiceServer) mustEmbedUnimplementedPeerServiceServer() {}

// UnsafePeerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PeerServiceServer will
// result in compilation errors.
type UnsafePeerServiceServer interface {
	mustEmbedUnimplementedPeerServiceServer()
}

func RegisterPeerServiceServer(s grpc.ServiceRegistrar, srv PeerServiceServer) {
	s.RegisterService(&PeerService_ServiceDesc, srv)
}

func _PeerService_GetPeers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPeersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerServiceServer).GetPeers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PeerService_GetPeers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerServiceServer).GetPeers(ctx, req.(*GetPeersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PeerService_ServiceDesc is the grpc.ServiceDesc for PeerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PeerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "p2p_proto.PeerService",
	HandlerType: (*PeerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPeers",
			Handler:    _PeerService_GetPeers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "p2p.proto",
}
