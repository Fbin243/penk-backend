// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.21.12
// source: currency/currency_message.proto

package currency

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

type CatchFishResp_FishType int32

const (
	CatchFishResp_None   CatchFishResp_FishType = 0
	CatchFishResp_Normal CatchFishResp_FishType = 1
	CatchFishResp_Gold   CatchFishResp_FishType = 2
)

// Enum value maps for CatchFishResp_FishType.
var (
	CatchFishResp_FishType_name = map[int32]string{
		0: "None",
		1: "Normal",
		2: "Gold",
	}
	CatchFishResp_FishType_value = map[string]int32{
		"None":   0,
		"Normal": 1,
		"Gold":   2,
	}
)

func (x CatchFishResp_FishType) Enum() *CatchFishResp_FishType {
	p := new(CatchFishResp_FishType)
	*p = x
	return p
}

func (x CatchFishResp_FishType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CatchFishResp_FishType) Descriptor() protoreflect.EnumDescriptor {
	return file_currency_currency_message_proto_enumTypes[0].Descriptor()
}

func (CatchFishResp_FishType) Type() protoreflect.EnumType {
	return &file_currency_currency_message_proto_enumTypes[0]
}

func (x CatchFishResp_FishType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CatchFishResp_FishType.Descriptor instead.
func (CatchFishResp_FishType) EnumDescriptor() ([]byte, []int) {
	return file_currency_currency_message_proto_rawDescGZIP(), []int{5, 0}
}

type CreateFishReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProfileId     string                 `protobuf:"bytes,1,opt,name=profile_id,json=profileId,proto3" json:"profile_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateFishReq) Reset() {
	*x = CreateFishReq{}
	mi := &file_currency_currency_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateFishReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateFishReq) ProtoMessage() {}

func (x *CreateFishReq) ProtoReflect() protoreflect.Message {
	mi := &file_currency_currency_message_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateFishReq.ProtoReflect.Descriptor instead.
func (*CreateFishReq) Descriptor() ([]byte, []int) {
	return file_currency_currency_message_proto_rawDescGZIP(), []int{0}
}

func (x *CreateFishReq) GetProfileId() string {
	if x != nil {
		return x.ProfileId
	}
	return ""
}

type CreateFishResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateFishResp) Reset() {
	*x = CreateFishResp{}
	mi := &file_currency_currency_message_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateFishResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateFishResp) ProtoMessage() {}

func (x *CreateFishResp) ProtoReflect() protoreflect.Message {
	mi := &file_currency_currency_message_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateFishResp.ProtoReflect.Descriptor instead.
func (*CreateFishResp) Descriptor() ([]byte, []int) {
	return file_currency_currency_message_proto_rawDescGZIP(), []int{1}
}

func (x *CreateFishResp) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type UpdateFishReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProfileId     string                 `protobuf:"bytes,1,opt,name=profile_id,json=profileId,proto3" json:"profile_id,omitempty"`
	Gold          int32                  `protobuf:"varint,2,opt,name=gold,proto3" json:"gold,omitempty"`
	Normal        int32                  `protobuf:"varint,3,opt,name=normal,proto3" json:"normal,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateFishReq) Reset() {
	*x = UpdateFishReq{}
	mi := &file_currency_currency_message_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateFishReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateFishReq) ProtoMessage() {}

func (x *UpdateFishReq) ProtoReflect() protoreflect.Message {
	mi := &file_currency_currency_message_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateFishReq.ProtoReflect.Descriptor instead.
func (*UpdateFishReq) Descriptor() ([]byte, []int) {
	return file_currency_currency_message_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateFishReq) GetProfileId() string {
	if x != nil {
		return x.ProfileId
	}
	return ""
}

func (x *UpdateFishReq) GetGold() int32 {
	if x != nil {
		return x.Gold
	}
	return 0
}

func (x *UpdateFishReq) GetNormal() int32 {
	if x != nil {
		return x.Normal
	}
	return 0
}

type UpdateFishResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FishId        string                 `protobuf:"bytes,1,opt,name=fish_id,json=fishId,proto3" json:"fish_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateFishResp) Reset() {
	*x = UpdateFishResp{}
	mi := &file_currency_currency_message_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateFishResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateFishResp) ProtoMessage() {}

func (x *UpdateFishResp) ProtoReflect() protoreflect.Message {
	mi := &file_currency_currency_message_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateFishResp.ProtoReflect.Descriptor instead.
func (*UpdateFishResp) Descriptor() ([]byte, []int) {
	return file_currency_currency_message_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateFishResp) GetFishId() string {
	if x != nil {
		return x.FishId
	}
	return ""
}

type CatchFishReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CatchFishReq) Reset() {
	*x = CatchFishReq{}
	mi := &file_currency_currency_message_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CatchFishReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CatchFishReq) ProtoMessage() {}

func (x *CatchFishReq) ProtoReflect() protoreflect.Message {
	mi := &file_currency_currency_message_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CatchFishReq.ProtoReflect.Descriptor instead.
func (*CatchFishReq) Descriptor() ([]byte, []int) {
	return file_currency_currency_message_proto_rawDescGZIP(), []int{4}
}

type CatchFishResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FishType      CatchFishResp_FishType `protobuf:"varint,1,opt,name=fish_type,json=fishType,proto3,enum=currency.CatchFishResp_FishType" json:"fish_type,omitempty"`
	Number        int32                  `protobuf:"varint,2,opt,name=number,proto3" json:"number,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CatchFishResp) Reset() {
	*x = CatchFishResp{}
	mi := &file_currency_currency_message_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CatchFishResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CatchFishResp) ProtoMessage() {}

