import React from 'react';
import { DownOutlined, UpOutlined } from '@ant-design/icons';
import { PermissionModule, PermissionItem } from '../types/auth';

interface PermissionModuleProps {
  module: PermissionModule;
  onToggle: () => void;
  onPermissionChange: (permissionId: string, checked: boolean) => void;
  onSelectAll: (select: boolean) => void;
}

const PermissionModuleComponent: React.FC<PermissionModuleProps> = ({
  module,
  onToggle,
  onPermissionChange,
  onSelectAll
}) => {
  const handlePermissionChange = (permission: PermissionItem, checked: boolean) => {
    onPermissionChange(permission.id, checked);
  };

  const handleSelectAll = () => {
    onSelectAll(true);
  };

  const handleClearAll = () => {
    onSelectAll(false);
  };

  return (
    <div className="user-edit-permission-module">
      <div className="user-edit-module-header" onClick={onToggle}>
        <div className="user-edit-module-info">
          <div className="user-edit-module-icon">
            {module.icon}
          </div>
          <div className="user-edit-module-details">
            <h4>{module.name}</h4>
            <p>{module.description}</p>
          </div>
        </div>
        <div className="user-edit-module-controls">
          <div className="user-edit-module-actions">
            <button 
              className="user-edit-action-btn select-all" 
              onClick={(e) => {
                e.stopPropagation();
                handleSelectAll();
              }}
              type="button"
            >
              全选
            </button>
            <button 
              className="user-edit-action-btn clear-all" 
              onClick={(e) => {
                e.stopPropagation();
                handleClearAll();
              }}
              type="button"
            >
              取消
            </button>
          </div>
          <div className={`user-edit-expand-arrow ${module.expanded ? 'expanded' : ''}`}>
            {module.expanded ? <UpOutlined /> : <DownOutlined />}
          </div>
        </div>
      </div>
      
      <div className={`user-edit-module-permissions ${module.expanded ? 'expanded' : ''}`}>
        <div className="user-edit-permission-grid">
          {module.permissions.map(permission => (
            <div key={permission.id} className="user-edit-permission-item">
              <input
                type="checkbox"
                className="user-edit-permission-checkbox"
                id={permission.id}
                checked={permission.checked || false}
                onChange={(e) => handlePermissionChange(permission, e.target.checked)}
              />
              <label 
                htmlFor={permission.id}
                className="user-edit-permission-label"
              >
                {permission.name}
              </label>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default PermissionModuleComponent;