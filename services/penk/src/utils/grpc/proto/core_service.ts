import type * as grpc from '@grpc/grpc-js';
import type { EnumTypeDefinition, MessageTypeDefinition } from '@grpc/proto-loader';

import type { CoreClient as _core_CoreClient, CoreDefinition as _core_CoreDefinition } from './core/Core';

type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  core: {
    Character: MessageTypeDefinition
    CharacterInput: MessageTypeDefinition
    Checkbox: MessageTypeDefinition
    CheckboxInput: MessageTypeDefinition
    Core: SubtypeConstructor<typeof grpc.Client, _core_CoreClient> & { service: _core_CoreDefinition }
    EntityType: EnumTypeDefinition
    Goal: MessageTypeDefinition
    GoalInput: MessageTypeDefinition
    GoalMetric: MessageTypeDefinition
    GoalMetricInput: MessageTypeDefinition
    GoalStatus: EnumTypeDefinition
    IntrospectReq: MessageTypeDefinition
    IntrospectResp: MessageTypeDefinition
    MetricCondition: EnumTypeDefinition
    Range: MessageTypeDefinition
    RangeInput: MessageTypeDefinition
    TaskInput: MessageTypeDefinition
    TaskMsg: MessageTypeDefinition
    TaskSession: MessageTypeDefinition
    TaskSessionInput: MessageTypeDefinition
    TimeTracking: MessageTypeDefinition
    TimeTrackingInput: MessageTypeDefinition
  }
}

