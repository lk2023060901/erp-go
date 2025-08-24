/**
 * 仪表板页面
 */

import {
  Card,
  Row,
  Col,
  Statistic,
  Typography,
  Progress,
  Table,
  Tag,
  Space,
  Avatar,
  List,
} from 'antd';
import {
  UserOutlined,
  TeamOutlined,
  SafetyCertificateOutlined,
  ApartmentOutlined,
  RiseOutlined,
  ClockCircleOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';

const { Title, Text } = Typography;

interface UserActivity {
  key: string;
  user: string;
  action: string;
  time: string;
  status: 'success' | 'warning' | 'error';
}

interface SystemMetric {
  title: string;
  value: number;
  suffix?: string;
  precision?: number;
  trend?: 'up' | 'down';
  trendValue?: number;
}

export function Dashboard() {
  // 系统统计数据
  const systemStats: SystemMetric[] = [
    {
      title: '总用户数',
      value: 1284,
      trend: 'up',
      trendValue: 12.5,
    },
    {
      title: '活跃用户',
      value: 987,
      trend: 'up',
      trendValue: 8.2,
    },
    {
      title: '角色数量',
      value: 25,
      trend: 'up',
      trendValue: 4.0,
    },
    {
      title: '权限数量',
      value: 156,
    },
  ];

  // 用户活动数据
  const userActivities: UserActivity[] = [
    {
      key: '1',
      user: '张三',
      action: '登录系统',
      time: '2分钟前',
      status: 'success',
    },
    {
      key: '2',
      user: '李四',
      action: '修改用户权限',
      time: '5分钟前',
      status: 'warning',
    },
    {
      key: '3',
      user: '王五',
      action: '创建新角色',
      time: '10分钟前',
      status: 'success',
    },
    {
      key: '4',
      user: '赵六',
      action: '登录失败',
      time: '15分钟前',
      status: 'error',
    },
    {
      key: '5',
      user: '钱七',
      action: '更新组织架构',
      time: '20分钟前',
      status: 'success',
    },
  ];

  // 系统性能数据
  const systemPerformance = [
    { name: 'CPU使用率', value: 65, status: 'active' },
    { name: '内存使用率', value: 78, status: 'active' },
    { name: '磁盘使用率', value: 45, status: 'active' },
    { name: '网络带宽', value: 32, status: 'active' },
  ];

  // 快捷操作
  const quickActions = [
    {
      title: '新建用户',
      description: '添加新用户到系统',
      icon: <UserOutlined />,
      color: '#1890ff',
    },
    {
      title: '创建角色',
      description: '定义新的用户角色',
      icon: <TeamOutlined />,
      color: '#52c41a',
    },
    {
      title: '权限配置',
      description: '配置系统权限规则',
      icon: <SafetyCertificateOutlined />,
      color: '#faad14',
    },
    {
      title: '组织管理',
      description: '管理组织架构',
      icon: <ApartmentOutlined />,
      color: '#eb2f96',
    },
  ];

  const activityColumns: ColumnsType<UserActivity> = [
    {
      title: '用户',
      dataIndex: 'user',
      key: 'user',
      render: (name: string) => (
        <Space>
          <Avatar size="small" icon={<UserOutlined />} />
          <Text>{name}</Text>
        </Space>
      ),
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
    },
    {
      title: '时间',
      dataIndex: 'time',
      key: 'time',
      render: (time: string) => (
        <Text type="secondary">
          <ClockCircleOutlined style={{ marginRight: 4 }} />
          {time}
        </Text>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const statusConfig = {
          success: { color: 'success', icon: <CheckCircleOutlined />, text: '成功' },
          warning: { color: 'warning', icon: <ExclamationCircleOutlined />, text: '警告' },
          error: { color: 'error', icon: <ExclamationCircleOutlined />, text: '失败' },
        };
        const config = statusConfig[status as keyof typeof statusConfig];
        return (
          <Tag color={config.color} icon={config.icon}>
            {config.text}
          </Tag>
        );
      },
    },
  ];

  const getProgressColor = (value: number) => {
    if (value >= 80) return '#ff4d4f';
    if (value >= 60) return '#faad14';
    return '#52c41a';
  };

  return (
    <div>
      <div style={{ marginBottom: 24 }}>
        <Title level={2} style={{ margin: 0, color: '#1f2937' }}>
          仪表板
        </Title>
        <Text type="secondary">
          欢迎回到ERP管理系统，这里是系统概览
        </Text>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        {systemStats.map((stat, index) => {
          const icons = [
            <UserOutlined style={{ color: '#1890ff' }} />,
            <RiseOutlined style={{ color: '#52c41a' }} />,
            <TeamOutlined style={{ color: '#faad14' }} />,
            <SafetyCertificateOutlined style={{ color: '#eb2f96' }} />,
          ];
          
          return (
            <Col xs={24} sm={12} lg={6} key={index}>
              <Card>
                <Statistic
                  title={stat.title}
                  value={stat.value}
                  precision={stat.precision}
                  suffix={stat.suffix}
                  prefix={icons[index]}
                />
                {stat.trend && (
                  <div style={{ marginTop: 8 }}>
                    <Text
                      type={stat.trend === 'up' ? 'success' : 'danger'}
                      style={{ fontSize: 12 }}
                    >
                      {stat.trend === 'up' ? '↗' : '↘'} {stat.trendValue}%
                    </Text>
                    <Text type="secondary" style={{ marginLeft: 8, fontSize: 12 }}>
                      较上月
                    </Text>
                  </div>
                )}
              </Card>
            </Col>
          );
        })}
      </Row>

      <Row gutter={[16, 16]}>
        {/* 系统性能 */}
        <Col xs={24} lg={8}>
          <Card title="系统性能" style={{ height: '100%' }}>
            <Space direction="vertical" style={{ width: '100%' }} size="large">
              {systemPerformance.map((item, index) => (
                <div key={index}>
                  <div style={{ marginBottom: 8 }}>
                    <Text strong>{item.name}</Text>
                    <Text style={{ float: 'right' }}>{item.value}%</Text>
                  </div>
                  <Progress
                    percent={item.value}
                    strokeColor={getProgressColor(item.value)}
                    size="small"
                  />
                </div>
              ))}
            </Space>
          </Card>
        </Col>

        {/* 快捷操作 */}
        <Col xs={24} lg={8}>
          <Card title="快捷操作" style={{ height: '100%' }}>
            <List
              dataSource={quickActions}
              renderItem={(item) => (
                <List.Item
                  style={{ cursor: 'pointer', padding: '12px 0' }}
                  onClick={() => console.log('点击:', item.title)}
                >
                  <List.Item.Meta
                    avatar={
                      <Avatar
                        style={{ backgroundColor: item.color }}
                        icon={item.icon}
                      />
                    }
                    title={
                      <Text strong style={{ color: '#1f2937' }}>
                        {item.title}
                      </Text>
                    }
                    description={
                      <Text type="secondary" style={{ fontSize: 12 }}>
                        {item.description}
                      </Text>
                    }
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>

        {/* 最近活动 */}
        <Col xs={24} lg={8}>
          <Card title="最近活动" style={{ height: '100%' }}>
            <Table
              dataSource={userActivities}
              columns={activityColumns}
              pagination={false}
              size="small"
              scroll={{ y: 300 }}
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
}

export default Dashboard;