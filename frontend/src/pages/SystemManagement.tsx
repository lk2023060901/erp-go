import React, { useState, useEffect } from 'react';
import {
  Card,
  Row,
  Col,
  Tabs,
  Form,
  Input,
  InputNumber,
  Switch,
  // Select,
  Button,
  Space,
  message,
  Divider,
  Alert,
  Table,
  Tag,
  Modal,
  // Upload,
  Progress,
  Statistic,
  Typography,
  // List,
  // Badge,
  // Tooltip,
} from 'antd';
import {
  SettingOutlined,
  DatabaseOutlined,
  // SafetyCertificateOutlined,
  // MailOutlined,
  // BellOutlined,
  FileTextOutlined,
  UploadOutlined,
  DownloadOutlined,
  ReloadOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  WarningOutlined,
} from '@ant-design/icons';
import type { TabsProps, TableColumnsType } from 'antd';
import dayjs from 'dayjs';

const { TextArea } = Input;
// const { Option } = Select;
const { Title } = Typography;
// const { Text } = Typography;

// 系统配置接口
interface SystemConfig {
  app_name: string;
  app_version: string;
  app_description: string;
  max_upload_size: number;
  session_timeout: number;
  password_min_length: number;
  password_require_special: boolean;
  password_require_number: boolean;
  password_require_uppercase: boolean;
  enable_email_notification: boolean;
  enable_sms_notification: boolean;
  enable_audit_log: boolean;
  max_login_attempts: number;
  lockout_duration: number;
}

// 系统状态接口
interface SystemStatus {
  uptime: number;
  cpu_usage: number;
  memory_usage: number;
  disk_usage: number;
  active_users: number;
  total_users: number;
  database_size: string;
  cache_size: string;
  last_backup: string;
}

// 系统日志接口
interface SystemLog {
  id: number;
  level: 'INFO' | 'WARN' | 'ERROR' | 'DEBUG';
  message: string;
  source: string;
  timestamp: string;
  user_id?: number;
  user_name?: string;
  ip_address?: string;
}

// 备份记录接口
interface BackupRecord {
  id: number;
  filename: string;
  size: string;
  type: 'auto' | 'manual';
  status: 'success' | 'failed' | 'running';
  created_at: string;
  created_by?: string;
}

