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
  JWTClaims,
  EnhancedUserPermission,
  FieldPermissionResponse,
  DocType,
  PermissionRule
} from '../types/auth';

// 超级管理员角色代码
export const SUPER_ADMIN_ROLE = 'SUPER_ADMIN';
export const ADMIN_ROLE = 'ADMIN';

/**
 * 检查用户是否为超级管理员
 */
export function isSuperAdmin(roles: string[] | Role[]): boolean {
  const roleNames = Array.isArray(roles) && typeof roles[0] === 'string' 
    ? roles as string[]
    : (roles as Role[]).map(role => role.name);
  
  return roleNames.includes(SUPER_ADMIN_ROLE);
}

/**
 * 检查用户是否为管理员
 */
export function isAdmin(roles: string[] | Role[]): boolean {
  const roleNames = Array.isArray(roles) && typeof roles[0] === 'string' 
    ? roles as string[]
    : (roles as Role[]).map(role => role.name);
  
  return roleNames.some(name => [SUPER_ADMIN_ROLE, ADMIN_ROLE].includes(name));
}

/**
 * 检查用户是否拥有指定角色
 */
export function hasRole(userRoles: string[] | Role[], targetRole: string | string[]): boolean {
  const roleNames = Array.isArray(userRoles) && typeof userRoles[0] === 'string' 
    ? userRoles as string[]
    : (userRoles as Role[]).map(role => role.name);
  
  const targetRoles = Array.isArray(targetRole) ? targetRole : [targetRole];
  
  return targetRoles.some(role => roleNames.includes(role));
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

// ========== ERP权限检查函数 ==========

/**
 * 检查用户在特定DocType上是否有特定权限
 */
export function hasErpPermission(
  userPermissions: EnhancedUserPermission[], 
  doc_type: string, 
  permission: keyof EnhancedUserPermission['effective_permissions']
): boolean {
  const userPerm = userPermissions.find(p => p.doc_type === doc_type);
  return userPerm?.effective_permissions[permission] || false;
}

/**
 * 检查用户是否有多个权限中的任意一个
 */
export function hasAnyErpPermission(
  userPermissions: EnhancedUserPermission[], 
  doc_type: string, 
  permissions: Array<keyof EnhancedUserPermission['effective_permissions']>
): boolean {
  return permissions.some(permission => 
    hasErpPermission(userPermissions, doc_type, permission)
  );
}

/**
 * 检查用户是否有所有指定权限
 */
export function hasAllErpPermissions(
  userPermissions: EnhancedUserPermission[], 
  doc_type: string, 
  permissions: Array<keyof EnhancedUserPermission['effective_permissions']>
): boolean {
  return permissions.every(permission => 
    hasErpPermission(userPermissions, doc_type, permission)
  );
}

/**
 * 获取用户在DocType上的所有权限
 */
export function getUserDocTypePermissions(
  userPermissions: EnhancedUserPermission[], 
  doc_type: string
): Partial<EnhancedUserPermission['effective_permissions']> {
  const userPerm = userPermissions.find(p => p.doc_type === doc_type);
  return userPerm?.effective_permissions || {};
}

/**
 * 检查用户是否可以访问特定字段
 */
export function canAccessField(
  fieldPermissions: FieldPermissionResponse,
  field_name: string,
  permission_type: 'read' | 'write'
): boolean {
  const field = fieldPermissions.accessible_fields.find(f => f.field_name === field_name);
  return permission_type === 'read' ? field?.can_read || false : field?.can_write || false;
}

/**
 * 获取用户可读取的字段列表
 */
export function getReadableFields(fieldPermissions: FieldPermissionResponse): string[] {
  return fieldPermissions.accessible_fields
    .filter(field => field.can_read)
    .map(field => field.field_name);
}

/**
 * 获取用户可编辑的字段列表
 */
export function getWritableFields(fieldPermissions: FieldPermissionResponse): string[] {
  return fieldPermissions.accessible_fields
    .filter(field => field.can_write)
    .map(field => field.field_name);
}

/**
 * 检查用户是否可以执行特定操作
 */
export function canPerformAction(
  userPermissions: EnhancedUserPermission[], 
  doc_type: string,
  action: 'create' | 'read' | 'update' | 'delete' | 'submit' | 'cancel' | 'amend'
): boolean {
  const actionMap: Record<string, keyof EnhancedUserPermission['effective_permissions']> = {
    create: 'create',
    read: 'read', 
    update: 'write',
    delete: 'delete',
    submit: 'submit',
    cancel: 'cancel',
    amend: 'amend'
  };

  const permission = actionMap[action];
  return permission ? hasErpPermission(userPermissions, doc_type, permission) : false;
}

/**
 * 获取所有ERP权限类型
 */
export function getErpPermissionTypes(): Array<{ 
  key: keyof EnhancedUserPermission['effective_permissions']; 
  label: string; 
  color?: string;
  description?: string;
}> {
  return [
    { key: 'read', label: '查看', color: 'blue', description: '查看文档内容' },
    { key: 'write', label: '修改', color: 'green', description: '修改文档内容' },
    { key: 'create', label: '创建', color: 'cyan', description: '创建新文档' },
    { key: 'delete', label: '删除', color: 'red', description: '删除文档' },
    { key: 'submit', label: '提交', color: 'purple', description: '提交文档进入工作流' },
    { key: 'cancel', label: '取消', color: 'orange', description: '取消已提交的文档' },
    { key: 'amend', label: '修正', color: 'magenta', description: '修正取消的文档' },
    { key: 'print', label: '打印', color: 'geekblue', description: '打印文档' },
    { key: 'email', label: '邮件', color: 'volcano', description: '通过邮件发送文档' },
    { key: 'import', label: '导入', color: 'gold', description: '导入文档数据' },
    { key: 'export', label: '导出', color: 'lime', description: '导出文档数据' },
    { key: 'share', label: '分享', color: 'cyan', description: '分享文档给其他用户' },
    { key: 'report', label: '报表', color: 'purple', description: '查看相关报表' }
  ];
}

/**
 * ERP权限助手类
 */
export class ErpPermissionHelper {
  /**
   * 创建权限检查器
   */
  static createChecker(userPermissions: EnhancedUserPermission[]) {
    return {
      // 检查单个权限
      hasPermission: (doc_type: string, permission: keyof EnhancedUserPermission['effective_permissions']) =>
        hasErpPermission(userPermissions, doc_type, permission),

      // 检查任意权限
      hasAnyPermission: (doc_type: string, permissions: Array<keyof EnhancedUserPermission['effective_permissions']>) =>
        hasAnyErpPermission(userPermissions, doc_type, permissions),

      // 检查所有权限
      hasAllPermissions: (doc_type: string, permissions: Array<keyof EnhancedUserPermission['effective_permissions']>) =>
        hasAllErpPermissions(userPermissions, doc_type, permissions),

      // 获取DocType权限
      getDocTypePermissions: (doc_type: string) =>
        getUserDocTypePermissions(userPermissions, doc_type),

      // 检查是否可以执行操作
      canPerformAction: (doc_type: string, action: 'create' | 'read' | 'update' | 'delete' | 'submit' | 'cancel' | 'amend') =>
        canPerformAction(userPermissions, doc_type, action),

      // 获取用户有权限的DocType列表
      getAccessibleDocTypes: () =>
        [...new Set(userPermissions.map(p => p.doc_type))],

      // 检查是否有任何读权限
      canRead: (doc_type: string) => hasErpPermission(userPermissions, doc_type, 'read'),
      canWrite: (doc_type: string) => hasErpPermission(userPermissions, doc_type, 'write'),
      canCreate: (doc_type: string) => hasErpPermission(userPermissions, doc_type, 'create'),
      canDelete: (doc_type: string) => hasErpPermission(userPermissions, doc_type, 'delete')
    };
  }

  /**
   * 生成权限矩阵（用于权限管理界面）
   */
  static generatePermissionMatrix(
    docTypes: DocType[],
    roles: Role[],
    permissionRules: PermissionRule[]
  ): Array<{
    docType: DocType;
    permissions: Array<{
      role: Role;
      rules: Partial<EnhancedUserPermission['effective_permissions']>;
    }>;
  }> {
    return docTypes.map(docType => ({
      docType,
      permissions: roles.map(role => {
        const rules = permissionRules.filter(
          rule => rule.document_type === docType.name && rule.role_id === role.id
        );

        // 合并同一角色的多个权限规则
        const mergedRules: Partial<EnhancedUserPermission['effective_permissions']> = {};
        rules.forEach(rule => {
          mergedRules.read = mergedRules.read || rule.read;
          mergedRules.write = mergedRules.write || rule.write;
          mergedRules.create = mergedRules.create || rule.create;
          mergedRules.delete = mergedRules.delete || rule.delete;
          mergedRules.submit = mergedRules.submit || rule.submit;
          mergedRules.cancel = mergedRules.cancel || rule.cancel;
          mergedRules.amend = mergedRules.amend || rule.amend;
          mergedRules.print = mergedRules.print || rule.print;
          mergedRules.email = mergedRules.email || rule.email;
          mergedRules.import = mergedRules.import || rule.import;
          mergedRules.export = mergedRules.export || rule.export;
          mergedRules.share = mergedRules.share || rule.share;
          mergedRules.report = mergedRules.report || rule.report;
        });

        return {
          role,
          rules: mergedRules
        };
      })
    }));
  }
}