package service

import (
	"context"
	"time"

	"erp-system/internal/biz"
	"erp-system/internal/middleware"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// SystemService 系统服务
type SystemService struct {
	auditUc *biz.AuditUsecase
	log     *log.Helper
}

// NewSystemService 创建系统服务
func NewSystemService(auditUc *biz.AuditUsecase, logger log.Logger) *SystemService {
	return &SystemService{
		auditUc: auditUc,
		log:     log.NewHelper(logger),
	}
}

// OperationLogListRequest 操作日志列表请求
type OperationLogListRequest struct {
	Page      int32     `json:"page" validate:"min=1"`
	Size      int32     `json:"size" validate:"min=1,max=100"`
	UserID    *int32    `json:"user_id"`
	Username  string    `json:"username"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Status    string    `json:"status"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// OperationLogListResponse 操作日志列表响应
type OperationLogListResponse struct {
	Logs  []*OperationLogInfo `json:"logs"`
	Total int32               `json:"total"`
	Page  int32               `json:"page"`
	Size  int32               `json:"size"`
}

// OperationLogInfo 操作日志信息
type OperationLogInfo struct {
	ID            int32     `json:"id"`
	UserID        *int32    `json:"user_id"`
	Username      string    `json:"username"`
	Action        string    `json:"action"`
	Resource      string    `json:"resource"`
	ResourceID    string    `json:"resource_id"`
	Description   string    `json:"description"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	RequestData   string    `json:"request_data,omitempty"`
	ResponseData  string    `json:"response_data,omitempty"`
	Status        string    `json:"status"`
	ErrorMessage  string    `json:"error_message"`
	ExecutionTime int32     `json:"execution_time"`
	CreatedAt     time.Time `json:"created_at"`
}

// StatisticsRequest 统计请求
type StatisticsRequest struct {
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
}

// OperationStatisticsInfo 操作统计信息
type OperationStatisticsInfo struct {
	TotalOperations      int32            `json:"total_operations"`
	SuccessOperations    int32            `json:"success_operations"`
	FailedOperations     int32            `json:"failed_operations"`
	SuccessRate          float64          `json:"success_rate"`
	AverageExecutionTime float64          `json:"average_execution_time"`
	ActiveUsers          int32            `json:"active_users"`
	ActionDistribution   map[string]int32 `json:"action_distribution"`
}

// UserActivityInfo 用户活动信息
type UserActivityInfo struct {
	UserID         int32  `json:"user_id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	OperationCount int32  `json:"operation_count"`
}

// CleanupLogsRequest 清理日志请求
type CleanupLogsRequest struct {
	BeforeTime time.Time `json:"before_time" validate:"required"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	Version       string    `json:"version"`
	StartTime     time.Time `json:"start_time"`
	Uptime        string    `json:"uptime"`
	DatabaseConnections int32  `json:"database_connections"`
	RedisConnections    int32  `json:"redis_connections"`
	MemoryUsage         string `json:"memory_usage"`
	CPUUsage            string `json:"cpu_usage"`
}

// GetOperationLogs 获取操作日志列表
func (s *SystemService) GetOperationLogs(ctx context.Context, req *OperationLogListRequest) (*OperationLogListResponse, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "AUDIT_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看操作日志")
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	// 构建查询请求
	bizReq := &biz.OperationLogListRequest{
		Page:      req.Page,
		Size:      req.Size,
		UserID:    req.UserID,
		Username:  req.Username,
		Action:    req.Action,
		Resource:  req.Resource,
		Status:    req.Status,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// 获取日志列表
	logs, total, err := s.auditUc.ListOperationLogs(ctx, bizReq)
	if err != nil {
		s.log.Errorf("Failed to list operation logs: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "操作日志获取失败")
	}

	// 转换为响应格式
	logInfos := make([]*OperationLogInfo, len(logs))
	for i, log := range logs {
		logInfos[i] = &OperationLogInfo{
			ID:            log.ID,
			UserID:        log.UserID,
			Username:      log.Username,
			Action:        log.Action,
			Resource:      log.Resource,
			ResourceID:    log.ResourceID,
			Description:   log.Description,
			IPAddress:     log.IPAddress,
			UserAgent:     log.UserAgent,
			Status:        log.Status,
			ErrorMessage:  log.ErrorMessage,
			ExecutionTime: log.ExecutionTime,
			CreatedAt:     log.CreatedAt,
		}

		// 只有超级管理员可以查看请求和响应数据
		if currentUser.HasRole("SUPER_ADMIN") {
			logInfos[i].RequestData = log.RequestData
			logInfos[i].ResponseData = log.ResponseData
		}
	}

	return &OperationLogListResponse{
		Logs:  logInfos,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// GetOperationLog 获取操作日志详情
func (s *SystemService) GetOperationLog(ctx context.Context, logID int32) (*OperationLogInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "AUDIT_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看操作日志")
	}

	// 获取日志
	log, err := s.auditUc.GetOperationLog(ctx, logID)
	if err != nil {
		return nil, errors.NotFound("OPERATION_LOG_NOT_FOUND", "操作日志不存在")
	}

	logInfo := &OperationLogInfo{
		ID:            log.ID,
		UserID:        log.UserID,
		Username:      log.Username,
		Action:        log.Action,
		Resource:      log.Resource,
		ResourceID:    log.ResourceID,
		Description:   log.Description,
		IPAddress:     log.IPAddress,
		UserAgent:     log.UserAgent,
		Status:        log.Status,
		ErrorMessage:  log.ErrorMessage,
		ExecutionTime: log.ExecutionTime,
		CreatedAt:     log.CreatedAt,
	}

	// 只有超级管理员可以查看请求和响应数据
	if currentUser.HasRole("SUPER_ADMIN") {
		logInfo.RequestData = log.RequestData
		logInfo.ResponseData = log.ResponseData
	}

	return logInfo, nil
}

// GetOperationStatistics 获取操作统计
func (s *SystemService) GetOperationStatistics(ctx context.Context, req *StatisticsRequest) (*OperationStatisticsInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "AUDIT_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看操作统计")
	}

	// 获取统计数据
	stats, err := s.auditUc.GetOperationStatistics(ctx, req.StartTime, req.EndTime)
	if err != nil {
		s.log.Errorf("Failed to get operation statistics: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "统计数据获取失败")
	}

	// 计算成功率
	var successRate float64
	if stats.TotalOperations > 0 {
		successRate = float64(stats.SuccessOperations) / float64(stats.TotalOperations) * 100
	}

	return &OperationStatisticsInfo{
		TotalOperations:      stats.TotalOperations,
		SuccessOperations:    stats.SuccessOperations,
		FailedOperations:     stats.FailedOperations,
		SuccessRate:          successRate,
		AverageExecutionTime: stats.AverageExecutionTime,
		ActiveUsers:          stats.ActiveUsers,
		ActionDistribution:   stats.ActionDistribution,
	}, nil
}

// GetTopActiveUsers 获取最活跃用户
func (s *SystemService) GetTopActiveUsers(ctx context.Context, req *StatisticsRequest, limit int32) ([]*UserActivityInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "AUDIT_MANAGER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看用户活动")
	}

	if limit <= 0 {
		limit = 10
	}

	// 获取活跃用户
	activities, err := s.auditUc.GetTopActiveUsers(ctx, req.StartTime, req.EndTime, limit)
	if err != nil {
		s.log.Errorf("Failed to get top active users: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "活跃用户数据获取失败")
	}

	// 转换为响应格式
	activityInfos := make([]*UserActivityInfo, len(activities))
	for i, activity := range activities {
		activityInfos[i] = &UserActivityInfo{
			UserID:         activity.UserID,
			Username:       activity.Username,
			Email:          activity.Email,
			OperationCount: activity.OperationCount,
		}
	}

	return activityInfos, nil
}

// CleanupOperationLogs 清理操作日志
func (s *SystemService) CleanupOperationLogs(ctx context.Context, req *CleanupLogsRequest) (int64, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasRole("SUPER_ADMIN") {
		return 0, errors.Forbidden("PERMISSION_DENIED", "无权限清理操作日志")
	}

	s.log.Infof("Cleaning up operation logs before %v by %s", req.BeforeTime, currentUser.Username)

	// 清理日志
	affected, err := s.auditUc.DeleteOperationLogs(ctx, req.BeforeTime)
	if err != nil {
		s.log.Errorf("Failed to cleanup operation logs: %v", err)
		return 0, errors.InternalServer("INTERNAL_ERROR", "日志清理失败")
	}

	s.log.Infof("Cleaned up %d operation logs", affected)
	return affected, nil
}

// GetSystemInfo 获取系统信息
func (s *SystemService) GetSystemInfo(ctx context.Context) (*SystemInfo, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看系统信息")
	}

	// TODO: 实现系统信息获取逻辑
	// 这里应该实现实际的系统监控逻辑
	startTime := time.Now().Add(-time.Hour * 24) // 模拟24小时前启动
	uptime := time.Since(startTime)

	return &SystemInfo{
		Version:             "1.0.0",
		StartTime:           startTime,
		Uptime:              uptime.String(),
		DatabaseConnections: 10,  // TODO: 实际获取数据库连接数
		RedisConnections:    5,   // TODO: 实际获取Redis连接数
		MemoryUsage:         "512MB", // TODO: 实际获取内存使用情况
		CPUUsage:            "15%",   // TODO: 实际获取CPU使用情况
	}, nil
}

// CreateOperationLog 创建操作日志
func (s *SystemService) CreateOperationLog(ctx context.Context, log *biz.OperationLog) error {
	return s.auditUc.CreateOperationLog(ctx, log)
}

// GetDashboardData 获取仪表板数据
func (s *SystemService) GetDashboardData(ctx context.Context) (map[string]interface{}, error) {
	// 检查权限
	currentUser := middleware.GetCurrentUser(ctx)
	if !currentUser.HasAnyRole("SUPER_ADMIN", "ADMIN", "DASHBOARD_VIEWER") {
		return nil, errors.Forbidden("PERMISSION_DENIED", "无权限查看仪表板")
	}

	// 获取最近24小时的统计数据
	endTime := time.Now()
	startTime := endTime.Add(-time.Hour * 24)

	stats, err := s.auditUc.GetOperationStatistics(ctx, startTime, endTime)
	if err != nil {
		s.log.Errorf("Failed to get operation statistics for dashboard: %v", err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "仪表板数据获取失败")
	}

	// 获取最活跃用户
	topUsers, err := s.auditUc.GetTopActiveUsers(ctx, startTime, endTime, 5)
	if err != nil {
		s.log.Errorf("Failed to get top active users for dashboard: %v", err)
		topUsers = []*biz.UserActivity{} // 如果获取失败，使用空数组
	}

	// 计算成功率
	var successRate float64
	if stats.TotalOperations > 0 {
		successRate = float64(stats.SuccessOperations) / float64(stats.TotalOperations) * 100
	}

	// 构建仪表板数据
	dashboardData := map[string]interface{}{
		"statistics": map[string]interface{}{
			"total_operations":       stats.TotalOperations,
			"success_operations":     stats.SuccessOperations,
			"failed_operations":      stats.FailedOperations,
			"success_rate":           successRate,
			"average_execution_time": stats.AverageExecutionTime,
			"active_users":           stats.ActiveUsers,
		},
		"action_distribution": stats.ActionDistribution,
		"top_active_users":    topUsers,
		"time_range": map[string]interface{}{
			"start_time": startTime,
			"end_time":   endTime,
		},
	}

	return dashboardData, nil
}