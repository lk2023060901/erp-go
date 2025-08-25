/**
 * DocType管理服务
 */

import { apiClient } from './authService';
import type { DocType, PageRequest, PageResponse } from '../types/auth';

// DocType列表响应类型
export interface DocTypeListResponse extends PageResponse<DocType> {}

// DocType查询参数
export interface DocTypeQueryParams extends PageRequest {
  module?: string;
  is_submittable?: boolean;
  is_child_table?: boolean;
  has_workflow?: boolean;
}

// 创建DocType请求
export interface CreateDocTypeRequest {
  name: string;
  label: string;
  module: string;
  description?: string;
  is_submittable?: boolean;
  is_child_table?: boolean;
  has_workflow?: boolean;
  track_changes?: boolean;
  applies_to_all_users?: boolean;
  max_attachments?: number;
  naming_rule?: string;
  title_field?: string;
  search_fields?: string[];
  sort_field?: string;
  sort_order?: 'ASC' | 'DESC';
}

// 更新DocType请求
export interface UpdateDocTypeRequest {
  label?: string;
  module?: string;
  description?: string;
  is_submittable?: boolean;
  is_child_table?: boolean;
  has_workflow?: boolean;
  track_changes?: boolean;
  applies_to_all_users?: boolean;
  max_attachments?: number;
  naming_rule?: string;
  title_field?: string;
  search_fields?: string[];
  sort_field?: string;
  sort_order?: 'ASC' | 'DESC';
}

