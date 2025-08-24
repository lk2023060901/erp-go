import React, { useState, useEffect } from 'react';
import {
  Card,
  Tree,
  Button,
  Space,
  Modal,
  Form,
  Input,
  Select,
  Switch,
  message,
  Popconfirm,
  Tag,
  Row,
  Col,
  Tabs,
  Alert,
  Spin,
  Badge,
  // Divider,
  Table,
  Typography,
  // Tooltip,
  // Dropdown,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  ApartmentOutlined,
  TeamOutlined,
  UserOutlined,
  // MoreOutlined,
  // CopyOutlined,
  // DownOutlined,
  // SearchOutlined,
} from '@ant-design/icons';
import type { TableColumnsType, TabsProps } from 'antd';
// import type { MenuProps } from 'antd';
import type { DataNode } from 'antd/es/tree';

const { Search } = Input;
const { TextArea } = Input;
const { Option } = Select;
const { Text } = Typography;

// 组织类型枚举
enum OrganizationType {
  COMPANY = 'company',
  DEPARTMENT = 'department',
  TEAM = 'team',
  GROUP = 'group'
}

// 组织数据结构
interface Organization {
  id: number;
  parent_id?: number;
  name: string;
  code: string;
  type: OrganizationType;
  description?: string;
  leader_id?: number;
  leader_name?: string;
  employee_count: number;
  level: number;
  sort_order: number;
  is_enabled: boolean;
  created_at: string;
  updated_at: string;
  children?: Organization[];
}

// 组织树节点
interface OrganizationTreeNode extends DataNode {
  id: number;
  name: string;
  type: OrganizationType;
  employee_count: number;
  is_enabled: boolean;
  children?: OrganizationTreeNode[];
}

// 创建组织请求
interface CreateOrganizationRequest {
  parent_id?: number;
  name: string;
  code: string;
  type: OrganizationType;
  description?: string;
  leader_id?: number;
  sort_order?: number;
  is_enabled: boolean;
}

// 更新组织请求
interface UpdateOrganizationRequest {
  name?: string;
  description?: string;
  leader_id?: number;
  sort_order?: number;
  is_enabled?: boolean;
}

// 组织列表响应
interface OrganizationListResponse {
  organizations: Organization[];
  total: number;
  page: number;
  size: number;
}

// 模拟API服务
const organizationService = {
  // 获取组织列表
  async getOrganizations(_params?: any): Promise<OrganizationListResponse> {
    await new Promise(resolve => setTimeout(resolve, 800));
    
    const mockOrganizations: Organization[] = [
      {
        id: 1,
        name: '总公司',
        code: 'COMPANY_ROOT',
        type: OrganizationType.COMPANY,
        description: '公司总部',
        employee_count: 150,
        level: 1,
        sort_order: 1,
        is_enabled: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      },
      {
        id: 2,
        parent_id: 1,
        name: '技术部',
        code: 'TECH_DEPT',
        type: OrganizationType.DEPARTMENT,
        description: '负责技术研发',
        leader_id: 2,
        leader_name: '张技术',
        employee_count: 50,
        level: 2,
        sort_order: 1,
        is_enabled: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      },
      {
        id: 3,
        parent_id: 1,
        name: '销售部',
        code: 'SALES_DEPT',
        type: OrganizationType.DEPARTMENT,
        description: '负责销售业务',
        leader_id: 3,
        leader_name: '李销售',
        employee_count: 30,
        level: 2,
        sort_order: 2,
        is_enabled: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      },
      {
        id: 4,
        parent_id: 2,
        name: '前端团队',
        code: 'FRONTEND_TEAM',
        type: OrganizationType.TEAM,
        description: '前端开发团队',
        leader_id: 4,
        leader_name: '王前端',
        employee_count: 15,
        level: 3,
        sort_order: 1,
        is_enabled: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      },
      {
        id: 5,
        parent_id: 2,
        name: '后端团队',
        code: 'BACKEND_TEAM',
        type: OrganizationType.TEAM,
        description: '后端开发团队',
        leader_id: 5,
        leader_name: '赵后端',
        employee_count: 20,
        level: 3,
        sort_order: 2,
        is_enabled: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      },
    ];

    return {
      organizations: mockOrganizations,
      total: mockOrganizations.length,
      page: 1,
      size: 100,
    };
  },

  // 获取组织树
  async getOrganizationTree(): Promise<Organization[]> {
    const response = await this.getOrganizations();
    return this.buildTree(response.organizations);
  },

  // 构建树结构
  buildTree(organizations: Organization[]): Organization[] {
    const map = new Map<number, Organization>();
    const roots: Organization[] = [];

    organizations.forEach(org => {
      map.set(org.id, { ...org, children: [] });
    });

    organizations.forEach(org => {
      const node = map.get(org.id)!;
      if (org.parent_id && map.has(org.parent_id)) {
        const parent = map.get(org.parent_id)!;
        if (!parent.children) parent.children = [];
        parent.children.push(node);
      } else {
        roots.push(node);
      }
    });

    return roots;
  },

  // 创建组织
  async createOrganization(data: CreateOrganizationRequest): Promise<Organization> {
    await new Promise(resolve => setTimeout(resolve, 500));
    
    const newOrg: Organization = {
      id: Date.now(),
      ...data,
      sort_order: data.sort_order || 0,
      employee_count: 0,
      level: data.parent_id ? 2 : 1,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };

    return newOrg;
  },

  // 更新组织
  async updateOrganization(id: number, data: UpdateOrganizationRequest): Promise<Organization> {
    await new Promise(resolve => setTimeout(resolve, 500));
    
    return {
      id,
      name: data.name || '更新的组织',
      code: 'UPDATED_ORG',
      type: OrganizationType.DEPARTMENT,
      employee_count: 10,
      level: 2,
      sort_order: data.sort_order || 0,
      is_enabled: data.is_enabled ?? true,
      created_at: '2024-01-01T00:00:00Z',
      updated_at: new Date().toISOString(),
    };
  },

  // 删除组织
  async deleteOrganization(_id: number): Promise<void> {
    await new Promise(resolve => setTimeout(resolve, 500));
  },

  // 移动组织
  async moveOrganization(_id: number, _targetParentId?: number): Promise<void> {
    await new Promise(resolve => setTimeout(resolve, 500));
  },
};

