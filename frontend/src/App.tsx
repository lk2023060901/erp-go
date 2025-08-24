/**
 * 主应用组件
 */

import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider, App as AntApp } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';
import { AuthProvider } from './contexts/AuthContext';
import { RouteGuard } from './components/RouteGuard';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import UserCenter from './pages/UserCenter';
import UserManagement from './pages/UserManagement';
import RoleManagement from './pages/RoleManagement';
import OrganizationManagement from './pages/OrganizationManagement';
import SystemManagement from './pages/SystemManagement';
import Layout from './components/Layout';
import UserEditPage from './pages/UserEditPage';
import './App.css';

// 设置dayjs中文语言
dayjs.locale('zh-cn');

// Ant Design 主题配置
const theme = {
  token: {
    colorPrimary: '#667eea',
    borderRadius: 8,
    wireframe: false,
  },
  components: {
    Button: {
      borderRadius: 8,
    },
    Input: {
      borderRadius: 8,
    },
    Card: {
      borderRadius: 12,
    },
  },
};

function App() {
  return (
    <ConfigProvider locale={zhCN} theme={theme}>
      <AntApp>
        <AuthProvider>
          <Router>
            <div className="app">
              <Routes>
                {/* 公共路由 */}
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                
                {/* 受保护的路由 */}
                <Route
                  path="/*"
                  element={
                    <RouteGuard requireAuth>
                      <Layout>
                        <Routes>
                          <Route path="/dashboard" element={<Dashboard />} />
                          <Route path="/users" element={<UserCenter />} />
                          <Route path="/users/list" element={<UserManagement />} />
                          <Route path="/users/create" element={<UserEditPage />} />
                          <Route path="/users/edit/:id" element={<UserEditPage />} />
                          <Route path="/users/profile" element={<UserEditPage />} />
                          <Route path="/users/import" element={<div>批量导入用户</div>} />
                          <Route path="/roles" element={<RoleManagement />} />
                          <Route path="/roles/list" element={<RoleManagement />} />
                          <Route path="/roles/create" element={<div>创建角色</div>} />
                          <Route path="/roles/edit/:id" element={<div>编辑角色</div>} />
                          <Route path="/roles/import" element={<div>批量导入角色</div>} />
                          <Route path="/permissions" element={<div>权限管理</div>} />
                          <Route path="/organizations" element={<OrganizationManagement />} />
                          <Route path="/settings" element={<SystemManagement />} />
                          <Route path="/profile" element={<div>个人资料</div>} />
                          
                          {/* 默认重定向到仪表板 */}
                          <Route path="/" element={<Navigate to="/dashboard" replace />} />
                          <Route path="*" element={<div>页面不存在</div>} />
                        </Routes>
                      </Layout>
                    </RouteGuard>
                  }
                />
              </Routes>
            </div>
          </Router>
        </AuthProvider>
      </AntApp>
    </ConfigProvider>
  );
}

export default App;