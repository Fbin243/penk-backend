// Original file: ../../proto/currency/currency_service.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { CreateFishReq as _currency_CreateFishReq, CreateFishReq__Output as _currency_CreateFishReq__Output } from '../currency/CreateFishReq';
import type { CreateFishResp as _currency_CreateFishResp, CreateFishResp__Output as _currency_CreateFishResp__Output } from '../currency/CreateFishResp';
import type { DeleteFishReq as _currency_DeleteFishReq, DeleteFishReq__Output as _currency_DeleteFishReq__Output } from '../currency/DeleteFishReq';
import type { DeleteFishResp as _currency_DeleteFishResp, DeleteFishResp__Output as _currency_DeleteFishResp__Output } from '../currency/DeleteFishResp';

export interface CurrencyClient extends grpc.Client {
  CreateFish(argument: _currency_CreateFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  CreateFish(argument: _currency_CreateFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  CreateFish(argument: _currency_CreateFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  CreateFish(argument: _currency_CreateFishReq, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  createFish(argument: _currency_CreateFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  createFish(argument: _currency_CreateFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  createFish(argument: _currency_CreateFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  createFish(argument: _currency_CreateFishReq, callback: grpc.requestCallback<_currency_CreateFishResp__Output>): grpc.ClientUnaryCall;
  
  DeleteFish(argument: _currency_DeleteFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_DeleteFishResp__Output>): grpc.ClientUnaryCall;
  DeleteFish(argument: _currency_DeleteFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_currency_DeleteFishResp__Output>): grpc.ClientUnaryCall;
  DeleteFish(argument: _currency_DeleteFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_DeleteFishResp__Output>): grpc.ClientUnaryCall;
  DeleteFish(argument: _currency_DeleteFishReq, callback: grpc.requestCallback<_currency_DeleteFishResp__Output>): grpc.ClientUnaryCall;
  deleteFish(argument: _currency_DeleteFishReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_DeleteFishResp__Output>): grpc.ClientUnaryCall;
  deleteFish(argument: _currency_DeleteFishReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_currency_DeleteFishResp__Output>): grpc.ClientUnaryCall;
  deleteFish(argument: _currency_DeleteFishReq, options: grpc.CallOptions, callback: grpc.requestCallback<_currency_DeleteFishResp__Output>): grpc.ClientUnaryCall;
  deleteFish(argument: _currency_DeleteFishReq, callback: grpc.requestCallback<_currency_DeleteFishResp__Output>): grpc.ClientUnaryCall;
  
}

export interface CurrencyHandlers extends grpc.UntypedServiceImplementation {
  CreateFish: grpc.handleUnaryCall<_currency_CreateFishReq__Output, _currency_CreateFishResp>;
  
  DeleteFish: grpc.handleUnaryCall<_currency_DeleteFishReq__Output, _currency_DeleteFishResp>;
  
}

export interface CurrencyDefinition extends grpc.ServiceDefinition {
  CreateFish: MethodDefinition<_currency_CreateFishReq, _currency_CreateFishResp, _currency_CreateFishReq__Output, _currency_CreateFishResp__Output>
  DeleteFish: MethodDefinition<_currency_DeleteFishReq, _currency_DeleteFishResp, _currency_DeleteFishReq__Output, _currency_DeleteFishResp__Output>
}
