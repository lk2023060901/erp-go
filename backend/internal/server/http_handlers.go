package server

import (
	"net/http"
	"strconv"
	"time"

	"erp-system/internal/service"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/gorilla/mux"
)

// ========== 用户管理处理器 ==========

// handleListUsers 获取用户列表
func (s *HTTPServer) handleListUsers(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")
	search := r.URL.Query().Get("search")

	page, _ := strconv.ParseInt(pageStr, 10, 32)
	size, _ := strconv.ParseInt(sizeStr, 10, 32)

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	req := &service.UserListRequest{
		Page:   int32(page),
		Size:   int32(size),
		Search: search,
	}

	resp, err := s.userService.ListUsers(r.Context(), req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleCreateUser 创建用户
func (s *HTTPServer) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req service.CreateUserRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}

	resp, err := s.userService.CreateUser(r.Context(), &req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusCreated, resp)
}

// handleGetUser 获取用户详情
func (s *HTTPServer) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	resp, err := s.userService.GetUser(r.Context(), int32(id))
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleUpdateUser 更新用户
func (s *HTTPServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	var req service.UpdateUserRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}
	req.ID = int32(id)

	resp, err := s.userService.UpdateUser(r.Context(), &req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleDeleteUser 删除用户
func (s *HTTPServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	if err := s.userService.DeleteUser(r.Context(), int32(id)); err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, map[string]string{
		"message": "用户删除成功",
	})
}

// handleAssignUserRoles 分配用户角色
func (s *HTTPServer) handleAssignUserRoles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	var req service.AssignRolesRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}
	req.UserID = int32(id)

	if err := s.userService.AssignRoles(r.Context(), &req); err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, map[string]string{
		"message": "角色分配成功",
	})
}

// handleResetUserPassword 重置用户密码
func (s *HTTPServer) handleResetUserPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	var req service.ResetPasswordRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}
	req.UserID = int32(id)

	if err := s.userService.ResetPassword(r.Context(), &req); err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, map[string]string{
		"message": "密码重置成功",
	})
}

// handleToggleUser2FA 切换用户2FA
func (s *HTTPServer) handleToggleUser2FA(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	var req service.ToggleTwoFactorRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}
	req.UserID = int32(id)

	if err := s.userService.ToggleTwoFactor(r.Context(), &req); err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, map[string]string{
		"message": "2FA设置更新成功",
	})
}

// ========== 角色管理处理器 ==========

// handleListRoles 获取角色列表
func (s *HTTPServer) handleListRoles(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")
	search := r.URL.Query().Get("search")
	isEnabledStr := r.URL.Query().Get("is_enabled")
	sortField := r.URL.Query().Get("sortField")
	sortOrder := r.URL.Query().Get("sortOrder")

	page, _ := strconv.ParseInt(pageStr, 10, 32)
	size, _ := strconv.ParseInt(sizeStr, 10, 32)

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	var isEnabled *bool
	if isEnabledStr != "" {
		if enabled, err := strconv.ParseBool(isEnabledStr); err == nil {
			isEnabled = &enabled
		}
	}

	req := &service.RoleListRequest{
		Page:      int32(page),
		Size:      int32(size),
		Search:    search,
		IsEnabled: isEnabled,
		SortField: sortField,
		SortOrder: sortOrder,
	}

	resp, err := s.roleService.ListRoles(r.Context(), req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleCreateRole 创建角色
func (s *HTTPServer) handleCreateRole(w http.ResponseWriter, r *http.Request) {
	var req service.CreateRoleRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}

	resp, err := s.roleService.CreateRole(r.Context(), &req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusCreated, resp)
}

// handleGetRole 获取角色详情
func (s *HTTPServer) handleGetRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	resp, err := s.roleService.GetRole(r.Context(), int32(id))
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleUpdateRole 更新角色
func (s *HTTPServer) handleUpdateRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	var req service.UpdateRoleRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}
	req.ID = int32(id)

	resp, err := s.roleService.UpdateRole(r.Context(), &req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleDeleteRole 删除角色
func (s *HTTPServer) handleDeleteRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	if err := s.roleService.DeleteRole(r.Context(), int32(id)); err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, map[string]string{
		"message": "角色删除成功",
	})
}

