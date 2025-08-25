import React, { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Space,
  Input,
  Select,
  Card,
  Modal,
  Form,
  Popconfirm,
  message,
  Tag,
  Tooltip,
  Row,
  Col,
  Switch,
  InputNumber,
  Typography
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  CopyOutlined,
  ReloadOutlined,
  DatabaseOutlined,
  QuestionCircleOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { DocType } from '../types/auth';
import {
  getDocTypeList,
  createDocType,
  updateDocType,
  deleteDocType,
  batchDeleteDocTypes,
  copyDocType,
  validateDocTypeName,
  getDocTypeModules,
  DOCTYPE_MODULES,
  type DocTypeQueryParams,
  type CreateDocTypeRequest,
  type UpdateDocTypeRequest
} from '../services/doctypeService';
import DocTypeTemplates, { type DocTypeTemplate } from '../components/DocTypeTemplates';
import './DocTypeManagement.css';
import '../components/DocTypeTemplates.css';

const { Option } = Select;
const { Title } = Typography;
const { TextArea } = Input;

// 表单模式
type FormMode = 'create' | 'edit' | 'copy';

const DocTypeManagement: React.FC = () => {
  // 数据状态
  const [doctypes, setDocTypes] = useState<DocType[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  
  // 搜索状态
  const [searchKeyword, setSearchKeyword] = useState('');
  const [filterModule, setFilterModule] = useState<string | undefined>(undefined);
  const [filterSubmittable, setFilterSubmittable] = useState<boolean | undefined>(undefined);
  const [filterChildTable, setFilterChildTable] = useState<boolean | undefined>(undefined);
  
  // 模态框状态
  const [modalVisible, setModalVisible] = useState(false);
  const [templateModalVisible, setTemplateModalVisible] = useState(false);
  const [modalMode, setModalMode] = useState<FormMode>('create');
  const [editingRecord, setEditingRecord] = useState<DocType | null>(null);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  
  // 其他状态
  const [modules, setModules] = useState<string[]>([]);
  const [form] = Form.useForm();

  // 初始化
  useEffect(() => {
    loadDocTypes();
    loadModules();
  }, [currentPage, pageSize, searchKeyword, filterModule, filterSubmittable, filterChildTable]);

  // 加载DocType列表
  const loadDocTypes = async () => {
    try {
      setLoading(true);
      
      const params: DocTypeQueryParams = {
        page: currentPage,
        page_size: pageSize,
        keyword: searchKeyword || undefined,
        module: filterModule,
        is_submittable: filterSubmittable,
        is_child_table: filterChildTable,
        sort_by: 'created_at',
        sort_order: 'desc'
      };

      const response = await getDocTypeList(params);
      setDocTypes(response.items || []);
      setTotal(response.total || 0);
    } catch (error: any) {
      message.error(error.message || '加载DocType列表失败');
    } finally {
      setLoading(false);
    }
  };

  // 加载模块列表
  const loadModules = async () => {
    try {
      const moduleList = await getDocTypeModules();
      setModules(moduleList);
    } catch (error) {
      // 使用默认模块
      setModules([...DOCTYPE_MODULES]);
    }
  };

  // 处理搜索
  const handleSearch = (value: string) => {
    setSearchKeyword(value);
    setCurrentPage(1);
  };

  // 处理筛选变化
  const handleFilterChange = () => {
    setCurrentPage(1);
  };

  // 重置搜索和筛选
  const handleReset = () => {
    setSearchKeyword('');
    setFilterModule(undefined);
    setFilterSubmittable(undefined);
    setFilterChildTable(undefined);
    setCurrentPage(1);
  };

  // 刷新数据
  const handleRefresh = () => {
    loadDocTypes();
  };

  // 打开模板选择模态框
  const handleCreate = () => {
    setTemplateModalVisible(true);
  };

  // 选择模板后打开创建模态框
  const handleSelectTemplate = (template: DocTypeTemplate) => {
    setModalMode('create');
    setEditingRecord(null);
    form.resetFields();
    
    // 根据模板预填充字段
    form.setFieldsValue({
      name: template.name,
      label: template.label,
      module: template.module,
      description: template.description,
      is_submittable: template.config.is_submittable,
      is_child_table: template.config.is_child_table,
      has_workflow: template.config.has_workflow,
      track_changes: template.config.track_changes,
      applies_to_all_users: template.config.applies_to_all_users,
      max_attachments: 0,
      sort_order: 'DESC'
    });
    
    setTemplateModalVisible(false);
    setModalVisible(true);
  };

  // 创建空白文档类型
  const handleCreateBlank = () => {
    setModalMode('create');
    setEditingRecord(null);
    form.resetFields();
    setTemplateModalVisible(false);
    setModalVisible(true);
  };

  // 根据布尔值确定文档类型
  const getDocTypeCategory = (is_submittable: boolean, is_child_table: boolean) => {
    if (is_child_table) return 'child';
    if (is_submittable) return 'submittable';
    return 'normal';
  };

  // 打开编辑模态框
  const handleEdit = (record: DocType) => {
    setModalMode('edit');
    setEditingRecord(record);
    
    const docTypeCategory = getDocTypeCategory(record.is_submittable || false, record.is_child_table || false);
    
    form.setFieldsValue({
      name: record.name,
      label: record.label,
      module: record.module,
      description: record.description,
      doctype_category: docTypeCategory, // 新增：文档类型分类
      is_submittable: record.is_submittable,
      is_child_table: record.is_child_table,
      has_workflow: record.has_workflow,
      track_changes: record.track_changes,
      applies_to_all_users: record.applies_to_all_users,
      max_attachments: record.max_attachments,
      naming_rule: record.naming_rule,
      title_field: record.title_field,
      search_fields: record.search_fields ? JSON.parse(record.search_fields) : [],
      sort_field: record.sort_field,
      sort_order: record.sort_order
    });
    setModalVisible(true);
  };

  // 打开复制模态框
  const handleCopy = (record: DocType) => {
    setModalMode('copy');
    setEditingRecord(record);
    
    const docTypeCategory = getDocTypeCategory(record.is_submittable || false, record.is_child_table || false);
    
    form.setFieldsValue({
      name: `${record.name}_copy`,
      label: `${record.label} (副本)`,
      module: record.module,
      description: record.description,
      doctype_category: docTypeCategory, // 新增：文档类型分类
      is_submittable: record.is_submittable,
      is_child_table: record.is_child_table,
      has_workflow: record.has_workflow,
      track_changes: record.track_changes,
      applies_to_all_users: record.applies_to_all_users,
      max_attachments: record.max_attachments,
      naming_rule: record.naming_rule,
      title_field: record.title_field,
      search_fields: record.search_fields ? JSON.parse(record.search_fields) : [],
      sort_field: record.sort_field,
      sort_order: record.sort_order
    });
    setModalVisible(true);
  };

  // 处理删除
  const handleDelete = async (id: number) => {
    try {
      await deleteDocType(id);
      message.success('删除成功');
      loadDocTypes();
    } catch (error: any) {
      message.error(error.message || '删除失败');
    }
  };

  // 处理批量删除
  const handleBatchDelete = async () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要删除的记录');
      return;
    }

    try {
      await batchDeleteDocTypes(selectedRowKeys as number[]);
      message.success('批量删除成功');
      setSelectedRowKeys([]);
      loadDocTypes();
    } catch (error: any) {
      message.error(error.message || '批量删除失败');
    }
  };

  // 提交表单
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      
      // 处理搜索字段数组
      const searchFields = values.search_fields || [];
      const formData = {
        ...values,
        search_fields: searchFields.length > 0 ? searchFields : undefined
      };

      if (modalMode === 'create') {
        await createDocType(formData as CreateDocTypeRequest);
        message.success('创建成功');
      } else if (modalMode === 'edit' && editingRecord) {
        const { name, ...updateData } = formData; // name字段不允许修改
        await updateDocType(editingRecord.id, updateData as UpdateDocTypeRequest);
        message.success('更新成功');
      } else if (modalMode === 'copy' && editingRecord) {
        await copyDocType(editingRecord.id, formData.name, formData.label);
        message.success('复制成功');
      }

      setModalVisible(false);
      loadDocTypes();
    } catch (error: any) {
      if (error.errorFields) {
        message.error('请检查表单输入');
      } else {
        message.error(error.message || '操作失败');
      }
    }
  };

  // 验证DocType名称
  const validateName = async (_: any, value: string) => {
    if (!value) {
      return Promise.reject(new Error('请输入DocType名称'));
    }

    // 名称格式验证
    if (!/^[A-Za-z][A-Za-z0-9_\s]*$/.test(value)) {
      return Promise.reject(new Error('名称必须以字母开头，只能包含字母、数字、下划线和空格'));
    }

    // 编辑模式时，如果名称未变化则跳过验证
    if (modalMode === 'edit' && editingRecord && value === editingRecord.name) {
      return Promise.resolve();
    }

    try {
      const result = await validateDocTypeName(value);
      if (!result.available) {
        return Promise.reject(new Error(result.message || '名称已存在'));
      }
      return Promise.resolve();
    } catch (error) {
      return Promise.reject(new Error('名称验证失败'));
    }
  };

  // 表格列定义
  const columns: ColumnsType<DocType> = [
    {
      title: '文档类型',
      dataIndex: 'name',
      key: 'name',
      width: 250,
      render: (text: string, record: DocType) => (
        <Space>
          <DatabaseOutlined style={{ color: '#1890ff' }} />
          <div>
            <div style={{ fontWeight: 500, fontSize: '14px' }}>{record.label}</div>
            <div style={{ fontSize: '12px', color: '#999' }}>系统名称: {text}</div>
          </div>
        </Space>
      ),
    },
    {
      title: '所属模块',
      dataIndex: 'module',
      key: 'module',
      width: 120,
      render: (text: string) => <Tag color="blue">{text}</Tag>,
    },
    {
      title: '类型',
      key: 'type',
      width: 120,
      render: (_, record: DocType) => {
        if (record.is_child_table) return <Tag color="orange">子表</Tag>;
        if (record.is_submittable) return <Tag color="green">审批文档</Tag>;
        return <Tag color="blue">普通文档</Tag>;
      },
    },
    {
      title: '状态',
      key: 'status',
      width: 100,
      render: (_, record: DocType) => {
        const hasWorkflow = record.has_workflow;
        const trackChanges = record.track_changes;
        if (hasWorkflow) return <Tag color="purple">工作流</Tag>;
        if (trackChanges) return <Tag color="cyan">跟踪变更</Tag>;
        return <Tag color="default">基础</Tag>;
      },
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: {
        showTitle: false,
      },
      render: (text: string) => (
        <Tooltip placement="topLeft" title={text}>
          {text || '暂无描述'}
        </Tooltip>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 150,
      render: (text: string) => new Date(text).toLocaleDateString(),
    },
    {
      title: '操作',
      key: 'action',
      width: 160,
      fixed: 'right',
      render: (_, record: DocType) => (
        <Space>
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Button
            type="link"
            size="small"
            icon={<CopyOutlined />}
            onClick={() => handleCopy(record)}
          >
            复制
          </Button>
          <Popconfirm
            title="确定要删除这个文档类型吗？"
            description="删除后将无法恢复"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              size="small"
              danger
              icon={<DeleteOutlined />}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  // 行选择配置
  const rowSelection = {
    selectedRowKeys,
    onChange: (selectedRowKeys: React.Key[]) => {
      setSelectedRowKeys(selectedRowKeys);
    },
  };

  return (
    <div className="doctype-management">
      <Card>
        <div className="doctype-management-header">
          <Title level={4}>
            <DatabaseOutlined /> 文档类型管理
          </Title>
          <div style={{ marginBottom: 24 }}>
            <p style={{ color: '#666', marginBottom: 16 }}>
              管理系统中的各种业务文档类型，如用户信息、订单记录、产品目录等。文档类型定义了数据的结构和权限规则。
            </p>
            
            {doctypes.length === 0 && !loading && (
              <div style={{ 
                background: '#f0f2ff', 
                border: '1px solid #d6e4ff', 
                borderRadius: '6px',
                padding: '16px',
                marginBottom: '16px'
              }}>
                <Title level={5} style={{ color: '#1890ff', marginBottom: 8 }}>
                  🎉 欢迎使用文档类型管理
                </Title>
                <p style={{ color: '#595959', marginBottom: 12 }}>
                  您还没有创建任何文档类型。文档类型是系统的基础，用来定义不同业务数据的结构。
                </p>
                <p style={{ color: '#595959', marginBottom: 8 }}>
                  <strong>快速开始：</strong>
                </p>
                <ul style={{ color: '#595959', paddingLeft: '20px', marginBottom: 12 }}>
                  <li>点击 "新建文档类型" 按钮</li>
                  <li>从预设模板中选择合适的类型</li>
                  <li>填写基本信息即可完成创建</li>
                </ul>
                <Button 
                  type="primary" 
                  size="small"
                  onClick={handleCreate}
                  style={{ marginTop: '8px' }}
                >
                  立即创建第一个文档类型
                </Button>
              </div>
            )}
          </div>
          
          {/* 简化的搜索区域 */}
          <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
            <Col span={8}>
              <Input.Search
                placeholder="搜索文档类型名称"
                allowClear
                onSearch={handleSearch}
                onChange={(e) => {
                  if (!e.target.value) {
                    handleSearch('');
                  }
                }}
              />
            </Col>
            <Col span={6}>
              <Select
                placeholder="选择模块"
                allowClear
                value={filterModule}
                onChange={(value) => {
                  setFilterModule(value);
                  handleFilterChange();
                }}
                style={{ width: '100%' }}
              >
                {modules.map(module => (
                  <Option key={module} value={module}>{module}</Option>
                ))}
              </Select>
            </Col>
            <Col span={10}>
              <Space>
                <Button onClick={handleReset}>重置筛选</Button>
                <Button icon={<ReloadOutlined />} onClick={handleRefresh}>
                  刷新
                </Button>
              </Space>
            </Col>
          </Row>

          {/* 操作按钮区域 */}
          <div style={{ marginBottom: 16 }}>
            <Space>
              <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={handleCreate}
                size="large"
              >
                新建文档类型
              </Button>
              {selectedRowKeys.length > 0 && (
                <Popconfirm
                  title="确定要删除选中的文档类型吗？"
                  description="删除后将无法恢复，相关权限配置也会被清除"
                  onConfirm={handleBatchDelete}
                  okText="确定"
                  cancelText="取消"
                >
                  <Button
                    danger
                    icon={<DeleteOutlined />}
                  >
                    删除选中项 ({selectedRowKeys.length})
                  </Button>
                </Popconfirm>
              )}
            </Space>
          </div>
        </div>

        {/* 数据表格 */}
        <Table
          columns={columns}
          dataSource={doctypes}
          rowKey="id"
          rowSelection={rowSelection}
          loading={loading}
          scroll={{ x: 1200 }}
          pagination={{
            current: currentPage,
            pageSize: pageSize,
            total: total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total, range) =>
              `第 ${range[0]}-${range[1]} 条/共 ${total} 条`,
            onChange: (page, size) => {
              setCurrentPage(page);
              setPageSize(size || 20);
            },
          }}
        />
      </Card>

      {/* 简化的创建/编辑模态框 */}
      <Modal
        title={
          modalMode === 'create' 
            ? '新建文档类型' 
            : modalMode === 'edit' 
              ? '编辑文档类型' 
              : '复制文档类型'
        }
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={600}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={{
            is_submittable: false,
            is_child_table: false,
            has_workflow: false,
            track_changes: true,
            applies_to_all_users: false,
            max_attachments: 0,
            sort_order: 'DESC'
          }}
        >
          {/* 基础信息 */}
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label={
                  <Space>
                    显示名称
                    <Tooltip title="用户看到的友好名称，如：用户信息、订单记录">
                      <QuestionCircleOutlined style={{ color: '#999', fontSize: '14px' }} />
                    </Tooltip>
                  </Space>
                }
                name="label"
                rules={[{ required: true, message: '请输入显示名称' }]}
              >
                <Input placeholder="如：用户信息" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label={
                  <Space>
                    系统名称
                    <Tooltip title="系统内部使用的英文名称，如：User、Order">
                      <QuestionCircleOutlined style={{ color: '#999', fontSize: '14px' }} />
                    </Tooltip>
                  </Space>
                }
                name="name"
                rules={[{ validator: validateName }]}
              >
                <Input 
                  placeholder="如：User"
                  disabled={modalMode === 'edit'}
                />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label="所属模块"
                name="module"
                rules={[{ required: true, message: '请选择所属模块' }]}
              >
                <Select placeholder="选择所属模块">
                  {modules.map(module => (
                    <Option key={module} value={module}>{module}</Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label={
                  <Space>
                    文档类型
                    <Tooltip title="选择文档的业务类型">
                      <QuestionCircleOutlined style={{ color: '#999', fontSize: '14px' }} />
                    </Tooltip>
                  </Space>
                }
                name="doctype_category"
              >
                <Select placeholder="选择文档类型" onChange={(value) => {
                  if (value === 'submittable') {
                    form.setFieldsValue({ is_submittable: true, is_child_table: false });
                  } else if (value === 'child') {
                    form.setFieldsValue({ is_submittable: false, is_child_table: true });
                  } else {
                    form.setFieldsValue({ is_submittable: false, is_child_table: false });
                  }
                }}>
                  <Option value="normal">普通文档</Option>
                  <Option value="submittable">审批文档</Option>
                  <Option value="child">子表文档</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            label="描述"
            name="description"
          >
            <TextArea rows={2} placeholder="简单描述这个文档类型的用途" />
          </Form.Item>

          {/* 隐藏的技术字段 */}
          <Form.Item name="is_submittable" style={{ display: 'none' }}>
            <Switch />
          </Form.Item>
          <Form.Item name="is_child_table" style={{ display: 'none' }}>
            <Switch />
          </Form.Item>
          <Form.Item name="has_workflow" style={{ display: 'none' }}>
            <Switch />
          </Form.Item>
          <Form.Item name="track_changes" style={{ display: 'none' }}>
            <Switch />
          </Form.Item>
          <Form.Item name="applies_to_all_users" style={{ display: 'none' }}>
            <Switch />
          </Form.Item>
          <Form.Item name="max_attachments" style={{ display: 'none' }}>
            <InputNumber />
          </Form.Item>
          <Form.Item name="naming_rule" style={{ display: 'none' }}>
            <Input />
          </Form.Item>
          <Form.Item name="title_field" style={{ display: 'none' }}>
            <Input />
          </Form.Item>
          <Form.Item name="sort_field" style={{ display: 'none' }}>
            <Input />
          </Form.Item>
          <Form.Item name="sort_order" style={{ display: 'none' }}>
            <Select />
          </Form.Item>
          <Form.Item name="search_fields" style={{ display: 'none' }}>
            <Select />
          </Form.Item>
        </Form>
      </Modal>

      {/* 模板选择模态框 */}
      <Modal
        title="选择文档类型模板"
        open={templateModalVisible}
        onCancel={() => setTemplateModalVisible(false)}
        footer={null}
        width={800}
        destroyOnClose
      >
        <DocTypeTemplates 
          onSelectTemplate={handleSelectTemplate}
          onCreateBlank={handleCreateBlank}
        />
      </Modal>
    </div>
  );
};

export default DocTypeManagement;