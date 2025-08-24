package cache

import (
	"context"
	"time"

	"erp-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

// NewRedisClient 创建Redis客户端
func NewRedisClient(config *conf.CacheConfig, logger log.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         config.Redis.Addr,
		Password:     config.Redis.Password,
		DB:           config.Redis.DB,
		PoolSize:     config.Redis.PoolSize,
		MinIdleConns: config.Redis.MinIdleConns,
		DialTimeout:  config.Redis.DialTimeout,
		ReadTimeout:  config.Redis.ReadTimeout,
		WriteTimeout: config.Redis.WriteTimeout,
		PoolTimeout:  config.Redis.PoolTimeout,
		IdleTimeout:  config.Redis.IdleTimeout,
	})

	helper := log.NewHelper(logger)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		helper.Fatalf("Failed to connect to Redis: %v", err)
	}

	helper.Infof("Connected to Redis at %s", config.Redis.Addr)

	return rdb
}

// RedisHealthCheck Redis健康检查
type RedisHealthCheck struct {
	client *redis.Client
	logger *log.Helper
}

// NewRedisHealthCheck 创建Redis健康检查
func NewRedisHealthCheck(client *redis.Client, logger log.Logger) *RedisHealthCheck {
	return &RedisHealthCheck{
		client: client,
		logger: log.NewHelper(logger),
	}
}

// CheckHealth 检查Redis健康状态
func (h *RedisHealthCheck) CheckHealth(ctx context.Context) error {
	// Ping测试
	_, err := h.client.Ping(ctx).Result()
	if err != nil {
		h.logger.Errorf("Redis ping failed: %v", err)
		return err
	}

	// 简单读写测试
	testKey := "health_check_" + time.Now().Format("20060102150405")
	testValue := "ok"

	// 写测试
	err = h.client.Set(ctx, testKey, testValue, time.Minute).Err()
	if err != nil {
		h.logger.Errorf("Redis write test failed: %v", err)
		return err
	}

	// 读测试
	result, err := h.client.Get(ctx, testKey).Result()
	if err != nil {
		h.logger.Errorf("Redis read test failed: %v", err)
		return err
	}

	if result != testValue {
		h.logger.Errorf("Redis read/write consistency check failed: expected %s, got %s", testValue, result)
		return err
	}

	// 清理测试数据
	h.client.Del(ctx, testKey)

	return nil
}

// GetConnectionStats 获取连接统计信息
func (h *RedisHealthCheck) GetConnectionStats() *redis.PoolStats {
	return h.client.PoolStats()
}

// RedisMetrics Redis指标收集器
type RedisMetrics struct {
	client *redis.Client
	logger *log.Helper
}

// NewRedisMetrics 创建Redis指标收集器
func NewRedisMetrics(client *redis.Client, logger log.Logger) *RedisMetrics {
	return &RedisMetrics{
		client: client,
		logger: log.NewHelper(logger),
	}
}

// CollectMetrics 收集Redis指标
func (m *RedisMetrics) CollectMetrics(ctx context.Context) (map[string]interface{}, error) {
	stats := m.client.PoolStats()

	metrics := map[string]interface{}{
		"pool_hits":        stats.Hits,
		"pool_misses":      stats.Misses,
		"pool_timeouts":    stats.Timeouts,
		"pool_total_conns": stats.TotalConns,
		"pool_idle_conns":  stats.IdleConns,
		"pool_stale_conns": stats.StaleConns,
	}

	// 获取Redis服务器信息
	info, err := m.client.Info(ctx, "server", "memory", "stats").Result()
	if err != nil {
		m.logger.Warnf("Failed to get Redis info: %v", err)
		return metrics, nil
	}

	// 这里应该解析info字符串，提取关键指标
	_ = info // 简化处理

	// 添加更多指标
	metrics["redis_version"] = "unknown"
	metrics["used_memory"] = "unknown"
	metrics["connected_clients"] = "unknown"
	metrics["total_commands_processed"] = "unknown"

	return metrics, nil
}
