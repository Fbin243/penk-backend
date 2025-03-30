import type * as grpc from '@grpc/grpc-js';
import type { MessageTypeDefinition } from '@grpc/proto-loader';

import type { CoreClient as _core_CoreClient, CoreDefinition as _core_CoreDefinition } from './core/Core';

type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  core: {
    Category: MessageTypeDefinition
    CategoryInput: MessageTypeDefinition
    CategoryStyle: MessageTypeDefinition
    CategoryStyleInput: MessageTypeDefinition
    Character: MessageTypeDefinition
    CharacterInput: MessageTypeDefinition
    CheckPermissionReq: MessageTypeDefinition
    CheckPermissionResp: MessageTypeDefinition
    Core: SubtypeConstructor<typeof grpc.Client, _core_CoreClient> & { service: _core_CoreDefinition }
    IntrospectReq: MessageTypeDefinition
    IntrospectResp: MessageTypeDefinition
    Metric: MessageTypeDefinition
    MetricInput: MessageTypeDefinition
  }
}

