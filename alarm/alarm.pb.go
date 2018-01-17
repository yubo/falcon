// Code generated by protoc-gen-go. DO NOT EDIT.
// source: alarm/alarm.proto

/*
Package alarm is a generated protocol buffer package.

It is generated from these files:
	alarm/alarm.proto

It has these top-level messages:
	Event
	PutRequest
	PutResponse
	Empty
	Stats
	StatsName
*/
package alarm

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

type Event struct {
	TagId     int64   `protobuf:"varint,1,opt,name=tagId" json:"tagId,omitempty"`
	Key       []byte  `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Expr      []byte  `protobuf:"bytes,3,opt,name=expr,proto3" json:"expr,omitempty"`
	Msg       []byte  `protobuf:"bytes,4,opt,name=msg,proto3" json:"msg,omitempty"`
	Timestamp int64   `protobuf:"varint,5,opt,name=timestamp" json:"timestamp,omitempty"`
	Value     float64 `protobuf:"fixed64,6,opt,name=value" json:"value,omitempty"`
	Priority  int32   `protobuf:"varint,7,opt,name=Priority" json:"Priority,omitempty"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Event) GetTagId() int64 {
	if m != nil {
		return m.TagId
	}
	return 0
}

func (m *Event) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Event) GetExpr() []byte {
	if m != nil {
		return m.Expr
	}
	return nil
}

func (m *Event) GetMsg() []byte {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (m *Event) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Event) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *Event) GetPriority() int32 {
	if m != nil {
		return m.Priority
	}
	return 0
}

