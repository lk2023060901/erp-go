package service

import (
	"context"
	"testing"

	"erp-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock usecase for testing
type MockPermissionUsecase struct {
	mock.Mock
}

// 确保 MockPermissionUsecase 实现了 PermissionUsecaseInterface 接口
var _ biz.PermissionUsecaseInterface = (*MockPermissionUsecase)(nil)

// 文档类型管理
func (m *MockPermissionUsecase) CreateDocType(ctx context.Context, docType *biz.DocType) (*biz.DocType, error) {
	args := m.Called(ctx, docType)
	return args.Get(0).(*biz.DocType), args.Error(1)
}

func (m *MockPermissionUsecase) GetDocType(ctx context.Context, name string) (*biz.DocType, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*biz.DocType), args.Error(1)
}

func (m *MockPermissionUsecase) ListDocTypes(ctx context.Context, module string) ([]*biz.DocType, error) {
	args := m.Called(ctx, module)
	return args.Get(0).([]*biz.DocType), args.Error(1)
}

func (m *MockPermissionUsecase) UpdateDocType(ctx context.Context, docType *biz.DocType) (*biz.DocType, error) {
	args := m.Called(ctx, docType)
	return args.Get(0).(*biz.DocType), args.Error(1)
}

func (m *MockPermissionUsecase) DeleteDocType(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

// 权限规则管理
func (m *MockPermissionUsecase) CreatePermissionRule(ctx context.Context, rule *biz.PermissionRule) (*biz.PermissionRule, error) {
	args := m.Called(ctx, rule)
	return args.Get(0).(*biz.PermissionRule), args.Error(1)
}

func (m *MockPermissionUsecase) CreatePermissionRuleFromRequest(ctx context.Context, req *biz.CreatePermissionRuleRequest) (*biz.PermissionRule, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*biz.PermissionRule), args.Error(1)
}

func (m *MockPermissionUsecase) GetPermissionRule(ctx context.Context, id int64) (*biz.PermissionRule, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*biz.PermissionRule), args.Error(1)
}

func (m *MockPermissionUsecase) ListPermissionRules(ctx context.Context, roleID int64, docType string) ([]*biz.PermissionRule, error) {
	args := m.Called(ctx, roleID, docType)
	return args.Get(0).([]*biz.PermissionRule), args.Error(1)
}

func (m *MockPermissionUsecase) UpdatePermissionRule(ctx context.Context, rule *biz.PermissionRule) (*biz.PermissionRule, error) {
	args := m.Called(ctx, rule)
	return args.Get(0).(*biz.PermissionRule), args.Error(1)
}

func (m *MockPermissionUsecase) DeletePermissionRule(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPermissionUsecase) BatchCreatePermissionRules(ctx context.Context, rules []*biz.PermissionRule) error {
	args := m.Called(ctx, rules)
	return args.Error(0)
}

// 用户权限管理
func (m *MockPermissionUsecase) CreateUserPermission(ctx context.Context, userPermission *biz.UserPermission) (*biz.UserPermission, error) {
	args := m.Called(ctx, userPermission)
	return args.Get(0).(*biz.UserPermission), args.Error(1)
}

func (m *MockPermissionUsecase) CreateUserPermissionFromRequest(ctx context.Context, req *biz.CreateUserPermissionRequest) (*biz.UserPermission, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*biz.UserPermission), args.Error(1)
}

func (m *MockPermissionUsecase) GetUserPermission(ctx context.Context, id int64) (*biz.UserPermission, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*biz.UserPermission), args.Error(1)
}

func (m *MockPermissionUsecase) ListUserPermissions(ctx context.Context, userID int64, docType string, page, size int32) ([]*biz.UserPermission, error) {
	args := m.Called(ctx, userID, docType, page, size)
	return args.Get(0).([]*biz.UserPermission), args.Error(1)
}

func (m *MockPermissionUsecase) GetUserPermissionsCount(ctx context.Context, userID int64, docType string) (int32, error) {
	args := m.Called(ctx, userID, docType)
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockPermissionUsecase) UpdateUserPermission(ctx context.Context, userPerm *biz.UserPermission) (*biz.UserPermission, error) {
	args := m.Called(ctx, userPerm)
	return args.Get(0).(*biz.UserPermission), args.Error(1)
}

func (m *MockPermissionUsecase) DeleteUserPermission(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPermissionUsecase) BatchCreateUserPermissions(ctx context.Context, permissions []*biz.UserPermission) error {
	args := m.Called(ctx, permissions)
	return args.Error(0)
}

// 字段权限级别管理
func (m *MockPermissionUsecase) CreateFieldPermissionLevel(ctx context.Context, fieldPerm *biz.FieldPermissionLevel) (*biz.FieldPermissionLevel, error) {
	args := m.Called(ctx, fieldPerm)
	return args.Get(0).(*biz.FieldPermissionLevel), args.Error(1)
}

