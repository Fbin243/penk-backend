// Original file: ../../pkg/proto/analytic.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { DeleteCapturedRecordsReq as _pb_DeleteCapturedRecordsReq, DeleteCapturedRecordsReq__Output as _pb_DeleteCapturedRecordsReq__Output } from '../pb/DeleteCapturedRecordsReq';
import type { DeleteCapturedRecordsResp as _pb_DeleteCapturedRecordsResp, DeleteCapturedRecordsResp__Output as _pb_DeleteCapturedRecordsResp__Output } from '../pb/DeleteCapturedRecordsResp';

export interface AnalyticClient extends grpc.Client {
  DeleteCapturedRecords(argument: _pb_DeleteCapturedRecordsReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_DeleteCapturedRecordsResp__Output>): grpc.ClientUnaryCall;
  DeleteCapturedRecords(argument: _pb_DeleteCapturedRecordsReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_pb_DeleteCapturedRecordsResp__Output>): grpc.ClientUnaryCall;
  DeleteCapturedRecords(argument: _pb_DeleteCapturedRecordsReq, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_DeleteCapturedRecordsResp__Output>): grpc.ClientUnaryCall;
  DeleteCapturedRecords(argument: _pb_DeleteCapturedRecordsReq, callback: grpc.requestCallback<_pb_DeleteCapturedRecordsResp__Output>): grpc.ClientUnaryCall;
  deleteCapturedRecords(argument: _pb_DeleteCapturedRecordsReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_DeleteCapturedRecordsResp__Output>): grpc.ClientUnaryCall;
  deleteCapturedRecords(argument: _pb_DeleteCapturedRecordsReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_pb_DeleteCapturedRecordsResp__Output>): grpc.ClientUnaryCall;
  deleteCapturedRecords(argument: _pb_DeleteCapturedRecordsReq, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_DeleteCapturedRecordsResp__Output>): grpc.ClientUnaryCall;
  deleteCapturedRecords(argument: _pb_DeleteCapturedRecordsReq, callback: grpc.requestCallback<_pb_DeleteCapturedRecordsResp__Output>): grpc.ClientUnaryCall;
  
}

export interface AnalyticHandlers extends grpc.UntypedServiceImplementation {
  DeleteCapturedRecords: grpc.handleUnaryCall<_pb_DeleteCapturedRecordsReq__Output, _pb_DeleteCapturedRecordsResp>;
  
}

export interface AnalyticDefinition extends grpc.ServiceDefinition {
  DeleteCapturedRecords: MethodDefinition<_pb_DeleteCapturedRecordsReq, _pb_DeleteCapturedRecordsResp, _pb_DeleteCapturedRecordsReq__Output, _pb_DeleteCapturedRecordsResp__Output>
}