// handleAssignRolePermissions 分配角色权限
func (s *HTTPServer) handleAssignRolePermissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	var req service.AssignPermissionsRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}
	req.RoleID = int32(id)

	if err := s.roleService.AssignPermissions(r.Context(), &req); err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, map[string]string{
		"message": "权限分配成功",
	})
}

// handleGetEnabledRoles 获取启用的角色列表
func (s *HTTPServer) handleGetEnabledRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := s.roleService.GetEnabledRoles(r.Context())
	if err != nil {
		s.sendError(w, err)
		return
	}

	// 包装响应格式以匹配前端期望
	response := map[string]interface{}{
		"roles": roles,
	}

	s.sendResponse(w, http.StatusOK, response)
}

// ========== 权限管理处理器 ==========

// handleListPermissions 获取权限列表 - Temporarily disabled (using ERP permission system)
func (s *HTTPServer) handleListPermissions(w http.ResponseWriter, r *http.Request) {
	s.sendError(w, errors.BadRequest("NOT_IMPLEMENTED", "Old permission system disabled. Use ERP permission system."))
	return
	/*
		resp, err := s.permissionService.ListPermissions(r.Context())
		if err != nil {
			s.sendError(w, err)
			return
		}

		s.sendResponse(w, http.StatusOK, resp)
	*/
}

// handleCreatePermission 创建权限 - Temporarily disabled (using ERP permission system)
func (s *HTTPServer) handleCreatePermission(w http.ResponseWriter, r *http.Request) {
	s.sendError(w, errors.BadRequest("NOT_IMPLEMENTED", "Old permission system disabled. Use ERP permission system."))
	return

	/*var req service.CreatePermissionRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}

	resp, err := s.permissionService.CreatePermission(r.Context(), &req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusCreated, resp)
	*/
}

// handleGetPermission 获取权限详情 - Temporarily disabled
func (s *HTTPServer) handleGetPermission(w http.ResponseWriter, r *http.Request) {
	s.sendError(w, errors.BadRequest("NOT_IMPLEMENTED", "Old permission system disabled. Use ERP permission system."))
	return
	/*vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	resp, err := s.permissionService.GetPermission(r.Context(), int32(id))
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
	*/
}

// handleUpdatePermission 更新权限 - Temporarily disabled (using ERP permission system)
func (s *HTTPServer) handleUpdatePermission(w http.ResponseWriter, r *http.Request) {
	s.sendError(w, errors.BadRequest("NOT_IMPLEMENTED", "Old permission system disabled. Use ERP permission system."))
	return

	/*vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	var req service.UpdatePermissionRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}
	req.ID = int32(id)

	resp, err := s.permissionService.UpdatePermission(r.Context(), &req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
	*/
}

// handleDeletePermission 删除权限 - Temporarily disabled
func (s *HTTPServer) handleDeletePermission(w http.ResponseWriter, r *http.Request) {
	s.sendError(w, errors.BadRequest("NOT_IMPLEMENTED", "Old permission system disabled. Use ERP permission system."))
	return
	/*vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	if err := s.permissionService.DeletePermission(r.Context(), int32(id)); err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, map[string]string{
		"message": "权限删除成功",
	})
	*/
}

// handleGetPermissionTree 获取权限树 - Temporarily disabled
func (s *HTTPServer) handleGetPermissionTree(w http.ResponseWriter, r *http.Request) {
	s.sendError(w, errors.BadRequest("NOT_IMPLEMENTED", "Old permission system disabled. Use ERP permission system."))
	return
	/*resp, err := s.permissionService.GetPermissionTree(r.Context())
	if err != nil {
		s.sendError(w, err)
		return
	}

	response := map[string]interface{}{
		"tree": resp,
	}
	s.sendResponse(w, http.StatusOK, response)
	*/
}

