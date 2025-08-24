import React from 'react';
import { Form, Input, Switch, Tooltip } from 'antd';
import { EyeInvisibleOutlined } from '@ant-design/icons';
import { TabProps } from '../types/auth';
import { getFieldPermissions } from './UserEditForm';
import AvatarUpload from './AvatarUpload';

interface BasicInfoTabProps extends TabProps {
  // form instance is not used in this component
}

const BasicInfoTab: React.FC<BasicInfoTabProps> = ({ 
  mode, 
  initialData, 
  onDataChange
}) => {
  // 状态管理移动到父组件
  // const isActive = initialData?.is_active !== false;

  // 处理数据变化
  const handleDataChange = (field: string, value: any) => {
    const newData = { [field]: value };
    onDataChange?.(newData);
  };

  // 处理头像变化
  const handleAvatarChange = (file: File | null) => {
    handleDataChange('avatar', file);
  };

  // 获取字段权限
  const emailPermissions = getFieldPermissions(mode, 'email');
  const usernamePermissions = getFieldPermissions(mode, 'username');
  const firstNamePermissions = getFieldPermissions(mode, 'first_name');
  const lastNamePermissions = getFieldPermissions(mode, 'last_name');

  return (
    <div>
      {/* 头像上传区域 */}
      <div className="user-edit-avatar-section">
        <AvatarUpload
          value={(initialData as any)?.avatar_url || ''}
          onChange={handleAvatarChange}
          disabled={false}
        />
        <div className="user-edit-avatar-info">
          <h4>用户头像</h4>
          <p>点击头像上传新图片，建议尺寸：200x200px</p>
        </div>
      </div>

      {/* 基本信息表单 */}
      <div className="user-edit-form-row">
        <div className="user-edit-form-group">
          <label className="user-edit-form-label">
            用户名
            {usernamePermissions.required && <span className="user-edit-required">*</span>}
          </label>
          <Form.Item
            name="username"
            rules={usernamePermissions.required ? [
              { required: true, message: '请输入用户名' },
              { min: 3, message: '用户名至少3个字符' },
              { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线' }
            ] : []}
          >
            {usernamePermissions.readonly ? (
              <Tooltip title="此字段不可修改">
                <div style={{ position: 'relative' }}>
                  <Input
                    className="user-edit-form-input readonly"
                    disabled
                    value={initialData?.username}
                    placeholder="用户名"
                  />
                  <EyeInvisibleOutlined className="user-edit-readonly-icon" />
                </div>
              </Tooltip>
            ) : (
              <Input
                className="user-edit-form-input"
                placeholder="请输入用户名"
                onChange={(e) => handleDataChange('username', e.target.value)}
              />
            )}
          </Form.Item>
        </div>

        <div className="user-edit-form-group">
          <label className="user-edit-form-label">
            电子邮件
            {emailPermissions.required && <span className="user-edit-required">*</span>}
          </label>
          <Form.Item
            name="email"
            rules={emailPermissions.required ? [
              { required: true, message: '请输入邮箱' },
              { type: 'email', message: '请输入有效的邮箱地址' }
            ] : []}
          >
            {emailPermissions.readonly ? (
              <Tooltip title="此字段不可修改">
                <div style={{ position: 'relative' }}>
                  <Input
                    className="user-edit-form-input readonly"
                    disabled
                    value={initialData?.email}
                    placeholder="邮箱地址"
                  />
                  <EyeInvisibleOutlined className="user-edit-readonly-icon" />
                </div>
              </Tooltip>
            ) : (
              <Input
                className="user-edit-form-input"
                type="email"
                placeholder="请输入邮箱地址"
                onChange={(e) => handleDataChange('email', e.target.value)}
              />
            )}
          </Form.Item>
        </div>
      </div>

      <div className="user-edit-form-row">
        <div className="user-edit-form-group">
          <label className="user-edit-form-label">
            姓
            {lastNamePermissions.required && <span className="user-edit-required">*</span>}
          </label>
          <Form.Item
            name="last_name"
            rules={lastNamePermissions.required ? [
              { required: true, message: '请输入姓' }
            ] : []}
          >
            <Input
              className="user-edit-form-input"
              placeholder="请输入姓"
              onChange={(e) => handleDataChange('last_name', e.target.value)}
            />
          </Form.Item>
        </div>

        <div className="user-edit-form-group">
          <label className="user-edit-form-label">
            名
            {firstNamePermissions.required && <span className="user-edit-required">*</span>}
          </label>
          <Form.Item
            name="first_name"
            rules={firstNamePermissions.required ? [
              { required: true, message: '请输入名' }
            ] : []}
          >
            <Input
              className="user-edit-form-input"
              placeholder="请输入名"
              onChange={(e) => handleDataChange('first_name', e.target.value)}
            />
          </Form.Item>
        </div>
      </div>

      <div className="user-edit-form-row">
        <div className="user-edit-form-group">
          <div style={{ display: 'flex', alignItems: 'center', paddingTop: '12px' }}>
            <span style={{ marginRight: '8px', fontSize: '14px', color: '#333' }}>账户状态：</span>
            <Form.Item
              name="is_active"
              valuePropName="checked"
              style={{ margin: 0 }}
              initialValue={true}
            >
              <Switch
                checkedChildren="启用"
                unCheckedChildren="禁用"
                defaultChecked={true}
                onChange={(checked) => {
                  handleDataChange('is_active', checked);
                }}
              />
            </Form.Item>
          </div>
        </div>
        <div className="user-edit-form-group">
          {/* 占位，保持布局平衡 */}
        </div>
      </div>
    </div>
  );
};

export default BasicInfoTab;