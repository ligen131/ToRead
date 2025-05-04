import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  CircularProgress,
  Box,
  Typography
} from '@mui/material';
import { collectionAPI } from '../../services/api';
import ErrorMessage from '../Common/ErrorMessage';

interface AddCollectionProps {
  open: boolean;
  onClose: () => void;
  onAdded: () => void;
}

const AddCollection: React.FC<AddCollectionProps> = ({ open, onClose, onAdded }) => {
  const [url, setUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [urlError, setUrlError] = useState('');

  const validateUrl = (value: string) => {
    if (!value) {
      return '请输入URL';
    }
    
    try {
      new URL(value);
      return '';
    } catch (e) {
      return 'URL格式不正确';
    }
  };

  const handleSubmit = async () => {
    const urlValidationError = validateUrl(url);
    if (urlValidationError) {
      setUrlError(urlValidationError);
      return;
    }

    try {
      setLoading(true);
      setError(null);
      
      const response = await collectionAPI.addCollection({ url });
      
      if (response.code === 200) {
        setUrl('');
        setLoading(false); // 先停止加载
        onAdded(); // 调用回调函数
        onClose(); // 关闭对话框
      } else {
        setError(response.msg || '添加收藏失败');
        setLoading(false);
      }
    } catch (err: any) {
      setError(err.response?.data?.msg || '添加收藏失败，请检查网络连接');
    } finally {
      setLoading(false);
    }
  };

  const handleUrlChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUrl(e.target.value);
    if (urlError) {
      setUrlError('');
    }
  };

  return (
    <Dialog 
      open={open} 
      onClose={!loading ? onClose : undefined}
      fullWidth
      maxWidth="sm"
    >
      <DialogTitle>添加收藏链接</DialogTitle>
      <DialogContent>
        {error && <ErrorMessage message={error} />}
        
        <TextField
          autoFocus
          margin="dense"
          id="url"
          label="URL"
          type="url"
          fullWidth
          variant="outlined"
          value={url}
          onChange={handleUrlChange}
          error={!!urlError}
          helperText={urlError}
          disabled={loading}
          placeholder="https://example.com"
        />
        
        {loading && (
          <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', mt: 3 }}>
            <CircularProgress />
            <Typography variant="body2" sx={{ mt: 2 }}>
              正在处理链接，这可能需要一些时间...
            </Typography>
          </Box>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} disabled={loading}>取消</Button>
        <Button 
          onClick={handleSubmit} 
          variant="contained" 
          disabled={loading || !url}
        >
          添加
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default AddCollection;
