package conf

import (
	"time"
)

// CacheConfig 缓存配置
type CacheConfig struct {
	Redis       RedisConfig       `yaml:"redis"`
	TTL         TTLConfig         `yaml:"ttl"`
	KeyPrefix   string            `yaml:"key_prefix"`
	Warmup      WarmupConfig      `yaml:"warmup"`
	Cleanup     CleanupConfig     `yaml:"cleanup"`
	Monitoring  MonitoringConfig  `yaml:"monitoring"`
	Performance PerformanceConfig `yaml:"performance"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string        `yaml:"addr"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	PoolSize     int           `yaml:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	PoolTimeout  time.Duration `yaml:"pool_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// TTLConfig 缓存TTL配置
type TTLConfig struct {
	UserPermissions     int `yaml:"user_permissions"`
	UserRoles           int `yaml:"user_roles"`
	PermissionRules     int `yaml:"permission_rules"`
	FieldPermissions    int `yaml:"field_permissions"`
	DocTypes            int `yaml:"doc_types"`
	UserPermissionLevel int `yaml:"user_permission_level"`
}

// WarmupConfig 预热配置
type WarmupConfig struct {
	Enabled           bool             `yaml:"enabled"`
	OnStartup         bool             `yaml:"on_startup"`
	ScheduledInterval string           `yaml:"scheduled_interval"`
	Strategies        []WarmupStrategy `yaml:"strategies"`
}

// WarmupStrategy 预热策略
type WarmupStrategy struct {
	Name       string   `yaml:"name"`
	CacheTypes []string `yaml:"cache_types"`
	Schedule   string   `yaml:"schedule"`
}

// CleanupConfig 清理配置
type CleanupConfig struct {
	Enabled        bool   `yaml:"enabled"`
	Schedule       string `yaml:"schedule"`
	MaxMemoryUsage string `yaml:"max_memory_usage"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Enabled            bool         `yaml:"enabled"`
	MetricsInterval    string       `yaml:"metrics_interval"`
	SlowQueryThreshold string       `yaml:"slow_query_threshold"`
	Alerts             AlertsConfig `yaml:"alerts"`
}

// AlertsConfig 告警配置
type AlertsConfig struct {
	CacheHitRateThreshold float64 `yaml:"cache_hit_rate_threshold"`
	MemoryUsageThreshold  float64 `yaml:"memory_usage_threshold"`
	ErrorRateThreshold    float64 `yaml:"error_rate_threshold"`
}

// PerformanceConfig 性能配置
type PerformanceConfig struct {
	BatchSize             int                  `yaml:"batch_size"`
	MaxConcurrentRequests int                  `yaml:"max_concurrent_requests"`
	CircuitBreaker        CircuitBreakerConfig `yaml:"circuit_breaker"`
}

// CircuitBreakerConfig 熔断器配置
type CircuitBreakerConfig struct {
	Enabled          bool          `yaml:"enabled"`
	FailureThreshold int           `yaml:"failure_threshold"`
	Timeout          time.Duration `yaml:"timeout"`
	RecoveryTimeout  time.Duration `yaml:"recovery_timeout"`
}

// GetCacheTTL 获取指定类型的缓存TTL
func (c *CacheConfig) GetCacheTTL(cacheType string) time.Duration {
	var seconds int

	switch cacheType {
	case "user_permissions":
		seconds = c.TTL.UserPermissions
	case "user_roles":
		seconds = c.TTL.UserRoles
	case "permission_rules":
		seconds = c.TTL.PermissionRules
	case "field_permissions":
		seconds = c.TTL.FieldPermissions
	case "doc_types":
		seconds = c.TTL.DocTypes
	case "user_permission_level":
		seconds = c.TTL.UserPermissionLevel
	default:
		seconds = 900 // 默认15分钟
	}

	return time.Duration(seconds) * time.Second
}
