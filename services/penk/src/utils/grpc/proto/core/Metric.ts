// Original file: ../../proto/core/metric_message.proto

import type { Long } from '@grpc/proto-loader';

export interface Metric {
  'id'?: (string);
  'createdAt'?: (number | string | Long);
  'updatedAt'?: (number | string | Long);
  'characterId'?: (string);
  'categoryId'?: (string);
  'name'?: (string);
  'value'?: (number | string);
  'unit'?: (string);
  '_categoryId'?: "categoryId";
}

export interface Metric__Output {
  'id': (string);
  'createdAt': (string);
  'updatedAt': (string);
  'characterId': (string);
  'categoryId'?: (string);
  'name': (string);
  'value': (number);
  'unit': (string);
  '_categoryId': "categoryId";
}
