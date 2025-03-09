import type * as grpc from '@grpc/grpc-js';
import type { MessageTypeDefinition } from '@grpc/proto-loader';

import type { AnalyticClient as _analytic_AnalyticClient, AnalyticDefinition as _analytic_AnalyticDefinition } from './analytic/Analytic';

type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  analytic: {
    Analytic: SubtypeConstructor<typeof grpc.Client, _analytic_AnalyticClient> & { service: _analytic_AnalyticDefinition }
    DeleteCapturedRecordsReq: MessageTypeDefinition
    DeleteCapturedRecordsResp: MessageTypeDefinition
  }
}

