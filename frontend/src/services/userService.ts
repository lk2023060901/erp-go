/**
 * 用户管理服务
 */

import { apiClient } from './authService';
import type { User, Role, UserFormData, Permission } from '../types/auth';

// 用户列表响应类型
export interface UsersListResponse {
  users: User[];
  total: number;
  page: number;
  size: number;
}

// 用户查询参数
export interface UsersQueryParams {
  page?: number;
  size?: number;
  search?: string;
}

// 创建用户请求（前端格式，会在发送前转换为后端格式）
export interface CreateUserRequest {
  username: string;
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  phone?: string;
  gender?: 'M' | 'F' | 'O'; // 前端格式，发送时转换为 MALE/FEMALE/OTHER
  birth_date?: string;
  is_active?: boolean; // 默认为true
  role_ids?: number[]; // 默认为[2]（管理员角色）
}

// 更新用户请求
export interface UpdateUserRequest {
  email?: string;
  first_name?: string;
  last_name?: string;
  phone?: string;
  gender?: 'M' | 'F' | 'O';
  birth_date?: string;
  is_active?: boolean;
  role_ids?: number[];
}

/**
 * 获取用户列表
 */
export async function getUsersList(params?: UsersQueryParams): Promise<UsersListResponse> {
  try {
    const queryParams = {
      page: 1,
      size: 10,
      ...params,
    };
    
    console.log('查询用户列表参数:', queryParams);
    const response = await apiClient.get<UsersListResponse>('/users', queryParams);
    console.log('用户列表响应:', response);
    return response;
  } catch (error) {
    console.error('Get users list failed:', error);
    throw new Error('获取用户列表失败');
  }
}

/**
 * 获取用户详情
 */
export async function getUserById(id: number): Promise<User> {
  try {
    const response = await apiClient.get<{ user: User }>(`/users/${id}`);
    return response.user;
  } catch (error) {
    console.error('Get user failed:', error);
    throw new Error('获取用户详情失败');
  }
}

/**
 * 创建用户
 */
export async function createUser(data: CreateUserRequest): Promise<User> {
  try {
    // 转换数据格式以匹配后端API
    const requestData = {
      username: data.username,
      email: data.email,
      password: data.password,
      first_name: data.first_name,
      last_name: data.last_name,
      phone: data.phone || '',
      gender: data.gender ? (data.gender === 'M' ? 'MALE' : data.gender === 'F' ? 'FEMALE' : 'OTHER') : '',
      is_active: data.is_active !== false, // 默认为true
      role_ids: data.role_ids && data.role_ids.length > 0 ? data.role_ids : [2], // 默认普通用户角色ID为2
    };

    console.log('创建用户请求数据:', requestData);
    const response = await apiClient.post<{ user: User }>('/users', requestData);
    return response.user;
  } catch (error: any) {
    console.error('Create user failed:', error);
    console.error('Error details:', error.details);
    console.error('Error code:', error.code);
    throw new Error(error.message || '创建用户失败');
  }
}

/**
 * 更新用户
 */
export async function updateUser(id: number, data: UpdateUserRequest): Promise<User> {
  try {
    const response = await apiClient.put<{ user: User }>(`/users/${id}`, data);
    return response.user;
  } catch (error) {
    console.error('Update user failed:', error);
    throw new Error('更新用户失败');
  }
}

/**
 * 删除用户
 */
export async function deleteUser(id: number): Promise<void> {
  try {
    await apiClient.delete(`/users/${id}`);
  } catch (error) {
    console.error('Delete user failed:', error);
    throw new Error('删除用户失败');
  }
}

/**
 * 重置用户密码
 */
export async function resetUserPassword(id: number, newPassword: string): Promise<void> {
  try {
    await apiClient.post(`/users/${id}/reset-password`, {
      new_password: newPassword,
    });
  } catch (error) {
    console.error('Reset user password failed:', error);
    throw new Error('重置用户密码失败');
  }
}

/**
 * 分配用户角色
 */
export async function assignUserRoles(id: number, roleIds: number[]): Promise<void> {
  try {
    await apiClient.post(`/users/${id}/roles`, {
      role_ids: roleIds,
    });
  } catch (error) {
    console.error('Assign user roles failed:', error);
    throw new Error('分配用户角色失败');
  }
}

/**
 * 获取用户角色
 */
export async function getUserRoles(id: number): Promise<Role[]> {
  try {
    const response = await apiClient.get<{ roles: Role[] }>(`/users/${id}/roles`);
    return response.roles;
  } catch (error) {
    console.error('Get user roles failed:', error);
    throw new Error('获取用户角色失败');
  }
}

/**
 * 切换用户2FA状态
 */
export async function toggleUser2FA(id: number, enabled: boolean): Promise<void> {
  try {
    await apiClient.post(`/users/${id}/toggle-2fa`, {
      enable: enabled,
    });
  } catch (error) {
    console.error('Toggle user 2FA failed:', error);
    throw new Error('切换用户双重认证失败');
  }
}

