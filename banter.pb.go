// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v4.25.3
// source: banter.proto

package banter

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

type BusTopic int32

const (
	BusTopic_BANTER_EVENT   BusTopic = 0
	BusTopic_BANTER_REQUEST BusTopic = 1
	BusTopic_BANTER_COMMAND BusTopic = 2
)

// Enum value maps for BusTopic.
var (
	BusTopic_name = map[int32]string{
		0: "BANTER_EVENT",
		1: "BANTER_REQUEST",
		2: "BANTER_COMMAND",
	}
	BusTopic_value = map[string]int32{
		"BANTER_EVENT":   0,
		"BANTER_REQUEST": 1,
		"BANTER_COMMAND": 2,
	}
)

func (x BusTopic) Enum() *BusTopic {
	p := new(BusTopic)
	*p = x
	return p
}

func (x BusTopic) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BusTopic) Descriptor() protoreflect.EnumDescriptor {
	return file_banter_proto_enumTypes[0].Descriptor()
}

func (BusTopic) Type() protoreflect.EnumType {
	return &file_banter_proto_enumTypes[0]
}

func (x BusTopic) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BusTopic.Descriptor instead.
func (BusTopic) EnumDescriptor() ([]byte, []int) {
	return file_banter_proto_rawDescGZIP(), []int{0}
}

type MessageTypeRequest int32

const (
	MessageTypeRequest_CONFIG_GET_REQ  MessageTypeRequest = 0
	MessageTypeRequest_CONFIG_GET_RESP MessageTypeRequest = 1
)

// Enum value maps for MessageTypeRequest.
var (
	MessageTypeRequest_name = map[int32]string{
		0: "CONFIG_GET_REQ",
		1: "CONFIG_GET_RESP",
	}
	MessageTypeRequest_value = map[string]int32{
		"CONFIG_GET_REQ":  0,
		"CONFIG_GET_RESP": 1,
	}
)

func (x MessageTypeRequest) Enum() *MessageTypeRequest {
	p := new(MessageTypeRequest)
	*p = x
	return p
}

func (x MessageTypeRequest) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageTypeRequest) Descriptor() protoreflect.EnumDescriptor {
	return file_banter_proto_enumTypes[1].Descriptor()
}

func (MessageTypeRequest) Type() protoreflect.EnumType {
	return &file_banter_proto_enumTypes[1]
}

func (x MessageTypeRequest) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageTypeRequest.Descriptor instead.
func (MessageTypeRequest) EnumDescriptor() ([]byte, []int) {
	return file_banter_proto_rawDescGZIP(), []int{1}
}

type MessageTypeCommand int32

const (
	MessageTypeCommand_CONFIG_SET_REQ  MessageTypeCommand = 0
	MessageTypeCommand_CONFIG_SET_RESP MessageTypeCommand = 1
)

// Enum value maps for MessageTypeCommand.
var (
	MessageTypeCommand_name = map[int32]string{
		0: "CONFIG_SET_REQ",
		1: "CONFIG_SET_RESP",
	}
	MessageTypeCommand_value = map[string]int32{
		"CONFIG_SET_REQ":  0,
		"CONFIG_SET_RESP": 1,
	}
)

func (x MessageTypeCommand) Enum() *MessageTypeCommand {
	p := new(MessageTypeCommand)
	*p = x
	return p
}

func (x MessageTypeCommand) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageTypeCommand) Descriptor() protoreflect.EnumDescriptor {
	return file_banter_proto_enumTypes[2].Descriptor()
}

func (MessageTypeCommand) Type() protoreflect.EnumType {
	return &file_banter_proto_enumTypes[2]
}

func (x MessageTypeCommand) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageTypeCommand.Descriptor instead.
func (MessageTypeCommand) EnumDescriptor() ([]byte, []int) {
	return file_banter_proto_rawDescGZIP(), []int{2}
}

