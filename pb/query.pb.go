// Code generated by protoc-gen-go. DO NOT EDIT.
// source: query.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type QueryOptions_FilterType int32

const (
	QueryOptions_NO_FILTER  QueryOptions_FilterType = 0
	QueryOptions_HIDE_OLDER QueryOptions_FilterType = 1
)

var QueryOptions_FilterType_name = map[int32]string{
	0: "NO_FILTER",
	1: "HIDE_OLDER",
}

var QueryOptions_FilterType_value = map[string]int32{
	"NO_FILTER":  0,
	"HIDE_OLDER": 1,
}

func (x QueryOptions_FilterType) String() string {
	return proto.EnumName(QueryOptions_FilterType_name, int32(x))
}

func (QueryOptions_FilterType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{0, 0}
}

type Query_Type int32

const (
	Query_THREAD_SNAPSHOTS Query_Type = 0
	Query_CONTACTS         Query_Type = 1
	Query_VIDEO_CHUNKS     Query_Type = 50
	Query_VIDEO            Query_Type = 51
	Query_SYNC_FILE        Query_Type = 52
	Query_STREAM           Query_Type = 53
)

var Query_Type_name = map[int32]string{
	0:  "THREAD_SNAPSHOTS",
	1:  "CONTACTS",
	50: "VIDEO_CHUNKS",
	51: "VIDEO",
	52: "SYNC_FILE",
	53: "STREAM",
}

var Query_Type_value = map[string]int32{
	"THREAD_SNAPSHOTS": 0,
	"CONTACTS":         1,
	"VIDEO_CHUNKS":     50,
	"VIDEO":            51,
	"SYNC_FILE":        52,
	"STREAM":           53,
}

func (x Query_Type) String() string {
	return proto.EnumName(Query_Type_name, int32(x))
}

func (Query_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{1, 0}
}

type PubSubQuery_ResponseType int32

const (
	PubSubQuery_P2P    PubSubQuery_ResponseType = 0
	PubSubQuery_PUBSUB PubSubQuery_ResponseType = 1
)

var PubSubQuery_ResponseType_name = map[int32]string{
	0: "P2P",
	1: "PUBSUB",
}

var PubSubQuery_ResponseType_value = map[string]int32{
	"P2P":    0,
	"PUBSUB": 1,
}

func (x PubSubQuery_ResponseType) String() string {
	return proto.EnumName(PubSubQuery_ResponseType_name, int32(x))
}

func (PubSubQuery_ResponseType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{2, 0}
}

type QueryOptions struct {
	LocalOnly            bool                    `protobuf:"varint,1,opt,name=localOnly,proto3" json:"localOnly,omitempty"`
	RemoteOnly           bool                    `protobuf:"varint,6,opt,name=remoteOnly,proto3" json:"remoteOnly,omitempty"`
	Limit                int32                   `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Wait                 int32                   `protobuf:"varint,3,opt,name=wait,proto3" json:"wait,omitempty"`
	Filter               QueryOptions_FilterType `protobuf:"varint,4,opt,name=filter,proto3,enum=QueryOptions_FilterType" json:"filter,omitempty"`
	Exclude              []string                `protobuf:"bytes,5,rep,name=exclude,proto3" json:"exclude,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *QueryOptions) Reset()         { *m = QueryOptions{} }
func (m *QueryOptions) String() string { return proto.CompactTextString(m) }
func (*QueryOptions) ProtoMessage()    {}
func (*QueryOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{0}
}

func (m *QueryOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryOptions.Unmarshal(m, b)
}
func (m *QueryOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryOptions.Marshal(b, m, deterministic)
}
func (m *QueryOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryOptions.Merge(m, src)
}
func (m *QueryOptions) XXX_Size() int {
	return xxx_messageInfo_QueryOptions.Size(m)
}
func (m *QueryOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryOptions.DiscardUnknown(m)
}

var xxx_messageInfo_QueryOptions proto.InternalMessageInfo

func (m *QueryOptions) GetLocalOnly() bool {
	if m != nil {
		return m.LocalOnly
	}
	return false
}

func (m *QueryOptions) GetRemoteOnly() bool {
	if m != nil {
		return m.RemoteOnly
	}
	return false
}

func (m *QueryOptions) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *QueryOptions) GetWait() int32 {
	if m != nil {
		return m.Wait
	}
	return 0
}

