/**
 * 用户过滤器面板组件
 */

import React, { useState } from 'react';
import {
  Card,
  Button,
  Select,
  Input,
  DatePicker,
  Space,
  Row,
  Col,
  Modal,
  Form,
  App,
  Tag,
  Divider,
  Tooltip
} from 'antd';
import {
  FilterOutlined,
  SaveOutlined,
  ClearOutlined,
  CheckOutlined,
  // DeleteOutlined,
  SortAscendingOutlined,
  PlusOutlined,
  CloseOutlined
} from '@ant-design/icons';
import dayjs from 'dayjs';
import type { Dayjs } from 'dayjs';
import type { SavedFilter } from '../services/filterService';

const { Option } = Select;

// 过滤条件接口
export interface FilterCondition {
  field: string;
  operator: string;
  value: any;
}

// 过滤器配置接口
export interface FilterConfig {
  conditions: FilterCondition[];
  logic: 'AND' | 'OR';
}

// 排序配置接口
export interface SortConfig {
  field: string;
  direction: 'asc' | 'desc';
}

// 保存的过滤器接口已从 filterService 导入

// 组件属性接口
export interface UserFilterPanelProps {
  visible: boolean;
  onToggle: () => void;
  onApplyFilter: (filterConfig: FilterConfig, sortConfig?: SortConfig) => void;
  savedFilters: SavedFilter[];
  onSaveFilter: (name: string, filterConfig: FilterConfig, sortConfig?: SortConfig, isDefault?: boolean) => void;
  onDeleteFilter: (filterId: number) => void;
  onLoadFilter: (filter: SavedFilter) => void;
  loading?: boolean;
}

// 用户字段配置
const USER_FIELDS = [
  { value: 'username', label: '用户名', type: 'string' },
  { value: 'email', label: '邮箱', type: 'string' },
  { value: 'first_name', label: '名字', type: 'string' },
  { value: 'last_name', label: '姓氏', type: 'string' },
  { value: 'phone', label: '手机号', type: 'string' },
  { value: 'gender', label: '性别', type: 'select', options: [
    { value: 'M', label: '男' },
    { value: 'F', label: '女' },
    { value: 'O', label: '其他' }
  ]},
  { value: 'is_active', label: '状态', type: 'select', options: [
    { value: 'true', label: '启用' },
    { value: 'false', label: '禁用' }
  ]},
  { value: 'created_at', label: '创建时间', type: 'date' },
  { value: 'updated_at', label: '更新时间', type: 'date' },
  { value: 'last_login_time', label: '最后登录', type: 'date' }
];

// 操作符配置
const OPERATORS = [
  { value: 'equals', label: '等于', types: ['string', 'select', 'boolean'] },
  { value: 'contains', label: '包含', types: ['string'] },
  { value: 'starts_with', label: '开头是', types: ['string'] },
  { value: 'ends_with', label: '结尾是', types: ['string'] },
  { value: 'greater_than', label: '大于', types: ['date', 'number'] },
  { value: 'less_than', label: '小于', types: ['date', 'number'] },
  { value: 'between', label: '介于', types: ['date', 'number'] },
  { value: 'in', label: '属于', types: ['string', 'select'] },
  { value: 'not_in', label: '不属于', types: ['string', 'select'] }
];

// 排序字段配置
const SORT_FIELDS = [
  { value: 'username', label: '用户名' },
  { value: 'email', label: '邮箱' },
  { value: 'first_name', label: '名字' },
  { value: 'last_name', label: '姓氏' },
  { value: 'created_at', label: '创建时间' },
  { value: 'updated_at', label: '更新时间' },
  { value: 'last_login_time', label: '最后登录时间' }
];

