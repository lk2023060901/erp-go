import React, { useState, useEffect } from 'react';
import { Form, Input, Select, Card, Alert, Spin, Checkbox, Tag, Space, Button, Row, Col } from 'antd';
import { 
  UserOutlined, 
  SettingOutlined, 
  FileTextOutlined, 
  BarChartOutlined,
  DownOutlined,
  UpOutlined,
  TeamOutlined,
  SafetyCertificateOutlined
} from '@ant-design/icons';
import { 
  TabProps, 
  PermissionModule, 
  Role, 
  DocType 
} from '../types/auth';
import { getRolesList } from '../services/roleService';
import { 
  getDocTypesList, 
  getPermissionRulesList,
  getUserPermissionsList
} from '../services/erpPermissionService';
import PermissionModuleComponent from './PermissionModule';

const { TextArea } = Input;
const { Option } = Select;

interface RolePermissionsTabProps extends TabProps {
  form?: any;
  userId?: number; // 用于编辑模式时获取用户现有权限
}

// 权限类型图标映射
const getPermissionIcon = (type: string) => {
  const iconMap: Record<string, React.ReactNode> = {
    'read': <FileTextOutlined style={{ color: '#52c41a' }} />,
    'write': <SettingOutlined style={{ color: '#1890ff' }} />,
    'create': <UserOutlined style={{ color: '#722ed1' }} />,
    'delete': <BarChartOutlined style={{ color: '#f5222d' }} />,
    'submit': <SafetyCertificateOutlined style={{ color: '#fa8c16' }} />,
    'cancel': <DownOutlined style={{ color: '#faad14' }} />,
    'amend': <UpOutlined style={{ color: '#13c2c2' }} />,
    'print': <FileTextOutlined style={{ color: '#eb2f96' }} />,
    'email': <TeamOutlined style={{ color: '#52c41a' }} />,
    'import': <UpOutlined style={{ color: '#1890ff' }} />,
    'export': <DownOutlined style={{ color: '#722ed1' }} />,
    'share': <TeamOutlined style={{ color: '#fa8c16' }} />,
    'report': <BarChartOutlined style={{ color: '#faad14' }} />
  };
  return iconMap[type] || <SettingOutlined />;
};