type Banter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Command  string `protobuf:"bytes,1,opt,name=command,proto3" json:"command,omitempty"`
	Text     string `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	Disabled bool   `protobuf:"varint,3,opt,name=disabled,proto3" json:"disabled,omitempty"`
	Random   bool   `protobuf:"varint,4,opt,name=random,proto3" json:"random,omitempty"`
}

func (x *Banter) Reset() {
	*x = Banter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_banter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Banter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Banter) ProtoMessage() {}

func (x *Banter) ProtoReflect() protoreflect.Message {
	mi := &file_banter_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Banter.ProtoReflect.Descriptor instead.
func (*Banter) Descriptor() ([]byte, []int) {
	return file_banter_proto_rawDescGZIP(), []int{0}
}

func (x *Banter) GetCommand() string {
	if x != nil {
		return x.Command
	}
	return ""
}

func (x *Banter) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *Banter) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *Banter) GetRandom() bool {
	if x != nil {
		return x.Random
	}
	return false
}

type EventSettings struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Enabled bool   `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Text    string `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *EventSettings) Reset() {
	*x = EventSettings{}
	if protoimpl.UnsafeEnabled {
		mi := &file_banter_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventSettings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventSettings) ProtoMessage() {}

func (x *EventSettings) ProtoReflect() protoreflect.Message {
	mi := &file_banter_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventSettings.ProtoReflect.Descriptor instead.
func (*EventSettings) Descriptor() ([]byte, []int) {
	return file_banter_proto_rawDescGZIP(), []int{1}
}

func (x *EventSettings) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *EventSettings) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type GuestList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Members []*GuestList_Member `protobuf:"bytes,1,rep,name=members,proto3" json:"members,omitempty"`
}

func (x *GuestList) Reset() {
	*x = GuestList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_banter_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuestList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuestList) ProtoMessage() {}

func (x *GuestList) ProtoReflect() protoreflect.Message {
	mi := &file_banter_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuestList.ProtoReflect.Descriptor instead.
func (*GuestList) Descriptor() ([]byte, []int) {
	return file_banter_proto_rawDescGZIP(), []int{2}
}

func (x *GuestList) GetMembers() []*GuestList_Member {
	if x != nil {
		return x.Members
	}
	return nil
}

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IntervalSeconds uint32                `protobuf:"varint,1,opt,name=interval_seconds,json=intervalSeconds,proto3" json:"interval_seconds,omitempty"`
	CooldownSeconds uint32                `protobuf:"varint,2,opt,name=cooldown_seconds,json=cooldownSeconds,proto3" json:"cooldown_seconds,omitempty"`
	Banters         []*Banter             `protobuf:"bytes,3,rep,name=banters,proto3" json:"banters,omitempty"`
	ChannelRaid     *EventSettings        `protobuf:"bytes,4,opt,name=channel_raid,json=channelRaid,proto3" json:"channel_raid,omitempty"`
	ChannelFollow   *EventSettings        `protobuf:"bytes,5,opt,name=channel_follow,json=channelFollow,proto3" json:"channel_follow,omitempty"`
	ChannelCheer    *EventSettings        `protobuf:"bytes,6,opt,name=channel_cheer,json=channelCheer,proto3" json:"channel_cheer,omitempty"`
	GuestLists      map[string]*GuestList `protobuf:"bytes,8,rep,name=guest_lists,json=guestLists,proto3" json:"guest_lists,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_banter_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_banter_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_banter_proto_rawDescGZIP(), []int{3}
}

func (x *Config) GetIntervalSeconds() uint32 {
	if x != nil {
		return x.IntervalSeconds
	}
	return 0
}

func (x *Config) GetCooldownSeconds() uint32 {
	if x != nil {
		return x.CooldownSeconds
	}
	return 0
}

func (x *Config) GetBanters() []*Banter {
	if x != nil {
		return x.Banters
	}
	return nil
}

func (x *Config) GetChannelRaid() *EventSettings {
	if x != nil {
		return x.ChannelRaid
	}
	return nil
}

func (x *Config) GetChannelFollow() *EventSettings {
	if x != nil {
		return x.ChannelFollow
	}
	return nil
}

func (x *Config) GetChannelCheer() *EventSettings {
	if x != nil {
		return x.ChannelCheer
	}
	return nil
}

func (x *Config) GetGuestLists() map[string]*GuestList {
	if x != nil {
		return x.GuestLists
	}
	return nil
}

type ConfigGetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ConfigGetRequest) Reset() {
	*x = ConfigGetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_banter_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigGetRequest) ProtoMessage() {}

func (x *ConfigGetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_banter_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigGetRequest.ProtoReflect.Descriptor instead.
func (*ConfigGetRequest) Descriptor() ([]byte, []int) {
	return file_banter_proto_rawDescGZIP(), []int{4}
}

type ConfigGetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Config *Config `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *ConfigGetResponse) Reset() {
	*x = ConfigGetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_banter_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigGetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigGetResponse) ProtoMessage() {}

func (x *ConfigGetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_banter_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigGetResponse.ProtoReflect.Descriptor instead.
func (*ConfigGetResponse) Descriptor() ([]byte, []int) {
	return file_banter_proto_rawDescGZIP(), []int{5}
}

func (x *ConfigGetResponse) GetConfig() *Config {
	if x != nil {
		return x.Config
	}
	return nil
}

type ConfigSetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Config *Config `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *ConfigSetRequest) Reset() {
	*x = ConfigSetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_banter_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigSetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigSetRequest) ProtoMessage() {}

