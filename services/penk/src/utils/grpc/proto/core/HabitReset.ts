// Original file: ../../proto/core/habit_message.proto

export const HabitReset = {
  Daily: 'Daily',
  Weekly: 'Weekly',
  Monthly: 'Monthly',
} as const;

export type HabitReset =
  | 'Daily'
  | 0
  | 'Weekly'
  | 1
  | 'Monthly'
  | 2

export type HabitReset__Output = typeof HabitReset[keyof typeof HabitReset]
