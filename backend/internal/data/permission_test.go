package data

import (
	"testing"

	"erp-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
)

func TestNewPermissionRepo(t *testing.T) {
	logger := log.DefaultLogger

	// 创建权限仓库
	repo := &permissionRepo{
		log: log.NewHelper(logger),
	}

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.log)
}

// 测试业务实体的验证逻辑
func TestDocType_Validate(t *testing.T) {
	tests := []struct {
		name    string
		docType *biz.DocType
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid doctype",
			docType: &biz.DocType{
				Name:   "TestDoc",
				Label:  "Test Document",
				Module: "TestModule",
			},
			wantErr: false,
		},
		{
			name: "empty name should fail",
			docType: &biz.DocType{
				Name:   "",
				Label:  "Test Document",
				Module: "TestModule",
			},
			wantErr: true,
			errMsg:  "DocType name is required",
		},
		{
			name: "empty label should fail",
			docType: &biz.DocType{
				Name:   "TestDoc",
				Label:  "",
				Module: "TestModule",
			},
			wantErr: true,
			errMsg:  "DocType label is required",
		},
		{
			name: "empty module should fail",
			docType: &biz.DocType{
				Name:   "TestDoc",
				Label:  "Test Document",
				Module: "",
			},
			wantErr: true,
			errMsg:  "DocType module is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.docType.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPermissionRule_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rule    *biz.PermissionRule
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid document level permission",
			rule: &biz.PermissionRule{
				RoleID:          1,
				DocType:         "User",
				PermissionLevel: 0,
				CanRead:         true,
				CanWrite:        true,
			},
			wantErr: false,
		},
		{
			name: "invalid role id",
			rule: &biz.PermissionRule{
				RoleID:          0,
				DocType:         "User",
				PermissionLevel: 0,
				CanRead:         true,
			},
			wantErr: true,
			errMsg:  "RoleID must be positive",
		},
		{
			name: "empty doc type",
			rule: &biz.PermissionRule{
				RoleID:          1,
				DocType:         "",
				PermissionLevel: 0,
				CanRead:         true,
			},
			wantErr: true,
			errMsg:  "DocType is required",
		},
		{
			name: "invalid permission level",
			rule: &biz.PermissionRule{
				RoleID:          1,
				DocType:         "User",
				PermissionLevel: -1,
				CanRead:         true,
			},
			wantErr: true,
			errMsg:  "PermissionLevel must be between 0 and 10",
		},
		{
			name: "document level without basic permissions",
			rule: &biz.PermissionRule{
				RoleID:          1,
				DocType:         "User",
				PermissionLevel: 0,
				CanRead:         false,
				CanWrite:        false,
			},
			wantErr: true,
			errMsg:  "at least one of CanRead or CanWrite must be true",
		},
		{
			name: "field level with document permissions should fail",
			rule: &biz.PermissionRule{
				RoleID:          1,
				DocType:         "User",
				PermissionLevel: 1,
				CanRead:         true,
				CanCreate:       true, // 字段级不应该有创建权限
			},
			wantErr: true,
			errMsg:  "field level permissions (level > 0) should only use CanRead and CanWrite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFieldPermissionLevel_Validate(t *testing.T) {
	tests := []struct {
		name    string
		field   *biz.FieldPermissionLevel
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid field permission",
			field: &biz.FieldPermissionLevel{
				DocType:         "User",
				FieldName:       "email",
				PermissionLevel: 1,
				FieldType:       "Data",
			},
			wantErr: false,
		},
		{
			name: "empty doc type",
			field: &biz.FieldPermissionLevel{
				DocType:         "",
				FieldName:       "email",
				PermissionLevel: 1,
				FieldType:       "Data",
			},
			wantErr: true,
			errMsg:  "DocType is required",
		},
		{
			name: "empty field name",
			field: &biz.FieldPermissionLevel{
				DocType:         "User",
				FieldName:       "",
				PermissionLevel: 1,
				FieldType:       "Data",
			},
			wantErr: true,
			errMsg:  "FieldName is required",
		},
		{
			name: "invalid permission level",
			field: &biz.FieldPermissionLevel{
				DocType:         "User",
				FieldName:       "email",
				PermissionLevel: -1,
				FieldType:       "Data",
			},
			wantErr: true,
			errMsg:  "PermissionLevel must be between 0 and 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.field.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
