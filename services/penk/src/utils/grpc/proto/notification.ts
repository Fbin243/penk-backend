import type * as grpc from '@grpc/grpc-js';
import type { MessageTypeDefinition } from '@grpc/proto-loader';

import type { NotificationClient as _notification_NotificationClient, NotificationDefinition as _notification_NotificationDefinition } from './notification/Notification';

type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  notification: {
    Notification: SubtypeConstructor<typeof grpc.Client, _notification_NotificationClient> & { service: _notification_NotificationDefinition }
    SendPushNotificationReq: MessageTypeDefinition
    SendPushNotificationResp: MessageTypeDefinition
  }
}

