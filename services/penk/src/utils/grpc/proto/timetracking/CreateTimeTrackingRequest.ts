// Original file: ../../proto/timetracking/timetracking_message.proto

import type { Long } from '@grpc/proto-loader';

export interface CreateTimeTrackingRequest {
  'characterId'?: (string);
  'categoryId'?: (string);
  'startTime'?: (number | string | Long);
  '_categoryId'?: "categoryId";
}

export interface CreateTimeTrackingRequest__Output {
  'characterId': (string);
  'categoryId'?: (string);
  'startTime': (string);
  '_categoryId': "categoryId";
}
