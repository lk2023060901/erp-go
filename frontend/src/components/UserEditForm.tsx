import React, { useState, useEffect } from 'react';
import { Form, App } from 'antd';
import { 
  UserEditFormProps, 
  UserFormData, 
  FormMode, 
  FieldPermissions, 
  User 
} from '../types/auth';
import BasicInfoTab from './BasicInfoTab';
import DetailInfoTab from './DetailInfoTab';
import RolePermissionsTab from './RolePermissionsTab';
import './UserEditForm.css';

// 获取字段权限的工具函数
export const getFieldPermissions = (mode: FormMode, fieldName: string): FieldPermissions => {
  const permissions: Record<string, Record<FormMode, FieldPermissions>> = {
    email: {
      create: { readonly: false, required: true, visible: true },
      edit: { readonly: true, required: false, visible: true },
      profile: { readonly: true, required: false, visible: true }
    },
    username: {
      create: { readonly: false, required: true, visible: true },
      edit: { readonly: true, required: false, visible: true },
      profile: { readonly: true, required: false, visible: true }
    },
    first_name: {
      create: { readonly: false, required: true, visible: true },
      edit: { readonly: false, required: true, visible: true },
      profile: { readonly: false, required: true, visible: true }
    },
    last_name: {
      create: { readonly: false, required: true, visible: true },
      edit: { readonly: false, required: true, visible: true },
      profile: { readonly: false, required: true, visible: true }
    }
  };
  
  return permissions[fieldName]?.[mode] || { readonly: false, required: false, visible: true };
};

// 将User类型转换为UserFormData类型
const userToFormData = (user: User): UserFormData => {
  return {
    username: user.username,
    email: user.email,
    first_name: user.first_name,
    last_name: user.last_name,
    is_active: user.is_enabled, // 注意字段名映射
    gender: user.gender,
    birth_date: user.birth_date,
    phone: user.phone || '',
    address: '', // User类型中没有此字段，使用默认值
    emergency_contact: '',
    emergency_phone: '',
    biography: '',
    selected_permissions: [],
    permission_notes: ''
  };
};

const UserEditForm: React.FC<UserEditFormProps> = ({
  user,
  mode,
  onSubmit,
  onCancel,
  showActionButtons = true,
  initialTab
}) => {
  const { message } = App.useApp();
  const [form] = Form.useForm();
  
  // 状态管理
  const [activeTab, setActiveTab] = useState(initialTab || 'basic');
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<UserFormData>(() => {
    if (user && mode !== 'create') {
      return userToFormData(user);
    }
    return {
      username: '',
      email: '',
      first_name: '',
      last_name: '',
      is_active: true,
      gender: undefined,
      birth_date: '',
      phone: '',
      address: '',
      emergency_contact: '',
      emergency_phone: '',
      biography: '',
      selected_permissions: [],
      permission_notes: ''
    };
  });

  // 初始化表单数据
  useEffect(() => {
    if (user && mode !== 'create') {
      const initialData = userToFormData(user);
      setFormData(initialData);
      form.setFieldsValue(initialData);
    }
  }, [user, mode, form]);

  // 处理表单提交
  const handleSubmit = async () => {
    try {
      setLoading(true);
      const values = await form.validateFields();
      const submitData: UserFormData = { ...formData, ...values };
      await onSubmit(submitData);
      message.success(mode === 'create' ? '创建用户成功' : '更新用户信息成功');
    } catch (error: any) {
      if (error.errorFields) {
        // 表单验证错误
        message.error('请检查表单中的错误信息');
      } else {
        // API错误
        message.error(error.message || '操作失败');
      }
    } finally {
      setLoading(false);
    }
  };

  // 处理取消
  const handleCancel = () => {
    form.resetFields();
    onCancel?.();
  };

  // 页面标题
  const getPageTitle = () => {
    switch (mode) {
      case 'create':
        return '创建用户';
      case 'edit':
        return '编辑用户信息';
      case 'profile':
        return '个人资料';
      default:
        return '用户信息';
    }
  };

  // 选项卡配置
  const tabs = [
    { key: 'basic', label: '基本资料' },
    { key: 'details', label: '详细信息' },
    { key: 'roles', label: '角色权限' }
  ];

  return (
    <div className="user-edit-form">
      {/* 面包屑导航 - 隐藏不必要的导航 */}
      {/* <nav className="user-edit-breadcrumb">
        <div className="user-edit-breadcrumb-list">
          <span className="user-edit-breadcrumb-item">用户管理</span>
          <span className="user-edit-breadcrumb-separator">/</span>
          <span className="user-edit-breadcrumb-item active">{getPageTitle()}</span>
        </div>
      </nav> */}

      {/* 主内容卡片 */}
      <div className="user-edit-main-card">
        {/* 卡片头部 */}
        <div className="user-edit-card-header">
          <div className="user-edit-header-content">
            <h1 className="user-edit-page-title">{getPageTitle()}</h1>
          </div>

          {/* 选项卡 */}
          <div className="user-edit-tabs">
            {tabs.map(tab => (
              <button
                key={tab.key}
                className={`user-edit-tab-item ${activeTab === tab.key ? 'active' : ''}`}
                onClick={() => setActiveTab(tab.key)}
              >
                {tab.label}
              </button>
            ))}
          </div>
        </div>

        {/* 表单内容 */}
        <Form
          form={form}
          layout="vertical"
          initialValues={formData}
          onFinish={handleSubmit}
        >
          {/* 基本资料选项卡 */}
          <div className={`user-edit-tab-content ${activeTab === 'basic' ? 'active' : ''}`}>
            <BasicInfoTab
              mode={mode}
              initialData={formData}
              onDataChange={(data) => setFormData(prev => ({ ...prev, ...data }))}
            />
          </div>

          {/* 详细信息选项卡 */}
          <div className={`user-edit-tab-content ${activeTab === 'details' ? 'active' : ''}`}>
            <DetailInfoTab
              mode={mode}
              initialData={formData}
              onDataChange={(data) => setFormData(prev => ({ ...prev, ...data }))}
            />
          </div>

          {/* 角色权限选项卡 */}
          <div className={`user-edit-tab-content ${activeTab === 'roles' ? 'active' : ''}`}>
            <RolePermissionsTab
              mode={mode}
              initialData={formData}
              onDataChange={(data) => setFormData(prev => ({ ...prev, ...data }))}
              form={form}
              userId={user?.id} // 传递用户ID用于编辑模式加载现有权限
            />
          </div>
        </Form>

        {/* 底部操作栏 */}
        {showActionButtons && (
          <div className="user-edit-form-footer">
            <button 
              type="button" 
              className="user-edit-btn user-edit-btn-default"
              onClick={handleCancel}
              disabled={loading}
            >
              取消
            </button>
            <button 
              type="button" 
              className="user-edit-btn user-edit-btn-primary"
              onClick={handleSubmit}
              disabled={loading}
            >
              {loading ? '保存中...' : (mode === 'create' ? '创建' : '保存更改')}
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default UserEditForm;