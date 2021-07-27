// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpc.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

// ============= 注册节点 ============
type RegisterReq struct {
	PeerNode             string   `protobuf:"bytes,1,opt,name=PeerNode,proto3" json:"PeerNode,omitempty"`
	SlaveKey             string   `protobuf:"bytes,2,opt,name=SlaveKey,proto3" json:"SlaveKey,omitempty"`
	MasterKey            string   `protobuf:"bytes,3,opt,name=MasterKey,proto3" json:"MasterKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterReq) Reset()         { *m = RegisterReq{} }
func (m *RegisterReq) String() string { return proto.CompactTextString(m) }
func (*RegisterReq) ProtoMessage()    {}
func (*RegisterReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{1}
}

func (m *RegisterReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterReq.Unmarshal(m, b)
}
func (m *RegisterReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterReq.Marshal(b, m, deterministic)
}
func (m *RegisterReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterReq.Merge(m, src)
}
func (m *RegisterReq) XXX_Size() int {
	return xxx_messageInfo_RegisterReq.Size(m)
}
func (m *RegisterReq) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterReq.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterReq proto.InternalMessageInfo

func (m *RegisterReq) GetPeerNode() string {
	if m != nil {
		return m.PeerNode
	}
	return ""
}

func (m *RegisterReq) GetSlaveKey() string {
	if m != nil {
		return m.SlaveKey
	}
	return ""
}

func (m *RegisterReq) GetMasterKey() string {
	if m != nil {
		return m.MasterKey
	}
	return ""
}

type RegisterAck struct {
	MasterKey            string   `protobuf:"bytes,1,opt,name=MasterKey,proto3" json:"MasterKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterAck) Reset()         { *m = RegisterAck{} }
func (m *RegisterAck) String() string { return proto.CompactTextString(m) }
func (*RegisterAck) ProtoMessage()    {}
func (*RegisterAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{2}
}

func (m *RegisterAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterAck.Unmarshal(m, b)
}
func (m *RegisterAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterAck.Marshal(b, m, deterministic)
}
func (m *RegisterAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterAck.Merge(m, src)
}
func (m *RegisterAck) XXX_Size() int {
	return xxx_messageInfo_RegisterAck.Size(m)
}
func (m *RegisterAck) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterAck.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterAck proto.InternalMessageInfo

func (m *RegisterAck) GetMasterKey() string {
	if m != nil {
		return m.MasterKey
	}
	return ""
}

type MountReq struct {
	Name                 string   `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MountReq) Reset()         { *m = MountReq{} }
func (m *MountReq) String() string { return proto.CompactTextString(m) }
func (*MountReq) ProtoMessage()    {}
func (*MountReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{3}
}

func (m *MountReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MountReq.Unmarshal(m, b)
}
func (m *MountReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MountReq.Marshal(b, m, deterministic)
}
func (m *MountReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MountReq.Merge(m, src)
}
func (m *MountReq) XXX_Size() int {
	return xxx_messageInfo_MountReq.Size(m)
}
func (m *MountReq) XXX_DiscardUnknown() {
	xxx_messageInfo_MountReq.DiscardUnknown(m)
}

var xxx_messageInfo_MountReq proto.InternalMessageInfo

func (m *MountReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// ============= 同步挂载点 ===========
type SyncReq struct {
	Id                   int64    `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SyncReq) Reset()         { *m = SyncReq{} }
func (m *SyncReq) String() string { return proto.CompactTextString(m) }
func (*SyncReq) ProtoMessage()    {}
func (*SyncReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{4}
}

func (m *SyncReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncReq.Unmarshal(m, b)
}
func (m *SyncReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncReq.Marshal(b, m, deterministic)
}
func (m *SyncReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncReq.Merge(m, src)
}
func (m *SyncReq) XXX_Size() int {
	return xxx_messageInfo_SyncReq.Size(m)
}
func (m *SyncReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncReq.DiscardUnknown(m)
}

var xxx_messageInfo_SyncReq proto.InternalMessageInfo

func (m *SyncReq) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type MountLogItem struct {
	Id                   int64    `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	IsAdd                bool     `protobuf:"varint,2,opt,name=IsAdd,proto3" json:"IsAdd,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=Name,proto3" json:"Name,omitempty"`
	PeerAddr             string   `protobuf:"bytes,4,opt,name=PeerAddr,proto3" json:"PeerAddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MountLogItem) Reset()         { *m = MountLogItem{} }
func (m *MountLogItem) String() string { return proto.CompactTextString(m) }
func (*MountLogItem) ProtoMessage()    {}
func (*MountLogItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{5}
}

func (m *MountLogItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MountLogItem.Unmarshal(m, b)
}
func (m *MountLogItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MountLogItem.Marshal(b, m, deterministic)
}
func (m *MountLogItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MountLogItem.Merge(m, src)
}
func (m *MountLogItem) XXX_Size() int {
	return xxx_messageInfo_MountLogItem.Size(m)
}
func (m *MountLogItem) XXX_DiscardUnknown() {
	xxx_messageInfo_MountLogItem.DiscardUnknown(m)
}

var xxx_messageInfo_MountLogItem proto.InternalMessageInfo

func (m *MountLogItem) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *MountLogItem) GetIsAdd() bool {
	if m != nil {
		return m.IsAdd
	}
	return false
}

func (m *MountLogItem) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MountLogItem) GetPeerAddr() string {
	if m != nil {
		return m.PeerAddr
	}
	return ""
}

