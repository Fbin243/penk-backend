// Original file: ../../proto/core/category_message.proto

import type { CategoryStyleInput as _core_CategoryStyleInput, CategoryStyleInput__Output as _core_CategoryStyleInput__Output } from '../core/CategoryStyleInput';

export interface CategoryInput {
  'id'?: (string);
  'name'?: (string);
  'description'?: (string);
  'style'?: (_core_CategoryStyleInput | null);
  '_id'?: "id";
  '_description'?: "description";
}

export interface CategoryInput__Output {
  'id'?: (string);
  'name': (string);
  'description'?: (string);
  'style': (_core_CategoryStyleInput__Output | null);
  '_id': "id";
  '_description': "description";
}
