import React, { useState } from 'react';
import { 
  AppBar, 
  Toolbar, 
  Typography, 
  Button, 
  IconButton,
  TextField,
  InputAdornment,
  Box,
  Menu,
  MenuItem,
  Avatar
} from '@mui/material';
import { 
  Search as SearchIcon,
  AccountCircle,
  Add as AddIcon
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

interface HeaderProps {
  onSearch?: (query: string) => void;
  onAddClick?: () => void;
}

const Header: React.FC<HeaderProps> = ({ onSearch, onAddClick }) => {
  const { user, logout, isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    handleClose();
    logout();
  };

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log('Search submitted:', searchQuery);
    if (onSearch) {
      onSearch(searchQuery);
    }
  };

  return (
    <AppBar position="static">
      <Toolbar>
        <Typography 
          variant="h6" 
          component="div" 
          sx={{ flexGrow: 0, cursor: 'pointer', mr: 2 }}
          onClick={() => navigate('/')}
        >
          ToRead
        </Typography>

        {isAuthenticated && (
          <>
            <Box component="form" onSubmit={handleSearchSubmit} sx={{ flexGrow: 1 }}>
              <TextField
                placeholder="搜索收藏..."
                size="small"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    e.preventDefault();
                    if (onSearch) {
                      onSearch(searchQuery);
                    }
                  }
                }}
                sx={{ 
                  backgroundColor: 'rgba(255, 255, 255, 0.15)',
                  borderRadius: 1,
                  '& .MuiOutlinedInput-root': {
                    color: 'white',
                    '& fieldset': { border: 'none' },
                  },
                  width: { xs: '100%', sm: '50%', md: '40%' }
                }}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <SearchIcon 
                        sx={{ color: 'white', cursor: 'pointer' }} 
                        onClick={() => {
                          if (onSearch) {
                            onSearch(searchQuery);
                          }
                        }}
                      />
                    </InputAdornment>
                  ),
                }}
              />
            </Box>

            <IconButton 
              color="inherit" 
              onClick={onAddClick}
              sx={{ mr: 2 }}
            >
              <AddIcon />
            </IconButton>

            <div>
              <IconButton
                size="large"
                onClick={handleMenu}
                color="inherit"
              >
                <Avatar sx={{ width: 32, height: 32, bgcolor: 'secondary.main' }}>
                  {user?.user_name.charAt(0).toUpperCase()}
                </Avatar>
              </IconButton>
              <Menu
                id="menu-appbar"
                anchorEl={anchorEl}
                anchorOrigin={{
                  vertical: 'bottom',
                  horizontal: 'right',
                }}
                keepMounted
                transformOrigin={{
                  vertical: 'top',
                  horizontal: 'right',
                }}
                open={Boolean(anchorEl)}
                onClose={handleClose}
              >
                <MenuItem disabled>
                  {user?.user_name}
                </MenuItem>
                <MenuItem onClick={handleLogout}>退出登录</MenuItem>
              </Menu>
            </div>
          </>
        )}

        {!isAuthenticated && (
          <div>
            <Button color="inherit" onClick={() => navigate('/login')}>登录</Button>
            <Button color="inherit" onClick={() => navigate('/register')}>注册</Button>
          </div>
        )}
      </Toolbar>
    </AppBar>
  );
};

export default Header;