func (m *QueryOptions) GetFilter() QueryOptions_FilterType {
	if m != nil {
		return m.Filter
	}
	return QueryOptions_NO_FILTER
}

func (m *QueryOptions) GetExclude() []string {
	if m != nil {
		return m.Exclude
	}
	return nil
}

type Query struct {
	Id                   string        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Token                string        `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	Type                 Query_Type    `protobuf:"varint,3,opt,name=type,proto3,enum=Query_Type" json:"type,omitempty"`
	Options              *QueryOptions `protobuf:"bytes,4,opt,name=options,proto3" json:"options,omitempty"`
	Payload              *any.Any      `protobuf:"bytes,5,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Query) Reset()         { *m = Query{} }
func (m *Query) String() string { return proto.CompactTextString(m) }
func (*Query) ProtoMessage()    {}
func (*Query) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{1}
}

func (m *Query) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Query.Unmarshal(m, b)
}
func (m *Query) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Query.Marshal(b, m, deterministic)
}
func (m *Query) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Query.Merge(m, src)
}
func (m *Query) XXX_Size() int {
	return xxx_messageInfo_Query.Size(m)
}
func (m *Query) XXX_DiscardUnknown() {
	xxx_messageInfo_Query.DiscardUnknown(m)
}

var xxx_messageInfo_Query proto.InternalMessageInfo

func (m *Query) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Query) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *Query) GetType() Query_Type {
	if m != nil {
		return m.Type
	}
	return Query_THREAD_SNAPSHOTS
}

func (m *Query) GetOptions() *QueryOptions {
	if m != nil {
		return m.Options
	}
	return nil
}

func (m *Query) GetPayload() *any.Any {
	if m != nil {
		return m.Payload
	}
	return nil
}

type PubSubQuery struct {
	Id                   string                   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type                 Query_Type               `protobuf:"varint,2,opt,name=type,proto3,enum=Query_Type" json:"type,omitempty"`
	Payload              *any.Any                 `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	ResponseType         PubSubQuery_ResponseType `protobuf:"varint,4,opt,name=responseType,proto3,enum=PubSubQuery_ResponseType" json:"responseType,omitempty"`
	Exclude              []string                 `protobuf:"bytes,5,rep,name=exclude,proto3" json:"exclude,omitempty"`
	Topic                string                   `protobuf:"bytes,6,opt,name=topic,proto3" json:"topic,omitempty"`
	Timeout              int32                    `protobuf:"varint,7,opt,name=timeout,proto3" json:"timeout,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *PubSubQuery) Reset()         { *m = PubSubQuery{} }
func (m *PubSubQuery) String() string { return proto.CompactTextString(m) }
func (*PubSubQuery) ProtoMessage()    {}
func (*PubSubQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{2}
}

func (m *PubSubQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PubSubQuery.Unmarshal(m, b)
}
func (m *PubSubQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PubSubQuery.Marshal(b, m, deterministic)
}
func (m *PubSubQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PubSubQuery.Merge(m, src)
}
func (m *PubSubQuery) XXX_Size() int {
	return xxx_messageInfo_PubSubQuery.Size(m)
}
func (m *PubSubQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_PubSubQuery.DiscardUnknown(m)
}

var xxx_messageInfo_PubSubQuery proto.InternalMessageInfo

func (m *PubSubQuery) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PubSubQuery) GetType() Query_Type {
	if m != nil {
		return m.Type
	}
	return Query_THREAD_SNAPSHOTS
}

func (m *PubSubQuery) GetPayload() *any.Any {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *PubSubQuery) GetResponseType() PubSubQuery_ResponseType {
	if m != nil {
		return m.ResponseType
	}
	return PubSubQuery_P2P
}

func (m *PubSubQuery) GetExclude() []string {
	if m != nil {
		return m.Exclude
	}
	return nil
}

func (m *PubSubQuery) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *PubSubQuery) GetTimeout() int32 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

