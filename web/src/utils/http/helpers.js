/**********************************
 * @FilePath: helpers.js
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/04 22:46:22
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { useAuthStore } from '@/store'

let isConfirming = false
export function resolveResError(code, message, needTip = true) {
  switch (code) {
    case 401:
      if (isConfirming || !needTip)
        return
      isConfirming = true
      $dialog.confirm({
        title: 'Prompt',
        type: 'info',
        content: 'Login has expired, would you like to log in again?',
        confirm() {
          useAuthStore().logout()
          window.$message?.success('Logged out')
          isConfirming = false
        },
        cancel() {
          isConfirming = false
        },
      })
      return false
    case 11007:
    case 11008:
      if (isConfirming || !needTip)
        return
      isConfirming = true
      $dialog.confirm({
        title: 'Prompt',
        type: 'info',
        content: `${message}, would you like to log in again?`,
        confirm() {
          useAuthStore().logout()
          window.$message?.success('Logged out')
          isConfirming = false
        },
        cancel() {
          isConfirming = false
        },
      })
      return false
    case 403:
      message = 'Request was denied'
      break
    case 404:
      message = 'Requested resource or endpoint does not exist'
      break
    case 500:
      message = 'Server error occurred'
      break
    default:
      message = message ?? `【${code}】: Unknown error!`
      break
  }
  needTip && window.$message?.error(message)
  return message
}
