// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: proto/proposal.proto

package proposal

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
	ProposalService_ProposalRPC_FullMethodName = "/lerpc.ProposalService/ProposalRPC"
)

// ProposalServiceClient is the client API for ProposalService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProposalServiceClient interface {
	ProposalRPC(ctx context.Context, in *Proposal, opts ...grpc.CallOption) (*Proposal, error)
}

type proposalServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProposalServiceClient(cc grpc.ClientConnInterface) ProposalServiceClient {
	return &proposalServiceClient{cc}
}

func (c *proposalServiceClient) ProposalRPC(ctx context.Context, in *Proposal, opts ...grpc.CallOption) (*Proposal, error) {
	out := new(Proposal)
	err := c.cc.Invoke(ctx, ProposalService_ProposalRPC_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProposalServiceServer is the server API for ProposalService service.
// All implementations must embed UnimplementedProposalServiceServer
// for forward compatibility
type ProposalServiceServer interface {
	ProposalRPC(context.Context, *Proposal) (*Proposal, error)
	mustEmbedUnimplementedProposalServiceServer()
}

// UnimplementedProposalServiceServer must be embedded to have forward compatible implementations.
type UnimplementedProposalServiceServer struct {
}

func (UnimplementedProposalServiceServer) ProposalRPC(context.Context, *Proposal) (*Proposal, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProposalRPC not implemented")
}
func (UnimplementedProposalServiceServer) mustEmbedUnimplementedProposalServiceServer() {}

// UnsafeProposalServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProposalServiceServer will
// result in compilation errors.
type UnsafeProposalServiceServer interface {
	mustEmbedUnimplementedProposalServiceServer()
}

func RegisterProposalServiceServer(s grpc.ServiceRegistrar, srv ProposalServiceServer) {
	s.RegisterService(&ProposalService_ServiceDesc, srv)
}

func _ProposalService_ProposalRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Proposal)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProposalServiceServer).ProposalRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProposalService_ProposalRPC_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProposalServiceServer).ProposalRPC(ctx, req.(*Proposal))
	}
	return interceptor(ctx, in, info, handler)
}

// ProposalService_ServiceDesc is the grpc.ServiceDesc for ProposalService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProposalService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lerpc.ProposalService",
	HandlerType: (*ProposalServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProposalRPC",
			Handler:    _ProposalService_ProposalRPC_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/proposal.proto",
}