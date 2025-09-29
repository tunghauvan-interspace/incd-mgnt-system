<template>
  <aside 
    :class="[
      'fixed top-16 left-0 z-40 transition-all duration-300 ease-in-out bg-white border-r border-gray-200 shadow-sm',
      isCollapsed ? 'w-16' : 'w-72',
      isMobileOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
    ]"
    style="height: calc(100vh - 4rem)"
  >
    <div class="flex flex-col h-full">
      <!-- Sidebar Header -->
      <div class="flex items-center justify-between p-4 border-b border-gray-200">
        <h2 v-show="!isCollapsed" class="text-lg font-semibold text-gray-800 transition-opacity duration-200">
          Navigation
        </h2>
        <button 
          @click="toggleCollapsed"
          class="hidden lg:flex p-1.5 rounded-lg hover:bg-gray-100 text-gray-500 hover:text-gray-700"
        >
          <svg 
            :class="['w-5 h-5 transition-transform duration-200', isCollapsed ? 'rotate-180' : '']"
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
      </div>

      <!-- Navigation Menu -->
      <nav class="flex-1 p-4 space-y-2 overflow-y-auto">
        <template v-for="item in navigationItems" :key="item.name">
          <div v-if="item.divider" class="my-4 border-t border-gray-200"></div>
          
          <RouterLink
            v-else
            :to="item.href"
            :class="[
              'flex items-center rounded-lg text-sm font-medium transition-all duration-200 group',
              isCollapsed ? 'p-3 justify-center' : 'p-3',
              isActive(item.href) 
                ? 'bg-blue-50 text-blue-700 border-r-3 border-blue-700' 
                : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
            ]"
          >
            <!-- Icon -->
            <component 
              :is="item.icon" 
              :class="[
                'flex-shrink-0 transition-colors duration-200',
                isCollapsed ? 'w-6 h-6' : 'w-5 h-5 mr-3',
                isActive(item.href) ? 'text-blue-700' : 'text-gray-400 group-hover:text-gray-500'
              ]"
            />
            
            <!-- Text and Badge -->
            <div 
              v-show="!isCollapsed" 
              class="flex items-center justify-between w-full transition-opacity duration-200"
            >
              <span>{{ item.name }}</span>
              <span 
                v-if="item.badge" 
                :class="[
                  'inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium',
                  item.badgeType === 'error' ? 'bg-red-100 text-red-800' :
                  item.badgeType === 'warning' ? 'bg-yellow-100 text-yellow-800' :
                  'bg-gray-100 text-gray-700'
                ]"
              >
                {{ item.badge }}
              </span>
            </div>

            <!-- Tooltip for collapsed state -->
            <div 
              v-if="isCollapsed"
              class="absolute left-16 px-2 py-1 ml-2 text-xs font-medium text-white bg-gray-900 rounded-md opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none z-50 whitespace-nowrap"
            >
              {{ item.name }}
              <span v-if="item.badge" class="ml-1">({{ item.badge }})</span>
            </div>
          </RouterLink>
        </template>
      </nav>

      <!-- Sidebar Footer -->
      <div class="p-4 border-t border-gray-200">
        <div 
          :class="[
            'flex items-center text-xs text-gray-500',
            isCollapsed ? 'justify-center' : 'space-x-2'
          ]"
        >
          <div class="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
          <span v-show="!isCollapsed">System Online</span>
        </div>
      </div>
    </div>

    <!-- Mobile overlay -->
    <div 
      v-if="isMobileOpen" 
      @click="closeMobile"
      class="fixed inset-0 z-30 bg-gray-900 bg-opacity-50 lg:hidden"
    ></div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useRoute } from 'vue-router'
import { RouterLink } from 'vue-router'

// Props
interface Props {
  isCollapsed?: boolean
  isMobileOpen?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isCollapsed: false,
  isMobileOpen: false
})

// Emits
const emit = defineEmits(['toggle-collapsed', 'close-mobile'])

const route = useRoute()

// Icon components as functions that return SVG elements
const DashboardIcon = () => h('svg', {
  fill: 'none',
  stroke: 'currentColor',
  viewBox: '0 0 24 24',
  class: 'w-5 h-5'
}, [
  h('path', {
    'stroke-linecap': 'round',
    'stroke-linejoin': 'round',
    'stroke-width': '2',
    d: 'M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z'
  })
])

const IncidentsIcon = () => h('svg', {
  fill: 'none',
  stroke: 'currentColor',
  viewBox: '0 0 24 24',
  class: 'w-5 h-5'
}, [
  h('path', {
    'stroke-linecap': 'round',
    'stroke-linejoin': 'round',
    'stroke-width': '2',
    d: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z'
  })
])

const AlertsIcon = () => h('svg', {
  fill: 'none',
  stroke: 'currentColor',
  viewBox: '0 0 24 24',
  class: 'w-5 h-5'
}, [
  h('path', {
    'stroke-linecap': 'round',
    'stroke-linejoin': 'round',
    'stroke-width': '2',
    d: 'M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9'
  })
])

const UsersIcon = () => h('svg', {
  fill: 'none',
  stroke: 'currentColor',
  viewBox: '0 0 24 24',
  class: 'w-5 h-5'
}, [
  h('path', {
    'stroke-linecap': 'round',
    'stroke-linejoin': 'round',
    'stroke-width': '2',
    d: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z'
  })
])

const SettingsIcon = () => h('svg', {
  fill: 'none',
  stroke: 'currentColor',
  viewBox: '0 0 24 24',
  class: 'w-5 h-5'
}, [
  h('path', {
    'stroke-linecap': 'round',
    'stroke-linejoin': 'round',
    'stroke-width': '2',
    d: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z'
  }),
  h('path', {
    'stroke-linecap': 'round',
    'stroke-linejoin': 'round',
    'stroke-width': '2',
    d: 'M15 12a3 3 0 11-6 0 3 3 0 016 0z'
  })
])

// Navigation items
const navigationItems = [
  {
    name: 'Dashboard',
    href: '/',
    icon: DashboardIcon
  },
  {
    name: 'Incidents',
    href: '/incidents',
    icon: IncidentsIcon,
    badge: '12',
    badgeType: 'error'
  },
  {
    name: 'Alerts',
    href: '/alerts',
    icon: AlertsIcon,
    badge: '5',
    badgeType: 'warning'
  },
  { divider: true },
  {
    name: 'Users',
    href: '/users',
    icon: UsersIcon
  },
  {
    name: 'Settings',
    href: '/settings',
    icon: SettingsIcon
  }
]

// Methods
const toggleCollapsed = () => {
  emit('toggle-collapsed')
}

const closeMobile = () => {
  emit('close-mobile')
}

const isActive = (href: string) => {
  if (href === '/') {
    return route.path === '/'
  }
  return route.path.startsWith(href)
}
</script>

<style scoped>
.border-r-3 {
  border-right-width: 3px;
}
</style>