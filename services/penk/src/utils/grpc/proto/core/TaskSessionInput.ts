// Original file: ../../proto/core/task_session_message.proto

import type { Long } from '@grpc/proto-loader';

export interface TaskSessionInput {
  'id'?: (string);
  'taskId'?: (string);
  'startTime'?: (number | string | Long);
  'endTime'?: (number | string | Long);
  'completedTime'?: (number | string | Long);
  '_id'?: "id";
  '_completedTime'?: "completedTime";
}

export interface TaskSessionInput__Output {
  'id'?: (string);
  'taskId': (string);
  'startTime': (string);
  'endTime': (string);
  'completedTime'?: (string);
  '_id': "id";
  '_completedTime': "completedTime";
}
