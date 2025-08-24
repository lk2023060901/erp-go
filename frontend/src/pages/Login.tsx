/**
 * 登录页面
 */

import { useState } from 'react';
import { Card, Form, Input, Button, Checkbox, Alert, Divider, Space, Typography } from 'antd';
import { UserOutlined, LockOutlined, SafetyCertificateOutlined, EyeInvisibleOutlined, EyeTwoTone } from '@ant-design/icons';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import type { LoginRequest } from '../types/auth';
import './Login.css';

const { Title, Text } = Typography;

interface LoginFormData extends LoginRequest {
  remember: boolean;
  totp_code?: string;
}

export function Login() {
  const navigate = useNavigate();
  const location = useLocation();
  const { login, loading, error } = useAuth();
  
  const [form] = Form.useForm();
  const [, setShowTwoFactor] = useState(false);
  const [loginStep, setLoginStep] = useState<'credentials' | 'totp'>('credentials');

  // 获取重定向地址
  const from = (location.state as any)?.from?.pathname || '/dashboard';

  const handleSubmit = async (values: LoginFormData) => {
    try {
      const loginData: LoginRequest = {
        username: values.username,
        password: values.password,
        two_factor_code: values.totp_code,
      };

      await login(loginData);
      
      // 登录成功，跳转到目标页面
      navigate(from, { replace: true });
    } catch (err: any) {
      // 检查是否需要双重认证
      if (err.message?.includes('2FA') || err.code === 'TOTP_REQUIRED') {
        setShowTwoFactor(true);
        setLoginStep('totp');
      }
    }
  };

  const resetForm = () => {
    setShowTwoFactor(false);
    setLoginStep('credentials');
    form.resetFields(['totp_code']);
  };

  return (
    <div className="login-container">
      <div className="login-background" />
      
      <div className="login-content">
        <Card className="login-card" variant="borderless">
          <div className="login-header">
            <div className="login-logo">
              <SafetyCertificateOutlined style={{ fontSize: 32, color: '#1890ff' }} />
            </div>
            <Title level={2} className="login-title">
              ERP 管理系统
            </Title>
            <Text className="login-subtitle">
              {loginStep === 'credentials' ? '请登录您的账户' : '请输入验证码'}
            </Text>
          </div>

          {error && (
            <Alert
              message="登录失败"
              description={error}
              type="error"
              showIcon
              closable
              style={{ marginBottom: 24 }}
            />
          )}

          <Form
            form={form}
            name="login"
            size="large"
            onFinish={handleSubmit}
            autoComplete="off"
            layout="vertical"
          >
            {loginStep === 'credentials' && (
              <>
                <Form.Item
                  name="username"
                  label="用户名"
                  rules={[
                    { required: true, message: '请输入用户名' },
                    { min: 3, message: '用户名至少3个字符' },
                  ]}
                  initialValue="admin"
                >
                  <Input
                    prefix={<UserOutlined />}
                    placeholder="请输入用户名"
                    autoComplete="username"
                    defaultValue="admin"
                  />
                </Form.Item>

                <Form.Item
                  name="password"
                  label="密码"
                  rules={[
                    { required: true, message: '请输入密码' },
                    { min: 6, message: '密码至少6个字符' },
                  ]}
                  extra="默认管理员账户：admin / Admin123@"
                  initialValue="Admin123@"
                >
                  <Input.Password
                    prefix={<LockOutlined />}
                    placeholder="请输入密码"
                    autoComplete="current-password"
                    iconRender={(visible) => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)}
                    defaultValue="Admin123@"
                  />
                </Form.Item>

                <Form.Item>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Form.Item name="remember" valuePropName="checked" noStyle>
                      <Checkbox>记住我</Checkbox>
                    </Form.Item>
                    <Link to="/forgot-password" className="login-forgot">
                      忘记密码？
                    </Link>
                  </div>
                </Form.Item>
              </>
            )}

            {loginStep === 'totp' && (
              <>
                <Form.Item
                  name="totp_code"
                  label="验证码"
                  rules={[
                    { required: true, message: '请输入6位验证码' },
                    { len: 6, message: '验证码必须是6位数字' },
                    { pattern: /^\d{6}$/, message: '验证码只能包含数字' },
                  ]}
                >
                  <Input
                    placeholder="请输入6位验证码"
                    maxLength={6}
                    style={{ textAlign: 'center', fontSize: '18px', letterSpacing: '0.5em' }}
                  />
                </Form.Item>

                <Alert
                  message="双重认证"
                  description="请打开您的身份验证器应用(如Google Authenticator)获取6位验证码"
                  type="info"
                  showIcon
                  style={{ marginBottom: 16 }}
                />
              </>
            )}

            <Form.Item>
              <Space direction="vertical" style={{ width: '100%' }} size="middle">
                <Button
                  type="primary"
                  htmlType="submit"
                  loading={loading}
                  block
                  size="large"
                >
                  {loginStep === 'credentials' ? '登录' : '验证并登录'}
                </Button>

                {loginStep === 'totp' && (
                  <Button
                    type="default"
                    onClick={resetForm}
                    block
                    size="large"
                  >
                    返回登录
                  </Button>
                )}
              </Space>
            </Form.Item>
          </Form>

          <Divider plain>
            <Text type="secondary" style={{ fontSize: '12px' }}>
              其他登录方式
            </Text>
          </Divider>

          <div className="login-footer">
            <Space split={<Divider type="vertical" />}>
              <Link to="/register">注册账户</Link>
              <Link to="/help">帮助中心</Link>
              <Link to="/about">关于我们</Link>
            </Space>
          </div>
        </Card>

        <div className="login-copyright">
          <Text type="secondary">
            © 2024 ERP管理系统. All rights reserved.
          </Text>
        </div>
      </div>
    </div>
  );
}

export default Login;