/**
 * 角色创建/编辑页面
 */

import { useState, useEffect } from 'react';
import { useNavigate, useLocation, useParams } from 'react-router-dom';
import {
  Card,
  Form,
  Input,
  Switch,
  Button,
  Space,
  message,
  Typography,
  Row,
  Col,
  Divider,
  Breadcrumb
} from 'antd';
import { ArrowLeftOutlined, HomeOutlined, UserOutlined } from '@ant-design/icons';
import { createRole, updateRole, getRoleById } from '../services/roleService';
import type { CreateRoleRequest } from '../services/roleService';

const { Title, Text } = Typography;

export function RoleCreate() {
  const navigate = useNavigate();
  const location = useLocation();
  const { id } = useParams<{ id: string }>();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [isEdit, setIsEdit] = useState(false);
  const [roleId, setRoleId] = useState<number | null>(null);

  // 从URL参数或location state获取初始数据
  useEffect(() => {
    const name = location.state?.name;
    console.log('RoleCreate useEffect - id:', id, 'pathname:', location.pathname);
    
    if (id) {
      // 编辑模式：从URL路径参数获取ID
      console.log('Setting edit mode with ID:', id);
      setIsEdit(true);
      setRoleId(parseInt(id));
      loadRole(parseInt(id));
    } else if (name) {
      // 创建模式：从简化模态框传入的角色名称
      form.setFieldsValue({ name });
    }
  }, [id, location, form]);

  // 加载角色数据（编辑模式）
  const loadRole = async (id: number) => {
    try {
      setLoading(true);
      console.log('Loading role with ID:', id);
      console.log('Role ID type:', typeof id);
      
      // 检查认证状态
      const token = localStorage.getItem('access_token');
      console.log('Current auth token exists:', !!token);
      console.log('Token value:', token ? token.substring(0, 20) + '...' : 'null');
      if (!token) {
        throw new Error('用户未认证，请重新登录');
      }
      
      // 验证角色ID是有效数字
      if (!id || isNaN(id) || id <= 0) {
        throw new Error(`无效的角色ID: ${id}`);
      }
      
      const role = await getRoleById(id);
      console.log('Loaded role:', role);
      
      if (!role) {
        throw new Error('无法获取角色数据');
      }
      
      // 安全地设置表单值，处理可能为空的字段
      const formValues = {
        name: role.name || '',
        description: role.description || '',
        home_page: role.home_page || '',
        default_route: role.default_route || '',
        restrict_to_domain: role.restrict_to_domain || '',
        desk_access: role.desk_access || false,
        require_two_factor: role.require_two_factor || false,
        is_enabled: role.is_enabled !== undefined ? role.is_enabled : true,
      };
      
      console.log('Setting form values:', formValues);
      form.setFieldsValue(formValues);
      console.log('Form values set successfully');
    } catch (error) {
      console.error('Failed to load role:', error);
      message.error(`加载角色数据失败: ${error instanceof Error ? error.message : '未知错误'}`);
    } finally {
      setLoading(false);
    }
  };

  // 保存角色
  const handleSave = async (values: CreateRoleRequest) => {
    try {
      setLoading(true);
      if (isEdit && roleId) {
        await updateRole(roleId, values);
        message.success('更新角色成功');
      } else {
        await createRole(values);
        message.success('创建角色成功');
      }
      navigate('/roles');
    } catch (error) {
      message.error(isEdit ? '更新角色失败' : '创建角色失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: '24px', minHeight: '100vh', backgroundColor: '#f5f5f5' }}>
      {/* 面包屑导航 */}
      <Breadcrumb style={{ marginBottom: '16px' }}>
        <Breadcrumb.Item>
          <HomeOutlined />
        </Breadcrumb.Item>
        <Breadcrumb.Item>
          <UserOutlined />
          <span>角色管理</span>
        </Breadcrumb.Item>
        <Breadcrumb.Item>
          {isEdit ? '编辑角色' : '创建角色'}
        </Breadcrumb.Item>
      </Breadcrumb>

      {/* 页面标题 */}
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: '24px' }}>
        <Button 
          type="text" 
          icon={<ArrowLeftOutlined />}
          onClick={() => navigate('/roles')}
          style={{ marginRight: '12px' }}
        >
          返回
        </Button>
        <Title level={2} style={{ margin: 0 }}>
          {isEdit ? '编辑角色' : '创建角色'}
        </Title>
      </div>

      {/* 表单内容 */}
      <Card loading={loading}>
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSave}
          initialValues={{
            is_enabled: true,
            desk_access: false,
            require_two_factor: false,
            is_custom: true,
          }}
        >
          {/* 基础信息 */}
          <Title level={4}>基础信息</Title>
          <Form.Item
            name="name"
            label="角色名称"
            rules={[
              { required: true, message: '请输入角色名称' },
              { min: 2, message: '角色名称至少2个字符' },
              { max: 50, message: '角色名称最多50个字符' },
            ]}
          >
            <Input 
              placeholder="请输入角色名称（如：销售经理、Sales Manager）" 
              style={{ fontSize: '14px' }}
            />
          </Form.Item>

          <Form.Item
            name="description"
            label="角色描述"
            rules={[{ max: 200, message: '描述最多200个字符' }]}
          >
            <Input.TextArea 
              rows={3}
              placeholder="请输入角色描述"
            />
          </Form.Item>

          <Divider />

          {/* 页面访问设置 */}
          <Title level={4}>页面访问设置</Title>
          <Row gutter={24}>
            <Col span={12}>
              <Form.Item
                name="home_page"
                label="默认主页"
              >
                <Input placeholder="/dashboard" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="default_route"
                label="默认路由"
              >
                <Input placeholder="/app" />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="restrict_to_domain"
            label="域限制"
          >
            <Input placeholder="限制角色访问的域名或组织" />
          </Form.Item>

          <Divider />

          {/* 权限设置 */}
          <Title level={4}>权限设置</Title>
          <Row gutter={24}>
            <Col span={8}>
              <Form.Item name="desk_access" valuePropName="checked">
                <Space direction="vertical">
                  <Space>
                    <Switch />
                    桌面访问权限
                  </Space>
                  <Text type="secondary" style={{ fontSize: '12px' }}>
                    允许访问后台管理界面
                  </Text>
                </Space>
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="require_two_factor" valuePropName="checked">
                <Space direction="vertical">
                  <Space>
                    <Switch />
                    强制双因子认证
                  </Space>
                  <Text type="secondary" style={{ fontSize: '12px' }}>
                    要求用户启用2FA
                  </Text>
                </Space>
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="is_custom" valuePropName="checked">
                <Space direction="vertical">
                  <Space>
                    <Switch disabled />
                    自定义角色
                  </Space>
                  <Text type="secondary" style={{ fontSize: '12px' }}>
                    系统自动识别角色类型
                  </Text>
                </Space>
              </Form.Item>
            </Col>
          </Row>

          <Divider />

          {/* 状态设置 */}
          <Title level={4}>状态设置</Title>
          <Form.Item name="is_enabled" valuePropName="checked">
            <Space direction="vertical">
              <Space>
                <Switch />
                启用角色
              </Space>
              <Text type="secondary" style={{ fontSize: '12px' }}>
                禁用后将从所有用户移除此角色
              </Text>
            </Space>
          </Form.Item>

          {/* 操作按钮 */}
          <div style={{ 
            marginTop: '48px', 
            paddingTop: '24px', 
            borderTop: '1px solid #f0f0f0',
            display: 'flex',
            justifyContent: 'flex-end'
          }}>
            <Button 
              type="primary" 
              size="large"
              htmlType="submit"
              loading={loading}
            >
              {isEdit ? '更新' : '创建'}
            </Button>
          </div>
        </Form>
      </Card>
    </div>
  );
}

export default RoleCreate;