// handleGetPermissionModules 获取权限模块列表 - Temporarily disabled
func (s *HTTPServer) handleGetPermissionModules(w http.ResponseWriter, r *http.Request) {
	s.sendError(w, errors.BadRequest("NOT_IMPLEMENTED", "Old permission system disabled. Use ERP permission system."))
	return
	/*modules, err := s.permissionService.GetModules(r.Context())
	if err != nil {
		s.sendError(w, err)
		return
	}

	response := map[string]interface{}{
		"modules": modules,
	}
	s.sendResponse(w, http.StatusOK, response)
	*/
}

// handleSyncApiPermissions 同步API权限 - Temporarily disabled
func (s *HTTPServer) handleSyncApiPermissions(w http.ResponseWriter, r *http.Request) {
	s.sendError(w, errors.BadRequest("NOT_IMPLEMENTED", "Old permission system disabled. Use ERP permission system."))
	return
	/*result, err := s.permissionService.SyncApiPermissions(r.Context())
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, result)
	*/
}

// handleGetUserMenus 获取用户菜单 - Temporarily disabled
func (s *HTTPServer) handleGetUserMenus(w http.ResponseWriter, r *http.Request) {
	s.sendError(w, errors.BadRequest("NOT_IMPLEMENTED", "Old permission system disabled. Use ERP permission system."))
	return
	/*resp, err := s.permissionService.GetUserMenus(r.Context())
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
	*/
}

// handleCheckUserPermission 检查用户权限 - Temporarily disabled
func (s *HTTPServer) handleCheckUserPermission(w http.ResponseWriter, r *http.Request) {
	s.sendError(w, errors.BadRequest("NOT_IMPLEMENTED", "Old permission system disabled. Use ERP permission system."))
	return
	/*var req struct {
		PermissionCode string `json:"permission_code"`
	}
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}

	hasPermission, err := s.permissionService.CheckUserPermission(r.Context(), req.PermissionCode)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, map[string]bool{
		"has_permission": hasPermission,
	})
	*/
}

// ========== 组织管理处理器 ==========

// handleCreateOrganization 创建组织
func (s *HTTPServer) handleCreateOrganization(w http.ResponseWriter, r *http.Request) {
	var req service.CreateOrganizationRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}

	resp, err := s.organizationService.CreateOrganization(r.Context(), &req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusCreated, resp)
}

// handleGetOrganization 获取组织详情
func (s *HTTPServer) handleGetOrganization(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	resp, err := s.organizationService.GetOrganization(r.Context(), int32(id))
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleUpdateOrganization 更新组织
func (s *HTTPServer) handleUpdateOrganization(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	var req service.UpdateOrganizationRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}
	req.ID = int32(id)

	resp, err := s.organizationService.UpdateOrganization(r.Context(), &req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleDeleteOrganization 删除组织
func (s *HTTPServer) handleDeleteOrganization(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	if err := s.organizationService.DeleteOrganization(r.Context(), int32(id)); err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, map[string]string{
		"message": "组织删除成功",
	})
}

// handleGetOrganizationTree 获取组织树
func (s *HTTPServer) handleGetOrganizationTree(w http.ResponseWriter, r *http.Request) {
	resp, err := s.organizationService.GetOrganizationTree(r.Context())
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleGetEnabledOrganizations 获取启用的组织列表
func (s *HTTPServer) handleGetEnabledOrganizations(w http.ResponseWriter, r *http.Request) {
	resp, err := s.organizationService.GetEnabledOrganizations(r.Context())
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleAssignOrganizationUsers 分配组织用户
func (s *HTTPServer) handleAssignOrganizationUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	var req service.AssignUsersRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}
	req.OrganizationID = int32(id)

	if err := s.organizationService.AssignUsers(r.Context(), &req); err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, map[string]string{
		"message": "用户分配成功",
	})
}

