// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: p2p.proto

package p2p_pb

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

type GetPeersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetPeersRequest) Reset() {
	*x = GetPeersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p2p_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPeersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPeersRequest) ProtoMessage() {}

func (x *GetPeersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_p2p_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPeersRequest.ProtoReflect.Descriptor instead.
func (*GetPeersRequest) Descriptor() ([]byte, []int) {
	return file_p2p_proto_rawDescGZIP(), []int{0}
}

type Attribute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *Attribute) Reset() {
	*x = Attribute{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p2p_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Attribute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attribute) ProtoMessage() {}

func (x *Attribute) ProtoReflect() protoreflect.Message {
	mi := &file_p2p_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Attribute.ProtoReflect.Descriptor instead.
func (*Attribute) Descriptor() ([]byte, []int) {
	return file_p2p_proto_rawDescGZIP(), []int{1}
}

func (x *Attribute) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Attribute) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type Peer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address    string       `protobuf:"bytes,1,opt,name=Address,proto3" json:"Address,omitempty"`
	Attributes []*Attribute `protobuf:"bytes,2,rep,name=Attributes,proto3" json:"Attributes,omitempty"`
	State      string       `protobuf:"bytes,3,opt,name=State,proto3" json:"State,omitempty"`
}

func (x *Peer) Reset() {
	*x = Peer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p2p_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Peer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Peer) ProtoMessage() {}

func (x *Peer) ProtoReflect() protoreflect.Message {
	mi := &file_p2p_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Peer.ProtoReflect.Descriptor instead.
func (*Peer) Descriptor() ([]byte, []int) {
	return file_p2p_proto_rawDescGZIP(), []int{2}
}

func (x *Peer) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Peer) GetAttributes() []*Attribute {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *Peer) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

type GetPeersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Peers []*Peer `protobuf:"bytes,1,rep,name=Peers,proto3" json:"Peers,omitempty"`
}

func (x *GetPeersResponse) Reset() {
	*x = GetPeersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p2p_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPeersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPeersResponse) ProtoMessage() {}

func (x *GetPeersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_p2p_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPeersResponse.ProtoReflect.Descriptor instead.
func (*GetPeersResponse) Descriptor() ([]byte, []int) {
	return file_p2p_proto_rawDescGZIP(), []int{3}
}

func (x *GetPeersResponse) GetPeers() []*Peer {
	if x != nil {
		return x.Peers
	}
	return nil
}

var File_p2p_proto protoreflect.FileDescriptor

var file_p2p_proto_rawDesc = []byte{
	0x0a, 0x09, 0x70, 0x32, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x70, 0x32, 0x70,
	0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x11, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x50, 0x65, 0x65,
	0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x33, 0x0a, 0x09, 0x41, 0x74, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x6c,
	0x0a, 0x04, 0x50, 0x65, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x34, 0x0a, 0x0a, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x32, 0x70, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x52, 0x0a, 0x41, 0x74, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x22, 0x39, 0x0a, 0x10,
	0x47, 0x65, 0x74, 0x50, 0x65, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x25, 0x0a, 0x05, 0x50, 0x65, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x0f, 0x2e, 0x70, 0x32, 0x70, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x65, 0x65, 0x72,
	0x52, 0x05, 0x50, 0x65, 0x65, 0x72, 0x73, 0x32, 0x52, 0x0a, 0x0b, 0x50, 0x65, 0x65, 0x72, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x50, 0x65, 0x65,
	0x72, 0x73, 0x12, 0x1a, 0x2e, 0x70, 0x32, 0x70, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47,
	0x65, 0x74, 0x50, 0x65, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b,
	0x2e, 0x70, 0x32, 0x70, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x65,
	0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x7e, 0x0a, 0x0d, 0x63,
	0x6f, 0x6d, 0x2e, 0x70, 0x32, 0x70, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x42, 0x08, 0x50, 0x32,
	0x70, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x72, 0x2d, 0x73, 0x68, 0x69, 0x66, 0x75, 0x2f, 0x67, 0x72,
	0x70, 0x63, 0x2d, 0x70, 0x32, 0x70, 0x2f, 0x70, 0x32, 0x70, 0x5f, 0x70, 0x62, 0xa2, 0x02, 0x03,
	0x50, 0x58, 0x58, 0xaa, 0x02, 0x08, 0x50, 0x32, 0x70, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0xca, 0x02,
	0x08, 0x50, 0x32, 0x70, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0xe2, 0x02, 0x14, 0x50, 0x32, 0x70, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x08, 0x50, 0x32, 0x70, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_p2p_proto_rawDescOnce sync.Once
	file_p2p_proto_rawDescData = file_p2p_proto_rawDesc
)

func file_p2p_proto_rawDescGZIP() []byte {
	file_p2p_proto_rawDescOnce.Do(func() {
		file_p2p_proto_rawDescData = protoimpl.X.CompressGZIP(file_p2p_proto_rawDescData)
	})
	return file_p2p_proto_rawDescData
}

var file_p2p_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_p2p_proto_goTypes = []interface{}{
	(*GetPeersRequest)(nil),  // 0: p2p_proto.GetPeersRequest
	(*Attribute)(nil),        // 1: p2p_proto.Attribute
	(*Peer)(nil),             // 2: p2p_proto.Peer
	(*GetPeersResponse)(nil), // 3: p2p_proto.GetPeersResponse
}
var file_p2p_proto_depIdxs = []int32{
	1, // 0: p2p_proto.Peer.Attributes:type_name -> p2p_proto.Attribute
	2, // 1: p2p_proto.GetPeersResponse.Peers:type_name -> p2p_proto.Peer
	0, // 2: p2p_proto.PeerService.GetPeers:input_type -> p2p_proto.GetPeersRequest
	3, // 3: p2p_proto.PeerService.GetPeers:output_type -> p2p_proto.GetPeersResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_p2p_proto_init() }
func file_p2p_proto_init() {
	if File_p2p_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_p2p_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPeersRequest); i {
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
		file_p2p_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Attribute); i {
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
		file_p2p_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Peer); i {
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
		file_p2p_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPeersResponse); i {
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
			RawDescriptor: file_p2p_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_p2p_proto_goTypes,
		DependencyIndexes: file_p2p_proto_depIdxs,
		MessageInfos:      file_p2p_proto_msgTypes,
	}.Build()
	File_p2p_proto = out.File
	file_p2p_proto_rawDesc = nil
	file_p2p_proto_goTypes = nil
	file_p2p_proto_depIdxs = nil
}
