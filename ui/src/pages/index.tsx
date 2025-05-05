import React from 'react';
import { ListWorkflowsResponse } from '../models/api';
import { Box, Grid } from '@mui/material';
import { SchedulerStatus } from '../models';
import { statusColorMapping } from '../consts';
import DashboardMetric from '../components/molecules/DashboardMetric';
import DashboardTimechart from '../components/molecules/DashboardTimechart';
import Title from '../components/atoms/Title';
import { AppBarContext } from '../contexts/AppBarContext';
import useSWR from 'swr';
import { useConfig } from '../contexts/ConfigContext';

type metrics = Record<SchedulerStatus, number>;

const defaultMetrics: metrics = {} as metrics;
for (const value in SchedulerStatus) {
  if (!isNaN(Number(value))) {
    const status = Number(value) as SchedulerStatus;
    defaultMetrics[status] = 0;
  }
}

function Dashboard() {
  const [metrics, setMetrics] = React.useState<metrics>(defaultMetrics);
  const appBarContext = React.useContext(AppBarContext);
  const config = useConfig();
  const { data } = useSWR<ListWorkflowsResponse>(
    `/dags?remoteNode=${appBarContext.selectedRemoteNode || 'local'}`,
    null,
    {
      refreshInterval: 10000,
    }
  );

  React.useEffect(() => {
    if (!data) {
      return;
    }
    const m = { ...defaultMetrics };
    data.DAGs?.forEach((wf) => {
      if (wf.Status && wf.Status.Status) {
        const status = wf.Status.Status;
        m[status] += 1;
      }
    });
    setMetrics(m as metrics);
  }, [data]);

  React.useEffect(() => {
    appBarContext.setTitle('Dashboard');
  }, [appBarContext]);

  return (
    <Grid container spacing={1} sx={{ mx: 2, width: '100%' }} paddingTop={2}>
      {(
        [
          [SchedulerStatus.Success, 'Successful'],
          [SchedulerStatus.Error, 'Failed'],
          [SchedulerStatus.Running, 'Running'],
          [SchedulerStatus.Cancel, 'Canceled'],
        ] as Array<[SchedulerStatus, string]>
      ).map(([status, label]) => (
        <Grid item xs={12} sm={6} md={4} lg={3} key={label}>
          <Box
            sx={{
              display: 'flex',
              flexDirection: 'column',
              justifyContent: 'center',
              alignItems: 'center',

              backgroundColor: statusColorMapping[status].backgroundColor,
              color: statusColorMapping[status].color,
              borderRadius: 3,
              boxShadow: 3,
              transition: 'transform 0.25s ease, box-shadow 0.25s ease',
              p: { xs: 2, md: 3 },

              height: 'auto',
              width: 'auto',
              minHeight: 80,
              maxWidth: '100',

              '&:hover': {
                transform: 'scale(1.02)',
                boxShadow: 6,
              },
            }}
          >
            <DashboardMetric
              title={label}
              color={statusColorMapping[status].color}
              value={metrics[status]}
            />
          </Box>
        </Grid>
      ))}

      <Grid item xs={12}>
        <Box
          sx={{
            p: 2,
            height: '100%',
          }}
        >
          <Title>Timeline</Title>
          <DashboardTimechart data={data?.DAGs || []} />
        </Box>
      </Grid>
    </Grid>
  );
}
export default Dashboard;
