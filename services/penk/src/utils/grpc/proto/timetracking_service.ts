import type * as grpc from '@grpc/grpc-js';
import type { MessageTypeDefinition } from '@grpc/proto-loader';

import type { TimeTrackingServiceClient as _timetracking_TimeTrackingServiceClient, TimeTrackingServiceDefinition as _timetracking_TimeTrackingServiceDefinition } from './timetracking/TimeTrackingService';

type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  common: {
    EmptyReq: MessageTypeDefinition
    EmptyResp: MessageTypeDefinition
  }
  timetracking: {
    CreateTimeTrackingRequest: MessageTypeDefinition
    TimeTracking: MessageTypeDefinition
    TimeTrackingService: SubtypeConstructor<typeof grpc.Client, _timetracking_TimeTrackingServiceClient> & { service: _timetracking_TimeTrackingServiceDefinition }
    TimeTrackingWithFish: MessageTypeDefinition
    TotalTimeTrackingRequest: MessageTypeDefinition
    TotalTimeTrackingResponse: MessageTypeDefinition
  }
}

