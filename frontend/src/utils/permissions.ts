/**
 * 权限工具函数
 */

import type { 
  User, 
  Role, 
  Permission,
  PermissionCheckOptions, 
  PermissionResult,
  MenuItem,
  JWTClaims
} from '../types/auth';

// 超级管理员角色代码
export const SUPER_ADMIN_ROLE = 'SUPER_ADMIN';
export const ADMIN_ROLE = 'ADMIN';

/**
 * 检查用户是否为超级管理员
 */
export function isSuperAdmin(roles: string[] | Role[]): boolean {
  const roleCodes = Array.isArray(roles) && typeof roles[0] === 'string' 
    ? roles as string[]
    : (roles as Role[]).map(role => role.code);
  
  return roleCodes.includes(SUPER_ADMIN_ROLE);
}

/**
 * 检查用户是否为管理员
 */
export function isAdmin(roles: string[] | Role[]): boolean {
  const roleCodes = Array.isArray(roles) && typeof roles[0] === 'string' 
    ? roles as string[]
    : (roles as Role[]).map(role => role.code);
  
  return roleCodes.some(code => [SUPER_ADMIN_ROLE, ADMIN_ROLE].includes(code));
}

/**
 * 检查用户是否拥有指定角色
 */
export function hasRole(userRoles: string[] | Role[], targetRole: string | string[]): boolean {
  const roleCodes = Array.isArray(userRoles) && typeof userRoles[0] === 'string' 
    ? userRoles as string[]
    : (userRoles as Role[]).map(role => role.code);
  
  const targetRoles = Array.isArray(targetRole) ? targetRole : [targetRole];
  
  return targetRoles.some(role => roleCodes.includes(role));
}

/**
 * 检查用户是否拥有指定权限
 */
export function hasPermission(userPermissions: string[], targetPermission: string | string[]): boolean {
  if (!userPermissions || userPermissions.length === 0) {
    return false;
  }
  
  const targetPermissions = Array.isArray(targetPermission) ? targetPermission : [targetPermission];
  
  return targetPermissions.some(permission => userPermissions.includes(permission));
}

/**
 * 检查用户是否拥有所有指定权限
 */
export function hasAllPermissions(userPermissions: string[], targetPermissions: string[]): boolean {
  if (!userPermissions || userPermissions.length === 0 || !targetPermissions || targetPermissions.length === 0) {
    return false;
  }
  
  return targetPermissions.every(permission => userPermissions.includes(permission));
}

/**
 * 检查用户是否拥有任意一个指定权限
 */
export function hasAnyPermission(userPermissions: string[], targetPermissions: string[]): boolean {
  if (!userPermissions || userPermissions.length === 0 || !targetPermissions || targetPermissions.length === 0) {
    return false;
  }
  
  return targetPermissions.some(permission => userPermissions.includes(permission));
}

/**
 * 综合权限检查
 */
export function checkPermission(
  user: User | null,
  userPermissions: string[],
  userRoles: Role[],
  options: PermissionCheckOptions
): PermissionResult {
  // 用户未认证
  if (!user) {
    return { allowed: false, reason: 'User not authenticated' };
  }

  // 超级管理员跳过检查（默认行为）
  if (options.skipSuperAdmin !== false && isSuperAdmin(userRoles)) {
    return { allowed: true, reason: 'Super admin bypass' };
  }

  // 自定义检查函数
  if (options.customCheck) {
    const allowed = options.customCheck(user, userPermissions, userRoles);
    return { allowed, reason: allowed ? 'Custom check passed' : 'Custom check failed' };
  }

  // 角色检查
  if (options.roles && options.roles.length > 0) {
    const hasRequiredRole = hasRole(userRoles, options.roles);
    if (!hasRequiredRole) {
      return { allowed: false, reason: `Required roles: ${options.roles.join(', ')}` };
    }
  }

  // 权限检查
  if (options.permissions && options.permissions.length > 0) {
    if (options.requireAll) {
      // 需要所有权限
      const hasAllRequired = hasAllPermissions(userPermissions, options.permissions);
      if (!hasAllRequired) {
        return { allowed: false, reason: `Required all permissions: ${options.permissions.join(', ')}` };
      }
    } else {
      // 需要任意一个权限
      const hasAnyRequired = hasAnyPermission(userPermissions, options.permissions);
      if (!hasAnyRequired) {
        return { allowed: false, reason: `Required any permission: ${options.permissions.join(', ')}` };
      }
    }
  }

  return { allowed: true, reason: 'Permission check passed' };
}

