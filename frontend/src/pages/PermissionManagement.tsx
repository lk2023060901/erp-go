import React, { useState, useEffect } from 'react';
import {
  Card,
  Button,
  Table,
  Space,
  Modal,
  Form,
  Input,
  Select,
  Switch,
  message,
  Popconfirm,
  Tag,
  Tree,
  Row,
  Col,
  Tabs,
  Alert,
  Spin,
  Badge,
  Divider,
  Checkbox,
  Typography,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SyncOutlined,
  ApiOutlined,
  MenuOutlined,
  SettingOutlined,
  BranchesOutlined,
  ReloadOutlined,
  DatabaseOutlined,
  TeamOutlined,
  SafetyCertificateOutlined,
  FileTextOutlined,
} from '@ant-design/icons';
import type { TableColumnsType, TabsProps } from 'antd';
import { 
  Permission, 
  DocType, 
  PermissionRule, 
  UserPermission,
  Role
} from '../types/auth';
import {
  // 传统权限服务（保留向后兼容）
  getPermissionsList,
  getPermissionTree,
  createPermission,
  updatePermission,
  deletePermission,
  batchDeletePermissions,
  getModules,
  syncApiPermissions,
  generatePermissionCode,
  type PermissionsQueryParams,
  type CreatePermissionRequest,
  type UpdatePermissionRequest,
  type PermissionTreeNode,
} from '../services/permissionService';
import {
  // ERP权限服务
  getDocTypesList,
  createDocType,
  updateDocType,
  deleteDocType,
  getPermissionRulesList,
  createPermissionRule,
  updatePermissionRule,
  deletePermissionRule,
  getUserPermissionsList,
  createUserPermission,
  updateUserPermission,
  deleteUserPermission,
  type CreateDocTypeRequest,
  type UpdateDocTypeRequest,
  type CreatePermissionRuleRequest,
  type UpdatePermissionRuleRequest,
  type CreateUserPermissionRequest,
  type UpdateUserPermissionRequest,
} from '../services/erpPermissionService';
import { getRolesList } from '../services/roleService';

const { Search } = Input;
const { TextArea } = Input;
const { Option } = Select;
const { Title, Text } = Typography;

// 权限类型常量
const PERMISSION_TYPES = [
  'can_read', 'can_write', 'can_create', 'can_delete', 
  'can_submit', 'can_cancel', 'can_amend', 'can_print', 
  'can_email', 'can_import', 'can_export', 'can_share', 
  'can_report', 'can_set_user_permissions'
] as const;

const PERMISSION_LABELS: Record<string, string> = {
  can_read: '读取',
  can_write: '写入',
  can_create: '创建',
  can_delete: '删除',
  can_submit: '提交',
  can_cancel: '取消',
  can_amend: '修正',
  can_print: '打印',
  can_email: '邮件',
  can_import: '导入',
  can_export: '导出',
  can_share: '分享',
  can_report: '报告',
  can_set_user_permissions: '设置用户权限'
};

