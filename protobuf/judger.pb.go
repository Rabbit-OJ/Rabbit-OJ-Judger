// Code generated by protoc-gen-go. DO NOT EDIT.
// source: services/judger/protobuf/judger.proto

package protobuf

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type JudgeRequest struct {
	Sid                  uint32   `protobuf:"varint,1,opt,name=sid,proto3" json:"sid,omitempty"`
	Tid                  uint32   `protobuf:"varint,2,opt,name=tid,proto3" json:"tid,omitempty"`
	Version              string   `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	Language             string   `protobuf:"bytes,4,opt,name=language,proto3" json:"language,omitempty"`
	TimeLimit            uint32   `protobuf:"varint,5,opt,name=time_limit,json=timeLimit,proto3" json:"time_limit,omitempty"`
	SpaceLimit           uint32   `protobuf:"varint,6,opt,name=space_limit,json=spaceLimit,proto3" json:"space_limit,omitempty"`
	CompMode             string   `protobuf:"bytes,7,opt,name=comp_mode,json=compMode,proto3" json:"comp_mode,omitempty"`
	Code                 []byte   `protobuf:"bytes,8,opt,name=code,proto3" json:"code,omitempty"`
	Time                 int64    `protobuf:"varint,9,opt,name=time,proto3" json:"time,omitempty"`
	IsContest            bool     `protobuf:"varint,10,opt,name=is_contest,json=isContest,proto3" json:"is_contest,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JudgeRequest) Reset()         { *m = JudgeRequest{} }
func (m *JudgeRequest) String() string { return proto.CompactTextString(m) }
func (*JudgeRequest) ProtoMessage()    {}
func (*JudgeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fffdb9f2f3cf657, []int{0}
}

func (m *JudgeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JudgeRequest.Unmarshal(m, b)
}
func (m *JudgeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JudgeRequest.Marshal(b, m, deterministic)
}
func (m *JudgeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JudgeRequest.Merge(m, src)
}
func (m *JudgeRequest) XXX_Size() int {
	return xxx_messageInfo_JudgeRequest.Size(m)
}
func (m *JudgeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_JudgeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_JudgeRequest proto.InternalMessageInfo

func (m *JudgeRequest) GetSid() uint32 {
	if m != nil {
		return m.Sid
	}
	return 0
}

func (m *JudgeRequest) GetTid() uint32 {
	if m != nil {
		return m.Tid
	}
	return 0
}

func (m *JudgeRequest) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *JudgeRequest) GetLanguage() string {
	if m != nil {
		return m.Language
	}
	return ""
}

func (m *JudgeRequest) GetTimeLimit() uint32 {
	if m != nil {
		return m.TimeLimit
	}
	return 0
}

func (m *JudgeRequest) GetSpaceLimit() uint32 {
	if m != nil {
		return m.SpaceLimit
	}
	return 0
}

func (m *JudgeRequest) GetCompMode() string {
	if m != nil {
		return m.CompMode
	}
	return ""
}

func (m *JudgeRequest) GetCode() []byte {
	if m != nil {
		return m.Code
	}
	return nil
}

func (m *JudgeRequest) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *JudgeRequest) GetIsContest() bool {
	if m != nil {
		return m.IsContest
	}
	return false
}

