// Original file: ../../proto/core/goal_message.proto

export const GoalStatus = {
  Planned: 'Planned',
  InProgress: 'InProgress',
  Completed: 'Completed',
  Overdue: 'Overdue',
} as const;

export type GoalStatus =
  | 'Planned'
  | 0
  | 'InProgress'
  | 1
  | 'Completed'
  | 2
  | 'Overdue'
  | 3

export type GoalStatus__Output = typeof GoalStatus[keyof typeof GoalStatus]
