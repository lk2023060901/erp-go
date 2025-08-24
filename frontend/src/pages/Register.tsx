import React, { useState } from 'react';
import { Form, Input, Button, message, Row, Col, Card, Typography, Select, Space } from 'antd';
import { UserOutlined, MailOutlined, LockOutlined, PhoneOutlined, UserAddOutlined } from '@ant-design/icons';
import { Link, useNavigate } from 'react-router-dom';
import { authService } from '../services/auth';
import './Register.css';

const { Title, Text } = Typography;
const { Option } = Select;

interface RegisterForm {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
  firstName: string;
  lastName: string;
  phone?: string;
  gender?: string;
}

const Register: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);

  const onFinish = async (values: RegisterForm) => {
    setLoading(true);
    try {
      await authService.register({
        username: values.username,
        email: values.email,
        password: values.password,
        confirm_password: values.confirmPassword,
        first_name: values.firstName,
        last_name: values.lastName,
        phone: values.phone,
        gender: values.gender,
      });
      
      message.success('注册成功！请检查您的邮箱完成验证后登录');
      navigate('/login');
    } catch (error: any) {
      console.error('Registration failed:', error);
      message.error(error.message || '注册失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const validatePassword = (_: any, value: string) => {
    if (!value) {
      return Promise.reject(new Error('请输入密码'));
    }
    
    const rules = [
      { test: /.{8,}/, message: '密码长度至少8位' },
      { test: /[A-Z]/, message: '密码必须包含大写字母' },
      { test: /[a-z]/, message: '密码必须包含小写字母' },
      { test: /[0-9]/, message: '密码必须包含数字' },
      { test: /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>?~]/, message: '密码必须包含特殊字符' },
    ];

    for (const rule of rules) {
      if (!rule.test.test(value)) {
        return Promise.reject(new Error(rule.message));
      }
    }

    return Promise.resolve();
  };

  const validateConfirmPassword = (_: any, value: string) => {
    if (!value) {
      return Promise.reject(new Error('请确认密码'));
    }
    if (value !== form.getFieldValue('password')) {
      return Promise.reject(new Error('两次输入的密码不一致'));
    }
    return Promise.resolve();
  };

  return (
    <div className="register-container">
      <div className="register-background"></div>
      <div className="register-content">
        <Card className="register-card">
          <div className="register-header">
            <UserAddOutlined className="register-icon" />
            <Title level={2} className="register-title">
              用户注册
            </Title>
            <Text type="secondary" className="register-subtitle">
              创建您的ERP系统账户
            </Text>
          </div>

          <Form
            form={form}
            name="register"
            onFinish={onFinish}
            layout="vertical"
            requiredMark={false}
            size="large"
          >
            <Row gutter={16}>
              <Col span={12}>
                <Form.Item
                  name="firstName"
                  label="姓"
                  rules={[
                    { required: true, message: '请输入姓氏' },
                    { min: 1, max: 50, message: '姓氏长度为1-50个字符' },
                  ]}
                >
                  <Input 
                    prefix={<UserOutlined />} 
                    placeholder="请输入姓氏" 
                  />
                </Form.Item>
              </Col>
              <Col span={12}>
                <Form.Item
                  name="lastName"
                  label="名"
                  rules={[
                    { required: true, message: '请输入名字' },
                    { min: 1, max: 50, message: '名字长度为1-50个字符' },
                  ]}
                >
                  <Input 
                    prefix={<UserOutlined />} 
                    placeholder="请输入名字" 
                  />
                </Form.Item>
              </Col>
            </Row>

            <Form.Item
              name="username"
              label="用户名"
              rules={[
                { required: true, message: '请输入用户名' },
                { min: 3, max: 32, message: '用户名长度为3-32个字符' },
                { 
                  pattern: /^[a-zA-Z0-9_-]+$/, 
                  message: '用户名只能包含字母、数字、下划线和连字符' 
                },
                { 
                  pattern: /^[a-zA-Z_-]/, 
                  message: '用户名不能以数字开头' 
                },
              ]}
            >
              <Input 
                prefix={<UserOutlined />} 
                placeholder="请输入用户名" 
              />
            </Form.Item>

            <Form.Item
              name="email"
              label="邮箱"
              rules={[
                { required: true, message: '请输入邮箱地址' },
                { type: 'email', message: '请输入有效的邮箱地址' },
              ]}
            >
              <Input 
                prefix={<MailOutlined />} 
                placeholder="请输入邮箱地址" 
              />
            </Form.Item>

            <Row gutter={16}>
              <Col span={12}>
                <Form.Item
                  name="phone"
                  label="手机号"
                  rules={[
                    { 
                      pattern: /^1[3-9]\d{9}$/, 
                      message: '请输入有效的手机号码' 
                    },
                  ]}
                >
                  <Input 
                    prefix={<PhoneOutlined />} 
                    placeholder="手机号（可选）" 
                  />
                </Form.Item>
              </Col>
              <Col span={12}>
                <Form.Item
                  name="gender"
                  label="性别"
                >
                  <Select placeholder="请选择性别（可选）">
                    <Option value="MALE">男</Option>
                    <Option value="FEMALE">女</Option>
                    <Option value="OTHER">其他</Option>
                  </Select>
                </Form.Item>
              </Col>
            </Row>

            <Form.Item
              name="password"
              label="密码"
              rules={[{ validator: validatePassword }]}
            >
              <Input.Password 
                prefix={<LockOutlined />} 
                placeholder="请输入密码"
              />
            </Form.Item>

            <Form.Item
              name="confirmPassword"
              label="确认密码"
              rules={[{ validator: validateConfirmPassword }]}
            >
              <Input.Password 
                prefix={<LockOutlined />} 
                placeholder="请再次输入密码"
              />
            </Form.Item>

            <div className="password-requirements">
              <Text type="secondary" style={{ fontSize: '12px' }}>
                密码要求：至少8位，包含大小写字母、数字和特殊字符
              </Text>
            </div>

            <Form.Item>
              <Button 
                type="primary" 
                htmlType="submit" 
                className="register-button"
                loading={loading}
                block
              >
                立即注册
              </Button>
            </Form.Item>

            <div className="register-footer">
              <Space>
                <Text type="secondary">已有账户？</Text>
                <Link to="/login" className="login-link">
                  立即登录
                </Link>
              </Space>
            </div>
          </Form>
        </Card>
      </div>
    </div>
  );
};

export default Register;