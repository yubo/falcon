// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service/service.proto

/*
Package service is a generated protocol buffer package.

It is generated from these files:
	service/service.proto

It has these top-level messages:
	Item
	DataPoint
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

type Item struct {
	Key       []byte  `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	ShardId   int32   `protobuf:"varint,2,opt,name=shardId" json:"shardId,omitempty"`
	Value     float64 `protobuf:"fixed64,3,opt,name=value" json:"value,omitempty"`
	Timestamp int64   `protobuf:"varint,4,opt,name=timestamp" json:"timestamp,omitempty"`
}

func (m *Item) Reset()                    { *m = Item{} }
func (m *Item) String() string            { return proto.CompactTextString(m) }
func (*Item) ProtoMessage()               {}
func (*Item) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Item) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Item) GetShardId() int32 {
	if m != nil {
		return m.ShardId
	}
	return 0
}

func (m *Item) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *Item) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type DataPoint struct {
	Timestamp int64   `protobuf:"varint,1,opt,name=timestamp" json:"timestamp,omitempty"`
	Value     float64 `protobuf:"fixed64,2,opt,name=value" json:"value,omitempty"`
}

func (m *DataPoint) Reset()                    { *m = DataPoint{} }
func (m *DataPoint) String() string            { return proto.CompactTextString(m) }
func (*DataPoint) ProtoMessage()               {}
func (*DataPoint) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DataPoint) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *DataPoint) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type PutRequest struct {
	Items []*Item `protobuf:"bytes,1,rep,name=items" json:"items,omitempty"`
}

func (m *PutRequest) Reset()                    { *m = PutRequest{} }
func (m *PutRequest) String() string            { return proto.CompactTextString(m) }
func (*PutRequest) ProtoMessage()               {}
func (*PutRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *PutRequest) GetItems() []*Item {
	if m != nil {
		return m.Items
	}
	return nil
}

type PutResponse struct {
	N int32 `protobuf:"varint,1,opt,name=n" json:"n,omitempty"`
}

func (m *PutResponse) Reset()                    { *m = PutResponse{} }
func (m *PutResponse) String() string            { return proto.CompactTextString(m) }
func (*PutResponse) ProtoMessage()               {}
func (*PutResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *PutResponse) GetN() int32 {
	if m != nil {
		return m.N
	}
	return 0
}

type GetRequest struct {
	Start     int64  `protobuf:"varint,1,opt,name=start" json:"start,omitempty"`
	End       int64  `protobuf:"varint,2,opt,name=end" json:"end,omitempty"`
	Key       []byte `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	ShardId   int32  `protobuf:"varint,4,opt,name=shardId" json:"shardId,omitempty"`
	ConsolFun Cf     `protobuf:"varint,5,opt,name=consolFun,enum=service.Cf" json:"consolFun,omitempty"`
}

func (m *GetRequest) Reset()                    { *m = GetRequest{} }
func (m *GetRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()               {}
func (*GetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

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

func (m *GetRequest) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *GetRequest) GetShardId() int32 {
	if m != nil {
		return m.ShardId
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
	Key []byte       `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Dps []*DataPoint `protobuf:"bytes,2,rep,name=dps" json:"dps,omitempty"`
}

func (m *GetResponse) Reset()                    { *m = GetResponse{} }
func (m *GetResponse) String() string            { return proto.CompactTextString(m) }
func (*GetResponse) ProtoMessage()               {}
func (*GetResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *GetResponse) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *GetResponse) GetDps() []*DataPoint {
	if m != nil {
		return m.Dps
	}
	return nil
}

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type Stats struct {
	Counter []uint64 `protobuf:"varint,1,rep,packed,name=counter" json:"counter,omitempty"`
}

func (m *Stats) Reset()                    { *m = Stats{} }
func (m *Stats) String() string            { return proto.CompactTextString(m) }
func (*Stats) ProtoMessage()               {}
func (*Stats) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

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
func (*StatsName) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *StatsName) GetCounterName() [][]byte {
	if m != nil {
		return m.CounterName
	}
	return nil
}

func init() {
	proto.RegisterType((*Item)(nil), "service.Item")
	proto.RegisterType((*DataPoint)(nil), "service.DataPoint")
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
	Metadata: "service/service.proto",
}

func init() { proto.RegisterFile("service/service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 583 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x93, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xbb, 0x76, 0xdc, 0x34, 0xe3, 0xb4, 0x32, 0xdb, 0x82, 0xac, 0x52, 0x21, 0xb3, 0xe5,
	0x10, 0x2a, 0x11, 0xd3, 0x72, 0xe3, 0x52, 0x85, 0xd6, 0x44, 0x91, 0x20, 0x54, 0xdb, 0x3f, 0x42,
	0x5c, 0xd0, 0x36, 0xdd, 0xa4, 0x16, 0xf1, 0xae, 0xc9, 0xae, 0x2b, 0xe5, 0xca, 0x95, 0x23, 0x2f,
	0xc3, 0x7b, 0xf0, 0x0a, 0x3c, 0x08, 0xf2, 0x3a, 0x8e, 0x9d, 0x02, 0x27, 0xef, 0xcc, 0x7c, 0xfb,
	0xf3, 0xcc, 0x37, 0x36, 0x3c, 0x54, 0x7c, 0x76, 0x17, 0x8f, 0x78, 0xb8, 0x78, 0x76, 0xd3, 0x99,
	0xd4, 0x12, 0x37, 0x17, 0xe1, 0xee, 0xde, 0x44, 0xca, 0xc9, 0x94, 0x87, 0x2c, 0x8d, 0x43, 0x26,
	0x84, 0xd4, 0x4c, 0xc7, 0x52, 0xa8, 0x42, 0x46, 0xc6, 0xd0, 0x18, 0x68, 0x9e, 0x60, 0x0f, 0xec,
	0x2f, 0x7c, 0xee, 0xa3, 0x00, 0x75, 0xda, 0x34, 0x3f, 0x62, 0x1f, 0x9a, 0xea, 0x96, 0xcd, 0x6e,
	0x06, 0x37, 0xbe, 0x15, 0xa0, 0x8e, 0x43, 0xcb, 0x10, 0xef, 0x80, 0x73, 0xc7, 0xa6, 0x19, 0xf7,
	0xed, 0x00, 0x75, 0x10, 0x2d, 0x02, 0xbc, 0x07, 0x2d, 0x1d, 0x27, 0x5c, 0x69, 0x96, 0xa4, 0x7e,
	0x23, 0x40, 0x1d, 0x9b, 0x56, 0x09, 0x72, 0x0c, 0xad, 0x53, 0xa6, 0xd9, 0x99, 0x8c, 0x85, 0x5e,
	0x95, 0xa2, 0x7b, 0xd2, 0x0a, 0x6f, 0xd5, 0xf0, 0xe4, 0x10, 0xe0, 0x2c, 0xd3, 0x94, 0x7f, 0xcd,
	0xb8, 0xd2, 0x78, 0x1f, 0x9c, 0x58, 0xf3, 0x44, 0xf9, 0x28, 0xb0, 0x3b, 0xee, 0xd1, 0x66, 0xb7,
	0x1c, 0x3e, 0x1f, 0x86, 0x16, 0x35, 0xf2, 0x18, 0x5c, 0x73, 0x45, 0xa5, 0x52, 0x28, 0x8e, 0xdb,
	0x80, 0x84, 0x79, 0x9b, 0x43, 0x91, 0x20, 0xdf, 0x11, 0x40, 0x9f, 0x2f, 0x81, 0x3b, 0xe0, 0x28,
	0xcd, 0x66, 0x7a, 0xd1, 0x4e, 0x11, 0xe4, 0xae, 0x70, 0x51, 0xcc, 0x6f, 0xd3, 0xfc, 0x58, 0xfa,
	0x64, 0xff, 0xd3, 0xa7, 0xc6, 0xaa, 0x4f, 0xcf, 0xa1, 0x35, 0x92, 0x42, 0xc9, 0xe9, 0xdb, 0x4c,
	0xf8, 0x4e, 0x80, 0x3a, 0x5b, 0x47, 0xee, 0xb2, 0xd1, 0x93, 0x31, 0xad, 0xaa, 0x24, 0x02, 0xd7,
	0x34, 0xb3, 0x68, 0xf5, 0xef, 0x6d, 0x3c, 0x03, 0xfb, 0x26, 0x55, 0xbe, 0x65, 0xc6, 0xc5, 0x4b,
	0xca, 0xd2, 0x53, 0x9a, 0x97, 0x49, 0x13, 0x9c, 0x28, 0x49, 0xf5, 0x9c, 0x3c, 0x05, 0xe7, 0x5c,
	0x33, 0xad, 0xf2, 0xee, 0x46, 0x32, 0x13, 0x9a, 0xcf, 0x8c, 0x55, 0x0d, 0x5a, 0x86, 0xe4, 0x05,
	0xb4, 0x8c, 0x64, 0xc8, 0x12, 0x8e, 0x03, 0x70, 0x17, 0xf9, 0x3c, 0x34, 0xd2, 0x36, 0xad, 0xa7,
	0x0e, 0xf6, 0xc1, 0x3a, 0x19, 0xe3, 0x26, 0xd8, 0xef, 0x07, 0x43, 0x6f, 0xcd, 0x1c, 0x7a, 0x1f,
	0x3d, 0x84, 0x5d, 0x68, 0xf6, 0xae, 0x22, 0xda, 0xeb, 0x47, 0x9e, 0x75, 0xd0, 0x85, 0x8d, 0x7c,
	0x01, 0x17, 0xf3, 0x94, 0xe3, 0x16, 0x38, 0xfd, 0xde, 0x65, 0x3f, 0xf2, 0xd6, 0x30, 0xc0, 0xfa,
	0x69, 0x44, 0x07, 0x57, 0x51, 0xa1, 0x3f, 0xf9, 0x70, 0x39, 0xbc, 0x88, 0xa8, 0x67, 0x1d, 0xfd,
	0xb4, 0xa0, 0x79, 0x5e, 0x8c, 0x82, 0x87, 0x60, 0x9f, 0x65, 0x1a, 0x6f, 0x2f, 0x67, 0xab, 0xd6,
	0xbd, 0xbb, 0xb3, 0x9a, 0x2c, 0x5c, 0x22, 0x7b, 0xdf, 0x7e, 0xfd, 0xfe, 0x61, 0x3d, 0x22, 0x0f,
	0xc2, 0xbb, 0xc3, 0xee, 0xcb, 0xf2, 0xfb, 0x0f, 0xd3, 0x4c, 0xbf, 0x46, 0x07, 0x39, 0xaf, 0xcf,
	0xeb, 0xbc, 0x6a, 0xdb, 0x35, 0x5e, 0xcd, 0xf5, 0xff, 0xf1, 0x26, 0xdc, 0xf0, 0x8e, 0x61, 0xa3,
	0xcf, 0x75, 0xe1, 0xea, 0xd6, 0xf2, 0xbe, 0xb1, 0x7b, 0xb7, 0x8a, 0x4d, 0x9d, 0x6c, 0x1b, 0xd2,
	0x26, 0x76, 0x17, 0x24, 0x73, 0xe9, 0x1d, 0xb4, 0x4b, 0x80, 0xf1, 0xfc, 0x3e, 0x04, 0xaf, 0x42,
	0x72, 0x0d, 0xf1, 0x0d, 0x08, 0x63, 0xaf, 0x06, 0xfa, 0x2c, 0x58, 0xc2, 0xdf, 0x04, 0x9f, 0x9e,
	0x4c, 0x62, 0x7d, 0x9b, 0x5d, 0x77, 0x47, 0x32, 0x09, 0xe7, 0xd9, 0xb5, 0x0c, 0xc7, 0x6c, 0x3a,
	0x92, 0xa2, 0xec, 0xfb, 0x7a, 0xdd, 0xfc, 0xe1, 0xaf, 0xfe, 0x04, 0x00, 0x00, 0xff, 0xff, 0x36,
	0x4d, 0xcf, 0x3f, 0x21, 0x04, 0x00, 0x00,
}
