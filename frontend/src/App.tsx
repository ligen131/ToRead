import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider, createTheme, CssBaseline } from '@mui/material';
import { AuthProvider, useAuth } from './context/AuthContext';
import Login from './components/Auth/Login';
import Register from './components/Auth/Register';
import CollectionList from './components/Collection/CollectionList';
import Header from './components/Layout/Header';
import AddCollection from './components/Collection/AddCollection';

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
