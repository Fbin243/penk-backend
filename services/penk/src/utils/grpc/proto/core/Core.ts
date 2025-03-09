// Original file: ../../proto/core/core_service.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { Character as _core_Character, Character__Output as _core_Character__Output } from '../core/Character';
import type { CharacterInput as _core_CharacterInput, CharacterInput__Output as _core_CharacterInput__Output } from '../core/CharacterInput';
import type { CheckPermissionReq as _core_CheckPermissionReq, CheckPermissionReq__Output as _core_CheckPermissionReq__Output } from '../core/CheckPermissionReq';
import type { CheckPermissionResp as _core_CheckPermissionResp, CheckPermissionResp__Output as _core_CheckPermissionResp__Output } from '../core/CheckPermissionResp';
import type { IntrospectReq as _core_IntrospectReq, IntrospectReq__Output as _core_IntrospectReq__Output } from '../core/IntrospectReq';
import type { IntrospectResp as _core_IntrospectResp, IntrospectResp__Output as _core_IntrospectResp__Output } from '../core/IntrospectResp';

export interface CoreClient extends grpc.Client {
  CheckPermission(argument: _core_CheckPermissionReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  CheckPermission(argument: _core_CheckPermissionReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  CheckPermission(argument: _core_CheckPermissionReq, options: grpc.CallOptions, callback: grpc.requestCallback<_core_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  CheckPermission(argument: _core_CheckPermissionReq, callback: grpc.requestCallback<_core_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  checkPermission(argument: _core_CheckPermissionReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  checkPermission(argument: _core_CheckPermissionReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  checkPermission(argument: _core_CheckPermissionReq, options: grpc.CallOptions, callback: grpc.requestCallback<_core_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  checkPermission(argument: _core_CheckPermissionReq, callback: grpc.requestCallback<_core_CheckPermissionResp__Output>): grpc.ClientUnaryCall;
  
  IntrospectToken(argument: _core_IntrospectReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  IntrospectToken(argument: _core_IntrospectReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  IntrospectToken(argument: _core_IntrospectReq, options: grpc.CallOptions, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  IntrospectToken(argument: _core_IntrospectReq, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectToken(argument: _core_IntrospectReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectToken(argument: _core_IntrospectReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectToken(argument: _core_IntrospectReq, options: grpc.CallOptions, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectToken(argument: _core_IntrospectReq, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  
  UpsertCharacter(argument: _core_CharacterInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  UpsertCharacter(argument: _core_CharacterInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  UpsertCharacter(argument: _core_CharacterInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  UpsertCharacter(argument: _core_CharacterInput, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  upsertCharacter(argument: _core_CharacterInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  upsertCharacter(argument: _core_CharacterInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  upsertCharacter(argument: _core_CharacterInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  upsertCharacter(argument: _core_CharacterInput, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  
}

export interface CoreHandlers extends grpc.UntypedServiceImplementation {
  CheckPermission: grpc.handleUnaryCall<_core_CheckPermissionReq__Output, _core_CheckPermissionResp>;
  
  IntrospectToken: grpc.handleUnaryCall<_core_IntrospectReq__Output, _core_IntrospectResp>;
  
  UpsertCharacter: grpc.handleUnaryCall<_core_CharacterInput__Output, _core_Character>;
  
}

export interface CoreDefinition extends grpc.ServiceDefinition {
  CheckPermission: MethodDefinition<_core_CheckPermissionReq, _core_CheckPermissionResp, _core_CheckPermissionReq__Output, _core_CheckPermissionResp__Output>
  IntrospectToken: MethodDefinition<_core_IntrospectReq, _core_IntrospectResp, _core_IntrospectReq__Output, _core_IntrospectResp__Output>
  UpsertCharacter: MethodDefinition<_core_CharacterInput, _core_Character, _core_CharacterInput__Output, _core_Character__Output>
}
