// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.11.4
// source: chacerapp/v1/templates.proto

package serverpb

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Template represents a pre-configured message that can be used to generate messages.
type Template struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The resource name of the template.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The friendly name that should be used when displaying the template.
	DisplayName string `protobuf:"bytes,2,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	// The configuration that should be used as the default values
	// when creating a message.
	Message *Message `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	// The recipient
	Recipient string `protobuf:"bytes,4,opt,name=recipient,proto3" json:"recipient,omitempty"`
	// Required. The person that sent the message. It must be the full
	// resource name of the recipient. The sender must be exist within the
	// same location that the message is being sent to.
	Sender string `protobuf:"bytes,5,opt,name=sender,proto3" json:"sender,omitempty"`
	// Optional. The room in the location that the recipient is being
	// requested in. It must be the resource name of the room.
	RequestedRoom string `protobuf:"bytes,6,opt,name=requested_room,json=requestedRoom,proto3" json:"requested_room,omitempty"`
	// Output Only. The resource name of the location the template exists in.
	Location string `protobuf:"bytes,7,opt,name=location,proto3" json:"location,omitempty"`
	// The short reason that the recipient is needed in a room. This field
	// can be a maximum of 100 characters.
	Reason string `protobuf:"bytes,8,opt,name=reason,proto3" json:"reason,omitempty"`
	// A longer description of why the recipient is needed in a room. This
	// field supports a maximum length of 1024 characters.
	Description string `protobuf:"bytes,9,opt,name=description,proto3" json:"description,omitempty"`
	// Output Only. Server-defined URL for the resource.
	SelfLink string `protobuf:"bytes,100,opt,name=self_link,json=selfLink,proto3" json:"self_link,omitempty"`
	// Output Only. The time the resource was created.
	CreatedTime *timestamp.Timestamp `protobuf:"bytes,101,opt,name=created_time,json=createdTime,proto3" json:"created_time,omitempty"`
	// Output Only. The time the resource was updated.
	UpdatedTime *timestamp.Timestamp `protobuf:"bytes,102,opt,name=updated_time,json=updatedTime,proto3" json:"updated_time,omitempty"`
}

func (x *Template) Reset() {
	*x = Template{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chacerapp_v1_templates_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Template) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Template) ProtoMessage() {}

func (x *Template) ProtoReflect() protoreflect.Message {
	mi := &file_chacerapp_v1_templates_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Template.ProtoReflect.Descriptor instead.
func (*Template) Descriptor() ([]byte, []int) {
	return file_chacerapp_v1_templates_proto_rawDescGZIP(), []int{0}
}

func (x *Template) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Template) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *Template) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *Template) GetRecipient() string {
	if x != nil {
		return x.Recipient
	}
	return ""
}

func (x *Template) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *Template) GetRequestedRoom() string {
	if x != nil {
		return x.RequestedRoom
	}
	return ""
}

func (x *Template) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *Template) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *Template) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Template) GetSelfLink() string {
	if x != nil {
		return x.SelfLink
	}
	return ""
}

func (x *Template) GetCreatedTime() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedTime
	}
	return nil
}

func (x *Template) GetUpdatedTime() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedTime
	}
	return nil
}

type ListTemplatesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The parent (account and location) where the templates will be listed
	// Specified in the format 'accounts/*/locations/*'.
	// Location "-" can be used to list devices in all locations.
	Parent string `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
	// The max number of results per page that should be returned. If the number
	// of available results is larger than `page_size`, a `next_page_token` is
	// returned which can be used to get the next page of results in subsequent
	// requests. Acceptable values are 1 to 500, inclusive. (Default: 500)
	PageSize int32 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// Specifies a page token to use. Set this to the nextPageToken returned by
	// previous list requests to get the next page of results.
	PageToken string `protobuf:"bytes,3,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
}

func (x *ListTemplatesRequest) Reset() {
	*x = ListTemplatesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chacerapp_v1_templates_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListTemplatesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTemplatesRequest) ProtoMessage() {}

