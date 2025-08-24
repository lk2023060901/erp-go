import React, { useState, useRef } from 'react';
import { message } from 'antd';
import { AvatarUploadProps } from '../types/auth';

const AvatarUpload: React.FC<AvatarUploadProps> = ({
  value,
  onChange,
  disabled = false
}) => {
  const [preview, setPreview] = useState<string | null>(value || null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // 处理文件选择
  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    // 验证文件类型
    const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif'];
    if (!allowedTypes.includes(file.type)) {
      message.error('请选择 JPG、PNG 或 GIF 格式的图片');
      return;
    }

    // 验证文件大小（最大 2MB）
    const maxSize = 2 * 1024 * 1024;
    if (file.size > maxSize) {
      message.error('图片大小不能超过 2MB');
      return;
    }

    // 创建预览URL
    const reader = new FileReader();
    reader.onload = (e) => {
      const result = e.target?.result as string;
      setPreview(result);
    };
    reader.readAsDataURL(file);

    // 通知父组件
    onChange?.(file);

    // 清空input值，允许选择同一个文件
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  // 处理点击上传
  const handleClick = () => {
    if (!disabled && fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  return (
    <div 
      className={`user-edit-avatar-upload ${disabled ? 'disabled' : ''}`}
      onClick={handleClick}
      style={{ cursor: disabled ? 'not-allowed' : 'pointer' }}
    >
      <div className="user-edit-avatar-image">
        {preview ? (
          <img 
            src={preview} 
            alt="用户头像" 
            style={{ 
              width: '100%', 
              height: '100%', 
              objectFit: 'cover' 
            }} 
          />
        ) : (
          <span className="user-edit-avatar-placeholder">👤</span>
        )}
      </div>
      
      {!disabled && (
        <div className="user-edit-avatar-overlay">
          <div className="user-edit-avatar-overlay-text">
            点击上传<br />头像
          </div>
        </div>
      )}

      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        style={{ display: 'none' }}
        onChange={handleFileSelect}
        disabled={disabled}
      />
    </div>
  );
};

export default AvatarUpload;