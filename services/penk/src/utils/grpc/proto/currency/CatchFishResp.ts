// Original file: ../../proto/currency/currency_message.proto


// Original file: ../../proto/currency/currency_message.proto

export const _currency_CatchFishResp_FishType = {
  None: 'None',
  Normal: 'Normal',
  Gold: 'Gold',
} as const;

export type _currency_CatchFishResp_FishType =
  | 'None'
  | 0
  | 'Normal'
  | 1
  | 'Gold'
  | 2

export type _currency_CatchFishResp_FishType__Output = typeof _currency_CatchFishResp_FishType[keyof typeof _currency_CatchFishResp_FishType]

export interface CatchFishResp {
  'fishType'?: (_currency_CatchFishResp_FishType);
  'number'?: (number);
}

export interface CatchFishResp__Output {
  'fishType': (_currency_CatchFishResp_FishType__Output);
  'number': (number);
}
