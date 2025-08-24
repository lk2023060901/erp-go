/**
 * 过滤器服务 API
 */

// import { request } from './api'; // TODO: 实现API模块

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

// 保存的过滤器接口
export interface SavedFilter {
  id: number;
  user_id: number;
  module_type: string;
  filter_name: string;
  filter_conditions: FilterConfig;
  sort_config?: SortConfig;
  is_default: boolean;
  is_public: boolean;
  created_at: string;
  updated_at: string;
}

// 创建过滤器请求
export interface CreateFilterRequest {
  module_type: string;
  filter_name: string;
  filter_conditions: FilterConfig;
  sort_config?: SortConfig;
  is_default?: boolean;
  is_public?: boolean;
}

// 更新过滤器请求
export interface UpdateFilterRequest {
  filter_name: string;
  filter_conditions: FilterConfig;
  sort_config?: SortConfig;
  is_default?: boolean;
  is_public?: boolean;
}

// API 响应接口
export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
}

// 创建过滤器
export const createFilter = async (request: CreateFilterRequest): Promise<SavedFilter> => {
  const response = await fetch('/api/v1/filters', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('access_token')}`
    },
    body: JSON.stringify({
      module_type: request.module_type,
      filter_name: request.filter_name,
      filter_conditions: JSON.stringify(request.filter_conditions),
      sort_config: request.sort_config ? JSON.stringify(request.sort_config) : undefined,
      is_default: request.is_default || false,
      is_public: request.is_public || false
    })
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.message || '创建过滤器失败');
  }

  const data: ApiResponse<SavedFilter> = await response.json();
  
  if (!data.success) {
    throw new Error(data.message || '创建过滤器失败');
  }

  return data.data;
};

// 获取过滤器列表
export const getFilters = async (moduleType: string): Promise<SavedFilter[]> => {
  const response = await fetch(`/api/v1/filters?module_type=${encodeURIComponent(moduleType)}`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('access_token')}`
    }
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.message || '获取过滤器列表失败');
  }

  const data: ApiResponse<SavedFilter[]> = await response.json();
  
  if (!data.success) {
    throw new Error(data.message || '获取过滤器列表失败');
  }

  return data.data || [];
};

// 获取过滤器详情
export const getFilter = async (filterId: number): Promise<SavedFilter> => {
  const response = await fetch(`/api/v1/filters/${filterId}`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('access_token')}`
    }
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.message || '获取过滤器详情失败');
  }

  const data: ApiResponse<SavedFilter> = await response.json();
  
  if (!data.success) {
    throw new Error(data.message || '获取过滤器详情失败');
  }

  return data.data;
};

// 更新过滤器
export const updateFilter = async (filterId: number, request: UpdateFilterRequest): Promise<SavedFilter> => {
  const response = await fetch(`/api/v1/filters/${filterId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('access_token')}`
    },
    body: JSON.stringify({
      filter_name: request.filter_name,
      filter_conditions: JSON.stringify(request.filter_conditions),
      sort_config: request.sort_config ? JSON.stringify(request.sort_config) : undefined,
      is_default: request.is_default || false,
      is_public: request.is_public || false
    })
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.message || '更新过滤器失败');
  }

  const data: ApiResponse<SavedFilter> = await response.json();
  
  if (!data.success) {
    throw new Error(data.message || '更新过滤器失败');
  }

  return data.data;
};

// 删除过滤器
export const deleteFilter = async (filterId: number): Promise<void> => {
  const response = await fetch(`/api/v1/filters/${filterId}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('access_token')}`
    }
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.message || '删除过滤器失败');
  }

  const data: ApiResponse<any> = await response.json();
  
  if (!data.success) {
    throw new Error(data.message || '删除过滤器失败');
  }
};

// 设置默认过滤器
export const setDefaultFilter = async (filterId: number, moduleType: string): Promise<void> => {
  const response = await fetch(`/api/v1/filters/${filterId}/set-default`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('access_token')}`
    },
    body: JSON.stringify({
      module_type: moduleType
    })
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.message || '设置默认过滤器失败');
  }

  const data: ApiResponse<any> = await response.json();
  
  if (!data.success) {
    throw new Error(data.message || '设置默认过滤器失败');
  }
};

// 过滤器工具函数

// 验证过滤条件
export const validateFilterConditions = (conditions: FilterCondition[]): boolean => {
  return conditions.every(condition => {
    return condition.field && 
           condition.operator && 
           condition.value !== '' && 
           condition.value !== null && 
           condition.value !== undefined;
  });
};

// 构建过滤器查询参数
export const buildFilterParams = (
  page: number,
  size: number,
  filterConfig?: FilterConfig,
  sortConfig?: SortConfig,
  filterId?: number,
  search?: string
) => {
  const params = new URLSearchParams({
    page: page.toString(),
    size: size.toString()
  });

  // 兼容旧版搜索
  if (search) {
    params.append('search', search);
  }

  // 新的过滤器参数
  if (filterConfig) {
    params.append('filter_conditions', JSON.stringify(filterConfig));
  }

  if (sortConfig) {
    params.append('sort_config', JSON.stringify(sortConfig));
  }

  if (filterId) {
    params.append('filter_id', filterId.toString());
  }

  return params;
};

// 解析过滤器条件显示文本
export const parseFilterConditionText = (condition: FilterCondition): string => {
  const fieldLabels: Record<string, string> = {
    'username': '用户名',
    'email': '邮箱',
    'first_name': '名字',
    'last_name': '姓氏',
    'phone': '手机号',
    'gender': '性别',
    'is_active': '状态',
    'created_at': '创建时间',
    'updated_at': '更新时间',
    'last_login_time': '最后登录'
  };

  const operatorLabels: Record<string, string> = {
    'equals': '等于',
    'contains': '包含',
    'starts_with': '开头是',
    'ends_with': '结尾是',
    'greater_than': '大于',
    'less_than': '小于',
    'between': '介于',
    'in': '属于',
    'not_in': '不属于'
  };

  const fieldLabel = fieldLabels[condition.field] || condition.field;
  const operatorLabel = operatorLabels[condition.operator] || condition.operator;
  
  let valueLabel = condition.value;
  if (condition.field === 'gender') {
    const genderLabels: Record<string, string> = { 'M': '男', 'F': '女', 'O': '其他' };
    valueLabel = genderLabels[condition.value] || condition.value;
  } else if (condition.field === 'is_active') {
    valueLabel = condition.value ? '启用' : '禁用';
  }

  return `${fieldLabel} ${operatorLabel} ${valueLabel}`;
};

export default {
  createFilter,
  getFilters,
  getFilter,
  updateFilter,
  deleteFilter,
  setDefaultFilter,
  validateFilterConditions,
  buildFilterParams,
  parseFilterConditionText
};