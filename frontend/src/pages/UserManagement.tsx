/**
 * 用户管理页面
 */

import React, { useState, useEffect } from 'react';
import {
  Card,
  Table,
  Button,
  Input,
  Space,
  Modal,
  Form,
  Select,
  DatePicker,
  Switch,
  Tag,
  Dropdown,
  Popconfirm,
  App,
  Tooltip,
  Avatar,
  Typography,
  Row,
  Col
} from 'antd';
import { useNavigate } from 'react-router-dom';
import CreateUserModal from '../components/CreateUserModal';
import UserFilterPanel from '../components/UserFilterPanel';
import { useAuth } from '../hooks/useAuth';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  MoreOutlined,
  UserOutlined,
  ExportOutlined,
  ReloadOutlined,
  KeyOutlined,
  SafetyCertificateOutlined,
  MailOutlined,
  PhoneOutlined,
  FilterOutlined,
  SortAscendingOutlined,
  SortDescendingOutlined,
  DownOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { MenuProps } from 'antd';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';

dayjs.extend(relativeTime);
import { 
  getUsersList, 
  // createUser, 
  updateUser, 
  deleteUser, 
  resetUserPassword,
  // assignUserRoles,
  // getUserRoles,
  toggleUser2FA,
  batchDeleteUsers,
  exportUsers
} from '../services/userService';
import {
  getFilters,
  createFilter,
  deleteFilter,
  // buildFilterParams,
  parseFilterConditionText
} from '../services/filterService';
import type { 
  // CreateUserRequest, 
  UpdateUserRequest 
} from '../services/userService';
import type {
  SavedFilter,
  FilterConfig,
  SortConfig
} from '../services/filterService';
import type { User } from '../types/auth';
// import type { Role } from '../types/auth';

const { Search } = Input;
const { Text } = Typography;

