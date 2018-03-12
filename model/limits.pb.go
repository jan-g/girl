// Code generated by protoc-gen-go. DO NOT EDIT.
// source: limits.proto

/*
Package model is a generated protocol buffer package.

It is generated from these files:
	limits.proto

It has these top-level messages:
	IHave
	Push
	HostTraffic
	Traffic
	ConnectorHandshake
	ResponderHandshake
	UseRequest
	UseResponse
*/
package model

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/empty"

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
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type IHave struct {
	Epoch int64    `protobuf:"varint,1,opt,name=epoch" json:"epoch,omitempty"`
	Hosts []string `protobuf:"bytes,2,rep,name=hosts" json:"hosts,omitempty"`
}

func (m *IHave) Reset()                    { *m = IHave{} }
func (m *IHave) String() string            { return proto.CompactTextString(m) }
func (*IHave) ProtoMessage()               {}
func (*IHave) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *IHave) GetEpoch() int64 {
	if m != nil {
		return m.Epoch
	}
	return 0
}

func (m *IHave) GetHosts() []string {
	if m != nil {
		return m.Hosts
	}
	return nil
}

type Push struct {
	Traffic []*HostTraffic `protobuf:"bytes,1,rep,name=traffic" json:"traffic,omitempty"`
}

func (m *Push) Reset()                    { *m = Push{} }
func (m *Push) String() string            { return proto.CompactTextString(m) }
func (*Push) ProtoMessage()               {}
func (*Push) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Push) GetTraffic() []*HostTraffic {
	if m != nil {
		return m.Traffic
	}
	return nil
}

