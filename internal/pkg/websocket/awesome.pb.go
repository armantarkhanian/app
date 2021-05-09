// Package websocket ...
package websocket

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LikeObject struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Exist bool   `protobuf:"varint,2,opt,name=exist,proto3" json:"exist,omitempty"`
}

func (x *LikeObject) Reset() {
	*x = LikeObject{}
	if protoimpl.UnsafeEnabled {
		mi := &file_awesome_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LikeObject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LikeObject) ProtoMessage() {}

func (x *LikeObject) ProtoReflect() protoreflect.Message {
	mi := &file_awesome_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LikeObject.ProtoReflect.Descriptor instead.
func (*LikeObject) Descriptor() ([]byte, []int) {
	return file_awesome_proto_rawDescGZIP(), []int{0}
}

func (x *LikeObject) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *LikeObject) GetExist() bool {
	if x != nil {
		return x.Exist
	}
	return false
}

type LikeMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Objects []*LikeObject `protobuf:"bytes,1,rep,name=objects,proto3" json:"objects,omitempty"`
}

func (x *LikeMessage) Reset() {
	*x = LikeMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_awesome_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LikeMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LikeMessage) ProtoMessage() {}

func (x *LikeMessage) ProtoReflect() protoreflect.Message {
	mi := &file_awesome_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LikeMessage.ProtoReflect.Descriptor instead.
func (*LikeMessage) Descriptor() ([]byte, []int) {
	return file_awesome_proto_rawDescGZIP(), []int{1}
}

func (x *LikeMessage) GetObjects() []*LikeObject {
	if x != nil {
		return x.Objects
	}
	return nil
}

var File_awesome_proto protoreflect.FileDescriptor

var file_awesome_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x61, 0x77, 0x65, 0x73, 0x6f, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x09, 0x77, 0x65, 0x62, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x32, 0x0a, 0x0a, 0x4c, 0x69,
	0x6b, 0x65, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x78, 0x69, 0x73,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x65, 0x78, 0x69, 0x73, 0x74, 0x22, 0x3e,
	0x0a, 0x0b, 0x4c, 0x69, 0x6b, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2f, 0x0a,
	0x07, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15,
	0x2e, 0x77, 0x65, 0x62, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x4c, 0x69, 0x6b, 0x65, 0x4f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x07, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x42, 0x0c,
	0x5a, 0x0a, 0x2f, 0x77, 0x65, 0x62, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_awesome_proto_rawDescOnce sync.Once
	file_awesome_proto_rawDescData = file_awesome_proto_rawDesc
)

func file_awesome_proto_rawDescGZIP() []byte {
	file_awesome_proto_rawDescOnce.Do(func() {
		file_awesome_proto_rawDescData = protoimpl.X.CompressGZIP(file_awesome_proto_rawDescData)
	})
	return file_awesome_proto_rawDescData
}

var file_awesome_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_awesome_proto_goTypes = []interface{}{
	(*LikeObject)(nil),  // 0: websocket.LikeObject
	(*LikeMessage)(nil), // 1: websocket.LikeMessage
}
var file_awesome_proto_depIdxs = []int32{
	0, // 0: websocket.LikeMessage.objects:type_name -> websocket.LikeObject
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_awesome_proto_init() }
func file_awesome_proto_init() {
	if File_awesome_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_awesome_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LikeObject); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_awesome_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LikeMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_awesome_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_awesome_proto_goTypes,
		DependencyIndexes: file_awesome_proto_depIdxs,
		MessageInfos:      file_awesome_proto_msgTypes,
	}.Build()
	File_awesome_proto = out.File
	file_awesome_proto_rawDesc = nil
	file_awesome_proto_goTypes = nil
	file_awesome_proto_depIdxs = nil
}
