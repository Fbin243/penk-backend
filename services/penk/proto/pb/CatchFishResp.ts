// Original file: ../../pkg/proto/currency.proto


// Original file: ../../pkg/proto/currency.proto

export const _pb_CatchFishResp_FishType = {
  None: 'None',
  Normal: 'Normal',
  Gold: 'Gold',
} as const;

export type _pb_CatchFishResp_FishType =
  | 'None'
  | 0
  | 'Normal'
  | 1
  | 'Gold'
  | 2

export type _pb_CatchFishResp_FishType__Output = typeof _pb_CatchFishResp_FishType[keyof typeof _pb_CatchFishResp_FishType]

export interface CatchFishResp {
  'fishType'?: (_pb_CatchFishResp_FishType);
  'number'?: (number);
}

export interface CatchFishResp__Output {
  'fishType': (_pb_CatchFishResp_FishType__Output);
  'number': (number);
}
