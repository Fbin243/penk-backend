// Original file: ../../proto/core/habit_message.proto

import type { CompletionType as _core_CompletionType, CompletionType__Output as _core_CompletionType__Output } from '../core/CompletionType';
import type { HabitReset as _core_HabitReset, HabitReset__Output as _core_HabitReset__Output } from '../core/HabitReset';

export interface HabitInput {
  'id'?: (string);
  'categoryId'?: (string);
  'completionType'?: (_core_CompletionType);
  'name'?: (string);
  'value'?: (number | string);
  'unit'?: (string);
  'rrule'?: (string);
  'resetDuration'?: (_core_HabitReset);
  '_id'?: "id";
  '_categoryId'?: "categoryId";
  '_unit'?: "unit";
}

export interface HabitInput__Output {
  'id'?: (string);
  'categoryId'?: (string);
  'completionType': (_core_CompletionType__Output);
  'name': (string);
  'value': (number);
  'unit'?: (string);
  'rrule': (string);
  'resetDuration': (_core_HabitReset__Output);
  '_id': "id";
  '_categoryId': "categoryId";
  '_unit': "unit";
}
