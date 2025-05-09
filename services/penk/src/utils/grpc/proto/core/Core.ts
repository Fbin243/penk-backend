// Original file: ../../proto/core/core_service.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { Category as _core_Category, Category__Output as _core_Category__Output } from '../core/Category';
import type { CategoryInput as _core_CategoryInput, CategoryInput__Output as _core_CategoryInput__Output } from '../core/CategoryInput';
import type { Character as _core_Character, Character__Output as _core_Character__Output } from '../core/Character';
import type { CharacterInput as _core_CharacterInput, CharacterInput__Output as _core_CharacterInput__Output } from '../core/CharacterInput';
import type { Goal as _core_Goal, Goal__Output as _core_Goal__Output } from '../core/Goal';
import type { GoalInput as _core_GoalInput, GoalInput__Output as _core_GoalInput__Output } from '../core/GoalInput';
import type { Habit as _core_Habit, Habit__Output as _core_Habit__Output } from '../core/Habit';
import type { HabitInput as _core_HabitInput, HabitInput__Output as _core_HabitInput__Output } from '../core/HabitInput';
import type { HabitLog as _core_HabitLog, HabitLog__Output as _core_HabitLog__Output } from '../core/HabitLog';
import type { HabitLogInput as _core_HabitLogInput, HabitLogInput__Output as _core_HabitLogInput__Output } from '../core/HabitLogInput';
import type { IdReq as _common_IdReq, IdReq__Output as _common_IdReq__Output } from '../common/IdReq';
import type { IdResp as _common_IdResp, IdResp__Output as _common_IdResp__Output } from '../common/IdResp';
import type { IntrospectReq as _core_IntrospectReq, IntrospectReq__Output as _core_IntrospectReq__Output } from '../core/IntrospectReq';
import type { IntrospectResp as _core_IntrospectResp, IntrospectResp__Output as _core_IntrospectResp__Output } from '../core/IntrospectResp';
import type { Metric as _core_Metric, Metric__Output as _core_Metric__Output } from '../core/Metric';
import type { MetricInput as _core_MetricInput, MetricInput__Output as _core_MetricInput__Output } from '../core/MetricInput';
import type { TaskInput as _core_TaskInput, TaskInput__Output as _core_TaskInput__Output } from '../core/TaskInput';
import type { TaskInputs as _core_TaskInputs, TaskInputs__Output as _core_TaskInputs__Output } from '../core/TaskInputs';
import type { TaskMsg as _core_TaskMsg, TaskMsg__Output as _core_TaskMsg__Output } from '../core/TaskMsg';
import type { TaskMsgs as _core_TaskMsgs, TaskMsgs__Output as _core_TaskMsgs__Output } from '../core/TaskMsgs';
import type { TaskSession as _core_TaskSession, TaskSession__Output as _core_TaskSession__Output } from '../core/TaskSession';
import type { TaskSessionInput as _core_TaskSessionInput, TaskSessionInput__Output as _core_TaskSessionInput__Output } from '../core/TaskSessionInput';
import type { TaskSessionInputs as _core_TaskSessionInputs, TaskSessionInputs__Output as _core_TaskSessionInputs__Output } from '../core/TaskSessionInputs';
import type { TaskSessions as _core_TaskSessions, TaskSessions__Output as _core_TaskSessions__Output } from '../core/TaskSessions';
import type { TimeTracking as _core_TimeTracking, TimeTracking__Output as _core_TimeTracking__Output } from '../core/TimeTracking';
import type { TimeTrackingInput as _core_TimeTrackingInput, TimeTrackingInput__Output as _core_TimeTrackingInput__Output } from '../core/TimeTrackingInput';

