import { CSSProperties } from 'react';
import { NodeStatus } from './models';
import { SchedulerStatus } from './models';

type statusColorMapping = {
  [key: number]: CSSProperties;
};
export const statusColorMapping: statusColorMapping = {
  [SchedulerStatus.None]: { backgroundColor: '#e5e7eb', color : '#111827' },
  [SchedulerStatus.Running]: { backgroundColor: '#2563eb', color: "#ffffff"  },
  [SchedulerStatus.Error]: { backgroundColor: '#dc2626', color : "#ffffff" },
  [SchedulerStatus.Cancel]: { backgroundColor: 'gray', color : "#ffffff" },
  [SchedulerStatus.Success]: { backgroundColor: '#16a34a', color: '#ffffff' },
  [SchedulerStatus.Skipped_Unused]: { backgroundColor: '#6b7280', color: '#ffffff' },
};

export const nodeStatusColorMapping = {
  [NodeStatus.None]: statusColorMapping[SchedulerStatus.None],
  [NodeStatus.Running]: statusColorMapping[SchedulerStatus.Running],
  [NodeStatus.Error]: statusColorMapping[SchedulerStatus.Error],
  [NodeStatus.Cancel]: statusColorMapping[SchedulerStatus.Cancel],
  [NodeStatus.Success]: statusColorMapping[SchedulerStatus.Success],
  [NodeStatus.Skipped]: statusColorMapping[SchedulerStatus.Skipped_Unused],
};

export const stepTabColStyles = [
  { maxWidth: '60px' },
  { maxWidth: '200px' },
  { maxWidth: '150px' },
  { maxWidth: '150px' },
  { maxWidth: '150px' },
  { maxWidth: '130px' },
  { maxWidth: '130px' },
  { maxWidth: '100px' },
  { maxWidth: '100px' },
  {},
];
