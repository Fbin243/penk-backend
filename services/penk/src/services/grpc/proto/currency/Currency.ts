// Original file: ../../proto/currency/currency_service.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { CatchFishReq as _currency_CatchFishReq, CatchFishReq__Output as _currency_CatchFishReq__Output } from '../currency/CatchFishReq';
import type { CatchFishResp as _currency_CatchFishResp, CatchFishResp__Output as _currency_CatchFishResp__Output } from '../currency/CatchFishResp';
import type { CreateFishReq as _currency_CreateFishReq, CreateFishReq__Output as _currency_CreateFishReq__Output } from '../currency/CreateFishReq';
import type { CreateFishResp as _currency_CreateFishResp, CreateFishResp__Output as _currency_CreateFishResp__Output } from '../currency/CreateFishResp';
import type { UpdateFishReq as _currency_UpdateFishReq, UpdateFishReq__Output as _currency_UpdateFishReq__Output } from '../currency/UpdateFishReq';
import type { UpdateFishResp as _currency_UpdateFishResp, UpdateFishResp__Output as _currency_UpdateFishResp__Output } from '../currency/UpdateFishResp';

export interface CurrencyClient extends grpc.Client {
  CatchFish(argument: _currency_CatchFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_CatchFishResp__Output>): grpc.ClientUnaryCall;
  CatchFish(argument: _currency_CatchFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_currency_CatchFishResp__Output>): grpc.ClientUnaryCall;
  CatchFish(argument: _currency_CatchFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_CatchFishResp__Output>): grpc.ClientUnaryCall;
  CatchFish(argument: _currency_CatchFishReq, callback: grpc.requestCallback<_currency_CatchFishResp__Output>): grpc.ClientUnaryCall;
  catchFish(argument: _currency_CatchFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_CatchFishResp__Output>): grpc.ClientUnaryCall;
  catchFish(argument: _currency_CatchFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_currency_CatchFishResp__Output>): grpc.ClientUnaryCall;
  catchFish(argument: _currency_CatchFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_CatchFishResp__Output>): grpc.ClientUnaryCall;
  catchFish(argument: _currency_CatchFishReq, callback: grpc.requestCallback<_currency_CatchFishResp__Output>): grpc.ClientUnaryCall;
  
  CreateFish(argument: _currency_CreateFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  CreateFish(argument: _currency_CreateFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  CreateFish(argument: _currency_CreateFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  CreateFish(argument: _currency_CreateFishReq, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  createFish(argument: _currency_CreateFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  createFish(argument: _currency_CreateFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  createFish(argument: _currency_CreateFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  createFish(argument: _currency_CreateFishReq, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  
  UpdateFish(argument: _currency_UpdateFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  UpdateFish(argument: _currency_UpdateFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_currency_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  UpdateFish(argument: _currency_UpdateFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  UpdateFish(argument: _currency_UpdateFishReq, callback: grpc.requestCallback<_currency_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  updateFish(argument: _currency_UpdateFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  updateFish(argument: _currency_UpdateFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_currency_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  updateFish(argument: _currency_UpdateFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  updateFish(argument: _currency_UpdateFishReq, callback: grpc.requestCallback<_currency_UpdateFishResp__Output>): grpc.ClientUnaryCall;
  
}

export interface CurrencyHandlers extends grpc.UntypedServiceImplementation {
  CatchFish: grpc.handleUnaryCall<_currency_CatchFishReq__Output, _currency_CatchFishResp>;
  
  CreateFish: grpc.handleUnaryCall<_currency_CreateFishReq__Output, _currency_CreateFishResp>;
  
  UpdateFish: grpc.handleUnaryCall<_currency_UpdateFishReq__Output, _currency_UpdateFishResp>;
  
}

export interface CurrencyDefinition extends grpc.ServiceDefinition {
  CatchFish: MethodDefinition<_currency_CatchFishReq, _currency_CatchFishResp, _currency_CatchFishReq__Output, _currency_CatchFishResp__Output>
  CreateFish: MethodDefinition<_currency_CreateFishReq, _currency_CreateFishResp, _currency_CreateFishReq__Output, _currency_CreateFishResp__Output>
  UpdateFish: MethodDefinition<_currency_UpdateFishReq, _currency_UpdateFishResp, _currency_UpdateFishReq__Output, _currency_UpdateFishResp__Output>
}
