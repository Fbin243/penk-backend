// Original file: ../../proto/core/entity_type_message.proto

export const EntityType = {
  Task: 'Task',
  Habit: 'Habit',
} as const;

export type EntityType =
  | 'Task'
  | 0
  | 'Habit'
  | 1

export type EntityType__Output = typeof EntityType[keyof typeof EntityType]
