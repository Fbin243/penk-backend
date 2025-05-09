// Original file: ../../proto/core/entity_type_message.proto

export const EntityType = {
  TaskType: 'TaskType',
  HabitType: 'HabitType',
} as const;

export type EntityType =
  | 'TaskType'
  | 0
  | 'HabitType'
  | 1

export type EntityType__Output = typeof EntityType[keyof typeof EntityType]
