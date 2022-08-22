// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: marketplace/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type MsgPublishCollection struct {
	Creator         string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	DenomId         string `protobuf:"bytes,2,opt,name=denomId,proto3" json:"denomId,omitempty"`
	MintRoyalties   string `protobuf:"bytes,3,opt,name=mintRoyalties,proto3" json:"mintRoyalties,omitempty"`
	ResaleRoyalties string `protobuf:"bytes,4,opt,name=resaleRoyalties,proto3" json:"resaleRoyalties,omitempty"`
}

func (m *MsgPublishCollection) Reset()         { *m = MsgPublishCollection{} }
func (m *MsgPublishCollection) String() string { return proto.CompactTextString(m) }
func (*MsgPublishCollection) ProtoMessage()    {}
func (*MsgPublishCollection) Descriptor() ([]byte, []int) {
	return fileDescriptor_689d664ba3f09b75, []int{0}
}
func (m *MsgPublishCollection) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgPublishCollection) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgPublishCollection.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgPublishCollection) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgPublishCollection.Merge(m, src)
}
func (m *MsgPublishCollection) XXX_Size() int {
	return m.Size()
}
func (m *MsgPublishCollection) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgPublishCollection.DiscardUnknown(m)
}

var xxx_messageInfo_MsgPublishCollection proto.InternalMessageInfo

func (m *MsgPublishCollection) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *MsgPublishCollection) GetDenomId() string {
	if m != nil {
		return m.DenomId
	}
	return ""
}

func (m *MsgPublishCollection) GetMintRoyalties() string {
	if m != nil {
		return m.MintRoyalties
	}
	return ""
}

func (m *MsgPublishCollection) GetResaleRoyalties() string {
	if m != nil {
		return m.ResaleRoyalties
	}
	return ""
}

type MsgPublishCollectionResponse struct {
}

func (m *MsgPublishCollectionResponse) Reset()         { *m = MsgPublishCollectionResponse{} }
func (m *MsgPublishCollectionResponse) String() string { return proto.CompactTextString(m) }
func (*MsgPublishCollectionResponse) ProtoMessage()    {}
func (*MsgPublishCollectionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_689d664ba3f09b75, []int{1}
}
func (m *MsgPublishCollectionResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgPublishCollectionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgPublishCollectionResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgPublishCollectionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgPublishCollectionResponse.Merge(m, src)
}
func (m *MsgPublishCollectionResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgPublishCollectionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgPublishCollectionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgPublishCollectionResponse proto.InternalMessageInfo

type MsgPublishNft struct {
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	TokenId string `protobuf:"bytes,2,opt,name=tokenId,proto3" json:"tokenId,omitempty"`
	DenomId string `protobuf:"bytes,3,opt,name=denomId,proto3" json:"denomId,omitempty"`
	Price   string `protobuf:"bytes,4,opt,name=price,proto3" json:"price,omitempty"`
}

func (m *MsgPublishNft) Reset()         { *m = MsgPublishNft{} }
func (m *MsgPublishNft) String() string { return proto.CompactTextString(m) }
func (*MsgPublishNft) ProtoMessage()    {}
func (*MsgPublishNft) Descriptor() ([]byte, []int) {
	return fileDescriptor_689d664ba3f09b75, []int{2}
}
func (m *MsgPublishNft) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgPublishNft) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgPublishNft.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgPublishNft) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgPublishNft.Merge(m, src)
}
func (m *MsgPublishNft) XXX_Size() int {
	return m.Size()
}
func (m *MsgPublishNft) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgPublishNft.DiscardUnknown(m)
}

var xxx_messageInfo_MsgPublishNft proto.InternalMessageInfo

func (m *MsgPublishNft) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *MsgPublishNft) GetTokenId() string {
	if m != nil {
		return m.TokenId
	}
	return ""
}

func (m *MsgPublishNft) GetDenomId() string {
	if m != nil {
		return m.DenomId
	}
	return ""
}

func (m *MsgPublishNft) GetPrice() string {
	if m != nil {
		return m.Price
	}
	return ""
}

type MsgPublishNftResponse struct {
}

