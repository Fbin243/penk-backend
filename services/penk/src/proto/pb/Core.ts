// Original file: ../../pkg/proto/core.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { CheckPermissionReq as _pb_CheckPermissionReq, CheckPermissionReq__Output as _pb_CheckPermissionReq__Output } from '../pb/CheckPermissionReq';
import type { CheckPermissionResp as _pb_CheckPermissionResp, CheckPermissionResp__Output as _pb_CheckPermissionResp__Output } from '../pb/CheckPermissionResp';
import type { IntrospectReq as _pb_IntrospectReq, IntrospectReq__Output as _pb_IntrospectReq__Output } from '../pb/IntrospectReq';
import type { IntrospectResp as _pb_IntrospectResp, IntrospectResp__Output as _pb_IntrospectResp__Output } from '../pb/IntrospectResp';

export interface CoreClient extends grpc.Client {
  CheckPermission(argument: _pb_CheckPermissionReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  CheckPermission(argument: _pb_CheckPermissionReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_pb_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  CheckPermission(argument: _pb_CheckPermissionReq, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  CheckPermission(argument: _pb_CheckPermissionReq, callback: grpc.requestCallback<_pb_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  checkPermission(argument: _pb_CheckPermissionReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  checkPermission(argument: _pb_CheckPermissionReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_pb_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  checkPermission(argument: _pb_CheckPermissionReq, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  checkPermission(argument: _pb_CheckPermissionReq, callback: grpc.requestCallback<_pb_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  
  IntrospectProfile(argument: _pb_IntrospectReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_IntrospectResp__Output>): grpc.ClientUnaryCall;
  IntrospectProfile(argument: _pb_IntrospectReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_pb_IntrospectResp__Output>): grpc.ClientUnaryCall;
  IntrospectProfile(argument: _pb_IntrospectReq, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_IntrospectResp__Output>): grpc.ClientUnaryCall;
  IntrospectProfile(argument: _pb_IntrospectReq, callback: grpc.requestCallback<_pb_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectProfile(argument: _pb_IntrospectReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectProfile(argument: _pb_IntrospectReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_pb_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectProfile(argument: _pb_IntrospectReq, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectProfile(argument: _pb_IntrospectReq, callback: grpc.requestCallback<_pb_IntrospectResp__Output>): grpc.ClientUnaryCall;
  
}

export interface CoreHandlers extends grpc.UntypedServiceImplementation {
  CheckPermission: grpc.handleUnaryCall<_pb_CheckPermissionReq__Output, _pb_CheckPermissionResp>;
  
  IntrospectProfile: grpc.handleUnaryCall<_pb_IntrospectReq__Output, _pb_IntrospectResp>;
  
}

export interface CoreDefinition extends grpc.ServiceDefinition {
  CheckPermission: MethodDefinition<_pb_CheckPermissionReq, _pb_CheckPermissionResp, _pb_CheckPermissionReq__Output, _pb_CheckPermissionResp__Output>
  IntrospectProfile: MethodDefinition<_pb_IntrospectReq, _pb_IntrospectResp, _pb_IntrospectReq__Output, _pb_IntrospectResp__Output>
}
