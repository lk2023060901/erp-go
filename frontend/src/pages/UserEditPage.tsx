import React from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import { App } from 'antd';
import UserEditForm from '../components/UserEditForm';
import { createUserFromForm, updateUserFromForm, getUserById, uploadUserAvatar } from '../services/userService';
import { FormMode, UserFormData, User } from '../types/auth';

const UserEditPage: React.FC = () => {
  const { id } = useParams<{ id?: string }>();
  const navigate = useNavigate();
  const location = useLocation();
  const { message } = App.useApp();
  
  // 从URL参数中获取初始选项卡
  const urlParams = new URLSearchParams(location.search);
  const initialTab = urlParams.get('tab');
  
  // 根据路由路径确定模式
  const getMode = (): FormMode => {
    if (location.pathname.includes('/create')) {
      return 'create';
    } else if (location.pathname.includes('/profile')) {
      return 'profile';
    } else {
      return 'edit';
    }
  };
  
  const mode = getMode();
  
  const [user, setUser] = React.useState<User | null>(null);
  const [loading, setLoading] = React.useState(false);

  // 加载用户数据（编辑模式）
  React.useEffect(() => {
    if (mode !== 'create' && id) {
      loadUser();
    }
  }, [id, mode]);

  const loadUser = async () => {
    if (!id) return;
    
    try {
      setLoading(true);
      const userData = await getUserById(parseInt(id));
      setUser(userData);
    } catch (error: any) {
      message.error(error.message || '加载用户信息失败');
      navigate('/users');
    } finally {
      setLoading(false);
    }
  };

  // 处理表单提交
  const handleSubmit = async (formData: UserFormData) => {
    try {
      let resultUser: User;
      
      if (mode === 'create') {
        // 创建用户
        resultUser = await createUserFromForm(formData);
        message.success('创建用户成功');
      } else if (id) {
        // 更新用户
        resultUser = await updateUserFromForm(parseInt(id), formData);
        message.success('更新用户信息成功');
      } else {
        throw new Error('无效的操作模式');
      }

      // 如果有头像文件，上传头像
      if (formData.avatar && resultUser.id) {
        try {
          const avatarUrl = await uploadUserAvatar(resultUser.id, formData.avatar);
          console.log('头像上传成功:', avatarUrl);
        } catch (avatarError) {
          console.warn('头像上传失败:', avatarError);
          message.warning('用户信息保存成功，但头像上传失败');
        }
      }

      // 返回用户列表页面
      navigate('/users');
    } catch (error: any) {
      throw error; // 让UserEditForm处理错误显示
    }
  };

  // 处理取消
  const handleCancel = () => {
    navigate('/users');
  };

  if (loading) {
    return <div>加载中...</div>;
  }

  return (
    <UserEditForm
      user={user}
      mode={mode || 'create'}
      onSubmit={handleSubmit}
      onCancel={handleCancel}
      showActionButtons={true}
      initialTab={initialTab || undefined}
    />
  );
};

export default UserEditPage;