const PermissionManagement: React.FC = () => {
  // 通用状态
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState('doctypes');
  const [searchText, setSearchText] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form] = Form.useForm();
  
  // DocType 管理状态
  const [docTypes, setDocTypes] = useState<DocType[]>([]);
  const [docTypesTotal, setDocTypesTotal] = useState(0);
  const [editingDocType, setEditingDocType] = useState<DocType | null>(null);
  
  // 权限规则管理状态
  const [permissionRules, setPermissionRules] = useState<PermissionRule[]>([]);
  const [permissionRulesTotal, setPermissionRulesTotal] = useState(0);
  const [editingPermissionRule, setEditingPermissionRule] = useState<PermissionRule | null>(null);
  const [roles, setRoles] = useState<Role[]>([]);
  
  // 用户权限管理状态
  const [userPermissions, setUserPermissions] = useState<UserPermission[]>([]);
  const [userPermissionsTotal, setUserPermissionsTotal] = useState(0);
  const [editingUserPermission, setEditingUserPermission] = useState<UserPermission | null>(null);
  
  // 传统权限系统状态（保持向后兼容）
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [permissionTree, setPermissionTree] = useState<PermissionTreeNode[]>([]);
  const [modules, setModules] = useState<string[]>([]);
  const [treeLoading, setTreeLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [selectedModule, setSelectedModule] = useState<string>();
  const [statusFilter, setStatusFilter] = useState<boolean>();
  const [editingPermission, setEditingPermission] = useState<Permission | null>(null);

  useEffect(() => {
    loadRoles(); // 加载角色列表供权限规则使用
  }, []);
  
  useEffect(() => {
    // 根据当前选项卡加载对应数据
    switch (activeTab) {
      case 'doctypes':
        loadDocTypes();
        break;
      case 'permission-rules':
        loadPermissionRules();
        break;
      case 'user-permissions':
        loadUserPermissions();
        break;
      case 'legacy-permissions':
        loadPermissions();
        loadModules();
        break;
      case 'tree':
        loadPermissionTree();
        break;
      default:
        break;
    }
  }, [activeTab, currentPage, pageSize, searchText, selectedModule, statusFilter]);

  // ========== DocType 管理函数 ==========
  const loadDocTypes = async () => {
    setLoading(true);
    try {
      const response = await getDocTypesList({
        page: currentPage,
        size: pageSize,
        ...(selectedModule && { module: selectedModule }),
      });
      setDocTypes(response.doctypes);
      setDocTypesTotal(response.total);
    } catch (error) {
      message.error('加载文档类型列表失败');
    } finally {
      setLoading(false);
    }
  };
  
  // ========== 权限规则管理函数 ==========
  const loadPermissionRules = async () => {
    setLoading(true);
    try {
      const response = await getPermissionRulesList({
        page: currentPage,
        size: pageSize,
      });
      setPermissionRules(response.rules);
      setPermissionRulesTotal(response.total);
    } catch (error) {
      message.error('加载权限规则列表失败');
    } finally {
      setLoading(false);
    }
  };
  
  // ========== 用户权限管理函数 ==========
  const loadUserPermissions = async () => {
    setLoading(true);
    try {
      const response = await getUserPermissionsList({
        page: currentPage,
        size: pageSize,
      });
      setUserPermissions(response.permissions);
      setUserPermissionsTotal(response.total);
    } catch (error) {
      message.error('加载用户权限列表失败');
    } finally {
      setLoading(false);
    }
  };
  
  const loadRoles = async () => {
    try {
      const response = await getRolesList({ page: 1, size: 1000 });
      setRoles(response.roles);
    } catch (error) {
      console.error('Load roles failed:', error);
    }
  };
  
  // ========== 传统权限系统函数（保持向后兼容） ==========
  const loadPermissions = async () => {
    setLoading(true);
    try {
      const params: PermissionsQueryParams = {
        page: currentPage,
        size: pageSize,
        ...(searchText && { search: searchText }),
        ...(selectedModule && { module: selectedModule }),
        ...(statusFilter !== undefined && { is_enabled: statusFilter }),
      };
      
      const response = await getPermissionsList(params);
      setPermissions(response.permissions);
      setTotal(response.total);
    } catch (error) {
      message.error('加载权限列表失败');
    } finally {
      setLoading(false);
    }
  };

  const loadPermissionTree = async () => {
    setTreeLoading(true);
    try {
      const tree = await getPermissionTree();
      setPermissionTree(tree || []);
    } catch (error) {
      message.error('加载权限树失败');
      setPermissionTree([]);
    } finally {
      setTreeLoading(false);
    }
  };

  const loadModules = async () => {
    try {
      const moduleList = await getModules();
      setModules(moduleList);
    } catch (error) {
      console.error('Load modules failed:', error);
    }
  };

  // ========== 通用操作函数 ==========
  const handleCreate = () => {
    if (activeTab === 'doctypes') {
      setEditingDocType(null);
    } else if (activeTab === 'permission-rules') {
      setEditingPermissionRule(null);
    } else if (activeTab === 'user-permissions') {
      setEditingUserPermission(null);
    } else {
      setEditingPermission(null);
    }
    
    setIsModalOpen(true);
    form.resetFields();
    
    // 根据不同类型设置默认值
    if (activeTab === 'doctypes') {
      form.setFieldsValue({
        is_submittable: false,
        is_child_table: false,
        has_workflow: false,
        track_changes: false,
        applies_to_all_users: false,
        max_attachments: 10,
        sort_order: 'ASC',
      });
    } else if (activeTab === 'permission-rules') {
      form.setFieldsValue({
        permission_level: 0,
        ...Object.fromEntries(PERMISSION_TYPES.map(type => [type, false])),
        if_owner: false,
      });
    } else if (activeTab === 'user-permissions') {
      form.setFieldsValue({
        hide_descendants: false,
        is_default: false,
      });
    } else {
      form.setFieldsValue({
        is_enabled: true,
        is_menu: false,
        is_button: false,
        is_api: true,
        level: 1,
        sort_order: 0,
      });
    }
  };

  const handleEdit = (record: any) => {
    if (activeTab === 'doctypes') {
      setEditingDocType(record as DocType);
      form.setFieldsValue({
        ...record,
        search_fields: typeof record.search_fields === 'string' 
          ? JSON.parse(record.search_fields || '[]') 
          : record.search_fields,
      });
    } else if (activeTab === 'permission-rules') {
      setEditingPermissionRule(record as PermissionRule);
      form.setFieldsValue(record);
    } else if (activeTab === 'user-permissions') {
      setEditingUserPermission(record as UserPermission);
      form.setFieldsValue(record);
    } else {
      setEditingPermission(record as Permission);
      form.setFieldsValue({
        ...record,
        parent_id: record.parent_id || undefined,
      });
    }
    setIsModalOpen(true);
  };

  const handleDelete = async (id: number) => {
    try {
      if (activeTab === 'doctypes') {
        const docType = docTypes.find(dt => dt.id === id);
        if (docType) {
          await deleteDocType(docType.name);
        }
      } else if (activeTab === 'permission-rules') {
        await deletePermissionRule(id);
      } else if (activeTab === 'user-permissions') {
        await deleteUserPermission(id);
      } else {
        await deletePermission(id);
      }
      
      message.success('删除成功');
      
      // 重新加载对应的数据
      switch (activeTab) {
        case 'doctypes':
          loadDocTypes();
          break;
        case 'permission-rules':
          loadPermissionRules();
          break;
        case 'user-permissions':
          loadUserPermissions();
          break;
        case 'legacy-permissions':
          loadPermissions();
          break;
        case 'tree':
          loadPermissionTree();
          break;
      }
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleBatchDelete = async () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要删除的项目');
      return;
    }

    try {
      // 目前只有传统权限系统支持批量删除
      if (activeTab === 'legacy-permissions') {
        await batchDeletePermissions(selectedRowKeys as number[]);
        message.success('批量删除成功');
        setSelectedRowKeys([]);
        loadPermissions();
      } else {
        message.info('该功能暂不支持批量删除');
      }
    } catch (error) {
      message.error('批量删除失败');
    }
  };

  const handleSyncApi = async () => {
    try {
      setLoading(true);
      const result = await syncApiPermissions();
      message.success(`同步完成：新建 ${result.created} 个，更新 ${result.updated} 个权限`);
      loadPermissions();
      if (activeTab === 'tree') {
        loadPermissionTree();
      }
    } catch (error) {
      message.error('同步API权限失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      
      if (activeTab === 'doctypes') {
        if (editingDocType) {
          const updateData: UpdateDocTypeRequest = {
            label: values.label,
            description: values.description,
            is_submittable: values.is_submittable,
            is_child_table: values.is_child_table,
            has_workflow: values.has_workflow,
            track_changes: values.track_changes,
            applies_to_all_users: values.applies_to_all_users,
            max_attachments: values.max_attachments,
            naming_rule: values.naming_rule,
            title_field: values.title_field,
            search_fields: values.search_fields,
            sort_field: values.sort_field,
            sort_order: values.sort_order,
          };
          await updateDocType(editingDocType.name, updateData);
        } else {
          const createData: CreateDocTypeRequest = {
            name: values.name,
            label: values.label,
            module: values.module,
            description: values.description,
            is_submittable: values.is_submittable,
            is_child_table: values.is_child_table,
            has_workflow: values.has_workflow,
            track_changes: values.track_changes,
            applies_to_all_users: values.applies_to_all_users,
            max_attachments: values.max_attachments,
            naming_rule: values.naming_rule,
            title_field: values.title_field,
            search_fields: values.search_fields,
            sort_field: values.sort_field,
            sort_order: values.sort_order,
          };
          await createDocType(createData);
        }
        loadDocTypes();
      } else if (activeTab === 'permission-rules') {
        if (editingPermissionRule) {
          const updateData: UpdatePermissionRuleRequest = Object.fromEntries(
            Object.entries(values).filter(([key]) => key !== 'role_id' && key !== 'document_type')
          );
          await updatePermissionRule(editingPermissionRule.id, updateData);
        } else {
          const createData: CreatePermissionRuleRequest = values;
          await createPermissionRule(createData);
        }
        loadPermissionRules();
      } else if (activeTab === 'user-permissions') {
        if (editingUserPermission) {
          const updateData: UpdateUserPermissionRequest = {
            condition: values.condition,
            applicable_for: values.applicable_for,
            hide_descendants: values.hide_descendants,
            is_default: values.is_default,
          };
          await updateUserPermission(editingUserPermission.id, updateData);
        } else {
          const createData: CreateUserPermissionRequest = {
            user_id: values.user_id,
            doc_type: values.doc_type,
            document_name: values.document_name,
            condition: values.condition,
            applicable_for: values.applicable_for,
            hide_descendants: values.hide_descendants,
            is_default: values.is_default,
          };
          await createUserPermission(createData);
        }
        loadUserPermissions();
      } else {
        // 传统权限系统逻辑保持不变
        if (editingPermission) {
          const updateData: UpdatePermissionRequest = {
            name: values.name,
            description: values.description,
            is_menu: values.is_menu,
            is_button: values.is_button,
            is_api: values.is_api,
            menu_url: values.menu_url,
            menu_icon: values.menu_icon,
            api_path: values.api_path,
            api_method: values.api_method,
            sort_order: values.sort_order,
            is_enabled: values.is_enabled,
          };
          await updatePermission(editingPermission.id, updateData);
        } else {
          const createData: CreatePermissionRequest = {
            parent_id: values.parent_id,
            name: values.name,
            code: values.code,
            resource: values.resource,
            action: values.action,
            module: values.module,
            description: values.description,
            is_menu: values.is_menu,
            is_button: values.is_button,
            is_api: values.is_api,
            menu_url: values.menu_url,
            menu_icon: values.menu_icon,
            api_path: values.api_path,
            api_method: values.api_method,
            level: values.level,
            sort_order: values.sort_order,
            is_enabled: values.is_enabled,
          };
          await createPermission(createData);
        }
        loadPermissions();
        if (activeTab === 'tree') {
          loadPermissionTree();
        }
      }
      
      message.success(editingDocType || editingPermissionRule || editingUserPermission || editingPermission ? '更新成功' : '创建成功');
      setIsModalOpen(false);
      form.resetFields();
    } catch (error) {
      const isEditing = editingDocType || editingPermissionRule || editingUserPermission || editingPermission;
      message.error(isEditing ? '更新失败' : '创建失败');
    }
  };

  const handleGenerateCode = () => {
    const resource = form.getFieldValue('resource');
    const action = form.getFieldValue('action');
    
    if (resource && action) {
      const code = generatePermissionCode(resource, action);
      form.setFieldValue('code', code);
    }
  };
  
  const handleRefresh = () => {
    switch (activeTab) {
      case 'doctypes':
        loadDocTypes();
        break;
      case 'permission-rules':
        loadPermissionRules();
        break;
      case 'user-permissions':
        loadUserPermissions();
        break;
      case 'legacy-permissions':
        loadPermissions();
        break;
      case 'tree':
        loadPermissionTree();
        break;
    }
  };

  // ========== 表格列定义 ==========
  const docTypeColumns: TableColumnsType<DocType> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 60,
    },
    {
      title: '文档类型名称',
      dataIndex: 'name',
      width: 150,
      render: (text: string) => <code style={{ background: '#f5f5f5', padding: '2px 6px', borderRadius: '4px' }}>{text}</code>,
    },
    {
      title: '显示标签',
      dataIndex: 'label',
      width: 150,
    },
    {
      title: '所属模块',
      dataIndex: 'module',
      width: 100,
      render: (text: string) => <Tag color="cyan">{text}</Tag>,
    },
    {
      title: '特性',
      width: 200,
      render: (_, record: DocType) => (
        <Space wrap>
          {record.is_submittable && <Tag color="green">可提交</Tag>}
          {record.is_child_table && <Tag color="blue">子表</Tag>}
          {record.has_workflow && <Tag color="purple">工作流</Tag>}
          {record.track_changes && <Tag color="orange">变更追踪</Tag>}
        </Space>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      ellipsis: true,
    },
    {
      title: '版本',
      dataIndex: 'version',
      width: 60,
      render: (version: number) => <Badge count={version} color="blue" />,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      width: 160,
      render: (time: string) => time?.slice(0, 19).replace('T', ' '),
    },
    {
      title: '操作',
      width: 120,
      fixed: 'right',
      render: (_, record: DocType) => (
        <Space size="small">
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
            size="small"
          >
            编辑
          </Button>
          <Popconfirm
            title="确定删除这个文档类型吗？"
            description="删除后不可恢复"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              danger
              icon={<DeleteOutlined />}
              size="small"
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const permissionRuleColumns: TableColumnsType<PermissionRule> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 60,
    },
    {
      title: '角色',
      dataIndex: 'role_id',
      width: 120,
      render: (roleId: number) => {
        const role = roles.find(r => r.id === roleId);
        return role ? <Tag color="blue">{role.name}</Tag> : roleId;
      },
    },
    {
      title: '文档类型',
      dataIndex: 'document_type',
      width: 120,
      render: (text: string) => <code style={{ background: '#f5f5f5', padding: '2px 6px', borderRadius: '4px' }}>{text}</code>,
    },
    {
      title: '权限级别',
      dataIndex: 'permission_level',
      width: 80,
      render: (level: number) => (
        <Badge 
          count={level} 
          color={level === 0 ? 'gold' : 'purple'} 
          title={level === 0 ? '文档级' : `字段级 L${level}`}
        />
      ),
    },
    {
      title: '权限矩阵',
      width: 400,
      render: (_, record: PermissionRule) => (
        <Space wrap>
          {PERMISSION_TYPES.map(perm => (
            record[perm as keyof PermissionRule] && (
              <Tag key={perm} color="green" style={{ margin: '1px' }}>
                {PERMISSION_LABELS[perm]}
              </Tag>
            )
          ))}
          {record.if_owner && <Tag color="orange">仅拥有者</Tag>}
        </Space>
      ),
    },
    {
      title: '条件',
      width: 150,
      render: (_, record: PermissionRule) => (
        <div>
          {record.match && <div><Text type="secondary">匹配: {record.match}</Text></div>}
          {record.select_condition && <div><Text type="secondary">查询: {record.select_condition.slice(0, 20)}...</Text></div>}
        </div>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      width: 160,
      render: (time: string) => time?.slice(0, 19).replace('T', ' '),
    },
    {
      title: '操作',
      width: 120,
      fixed: 'right',
      render: (_, record: PermissionRule) => (
        <Space size="small">
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
            size="small"
          >
            编辑
          </Button>
          <Popconfirm
            title="确定删除这个权限规则吗？"
            description="删除后不可恢复"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              danger
              icon={<DeleteOutlined />}
              size="small"
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const userPermissionColumns: TableColumnsType<UserPermission> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 60,
    },
    {
      title: '用户ID',
      dataIndex: 'user_id',
      width: 80,
    },
    {
      title: '文档类型',
      dataIndex: 'doc_type',
      width: 120,
      render: (text: string) => <code style={{ background: '#f5f5f5', padding: '2px 6px', borderRadius: '4px' }}>{text}</code>,
    },
    {
      title: '文档名称',
      dataIndex: 'doc_name',
      width: 120,
      render: (text: string) => text ? <Tag color="blue">{text}</Tag> : '-',
    },
    {
      title: '权限值',
      dataIndex: 'permission_value',
      width: 150,
    },
    {
      title: '适用于',
      dataIndex: 'applicable_for',
      width: 120,
      render: (text: string) => text || '-',
    },
    {
      title: '特性',
      width: 120,
      render: (_, record: UserPermission) => (
        <Space direction="vertical" size="small">
          {record.hide_descendants && <Tag color="orange">隐藏子级</Tag>}
          {record.is_default && <Tag color="green">默认</Tag>}
        </Space>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      width: 160,
      render: (time: string) => time?.slice(0, 19).replace('T', ' '),
    },
    {
      title: '操作',
      width: 120,
      fixed: 'right',
      render: (_, record: UserPermission) => (
        <Space size="small">
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
            size="small"
          >
            编辑
          </Button>
          <Popconfirm
            title="确定删除这个用户权限吗？"
            description="删除后不可恢复"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              danger
              icon={<DeleteOutlined />}
              size="small"
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const legacyPermissionColumns: TableColumnsType<Permission> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 60,
    },
    {
      title: '权限名称',
      dataIndex: 'name',
      width: 150,
      render: (text: string, record: Permission) => (
        <Space>
          <span>{text}</span>
          {record.is_menu && <Tag color="blue" icon={<MenuOutlined />}>菜单</Tag>}
          {record.is_button && <Tag color="green" icon={<SettingOutlined />}>按钮</Tag>}
          {record.is_api && <Tag color="orange" icon={<ApiOutlined />}>接口</Tag>}
        </Space>
      ),
    },
    {
      title: '权限代码',
      dataIndex: 'code',
      width: 180,
      render: (text: string) => <code style={{ background: '#f5f5f5', padding: '2px 6px', borderRadius: '4px' }}>{text}</code>,
    },
    {
      title: '所属模块',
      dataIndex: 'module',
      width: 100,
      render: (text: string) => <Tag color="cyan">{text}</Tag>,
    },
    {
      title: '资源/操作',
      width: 150,
      render: (_, record: Permission) => (
        <Space direction="vertical" size="small">
          <span>资源: {record.resource}</span>
          <span>操作: {record.action}</span>
        </Space>
      ),
    },
    {
      title: '层级',
      dataIndex: 'level',
      width: 60,
      render: (level: number) => <Badge count={level} color="purple" />,
    },
    {
      title: '排序',
      dataIndex: 'sort_order',
      width: 60,
    },
    {
      title: '状态',
      dataIndex: 'is_enabled',
      width: 80,
      render: (enabled: boolean) => (
        <Tag color={enabled ? 'success' : 'default'}>
          {enabled ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      ellipsis: true,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      width: 160,
      render: (time: string) => time?.slice(0, 19).replace('T', ' '),
    },
    {
      title: '操作',
      width: 120,
      fixed: 'right',
      render: (_, record: Permission) => (
        <Space size="small">
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
            size="small"
          >
            编辑
          </Button>
          <Popconfirm
            title="确定删除这个权限吗？"
            description="删除后不可恢复"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              danger
              icon={<DeleteOutlined />}
              size="small"
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];
  
  // 根据当前选项卡获取对应的表格列
  const getCurrentColumns = () => {
    switch (activeTab) {
      case 'doctypes':
        return docTypeColumns;
      case 'permission-rules':
        return permissionRuleColumns;
      case 'user-permissions':
        return userPermissionColumns;
      case 'legacy-permissions':
        return legacyPermissionColumns;
      default:
        return [];
    }
  };
  
  // 根据当前选项卡获取对应的数据源
  const getCurrentDataSource = () => {
    switch (activeTab) {
      case 'doctypes':
        return docTypes;
      case 'permission-rules':
        return permissionRules;
      case 'user-permissions':
        return userPermissions;
      case 'legacy-permissions':
        return permissions;
      default:
        return [];
    }
  };
  
  // 根据当前选项卡获取对应的总数
  const getCurrentTotal = () => {
    switch (activeTab) {
      case 'doctypes':
        return docTypesTotal;
      case 'permission-rules':
        return permissionRulesTotal;
      case 'user-permissions':
        return userPermissionsTotal;
      case 'legacy-permissions':
        return total;
      default:
        return 0;
    }
  };

  const renderTreeNode = (nodes: PermissionTreeNode[]): any[] => {
    return nodes.map((node) => ({
      title: (
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <span style={{ fontWeight: 500 }}>{node.name}</span>
          <code style={{ fontSize: '12px', color: '#666', background: '#f5f5f5', padding: '1px 4px', borderRadius: '2px' }}>
            {node.code}
          </code>
          {node.is_menu && <Tag color="blue">菜单</Tag>}
          {node.is_button && <Tag color="green">按钮</Tag>}
          {node.is_api && <Tag color="orange">接口</Tag>}
          {!node.is_enabled && <Tag color="default">禁用</Tag>}
        </div>
      ),
      key: node.id,
      children: node.children ? renderTreeNode(node.children) : undefined,
    }));
  };
  
  // 渲染权限矩阵（用于权限规则表单）
  const renderPermissionMatrix = () => {
    return (
      <div>
        <Title level={5}>权限矩阵 (13种权限类型)</Title>
        <Row gutter={[8, 8]}>
          {PERMISSION_TYPES.map(perm => (
            <Col span={6} key={perm}>
              <Form.Item
                name={perm}
                valuePropName="checked"
                style={{ marginBottom: 8 }}
              >
                <Checkbox>{PERMISSION_LABELS[perm]}</Checkbox>
              </Form.Item>
            </Col>
          ))}
        </Row>
      </div>
    );
  };

  // 通用表格渲染
  const renderTableContent = (tabKey: string, buttonText: string, searchPlaceholder: string) => (
    <div>
      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Search
            placeholder={searchPlaceholder}
            allowClear
            onSearch={setSearchText}
            style={{ width: '100%' }}
          />
        </Col>
        {tabKey === 'legacy-permissions' && (
          <>
            <Col span={4}>
              <Select
                placeholder="选择模块"
                allowClear
                value={selectedModule}
                onChange={setSelectedModule}
                style={{ width: '100%' }}
              >
                {modules.map(module => (
                  <Option key={module} value={module}>{module}</Option>
                ))}
              </Select>
            </Col>
            <Col span={4}>
              <Select
                placeholder="状态筛选"
                allowClear
                value={statusFilter}
                onChange={setStatusFilter}
                style={{ width: '100%' }}
              >
                <Option value={true}>启用</Option>
                <Option value={false}>禁用</Option>
              </Select>
            </Col>
          </>
        )}
        <Col span={tabKey === 'legacy-permissions' ? 10 : 18}>
          <Space>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={handleCreate}
            >
              {buttonText}
            </Button>
            {tabKey === 'legacy-permissions' && (
              <>
                <Button
                  icon={<SyncOutlined />}
                  onClick={handleSyncApi}
                  loading={loading}
                >
                  同步API权限
                </Button>
                <Popconfirm
                  title="确定批量删除选中的项目吗？"
                  description="删除后不可恢复"
                  onConfirm={handleBatchDelete}
                  okText="确定"
                  cancelText="取消"
                >
                  <Button
                    danger
                    icon={<DeleteOutlined />}
                    disabled={selectedRowKeys.length === 0}
                  >
                    批量删除 ({selectedRowKeys.length})
                  </Button>
                </Popconfirm>
              </>
            )}
            <Button
              icon={<ReloadOutlined />}
              onClick={handleRefresh}
            >
              刷新
            </Button>
          </Space>
        </Col>
      </Row>

      <Table
        columns={getCurrentColumns()}
        dataSource={getCurrentDataSource()}
        rowKey="id"
        loading={loading}
        scroll={{ x: 1400 }}
        rowSelection={tabKey === 'legacy-permissions' ? {
          selectedRowKeys,
          onChange: setSelectedRowKeys,
          getCheckboxProps: (record: any) => ({
            disabled: record.code === 'system.admin',
          }),
        } : undefined}
        pagination={{
          current: currentPage,
          pageSize: pageSize,
          total: getCurrentTotal(),
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total, range) => `共 ${total} 条记录，显示第 ${range[0]}-${range[1]} 条`,
          onChange: setCurrentPage,
          onShowSizeChange: (_, size) => {
            setPageSize(size);
            setCurrentPage(1);
          },
        }}
      />
    </div>
  );

  const tabItems: TabsProps['items'] = [
    {
      key: 'doctypes',
      label: (
        <span>
          <DatabaseOutlined />
          文档类型管理
        </span>
      ),
      children: (
        <div>
          <Alert
            message="文档类型管理"
            description="管理系统中的所有文档类型（DocType），包括其基本属性、特性配置等。这是ERP权限系统的基础层。"
            type="info"
            showIcon
            style={{ marginBottom: 16 }}
          />
          {renderTableContent('doctypes', '新建文档类型', '搜索文档类型名称或标签')}
        </div>
      ),
    },
    {
      key: 'permission-rules',
      label: (
        <span>
          <SafetyCertificateOutlined />
          权限规则管理
        </span>
      ),
      children: (
        <div>
          <Alert
            message="权限规则管理"
            description="配置角色对特定文档类型的权限规则，支持13种权限类型和条件设置。这是ERP权限系统的核心层。"
            type="info"
            showIcon
            style={{ marginBottom: 16 }}
          />
          {renderTableContent('permission-rules', '新建权限规则', '搜索权限规则')}
        </div>
      ),
    },
    {
      key: 'user-permissions',
      label: (
        <span>
          <TeamOutlined />
          用户权限管理
        </span>
      ),
      children: (
        <div>
          <Alert
            message="用户权限管理"
            description="为特定用户配置对特定文档或文档类型的权限覆盖，实现更精细的权限控制。这是ERP权限系统的用户层。"
            type="info"
            showIcon
            style={{ marginBottom: 16 }}
          />
          {renderTableContent('user-permissions', '新建用户权限', '搜索用户权限')}
        </div>
      ),
    },
    {
      key: 'legacy-permissions',
      label: (
        <span>
          <FileTextOutlined />
          传统权限系统
        </span>
      ),
      children: (
        <div>
          <Alert
            message="传统权限系统（兼容模式）"
            description="传统的RBAC权限管理，用于向后兼容和特殊场景。建议优先使用上述ERP的权限系统。"
            type="warning"
            showIcon
            style={{ marginBottom: 16 }}
          />
          {renderTableContent('legacy-permissions', '新建权限', '搜索权限名称或代码')}
        </div>
      ),
    },
    {
      key: 'tree',
      label: (
        <span>
          <BranchesOutlined />
          权限树视图
        </span>
      ),
      children: (
        <Card>
          <div style={{ marginBottom: 16 }}>
            <Alert
              message="传统权限树视图"
              description="以树状结构展示传统权限系统的层级关系，便于理解权限体系结构。"
              type="info"
              showIcon
            />
          </div>
          <Spin spinning={treeLoading}>
            {permissionTree.length > 0 ? (
              <Tree
                treeData={renderTreeNode(permissionTree)}
                defaultExpandAll
                showLine={{ showLeafIcon: false }}
                style={{ marginTop: 16 }}
              />
            ) : (
              <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
                暂无权限数据
              </div>
            )}
          </Spin>
        </Card>
      ),
    },
  ];

  // 获取模态框标题
  const getModalTitle = () => {
    if (activeTab === 'doctypes') {
      return editingDocType ? '编辑文档类型' : '新建文档类型';
    } else if (activeTab === 'permission-rules') {
      return editingPermissionRule ? '编辑权限规则' : '新建权限规则';
    } else if (activeTab === 'user-permissions') {
      return editingUserPermission ? '编辑用户权限' : '新建用户权限';
    } else {
      return editingPermission ? '编辑权限' : '新建权限';
    }
  };

  return (
    <div>
      <Card title="权限管理系统" extra={
        <Text type="secondary">
          基于现代ERP架构的4层权限管理系统
        </Text>
      }>
        <Tabs
          activeKey={activeTab}
          onChange={(key) => {
            setActiveTab(key);
            setCurrentPage(1); // 切换tab时重置分页
            setSelectedRowKeys([]); // 清空选择
          }}
          items={tabItems}
        />
      </Card>

      <Modal
        title={getModalTitle()}
        open={isModalOpen}
        onOk={handleSubmit}
        onCancel={() => {
          setIsModalOpen(false);
          form.resetFields();
        }}
        width={activeTab === 'permission-rules' ? 1000 : 800}
        okText="确定"
        cancelText="取消"
      >
        <Form
          form={form}
          layout="vertical"
          style={{ marginTop: 16 }}
        >
          {/* DocType 表单 */}
          {activeTab === 'doctypes' && (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label="文档类型名称"
                    name="name"
                    rules={[{ required: true, message: '请输入文档类型名称' }]}
                  >
                    <Input 
                      placeholder="如：User" 
                      disabled={!!editingDocType}
                      style={{ fontFamily: 'monospace' }}
                    />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label="显示标签"
                    name="label"
                    rules={[{ required: true, message: '请输入显示标签' }]}
                  >
                    <Input placeholder="如：用户" />
                  </Form.Item>
                </Col>
              </Row>
              
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label="所属模块"
                    name="module"
                    rules={[{ required: true, message: '请输入所属模块' }]}
                  >
                    <Input placeholder="如：Core" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label="最大附件数"
                    name="max_attachments"
                  >
                    <Input type="number" placeholder="默认：10" />
                  </Form.Item>
                </Col>
              </Row>
              
              <Form.Item
                label="描述"
                name="description"
              >
                <TextArea rows={3} placeholder="文档类型的详细描述" />
              </Form.Item>
              
              <Divider orientation="left">特性配置</Divider>
              <Row gutter={16}>
                <Col span={6}>
                  <Form.Item
                    label="可提交"
                    name="is_submittable"
                    valuePropName="checked"
                  >
                    <Switch />
                  </Form.Item>
                </Col>
                <Col span={6}>
                  <Form.Item
                    label="子表"
                    name="is_child_table"
                    valuePropName="checked"
                  >
                    <Switch />
                  </Form.Item>
                </Col>
                <Col span={6}>
                  <Form.Item
                    label="工作流"
                    name="has_workflow"
                    valuePropName="checked"
                  >
                    <Switch />
                  </Form.Item>
                </Col>
                <Col span={6}>
                  <Form.Item
                    label="变更追踪"
                    name="track_changes"
                    valuePropName="checked"
                  >
                    <Switch />
                  </Form.Item>
                </Col>
              </Row>
            </>
          )}
          
          {/* 权限规则表单 */}
          {activeTab === 'permission-rules' && (
            <>
              <Row gutter={16}>
                <Col span={8}>
                  <Form.Item
                    label="角色"
                    name="role_id"
                    rules={[{ required: true, message: '请选择角色' }]}
                  >
                    <Select placeholder="请选择角色" disabled={!!editingPermissionRule}>
                      {roles.map(role => (
                        <Option key={role.id} value={role.id}>
                          {role.name}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item
                    label="文档类型"
                    name="document_type"
                    rules={[{ required: true, message: '请输入文档类型' }]}
                  >
                    <Select placeholder="请选择文档类型" disabled={!!editingPermissionRule}>
                      {docTypes.map(dt => (
                        <Option key={dt.id} value={dt.name}>
                          {dt.label} ({dt.name})
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item
                    label="权限级别"
                    name="permission_level"
                  >
                    <Select placeholder="请选择权限级别">
                      <Option value={0}>0 - 文档级</Option>
                      {Array.from({length: 10}, (_, i) => i + 1).map(level => (
                        <Option key={level} value={level}>L{level} - 字段级</Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
              </Row>
              
              {renderPermissionMatrix()}
              
              <Divider orientation="left">条件设置</Divider>
              <Row gutter={16}>
                <Col span={8}>
                  <Form.Item
                    label="仅拥有者"
                    name="if_owner"
                    valuePropName="checked"
                  >
                    <Switch />
                  </Form.Item>
                </Col>
                <Col span={16}>
                  <Form.Item
                    label="匹配条件"
                    name="match"
                  >
                    <Input placeholder="如：company" />
                  </Form.Item>
                </Col>
              </Row>
              
              <Form.Item
                label="查询条件"
                name="select_condition"
              >
                <TextArea rows={2} placeholder="SQL WHERE 条件" />
              </Form.Item>
            </>
          )}
          
          {/* 用户权限表单 */}
          {activeTab === 'user-permissions' && (
            <>
              <Row gutter={16}>
                <Col span={8}>
                  <Form.Item
                    label="用户ID"
                    name="user_id"
                    rules={[{ required: true, message: '请输入用户ID' }]}
                  >
                    <Input type="number" placeholder="用户ID" disabled={!!editingUserPermission} />
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item
                    label="文档类型"
                    name="doc_type"
                    rules={[{ required: true, message: '请选择文档类型' }]}
                  >
                    <Select placeholder="请选择文档类型" disabled={!!editingUserPermission}>
                      {docTypes.map(dt => (
                        <Option key={dt.id} value={dt.name}>
                          {dt.label} ({dt.name})
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item
                    label="文档名称"
                    name="document_name"
                  >
                    <Input placeholder="特定文档名称（可选）" />
                  </Form.Item>
                </Col>
              </Row>
              
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label="权限值"
                    name="condition"
                    rules={[{ required: true, message: '请输入权限值' }]}
                  >
                    <Input placeholder="如：特定公司名称" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label="适用于"
                    name="applicable_for"
                  >
                    <Input placeholder="适用的文档类型（可选）" />
                  </Form.Item>
                </Col>
              </Row>
              
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label="隐藏子级"
                    name="hide_descendants"
                    valuePropName="checked"
                  >
                    <Switch />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label="设为默认"
                    name="is_default"
                    valuePropName="checked"
                  >
                    <Switch />
                  </Form.Item>
                </Col>
              </Row>
            </>
          )}
          
          {/* 传统权限表单 */}
          {activeTab === 'legacy-permissions' && (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label="权限名称"
                    name="name"
                    rules={[{ required: true, message: '请输入权限名称' }]}
                  >
                    <Input placeholder="如：用户管理" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label="权限代码"
                    name="code"
                    rules={[
                      { required: true, message: '请输入权限代码' },
                      { pattern: /^[a-z][a-z0-9_]*\.[a-z][a-z0-9_]*$/, message: '格式：resource.action，如 user.create' }
                    ]}
                  >
                    <Input placeholder="如：user.create" />
                  </Form.Item>
                </Col>
              </Row>

              <Row gutter={16}>
                <Col span={8}>
                  <Form.Item
                    label="资源"
                    name="resource"
                    rules={[{ required: true, message: '请输入资源名称' }]}
                  >
                    <Input placeholder="如：user" onChange={handleGenerateCode} />
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item
                    label="操作"
                    name="action"
                    rules={[{ required: true, message: '请输入操作名称' }]}
                  >
                    <Input placeholder="如：create" onChange={handleGenerateCode} />
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item
                    label="所属模块"
                    name="module"
                    rules={[{ required: true, message: '请选择所属模块' }]}
                  >
                    <Select placeholder="请选择模块">
                      {modules.map(module => (
                        <Option key={module} value={module}>{module}</Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
              </Row>

              <Form.Item
                label="权限描述"
                name="description"
              >
                <TextArea rows={3} placeholder="权限的详细描述" />
              </Form.Item>

              <Divider orientation="left">权限类型</Divider>
              <Row gutter={16}>
                <Col span={8}>
                  <Form.Item
                    label="菜单权限"
                    name="is_menu"
                    valuePropName="checked"
                  >
                    <Switch />
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item
                    label="按钮权限"
                    name="is_button"
                    valuePropName="checked"
                  >
                    <Switch />
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item
                    label="接口权限"
                    name="is_api"
                    valuePropName="checked"
                  >
                    <Switch />
                  </Form.Item>
                </Col>
              </Row>

              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label="菜单路径"
                    name="menu_url"
                  >
                    <Input placeholder="如：/users" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label="菜单图标"
                    name="menu_icon"
                  >
                    <Input placeholder="如：UserOutlined" />
                  </Form.Item>
                </Col>
              </Row>

              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label="API路径"
                    name="api_path"
                  >
                    <Input placeholder="如：/api/v1/users" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label="API方法"
                    name="api_method"
                  >
                    <Select placeholder="请选择HTTP方法">
                      <Option value="GET">GET</Option>
                      <Option value="POST">POST</Option>
                      <Option value="PUT">PUT</Option>
                      <Option value="DELETE">DELETE</Option>
                      <Option value="PATCH">PATCH</Option>
                    </Select>
                  </Form.Item>
                </Col>
              </Row>

              <Row gutter={16}>
                <Col span={8}>
                  <Form.Item
                    label="权限层级"
                    name="level"
                    rules={[{ required: true, message: '请输入权限层级' }]}
                  >
                    <Select placeholder="请选择层级">
                      <Option value={1}>一级权限</Option>
                      <Option value={2}>二级权限</Option>
                      <Option value={3}>三级权限</Option>
                      <Option value={4}>四级权限</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item
                    label="排序顺序"
                    name="sort_order"
                  >
                    <Input type="number" placeholder="数字越小排序越靠前" />
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item
                    label="启用状态"
                    name="is_enabled"
                    valuePropName="checked"
                  >
                    <Switch />
                  </Form.Item>
                </Col>
              </Row>
            </>
          )}
        </Form>
      </Modal>
    </div>
  );
};

export default PermissionManagement;