func (x *ConfigSetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_banter_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigSetRequest.ProtoReflect.Descriptor instead.
func (*ConfigSetRequest) Descriptor() ([]byte, []int) {
	return file_banter_proto_rawDescGZIP(), []int{6}
}

func (x *ConfigSetRequest) GetConfig() *Config {
	if x != nil {
		return x.Config
	}
	return nil
}

type ConfigSetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Config *Config `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *ConfigSetResponse) Reset() {
	*x = ConfigSetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_banter_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigSetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigSetResponse) ProtoMessage() {}

func (x *ConfigSetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_banter_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigSetResponse.ProtoReflect.Descriptor instead.
func (*ConfigSetResponse) Descriptor() ([]byte, []int) {
	return file_banter_proto_rawDescGZIP(), []int{7}
}

func (x *ConfigSetResponse) GetConfig() *Config {
	if x != nil {
		return x.Config
	}
	return nil
}

type GuestList_Member struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login string `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Id    string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GuestList_Member) Reset() {
	*x = GuestList_Member{}
	if protoimpl.UnsafeEnabled {
		mi := &file_banter_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuestList_Member) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuestList_Member) ProtoMessage() {}

func (x *GuestList_Member) ProtoReflect() protoreflect.Message {
	mi := &file_banter_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuestList_Member.ProtoReflect.Descriptor instead.
func (*GuestList_Member) Descriptor() ([]byte, []int) {
	return file_banter_proto_rawDescGZIP(), []int{2, 0}
}

func (x *GuestList_Member) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *GuestList_Member) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_banter_proto protoreflect.FileDescriptor

var file_banter_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x62, 0x61, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x62, 0x61, 0x6e, 0x74, 0x65, 0x72, 0x22, 0x6a, 0x0a, 0x06, 0x42, 0x61, 0x6e, 0x74, 0x65, 0x72,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65,
	0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x61,
	0x6e, 0x64, 0x6f, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x72, 0x61, 0x6e, 0x64,
	0x6f, 0x6d, 0x22, 0x3d, 0x0a, 0x0d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x74, 0x74, 0x69,
	0x6e, 0x67, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78,
	0x74, 0x22, 0x6f, 0x0a, 0x09, 0x47, 0x75, 0x65, 0x73, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x32,
	0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x18, 0x2e, 0x62, 0x61, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x47, 0x75, 0x65, 0x73, 0x74, 0x4c, 0x69,
	0x73, 0x74, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x73, 0x1a, 0x2e, 0x0a, 0x06, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05,
	0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x6f, 0x67,
	0x69, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x22, 0xd5, 0x03, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x29, 0x0a,
	0x10, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61,
	0x6c, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x63, 0x6f, 0x6f, 0x6c,
	0x64, 0x6f, 0x77, 0x6e, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0f, 0x63, 0x6f, 0x6f, 0x6c, 0x64, 0x6f, 0x77, 0x6e, 0x53, 0x65, 0x63, 0x6f,
	0x6e, 0x64, 0x73, 0x12, 0x28, 0x0a, 0x07, 0x62, 0x61, 0x6e, 0x74, 0x65, 0x72, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x62, 0x61, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x42, 0x61,
	0x6e, 0x74, 0x65, 0x72, 0x52, 0x07, 0x62, 0x61, 0x6e, 0x74, 0x65, 0x72, 0x73, 0x12, 0x38, 0x0a,
	0x0c, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x72, 0x61, 0x69, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x62, 0x61, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x52, 0x61, 0x69, 0x64, 0x12, 0x3c, 0x0a, 0x0e, 0x63, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x5f, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x62, 0x61, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x65,
	0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x0d, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x46,
	0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x12, 0x3a, 0x0a, 0x0d, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x5f, 0x63, 0x68, 0x65, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x62,
	0x61, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x74, 0x74, 0x69,
	0x6e, 0x67, 0x73, 0x52, 0x0c, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x43, 0x68, 0x65, 0x65,
	0x72, 0x12, 0x3f, 0x0a, 0x0b, 0x67, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x73,
	0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x62, 0x61, 0x6e, 0x74, 0x65, 0x72, 0x2e,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x47, 0x75, 0x65, 0x73, 0x74, 0x4c, 0x69, 0x73, 0x74,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x67, 0x75, 0x65, 0x73, 0x74, 0x4c, 0x69, 0x73,
	0x74, 0x73, 0x1a, 0x50, 0x0a, 0x0f, 0x47, 0x75, 0x65, 0x73, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x27, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x62, 0x61, 0x6e, 0x74, 0x65, 0x72, 0x2e,
	0x47, 0x75, 0x65, 0x73, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x4a, 0x04, 0x08, 0x07, 0x10, 0x08, 0x22, 0x12, 0x0a, 0x10, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3b,
	0x0a, 0x11, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x62, 0x61, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x3a, 0x0a, 0x10, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x26, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x62, 0x61, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52,
	0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x3b, 0x0a, 0x11, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x06,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x62,
	0x61, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x06, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2a, 0x44, 0x0a, 0x08, 0x42, 0x75, 0x73, 0x54, 0x6f, 0x70, 0x69, 0x63,
	0x12, 0x10, 0x0a, 0x0c, 0x42, 0x41, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54,
	0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x42, 0x41, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x52, 0x45, 0x51,
	0x55, 0x45, 0x53, 0x54, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x42, 0x41, 0x4e, 0x54, 0x45, 0x52,
	0x5f, 0x43, 0x4f, 0x4d, 0x4d, 0x41, 0x4e, 0x44, 0x10, 0x02, 0x2a, 0x3d, 0x0a, 0x12, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x0e, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x5f, 0x47, 0x45, 0x54, 0x5f, 0x52,
	0x45, 0x51, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x5f, 0x47,
	0x45, 0x54, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x10, 0x01, 0x2a, 0x3d, 0x0a, 0x12, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12,
	0x12, 0x0a, 0x0e, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x5f, 0x53, 0x45, 0x54, 0x5f, 0x52, 0x45,
	0x51, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x5f, 0x53, 0x45,
	0x54, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x10, 0x01, 0x42, 0x21, 0x5a, 0x1f, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x75, 0x74, 0x6f, 0x6e, 0x6f, 0x6d, 0x6f, 0x75,
	0x73, 0x6b, 0x6f, 0x69, 0x2f, 0x62, 0x61, 0x6e, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_banter_proto_rawDescOnce sync.Once
	file_banter_proto_rawDescData = file_banter_proto_rawDesc
)

