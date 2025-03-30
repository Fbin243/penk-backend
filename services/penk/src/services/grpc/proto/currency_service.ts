import type * as grpc from '@grpc/grpc-js';
import type { MessageTypeDefinition } from '@grpc/proto-loader';

import type { CurrencyClient as _currency_CurrencyClient, CurrencyDefinition as _currency_CurrencyDefinition } from './currency/Currency';

type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  currency: {
    CatchFishReq: MessageTypeDefinition
    CatchFishResp: MessageTypeDefinition
    CreateFishReq: MessageTypeDefinition
    CreateFishResp: MessageTypeDefinition
    Currency: SubtypeConstructor<typeof grpc.Client, _currency_CurrencyClient> & { service: _currency_CurrencyDefinition }
    UpdateFishReq: MessageTypeDefinition
    UpdateFishResp: MessageTypeDefinition
  }
}

