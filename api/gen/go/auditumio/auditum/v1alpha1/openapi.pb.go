// Copyright 2023 Igor Zibarev
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: auditumio/auditum/v1alpha1/openapi.proto

package auditumv1alpha1

import (
	_ "github.com/auditumio/auditum/api/gen/go/protoc-gen-openapiv2/options"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_auditumio_auditum_v1alpha1_openapi_proto protoreflect.FileDescriptor

var file_auditumio_auditum_v1alpha1_openapi_proto_rawDesc = []byte{
	0x0a, 0x28, 0x61, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x69, 0x6f, 0x2f, 0x61, 0x75, 0x64, 0x69,
	0x74, 0x75, 0x6d, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x6f, 0x70, 0x65,
	0x6e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x61, 0x75, 0x64, 0x69,
	0x74, 0x75, 0x6d, 0x69, 0x6f, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x2e, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67,
	0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x42, 0xce, 0x09, 0x92, 0x41, 0xbe, 0x07, 0x12, 0xf2, 0x02,
	0x0a, 0x0b, 0x41, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x20, 0x41, 0x50, 0x49, 0x12, 0xd8, 0x02,
	0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x70, 0x65, 0x63,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x41, 0x75,
	0x64, 0x69, 0x74, 0x75, 0x6d, 0x20, 0x48, 0x54, 0x54, 0x50, 0x20, 0x41, 0x50, 0x49, 0x2e, 0x0a,
	0x0a, 0x46, 0x6f, 0x72, 0x20, 0x67, 0x52, 0x50, 0x43, 0x20, 0x41, 0x50, 0x49, 0x2c, 0x20, 0x73,
	0x65, 0x65, 0x20, 0x5b, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x5d, 0x28,
	0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x69, 0x6f, 0x2f, 0x61, 0x75, 0x64,
	0x69, 0x74, 0x75, 0x6d, 0x2f, 0x74, 0x72, 0x65, 0x65, 0x2f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d,
	0x69, 0x6f, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x29, 0x2e, 0x0a, 0x0a, 0x46, 0x6f, 0x72, 0x20, 0x6d, 0x6f, 0x72, 0x65, 0x20,
	0x69, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2c, 0x20, 0x73, 0x65, 0x65,
	0x20, 0x5b, 0x55, 0x73, 0x61, 0x67, 0x65, 0x20, 0x47, 0x75, 0x69, 0x64, 0x65, 0x5d, 0x28, 0x2f,
	0x64, 0x6f, 0x63, 0x73, 0x2f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2d, 0x67, 0x75, 0x69, 0x64, 0x65,
	0x2f, 0x29, 0x2e, 0x0a, 0x0a, 0x46, 0x6f, 0x72, 0x20, 0x61, 0x6e, 0x79, 0x20, 0x69, 0x73, 0x73,
	0x75, 0x65, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x20,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x2c, 0x20, 0x70, 0x6c, 0x65, 0x61, 0x73, 0x65,
	0x20, 0x75, 0x73, 0x65, 0x20, 0x5b, 0x47, 0x69, 0x74, 0x48, 0x75, 0x62, 0x20, 0x49, 0x73, 0x73,
	0x75, 0x65, 0x20, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x5d, 0x28, 0x68, 0x74, 0x74, 0x70,
	0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61,
	0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x69, 0x6f, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d,
	0x2f, 0x69, 0x73, 0x73, 0x75, 0x65, 0x73, 0x29, 0x32, 0x08, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x22, 0x0d, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x32, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a,
	0x73, 0x6f, 0x6e, 0x3a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x3a, 0x17, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x2b, 0x70, 0x72, 0x65, 0x74, 0x74, 0x79, 0x6a, 0xbf,
	0x02, 0x0a, 0x08, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x12, 0xef, 0x01, 0x2a, 0x2a,
	0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2a, 0x2a, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68, 0x65,
	0x20, 0x74, 0x6f, 0x70, 0x2d, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x20, 0x69, 0x6e, 0x20, 0x2a, 0x41, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x2a,
	0x2e, 0x20, 0x2a, 0x2a, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2a, 0x2a, 0x20, 0x63,
	0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x20, 0x61, 0x6c, 0x6c, 0x20, 0x6f, 0x74, 0x68, 0x65, 0x72,
	0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x20, 0x6c, 0x69, 0x6b, 0x65, 0x20,
	0x2a, 0x2a, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x2a, 0x2a, 0x2e, 0x0a, 0x0a, 0x41, 0x20,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x20, 0x69, 0x73, 0x20, 0x61, 0x20, 0x6c, 0x6f, 0x67,
	0x69, 0x63, 0x61, 0x6c, 0x20, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x69, 0x6e, 0x67, 0x20, 0x6f, 0x66,
	0x20, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x2e, 0x20, 0x54, 0x79, 0x70, 0x69, 0x63, 0x61,
	0x6c, 0x6c, 0x79, 0x2c, 0x20, 0x61, 0x20, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x20, 0x69,
	0x73, 0x20, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x65, 0x61,
	0x63, 0x68, 0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x20, 0x74, 0x68, 0x61, 0x74, 0x20, 0x69, 0x73, 0x20, 0x62,
	0x65, 0x69, 0x6e, 0x67, 0x20, 0x61, 0x75, 0x64, 0x69, 0x74, 0x65, 0x64, 0x2e, 0x1a, 0x41, 0x0a,
	0x1d, 0x55, 0x73, 0x61, 0x67, 0x65, 0x20, 0x47, 0x75, 0x69, 0x64, 0x65, 0x20, 0x3a, 0x3a, 0x20,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x20, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x20,
	0x2f, 0x64, 0x6f, 0x63, 0x73, 0x2f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2d, 0x67, 0x75, 0x69, 0x64,
	0x65, 0x2f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x6a, 0xb8, 0x01, 0x0a, 0x07, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x6a, 0x2a, 0x2a,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x2a, 0x2a, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68, 0x65, 0x20,
	0x63, 0x6f, 0x72, 0x65, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x20, 0x69, 0x6e,
	0x20, 0x2a, 0x41, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x2a, 0x2e, 0x20, 0x2a, 0x2a, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x73, 0x2a, 0x2a, 0x20, 0x72, 0x65, 0x70, 0x72, 0x65, 0x73, 0x65, 0x6e,
	0x74, 0x20, 0x61, 0x75, 0x64, 0x69, 0x74, 0x20, 0x74, 0x72, 0x61, 0x69, 0x6c, 0x20, 0x72, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x73, 0x20, 0x61, 0x2e, 0x6b, 0x2e, 0x61, 0x2e, 0x20, 0x61, 0x75, 0x64,
	0x69, 0x74, 0x20, 0x6c, 0x6f, 0x67, 0x73, 0x2e, 0x1a, 0x41, 0x0a, 0x1d, 0x55, 0x73, 0x61, 0x67,
	0x65, 0x20, 0x47, 0x75, 0x69, 0x64, 0x65, 0x20, 0x3a, 0x3a, 0x20, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x20, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x20, 0x2f, 0x64, 0x6f, 0x63, 0x73,
	0x2f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2d, 0x67, 0x75, 0x69, 0x64, 0x65, 0x2f, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x2d, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x0a, 0x1e, 0x63, 0x6f, 0x6d,
	0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x69, 0x6f, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74,
	0x75, 0x6d, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x42, 0x0c, 0x4f, 0x70, 0x65,
	0x6e, 0x61, 0x70, 0x69, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x52, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x69,
	0x6f, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65,
	0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x69, 0x6f, 0x2f, 0x61,
	0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x3b,
	0x61, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xa2,
	0x02, 0x03, 0x41, 0x41, 0x58, 0xaa, 0x02, 0x1a, 0x41, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x69,
	0x6f, 0x2e, 0x41, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x2e, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0xca, 0x02, 0x1a, 0x41, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x69, 0x6f, 0x5c, 0x41,
	0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xe2,
	0x02, 0x26, 0x41, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x69, 0x6f, 0x5c, 0x41, 0x75, 0x64, 0x69,
	0x74, 0x75, 0x6d, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1c, 0x41, 0x75, 0x64, 0x69, 0x74,
	0x75, 0x6d, 0x69, 0x6f, 0x3a, 0x3a, 0x41, 0x75, 0x64, 0x69, 0x74, 0x75, 0x6d, 0x3a, 0x3a, 0x56,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_auditumio_auditum_v1alpha1_openapi_proto_goTypes = []any{}
var file_auditumio_auditum_v1alpha1_openapi_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_auditumio_auditum_v1alpha1_openapi_proto_init() }
func file_auditumio_auditum_v1alpha1_openapi_proto_init() {
	if File_auditumio_auditum_v1alpha1_openapi_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_auditumio_auditum_v1alpha1_openapi_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_auditumio_auditum_v1alpha1_openapi_proto_goTypes,
		DependencyIndexes: file_auditumio_auditum_v1alpha1_openapi_proto_depIdxs,
	}.Build()
	File_auditumio_auditum_v1alpha1_openapi_proto = out.File
	file_auditumio_auditum_v1alpha1_openapi_proto_rawDesc = nil
	file_auditumio_auditum_v1alpha1_openapi_proto_goTypes = nil
	file_auditumio_auditum_v1alpha1_openapi_proto_depIdxs = nil
}
