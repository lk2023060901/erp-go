import React, { useState } from 'react';
import { Form, Input } from 'antd';
import { 
  UserOutlined, 
  SettingOutlined, 
  FileTextOutlined, 
  BarChartOutlined,
  DownOutlined,
  UpOutlined
} from '@ant-design/icons';
import { TabProps, PermissionModule } from '../types/auth';
import PermissionModuleComponent from './PermissionModule';

const { TextArea } = Input;

interface RolePermissionsTabProps extends TabProps {
  form?: any;
}

// 权限模块数据定义
const permissionModulesData: PermissionModule[] = [
  {
    id: 'user-management',
    name: '用户管理',
    description: '管理用户账户、权限分配和用户信息',
    icon: <UserOutlined />,
    permissions: [
      { id: 'user-list', name: '查看用户列表', code: 'user.list', checked: true },
      { id: 'user-create', name: '创建用户', code: 'user.create', checked: true },
      { id: 'user-edit', name: '编辑用户信息', code: 'user.edit', checked: true },
      { id: 'user-delete', name: '删除用户', code: 'user.delete', checked: false },
      { id: 'user-reset-password', name: '重置密码', code: 'user.reset-password', checked: true },
      { id: 'user-assign-roles', name: '分配角色', code: 'user.assign-roles', checked: true },
    ],
    expanded: false
  },
  {
    id: 'system-management',
    name: '系统管理',
    description: '系统配置、参数设置和维护功能',
    icon: <SettingOutlined />,
    permissions: [
      { id: 'system-config', name: '系统配置', code: 'system.config', checked: false },
      { id: 'system-backup', name: '数据备份', code: 'system.backup', checked: false },
      { id: 'system-logs', name: '系统日志', code: 'system.logs', checked: false },
      { id: 'system-monitor', name: '系统监控', code: 'system.monitor', checked: false },
    ],
    expanded: false
  },
  {
    id: 'content-management',
    name: '内容管理',
    description: '文章、页面和媒体内容管理',
    icon: <FileTextOutlined />,
    permissions: [
      { id: 'content-create', name: '创建内容', code: 'content.create', checked: true },
      { id: 'content-edit', name: '编辑内容', code: 'content.edit', checked: true },
      { id: 'content-delete', name: '删除内容', code: 'content.delete', checked: false },
      { id: 'content-publish', name: '发布内容', code: 'content.publish', checked: true },
      { id: 'media-upload', name: '上传媒体', code: 'media.upload', checked: true },
      { id: 'media-manage', name: '媒体管理', code: 'media.manage', checked: false },
    ],
    expanded: false
  },
  {
    id: 'analytics',
    name: '数据分析',
    description: '报表查看、数据统计和分析功能',
    icon: <BarChartOutlined />,
    permissions: [
      { id: 'analytics-view', name: '查看报表', code: 'analytics.view', checked: true },
      { id: 'analytics-export', name: '导出数据', code: 'analytics.export', checked: true },
      { id: 'analytics-dashboard', name: '数据大屏', code: 'analytics.dashboard', checked: false },
      { id: 'analytics-custom', name: '自定义报表', code: 'analytics.custom', checked: false },
    ],
    expanded: false
  }
];

