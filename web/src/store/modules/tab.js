/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:25:52
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { defineStore } from 'pinia'
import { useRouterStore } from './router'

export const useTabStore = defineStore('tab', {
  state: () => ({
    tabs: [],
    activeTab: '',
    reloading: false,
  }),
  getters: {
    activeIndex() {
      return this.tabs.findIndex(item => item.path === this.activeTab)
    },
  },
  actions: {
    async setActiveTab(path) {
  await nextTick() // Wait for tab DOM update before setting active so the tab bar scrolls to the new tab
      this.activeTab = path
    },
    setTabs(tabs) {
      this.tabs = tabs
    },
    addTab(tab = {}) {
      const findIndex = this.tabs.findIndex(item => item.path === tab.path)
      if (findIndex !== -1) {
        this.tabs.splice(findIndex, 1, tab)
      }
      else {
        this.setTabs([...this.tabs, tab])
      }
      this.setActiveTab(tab.path)
    },
    async reloadTab(path, keepAlive) {
      const findItem = this.tabs.find(item => item.path === path)
      if (!findItem)
        return
  // Updating the key can invalidate keepAlive
      if (keepAlive)
        findItem.keepAlive = false
      $loadingBar.start()
      this.reloading = true
      await nextTick()
      this.reloading = false
      findItem.keepAlive = !!keepAlive
      setTimeout(() => {
        document.documentElement.scrollTo({ left: 0, top: 0 })
        $loadingBar.finish()
      }, 100)
    },
    async removeTab(path) {
      this.setTabs(this.tabs.filter(tab => tab.path !== path))
      if (path === this.activeTab) {
        useRouterStore().router?.push(this.tabs[this.tabs.length - 1].path)
      }
    },
    removeOther(curPath = this.activeTab) {
      this.setTabs(this.tabs.filter(tab => tab.path === curPath))
      if (curPath !== this.activeTab) {
        useRouterStore().router?.push(this.tabs[this.tabs.length - 1].path)
      }
    },
    removeLeft(curPath) {
      const curIndex = this.tabs.findIndex(item => item.path === curPath)
      const filterTabs = this.tabs.filter((item, index) => index >= curIndex)
      this.setTabs(filterTabs)
      if (!filterTabs.find(item => item.path === this.activeTab)) {
        useRouterStore().router?.push(filterTabs[filterTabs.length - 1].path)
      }
    },
    removeRight(curPath) {
      const curIndex = this.tabs.findIndex(item => item.path === curPath)
      const filterTabs = this.tabs.filter((item, index) => index <= curIndex)
      this.setTabs(filterTabs)
      if (!filterTabs.find(item => item.path === this.activeTab.value)) {
        useRouterStore().router?.push(filterTabs[filterTabs.length - 1].path)
      }
    },
    resetTabs() {
      this.$reset()
    },
  },
  persist: {
    pick: ['tabs'],
    storage: sessionStorage,
  },
})
