/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/04 22:50:38
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { request } from '@/utils'

export default {
  // Get user info
  getUser: () => request.get('/user/detail'),
  // Refresh token
  refreshToken: () => request.get('/auth/refresh/token'),
  // Logout
  logout: () => request.post('/auth/logout', {}, { needTip: false }),
  // Switch current role
  switchCurrentRole: role => request.post(`/auth/current-role/switch/${role}`),
  // Get role permissions
  getRolePermissions: () => request.get('/role/permissions/tree'),
  // Validate menu path
  validateMenuPath: path => request.get(`/permission/menu/validate?path=${path}`),
}
