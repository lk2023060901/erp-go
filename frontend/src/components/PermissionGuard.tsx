/**
 * 权限守卫组件
 */

import React, { ReactNode } from 'react';
import { usePermission } from '../hooks/useAuth';
import type { PermissionCheckOptions } from '../types/auth';

interface PermissionGuardProps {
  children: ReactNode;
  // 权限检查选项
  permissions?: string[];
  roles?: string[];
  requireAll?: boolean;
  skipSuperAdmin?: boolean;
  customCheck?: PermissionCheckOptions['customCheck'];
  // 渲染选项
  fallback?: ReactNode;
  loading?: ReactNode;
  // 行为选项
  hideIfNoPermission?: boolean;
  disableIfNoPermission?: boolean;
  // 调试选项
  debug?: boolean;
}

/**
 * 权限守卫组件
 * 
 * 用法示例：
 * ```tsx
 * // 基础用法 - 检查权限
 * <PermissionGuard permissions={['user.create']}>
 *   <CreateUserButton />
 * </PermissionGuard>
 * 
 * // 检查角色
 * <PermissionGuard roles={['ADMIN']}>
 *   <AdminPanel />
 * </PermissionGuard>
 * 
 * // 组合检查
 * <PermissionGuard 
 *   permissions={['user.create', 'user.edit']} 
 *   requireAll={false}
 *   fallback={<div>权限不足</div>}
 * >
 *   <UserManagement />
 * </PermissionGuard>
 * 
 * // 禁用而非隐藏
 * <PermissionGuard 
 *   permissions={['user.delete']}
 *   disableIfNoPermission={true}
 * >
 *   <DeleteButton />
 * </PermissionGuard>
 * ```
 */
export function PermissionGuard({
  children,
  permissions,
  roles,
  requireAll = false,
  skipSuperAdmin = true,
  customCheck,
  fallback = null,
  loading = null,
  hideIfNoPermission = true,
  disableIfNoPermission = false,
  debug = false,
}: PermissionGuardProps) {
  const { user, checkPermission, permissions: userPermissions, roles: userRoles } = usePermission();

  // 构建权限检查选项
  const checkOptions: PermissionCheckOptions = {
    permissions,
    roles,
    requireAll,
    skipSuperAdmin,
    customCheck,
  };

  // 执行权限检查
  const permissionResult = checkPermission(checkOptions);

  // 调试模式输出
  if (debug) {
    console.group(`PermissionGuard Debug`);
    console.log('User:', user);
    console.log('User Permissions:', userPermissions);
    console.log('User Roles:', userRoles);
    console.log('Check Options:', checkOptions);
    console.log('Permission Result:', permissionResult);
    console.groupEnd();
  }

  // 用户未加载完成时显示加载状态
  if (!user && loading) {
    return <>{loading}</>;
  }

  // 权限检查未通过
  if (!permissionResult.allowed) {
    // 隐藏组件
    if (hideIfNoPermission) {
      return <>{fallback}</>;
    }
    
    // 禁用组件
    if (disableIfNoPermission && React.isValidElement(children)) {
      return React.cloneElement(children, {
        disabled: true,
        title: `权限不足: ${permissionResult.reason}`,
        ...children.props,
      });
    }

    return <>{fallback}</>;
  }

  return <>{children}</>;
}

/**
 * 权限高阶组件
 */
export function withPermission(checkOptions: Omit<PermissionGuardProps, 'children'>) {
  return function <P extends object>(Component: React.ComponentType<P>) {
    const WrappedComponent = (props: P) => (
      <PermissionGuard {...checkOptions}>
        <Component {...props} />
      </PermissionGuard>
    );

    WrappedComponent.displayName = `withPermission(${Component.displayName || Component.name})`;
    
    return WrappedComponent;
  };
}

/**
 * 权限检查渲染函数组件
 */
interface PermissionRenderProps {
  children: (result: {
    allowed: boolean;
    reason?: string;
    hasPermission: (permission: string | string[]) => boolean;
    hasRole: (role: string | string[]) => boolean;
  }) => ReactNode;
  permissions?: string[];
  roles?: string[];
  requireAll?: boolean;
  skipSuperAdmin?: boolean;
  customCheck?: PermissionCheckOptions['customCheck'];
}

