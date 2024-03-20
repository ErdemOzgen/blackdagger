import { Box } from '@mui/material';
import moment from 'moment';
import React from 'react';
import {
  Bar,
  BarChart,
  Cell,
  LabelList,
  ResponsiveContainer,
  XAxis,
  YAxis,
} from 'recharts';
import { statusColorMapping } from '../../consts';
import { DAGStatus, SchedulerStatus, Status } from '../../models';

type Props = { data: DAGStatus[] };

type DataFrame = {
  name: string;
  status: SchedulerStatus;
  values: [number, number];
};

function DashboardTimechart({ data: input }: Props) {
  const [data, setData] = React.useState<DataFrame[]>([]);
  React.useEffect(() => {
    const transformedData: DataFrame[] = input.map((dagStatus) => {
      const status: Status | undefined = dagStatus.Status;
      if (!status || status.StartedAt === '-' || !status.FinishedAt) {
        return null;
      }
      const startUnix = Math.max(moment(status.StartedAt).unix(), moment().startOf('day').unix());
      const endUnix = status.FinishedAt && status.FinishedAt !== '-' ? moment(status.FinishedAt).unix() : moment().unix();

      return {
        name: status.Name,
        status: status.Status,
        values: [startUnix, endUnix],
      };
    }).filter((item): item is DataFrame => item !== null);

    setData(transformedData.sort((a, b) => a.values[0] - b.values[0]));
  }, [input]);

  const now = moment();
  const shouldScroll = data.length >= 40;
  // Generating tick values for the XAxis
  const tickValues = Array.from({ length: 24 }, (_, index) => now.startOf('day').add(index, 'hours').unix());

  return (
    <TimelineWrapper shouldScroll={shouldScroll}>
      <ResponsiveContainer
        width="100%"
        minHeight={shouldScroll ? data.length * 12 : undefined}
        height={shouldScroll ? undefined : '90%'}
      >
        <BarChart data={data} layout="vertical">
          <XAxis
            tickFormatter={(unixTime) => moment.unix(unixTime).format('HH:mm')}
            type="number"
            dataKey="values"
            domain={[now.startOf('day').unix(), now.endOf('day').unix()]}
            ticks={tickValues}
          />
          <YAxis dataKey="name" type="category" hide />
          <Bar dataKey="values" fill="lightblue" minPointSize={2}>
            {data.map((entry, index) => (
              // Use statusColorMapping for setting the cell colors based on status
              <Cell key={`cell-${index}`} fill={statusColorMapping[entry.status]?.backgroundColor || '#ccc'} />
            ))}
            {/* Set the LabelList fill color to bone white */}
            <LabelList dataKey="name" position="center" fill="#E3DAC9" />
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </TimelineWrapper>
  );
}




function TimelineWrapper({
  children,
  shouldScroll,
}: {
  children: React.ReactNode;
  shouldScroll: boolean;
}) {
  return shouldScroll ? (
    <Box
      sx={{
        width: '100%',
        maxWidth: '100%',
        height: '90%',
        overflow: 'auto',
      }}
    >
      {children}
    </Box>
  ) : (
    <React.Fragment>{children}</React.Fragment>
  );
}

export default DashboardTimechart;