func (m *MsgPublishNftResponse) Reset()         { *m = MsgPublishNftResponse{} }
func (m *MsgPublishNftResponse) String() string { return proto.CompactTextString(m) }
func (*MsgPublishNftResponse) ProtoMessage()    {}
func (*MsgPublishNftResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_689d664ba3f09b75, []int{3}
}
func (m *MsgPublishNftResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgPublishNftResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgPublishNftResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgPublishNftResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgPublishNftResponse.Merge(m, src)
}
func (m *MsgPublishNftResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgPublishNftResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgPublishNftResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgPublishNftResponse proto.InternalMessageInfo

type MsgBuyNft struct {
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Id      uint64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (m *MsgBuyNft) Reset()         { *m = MsgBuyNft{} }
func (m *MsgBuyNft) String() string { return proto.CompactTextString(m) }
func (*MsgBuyNft) ProtoMessage()    {}
func (*MsgBuyNft) Descriptor() ([]byte, []int) {
	return fileDescriptor_689d664ba3f09b75, []int{4}
}
func (m *MsgBuyNft) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgBuyNft) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgBuyNft.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgBuyNft) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgBuyNft.Merge(m, src)
}
func (m *MsgBuyNft) XXX_Size() int {
	return m.Size()
}
func (m *MsgBuyNft) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgBuyNft.DiscardUnknown(m)
}

var xxx_messageInfo_MsgBuyNft proto.InternalMessageInfo

func (m *MsgBuyNft) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *MsgBuyNft) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type MsgBuyNftResponse struct {
}

func (m *MsgBuyNftResponse) Reset()         { *m = MsgBuyNftResponse{} }
func (m *MsgBuyNftResponse) String() string { return proto.CompactTextString(m) }
func (*MsgBuyNftResponse) ProtoMessage()    {}
func (*MsgBuyNftResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_689d664ba3f09b75, []int{5}
}
func (m *MsgBuyNftResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgBuyNftResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgBuyNftResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgBuyNftResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgBuyNftResponse.Merge(m, src)
}
func (m *MsgBuyNftResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgBuyNftResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgBuyNftResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgBuyNftResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgPublishCollection)(nil), "cudoventures.cudosnode.marketplace.MsgPublishCollection")
	proto.RegisterType((*MsgPublishCollectionResponse)(nil), "cudoventures.cudosnode.marketplace.MsgPublishCollectionResponse")
	proto.RegisterType((*MsgPublishNft)(nil), "cudoventures.cudosnode.marketplace.MsgPublishNft")
	proto.RegisterType((*MsgPublishNftResponse)(nil), "cudoventures.cudosnode.marketplace.MsgPublishNftResponse")
	proto.RegisterType((*MsgBuyNft)(nil), "cudoventures.cudosnode.marketplace.MsgBuyNft")
	proto.RegisterType((*MsgBuyNftResponse)(nil), "cudoventures.cudosnode.marketplace.MsgBuyNftResponse")
}

func init() { proto.RegisterFile("marketplace/tx.proto", fileDescriptor_689d664ba3f09b75) }