type SyncAck struct {
	List                 []*MountLogItem `protobuf:"bytes,1,rep,name=List,proto3" json:"List,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *SyncAck) Reset()         { *m = SyncAck{} }
func (m *SyncAck) String() string { return proto.CompactTextString(m) }
func (*SyncAck) ProtoMessage()    {}
func (*SyncAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{6}
}

func (m *SyncAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncAck.Unmarshal(m, b)
}
func (m *SyncAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncAck.Marshal(b, m, deterministic)
}
func (m *SyncAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncAck.Merge(m, src)
}
func (m *SyncAck) XXX_Size() int {
	return xxx_messageInfo_SyncAck.Size(m)
}
func (m *SyncAck) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncAck.DiscardUnknown(m)
}

var xxx_messageInfo_SyncAck proto.InternalMessageInfo

func (m *SyncAck) GetList() []*MountLogItem {
	if m != nil {
		return m.List
	}
	return nil
}

// ============= 远程调用 =============
type CheckKeyReq struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CheckKeyReq) Reset()         { *m = CheckKeyReq{} }
func (m *CheckKeyReq) String() string { return proto.CompactTextString(m) }
func (*CheckKeyReq) ProtoMessage()    {}
func (*CheckKeyReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{7}
}

func (m *CheckKeyReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckKeyReq.Unmarshal(m, b)
}
func (m *CheckKeyReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckKeyReq.Marshal(b, m, deterministic)
}
func (m *CheckKeyReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckKeyReq.Merge(m, src)
}
func (m *CheckKeyReq) XXX_Size() int {
	return xxx_messageInfo_CheckKeyReq.Size(m)
}
func (m *CheckKeyReq) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckKeyReq.DiscardUnknown(m)
}

var xxx_messageInfo_CheckKeyReq proto.InternalMessageInfo

func (m *CheckKeyReq) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type CallReq struct {
	ServiceName          string   `protobuf:"bytes,1,opt,name=ServiceName,proto3" json:"ServiceName,omitempty"`
	NodeName             string   `protobuf:"bytes,2,opt,name=NodeName,proto3" json:"NodeName,omitempty"`
	Method               string   `protobuf:"bytes,3,opt,name=Method,proto3" json:"Method,omitempty"`
	Args                 [][]byte `protobuf:"bytes,4,rep,name=Args,proto3" json:"Args,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CallReq) Reset()         { *m = CallReq{} }
func (m *CallReq) String() string { return proto.CompactTextString(m) }
func (*CallReq) ProtoMessage()    {}
func (*CallReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{8}
}

func (m *CallReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CallReq.Unmarshal(m, b)
}
func (m *CallReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CallReq.Marshal(b, m, deterministic)
}
func (m *CallReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CallReq.Merge(m, src)
}
func (m *CallReq) XXX_Size() int {
	return xxx_messageInfo_CallReq.Size(m)
}
func (m *CallReq) XXX_DiscardUnknown() {
	xxx_messageInfo_CallReq.DiscardUnknown(m)
}

var xxx_messageInfo_CallReq proto.InternalMessageInfo

func (m *CallReq) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *CallReq) GetNodeName() string {
	if m != nil {
		return m.NodeName
	}
	return ""
}