// ========== 系统管理处理器 ==========

// handleGetOperationLogs 获取操作日志列表
func (s *HTTPServer) handleGetOperationLogs(w http.ResponseWriter, r *http.Request) {
	var req service.OperationLogListRequest

	// 解析查询参数
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.ParseInt(pageStr, 10, 32); err == nil {
			req.Page = int32(page)
		}
	}
	if sizeStr := r.URL.Query().Get("size"); sizeStr != "" {
		if size, err := strconv.ParseInt(sizeStr, 10, 32); err == nil {
			req.Size = int32(size)
		}
	}
	req.Username = r.URL.Query().Get("username")
	req.Action = r.URL.Query().Get("action")
	req.Resource = r.URL.Query().Get("resource")
	req.Status = r.URL.Query().Get("status")

	if startTimeStr := r.URL.Query().Get("start_time"); startTimeStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			req.StartTime = startTime
		}
	}
	if endTimeStr := r.URL.Query().Get("end_time"); endTimeStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			req.EndTime = endTime
		}
	}

	resp, err := s.systemService.GetOperationLogs(r.Context(), &req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleGetOperationLog 获取操作日志详情
func (s *HTTPServer) handleGetOperationLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		s.sendError(w, err)
		return
	}

	resp, err := s.systemService.GetOperationLog(r.Context(), int32(id))
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleGetOperationStatistics 获取操作统计
func (s *HTTPServer) handleGetOperationStatistics(w http.ResponseWriter, r *http.Request) {
	var req service.StatisticsRequest

	if startTimeStr := r.URL.Query().Get("start_time"); startTimeStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			req.StartTime = startTime
		}
	}
	if endTimeStr := r.URL.Query().Get("end_time"); endTimeStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			req.EndTime = endTime
		}
	}

	// 默认查询最近24小时
	if req.StartTime.IsZero() || req.EndTime.IsZero() {
		req.EndTime = time.Now()
		req.StartTime = req.EndTime.Add(-24 * time.Hour)
	}

	resp, err := s.systemService.GetOperationStatistics(r.Context(), &req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleGetTopActiveUsers 获取最活跃用户
func (s *HTTPServer) handleGetTopActiveUsers(w http.ResponseWriter, r *http.Request) {
	var req service.StatisticsRequest
	limit := int32(10)

	if startTimeStr := r.URL.Query().Get("start_time"); startTimeStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			req.StartTime = startTime
		}
	}
	if endTimeStr := r.URL.Query().Get("end_time"); endTimeStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			req.EndTime = endTime
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil && l > 0 {
			limit = int32(l)
		}
	}

	// 默认查询最近24小时
	if req.StartTime.IsZero() || req.EndTime.IsZero() {
		req.EndTime = time.Now()
		req.StartTime = req.EndTime.Add(-24 * time.Hour)
	}

	resp, err := s.systemService.GetTopActiveUsers(r.Context(), &req, limit)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleCleanupLogs 清理操作日志
func (s *HTTPServer) handleCleanupLogs(w http.ResponseWriter, r *http.Request) {
	var req service.CleanupLogsRequest
	if err := s.parseJSON(r, &req); err != nil {
		s.sendError(w, err)
		return
	}

	affected, err := s.systemService.CleanupOperationLogs(r.Context(), &req)
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, map[string]interface{}{
		"message":       "日志清理完成",
		"affected_rows": affected,
	})
}

// handleGetSystemInfo 获取系统信息
func (s *HTTPServer) handleGetSystemInfo(w http.ResponseWriter, r *http.Request) {
	resp, err := s.systemService.GetSystemInfo(r.Context())
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}

// handleGetDashboardData 获取仪表板数据
func (s *HTTPServer) handleGetDashboardData(w http.ResponseWriter, r *http.Request) {
	resp, err := s.systemService.GetDashboardData(r.Context())
	if err != nil {
		s.sendError(w, err)
		return
	}

	s.sendResponse(w, http.StatusOK, resp)
}
