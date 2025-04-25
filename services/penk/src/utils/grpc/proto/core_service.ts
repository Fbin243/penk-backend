import type * as grpc from '@grpc/grpc-js';
import type { EnumTypeDefinition, MessageTypeDefinition } from '@grpc/proto-loader';

import type { CoreClient as _core_CoreClient, CoreDefinition as _core_CoreDefinition } from './core/Core';

type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  common: {
    EmptyReq: MessageTypeDefinition
    EmptyResp: MessageTypeDefinition
    IdReq: MessageTypeDefinition
    IdResp: MessageTypeDefinition
  }
  core: {
    Category: MessageTypeDefinition
    CategoryInput: MessageTypeDefinition
    CategoryStyle: MessageTypeDefinition
    CategoryStyleInput: MessageTypeDefinition
    Character: MessageTypeDefinition
    CharacterInput: MessageTypeDefinition
    Checkbox: MessageTypeDefinition
    CheckboxInput: MessageTypeDefinition
    CompletionType: EnumTypeDefinition
    Core: SubtypeConstructor<typeof grpc.Client, _core_CoreClient> & { service: _core_CoreDefinition }
    EntityType: EnumTypeDefinition
    Goal: MessageTypeDefinition
    GoalInput: MessageTypeDefinition
    GoalMetric: MessageTypeDefinition
    GoalMetricInput: MessageTypeDefinition
    GoalStatus: EnumTypeDefinition
    Habit: MessageTypeDefinition
    HabitInput: MessageTypeDefinition
    HabitReset: EnumTypeDefinition
    IntrospectReq: MessageTypeDefinition
    IntrospectResp: MessageTypeDefinition
    Metric: MessageTypeDefinition
    MetricCondition: EnumTypeDefinition
    MetricInput: MessageTypeDefinition
    Range: MessageTypeDefinition
    RangeInput: MessageTypeDefinition
    TaskInput: MessageTypeDefinition
    TaskInputs: MessageTypeDefinition
    TaskMsg: MessageTypeDefinition
    TaskMsgs: MessageTypeDefinition
    TaskSession: MessageTypeDefinition
    TaskSessionInput: MessageTypeDefinition
    TaskSessionInputs: MessageTypeDefinition
    TaskSessions: MessageTypeDefinition
    TimeTracking: MessageTypeDefinition
    TimeTrackingInput: MessageTypeDefinition
  }
}

