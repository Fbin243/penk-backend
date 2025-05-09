// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.21.12
// source: core/task_session_message.proto

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

type TaskSession struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	TaskId        string                 `protobuf:"bytes,2,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	StartTime     int64                  `protobuf:"varint,3,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime       int64                  `protobuf:"varint,4,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	CompletedTime *int64                 `protobuf:"varint,5,opt,name=completed_time,json=completedTime,proto3,oneof" json:"completed_time,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TaskSession) Reset() {
	*x = TaskSession{}
	mi := &file_core_task_session_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TaskSession) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskSession) ProtoMessage() {}

func (x *TaskSession) ProtoReflect() protoreflect.Message {
	mi := &file_core_task_session_message_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskSession.ProtoReflect.Descriptor instead.
func (*TaskSession) Descriptor() ([]byte, []int) {
	return file_core_task_session_message_proto_rawDescGZIP(), []int{0}
}

func (x *TaskSession) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TaskSession) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

func (x *TaskSession) GetStartTime() int64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *TaskSession) GetEndTime() int64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

func (x *TaskSession) GetCompletedTime() int64 {
	if x != nil && x.CompletedTime != nil {
		return *x.CompletedTime
	}
	return 0
}

type TaskSessions struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TaskSessions  []*TaskSession         `protobuf:"bytes,1,rep,name=task_sessions,json=taskSessions,proto3" json:"task_sessions,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TaskSessions) Reset() {
	*x = TaskSessions{}
	mi := &file_core_task_session_message_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TaskSessions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskSessions) ProtoMessage() {}

func (x *TaskSessions) ProtoReflect() protoreflect.Message {
	mi := &file_core_task_session_message_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskSessions.ProtoReflect.Descriptor instead.
func (*TaskSessions) Descriptor() ([]byte, []int) {
	return file_core_task_session_message_proto_rawDescGZIP(), []int{1}
}

func (x *TaskSessions) GetTaskSessions() []*TaskSession {
	if x != nil {
		return x.TaskSessions
	}
	return nil
}

type TaskSessionInput struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            *string                `protobuf:"bytes,1,opt,name=id,proto3,oneof" json:"id,omitempty"`
	TaskId        string                 `protobuf:"bytes,2,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	StartTime     int64                  `protobuf:"varint,3,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime       int64                  `protobuf:"varint,4,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	CompletedTime *int64                 `protobuf:"varint,5,opt,name=completed_time,json=completedTime,proto3,oneof" json:"completed_time,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TaskSessionInput) Reset() {
	*x = TaskSessionInput{}
	mi := &file_core_task_session_message_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TaskSessionInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskSessionInput) ProtoMessage() {}

func (x *TaskSessionInput) ProtoReflect() protoreflect.Message {
	mi := &file_core_task_session_message_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskSessionInput.ProtoReflect.Descriptor instead.
func (*TaskSessionInput) Descriptor() ([]byte, []int) {
	return file_core_task_session_message_proto_rawDescGZIP(), []int{2}
}

func (x *TaskSessionInput) GetId() string {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return ""
}

func (x *TaskSessionInput) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

func (x *TaskSessionInput) GetStartTime() int64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *TaskSessionInput) GetEndTime() int64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

func (x *TaskSessionInput) GetCompletedTime() int64 {
	if x != nil && x.CompletedTime != nil {
		return *x.CompletedTime
	}
	return 0
}

type TaskSessionInputs struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	TaskSessionInputs []*TaskSessionInput    `protobuf:"bytes,1,rep,name=task_session_inputs,json=taskSessionInputs,proto3" json:"task_session_inputs,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *TaskSessionInputs) Reset() {
	*x = TaskSessionInputs{}
	mi := &file_core_task_session_message_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TaskSessionInputs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskSessionInputs) ProtoMessage() {}

