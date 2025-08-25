/**
 * 角色管理服务
 */

import { apiClient } from './authService';
import type { Role } from '../types/auth';

// 角色列表响应类型
export interface RolesListResponse {
  roles: Role[];
  total: number;
  page: number;
  size: number;
}

// 角色查询参数
export interface RolesQueryParams {
  page?: number;
  size?: number;
  search?: string;
  is_enabled?: boolean;
  sortField?: string;
  sortOrder?: 'asc' | 'desc';
}

// 创建角色请求
export interface CreateRoleRequest {
  name: string;
  description?: string;
  is_enabled?: boolean;
  sort_order?: number;
  
  // 新增字段
  home_page?: string;
  default_route?: string;
  restrict_to_domain?: string;
  is_custom?: boolean;
  desk_access?: boolean;
  require_two_factor?: boolean;
}

// 更新角色请求
export interface UpdateRoleRequest {
  name?: string;
  description?: string;
  is_enabled?: boolean;
  sort_order?: number;
  
  // 新增字段
  home_page?: string;
  default_route?: string;
  restrict_to_domain?: string;
  desk_access?: boolean;
  require_two_factor?: boolean;
}

/**
 * 获取角色列表
 */
export async function getRolesList(params?: RolesQueryParams): Promise<RolesListResponse> {
  try {
    const queryParams = {
      page: 1,
      size: 10,
      ...params,
    };
    
    const response = await apiClient.get<RolesListResponse>('/roles', queryParams);
    return response;
  } catch (error) {
    console.error('Get roles list failed:', error);
    throw new Error('获取角色列表失败');
  }
}

/**
 * 获取所有启用的角色
 */
export async function getEnabledRoles(): Promise<Role[]> {
  try {
    console.log('🔍 Making API call to /roles/enabled');
    const response = await apiClient.get<Role[] | { roles: Role[] }>('/roles/enabled');
    console.log('✅ API call successful, response:', response);
    
    // 处理两种可能的响应格式
    let roles: Role[] = [];
    if (Array.isArray(response)) {
      // 直接返回数组格式
      roles = response;
      console.log('🔍 Response is array format, length:', roles.length);
    } else if (response && response.roles && Array.isArray(response.roles)) {
      // 包装在 roles 属性中的格式
      roles = response.roles;
      console.log('🔍 Response is object format with roles property, length:', roles.length);
    } else {
      console.warn('⚠️ Unexpected response format:', response);
    }
    
    console.log('📋 Final roles:', roles);
    return roles;
  } catch (error) {
    console.error('❌ Get enabled roles failed:', error);
    throw error;
  }
}

/**
 * 获取角色详情
 */
export async function getRoleById(id: number): Promise<Role> {
  try {
    console.log('getRoleById called with id:', id);
    const response = await apiClient.get<Role>(`/roles/${id}`);
    console.log('getRoleById raw response:', response);
    
    // 如果响应被包装在额外的对象中，需要解包
    const role = (response as any)?.role || response;
    console.log('getRoleById processed role:', role);
    
    if (!role || !role.name) {
      console.error('Invalid role data structure:', role);
      throw new Error('服务器返回的角色数据格式不正确');
    }
    
    return role;
  } catch (error: any) {
    console.error('Get role failed - Full error:', error);
    console.error('Error response:', error.response);
    console.error('Error status:', error.response?.status);
    console.error('Error data:', error.response?.data);
    throw new Error(`获取角色详情失败: ${error.response?.data?.error?.message || error.message}`);
  }
}

/**
 * 创建角色
 */
export async function createRole(data: CreateRoleRequest): Promise<Role> {
  try {
    const response = await apiClient.post<{ role: Role }>('/roles', data);
    return response.role;
  } catch (error) {
    console.error('Create role failed:', error);
    throw new Error('创建角色失败');
  }
}

/**
 * 更新角色
 */
export async function updateRole(id: number, data: UpdateRoleRequest): Promise<Role> {
  try {
    const response = await apiClient.put<{ role: Role }>(`/roles/${id}`, data);
    return response.role;
  } catch (error) {
    console.error('Update role failed:', error);
    throw new Error('更新角色失败');
  }
}

/**
 * 删除角色
 */
export async function deleteRole(id: number): Promise<void> {
  try {
    await apiClient.delete(`/roles/${id}`);
  } catch (error) {
    console.error('Delete role failed:', error);
    throw new Error('删除角色失败');
  }
}

/**
 * 获取角色权限
 */
export async function getRolePermissions(id: number): Promise<number[]> {
  try {
    const response = await apiClient.get<{ permission_ids: number[] }>(`/roles/${id}/permissions`);
    return response.permission_ids;
  } catch (error) {
    console.error('Get role permissions failed:', error);
    throw new Error('获取角色权限失败');
  }
}

/**
 * 分配角色权限
 */
export async function assignRolePermissions(id: number, permissionIds: number[]): Promise<void> {
  try {
    await apiClient.post(`/roles/${id}/permissions`, {
      permission_ids: permissionIds,
    });
  } catch (error) {
    console.error('Assign role permissions failed:', error);
    throw new Error('分配角色权限失败');
  }
}

/**
 * 批量删除角色
 */
export async function batchDeleteRoles(ids: number[]): Promise<void> {
  try {
    await apiClient.post('/roles/batch-delete', {
      role_ids: ids,
    });
  } catch (error) {
    console.error('Batch delete roles failed:', error);
    throw new Error('批量删除角色失败');
  }
}

/**
 * 复制角色
 */
export async function copyRole(id: number, newName: string): Promise<Role> {
  try {
    const response = await apiClient.post<{ role: Role }>(`/roles/${id}/copy`, {
      name: newName,
    });
    return response.role;
  } catch (error) {
    console.error('Copy role failed:', error);
    throw new Error('复制角色失败');
  }
}