type JudgeCaseResult struct {
	Status               string   `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	SpaceUsed            uint32   `protobuf:"varint,2,opt,name=space_used,json=spaceUsed,proto3" json:"space_used,omitempty"`
	TimeUsed             uint32   `protobuf:"varint,3,opt,name=time_used,json=timeUsed,proto3" json:"time_used,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JudgeCaseResult) Reset()         { *m = JudgeCaseResult{} }
func (m *JudgeCaseResult) String() string { return proto.CompactTextString(m) }
func (*JudgeCaseResult) ProtoMessage()    {}
func (*JudgeCaseResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fffdb9f2f3cf657, []int{1}
}

func (m *JudgeCaseResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JudgeCaseResult.Unmarshal(m, b)
}
func (m *JudgeCaseResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JudgeCaseResult.Marshal(b, m, deterministic)
}
func (m *JudgeCaseResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JudgeCaseResult.Merge(m, src)
}
func (m *JudgeCaseResult) XXX_Size() int {
	return xxx_messageInfo_JudgeCaseResult.Size(m)
}
func (m *JudgeCaseResult) XXX_DiscardUnknown() {
	xxx_messageInfo_JudgeCaseResult.DiscardUnknown(m)
}

var xxx_messageInfo_JudgeCaseResult proto.InternalMessageInfo

func (m *JudgeCaseResult) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *JudgeCaseResult) GetSpaceUsed() uint32 {
	if m != nil {
		return m.SpaceUsed
	}
	return 0
}

func (m *JudgeCaseResult) GetTimeUsed() uint32 {
	if m != nil {
		return m.TimeUsed
	}
	return 0
}

