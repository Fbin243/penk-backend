// Original file: ../../pkg/proto/currency.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { CatchFishReq as _pb_CatchFishReq, CatchFishReq__Output as _pb_CatchFishReq__Output } from '../pb/CatchFishReq';
import type { CatchFishResp as _pb_CatchFishResp, CatchFishResp__Output as _pb_CatchFishResp__Output } from '../pb/CatchFishResp';
import type { CreateFishReq as _pb_CreateFishReq, CreateFishReq__Output as _pb_CreateFishReq__Output } from '../pb/CreateFishReq';
import type { CreateFishResp as _pb_CreateFishResp, CreateFishResp__Output as _pb_CreateFishResp__Output } from '../pb/CreateFishResp';
import type { UpdateFishReq as _pb_UpdateFishReq, UpdateFishReq__Output as _pb_UpdateFishReq__Output } from '../pb/UpdateFishReq';
import type { UpdateFishResp as _pb_UpdateFishResp, UpdateFishResp__Output as _pb_UpdateFishResp__Output } from '../pb/UpdateFishResp';

export interface CurrencyClient extends grpc.Client {
  CatchFish(argument: _pb_CatchFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_CatchFishResp__Output>): grpc.ClientUnaryCall;
  CatchFish(argument: _pb_CatchFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_pb_CatchFishResp__Output>): grpc.ClientUnaryCall;
  CatchFish(argument: _pb_CatchFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_CatchFishResp__Output>): grpc.ClientUnaryCall;
  CatchFish(argument: _pb_CatchFishReq, callback: grpc.requestCallback<_pb_CatchFishResp__Output>): grpc.ClientUnaryCall;
  catchFish(argument: _pb_CatchFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_CatchFishResp__Output>): grpc.ClientUnaryCall;
  catchFish(argument: _pb_CatchFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_pb_CatchFishResp__Output>): grpc.ClientUnaryCall;
  catchFish(argument: _pb_CatchFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_CatchFishResp__Output>): grpc.ClientUnaryCall;
  catchFish(argument: _pb_CatchFishReq, callback: grpc.requestCallback<_pb_CatchFishResp__Output>): grpc.ClientUnaryCall;
  
  CreateFish(argument: _pb_CreateFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_CreateFishResp__Output>): grpc.ClientUnaryCall;
  CreateFish(argument: _pb_CreateFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_pb_CreateFishResp__Output>): grpc.ClientUnaryCall;
  CreateFish(argument: _pb_CreateFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_CreateFishResp__Output>): grpc.ClientUnaryCall;
  CreateFish(argument: _pb_CreateFishReq, callback: grpc.requestCallback<_pb_CreateFishResp__Output>): grpc.ClientUnaryCall;
  createFish(argument: _pb_CreateFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_CreateFishResp__Output>): grpc.ClientUnaryCall;
  createFish(argument: _pb_CreateFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_pb_CreateFishResp__Output>): grpc.ClientUnaryCall;
  createFish(argument: _pb_CreateFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_CreateFishResp__Output>): grpc.ClientUnaryCall;
  createFish(argument: _pb_CreateFishReq, callback: grpc.requestCallback<_pb_CreateFishResp__Output>): grpc.ClientUnaryCall;
  
  UpdateFish(argument: _pb_UpdateFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  UpdateFish(argument: _pb_UpdateFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_pb_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  UpdateFish(argument: _pb_UpdateFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  UpdateFish(argument: _pb_UpdateFishReq, callback: grpc.requestCallback<_pb_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  updateFish(argument: _pb_UpdateFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  updateFish(argument: _pb_UpdateFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_pb_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  updateFish(argument: _pb_UpdateFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_pb_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  updateFish(argument: _pb_UpdateFishReq, callback: grpc.requestCallback<_pb_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  
}

export interface CurrencyHandlers extends grpc.UntypedServiceImplementation {
  CatchFish: grpc.handleUnaryCall<_pb_CatchFishReq__Output, _pb_CatchFishResp>;
  
  CreateFish: grpc.handleUnaryCall<_pb_CreateFishReq__Output, _pb_CreateFishResp>;
  
  UpdateFish: grpc.handleUnaryCall<_pb_UpdateFishReq__Output, _pb_UpdateFishResp>;
  
}

export interface CurrencyDefinition extends grpc.ServiceDefinition {
  CatchFish: MethodDefinition<_pb_CatchFishReq, _pb_CatchFishResp, _pb_CatchFishReq__Output, _pb_CatchFishResp__Output>
  CreateFish: MethodDefinition<_pb_CreateFishReq, _pb_CreateFishResp, _pb_CreateFishReq__Output, _pb_CreateFishResp__Output>
  UpdateFish: MethodDefinition<_pb_UpdateFishReq, _pb_UpdateFishResp, _pb_UpdateFishReq__Output, _pb_UpdateFishResp__Output>
}
