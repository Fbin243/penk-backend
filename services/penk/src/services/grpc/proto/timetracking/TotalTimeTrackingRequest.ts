// Original file: ../../proto/timetracking/timetracking_message.proto

import type { Long } from '@grpc/proto-loader';

export interface TotalTimeTrackingRequest {
  'characterId'?: (string);
  'timestamp'?: (number | string | Long);
}

export interface TotalTimeTrackingRequest__Output {
  'characterId': (string);
  'timestamp': (string);
}
