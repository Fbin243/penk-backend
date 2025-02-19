import type * as grpc from '@grpc/grpc-js';
import type { MessageTypeDefinition } from '@grpc/proto-loader';

import type { CoreClient as _pb_CoreClient, CoreDefinition as _pb_CoreDefinition } from './pb/Core';

type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  pb: {
    CheckPermissionReq: MessageTypeDefinition
    CheckPermissionResp: MessageTypeDefinition
    Core: SubtypeConstructor<typeof grpc.Client, _pb_CoreClient> & { service: _pb_CoreDefinition }
    IntrospectReq: MessageTypeDefinition
    IntrospectResp: MessageTypeDefinition
  }
}