export const UserFilterPanel: React.FC<UserFilterPanelProps> = ({
  visible,
  // onToggle,
  onApplyFilter,
  savedFilters,
  onSaveFilter,
  onDeleteFilter,
  onLoadFilter,
  loading = false
}) => {
  const { message } = App.useApp();
  
  // 状态管理
  const [filterConfig, setFilterConfig] = useState<FilterConfig>({
    conditions: [{ field: 'username', operator: 'equals', value: '' }],
    logic: 'AND'
  });
  const [sortConfig, setSortConfig] = useState<SortConfig>({
    field: 'username',
    direction: 'asc'
  });
  const [saveModalVisible, setSaveModalVisible] = useState(false);
  const [activeFilter, setActiveFilter] = useState<number | null>(null);
  
  // 表单实例
  const [saveForm] = Form.useForm();

  // 添加过滤条件
  const addCondition = () => {
    setFilterConfig(prev => ({
      ...prev,
      conditions: [...prev.conditions, { field: 'username', operator: 'equals', value: '' }]
    }));
  };

  // 移除过滤条件
  const removeCondition = (index: number) => {
    if (filterConfig.conditions.length <= 1) {
      message.warning('至少需要保留一个过滤条件');
      return;
    }
    
    setFilterConfig(prev => ({
      ...prev,
      conditions: prev.conditions.filter((_, i) => i !== index)
    }));
  };

  // 更新过滤条件
  const updateCondition = (index: number, field: keyof FilterCondition, value: any) => {
    setFilterConfig(prev => ({
      ...prev,
      conditions: prev.conditions.map((condition, i) => 
        i === index ? { ...condition, [field]: value } : condition
      )
    }));
  };

  // 清空过滤器
  const clearFilters = () => {
    setFilterConfig({
      conditions: [{ field: 'username', operator: 'equals', value: '' }],
      logic: 'AND'
    });
    setSortConfig({
      field: 'username',
      direction: 'asc'
    });
    setActiveFilter(null);
  };

  // 应用过滤器
  const applyFilters = () => {
    // 验证过滤条件
    const validConditions = filterConfig.conditions.filter(condition => {
      if (!condition.field || !condition.operator) return false;
      if (condition.value === '' || condition.value === null || condition.value === undefined) return false;
      return true;
    });

    if (validConditions.length === 0) {
      message.warning('请至少设置一个有效的过滤条件');
      return;
    }

    onApplyFilter({ ...filterConfig, conditions: validConditions }, sortConfig);
    message.success('过滤器已应用');
  };

  // 保存过滤器
  const saveFilter = () => {
    setSaveModalVisible(true);
  };

  // 确认保存过滤器
  const handleSaveFilter = async () => {
    try {
      const values = await saveForm.validateFields();
      onSaveFilter(values.filterName, filterConfig, sortConfig, values.isDefault);
      setSaveModalVisible(false);
      saveForm.resetFields();
      message.success('过滤器保存成功');
    } catch (error) {
      console.error('保存过滤器失败:', error);
    }
  };

  // 加载保存的过滤器
  const loadSavedFilter = (filter: SavedFilter) => {
    setFilterConfig(filter.filter_conditions);
    if (filter.sort_config) {
      setSortConfig(filter.sort_config);
    }
    setActiveFilter(filter.id);
    onLoadFilter(filter);
    message.success(`已加载过滤器: ${filter.filter_name}`);
  };

  // 删除保存的过滤器
  const handleDeleteFilter = (filterId: number, filterName: string) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除过滤器 "${filterName}" 吗？`,
      onOk: () => {
        onDeleteFilter(filterId);
        if (activeFilter === filterId) {
          setActiveFilter(null);
        }
        message.success('过滤器已删除');
      }
    });
  };

  // 获取字段类型
  const getFieldType = (fieldValue: string) => {
    const field = USER_FIELDS.find(f => f.value === fieldValue);
    return field?.type || 'string';
  };

  // 获取字段选项
  const getFieldOptions = (fieldValue: string) => {
    const field = USER_FIELDS.find(f => f.value === fieldValue);
    return field?.options || [];
  };

  // 获取可用操作符
  const getAvailableOperators = (fieldType: string) => {
    return OPERATORS.filter(op => op.types.includes(fieldType));
  };

  // 渲染值输入组件
  const renderValueInput = (condition: FilterCondition, index: number) => {
    const fieldType = getFieldType(condition.field);
    const fieldOptions = getFieldOptions(condition.field);

    switch (fieldType) {
      case 'select':
        return (
          <Select
            value={condition.value}
            onChange={(value) => updateCondition(index, 'value', value)}
            placeholder="请选择"
            style={{ width: '100%' }}
          >
            {fieldOptions.map(option => (
              <Option key={String(option.value)} value={String(option.value)}>
                {option.label}
              </Option>
            ))}
          </Select>
        );
      
      case 'date':
        return (
          <DatePicker
            value={condition.value ? dayjs(condition.value) : null}
            onChange={(date: Dayjs | null) => 
              updateCondition(index, 'value', date ? date.format('YYYY-MM-DD') : '')
            }
            placeholder="请选择日期"
            style={{ width: '100%' }}
          />
        );
      
      default:
        return (
          <Input
            value={condition.value}
            onChange={(e) => updateCondition(index, 'value', e.target.value)}
            placeholder="请输入筛选值"
          />
        );
    }
  };

  if (!visible) {
    return null;
  }

  return (
    <Card 
      className="user-filter-panel"
      style={{ 
        marginBottom: 16,
        backgroundColor: '#f9fafb',
        borderColor: '#e5e7eb'
      }}
    >
      {/* 过滤器头部 */}
      <Row justify="space-between" align="middle" style={{ marginBottom: 20 }}>
        <Col>
          <Space>
            <FilterOutlined />
            <span style={{ fontSize: 16, fontWeight: 500, color: '#1f2937' }}>过滤器设置</span>
          </Space>
        </Col>
        <Col>
          <Space>
            <Button
              size="small"
              icon={<SaveOutlined />}
              onClick={saveFilter}
            >
              保存过滤器
            </Button>
            <Button
              size="small"
              icon={<ClearOutlined />}
              onClick={clearFilters}
            >
              清空过滤器
            </Button>
            <Button
              type="primary"
              size="small"
              icon={<CheckOutlined />}
              onClick={applyFilters}
              loading={loading}
            >
              应用过滤器
            </Button>
          </Space>
        </Col>
      </Row>

      {/* 已保存的过滤器 */}
      {savedFilters.length > 0 && (
        <>
          <div style={{ marginBottom: 12, fontSize: 14, fontWeight: 500, color: '#374151' }}>
            已保存的过滤器
          </div>
          <div style={{ marginBottom: 20 }}>
            <Space wrap>
              {savedFilters.map(filter => (
                <Tag
                  key={filter.id}
                  color={activeFilter === filter.id ? '#1d4ed8' : '#eff6ff'}
                  style={{
                    color: activeFilter === filter.id ? 'white' : '#1d4ed8',
                    cursor: 'pointer',
                    padding: '4px 12px',
                    borderRadius: 16,
                    border: `1px solid ${activeFilter === filter.id ? '#1d4ed8' : '#dbeafe'}`
                  }}
                  onClick={() => loadSavedFilter(filter)}
                >
                  {filter.filter_name}
                  {filter.is_default && <span style={{ marginLeft: 4 }}>（默认）</span>}
                  <CloseOutlined
                    style={{ marginLeft: 6, fontSize: 10 }}
                    onClick={(e) => {
                      e.stopPropagation();
                      handleDeleteFilter(filter.id, filter.filter_name);
                    }}
                  />
                </Tag>
              ))}
            </Space>
          </div>
          <Divider style={{ margin: '12px 0' }} />
        </>
      )}

      {/* 排序配置 */}
      <Row style={{ marginBottom: 20 }}>
        <Col span={24}>
          <Space>
            <SortAscendingOutlined />
            <span style={{ fontSize: 14, fontWeight: 500 }}>排序设置:</span>
            <Select
              value={sortConfig.field}
              onChange={(value) => setSortConfig(prev => ({ ...prev, field: value }))}
              style={{ width: 150 }}
            >
              {SORT_FIELDS.map(field => (
                <Option key={field.value} value={field.value}>
                  {field.label}
                </Option>
              ))}
            </Select>
            <Select
              value={sortConfig.direction}
              onChange={(value) => setSortConfig(prev => ({ ...prev, direction: value }))}
              style={{ width: 100 }}
            >
              <Option value="asc">升序</Option>
              <Option value="desc">降序</Option>
            </Select>
          </Space>
        </Col>
      </Row>

      <Divider style={{ margin: '12px 0' }} />

      {/* 过滤条件 */}
      <div style={{ marginBottom: 20 }}>
        {filterConfig.conditions.map((condition, index) => (
          <Card
            key={index}
            size="small"
            style={{
              backgroundColor: 'white',
              border: '1px solid #e5e7eb',
              borderRadius: 6,
              marginBottom: 12
            }}
          >
            <Row justify="space-between" align="middle" style={{ marginBottom: 12 }}>
              <Col>
                <span style={{
                  fontSize: 12,
                  color: '#6b7280',
                  backgroundColor: '#f3f4f6',
                  padding: '2px 8px',
                  borderRadius: 4
                }}>
                  {index === 0 ? 'WHERE' : filterConfig.logic}
                </span>
              </Col>
              <Col>
                {filterConfig.conditions.length > 1 && (
                  <Tooltip title="删除条件">
                    <Button
                      type="text"
                      size="small"
                      icon={<CloseOutlined />}
                      onClick={() => removeCondition(index)}
                      style={{ color: '#ef4444' }}
                    />
                  </Tooltip>
                )}
              </Col>
            </Row>

            <Row gutter={16}>
              <Col span={8}>
                <div style={{ marginBottom: 4, fontSize: 12, color: '#374151' }}>字段</div>
                <Select
                  value={condition.field}
                  onChange={(value) => {
                    updateCondition(index, 'field', value);
                    // 重置操作符和值
                    const fieldType = getFieldType(value);
                    const availableOperators = getAvailableOperators(fieldType);
                    if (availableOperators.length > 0) {
                      updateCondition(index, 'operator', availableOperators[0].value);
                    }
                    updateCondition(index, 'value', '');
                  }}
                  style={{ width: '100%' }}
                >
                  {USER_FIELDS.map(field => (
                    <Option key={field.value} value={field.value}>
                      {field.label}
                    </Option>
                  ))}
                </Select>
              </Col>
              <Col span={8}>
                <div style={{ marginBottom: 4, fontSize: 12, color: '#374151' }}>条件</div>
                <Select
                  value={condition.operator}
                  onChange={(value) => updateCondition(index, 'operator', value)}
                  style={{ width: '100%' }}
                >
                  {getAvailableOperators(getFieldType(condition.field)).map(operator => (
                    <Option key={operator.value} value={operator.value}>
                      {operator.label}
                    </Option>
                  ))}
                </Select>
              </Col>
              <Col span={8}>
                <div style={{ marginBottom: 4, fontSize: 12, color: '#374151' }}>值</div>
                {renderValueInput(condition, index)}
              </Col>
            </Row>
          </Card>
        ))}
      </div>

      {/* 添加条件按钮 */}
      <Button
        type="dashed"
        icon={<PlusOutlined />}
        onClick={addCondition}
        style={{
          width: '100%',
          borderColor: '#6366f1',
          color: '#6366f1',
          borderStyle: 'dashed',
          borderWidth: 2,
          height: 40
        }}
      >
        添加筛选条件
      </Button>

      {/* 保存过滤器模态框 */}
      <Modal
        title="保存过滤器"
        open={saveModalVisible}
        onCancel={() => {
          setSaveModalVisible(false);
          saveForm.resetFields();
        }}
        footer={null}
        width={400}
      >
        <Form
          form={saveForm}
          layout="vertical"
          onFinish={handleSaveFilter}
        >
          <Form.Item
            name="filterName"
            label="过滤器名称"
            rules={[
              { required: true, message: '请输入过滤器名称' },
              { max: 100, message: '名称不能超过100个字符' }
            ]}
          >
            <Input placeholder="请输入过滤器名称" />
          </Form.Item>

          <Form.Item
            name="isDefault"
            valuePropName="checked"
          >
            <Space>
              <input type="checkbox" />
              <span>设为默认过滤器</span>
            </Space>
          </Form.Item>

          <Row justify="end">
            <Space>
              <Button onClick={() => {
                setSaveModalVisible(false);
                saveForm.resetFields();
              }}>
                取消
              </Button>
              <Button type="primary" htmlType="submit">
                保存
              </Button>
            </Space>
          </Row>
        </Form>
      </Modal>
    </Card>
  );
};

export default UserFilterPanel;