// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        v3.21.12
// source: timetracking/timetracking_message.proto

package timetracking

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReferenceType int32

const (
	ReferenceType_Habit ReferenceType = 0
	ReferenceType_Task  ReferenceType = 1
)

// Enum value maps for ReferenceType.
var (
	ReferenceType_name = map[int32]string{
		0: "Habit",
		1: "Task",
	}
	ReferenceType_value = map[string]int32{
		"Habit": 0,
		"Task":  1,
	}
)

func (x ReferenceType) Enum() *ReferenceType {
	p := new(ReferenceType)
	*p = x
	return p
}

func (x ReferenceType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ReferenceType) Descriptor() protoreflect.EnumDescriptor {
	return file_timetracking_timetracking_message_proto_enumTypes[0].Descriptor()
}

func (ReferenceType) Type() protoreflect.EnumType {
	return &file_timetracking_timetracking_message_proto_enumTypes[0]
}

func (x ReferenceType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ReferenceType.Descriptor instead.
func (ReferenceType) EnumDescriptor() ([]byte, []int) {
	return file_timetracking_timetracking_message_proto_rawDescGZIP(), []int{0}
}

type TimeTracking struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CharacterId   string                 `protobuf:"bytes,2,opt,name=character_id,json=characterId,proto3" json:"character_id,omitempty"`
	CategoryId    *string                `protobuf:"bytes,3,opt,name=category_id,json=categoryId,proto3,oneof" json:"category_id,omitempty"`
	ReferenceId   *string                `protobuf:"bytes,4,opt,name=reference_id,json=referenceId,proto3,oneof" json:"reference_id,omitempty"`
	ReferenceType *ReferenceType         `protobuf:"varint,5,opt,name=reference_type,json=referenceType,proto3,enum=timetracking.ReferenceType,oneof" json:"reference_type,omitempty"`
	StartTime     int64                  `protobuf:"varint,6,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime       *int64                 `protobuf:"varint,7,opt,name=end_time,json=endTime,proto3,oneof" json:"end_time,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TimeTracking) Reset() {
	*x = TimeTracking{}
	mi := &file_timetracking_timetracking_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TimeTracking) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeTracking) ProtoMessage() {}

func (x *TimeTracking) ProtoReflect() protoreflect.Message {
	mi := &file_timetracking_timetracking_message_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeTracking.ProtoReflect.Descriptor instead.
func (*TimeTracking) Descriptor() ([]byte, []int) {
	return file_timetracking_timetracking_message_proto_rawDescGZIP(), []int{0}
}

func (x *TimeTracking) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TimeTracking) GetCharacterId() string {
	if x != nil {
		return x.CharacterId
	}
	return ""
}

func (x *TimeTracking) GetCategoryId() string {
	if x != nil && x.CategoryId != nil {
		return *x.CategoryId
	}
	return ""
}

func (x *TimeTracking) GetReferenceId() string {
	if x != nil && x.ReferenceId != nil {
		return *x.ReferenceId
	}
	return ""
}

func (x *TimeTracking) GetReferenceType() ReferenceType {
	if x != nil && x.ReferenceType != nil {
		return *x.ReferenceType
	}
	return ReferenceType_Habit
}

func (x *TimeTracking) GetStartTime() int64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *TimeTracking) GetEndTime() int64 {
	if x != nil && x.EndTime != nil {
		return *x.EndTime
	}
	return 0
}

type TimeTrackingWithFish struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TimeTracking  *TimeTracking          `protobuf:"bytes,1,opt,name=time_tracking,json=timeTracking,proto3" json:"time_tracking,omitempty"`
	Normal        int32                  `protobuf:"varint,2,opt,name=normal,proto3" json:"normal,omitempty"`
	Gold          int32                  `protobuf:"varint,3,opt,name=gold,proto3" json:"gold,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TimeTrackingWithFish) Reset() {
	*x = TimeTrackingWithFish{}
	mi := &file_timetracking_timetracking_message_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TimeTrackingWithFish) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeTrackingWithFish) ProtoMessage() {}

func (x *TimeTrackingWithFish) ProtoReflect() protoreflect.Message {
	mi := &file_timetracking_timetracking_message_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeTrackingWithFish.ProtoReflect.Descriptor instead.
func (*TimeTrackingWithFish) Descriptor() ([]byte, []int) {
	return file_timetracking_timetracking_message_proto_rawDescGZIP(), []int{1}
}

func (x *TimeTrackingWithFish) GetTimeTracking() *TimeTracking {
	if x != nil {
		return x.TimeTracking
	}
	return nil
}

func (x *TimeTrackingWithFish) GetNormal() int32 {
	if x != nil {
		return x.Normal
	}
	return 0
}

func (x *TimeTrackingWithFish) GetGold() int32 {
	if x != nil {
		return x.Gold
	}
	return 0
}

type TotalTimeTrackingReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CharacterId   string                 `protobuf:"bytes,1,opt,name=character_id,json=characterId,proto3" json:"character_id,omitempty"`
	Timestamp     int64                  `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TotalTimeTrackingReq) Reset() {
	*x = TotalTimeTrackingReq{}
	mi := &file_timetracking_timetracking_message_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TotalTimeTrackingReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TotalTimeTrackingReq) ProtoMessage() {}

