// Original file: ../../proto/core/habit_log_message.proto


export interface HabitLogInput {
  'timestamp'?: (string);
  'habitId'?: (string);
  'value'?: (number | string);
}

export interface HabitLogInput__Output {
  'timestamp': (string);
  'habitId': (string);
  'value': (number);
}
