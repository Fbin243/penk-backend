export * from "./__generated__";

export interface Profile {
  id: string;
  email: string;
  name: string;
  currentCharacterId: string;
}

export interface TimeTracking {
  id: string;
  characterId: string;
  categoryId?: string;
  startTime: string;
  endTime?: string;
}

export interface TimeTrackingWithFish {
  timeTracking: TimeTracking;
  normal: number;
  gold: number;
}

export interface Category {
  id: string;
  name: string;
}

export interface Character {
  id: string;
  name: string;
  categories: Category[];
}
