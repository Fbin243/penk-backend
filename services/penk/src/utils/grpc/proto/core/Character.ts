// Original file: ../../proto/core/character_message.proto

import type { Long } from '@grpc/proto-loader';

export interface Character {
  'id'?: (string);
  'createdAt'?: (number | string | Long);
  'updatedAt'?: (number | string | Long);
  'profileId'?: (string);
  'name'?: (string);
}

export interface Character__Output {
  'id': (string);
  'createdAt': (string);
  'updatedAt': (string);
  'profileId': (string);
  'name': (string);
}
