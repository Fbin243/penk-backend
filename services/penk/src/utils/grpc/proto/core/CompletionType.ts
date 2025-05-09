// Original file: ../../proto/core/habit_message.proto

export const CompletionType = {
  Number: 'Number',
  Time: 'Time',
} as const;

export type CompletionType =
  | 'Number'
  | 0
  | 'Time'
  | 1

export type CompletionType__Output = typeof CompletionType[keyof typeof CompletionType]
