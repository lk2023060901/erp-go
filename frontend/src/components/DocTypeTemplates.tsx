import React from 'react';
import { Card, Button, Row, Col, Tag, Typography } from 'antd';
import { 
  UserOutlined, 
  ShoppingCartOutlined, 
  DatabaseOutlined,
  FileTextOutlined,
  SettingOutlined,
  TeamOutlined
} from '@ant-design/icons';

const { Title, Text } = Typography;

export interface DocTypeTemplate {
  id: string;
  name: string;
  label: string;
  description: string;
  module: string;
  icon: React.ReactNode;
  type: 'normal' | 'submittable' | 'child';
  config: {
    is_submittable: boolean;
    is_child_table: boolean;
    has_workflow: boolean;
    track_changes: boolean;
    applies_to_all_users: boolean;
  };
}

const templates: DocTypeTemplate[] = [
  {
    id: 'user',
    name: 'User',
    label: '用户信息',
    description: '存储系统用户的基本信息，如姓名、邮箱、电话等',
    module: 'Core',
    icon: <UserOutlined style={{ fontSize: '24px', color: '#1890ff' }} />,
    type: 'normal',
    config: {
      is_submittable: false,
      is_child_table: false,
      has_workflow: false,
      track_changes: true,
      applies_to_all_users: true,
    }
  },
  {
    id: 'order',
    name: 'Order',
    label: '订单记录',
    description: '客户订单信息，包含订单状态、金额、商品等信息',
    module: 'Selling',
    icon: <ShoppingCartOutlined style={{ fontSize: '24px', color: '#52c41a' }} />,
    type: 'submittable',
    config: {
      is_submittable: true,
      is_child_table: false,
      has_workflow: true,
      track_changes: true,
      applies_to_all_users: false,
    }
  },
  {
    id: 'product',
    name: 'Product',
    label: '产品目录',
    description: '产品信息管理，包括产品名称、价格、库存等',
    module: 'Stock',
    icon: <DatabaseOutlined style={{ fontSize: '24px', color: '#fa8c16' }} />,
    type: 'normal',
    config: {
      is_submittable: false,
      is_child_table: false,
      has_workflow: false,
      track_changes: true,
      applies_to_all_users: false,
    }
  },
  {
    id: 'document',
    name: 'Document',
    label: '文档管理',
    description: '各类业务文档，如合同、报告、说明书等',
    module: 'Custom',
    icon: <FileTextOutlined style={{ fontSize: '24px', color: '#722ed1' }} />,
    type: 'submittable',
    config: {
      is_submittable: true,
      is_child_table: false,
      has_workflow: true,
      track_changes: true,
      applies_to_all_users: false,
    }
  },
  {
    id: 'config',
    name: 'SystemConfig',
    label: '系统配置',
    description: '系统参数和配置信息管理',
    module: 'Setup',
    icon: <SettingOutlined style={{ fontSize: '24px', color: '#13c2c2' }} />,
    type: 'normal',
    config: {
      is_submittable: false,
      is_child_table: false,
      has_workflow: false,
      track_changes: true,
      applies_to_all_users: true,
    }
  },
  {
    id: 'team',
    name: 'Team',
    label: '团队信息',
    description: '团队或部门信息，包含成员、职责等',
    module: 'Core',
    icon: <TeamOutlined style={{ fontSize: '24px', color: '#eb2f96' }} />,
    type: 'normal',
    config: {
      is_submittable: false,
      is_child_table: false,
      has_workflow: false,
      track_changes: true,
      applies_to_all_users: false,
    }
  }
];

interface DocTypeTemplatesProps {
  onSelectTemplate: (template: DocTypeTemplate) => void;
  onCreateBlank: () => void;
}

const DocTypeTemplates: React.FC<DocTypeTemplatesProps> = ({ onSelectTemplate, onCreateBlank }) => {
  const getTypeTag = (type: string) => {
    switch (type) {
      case 'submittable':
        return <Tag color="green">审批文档</Tag>;
      case 'child':
        return <Tag color="orange">子表</Tag>;
      default:
        return <Tag color="blue">普通文档</Tag>;
    }
  };

  return (
    <div className="doctype-templates">
      <Title level={4} style={{ marginBottom: 16 }}>
        选择文档类型模板
      </Title>
      <Text type="secondary" style={{ marginBottom: 24, display: 'block' }}>
        选择一个预设模板快速开始，或创建空白文档类型自定义配置
      </Text>

      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        {templates.map(template => (
          <Col span={8} key={template.id}>
            <Card 
              hoverable
              className="template-card"
              onClick={() => onSelectTemplate(template)}
              style={{ height: '140px' }}
            >
              <div style={{ display: 'flex', alignItems: 'flex-start', height: '100%' }}>
                <div style={{ marginRight: 12 }}>
                  {template.icon}
                </div>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div style={{ marginBottom: 4 }}>
                    <Text strong style={{ fontSize: '14px' }}>{template.label}</Text>
                    <div>{getTypeTag(template.type)}</div>
                  </div>
                  <Text 
                    type="secondary" 
                    style={{ 
                      fontSize: '12px',
                      lineHeight: '16px',
                      display: '-webkit-box',
                      WebkitLineClamp: 3,
                      WebkitBoxOrient: 'vertical',
                      overflow: 'hidden'
                    }}
                  >
                    {template.description}
                  </Text>
                </div>
              </div>
            </Card>
          </Col>
        ))}
      </Row>

      <div style={{ textAlign: 'center', marginTop: 24 }}>
        <Text type="secondary">或者</Text>
        <br />
        <Button 
          type="default" 
          onClick={onCreateBlank}
          style={{ marginTop: 8 }}
        >
          创建空白文档类型
        </Button>
      </div>
    </div>
  );
};

export default DocTypeTemplates;
export { templates };