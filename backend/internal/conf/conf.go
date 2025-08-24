package conf

import "time"

// Bootstrap 启动配置
type Bootstrap struct {
	Server   *Server   `yaml:"server"`
	Data     *Data     `yaml:"data"`
	Auth     *Auth     `yaml:"auth"`
	Log      *Log      `yaml:"log"`
	Upload   *Upload   `yaml:"upload"`
	Cors     *Cors     `yaml:"cors"`
	Security *Security `yaml:"security"`
}

// Server 服务器配置
type Server struct {
	Http *HTTP `yaml:"http"`
	Grpc *GRPC `yaml:"grpc"`
}

// HTTP HTTP服务器配置
type HTTP struct {
	Network string `yaml:"network"`
	Addr    string `yaml:"addr"`
	Timeout string `yaml:"timeout"`
}

// AsDuration 返回超时时间
func (h *HTTP) GetTimeout() time.Duration {
	if h.Timeout == "" {
		return 30 * time.Second
	}
	duration, err := time.ParseDuration(h.Timeout)
	if err != nil {
		return 30 * time.Second
	}
	return duration
}

// GRPC gRPC服务器配置
type GRPC struct {
	Network string `yaml:"network"`
	Addr    string `yaml:"addr"`
	Timeout string `yaml:"timeout"`
}

// AsDuration 返回超时时间
func (g *GRPC) GetTimeout() time.Duration {
	if g.Timeout == "" {
		return 30 * time.Second
	}
	duration, err := time.ParseDuration(g.Timeout)
	if err != nil {
		return 30 * time.Second
	}
	return duration
}

// Data 数据配置
type Data struct {
	Database *Database `yaml:"database"`
	Redis    *Redis    `yaml:"redis"`
	Jwt      *JWT      `yaml:"jwt"`
}

// JWT JWT配置
type JWT struct {
	SecretKey          string `yaml:"secret_key"`
	AccessTokenExpire  int64  `yaml:"access_token_expire"`
	RefreshTokenExpire int64  `yaml:"refresh_token_expire"`
}

// Database 数据库配置
type Database struct {
	Driver string `yaml:"driver"`
	Source string `yaml:"source"`
}

// Redis Redis配置
type Redis struct {
	Network  string `yaml:"network"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// Auth 认证配置
type Auth struct {
	JwtSecret          string `yaml:"jwt_secret"`
	JwtExpire          string `yaml:"jwt_expire"`
	RefreshTokenExpire string `yaml:"refresh_token_expire"`
}

// Log 日志配置
type Log struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Upload 文件上传配置
type Upload struct {
	MaxSize      string   `yaml:"max_size"`
	Path         string   `yaml:"path"`
	AllowedTypes []string `yaml:"allowed_types"`
}

// Cors 跨域配置
type Cors struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	AllowHeaders     []string `yaml:"allow_headers"`
	AllowMethods     []string `yaml:"allow_methods"`
}

// Security 安全配置
type Security struct {
	RateLimit *RateLimit `yaml:"rate_limit"`
	TOTP      *TOTP      `yaml:"totp"`
	Password  *Password  `yaml:"password"`
}

// RateLimit 限流配置
type RateLimit struct {
	Enabled           bool `yaml:"enabled"`
	RequestsPerMinute int  `yaml:"requests_per_minute"`
}

// TOTP 双因素认证配置
type TOTP struct {
	Issuer string `yaml:"issuer"`
}

// Password 密码策略配置
type Password struct {
	MinLength      int  `yaml:"min_length"`
	RequireSpecial bool `yaml:"require_special"`
	RequireNumber  bool `yaml:"require_number"`
}