/**
 * 批量删除用户
 */
export async function batchDeleteUsers(ids: number[]): Promise<void> {
  try {
    await apiClient.post('/users/batch-delete', {
      user_ids: ids,
    });
  } catch (error) {
    console.error('Batch delete users failed:', error);
    throw new Error('批量删除用户失败');
  }
}

/**
 * 导出用户数据
 */
export async function exportUsers(params?: UsersQueryParams): Promise<Blob> {
  try {
    const queryParams = new URLSearchParams(params as any).toString();
    const token = localStorage.getItem('access_token');
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';
    
    const response = await fetch(`${API_BASE_URL}/users/export?${queryParams}`, {
      method: 'GET',
      headers: {
        ...(token && { 'Authorization': `Bearer ${token}` }),
      },
    });
    
    if (!response.ok) {
      throw new Error('导出失败');
    }
    
    return response.blob();
  } catch (error) {
    console.error('Export users failed:', error);
    throw new Error('导出用户数据失败');
  }
}

/**
 * 创建用户（从UserFormData）
 */
export async function createUserFromForm(data: UserFormData): Promise<User> {
  try {
    // 转换数据格式以匹配后端API
    const requestData = {
      username: data.username,
      email: data.email,
      password: 'temp123456', // 临时密码，实际应该从表单获取
      first_name: data.first_name,
      last_name: data.last_name,
      phone: data.phone || '',
      gender: data.gender ? (data.gender === 'M' ? 'MALE' : data.gender === 'F' ? 'FEMALE' : 'OTHER') : '',
      birth_date: data.birth_date || '',
      is_active: data.is_active !== false,
      role_ids: data.selected_permissions && data.selected_permissions.length > 0 
        ? [2] // 默认普通用户角色ID为2
        : [2],
    };

    console.log('创建用户请求数据:', requestData);
    const response = await apiClient.post<{ user: User }>('/users', requestData);
    return response.user;
  } catch (error: any) {
    console.error('Create user from form failed:', error);
    throw new Error(error.message || '创建用户失败');
  }
}

/**
 * 更新用户（从UserFormData）
 */
export async function updateUserFromForm(id: number, data: UserFormData): Promise<User> {
  try {
    // 转换数据格式，排除不可修改字段
    const requestData = {
      first_name: data.first_name,
      last_name: data.last_name,
      phone: data.phone || '',
      gender: data.gender ? (data.gender === 'M' ? 'MALE' : data.gender === 'F' ? 'FEMALE' : 'OTHER') : '',
      birth_date: data.birth_date || '',
      is_active: data.is_active,
      address: data.address || '',
      emergency_contact: data.emergency_contact || '',
      emergency_phone: data.emergency_phone || '',
      biography: data.biography || '',
    };

    console.log('更新用户请求数据:', requestData);
    const response = await apiClient.put<{ user: User }>(`/users/${id}`, requestData);
    return response.user;
  } catch (error: any) {
    console.error('Update user from form failed:', error);
    throw new Error(error.message || '更新用户失败');
  }
}

/**
 * 上传用户头像
 */
export async function uploadUserAvatar(userId: number, file: File): Promise<string> {
  try {
    const formData = new FormData();
    formData.append('avatar', file);

    const token = localStorage.getItem('access_token');
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';
    
    const response = await fetch(`${API_BASE_URL}/users/${userId}/avatar`, {
      method: 'POST',
      headers: {
        ...(token && { 'Authorization': `Bearer ${token}` }),
      },
      body: formData,
    });
    
    if (!response.ok) {
      throw new Error('头像上传失败');
    }
    
    const result = await response.json();
    return result.avatar_url;
  } catch (error) {
    console.error('Upload avatar failed:', error);
    throw new Error('头像上传失败');
  }
}

/**
 * 获取权限列表
 */
export async function getPermissionsList(): Promise<Permission[]> {
  try {
    const response = await apiClient.get<{ permissions: Permission[] }>('/permissions');
    return response.permissions;
  } catch (error) {
    console.error('Get permissions list failed:', error);
    throw new Error('获取权限列表失败');
  }
}

/**
 * 获取用户权限
 */
export async function getUserPermissions(id: number): Promise<Permission[]> {
  try {
    const response = await apiClient.get<{ permissions: Permission[] }>(`/users/${id}/permissions`);
    return response.permissions;
  } catch (error) {
    console.error('Get user permissions failed:', error);
    throw new Error('获取用户权限失败');
  }
}

/**
 * 分配用户权限
 */
export async function assignUserPermissions(id: number, permissionCodes: string[]): Promise<void> {
  try {
    await apiClient.post(`/users/${id}/permissions`, {
      permission_codes: permissionCodes,
    });
  } catch (error) {
    console.error('Assign user permissions failed:', error);
    throw new Error('分配用户权限失败');
  }
}