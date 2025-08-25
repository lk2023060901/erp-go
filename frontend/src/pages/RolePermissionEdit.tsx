/**
 * Permission Manager - 1:1 Frappe 复刻版本
 */

import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Card,
  Table,
  Button,
  Space,
  Typography,
  message,
  Spin,
  Row,
  Col,
  Form,
  Modal,
  Popconfirm,
  Select,
  Checkbox,
  AutoComplete,
  Tooltip,
  Input
} from 'antd';
import {
  PlusOutlined,
  DeleteOutlined,
  ReloadOutlined,
  UserOutlined,
  TeamOutlined
} from '@ant-design/icons';
import { getRoleById, getEnabledRoles } from '../services/roleService';
import { 
  getDocTypesList, 
  getPermissionRulesList, 
  createPermissionRule, 
  updatePermissionRule, 
  deletePermissionRule
} from '../services/erpPermissionService';
import type { Role, DocType, PermissionRule } from '../types/auth';

const { Title, Text } = Typography;

export function RolePermissionEdit() {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [form] = Form.useForm();
  
  // 基础状态
  const [loading, setLoading] = useState(false);
  const [role, setRole] = useState<Role | null>(null);
  const [selectedDocType, setSelectedDocType] = useState<string>('');
  const [docTypeSearchValue, setDocTypeSearchValue] = useState<string>('');
  const [selectedRoleName, setSelectedRoleName] = useState<string>('');
  
  // 数据状态
  const [docTypes, setDocTypes] = useState<DocType[]>([]);
  const [permissionRules, setPermissionRules] = useState<PermissionRule[]>([]);
  const [allRoles, setAllRoles] = useState<Role[]>([]);
  
  // 编辑状态
  const [editingRule, setEditingRule] = useState<PermissionRule | null>(null);
  const [editFormVisible, setEditFormVisible] = useState(false);
  

  // 权限类型常量
  const PERMISSION_TYPES = [
    'can_read', 'can_write', 'can_create', 'can_delete',
    'can_submit', 'can_cancel', 'can_amend', 'can_print',
    'can_email', 'can_import', 'can_export', 'can_share', 'can_report'
  ];

  const PERMISSION_LABELS = {
    can_read: '读取',
    can_write: '写入', 
    can_create: '创建',
    can_delete: '删除',
    can_submit: '提交',
    can_cancel: '取消',
    can_amend: '修改',
    can_print: '打印',
    can_email: '邮件',
    can_import: '导入',
    can_export: '导出',
    can_share: '分享',
    can_report: '报表'
  };

  useEffect(() => {
    console.log('🔄 useEffect triggered with id:', id);
    if (id) {
      console.log('📞 Calling loadRoleData, loadPermissionData, and loadAllRoles');
      loadRoleData(parseInt(id));
      loadPermissionData(parseInt(id));
      loadAllRoles();
    } else {
      console.log('⚠️ No id provided, skipping data loading');
    }
  }, [id]);

  // 设置选中的角色名称
  useEffect(() => {
    if (role?.name) {
      setSelectedRoleName(role.name);
    }
  }, [role]);


  // 加载所有角色
  const loadAllRoles = async () => {
    try {
      console.log('🔍 Loading all roles...');
      const rolesData = await getEnabledRoles();
      console.log('📋 Loaded roles response:', rolesData);
      console.log('📋 Response type:', typeof rolesData);
      console.log('📋 Response constructor:', rolesData?.constructor?.name);
      
      // 数据可能直接是数组或者包装在 roles 属性中
      let roles = [];
      if (Array.isArray(rolesData)) {
        roles = rolesData;
      } else if (rolesData && rolesData.roles && Array.isArray(rolesData.roles)) {
        roles = rolesData.roles;
      } else {
        console.warn('⚠️ Unexpected roles data format:', rolesData);
      }
      
      console.log('📋 Final roles array:', roles);
      if (roles.length > 0) {
        setAllRoles(roles);
        console.log('✅ Successfully loaded', roles.length, 'roles');
      } else {
        console.warn('⚠️ No roles found in response');
      }
    } catch (error) {
      console.error('❌ Failed to load roles:', error);
      console.error('❌ Error details:', error.response || error.message);
    }
  };

  // 加载角色数据
  const loadRoleData = async (roleId: number) => {
    try {
      setLoading(true);
      const roleData = await getRoleById(roleId);
      setRole(roleData);
    } catch (error) {
      message.error('加载角色数据失败');
      console.error('Failed to load role:', error);
    } finally {
      setLoading(false);
    }
  };

  // 加载权限数据
  const loadPermissionData = async (roleId: number) => {
    try {
      setLoading(true);
      
      console.log('🔍 Loading permission data for role:', roleId, role?.name);
      
      // 并行加载DocTypes和权限规则
      const [docTypesResponse, permissionsResponse] = await Promise.all([
        getDocTypesList(),
        getPermissionRulesList(roleId)
      ]);
      
      console.log('📋 DocTypes response:', docTypesResponse);
      console.log('🔐 Permissions response:', permissionsResponse);
      
      const docTypesData = docTypesResponse.doctypes || [];
      const permissionsData = permissionsResponse.rules || [];
      
      console.log('📋 Processed docTypes:', docTypesData);
      console.log('🔐 Processed permissions:', permissionsData);
      
      setDocTypes(docTypesData);
      setPermissionRules(permissionsData);
      
      // 自动选择第一个DocType
      if (docTypesData.length > 0 && !selectedDocType) {
        setSelectedDocType(docTypesData[0].name);
      }
    } catch (error: any) {
      console.error('❌ Failed to load data:', error);
      console.error('Error details:', {
        message: error.message,
        stack: error.stack,
        response: error.response?.data
      });
      message.error(`加载权限数据失败: ${error.message}`);
      setDocTypes([]);
      setPermissionRules([]);
    } finally {
      setLoading(false);
    }
  };


  // 获取过滤后的DocType选项
  const getFilteredDocTypes = () => {
    if (!docTypeSearchValue) return docTypes;
    return docTypes.filter(docType => 
      docType.name.toLowerCase().includes(docTypeSearchValue.toLowerCase()) ||
      docType.label.toLowerCase().includes(docTypeSearchValue.toLowerCase()) ||
      docType.module.toLowerCase().includes(docTypeSearchValue.toLowerCase())
    );
  };

  // 按角色分组权限规则
  const getPermissionsByRole = () => {
    const rolePermissionMap: { [roleName: string]: PermissionRule } = {};
    
    // 为每个角色创建空的权限对象
    allRoles.forEach(role => {
      rolePermissionMap[role.name] = {
        id: 0,
        role_id: role.id,
        doc_type: selectedDocType,
        permission_level: 0,
        can_read: false,
        can_write: false,
        can_create: false,
        can_delete: false,
        can_submit: false,
        can_cancel: false,
        can_amend: false,
        can_print: false,
        can_email: false,
        can_import: false,
        can_export: false,
        can_share: false,
        can_report: false,
        only_if_creator: false,
        created_at: '',
        updated_at: ''
      };
    });
    
    // 将实际的权限规则填入对应的角色
    permissionRules
      .filter(rule => rule.doc_type === selectedDocType)
      .forEach(rule => {
        const roleName = allRoles.find(r => r.id === rule.role_id)?.name;
        if (roleName) {
          rolePermissionMap[roleName] = rule;
        }
      });
    
    return rolePermissionMap;
  };

  // 角色选择变更处理
  const handleRoleChange = async (roleName: string) => {
    setSelectedRoleName(roleName);
  };

  // 角色选择变更处理（用于AutoComplete的onSelect）
  const handleRoleSelectionChange = async (roleName: string) => {
    console.log('🔄 Role selection changed to:', roleName);
    
    // 查找选中的角色信息
    const selectedRole = allRoles.find(r => r.name === roleName);
    if (selectedRole) {
      console.log('📋 Loading permissions for role:', selectedRole);
      
      // 更新当前角色信息
      setRole(selectedRole);
      
      // 加载该角色的权限数据
      try {
        setLoading(true);
        await loadPermissionData(selectedRole.id);
      } catch (error) {
        console.error('❌ Failed to load permissions for selected role:', error);
      } finally {
        setLoading(false);
      }
    } else {
      console.warn('⚠️ Selected role not found in allRoles:', roleName);
    }
  };

  // 编辑权限规则
  const handleEditPermissionRule = (rule: PermissionRule) => {
    setEditingRule(rule);
    form.setFieldsValue({
      document_type: rule.document_type,
      role: role?.name,
      permission_level: rule.permission_level,
      ...PERMISSION_TYPES.reduce((acc, type) => {
        acc[type] = rule[type as keyof PermissionRule] || false;
        return acc;
      }, {} as Record<string, boolean>),
      if_owner: rule.if_owner || false
    });
    setEditFormVisible(true);
  };

  // 新建权限规则
  const handleCreatePermissionRule = (docTypeName?: string) => {
    const targetDocType = docTypeName || selectedDocType;
    setEditingRule(null);
    form.setFieldsValue({
      document_type: targetDocType,
      role: role?.name,
      permission_level: 0,
      ...PERMISSION_TYPES.reduce((acc, type) => {
        acc[type] = false;
        return acc;
      }, {} as Record<string, boolean>),
      if_owner: false
    });
    setEditFormVisible(true);
  };

  // 保存权限规则
  const handleSavePermissionRule = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);

      const ruleData = {
        role_id: parseInt(id!),
        doc_type: values.document_type,
        permission_level: values.permission_level,
        can_read: true,  // 默认设置读权限
        can_write: true, // 默认设置写权限
        can_create: values.can_create || false,
        can_delete: values.can_delete || false,
        can_submit: values.can_submit || false,
        can_cancel: values.can_cancel || false,
        can_amend: values.can_amend || false,
        can_print: values.can_print || false,
        can_email: values.can_email || false,
        can_import: values.can_import || false,
        can_export: values.can_export || false,
        can_share: values.can_share || false,
        can_report: values.can_report || false,
        can_set_user_permissions: values.can_set_user_permissions || false,
        only_if_creator: values.if_owner || false
      };

      if (editingRule) {
        await updatePermissionRule(editingRule.id, ruleData);
        message.success('更新权限规则成功');
      } else {
        await createPermissionRule(ruleData);
        message.success('创建权限规则成功');
      }

      setEditFormVisible(false);
      loadPermissionData(parseInt(id!));
    } catch (error) {
      message.error('保存权限规则失败');
      console.error('Save permission rule failed:', error);
    } finally {
      setLoading(false);
    }
  };

  // 删除权限规则
  const handleDeletePermissionRule = async (ruleId: number) => {
    try {
      setLoading(true);
      await deletePermissionRule(ruleId);
      message.success('删除权限规则成功');
      loadPermissionData(parseInt(id!));
    } catch (error) {
      message.error('删除权限规则失败');
      console.error('Delete permission rule failed:', error);
    } finally {
      setLoading(false);
    }
  };


  // 权限更改处理函数
  const handlePermissionChange = async (ruleId: number, permissionType: string, checked: boolean) => {
    try {
      const rule = permissionRules.find(r => r.id === ruleId);
      if (!rule) return;
      
      const updateData = {
        role_id: rule.role_id,
        doc_type: rule.doc_type,
        permission_level: rule.permission_level,
        can_read: rule.can_read || false,
        can_write: rule.can_write || false,
        can_create: rule.can_create || false,
        can_delete: rule.can_delete || false,
        can_submit: rule.can_submit || false,
        can_cancel: rule.can_cancel || false,
        can_amend: rule.can_amend || false,
        can_print: rule.can_print || false,
        can_email: rule.can_email || false,
        can_import: rule.can_import || false,
        can_export: rule.can_export || false,
        can_share: rule.can_share || false,
        can_report: rule.can_report || false,
        can_set_user_permissions: rule.can_set_user_permissions || false,
        only_if_creator: rule.only_if_creator || false,
        [permissionType]: checked
      };

      await updatePermissionRule(ruleId, updateData);
      // 重新加载权限数据
      if (id) {
        loadPermissionData(parseInt(id));
      } else {
        reloadCurrentPermissions();
      }
    } catch (error) {
      message.error('更新权限失败');
      console.error('Update permission failed:', error);
    }
  };




  if (!role) {
    return (
      <div style={{ padding: '24px', textAlign: 'center' }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div style={{ padding: '24px' }}>
      <Spin spinning={loading}>
        {/* Header左侧2个输入框 */}
        <div style={{
          display: 'flex',
          alignItems: 'center',
          gap: '16px',
          marginBottom: '24px',
          padding: '16px 24px',
          backgroundColor: '#fafafa',
          borderRadius: 8,
          border: '1px solid #f0f0f0'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            <Typography.Text>角色名称</Typography.Text>
            <AutoComplete
              style={{ width: 200 }}
              placeholder="角色名称"
              value={selectedRoleName}
              onChange={setSelectedRoleName}
              onSelect={(value) => {
                console.log('🎯 Selected role:', value);
                setSelectedRoleName(value);
                handleRoleSelectionChange(value);
              }}
              options={(allRoles || []).map(role => ({
                value: role.name,
                label: role.name
              }))}
              filterOption={(input, option) => {
                if (!input) return true; // 如果没有输入，显示所有选项
                return option!.label.toLowerCase().includes(input.toLowerCase());
              }}
              allowClear
              popupMatchSelectWidth={false}
              styles={{
                popup: {
                  minWidth: 200
                }
              }}
              onDropdownVisibleChange={(open) => {
                console.log('🔽 Dropdown visibility changed:', open);
                console.log('📋 Available roles:', allRoles);
              }}
            />
          </div>

          <Space>
            <Button
              icon={<UserOutlined />}
              onClick={() => navigate('/user-permissions')}
            >
              设置用户权限
            </Button>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => handleCreatePermissionRule()}
              style={{ backgroundColor: '#667eea', borderColor: '#667eea' }}
            >
              添加权限规则
            </Button>
          </Space>
        </div>

        {/* 权限规则表格 */}
        <div style={{ backgroundColor: 'white', borderRadius: 8, overflow: 'hidden' }}>
          <Table
            dataSource={permissionRules}
            rowKey="id"
            pagination={false}
            columns={[
              {
                title: '文档类型',
                dataIndex: 'doc_type',
                key: 'doc_type',
                width: 120,
                render: (docTypeName: string, record: PermissionRule) => {
                  console.log('🏷️ Rendering docType:', { docTypeName, record, docTypes });
                  
                  // 如果 doc_type 为空，显示提示信息
                  if (!docTypeName || docTypeName.trim() === '') {
                    return <span style={{ color: '#999', fontStyle: 'italic' }}>未指定文档类型</span>;
                  }
                  
                  const docType = docTypes.find(dt => dt.name === docTypeName);
                  return docType?.label || docTypeName;
                }
              },
              {
                title: '角色',
                dataIndex: 'role',
                key: 'role',
                width: 120,
                render: (_, record: PermissionRule) => {
                  const role = allRoles.find(r => r.id === record.role_id);
                  return role?.name || '未知角色';
                }
              },
              {
                title: '等级',
                dataIndex: 'permission_level',
                key: 'permission_level',
                width: 80,
                render: (level: number) => level
              },
              {
                title: '权限',
                key: 'permissions',
                width: 400,
                render: (_, record: PermissionRule) => {
                  const permissions = PERMISSION_TYPES.map(type => ({
                    key: type,
                    label: PERMISSION_LABELS[type as keyof typeof PERMISSION_LABELS],
                    checked: record[type as keyof PermissionRule] as boolean || false
                  }));

                  return (
                    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '8px' }}>
                      {permissions.map((perm, index) => (
                        <div key={index} style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                          <Checkbox 
                            checked={perm.checked} 
                            onChange={(e) => handlePermissionChange(record.id, perm.key, e.target.checked)}
                          />
                          <span style={{ fontSize: '12px' }}>{perm.label}</span>
                        </div>
                      ))}
                      <div style={{ display: 'flex', alignItems: 'center', gap: '4px', gridColumn: 'span 3' }}>
                        <Checkbox 
                          checked={record.only_if_creator || false}
                          onChange={(e) => handlePermissionChange(record.id, 'only_if_creator', e.target.checked)}
                        />
                        <span style={{ fontSize: '12px' }}>仅当创造者</span>
                      </div>
                    </div>
                  );
                }
              },
              {
                title: '',
                key: 'action',
                width: 60,
                render: (_, record: PermissionRule) => (
                  <Popconfirm
                    title="确定要删除此权限规则吗？"
                    onConfirm={() => handleDeletePermissionRule(record.id)}
                    okText="确定"
                    cancelText="取消"
                  >
                    <Button 
                      type="text" 
                      danger 
                      icon={<DeleteOutlined />}
                      size="small"
                      shape="circle"
                    />
                  </Popconfirm>
                )
              }
            ]}
          />
        </div>
      </Spin>

      {/* 权限规则编辑模态框 */}
      <Modal
        title={editingRule && editingRule.id > 0 ? '编辑权限规则' : '添加新的权限规则'}
        open={editFormVisible}
        onOk={handleSavePermissionRule}
        onCancel={() => setEditFormVisible(false)}
        width={800}
        destroyOnClose
        okText={editingRule && editingRule.id > 0 ? '更新' : '添加'}
        cancelText="取消"
        okButtonProps={{ 
          style: { backgroundColor: '#667eea', borderColor: '#667eea' }
        }}
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={{
            permission_level: 0
          }}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label="文档类型"
                name="document_type"
                rules={[{ required: true, message: '请选择文档类型' }]}
              >
                <Select
                  placeholder="选择文档类型"
                  disabled={!!editingRule}
                  showSearch
                  filterOption={(input, option) =>
                    (option?.children as string)?.toLowerCase().includes(input.toLowerCase())
                  }
                >
                  {docTypes.map(docType => (
                    <Select.Option key={docType.name} value={docType.name}>
                      {docType.label}
                    </Select.Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label="角色"
                name="role"
                rules={[{ required: true, message: '请选择角色' }]}
              >
                <Select 
                  disabled
                >
                  <Select.Option value={role?.name}>{role?.name}</Select.Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            label="权限级别"
            name="permission_level"
            rules={[{ required: true, message: '请选择权限级别' }]}
          >
            <Select placeholder="0">
              <Select.Option value={0}>0</Select.Option>
              <Select.Option value={1}>1</Select.Option>
              <Select.Option value={2}>2</Select.Option>
            </Select>
          </Form.Item>

          <Typography.Title level={5}>权限设置</Typography.Title>
          <Row gutter={[16, 8]}>
            {PERMISSION_TYPES.map(type => (
              <Col key={type} span={6}>
                <Form.Item name={type} valuePropName="checked" noStyle>
                  <Checkbox>{PERMISSION_LABELS[type as keyof typeof PERMISSION_LABELS]}</Checkbox>
                </Form.Item>
              </Col>
            ))}
            <Col span={12}>
              <Form.Item name="if_owner" valuePropName="checked" noStyle>
                <Checkbox>仅当创造者</Checkbox>
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>

    </div>
  );
}

export default RolePermissionEdit;