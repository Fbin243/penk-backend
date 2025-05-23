import type * as grpc from '@grpc/grpc-js';
import type { EnumTypeDefinition, MessageTypeDefinition } from '@grpc/proto-loader';


type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  core: {
    CompletionType: EnumTypeDefinition
    Habit: MessageTypeDefinition
    HabitInput: MessageTypeDefinition
    HabitReset: EnumTypeDefinition
  }
}

