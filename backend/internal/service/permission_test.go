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

func TestNewPermissionService(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	logger := log.DefaultLogger

	service := NewPermissionService(mockUsecase, logger)

	assert.NotNil(t, service)
	assert.NotNil(t, service.permissionUc)
	assert.NotNil(t, service.log)
}

// Service层验证逻辑测试

func TestPermissionService_CreateDocType_Validation(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	logger := log.DefaultLogger
	service := NewPermissionService(mockUsecase, logger)
	ctx := context.Background()

	t.Run("valid request should pass validation", func(t *testing.T) {
		req := &CreateDocTypeRequest{
			Name:   "TestDoc",
			Label:  "Test Document",
			Module: "Core",
		}

		// Mock the usecase call
		expectedDocType := &biz.DocType{
			ID:     1,
			Name:   req.Name,
			Label:  req.Label,
			Module: req.Module,
		}

		mockUsecase.On("CreateDocType", ctx, mock.AnythingOfType("*biz.DocType")).Return(expectedDocType, nil)

		result, err := service.CreateDocType(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
		assert.Equal(t, req.Label, result.Label)
		assert.Equal(t, req.Module, result.Module)
		mockUsecase.AssertExpectations(t)
	})
}

func TestPermissionService_CreatePermissionRule_Validation(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	logger := log.DefaultLogger
	service := NewPermissionService(mockUsecase, logger)
	ctx := context.Background()

	t.Run("valid permission rule request", func(t *testing.T) {
		req := &CreatePermissionRuleRequest{
			RoleID:          1,
			DocType:         "User",
			PermissionLevel: 0,
			CanRead:         true,
			CanWrite:        true,
		}

		expectedRule := &biz.PermissionRule{
			ID:              1,
			RoleID:          req.RoleID,
			DocType:         req.DocType,
			PermissionLevel: req.PermissionLevel,
			CanRead:         req.CanRead,
			CanWrite:        req.CanWrite,
		}

		mockUsecase.On("CreatePermissionRule", ctx, mock.AnythingOfType("*biz.PermissionRule")).Return(expectedRule, nil)

		result, err := service.CreatePermissionRule(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.RoleID, result.RoleID)
		assert.Equal(t, req.DocType, result.DocType)
		assert.Equal(t, req.CanRead, result.CanRead)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("invalid role ID should fail", func(t *testing.T) {
		req := &CreatePermissionRuleRequest{
			RoleID:          0, // Invalid: should be positive
			DocType:         "User",
			PermissionLevel: 0,
		}

		result, err := service.CreatePermissionRule(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "role_id")
	})

	t.Run("empty doc type should fail", func(t *testing.T) {
		req := &CreatePermissionRuleRequest{
			RoleID:          1,
			DocType:         "", // Invalid: empty doc type
			PermissionLevel: 0,
		}

		result, err := service.CreatePermissionRule(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "doc_type")
	})

	t.Run("invalid permission level should fail", func(t *testing.T) {
		req := &CreatePermissionRuleRequest{
			RoleID:          1,
			DocType:         "User",
			PermissionLevel: 10, // Invalid: should be 0-9
		}

		result, err := service.CreatePermissionRule(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "permission_level")
	})
}