// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        (unknown)
// source: hooks/slack/v1/slack.proto

package slackv1

import (
	_ "go.breu.io/quantm/internal/proto/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OauthRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          string                 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Kind          string                 `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	LinkTo        string                 `protobuf:"bytes,3,opt,name=link_to,json=linkTo,proto3" json:"link_to,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OauthRequest) Reset() {
	*x = OauthRequest{}
	mi := &file_hooks_slack_v1_slack_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OauthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OauthRequest) ProtoMessage() {}

func (x *OauthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_hooks_slack_v1_slack_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OauthRequest.ProtoReflect.Descriptor instead.
func (*OauthRequest) Descriptor() ([]byte, []int) {
	return file_hooks_slack_v1_slack_proto_rawDescGZIP(), []int{0}
}

func (x *OauthRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *OauthRequest) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *OauthRequest) GetLinkTo() string {
	if x != nil {
		return x.LinkTo
	}
	return ""
}

// @ysf - we might need a db entity for this. just like we have for github
type Install struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	WorkspaceId   string                 `protobuf:"bytes,1,opt,name=workspace_id,json=workspaceId,proto3" json:"workspace_id,omitempty"`
	WorkspaceName string                 `protobuf:"bytes,2,opt,name=workspace_name,json=workspaceName,proto3" json:"workspace_name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Install) Reset() {
	*x = Install{}
	mi := &file_hooks_slack_v1_slack_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Install) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Install) ProtoMessage() {}

func (x *Install) ProtoReflect() protoreflect.Message {
	mi := &file_hooks_slack_v1_slack_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Install.ProtoReflect.Descriptor instead.
func (*Install) Descriptor() ([]byte, []int) {
	return file_hooks_slack_v1_slack_proto_rawDescGZIP(), []int{1}
}

func (x *Install) GetWorkspaceId() string {
	if x != nil {
		return x.WorkspaceId
	}
	return ""
}

func (x *Install) GetWorkspaceName() string {
	if x != nil {
		return x.WorkspaceName
	}
	return ""
}

var File_hooks_slack_v1_slack_proto protoreflect.FileDescriptor

var file_hooks_slack_v1_slack_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x2f, 0x73, 0x6c, 0x61, 0x63, 0x6b, 0x2f, 0x76, 0x31,
	0x2f, 0x73, 0x6c, 0x61, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x68, 0x6f,
	0x6f, 0x6b, 0x73, 0x2e, 0x73, 0x6c, 0x61, 0x63, 0x6b, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75,
	0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x59, 0x0a, 0x0c, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69,
	0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x21,
	0x0a, 0x07, 0x6c, 0x69, 0x6e, 0x6b, 0x5f, 0x74, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x06, 0x6c, 0x69, 0x6e, 0x6b, 0x54,
	0x6f, 0x22, 0x53, 0x0a, 0x07, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x12, 0x21, 0x0a, 0x0c,
	0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x64, 0x12,
	0x25, 0x0a, 0x0e, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x32, 0x4d, 0x0a, 0x0c, 0x53, 0x6c, 0x61, 0x63, 0x6b, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3d, 0x0a, 0x05, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x12,
	0x1c, 0x2e, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x2e, 0x73, 0x6c, 0x61, 0x63, 0x6b, 0x2e, 0x76, 0x31,
	0x2e, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0xb3, 0x01, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x6f,
	0x6f, 0x6b, 0x73, 0x2e, 0x73, 0x6c, 0x61, 0x63, 0x6b, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x53, 0x6c,
	0x61, 0x63, 0x6b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x37, 0x67, 0x6f, 0x2e, 0x62,
	0x72, 0x65, 0x75, 0x2e, 0x69, 0x6f, 0x2f, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x6d, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x68, 0x6f, 0x6f,
	0x6b, 0x73, 0x2f, 0x73, 0x6c, 0x61, 0x63, 0x6b, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x6c, 0x61, 0x63,
	0x6b, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x48, 0x53, 0x58, 0xaa, 0x02, 0x0e, 0x48, 0x6f, 0x6f, 0x6b,
	0x73, 0x2e, 0x53, 0x6c, 0x61, 0x63, 0x6b, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0e, 0x48, 0x6f, 0x6f,
	0x6b, 0x73, 0x5c, 0x53, 0x6c, 0x61, 0x63, 0x6b, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1a, 0x48, 0x6f,
	0x6f, 0x6b, 0x73, 0x5c, 0x53, 0x6c, 0x61, 0x63, 0x6b, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x10, 0x48, 0x6f, 0x6f, 0x6b, 0x73,
	0x3a, 0x3a, 0x53, 0x6c, 0x61, 0x63, 0x6b, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_hooks_slack_v1_slack_proto_rawDescOnce sync.Once
	file_hooks_slack_v1_slack_proto_rawDescData = file_hooks_slack_v1_slack_proto_rawDesc
)

func file_hooks_slack_v1_slack_proto_rawDescGZIP() []byte {
	file_hooks_slack_v1_slack_proto_rawDescOnce.Do(func() {
		file_hooks_slack_v1_slack_proto_rawDescData = protoimpl.X.CompressGZIP(file_hooks_slack_v1_slack_proto_rawDescData)
	})
	return file_hooks_slack_v1_slack_proto_rawDescData
}

var file_hooks_slack_v1_slack_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_hooks_slack_v1_slack_proto_goTypes = []any{
	(*OauthRequest)(nil),  // 0: hooks.slack.v1.OauthRequest
	(*Install)(nil),       // 1: hooks.slack.v1.Install
	(*emptypb.Empty)(nil), // 2: google.protobuf.Empty
}
var file_hooks_slack_v1_slack_proto_depIdxs = []int32{
	0, // 0: hooks.slack.v1.SlackService.Oauth:input_type -> hooks.slack.v1.OauthRequest
	2, // 1: hooks.slack.v1.SlackService.Oauth:output_type -> google.protobuf.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_hooks_slack_v1_slack_proto_init() }
func file_hooks_slack_v1_slack_proto_init() {
	if File_hooks_slack_v1_slack_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_hooks_slack_v1_slack_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_hooks_slack_v1_slack_proto_goTypes,
		DependencyIndexes: file_hooks_slack_v1_slack_proto_depIdxs,
		MessageInfos:      file_hooks_slack_v1_slack_proto_msgTypes,
	}.Build()
	File_hooks_slack_v1_slack_proto = out.File
	file_hooks_slack_v1_slack_proto_rawDesc = nil
	file_hooks_slack_v1_slack_proto_goTypes = nil
	file_hooks_slack_v1_slack_proto_depIdxs = nil
}