func (x *CatchFishResp) ProtoReflect() protoreflect.Message {
	mi := &file_currency_currency_message_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CatchFishResp.ProtoReflect.Descriptor instead.
func (*CatchFishResp) Descriptor() ([]byte, []int) {
	return file_currency_currency_message_proto_rawDescGZIP(), []int{5}
}

func (x *CatchFishResp) GetFishType() CatchFishResp_FishType {
	if x != nil {
		return x.FishType
	}
	return CatchFishResp_None
}

func (x *CatchFishResp) GetNumber() int32 {
	if x != nil {
		return x.Number
	}
	return 0
}

type DeleteFishReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProfileId     string                 `protobuf:"bytes,1,opt,name=profile_id,json=profileId,proto3" json:"profile_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteFishReq) Reset() {
	*x = DeleteFishReq{}
	mi := &file_currency_currency_message_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteFishReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFishReq) ProtoMessage() {}

func (x *DeleteFishReq) ProtoReflect() protoreflect.Message {
	mi := &file_currency_currency_message_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFishReq.ProtoReflect.Descriptor instead.
func (*DeleteFishReq) Descriptor() ([]byte, []int) {
	return file_currency_currency_message_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteFishReq) GetProfileId() string {
	if x != nil {
		return x.ProfileId
	}
	return ""
}

type DeleteFishResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteFishResp) Reset() {
	*x = DeleteFishResp{}
	mi := &file_currency_currency_message_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteFishResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFishResp) ProtoMessage() {}

func (x *DeleteFishResp) ProtoReflect() protoreflect.Message {
	mi := &file_currency_currency_message_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFishResp.ProtoReflect.Descriptor instead.
func (*DeleteFishResp) Descriptor() ([]byte, []int) {
	return file_currency_currency_message_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteFishResp) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_currency_currency_message_proto protoreflect.FileDescriptor

const file_currency_currency_message_proto_rawDesc = "" +
	"\n" +
	"\x1fcurrency/currency_message.proto\x12\bcurrency\".\n" +
	"\rCreateFishReq\x12\x1d\n" +
	"\n" +
	"profile_id\x18\x01 \x01(\tR\tprofileId\"*\n" +
	"\x0eCreateFishResp\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\"Z\n" +
	"\rUpdateFishReq\x12\x1d\n" +
	"\n" +
	"profile_id\x18\x01 \x01(\tR\tprofileId\x12\x12\n" +
	"\x04gold\x18\x02 \x01(\x05R\x04gold\x12\x16\n" +
	"\x06normal\x18\x03 \x01(\x05R\x06normal\")\n" +
	"\x0eUpdateFishResp\x12\x17\n" +
	"\afish_id\x18\x01 \x01(\tR\x06fishId\"\x0e\n" +
	"\fCatchFishReq\"\x92\x01\n" +
	"\rCatchFishResp\x12=\n" +
	"\tfish_type\x18\x01 \x01(\x0e2 .currency.CatchFishResp.FishTypeR\bfishType\x12\x16\n" +
	"\x06number\x18\x02 \x01(\x05R\x06number\"*\n" +
	"\bFishType\x12\b\n" +
	"\x04None\x10\x00\x12\n" +
	"\n" +
	"\x06Normal\x10\x01\x12\b\n" +
	"\x04Gold\x10\x02\".\n" +
	"\rDeleteFishReq\x12\x1d\n" +
	"\n" +
	"profile_id\x18\x01 \x01(\tR\tprofileId\"*\n" +
	"\x0eDeleteFishResp\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccessB\x1dZ\x1btenkhours/proto/pb/currencyb\x06proto3"

var (
	file_currency_currency_message_proto_rawDescOnce sync.Once
	file_currency_currency_message_proto_rawDescData []byte
)

func file_currency_currency_message_proto_rawDescGZIP() []byte {
	file_currency_currency_message_proto_rawDescOnce.Do(func() {
		file_currency_currency_message_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_currency_currency_message_proto_rawDesc), len(file_currency_currency_message_proto_rawDesc)))
	})
	return file_currency_currency_message_proto_rawDescData
}

var file_currency_currency_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_currency_currency_message_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_currency_currency_message_proto_goTypes = []any{
	(CatchFishResp_FishType)(0), // 0: currency.CatchFishResp.FishType
	(*CreateFishReq)(nil),       // 1: currency.CreateFishReq
	(*CreateFishResp)(nil),      // 2: currency.CreateFishResp
	(*UpdateFishReq)(nil),       // 3: currency.UpdateFishReq
	(*UpdateFishResp)(nil),      // 4: currency.UpdateFishResp
	(*CatchFishReq)(nil),        // 5: currency.CatchFishReq
	(*CatchFishResp)(nil),       // 6: currency.CatchFishResp
	(*DeleteFishReq)(nil),       // 7: currency.DeleteFishReq
	(*DeleteFishResp)(nil),      // 8: currency.DeleteFishResp
}
var file_currency_currency_message_proto_depIdxs = []int32{
	0, // 0: currency.CatchFishResp.fish_type:type_name -> currency.CatchFishResp.FishType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_currency_currency_message_proto_init() }
func file_currency_currency_message_proto_init() {
	if File_currency_currency_message_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_currency_currency_message_proto_rawDesc), len(file_currency_currency_message_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_currency_currency_message_proto_goTypes,
		DependencyIndexes: file_currency_currency_message_proto_depIdxs,
		EnumInfos:         file_currency_currency_message_proto_enumTypes,
		MessageInfos:      file_currency_currency_message_proto_msgTypes,
	}.Build()
	File_currency_currency_message_proto = out.File
	file_currency_currency_message_proto_goTypes = nil
	file_currency_currency_message_proto_depIdxs = nil
}
