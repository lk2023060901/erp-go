package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"erp-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// auditRepo 审计仓储实现
type auditRepo struct {
	data *Data
	log  *log.Helper
}

// NewAuditRepo 创建审计仓储
func NewAuditRepo(data *Data, logger log.Logger) biz.AuditRepo {
	return &auditRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateOperationLog 创建操作日志
func (r *auditRepo) CreateOperationLog(ctx context.Context, log *biz.OperationLog) error {
	query := `
		INSERT INTO operation_logs (user_id, username, action, resource, resource_id, description,
		                           ip_address, user_agent, request_data, response_data, 
		                           status, error_message, execution_time, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err := r.data.db.ExecContext(ctx, query,
		log.UserID, log.Username, log.Action, log.Resource, log.ResourceID,
		log.Description, log.IPAddress, log.UserAgent, log.RequestData,
		log.ResponseData, log.Status, log.ErrorMessage, log.ExecutionTime,
		log.CreatedAt,
	)

	if err != nil {
		r.log.Errorf("failed to create operation log: %v", err)
		return err
	}

	return nil
}

// GetOperationLog 获取操作日志
func (r *auditRepo) GetOperationLog(ctx context.Context, id int32) (*biz.OperationLog, error) {
	var log biz.OperationLog
	var userID sql.NullInt32

	query := `
		SELECT id, user_id, username, action, resource, resource_id, description,
		       ip_address, user_agent, request_data, response_data,
		       status, error_message, execution_time, created_at
		FROM operation_logs WHERE id = $1`

	err := r.data.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID, &userID, &log.Username, &log.Action, &log.Resource,
		&log.ResourceID, &log.Description, &log.IPAddress, &log.UserAgent,
		&log.RequestData, &log.ResponseData, &log.Status, &log.ErrorMessage,
		&log.ExecutionTime, &log.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("operation log not found")
		}
		r.log.Errorf("failed to get operation log: %v", err)
		return nil, err
	}

	if userID.Valid {
		log.UserID = &userID.Int32
	}

	return &log, nil
}

// ListOperationLogs 操作日志列表
func (r *auditRepo) ListOperationLogs(ctx context.Context, req *biz.OperationLogListRequest) ([]*biz.OperationLog, int32, error) {
	offset := (req.Page - 1) * req.Size
	var logs []*biz.OperationLog
	var total int32

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if req.UserID != nil {
		whereClause += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *req.UserID)
		argIndex++
	}

	if req.Username != "" {
		whereClause += fmt.Sprintf(" AND username ILIKE $%d", argIndex)
		args = append(args, "%"+req.Username+"%")
		argIndex++
	}

	if req.Action != "" {
		whereClause += fmt.Sprintf(" AND action = $%d", argIndex)
		args = append(args, req.Action)
		argIndex++
	}

	if req.Resource != "" {
		whereClause += fmt.Sprintf(" AND resource = $%d", argIndex)
		args = append(args, req.Resource)
		argIndex++
	}

	if req.Status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, req.Status)
		argIndex++
	}

	if !req.StartTime.IsZero() {
		whereClause += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, req.StartTime)
		argIndex++
	}

	if !req.EndTime.IsZero() {
		whereClause += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, req.EndTime)
		argIndex++
	}

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM operation_logs " + whereClause
	err := r.data.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		r.log.Errorf("failed to count operation logs: %v", err)
		return nil, 0, err
	}

	// 查询数据
	query := fmt.Sprintf(`
		SELECT id, user_id, username, action, resource, resource_id, description,
		       ip_address, user_agent, request_data, response_data,
		       status, error_message, execution_time, created_at
		FROM operation_logs %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, req.Size, offset)

	rows, err := r.data.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.log.Errorf("failed to list operation logs: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var log biz.OperationLog
		var userID sql.NullInt32

		err := rows.Scan(
			&log.ID, &userID, &log.Username, &log.Action, &log.Resource,
			&log.ResourceID, &log.Description, &log.IPAddress, &log.UserAgent,
			&log.RequestData, &log.ResponseData, &log.Status, &log.ErrorMessage,
			&log.ExecutionTime, &log.CreatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan operation log: %v", err)
			return nil, 0, err
		}

		if userID.Valid {
			log.UserID = &userID.Int32
		}

		logs = append(logs, &log)
	}

	return logs, total, nil
}