var fileDescriptor_689d664ba3f09b75 = []byte{
	// 394 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0xcf, 0x6e, 0xda, 0x40,
	0x10, 0xc6, 0x31, 0xa6, 0x54, 0x8c, 0x44, 0x2b, 0x5c, 0xaa, 0x5a, 0xa8, 0xb2, 0x2a, 0xab, 0x07,
	0x2e, 0xd8, 0x6a, 0x2b, 0xaa, 0xf6, 0x56, 0xc1, 0xa9, 0x07, 0x10, 0xf2, 0x21, 0x87, 0xdc, 0x8c,
	0x3d, 0x98, 0x15, 0xf6, 0xae, 0xe5, 0x5d, 0x47, 0x90, 0x6b, 0x1e, 0x20, 0x39, 0xe7, 0x89, 0x72,
	0xe4, 0x98, 0x63, 0x04, 0x2f, 0x12, 0xe1, 0x3f, 0x80, 0x13, 0x44, 0x42, 0x8e, 0xdf, 0xce, 0x37,
	0xdf, 0xfc, 0x76, 0x56, 0x0b, 0xcd, 0xc0, 0x8e, 0x66, 0x28, 0x42, 0xdf, 0x76, 0xd0, 0x14, 0x73,
	0x23, 0x8c, 0x98, 0x60, 0x8a, 0xee, 0xc4, 0x2e, 0xbb, 0x40, 0x2a, 0xe2, 0x08, 0xb9, 0xb1, 0x11,
	0x9c, 0x32, 0x17, 0x8d, 0x3d, 0xb3, 0x7e, 0x2b, 0x41, 0x73, 0xc0, 0xbd, 0x51, 0x3c, 0xf6, 0x09,
	0x9f, 0xf6, 0x99, 0xef, 0xa3, 0x23, 0x08, 0xa3, 0x8a, 0x0a, 0xef, 0x9d, 0x08, 0x6d, 0xc1, 0x22,
	0x55, 0xfa, 0x26, 0xb5, 0x6b, 0x56, 0x2e, 0x37, 0x15, 0x17, 0x29, 0x0b, 0xfe, 0xbb, 0x6a, 0x39,
	0xad, 0x64, 0x52, 0xf9, 0x0e, 0xf5, 0x80, 0x50, 0x61, 0xb1, 0x85, 0xed, 0x0b, 0x82, 0x5c, 0x95,
	0x93, 0x7a, 0xf1, 0x50, 0x69, 0xc3, 0xc7, 0x08, 0xb9, 0xed, 0xe3, 0xce, 0x57, 0x49, 0x7c, 0x4f,
	0x8f, 0x75, 0x0d, 0xbe, 0x1e, 0x62, 0xb3, 0x90, 0x87, 0x8c, 0x72, 0xd4, 0x39, 0xd4, 0x77, 0xf5,
	0xe1, 0x44, 0x1c, 0x87, 0x16, 0x6c, 0x86, 0x74, 0x07, 0x9d, 0xc9, 0xfd, 0xeb, 0xc8, 0xc5, 0xeb,
	0x34, 0xe1, 0x5d, 0x18, 0x11, 0x07, 0x33, 0xbc, 0x54, 0xe8, 0x5f, 0xe0, 0x73, 0x61, 0xe8, 0x96,
	0xa6, 0x0b, 0xb5, 0x01, 0xf7, 0x7a, 0xf1, 0xe2, 0x38, 0xc9, 0x07, 0x28, 0x93, 0x14, 0xa2, 0x62,
	0x95, 0x89, 0xab, 0x7f, 0x82, 0xc6, 0xb6, 0x2d, 0xcf, 0xfa, 0x79, 0x25, 0x83, 0x3c, 0xe0, 0x9e,
	0x72, 0x2d, 0x41, 0xe3, 0xf9, 0xdb, 0xfc, 0x31, 0x5e, 0x7e, 0x59, 0xe3, 0xd0, 0xe6, 0x5a, 0xff,
	0xde, 0xda, 0x99, 0x93, 0x29, 0x97, 0x00, 0x7b, 0x0b, 0xff, 0x71, 0x5a, 0xde, 0x70, 0x22, 0x5a,
	0x7f, 0x4f, 0x6e, 0xd9, 0xce, 0xa6, 0x50, 0xcd, 0xd6, 0xdb, 0x79, 0x65, 0x48, 0x6a, 0x6f, 0x75,
	0x4f, 0xb2, 0xe7, 0xf3, 0x7a, 0xa3, 0xbb, 0x95, 0x26, 0x2d, 0x57, 0x9a, 0xf4, 0xb0, 0xd2, 0xa4,
	0x9b, 0xb5, 0x56, 0x5a, 0xae, 0xb5, 0xd2, 0xfd, 0x5a, 0x2b, 0x9d, 0xff, 0xf6, 0x88, 0x98, 0xc6,
	0x63, 0xc3, 0x61, 0x81, 0xd9, 0x8f, 0x5d, 0x76, 0x96, 0x45, 0x9b, 0x49, 0x74, 0x67, 0x93, 0x6d,
	0xce, 0xcd, 0xc2, 0xaf, 0x5c, 0x84, 0xc8, 0xc7, 0xd5, 0xe4, 0x67, 0xfe, 0x7a, 0x0c, 0x00, 0x00,
	0xff, 0xff, 0x20, 0x6f, 0x01, 0x0f, 0xb1, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	PublishCollection(ctx context.Context, in *MsgPublishCollection, opts ...grpc.CallOption) (*MsgPublishCollectionResponse, error)
	PublishNft(ctx context.Context, in *MsgPublishNft, opts ...grpc.CallOption) (*MsgPublishNftResponse, error)
	BuyNft(ctx context.Context, in *MsgBuyNft, opts ...grpc.CallOption) (*MsgBuyNftResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) PublishCollection(ctx context.Context, in *MsgPublishCollection, opts ...grpc.CallOption) (*MsgPublishCollectionResponse, error) {
	out := new(MsgPublishCollectionResponse)
	err := c.cc.Invoke(ctx, "/cudoventures.cudosnode.marketplace.Msg/PublishCollection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) PublishNft(ctx context.Context, in *MsgPublishNft, opts ...grpc.CallOption) (*MsgPublishNftResponse, error) {
	out := new(MsgPublishNftResponse)
	err := c.cc.Invoke(ctx, "/cudoventures.cudosnode.marketplace.Msg/PublishNft", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) BuyNft(ctx context.Context, in *MsgBuyNft, opts ...grpc.CallOption) (*MsgBuyNftResponse, error) {
	out := new(MsgBuyNftResponse)
	err := c.cc.Invoke(ctx, "/cudoventures.cudosnode.marketplace.Msg/BuyNft", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	PublishCollection(context.Context, *MsgPublishCollection) (*MsgPublishCollectionResponse, error)
	PublishNft(context.Context, *MsgPublishNft) (*MsgPublishNftResponse, error)
	BuyNft(context.Context, *MsgBuyNft) (*MsgBuyNftResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) PublishCollection(ctx context.Context, req *MsgPublishCollection) (*MsgPublishCollectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishCollection not implemented")
}
func (*UnimplementedMsgServer) PublishNft(ctx context.Context, req *MsgPublishNft) (*MsgPublishNftResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishNft not implemented")
}
func (*UnimplementedMsgServer) BuyNft(ctx context.Context, req *MsgBuyNft) (*MsgBuyNftResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuyNft not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_PublishCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgPublishCollection)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).PublishCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cudoventures.cudosnode.marketplace.Msg/PublishCollection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).PublishCollection(ctx, req.(*MsgPublishCollection))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_PublishNft_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgPublishNft)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).PublishNft(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cudoventures.cudosnode.marketplace.Msg/PublishNft",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).PublishNft(ctx, req.(*MsgPublishNft))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_BuyNft_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgBuyNft)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).BuyNft(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cudoventures.cudosnode.marketplace.Msg/BuyNft",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).BuyNft(ctx, req.(*MsgBuyNft))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cudoventures.cudosnode.marketplace.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PublishCollection",
			Handler:    _Msg_PublishCollection_Handler,
		},
		{
			MethodName: "PublishNft",
			Handler:    _Msg_PublishNft_Handler,
		},
		{
			MethodName: "BuyNft",
			Handler:    _Msg_BuyNft_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "marketplace/tx.proto",
}

