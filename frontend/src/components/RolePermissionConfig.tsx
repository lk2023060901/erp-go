/**
 * 角色权限配置组件 - 支持ERP权限系统
 */

import React, { useState, useEffect } from 'react';
import {
  Modal,
  Table,
  Button,
  Space,
  Select,
  Form,
  Input,
  message,
  Tabs,
  Alert,
  Spin,
  Tag,
  Collapse,
  Checkbox,
  Row,
  Col,
  Divider,
  Typography
} from 'antd';
import {
  DatabaseOutlined,
  SafetyCertificateOutlined,
  SettingOutlined,
  FileTextOutlined,
  DeleteOutlined,
  PlusOutlined,
  EditOutlined,
  BulbOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { Role, DocType, PermissionRule } from '../types/auth';
import {
  getDocTypesList,
  getPermissionRulesList,
  createPermissionRule,
  updatePermissionRule,
  deletePermissionRule
} from '../services/erpPermissionService';
import './RolePermissionConfig.css';

const { Option } = Select;
const { Text } = Typography;

// ERP权限类型常量
const PERMISSION_TYPES = [
  'can_read',
  'can_write', 
  'can_create',
  'can_delete',
  'can_submit',
  'can_cancel',
  'can_amend',
  'can_print',
  'can_email',
  'can_import',
  'can_export',
  'can_share',
  'can_report'
];

// 权限标签映射
const PERMISSION_LABELS: Record<string, string> = {
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

interface RolePermissionConfigProps {
  role: Role;
  visible: boolean;
  onClose: () => void;
  onSave?: () => void;
}

// 权限类型图标映射
const getPermissionIcon = (type: string) => {
  const iconMap: Record<string, React.ReactNode> = {
    'read': <FileTextOutlined style={{ color: '#52c41a' }} />,
    'write': <EditOutlined style={{ color: '#1890ff' }} />,
    'create': <PlusOutlined style={{ color: '#722ed1' }} />,
    'delete': <DeleteOutlined style={{ color: '#f5222d' }} />,
    'submit': <SafetyCertificateOutlined style={{ color: '#fa8c16' }} />,
    'cancel': <SafetyCertificateOutlined style={{ color: '#faad14' }} />,
    'amend': <EditOutlined style={{ color: '#13c2c2' }} />,
    'print': <FileTextOutlined style={{ color: '#eb2f96' }} />,
    'email': <FileTextOutlined style={{ color: '#52c41a' }} />,
    'import': <PlusOutlined style={{ color: '#1890ff' }} />,
    'export': <FileTextOutlined style={{ color: '#722ed1' }} />,
    'share': <SafetyCertificateOutlined style={{ color: '#fa8c16' }} />,
    'report': <FileTextOutlined style={{ color: '#faad14' }} />
  };
  return iconMap[type] || <SettingOutlined />;
};

const RolePermissionConfig: React.FC<RolePermissionConfigProps> = ({
  role,
  visible,
  onClose,
  onSave
}) => {
  // 状态管理
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState('doctype-permissions');
  
  // 数据状态
  const [docTypes, setDocTypes] = useState<DocType[]>([]);
  const [permissionRules, setPermissionRules] = useState<PermissionRule[]>([]);
  
  // 编辑状态
  const [editingRule, setEditingRule] = useState<PermissionRule | null>(null);
  const [editFormVisible, setEditFormVisible] = useState(false);
  const [form] = Form.useForm();

  // 加载数据
  useEffect(() => {
    if (visible) {
      loadData();
    }
  }, [visible, role.id]);

  const loadData = async () => {
    try {
      setLoading(true);
      
      console.log('🔍 Loading permission data for role:', role.id, role.name);
      
      const [docTypesResponse, permissionsResponse] = await Promise.all([
        getDocTypesList(), // 不传分页参数，因为该函数不接受这些参数
        getPermissionRulesList(role.id) // 只传role_id，因为函数签名是 getPermissionRulesList(roleId?, docType?)
      ]);
      
      console.log('📋 DocTypes response:', docTypesResponse);
      console.log('🔐 Permissions response:', permissionsResponse);
      
      setDocTypes(docTypesResponse.doctypes || []);
      setPermissionRules(permissionsResponse.rules || []);
    } catch (error: any) {
      console.error('❌ Failed to load data:', error);
      console.error('Error details:', {
        message: error.message,
        response: error.response,
        status: error.response?.status,
        data: error.response?.data
      });
      
      // 提供更详细的错误信息
      const errorMessage = error.response?.data?.error?.message || error.message || '未知错误';
      message.error(`加载权限数据失败: ${errorMessage}`);
    } finally {
      setLoading(false);
    }
  };

  // 按模块分组DocType
  const docTypesByModule = docTypes.reduce((acc, docType) => {
    if (!acc[docType.module]) {
      acc[docType.module] = [];
    }
    acc[docType.module].push(docType);
    return acc;
  }, {} as Record<string, DocType[]>);

  // 获取特定DocType的权限规则
  const getDocTypePermissions = (docTypeName: string) => {
    return permissionRules.filter(rule => rule.document_type === docTypeName);
  };

  // 创建新的权限规则
  const handleCreatePermissionRule = async (docType: DocType) => {
    setEditingRule(null);
    form.resetFields();
    form.setFieldsValue({
      document_type: docType.name,
      permission_level: 0,
      read: true,
      write: false,
      create: false,
      delete: false,
      submit: false,
      cancel: false,
      amend: false,
      print: false,
      email: false,
      import: false,
      export: false,
      share: false,
      report: false,
      set_user_permissions: false,
      if_owner: false
    });
    setEditFormVisible(true);
  };

  // 编辑权限规则
  const handleEditPermissionRule = (rule: PermissionRule) => {
    setEditingRule(rule);
    form.setFieldsValue({
      document_type: rule.document_type,
      permission_level: rule.permission_level,
      read: rule.read,
      write: rule.write,
      create: rule.create,
      delete: rule.delete,
      submit: rule.submit,
      cancel: rule.cancel,
      amend: rule.amend,
      print: rule.print,
      email: rule.email,
      import: rule.import,
      export: rule.export,
      share: rule.share,
      report: rule.report,
      set_user_permissions: rule.set_user_permissions,
      if_owner: rule.if_owner,
      match: rule.match,
      select_condition: rule.select_condition,
      delete_condition: rule.delete_condition,
      amend_condition: rule.amend_condition
    });
    setEditFormVisible(true);
  };

  // 删除权限规则
  const handleDeletePermissionRule = async (ruleId: number) => {
    try {
      await deletePermissionRule(ruleId);
      message.success('权限规则删除成功');
      loadData();
    } catch (error: any) {
      message.error(error.message || '删除失败');
    }
  };

  // 保存权限规则
  const handleSavePermissionRule = async () => {
    try {
      const values = await form.validateFields();
      
      const ruleData = {
        role_id: role.id,
        document_type: values.document_type,
        permission_level: values.permission_level || 0,
        read: values.read || false,
        write: values.write || false,
        create: values.create || false,
        delete: values.delete || false,
        submit: values.submit || false,
        cancel: values.cancel || false,
        amend: values.amend || false,
        print: values.print || false,
        email: values.email || false,
        import: values.import || false,
        export: values.export || false,
        share: values.share || false,
        report: values.report || false,
        set_user_permissions: values.set_user_permissions || false,
        if_owner: values.if_owner || false,
        match: values.match,
        select_condition: values.select_condition,
        delete_condition: values.delete_condition,
        amend_condition: values.amend_condition
      };

      if (editingRule) {
        await updatePermissionRule(editingRule.id, ruleData);
        message.success('权限规则更新成功');
      } else {
        await createPermissionRule(ruleData);
        message.success('权限规则创建成功');
      }

      setEditFormVisible(false);
      loadData();
      onSave?.();
    } catch (error: any) {
      if (error.errorFields) {
        message.error('请检查表单输入');
      } else {
        message.error(error.message || '保存失败');
      }
    }
  };

  // DocType权限表格列定义
  const docTypeColumns: ColumnsType<DocType> = [
    {
      title: 'DocType',
      key: 'name',
      render: (_, record: DocType) => (
        <Space>
          <DatabaseOutlined style={{ color: '#1890ff' }} />
          <div>
            <div style={{ fontWeight: 500 }}>{record.name}</div>
            <div style={{ fontSize: '12px', color: '#999' }}>{record.label}</div>
          </div>
        </Space>
      ),
    },
    {
      title: '权限规则',
      key: 'permissions',
      render: (_, record: DocType) => {
        const rules = getDocTypePermissions(record.name);
        return (
          <Space wrap>
            {rules.length > 0 ? (
              rules.map(rule => (
                <Tag key={rule.id} color="blue">
                  级别 {rule.permission_level}
                </Tag>
              ))
            ) : (
              <Text type="secondary">未配置</Text>
            )}
          </Space>
        );
      },
    },
    {
      title: '操作',
      key: 'action',
      width: 160,
      render: (_, record: DocType) => (
        <Space>
          <Button
            type="text"
            size="small"
            icon={<PlusOutlined />}
            onClick={() => handleCreatePermissionRule(record)}
          >
            添加权限
          </Button>
        </Space>
      ),
    },
  ];

  // 权限规则表格列定义
  const permissionRuleColumns: ColumnsType<PermissionRule> = [
    {
      title: 'DocType',
      dataIndex: 'document_type',
      key: 'document_type',
    },
    {
      title: '权限级别',
      dataIndex: 'permission_level',
      key: 'permission_level',
      render: (level: number) => <Tag color="blue">级别 {level}</Tag>,
    },
    {
      title: '权限',
      key: 'permissions',
      render: (_, record: PermissionRule) => (
        <Space wrap>
          {PERMISSION_TYPES.filter(type => record[type as keyof PermissionRule] === true).map(type => (
            <Tag key={type} icon={getPermissionIcon(type)} color="green">
              {PERMISSION_LABELS[type]}
            </Tag>
          ))}
        </Space>
      ),
    },
    {
      title: '条件',
      key: 'conditions',
      render: (_, record: PermissionRule) => (
        <Space direction="vertical" size="small">
          {record.if_owner && <Tag color="orange">仅限所有者</Tag>}
          {record.match && <Tag color="purple">匹配: {record.match}</Tag>}
        </Space>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      render: (_, record: PermissionRule) => (
        <Space>
          <Button
            type="text"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEditPermissionRule(record)}
          />
          <Button
            type="text"
            size="small"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDeletePermissionRule(record.id)}
          />
        </Space>
      ),
    },
  ];

  // 选项卡配置
  const tabItems = [
    {
      key: 'doctype-permissions',
      label: (
        <span>
          <DatabaseOutlined />
          DocType权限
        </span>
      ),
      children: (
        <div>
          <Alert
            message="DocType权限配置"
            description="为该角色配置不同文档类型的权限。可以设置不同的权限级别和条件。"
            type="info"
            showIcon
            style={{ marginBottom: 16 }}
          />
          
          <Collapse
            ghost
            items={Object.keys(docTypesByModule).map(module => ({
              key: module,
              label: (
                <Space>
                  <DatabaseOutlined />
                  <strong>{module}</strong>
                  <Tag color="blue">{docTypesByModule[module].length}</Tag>
                </Space>
              ),
              children: (
                <Table
                  columns={docTypeColumns}
                  dataSource={docTypesByModule[module]}
                  rowKey="id"
                  size="small"
                  pagination={false}
                />
              )
            }))}
            defaultActiveKey={Object.keys(docTypesByModule)}
          />
        </div>
      ),
    },
    {
      key: 'permission-rules',
      label: (
        <span>
          <SafetyCertificateOutlined />
          权限规则 ({permissionRules?.length || 0})
        </span>
      ),
      children: (
        <div>
          <Alert
            message="权限规则管理"
            description="查看和编辑该角色的所有权限规则。每个规则定义了对特定DocType的访问权限。"
            type="info"
            showIcon
            style={{ marginBottom: 16 }}
          />
          
          <Table
            columns={permissionRuleColumns}
            dataSource={permissionRules}
            rowKey="id"
            loading={loading}
            pagination={{
              pageSize: 10,
              showSizeChanger: true,
              showQuickJumper: true,
              showTotal: (total, range) =>
                `第 ${range[0]}-${range[1]} 条/共 ${total} 条`,
            }}
          />
        </div>
      ),
    },
  ];

  return (
    <>
      <Modal
        title={
          <Space>
            <SafetyCertificateOutlined />
            配置角色权限 - {role.name}
          </Space>
        }
        open={visible}
        onCancel={onClose}
        width={1200}
        footer={[
          <Button key="cancel" onClick={onClose}>
            关闭
          </Button>,
        ]}
        destroyOnClose
      >
        <Spin spinning={loading}>
          <Tabs
            activeKey={activeTab}
            onChange={setActiveTab}
            items={tabItems}
          />
        </Spin>
      </Modal>

      {/* 权限规则编辑模态框 */}
      <Modal
        title={editingRule ? '编辑权限规则' : '新建权限规则'}
        open={editFormVisible}
        onOk={handleSavePermissionRule}
        onCancel={() => setEditFormVisible(false)}
        width={800}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={{
            permission_level: 0,
            if_owner: false
          }}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label="DocType"
                name="document_type"
                rules={[{ required: true, message: '请输入DocType名称' }]}
              >
                <Input disabled={!!editingRule} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label="权限级别"
                name="permission_level"
                rules={[{ required: true, message: '请选择权限级别' }]}
              >
                <Select>
                  <Option value={0}>文档级 (0)</Option>
                  {[1, 2, 3, 4, 5].map(level => (
                    <Option key={level} value={level}>字段级 ({level})</Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Divider>权限设置</Divider>

          <Row gutter={[16, 16]}>
            {PERMISSION_TYPES.map(type => (
              <Col span={6} key={type}>
                <Form.Item name={type} valuePropName="checked">
                  <Space>
                    <Checkbox />
                    {getPermissionIcon(type)}
                    {PERMISSION_LABELS[type]}
                  </Space>
                </Form.Item>
              </Col>
            ))}
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="set_user_permissions" valuePropName="checked">
                <Space>
                  <Checkbox />
                  <BulbOutlined />
                  设置用户权限
                </Space>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="if_owner" valuePropName="checked">
                <Space>
                  <Checkbox />
                  <SafetyCertificateOutlined />
                  仅限所有者
                </Space>
              </Form.Item>
            </Col>
          </Row>

          <Divider>条件设置</Divider>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label="匹配条件"
                name="match"
              >
                <Input placeholder="字段匹配条件，如: company=user.company" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label="选择条件"
                name="select_condition"
              >
                <Input placeholder="SQL WHERE 条件" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label="删除条件"
                name="delete_condition"
              >
                <Input placeholder="删除时的条件检查" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label="修正条件"
                name="amend_condition"
              >
                <Input placeholder="修正时的条件检查" />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>
    </>
  );
};

export default RolePermissionConfig;