func (m *CallReq) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *CallReq) GetArgs() [][]byte {
	if m != nil {
		return m.Args
	}
	return nil
}

type CallAck struct {
	Args                 [][]byte `protobuf:"bytes,1,rep,name=Args,proto3" json:"Args,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CallAck) Reset()         { *m = CallAck{} }
func (m *CallAck) String() string { return proto.CompactTextString(m) }
func (*CallAck) ProtoMessage()    {}
func (*CallAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{9}
}

func (m *CallAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CallAck.Unmarshal(m, b)
}
func (m *CallAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CallAck.Marshal(b, m, deterministic)
}
func (m *CallAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CallAck.Merge(m, src)
}
func (m *CallAck) XXX_Size() int {
	return xxx_messageInfo_CallAck.Size(m)
}
func (m *CallAck) XXX_DiscardUnknown() {
	xxx_messageInfo_CallAck.DiscardUnknown(m)
}

var xxx_messageInfo_CallAck proto.InternalMessageInfo

func (m *CallAck) GetArgs() [][]byte {
	if m != nil {
		return m.Args
	}
	return nil
}

// ============= 流抽象 ==============
type StreamMsg struct {
	// Types that are valid to be assigned to StreamType:
	//	*StreamMsg_Req
	//	*StreamMsg_Data
	StreamType           isStreamMsg_StreamType `protobuf_oneof:"StreamType"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *StreamMsg) Reset()         { *m = StreamMsg{} }
func (m *StreamMsg) String() string { return proto.CompactTextString(m) }
func (*StreamMsg) ProtoMessage()    {}
func (*StreamMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{10}
}

func (m *StreamMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamMsg.Unmarshal(m, b)
}
func (m *StreamMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamMsg.Marshal(b, m, deterministic)
}
func (m *StreamMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamMsg.Merge(m, src)
}
func (m *StreamMsg) XXX_Size() int {
	return xxx_messageInfo_StreamMsg.Size(m)
}
func (m *StreamMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamMsg.DiscardUnknown(m)
}

var xxx_messageInfo_StreamMsg proto.InternalMessageInfo

type isStreamMsg_StreamType interface {
	isStreamMsg_StreamType()
}

type StreamMsg_Req struct {
	Req *CallReq `protobuf:"bytes,1,opt,name=Req,proto3,oneof"`
}

type StreamMsg_Data struct {
	Data *Bytes `protobuf:"bytes,2,opt,name=Data,proto3,oneof"`
}

func (*StreamMsg_Req) isStreamMsg_StreamType() {}

func (*StreamMsg_Data) isStreamMsg_StreamType() {}

func (m *StreamMsg) GetStreamType() isStreamMsg_StreamType {
	if m != nil {
		return m.StreamType
	}
	return nil
}

func (m *StreamMsg) GetReq() *CallReq {
	if x, ok := m.GetStreamType().(*StreamMsg_Req); ok {
		return x.Req
	}
	return nil
}

func (m *StreamMsg) GetData() *Bytes {
	if x, ok := m.GetStreamType().(*StreamMsg_Data); ok {
		return x.Data
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*StreamMsg) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*StreamMsg_Req)(nil),
		(*StreamMsg_Data)(nil),
	}
}

