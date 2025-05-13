import React, { useState, useEffect } from 'react';
import { 
  Container, 
  Typography, 
  TextField, 
  Button, 
  Box, 
  Paper,
  Link,
  InputAdornment,
  IconButton
} from '@mui/material';
import { Visibility, VisibilityOff } from '@mui/icons-material';
import { Link as RouterLink } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import ErrorMessage from '../Common/ErrorMessage';
import Loading from '../Common/Loading';
import NoticeDialog from '../Common/NoticeDialog';

const Login: React.FC = () => {
  const { login, loading, error } = useAuth();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [showNotice, setShowNotice] = useState(false);
  const [formErrors, setFormErrors] = useState({
    username: '',
    password: ''
  });

  useEffect(() => {
    const hasSeenNotice = localStorage.getItem('hasSeenNotice');
    const currentTime = new Date().getTime();
    const lastNoticeTime = parseInt(localStorage.getItem('lastNoticeTime') || '0');
    
    if (!hasSeenNotice || (currentTime - lastNoticeTime > 24 * 60 * 60 * 1000)) {
      setShowNotice(true);
    }
  }, []);

  const validateForm = () => {
    let isValid = true;
    const errors = {
      username: '',
      password: ''
    };

    if (!username.trim()) {
      errors.username = '用户名不能为空';
      isValid = false;
    }

    if (!password) {
      errors.password = '密码不能为空';
      isValid = false;
    }

    setFormErrors(errors);
    return isValid;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (validateForm()) {
      await login(username, password);
    }
  };

  const handleNoticeClose = () => {
    setShowNotice(false);
    localStorage.setItem('hasSeenNotice', 'true');
    localStorage.setItem('lastNoticeTime', new Date().getTime().toString());
  };

  return (
    <Container maxWidth="sm">
      <Box 
        sx={{ 
          mt: 8, 
          display: 'flex', 
          flexDirection: 'column', 
          alignItems: 'center' 
        }}
      >
        <Paper 
          elevation={3} 
          sx={{ 
            p: 4, 
            width: '100%', 
            borderRadius: 2 
          }}
        >
          <Typography component="h1" variant="h5" align="center" gutterBottom>
            登录 ToRead
          </Typography>
          
          {error && <ErrorMessage message={error} />}
          
          {loading ? (
            <Loading message="登录中..." />
          ) : (
            <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2 }}>
              <TextField
                margin="normal"
                required
                fullWidth
                id="username"
                label="用户名"
                name="username"
                autoComplete="username"
                autoFocus
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                error={!!formErrors.username}
                helperText={formErrors.username}
              />
              
              <TextField
                margin="normal"
                required
                fullWidth
                name="password"
                label="密码"
                type={showPassword ? 'text' : 'password'}
                id="password"
                autoComplete="current-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                error={!!formErrors.password}
                helperText={formErrors.password}
                InputProps={{
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        onClick={() => setShowPassword(!showPassword)}
                        edge="end"
                      >
                        {showPassword ? <VisibilityOff /> : <Visibility />}
                      </IconButton>
                    </InputAdornment>
                  )
                }}
              />
              
              <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
              >
                登录
              </Button>
              
              <Box textAlign="center">
                <Link component={RouterLink} to="/register" variant="body2">
                  {"没有账号？立即注册"}
                </Link>
              </Box>
            </Box>
          )}
        </Paper>
      </Box>
      
      {/* 通知对话框 */}
      <NoticeDialog 
        open={showNotice} 
        onClose={handleNoticeClose} 
      />
    </Container>
  );
};

export default Login;