export function PermissionRender({
  children,
  permissions,
  roles,
  requireAll = false,
  skipSuperAdmin = true,
  customCheck,
}: PermissionRenderProps) {
  const { checkPermission, hasPermission, hasRole } = usePermission();

  const checkOptions: PermissionCheckOptions = {
    permissions,
    roles,
    requireAll,
    skipSuperAdmin,
    customCheck,
  };

  const result = checkPermission(checkOptions);

  return (
    <>
      {children({
        allowed: result.allowed,
        reason: result.reason,
        hasPermission,
        hasRole,
      })}
    </>
  );
}

/**
 * 权限提示组件
 */
interface PermissionTooltipProps {
  children: ReactNode;
  permissions?: string[];
  roles?: string[];
  requireAll?: boolean;
  showTooltip?: boolean;
  tooltipContent?: ReactNode;
}

export function PermissionTooltip({
  children,
  permissions,
  roles,
  requireAll = false,
  showTooltip = true,
  tooltipContent,
}: PermissionTooltipProps) {
  const { checkPermission } = usePermission();

  const result = checkPermission({
    permissions,
    roles,
    requireAll,
  });

  // 如果有权限或不显示提示，直接返回子组件
  if (result.allowed || !showTooltip) {
    return <>{children}</>;
  }

  // 这里应该使用实际的Tooltip组件，比如antd的Tooltip
  const tooltip = tooltipContent || `权限不足: ${result.reason}`;

  return (
    <span title={typeof tooltip === 'string' ? tooltip : undefined}>
      {children}
    </span>
  );
}

/**
 * 角色守卫组件（简化版本）
 */
interface RoleGuardProps {
  children: ReactNode;
  roles: string | string[];
  fallback?: ReactNode;
}

export function RoleGuard({ children, roles, fallback = null }: RoleGuardProps) {
  return (
    <PermissionGuard roles={Array.isArray(roles) ? roles : [roles]} fallback={fallback}>
      {children}
    </PermissionGuard>
  );
}

/**
 * 超级管理员守卫组件
 */
interface SuperAdminGuardProps {
  children: ReactNode;
  fallback?: ReactNode;
}

export function SuperAdminGuard({ children, fallback = null }: SuperAdminGuardProps) {
  return (
    <PermissionGuard roles={['SUPER_ADMIN']} fallback={fallback}>
      {children}
    </PermissionGuard>
  );
}

/**
 * 管理员守卫组件
 */
interface AdminGuardProps {
  children: ReactNode;
  fallback?: ReactNode;
}

export function AdminGuard({ children, fallback = null }: AdminGuardProps) {
  return (
    <PermissionGuard roles={['SUPER_ADMIN', 'ADMIN']} fallback={fallback}>
      {children}
    </PermissionGuard>
  );
}

/**
 * 权限按钮组件
 */
interface PermissionButtonProps extends Omit<React.ButtonHTMLAttributes<HTMLButtonElement>, 'disabled'> {
  permissions?: string[];
  roles?: string[];
  requireAll?: boolean;
  hideIfNoPermission?: boolean;
  forceDisabled?: boolean;
  children: ReactNode;
}

export function PermissionButton({
  permissions,
  roles,
  requireAll = false,
  hideIfNoPermission = false,
  forceDisabled = false,
  children,
  ...buttonProps
}: PermissionButtonProps) {
  const { checkPermission } = usePermission();

  const result = checkPermission({
    permissions,
    roles,
    requireAll,
  });

  // 如果没有权限且需要隐藏，返回null
  if (!result.allowed && hideIfNoPermission) {
    return null;
  }

  return (
    <button
      {...buttonProps}
      disabled={forceDisabled || !result.allowed || (buttonProps as any).disabled}
      title={!result.allowed ? `权限不足: ${result.reason}` : buttonProps.title}
    >
      {children}
    </button>
  );
}

/**
 * 条件渲染组件
 */
interface ConditionalRenderProps {
  condition: boolean;
  children: ReactNode;
  fallback?: ReactNode;
}

export function ConditionalRender({ condition, children, fallback = null }: ConditionalRenderProps) {
  return condition ? <>{children}</> : <>{fallback}</>;
}

/**
 * 权限菜单项组件
 */
interface PermissionMenuItemProps {
  children: ReactNode;
  permissions?: string[];
  roles?: string[];
  requireAll?: boolean;
  visible?: boolean;
}

export function PermissionMenuItem({
  children,
  permissions,
  roles,
  requireAll = false,
  visible = true,
}: PermissionMenuItemProps) {
  if (!visible) {
    return null;
  }

  return (
    <PermissionGuard
      permissions={permissions}
      roles={roles}
      requireAll={requireAll}
      hideIfNoPermission={true}
    >
      {children}
    </PermissionGuard>
  );
}