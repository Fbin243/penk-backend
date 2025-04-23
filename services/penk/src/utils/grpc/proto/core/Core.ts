// Original file: ../../proto/core/core_service.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { Character as _core_Character, Character__Output as _core_Character__Output } from '../core/Character';
import type { CharacterInput as _core_CharacterInput, CharacterInput__Output as _core_CharacterInput__Output } from '../core/CharacterInput';
import type { Goal as _core_Goal, Goal__Output as _core_Goal__Output } from '../core/Goal';
import type { GoalInput as _core_GoalInput, GoalInput__Output as _core_GoalInput__Output } from '../core/GoalInput';
import type { IntrospectReq as _core_IntrospectReq, IntrospectReq__Output as _core_IntrospectReq__Output } from '../core/IntrospectReq';
import type { IntrospectResp as _core_IntrospectResp, IntrospectResp__Output as _core_IntrospectResp__Output } from '../core/IntrospectResp';
import type { TaskInput as _core_TaskInput, TaskInput__Output as _core_TaskInput__Output } from '../core/TaskInput';
import type { TaskMsg as _core_TaskMsg, TaskMsg__Output as _core_TaskMsg__Output } from '../core/TaskMsg';
import type { TaskSession as _core_TaskSession, TaskSession__Output as _core_TaskSession__Output } from '../core/TaskSession';
import type { TaskSessionInput as _core_TaskSessionInput, TaskSessionInput__Output as _core_TaskSessionInput__Output } from '../core/TaskSessionInput';
import type { TimeTracking as _core_TimeTracking, TimeTracking__Output as _core_TimeTracking__Output } from '../core/TimeTracking';
import type { TimeTrackingInput as _core_TimeTrackingInput, TimeTrackingInput__Output as _core_TimeTrackingInput__Output } from '../core/TimeTrackingInput';