func (m *MockPermissionUsecase) GetFieldPermissionLevel(ctx context.Context, id int64) (*biz.FieldPermissionLevel, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*biz.FieldPermissionLevel), args.Error(1)
}

func (m *MockPermissionUsecase) UpdateFieldPermissionLevel(ctx context.Context, fieldPerm *biz.FieldPermissionLevel) (*biz.FieldPermissionLevel, error) {
	args := m.Called(ctx, fieldPerm)
	return args.Get(0).(*biz.FieldPermissionLevel), args.Error(1)
}

func (m *MockPermissionUsecase) DeleteFieldPermissionLevel(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPermissionUsecase) ListFieldPermissionLevels(ctx context.Context, docType string, page, size int32) ([]*biz.FieldPermissionLevel, error) {
	args := m.Called(ctx, docType, page, size)
	return args.Get(0).([]*biz.FieldPermissionLevel), args.Error(1)
}

func (m *MockPermissionUsecase) GetFieldPermissionLevelsCount(ctx context.Context, docType string) (int32, error) {
	args := m.Called(ctx, docType)
	return args.Get(0).(int32), args.Error(1)
}

// 文档工作流状态管理
func (m *MockPermissionUsecase) CreateDocumentWorkflowState(ctx context.Context, state *biz.DocumentWorkflowState) (*biz.DocumentWorkflowState, error) {
	args := m.Called(ctx, state)
	return args.Get(0).(*biz.DocumentWorkflowState), args.Error(1)
}

func (m *MockPermissionUsecase) GetDocumentWorkflowState(ctx context.Context, stateID int64) (*biz.DocumentWorkflowState, error) {
	args := m.Called(ctx, stateID)
	return args.Get(0).(*biz.DocumentWorkflowState), args.Error(1)
}

func (m *MockPermissionUsecase) UpdateDocumentWorkflowState(ctx context.Context, state *biz.DocumentWorkflowState) (*biz.DocumentWorkflowState, error) {
	args := m.Called(ctx, state)
	return args.Get(0).(*biz.DocumentWorkflowState), args.Error(1)
}

func (m *MockPermissionUsecase) DeleteDocumentWorkflowState(ctx context.Context, stateID int64) error {
	args := m.Called(ctx, stateID)
	return args.Error(0)
}

func (m *MockPermissionUsecase) ListDocumentWorkflowStates(ctx context.Context, docType, documentName, state string, userID int64, page, size int32) ([]*biz.DocumentWorkflowState, error) {
	args := m.Called(ctx, docType, documentName, state, userID, page, size)
	return args.Get(0).([]*biz.DocumentWorkflowState), args.Error(1)
}

func (m *MockPermissionUsecase) GetDocumentWorkflowStatesCount(ctx context.Context, docType, documentName, state string, userID int64) (int32, error) {
	args := m.Called(ctx, docType, documentName, state, userID)
	return args.Get(0).(int32), args.Error(1)
}

// 权限检查和查询
func (m *MockPermissionUsecase) CheckDocumentPermission(ctx context.Context, req *biz.PermissionCheckRequest) (*biz.PermissionCheckResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*biz.PermissionCheckResponse), args.Error(1)
}

func (m *MockPermissionUsecase) CheckUserPermission(ctx context.Context, userID int64, permissionCode string) (bool, error) {
	args := m.Called(ctx, userID, permissionCode)
	return args.Bool(0), args.Error(1)
}

func (m *MockPermissionUsecase) CheckPermission(ctx context.Context, userID int64, documentType, action string, permissionLevel int) (bool, error) {
	args := m.Called(ctx, userID, documentType, action, permissionLevel)
	return args.Bool(0), args.Error(1)
}

func (m *MockPermissionUsecase) GetUserPermissionLevel(ctx context.Context, userID int64, documentType string) (int, error) {
	args := m.Called(ctx, userID, documentType)
	return args.Int(0), args.Error(1)
}

func (m *MockPermissionUsecase) GetUserEnhancedPermissions(ctx context.Context, userID int64, docType string) ([]*biz.EnhancedUserPermission, error) {
	args := m.Called(ctx, userID, docType)
	return args.Get(0).([]*biz.EnhancedUserPermission), args.Error(1)
}

func (m *MockPermissionUsecase) GetAccessibleFields(ctx context.Context, req *biz.FieldPermissionRequest) (*biz.FieldPermissionResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*biz.FieldPermissionResponse), args.Error(1)
}

func (m *MockPermissionUsecase) FilterDocumentsByPermission(ctx context.Context, userID int64, documentType string, documents []map[string]interface{}) ([]map[string]interface{}, error) {
	args := m.Called(ctx, userID, documentType, documents)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockPermissionUsecase) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func TestNewFrappePermissionService(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	logger := log.DefaultLogger

	service := NewFrappePermissionService(mockUsecase, logger)
	
	assert.NotNil(t, service)
	assert.NotNil(t, service.permissionUc)
	assert.NotNil(t, service.log)
}

// Service层验证逻辑测试