// DeleteOperationLogs 删除操作日志
func (r *auditRepo) DeleteOperationLogs(ctx context.Context, beforeTime time.Time) (int64, error) {
	query := "DELETE FROM operation_logs WHERE created_at < $1"
	
	result, err := r.data.db.ExecContext(ctx, query, beforeTime)
	if err != nil {
		r.log.Errorf("failed to delete operation logs: %v", err)
		return 0, err
	}

	affected, _ := result.RowsAffected()
	return affected, nil
}

// GetOperationStatistics 获取操作统计
func (r *auditRepo) GetOperationStatistics(ctx context.Context, startTime, endTime time.Time) (*biz.OperationStatistics, error) {
	var stats biz.OperationStatistics

	// 总操作数
	totalQuery := "SELECT COUNT(*) FROM operation_logs WHERE created_at BETWEEN $1 AND $2"
	err := r.data.db.QueryRowContext(ctx, totalQuery, startTime, endTime).Scan(&stats.TotalOperations)
	if err != nil {
		r.log.Errorf("failed to get total operations: %v", err)
		return nil, err
	}

	// 成功操作数
	successQuery := "SELECT COUNT(*) FROM operation_logs WHERE created_at BETWEEN $1 AND $2 AND status = 'success'"
	err = r.data.db.QueryRowContext(ctx, successQuery, startTime, endTime).Scan(&stats.SuccessOperations)
	if err != nil {
		r.log.Errorf("failed to get success operations: %v", err)
		return nil, err
	}

	// 失败操作数
	failedQuery := "SELECT COUNT(*) FROM operation_logs WHERE created_at BETWEEN $1 AND $2 AND status = 'failed'"
	err = r.data.db.QueryRowContext(ctx, failedQuery, startTime, endTime).Scan(&stats.FailedOperations)
	if err != nil {
		r.log.Errorf("failed to get failed operations: %v", err)
		return nil, err
	}

	// 平均执行时间
	avgTimeQuery := "SELECT COALESCE(AVG(execution_time), 0) FROM operation_logs WHERE created_at BETWEEN $1 AND $2"
	err = r.data.db.QueryRowContext(ctx, avgTimeQuery, startTime, endTime).Scan(&stats.AverageExecutionTime)
	if err != nil {
		r.log.Errorf("failed to get average execution time: %v", err)
		return nil, err
	}

	// 活跃用户数
	activeUsersQuery := "SELECT COUNT(DISTINCT user_id) FROM operation_logs WHERE created_at BETWEEN $1 AND $2 AND user_id IS NOT NULL"
	err = r.data.db.QueryRowContext(ctx, activeUsersQuery, startTime, endTime).Scan(&stats.ActiveUsers)
	if err != nil {
		r.log.Errorf("failed to get active users: %v", err)
		return nil, err
	}

	// 操作类型分布
	actionQuery := `
		SELECT action, COUNT(*) as count 
		FROM operation_logs 
		WHERE created_at BETWEEN $1 AND $2 
		GROUP BY action 
		ORDER BY count DESC`

	rows, err := r.data.db.QueryContext(ctx, actionQuery, startTime, endTime)
	if err != nil {
		r.log.Errorf("failed to get action distribution: %v", err)
		return nil, err
	}
	defer rows.Close()

	stats.ActionDistribution = make(map[string]int32)
	for rows.Next() {
		var action string
		var count int32
		if err := rows.Scan(&action, &count); err != nil {
			r.log.Errorf("failed to scan action distribution: %v", err)
			continue
		}
		stats.ActionDistribution[action] = count
	}

	return &stats, nil
}

// GetTopActiveUsers 获取最活跃用户
func (r *auditRepo) GetTopActiveUsers(ctx context.Context, startTime, endTime time.Time, limit int32) ([]*biz.UserActivity, error) {
	query := `
		SELECT u.id, u.username, u.email, COUNT(ol.id) as operation_count
		FROM users u
		INNER JOIN operation_logs ol ON u.id = ol.user_id
		WHERE ol.created_at BETWEEN $1 AND $2
		GROUP BY u.id, u.username, u.email
		ORDER BY operation_count DESC
		LIMIT $3`

	rows, err := r.data.db.QueryContext(ctx, query, startTime, endTime, limit)
	if err != nil {
		r.log.Errorf("failed to get top active users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var activities []*biz.UserActivity
	for rows.Next() {
		var activity biz.UserActivity
		err := rows.Scan(&activity.UserID, &activity.Username, &activity.Email, &activity.OperationCount)
		if err != nil {
			r.log.Errorf("failed to scan user activity: %v", err)
			return nil, err
		}
		activities = append(activities, &activity)
	}

	return activities, nil
}