const OrganizationManagement: React.FC = () => {
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [organizationTree, setOrganizationTree] = useState<Organization[]>([]);
  const [loading, setLoading] = useState(false);
  const [treeLoading, setTreeLoading] = useState(false);
  // const [total, setTotal] = useState(0);
  // const [selectedOrg, setSelectedOrg] = useState<Organization | null>(null);
  const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([]);
  const [selectedKeys, setSelectedKeys] = useState<React.Key[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingOrg, setEditingOrg] = useState<Organization | null>(null);
  const [form] = Form.useForm();
  const [activeTab, setActiveTab] = useState('tree');

  useEffect(() => {
    if (activeTab === 'tree') {
      loadOrganizationTree();
    } else {
      loadOrganizations();
    }
  }, [activeTab]);

  const loadOrganizations = async () => {
    setLoading(true);
    try {
      const response = await organizationService.getOrganizations();
      setOrganizations(response.organizations);
      // setTotal(response.total);
    } catch (error) {
      message.error('加载组织列表失败');
    } finally {
      setLoading(false);
    }
  };

  const loadOrganizationTree = async () => {
    setTreeLoading(true);
    try {
      const tree = await organizationService.getOrganizationTree();
      setOrganizationTree(tree);
      
      // 默认展开第一层
      const firstLevelKeys = tree.map(org => org.id);
      setExpandedKeys(firstLevelKeys);
    } catch (error) {
      message.error('加载组织树失败');
    } finally {
      setTreeLoading(false);
    }
  };

  const handleCreate = (parentOrg?: Organization) => {
    setEditingOrg(null);
    setIsModalOpen(true);
    form.resetFields();
    form.setFieldsValue({
      parent_id: parentOrg?.id,
      type: OrganizationType.DEPARTMENT,
      is_enabled: true,
      sort_order: 0,
    });
  };

  const handleEdit = (org: Organization) => {
    setEditingOrg(org);
    setIsModalOpen(true);
    form.setFieldsValue({
      name: org.name,
      description: org.description,
      leader_id: org.leader_id,
      sort_order: org.sort_order,
      is_enabled: org.is_enabled,
    });
  };

  const handleDelete = async (org: Organization) => {
    if (org.employee_count > 0) {
      message.warning('该组织下还有员工，无法删除');
      return;
    }

    try {
      await organizationService.deleteOrganization(org.id);
      message.success('删除成功');
      if (activeTab === 'tree') {
        loadOrganizationTree();
      } else {
        loadOrganizations();
      }
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      
      if (editingOrg) {
        const updateData: UpdateOrganizationRequest = {
          name: values.name,
          description: values.description,
          leader_id: values.leader_id,
          sort_order: values.sort_order,
          is_enabled: values.is_enabled,
        };
        
        await organizationService.updateOrganization(editingOrg.id, updateData);
        message.success('更新成功');
      } else {
        const createData: CreateOrganizationRequest = {
          parent_id: values.parent_id,
          name: values.name,
          code: values.code,
          type: values.type,
          description: values.description,
          leader_id: values.leader_id,
          sort_order: values.sort_order,
          is_enabled: values.is_enabled,
        };
        
        await organizationService.createOrganization(createData);
        message.success('创建成功');
      }
      
      setIsModalOpen(false);
      form.resetFields();
      if (activeTab === 'tree') {
        loadOrganizationTree();
      } else {
        loadOrganizations();
      }
    } catch (error) {
      message.error(editingOrg ? '更新失败' : '创建失败');
    }
  };

  const getTypeTag = (type: OrganizationType) => {
    const typeConfig = {
      [OrganizationType.COMPANY]: { color: 'red', icon: <ApartmentOutlined />, text: '公司' },
      [OrganizationType.DEPARTMENT]: { color: 'blue', icon: <TeamOutlined />, text: '部门' },
      [OrganizationType.TEAM]: { color: 'green', icon: <UserOutlined />, text: '团队' },
      [OrganizationType.GROUP]: { color: 'orange', icon: <UserOutlined />, text: '小组' },
    };

    const config = typeConfig[type];
    return (
      <Tag color={config.color} icon={config.icon}>
        {config.text}
      </Tag>
    );
  };

  const renderTreeNode = (nodes: Organization[]): OrganizationTreeNode[] => {
    return nodes.map((node) => ({
      key: node.id,
      id: node.id,
      name: node.name,
      type: node.type,
      employee_count: node.employee_count,
      is_enabled: node.is_enabled,
      title: (
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, width: '100%' }}>
          <div style={{ flex: 1, display: 'flex', alignItems: 'center', gap: 8 }}>
            <span style={{ fontWeight: 500 }}>{node.name}</span>
            {getTypeTag(node.type)}
            <Badge count={node.employee_count} color="cyan" title="员工数量" />
            {!node.is_enabled && <Tag color="default">禁用</Tag>}
            {node.leader_name && (
              <Text type="secondary" style={{ fontSize: '12px' }}>
                负责人: {node.leader_name}
              </Text>
            )}
          </div>
          <Space size="small">
            <Button
              type="text"
              size="small"
              icon={<PlusOutlined />}
              onClick={(e) => {
                e.stopPropagation();
                handleCreate(node);
              }}
              title="添加下级组织"
            />
            <Button
              type="text"
              size="small"
              icon={<EditOutlined />}
              onClick={(e) => {
                e.stopPropagation();
                handleEdit(node);
              }}
              title="编辑"
            />
            <Popconfirm
              title="确定删除这个组织吗？"
              description="删除后不可恢复，请确保该组织下无员工"
              onConfirm={(e) => {
                e?.stopPropagation();
                handleDelete(node);
              }}
              onCancel={(e) => e?.stopPropagation()}
              okText="确定"
              cancelText="取消"
            >
              <Button
                type="text"
                size="small"
                danger
                icon={<DeleteOutlined />}
                onClick={(e) => e.stopPropagation()}
                title="删除"
              />
            </Popconfirm>
          </Space>
        </div>
      ),
      children: node.children ? renderTreeNode(node.children) : undefined,
    }));
  };

  const columns: TableColumnsType<Organization> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 60,
    },
    {
      title: '组织名称',
      dataIndex: 'name',
      width: 200,
      render: (text: string, record: Organization) => (
        <div>
          <div style={{ fontWeight: 500 }}>{text}</div>
          <Text type="secondary" style={{ fontSize: '12px' }}>
            {record.code}
          </Text>
        </div>
      ),
    },
    {
      title: '类型',
      dataIndex: 'type',
      width: 100,
      render: (type: OrganizationType) => getTypeTag(type),
    },
    {
      title: '负责人',
      dataIndex: 'leader_name',
      width: 120,
      render: (name: string) => name || '-',
    },
    {
      title: '员工数量',
      dataIndex: 'employee_count',
      width: 100,
      render: (count: number) => (
        <Badge count={count} showZero color="cyan" />
      ),
    },
    {
      title: '层级',
      dataIndex: 'level',
      width: 80,
      render: (level: number) => <Badge count={level} color="purple" />,
    },
    {
      title: '排序',
      dataIndex: 'sort_order',
      width: 80,
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
      render: (desc: string) => desc || '-',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      width: 160,
      render: (time: string) => time?.slice(0, 19).replace('T', ' '),
    },
    {
      title: '操作',
      width: 150,
      fixed: 'right',
      render: (_, record: Organization) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<PlusOutlined />}
            onClick={() => handleCreate(record)}
          >
            添加下级
          </Button>
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定删除这个组织吗？"
            description="删除后不可恢复"
            onConfirm={() => handleDelete(record)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              danger
              size="small"
              icon={<DeleteOutlined />}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const tabItems: TabsProps['items'] = [
    {
      key: 'tree',
      label: (
        <span>
          <ApartmentOutlined />
          组织树
        </span>
      ),
      children: (
        <Card>
          <div style={{ marginBottom: 16 }}>
            <Row justify="space-between" align="middle">
              <Col>
                <Alert
                  message="组织架构"
                  description="以树状结构展示组织架构，支持拖拽调整层级关系。"
                  type="info"
                  showIcon
                />
              </Col>
              <Col>
                <Space>
                  <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={() => handleCreate()}
                  >
                    新建组织
                  </Button>
                  <Button
                    icon={<ReloadOutlined />}
                    onClick={loadOrganizationTree}
                  >
                    刷新
                  </Button>
                </Space>
              </Col>
            </Row>
          </div>
          <Spin spinning={treeLoading}>
            {organizationTree.length > 0 ? (
              <Tree
                treeData={renderTreeNode(organizationTree)}
                expandedKeys={expandedKeys}
                selectedKeys={selectedKeys}
                onExpand={setExpandedKeys}
                onSelect={setSelectedKeys}
                showLine={{ showLeafIcon: false }}
                style={{ marginTop: 16 }}
                draggable
                onDrop={(info) => {
                  console.log('Drop:', info);
                  // 这里可以实现拖拽移动组织的逻辑
                }}
              />
            ) : (
              <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
                暂无组织数据
              </div>
            )}
          </Spin>
        </Card>
      ),
    },
    {
      key: 'list',
      label: '组织列表',
      children: (
        <div>
          <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
            <Col span={8}>
              <Search
                placeholder="搜索组织名称或代码"
                allowClear
                style={{ width: '100%' }}
              />
            </Col>
            <Col span={16}>
              <Space style={{ float: 'right' }}>
                <Button
                  type="primary"
                  icon={<PlusOutlined />}
                  onClick={() => handleCreate()}
                >
                  新建组织
                </Button>
                <Button
                  icon={<ReloadOutlined />}
                  onClick={loadOrganizations}
                >
                  刷新
                </Button>
              </Space>
            </Col>
          </Row>

          <Table
            columns={columns}
            dataSource={organizations}
            rowKey="id"
            loading={loading}
            scroll={{ x: 1200 }}
            pagination={{
              showSizeChanger: true,
              showQuickJumper: true,
              showTotal: (total, range) => `共 ${total} 条记录，显示第 ${range[0]}-${range[1]} 条`,
            }}
          />
        </div>
      ),
    },
  ];

  return (
    <div>
      <Card title="组织架构管理">
        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          items={tabItems}
        />
      </Card>

      <Modal
        title={editingOrg ? '编辑组织' : '新建组织'}
        open={isModalOpen}
        onOk={handleSubmit}
        onCancel={() => {
          setIsModalOpen(false);
          form.resetFields();
        }}
        width={600}
        okText="确定"
        cancelText="取消"
      >
        <Form
          form={form}
          layout="vertical"
          style={{ marginTop: 16 }}
        >
          {!editingOrg && (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label="组织名称"
                    name="name"
                    rules={[{ required: true, message: '请输入组织名称' }]}
                  >
                    <Input placeholder="如：技术部" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label="组织代码"
                    name="code"
                    rules={[
                      { required: true, message: '请输入组织代码' },
                      { pattern: /^[A-Z][A-Z0-9_]*$/, message: '只能包含大写字母、数字和下划线' }
                    ]}
                  >
                    <Input placeholder="如：TECH_DEPT" />
                  </Form.Item>
                </Col>
              </Row>

              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label="组织类型"
                    name="type"
                    rules={[{ required: true, message: '请选择组织类型' }]}
                  >
                    <Select placeholder="请选择组织类型">
                      <Option value={OrganizationType.COMPANY}>公司</Option>
                      <Option value={OrganizationType.DEPARTMENT}>部门</Option>
                      <Option value={OrganizationType.TEAM}>团队</Option>
                      <Option value={OrganizationType.GROUP}>小组</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    label="上级组织"
                    name="parent_id"
                  >
                    <Select placeholder="选择上级组织（可选）" allowClear>
                      {(organizations as Organization[])
                        .filter((org: Organization) => !editingOrg || org.id !== (editingOrg as Organization).id)
                        .map(org => (
                          <Option key={org.id} value={org.id}>
                            {org.name}
                          </Option>
                        ))}
                    </Select>
                  </Form.Item>
                </Col>
              </Row>
            </>
          )}

          {editingOrg && (
            <Form.Item
              label="组织名称"
              name="name"
              rules={[{ required: true, message: '请输入组织名称' }]}
            >
              <Input placeholder="如：技术部" />
            </Form.Item>
          )}

          <Form.Item
            label="组织描述"
            name="description"
          >
            <TextArea rows={3} placeholder="组织的详细描述" />
          </Form.Item>

          <Row gutter={16}>
            <Col span={8}>
              <Form.Item
                label="负责人"
                name="leader_id"
              >
                <Select placeholder="选择负责人" allowClear showSearch>
                  <Option value={1}>张三</Option>
                  <Option value={2}>李四</Option>
                  <Option value={3}>王五</Option>
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

export default OrganizationManagement;