const RolePermissionsTab: React.FC<RolePermissionsTabProps> = ({ 
  onDataChange 
}) => {
  const [modules, setModules] = useState<PermissionModule[]>(permissionModulesData);
  const [allExpanded, setAllExpanded] = useState(false);

  // 处理数据变化
  const handleDataChange = (field: string, value: any) => {
    const newData = { [field]: value };
    onDataChange?.(newData);
  };

  // 处理权限变化
  const handlePermissionChange = (moduleId: string, permissionId: string, checked: boolean) => {
    const updatedModules = modules.map(module => {
      if (module.id === moduleId) {
        const updatedPermissions = module.permissions.map(permission => 
          permission.id === permissionId 
            ? { ...permission, checked }
            : permission
        );
        return { ...module, permissions: updatedPermissions };
      }
      return module;
    });
    
    setModules(updatedModules);
    
    // 收集所有选中的权限
    const selectedPermissions: string[] = [];
    updatedModules.forEach(module => {
      module.permissions.forEach(permission => {
        if (permission.checked) {
          selectedPermissions.push(permission.code);
        }
      });
    });
    
    handleDataChange('selected_permissions', selectedPermissions);
  };

  // 切换模块展开状态
  const toggleModule = (moduleId: string) => {
    setModules(prev => prev.map(module => 
      module.id === moduleId 
        ? { ...module, expanded: !module.expanded }
        : module
    ));
  };

  // 全选权限
  const selectAllPermissions = () => {
    const updatedModules = modules.map(module => ({
      ...module,
      permissions: module.permissions.map(permission => ({ ...permission, checked: true }))
    }));
    setModules(updatedModules);
    
    // 收集所有权限
    const allPermissions: string[] = [];
    updatedModules.forEach(module => {
      module.permissions.forEach(permission => {
        allPermissions.push(permission.code);
      });
    });
    
    handleDataChange('selected_permissions', allPermissions);
  };

  // 清除所有权限
  const clearAllPermissions = () => {
    const updatedModules = modules.map(module => ({
      ...module,
      permissions: module.permissions.map(permission => ({ ...permission, checked: false }))
    }));
    setModules(updatedModules);
    handleDataChange('selected_permissions', []);
  };

  // 模块权限全选/取消
  const selectModulePermissions = (moduleId: string, select: boolean) => {
    const updatedModules = modules.map(module => {
      if (module.id === moduleId) {
        return {
          ...module,
          permissions: module.permissions.map(permission => ({ ...permission, checked: select }))
        };
      }
      return module;
    });
    setModules(updatedModules);
    
    // 收集所有选中的权限
    const selectedPermissions: string[] = [];
    updatedModules.forEach(module => {
      module.permissions.forEach(permission => {
        if (permission.checked) {
          selectedPermissions.push(permission.code);
        }
      });
    });
    
    handleDataChange('selected_permissions', selectedPermissions);
  };

  // 展开/收起所有模块
  const toggleAllModules = () => {
    const newExpandedState = !allExpanded;
    setAllExpanded(newExpandedState);
    setModules(prev => prev.map(module => ({ ...module, expanded: newExpandedState })));
  };

  return (
    <div>
      {/* 权限管理头部 */}
      <div className="user-edit-permission-header">
        <h3 className="user-edit-permission-title">权限设置</h3>
        <div className="user-edit-permission-actions">
          <button 
            className="user-edit-action-btn select-all" 
            onClick={selectAllPermissions}
            type="button"
          >
            全选
          </button>
          <button 
            className="user-edit-action-btn clear-all" 
            onClick={clearAllPermissions}
            type="button"
          >
            全部取消
          </button>
          <button 
            className="user-edit-action-btn expand-all" 
            onClick={toggleAllModules}
            type="button"
          >
            {allExpanded ? <UpOutlined /> : <DownOutlined />}
            <span style={{ marginLeft: '4px' }}>
              {allExpanded ? '收起全部' : '展开全部'}
            </span>
          </button>
        </div>
      </div>

      {/* 权限模块列表 */}
      {modules.map(module => (
        <PermissionModuleComponent
          key={module.id}
          module={module}
          onToggle={() => toggleModule(module.id)}
          onPermissionChange={(permissionId, checked) => 
            handlePermissionChange(module.id, permissionId, checked)
          }
          onSelectAll={(select) => selectModulePermissions(module.id, select)}
        />
      ))}

      {/* 备注说明 */}
      <div className="user-edit-form-row" style={{ marginTop: '32px' }}>
        <div className="user-edit-form-group full-width">
          <label className="user-edit-form-label">权限备注</label>
          <Form.Item name="permission_notes">
            <TextArea
              className="user-edit-form-textarea"
              placeholder="请输入权限分配的备注说明..."
              rows={3}
              onChange={(e) => handleDataChange('permission_notes', e.target.value)}
            />
          </Form.Item>
        </div>
      </div>
    </div>
  );
};

export default RolePermissionsTab;