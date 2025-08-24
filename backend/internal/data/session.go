package data

import (
	"context"
	"database/sql"
	"time"

	"erp-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// sessionRepo 会话仓储实现
type sessionRepo struct {
	data *Data
	log  *log.Helper
}

// NewSessionRepo 创建会话仓储
func NewSessionRepo(data *Data, logger log.Logger) biz.SessionRepo {
	return &sessionRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateSession 创建会话
func (r *sessionRepo) CreateSession(ctx context.Context, session *biz.UserSession) (*biz.UserSession, error) {
	query := `
		INSERT INTO user_sessions (id, user_id, device_type, ip_address, user_agent, 
		                          location, is_active, last_activity_at, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.data.db.ExecContext(ctx, query,
		session.ID, session.UserID, session.DeviceType, session.IPAddress,
		session.UserAgent, session.Location, session.IsActive,
		session.LastActivity, session.ExpiresAt, session.CreatedAt,
	)

	if err != nil {
		r.log.Errorf("failed to create session: %v", err)
		return nil, err
	}

	return session, nil
}

// GetSession 获取会话
func (r *sessionRepo) GetSession(ctx context.Context, sessionID string) (*biz.UserSession, error) {
	var session biz.UserSession
	var lastActivity sql.NullTime

	query := `
		SELECT id, user_id, device_type, ip_address, user_agent, location,
		       is_active, last_activity_at, expires_at, created_at
		FROM user_sessions WHERE id = $1`

	err := r.data.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID, &session.UserID, &session.DeviceType, &session.IPAddress,
		&session.UserAgent, &session.Location, &session.IsActive,
		&lastActivity, &session.ExpiresAt, &session.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.log.Errorf("failed to get session: %v", err)
		return nil, err
	}

	if lastActivity.Valid {
		session.LastActivity = lastActivity.Time
	}

	return &session, nil
}

// GetUserSession 获取用户会话
func (r *sessionRepo) GetUserSession(ctx context.Context, userID int64, sessionID string) (*biz.UserSession, error) {
	var session biz.UserSession
	var lastActivity sql.NullTime

	query := `
		SELECT id, user_id, device_type, ip_address, user_agent, location,
		       is_active, last_activity_at, expires_at, created_at
		FROM user_sessions WHERE user_id = $1 AND id = $2`

	err := r.data.db.QueryRowContext(ctx, query, userID, sessionID).Scan(
		&session.ID, &session.UserID, &session.DeviceType, &session.IPAddress,
		&session.UserAgent, &session.Location, &session.IsActive,
		&lastActivity, &session.ExpiresAt, &session.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.log.Errorf("failed to get user session: %v", err)
		return nil, err
	}

	if lastActivity.Valid {
		session.LastActivity = lastActivity.Time
	}

	return &session, nil
}

// UpdateSessionActivity 更新会话活动时间
func (r *sessionRepo) UpdateSessionActivity(ctx context.Context, sessionID string) error {
	query := `UPDATE user_sessions SET last_activity_at = $1 WHERE id = $2`
	
	_, err := r.data.db.ExecContext(ctx, query, time.Now(), sessionID)
	if err != nil {
		r.log.Errorf("failed to update session activity: %v", err)
		return err
	}

	return nil
}

// DeactivateSession 停用会话
func (r *sessionRepo) DeactivateSession(ctx context.Context, sessionID string) error {
	query := `UPDATE user_sessions SET is_active = false WHERE id = $1`
	
	_, err := r.data.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		r.log.Errorf("failed to deactivate session: %v", err)
		return err
	}

	return nil
}

// DeactivateUserSessions 停用用户所有会话
func (r *sessionRepo) DeactivateUserSessions(ctx context.Context, userID int64) error {
	query := `UPDATE user_sessions SET is_active = false WHERE user_id = $1`
	
	_, err := r.data.db.ExecContext(ctx, query, userID)
	if err != nil {
		r.log.Errorf("failed to deactivate user sessions: %v", err)
		return err
	}

	return nil
}

// ListUserSessions 获取用户会话列表
func (r *sessionRepo) ListUserSessions(ctx context.Context, userID int64) ([]*biz.UserSession, error) {
	query := `
		SELECT id, user_id, device_type, ip_address, user_agent, location,
		       is_active, last_activity_at, expires_at, created_at
		FROM user_sessions 
		WHERE user_id = $1 AND expires_at > $2
		ORDER BY last_activity_at DESC`

	rows, err := r.data.db.QueryContext(ctx, query, userID, time.Now())
	if err != nil {
		r.log.Errorf("failed to list user sessions: %v", err)
		return nil, err
	}
	defer rows.Close()

	var sessions []*biz.UserSession
	for rows.Next() {
		var session biz.UserSession
		var lastActivity sql.NullTime

		err := rows.Scan(
			&session.ID, &session.UserID, &session.DeviceType, &session.IPAddress,
			&session.UserAgent, &session.Location, &session.IsActive,
			&lastActivity, &session.ExpiresAt, &session.CreatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan session: %v", err)
			return nil, err
		}

		if lastActivity.Valid {
			session.LastActivity = lastActivity.Time
		}

		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// CleanupExpiredSessions 清理过期会话
func (r *sessionRepo) CleanupExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM user_sessions WHERE expires_at < $1`
	
	result, err := r.data.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		r.log.Errorf("failed to cleanup expired sessions: %v", err)
		return err
	}

	affected, _ := result.RowsAffected()
	if affected > 0 {
		r.log.Infof("Cleaned up %d expired sessions", affected)
	}

	return nil
}