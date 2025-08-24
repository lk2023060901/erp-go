/**
 * 主布局组件
 */

import React, { useState } from 'react';
import {
  Layout as AntLayout,
  Menu,
  Avatar,
  Dropdown,
  Button,
  Breadcrumb,
  Typography,
  Badge,
  Divider
} from 'antd';
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  DashboardOutlined,
  UserOutlined,
  SafetyCertificateOutlined,
  ApartmentOutlined,
  SettingOutlined,
  LogoutOutlined,
  BellOutlined,
  QuestionCircleOutlined
} from '@ant-design/icons';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

const { Header, Sider, Content } = AntLayout;
const { Text } = Typography;

interface LayoutProps {
  children: React.ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout } = useAuth();
  
  const [collapsed, setCollapsed] = useState(false);

  // 菜单项配置
  const menuItems = [
    {
      key: '/dashboard',
      icon: <DashboardOutlined />,
      label: '仪表板',
    },
    {
      key: '/users',
      icon: <UserOutlined />,
      label: '用户中心',
    },
    {
      key: '/organizations',
      icon: <ApartmentOutlined />,
      label: '组织架构',
    },
    {
      key: '/settings',
      icon: <SettingOutlined />,
      label: '系统设置',
    },
  ];

  // 用户下拉菜单
  const userMenuItems = [
    {
      key: 'profile',
      label: '个人资料',
      icon: <UserOutlined />,
      onClick: () => navigate('/profile'),
    },
    {
      key: 'settings',
      label: '账户设置',
      icon: <SettingOutlined />,
      onClick: () => navigate('/account-settings'),
    },
    {
      type: 'divider' as const,
    },
    {
      key: 'help',
      label: '帮助中心',
      icon: <QuestionCircleOutlined />,
      onClick: () => window.open('/help', '_blank'),
    },
    {
      key: 'logout',
      label: '退出登录',
      icon: <LogoutOutlined />,
      onClick: handleLogout,
    },
  ];

  // 处理菜单点击
  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key);
  };

  // 处理登出
  async function handleLogout() {
    try {
      await logout();
      navigate('/login');
    } catch (error) {
      console.error('Logout failed:', error);
    }
  }

  // 获取面包屑
  const getBreadcrumb = () => {
    const pathMap: { [key: string]: string } = {
      '/dashboard': '仪表板',
      '/users': '用户中心',
      '/users/create': '创建用户',
      '/users/list': '用户管理',
      '/users/edit': '编辑用户',
      '/users/profile': '个人资料',
      '/roles': '角色管理',
      '/roles/list': '角色列表',
      '/roles/create': '创建角色',
      '/roles/edit': '编辑角色', 
      '/permissions': '权限管理',
      '/organizations': '组织架构',
      '/settings': '系统设置',
      '/profile': '个人资料',
    };

    const pathSegments = location.pathname.split('/').filter(Boolean);
    const breadcrumbItems = pathSegments.map((segment, index) => {
      const path = '/' + pathSegments.slice(0, index + 1).join('/');
      const isLast = index === pathSegments.length - 1;
      
      return {
        title: isLast ? pathMap[path] || segment : (
          <a 
            href="#" 
            onClick={(e) => {
              e.preventDefault();
              navigate(path);
            }}
            style={{ color: '#1890ff', textDecoration: 'none' }}
          >
            {pathMap[path] || segment}
          </a>
        ),
      };
    });

    return [
      { 
        title: (
          <a 
            href="#" 
            onClick={(e) => {
              e.preventDefault();
              navigate('/dashboard');
            }}
            style={{ color: '#1890ff', textDecoration: 'none' }}
          >
            首页
          </a>
        )
      }, 
      ...breadcrumbItems
    ];
  };

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        width={240}
        style={{
          background: '#fff',
          boxShadow: '2px 0 6px rgba(0,21,41,.05)',
        }}
      >
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: collapsed ? 'center' : 'flex-start',
            padding: collapsed ? 0 : '0 24px',
            borderBottom: '1px solid #f0f0f0',
          }}
        >
          <SafetyCertificateOutlined
            style={{
              fontSize: 24,
              color: '#667eea',
              marginRight: collapsed ? 0 : 12,
            }}
          />
          {!collapsed && (
            <Text strong style={{ fontSize: 16, color: '#1f2937' }}>
              ERP系统
            </Text>
          )}
        </div>
        
        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={handleMenuClick}
          style={{
            border: 'none',
            padding: '8px',
          }}
        />
      </Sider>
      
      <AntLayout>
        <Header
          style={{
            background: '#fff',
            padding: '0 24px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            boxShadow: '0 2px 8px rgba(0,0,0,.06)',
            borderBottom: '1px solid #f0f0f0',
            height: 64,
          }}
        >
          <div style={{ 
            display: 'flex', 
            alignItems: 'center', 
            height: '100%' 
          }}>
            <Button
              type="text"
              icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
              onClick={() => setCollapsed(!collapsed)}
              style={{
                fontSize: '16px',
                width: 40,
                height: 40,
                marginRight: 8,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            />
            
            <div style={{
              display: 'flex',
              alignItems: 'center',
              height: 40,
            }}>
              <Breadcrumb 
                items={getBreadcrumb()} 
                style={{
                  lineHeight: '40px',
                  margin: 0,
                }}
              />
            </div>
          </div>
          
          <div style={{
            display: 'flex',
            alignItems: 'center',
            height: '100%',
            gap: '16px',
          }}>
            <Badge count={5} size="small">
              <Button
                type="text"
                icon={<BellOutlined />}
                style={{ 
                  fontSize: '16px',
                  width: 40,
                  height: 40,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}
              />
            </Badge>
            
            <Divider 
              type="vertical" 
              style={{ 
                height: '24px',
                margin: 0,
              }} 
            />
            
            <Dropdown
              menu={{ items: userMenuItems }}
              placement="bottomRight"
              trigger={['click']}
            >
              <div style={{ 
                cursor: 'pointer', 
                display: 'flex', 
                alignItems: 'center', 
                gap: '8px',
                height: 40,
              }}>
                <Avatar
                  size="small"
                  src={user?.avatar_url}
                  icon={!user?.avatar_url ? <UserOutlined /> : undefined}
                  style={{ backgroundColor: '#667eea' }}
                />
                <Text style={{ 
                  color: '#1f2937',
                  lineHeight: '40px',
                }}>
                  {user?.first_name || user?.username || '用户'}
                </Text>
              </div>
            </Dropdown>
          </div>
        </Header>
        
        <Content
          style={{
            margin: '24px',
            padding: '24px',
            background: '#fff',
            borderRadius: '8px',
            boxShadow: '0 2px 8px rgba(0,0,0,.06)',
            overflow: 'auto',
          }}
        >
          {children}
        </Content>
      </AntLayout>
    </AntLayout>
  );
}

export default Layout;