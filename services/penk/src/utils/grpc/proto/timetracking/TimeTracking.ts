// Original file: ../../proto/timetracking/timetracking_message.proto

import type { Long } from '@grpc/proto-loader';

export interface TimeTracking {
  'id'?: (string);
  'characterId'?: (string);
  'categoryId'?: (string);
  'startTime'?: (number | string | Long);
  'endTime'?: (number | string | Long);
  '_categoryId'?: "categoryId";
  '_endTime'?: "endTime";
}

export interface TimeTracking__Output {
  'id': (string);
  'characterId': (string);
  'categoryId'?: (string);
  'startTime': (string);
  'endTime'?: (string);
  '_categoryId': "categoryId";
  '_endTime': "endTime";
}