const SystemManagement: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [systemConfig, setSystemConfig] = useState<SystemConfig>({
    app_name: 'ERP系统',
    app_version: '1.0.0',
    app_description: '企业资源计划管理系统',
    max_upload_size: 10,
    session_timeout: 30,
    password_min_length: 8,
    password_require_special: true,
    password_require_number: true,
    password_require_uppercase: true,
    enable_email_notification: true,
    enable_sms_notification: false,
    enable_audit_log: true,
    max_login_attempts: 5,
    lockout_duration: 15,
  });
  const [systemStatus, setSystemStatus] = useState<SystemStatus>({
    uptime: 86400,
    cpu_usage: 45.2,
    memory_usage: 67.8,
    disk_usage: 34.5,
    active_users: 15,
    total_users: 120,
    database_size: '245.6 MB',
    cache_size: '128.3 MB',
    last_backup: '2024-08-23T02:00:00Z',
  });
  const [systemLogs, setSystemLogs] = useState<SystemLog[]>([]);
  const [backupRecords, setBackupRecords] = useState<BackupRecord[]>([]);
  const [form] = Form.useForm();
  const [activeTab, setActiveTab] = useState('config');

  useEffect(() => {
    loadSystemLogs();
    loadBackupRecords();
    // 设置定时器更新系统状态
    const timer = setInterval(() => {
      updateSystemStatus();
    }, 30000); // 每30秒更新一次

    return () => clearInterval(timer);
  }, []);

  const loadSystemLogs = async () => {
    // 模拟加载系统日志
    const mockLogs: SystemLog[] = [
      {
        id: 1,
        level: 'INFO',
        message: '用户登录成功',
        source: 'auth.service',
        timestamp: dayjs().subtract(5, 'minute').toISOString(),
        user_id: 1,
        user_name: 'admin',
        ip_address: '192.168.1.100',
      },
      {
        id: 2,
        level: 'WARN',
        message: '数据库连接池使用率过高',
        source: 'database.pool',
        timestamp: dayjs().subtract(15, 'minute').toISOString(),
      },
      {
        id: 3,
        level: 'ERROR',
        message: '邮件发送失败',
        source: 'mail.service',
        timestamp: dayjs().subtract(30, 'minute').toISOString(),
      },
      {
        id: 4,
        level: 'INFO',
        message: '定时备份任务开始',
        source: 'backup.scheduler',
        timestamp: dayjs().subtract(2, 'hour').toISOString(),
      },
      {
        id: 5,
        level: 'INFO',
        message: '定时备份任务完成',
        source: 'backup.scheduler',
        timestamp: dayjs().subtract(2, 'hour').add(15, 'minute').toISOString(),
      },
    ];
    setSystemLogs(mockLogs);
  };

  const loadBackupRecords = async () => {
    // 模拟加载备份记录
    const mockBackups: BackupRecord[] = [
      {
        id: 1,
        filename: 'backup_2024-08-23_02-00-00.sql',
        size: '245.6 MB',
        type: 'auto',
        status: 'success',
        created_at: dayjs().subtract(6, 'hour').toISOString(),
        created_by: 'system',
      },
      {
        id: 2,
        filename: 'backup_2024-08-22_14-30-15.sql',
        size: '243.2 MB',
        type: 'manual',
        status: 'success',
        created_at: dayjs().subtract(1, 'day').subtract(6, 'hour').toISOString(),
        created_by: 'admin',
      },
      {
        id: 3,
        filename: 'backup_2024-08-22_02-00-00.sql',
        size: '241.8 MB',
        type: 'auto',
        status: 'success',
        created_at: dayjs().subtract(1, 'day').subtract(18, 'hour').toISOString(),
        created_by: 'system',
      },
    ];
    setBackupRecords(mockBackups);
  };

  const updateSystemStatus = async () => {
    // 模拟更新系统状态
    setSystemStatus(prev => ({
      ...prev,
      cpu_usage: Math.random() * 100,
      memory_usage: Math.random() * 100,
      active_users: Math.floor(Math.random() * 50) + 10,
    }));
  };

  const handleSaveConfig = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);
      
      // 模拟保存配置
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setSystemConfig({ ...systemConfig, ...values });
      message.success('系统配置保存成功');
    } catch (error) {
      message.error('保存配置失败');
    } finally {
      setLoading(false);
    }
  };

  const handleBackup = async () => {
    try {
      setLoading(true);
      // 模拟创建备份
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      const newBackup: BackupRecord = {
        id: Date.now(),
        filename: `backup_${dayjs().format('YYYY-MM-DD_HH-mm-ss')}.sql`,
        size: '246.3 MB',
        type: 'manual',
        status: 'success',
        created_at: dayjs().toISOString(),
        created_by: 'admin',
      };
      
      setBackupRecords([newBackup, ...backupRecords]);
      message.success('数据备份创建成功');
    } catch (error) {
      message.error('备份失败');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteBackup = (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个备份文件吗？删除后无法恢复。',
      okText: '确定',
      cancelText: '取消',
      onOk: () => {
        setBackupRecords(backupRecords.filter(record => record.id !== id));
        message.success('备份文件删除成功');
      },
    });
  };

  const formatUptime = (seconds: number) => {
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${days}天 ${hours}时 ${minutes}分`;
  };

  const getLogLevelColor = (level: string) => {
    const colors = {
      INFO: 'blue',
      WARN: 'orange',
      ERROR: 'red',
      DEBUG: 'gray',
    };
    return colors[level as keyof typeof colors] || 'default';
  };

  const getLogLevelIcon = (level: string) => {
    const icons = {
      INFO: <CheckCircleOutlined />,
      WARN: <WarningOutlined />,
      ERROR: <ExclamationCircleOutlined />,
      DEBUG: <ClockCircleOutlined />,
    };
    return icons[level as keyof typeof icons];
  };

  const logColumns: TableColumnsType<SystemLog> = [
    {
      title: '级别',
      dataIndex: 'level',
      width: 80,
      render: (level: string) => (
        <Tag color={getLogLevelColor(level)} icon={getLogLevelIcon(level)}>
          {level}
        </Tag>
      ),
    },
    {
      title: '消息',
      dataIndex: 'message',
      ellipsis: true,
    },
    {
      title: '来源',
      dataIndex: 'source',
      width: 150,
      render: (source: string) => <code>{source}</code>,
    },
    {
      title: '用户',
      width: 120,
      render: (_, record) => 
        record.user_name ? `${record.user_name}` : '系统',
    },
    {
      title: 'IP地址',
      dataIndex: 'ip_address',
      width: 130,
      render: (ip: string) => ip || '-',
    },
    {
      title: '时间',
      dataIndex: 'timestamp',
      width: 160,
      render: (time: string) => dayjs(time).format('MM-DD HH:mm:ss'),
    },
  ];

  const backupColumns: TableColumnsType<BackupRecord> = [
    {
      title: '文件名',
      dataIndex: 'filename',
      render: (filename: string) => <code>{filename}</code>,
    },
    {
      title: '大小',
      dataIndex: 'size',
      width: 100,
    },
    {
      title: '类型',
      dataIndex: 'type',
      width: 80,
      render: (type: string) => (
        <Tag color={type === 'auto' ? 'blue' : 'green'}>
          {type === 'auto' ? '自动' : '手动'}
        </Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 80,
      render: (status: string) => (
        <Tag color={status === 'success' ? 'success' : status === 'failed' ? 'error' : 'processing'}>
          {status === 'success' ? '成功' : status === 'failed' ? '失败' : '运行中'}
        </Tag>
      ),
    },
    {
      title: '创建者',
      dataIndex: 'created_by',
      width: 100,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      width: 160,
      render: (time: string) => dayjs(time).format('MM-DD HH:mm:ss'),
    },
    {
      title: '操作',
      width: 150,
      render: (_, record) => (
        <Space size="small">
          <Button 
            type="link" 
            size="small"
            icon={<DownloadOutlined />}
          >
            下载
          </Button>
          <Button 
            type="link" 
            size="small" 
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDeleteBackup(record.id)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ];

  const tabItems: TabsProps['items'] = [
    {
      key: 'config',
      label: (
        <span>
          <SettingOutlined />
          系统配置
        </span>
      ),
      children: (
        <Card>
          <Form
            form={form}
            layout="vertical"
            initialValues={systemConfig}
            onValuesChange={(_, allValues) => setSystemConfig({ ...systemConfig, ...allValues })}
          >
            <Title level={4}>基本设置</Title>
            <Row gutter={16}>
              <Col span={8}>
                <Form.Item
                  label="应用名称"
                  name="app_name"
                  rules={[{ required: true, message: '请输入应用名称' }]}
                >
                  <Input placeholder="请输入应用名称" />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item
                  label="应用版本"
                  name="app_version"
                  rules={[{ required: true, message: '请输入应用版本' }]}
                >
                  <Input placeholder="如：1.0.0" />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item
                  label="最大上传大小(MB)"
                  name="max_upload_size"
                  rules={[{ required: true, message: '请输入最大上传大小' }]}
                >
                  <InputNumber min={1} max={1000} style={{ width: '100%' }} />
                </Form.Item>
              </Col>
            </Row>

            <Form.Item
              label="应用描述"
              name="app_description"
            >
              <TextArea rows={3} placeholder="请输入应用描述" />
            </Form.Item>

            <Divider />
            <Title level={4}>安全设置</Title>
            <Row gutter={16}>
              <Col span={8}>
                <Form.Item
                  label="会话超时(分钟)"
                  name="session_timeout"
                  rules={[{ required: true, message: '请输入会话超时时间' }]}
                >
                  <InputNumber min={5} max={1440} style={{ width: '100%' }} />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item
                  label="密码最小长度"
                  name="password_min_length"
                  rules={[{ required: true, message: '请输入密码最小长度' }]}
                >
                  <InputNumber min={6} max={20} style={{ width: '100%' }} />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item
                  label="最大登录尝试次数"
                  name="max_login_attempts"
                  rules={[{ required: true, message: '请输入最大登录尝试次数' }]}
                >
                  <InputNumber min={3} max={10} style={{ width: '100%' }} />
                </Form.Item>
              </Col>
            </Row>

            <Row gutter={16}>
              <Col span={8}>
                <Form.Item
                  label="锁定时长(分钟)"
                  name="lockout_duration"
                  rules={[{ required: true, message: '请输入锁定时长' }]}
                >
                  <InputNumber min={5} max={60} style={{ width: '100%' }} />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item
                  label="密码需要特殊字符"
                  name="password_require_special"
                  valuePropName="checked"
                >
                  <Switch />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item
                  label="密码需要数字"
                  name="password_require_number"
                  valuePropName="checked"
                >
                  <Switch />
                </Form.Item>
              </Col>
            </Row>

            <Form.Item
              label="密码需要大写字母"
              name="password_require_uppercase"
              valuePropName="checked"
            >
              <Switch />
            </Form.Item>

            <Divider />
            <Title level={4}>通知设置</Title>
            <Row gutter={16}>
              <Col span={8}>
                <Form.Item
                  label="启用邮件通知"
                  name="enable_email_notification"
                  valuePropName="checked"
                >
                  <Switch />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item
                  label="启用短信通知"
                  name="enable_sms_notification"
                  valuePropName="checked"
                >
                  <Switch />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item
                  label="启用审计日志"
                  name="enable_audit_log"
                  valuePropName="checked"
                >
                  <Switch />
                </Form.Item>
              </Col>
            </Row>

            <Divider />
            <Button 
              type="primary" 
              size="large"
              loading={loading}
              onClick={handleSaveConfig}
            >
              保存配置
            </Button>
          </Form>
        </Card>
      ),
    },
    {
      key: 'status',
      label: (
        <span>
          <DatabaseOutlined />
          系统状态
        </span>
      ),
      children: (
        <div>
          <Row gutter={[16, 16]}>
            <Col span={6}>
              <Card>
                <Statistic
                  title="系统运行时间"
                  value={formatUptime(systemStatus.uptime)}
                  prefix={<ClockCircleOutlined />}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="在线用户"
                  value={systemStatus.active_users}
                  suffix={`/ ${systemStatus.total_users}`}
                  prefix={<CheckCircleOutlined style={{ color: '#52c41a' }} />}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="数据库大小"
                  value={systemStatus.database_size}
                  prefix={<DatabaseOutlined />}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="缓存大小"
                  value={systemStatus.cache_size}
                  prefix={<DatabaseOutlined />}
                />
              </Card>
            </Col>
          </Row>

          <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
            <Col span={8}>
              <Card title="CPU使用率" extra={`${systemStatus.cpu_usage.toFixed(1)}%`}>
                <Progress 
                  percent={systemStatus.cpu_usage} 
                  status={systemStatus.cpu_usage > 80 ? 'exception' : 'active'}
                  showInfo={false}
                />
              </Card>
            </Col>
            <Col span={8}>
              <Card title="内存使用率" extra={`${systemStatus.memory_usage.toFixed(1)}%`}>
                <Progress 
                  percent={systemStatus.memory_usage} 
                  status={systemStatus.memory_usage > 80 ? 'exception' : 'active'}
                  showInfo={false}
                />
              </Card>
            </Col>
            <Col span={8}>
              <Card title="磁盘使用率" extra={`${systemStatus.disk_usage.toFixed(1)}%`}>
                <Progress 
                  percent={systemStatus.disk_usage} 
                  status={systemStatus.disk_usage > 80 ? 'exception' : 'active'}
                  showInfo={false}
                />
              </Card>
            </Col>
          </Row>

          <Card title="最后备份时间" style={{ marginTop: 16 }}>
            <Alert
              message={`最后一次备份: ${dayjs(systemStatus.last_backup).format('YYYY-MM-DD HH:mm:ss')}`}
              description={`距离现在: ${dayjs().diff(dayjs(systemStatus.last_backup), 'hour')} 小时前`}
              type="info"
              showIcon
            />
          </Card>
        </div>
      ),
    },
    {
      key: 'logs',
      label: (
        <span>
          <FileTextOutlined />
          系统日志
        </span>
      ),
      children: (
        <Card
          title="系统日志"
          extra={
            <Button 
              icon={<ReloadOutlined />} 
              onClick={loadSystemLogs}
            >
              刷新
            </Button>
          }
        >
          <Table
            columns={logColumns}
            dataSource={systemLogs}
            rowKey="id"
            pagination={{
              pageSize: 10,
              showQuickJumper: true,
              showSizeChanger: true,
            }}
            scroll={{ x: 1000 }}
          />
        </Card>
      ),
    },
    {
      key: 'backup',
      label: (
        <span>
          <DatabaseOutlined />
          数据备份
        </span>
      ),
      children: (
        <Card
          title="数据备份"
          extra={
            <Space>
              <Button 
                type="primary"
                icon={<UploadOutlined />}
                loading={loading}
                onClick={handleBackup}
              >
                创建备份
              </Button>
              <Button 
                icon={<ReloadOutlined />} 
                onClick={loadBackupRecords}
              >
                刷新
              </Button>
            </Space>
          }
        >
          <Alert
            message="数据备份说明"
            description="系统会每天凌晨2点自动创建数据备份。您也可以手动创建备份。备份文件包含完整的数据库数据。"
            type="info"
            showIcon
            style={{ marginBottom: 16 }}
          />
          <Table
            columns={backupColumns}
            dataSource={backupRecords}
            rowKey="id"
            pagination={{
              pageSize: 10,
              showQuickJumper: true,
            }}
            scroll={{ x: 1000 }}
          />
        </Card>
      ),
    },
  ];

  return (
    <div>
      <Card title="系统管理">
        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          items={tabItems}
        />
      </Card>
    </div>
  );
};

export default SystemManagement;