type Bytes struct {
	Bytes                []byte   `protobuf:"bytes,1,opt,name=Bytes,proto3" json:"Bytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Bytes) Reset()         { *m = Bytes{} }
func (m *Bytes) String() string { return proto.CompactTextString(m) }
func (*Bytes) ProtoMessage()    {}
func (*Bytes) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{11}
}

func (m *Bytes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Bytes.Unmarshal(m, b)
}
func (m *Bytes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Bytes.Marshal(b, m, deterministic)
}
func (m *Bytes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bytes.Merge(m, src)
}
func (m *Bytes) XXX_Size() int {
	return xxx_messageInfo_Bytes.Size(m)
}
func (m *Bytes) XXX_DiscardUnknown() {
	xxx_messageInfo_Bytes.DiscardUnknown(m)
}

var xxx_messageInfo_Bytes proto.InternalMessageInfo

func (m *Bytes) GetBytes() []byte {
	if m != nil {
		return m.Bytes
	}
	return nil
}

type NodeName struct {
	Str                  string   `protobuf:"bytes,1,opt,name=Str,proto3" json:"Str,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeName) Reset()         { *m = NodeName{} }
func (m *NodeName) String() string { return proto.CompactTextString(m) }
func (*NodeName) ProtoMessage()    {}
func (*NodeName) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{12}
}

func (m *NodeName) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeName.Unmarshal(m, b)
}
func (m *NodeName) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeName.Marshal(b, m, deterministic)
}
func (m *NodeName) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeName.Merge(m, src)
}
func (m *NodeName) XXX_Size() int {
	return xxx_messageInfo_NodeName.Size(m)
}
func (m *NodeName) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeName.DiscardUnknown(m)
}

var xxx_messageInfo_NodeName proto.InternalMessageInfo

func (m *NodeName) GetStr() string {
	if m != nil {
		return m.Str
	}
	return ""
}

