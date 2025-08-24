import React, { useState } from 'react';
import { Form, Input, Switch, App } from 'antd';
import { useNavigate } from 'react-router-dom';
import { createUser } from '../services/userService';
import type { CreateUserRequest } from '../services/userService';
import './CreateUserModal.css';

interface CreateUserModalProps {
  visible: boolean;
  onClose: () => void;
  onSuccess?: () => void;
}

const CreateUserModal: React.FC<CreateUserModalProps> = ({ visible, onClose, onSuccess }) => {
  const { message } = App.useApp();
  const [form] = Form.useForm();
  const navigate = useNavigate();
  
  // 状态管理
  const [passwordModalVisible, setPasswordModalVisible] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [generatedPassword, setGeneratedPassword] = useState('');

  // 生成强密码函数
  const generateStrongPassword = (length: number = 12): string => {
    const uppercase = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
    const lowercase = 'abcdefghijklmnopqrstuvwxyz';
    const numbers = '0123456789';
    const symbols = '!@#$%^&*()_+-=[]{}|;:,.<>?';
    
    const allChars = uppercase + lowercase + numbers + symbols;
    let password = '';
    
    // 确保包含各种字符类型
    password += uppercase[Math.floor(Math.random() * uppercase.length)];
    password += lowercase[Math.floor(Math.random() * lowercase.length)];
    password += numbers[Math.floor(Math.random() * numbers.length)];
    password += symbols[Math.floor(Math.random() * symbols.length)];
    
    // 填充剩余长度
    for (let i = 4; i < length; i++) {
      password += allChars[Math.floor(Math.random() * allChars.length)];
    }
    
    // 打乱字符顺序
    return password.split('').sort(() => 0.5 - Math.random()).join('');
  };

  // 生成密码
  const handleGeneratePassword = () => {
    const newPassword = generateStrongPassword();
    setGeneratedPassword(newPassword);
    setPasswordModalVisible(true);
  };

  // 复制密码
  const handleCopyPassword = async () => {
    try {
      await navigator.clipboard.writeText(generatedPassword);
      message.success('已复制');
      
      // 自动填入表单
      form.setFieldsValue({ password: generatedPassword });
      
      // 延时关闭模态框
      setTimeout(() => {
        setPasswordModalVisible(false);
      }, 1500);
    } catch (error) {
      message.error('复制失败');
    }
  };

  // 创建用户
  const handleCreateUser = async (values: any) => {
    try {
      // 构造请求数据，使用新的字段结构
      const createData: CreateUserRequest = {
        username: values.username,
        email: values.email,
        first_name: '',  // 空字符串作为默认值
        last_name: '',   // 空字符串作为默认值
        password: values.password,
        is_active: values.is_active !== false, // 默认启用
      };
      
      await createUser(createData);
      message.success('创建用户成功');
      form.resetFields();
      onClose();
      onSuccess?.();
    } catch (error: any) {
      const errorMessage = error.message || '创建用户失败';
      message.error(errorMessage);
    }
  };

  const handleClose = () => {
    form.resetFields();
    onClose();
  };

  if (!visible) return null;

  return (
    <>
      {/* 创建用户模态框 */}
      <div className="user-modal-overlay">
        <div className="user-modal">
          <div className="user-modal-header">
            <div className="user-modal-title">创建用户</div>
            <button className="user-close-btn" onClick={handleClose}>
              ×
            </button>
          </div>

          <div className="user-modal-body">
            <Form
              form={form}
              layout="vertical"
              onFinish={handleCreateUser}
              initialValues={{ is_active: true }}
            >
              <div className="user-form-group">
                <label className="user-form-label">用户名<span className="user-required">*</span></label>
                <Form.Item
                  name="username"
                  rules={[
                    { required: true, message: '请输入用户名' },
                    { min: 3, message: '用户名至少3个字符' },
                    { max: 20, message: '用户名最多20个字符' },
                    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线' },
                  ]}
                >
                  <Input className="user-form-input user-wide-field" placeholder="请输入用户名" />
                </Form.Item>
              </div>

              <div className="user-form-group">
                <label className="user-form-label">邮件地址<span className="user-required">*</span></label>
                <Form.Item
                  name="email"
                  rules={[
                    { required: true, message: '请输入邮箱' },
                    { type: 'email', message: '请输入有效的邮箱地址' },
                  ]}
                >
                  <Input className="user-form-input user-wide-field" placeholder="请输入邮件地址" />
                </Form.Item>
              </div>


              <div className="user-form-group">
                <label className="user-form-label">密码<span className="user-required">*</span></label>
                <div className="user-password-wrapper">
                  <Form.Item
                    name="password"
                    rules={[
                      { required: true, message: '请输入密码' },
                      { min: 6, message: '密码至少6个字符' },
                    ]}
                  >
                    <Input 
                      className="user-form-input user-wide-field"
                      type={showPassword ? 'text' : 'password'}
                      placeholder=""
                    />
                  </Form.Item>
                  <div className="user-password-icons">
                    <svg 
                      className="user-password-icon" 
                      onClick={() => setShowPassword(!showPassword)}
                      viewBox="0 0 24 24" 
                      fill="currentColor"
                    >
                      {showPassword ? (
                        <path d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21m-12-12h.01M12 9a3 3 0 013 3m-3-3a3 3 0 00-3 3m3-3V6a9 9 0 019 9v1.5m-9-10.5a9 9 0 00-9 9v1.5" />
                      ) : (
                        <path d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zM12 17c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z"/>
                      )}
                    </svg>
                    <svg 
                      className="user-password-icon" 
                      onClick={handleGeneratePassword}
                      viewBox="0 0 24 24" 
                      fill="currentColor"
                    >
                      <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                    </svg>
                  </div>
                </div>
              </div>

              <div className="user-form-group">
                <div style={{ display: 'flex', alignItems: 'center', paddingTop: '12px' }}>
                  <span style={{ marginRight: '8px', fontSize: '14px', color: '#333' }}>账户状态：</span>
                  <Form.Item
                    name="is_active"
                    valuePropName="checked"
                    style={{ margin: 0 }}
                  >
                    <Switch
                      checkedChildren="启用"
                      unCheckedChildren="禁用"
                      defaultChecked={true}
                    />
                  </Form.Item>
                </div>
              </div>
            </Form>
          </div>

          <div className="user-modal-footer">
            <button 
              type="button" 
              className="user-btn user-btn-default"
              onClick={() => {
                onClose();
                navigate('/users/create');
              }}
            >
              编辑完整信息
            </button>
            <button 
              type="button" 
              className="user-btn user-btn-primary"
              onClick={() => form.submit()}
            >
              创建
            </button>
          </div>
        </div>
      </div>

      {/* 密码生成模态框 */}
      {passwordModalVisible && (
        <div className="user-password-modal-overlay">
          <div className="user-password-modal">
            <div className="user-password-modal-header">
              <div className="user-password-modal-title">生成密码</div>
              <button 
                className="user-close-btn" 
                onClick={() => setPasswordModalVisible(false)}
              >
                ×
              </button>
            </div>

            <div className="user-password-modal-body">
              <div className="user-password-tip">
                <p>系统为您生成了一个安全的密码，请复制使用：</p>
              </div>
              <div className="user-generated-password">
                <div className="user-password-text">{generatedPassword}</div>
                <button 
                  type="button" 
                  className="user-copy-btn" 
                  onClick={handleCopyPassword}
                >
                  复制密码
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </>
  );
};

export default CreateUserModal;