const RolePermissionsTab: React.FC<RolePermissionsTabProps> = ({ 
  onDataChange,
  userId,
  mode
}) => {
  // 状态管理
  const [loading, setLoading] = useState(false);
  const [activeTabKey, setActiveTabKey] = useState('roles');
  
  // ERP权限系统状态
  const [roles, setRoles] = useState<Role[]>([]);
  const [docTypes, setDocTypes] = useState<DocType[]>([]);
  
  // 选择状态
  const [selectedRoles, setSelectedRoles] = useState<number[]>([]);
  const [selectedDocTypes, setSelectedDocTypes] = useState<string[]>([]);
  const [customPermissions, setCustomPermissions] = useState<Record<string, Record<string, boolean>>>({});
  
  // 传统权限系统状态（保持兼容）
  const [modules, setModules] = useState<PermissionModule[]>([]);
  const [allExpanded, setAllExpanded] = useState(false);

  // 初始化数据加载
  useEffect(() => {
    loadInitialData();
  }, []);
  
  // 加载用户现有权限（编辑模式）
  useEffect(() => {
    if (userId && mode === 'edit') {
      loadUserPermissions();
    }
  }, [userId, mode]);
  
  const loadInitialData = async () => {
    setLoading(true);
    try {
      const [rolesResponse, docTypesResponse] = await Promise.all([
        getRolesList({ page: 1, size: 1000 }),
        getDocTypesList({ page: 1, size: 1000 })
      ]);
      
      setRoles(rolesResponse.roles);
      setDocTypes(docTypesResponse.doctypes);
      
      // 构建传统权限模块（保持兼容性）
      buildLegacyModules(docTypesResponse.doctypes);
    } catch (error) {
      console.error('Load permission data failed:', error);
    } finally {
      setLoading(false);
    }
  };
  
  const loadUserPermissions = async () => {
    if (!userId) return;
    
    try {
      const [permissionsResponse, userPermsResponse] = await Promise.all([
        getPermissionRulesList({ page: 1, size: 1000 }),
        getUserPermissionsList({ page: 1, size: 1000, user_id: userId })
      ]);
      
      // TODO: 处理权限规则和用户权限数据
      console.log('Permission rules loaded:', permissionsResponse.rules);
      console.log('User permissions loaded:', userPermsResponse.permissions);
      
      // 根据现有权限设置选择状态
      // TODO: 这里需要根据实际的用户角色关联来设置selectedRoles
      // 目前数据结构中没有直接的用户角色关联，需要后端支持
    } catch (error) {
      console.error('Load user permissions failed:', error);
    }
  };
  
  const buildLegacyModules = (docTypes: DocType[]) => {
    const moduleMap: Record<string, PermissionModule> = {};
    
    docTypes.forEach(docType => {
      const moduleId = docType.module.toLowerCase().replace(/\s+/g, '-');
      
      if (!moduleMap[moduleId]) {
        moduleMap[moduleId] = {
          id: moduleId,
          name: docType.module,
          description: `${docType.module}模块的文档类型权限管理`,
          icon: getModuleIcon(docType.module),
          permissions: [],
          expanded: false
        };
      }
      
      // 为每个文档类型创建权限项
      PERMISSION_TYPES.forEach(permType => {
        moduleMap[moduleId].permissions.push({
          id: `${docType.name}-${permType}`,
          name: `${PERMISSION_LABELS[permType]} ${docType.label}`,
          code: `${docType.name}.${permType}`,
          description: `对${docType.label}进行${PERMISSION_LABELS[permType]}操作`,
          checked: false
        });
      });
    });
    
    setModules(Object.values(moduleMap));
  };
  
  const getModuleIcon = (moduleName: string) => {
    const iconMap: Record<string, React.ReactNode> = {
      'Core': <UserOutlined />,
      'System': <SettingOutlined />,
      'Content': <FileTextOutlined />,
      'Analytics': <BarChartOutlined />
    };
    return iconMap[moduleName] || <SettingOutlined />;
  };
  
  // 处理数据变化
  const handleDataChange = (field: string, value: any) => {
    const newData = { [field]: value };
    onDataChange?.(newData);
  };

  // 处理角色选择变化
  const handleRoleChange = (roleIds: number[]) => {
    setSelectedRoles(roleIds);
    handleDataChange('selected_roles', roleIds);
  };
  
  // 处理文档类型选择变化  
  const handleDocTypeChange = (docTypeNames: string[]) => {
    setSelectedDocTypes(docTypeNames);
    handleDataChange('selected_doctypes', docTypeNames);
  };
  
  // 处理自定义权限变化
  const handleCustomPermissionChange = (docType: string, permission: string, checked: boolean) => {
    setCustomPermissions(prev => ({
      ...prev,
      [docType]: {
        ...prev[docType],
        [permission]: checked
      }
    }));
    
    handleDataChange('custom_permissions', {
      ...customPermissions,
      [docType]: {
        ...customPermissions[docType],
        [permission]: checked
      }
    });
  };
  
  // 处理传统权限变化
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

  const tabItems = [
    {
      key: 'roles',
      label: (
        <span>
          <TeamOutlined />
          角色分配
        </span>
      ),
      children: (
        <Card title="角色分配" size="small">
          <Alert
            message="角色权限管理"
            description="通过分配角色来继承对应的权限规则。用户将获得所选角色的所有权限。"
            type="info"
            showIcon
            style={{ marginBottom: 16 }}
          />
          
          <Form.Item label="选择角色" name="selected_roles">
            <Select
              mode="multiple"
              placeholder="请选择用户角色"
              value={selectedRoles}
              onChange={handleRoleChange}
              style={{ width: '100%' }}
            >
              {roles.map(role => (
                <Option key={role.id} value={role.id}>
                  <Space>
                    <Tag color={role.is_system_role ? 'red' : 'blue'}>
                      {role.name}
                    </Tag>
                    {role.description && <span style={{ color: '#999' }}>- {role.description}</span>}
                  </Space>
                </Option>
              ))}
            </Select>
          </Form.Item>
          
          {selectedRoles.length > 0 && (
            <div style={{ marginTop: 16 }}>
              <h4>已选择的角色：</h4>
              <Space wrap>
                {selectedRoles.map(roleId => {
                  const role = roles.find(r => r.id === roleId);
                  return role ? (
                    <Tag 
                      key={roleId} 
                      color={role.is_system_role ? 'red' : 'blue'}
                      closable
                      onClose={() => handleRoleChange(selectedRoles.filter(id => id !== roleId))}
                    >
                      {role.name}
                    </Tag>
                  ) : null;
                })}
              </Space>
            </div>
          )}
        </Card>
      )
    },
    {
      key: 'custom',
      label: (
        <span>
          <SafetyCertificateOutlined />
          自定义权限
        </span>
      ),
      children: (
        <Card title="自定义权限" size="small">
          <Alert
            message="文档类型权限"
            description="为用户分配特定文档类型的详细权限。这将覆盖角色权限的设置。"
            type="warning"
            showIcon
            style={{ marginBottom: 16 }}
          />
          
          <Form.Item label="选择文档类型" name="selected_doctypes">
            <Select
              mode="multiple"
              placeholder="请选择要配置权限的文档类型"
              value={selectedDocTypes}
              onChange={handleDocTypeChange}
              style={{ width: '100%' }}
            >
              {docTypes.map(docType => (
                <Option key={docType.name} value={docType.name}>
                  <Space>
                    <Tag color="cyan">{docType.module}</Tag>
                    {docType.label} ({docType.name})
                  </Space>
                </Option>
              ))}
            </Select>
          </Form.Item>
          
          {selectedDocTypes.length > 0 && (
            <div style={{ marginTop: 20 }}>
              <h4>权限配置：</h4>
              {selectedDocTypes.map(docTypeName => {
                const docType = docTypes.find(dt => dt.name === docTypeName);
                if (!docType) return null;
                
                return (
                  <Card key={docTypeName} size="small" style={{ marginBottom: 16 }}>
                    <h5>
                      <Space>
                        <Tag color="cyan">{docType.module}</Tag>
                        {docType.label} ({docType.name})
                      </Space>
                    </h5>
                    <Row gutter={[8, 8]}>
                      {PERMISSION_TYPES.map(permType => (
                        <Col span={6} key={permType}>
                          <Checkbox
                            checked={customPermissions[docTypeName]?.[permType] || false}
                            onChange={(e) => handleCustomPermissionChange(docTypeName, permType, e.target.checked)}
                          >
                            <Space>
                              {getPermissionIcon(permType)}
                              {PERMISSION_LABELS[permType]}
                            </Space>
                          </Checkbox>
                        </Col>
                      ))}
                    </Row>
                  </Card>
                );
              })}
            </div>
          )}
        </Card>
      )
    },
    {
      key: 'legacy',
      label: (
        <span>
          <FileTextOutlined />
          传统权限
        </span>
      ),
      children: (
        <div>
          <Alert
            message="传统权限系统（兼容模式）"
            description="使用传统的模块化权限管理。建议优先使用上述角色分配和自定义权限。"
            type="info"
            showIcon
            style={{ marginBottom: 16 }}
          />
          
          {loading ? (
            <div style={{ textAlign: 'center', padding: 40 }}>
              <Spin size="large" />
            </div>
          ) : (
            <>
              {/* 权限管理头部 */}
              <div className="user-edit-permission-header">
                <h4 className="user-edit-permission-title">模块权限设置</h4>
                <div className="user-edit-permission-actions">
                  <Button size="small" onClick={selectAllPermissions}>
                    全选
                  </Button>
                  <Button size="small" onClick={clearAllPermissions}>
                    全部取消
                  </Button>
                  <Button size="small" onClick={toggleAllModules}>
                    {allExpanded ? <UpOutlined /> : <DownOutlined />}
                    {allExpanded ? '收起全部' : '展开全部'}
                  </Button>
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
            </>
          )}
        </div>
      )
    }
  ];
  
  return (
    <Spin spinning={loading}>
      <div>
        {/* 权限管理选项卡 */}
        <Card>
          <div style={{ marginBottom: 16 }}>
            <Space>
              {tabItems.map(item => (
                <Button
                  key={item.key}
                  type={activeTabKey === item.key ? 'primary' : 'default'}
                  onClick={() => setActiveTabKey(item.key)}
                  style={{ marginRight: 8 }}
                >
                  {item.label}
                </Button>
              ))}
            </Space>
          </div>
          
          {/* 当前选项卡内容 */}
          {tabItems.find(item => item.key === activeTabKey)?.children}
        </Card>

        {/* 备注说明 */}
        <Card title="权限备注" style={{ marginTop: 16 }}>
          <Form.Item name="permission_notes">
            <TextArea
              placeholder="请输入权限分配的备注说明..."
              rows={3}
              onChange={(e) => handleDataChange('permission_notes', e.target.value)}
            />
          </Form.Item>
        </Card>
      </div>
    </Spin>
  );
};

export default RolePermissionsTab;