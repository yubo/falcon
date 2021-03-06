// Code generated by protoc-gen-go. DO NOT EDIT.
// source: modules/service/service.proto

/*
Package service is a generated protocol buffer package.

It is generated from these files:
	modules/service/service.proto

It has these top-level messages:
	PutRequest
	PutResponse
	GetRequest
	GetResponse
	Empty
	Stats
	StatsName
*/
package service

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import tsdb "github.com/yubo/falcon/lib/tsdb"

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

type Cf int32

const (
	Cf_MIN     Cf = 0
	Cf_MAX     Cf = 1
	Cf_AVERAGE Cf = 2
)

var Cf_name = map[int32]string{
	0: "MIN",
	1: "MAX",
	2: "AVERAGE",
}
var Cf_value = map[string]int32{
	"MIN":     0,
	"MAX":     1,
	"AVERAGE": 2,
}

func (x Cf) String() string {
	return proto.EnumName(Cf_name, int32(x))
}
func (Cf) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ItemType int32

const (
	ItemType_GAUGE   ItemType = 0
	ItemType_DERIVE  ItemType = 1
	ItemType_COUNTER ItemType = 2
)

var ItemType_name = map[int32]string{
	0: "GAUGE",
	1: "DERIVE",
	2: "COUNTER",
}
var ItemType_value = map[string]int32{
	"GAUGE":   0,
	"DERIVE":  1,
	"COUNTER": 2,
}

func (x ItemType) String() string {
	return proto.EnumName(ItemType_name, int32(x))
}
func (ItemType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type PutRequest struct {
	Data []*tsdb.DataPoint `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
}

func (m *PutRequest) Reset()                    { *m = PutRequest{} }
func (m *PutRequest) String() string            { return proto.CompactTextString(m) }
func (*PutRequest) ProtoMessage()               {}
func (*PutRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *PutRequest) GetData() []*tsdb.DataPoint {
	if m != nil {
		return m.Data
	}
	return nil
}

type PutResponse struct {
	N int32 `protobuf:"varint,1,opt,name=n" json:"n,omitempty"`
}

func (m *PutResponse) Reset()                    { *m = PutResponse{} }
func (m *PutResponse) String() string            { return proto.CompactTextString(m) }
func (*PutResponse) ProtoMessage()               {}
func (*PutResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *PutResponse) GetN() int32 {
	if m != nil {
		return m.N
	}
	return 0
}

type GetRequest struct {
	Keys      []*tsdb.Key `protobuf:"bytes,1,rep,name=keys" json:"keys,omitempty"`
	Start     int64       `protobuf:"varint,2,opt,name=start" json:"start,omitempty"`
	End       int64       `protobuf:"varint,3,opt,name=end" json:"end,omitempty"`
	ConsolFun Cf          `protobuf:"varint,4,opt,name=consolFun,enum=service.Cf" json:"consolFun,omitempty"`
}

func (m *GetRequest) Reset()                    { *m = GetRequest{} }
func (m *GetRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()               {}
func (*GetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GetRequest) GetKeys() []*tsdb.Key {
	if m != nil {
		return m.Keys
	}
	return nil
}

func (m *GetRequest) GetStart() int64 {
	if m != nil {
		return m.Start
	}
	return 0
}

func (m *GetRequest) GetEnd() int64 {
	if m != nil {
		return m.End
	}
	return 0
}

func (m *GetRequest) GetConsolFun() Cf {
	if m != nil {
		return m.ConsolFun
	}
	return Cf_MIN
}

type GetResponse struct {
	Data []*tsdb.DataPoints `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
}

func (m *GetResponse) Reset()                    { *m = GetResponse{} }
func (m *GetResponse) String() string            { return proto.CompactTextString(m) }
func (*GetResponse) ProtoMessage()               {}
func (*GetResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GetResponse) GetData() []*tsdb.DataPoints {
	if m != nil {
		return m.Data
	}
	return nil
}

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type Stats struct {
	Counter []uint64 `protobuf:"varint,1,rep,packed,name=counter" json:"counter,omitempty"`
}

func (m *Stats) Reset()                    { *m = Stats{} }
func (m *Stats) String() string            { return proto.CompactTextString(m) }
func (*Stats) ProtoMessage()               {}
func (*Stats) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Stats) GetCounter() []uint64 {
	if m != nil {
		return m.Counter
	}
	return nil
}

type StatsName struct {
	CounterName [][]byte `protobuf:"bytes,1,rep,name=counterName,proto3" json:"counterName,omitempty"`
}

