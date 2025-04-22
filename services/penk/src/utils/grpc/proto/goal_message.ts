import type * as grpc from '@grpc/grpc-js';
import type { EnumTypeDefinition, MessageTypeDefinition } from '@grpc/proto-loader';


type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  core: {
    Checkbox: MessageTypeDefinition
    CheckboxInput: MessageTypeDefinition
    Goal: MessageTypeDefinition
    GoalInput: MessageTypeDefinition
    GoalMetric: MessageTypeDefinition
    GoalMetricInput: MessageTypeDefinition
    GoalStatus: EnumTypeDefinition
    MetricCondition: EnumTypeDefinition
    Range: MessageTypeDefinition
    RangeInput: MessageTypeDefinition
  }
}