func TestFrappePermissionService_CreateDocType_Validation(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	logger := log.DefaultLogger
	service := NewFrappePermissionService(mockUsecase, logger)
	ctx := context.Background()

	t.Run("valid request should pass validation", func(t *testing.T) {
		req := &CreateDocTypeRequest{
			Name:   "TestDoc",
			Label:  "Test Document",
			Module: "TestModule",
		}

		expectedDocType := &biz.DocType{
			ID:      1,
			Name:    "TestDoc",
			Label:   "TestDoc",
			Module:  "TestModule",
			Version: 1,
		}

		mockUsecase.On("CreateDocType", ctx, mock.AnythingOfType("*biz.DocType")).Return(expectedDocType, nil)

		result, err := service.CreateDocType(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "TestDoc", result.Name)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("empty name should fail validation", func(t *testing.T) {
		req := &CreateDocTypeRequest{
			Name:   "", // 空名称应该失败
			Label:  "Test Document",
			Module: "TestModule",
		}

		result, err := service.CreateDocType(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "数据验证失败")
	})

	t.Run("empty module should fail validation", func(t *testing.T) {
		req := &CreateDocTypeRequest{
			Name:   "TestDoc",
			Label:  "Test Document",
			Module: "", // 空模块应该失败
		}

		result, err := service.CreateDocType(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "数据验证失败")
	})
}

func TestFrappePermissionService_CreatePermissionRule_Validation(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	logger := log.DefaultLogger
	service := NewFrappePermissionService(mockUsecase, logger)
	ctx := context.Background()

	t.Run("valid document level request should pass", func(t *testing.T) {
		req := &CreatePermissionRuleRequest{
			RoleID:          1,
			DocType:         "User",
			PermissionLevel: 0,
			CanRead:         true,
			CanWrite:        true,
		}

		expectedRule := &biz.PermissionRule{
			ID:              1,
			RoleID:          1,
			DocType:         "User",
			PermissionLevel: 0,
			CanRead:         true,
			CanWrite:        true,
		}

		mockUsecase.On("CreatePermissionRule", ctx, mock.AnythingOfType("*biz.PermissionRule")).Return(expectedRule, nil)

		result, err := service.CreatePermissionRule(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("invalid role id should fail", func(t *testing.T) {
		req := &CreatePermissionRuleRequest{
			RoleID:          0, // 无效的角色ID
			DocType:         "User",
			PermissionLevel: 0,
			CanRead:         true,
		}

		result, err := service.CreatePermissionRule(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "请求数据验证失败")
	})

	t.Run("document level without basic permissions should fail", func(t *testing.T) {
		req := &CreatePermissionRuleRequest{
			RoleID:          1,
			DocType:         "User",
			PermissionLevel: 0,
			CanRead:         false,
			CanWrite:        false, // 文档级权限至少需要读或写
		}

		result, err := service.CreatePermissionRule(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "权限规则验证失败")
	})

	t.Run("field level with document permissions should fail", func(t *testing.T) {
		req := &CreatePermissionRuleRequest{
			RoleID:          1,
			DocType:         "User",
			PermissionLevel: 1, // 字段级权限
			CanRead:         true,
			CanCreate:       true, // 字段级不应该有创建权限
		}

		result, err := service.CreatePermissionRule(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "权限规则验证失败")
	})
}

func TestFrappePermissionService_CheckDocumentPermission_Validation(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	logger := log.DefaultLogger
	service := NewFrappePermissionService(mockUsecase, logger)
	ctx := context.Background()

	t.Run("valid request should pass validation", func(t *testing.T) {
		req := &CheckDocumentPermissionRequest{
			UserID:     1,
			DocType:    "User",
			Permission: "read",
		}

		expectedResponse := &biz.PermissionCheckResponse{
			HasPermission: true,
		}

		mockUsecase.On("CheckDocumentPermission", ctx, mock.AnythingOfType("*biz.PermissionCheckRequest")).Return(expectedResponse, nil)

		result, err := service.CheckDocumentPermission(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.HasPermission)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("invalid user id should fail validation", func(t *testing.T) {
		req := &CheckDocumentPermissionRequest{
			UserID:     0, // 无效的用户ID
			DocType:    "User",
			Permission: "read",
		}

		result, err := service.CheckDocumentPermission(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "权限检查请求验证失败")
	})

	t.Run("empty doc type should fail validation", func(t *testing.T) {
		req := &CheckDocumentPermissionRequest{
			UserID:     1,
			DocType:    "", // 空文档类型
			Permission: "read",
		}

		result, err := service.CheckDocumentPermission(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "权限检查请求验证失败")
	})

	t.Run("invalid permission type should fail validation", func(t *testing.T) {
		req := &CheckDocumentPermissionRequest{
			UserID:     1,
			DocType:    "User",
			Permission: "invalid_permission", // 无效的权限类型
		}

		result, err := service.CheckDocumentPermission(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "权限检查请求验证失败")
	})
}