export interface CoreClient extends grpc.Client {
  DeleteCategory(argument: _common_IdReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteCategory(argument: _common_IdReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteCategory(argument: _common_IdReq, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteCategory(argument: _common_IdReq, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteCategory(argument: _common_IdReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteCategory(argument: _common_IdReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteCategory(argument: _common_IdReq, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteCategory(argument: _common_IdReq, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  
  DeleteGoal(argument: _common_IdReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteGoal(argument: _common_IdReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteGoal(argument: _common_IdReq, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteGoal(argument: _common_IdReq, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteGoal(argument: _common_IdReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteGoal(argument: _common_IdReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteGoal(argument: _common_IdReq, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteGoal(argument: _common_IdReq, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  
  DeleteHabit(argument: _common_IdReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteHabit(argument: _common_IdReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteHabit(argument: _common_IdReq, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteHabit(argument: _common_IdReq, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteHabit(argument: _common_IdReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteHabit(argument: _common_IdReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteHabit(argument: _common_IdReq, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteHabit(argument: _common_IdReq, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  
  DeleteMetric(argument: _common_IdReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteMetric(argument: _common_IdReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteMetric(argument: _common_IdReq, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteMetric(argument: _common_IdReq, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteMetric(argument: _common_IdReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteMetric(argument: _common_IdReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteMetric(argument: _common_IdReq, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteMetric(argument: _common_IdReq, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  
  DeleteTask(argument: _common_IdReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteTask(argument: _common_IdReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteTask(argument: _common_IdReq, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteTask(argument: _common_IdReq, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteTask(argument: _common_IdReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteTask(argument: _common_IdReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteTask(argument: _common_IdReq, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteTask(argument: _common_IdReq, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  
  DeleteTaskSession(argument: _common_IdReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteTaskSession(argument: _common_IdReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteTaskSession(argument: _common_IdReq, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  DeleteTaskSession(argument: _common_IdReq, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteTaskSession(argument: _common_IdReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteTaskSession(argument: _common_IdReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteTaskSession(argument: _common_IdReq, options: grpc.CallOptions, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  deleteTaskSession(argument: _common_IdReq, callback: grpc.requestCallback<_common_IdResp__Output>): grpc.ClientUnaryCall;
  
  IntrospectUser(argument: _core_IntrospectReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  IntrospectUser(argument: _core_IntrospectReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  IntrospectUser(argument: _core_IntrospectReq, options: grpc.CallOptions, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  IntrospectUser(argument: _core_IntrospectReq, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectUser(argument: _core_IntrospectReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectUser(argument: _core_IntrospectReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectUser(argument: _core_IntrospectReq, options: grpc.CallOptions, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  introspectUser(argument: _core_IntrospectReq, callback: grpc.requestCallback<_core_IntrospectResp__Output>): grpc.ClientUnaryCall;
  
  UpsertCategory(argument: _core_CategoryInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Category__Output>): grpc.ClientUnaryCall;
  UpsertCategory(argument: _core_CategoryInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_Category__Output>): grpc.ClientUnaryCall;
  UpsertCategory(argument: _core_CategoryInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Category__Output>): grpc.ClientUnaryCall;
  UpsertCategory(argument: _core_CategoryInput, callback: grpc.requestCallback<_core_Category__Output>): grpc.ClientUnaryCall;
  upsertCategory(argument: _core_CategoryInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Category__Output>): grpc.ClientUnaryCall;
  upsertCategory(argument: _core_CategoryInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_Category__Output>): grpc.ClientUnaryCall;
  upsertCategory(argument: _core_CategoryInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Category__Output>): grpc.ClientUnaryCall;
  upsertCategory(argument: _core_CategoryInput, callback: grpc.requestCallback<_core_Category__Output>): grpc.ClientUnaryCall;
  
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
  
  UpsertHabit(argument: _core_HabitInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Habit__Output>): grpc.ClientUnaryCall;
  UpsertHabit(argument: _core_HabitInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_Habit__Output>): grpc.ClientUnaryCall;
  UpsertHabit(argument: _core_HabitInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Habit__Output>): grpc.ClientUnaryCall;
  UpsertHabit(argument: _core_HabitInput, callback: grpc.requestCallback<_core_Habit__Output>): grpc.ClientUnaryCall;
  upsertHabit(argument: _core_HabitInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Habit__Output>): grpc.ClientUnaryCall;
  upsertHabit(argument: _core_HabitInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_Habit__Output>): grpc.ClientUnaryCall;
  upsertHabit(argument: _core_HabitInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Habit__Output>): grpc.ClientUnaryCall;
  upsertHabit(argument: _core_HabitInput, callback: grpc.requestCallback<_core_Habit__Output>): grpc.ClientUnaryCall;
  
  UpsertHabitLog(argument: _core_HabitLogInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_HabitLog__Output>): grpc.ClientUnaryCall;
  UpsertHabitLog(argument: _core_HabitLogInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_HabitLog__Output>): grpc.ClientUnaryCall;
  UpsertHabitLog(argument: _core_HabitLogInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_HabitLog__Output>): grpc.ClientUnaryCall;
  UpsertHabitLog(argument: _core_HabitLogInput, callback: grpc.requestCallback<_core_HabitLog__Output>): grpc.ClientUnaryCall;
  upsertHabitLog(argument: _core_HabitLogInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_HabitLog__Output>): grpc.ClientUnaryCall;
  upsertHabitLog(argument: _core_HabitLogInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_HabitLog__Output>): grpc.ClientUnaryCall;
  upsertHabitLog(argument: _core_HabitLogInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_HabitLog__Output>): grpc.ClientUnaryCall;
  upsertHabitLog(argument: _core_HabitLogInput, callback: grpc.requestCallback<_core_HabitLog__Output>): grpc.ClientUnaryCall;
  
  UpsertMetric(argument: _core_MetricInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Metric__Output>): grpc.ClientUnaryCall;
  UpsertMetric(argument: _core_MetricInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_Metric__Output>): grpc.ClientUnaryCall;
  UpsertMetric(argument: _core_MetricInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Metric__Output>): grpc.ClientUnaryCall;
  UpsertMetric(argument: _core_MetricInput, callback: grpc.requestCallback<_core_Metric__Output>): grpc.ClientUnaryCall;
  upsertMetric(argument: _core_MetricInput, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Metric__Output>): grpc.ClientUnaryCall;
  upsertMetric(argument: _core_MetricInput, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_Metric__Output>): grpc.ClientUnaryCall;
  upsertMetric(argument: _core_MetricInput, options: grpc.CallOptions, callback: grpc.requestCallback<_core_Metric__Output>): grpc.ClientUnaryCall;
  upsertMetric(argument: _core_MetricInput, callback: grpc.requestCallback<_core_Metric__Output>): grpc.ClientUnaryCall;
  
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
  
  UpsertTaskSessions(argument: _core_TaskSessionInputs, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskSessions__Output>): grpc.ClientUnaryCall;
  UpsertTaskSessions(argument: _core_TaskSessionInputs, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_TaskSessions__Output>): grpc.ClientUnaryCall;
  UpsertTaskSessions(argument: _core_TaskSessionInputs, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskSessions__Output>): grpc.ClientUnaryCall;
  UpsertTaskSessions(argument: _core_TaskSessionInputs, callback: grpc.requestCallback<_core_TaskSessions__Output>): grpc.ClientUnaryCall;
  upsertTaskSessions(argument: _core_TaskSessionInputs, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskSessions__Output>): grpc.ClientUnaryCall;
  upsertTaskSessions(argument: _core_TaskSessionInputs, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_TaskSessions__Output>): grpc.ClientUnaryCall;
  upsertTaskSessions(argument: _core_TaskSessionInputs, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskSessions__Output>): grpc.ClientUnaryCall;
  upsertTaskSessions(argument: _core_TaskSessionInputs, callback: grpc.requestCallback<_core_TaskSessions__Output>): grpc.ClientUnaryCall;
  
  UpsertTasks(argument: _core_TaskInputs, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskMsgs__Output>): grpc.ClientUnaryCall;
  UpsertTasks(argument: _core_TaskInputs, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_TaskMsgs__Output>): grpc.ClientUnaryCall;
  UpsertTasks(argument: _core_TaskInputs, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskMsgs__Output>): grpc.ClientUnaryCall;
  UpsertTasks(argument: _core_TaskInputs, callback: grpc.requestCallback<_core_TaskMsgs__Output>): grpc.ClientUnaryCall;
  upsertTasks(argument: _core_TaskInputs, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskMsgs__Output>): grpc.ClientUnaryCall;
  upsertTasks(argument: _core_TaskInputs, metadata: grpc.Metadata, callback: grpc.requestCallback<_core_TaskMsgs__Output>): grpc.ClientUnaryCall;
  upsertTasks(argument: _core_TaskInputs, options: grpc.CallOptions, callback: grpc.requestCallback<_core_TaskMsgs__Output>): grpc.ClientUnaryCall;
  upsertTasks(argument: _core_TaskInputs, callback: grpc.requestCallback<_core_TaskMsgs__Output>): grpc.ClientUnaryCall;
  
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
  DeleteCategory: grpc.handleUnaryCall<_common_IdReq__Output, _common_IdResp>;
  
  DeleteGoal: grpc.handleUnaryCall<_common_IdReq__Output, _common_IdResp>;
  
  DeleteHabit: grpc.handleUnaryCall<_common_IdReq__Output, _common_IdResp>;
  
  DeleteMetric: grpc.handleUnaryCall<_common_IdReq__Output, _common_IdResp>;
  
  DeleteTask: grpc.handleUnaryCall<_common_IdReq__Output, _common_IdResp>;
  
  DeleteTaskSession: grpc.handleUnaryCall<_common_IdReq__Output, _common_IdResp>;
  
  IntrospectUser: grpc.handleUnaryCall<_core_IntrospectReq__Output, _core_IntrospectResp>;
  
  UpsertCategory: grpc.handleUnaryCall<_core_CategoryInput__Output, _core_Category>;
  
  UpsertCharacter: grpc.handleUnaryCall<_core_CharacterInput__Output, _core_Character>;
  
  UpsertGoal: grpc.handleUnaryCall<_core_GoalInput__Output, _core_Goal>;
  
  UpsertHabit: grpc.handleUnaryCall<_core_HabitInput__Output, _core_Habit>;
  
  UpsertHabitLog: grpc.handleUnaryCall<_core_HabitLogInput__Output, _core_HabitLog>;
  
  UpsertMetric: grpc.handleUnaryCall<_core_MetricInput__Output, _core_Metric>;
  
  UpsertTask: grpc.handleUnaryCall<_core_TaskInput__Output, _core_TaskMsg>;
  
  UpsertTaskSession: grpc.handleUnaryCall<_core_TaskSessionInput__Output, _core_TaskSession>;
  
  UpsertTaskSessions: grpc.handleUnaryCall<_core_TaskSessionInputs__Output, _core_TaskSessions>;
  
  UpsertTasks: grpc.handleUnaryCall<_core_TaskInputs__Output, _core_TaskMsgs>;
  
  UpsertTimeTracking: grpc.handleUnaryCall<_core_TimeTrackingInput__Output, _core_TimeTracking>;
  
}

export interface CoreDefinition extends grpc.ServiceDefinition {
  DeleteCategory: MethodDefinition<_common_IdReq, _common_IdResp, _common_IdReq__Output, _common_IdResp__Output>
  DeleteGoal: MethodDefinition<_common_IdReq, _common_IdResp, _common_IdReq__Output, _common_IdResp__Output>
  DeleteHabit: MethodDefinition<_common_IdReq, _common_IdResp, _common_IdReq__Output, _common_IdResp__Output>
  DeleteMetric: MethodDefinition<_common_IdReq, _common_IdResp, _common_IdReq__Output, _common_IdResp__Output>
  DeleteTask: MethodDefinition<_common_IdReq, _common_IdResp, _common_IdReq__Output, _common_IdResp__Output>
  DeleteTaskSession: MethodDefinition<_common_IdReq, _common_IdResp, _common_IdReq__Output, _common_IdResp__Output>
  IntrospectUser: MethodDefinition<_core_IntrospectReq, _core_IntrospectResp, _core_IntrospectReq__Output, _core_IntrospectResp__Output>
  UpsertCategory: MethodDefinition<_core_CategoryInput, _core_Category, _core_CategoryInput__Output, _core_Category__Output>
  UpsertCharacter: MethodDefinition<_core_CharacterInput, _core_Character, _core_CharacterInput__Output, _core_Character__Output>
  UpsertGoal: MethodDefinition<_core_GoalInput, _core_Goal, _core_GoalInput__Output, _core_Goal__Output>
  UpsertHabit: MethodDefinition<_core_HabitInput, _core_Habit, _core_HabitInput__Output, _core_Habit__Output>
  UpsertHabitLog: MethodDefinition<_core_HabitLogInput, _core_HabitLog, _core_HabitLogInput__Output, _core_HabitLog__Output>
  UpsertMetric: MethodDefinition<_core_MetricInput, _core_Metric, _core_MetricInput__Output, _core_Metric__Output>
  UpsertTask: MethodDefinition<_core_TaskInput, _core_TaskMsg, _core_TaskInput__Output, _core_TaskMsg__Output>
  UpsertTaskSession: MethodDefinition<_core_TaskSessionInput, _core_TaskSession, _core_TaskSessionInput__Output, _core_TaskSession__Output>
  UpsertTaskSessions: MethodDefinition<_core_TaskSessionInputs, _core_TaskSessions, _core_TaskSessionInputs__Output, _core_TaskSessions__Output>
  UpsertTasks: MethodDefinition<_core_TaskInputs, _core_TaskMsgs, _core_TaskInputs__Output, _core_TaskMsgs__Output>
  UpsertTimeTracking: MethodDefinition<_core_TimeTrackingInput, _core_TimeTracking, _core_TimeTrackingInput__Output, _core_TimeTracking__Output>
}