export function UserManagement() {
  const { message } = App.useApp();
  const navigate = useNavigate();
  const { user: currentUser } = useAuth();
  
  // 状态管理
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchText, setSearchText] = useState('');
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  // 过滤器状态
  const [filterVisible, setFilterVisible] = useState(false);
  const [savedFilters, setSavedFilters] = useState<SavedFilter[]>([]);
  const [currentFilter, setCurrentFilter] = useState<FilterConfig | null>(null);
  const [currentSort, setCurrentSort] = useState<SortConfig>({
    field: 'username',
    direction: 'asc'
  });
  const [sortField, setSortField] = useState<string>('username');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');
  const [activeFilterId, setActiveFilterId] = useState<number | null>(null);

  // 模态框状态
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  // const [roleModalVisible, setRoleModalVisible] = useState(false);
  const [passwordModalVisible, setPasswordModalVisible] = useState(false);
  
  // 编辑状态
  const [editingUser, setEditingUser] = useState<User | null>(null);

  // 表单实例
  const [editForm] = Form.useForm();
  const [passwordForm] = Form.useForm();

  // 初始化数据
  useEffect(() => {
    loadUsers();
    loadSavedFilters();
  }, [page, pageSize, searchText, currentFilter, currentSort, activeFilterId]);

  // 加载用户列表
  const loadUsers = async () => {
    try {
      // 检查认证状态
      if (!checkAuthStatus()) {
        setLoading(false);
        return;
      }

      setLoading(true);
      
      // 构建查询参数（兼容原有API）
      const queryParams: any = {
        page,
        size: pageSize,
      };

      // 优先使用传统搜索（兼容性）
      if (searchText) {
        queryParams.search = searchText;
      }

      console.log('加载用户列表，参数:', queryParams);
      const response = await getUsersList(queryParams);
      
      console.log('获取到的用户列表:', response);
      setUsers(response.users);
      setTotal(response.total);
    } catch (error: any) {
      console.error('Failed to load users:', error);
      
      // 增强错误处理
      if (error.message?.includes('401') || error.message?.includes('Unauthorized')) {
        message.error('认证已过期，请重新登录');
        localStorage.removeItem('access_token');
        // 这里可以重定向到登录页面
      } else if (error.message?.includes('Network')) {
        message.error('网络连接失败，请检查网络后重试');
      } else {
        message.error(error.message || '加载用户列表失败');
      }
    } finally {
      setLoading(false);
    }
  };

  // 加载保存的过滤器
  const loadSavedFilters = async () => {
    try {
      const filters = await getFilters('users');
      setSavedFilters(filters);
    } catch (error) {
      console.error('Failed to load filters:', error);
      // 不显示错误消息，因为这是后台操作
    }
  };

  // 搜索处理
  const handleSearch = (value: string) => {
    setSearchText(value);
    setPage(1);
    // 清空过滤器状态，使用传统搜索
    setCurrentFilter(null);
    setActiveFilterId(null);
  };

  // 检查认证状态
  const checkAuthStatus = () => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      message.warning('请先登录以查看用户列表');
      // 这里可以重定向到登录页面，或显示登录模态框
      return false;
    }
    return true;
  };

  // 检查是否可以删除用户（不能删除自己）
  const canDeleteUser = (targetUser: User) => {
    return currentUser && currentUser.id !== targetUser.id;
  };

  // 切换过滤器面板
  const toggleFilterPanel = () => {
    setFilterVisible(!filterVisible);
  };

  // 应用过滤器
  const handleApplyFilter = (filterConfig: FilterConfig, sortConfig?: SortConfig) => {
    setCurrentFilter(filterConfig);
    if (sortConfig) {
      setCurrentSort(sortConfig);
    }
    setPage(1);
    setSearchText(''); // 清空传统搜索
    setActiveFilterId(null); // 清空保存的过滤器ID
  };

  // 保存过滤器
  const handleSaveFilter = async (
    name: string, 
    filterConfig: FilterConfig, 
    sortConfig?: SortConfig,
    isDefault?: boolean
  ) => {
    try {
      await createFilter({
        module_type: 'users',
        filter_name: name,
        filter_conditions: filterConfig,
        sort_config: sortConfig,
        is_default: isDefault
      });
      
      message.success('过滤器保存成功');
      loadSavedFilters(); // 重新加载保存的过滤器
    } catch (error: any) {
      message.error(error.message || '保存过滤器失败');
    }
  };

  // 删除过滤器
  const handleDeleteFilter = async (filterId: number) => {
    try {
      await deleteFilter(filterId);
      setSavedFilters(prev => prev.filter(f => f.id !== filterId));
      if (activeFilterId === filterId) {
        setActiveFilterId(null);
        setCurrentFilter(null);
      }
      message.success('过滤器删除成功');
    } catch (error: any) {
      message.error(error.message || '删除过滤器失败');
    }
  };

  // 加载保存的过滤器
  const handleLoadFilter = (filter: SavedFilter) => {
    setCurrentFilter(filter.filter_conditions);
    if (filter.sort_config) {
      setCurrentSort(filter.sort_config);
    }
    setActiveFilterId(filter.id);
    setPage(1);
    setSearchText(''); // 清空传统搜索
  };

  // 排序菜单选项
  const getSortMenuItems = () => [
    { key: 'username', label: '用户名' },
    { key: 'email', label: '邮箱' }, 
    { key: 'first_name', label: '姓名' },
    { key: 'is_active', label: '状态' },
    { key: 'created_at', label: '创建时间' },
    { key: 'last_login_at', label: '最后登录时间' },
  ];

  // 处理排序字段变化
  const handleSortChange = ({ key }: { key: string }) => {
    setSortField(key);
    setCurrentSort({ field: key, direction: sortOrder });
    setPage(1);
  };

  // 处理排序方向变化
  const handleSortOrderChange = () => {
    const newOrder = sortOrder === 'asc' ? 'desc' : 'asc';
    setSortOrder(newOrder);
    setCurrentSort({ field: sortField, direction: newOrder });
    setPage(1);
    message.success(`已切换到${newOrder === 'asc' ? '升序' : '降序'}排序`);
  };

  // 获取当前排序显示文本
  const getCurrentSortLabel = () => {
    const sortOptions = getSortMenuItems();
    const currentOption = sortOptions.find(item => item.key === sortField);
    return currentOption?.label || '用户名';
  };


  // 编辑用户
  const handleEditUser = async (values: UpdateUserRequest) => {
    if (!editingUser) return;
    
    try {
      console.log('更新用户数据:', values);
      console.log('用户ID:', editingUser.id);
      
      // 转换日期格式
      const updateData = {
        ...values,
        birth_date: values.birth_date ? dayjs(values.birth_date).format('YYYY-MM-DD') : undefined,
      };
      
      console.log('转换后的更新数据:', updateData);
      await updateUser(editingUser.id, updateData);
      message.success('更新用户成功');
      setEditModalVisible(false);
      setEditingUser(null);
      editForm.resetFields();
      loadUsers();
    } catch (error: any) {
      console.error('更新用户失败:', error);
      const errorMessage = error.message || '更新用户失败';
      message.error(errorMessage);
    }
  };

  // 删除用户
  const handleDeleteUser = async (user: User) => {
    try {
      await deleteUser(user.id);
      message.success('删除用户成功');
      loadUsers();
    } catch (error) {
      message.error('删除用户失败');
    }
  };

  // 批量删除用户
  const handleBatchDelete = async () => {
    try {
      // 过滤掉当前用户ID，防止自删除
      const validUserIds = (selectedRowKeys as number[]).filter(id => 
        currentUser && id !== currentUser.id
      );
      
      if (validUserIds.length === 0) {
        message.warning('没有可以删除的用户');
        return;
      }
      
      if (validUserIds.length < selectedRowKeys.length) {
        message.info('已自动排除当前用户，不能删除自己');
      }
      
      await batchDeleteUsers(validUserIds);
      message.success(`成功删除 ${validUserIds.length} 个用户`);
      setSelectedRowKeys([]);
      loadUsers();
    } catch (error) {
      message.error('批量删除失败');
    }
  };



  // 切换2FA
  const handleToggle2FA = async (user: User, enabled: boolean) => {
    try {
      await toggleUser2FA(user.id, enabled);
      message.success(`${enabled ? '启用' : '禁用'}双重认证成功`);
      loadUsers();
    } catch (error) {
      message.error('操作失败');
    }
  };

  // 导出用户
  const handleExportUsers = async () => {
    try {
      const blob = await exportUsers({
        search: searchText || undefined,
      });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.style.display = 'none';
      a.href = url;
      a.download = `users_${dayjs().format('YYYY-MM-DD_HH-mm-ss')}.xlsx`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      message.success('导出成功');
    } catch (error) {
      message.error('导出失败');
    }
  };

  // 打开编辑模态框
  const openEditModal = (user: User) => {
    console.log('打开编辑模态框，用户数据:', user);
    console.log('生日数据:', user.birth_date);
    
    setEditingUser(user);
    const formData = {
      first_name: user.first_name,
      last_name: user.last_name,
      phone: user.phone,
      gender: user.gender,
      birth_date: user.birth_date ? dayjs(user.birth_date) : null,
      is_active: (user as any).is_active,
    };
    
    console.log('表单数据:', formData);
    editForm.setFieldsValue(formData);
    setEditModalVisible(true);
  };

  // 打开角色模态框
  const openRoleModal = async (user: User) => {
    try {
      setEditingUser(user);
      // const roles = await getUserRoles(user.id);
      // setUserRoles(roles);
      // roleForm.setFieldsValue({
      //   role_ids: roles.map(role => role.id),
      // });
      // setRoleModalVisible(true);
    } catch (error) {
      message.error('获取用户角色失败');
    }
  };

  // 打开密码重置模态框
  const openPasswordModal = (user: User) => {
    setEditingUser(user);
    setPasswordModalVisible(true);
  };

  // 重置密码
  const handleResetPassword = async (values: { newPassword: string; confirmPassword: string }) => {
    if (!editingUser) return;
    
    try {
      await resetUserPassword(editingUser.id, values.newPassword);
      message.success('密码重置成功');
      setPasswordModalVisible(false);
      setEditingUser(null);
      passwordForm.resetFields();
    } catch (error: any) {
      const errorMessage = error.message || '密码重置失败';
      message.error(errorMessage);
    }
  };

  // 表格列定义
  const columns: ColumnsType<User> = [
    {
      title: '用户',
      dataIndex: 'username',
      key: 'username',
      render: (_, record) => (
        <Space>
          <Avatar 
            src={record.avatar_url} 
            icon={<UserOutlined />}
            size="small"
          />
          <div>
            <div>{record.username}</div>
            <Text type="secondary" style={{ fontSize: '12px' }}>
              {record.first_name} {record.last_name}
            </Text>
          </div>
        </Space>
      ),
    },
    {
      title: '联系信息',
      dataIndex: 'email',
      key: 'email',
      render: (_, record) => (
        <div>
          <div>
            <MailOutlined style={{ marginRight: 4 }} />
            {record.email}
          </div>
          {record.phone && (
            <div style={{ marginTop: 4 }}>
              <PhoneOutlined style={{ marginRight: 4 }} />
              <Text type="secondary">{record.phone}</Text>
            </div>
          )}
        </div>
      ),
    },
    {
      title: '状态',
      dataIndex: 'is_active',
      key: 'is_active',
      render: (isActive: boolean) => (
        <Tag color={isActive ? 'green' : 'red'}>
          {isActive ? '正常' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '角色',
      dataIndex: 'roles',
      key: 'roles',
      render: (roles: string[]) => (
        <div>
          {roles?.map(role => (
            <Tag key={role} color="blue" style={{ margin: '2px' }}>
              {role}
            </Tag>
          ))}
        </div>
      ),
    },
    {
      title: '2FA',
      dataIndex: 'two_factor_enabled',
      key: 'two_factor_enabled',
      render: (enabled: boolean, record) => (
        <Switch
          size="small"
          checked={enabled}
          onChange={(checked) => handleToggle2FA(record, checked)}
        />
      ),
    },
    {
      title: '最后登录',
      dataIndex: 'last_login_at',
      key: 'last_login_at',
      render: (time: string) => time && time !== '0001-01-01T00:00:00Z' ? (
        <Tooltip title={dayjs(time).format('YYYY-MM-DD HH:mm:ss')}>
          {dayjs(time).fromNow()}
        </Tooltip>
      ) : '从未登录',
    },
    {
      title: '操作',
      key: 'actions',
      width: 80,
      render: (_, record) => {
        const moreItems: MenuProps['items'] = [
          {
            key: 'edit',
            label: '编辑资料',
            icon: <EditOutlined />,
            onClick: () => openEditModal(record),
          },
          {
            key: 'userSettings',
            label: '用户设置',
            icon: <UserOutlined />,
            onClick: () => navigate(`/users/edit/${record.id}?tab=basic`),
          },
          {
            key: 'editPermissions',
            label: '编辑权限',
            icon: <SafetyCertificateOutlined />,
            onClick: () => navigate(`/users/edit/${record.id}?tab=roles`),
          },
          {
            type: 'divider',
          },
          {
            key: 'roles',
            label: '分配角色',
            icon: <SafetyCertificateOutlined />,
            onClick: () => openRoleModal(record),
          },
          {
            key: 'password',
            label: '重置密码',
            icon: <KeyOutlined />,
            onClick: () => openPasswordModal(record),
          },
          {
            type: 'divider',
          },
          {
            key: 'delete',
            label: canDeleteUser(record) ? '删除用户' : (
              <Tooltip title="不能删除自己">
                <span style={{ color: '#ccc' }}>删除用户</span>
              </Tooltip>
            ),
            icon: <DeleteOutlined />,
            danger: true,
            disabled: !canDeleteUser(record),
            onClick: canDeleteUser(record) ? () => {
              Modal.confirm({
                title: '确认删除',
                content: `确定要删除用户 "${record.username}" 吗？`,
                onOk: () => handleDeleteUser(record),
              });
            } : undefined,
          },
        ];

        return (
          <Dropdown menu={{ items: moreItems }} trigger={['click']}>
            <Button type="text" size="small" icon={<MoreOutlined />} title="更多操作" />
          </Dropdown>
        );
      },
    },
  ];

  // 行选择配置
  const rowSelection = {
    selectedRowKeys,
    onChange: setSelectedRowKeys,
    getCheckboxProps: (record: User) => ({
      disabled: record.username === 'admin' || !canDeleteUser(record), // 保护管理员账户和当前用户
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
                placeholder="搜索用户名、邮箱、姓名"
                allowClear
                onSearch={handleSearch}
                style={{ width: 300 }}
                value={searchText}
              />
              <Button 
                icon={<ReloadOutlined />} 
                onClick={loadUsers}
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
                创建用户
              </Button>
              <Button 
                icon={<FilterOutlined />}
                onClick={toggleFilterPanel}
                type={filterVisible ? 'primary' : 'default'}
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
                onClick={handleExportUsers}
              >
                导出
              </Button>
              {selectedRowKeys.length > 0 && (() => {
                // 计算实际可删除的用户数量
                const validUserIds = (selectedRowKeys as number[]).filter(id => 
                  currentUser && id !== currentUser.id
                );
                return validUserIds.length > 0 ? (
                  <Popconfirm
                    title={`确定要删除选中的 ${validUserIds.length} 个用户吗？${validUserIds.length < selectedRowKeys.length ? '（已排除当前用户）' : ''}`}
                    onConfirm={handleBatchDelete}
                    okText="确定"
                    cancelText="取消"
                  >
                    <Button danger icon={<DeleteOutlined />}>
                      批量删除 ({validUserIds.length})
                    </Button>
                  </Popconfirm>
                ) : (
                  <Tooltip title="选中的用户中没有可删除的用户">
                    <Button danger disabled icon={<DeleteOutlined />}>
                      批量删除 ({selectedRowKeys.length})
                    </Button>
                  </Tooltip>
                );
              })()}
            </Space>
          </Col>
        </Row>
      </Card>

      {/* 过滤器面板 */}
      <UserFilterPanel
        visible={filterVisible}
        onToggle={toggleFilterPanel}
        onApplyFilter={handleApplyFilter}
        savedFilters={savedFilters}
        onSaveFilter={handleSaveFilter}
        onDeleteFilter={handleDeleteFilter}
        onLoadFilter={handleLoadFilter}
        loading={loading}
      />

      {/* 当前过滤器状态显示 */}
      {(currentFilter || searchText) && (
        <Card style={{ marginBottom: '16px' }}>
          <Row justify="space-between" align="middle">
            <Col>
              <Space>
                <FilterOutlined style={{ color: '#1890ff' }} />
                <span style={{ fontWeight: 500 }}>当前过滤条件：</span>
                {searchText && (
                  <Tag color="blue">搜索: {searchText}</Tag>
                )}
                {currentFilter?.conditions.map((condition, index) => (
                  <Tag key={index} color="blue">
                    {parseFilterConditionText(condition)}
                  </Tag>
                ))}
                {currentSort && (
                  <Tag color="green">
                    排序: {currentSort.field} ({currentSort.direction === 'asc' ? '升序' : '降序'})
                  </Tag>
                )}
              </Space>
            </Col>
            <Col>
              <Button 
                size="small"
                onClick={() => {
                  setCurrentFilter(null);
                  setSearchText('');
                  setActiveFilterId(null);
                  setCurrentSort({ field: 'username', direction: 'asc' });
                }}
              >
                清除所有筛选
              </Button>
            </Col>
          </Row>
        </Card>
      )}

      {/* 用户表格 */}
      <Card>
        <Table
          columns={columns}
          dataSource={users}
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

      {/* 创建用户模态框 */}
      <CreateUserModal 
        visible={createModalVisible}
        onClose={() => setCreateModalVisible(false)}
        onSuccess={() => loadUsers()}
      />

      {/* 编辑用户模态框 */}
      <Modal
        title="编辑用户"
        open={editModalVisible}
        onCancel={() => {
          setEditModalVisible(false);
          setEditingUser(null);
          editForm.resetFields();
        }}
        footer={null}
        width={600}
      >
        <Form
          form={editForm}
          layout="vertical"
          onFinish={handleEditUser}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="first_name"
                label="名字"
                rules={[{ required: true, message: '请输入名字' }]}
              >
                <Input placeholder="请输入名字" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="last_name"
                label="姓氏"
                rules={[{ required: true, message: '请输入姓氏' }]}
              >
                <Input placeholder="请输入姓氏" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="gender"
                label="性别"
              >
                <Select placeholder="请选择性别" allowClear>
                  <Select.Option value="M">男</Select.Option>
                  <Select.Option value="F">女</Select.Option>
                  <Select.Option value="O">不愿透露</Select.Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="birth_date"
                label="出生日期"
              >
                <DatePicker 
                  style={{ width: '100%' }}
                  placeholder="请选择出生日期" 
                />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="phone"
                label="手机号"
                rules={[
                  { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确' },
                ]}
              >
                <Input placeholder="请输入手机号" />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="is_active"
            valuePropName="checked"
          >
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <span>账户状态：</span>
              <Switch checkedChildren="启用" unCheckedChildren="禁用" />
            </div>
          </Form.Item>

          <Row justify="end">
            <Button type="primary" htmlType="submit">
              保存
            </Button>
          </Row>
        </Form>
      </Modal>

      {/* 重置密码模态框 */}
      <Modal
        title="重置密码"
        open={passwordModalVisible}
        onCancel={() => {
          setPasswordModalVisible(false);
          setEditingUser(null);
          passwordForm.resetFields();
        }}
        footer={null}
        width={400}
      >
        <Form
          form={passwordForm}
          layout="vertical"
          onFinish={handleResetPassword}
        >
          <Form.Item
            name="newPassword"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 8, message: '密码至少8个字符' },
              { 
                pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]/,
                message: '密码需包含大小写字母、数字和特殊字符'
              },
            ]}
          >
            <Input.Password placeholder="请输入新密码" />
          </Form.Item>

          <Form.Item
            name="confirmPassword"
            label="确认密码"
            dependencies={['newPassword']}
            rules={[
              { required: true, message: '请确认新密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('newPassword') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'));
                },
              }),
            ]}
          >
            <Input.Password placeholder="请再次输入新密码" />
          </Form.Item>

          <Row justify="end">
            <Space>
              <Button onClick={() => {
                setPasswordModalVisible(false);
                setEditingUser(null);
                passwordForm.resetFields();
              }}>
                取消
              </Button>
              <Button type="primary" htmlType="submit">
                重置密码
              </Button>
            </Space>
          </Row>
        </Form>
      </Modal>

      {/* 其他模态框可以继续添加... */}
    </div>
  );
}

export default UserManagement;