/**
 * 权限编码生成器
 */
export class PermissionCodeGenerator {
  /**
   * 生成标准权限编码
   * @param module 模块名
   * @param resource 资源名
   * @param action 动作
   */
  static generate(module: string, resource: string, action: string): string {
    return `${module}.${resource}.${action}`;
  }

  /**
   * 生成系统管理权限编码
   */
  static system(resource: string, action: string): string {
    return this.generate('system', resource, action);
  }

  /**
   * 生成用户管理权限编码
   */
  static user(action: string): string {
    return this.system('user', action);
  }

  /**
   * 生成角色管理权限编码
   */
  static role(action: string): string {
    return this.system('role', action);
  }

  /**
   * 生成权限管理权限编码
   */
  static permission(action: string): string {
    return this.system('permission', action);
  }

  /**
   * 生成组织管理权限编码
   */
  static organization(action: string): string {
    return this.system('organization', action);
  }

  /**
   * 常用权限编码
   */
  static get USER_VIEW() { return this.user('view'); }
  static get USER_CREATE() { return this.user('create'); }
  static get USER_EDIT() { return this.user('edit'); }
  static get USER_DELETE() { return this.user('delete'); }
  static get USER_ASSIGN_ROLE() { return this.user('assign_role'); }

  static get ROLE_VIEW() { return this.role('view'); }
  static get ROLE_CREATE() { return this.role('create'); }
  static get ROLE_EDIT() { return this.role('edit'); }
  static get ROLE_DELETE() { return this.role('delete'); }
  static get ROLE_ASSIGN_PERMISSION() { return this.role('assign_permission'); }

  static get PERMISSION_VIEW() { return this.permission('view'); }
  static get PERMISSION_CREATE() { return this.permission('create'); }
  static get PERMISSION_EDIT() { return this.permission('edit'); }
  static get PERMISSION_DELETE() { return this.permission('delete'); }

  static get ORG_VIEW() { return this.organization('view'); }
  static get ORG_CREATE() { return this.organization('create'); }
  static get ORG_EDIT() { return this.organization('edit'); }
  static get ORG_DELETE() { return this.organization('delete'); }
  static get ORG_MANAGE_MEMBER() { return this.organization('manage_member'); }
}

/**
 * 菜单过滤器 - 根据权限过滤菜单项
 */
export function filterMenuByPermission(
  menus: MenuItem[],
  userPermissions: string[],
  userRoles: Role[]
): MenuItem[] {
  return menus
    .filter(menu => {
      // 如果没有权限要求，则显示
      if (!menu.meta?.permissions && !menu.meta?.roles) {
        return true;
      }

      // 检查权限要求
      if (menu.meta.permissions && menu.meta.permissions.length > 0) {
        if (!hasAnyPermission(userPermissions, menu.meta.permissions)) {
          return false;
        }
      }

      // 检查角色要求
      if (menu.meta.roles && menu.meta.roles.length > 0) {
        if (!hasRole(userRoles, menu.meta.roles)) {
          return false;
        }
      }

      return true;
    })
    .map(menu => ({
      ...menu,
      children: menu.children ? filterMenuByPermission(menu.children, userPermissions, userRoles) : []
    }));
}

/**
 * 构建权限树
 */
export function buildPermissionTree(permissions: Permission[]): Permission[] {
  const permissionMap = new Map<number, Permission & { children: Permission[] }>();
  
  // 创建权限映射
  permissions.forEach(permission => {
    permissionMap.set(permission.id, { ...permission, children: [] });
  });

  const tree: Permission[] = [];

  // 构建树结构
  permissions.forEach(permission => {
    const node = permissionMap.get(permission.id)!;
    
    if (permission.parent_id && permissionMap.has(permission.parent_id)) {
      const parent = permissionMap.get(permission.parent_id)!;
      parent.children.push(node);
    } else {
      tree.push(node);
    }
  });

  return tree;
}

/**
 * 权限分组器
 */
