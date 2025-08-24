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
} from '@ant-design/icons';
import type { TableColumnsType, TabsProps } from 'antd';
import { Permission } from '../types/auth';
import {
  getPermissionsList,
  getPermissionTree,
  createPermission,
  updatePermission,
  deletePermission,
  batchDeletePermissions,
  getModules,
  // checkPermissionCode,
  syncApiPermissions,
  generatePermissionCode,
  // formatPermissionTreeForAntd,
  // type PermissionsListResponse,
  type PermissionsQueryParams,
  type CreatePermissionRequest,
  type UpdatePermissionRequest,
  type PermissionTreeNode,
} from '../services/permissionService';

const { Search } = Input;
const { TextArea } = Input;
const { Option } = Select;

const PermissionManagement: React.FC = () => {
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [permissionTree, setPermissionTree] = useState<PermissionTreeNode[]>([]);
  const [modules, setModules] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const [treeLoading, setTreeLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchText, setSearchText] = useState('');
  const [selectedModule, setSelectedModule] = useState<string>();
  const [statusFilter, setStatusFilter] = useState<boolean>();
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingPermission, setEditingPermission] = useState<Permission | null>(null);
  const [form] = Form.useForm();
  const [activeTab, setActiveTab] = useState('list');

  useEffect(() => {
    loadPermissions();
    loadModules();
  }, [currentPage, pageSize, searchText, selectedModule, statusFilter]);

  useEffect(() => {
    if (activeTab === 'tree') {
      loadPermissionTree();
    }
  }, [activeTab]);

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

  const handleCreate = () => {
    setEditingPermission(null);
    setIsModalOpen(true);
    form.resetFields();
    form.setFieldsValue({
      is_enabled: true,
      is_menu: false,
      is_button: false,
      is_api: true,
      level: 1,
      sort_order: 0,
    });
  };

  const handleEdit = (record: Permission) => {
    setEditingPermission(record);
    setIsModalOpen(true);
    form.setFieldsValue({
      ...record,
      parent_id: record.parent_id || undefined,
    });
  };

  const handleDelete = async (id: number) => {
    try {
      await deletePermission(id);
      message.success('删除成功');
      loadPermissions();
      if (activeTab === 'tree') {
        loadPermissionTree();
      }
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleBatchDelete = async () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要删除的权限');
      return;
    }

    try {
      await batchDeletePermissions(selectedRowKeys as number[]);
      message.success('批量删除成功');
      setSelectedRowKeys([]);
      loadPermissions();
      if (activeTab === 'tree') {
        loadPermissionTree();
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
        message.success('更新成功');
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
        message.success('创建成功');
      }
      
      setIsModalOpen(false);
      form.resetFields();
      loadPermissions();
      if (activeTab === 'tree') {
        loadPermissionTree();
      }
    } catch (error) {
      message.error(editingPermission ? '更新失败' : '创建失败');
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

  const columns: TableColumnsType<Permission> = [
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

  const tabItems: TabsProps['items'] = [
    {
      key: 'list',
      label: '权限列表',
      children: (
        <div>
          <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
            <Col span={6}>
              <Search
                placeholder="搜索权限名称或代码"
                allowClear
                onSearch={setSearchText}
                style={{ width: '100%' }}
              />
            </Col>
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
            <Col span={10}>
              <Space>
                <Button
                  type="primary"
                  icon={<PlusOutlined />}
                  onClick={handleCreate}
                >
                  新建权限
                </Button>
                <Button
                  icon={<SyncOutlined />}
                  onClick={handleSyncApi}
                  loading={loading}
                >
                  同步API权限
                </Button>
                <Popconfirm
                  title="确定批量删除选中的权限吗？"
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
                <Button
                  icon={<ReloadOutlined />}
                  onClick={loadPermissions}
                >
                  刷新
                </Button>
              </Space>
            </Col>
          </Row>

          <Table
            columns={columns}
            dataSource={permissions}
            rowKey="id"
            loading={loading}
            scroll={{ x: 1400 }}
            rowSelection={{
              selectedRowKeys,
              onChange: setSelectedRowKeys,
              getCheckboxProps: (record: Permission) => ({
                disabled: record.code === 'system.admin',
              }),
            }}
            pagination={{
              current: currentPage,
              pageSize: pageSize,
              total: total,
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
      ),
    },
    {
      key: 'tree',
      label: (
        <span>
          <BranchesOutlined />
          权限树
        </span>
      ),
      children: (
        <Card>
          <div style={{ marginBottom: 16 }}>
            <Alert
              message="权限树视图"
              description="以树状结构展示所有权限的层级关系，便于理解权限体系结构。"
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

  return (
    <div>
      <Card title="权限管理">
        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          items={tabItems}
        />
      </Card>

      <Modal
        title={editingPermission ? '编辑权限' : '新建权限'}
        open={isModalOpen}
        onOk={handleSubmit}
        onCancel={() => {
          setIsModalOpen(false);
          form.resetFields();
        }}
        width={800}
        okText="确定"
        cancelText="取消"
      >
        <Form
          form={form}
          layout="vertical"
          style={{ marginTop: 16 }}
        >
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
        </Form>
      </Modal>
    </div>
  );
};

export default PermissionManagement;