func (x *TotalTimeTrackingReq) ProtoReflect() protoreflect.Message {
	mi := &file_timetracking_timetracking_message_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TotalTimeTrackingReq.ProtoReflect.Descriptor instead.
func (*TotalTimeTrackingReq) Descriptor() ([]byte, []int) {
	return file_timetracking_timetracking_message_proto_rawDescGZIP(), []int{2}
}

func (x *TotalTimeTrackingReq) GetCharacterId() string {
	if x != nil {
		return x.CharacterId
	}
	return ""
}

func (x *TotalTimeTrackingReq) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type TotalTimeTrackingResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TotalTime     int32                  `protobuf:"varint,1,opt,name=total_time,json=totalTime,proto3" json:"total_time,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TotalTimeTrackingResp) Reset() {
	*x = TotalTimeTrackingResp{}
	mi := &file_timetracking_timetracking_message_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TotalTimeTrackingResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TotalTimeTrackingResp) ProtoMessage() {}

func (x *TotalTimeTrackingResp) ProtoReflect() protoreflect.Message {
	mi := &file_timetracking_timetracking_message_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TotalTimeTrackingResp.ProtoReflect.Descriptor instead.
func (*TotalTimeTrackingResp) Descriptor() ([]byte, []int) {
	return file_timetracking_timetracking_message_proto_rawDescGZIP(), []int{3}
}

func (x *TotalTimeTrackingResp) GetTotalTime() int32 {
	if x != nil {
		return x.TotalTime
	}
	return 0
}

type CreateTimeTrackingReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CharacterId   string                 `protobuf:"bytes,1,opt,name=character_id,json=characterId,proto3" json:"character_id,omitempty"`
	CategoryId    *string                `protobuf:"bytes,2,opt,name=category_id,json=categoryId,proto3,oneof" json:"category_id,omitempty"`
	ReferenceId   *string                `protobuf:"bytes,3,opt,name=reference_id,json=referenceId,proto3,oneof" json:"reference_id,omitempty"`
	ReferenceType *ReferenceType         `protobuf:"varint,4,opt,name=reference_type,json=referenceType,proto3,enum=timetracking.ReferenceType,oneof" json:"reference_type,omitempty"`
	StartTime     int64                  `protobuf:"varint,5,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateTimeTrackingReq) Reset() {
	*x = CreateTimeTrackingReq{}
	mi := &file_timetracking_timetracking_message_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateTimeTrackingReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTimeTrackingReq) ProtoMessage() {}

