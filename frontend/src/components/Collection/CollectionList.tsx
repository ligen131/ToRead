import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Card, 
  CardContent, 
  CardActionArea, 
  Chip, 
  Grid,
  Button,
  Container,
  Divider
} from '@mui/material';
import { 
  Link as LinkIcon, 
  Image as ImageIcon, 
  VideoLibrary as VideoIcon 
} from '@mui/icons-material';
import moment from 'moment';
import { collectionAPI } from '../../services/api';
import Loading from '../Common/Loading';
import ErrorMessage from '../Common/ErrorMessage';
import Summary from './Summary';
import AddCollection from './AddCollection';

interface Collection {
  collection_id: number;
  url: string;
  type: 'text' | 'image' | 'video';
  title: string;
  description: string;
  tags: string[];
  created_at: number;
}

interface CollectionListProps {
  searchQuery?: string;
}

const CollectionList: React.FC<CollectionListProps> = ({ searchQuery: externalSearchQuery = '' }) => {
  const [collections, setCollections] = useState<Collection[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [tags, setTags] = useState<string[]>([]);
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [internalSearchQuery, setInternalSearchQuery] = useState('');
  const [showSummary, setShowSummary] = useState(false);
  const [showAddCollection, setShowAddCollection] = useState(false);

  // 使用外部搜索查询更新内部状态
  useEffect(() => {
    setInternalSearchQuery(externalSearchQuery);
  }, [externalSearchQuery]);

  // 获取所有标签
  const fetchTags = async () => {
    try {
      const response = await collectionAPI.getTags();
      if (response.code === 200) {
        setTags(response.data.tags);
      }
    } catch (err: any) {
      console.error('获取标签失败:', err);
    }
  };

  // 获取收藏列表
  const fetchCollections = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const params: { search?: string; tags?: string[] } = {};
      if (internalSearchQuery) {
        params.search = internalSearchQuery;
      }
      if (selectedTags.length > 0) {
        params.tags = selectedTags;
      }
      
      const response = await collectionAPI.getList(params);
      
      if (response.code === 200) {
        // 按创建时间从新到旧排序
        const sortedCollections = [...(response.data.collections || [])].sort(
          (a, b) => b.created_at - a.created_at
        );
        setCollections(sortedCollections);
      } else {
        setError(response.msg || '获取收藏列表失败');
      }
    } catch (err: any) {
      setError(err.response?.data?.msg || '获取收藏列表失败，请检查网络连接');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTags();
    fetchCollections();
  }, []);

  useEffect(() => {
    fetchCollections();
  }, [internalSearchQuery, selectedTags]);

  const handleTagClick = (tag: string) => {
    setSelectedTags(prev => {
      if (prev.includes(tag)) {
        return prev.filter(t => t !== tag);
      } else {
        return [...prev, tag];
      }
    });
  };

  const handleAddCollection = () => {
    setShowAddCollection(true);
  };

  const handleCollectionAdded = () => {
    setShowAddCollection(false);
    fetchCollections();
    fetchTags();
  };

  const getIconByType = (type: string) => {
    switch (type) {
      case 'image':
        return <ImageIcon />;
      case 'video':
        return <VideoIcon />;
      default:
        return <LinkIcon />;
    }
  };

  return (
    <Container>
      <Box sx={{ py: 3 }}>
        {/* 标签筛选区域 */}
        <Box sx={{ mb: 3 }}>
          <Typography variant="h6" gutterBottom>
            标签筛选
          </Typography>
          <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
            {tags.map((tag) => (
              <Chip
                key={tag}
                label={tag}
                clickable
                color={selectedTags.includes(tag) ? "primary" : "default"}
                onClick={() => handleTagClick(tag)}
                sx={{ mb: 1 }}
              />
            ))}
          </Box>
        </Box>

        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
          <Typography variant="h5" component="h2">
            {selectedTags.length > 0 || internalSearchQuery ? '筛选结果' : '我的收藏'}
            {internalSearchQuery && (
              <Typography variant="body1" component="span" sx={{ ml: 2, color: 'text.secondary' }}>
                搜索: "{internalSearchQuery}"
              </Typography>
            )}
          </Typography>
          <Button 
            variant="contained" 
            color="primary"
            onClick={() => setShowSummary(!showSummary)}
          >
            {showSummary ? '隐藏总结' : '显示总结'}
          </Button>
        </Box>

        {/* 总结区域 */}
        {showSummary && (
          <Summary 
            searchQuery={internalSearchQuery} 
            selectedTags={selectedTags} 
          />
        )}

        {/* 添加收藏对话框 */}
        <AddCollection 
          open={showAddCollection} 
          onClose={() => setShowAddCollection(false)}
          onAdded={handleCollectionAdded}
        />

        {/* 收藏列表 */}
        {loading ? (
          <Loading />
        ) : error ? (
          <ErrorMessage message={error} />
        ) : collections.length === 0 ? (
          <Box sx={{ textAlign: 'center', py: 5 }}>
            <Typography variant="h6" color="text.secondary">
              暂无收藏内容
            </Typography>
            <Button 
              variant="contained" 
              sx={{ mt: 2 }}
              onClick={handleAddCollection}
            >
              添加收藏
            </Button>
          </Box>
        ) : (
          <Grid container spacing={3}>
            {collections.map((collection) => (
              <Grid item xs={12} sm={6} md={4} key={collection.collection_id}>
                <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
                  <CardActionArea 
                    component="a"
                    href={collection.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column', alignItems: 'stretch' }}
                  >
                    <CardContent sx={{ flexGrow: 1 }}>
                      <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                        {getIconByType(collection.type)}
                        <Typography variant="h6" component="div" sx={{ ml: 1 }}>
                          {collection.title}
                        </Typography>
                      </Box>
                      <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                        {collection.description}
                      </Typography>
                      <Typography 
                        variant="body2" 
                        color="primary"
                        sx={{ 
                          mb: 2, 
                          overflow: 'hidden', 
                          textOverflow: 'ellipsis', 
                          whiteSpace: 'nowrap',
                          fontStyle: 'italic'
                        }}
                      >
                        {collection.url}
                      </Typography>
                      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, mb: 2 }}>
                        {collection.tags.map((tag) => (
                          <Chip
                            key={tag}
                            label={`#${tag}`}
                            size="small"
                            onClick={(e) => {
                              e.preventDefault();
                              e.stopPropagation();
                              handleTagClick(tag);
                            }}
                          />
                        ))}
                      </Box>
                      <Divider sx={{ mb: 1 }} />
                      <Typography variant="caption" color="text.secondary">
                        {moment.unix(collection.created_at).format('YYYY-MM-DD HH:mm')}
                      </Typography>
                    </CardContent>
                  </CardActionArea>
                </Card>
              </Grid>
            ))}
          </Grid>
        )}
      </Box>
    </Container>
  );
};

export default CollectionList;
