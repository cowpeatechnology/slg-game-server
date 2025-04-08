// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v4.25.3
// source: proto/message.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 消息类型枚举
type MessageType int32

const (
	MessageType_UNKNOWN   MessageType = 0
	MessageType_HEARTBEAT MessageType = 1
	MessageType_ERROR     MessageType = 2
	MessageType_TEXT      MessageType = 3 // 文本消息类型
)

// Enum value maps for MessageType.
var (
	MessageType_name = map[int32]string{
		0: "UNKNOWN",
		1: "HEARTBEAT",
		2: "ERROR",
		3: "TEXT",
	}
	MessageType_value = map[string]int32{
		"UNKNOWN":   0,
		"HEARTBEAT": 1,
		"ERROR":     2,
		"TEXT":      3,
	}
)

func (x MessageType) Enum() *MessageType {
	p := new(MessageType)
	*p = x
	return p
}

func (x MessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_message_proto_enumTypes[0].Descriptor()
}

func (MessageType) Type() protoreflect.EnumType {
	return &file_proto_message_proto_enumTypes[0]
}

func (x MessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageType.Descriptor instead.
func (MessageType) EnumDescriptor() ([]byte, []int) {
	return file_proto_message_proto_rawDescGZIP(), []int{0}
}

// 基础消息结构
type Message struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Type          MessageType            `protobuf:"varint,1,opt,name=type,proto3,enum=proto.MessageType" json:"type,omitempty"`
	Payload       []byte                 `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Message) Reset() {
	*x = Message{}
	mi := &file_proto_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_proto_message_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_proto_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetType() MessageType {
	if x != nil {
		return x.Type
	}
	return MessageType_UNKNOWN
}

func (x *Message) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

// 心跳消息
type Heartbeat struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Timestamp     int64                  `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Heartbeat) Reset() {
	*x = Heartbeat{}
	mi := &file_proto_message_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Heartbeat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Heartbeat) ProtoMessage() {}

func (x *Heartbeat) ProtoReflect() protoreflect.Message {
	mi := &file_proto_message_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Heartbeat.ProtoReflect.Descriptor instead.
func (*Heartbeat) Descriptor() ([]byte, []int) {
	return file_proto_message_proto_rawDescGZIP(), []int{1}
}

func (x *Heartbeat) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

// 错误消息
type Error struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          int32                  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Error) Reset() {
	*x = Error{}
	mi := &file_proto_message_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_proto_message_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_proto_message_proto_rawDescGZIP(), []int{2}
}

func (x *Error) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Error) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// 文本消息
type TextMessage struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Content       string                 `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`      // 消息内容
	Timestamp     int64                  `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"` // 时间戳
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TextMessage) Reset() {
	*x = TextMessage{}
	mi := &file_proto_message_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TextMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TextMessage) ProtoMessage() {}

func (x *TextMessage) ProtoReflect() protoreflect.Message {
	mi := &file_proto_message_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TextMessage.ProtoReflect.Descriptor instead.
func (*TextMessage) Descriptor() ([]byte, []int) {
	return file_proto_message_proto_rawDescGZIP(), []int{3}
}

func (x *TextMessage) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *TextMessage) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

var File_proto_message_proto protoreflect.FileDescriptor

const file_proto_message_proto_rawDesc = "" +
	"\n" +
	"\x13proto/message.proto\x12\x05proto\"K\n" +
	"\aMessage\x12&\n" +
	"\x04type\x18\x01 \x01(\x0e2\x12.proto.MessageTypeR\x04type\x12\x18\n" +
	"\apayload\x18\x02 \x01(\fR\apayload\")\n" +
	"\tHeartbeat\x12\x1c\n" +
	"\ttimestamp\x18\x01 \x01(\x03R\ttimestamp\"5\n" +
	"\x05Error\x12\x12\n" +
	"\x04code\x18\x01 \x01(\x05R\x04code\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\"E\n" +
	"\vTextMessage\x12\x18\n" +
	"\acontent\x18\x01 \x01(\tR\acontent\x12\x1c\n" +
	"\ttimestamp\x18\x02 \x01(\x03R\ttimestamp*>\n" +
	"\vMessageType\x12\v\n" +
	"\aUNKNOWN\x10\x00\x12\r\n" +
	"\tHEARTBEAT\x10\x01\x12\t\n" +
	"\x05ERROR\x10\x02\x12\b\n" +
	"\x04TEXT\x10\x03B\x17Z\x15slg-game-server/protob\x06proto3"

var (
	file_proto_message_proto_rawDescOnce sync.Once
	file_proto_message_proto_rawDescData []byte
)

func file_proto_message_proto_rawDescGZIP() []byte {
	file_proto_message_proto_rawDescOnce.Do(func() {
		file_proto_message_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_message_proto_rawDesc), len(file_proto_message_proto_rawDesc)))
	})
	return file_proto_message_proto_rawDescData
}

var file_proto_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_message_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_message_proto_goTypes = []any{
	(MessageType)(0),    // 0: proto.MessageType
	(*Message)(nil),     // 1: proto.Message
	(*Heartbeat)(nil),   // 2: proto.Heartbeat
	(*Error)(nil),       // 3: proto.Error
	(*TextMessage)(nil), // 4: proto.TextMessage
}
var file_proto_message_proto_depIdxs = []int32{
	0, // 0: proto.Message.type:type_name -> proto.MessageType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_message_proto_init() }
func file_proto_message_proto_init() {
	if File_proto_message_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_message_proto_rawDesc), len(file_proto_message_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_message_proto_goTypes,
		DependencyIndexes: file_proto_message_proto_depIdxs,
		EnumInfos:         file_proto_message_proto_enumTypes,
		MessageInfos:      file_proto_message_proto_msgTypes,
	}.Build()
	File_proto_message_proto = out.File
	file_proto_message_proto_goTypes = nil
	file_proto_message_proto_depIdxs = nil
}