func (x *TaskSessionInputs) ProtoReflect() protoreflect.Message {
	mi := &file_core_task_session_message_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskSessionInputs.ProtoReflect.Descriptor instead.
func (*TaskSessionInputs) Descriptor() ([]byte, []int) {
	return file_core_task_session_message_proto_rawDescGZIP(), []int{3}
}

func (x *TaskSessionInputs) GetTaskSessionInputs() []*TaskSessionInput {
	if x != nil {
		return x.TaskSessionInputs
	}
	return nil
}

var File_core_task_session_message_proto protoreflect.FileDescriptor

const file_core_task_session_message_proto_rawDesc = "" +
	"\n" +
	"\x1fcore/task_session_message.proto\x12\x04core\"\xaf\x01\n" +
	"\vTaskSession\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x17\n" +
	"\atask_id\x18\x02 \x01(\tR\x06taskId\x12\x1d\n" +
	"\n" +
	"start_time\x18\x03 \x01(\x03R\tstartTime\x12\x19\n" +
	"\bend_time\x18\x04 \x01(\x03R\aendTime\x12*\n" +
	"\x0ecompleted_time\x18\x05 \x01(\x03H\x00R\rcompletedTime\x88\x01\x01B\x11\n" +
	"\x0f_completed_time\"F\n" +
	"\fTaskSessions\x126\n" +
	"\rtask_sessions\x18\x01 \x03(\v2\x11.core.TaskSessionR\ftaskSessions\"\xc0\x01\n" +
	"\x10TaskSessionInput\x12\x13\n" +
	"\x02id\x18\x01 \x01(\tH\x00R\x02id\x88\x01\x01\x12\x17\n" +
	"\atask_id\x18\x02 \x01(\tR\x06taskId\x12\x1d\n" +
	"\n" +
	"start_time\x18\x03 \x01(\x03R\tstartTime\x12\x19\n" +
	"\bend_time\x18\x04 \x01(\x03R\aendTime\x12*\n" +
	"\x0ecompleted_time\x18\x05 \x01(\x03H\x01R\rcompletedTime\x88\x01\x01B\x05\n" +
	"\x03_idB\x11\n" +
	"\x0f_completed_time\"[\n" +
	"\x11TaskSessionInputs\x12F\n" +
	"\x13task_session_inputs\x18\x01 \x03(\v2\x16.core.TaskSessionInputR\x11taskSessionInputsB\x19Z\x17tenkhours/proto/pb/coreb\x06proto3"

var (
	file_core_task_session_message_proto_rawDescOnce sync.Once
	file_core_task_session_message_proto_rawDescData []byte
)

func file_core_task_session_message_proto_rawDescGZIP() []byte {
	file_core_task_session_message_proto_rawDescOnce.Do(func() {
		file_core_task_session_message_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_core_task_session_message_proto_rawDesc), len(file_core_task_session_message_proto_rawDesc)))
	})
	return file_core_task_session_message_proto_rawDescData
}

var file_core_task_session_message_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_core_task_session_message_proto_goTypes = []any{
	(*TaskSession)(nil),       // 0: core.TaskSession
	(*TaskSessions)(nil),      // 1: core.TaskSessions
	(*TaskSessionInput)(nil),  // 2: core.TaskSessionInput
	(*TaskSessionInputs)(nil), // 3: core.TaskSessionInputs
}
var file_core_task_session_message_proto_depIdxs = []int32{
	0, // 0: core.TaskSessions.task_sessions:type_name -> core.TaskSession
	2, // 1: core.TaskSessionInputs.task_session_inputs:type_name -> core.TaskSessionInput
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_core_task_session_message_proto_init() }
func file_core_task_session_message_proto_init() {
	if File_core_task_session_message_proto != nil {
		return
	}
	file_core_task_session_message_proto_msgTypes[0].OneofWrappers = []any{}
	file_core_task_session_message_proto_msgTypes[2].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_core_task_session_message_proto_rawDesc), len(file_core_task_session_message_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_core_task_session_message_proto_goTypes,
		DependencyIndexes: file_core_task_session_message_proto_depIdxs,
		MessageInfos:      file_core_task_session_message_proto_msgTypes,
	}.Build()
	File_core_task_session_message_proto = out.File
	file_core_task_session_message_proto_goTypes = nil
	file_core_task_session_message_proto_depIdxs = nil
}
