import React from 'react';
import { Alert, AlertTitle } from '@mui/material';

interface ErrorMessageProps {
  message: string;
  severity?: 'error' | 'warning' | 'info' | 'success';
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({
  message,
  severity = 'error',
}) => {
  return (
    <Alert severity={severity} sx={{ mt: 2, mb: 2 }}>
      <AlertTitle>{severity === 'error' ? '错误' : '提示'}</AlertTitle>
      {message}
    </Alert>
  );
};

export default ErrorMessage;