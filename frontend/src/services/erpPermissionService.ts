/**
 * ERP权限管理服务
 */

import { apiClient } from './authService';
import type { 
  DocType, 
  PermissionRule, 
  UserPermission, 
  ErpPermissionCheckOptions,
  PermissionCheckResponse,
  UserPermissionLevelResponse,
  EnhancedUserPermission,
  FieldPermissionResponse
} from '../types/auth';

// ========== DocType管理 ==========

export interface DocTypeListResponse {
  doctypes: DocType[];
  total: number;
  page: number;
  size: number;
}

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

export interface UpdateDocTypeRequest {
  label?: string;
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

// 获取DocType列表
export const getDocTypesList = async (module?: string): Promise<DocTypeListResponse> => {
  try {
    const params = module ? { module } : {};
    const response = await apiClient.get<{ doc_types: DocType[]; total: number }>('/erp-permissions/doctypes', params);
    return {
      doctypes: response.doc_types,
      total: response.total,
      page: 1,
      size: response.doc_types.length
    };
  } catch (error) {
    console.error('获取DocType列表失败:', error);
    throw new Error('获取DocType列表失败');
  }
};

// 获取单个DocType
export const getDocType = async (name: string): Promise<DocType> => {
  try {
    const response = await apiClient.get<{ doctype: DocType }>(`/erp-permissions/doctypes/${name}`);
    return response.doctype;
  } catch (error) {
    console.error('获取DocType详情失败:', error);
    throw new Error('获取DocType详情失败');
  }
};

// 创建DocType
export const createDocType = async (data: CreateDocTypeRequest): Promise<DocType> => {
  try {
    const response = await apiClient.post<{ doctype: DocType }>('/erp-permissions/doctypes', data);
    return response.doctype;
  } catch (error) {
    console.error('创建DocType失败:', error);
    throw new Error('创建DocType失败');
  }
};

// 更新DocType
export const updateDocType = async (name: string, data: UpdateDocTypeRequest): Promise<DocType> => {
  try {
    const response = await apiClient.put<{ doctype: DocType }>(`/erp-permissions/doctypes/${name}`, data);
    return response.doctype;
  } catch (error) {
    console.error('更新DocType失败:', error);
    throw new Error('更新DocType失败');
  }
};

// 删除DocType
export const deleteDocType = async (name: string): Promise<void> => {
  try {
    await apiClient.delete(`/erp-permissions/doctypes/${name}`);
  } catch (error) {
    console.error('删除DocType失败:', error);
    throw new Error('删除DocType失败');
  }
};

// ========== 权限规则管理 ==========

export interface PermissionRuleListResponse {
  rules: PermissionRule[];
  total: number;
}

export interface CreatePermissionRuleRequest {
  role_id: number;
  doc_type: string;
  permission_level: number;
  can_read: boolean;
  can_write: boolean;
  can_create: boolean;
  can_delete: boolean;
  can_submit?: boolean;
  can_cancel?: boolean;
  can_amend?: boolean;
  can_print?: boolean;
  can_email?: boolean;
  can_import?: boolean;
  can_export?: boolean;
  can_share?: boolean;
  can_report?: boolean;
  can_set_user_permissions?: boolean;
  only_if_creator?: boolean;
}

export interface UpdatePermissionRuleRequest extends CreatePermissionRuleRequest {
  id: number;
}

// 获取权限规则列表
export const getPermissionRulesList = async (roleId?: number, docType?: string): Promise<PermissionRuleListResponse> => {
  try {
    const params: any = {};
    if (roleId) params.role_id = roleId;
    if (docType) params.doc_type = docType;
    
    const response = await apiClient.get<{ rules: PermissionRule[]; total: number }>('/erp-permissions/permission-rules', params);
    return response;
  } catch (error) {
    console.error('获取权限规则列表失败:', error);
    throw new Error('获取权限规则列表失败');
  }
};

// 获取单个权限规则
export const getPermissionRule = async (id: number): Promise<PermissionRule> => {
  try {
    const response = await apiClient.get<{ rule: PermissionRule }>(`/erp-permissions/permission-rules/${id}`);
    return response.rule;
  } catch (error) {
    console.error('获取权限规则详情失败:', error);
    throw new Error('获取权限规则详情失败');
  }
};

// 创建权限规则
export const createPermissionRule = async (data: CreatePermissionRuleRequest): Promise<PermissionRule> => {
  try {
    const response = await apiClient.post<{ rule: PermissionRule }>('/erp-permissions/permission-rules', data);
    return response.rule;
  } catch (error) {
    console.error('创建权限规则失败:', error);
    throw new Error('创建权限规则失败');
  }
};

// 更新权限规则
export const updatePermissionRule = async (id: number, data: Omit<UpdatePermissionRuleRequest, 'id'>): Promise<PermissionRule> => {
  try {
    const response = await apiClient.put<{ rule: PermissionRule }>(`/erp-permissions/permission-rules/${id}`, { ...data, id });
    return response.rule;
  } catch (error) {
    console.error('更新权限规则失败:', error);
    throw new Error('更新权限规则失败');
  }
};

// 删除权限规则
export const deletePermissionRule = async (id: number): Promise<void> => {
  try {
    await apiClient.delete(`/erp-permissions/permission-rules/${id}`);
  } catch (error) {
    console.error('删除权限规则失败:', error);
    throw new Error('删除权限规则失败');
  }
};

// 批量创建权限规则
export const batchCreatePermissionRules = async (rules: CreatePermissionRuleRequest[]): Promise<void> => {
  try {
    await apiClient.post('/erp-permissions/permission-rules/batch', { rules });
  } catch (error) {
    console.error('批量创建权限规则失败:', error);
    throw new Error('批量创建权限规则失败');
  }
};

// ========== 权限检查 ==========

// 检查用户权限
export const checkPermission = async (options: ErpPermissionCheckOptions): Promise<PermissionCheckResponse> => {
  try {
    const response = await apiClient.post<PermissionCheckResponse>('/erp-permissions/check-permission', options);
    return response;
  } catch (error) {
    console.error('权限检查失败:', error);
    throw new Error('权限检查失败');
  }
};

// 检查文档权限
export const checkDocumentPermission = async (
  userId: number, 
  docType: string, 
  permission: string, 
  docId?: number
): Promise<PermissionCheckResponse> => {
  try {
    const data = {
      user_id: userId,
      doc_type: docType,
      permission,
      ...(docId && { doc_id: docId })
    };
    const response = await apiClient.post<PermissionCheckResponse>('/erp-permissions/check-document-permission', data);
    return response;
  } catch (error) {
    console.error('文档权限检查失败:', error);
    throw new Error('文档权限检查失败');
  }
};

// 获取用户权限级别
export const getUserPermissionLevel = async (userId: number, docType: string): Promise<UserPermissionLevelResponse> => {
  try {
    const response = await apiClient.post<UserPermissionLevelResponse>('/erp-permissions/get-user-permission-level', {
      user_id: userId,
      doc_type: docType
    });
    return response;
  } catch (error) {
    console.error('获取用户权限级别失败:', error);
    throw new Error('获取用户权限级别失败');
  }
};

// 获取用户增强权限
export const getUserEnhancedPermissions = async (userId: number, docType: string): Promise<EnhancedUserPermission[]> => {
  try {
    const response = await apiClient.post<{ permissions: EnhancedUserPermission[] }>('/erp-permissions/get-user-enhanced-permissions', {
      user_id: userId,
      doc_type: docType
    });
    return response.permissions;
  } catch (error) {
    console.error('获取用户增强权限失败:', error);
    throw new Error('获取用户增强权限失败');
  }
};

// 获取可访问字段
export const getAccessibleFields = async (userId: number, docType: string): Promise<FieldPermissionResponse> => {
  try {
    const response = await apiClient.post<FieldPermissionResponse>('/erp-permissions/get-accessible-fields', {
      user_id: userId,
      doc_type: docType
    });
    return response;
  } catch (error) {
    console.error('获取可访问字段失败:', error);
    throw new Error('获取可访问字段失败');
  }
};

// 按权限过滤文档
export const filterDocumentsByPermission = async (
  userId: number, 
  documentType: string, 
  documents: any[]
): Promise<any[]> => {
  try {
    const response = await apiClient.post<{ filtered_documents: any[] }>('/erp-permissions/filter-documents-by-permission', {
      user_id: userId,
      document_type: documentType,
      documents
    });
    return response.filtered_documents;
  } catch (error) {
    console.error('按权限过滤文档失败:', error);
    throw new Error('按权限过滤文档失败');
  }
};

// 获取用户角色
export const getUserRoles = async (userId: number): Promise<string[]> => {
  try {
    const response = await apiClient.post<{ roles: string[] }>('/erp-permissions/get-user-roles', {
      user_id: userId
    });
    return response.roles;
  } catch (error) {
    console.error('获取用户角色失败:', error);
    throw new Error('获取用户角色失败');
  }
};

// ========== 用户权限管理 ==========

export interface UserPermissionListResponse {
  user_permissions: UserPermission[];
  total: number;
  page: number;
  size: number;
}

export interface CreateUserPermissionRequest {
  user_id: number;
  doc_type: string;
  document_name: string;
  condition?: string;
  applicable_for?: string;
  hide_descendants?: boolean;
  is_default?: boolean;
}

// 获取用户权限列表
export const getUserPermissionsList = async (
  userId?: number, 
  docType?: string,
  page: number = 1,
  size: number = 20
): Promise<UserPermissionListResponse> => {
  try {
    const params: any = { page, size };
    if (userId) params.user_id = userId;
    if (docType) params.doc_type = docType;
    
    const response = await apiClient.get<UserPermissionListResponse>('/erp-permissions/user-permissions', params);
    return response;
  } catch (error) {
    console.error('获取用户权限列表失败:', error);
    throw new Error('获取用户权限列表失败');
  }
};

// 创建用户权限
export const createUserPermission = async (data: CreateUserPermissionRequest): Promise<UserPermission> => {
  try {
    const response = await apiClient.post<{ user_permission: UserPermission }>('/erp-permissions/user-permissions', data);
    return response.user_permission;
  } catch (error) {
    console.error('创建用户权限失败:', error);
    throw new Error('创建用户权限失败');
  }
};

// 获取用户权限详情
export const getUserPermission = async (id: number): Promise<UserPermission> => {
  try {
    const response = await apiClient.get<{ user_permission: UserPermission }>(`/erp-permissions/user-permissions/${id}`);
    return response.user_permission;
  } catch (error) {
    console.error('获取用户权限详情失败:', error);
    throw new Error('获取用户权限详情失败');
  }
};

// 更新用户权限
export const updateUserPermission = async (id: number, data: CreateUserPermissionRequest): Promise<UserPermission> => {
  try {
    const response = await apiClient.put<{ user_permission: UserPermission }>(`/erp-permissions/user-permissions/${id}`, data);
    return response.user_permission;
  } catch (error) {
    console.error('更新用户权限失败:', error);
    throw new Error('更新用户权限失败');
  }
};

// 删除用户权限
export const deleteUserPermission = async (id: number): Promise<void> => {
  try {
    await apiClient.delete(`/erp-permissions/user-permissions/${id}`);
  } catch (error) {
    console.error('删除用户权限失败:', error);
    throw new Error('删除用户权限失败');
  }
};

// ========== 字段权限级别管理 ==========

export interface FieldPermissionLevel {
  id: number;
  doc_type: string;
  field_name: string;
  permission_level: number;
  created_at: string;
  updated_at: string;
}

export interface FieldPermissionLevelListResponse {
  field_permission_levels: FieldPermissionLevel[];
  total: number;
  page: number;
  size: number;
}

export interface CreateFieldPermissionLevelRequest {
  doc_type: string;
  field_name: string;
  permission_level: number;
}

// 获取字段权限级别列表
export const getFieldPermissionLevelsList = async (
  docType?: string,
  page: number = 1,
  size: number = 20
): Promise<FieldPermissionLevelListResponse> => {
  try {
    const params: any = { page, size };
    if (docType) params.doc_type = docType;
    
    const response = await apiClient.get<FieldPermissionLevelListResponse>('/erp-permissions/field-permission-levels', params);
    return response;
  } catch (error) {
    console.error('获取字段权限级别列表失败:', error);
    throw new Error('获取字段权限级别列表失败');
  }
};

// 创建字段权限级别
export const createFieldPermissionLevel = async (data: CreateFieldPermissionLevelRequest): Promise<FieldPermissionLevel> => {
  try {
    const response = await apiClient.post<{ field_permission_level: FieldPermissionLevel }>('/erp-permissions/field-permission-levels', data);
    return response.field_permission_level;
  } catch (error) {
    console.error('创建字段权限级别失败:', error);
    throw new Error('创建字段权限级别失败');
  }
};

// 更新字段权限级别
export const updateFieldPermissionLevel = async (id: number, data: CreateFieldPermissionLevelRequest): Promise<FieldPermissionLevel> => {
  try {
    const response = await apiClient.put<{ field_permission_level: FieldPermissionLevel }>(`/erp-permissions/field-permission-levels/${id}`, data);
    return response.field_permission_level;
  } catch (error) {
    console.error('更新字段权限级别失败:', error);
    throw new Error('更新字段权限级别失败');
  }
};

// 删除字段权限级别
export const deleteFieldPermissionLevel = async (id: number): Promise<void> => {
  try {
    await apiClient.delete(`/erp-permissions/field-permission-levels/${id}`);
  } catch (error) {
    console.error('删除字段权限级别失败:', error);
    throw new Error('删除字段权限级别失败');
  }
};