func (m *MsgPublishCollection) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgPublishCollection) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgPublishCollection) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ResaleRoyalties) > 0 {
		i -= len(m.ResaleRoyalties)
		copy(dAtA[i:], m.ResaleRoyalties)
		i = encodeVarintTx(dAtA, i, uint64(len(m.ResaleRoyalties)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.MintRoyalties) > 0 {
		i -= len(m.MintRoyalties)
		copy(dAtA[i:], m.MintRoyalties)
		i = encodeVarintTx(dAtA, i, uint64(len(m.MintRoyalties)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.DenomId) > 0 {
		i -= len(m.DenomId)
		copy(dAtA[i:], m.DenomId)
		i = encodeVarintTx(dAtA, i, uint64(len(m.DenomId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgPublishCollectionResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgPublishCollectionResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgPublishCollectionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgPublishNft) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgPublishNft) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgPublishNft) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Price) > 0 {
		i -= len(m.Price)
		copy(dAtA[i:], m.Price)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Price)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.DenomId) > 0 {
		i -= len(m.DenomId)
		copy(dAtA[i:], m.DenomId)
		i = encodeVarintTx(dAtA, i, uint64(len(m.DenomId)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.TokenId) > 0 {
		i -= len(m.TokenId)
		copy(dAtA[i:], m.TokenId)
		i = encodeVarintTx(dAtA, i, uint64(len(m.TokenId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgPublishNftResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgPublishNftResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgPublishNftResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgBuyNft) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgBuyNft) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgBuyNft) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgBuyNftResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgBuyNftResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgBuyNftResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgPublishCollection) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.DenomId)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.MintRoyalties)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.ResaleRoyalties)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgPublishCollectionResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgPublishNft) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.TokenId)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.DenomId)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Price)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgPublishNftResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgBuyNft) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.Id != 0 {
		n += 1 + sovTx(uint64(m.Id))
	}
	return n
}

func (m *MsgBuyNftResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgPublishCollection) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgPublishCollection: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgPublishCollection: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DenomId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DenomId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MintRoyalties", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MintRoyalties = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ResaleRoyalties", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ResaleRoyalties = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgPublishCollectionResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgPublishCollectionResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgPublishCollectionResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgPublishNft) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgPublishNft: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgPublishNft: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DenomId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DenomId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Price = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgPublishNftResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgPublishNftResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgPublishNftResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgBuyNft) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgBuyNft: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgBuyNft: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgBuyNftResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgBuyNftResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgBuyNftResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
