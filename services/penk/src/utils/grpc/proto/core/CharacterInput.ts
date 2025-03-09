// Original file: ../../proto/core/character_message.proto

import type { CategoryInput as _core_CategoryInput, CategoryInput__Output as _core_CategoryInput__Output } from '../core/CategoryInput';
import type { MetricInput as _core_MetricInput, MetricInput__Output as _core_MetricInput__Output } from '../core/MetricInput';

export interface CharacterInput {
  'id'?: (string);
  'name'?: (string);
  'gender'?: (boolean);
  'tags'?: (string)[];
  'categories'?: (_core_CategoryInput)[];
  'metrics'?: (_core_MetricInput)[];
  '_id'?: "id";
}

export interface CharacterInput__Output {
  'id'?: (string);
  'name': (string);
  'gender': (boolean);
  'tags': (string)[];
  'categories': (_core_CategoryInput__Output)[];
  'metrics': (_core_MetricInput__Output)[];
  '_id': "id";
}
