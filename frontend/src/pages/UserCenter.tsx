import React, { useState } from 'react';
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
import CreateUserModal from '../components/CreateUserModal';

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
      justifyContent: 'center',
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
    <span style={{ fontSize: '14px', opacity: 0.8 }}>{icon}</span>
    <span style={{ fontSize: '12px', fontWeight: 500, color: '#1a1a1a' }}>{title}</span>
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
    <div style={{ display: 'flex', gap: '4px', flexWrap: 'nowrap' }}>
      {children}
    </div>
  </Card>
);

interface FunctionCardProps {
  icon: React.ReactNode;
  title: string;
  description: string;
  metrics: Array<{
    number: string;
    label: string;
  }>;
  onClick: () => void;
}

const FunctionCard: React.FC<FunctionCardProps> = ({ 
  icon, 
  title, 
  description, 
  metrics, 
  onClick 
}) => (
  <Card
    hoverable
    onClick={onClick}
    style={{
      borderRadius: '16px',
      border: '1px solid #e5e7eb',
      boxShadow: '0 1px 3px rgba(0, 0, 0, 0.05)',
      cursor: 'pointer',
      transition: 'all 0.2s cubic-bezier(0.4, 0, 0.2, 1)',
    }}
    bodyStyle={{ padding: '24px' }}
    onMouseEnter={(e) => {
      e.currentTarget.style.transform = 'translateY(-2px)';
      e.currentTarget.style.borderColor = '#6366f1';
      e.currentTarget.style.boxShadow = '0 10px 25px rgba(0, 0, 0, 0.1)';
    }}
    onMouseLeave={(e) => {
      e.currentTarget.style.transform = 'translateY(0)';
      e.currentTarget.style.borderColor = '#e5e7eb';
      e.currentTarget.style.boxShadow = '0 1px 3px rgba(0, 0, 0, 0.05)';
    }}
  >
    <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginBottom: '20px' }}>
      <div style={{
        width: '48px',
        height: '48px',
        borderRadius: '12px',
        background: 'linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}>
        <span style={{ fontSize: '24px', color: '#ffffff' }}>{icon}</span>
      </div>
      <Title level={4} style={{ margin: 0, fontSize: '18px', fontWeight: 600, color: '#1a1a1a' }}>
        {title}
      </Title>
    </div>
    
    <div style={{ 
      fontSize: '14px', 
      color: '#6b7280', 
      lineHeight: 1.6, 
      display: 'block', 
      marginBottom: '20px' 
    }}>
      {description}
    </div>
    
    <div style={{ display: 'flex', gap: '20px', fontSize: '13px' }}>
      {metrics.map((metric, index) => (
        <div key={index} style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
          <div style={{ 
            fontSize: '20px', 
            fontWeight: 700, 
            color: '#1a1a1a', 
            lineHeight: 1 
          }}>
            {metric.number}
          </div>
          <div style={{ 
            color: '#6b7280', 
            textTransform: 'uppercase', 
            letterSpacing: '0.5px', 
            fontSize: '11px' 
          }}>
            {metric.label}
          </div>
        </div>
      ))}
    </div>
  </Card>
);

const UserCenter: React.FC = () => {
  const navigate = useNavigate();
  const [createModalVisible, setCreateModalVisible] = useState(false);

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
        
        <Row gutter={[16, 16]}>
          <Col xs={24} lg={12}>
            <CategoryCard icon={<UserOutlined />} title="用户管理">
              <QuickAction
                icon={<UnorderedListOutlined />}
                title="用户列表"
                onClick={() => handleNavigate('/users/list')}
              />
              <QuickAction
                icon={<PlusOutlined />}
                title="创建用户"
                onClick={() => setCreateModalVisible(true)}
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
          
          <Col xs={24} lg={12}>
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
                icon={<ImportOutlined />}
                title="批量导入"
                onClick={() => handleNavigate('/roles/import')}
              />
              <div style={{ width: '100%', minHeight: '36px' }} />
            </CategoryCard>
          </Col>
        </Row>
      </div>
      
      {/* 分割线 */}
      <Divider style={{ margin: '24px 0', backgroundColor: '#e5e7eb' }} />
      
      {/* 核心功能模块 */}
      <div>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '24px' }}>
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
            核心功能
          </Title>
        </div>
        
        <Row gutter={[20, 20]}>
          <Col xs={24} lg={8}>
            <FunctionCard
              icon={<UserOutlined />}
              title="用户管理"
              description="全面的用户生命周期管理，支持批量操作、状态监控和用户分析统计等高级功能。"
              metrics={[
                { number: '142', label: '活跃用户' },
                { number: '14', label: '禁用用户' },
                { number: '8', label: '待审核' }
              ]}
              onClick={() => handleNavigate('/users/list')}
            />
          </Col>
          
          <Col xs={24} lg={8}>
            <FunctionCard
              icon={<TeamOutlined />}
              title="角色管理"
              description="基于角色的访问控制系统，支持层级权限管理、权限继承和动态角色分配。"
              metrics={[
                { number: '5', label: '管理员角色' },
                { number: '7', label: '普通角色' },
                { number: '12', label: '总计' }
              ]}
              onClick={() => handleNavigate('/roles')}
            />
          </Col>
          
          <Col xs={24} lg={8}>
            <FunctionCard
              icon={<SafetyCertificateOutlined />}
              title="权限管理"
              description="精细化权限控制，提供API级别的安全管理、资源访问控制和完整的审计追踪。"
              metrics={[
                { number: '25', label: '功能权限' },
                { number: '64', label: 'API权限' },
                { number: '89', label: '总计' }
              ]}
              onClick={() => handleNavigate('/permissions')}
            />
          </Col>
        </Row>
      </div>

      {/* 创建用户模态框 */}
      <CreateUserModal 
        visible={createModalVisible}
        onClose={() => setCreateModalVisible(false)}
        onSuccess={() => {
          // 可以在这里添加成功后的操作，比如刷新数据
        }}
      />
    </div>
  );
};

export default UserCenter;