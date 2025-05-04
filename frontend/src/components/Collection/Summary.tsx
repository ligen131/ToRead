import React, { useState, useEffect } from 'react';
import { 
  Paper, 
  Typography, 
  Box, 
  CircularProgress 
} from '@mui/material';
import { collectionAPI } from '../../services/api';
import ErrorMessage from '../Common/ErrorMessage';

interface SummaryProps {
  searchQuery?: string;
  selectedTags?: string[];
}

const Summary: React.FC<SummaryProps> = ({ searchQuery, selectedTags = [] }) => {
  const [summary, setSummary] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchSummary = async () => {
      try {
        setLoading(true);
        setError(null);
        
        const params: { search?: string; tags?: string[] } = {};
        if (searchQuery) {
          params.search = searchQuery;
        }
        if (selectedTags.length > 0) {
          params.tags = selectedTags;
        }
        
        const response = await collectionAPI.getSummary(params);
        
        if (response.code === 200) {
          setSummary(response.data.summary || '暂无总结内容');
        } else {
          setError(response.msg || '获取总结失败');
        }
      } catch (err: any) {
        setError(err.response?.data?.msg || '获取总结失败，请检查网络连接');
      } finally {
        setLoading(false);
      }
    };

    fetchSummary();
  }, [searchQuery, selectedTags]);

  return (
    <Paper elevation={2} sx={{ p: 3, mb: 4 }}>
      <Typography variant="h6" gutterBottom>
        内容总结
      </Typography>
      
      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
          <CircularProgress />
        </Box>
      ) : error ? (
        <ErrorMessage message={error} />
      ) : (
        <Typography variant="body1" component="div" sx={{ whiteSpace: 'pre-wrap' }}>
          {summary}
        </Typography>
      )}
    </Paper>
  );
};

export default Summary;