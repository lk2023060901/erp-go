import React from 'react';
import { Card, Row, Col, Button, Typography, Divider } from 'antd';
import {
  UserOutlined,
  TeamOutlined,
  PlusOutlined,
  UnorderedListOutlined,
  EditOutlined,
  ImportOutlined,
  SafetyCertificateOutlined
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';

const { Title } = Typography;

interface QuickActionProps {
  icon: React.ReactNode;
  title: string;
  onClick: () => void;
}

const QuickAction: React.FC<QuickActionProps> = ({ icon, title, onClick }) => (
  <Button
    type="text"
    className="quick-action-btn"
    onClick={onClick}
    style={{
      height: 'auto',
      padding: '8px 6px',
      display: 'flex',
      alignItems: 'center',
      gap: '4px',
      border: '1px solid #e5e7eb',
      borderRadius: '6px',
      backgroundColor: '#ffffff',
      transition: 'all 0.2s cubic-bezier(0.4, 0, 0.2, 1)',
      boxShadow: '0 1px 2px rgba(0, 0, 0, 0.05)',
      width: '100%',
      minHeight: '36px',
    }}
    onMouseEnter={(e) => {
      e.currentTarget.style.transform = 'translateY(-1px)';
      e.currentTarget.style.borderColor = '#6366f1';
      e.currentTarget.style.boxShadow = '0 4px 12px rgba(0, 0, 0, 0.1)';
      e.currentTarget.style.backgroundColor = '#f8fafc';
    }}
    onMouseLeave={(e) => {
      e.currentTarget.style.transform = 'translateY(0)';
      e.currentTarget.style.borderColor = '#e5e7eb';
      e.currentTarget.style.boxShadow = '0 1px 2px rgba(0, 0, 0, 0.05)';
      e.currentTarget.style.backgroundColor = '#ffffff';
    }}
  >
    <span style={{ fontSize: '12px', opacity: 0.8 }}>{icon}</span>
    <span style={{ fontSize: '11px', fontWeight: 500, color: '#1a1a1a' }}>{title}</span>
  </Button>
);

interface CategoryCardProps {
  icon: React.ReactNode;
  title: string;
  children: React.ReactNode;
}

const CategoryCard: React.FC<CategoryCardProps> = ({ icon, title, children }) => (
  <Card
    style={{
      borderRadius: '16px',
      border: '1px solid #e5e7eb',
      boxShadow: '0 1px 3px rgba(0, 0, 0, 0.05)',
    }}
    bodyStyle={{ padding: '20px' }}
  >
    <div style={{ 
      display: 'flex', 
      alignItems: 'center', 
      gap: '10px', 
      marginBottom: '16px',
      paddingBottom: '12px',
      borderBottom: '1px solid #f3f4f6'
    }}>
      <span style={{ fontSize: '20px', opacity: 0.7 }}>{icon}</span>
      <Title level={5} style={{ margin: 0, fontSize: '16px', fontWeight: 600, color: '#1a1a1a' }}>
        {title}
      </Title>
    </div>
    <div style={{ display: 'flex', gap: '6px', flexWrap: 'nowrap' }}>
      {children}
    </div>
  </Card>
);

const UserCenterNew: React.FC = () => {
  const navigate = useNavigate();

  const handleNavigate = (path: string) => {
    navigate(path);
  };

  return (
    <div style={{ margin: '16px 0 32px 0' }}>
      {/* 快速操作区域 */}
      <div style={{ marginBottom: '32px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '20px' }}>
          <div style={{
            width: '6px',
            height: '6px',
            borderRadius: '50%',
            backgroundColor: '#6366f1',
          }} />
          <Title level={3} style={{ 
            margin: 0, 
            fontSize: '18px', 
            fontWeight: 600, 
            color: '#1a1a1a',
            letterSpacing: '-0.01em'
          }}>
            快速操作
          </Title>
        </div>
        
        <Row gutter={[0, 16]}>
          <Col span={24}>
            <CategoryCard icon={<UserOutlined />} title="用户管理">
              <QuickAction
                icon={<UnorderedListOutlined />}
                title="用户列表"
                onClick={() => handleNavigate('/users/list')}
              />
              <QuickAction
                icon={<PlusOutlined />}
                title="创建用户"
                onClick={() => handleNavigate('/users/create')}
              />
              <QuickAction
                icon={<EditOutlined />}
                title="用户设置"
                onClick={() => handleNavigate('/users/profile')}
              />
              <QuickAction
                icon={<ImportOutlined />}
                title="批量导入"
                onClick={() => handleNavigate('/users/import')}
              />
            </CategoryCard>
          </Col>
          
          <Col span={24}>
            <CategoryCard icon={<TeamOutlined />} title="角色管理">
              <QuickAction
                icon={<UnorderedListOutlined />}
                title="角色列表"
                onClick={() => handleNavigate('/roles/list')}
              />
              <QuickAction
                icon={<PlusOutlined />}
                title="创建角色"
                onClick={() => handleNavigate('/roles/create')}
              />
              <QuickAction
                icon={<EditOutlined />}
                title="编辑角色"
                onClick={() => handleNavigate('/roles/edit')}
              />
              <QuickAction
                icon={<ImportOutlined />}
                title="批量导入"
                onClick={() => handleNavigate('/roles/import')}
              />
            </CategoryCard>
          </Col>
        </Row>
      </div>
      
      {/* 分割线 */}
      <Divider />
      
      {/* 核心功能模块 */}
      <div>
        <Title level={4} style={{ marginBottom: '24px', color: '#1a1a1a' }}>
          核心功能
        </Title>
        
        <Row gutter={[24, 24]}>
          <Col xs={24} lg={8}>
            <Card
              hoverable
              onClick={() => handleNavigate('/users/list')}
              style={{ cursor: 'pointer', height: '200px' }}
              bodyStyle={{ padding: '24px' }}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginBottom: '16px' }}>
                <div style={{
                  width: '48px',
                  height: '48px',
                  borderRadius: '12px',
                  background: 'linear-gradient(135deg, #1890ff 0%, #36cfc9 100%)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}>
                  <UserOutlined style={{ fontSize: '24px', color: '#ffffff' }} />
                </div>
                <Title level={4} style={{ margin: 0 }}>用户管理</Title>
              </div>
              
              <p style={{ color: '#666', marginBottom: '16px' }}>
                全面的用户生命周期管理，支持批量操作、状态监控和用户分析统计等高级功能。
              </p>
              
              <div style={{ display: 'flex', gap: '20px' }}>
                <div>
                  <div style={{ fontSize: '20px', fontWeight: 'bold', color: '#1a1a1a' }}>142</div>
                  <div style={{ fontSize: '12px', color: '#666' }}>活跃用户</div>
                </div>
                <div>
                  <div style={{ fontSize: '20px', fontWeight: 'bold', color: '#1a1a1a' }}>14</div>
                  <div style={{ fontSize: '12px', color: '#666' }}>禁用用户</div>
                </div>
                <div>
                  <div style={{ fontSize: '20px', fontWeight: 'bold', color: '#1a1a1a' }}>8</div>
                  <div style={{ fontSize: '12px', color: '#666' }}>待审核</div>
                </div>
              </div>
            </Card>
          </Col>
          
          <Col xs={24} lg={8}>
            <Card
              hoverable
              onClick={() => handleNavigate('/roles')}
              style={{ cursor: 'pointer', height: '200px' }}
              bodyStyle={{ padding: '24px' }}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginBottom: '16px' }}>
                <div style={{
                  width: '48px',
                  height: '48px',
                  borderRadius: '12px',
                  background: 'linear-gradient(135deg, #52c41a 0%, #73d13d 100%)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}>
                  <TeamOutlined style={{ fontSize: '24px', color: '#ffffff' }} />
                </div>
                <Title level={4} style={{ margin: 0 }}>角色管理</Title>
              </div>
              
              <p style={{ color: '#666', marginBottom: '16px' }}>
                基于角色的访问控制系统，支持层级权限管理、权限继承和动态角色分配。
              </p>
              
              <div style={{ display: 'flex', gap: '20px' }}>
                <div>
                  <div style={{ fontSize: '20px', fontWeight: 'bold', color: '#1a1a1a' }}>5</div>
                  <div style={{ fontSize: '12px', color: '#666' }}>管理员角色</div>
                </div>
                <div>
                  <div style={{ fontSize: '20px', fontWeight: 'bold', color: '#1a1a1a' }}>7</div>
                  <div style={{ fontSize: '12px', color: '#666' }}>普通角色</div>
                </div>
                <div>
                  <div style={{ fontSize: '20px', fontWeight: 'bold', color: '#1a1a1a' }}>12</div>
                  <div style={{ fontSize: '12px', color: '#666' }}>总计</div>
                </div>
              </div>
            </Card>
          </Col>
          
          <Col xs={24} lg={8}>
            <Card
              hoverable
              onClick={() => handleNavigate('/permissions')}
              style={{ cursor: 'pointer', height: '200px' }}
              bodyStyle={{ padding: '24px' }}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginBottom: '16px' }}>
                <div style={{
                  width: '48px',
                  height: '48px',
                  borderRadius: '12px',
                  background: 'linear-gradient(135deg, #722ed1 0%, #9254de 100%)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}>
                  <SafetyCertificateOutlined style={{ fontSize: '24px', color: '#ffffff' }} />
                </div>
                <Title level={4} style={{ margin: 0 }}>权限管理</Title>
              </div>
              
              <p style={{ color: '#666', marginBottom: '16px' }}>
                精细化权限控制，提供API级别的安全管理、资源访问控制和完整的审计追踪。
              </p>
              
              <div style={{ display: 'flex', gap: '20px' }}>
                <div>
                  <div style={{ fontSize: '20px', fontWeight: 'bold', color: '#1a1a1a' }}>25</div>
                  <div style={{ fontSize: '12px', color: '#666' }}>功能权限</div>
                </div>
                <div>
                  <div style={{ fontSize: '20px', fontWeight: 'bold', color: '#1a1a1a' }}>64</div>
                  <div style={{ fontSize: '12px', color: '#666' }}>API权限</div>
                </div>
                <div>
                  <div style={{ fontSize: '20px', fontWeight: 'bold', color: '#1a1a1a' }}>89</div>
                  <div style={{ fontSize: '12px', color: '#666' }}>总计</div>
                </div>
              </div>
            </Card>
          </Col>
        </Row>
      </div>
    </div>
  );
};

export default UserCenterNew;