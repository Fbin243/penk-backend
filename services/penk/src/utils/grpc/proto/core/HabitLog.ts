// Original file: ../../proto/core/habit_log_message.proto


export interface HabitLog {
  'id'?: (string);
  'timestamp'?: (string);
  'habitId'?: (string);
  'value'?: (number | string);
}

export interface HabitLog__Output {
  'id': (string);
  'timestamp': (string);
  'habitId': (string);
  'value': (number);
}
