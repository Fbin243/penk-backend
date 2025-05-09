// Original file: ../../proto/notification/notification.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { SendPushNotificationReq as _notification_SendPushNotificationReq, SendPushNotificationReq__Output as _notification_SendPushNotificationReq__Output } from '../notification/SendPushNotificationReq';
import type { SendPushNotificationResp as _notification_SendPushNotificationResp, SendPushNotificationResp__Output as _notification_SendPushNotificationResp__Output } from '../notification/SendPushNotificationResp';

export interface NotificationClient extends grpc.Client {
  SendPushNotification(argument: _notification_SendPushNotificationReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_notification_SendPushNotificationResp__Output>): grpc.ClientUnaryCall;
  SendPushNotification(argument: _notification_SendPushNotificationReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_notification_SendPushNotificationResp__Output>): grpc.ClientUnaryCall;
  SendPushNotification(argument: _notification_SendPushNotificationReq, options: grpc.CallOptions, callback: grpc.requestCallback<_notification_SendPushNotificationResp__Output>): grpc.ClientUnaryCall;
  SendPushNotification(argument: _notification_SendPushNotificationReq, callback: grpc.requestCallback<_notification_SendPushNotificationResp__Output>): grpc.ClientUnaryCall;
  sendPushNotification(argument: _notification_SendPushNotificationReq, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_notification_SendPushNotificationResp__Output>): grpc.ClientUnaryCall;
  sendPushNotification(argument: _notification_SendPushNotificationReq, metadata: grpc.Metadata, callback: grpc.requestCallback<_notification_SendPushNotificationResp__Output>): grpc.ClientUnaryCall;
  sendPushNotification(argument: _notification_SendPushNotificationReq, options: grpc.CallOptions, callback: grpc.requestCallback<_notification_SendPushNotificationResp__Output>): grpc.ClientUnaryCall;
  sendPushNotification(argument: _notification_SendPushNotificationReq, callback: grpc.requestCallback<_notification_SendPushNotificationResp__Output>): grpc.ClientUnaryCall;
  
}

export interface NotificationHandlers extends grpc.UntypedServiceImplementation {
  SendPushNotification: grpc.handleUnaryCall<_notification_SendPushNotificationReq__Output, _notification_SendPushNotificationResp>;
  
}

export interface NotificationDefinition extends grpc.ServiceDefinition {
  SendPushNotification: MethodDefinition<_notification_SendPushNotificationReq, _notification_SendPushNotificationResp, _notification_SendPushNotificationReq__Output, _notification_SendPushNotificationResp__Output>
}