type QueryResult struct {
	Id                   string               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Date                 *timestamp.Timestamp `protobuf:"bytes,2,opt,name=date,proto3" json:"date,omitempty"`
	Local                bool                 `protobuf:"varint,3,opt,name=local,proto3" json:"local,omitempty"`
	Value                *any.Any             `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *QueryResult) Reset()         { *m = QueryResult{} }
func (m *QueryResult) String() string { return proto.CompactTextString(m) }
func (*QueryResult) ProtoMessage()    {}
func (*QueryResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{3}
}

func (m *QueryResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryResult.Unmarshal(m, b)
}
func (m *QueryResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryResult.Marshal(b, m, deterministic)
}
func (m *QueryResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResult.Merge(m, src)
}
func (m *QueryResult) XXX_Size() int {
	return xxx_messageInfo_QueryResult.Size(m)
}
func (m *QueryResult) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryResult.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResult proto.InternalMessageInfo

func (m *QueryResult) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *QueryResult) GetDate() *timestamp.Timestamp {
	if m != nil {
		return m.Date
	}
	return nil
}

func (m *QueryResult) GetLocal() bool {
	if m != nil {
		return m.Local
	}
	return false
}

func (m *QueryResult) GetValue() *any.Any {
	if m != nil {
		return m.Value
	}
	return nil
}

type QueryResults struct {
	Type                 Query_Type     `protobuf:"varint,1,opt,name=type,proto3,enum=Query_Type" json:"type,omitempty"`
	Items                []*QueryResult `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *QueryResults) Reset()         { *m = QueryResults{} }
func (m *QueryResults) String() string { return proto.CompactTextString(m) }
func (*QueryResults) ProtoMessage()    {}
func (*QueryResults) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{4}
}

func (m *QueryResults) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryResults.Unmarshal(m, b)
}
func (m *QueryResults) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryResults.Marshal(b, m, deterministic)
}
func (m *QueryResults) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResults.Merge(m, src)
}
func (m *QueryResults) XXX_Size() int {
	return xxx_messageInfo_QueryResults.Size(m)
}
func (m *QueryResults) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryResults.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResults proto.InternalMessageInfo

func (m *QueryResults) GetType() Query_Type {
	if m != nil {
		return m.Type
	}
	return Query_THREAD_SNAPSHOTS
}

func (m *QueryResults) GetItems() []*QueryResult {
	if m != nil {
		return m.Items
	}
	return nil
}

type PubSubQueryResults struct {
	Id                   string        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Results              *QueryResults `protobuf:"bytes,2,opt,name=results,proto3" json:"results,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *PubSubQueryResults) Reset()         { *m = PubSubQueryResults{} }
func (m *PubSubQueryResults) String() string { return proto.CompactTextString(m) }
func (*PubSubQueryResults) ProtoMessage()    {}
func (*PubSubQueryResults) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{5}
}

func (m *PubSubQueryResults) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PubSubQueryResults.Unmarshal(m, b)
}
func (m *PubSubQueryResults) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PubSubQueryResults.Marshal(b, m, deterministic)
}
func (m *PubSubQueryResults) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PubSubQueryResults.Merge(m, src)
}
func (m *PubSubQueryResults) XXX_Size() int {
	return xxx_messageInfo_PubSubQueryResults.Size(m)
}
func (m *PubSubQueryResults) XXX_DiscardUnknown() {
	xxx_messageInfo_PubSubQueryResults.DiscardUnknown(m)
}

var xxx_messageInfo_PubSubQueryResults proto.InternalMessageInfo

func (m *PubSubQueryResults) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PubSubQueryResults) GetResults() *QueryResults {
	if m != nil {
		return m.Results
	}
	return nil
}

type StreamQuery struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Startindex           int64    `protobuf:"varint,2,opt,name=startindex,proto3" json:"startindex,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StreamQuery) Reset()         { *m = StreamQuery{} }
func (m *StreamQuery) String() string { return proto.CompactTextString(m) }
func (*StreamQuery) ProtoMessage()    {}
func (*StreamQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{6}
}

func (m *StreamQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamQuery.Unmarshal(m, b)
}
func (m *StreamQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamQuery.Marshal(b, m, deterministic)
}
func (m *StreamQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamQuery.Merge(m, src)
}
func (m *StreamQuery) XXX_Size() int {
	return xxx_messageInfo_StreamQuery.Size(m)
}
func (m *StreamQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamQuery.DiscardUnknown(m)
}

var xxx_messageInfo_StreamQuery proto.InternalMessageInfo

func (m *StreamQuery) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *StreamQuery) GetStartindex() int64 {
	if m != nil {
		return m.Startindex
	}
	return 0
}

