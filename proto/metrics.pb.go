// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.1
// source: proto/metrics.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Metrics represents a single metric data point, including its identifier,
// type, and value. It supports both gauge and counter metric types.
type Metrics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Delta is used for counter metrics to represent a change in value.
	// 'optional' in proto3 is the default and need not be explicitly specified.
	Delta *wrapperspb.Int64Value `protobuf:"bytes,1,opt,name=delta,proto3" json:"delta,omitempty"`
	// Value is used for gauge metrics to represent a measurement value.
	Value *wrapperspb.DoubleValue `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	// ID is the unique identifier of the metric.
	Id string `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	// MType is the type of the metric, such as gauge or counter.
	MType string `protobuf:"bytes,4,opt,name=mType,proto3" json:"mType,omitempty"`
}

func (x *Metrics) Reset() {
	*x = Metrics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metrics) ProtoMessage() {}

func (x *Metrics) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metrics.ProtoReflect.Descriptor instead.
func (*Metrics) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{0}
}

func (x *Metrics) GetDelta() *wrapperspb.Int64Value {
	if x != nil {
		return x.Delta
	}
	return nil
}

func (x *Metrics) GetValue() *wrapperspb.DoubleValue {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *Metrics) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Metrics) GetMType() string {
	if x != nil {
		return x.MType
	}
	return ""
}

type MetricsArray struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Objects []*Metrics `protobuf:"bytes,1,rep,name=objects,proto3" json:"objects,omitempty"`
}

func (x *MetricsArray) Reset() {
	*x = MetricsArray{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricsArray) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricsArray) ProtoMessage() {}

func (x *MetricsArray) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricsArray.ProtoReflect.Descriptor instead.
func (*MetricsArray) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{1}
}

func (x *MetricsArray) GetObjects() []*Metrics {
	if x != nil {
		return x.Objects
	}
	return nil
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{2}
}

var File_proto_metrics_proto protoreflect.FileDescriptor

var file_proto_metrics_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x1a, 0x1e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x96,
	0x01, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x31, 0x0a, 0x05, 0x64, 0x65,
	0x6c, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x6e, 0x74, 0x36,
	0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x12, 0x32, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44,
	0x6f, 0x75, 0x62, 0x6c, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x0a, 0x0c, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x41, 0x72, 0x72, 0x61, 0x79, 0x12, 0x2a, 0x0a, 0x07, 0x6f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x07, 0x6f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x73, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x8d, 0x02, 0x0a,
	0x0a, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x26, 0x0a, 0x04, 0x50,
	0x69, 0x6e, 0x67, 0x12, 0x0e, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x0e, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x12, 0x32, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x12, 0x10, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x1a, 0x10, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x2f, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x12, 0x10, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x1a, 0x10, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x3a, 0x0a, 0x11, 0x4d, 0x61, 0x73, 0x73,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x15, 0x2e,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x41,
	0x72, 0x72, 0x61, 0x79, 0x1a, 0x0e, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x12, 0x36, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x0e, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x15, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x41, 0x72, 0x72, 0x61, 0x79, 0x42, 0x12, 0x5a, 0x10,
	0x6d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_metrics_proto_rawDescOnce sync.Once
	file_proto_metrics_proto_rawDescData = file_proto_metrics_proto_rawDesc
)

func file_proto_metrics_proto_rawDescGZIP() []byte {
	file_proto_metrics_proto_rawDescOnce.Do(func() {
		file_proto_metrics_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_metrics_proto_rawDescData)
	})
	return file_proto_metrics_proto_rawDescData
}

var file_proto_metrics_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_metrics_proto_goTypes = []interface{}{
	(*Metrics)(nil),                // 0: metrics.Metrics
	(*MetricsArray)(nil),           // 1: metrics.MetricsArray
	(*Empty)(nil),                  // 2: metrics.Empty
	(*wrapperspb.Int64Value)(nil),  // 3: google.protobuf.Int64Value
	(*wrapperspb.DoubleValue)(nil), // 4: google.protobuf.DoubleValue
}
var file_proto_metrics_proto_depIdxs = []int32{
	3, // 0: metrics.Metrics.delta:type_name -> google.protobuf.Int64Value
	4, // 1: metrics.Metrics.value:type_name -> google.protobuf.DoubleValue
	0, // 2: metrics.MetricsArray.objects:type_name -> metrics.Metrics
	2, // 3: metrics.Monitoring.Ping:input_type -> metrics.Empty
	0, // 4: metrics.Monitoring.UpdateMetric:input_type -> metrics.Metrics
	0, // 5: metrics.Monitoring.GetMetric:input_type -> metrics.Metrics
	1, // 6: metrics.Monitoring.MassUpdateMetrics:input_type -> metrics.MetricsArray
	2, // 7: metrics.Monitoring.GetAllMetrics:input_type -> metrics.Empty
	2, // 8: metrics.Monitoring.Ping:output_type -> metrics.Empty
	0, // 9: metrics.Monitoring.UpdateMetric:output_type -> metrics.Metrics
	0, // 10: metrics.Monitoring.GetMetric:output_type -> metrics.Metrics
	2, // 11: metrics.Monitoring.MassUpdateMetrics:output_type -> metrics.Empty
	1, // 12: metrics.Monitoring.GetAllMetrics:output_type -> metrics.MetricsArray
	8, // [8:13] is the sub-list for method output_type
	3, // [3:8] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_metrics_proto_init() }
func file_proto_metrics_proto_init() {
	if File_proto_metrics_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_metrics_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metrics); i {
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
		file_proto_metrics_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetricsArray); i {
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
		file_proto_metrics_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
			RawDescriptor: file_proto_metrics_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_metrics_proto_goTypes,
		DependencyIndexes: file_proto_metrics_proto_depIdxs,
		MessageInfos:      file_proto_metrics_proto_msgTypes,
	}.Build()
	File_proto_metrics_proto = out.File
	file_proto_metrics_proto_rawDesc = nil
	file_proto_metrics_proto_goTypes = nil
	file_proto_metrics_proto_depIdxs = nil
}