type PutRequest struct {
	Events []*Event `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
}

func (m *PutRequest) Reset()                    { *m = PutRequest{} }
func (m *PutRequest) String() string            { return proto.CompactTextString(m) }
func (*PutRequest) ProtoMessage()               {}
func (*PutRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *PutRequest) GetEvents() []*Event {
	if m != nil {
		return m.Events
	}
	return nil
}

type PutResponse struct {
	N int32 `protobuf:"varint,1,opt,name=n" json:"n,omitempty"`
}

func (m *PutResponse) Reset()                    { *m = PutResponse{} }
func (m *PutResponse) String() string            { return proto.CompactTextString(m) }
func (*PutResponse) ProtoMessage()               {}
func (*PutResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *PutResponse) GetN() int32 {
	if m != nil {
		return m.N
	}
	return 0
}

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type Stats struct {
	Counter []uint64 `protobuf:"varint,1,rep,packed,name=counter" json:"counter,omitempty"`
}

func (m *Stats) Reset()                    { *m = Stats{} }
func (m *Stats) String() string            { return proto.CompactTextString(m) }
func (*Stats) ProtoMessage()               {}
func (*Stats) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

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
func (*StatsName) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *StatsName) GetCounterName() [][]byte {
	if m != nil {
		return m.CounterName
	}
	return nil
}

func init() {
	proto.RegisterType((*Event)(nil), "alarm.Event")
	proto.RegisterType((*PutRequest)(nil), "alarm.PutRequest")
	proto.RegisterType((*PutResponse)(nil), "alarm.PutResponse")
	proto.RegisterType((*Empty)(nil), "alarm.Empty")
	proto.RegisterType((*Stats)(nil), "alarm.Stats")
	proto.RegisterType((*StatsName)(nil), "alarm.StatsName")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Alarm service

type AlarmClient interface {
	Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error)
	GetStats(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Stats, error)
	GetStatsName(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*StatsName, error)
}

type alarmClient struct {
	cc *grpc.ClientConn
}

func NewAlarmClient(cc *grpc.ClientConn) AlarmClient {
	return &alarmClient{cc}
}

func (c *alarmClient) Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error) {
	out := new(PutResponse)
	err := grpc.Invoke(ctx, "/alarm.Alarm/Put", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alarmClient) GetStats(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Stats, error) {
	out := new(Stats)
	err := grpc.Invoke(ctx, "/alarm.Alarm/GetStats", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alarmClient) GetStatsName(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*StatsName, error) {
	out := new(StatsName)
	err := grpc.Invoke(ctx, "/alarm.Alarm/GetStatsName", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Alarm service

type AlarmServer interface {
	Put(context.Context, *PutRequest) (*PutResponse, error)
	GetStats(context.Context, *Empty) (*Stats, error)
	GetStatsName(context.Context, *Empty) (*StatsName, error)
}

func RegisterAlarmServer(s *grpc.Server, srv AlarmServer) {
	s.RegisterService(&_Alarm_serviceDesc, srv)
}

func _Alarm_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlarmServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alarm.Alarm/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlarmServer).Put(ctx, req.(*PutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Alarm_GetStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlarmServer).GetStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alarm.Alarm/GetStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlarmServer).GetStats(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Alarm_GetStatsName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlarmServer).GetStatsName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alarm.Alarm/GetStatsName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlarmServer).GetStatsName(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Alarm_serviceDesc = grpc.ServiceDesc{
	ServiceName: "alarm.Alarm",
	HandlerType: (*AlarmServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Put",
			Handler:    _Alarm_Put_Handler,
		},
		{
			MethodName: "GetStats",
			Handler:    _Alarm_GetStats_Handler,
		},
		{
			MethodName: "GetStatsName",
			Handler:    _Alarm_GetStatsName_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "alarm/alarm.proto",
}

func init() { proto.RegisterFile("alarm/alarm.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 419 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x54, 0x52, 0xc1, 0x8e, 0xd3, 0x30,
	0x14, 0x94, 0x37, 0x4d, 0xbb, 0x7d, 0x0d, 0xa8, 0xfb, 0xd8, 0x83, 0x55, 0x2a, 0x14, 0x2c, 0x0e,
	0x11, 0x12, 0x0d, 0x94, 0x1b, 0x9c, 0x58, 0x69, 0x05, 0x5c, 0x50, 0x65, 0x6e, 0x5c, 0x90, 0x5b,
	0x4c, 0x88, 0xa8, 0xed, 0x10, 0xbf, 0x54, 0xf4, 0xca, 0x2f, 0xf0, 0x0d, 0xfc, 0x10, 0xfc, 0x02,
	0x1f, 0x82, 0x62, 0xb7, 0xbb, 0xe5, 0x12, 0x79, 0xc6, 0x33, 0x93, 0x19, 0xc9, 0x70, 0xa1, 0xb6,
	0xaa, 0x35, 0x65, 0xf8, 0x2e, 0x9a, 0xd6, 0x91, 0xc3, 0x34, 0x80, 0xd9, 0xbc, 0x72, 0xae, 0xda,
	0xea, 0x52, 0x35, 0x75, 0xa9, 0xac, 0x75, 0xa4, 0xa8, 0x76, 0xd6, 0x47, 0x91, 0xf8, 0xc5, 0x20,
	0xbd, 0xde, 0x69, 0x4b, 0x78, 0x09, 0x29, 0xa9, 0xea, 0xed, 0x27, 0xce, 0x72, 0x56, 0x24, 0x32,
	0x02, 0x9c, 0x42, 0xf2, 0x55, 0xef, 0xf9, 0x59, 0xce, 0x8a, 0x4c, 0xf6, 0x47, 0x44, 0x18, 0xe8,
	0xef, 0x4d, 0xcb, 0x93, 0x40, 0x85, 0x73, 0xaf, 0x32, 0xbe, 0xe2, 0x83, 0xa8, 0x32, 0xbe, 0xc2,
	0x39, 0x8c, 0xa9, 0x36, 0xda, 0x93, 0x32, 0x0d, 0x4f, 0x43, 0xe2, 0x2d, 0xd1, 0xff, 0x6b, 0xa7,
	0xb6, 0x9d, 0xe6, 0xc3, 0x9c, 0x15, 0x4c, 0x46, 0x80, 0x33, 0x38, 0x5f, 0xb5, 0xb5, 0x6b, 0x6b,
	0xda, 0xf3, 0x51, 0xce, 0x8a, 0x54, 0xde, 0x60, 0xb1, 0x04, 0x58, 0x75, 0x24, 0xf5, 0xb7, 0x4e,
	0x7b, 0xc2, 0x47, 0x30, 0xd4, 0x7d, 0x69, 0xcf, 0x59, 0x9e, 0x14, 0x93, 0x65, 0xb6, 0x88, 0xc3,
	0xc3, 0x12, 0x79, 0xb8, 0x13, 0xf7, 0x61, 0x12, 0x3c, 0xbe, 0x71, 0xd6, 0x6b, 0xcc, 0x80, 0xd9,
	0x30, 0x2e, 0x95, 0xcc, 0x8a, 0x11, 0xa4, 0xd7, 0xa6, 0xa1, 0xbd, 0x78, 0x08, 0xe9, 0x7b, 0x52,
	0xe4, 0x91, 0xc3, 0x68, 0xe3, 0x3a, 0x4b, 0xba, 0x0d, 0xa9, 0x03, 0x79, 0x84, 0xe2, 0x09, 0x8c,
	0x83, 0xe4, 0x9d, 0x32, 0x1a, 0x73, 0x98, 0x1c, 0xf8, 0x1e, 0x06, 0x69, 0x26, 0x4f, 0xa9, 0xe5,
	0x6f, 0x06, 0xe9, 0xab, 0xbe, 0x0f, 0x5e, 0x41, 0xb2, 0xea, 0x08, 0x2f, 0x0e, 0xf5, 0x6e, 0x17,
	0xcc, 0xf0, 0x94, 0x8a, 0x05, 0xc5, 0xe5, 0x8f, 0x3f, 0x7f, 0x7f, 0x9e, 0xdd, 0x15, 0xe3, 0x72,
	0xf7, 0x6c, 0xf1, 0xb4, 0x6c, 0x3a, 0x7a, 0xc1, 0x1e, 0xe3, 0x4b, 0x38, 0x7f, 0xad, 0x29, 0x56,
	0xbc, 0xd9, 0xd9, 0x37, 0x9f, 0x1d, 0x51, 0xb8, 0x13, 0xf7, 0x82, 0xfb, 0x0e, 0x4e, 0xa2, 0xdb,
	0x07, 0xc3, 0x1b, 0xc8, 0x8e, 0xe6, 0x50, 0xfe, 0xff, 0x80, 0xe9, 0x69, 0x40, 0x7f, 0x2f, 0x78,
	0x08, 0x41, 0x9c, 0x9e, 0x84, 0x7c, 0xb4, 0xca, 0xe8, 0xab, 0x07, 0x1f, 0xe6, 0x55, 0x4d, 0x5f,
	0xba, 0xf5, 0x62, 0xe3, 0x4c, 0xb9, 0xef, 0xd6, 0xae, 0xfc, 0xac, 0xb6, 0x1b, 0x67, 0xe3, 0x9b,
	0x5b, 0x0f, 0xc3, 0x7b, 0x7a, 0xfe, 0x2f, 0x00, 0x00, 0xff, 0xff, 0xf9, 0x27, 0xe4, 0xf3, 0x89,
	0x02, 0x00, 0x00,
}