type StreamQueryResultItem struct {
	Pid                  string   `protobuf:"bytes,1,opt,name=pid,proto3" json:"pid,omitempty"`
	Hopcnt               int32    `protobuf:"varint,2,opt,name=hopcnt,proto3" json:"hopcnt,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StreamQueryResultItem) Reset()         { *m = StreamQueryResultItem{} }
func (m *StreamQueryResultItem) String() string { return proto.CompactTextString(m) }
func (*StreamQueryResultItem) ProtoMessage()    {}
func (*StreamQueryResultItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{7}
}

func (m *StreamQueryResultItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamQueryResultItem.Unmarshal(m, b)
}
func (m *StreamQueryResultItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamQueryResultItem.Marshal(b, m, deterministic)
}
func (m *StreamQueryResultItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamQueryResultItem.Merge(m, src)
}
func (m *StreamQueryResultItem) XXX_Size() int {
	return xxx_messageInfo_StreamQueryResultItem.Size(m)
}
func (m *StreamQueryResultItem) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamQueryResultItem.DiscardUnknown(m)
}

var xxx_messageInfo_StreamQueryResultItem proto.InternalMessageInfo

func (m *StreamQueryResultItem) GetPid() string {
	if m != nil {
		return m.Pid
	}
	return ""
}

func (m *StreamQueryResultItem) GetHopcnt() int32 {
	if m != nil {
		return m.Hopcnt
	}
	return 0
}

func init() {
	proto.RegisterEnum("QueryOptions_FilterType", QueryOptions_FilterType_name, QueryOptions_FilterType_value)
	proto.RegisterEnum("Query_Type", Query_Type_name, Query_Type_value)
	proto.RegisterEnum("PubSubQuery_ResponseType", PubSubQuery_ResponseType_name, PubSubQuery_ResponseType_value)
	proto.RegisterType((*QueryOptions)(nil), "QueryOptions")
	proto.RegisterType((*Query)(nil), "Query")
	proto.RegisterType((*PubSubQuery)(nil), "PubSubQuery")
	proto.RegisterType((*QueryResult)(nil), "QueryResult")
	proto.RegisterType((*QueryResults)(nil), "QueryResults")
	proto.RegisterType((*PubSubQueryResults)(nil), "PubSubQueryResults")
	proto.RegisterType((*StreamQuery)(nil), "StreamQuery")
	proto.RegisterType((*StreamQueryResultItem)(nil), "StreamQueryResultItem")
}

func init() { proto.RegisterFile("query.proto", fileDescriptor_5c6ac9b241082464) }

var fileDescriptor_5c6ac9b241082464 = []byte{
	// 695 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xcd, 0x6e, 0xda, 0x4c,
	0x14, 0x8d, 0x0d, 0x86, 0x70, 0x4d, 0x90, 0x35, 0xca, 0x17, 0x39, 0x51, 0x94, 0x20, 0x7f, 0x8b,
	0xa0, 0xef, 0x93, 0x26, 0x15, 0x69, 0x97, 0x59, 0x10, 0x20, 0x02, 0x35, 0x01, 0x3a, 0x76, 0x2a,
	0xb5, 0x1b, 0x64, 0x60, 0x92, 0xba, 0x35, 0x9e, 0xa9, 0x3d, 0x6e, 0xc3, 0x33, 0x74, 0xd1, 0x17,
	0xe8, 0x3b, 0xf6, 0x15, 0x2a, 0xcf, 0xe0, 0xc4, 0xf9, 0x93, 0xba, 0xf3, 0x9d, 0x73, 0x7c, 0x7d,
	0xce, 0xb9, 0x77, 0x0c, 0xe6, 0xd7, 0x94, 0xc6, 0x2b, 0xcc, 0x63, 0x26, 0xd8, 0xde, 0xee, 0x0d,
	0x63, 0x37, 0x21, 0x3d, 0x96, 0xd5, 0x2c, 0xbd, 0x3e, 0xf6, 0xa3, 0x1c, 0x3a, 0x7c, 0x0c, 0x89,
	0x60, 0x49, 0x13, 0xe1, 0x2f, 0xb9, 0x22, 0x38, 0xbf, 0x35, 0xa8, 0xbf, 0xcb, 0x7a, 0x8d, 0xb9,
	0x08, 0x58, 0x94, 0xa0, 0x7d, 0xa8, 0x85, 0x6c, 0xee, 0x87, 0xe3, 0x28, 0x5c, 0xd9, 0x5a, 0x53,
	0x6b, 0x6d, 0x92, 0xfb, 0x03, 0x74, 0x00, 0x10, 0xd3, 0x25, 0x13, 0x54, 0xc2, 0x15, 0x09, 0x17,
	0x4e, 0xd0, 0x36, 0x18, 0x61, 0xb0, 0x0c, 0x84, 0xad, 0x37, 0xb5, 0x96, 0x41, 0x54, 0x81, 0x10,
	0x94, 0xbf, 0xfb, 0x81, 0xb0, 0x4b, 0xf2, 0x50, 0x3e, 0xa3, 0x57, 0x50, 0xb9, 0x0e, 0x42, 0x41,
	0x63, 0xbb, 0xdc, 0xd4, 0x5a, 0x8d, 0xb6, 0x8d, 0x8b, 0x32, 0xf0, 0xb9, 0xc4, 0xbc, 0x15, 0xa7,
	0x64, 0xcd, 0x43, 0x36, 0x54, 0xe9, 0xed, 0x3c, 0x4c, 0x17, 0xd4, 0x36, 0x9a, 0xa5, 0x56, 0x8d,
	0xe4, 0xa5, 0xf3, 0x3f, 0xc0, 0x3d, 0x1f, 0x6d, 0x41, 0x6d, 0x34, 0x9e, 0x9e, 0x0f, 0x2f, 0xbc,
	0x3e, 0xb1, 0x36, 0x50, 0x03, 0x60, 0x30, 0xec, 0xf5, 0xa7, 0xe3, 0x8b, 0x5e, 0x9f, 0x58, 0x9a,
	0xf3, 0x43, 0x07, 0x43, 0x7e, 0x0a, 0x35, 0x40, 0x0f, 0x16, 0xd2, 0x63, 0x8d, 0xe8, 0xc1, 0x22,
	0x13, 0x2f, 0xd8, 0x17, 0x1a, 0x49, 0xf1, 0x35, 0xa2, 0x0a, 0x74, 0x08, 0x65, 0xb1, 0xe2, 0x54,
	0x8a, 0x6f, 0xb4, 0x4d, 0x25, 0x13, 0x4b, 0x65, 0x12, 0x40, 0x47, 0x50, 0x65, 0x4a, 0xb5, 0xb4,
	0x62, 0xb6, 0xb7, 0x1e, 0x58, 0x21, 0x39, 0x8a, 0x30, 0x54, 0xb9, 0xbf, 0x0a, 0x99, 0xbf, 0xb0,
	0x0d, 0x49, 0xdc, 0xc6, 0x6a, 0x3c, 0x38, 0x1f, 0x0f, 0xee, 0x44, 0x2b, 0x92, 0x93, 0x9c, 0x19,
	0x94, 0xa5, 0xa1, 0x6d, 0xb0, 0xbc, 0x01, 0xe9, 0x77, 0x7a, 0x53, 0x77, 0xd4, 0x99, 0xb8, 0x83,
	0xb1, 0xe7, 0x5a, 0x1b, 0xa8, 0x0e, 0x9b, 0xdd, 0xf1, 0xc8, 0xeb, 0x74, 0x3d, 0xd7, 0xd2, 0x90,
	0x05, 0xf5, 0xf7, 0xc3, 0x5e, 0x7f, 0x3c, 0xed, 0x0e, 0xae, 0x46, 0x6f, 0x5d, 0xab, 0x8d, 0x6a,
	0x60, 0xc8, 0x13, 0xeb, 0x24, 0x4b, 0xc4, 0xfd, 0x30, 0xea, 0x66, 0x99, 0xf4, 0xad, 0xd7, 0x08,
	0xa0, 0xe2, 0x7a, 0xa4, 0xdf, 0xb9, 0xb4, 0xde, 0x38, 0xbf, 0x74, 0x30, 0x27, 0xe9, 0xcc, 0x4d,
	0x67, 0xcf, 0x67, 0x92, 0xbb, 0xd7, 0x5f, 0x72, 0x5f, 0x30, 0x55, 0xfa, 0x0b, 0x53, 0xe8, 0x14,
	0xea, 0x31, 0x4d, 0x38, 0x8b, 0x12, 0x9a, 0x75, 0x59, 0x4f, 0x7f, 0x17, 0x17, 0x44, 0x60, 0x52,
	0x20, 0x90, 0x07, 0xf4, 0x97, 0x97, 0x40, 0x4d, 0x8f, 0x07, 0x73, 0xb9, 0x95, 0x72, 0x7a, 0x3c,
	0x98, 0x67, 0xfc, 0x6c, 0xe5, 0x59, 0x2a, 0xec, 0xaa, 0xdc, 0xbe, 0xbc, 0x74, 0xfe, 0x85, 0x7a,
	0xf1, 0x3b, 0xa8, 0x0a, 0xa5, 0x49, 0x7b, 0x62, 0x6d, 0x64, 0xf1, 0x4c, 0xae, 0xce, 0xdc, 0xab,
	0x33, 0x4b, 0x73, 0x7e, 0x6a, 0x60, 0x4a, 0x4d, 0x84, 0x26, 0x69, 0x28, 0x9e, 0xc4, 0x83, 0xa1,
	0xbc, 0xf0, 0x85, 0x8a, 0xc7, 0x6c, 0xef, 0x3d, 0xb1, 0xee, 0xe5, 0xd7, 0x8d, 0x48, 0x9e, 0xbc,
	0x1f, 0xd9, 0x65, 0x92, 0x59, 0x6d, 0x12, 0x55, 0xa0, 0xff, 0xc0, 0xf8, 0xe6, 0x87, 0x29, 0x5d,
	0xef, 0xcf, 0xf3, 0x09, 0x2a, 0x8a, 0xe3, 0xae, 0xef, 0xab, 0x12, 0x94, 0xdc, 0x0d, 0x48, 0x7b,
	0x69, 0x40, 0x0e, 0x18, 0x81, 0xa0, 0xcb, 0xc4, 0xd6, 0x9b, 0xa5, 0x96, 0xd9, 0xae, 0xe3, 0xc2,
	0xeb, 0x44, 0x41, 0xce, 0x25, 0xa0, 0x42, 0xfe, 0x79, 0xeb, 0xc7, 0x66, 0x8f, 0xa0, 0x1a, 0x2b,
	0x68, 0xed, 0x77, 0xab, 0xd8, 0x2b, 0x21, 0x39, 0xea, 0x9c, 0x82, 0xe9, 0x8a, 0x98, 0xfa, 0xcb,
	0xe7, 0x77, 0xea, 0x00, 0x20, 0x11, 0x7e, 0x2c, 0x82, 0x68, 0x41, 0x6f, 0x65, 0xab, 0x12, 0x29,
	0x9c, 0x38, 0x1d, 0xf8, 0xa7, 0xf0, 0xba, 0xea, 0x3e, 0x14, 0x74, 0x89, 0x2c, 0x28, 0xf1, 0xbb,
	0x4e, 0xd9, 0x23, 0xda, 0x81, 0xca, 0x27, 0xc6, 0xe7, 0x51, 0xfe, 0xc3, 0x59, 0x57, 0x67, 0xfb,
	0xb0, 0x93, 0x7c, 0x16, 0x29, 0x66, 0x9c, 0x46, 0x11, 0x15, 0x58, 0xd0, 0x5b, 0x11, 0x84, 0x94,
	0xcf, 0x3e, 0xea, 0x7c, 0x36, 0xab, 0xc8, 0x60, 0x4f, 0xfe, 0x04, 0x00, 0x00, 0xff, 0xff, 0x14,
	0x76, 0xd8, 0xba, 0x46, 0x05, 0x00, 0x00,
}
