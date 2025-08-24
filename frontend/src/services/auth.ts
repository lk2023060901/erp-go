// 重新导出认证服务
export * from './authService';

// 创建一个authService对象以便向后兼容
import * as authFunctions from './authService';
export const authService = authFunctions;