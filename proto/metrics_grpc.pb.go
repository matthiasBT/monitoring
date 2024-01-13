// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: proto/metrics.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Monitoring_Ping_FullMethodName                       = "/metrics.Monitoring/Ping"
	Monitoring_UpdateMetric_FullMethodName               = "/metrics.Monitoring/UpdateMetric"
	Monitoring_GetMetric_FullMethodName                  = "/metrics.Monitoring/GetMetric"
	Monitoring_MassUpdateMetrics_FullMethodName          = "/metrics.Monitoring/MassUpdateMetrics"
	Monitoring_MassUpdateMetricsEncrypted_FullMethodName = "/metrics.Monitoring/MassUpdateMetricsEncrypted"
	Monitoring_GetAllMetrics_FullMethodName              = "/metrics.Monitoring/GetAllMetrics"
)

// MonitoringClient is the client API for Monitoring service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MonitoringClient interface {
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	UpdateMetric(ctx context.Context, in *Metrics, opts ...grpc.CallOption) (*Metrics, error)
	GetMetric(ctx context.Context, in *Metrics, opts ...grpc.CallOption) (*Metrics, error)
	MassUpdateMetrics(ctx context.Context, in *MetricsArray, opts ...grpc.CallOption) (*Empty, error)
	MassUpdateMetricsEncrypted(ctx context.Context, in *EncryptedMetricsArray, opts ...grpc.CallOption) (*Empty, error)
	GetAllMetrics(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*MetricsArray, error)
}

type monitoringClient struct {
	cc grpc.ClientConnInterface
}

func NewMonitoringClient(cc grpc.ClientConnInterface) MonitoringClient {
	return &monitoringClient{cc}
}

func (c *monitoringClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Monitoring_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *monitoringClient) UpdateMetric(ctx context.Context, in *Metrics, opts ...grpc.CallOption) (*Metrics, error) {
	out := new(Metrics)
	err := c.cc.Invoke(ctx, Monitoring_UpdateMetric_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *monitoringClient) GetMetric(ctx context.Context, in *Metrics, opts ...grpc.CallOption) (*Metrics, error) {
	out := new(Metrics)
	err := c.cc.Invoke(ctx, Monitoring_GetMetric_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *monitoringClient) MassUpdateMetrics(ctx context.Context, in *MetricsArray, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Monitoring_MassUpdateMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *monitoringClient) MassUpdateMetricsEncrypted(ctx context.Context, in *EncryptedMetricsArray, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Monitoring_MassUpdateMetricsEncrypted_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *monitoringClient) GetAllMetrics(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*MetricsArray, error) {
	out := new(MetricsArray)
	err := c.cc.Invoke(ctx, Monitoring_GetAllMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MonitoringServer is the server API for Monitoring service.
// All implementations must embed UnimplementedMonitoringServer
// for forward compatibility
type MonitoringServer interface {
	Ping(context.Context, *Empty) (*Empty, error)
	UpdateMetric(context.Context, *Metrics) (*Metrics, error)
	GetMetric(context.Context, *Metrics) (*Metrics, error)
	MassUpdateMetrics(context.Context, *MetricsArray) (*Empty, error)
	MassUpdateMetricsEncrypted(context.Context, *EncryptedMetricsArray) (*Empty, error)
	GetAllMetrics(context.Context, *Empty) (*MetricsArray, error)
	mustEmbedUnimplementedMonitoringServer()
}

// UnimplementedMonitoringServer must be embedded to have forward compatible implementations.
type UnimplementedMonitoringServer struct {
}

func (UnimplementedMonitoringServer) Ping(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedMonitoringServer) UpdateMetric(context.Context, *Metrics) (*Metrics, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetric not implemented")
}
func (UnimplementedMonitoringServer) GetMetric(context.Context, *Metrics) (*Metrics, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetric not implemented")
}
func (UnimplementedMonitoringServer) MassUpdateMetrics(context.Context, *MetricsArray) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MassUpdateMetrics not implemented")
}
func (UnimplementedMonitoringServer) MassUpdateMetricsEncrypted(context.Context, *EncryptedMetricsArray) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MassUpdateMetricsEncrypted not implemented")
}
func (UnimplementedMonitoringServer) GetAllMetrics(context.Context, *Empty) (*MetricsArray, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllMetrics not implemented")
}
func (UnimplementedMonitoringServer) mustEmbedUnimplementedMonitoringServer() {}

// UnsafeMonitoringServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MonitoringServer will
// result in compilation errors.
type UnsafeMonitoringServer interface {
	mustEmbedUnimplementedMonitoringServer()
}

func RegisterMonitoringServer(s grpc.ServiceRegistrar, srv MonitoringServer) {
	s.RegisterService(&Monitoring_ServiceDesc, srv)
}

func _Monitoring_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitoringServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Monitoring_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitoringServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Monitoring_UpdateMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Metrics)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitoringServer).UpdateMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Monitoring_UpdateMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitoringServer).UpdateMetric(ctx, req.(*Metrics))
	}
	return interceptor(ctx, in, info, handler)
}

func _Monitoring_GetMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Metrics)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitoringServer).GetMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Monitoring_GetMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitoringServer).GetMetric(ctx, req.(*Metrics))
	}
	return interceptor(ctx, in, info, handler)
}

func _Monitoring_MassUpdateMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetricsArray)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitoringServer).MassUpdateMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Monitoring_MassUpdateMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitoringServer).MassUpdateMetrics(ctx, req.(*MetricsArray))
	}
	return interceptor(ctx, in, info, handler)
}

func _Monitoring_MassUpdateMetricsEncrypted_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EncryptedMetricsArray)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitoringServer).MassUpdateMetricsEncrypted(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Monitoring_MassUpdateMetricsEncrypted_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitoringServer).MassUpdateMetricsEncrypted(ctx, req.(*EncryptedMetricsArray))
	}
	return interceptor(ctx, in, info, handler)
}

func _Monitoring_GetAllMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitoringServer).GetAllMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Monitoring_GetAllMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitoringServer).GetAllMetrics(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Monitoring_ServiceDesc is the grpc.ServiceDesc for Monitoring service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Monitoring_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "metrics.Monitoring",
	HandlerType: (*MonitoringServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Monitoring_Ping_Handler,
		},
		{
			MethodName: "UpdateMetric",
			Handler:    _Monitoring_UpdateMetric_Handler,
		},
		{
			MethodName: "GetMetric",
			Handler:    _Monitoring_GetMetric_Handler,
		},
		{
			MethodName: "MassUpdateMetrics",
			Handler:    _Monitoring_MassUpdateMetrics_Handler,
		},
		{
			MethodName: "MassUpdateMetricsEncrypted",
			Handler:    _Monitoring_MassUpdateMetricsEncrypted_Handler,
		},
		{
			MethodName: "GetAllMetrics",
			Handler:    _Monitoring_GetAllMetrics_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/metrics.proto",
}