func (x *CreateTimeTrackingReq) ProtoReflect() protoreflect.Message {
	mi := &file_timetracking_timetracking_message_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTimeTrackingReq.ProtoReflect.Descriptor instead.
func (*CreateTimeTrackingReq) Descriptor() ([]byte, []int) {
	return file_timetracking_timetracking_message_proto_rawDescGZIP(), []int{4}
}

func (x *CreateTimeTrackingReq) GetCharacterId() string {
	if x != nil {
		return x.CharacterId
	}
	return ""
}

func (x *CreateTimeTrackingReq) GetCategoryId() string {
	if x != nil && x.CategoryId != nil {
		return *x.CategoryId
	}
	return ""
}

func (x *CreateTimeTrackingReq) GetReferenceId() string {
	if x != nil && x.ReferenceId != nil {
		return *x.ReferenceId
	}
	return ""
}

func (x *CreateTimeTrackingReq) GetReferenceType() ReferenceType {
	if x != nil && x.ReferenceType != nil {
		return *x.ReferenceType
	}
	return ReferenceType_Habit
}

func (x *CreateTimeTrackingReq) GetStartTime() int64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

var File_timetracking_timetracking_message_proto protoreflect.FileDescriptor

var file_timetracking_timetracking_message_proto_rawDesc = []byte{
	0x0a, 0x27, 0x74, 0x69, 0x6d, 0x65, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x5f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x74, 0x69, 0x6d, 0x65, 0x74,
	0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x22, 0xd8, 0x02, 0x0a, 0x0c, 0x54, 0x69, 0x6d, 0x65,
	0x54, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x72,
	0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0b, 0x63,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x49, 0x64, 0x88, 0x01,
	0x01, 0x12, 0x26, 0x0a, 0x0c, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x0b, 0x72, 0x65, 0x66, 0x65, 0x72,
	0x65, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x47, 0x0a, 0x0e, 0x72, 0x65, 0x66,
	0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x1b, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67,
	0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x48, 0x02,
	0x52, 0x0d, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x88,
	0x01, 0x01, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x1e, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x03, 0x48, 0x03, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x88, 0x01,
	0x01, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x5f, 0x69,
	0x64, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x5f,
	0x69, 0x64, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x22, 0x83, 0x01, 0x0a, 0x14, 0x54, 0x69, 0x6d, 0x65, 0x54, 0x72, 0x61, 0x63, 0x6b,
	0x69, 0x6e, 0x67, 0x57, 0x69, 0x74, 0x68, 0x46, 0x69, 0x73, 0x68, 0x12, 0x3f, 0x0a, 0x0d, 0x74,
	0x69, 0x6d, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e,
	0x67, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x0c,
	0x74, 0x69, 0x6d, 0x65, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x12, 0x16, 0x0a, 0x06,
	0x6e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6e, 0x6f,
	0x72, 0x6d, 0x61, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x67, 0x6f, 0x6c, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x67, 0x6f, 0x6c, 0x64, 0x22, 0x57, 0x0a, 0x14, 0x54, 0x6f, 0x74, 0x61,
	0x6c, 0x54, 0x69, 0x6d, 0x65, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71,
	0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x22, 0x36, 0x0a, 0x15, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x54, 0x72,
	0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x22, 0xa4, 0x02, 0x0a, 0x15, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x72, 0x61,
	0x63, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0b, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0a, 0x63,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x26, 0x0a, 0x0c,
	0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x01, 0x52, 0x0b, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x49,
	0x64, 0x88, 0x01, 0x01, 0x12, 0x47, 0x0a, 0x0e, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x74,
	0x69, 0x6d, 0x65, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x52, 0x65, 0x66, 0x65,
	0x72, 0x65, 0x6e, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x48, 0x02, 0x52, 0x0d, 0x72, 0x65, 0x66,
	0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a,
	0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x0e, 0x0a, 0x0c,
	0x5f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x42, 0x0f, 0x0a, 0x0d,
	0x5f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x42, 0x11, 0x0a,
	0x0f, 0x5f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x2a, 0x24, 0x0a, 0x0d, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x09, 0x0a, 0x05, 0x48, 0x61, 0x62, 0x69, 0x74, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04,
	0x54, 0x61, 0x73, 0x6b, 0x10, 0x01, 0x42, 0x1e, 0x5a, 0x1c, 0x74, 0x65, 0x6e, 0x6b, 0x68, 0x6f,
	0x75, 0x72, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x74, 0x72,
	0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_timetracking_timetracking_message_proto_rawDescOnce sync.Once
	file_timetracking_timetracking_message_proto_rawDescData = file_timetracking_timetracking_message_proto_rawDesc
)

func file_timetracking_timetracking_message_proto_rawDescGZIP() []byte {
	file_timetracking_timetracking_message_proto_rawDescOnce.Do(func() {
		file_timetracking_timetracking_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_timetracking_timetracking_message_proto_rawDescData)
	})
	return file_timetracking_timetracking_message_proto_rawDescData
}

var file_timetracking_timetracking_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_timetracking_timetracking_message_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_timetracking_timetracking_message_proto_goTypes = []any{
	(ReferenceType)(0),            // 0: timetracking.ReferenceType
	(*TimeTracking)(nil),          // 1: timetracking.TimeTracking
	(*TimeTrackingWithFish)(nil),  // 2: timetracking.TimeTrackingWithFish
	(*TotalTimeTrackingReq)(nil),  // 3: timetracking.TotalTimeTrackingReq
	(*TotalTimeTrackingResp)(nil), // 4: timetracking.TotalTimeTrackingResp
	(*CreateTimeTrackingReq)(nil), // 5: timetracking.CreateTimeTrackingReq
}
var file_timetracking_timetracking_message_proto_depIdxs = []int32{
	0, // 0: timetracking.TimeTracking.reference_type:type_name -> timetracking.ReferenceType
	1, // 1: timetracking.TimeTrackingWithFish.time_tracking:type_name -> timetracking.TimeTracking
	0, // 2: timetracking.CreateTimeTrackingReq.reference_type:type_name -> timetracking.ReferenceType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_timetracking_timetracking_message_proto_init() }
func file_timetracking_timetracking_message_proto_init() {
	if File_timetracking_timetracking_message_proto != nil {
		return
	}
	file_timetracking_timetracking_message_proto_msgTypes[0].OneofWrappers = []any{}
	file_timetracking_timetracking_message_proto_msgTypes[4].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_timetracking_timetracking_message_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_timetracking_timetracking_message_proto_goTypes,
		DependencyIndexes: file_timetracking_timetracking_message_proto_depIdxs,
		EnumInfos:         file_timetracking_timetracking_message_proto_enumTypes,
		MessageInfos:      file_timetracking_timetracking_message_proto_msgTypes,
	}.Build()
	File_timetracking_timetracking_message_proto = out.File
	file_timetracking_timetracking_message_proto_rawDesc = nil
	file_timetracking_timetracking_message_proto_goTypes = nil
	file_timetracking_timetracking_message_proto_depIdxs = nil
}
