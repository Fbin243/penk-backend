// Original file: ../../proto/core/metric_message.proto


export interface MetricInput {
  'id'?: (string);
  'categoryId'?: (string);
  'name'?: (string);
  'value'?: (number | string);
  'unit'?: (string);
  '_id'?: "id";
  '_categoryId'?: "categoryId";
}

export interface MetricInput__Output {
  'id'?: (string);
  'categoryId'?: (string);
  'name': (string);
  'value': (number);
  'unit': (string);
  '_id': "id";
  '_categoryId': "categoryId";
}