// DocType字段类型
export interface DocTypeField {
  id: number;
  doctype_id: number;
  field_name: string;
  field_type: 'Data' | 'Text' | 'Int' | 'Float' | 'Currency' | 'Date' | 'Datetime' | 'Select' | 'Link' | 'Check' | 'Text Editor' | 'Code' | 'Password' | 'Read Only' | 'Attach' | 'Attach Image' | 'Table' | 'Dynamic Link' | 'Small Text';
  field_label: string;
  description?: string;
  options?: string;
  is_mandatory: boolean;
  is_unique: boolean;
  default_value?: string;
  permissions_level: number;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

// DocType模块选项
export const DOCTYPE_MODULES = [
  'Core',
  'Setup',
  'Custom',
  'Accounts',
  'CRM',
  'Selling',
  'Buying',
  'Stock',
  'Manufacturing',
  'Projects',
  'Support',
  'Website',
  'Portal',
  'Utilities',
  'Integrations'
] as const;

// 命名规则选项
export const NAMING_RULES = [
  { label: '自动序号', value: 'autoname' },
  { label: '基于字段', value: 'field' },
  { label: '表达式', value: 'expression' },
  { label: '提示输入', value: 'prompt' },
  { label: '系列编号', value: 'naming_series' }
] as const;

/**
 * 获取DocType列表
 */
export async function getDocTypeList(params?: DocTypeQueryParams): Promise<DocTypeListResponse> {
  try {
    const queryParams = {
      page: 1,
      page_size: 20,
      ...params,
    };
    
    const response = await apiClient.get<DocTypeListResponse>('/doctypes', queryParams);
    return response;
  } catch (error) {
    console.error('Get doctype list failed:', error);
    throw new Error('获取DocType列表失败');
  }
}

/**
 * 获取所有模块的DocType选项
 */
export async function getDocTypeOptions(module?: string): Promise<{ value: string; label: string }[]> {
  try {
    const params = module ? { module } : undefined;
    const response = await apiClient.get<{ doctypes: Array<{ name: string; label: string }> }>('/doctypes/options', params);
    return response.doctypes.map(dt => ({ value: dt.name, label: dt.label }));
  } catch (error) {
    console.error('Get doctype options failed:', error);
    throw new Error('获取DocType选项失败');
  }
}

/**
 * 获取DocType详情
 */
export async function getDocTypeById(id: number): Promise<DocType> {
  try {
    const response = await apiClient.get<{ doctype: DocType }>(`/doctypes/${id}`);
    return response.doctype;
  } catch (error) {
    console.error('Get doctype failed:', error);
    throw new Error('获取DocType详情失败');
  }
}

/**
 * 根据名称获取DocType
 */
export async function getDocTypeByName(name: string): Promise<DocType> {
  try {
    const response = await apiClient.get<{ doctype: DocType }>(`/doctypes/name/${name}`);
    return response.doctype;
  } catch (error) {
    console.error('Get doctype by name failed:', error);
    throw new Error('获取DocType失败');
  }
}

/**
 * 创建DocType
 */
export async function createDocType(data: CreateDocTypeRequest): Promise<DocType> {
  try {
    const response = await apiClient.post<{ doctype: DocType }>('/doctypes', data);
    return response.doctype;
  } catch (error) {
    console.error('Create doctype failed:', error);
    throw new Error('创建DocType失败');
  }
}

/**
 * 更新DocType
 */
export async function updateDocType(id: number, data: UpdateDocTypeRequest): Promise<DocType> {
  try {
    const response = await apiClient.put<{ doctype: DocType }>(`/doctypes/${id}`, data);
    return response.doctype;
  } catch (error) {
    console.error('Update doctype failed:', error);
    throw new Error('更新DocType失败');
  }
}

/**
 * 删除DocType
 */
export async function deleteDocType(id: number): Promise<void> {
  try {
    await apiClient.delete(`/doctypes/${id}`);
  } catch (error) {
    console.error('Delete doctype failed:', error);
    throw new Error('删除DocType失败');
  }
}

/**
 * 批量删除DocType
 */
export async function batchDeleteDocTypes(ids: number[]): Promise<void> {
  try {
    await apiClient.post('/doctypes/batch-delete', {
      doctype_ids: ids,
    });
  } catch (error) {
    console.error('Batch delete doctypes failed:', error);
    throw new Error('批量删除DocType失败');
  }
}

/**
 * 复制DocType
 */
export async function copyDocType(id: number, newName: string, newLabel: string): Promise<DocType> {
  try {
    const response = await apiClient.post<{ doctype: DocType }>(`/doctypes/${id}/copy`, {
      name: newName,
      label: newLabel,
    });
    return response.doctype;
  } catch (error) {
    console.error('Copy doctype failed:', error);
    throw new Error('复制DocType失败');
  }
}

/**
 * 获取DocType字段列表
 */
export async function getDocTypeFields(doctypeId: number): Promise<DocTypeField[]> {
  try {
    const response = await apiClient.get<{ fields: DocTypeField[] }>(`/doctypes/${doctypeId}/fields`);
    return response.fields;
  } catch (error) {
    console.error('Get doctype fields failed:', error);
    throw new Error('获取DocType字段失败');
  }
}

/**
 * 更新DocType字段
 */
export async function updateDocTypeFields(doctypeId: number, fields: Partial<DocTypeField>[]): Promise<void> {
  try {
    await apiClient.put(`/doctypes/${doctypeId}/fields`, {
      fields,
    });
  } catch (error) {
    console.error('Update doctype fields failed:', error);
    throw new Error('更新DocType字段失败');
  }
}

/**
 * 获取可用模块列表
 */
export async function getDocTypeModules(): Promise<string[]> {
  try {
    const response = await apiClient.get<{ modules: string[] }>('/doctypes/modules');
    return response.modules;
  } catch (error) {
    console.error('Get doctype modules failed:', error);
    return [...DOCTYPE_MODULES];
  }
}

/**
 * 验证DocType名称是否可用
 */
export async function validateDocTypeName(name: string): Promise<{ available: boolean; message?: string }> {
  try {
    const response = await apiClient.get<{ available: boolean; message?: string }>(`/doctypes/validate-name/${name}`);
    return response;
  } catch (error) {
    console.error('Validate doctype name failed:', error);
    throw new Error('验证DocType名称失败');
  }
}