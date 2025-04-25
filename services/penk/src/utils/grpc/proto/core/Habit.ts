// Original file: ../../proto/core/habit_message.proto

import type { CompletionType as _core_CompletionType, CompletionType__Output as _core_CompletionType__Output } from '../core/CompletionType';
import type { HabitReset as _core_HabitReset, HabitReset__Output as _core_HabitReset__Output } from '../core/HabitReset';
import type { Long } from '@grpc/proto-loader';

export interface Habit {
  'id'?: (string);
  'createdAt'?: (number | string | Long);
  'updatedAt'?: (number | string | Long);
  'characterId'?: (string);
  'categoryId'?: (string);
  'completionType'?: (_core_CompletionType);
  'name'?: (string);
  'value'?: (number | string);
  'unit'?: (string);
  'rrule'?: (string);
  'reset'?: (_core_HabitReset);
}

export interface Habit__Output {
  'id': (string);
  'createdAt': (string);
  'updatedAt': (string);
  'characterId': (string);
  'categoryId': (string);
  'completionType': (_core_CompletionType__Output);
  'name': (string);
  'value': (number);
  'unit': (string);
  'rrule': (string);
  'reset': (_core_HabitReset__Output);
}
