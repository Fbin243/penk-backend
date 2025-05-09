// Original file: ../../proto/core/category_message.proto

import type { CategoryStyle as _core_CategoryStyle, CategoryStyle__Output as _core_CategoryStyle__Output } from '../core/CategoryStyle';
import type { Long } from '@grpc/proto-loader';

export interface Category {
  'id'?: (string);
  'createdAt'?: (number | string | Long);
  'updatedAt'?: (number | string | Long);
  'characterId'?: (string);
  'name'?: (string);
  'description'?: (string);
  'style'?: (_core_CategoryStyle | null);
  '_description'?: "description";
}

export interface Category__Output {
  'id': (string);
  'createdAt': (string);
  'updatedAt': (string);
  'characterId': (string);
  'name': (string);
  'description'?: (string);
  'style': (_core_CategoryStyle__Output | null);
  '_description': "description";
}
