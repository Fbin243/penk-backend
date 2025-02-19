import type * as grpc from '@grpc/grpc-js';
import type { MessageTypeDefinition } from '@grpc/proto-loader';

import type { AnalyticClient as _pb_AnalyticClient, AnalyticDefinition as _pb_AnalyticDefinition } from './pb/Analytic';

type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  pb: {
    Analytic: SubtypeConstructor<typeof grpc.Client, _pb_AnalyticClient> & { service: _pb_AnalyticDefinition }
    DeleteCapturedRecordsReq: MessageTypeDefinition
    DeleteCapturedRecordsResp: MessageTypeDefinition
  }
}

