import * as React from 'react';
import Typography from '@mui/material/Typography';
import { TypographyVariant } from '@mui/material';

interface TitleProps {
  children?: React.ReactNode;
  variant?: TypographyVariant;
}

export default function Title({ children, variant = 'h4' }: TitleProps) {
  return (
    <Typography
      component="h2"
      variant={variant}
      gutterBottom
      sx={{
        fontWeight: '800',
        color: '#FFFEFE',
      }}
    >
      {children}
    </Typography>
  );
}
