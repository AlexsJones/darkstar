// Code generated by protoc-gen-go. DO NOT EDIT.
// source: data/message.proto

/*
Package data is a generated protocol buffer package.

It is generated from these files:
	data/message.proto
	data/uuid.proto

It has these top-level messages:
	Message
	UUID
*/
package data

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Message struct {
	Uuid *UUID `protobuf:"bytes,1,opt,name=uuid" json:"uuid,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Message) GetUuid() *UUID {
	if m != nil {
		return m.Uuid
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "data.Message")
}

func init() { proto.RegisterFile("data/message.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 94 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4a, 0x49, 0x2c, 0x49,
	0xd4, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0x01, 0x89, 0x49, 0xf1, 0x83, 0x65, 0x4a, 0x4b, 0x33, 0x53, 0x20, 0xc2, 0x4a, 0x9a, 0x5c, 0xec,
	0xbe, 0x10, 0x75, 0x42, 0x72, 0x5c, 0x2c, 0x20, 0x09, 0x09, 0x46, 0x05, 0x46, 0x0d, 0x6e, 0x23,
	0x2e, 0x3d, 0x90, 0x52, 0xbd, 0xd0, 0x50, 0x4f, 0x97, 0x20, 0xb0, 0x78, 0x12, 0x1b, 0x58, 0x87,
	0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x5a, 0x79, 0x18, 0x6c, 0x5e, 0x00, 0x00, 0x00,
}
