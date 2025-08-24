package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager JWT管理器
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// CustomClaims 自定义JWT声明
type CustomClaims struct {
	UserID       int64    `json:"user_id"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
	Permissions  []string `json:"permissions"`
	SessionID    string   `json:"session_id"`
	TokenType    string   `json:"token_type"` // access, refresh
	jwt.RegisteredClaims
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// Generate 生成JWT令牌
func (manager *JWTManager) Generate(userID int64, username, email string, roles, permissions []string, sessionID, tokenType string) (string, error) {
	claims := CustomClaims{
		UserID:      userID,
		Username:    username,
		Email:       email,
		Roles:       roles,
		Permissions: permissions,
		SessionID:   sessionID,
		TokenType:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "erp-system",
			Subject:   username,
			ID:        sessionID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

// Verify 验证JWT令牌
func (manager *JWTManager) Verify(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected token signing method")
			}
			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// GenerateRefreshToken 生成刷新令牌
func (manager *JWTManager) GenerateRefreshToken(userID int64, username, sessionID string) (string, error) {
	refreshClaims := CustomClaims{
		UserID:    userID,
		Username:  username,
		SessionID: sessionID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)), // 30 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "erp-system",
			Subject:   username,
			ID:        sessionID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	return token.SignedString([]byte(manager.secretKey))
}

// ParseToken 解析令牌但不验证过期
func (manager *JWTManager) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		// 如果只是过期错误，仍然返回claims (jwt v5 uses different error handling)
		if token != nil && token.Claims != nil {
			if claims, ok := token.Claims.(*CustomClaims); ok {
				return claims, nil
			}
		}
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}