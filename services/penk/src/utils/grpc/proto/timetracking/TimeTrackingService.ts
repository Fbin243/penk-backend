// Original file: ../../proto/timetracking/timetracking_service.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { CreateTimeTrackingRequest as _timetracking_CreateTimeTrackingRequest, CreateTimeTrackingRequest__Output as _timetracking_CreateTimeTrackingRequest__Output } from '../timetracking/CreateTimeTrackingRequest';
import type { EmptyReq as _common_EmptyReq, EmptyReq__Output as _common_EmptyReq__Output } from '../common/EmptyReq';
import type { TimeTracking as _timetracking_TimeTracking, TimeTracking__Output as _timetracking_TimeTracking__Output } from '../timetracking/TimeTracking';
import type { TimeTrackingWithFish as _timetracking_TimeTrackingWithFish, TimeTrackingWithFish__Output as _timetracking_TimeTrackingWithFish__Output } from '../timetracking/TimeTrackingWithFish';
import type { TotalTimeTrackingRequest as _timetracking_TotalTimeTrackingRequest, TotalTimeTrackingRequest__Output as _timetracking_TotalTimeTrackingRequest__Output } from '../timetracking/TotalTimeTrackingRequest';
import type { TotalTimeTrackingResponse as _timetracking_TotalTimeTrackingResponse, TotalTimeTrackingResponse__Output as _timetracking_TotalTimeTrackingResponse__Output } from '../timetracking/TotalTimeTrackingResponse';

export interface TimeTrackingServiceClient extends grpc.Client {
  CreateTimeTracking(argument: _timetracking_CreateTimeTrackingRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  CreateTimeTracking(argument: _timetracking_CreateTimeTrackingRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  CreateTimeTracking(argument: _timetracking_CreateTimeTrackingRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  CreateTimeTracking(argument: _timetracking_CreateTimeTrackingRequest, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  createTimeTracking(argument: _timetracking_CreateTimeTrackingRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  createTimeTracking(argument: _timetracking_CreateTimeTrackingRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  createTimeTracking(argument: _timetracking_CreateTimeTrackingRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  createTimeTracking(argument: _timetracking_CreateTimeTrackingRequest, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  
  GetCurrentTimeTracking(argument: _common_EmptyReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  GetCurrentTimeTracking(argument: _common_EmptyReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  GetCurrentTimeTracking(argument: _common_EmptyReq, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  GetCurrentTimeTracking(argument: _common_EmptyReq, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  getCurrentTimeTracking(argument: _common_EmptyReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  getCurrentTimeTracking(argument: _common_EmptyReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  getCurrentTimeTracking(argument: _common_EmptyReq, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  getCurrentTimeTracking(argument: _common_EmptyReq, callback: grpc.requestCallback<_timetracking_TimeTracking__Output>): grpc.ClientUnaryCall;
  
  GetTotalCurrentTimeTracking(argument: _timetracking_TotalTimeTrackingRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TotalTimeTrackingResponse__Output>): grpc.ClientUnaryCall;
  GetTotalCurrentTimeTracking(argument: _timetracking_TotalTimeTrackingRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_timetracking_TotalTimeTrackingResponse__Output>): grpc.ClientUnaryCall;
  GetTotalCurrentTimeTracking(argument: _timetracking_TotalTimeTrackingRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TotalTimeTrackingResponse__Output>): grpc.ClientUnaryCall;
  GetTotalCurrentTimeTracking(argument: _timetracking_TotalTimeTrackingRequest, callback: grpc.requestCallback<_timetracking_TotalTimeTrackingResponse__Output>): grpc.ClientUnaryCall;
  getTotalCurrentTimeTracking(argument: _timetracking_TotalTimeTrackingRequest, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TotalTimeTrackingResponse__Output>): grpc.ClientUnaryCall;
  getTotalCurrentTimeTracking(argument: _timetracking_TotalTimeTrackingRequest, metadata: grpc.Metadata, callback: grpc.requestCallback<_timetracking_TotalTimeTrackingResponse__Output>): grpc.ClientUnaryCall;
  getTotalCurrentTimeTracking(argument: _timetracking_TotalTimeTrackingRequest, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TotalTimeTrackingResponse__Output>): grpc.ClientUnaryCall;
  getTotalCurrentTimeTracking(argument: _timetracking_TotalTimeTrackingRequest, callback: grpc.requestCallback<_timetracking_TotalTimeTrackingResponse__Output>): grpc.ClientUnaryCall;
  
  UpdateTimeTracking(argument: _common_EmptyReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TimeTrackingWithFish__Output>): grpc.ClientUnaryCall;
  UpdateTimeTracking(argument: _common_EmptyReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_timetracking_TimeTrackingWithFish__Output>): grpc.ClientUnaryCall;
  UpdateTimeTracking(argument: _common_EmptyReq, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TimeTrackingWithFish__Output>): grpc.ClientUnaryCall;
  UpdateTimeTracking(argument: _common_EmptyReq, callback: grpc.requestCallback<_timetracking_TimeTrackingWithFish__Output>): grpc.ClientUnaryCall;
  updateTimeTracking(argument: _common_EmptyReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TimeTrackingWithFish__Output>): grpc.ClientUnaryCall;
  updateTimeTracking(argument: _common_EmptyReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_timetracking_TimeTrackingWithFish__Output>): grpc.ClientUnaryCall;
  updateTimeTracking(argument: _common_EmptyReq, options: grpc.CallOptions, callback: grpc.requestCallback<_timetracking_TimeTrackingWithFish__Output>): grpc.ClientUnaryCall;
  updateTimeTracking(argument: _common_EmptyReq, callback: grpc.requestCallback<_timetracking_TimeTrackingWithFish__Output>): grpc.ClientUnaryCall;
  
}

export interface TimeTrackingServiceHandlers extends grpc.UntypedServiceImplementation {
  CreateTimeTracking: grpc.handleUnaryCall<_timetracking_CreateTimeTrackingRequest__Output, _timetracking_TimeTracking>;
  
  GetCurrentTimeTracking: grpc.handleUnaryCall<_common_EmptyReq__Output, _timetracking_TimeTracking>;
  
  GetTotalCurrentTimeTracking: grpc.handleUnaryCall<_timetracking_TotalTimeTrackingRequest__Output, _timetracking_TotalTimeTrackingResponse>;
  
  UpdateTimeTracking: grpc.handleUnaryCall<_common_EmptyReq__Output, _timetracking_TimeTrackingWithFish>;
  
}

export interface TimeTrackingServiceDefinition extends grpc.ServiceDefinition {
  CreateTimeTracking: MethodDefinition<_timetracking_CreateTimeTrackingRequest, _timetracking_TimeTracking, _timetracking_CreateTimeTrackingRequest__Output, _timetracking_TimeTracking__Output>
  GetCurrentTimeTracking: MethodDefinition<_common_EmptyReq, _timetracking_TimeTracking, _common_EmptyReq__Output, _timetracking_TimeTracking__Output>
  GetTotalCurrentTimeTracking: MethodDefinition<_timetracking_TotalTimeTrackingRequest, _timetracking_TotalTimeTrackingResponse, _timetracking_TotalTimeTrackingRequest__Output, _timetracking_TotalTimeTrackingResponse__Output>
  UpdateTimeTracking: MethodDefinition<_common_EmptyReq, _timetracking_TimeTrackingWithFish, _common_EmptyReq__Output, _timetracking_TimeTrackingWithFish__Output>
}