func (m *StatsName) Reset()                    { *m = StatsName{} }
func (m *StatsName) String() string            { return proto.CompactTextString(m) }
func (*StatsName) ProtoMessage()               {}
func (*StatsName) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *StatsName) GetCounterName() [][]byte {
	if m != nil {
		return m.CounterName
	}
	return nil
}

func init() {
	proto.RegisterType((*PutRequest)(nil), "service.PutRequest")
	proto.RegisterType((*PutResponse)(nil), "service.PutResponse")
	proto.RegisterType((*GetRequest)(nil), "service.GetRequest")
	proto.RegisterType((*GetResponse)(nil), "service.GetResponse")
	proto.RegisterType((*Empty)(nil), "service.Empty")
	proto.RegisterType((*Stats)(nil), "service.Stats")
	proto.RegisterType((*StatsName)(nil), "service.StatsName")
	proto.RegisterEnum("service.Cf", Cf_name, Cf_value)
	proto.RegisterEnum("service.ItemType", ItemType_name, ItemType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Service service

type ServiceClient interface {
	Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	GetStats(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Stats, error)
	GetStatsName(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*StatsName, error)
}

type serviceClient struct {
	cc *grpc.ClientConn
}

func NewServiceClient(cc *grpc.ClientConn) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error) {
	out := new(PutResponse)
	err := grpc.Invoke(ctx, "/service.Service/Put", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := grpc.Invoke(ctx, "/service.Service/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetStats(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Stats, error) {
	out := new(Stats)
	err := grpc.Invoke(ctx, "/service.Service/GetStats", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetStatsName(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*StatsName, error) {
	out := new(StatsName)
	err := grpc.Invoke(ctx, "/service.Service/GetStatsName", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Service service

type ServiceServer interface {
	Put(context.Context, *PutRequest) (*PutResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	GetStats(context.Context, *Empty) (*Stats, error)
	GetStatsName(context.Context, *Empty) (*StatsName, error)
}

func RegisterServiceServer(s *grpc.Server, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Service/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Put(ctx, req.(*PutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Service/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Service/GetStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetStats(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetStatsName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetStatsName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Service/GetStatsName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetStatsName(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "service.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Put",
			Handler:    _Service_Put_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Service_Get_Handler,
		},
		{
			MethodName: "GetStats",
			Handler:    _Service_GetStats_Handler,
		},
		{
			MethodName: "GetStatsName",
			Handler:    _Service_GetStatsName_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "modules/service/service.proto",
}

func init() { proto.RegisterFile("modules/service/service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 507 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x53, 0xdd, 0x8e, 0xd2, 0x40,
	0x18, 0xa5, 0x2d, 0x6c, 0x97, 0xaf, 0xb8, 0x36, 0xc3, 0x5e, 0x34, 0xe8, 0x26, 0x38, 0x6b, 0x0c,
	0x6e, 0x62, 0xeb, 0xb2, 0x0f, 0x60, 0x90, 0xad, 0x0d, 0x51, 0x91, 0xcc, 0xfe, 0xc4, 0x78, 0x63,
	0x06, 0x18, 0xb0, 0x91, 0xce, 0x20, 0x33, 0xdd, 0xa4, 0x37, 0xc6, 0xf8, 0x0a, 0x3e, 0x9a, 0xaf,
	0xe0, 0x83, 0x18, 0xa6, 0x05, 0x0a, 0x37, 0xed, 0x9c, 0x33, 0xe7, 0x9c, 0xf6, 0x3b, 0xd3, 0xc2,
	0x59, 0x22, 0xa6, 0xe9, 0x82, 0xc9, 0x40, 0xb2, 0xd5, 0x43, 0x3c, 0x61, 0x9b, 0xbb, 0xbf, 0x5c,
	0x09, 0x25, 0x90, 0x5d, 0xc0, 0xd6, 0xd3, 0xb9, 0x10, 0xf3, 0x05, 0x0b, 0xe8, 0x32, 0x0e, 0x28,
	0xe7, 0x42, 0x51, 0x15, 0x0b, 0x2e, 0x73, 0x59, 0xab, 0xb9, 0x88, 0xc7, 0x81, 0x92, 0xd3, 0xfc,
	0x92, 0x93, 0xf8, 0x12, 0x60, 0x94, 0x2a, 0xc2, 0x7e, 0xa4, 0x4c, 0x2a, 0x74, 0x0e, 0xd5, 0x29,
	0x55, 0xd4, 0x33, 0xda, 0x56, 0xc7, 0xe9, 0x3e, 0xf6, 0xb5, 0xf0, 0x9a, 0x2a, 0x3a, 0x12, 0x31,
	0x57, 0x44, 0x6f, 0xe2, 0x27, 0xe0, 0x68, 0x8b, 0x5c, 0x0a, 0x2e, 0x19, 0x6a, 0x80, 0xc1, 0x3d,
	0xa3, 0x6d, 0x74, 0x6a, 0xc4, 0xe0, 0xf8, 0x27, 0x40, 0xc4, 0xb6, 0x79, 0x67, 0x50, 0xfd, 0xce,
	0x32, 0x59, 0xe4, 0xd5, 0xf3, 0xbc, 0xf7, 0x2c, 0x23, 0x9a, 0x46, 0xa7, 0x50, 0x93, 0x8a, 0xae,
	0x94, 0x67, 0xb6, 0x8d, 0x8e, 0x45, 0x72, 0x80, 0x5c, 0xb0, 0x18, 0x9f, 0x7a, 0x96, 0xe6, 0xd6,
	0x4b, 0xf4, 0x12, 0xea, 0x13, 0xc1, 0xa5, 0x58, 0xbc, 0x4b, 0xb9, 0x57, 0x6d, 0x1b, 0x9d, 0x93,
	0xae, 0xe3, 0x6f, 0x3a, 0xe8, 0xcf, 0xc8, 0x6e, 0x17, 0x5f, 0x81, 0xa3, 0x9f, 0x5f, 0xbc, 0xdc,
	0xf3, 0xbd, 0x81, 0xdc, 0x83, 0x81, 0x64, 0x31, 0x91, 0x0d, 0xb5, 0x30, 0x59, 0xaa, 0x0c, 0x3f,
	0x83, 0xda, 0x8d, 0xa2, 0x4a, 0x22, 0x0f, 0xec, 0x89, 0x48, 0xb9, 0x62, 0x2b, 0x6d, 0xad, 0x92,
	0x0d, 0xc4, 0xaf, 0xa0, 0xae, 0x25, 0x43, 0x9a, 0x30, 0xd4, 0x06, 0xa7, 0xe0, 0xd7, 0x50, 0x4b,
	0x1b, 0xa4, 0x4c, 0x5d, 0x9c, 0x83, 0xd9, 0x9f, 0x21, 0x1b, 0xac, 0x8f, 0x83, 0xa1, 0x5b, 0xd1,
	0x8b, 0xde, 0x67, 0xd7, 0x40, 0x0e, 0xd8, 0xbd, 0xfb, 0x90, 0xf4, 0xa2, 0xd0, 0x35, 0x2f, 0x7c,
	0x38, 0x1e, 0x28, 0x96, 0xdc, 0x66, 0x4b, 0x86, 0xea, 0x50, 0x8b, 0x7a, 0x77, 0x51, 0xe8, 0x56,
	0x10, 0xc0, 0xd1, 0x75, 0x48, 0x06, 0xf7, 0x61, 0xae, 0xef, 0x7f, 0xba, 0x1b, 0xde, 0x86, 0xc4,
	0x35, 0xbb, 0xbf, 0x4c, 0xb0, 0x6f, 0xf2, 0xf1, 0x51, 0x17, 0xac, 0x51, 0xaa, 0x50, 0x73, 0xdb,
	0xc7, 0xee, 0x38, 0x5b, 0xa7, 0xfb, 0x64, 0xde, 0x09, 0xae, 0xac, 0x3d, 0x11, 0x2b, 0x7b, 0x76,
	0x47, 0x56, 0xf2, 0x94, 0x7a, 0xc4, 0x15, 0xf4, 0x06, 0x8e, 0x23, 0xa6, 0xf2, 0x76, 0x4e, 0xb6,
	0x1a, 0x5d, 0x5b, 0x6b, 0x87, 0xf5, 0x3e, 0x6e, 0xfe, 0xfe, 0xfb, 0xef, 0x8f, 0xf9, 0x08, 0x39,
	0xc1, 0xc3, 0xa5, 0xff, 0x3a, 0x90, 0xda, 0xf4, 0x01, 0x1a, 0x9b, 0x00, 0xdd, 0xdd, 0x61, 0x08,
	0xda, 0x0f, 0x59, 0x6b, 0xb0, 0xa7, 0x83, 0x10, 0x72, 0x4b, 0x41, 0x5f, 0x39, 0x4d, 0xd8, 0xdb,
	0xce, 0x97, 0x17, 0xf3, 0x58, 0x7d, 0x4b, 0xc7, 0xfe, 0x44, 0x24, 0x41, 0x96, 0x8e, 0x45, 0x30,
	0xa3, 0x8b, 0x89, 0xe0, 0xc1, 0xc1, 0xbf, 0x32, 0x3e, 0xd2, 0x1f, 0xfa, 0xd5, 0xff, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x4c, 0xc6, 0x6e, 0x4d, 0x45, 0x03, 0x00, 0x00,
}
