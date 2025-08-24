package pkg

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// PasswordManager 密码管理器
type PasswordManager struct {
	cost int
}

// NewPasswordManager 创建密码管理器
func NewPasswordManager() *PasswordManager {
	return &PasswordManager{
		cost: bcrypt.DefaultCost,
	}
}

// HashPassword 密码加密
func (pm *PasswordManager) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), pm.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ValidatePassword 密码验证
func (pm *PasswordManager) ValidatePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// ValidatePasswordStrength 验证密码强度
func (pm *PasswordManager) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	// 检查是否包含大写字母
	if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	// 检查是否包含小写字母
	if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	// 检查是否包含数字
	if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
		return fmt.Errorf("password must contain at least one number")
	}

	// 检查是否包含特殊字符
	if matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]`, password); !matched {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// GenerateRandomPassword 生成随机密码
func (pm *PasswordManager) GenerateRandomPassword(length int) string {
	if length < 8 {
		length = 12
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)

	for i := range password {
		randomBytes := make([]byte, 1)
		rand.Read(randomBytes)
		password[i] = charset[int(randomBytes[0])%len(charset)]
	}

	return string(password)
}

// GenerateTOTPSecret 生成TOTP密钥
func (pm *PasswordManager) GenerateTOTPSecret() string {
	secret := make([]byte, 20)
	rand.Read(secret)
	return base32.StdEncoding.EncodeToString(secret)
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateUsername 验证用户名格式
func ValidateUsername(username string) error {
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}

	if len(username) > 32 {
		return fmt.Errorf("username must not exceed 32 characters")
	}

	// 只允许字母、数字、下划线和连字符
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username); !matched {
		return fmt.Errorf("username can only contain letters, numbers, underscores and hyphens")
	}

	// 不能以数字开头
	if matched, _ := regexp.MatchString(`^[0-9]`, username); matched {
		return fmt.Errorf("username cannot start with a number")
	}

	return nil
}

// ValidatePhone 验证手机号格式
func ValidatePhone(phone string) bool {
	// 简单的手机号验证，可根据需要调整
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(phone)
}
