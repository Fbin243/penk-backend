export enum PenKFunction {
  createTimeTracking = "create_time_tracking",
  updateTimeTracking = "update_time_tracking",
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

export interface User {
  id: string;
  name: string;
  currentCharacterId?: string;
  characters: Character[];
}
