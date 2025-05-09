// Original file: ../../proto/core/goal_message.proto

export const MetricCondition = {
  lt: 'lt',
  lte: 'lte',
  eq: 'eq',
  gte: 'gte',
  gt: 'gt',
  ir: 'ir',
} as const;

export type MetricCondition =
  | 'lt'
  | 0
  | 'lte'
  | 1
  | 'eq'
  | 2
  | 'gte'
  | 3
  | 'gt'
  | 4
  | 'ir'
  | 5

export type MetricCondition__Output = typeof MetricCondition[keyof typeof MetricCondition]