// ============= 挂载时间 ===========
type MountTimeAck struct {
	Unix                 int64    `protobuf:"varint,1,opt,name=Unix,proto3" json:"Unix,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MountTimeAck) Reset()         { *m = MountTimeAck{} }
func (m *MountTimeAck) String() string { return proto.CompactTextString(m) }
func (*MountTimeAck) ProtoMessage()    {}
func (*MountTimeAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{13}
}

func (m *MountTimeAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MountTimeAck.Unmarshal(m, b)
}
func (m *MountTimeAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MountTimeAck.Marshal(b, m, deterministic)
}
func (m *MountTimeAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MountTimeAck.Merge(m, src)
}
func (m *MountTimeAck) XXX_Size() int {
	return xxx_messageInfo_MountTimeAck.Size(m)
}
func (m *MountTimeAck) XXX_DiscardUnknown() {
	xxx_messageInfo_MountTimeAck.DiscardUnknown(m)
}

var xxx_messageInfo_MountTimeAck proto.InternalMessageInfo

func (m *MountTimeAck) GetUnix() int64 {
	if m != nil {
		return m.Unix
	}
	return 0
}

func init() {
	proto.RegisterType((*Empty)(nil), "pb.Empty")
	proto.RegisterType((*RegisterReq)(nil), "pb.RegisterReq")
	proto.RegisterType((*RegisterAck)(nil), "pb.RegisterAck")
	proto.RegisterType((*MountReq)(nil), "pb.MountReq")
	proto.RegisterType((*SyncReq)(nil), "pb.SyncReq")
	proto.RegisterType((*MountLogItem)(nil), "pb.MountLogItem")
	proto.RegisterType((*SyncAck)(nil), "pb.SyncAck")
	proto.RegisterType((*CheckKeyReq)(nil), "pb.CheckKeyReq")
	proto.RegisterType((*CallReq)(nil), "pb.CallReq")
	proto.RegisterType((*CallAck)(nil), "pb.CallAck")
	proto.RegisterType((*StreamMsg)(nil), "pb.StreamMsg")
	proto.RegisterType((*Bytes)(nil), "pb.Bytes")
	proto.RegisterType((*NodeName)(nil), "pb.NodeName")
	proto.RegisterType((*MountTimeAck)(nil), "pb.MountTimeAck")
}

func init() { proto.RegisterFile("rpc.proto", fileDescriptor_77a6da22d6a3feb1) }

var fileDescriptor_77a6da22d6a3feb1 = []byte{
	// 565 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xad, 0x63, 0xe7, 0xc3, 0xe3, 0x00, 0xd5, 0x0a, 0xa1, 0x10, 0xb5, 0xc4, 0x5a, 0x21, 0x88,
	0x28, 0x0a, 0x55, 0x38, 0x72, 0x4a, 0x0a, 0x52, 0xa3, 0x36, 0x15, 0x5a, 0xb7, 0x17, 0x38, 0x39,
	0xf6, 0x2a, 0x8d, 0xf2, 0x61, 0xd7, 0x5e, 0x2a, 0xfc, 0x8f, 0xf8, 0x39, 0xfc, 0x24, 0x34, 0x63,
	0xaf, 0xed, 0x84, 0x43, 0x6f, 0x3b, 0xfb, 0x26, 0xb3, 0xef, 0xcd, 0x7b, 0x31, 0xd8, 0x49, 0x1c,
	0x8c, 0xe2, 0x24, 0x52, 0x11, 0x6b, 0xc4, 0x0b, 0xde, 0x86, 0xe6, 0xb7, 0x6d, 0xac, 0x32, 0x1e,
	0x80, 0x23, 0xe4, 0x72, 0x95, 0x2a, 0x99, 0x08, 0xf9, 0xc0, 0xfa, 0xd0, 0xf9, 0x2e, 0x65, 0x72,
	0x13, 0x85, 0xb2, 0x67, 0xb8, 0xc6, 0xd0, 0x16, 0x65, 0x8d, 0x98, 0xb7, 0xf1, 0x1f, 0xe5, 0x95,
	0xcc, 0x7a, 0x8d, 0x1c, 0xd3, 0x35, 0x3b, 0x01, 0x7b, 0xee, 0xe3, 0x10, 0x04, 0x4d, 0x02, 0xab,
	0x0b, 0x7e, 0x56, 0x3d, 0x32, 0x09, 0xd6, 0xfb, 0xcd, 0xc6, 0x61, 0xf3, 0x1b, 0xe8, 0xcc, 0xa3,
	0x5f, 0x3b, 0x85, 0x74, 0x18, 0x58, 0x37, 0xfe, 0x56, 0x53, 0xa1, 0x33, 0x7f, 0x0d, 0x6d, 0x2f,
	0xdb, 0x05, 0x08, 0x3f, 0x87, 0xc6, 0x2c, 0x24, 0xd0, 0x14, 0x8d, 0x59, 0xc8, 0x43, 0xe8, 0xd2,
	0x4f, 0xaf, 0xa3, 0xe5, 0x4c, 0xc9, 0xed, 0x21, 0xce, 0x5e, 0x42, 0x73, 0x96, 0x4e, 0xc2, 0x90,
	0xe8, 0x77, 0x44, 0x5e, 0x94, 0x8f, 0x98, 0xd5, 0x23, 0x7a, 0x0f, 0x93, 0x30, 0x4c, 0x7a, 0x56,
	0xb5, 0x07, 0xac, 0xf9, 0xa7, 0x9c, 0x00, 0x2a, 0x79, 0x0b, 0xd6, 0xf5, 0x2a, 0x55, 0x3d, 0xc3,
	0x35, 0x87, 0xce, 0xf8, 0x78, 0x14, 0x2f, 0x46, 0x75, 0x02, 0x82, 0x50, 0x3e, 0x00, 0xe7, 0xe2,
	0x5e, 0x06, 0xeb, 0x2b, 0x99, 0x21, 0xeb, 0x63, 0x30, 0xd7, 0xa5, 0x70, 0x3c, 0xf2, 0x14, 0xda,
	0x17, 0xfe, 0x66, 0x83, 0xa0, 0x0b, 0x8e, 0x27, 0x93, 0xc7, 0x55, 0x20, 0x6b, 0xc2, 0xeb, 0x57,
	0x48, 0x0d, 0xed, 0x20, 0xb8, 0xb0, 0x41, 0xd7, 0xec, 0x15, 0xb4, 0xe6, 0x52, 0xdd, 0x47, 0x61,
	0x21, 0xa6, 0xa8, 0x50, 0xe2, 0x24, 0x59, 0xa6, 0x3d, 0xcb, 0x35, 0x87, 0x5d, 0x41, 0x67, 0x7e,
	0x9a, 0x3f, 0x8a, 0x32, 0x34, 0x6c, 0xd4, 0xe0, 0x9f, 0x60, 0x7b, 0x2a, 0x91, 0xfe, 0x76, 0x9e,
	0x2e, 0xd9, 0x00, 0x4c, 0x21, 0x1f, 0x88, 0x8d, 0x33, 0x76, 0x50, 0x66, 0xc1, 0xf7, 0xf2, 0x48,
	0x20, 0xc2, 0x06, 0x60, 0x7d, 0xf5, 0x95, 0x4f, 0x84, 0x9c, 0xb1, 0x8d, 0x1d, 0xd3, 0x4c, 0xc9,
	0xf4, 0xf2, 0x48, 0x10, 0x30, 0xed, 0x02, 0xe4, 0xe3, 0x6e, 0xb3, 0x58, 0xf2, 0x53, 0x68, 0x12,
	0x8c, 0x8e, 0xd0, 0x81, 0x46, 0x77, 0x45, 0x5e, 0xf0, 0x93, 0x4a, 0x22, 0x6e, 0xcb, 0x53, 0x89,
	0xde, 0x96, 0xa7, 0x12, 0xce, 0x0b, 0x97, 0x6f, 0x57, 0x5b, 0x59, 0xb0, 0xbf, 0xdb, 0xad, 0x7e,
	0x17, 0x3e, 0xd3, 0x79, 0xfc, 0xc7, 0xd0, 0x19, 0x13, 0x71, 0xc0, 0x3e, 0x42, 0x47, 0xe7, 0x8f,
	0xbd, 0x40, 0x6e, 0xb5, 0xc8, 0xf7, 0xf7, 0x2e, 0x70, 0x9e, 0x0b, 0x4d, 0x9a, 0xcf, 0xba, 0xa5,
	0x9f, 0xd8, 0x47, 0xa2, 0xe8, 0x4f, 0xc3, 0x38, 0xb4, 0xef, 0x76, 0x4f, 0xf4, 0xbc, 0x07, 0x0b,
	0x53, 0xc2, 0x68, 0x5b, 0x45, 0x60, 0xfb, 0xff, 0x25, 0xe4, 0xdc, 0x18, 0xff, 0x35, 0xa0, 0x8d,
	0xd9, 0x42, 0xa2, 0xef, 0xa0, 0xa3, 0x93, 0x92, 0x13, 0xad, 0xe5, 0xa6, 0x3e, 0xdc, 0x05, 0x0b,
	0x0d, 0x60, 0x75, 0x2b, 0xfa, 0x65, 0x81, 0x22, 0x3e, 0x40, 0x2b, 0xdf, 0x37, 0x7b, 0x46, 0x04,
	0xb4, 0x95, 0xfd, 0xfd, 0x72, 0x68, 0x9c, 0x1b, 0xec, 0x0c, 0xec, 0x72, 0xa1, 0xb9, 0x20, 0xbd,
	0xfd, 0x1a, 0x61, 0xbd, 0xed, 0x43, 0xed, 0x65, 0x6b, 0x45, 0x6f, 0xda, 0xfa, 0x61, 0x8d, 0xbe,
	0xc4, 0x8b, 0x45, 0x8b, 0x3e, 0x38, 0x9f, 0xff, 0x05, 0x00, 0x00, 0xff, 0xff, 0x13, 0xbb, 0x01,
	0xa9, 0x7d, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MasterRpcClient is the client API for MasterRpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MasterRpcClient interface {
	// // 心跳
	// rpc Ping(Empty) returns(Empty);
	// 注册
	Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterAck, error)
	// 挂载节点
	Mount(ctx context.Context, in *MountReq, opts ...grpc.CallOption) (*Empty, error)
	// 移除节点
	UnMount(ctx context.Context, in *MountReq, opts ...grpc.CallOption) (*Empty, error)
	// 同步挂载点
	Sync(ctx context.Context, in *SyncReq, opts ...grpc.CallOption) (MasterRpc_SyncClient, error)
}

type masterRpcClient struct {
	cc *grpc.ClientConn
}

func NewMasterRpcClient(cc *grpc.ClientConn) MasterRpcClient {
	return &masterRpcClient{cc}
}

func (c *masterRpcClient) Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterAck, error) {
	out := new(RegisterAck)
	err := c.cc.Invoke(ctx, "/pb.MasterRpc/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterRpcClient) Mount(ctx context.Context, in *MountReq, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/pb.MasterRpc/Mount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterRpcClient) UnMount(ctx context.Context, in *MountReq, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/pb.MasterRpc/UnMount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterRpcClient) Sync(ctx context.Context, in *SyncReq, opts ...grpc.CallOption) (MasterRpc_SyncClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MasterRpc_serviceDesc.Streams[0], "/pb.MasterRpc/Sync", opts...)
	if err != nil {
		return nil, err
	}
	x := &masterRpcSyncClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MasterRpc_SyncClient interface {
	Recv() (*MountLogItem, error)
	grpc.ClientStream
}

type masterRpcSyncClient struct {
	grpc.ClientStream
}

func (x *masterRpcSyncClient) Recv() (*MountLogItem, error) {
	m := new(MountLogItem)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MasterRpcServer is the server API for MasterRpc service.
type MasterRpcServer interface {
	// // 心跳
	// rpc Ping(Empty) returns(Empty);
	// 注册
	Register(context.Context, *RegisterReq) (*RegisterAck, error)
	// 挂载节点
	Mount(context.Context, *MountReq) (*Empty, error)
	// 移除节点
	UnMount(context.Context, *MountReq) (*Empty, error)
	// 同步挂载点
	Sync(*SyncReq, MasterRpc_SyncServer) error
}

func RegisterMasterRpcServer(s *grpc.Server, srv MasterRpcServer) {
	s.RegisterService(&_MasterRpc_serviceDesc, srv)
}

func _MasterRpc_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterRpcServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.MasterRpc/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterRpcServer).Register(ctx, req.(*RegisterReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterRpc_Mount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MountReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterRpcServer).Mount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.MasterRpc/Mount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterRpcServer).Mount(ctx, req.(*MountReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterRpc_UnMount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MountReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterRpcServer).UnMount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.MasterRpc/UnMount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterRpcServer).UnMount(ctx, req.(*MountReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterRpc_Sync_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SyncReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MasterRpcServer).Sync(m, &masterRpcSyncServer{stream})
}

type MasterRpc_SyncServer interface {
	Send(*MountLogItem) error
	grpc.ServerStream
}

type masterRpcSyncServer struct {
	grpc.ServerStream
}

func (x *masterRpcSyncServer) Send(m *MountLogItem) error {
	return x.ServerStream.SendMsg(m)
}

var _MasterRpc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.MasterRpc",
	HandlerType: (*MasterRpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _MasterRpc_Register_Handler,
		},
		{
			MethodName: "Mount",
			Handler:    _MasterRpc_Mount_Handler,
		},
		{
			MethodName: "UnMount",
			Handler:    _MasterRpc_UnMount_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Sync",
			Handler:       _MasterRpc_Sync_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "rpc.proto",
}

// PeerRpcClient is the client API for PeerRpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PeerRpcClient interface {
	CheckKey(ctx context.Context, in *CheckKeyReq, opts ...grpc.CallOption) (*Empty, error)
	// 远程调用
	Call(ctx context.Context, in *CallReq, opts ...grpc.CallOption) (*CallAck, error)
	// 双向流抽象
	Stream(ctx context.Context, opts ...grpc.CallOption) (PeerRpc_StreamClient, error)
	// 挂载时间
	MountTime(ctx context.Context, in *NodeName, opts ...grpc.CallOption) (*MountTimeAck, error)
	// 取消挂载
	UnMount(ctx context.Context, in *NodeName, opts ...grpc.CallOption) (*Empty, error)
}

type peerRpcClient struct {
	cc *grpc.ClientConn
}

func NewPeerRpcClient(cc *grpc.ClientConn) PeerRpcClient {
	return &peerRpcClient{cc}
}

func (c *peerRpcClient) CheckKey(ctx context.Context, in *CheckKeyReq, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/pb.PeerRpc/CheckKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *peerRpcClient) Call(ctx context.Context, in *CallReq, opts ...grpc.CallOption) (*CallAck, error) {
	out := new(CallAck)
	err := c.cc.Invoke(ctx, "/pb.PeerRpc/Call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *peerRpcClient) Stream(ctx context.Context, opts ...grpc.CallOption) (PeerRpc_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_PeerRpc_serviceDesc.Streams[0], "/pb.PeerRpc/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &peerRpcStreamClient{stream}
	return x, nil
}

type PeerRpc_StreamClient interface {
	Send(*StreamMsg) error
	Recv() (*StreamMsg, error)
	grpc.ClientStream
}

type peerRpcStreamClient struct {
	grpc.ClientStream
}

func (x *peerRpcStreamClient) Send(m *StreamMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *peerRpcStreamClient) Recv() (*StreamMsg, error) {
	m := new(StreamMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *peerRpcClient) MountTime(ctx context.Context, in *NodeName, opts ...grpc.CallOption) (*MountTimeAck, error) {
	out := new(MountTimeAck)
	err := c.cc.Invoke(ctx, "/pb.PeerRpc/MountTime", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *peerRpcClient) UnMount(ctx context.Context, in *NodeName, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/pb.PeerRpc/UnMount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PeerRpcServer is the server API for PeerRpc service.
type PeerRpcServer interface {
	CheckKey(context.Context, *CheckKeyReq) (*Empty, error)
	// 远程调用
	Call(context.Context, *CallReq) (*CallAck, error)
	// 双向流抽象
	Stream(PeerRpc_StreamServer) error
	// 挂载时间
	MountTime(context.Context, *NodeName) (*MountTimeAck, error)
	// 取消挂载
	UnMount(context.Context, *NodeName) (*Empty, error)
}

func RegisterPeerRpcServer(s *grpc.Server, srv PeerRpcServer) {
	s.RegisterService(&_PeerRpc_serviceDesc, srv)
}

func _PeerRpc_CheckKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckKeyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerRpcServer).CheckKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.PeerRpc/CheckKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerRpcServer).CheckKey(ctx, req.(*CheckKeyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PeerRpc_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CallReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerRpcServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.PeerRpc/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerRpcServer).Call(ctx, req.(*CallReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PeerRpc_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PeerRpcServer).Stream(&peerRpcStreamServer{stream})
}

type PeerRpc_StreamServer interface {
	Send(*StreamMsg) error
	Recv() (*StreamMsg, error)
	grpc.ServerStream
}

type peerRpcStreamServer struct {
	grpc.ServerStream
}

func (x *peerRpcStreamServer) Send(m *StreamMsg) error {
	return x.ServerStream.SendMsg(m)
}

func (x *peerRpcStreamServer) Recv() (*StreamMsg, error) {
	m := new(StreamMsg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _PeerRpc_MountTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerRpcServer).MountTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.PeerRpc/MountTime",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerRpcServer).MountTime(ctx, req.(*NodeName))
	}
	return interceptor(ctx, in, info, handler)
}

func _PeerRpc_UnMount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerRpcServer).UnMount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.PeerRpc/UnMount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerRpcServer).UnMount(ctx, req.(*NodeName))
	}
	return interceptor(ctx, in, info, handler)
}

var _PeerRpc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.PeerRpc",
	HandlerType: (*PeerRpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckKey",
			Handler:    _PeerRpc_CheckKey_Handler,
		},
		{
			MethodName: "Call",
			Handler:    _PeerRpc_Call_Handler,
		},
		{
			MethodName: "MountTime",
			Handler:    _PeerRpc_MountTime_Handler,
		},
		{
			MethodName: "UnMount",
			Handler:    _PeerRpc_UnMount_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _PeerRpc_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "rpc.proto",
}
