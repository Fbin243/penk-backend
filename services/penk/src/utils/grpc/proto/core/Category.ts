// Original file: ../../proto/core/character_message.proto

import type { CategoryStyle as _core_CategoryStyle, CategoryStyle__Output as _core_CategoryStyle__Output } from '../core/CategoryStyle';
import type { Metric as _core_Metric, Metric__Output as _core_Metric__Output } from '../core/Metric';

export interface Category {
  'id'?: (string);
  'name'?: (string);
  'description'?: (string);
  'style'?: (_core_CategoryStyle | null);
  'metrics'?: (_core_Metric)[];
}

export interface Category__Output {
  'id': (string);
  'name': (string);
  'description': (string);
  'style': (_core_CategoryStyle__Output | null);
  'metrics': (_core_Metric__Output)[];
}