export function groupPermissionsByModule(permissions: Permission[]): Record<string, Permission[]> {
  return permissions.reduce((acc, permission) => {
    if (!acc[permission.module]) {
      acc[permission.module] = [];
    }
    acc[permission.module].push(permission);
    return acc;
  }, {} as Record<string, Permission[]>);
}

/**
 * JWT Token 解析器
 */
export class TokenUtils {
  /**
   * 解析JWT Token（不验证签名）
   */
  static parseJWT(token: string): JWTClaims | null {
    try {
      const parts = token.split('.');
      if (parts.length !== 3) {
        return null;
      }

      const payload = parts[1];
      const decoded = atob(payload);
      return JSON.parse(decoded) as JWTClaims;
    } catch (error) {
      console.warn('Failed to parse JWT token:', error);
      return null;
    }
  }

  /**
   * 检查Token是否过期
   */
  static isTokenExpired(token: string): boolean {
    const claims = this.parseJWT(token);
    if (!claims) {
      return true;
    }

    const now = Math.floor(Date.now() / 1000);
    return claims.exp <= now;
  }

  /**
   * 获取Token剩余时间（秒）
   */
  static getTokenRemainingTime(token: string): number {
    const claims = this.parseJWT(token);
    if (!claims) {
      return 0;
    }

    const now = Math.floor(Date.now() / 1000);
    return Math.max(0, claims.exp - now);
  }

  /**
   * 检查Token是否即将过期（默认5分钟内）
   */
  static isTokenExpiringSoon(token: string, thresholdMinutes: number = 5): boolean {
    const remainingTime = this.getTokenRemainingTime(token);
    return remainingTime <= thresholdMinutes * 60;
  }
}

/**
 * 权限检查装饰器工厂
 */
export function withPermission(_options: PermissionCheckOptions) {
  return function <T extends (...args: any[]) => any>(target: T): T {
    return ((...args: any[]) => {
      // 这里应该从上下文获取当前用户信息
      // 实际实现需要结合具体的状态管理方案
      console.warn('Permission decorator used, please implement permission check');
      return target(...args);
    }) as T;
  };
}

/**
 * 权限验证Hook帮助函数
 */
export class PermissionHelper {
  /**
   * 创建权限检查函数
   */
  static createChecker(
    userPermissions: string[],
    userRoles: Role[],
    user: User | null
  ) {
    return {
      hasPermission: (permission: string | string[], options?: Omit<PermissionCheckOptions, 'permissions'>) =>
        checkPermission(user, userPermissions, userRoles, {
          ...options,
          permissions: Array.isArray(permission) ? permission : [permission]
        }).allowed,

      hasRole: (role: string | string[]) =>
        hasRole(userRoles, role),

      hasAnyPermission: (permissions: string[]) =>
        hasAnyPermission(userPermissions, permissions),

      hasAllPermissions: (permissions: string[]) =>
        hasAllPermissions(userPermissions, permissions),

      checkPermission: (options: PermissionCheckOptions) =>
        checkPermission(user, userPermissions, userRoles, options),

      isSuperAdmin: () => isSuperAdmin(userRoles),
      isAdmin: () => isAdmin(userRoles)
    };
  }
}

/**
 * 权限缓存管理器
 */
export class PermissionCache {
  private static cache = new Map<string, { data: any; timestamp: number; ttl: number }>();

  /**
   * 设置缓存
   */
  static set(key: string, data: any, ttlSeconds: number = 300): void {
    this.cache.set(key, {
      data,
      timestamp: Date.now(),
      ttl: ttlSeconds * 1000
    });
  }

  /**
   * 获取缓存
   */
  static get<T = any>(key: string): T | null {
    const item = this.cache.get(key);
    if (!item) {
      return null;
    }

    if (Date.now() - item.timestamp > item.ttl) {
      this.cache.delete(key);
      return null;
    }

    return item.data as T;
  }

  /**
   * 删除缓存
   */
  static delete(key: string): void {
    this.cache.delete(key);
  }

  /**
   * 清除所有缓存
   */
  static clear(): void {
    this.cache.clear();
  }

  /**
   * 清除过期缓存
   */
  static cleanup(): void {
    const now = Date.now();
    for (const [key, item] of this.cache.entries()) {
      if (now - item.timestamp > item.ttl) {
        this.cache.delete(key);
      }
    }
  }
}