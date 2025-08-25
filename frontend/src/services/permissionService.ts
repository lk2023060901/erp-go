/**
 * 权限管理服务
 */

import { apiClient } from './authService';
import type { Permission } from '../types/auth';

// 权限列表响应类型
export interface PermissionsListResponse {
  permissions: Permission[];
  total: number;
  page: number;
  size: number;
}

// 权限查询参数
export interface PermissionsQueryParams {
  page?: number;
  size?: number;
  search?: string;
  module?: string;
  is_enabled?: boolean;
  parent_id?: number;
}

// 创建权限请求
export interface CreatePermissionRequest {
  parent_id?: number;
  name: string;
  code: string;
  resource: string;
  action: string;
  module: string;
  description?: string;
  is_menu?: boolean;
  is_button?: boolean;
  is_api?: boolean;
  menu_url?: string;
  menu_icon?: string;
  api_path?: string;
  api_method?: string;
  level?: number;
  sort_order?: number;
  is_enabled?: boolean;
}

// 更新权限请求
export interface UpdatePermissionRequest {
  name?: string;
  description?: string;
  is_menu?: boolean;
  is_button?: boolean;
  is_api?: boolean;
  menu_url?: string;
  menu_icon?: string;
  api_path?: string;
  api_method?: string;
  sort_order?: number;
  is_enabled?: boolean;
}

// 权限树节点
export interface PermissionTreeNode {
  id: number;
  key: string;
  title: string;
  name: string;
  code: string;
  module: string;
  level: number;
  is_menu: boolean;
  is_button: boolean;
  is_api: boolean;
  is_enabled: boolean;
  children?: PermissionTreeNode[];
}

/**
 * 获取权限列表
 * @deprecated 传统权限系统已禁用，请使用ERP权限系统
 */
export async function getPermissionsList(params?: PermissionsQueryParams): Promise<PermissionsListResponse> {
  console.warn('getPermissionsList is deprecated. Use ERP permission system instead.');
  throw new Error('传统权限系统已禁用，请使用ERP权限系统');
}

/**
 * 获取权限树
 * @deprecated 传统权限系统已禁用，请使用ERP权限系统
 */
export async function getPermissionTree(): Promise<PermissionTreeNode[]> {
  console.warn('getPermissionTree is deprecated. Use ERP permission system instead.');
  throw new Error('传统权限系统已禁用，请使用ERP权限系统');
}

/**
 * 获取权限详情
 */
export async function getPermissionById(id: number): Promise<Permission> {
  try {
    const response = await apiClient.get<{ permission: Permission }>(`/permissions/${id}`);
    return response.permission;
  } catch (error) {
    console.error('Get permission failed:', error);
    throw new Error('获取权限详情失败');
  }
}

/**
 * 创建权限
 */
export async function createPermission(data: CreatePermissionRequest): Promise<Permission> {
  try {
    const response = await apiClient.post<{ permission: Permission }>('/permissions', data);
    return response.permission;
  } catch (error) {
    console.error('Create permission failed:', error);
    throw new Error('创建权限失败');
  }
}

/**
 * 更新权限
 */
export async function updatePermission(id: number, data: UpdatePermissionRequest): Promise<Permission> {
  try {
    const response = await apiClient.put<{ permission: Permission }>(`/permissions/${id}`, data);
    return response.permission;
  } catch (error) {
    console.error('Update permission failed:', error);
    throw new Error('更新权限失败');
  }
}

/**
 * 删除权限
 */
export async function deletePermission(id: number): Promise<void> {
  try {
    await apiClient.delete(`/permissions/${id}`);
  } catch (error) {
    console.error('Delete permission failed:', error);
    throw new Error('删除权限失败');
  }
}

/**
 * 获取模块列表
 */
export async function getModules(): Promise<string[]> {
  try {
    const response = await apiClient.get<{ modules: string[] }>('/permissions/modules');
    return response.modules;
  } catch (error) {
    console.error('Get modules failed:', error);
    throw new Error('获取模块列表失败');
  }
}

/**
 * 检查权限代码是否可用
 */
export async function checkPermissionCode(code: string, excludeId?: number): Promise<boolean> {
  try {
    const response = await apiClient.post<{ available: boolean }>('/permissions/check-code', {
      code,
      exclude_id: excludeId,
    });
    return response.available;
  } catch (error) {
    console.error('Check permission code failed:', error);
    throw new Error('检查权限代码失败');
  }
}

/**
 * 批量删除权限
 */
export async function batchDeletePermissions(ids: number[]): Promise<void> {
  try {
    await apiClient.post('/permissions/batch-delete', {
      permission_ids: ids,
    });
  } catch (error) {
    console.error('Batch delete permissions failed:', error);
    throw new Error('批量删除权限失败');
  }
}

/**
 * 同步API权限
 */
export async function syncApiPermissions(): Promise<{ created: number; updated: number }> {
  try {
    const response = await apiClient.post<{ created: number; updated: number }>('/permissions/sync-api');
    return response;
  } catch (error) {
    console.error('Sync API permissions failed:', error);
    throw new Error('同步API权限失败');
  }
}

/**
 * 生成权限代码
 */
export function generatePermissionCode(resource: string, action: string): string {
  return `${resource.toLowerCase()}.${action.toLowerCase()}`;
}

/**
 * 格式化权限树为Ant Design Tree组件数据
 * @deprecated 传统权限系统已禁用，请使用ERP权限系统
 */
export function formatPermissionTreeForAntd(tree: PermissionTreeNode[], checkedKeys: number[] = []): any[] {
  console.warn('formatPermissionTreeForAntd is deprecated. Use ERP permission system instead.');
  return [];
}