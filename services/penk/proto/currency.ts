import type * as grpc from '@grpc/grpc-js';
import type { MessageTypeDefinition } from '@grpc/proto-loader';

import type { CurrencyClient as _pb_CurrencyClient, CurrencyDefinition as _pb_CurrencyDefinition } from './pb/Currency';

type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  pb: {
    CatchFishReq: MessageTypeDefinition
    CatchFishResp: MessageTypeDefinition
    CreateFishReq: MessageTypeDefinition
    CreateFishResp: MessageTypeDefinition
    Currency: SubtypeConstructor<typeof grpc.Client, _pb_CurrencyClient> & { service: _pb_CurrencyDefinition }
    UpdateFishReq: MessageTypeDefinition
    UpdateFishResp: MessageTypeDefinition
  }
}

