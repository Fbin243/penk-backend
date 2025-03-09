// Original file: ../../proto/timetracking/timetracking_message.proto

import type { TimeTracking as _timetracking_TimeTracking, TimeTracking__Output as _timetracking_TimeTracking__Output } from '../timetracking/TimeTracking';

export interface TimeTrackingWithFish {
  'timeTracking'?: (_timetracking_TimeTracking | null);
  'normal'?: (number);
  'gold'?: (number);
}

export interface TimeTrackingWithFish__Output {
  'timeTracking': (_timetracking_TimeTracking__Output | null);
  'normal': (number);
  'gold': (number);
}
