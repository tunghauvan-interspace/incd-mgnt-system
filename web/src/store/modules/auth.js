/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:25:39
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { defineStore } from 'pinia'
import { usePermissionStore, useRouterStore, useTabStore, useUserStore } from '@/store'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    accessToken: undefined,
  }),
  actions: {
    setToken({ accessToken }) {
      this.accessToken = accessToken
    },
    resetToken() {
      this.$reset()
    },
    toLogin() {
      const { router, route } = useRouterStore()
      router.replace({
        path: '/login',
        query: route.query,
      })
    },
    async switchCurrentRole(data) {
      this.resetLoginState()
      await nextTick()
      this.setToken(data)
    },
    resetLoginState() {
      const { resetUser } = useUserStore()
      const { resetRouter } = useRouterStore()
      const { resetPermission, accessRoutes } = usePermissionStore()
      const { resetTabs } = useTabStore()
  // Reset router
  resetRouter(accessRoutes)
  // Reset user
  resetUser()
  // Reset permissions
  resetPermission()
  // Reset tabs
  resetTabs()
  // Reset token
  this.resetToken()
    },
    async logout() {
      this.resetLoginState()
      this.toLogin()
    },
  },
  persist: {
    key: 'vue-naivue-admin_auth',
  },
})