export interface CoreClient extends grpc.Client {
  IntrospectUser(argument: _core_IntrospectReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  IntrospectUser(argument: _core_IntrospectReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  IntrospectUser(argument: _core_IntrospectReq, options: grpc.CallOptions, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  IntrospectUser(argument: _core_IntrospectReq, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectUser(argument: _core_IntrospectReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectUser(argument: _core_IntrospectReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectUser(argument: _core_IntrospectReq, options: grpc.CallOptions, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectUser(argument: _core_IntrospectReq, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  
  UpsertCharacter(argument: _core_CharacterInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  UpsertCharacter(argument: _core_CharacterInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  UpsertCharacter(argument: _core_CharacterInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  UpsertCharacter(argument: _core_CharacterInput, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  upsertCharacter(argument: _core_CharacterInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  upsertCharacter(argument: _core_CharacterInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  upsertCharacter(argument: _core_CharacterInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  upsertCharacter(argument: _core_CharacterInput, callback: grpc.requestCallback<_core_Character__Output>): grpc.ClientUnaryCall;
  
  UpsertGoal(argument: _core_GoalInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Goal__Output>): grpc.ClientUnaryCall;
  UpsertGoal(argument: _core_GoalInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_Goal__Output>): grpc.ClientUnaryCall;
  UpsertGoal(argument: _core_GoalInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Goal__Output>): grpc.ClientUnaryCall;
  UpsertGoal(argument: _core_GoalInput, callback: grpc.requestCallback<_core_Goal__Output>): grpc.ClientUnaryCall;
  upsertGoal(argument: _core_GoalInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Goal__Output>): grpc.ClientUnaryCall;
  upsertGoal(argument: _core_GoalInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_Goal__Output>): grpc.ClientUnaryCall;
  upsertGoal(argument: _core_GoalInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Goal__Output>): grpc.ClientUnaryCall;
  upsertGoal(argument: _core_GoalInput, callback: grpc.requestCallback<_core_Goal__Output>): grpc.ClientUnaryCall;
  
  UpsertTask(argument: _core_TaskInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskMsg__Output>): grpc.ClientUnaryCall;
  UpsertTask(argument: _core_TaskInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_TaskMsg__Output>): grpc.ClientUnaryCall;
  UpsertTask(argument: _core_TaskInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskMsg__Output>): grpc.ClientUnaryCall;
  UpsertTask(argument: _core_TaskInput, callback: grpc.requestCallback<_core_TaskMsg__Output>): grpc.ClientUnaryCall;
  upsertTask(argument: _core_TaskInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskMsg__Output>): grpc.ClientUnaryCall;
  upsertTask(argument: _core_TaskInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_TaskMsg__Output>): grpc.ClientUnaryCall;
  upsertTask(argument: _core_TaskInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskMsg__Output>): grpc.ClientUnaryCall;
  upsertTask(argument: _core_TaskInput, callback: grpc.requestCallback<_core_TaskMsg__Output>): grpc.ClientUnaryCall;
  
  UpsertTaskSession(argument: _core_TaskSessionInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskSession__Output>): grpc.ClientUnaryCall;
  UpsertTaskSession(argument: _core_TaskSessionInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_TaskSession__Output>): grpc.ClientUnaryCall;
  UpsertTaskSession(argument: _core_TaskSessionInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskSession__Output>): grpc.ClientUnaryCall;
  UpsertTaskSession(argument: _core_TaskSessionInput, callback: grpc.requestCallback<_core_TaskSession__Output>): grpc.ClientUnaryCall;
  upsertTaskSession(argument: _core_TaskSessionInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskSession__Output>): grpc.ClientUnaryCall;
  upsertTaskSession(argument: _core_TaskSessionInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_TaskSession__Output>): grpc.ClientUnaryCall;
  upsertTaskSession(argument: _core_TaskSessionInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskSession__Output>): grpc.ClientUnaryCall;
  upsertTaskSession(argument: _core_TaskSessionInput, callback: grpc.requestCallback<_core_TaskSession__Output>): grpc.ClientUnaryCall;
  
  UpsertTimeTracking(argument: _core_TimeTrackingInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TimeTracking__Output>): grpc.ClientUnaryCall;
  UpsertTimeTracking(argument: _core_TimeTrackingInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_TimeTracking__Output>): grpc.ClientUnaryCall;
  UpsertTimeTracking(argument: _core_TimeTrackingInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TimeTracking__Output>): grpc.ClientUnaryCall;
  UpsertTimeTracking(argument: _core_TimeTrackingInput, callback: grpc.requestCallback<_core_TimeTracking__Output>): grpc.ClientUnaryCall;
  upsertTimeTracking(argument: _core_TimeTrackingInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TimeTracking__Output>): grpc.ClientUnaryCall;
  upsertTimeTracking(argument: _core_TimeTrackingInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_TimeTracking__Output>): grpc.ClientUnaryCall;
  upsertTimeTracking(argument: _core_TimeTrackingInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TimeTracking__Output>): grpc.ClientUnaryCall;
  upsertTimeTracking(argument: _core_TimeTrackingInput, callback: grpc.requestCallback<_core_TimeTracking__Output>): grpc.ClientUnaryCall;
  
}

export interface CoreHandlers extends grpc.UntypedServiceImplementation {
  IntrospectUser: grpc.handleUnaryCall<_core_IntrospectReq__Output, _core_IntrospectResp>;
  
  UpsertCharacter: grpc.handleUnaryCall<_core_CharacterInput__Output, _core_Character>;
  
  UpsertGoal: grpc.handleUnaryCall<_core_GoalInput__Output, _core_Goal>;
  
  UpsertTask: grpc.handleUnaryCall<_core_TaskInput__Output, _core_TaskMsg>;
  
  UpsertTaskSession: grpc.handleUnaryCall<_core_TaskSessionInput__Output, _core_TaskSession>;
  
  UpsertTimeTracking: grpc.handleUnaryCall<_core_TimeTrackingInput__Output, _core_TimeTracking>;
  
}

export interface CoreDefinition extends grpc.ServiceDefinition {
  IntrospectUser: MethodDefinition<_core_IntrospectReq, _core_IntrospectResp, _core_IntrospectReq__Output, _core_IntrospectResp__Output>
  UpsertCharacter: MethodDefinition<_core_CharacterInput, _core_Character, _core_CharacterInput__Output, _core_Character__Output>
  UpsertGoal: MethodDefinition<_core_GoalInput, _core_Goal, _core_GoalInput__Output, _core_Goal__Output>
  UpsertTask: MethodDefinition<_core_TaskInput, _core_TaskMsg, _core_TaskInput__Output, _core_TaskMsg__Output>
  UpsertTaskSession: MethodDefinition<_core_TaskSessionInput, _core_TaskSession, _core_TaskSessionInput__Output, _core_TaskSession__Output>
  UpsertTimeTracking: MethodDefinition<_core_TimeTrackingInput, _core_TimeTracking, _core_TimeTrackingInput__Output, _core_TimeTracking__Output>
}
