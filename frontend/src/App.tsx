import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider, createTheme, CssBaseline } from '@mui/material';
import { AuthProvider, useAuth } from './context/AuthContext';
import Login from './components/Auth/Login';
import Register from './components/Auth/Register';
import CollectionList from './components/Collection/CollectionList';
import Header from './components/Layout/Header';
import AddCollection from './components/Collection/AddCollection';
import NoticeDialog from './components/Common/NoticeDialog';

// 创建主题
const theme = createTheme({
  palette: {
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
  },
});

// 受保护的路由组件
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, loading } = useAuth();

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  return <>{children}</>;
};

// 主应用内容
const AppContent: React.FC = () => {
  const [showAddCollection, setShowAddCollection] = React.useState(false);
  const [searchQuery, setSearchQuery] = React.useState('');
  const [showNotice, setShowNotice] = React.useState(false);
  const { isAuthenticated } = useAuth();

  // 在组件加载时和认证状态变化时显示通知
  useEffect(() => {
    if (isAuthenticated) {
      const hasSeenNotice = localStorage.getItem('hasSeenNotice');
      const currentTime = new Date().getTime();
      const lastNoticeTime = parseInt(localStorage.getItem('lastNoticeTime') || '0');
      
      if (!hasSeenNotice || (currentTime - lastNoticeTime > 24 * 60 * 60 * 1000)) {
        setShowNotice(true);
      }
    }
  }, [isAuthenticated]);

  const handleSearch = (query: string) => {
    setSearchQuery(query);
  };

  const handleAddCollection = () => {
    setShowAddCollection(true);
  };

  const handleCollectionAdded = () => {
    setShowAddCollection(false);
    // 可以在这里刷新收藏列表
  };

  const handleNoticeClose = () => {
    setShowNotice(false);
    localStorage.setItem('hasSeenNotice', 'true');
    localStorage.setItem('lastNoticeTime', new Date().getTime().toString());
  };

  return (
    <>
      <Header 
        onSearch={handleSearch}
        onAddClick={handleAddCollection}
      />
      <Routes>
        <Route path="/" element={
          <ProtectedRoute>
            <CollectionList searchQuery={searchQuery} />
          </ProtectedRoute>
        } />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
      
      <AddCollection 
        open={showAddCollection} 
        onClose={() => setShowAddCollection(false)}
        onAdded={handleCollectionAdded}
      />

      {/* 通知对话框 */}
      <NoticeDialog 
        open={showNotice} 
        onClose={handleNoticeClose} 
      />
    </>
  );
};

// 主应用
const App: React.FC = () => {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <AuthProvider>
          <AppContent />
        </AuthProvider>
      </Router>
    </ThemeProvider>
  );
};

export default App;
