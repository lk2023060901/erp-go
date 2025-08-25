/**
 * 角色管理页面
 */

import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Card,
  Table,
  Button,
  Input,
  Space,
  Modal,
  Form,
  Switch,
  Tag,
  Dropdown,
  // Menu,
  Popconfirm,
  message,
  Tooltip,
  Typography,
  Row,
  Col,
  Divider,
  Tree,
  Spin,
  InputNumber
} from 'antd';
import {
  PlusOutlined,
  // SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  MoreOutlined,
  ReloadOutlined,
  CopyOutlined,
  FilterOutlined,
  SortAscendingOutlined,
  SortDescendingOutlined,
  DownOutlined,
  // SwapOutlined,
  SafetyCertificateOutlined,
  ExportOutlined,
  SettingOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { MenuProps } from 'antd';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';

dayjs.extend(relativeTime);
import { 
  getRolesList, 
  createRole, 
  updateRole, 
  deleteRole,
  assignRolePermissions,
  getRolePermissions,
  copyRole,
  batchDeleteRoles
} from '../services/roleService';
// TODO: Replace with ERP permission system
// import {
//   getPermissionTree,
//   formatPermissionTreeForAntd
// } from '../services/permissionService';
import type { 
  // RolesListResponse, 
  CreateRoleRequest, 
  UpdateRoleRequest 
} from '../services/roleService';
// TODO: Replace with ERP permission system types
// import type { PermissionTreeNode } from '../services/permissionService';
import type { Role } from '../types/auth';

const { Search } = Input;
const { Text } = Typography;

export function RoleManagement() {
  const navigate = useNavigate();
  
  // 状态管理
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchText, setSearchText] = useState('');
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  // 模态框状态
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [permissionModalVisible, setPermissionModalVisible] = useState(false);
  const [copyModalVisible, setCopyModalVisible] = useState(false);
  
  // 编辑状态
  const [editingRole, setEditingRole] = useState<Role | null>(null);
  // TODO: Replace with ERP permission system
  // const [permissionTree, setPermissionTree] = useState<PermissionTreeNode[]>([]);
  const [checkedPermissions, setCheckedPermissions] = useState<number[]>([]);
  const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([]);
  const [permissionLoading, setPermissionLoading] = useState(false);
  
  // 排序状态
  const [sortField, setSortField] = useState<string>('name');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');

  // 表单实例
  const [createForm] = Form.useForm();
  const [copyForm] = Form.useForm();

  // 初始化数据
  useEffect(() => {
    loadRoles();
    // TODO: Replace with ERP permission system
    // loadPermissionTree();
  }, [page, pageSize, searchText, sortField, sortOrder]);

  // 加载角色列表
  const loadRoles = async () => {
    try {
      setLoading(true);
      const response = await getRolesList({
        page,
        size: pageSize,
        search: searchText || undefined,
        sortField: sortField,
        sortOrder: sortOrder,
      });
      setRoles(response.roles);
      setTotal(response.total);
    } catch (error) {
      message.error('加载角色列表失败');
    } finally {
      setLoading(false);
    }
  };

  // TODO: Replace with ERP permission system
  // 加载权限树
  // const loadPermissionTree = async () => {
  //   try {
  //     const tree = await getPermissionTree();
  //     setPermissionTree(tree);
  //   } catch (error) {
  //     message.error('加载权限树失败');
  //   }
  // };

  // 搜索处理
  const handleSearch = (value: string) => {
    setSearchText(value);
    setPage(1);
  };

  // 创建角色
  const handleCreateRole = async (values: CreateRoleRequest) => {
    try {
      await createRole(values);
      message.success('创建角色成功');
      setCreateModalVisible(false);
      createForm.resetFields();
      loadRoles();
    } catch (error) {
      message.error('创建角色失败');
    }
  };

  // 编辑角色

  // 删除角色
  const handleDeleteRole = async (role: Role) => {
    try {
      await deleteRole(role.id);
      message.success('删除角色成功');
      loadRoles();
    } catch (error) {
      message.error('删除角色失败');
    }
  };

  // 批量删除角色
  const handleBatchDelete = async () => {
    try {
      await batchDeleteRoles(selectedRowKeys as number[]);
      message.success(`成功删除 ${selectedRowKeys.length} 个角色`);
      setSelectedRowKeys([]);
      loadRoles();
    } catch (error) {
      message.error('批量删除失败');
    }
  };

  // 复制角色
  const handleCopyRole = async (values: { name: string }) => {
    if (!editingRole) return;
    
    try {
      await copyRole(editingRole.id, values.name);
      message.success('复制角色成功');
      setCopyModalVisible(false);
      setEditingRole(null);
      copyForm.resetFields();
      loadRoles();
    } catch (error) {
      message.error('复制角色失败');
    }
  };

  // 分配权限
  const handleAssignPermissions = async () => {
    if (!editingRole) return;
    
    try {
      await assignRolePermissions(editingRole.id, checkedPermissions);
      message.success('分配权限成功');
      setPermissionModalVisible(false);
      setEditingRole(null);
      setCheckedPermissions([]);
    } catch (error) {
      message.error('分配权限失败');
    }
  };


  // 打开复制模态框
  const openCopyModal = (role: Role) => {
    setEditingRole(role);
    copyForm.setFieldsValue({
      name: `${role.name}_副本`,
    });
    setCopyModalVisible(true);
  };

  // 打开传统权限分配模态框（保留向后兼容）
  const openLegacyPermissionModal = async (role: Role) => {
    try {
      setPermissionLoading(true);
      setEditingRole(role);
      const permissions = await getRolePermissions(role.id);
      setCheckedPermissions(permissions);
      setPermissionModalVisible(true);
      
      // 展开所有节点
      const getAllKeys = (nodes: PermissionTreeNode[]): React.Key[] => {
        let keys: React.Key[] = [];
        nodes.forEach(node => {
          keys.push(node.id);
          if (node.children) {
            keys = keys.concat(getAllKeys(node.children));
          }
        });
        return keys;
      };
      // TODO: Replace with ERP permission system
      // setExpandedKeys(getAllKeys(permissionTree));
    } catch (error) {
      message.error('获取角色权限失败');
    } finally {
      setPermissionLoading(false);
    }
  };

  // 排序菜单配置
  const getSortMenuItems = () => {
    const sortOptions = [
      { key: 'created_at', label: '创建时间', icon: '📅' },
      { key: 'name', label: '角色名称', icon: '👤' },
      { key: 'user_count', label: '用户数量', icon: '#️⃣' },
      { key: 'is_system_role', label: '角色类型', icon: '⚙️' },
    ];

    return sortOptions.map(option => ({
      key: option.key,
      label: (
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <span>{option.icon}</span>
          <span>{option.label}</span>
        </div>
      ),
    }));
  };

  const handleSortChange = ({ key }: { key: string }) => {
    setSortField(key);
    const sortOptions = [
      { key: 'created_at', label: '创建时间' },
      { key: 'name', label: '角色名称' },
      { key: 'user_count', label: '用户数量' },
      { key: 'is_system_role', label: '角色类型' },
    ];
    const selectedOption = sortOptions.find(item => item.key === key);
    message.success(`已切换到按 "${selectedOption?.label}" 排序`);
  };

  // 切换排序方向
  const handleSortOrderChange = () => {
    const newOrder = sortOrder === 'asc' ? 'desc' : 'asc';
    setSortOrder(newOrder);
    message.success(`已切换到${newOrder === 'asc' ? '升序' : '降序'}排序`);
  };

  // 获取当前排序显示文本
  const getCurrentSortLabel = () => {
    const sortOptions = [
      { key: 'created_at', label: '创建时间' },
      { key: 'name', label: '角色名称' },
      { key: 'user_count', label: '用户数量' },
      { key: 'is_system_role', label: '角色类型' },
    ];
    const currentOption = sortOptions.find(item => item.key === sortField);
    return currentOption?.label || '角色名称';
  };

  // 表格列定义
  const columns: ColumnsType<Role> = [
    {
      title: '角色名称',
      dataIndex: 'name',
      key: 'name',
      render: (name, record) => (
        <div>
          <div style={{ fontWeight: 500 }}>{name}</div>
          {record.description && (
            <div>
              <Text type="secondary" style={{ fontSize: '12px' }}>
                {record.description.length > 50 ? `${record.description.substring(0, 50)}...` : record.description}
              </Text>
            </div>
          )}
        </div>
      ),
    },
    {
      title: '类型',
      dataIndex: 'is_system_role',
      key: 'is_system_role',
      render: (isSystem: boolean, record) => (
        <div>
          <Tag color={isSystem ? 'orange' : 'blue'}>
            {isSystem ? '系统角色' : '自定义角色'}
          </Tag>
          {record.is_custom && (
            <div style={{ marginTop: 4 }}>
              <Tag color="purple" style={{ fontSize: '11px' }}>自定义</Tag>
            </div>
          )}
        </div>
      ),
    },
    {
      title: '权限设置',
      key: 'permissions',
      render: (_, record) => (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '2px' }}>
          {record.desk_access && (
            <Tag color="green" style={{ fontSize: '11px' }}>桌面访问</Tag>
          )}
          {record.require_two_factor && (
            <Tag color="red" style={{ fontSize: '11px' }}>2FA</Tag>
          )}
          {record.restrict_to_domain && (
            <Tooltip title={`域限制: ${record.restrict_to_domain}`}>
              <Tag color="orange" style={{ fontSize: '11px' }}>域限制</Tag>
            </Tooltip>
          )}
        </div>
      ),
    },
    {
      title: '状态',
      dataIndex: 'is_enabled',
      key: 'is_enabled',
      render: (isEnabled: boolean) => (
        <Tag color={isEnabled ? 'green' : 'red'}>
          {isEnabled ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (time: string) => (
        <Tooltip title={dayjs(time).format('YYYY-MM-DD HH:mm:ss')}>
          {dayjs(time).fromNow()}
        </Tooltip>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record) => {
        const items: MenuProps['items'] = [
          {
            key: 'edit',
            label: '编辑角色',
            icon: <EditOutlined />,
            onClick: () => navigate(`/roles/edit/${record.id}`),
          },
          {
            key: 'permissions',
            label: '配置权限',
            icon: <SafetyCertificateOutlined />,
            onClick: () => navigate(`/roles/${record.id}/permissions`),
          },
          {
            key: 'legacy-permissions',
            label: '传统权限',
            icon: <SettingOutlined />,
            onClick: () => openLegacyPermissionModal(record),
          },
          {
            key: 'copy',
            label: '复制角色',
            icon: <CopyOutlined />,
            onClick: () => openCopyModal(record),
          },
          {
            type: 'divider',
          },
          {
            key: 'delete',
            label: '删除角色',
            icon: <DeleteOutlined />,
            danger: true,
            disabled: record.is_system_role,
            onClick: () => {
              if (record.is_system_role) {
                message.warning('系统角色不能删除');
                return;
              }
              Modal.confirm({
                title: '确认删除',
                content: `确定要删除角色 "${record.name}" 吗？`,
                onOk: () => handleDeleteRole(record),
              });
            },
          },
        ];

        return (
          <Dropdown menu={{ items }} trigger={['click']}>
            <Button type="text" icon={<MoreOutlined />} />
          </Dropdown>
        );
      },
    },
  ];

  // 行选择配置
  const rowSelection = {
    selectedRowKeys,
    onChange: setSelectedRowKeys,
    getCheckboxProps: (record: Role) => ({
      disabled: record.is_system_role,
    }),
  };

  return (
    <div style={{ padding: '24px' }}>
      {/* 页面头部 */}
      <Card style={{ marginBottom: '16px' }}>
        <Row justify="space-between" align="middle">
          <Col>
            <Space>
              <Search
                placeholder="搜索角色名称、代码、描述"
                allowClear
                onSearch={handleSearch}
                style={{ width: 300 }}
              />
              <Button 
                icon={<ReloadOutlined />} 
                onClick={loadRoles}
                loading={loading}
              >
                刷新
              </Button>
            </Space>
          </Col>
          <Col>
            <Space>
              <Button 
                type="primary"
                icon={<PlusOutlined />}
                onClick={() => setCreateModalVisible(true)}
              >
                创建角色
              </Button>
              <Button 
                icon={<FilterOutlined />}
                onClick={() => message.info('过滤器功能开发中')}
              >
                过滤器
              </Button>
              <Space.Compact>
                {/* 左侧排序方向按钮 */}
                <Button 
                  icon={sortOrder === 'asc' ? <SortAscendingOutlined /> : <SortDescendingOutlined />}
                  onClick={handleSortOrderChange}
                  style={{
                    backgroundColor: '#fff',
                    borderColor: '#d1d5db',
                    color: sortOrder === 'asc' ? '#1890ff' : '#ff4d4f'
                  }}
                  title={sortOrder === 'asc' ? '当前升序，点击切换为降序' : '当前降序，点击切换为升序'}
                />
                {/* 右侧排序字段下拉框 */}
                <Dropdown
                  menu={{
                    items: getSortMenuItems(),
                    onClick: handleSortChange,
                  }}
                  trigger={['click']}
                >
                  <Button 
                    style={{
                      backgroundColor: '#fff',
                      borderColor: '#d1d5db',
                      color: '#374151',
                      display: 'flex',
                      alignItems: 'center',
                      gap: '6px'
                    }}
                  >
                    <span>{getCurrentSortLabel()}</span>
                    <DownOutlined style={{ fontSize: '10px' }} />
                  </Button>
                </Dropdown>
              </Space.Compact>
              <Button 
                icon={<ExportOutlined />}
                onClick={() => message.info('导出功能开发中')}
              >
                导出
              </Button>
              {selectedRowKeys.length > 0 && (
                <Popconfirm
                  title={`确定要删除选中的 ${selectedRowKeys.length} 个角色吗？`}
                  onConfirm={handleBatchDelete}
                  okText="确定"
                  cancelText="取消"
                >
                  <Button danger icon={<DeleteOutlined />}>
                    批量删除 ({selectedRowKeys.length})
                  </Button>
                </Popconfirm>
              )}
            </Space>
          </Col>
        </Row>
      </Card>

      {/* 角色表格 */}
      <Card>
        <Table
          columns={columns}
          dataSource={roles}
          rowKey="id"
          loading={loading}
          rowSelection={rowSelection}
          pagination={{
            current: page,
            pageSize: pageSize,
            total: total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total, range) => 
              `第 ${range[0]}-${range[1]} 条/共 ${total} 条`,
            onChange: (page, size) => {
              setPage(page);
              setPageSize(size || 10);
            },
          }}
        />
      </Card>

      {/* 创建角色模态框 */}
      <Modal
        title="创建角色"
        open={createModalVisible}
        onCancel={() => {
          setCreateModalVisible(false);
          createForm.resetFields();
        }}
        footer={null}
        width={400}
      >
        <Form
          form={createForm}
          layout="vertical"
          onFinish={handleCreateRole}
        >
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

          <div style={{ 
            marginTop: '24px', 
            display: 'flex', 
            justifyContent: 'space-between', 
            alignItems: 'center' 
          }}>
            <Button 
              type="link" 
              onClick={() => {
                const formData = createForm.getFieldsValue();
                navigate('/roles/create', { 
                  state: { name: formData.name } 
                });
                setCreateModalVisible(false);
                createForm.resetFields();
              }}
              style={{ padding: 0 }}
            >
              编辑完整信息
            </Button>
            
            <Space>
              <Button onClick={() => {
                setCreateModalVisible(false);
                createForm.resetFields();
              }}>
                取消
              </Button>
              <Button type="primary" htmlType="submit">
                创建
              </Button>
            </Space>
          </div>
        </Form>
      </Modal>


      {/* 权限分配模态框 */}
      <Modal
        title={`分配权限 - ${editingRole?.name}`}
        open={permissionModalVisible}
        onCancel={() => {
          setPermissionModalVisible(false);
          setEditingRole(null);
          setCheckedPermissions([]);
        }}
        onOk={handleAssignPermissions}
        width={600}
        styles={{ body: { maxHeight: '500px', overflow: 'auto' } }}
      >
        {permissionLoading ? (
          <div style={{ textAlign: 'center', padding: '50px 0' }}>
            <Spin size="large" />
            <div style={{ marginTop: 16 }}>加载权限数据中...</div>
          </div>
        ) : (
          <div style={{ textAlign: 'center', padding: '50px', color: '#999' }}>
            <p>权限分配功能正在开发中</p>
            <p>TODO: 使用ERP权限系统实现权限树和分配功能</p>
          </div>
        )}
      </Modal>

      {/* 复制角色模态框 */}
      <Modal
        title={`复制角色 - ${editingRole?.name}`}
        open={copyModalVisible}
        onCancel={() => {
          setCopyModalVisible(false);
          setEditingRole(null);
          copyForm.resetFields();
        }}
        footer={null}
        width={400}
      >
        <Form
          form={copyForm}
          layout="vertical"
          onFinish={handleCopyRole}
        >
          <Form.Item
            name="name"
            label="新角色名称"
            rules={[
              { required: true, message: '请输入新角色名称' },
              { min: 2, message: '角色名称至少2个字符' },
              { max: 50, message: '角色名称最多50个字符' },
            ]}
          >
            <Input placeholder="请输入新角色名称" />
          </Form.Item>

          <div style={{ 
            marginTop: '24px', 
            display: 'flex', 
            justifyContent: 'flex-end' 
          }}>
            <Space>
              <Button onClick={() => {
                setCopyModalVisible(false);
                setEditingRole(null);
                copyForm.resetFields();
              }}>
                取消
              </Button>
              <Button type="primary" htmlType="submit">
                复制
              </Button>
            </Space>
          </div>
        </Form>
      </Modal>

    </div>
  );
}

export default RoleManagement;