func (x *ListTemplatesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chacerapp_v1_templates_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTemplatesRequest.ProtoReflect.Descriptor instead.
func (*ListTemplatesRequest) Descriptor() ([]byte, []int) {
	return file_chacerapp_v1_templates_proto_rawDescGZIP(), []int{1}
}

func (x *ListTemplatesRequest) GetParent() string {
	if x != nil {
		return x.Parent
	}
	return ""
}

func (x *ListTemplatesRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListTemplatesRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

// ListTemplatesResponse will list the templates for a location.
type ListTemplatesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A list of templates in the specified location.
	Templates []*Template `protobuf:"bytes,1,rep,name=templates,proto3" json:"templates,omitempty"`
	// This token allows you to get the next page of results for list requests.
	// If the number of results is larger than `page_size`, use the
	// `next_page_token` as a value for the query parameter `page_token` in the
	// next request. The value will become empty when there are no more pages.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *ListTemplatesResponse) Reset() {
	*x = ListTemplatesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chacerapp_v1_templates_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListTemplatesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTemplatesResponse) ProtoMessage() {}

func (x *ListTemplatesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chacerapp_v1_templates_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTemplatesResponse.ProtoReflect.Descriptor instead.
func (*ListTemplatesResponse) Descriptor() ([]byte, []int) {
	return file_chacerapp_v1_templates_proto_rawDescGZIP(), []int{2}
}

func (x *ListTemplatesResponse) GetTemplates() []*Template {
	if x != nil {
		return x.Templates
	}
	return nil
}

func (x *ListTemplatesResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

type CreateTemplateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The parent (account and location) where the templates will be created.
	// Specified in the format 'accounts/*/locations/*'.
	Parent string `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
	// The template that should be created.
	Template *Template `protobuf:"bytes,2,opt,name=template,proto3" json:"template,omitempty"`
}

func (x *CreateTemplateRequest) Reset() {
	*x = CreateTemplateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chacerapp_v1_templates_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTemplateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTemplateRequest) ProtoMessage() {}

func (x *CreateTemplateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chacerapp_v1_templates_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTemplateRequest.ProtoReflect.Descriptor instead.
func (*CreateTemplateRequest) Descriptor() ([]byte, []int) {
	return file_chacerapp_v1_templates_proto_rawDescGZIP(), []int{3}
}

func (x *CreateTemplateRequest) GetParent() string {
	if x != nil {
		return x.Parent
	}
	return ""
}

func (x *CreateTemplateRequest) GetTemplate() *Template {
	if x != nil {
		return x.Template
	}
	return nil
}

var File_chacerapp_v1_templates_proto protoreflect.FileDescriptor

var file_chacerapp_v1_templates_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x63, 0x68, 0x61, 0x63, 0x65, 0x72, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x2f, 0x74,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13,
	0x63, 0x68, 0x61, 0x63, 0x65, 0x72, 0x61, 0x70, 0x70, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x63, 0x68, 0x61, 0x63, 0x65, 0x72, 0x61, 0x70, 0x70, 0x2f, 0x76,
	0x31, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x88, 0x04, 0x0a, 0x08,
	0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c,
	0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x2f, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x15, 0x2e, 0x63, 0x68, 0x61, 0x63, 0x65, 0x72, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x40, 0x0a, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x22, 0xe0, 0x41, 0x02, 0xfa, 0x41, 0x1c, 0x0a, 0x1a, 0x61, 0x70, 0x69,
	0x73, 0x2e, 0x63, 0x68, 0x61, 0x63, 0x65, 0x72, 0x61, 0x70, 0x70, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65,
	0x6e, 0x74, 0x12, 0x3a, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x22, 0xe0, 0x41, 0x02, 0xfa, 0x41, 0x1c, 0x0a, 0x1a, 0x61, 0x70, 0x69, 0x73,
	0x2e, 0x63, 0x68, 0x61, 0x63, 0x65, 0x72, 0x61, 0x70, 0x70, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x43,
	0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x25,
	0x0a, 0x0e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x5f, 0x72, 0x6f, 0x6f, 0x6d,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65,
	0x64, 0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x73,
	0x65, 0x6c, 0x66, 0x5f, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x64, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x73, 0x65, 0x6c, 0x66, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x3d, 0x0a, 0x0c, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x65, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x3d, 0x0a, 0x0c, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x66, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x6a, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x65,
	0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73,
	0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53,
	0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x22, 0x7c, 0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x09, 0x74,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d,
	0x2e, 0x63, 0x68, 0x61, 0x63, 0x65, 0x72, 0x61, 0x70, 0x70, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x09, 0x74,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74,
	0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x22, 0x6a, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x61, 0x72,
	0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e,
	0x74, 0x12, 0x39, 0x0a, 0x08, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x63, 0x68, 0x61, 0x63, 0x65, 0x72, 0x61, 0x70, 0x70, 0x2e,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x52, 0x08, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x32, 0xc0, 0x02, 0x0a,
	0x09, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x12, 0x9d, 0x01, 0x0a, 0x0d, 0x4c,
	0x69, 0x73, 0x74, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x12, 0x29, 0x2e, 0x63,
	0x68, 0x61, 0x63, 0x65, 0x72, 0x61, 0x70, 0x70, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x63, 0x68, 0x61, 0x63, 0x65, 0x72,
	0x61, 0x70, 0x70, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x35, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2f, 0x12, 0x2d, 0x2f, 0x76, 0x31,
	0x2f, 0x7b, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x3d, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x73, 0x2f, 0x2a, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x2a, 0x7d,
	0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x12, 0x92, 0x01, 0x0a, 0x0e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x2a, 0x2e,
	0x63, 0x68, 0x61, 0x63, 0x65, 0x72, 0x61, 0x70, 0x70, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x63, 0x68, 0x61, 0x63,
	0x65, 0x72, 0x61, 0x70, 0x70, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e,
	0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x22, 0x35, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2f,
	0x22, 0x2d, 0x2f, 0x76, 0x31, 0x2f, 0x7b, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x3d, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x2f, 0x2a, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x2a, 0x7d, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x42,
	0x94, 0x01, 0x0a, 0x17, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x68, 0x61, 0x63, 0x65, 0x72, 0x61, 0x70,
	0x70, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x42, 0x15, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x63, 0x68, 0x61, 0x63, 0x65, 0x72, 0x61, 0x70, 0x70, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x70, 0x62, 0xaa, 0x02, 0x16, 0x43, 0x68, 0x61, 0x63, 0x65, 0x72, 0x61, 0x70, 0x70,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x16,
	0x43, 0x68, 0x61, 0x63, 0x65, 0x72, 0x61, 0x70, 0x70, 0x5c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x69, 0x6e, 0x67, 0x5c, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_chacerapp_v1_templates_proto_rawDescOnce sync.Once
	file_chacerapp_v1_templates_proto_rawDescData = file_chacerapp_v1_templates_proto_rawDesc
)

func file_chacerapp_v1_templates_proto_rawDescGZIP() []byte {
	file_chacerapp_v1_templates_proto_rawDescOnce.Do(func() {
		file_chacerapp_v1_templates_proto_rawDescData = protoimpl.X.CompressGZIP(file_chacerapp_v1_templates_proto_rawDescData)
	})
	return file_chacerapp_v1_templates_proto_rawDescData
}

var file_chacerapp_v1_templates_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_chacerapp_v1_templates_proto_goTypes = []interface{}{
	(*Template)(nil),              // 0: chacerapp.config.v1.Template
	(*ListTemplatesRequest)(nil),  // 1: chacerapp.config.v1.ListTemplatesRequest
	(*ListTemplatesResponse)(nil), // 2: chacerapp.config.v1.ListTemplatesResponse
	(*CreateTemplateRequest)(nil), // 3: chacerapp.config.v1.CreateTemplateRequest
	(*Message)(nil),               // 4: chacerapp.v1.Message
	(*timestamp.Timestamp)(nil),   // 5: google.protobuf.Timestamp
}
var file_chacerapp_v1_templates_proto_depIdxs = []int32{
	4, // 0: chacerapp.config.v1.Template.message:type_name -> chacerapp.v1.Message
	5, // 1: chacerapp.config.v1.Template.created_time:type_name -> google.protobuf.Timestamp
	5, // 2: chacerapp.config.v1.Template.updated_time:type_name -> google.protobuf.Timestamp
	0, // 3: chacerapp.config.v1.ListTemplatesResponse.templates:type_name -> chacerapp.config.v1.Template
	0, // 4: chacerapp.config.v1.CreateTemplateRequest.template:type_name -> chacerapp.config.v1.Template
	1, // 5: chacerapp.config.v1.Templates.ListTemplates:input_type -> chacerapp.config.v1.ListTemplatesRequest
	3, // 6: chacerapp.config.v1.Templates.CreateTemplate:input_type -> chacerapp.config.v1.CreateTemplateRequest
	2, // 7: chacerapp.config.v1.Templates.ListTemplates:output_type -> chacerapp.config.v1.ListTemplatesResponse
	0, // 8: chacerapp.config.v1.Templates.CreateTemplate:output_type -> chacerapp.config.v1.Template
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_chacerapp_v1_templates_proto_init() }
func file_chacerapp_v1_templates_proto_init() {
	if File_chacerapp_v1_templates_proto != nil {
		return
	}
	file_chacerapp_v1_messages_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_chacerapp_v1_templates_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Template); i {
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
		file_chacerapp_v1_templates_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListTemplatesRequest); i {
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
		file_chacerapp_v1_templates_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListTemplatesResponse); i {
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
		file_chacerapp_v1_templates_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTemplateRequest); i {
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
			RawDescriptor: file_chacerapp_v1_templates_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chacerapp_v1_templates_proto_goTypes,
		DependencyIndexes: file_chacerapp_v1_templates_proto_depIdxs,
		MessageInfos:      file_chacerapp_v1_templates_proto_msgTypes,
	}.Build()
	File_chacerapp_v1_templates_proto = out.File
	file_chacerapp_v1_templates_proto_rawDesc = nil
	file_chacerapp_v1_templates_proto_goTypes = nil
	file_chacerapp_v1_templates_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// TemplatesClient is the client API for Templates service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TemplatesClient interface {
	// Lists the templates in the system based on provided filters.
	ListTemplates(ctx context.Context, in *ListTemplatesRequest, opts ...grpc.CallOption) (*ListTemplatesResponse, error)
	// Create a new template that can be used to send new messages on the platform. The
	// display name of the template must be unique within it's parent. The template stores
	// the rendering configuration that should be used by default when sending the message.
	CreateTemplate(ctx context.Context, in *CreateTemplateRequest, opts ...grpc.CallOption) (*Template, error)
}

type templatesClient struct {
	cc grpc.ClientConnInterface
}

func NewTemplatesClient(cc grpc.ClientConnInterface) TemplatesClient {
	return &templatesClient{cc}
}

func (c *templatesClient) ListTemplates(ctx context.Context, in *ListTemplatesRequest, opts ...grpc.CallOption) (*ListTemplatesResponse, error) {
	out := new(ListTemplatesResponse)
	err := c.cc.Invoke(ctx, "/chacerapp.config.v1.Templates/ListTemplates", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *templatesClient) CreateTemplate(ctx context.Context, in *CreateTemplateRequest, opts ...grpc.CallOption) (*Template, error) {
	out := new(Template)
	err := c.cc.Invoke(ctx, "/chacerapp.config.v1.Templates/CreateTemplate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TemplatesServer is the server API for Templates service.
type TemplatesServer interface {
	// Lists the templates in the system based on provided filters.
	ListTemplates(context.Context, *ListTemplatesRequest) (*ListTemplatesResponse, error)
	// Create a new template that can be used to send new messages on the platform. The
	// display name of the template must be unique within it's parent. The template stores
	// the rendering configuration that should be used by default when sending the message.
	CreateTemplate(context.Context, *CreateTemplateRequest) (*Template, error)
}

// UnimplementedTemplatesServer can be embedded to have forward compatible implementations.
type UnimplementedTemplatesServer struct {
}

func (*UnimplementedTemplatesServer) ListTemplates(context.Context, *ListTemplatesRequest) (*ListTemplatesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTemplates not implemented")
}
func (*UnimplementedTemplatesServer) CreateTemplate(context.Context, *CreateTemplateRequest) (*Template, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTemplate not implemented")
}

func RegisterTemplatesServer(s *grpc.Server, srv TemplatesServer) {
	s.RegisterService(&_Templates_serviceDesc, srv)
}

func _Templates_ListTemplates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTemplatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TemplatesServer).ListTemplates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chacerapp.config.v1.Templates/ListTemplates",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TemplatesServer).ListTemplates(ctx, req.(*ListTemplatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Templates_CreateTemplate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTemplateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TemplatesServer).CreateTemplate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chacerapp.config.v1.Templates/CreateTemplate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TemplatesServer).CreateTemplate(ctx, req.(*CreateTemplateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Templates_serviceDesc = grpc.ServiceDesc{
	ServiceName: "chacerapp.config.v1.Templates",
	HandlerType: (*TemplatesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListTemplates",
			Handler:    _Templates_ListTemplates_Handler,
		},
		{
			MethodName: "CreateTemplate",
			Handler:    _Templates_CreateTemplate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chacerapp/v1/templates.proto",
}