type JudgeResponse struct {
	Sid                  uint32             `protobuf:"varint,1,opt,name=sid,proto3" json:"sid,omitempty"`
	IsContest            bool               `protobuf:"varint,2,opt,name=is_contest,json=isContest,proto3" json:"is_contest,omitempty"`
	Result               []*JudgeCaseResult `protobuf:"bytes,3,rep,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *JudgeResponse) Reset()         { *m = JudgeResponse{} }
func (m *JudgeResponse) String() string { return proto.CompactTextString(m) }
func (*JudgeResponse) ProtoMessage()    {}
func (*JudgeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fffdb9f2f3cf657, []int{2}
}

func (m *JudgeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JudgeResponse.Unmarshal(m, b)
}
func (m *JudgeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JudgeResponse.Marshal(b, m, deterministic)
}
func (m *JudgeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JudgeResponse.Merge(m, src)
}
func (m *JudgeResponse) XXX_Size() int {
	return xxx_messageInfo_JudgeResponse.Size(m)
}
func (m *JudgeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_JudgeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_JudgeResponse proto.InternalMessageInfo

func (m *JudgeResponse) GetSid() uint32 {
	if m != nil {
		return m.Sid
	}
	return 0
}

func (m *JudgeResponse) GetIsContest() bool {
	if m != nil {
		return m.IsContest
	}
	return false
}

func (m *JudgeResponse) GetResult() []*JudgeCaseResult {
	if m != nil {
		return m.Result
	}
	return nil
}

func init() {
	proto.RegisterType((*JudgeRequest)(nil), "protobuf.JudgeRequest")
	proto.RegisterType((*JudgeCaseResult)(nil), "protobuf.JudgeCaseResult")
	proto.RegisterType((*JudgeResponse)(nil), "protobuf.JudgeResponse")
}

func init() {
	proto.RegisterFile("services/judger/protobuf/judger.proto", fileDescriptor_2fffdb9f2f3cf657)
}

var fileDescriptor_2fffdb9f2f3cf657 = []byte{
	// 356 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0x51, 0x4b, 0xeb, 0x40,
	0x10, 0x85, 0x49, 0xd3, 0xa6, 0xc9, 0xb4, 0xbd, 0xf7, 0xb2, 0x0f, 0xd7, 0xb5, 0x22, 0x86, 0x82,
	0x90, 0xa7, 0x16, 0xeb, 0x2f, 0x90, 0x3e, 0x08, 0x45, 0x5f, 0x16, 0x7c, 0x2e, 0x69, 0x32, 0x96,
	0x95, 0x26, 0x1b, 0x33, 0x9b, 0xfe, 0x0a, 0x7f, 0xb4, 0xec, 0x24, 0x11, 0xac, 0xbe, 0xcd, 0xf9,
	0xce, 0x76, 0x38, 0x73, 0x1a, 0xb8, 0x25, 0xac, 0x4f, 0x3a, 0x43, 0x5a, 0xbd, 0x35, 0xf9, 0x01,
	0xeb, 0x55, 0x55, 0x1b, 0x6b, 0xf6, 0xcd, 0x6b, 0xa7, 0x97, 0xac, 0x45, 0xd8, 0xe3, 0xc5, 0xc7,
	0x00, 0xa6, 0x5b, 0x67, 0x29, 0x7c, 0x6f, 0x90, 0xac, 0xf8, 0x07, 0x3e, 0xe9, 0x5c, 0x7a, 0xb1,
	0x97, 0xcc, 0x94, 0x1b, 0x1d, 0xb1, 0x3a, 0x97, 0x83, 0x96, 0x58, 0x9d, 0x0b, 0x09, 0xe3, 0x13,
	0xd6, 0xa4, 0x4d, 0x29, 0xfd, 0xd8, 0x4b, 0x22, 0xd5, 0x4b, 0x31, 0x87, 0xf0, 0x98, 0x96, 0x87,
	0x26, 0x3d, 0xa0, 0x1c, 0xb2, 0xf5, 0xa5, 0xc5, 0x35, 0x80, 0xd5, 0x05, 0xee, 0x8e, 0xba, 0xd0,
	0x56, 0x8e, 0x78, 0x5d, 0xe4, 0xc8, 0x93, 0x03, 0xe2, 0x06, 0x26, 0x54, 0xa5, 0x59, 0xef, 0x07,
	0xec, 0x03, 0xa3, 0xf6, 0xc1, 0x15, 0x44, 0x99, 0x29, 0xaa, 0x5d, 0x61, 0x72, 0x94, 0xe3, 0x76,
	0xb9, 0x03, 0xcf, 0x26, 0x47, 0x21, 0x60, 0x98, 0x39, 0x1e, 0xc6, 0x5e, 0x32, 0x55, 0x3c, 0x3b,
	0xe6, 0xd6, 0xcb, 0x28, 0xf6, 0x12, 0x5f, 0xf1, 0xec, 0x42, 0x68, 0xda, 0x65, 0xa6, 0xb4, 0x48,
	0x56, 0x42, 0xec, 0x25, 0xa1, 0x8a, 0x34, 0x6d, 0x5a, 0xb0, 0x40, 0xf8, 0xcb, 0x6d, 0x6c, 0x52,
	0x42, 0x85, 0xd4, 0x1c, 0xad, 0xf8, 0x0f, 0x01, 0xd9, 0xd4, 0x36, 0xc4, 0x9d, 0x44, 0xaa, 0x53,
	0x6e, 0x53, 0x9b, 0xb7, 0x21, 0xec, 0xdb, 0x89, 0x98, 0xbc, 0x10, 0xe6, 0x2e, 0x2d, 0x5f, 0xcb,
	0xae, 0xcf, 0x6e, 0xe8, 0x80, 0x33, 0x17, 0x04, 0xb3, 0xae, 0x74, 0xaa, 0x4c, 0x49, 0xf8, 0x4b,
	0xeb, 0xdf, 0x83, 0x0e, 0xce, 0x82, 0x8a, 0x3b, 0x08, 0x6a, 0xce, 0x27, 0xfd, 0xd8, 0x4f, 0x26,
	0xeb, 0xcb, 0x65, 0xff, 0x97, 0x2e, 0xcf, 0x0e, 0x50, 0xdd, 0xc3, 0xf5, 0x16, 0x46, 0x6c, 0x89,
	0x07, 0xf8, 0xf3, 0x88, 0xb6, 0x0f, 0xc0, 0x37, 0x9e, 0xfd, 0xba, 0xfb, 0x18, 0xe6, 0x17, 0x3f,
	0x78, 0x9b, 0x77, 0x1f, 0x30, 0xbf, 0xff, 0x0c, 0x00, 0x00, 0xff, 0xff, 0x79, 0x0e, 0x57, 0x1b,
	0x70, 0x02, 0x00, 0x00,
}
