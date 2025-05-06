import React, { useState } from 'react';
import {
  Box,
  Stack,
  Typography,
  useTheme,
  InputAdornment,
  TextField,
  IconButton,
  Tooltip,
} from '@mui/material';
import { ContentCopy, Search } from '@mui/icons-material';
import { LogFile } from '../../models/api';
import BorderedBox from '../atoms/BorderedBox';
import LabeledItem from '../atoms/LabeledItem';
import LoadingIndicator from '../atoms/LoadingIndicator';
import NodeStatusChip from '../molecules/NodeStatusChip';

type Props = {
  log?: LogFile;
};

function ExecutionLog({ log }: Props) {
  const theme = useTheme();
  const [searchTerm, setSearchTerm] = useState('');

  if (!log) return <LoadingIndicator />;

  const filteredContent = searchTerm
    ? log.Content.split('\n')
        .filter((line: string) =>
          line.toLowerCase().includes(searchTerm.toLowerCase())
        )
        .join('\n')
    : log.Content;

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(log.Content);
    } catch (err) {
      console.error('Failed to copy log output:', err);
    }
  };

  return (
    <Box
      sx={{
        padding: 3,
        backgroundColor: theme.palette.background.paper,
        borderRadius: 2,
        boxShadow: theme.shadows[1],
      }}
    >
      <Typography variant="h5" fontWeight="bold" gutterBottom>
        Execution Log
      </Typography>

      <Stack spacing={2}>
        <LabeledItem label="Log File">
          <Typography variant="body2">{log.LogFile}</Typography>
        </LabeledItem>

        {log.Step && (
          <>
            <LabeledItem label="Step Name">
              <Typography variant="body2">{log.Step.Step.Name}</Typography>
            </LabeledItem>

            <Stack direction={{ xs: 'column', sm: 'row' }} spacing={4}>
              <LabeledItem label="Started At">
                <Typography variant="body2">{log.Step.StartedAt}</Typography>
              </LabeledItem>
              <LabeledItem label="Finished At">
                <Typography variant="body2">{log.Step.FinishedAt}</Typography>
              </LabeledItem>
            </Stack>

            <LabeledItem label="Status">
              <NodeStatusChip status={log.Step.Status}>
                {log.Step.StatusText}
              </NodeStatusChip>
            </LabeledItem>
          </>
        )}
      </Stack>

      <Stack direction="row" spacing={2} alignItems="center" mt={4} mb={1}>
        <TextField
          size="small"
          placeholder="Search logs..."
          variant="outlined"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          sx={{ width: '100%' }}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <Search fontSize="small" />
              </InputAdornment>
            ),
          }}
        />

        <Tooltip title="Copy full log to clipboard">
          <IconButton onClick={handleCopy}>
            <ContentCopy fontSize="small" />
          </IconButton>
        </Tooltip>
      </Stack>

      <BorderedBox
        sx={{
          backgroundColor: '#0f0f0f',
          color: '#f1f1f1',
          fontFamily: 'Roboto Mono, Courier New, monospace',
          fontSize: '1 rem',
          padding: 2,
          height: '60vh',
          overflowY: 'auto',
          borderRadius: 2,
          whiteSpace: 'pre-wrap',
          lineHeight: 1.5,
          wordBreak: 'break-word',
          border: `1px solid ${theme.palette.divider}`,
        }}
        dangerouslySetInnerHTML={{
          __html: filteredContent || '<No log output>',
        }}
      />
    </Box>
  );
}

export default React.memo(ExecutionLog, (prevProps, nextProps) => {
  return prevProps.log?.Content === nextProps.log?.Content;
});