type HostTraffic struct {
	Name    string     `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Traffic []*Traffic `protobuf:"bytes,2,rep,name=traffic" json:"traffic,omitempty"`
}

func (m *HostTraffic) Reset()                    { *m = HostTraffic{} }
func (m *HostTraffic) String() string            { return proto.CompactTextString(m) }
func (*HostTraffic) ProtoMessage()               {}
func (*HostTraffic) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *HostTraffic) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *HostTraffic) GetTraffic() []*Traffic {
	if m != nil {
		return m.Traffic
	}
	return nil
}

type Traffic struct {
	Facet string `protobuf:"bytes,1,opt,name=facet" json:"facet,omitempty"`
	Usage int64  `protobuf:"varint,2,opt,name=usage" json:"usage,omitempty"`
}

func (m *Traffic) Reset()                    { *m = Traffic{} }
func (m *Traffic) String() string            { return proto.CompactTextString(m) }
func (*Traffic) ProtoMessage()               {}
func (*Traffic) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Traffic) GetFacet() string {
	if m != nil {
		return m.Facet
	}
	return ""
}

func (m *Traffic) GetUsage() int64 {
	if m != nil {
		return m.Usage
	}
	return 0
}

type ConnectorHandshake struct {
	// Types that are valid to be assigned to Hs:
	//	*ConnectorHandshake_IHave
	//	*ConnectorHandshake_Push
	Hs isConnectorHandshake_Hs `protobuf_oneof:"hs"`
}

func (m *ConnectorHandshake) Reset()                    { *m = ConnectorHandshake{} }
func (m *ConnectorHandshake) String() string            { return proto.CompactTextString(m) }
func (*ConnectorHandshake) ProtoMessage()               {}
func (*ConnectorHandshake) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type isConnectorHandshake_Hs interface {
	isConnectorHandshake_Hs()
}

type ConnectorHandshake_IHave struct {
	IHave *IHave `protobuf:"bytes,1,opt,name=iHave,oneof"`
}
type ConnectorHandshake_Push struct {
	Push *Push `protobuf:"bytes,2,opt,name=push,oneof"`
}

func (*ConnectorHandshake_IHave) isConnectorHandshake_Hs() {}
func (*ConnectorHandshake_Push) isConnectorHandshake_Hs()  {}

func (m *ConnectorHandshake) GetHs() isConnectorHandshake_Hs {
	if m != nil {
		return m.Hs
	}
	return nil
}

func (m *ConnectorHandshake) GetIHave() *IHave {
	if x, ok := m.GetHs().(*ConnectorHandshake_IHave); ok {
		return x.IHave
	}
	return nil
}

func (m *ConnectorHandshake) GetPush() *Push {
	if x, ok := m.GetHs().(*ConnectorHandshake_Push); ok {
		return x.Push
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ConnectorHandshake) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ConnectorHandshake_OneofMarshaler, _ConnectorHandshake_OneofUnmarshaler, _ConnectorHandshake_OneofSizer, []interface{}{
		(*ConnectorHandshake_IHave)(nil),
		(*ConnectorHandshake_Push)(nil),
	}
}

func _ConnectorHandshake_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ConnectorHandshake)
	// hs
	switch x := m.Hs.(type) {
	case *ConnectorHandshake_IHave:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.IHave); err != nil {
			return err
		}
	case *ConnectorHandshake_Push:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Push); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ConnectorHandshake.Hs has unexpected type %T", x)
	}
	return nil
}

func _ConnectorHandshake_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ConnectorHandshake)
	switch tag {
	case 1: // hs.iHave
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(IHave)
		err := b.DecodeMessage(msg)
		m.Hs = &ConnectorHandshake_IHave{msg}
		return true, err
	case 2: // hs.push
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Push)
		err := b.DecodeMessage(msg)
		m.Hs = &ConnectorHandshake_Push{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ConnectorHandshake_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ConnectorHandshake)
	// hs
	switch x := m.Hs.(type) {
	case *ConnectorHandshake_IHave:
		s := proto.Size(x.IHave)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ConnectorHandshake_Push:
		s := proto.Size(x.Push)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type ResponderHandshake struct {
	IWant []string `protobuf:"bytes,1,rep,name=iWant" json:"iWant,omitempty"`
	Push  *Push    `protobuf:"bytes,2,opt,name=push" json:"push,omitempty"`
}

func (m *ResponderHandshake) Reset()                    { *m = ResponderHandshake{} }
func (m *ResponderHandshake) String() string            { return proto.CompactTextString(m) }
func (*ResponderHandshake) ProtoMessage()               {}
func (*ResponderHandshake) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ResponderHandshake) GetIWant() []string {
	if m != nil {
		return m.IWant
	}
	return nil
}

func (m *ResponderHandshake) GetPush() *Push {
	if m != nil {
		return m.Push
	}
	return nil
}

type UseRequest struct {
	Facet    string `protobuf:"bytes,1,opt,name=facet" json:"facet,omitempty"`
	Quantity int64  `protobuf:"varint,2,opt,name=quantity" json:"quantity,omitempty"`
}

func (m *UseRequest) Reset()                    { *m = UseRequest{} }
func (m *UseRequest) String() string            { return proto.CompactTextString(m) }
func (*UseRequest) ProtoMessage()               {}
func (*UseRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *UseRequest) GetFacet() string {
	if m != nil {
		return m.Facet
	}
	return ""
}

func (m *UseRequest) GetQuantity() int64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

type UseResponse struct {
	Facet     string `protobuf:"bytes,1,opt,name=facet" json:"facet,omitempty"`
	Quantity  int64  `protobuf:"varint,2,opt,name=quantity" json:"quantity,omitempty"`
	Remaining int64  `protobuf:"varint,3,opt,name=remaining" json:"remaining,omitempty"`
}

func (m *UseResponse) Reset()                    { *m = UseResponse{} }
func (m *UseResponse) String() string            { return proto.CompactTextString(m) }
func (*UseResponse) ProtoMessage()               {}
func (*UseResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *UseResponse) GetFacet() string {
	if m != nil {
		return m.Facet
	}
	return ""
}

func (m *UseResponse) GetQuantity() int64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func (m *UseResponse) GetRemaining() int64 {
	if m != nil {
		return m.Remaining
	}
	return 0
}

func init() {
	proto.RegisterType((*IHave)(nil), "model.IHave")
	proto.RegisterType((*Push)(nil), "model.Push")
	proto.RegisterType((*HostTraffic)(nil), "model.HostTraffic")
	proto.RegisterType((*Traffic)(nil), "model.Traffic")
	proto.RegisterType((*ConnectorHandshake)(nil), "model.ConnectorHandshake")
	proto.RegisterType((*ResponderHandshake)(nil), "model.ResponderHandshake")
	proto.RegisterType((*UseRequest)(nil), "model.UseRequest")
	proto.RegisterType((*UseResponse)(nil), "model.UseResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for GossipProtocol service

type GossipProtocolClient interface {
	Gossip(ctx context.Context, opts ...grpc.CallOption) (GossipProtocol_GossipClient, error)
	Sync(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
}

type gossipProtocolClient struct {
	cc *grpc.ClientConn
}

func NewGossipProtocolClient(cc *grpc.ClientConn) GossipProtocolClient {
	return &gossipProtocolClient{cc}
}

func (c *gossipProtocolClient) Gossip(ctx context.Context, opts ...grpc.CallOption) (GossipProtocol_GossipClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_GossipProtocol_serviceDesc.Streams[0], c.cc, "/model.GossipProtocol/Gossip", opts...)
	if err != nil {
		return nil, err
	}
	x := &gossipProtocolGossipClient{stream}
	return x, nil
}

type GossipProtocol_GossipClient interface {
	Send(*ConnectorHandshake) error
	Recv() (*ResponderHandshake, error)
	grpc.ClientStream
}

type gossipProtocolGossipClient struct {
	grpc.ClientStream
}

func (x *gossipProtocolGossipClient) Send(m *ConnectorHandshake) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gossipProtocolGossipClient) Recv() (*ResponderHandshake, error) {
	m := new(ResponderHandshake)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gossipProtocolClient) Sync(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/model.GossipProtocol/Sync", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GossipProtocol service

type GossipProtocolServer interface {
	Gossip(GossipProtocol_GossipServer) error
	Sync(context.Context, *google_protobuf.Empty) (*google_protobuf.Empty, error)
}

func RegisterGossipProtocolServer(s *grpc.Server, srv GossipProtocolServer) {
	s.RegisterService(&_GossipProtocol_serviceDesc, srv)
}

func _GossipProtocol_Gossip_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GossipProtocolServer).Gossip(&gossipProtocolGossipServer{stream})
}

type GossipProtocol_GossipServer interface {
	Send(*ResponderHandshake) error
	Recv() (*ConnectorHandshake, error)
	grpc.ServerStream
}

type gossipProtocolGossipServer struct {
	grpc.ServerStream
}

func (x *gossipProtocolGossipServer) Send(m *ResponderHandshake) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gossipProtocolGossipServer) Recv() (*ConnectorHandshake, error) {
	m := new(ConnectorHandshake)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GossipProtocol_Sync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GossipProtocolServer).Sync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/model.GossipProtocol/Sync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GossipProtocolServer).Sync(ctx, req.(*google_protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _GossipProtocol_serviceDesc = grpc.ServiceDesc{
	ServiceName: "model.GossipProtocol",
	HandlerType: (*GossipProtocolServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Sync",
			Handler:    _GossipProtocol_Sync_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Gossip",
			Handler:       _GossipProtocol_Gossip_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "limits.proto",
}

// Client API for LimiterProtocol service

type LimiterProtocolClient interface {
	Use(ctx context.Context, in *UseRequest, opts ...grpc.CallOption) (*UseResponse, error)
}

type limiterProtocolClient struct {
	cc *grpc.ClientConn
}

func NewLimiterProtocolClient(cc *grpc.ClientConn) LimiterProtocolClient {
	return &limiterProtocolClient{cc}
}

func (c *limiterProtocolClient) Use(ctx context.Context, in *UseRequest, opts ...grpc.CallOption) (*UseResponse, error) {
	out := new(UseResponse)
	err := grpc.Invoke(ctx, "/model.LimiterProtocol/Use", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for LimiterProtocol service

type LimiterProtocolServer interface {
	Use(context.Context, *UseRequest) (*UseResponse, error)
}

func RegisterLimiterProtocolServer(s *grpc.Server, srv LimiterProtocolServer) {
	s.RegisterService(&_LimiterProtocol_serviceDesc, srv)
}

func _LimiterProtocol_Use_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LimiterProtocolServer).Use(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/model.LimiterProtocol/Use",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LimiterProtocolServer).Use(ctx, req.(*UseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _LimiterProtocol_serviceDesc = grpc.ServiceDesc{
	ServiceName: "model.LimiterProtocol",
	HandlerType: (*LimiterProtocolServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Use",
			Handler:    _LimiterProtocol_Use_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "limits.proto",
}

func init() { proto.RegisterFile("limits.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 443 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0x41, 0x6f, 0xd3, 0x40,
	0x10, 0x85, 0xeb, 0xc4, 0x6e, 0xc9, 0xb8, 0x2a, 0x62, 0x14, 0x21, 0x63, 0x90, 0x08, 0x16, 0x07,
	0x1f, 0x2a, 0x17, 0xb9, 0xc0, 0x11, 0xa4, 0x22, 0x84, 0x51, 0x39, 0x54, 0x0b, 0x15, 0x27, 0x0e,
	0x5b, 0x67, 0x12, 0xaf, 0x88, 0x77, 0x5d, 0xef, 0x1a, 0x29, 0xff, 0x81, 0x1f, 0x8d, 0xbc, 0xeb,
	0x34, 0x91, 0xa2, 0x1c, 0xb8, 0xe5, 0xbd, 0x99, 0x7c, 0xb3, 0xf3, 0x34, 0x86, 0xd3, 0x95, 0xa8,
	0x85, 0xd1, 0x59, 0xd3, 0x2a, 0xa3, 0x30, 0xa8, 0xd5, 0x9c, 0x56, 0xf1, 0xf3, 0xa5, 0x52, 0xcb,
	0x15, 0x5d, 0x58, 0xf3, 0xae, 0x5b, 0x5c, 0x50, 0xdd, 0x98, 0xb5, 0xeb, 0x49, 0x2e, 0x21, 0xf8,
	0x5a, 0xf0, 0x3f, 0x84, 0x53, 0x08, 0xa8, 0x51, 0x65, 0x15, 0x79, 0x33, 0x2f, 0x1d, 0x33, 0x27,
	0x7a, 0xb7, 0x52, 0xda, 0xe8, 0x68, 0x34, 0x1b, 0xa7, 0x13, 0xe6, 0x44, 0xf2, 0x16, 0xfc, 0x9b,
	0x4e, 0x57, 0x78, 0x0e, 0x27, 0xa6, 0xe5, 0x8b, 0x85, 0x28, 0x23, 0x6f, 0x36, 0x4e, 0xc3, 0x1c,
	0x33, 0x3b, 0x32, 0x2b, 0x94, 0x36, 0x3f, 0x5c, 0x85, 0x6d, 0x5a, 0x92, 0x6b, 0x08, 0x77, 0x7c,
	0x44, 0xf0, 0x25, 0xaf, 0xc9, 0xce, 0x9b, 0x30, 0xfb, 0x1b, 0xd3, 0x2d, 0x70, 0x64, 0x81, 0x67,
	0x03, 0x70, 0x0f, 0xf6, 0x0e, 0x4e, 0x36, 0xa0, 0x29, 0x04, 0x0b, 0x5e, 0x92, 0x19, 0x48, 0x4e,
	0xf4, 0x6e, 0xa7, 0xf9, 0x92, 0xa2, 0x91, 0xdb, 0xc7, 0x8a, 0xa4, 0x04, 0xfc, 0xa4, 0xa4, 0xa4,
	0xd2, 0xa8, 0xb6, 0xe0, 0x72, 0xae, 0x2b, 0xfe, 0x9b, 0xf0, 0x35, 0x04, 0xa2, 0x0f, 0xc1, 0x12,
	0xc2, 0xfc, 0x74, 0x18, 0x6a, 0x83, 0x29, 0x8e, 0x98, 0x2b, 0xe2, 0x2b, 0xf0, 0x9b, 0x4e, 0x57,
	0x16, 0x18, 0xe6, 0xe1, 0xd0, 0xd4, 0x07, 0x51, 0x1c, 0x31, 0x5b, 0xba, 0xf2, 0x61, 0x54, 0xe9,
	0xe4, 0x1a, 0x90, 0x91, 0x6e, 0x94, 0x9c, 0xd3, 0xce, 0x90, 0x29, 0x04, 0xe2, 0x27, 0x97, 0xc6,
	0x46, 0x35, 0x61, 0x4e, 0xe0, 0xcb, 0x83, 0x50, 0x87, 0x4c, 0x3e, 0x00, 0xdc, 0x6a, 0x62, 0x74,
	0xdf, 0x91, 0x36, 0x07, 0x76, 0x8d, 0xe1, 0xd1, 0x7d, 0xc7, 0xa5, 0x11, 0x66, 0x3d, 0xac, 0xfb,
	0xa0, 0x93, 0x5f, 0x10, 0xda, 0xff, 0xf7, 0xef, 0xd1, 0xf4, 0xff, 0x00, 0x7c, 0x01, 0x93, 0x96,
	0x6a, 0x2e, 0xa4, 0x90, 0xcb, 0x68, 0x6c, 0x8b, 0x5b, 0x23, 0xff, 0xeb, 0xc1, 0xd9, 0x17, 0xa5,
	0xb5, 0x68, 0x6e, 0xfa, 0x7b, 0x2a, 0xd5, 0x0a, 0xaf, 0xe0, 0xd8, 0x39, 0xf8, 0x6c, 0x58, 0x67,
	0x3f, 0xf2, 0x78, 0x53, 0xda, 0x0f, 0x2a, 0xf5, 0xde, 0x78, 0xf8, 0x1e, 0xfc, 0xef, 0x6b, 0x59,
	0xe2, 0xd3, 0xcc, 0x1d, 0x6f, 0xb6, 0x39, 0xde, 0xec, 0x73, 0x7f, 0xbc, 0xf1, 0x01, 0x3f, 0xff,
	0x08, 0x8f, 0xbf, 0xf5, 0x9f, 0x00, 0xb5, 0x0f, 0xcf, 0x39, 0x87, 0xf1, 0xad, 0x26, 0x7c, 0x32,
	0x0c, 0xdc, 0x86, 0x19, 0xe3, 0xae, 0xe5, 0xf2, 0xb9, 0x3b, 0xb6, 0xc0, 0xcb, 0x7f, 0x01, 0x00,
	0x00, 0xff, 0xff, 0x15, 0xd4, 0x30, 0x25, 0x4a, 0x03, 0x00, 0x00,
}
