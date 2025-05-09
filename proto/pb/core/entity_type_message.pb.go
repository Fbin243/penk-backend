// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.21.12
// source: core/entity_type_message.proto

package core

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

type EntityType int32

const (
	EntityType_TaskType  EntityType = 0
	EntityType_HabitType EntityType = 1
)

// Enum value maps for EntityType.
var (
	EntityType_name = map[int32]string{
		0: "TaskType",
		1: "HabitType",
	}
	EntityType_value = map[string]int32{
		"TaskType":  0,
		"HabitType": 1,
	}
)

func (x EntityType) Enum() *EntityType {
	p := new(EntityType)
	*p = x
	return p
}

func (x EntityType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EntityType) Descriptor() protoreflect.EnumDescriptor {
	return file_core_entity_type_message_proto_enumTypes[0].Descriptor()
}

func (EntityType) Type() protoreflect.EnumType {
	return &file_core_entity_type_message_proto_enumTypes[0]
}

func (x EntityType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EntityType.Descriptor instead.
func (EntityType) EnumDescriptor() ([]byte, []int) {
	return file_core_entity_type_message_proto_rawDescGZIP(), []int{0}
}

var File_core_entity_type_message_proto protoreflect.FileDescriptor

const file_core_entity_type_message_proto_rawDesc = "" +
	"\n" +
	"\x1ecore/entity_type_message.proto\x12\x04core*)\n" +
	"\n" +
	"EntityType\x12\f\n" +
	"\bTaskType\x10\x00\x12\r\n" +
	"\tHabitType\x10\x01B\x19Z\x17tenkhours/proto/pb/coreb\x06proto3"

var (
	file_core_entity_type_message_proto_rawDescOnce sync.Once
	file_core_entity_type_message_proto_rawDescData []byte
)

func file_core_entity_type_message_proto_rawDescGZIP() []byte {
	file_core_entity_type_message_proto_rawDescOnce.Do(func() {
		file_core_entity_type_message_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_core_entity_type_message_proto_rawDesc), len(file_core_entity_type_message_proto_rawDesc)))
	})
	return file_core_entity_type_message_proto_rawDescData
}

var file_core_entity_type_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_core_entity_type_message_proto_goTypes = []any{
	(EntityType)(0), // 0: core.EntityType
}
var file_core_entity_type_message_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_core_entity_type_message_proto_init() }
func file_core_entity_type_message_proto_init() {
	if File_core_entity_type_message_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_core_entity_type_message_proto_rawDesc), len(file_core_entity_type_message_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_core_entity_type_message_proto_goTypes,
		DependencyIndexes: file_core_entity_type_message_proto_depIdxs,
		EnumInfos:         file_core_entity_type_message_proto_enumTypes,
	}.Build()
	File_core_entity_type_message_proto = out.File
	file_core_entity_type_message_proto_goTypes = nil
	file_core_entity_type_message_proto_depIdxs = nil
}
