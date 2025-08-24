package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"erp-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// userFilterRepo 用户过滤器仓储实现
type userFilterRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserFilterRepo 创建用户过滤器仓储
func NewUserFilterRepo(data *Data, logger log.Logger) biz.UserFilterRepo {
	return &userFilterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateFilter 创建过滤器
func (r *userFilterRepo) CreateFilter(ctx context.Context, filter *biz.UserFilter) (*biz.UserFilter, error) {
	var id int32
	
	conditionsJSON, err := json.Marshal(filter.FilterConditions)
	if err != nil {
		r.log.Errorf("failed to marshal filter conditions: %v", err)
		return nil, err
	}

	var sortConfigJSON []byte
	if filter.SortConfig != nil {
		sortConfigJSON, err = json.Marshal(filter.SortConfig)
		if err != nil {
			r.log.Errorf("failed to marshal sort config: %v", err)
			return nil, err
		}
	}

	query := `
		INSERT INTO user_filters (user_id, module_type, filter_name, filter_conditions, sort_config, is_default, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	now := time.Now()
	err = r.data.db.QueryRowContext(ctx, query,
		filter.UserID, filter.ModuleType, filter.FilterName, conditionsJSON,
		sortConfigJSON, filter.IsDefault, filter.IsPublic, now, now,
	).Scan(&id)

	if err != nil {
		r.log.Errorf("failed to create filter: %v", err)
		return nil, err
	}

	filter.ID = id
	filter.CreatedAt = now
	filter.UpdatedAt = now
	return filter, nil
}

// GetFilter 获取过滤器
func (r *userFilterRepo) GetFilter(ctx context.Context, id int32) (*biz.UserFilter, error) {
	var filter biz.UserFilter
	var conditionsJSON, sortConfigJSON sql.NullString

	query := `
		SELECT id, user_id, module_type, filter_name, filter_conditions, sort_config, 
		       is_default, is_public, created_at, updated_at
		FROM user_filters WHERE id = $1`

	err := r.data.db.QueryRowContext(ctx, query, id).Scan(
		&filter.ID, &filter.UserID, &filter.ModuleType, &filter.FilterName,
		&conditionsJSON, &sortConfigJSON, &filter.IsDefault, &filter.IsPublic,
		&filter.CreatedAt, &filter.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, biz.ErrFilterNotFound
		}
		r.log.Errorf("failed to get filter: %v", err)
		return nil, err
	}

	// 解析JSON字段
	if conditionsJSON.Valid {
		err = json.Unmarshal([]byte(conditionsJSON.String), &filter.FilterConditions)
		if err != nil {
			r.log.Errorf("failed to unmarshal filter conditions: %v", err)
			return nil, err
		}
	}

	if sortConfigJSON.Valid && sortConfigJSON.String != "" {
		err = json.Unmarshal([]byte(sortConfigJSON.String), &filter.SortConfig)
		if err != nil {
			r.log.Errorf("failed to unmarshal sort config: %v", err)
			return nil, err
		}
	}

	return &filter, nil
}

// UpdateFilter 更新过滤器
func (r *userFilterRepo) UpdateFilter(ctx context.Context, filter *biz.UserFilter) (*biz.UserFilter, error) {
	conditionsJSON, err := json.Marshal(filter.FilterConditions)
	if err != nil {
		r.log.Errorf("failed to marshal filter conditions: %v", err)
		return nil, err
	}

	var sortConfigJSON []byte
	if filter.SortConfig != nil {
		sortConfigJSON, err = json.Marshal(filter.SortConfig)
		if err != nil {
			r.log.Errorf("failed to marshal sort config: %v", err)
			return nil, err
		}
	}

	query := `
		UPDATE user_filters 
		SET filter_name = $1, filter_conditions = $2, sort_config = $3, 
		    is_default = $4, is_public = $5, updated_at = $6
		WHERE id = $7`

	now := time.Now()
	_, err = r.data.db.ExecContext(ctx, query,
		filter.FilterName, conditionsJSON, sortConfigJSON,
		filter.IsDefault, filter.IsPublic, now, filter.ID,
	)

	if err != nil {
		r.log.Errorf("failed to update filter: %v", err)
		return nil, err
	}

	filter.UpdatedAt = now
	return filter, nil
}

// DeleteFilter 删除过滤器
func (r *userFilterRepo) DeleteFilter(ctx context.Context, id int32) error {
	query := "DELETE FROM user_filters WHERE id = $1"
	_, err := r.data.db.ExecContext(ctx, query, id)
	if err != nil {
		r.log.Errorf("failed to delete filter: %v", err)
		return err
	}
	return nil
}

// ListFilters 获取过滤器列表
func (r *userFilterRepo) ListFilters(ctx context.Context, userID int32, moduleType string) ([]*biz.UserFilter, error) {
	query := `
		SELECT id, user_id, module_type, filter_name, filter_conditions, sort_config,
		       is_default, is_public, created_at, updated_at
		FROM user_filters 
		WHERE (user_id = $1 OR is_public = true) AND module_type = $2
		ORDER BY is_default DESC, created_at DESC`

	rows, err := r.data.db.QueryContext(ctx, query, userID, moduleType)
	if err != nil {
		r.log.Errorf("failed to list filters: %v", err)
		return nil, err
	}
	defer rows.Close()

	var filters []*biz.UserFilter
	for rows.Next() {
		var filter biz.UserFilter
		var conditionsJSON, sortConfigJSON sql.NullString

		err := rows.Scan(
			&filter.ID, &filter.UserID, &filter.ModuleType, &filter.FilterName,
			&conditionsJSON, &sortConfigJSON, &filter.IsDefault, &filter.IsPublic,
			&filter.CreatedAt, &filter.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan filter: %v", err)
			return nil, err
		}

		// 解析JSON字段
		if conditionsJSON.Valid {
			err = json.Unmarshal([]byte(conditionsJSON.String), &filter.FilterConditions)
			if err != nil {
				r.log.Errorf("failed to unmarshal filter conditions: %v", err)
				continue
			}
		}

		if sortConfigJSON.Valid && sortConfigJSON.String != "" {
			err = json.Unmarshal([]byte(sortConfigJSON.String), &filter.SortConfig)
			if err != nil {
				r.log.Errorf("failed to unmarshal sort config: %v", err)
				// 排序配置解析失败不影响过滤器使用
			}
		}

		filters = append(filters, &filter)
	}

	return filters, nil
}

// SetDefaultFilter 设置默认过滤器
func (r *userFilterRepo) SetDefaultFilter(ctx context.Context, userID int32, moduleType string, filterID int32) error {
	tx, err := r.data.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 清除用户在该模块的所有默认过滤器
	_, err = tx.ExecContext(ctx,
		"UPDATE user_filters SET is_default = false WHERE user_id = $1 AND module_type = $2",
		userID, moduleType)
	if err != nil {
		r.log.Errorf("failed to clear default filters: %v", err)
		return err
	}

	// 设置新的默认过滤器
	_, err = tx.ExecContext(ctx,
		"UPDATE user_filters SET is_default = true WHERE id = $1 AND user_id = $2",
		filterID, userID)
	if err != nil {
		r.log.Errorf("failed to set default filter: %v", err)
		return err
	}

	return tx.Commit()
}