func file_banter_proto_rawDescGZIP() []byte {
	file_banter_proto_rawDescOnce.Do(func() {
		file_banter_proto_rawDescData = protoimpl.X.CompressGZIP(file_banter_proto_rawDescData)
	})
	return file_banter_proto_rawDescData
}

var file_banter_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_banter_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_banter_proto_goTypes = []any{
	(BusTopic)(0),             // 0: banter.BusTopic
	(MessageTypeRequest)(0),   // 1: banter.MessageTypeRequest
	(MessageTypeCommand)(0),   // 2: banter.MessageTypeCommand
	(*Banter)(nil),            // 3: banter.Banter
	(*EventSettings)(nil),     // 4: banter.EventSettings
	(*GuestList)(nil),         // 5: banter.GuestList
	(*Config)(nil),            // 6: banter.Config
	(*ConfigGetRequest)(nil),  // 7: banter.ConfigGetRequest
	(*ConfigGetResponse)(nil), // 8: banter.ConfigGetResponse
	(*ConfigSetRequest)(nil),  // 9: banter.ConfigSetRequest
	(*ConfigSetResponse)(nil), // 10: banter.ConfigSetResponse
	(*GuestList_Member)(nil),  // 11: banter.GuestList.Member
	nil,                       // 12: banter.Config.GuestListsEntry
}
var file_banter_proto_depIdxs = []int32{
	11, // 0: banter.GuestList.members:type_name -> banter.GuestList.Member
	3,  // 1: banter.Config.banters:type_name -> banter.Banter
	4,  // 2: banter.Config.channel_raid:type_name -> banter.EventSettings
	4,  // 3: banter.Config.channel_follow:type_name -> banter.EventSettings
	4,  // 4: banter.Config.channel_cheer:type_name -> banter.EventSettings
	12, // 5: banter.Config.guest_lists:type_name -> banter.Config.GuestListsEntry
	6,  // 6: banter.ConfigGetResponse.config:type_name -> banter.Config
	6,  // 7: banter.ConfigSetRequest.config:type_name -> banter.Config
	6,  // 8: banter.ConfigSetResponse.config:type_name -> banter.Config
	5,  // 9: banter.Config.GuestListsEntry.value:type_name -> banter.GuestList
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_banter_proto_init() }
func file_banter_proto_init() {
	if File_banter_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_banter_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Banter); i {
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
		file_banter_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*EventSettings); i {
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
		file_banter_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*GuestList); i {
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
		file_banter_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Config); i {
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
		file_banter_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*ConfigGetRequest); i {
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
		file_banter_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*ConfigGetResponse); i {
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
		file_banter_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*ConfigSetRequest); i {
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
		file_banter_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*ConfigSetResponse); i {
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
		file_banter_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*GuestList_Member); i {
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
			RawDescriptor: file_banter_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_banter_proto_goTypes,
		DependencyIndexes: file_banter_proto_depIdxs,
		EnumInfos:         file_banter_proto_enumTypes,
		MessageInfos:      file_banter_proto_msgTypes,
	}.Build()
	File_banter_proto = out.File
	file_banter_proto_rawDesc = nil
	file_banter_proto_goTypes = nil
	file_banter_proto_depIdxs = nil
}
