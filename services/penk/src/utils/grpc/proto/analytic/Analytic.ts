// Original file: ../../proto/analytic/analytic_service.proto

import type * as grpc from "@grpc/grpc-js";
import type { MethodDefinition } from "@grpc/proto-loader";
import type {
  DeleteCapturedRecordsReq as _analytic_DeleteCapturedRecordsReq,
  DeleteCapturedRecordsReq__Output as _analytic_DeleteCapturedRecordsReq__Output,
} from "./DeleteCapturedRecordsReq";
import type {
  DeleteCapturedRecordsResp as _analytic_DeleteCapturedRecordsResp,
  DeleteCapturedRecordsResp__Output as _analytic_DeleteCapturedRecordsResp__Output,
} from "./DeleteCapturedRecordsResp";

export interface AnalyticClient extends grpc.Client {
  DeleteCapturedRecords(
    argument: _analytic_DeleteCapturedRecordsReq,
    metadata: grpc.Metadata,
    options: grpc.CallOptions,
    callback: grpc.requestCallback<_analytic_DeleteCapturedRecordsResp__Output>,
  ): grpc.ClientUnaryCall;
  DeleteCapturedRecords(
    argument: _analytic_DeleteCapturedRecordsReq,
    metadata: grpc.Metadata,
    callback: grpc.requestCallback<_analytic_DeleteCapturedRecordsResp__Output>,
  ): grpc.ClientUnaryCall;
  DeleteCapturedRecords(
    argument: _analytic_DeleteCapturedRecordsReq,
    options: grpc.CallOptions,
    callback: grpc.requestCallback<_analytic_DeleteCapturedRecordsResp__Output>,
  ): grpc.ClientUnaryCall;
  DeleteCapturedRecords(
    argument: _analytic_DeleteCapturedRecordsReq,
    callback: grpc.requestCallback<_analytic_DeleteCapturedRecordsResp__Output>,
  ): grpc.ClientUnaryCall;
  deleteCapturedRecords(
    argument: _analytic_DeleteCapturedRecordsReq,
    metadata: grpc.Metadata,
    options: grpc.CallOptions,
    callback: grpc.requestCallback<_analytic_DeleteCapturedRecordsResp__Output>,
  ): grpc.ClientUnaryCall;
  deleteCapturedRecords(
    argument: _analytic_DeleteCapturedRecordsReq,
    metadata: grpc.Metadata,
    callback: grpc.requestCallback<_analytic_DeleteCapturedRecordsResp__Output>,
  ): grpc.ClientUnaryCall;
  deleteCapturedRecords(
    argument: _analytic_DeleteCapturedRecordsReq,
    options: grpc.CallOptions,
    callback: grpc.requestCallback<_analytic_DeleteCapturedRecordsResp__Output>,
  ): grpc.ClientUnaryCall;
  deleteCapturedRecords(
    argument: _analytic_DeleteCapturedRecordsReq,
    callback: grpc.requestCallback<_analytic_DeleteCapturedRecordsResp__Output>,
  ): grpc.ClientUnaryCall;
}

export interface AnalyticHandlers extends grpc.UntypedServiceImplementation {
  DeleteCapturedRecords: grpc.handleUnaryCall<
    _analytic_DeleteCapturedRecordsReq__Output,
    _analytic_DeleteCapturedRecordsResp
  >;
}

export interface AnalyticDefinition extends grpc.ServiceDefinition {
  DeleteCapturedRecords: MethodDefinition<
    _analytic_DeleteCapturedRecordsReq,
    _analytic_DeleteCapturedRecordsResp,
    _analytic_DeleteCapturedRecordsReq__Output,
    _analytic_DeleteCapturedRecordsResp__Output
  >;
}
