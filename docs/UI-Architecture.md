# UI Architecture - Incident Management System

This document outlines the planned user interface architecture, component structure, and design system for the incident management system frontend.

## Table of Contents
- [Technology Stack](#technology-stack)
- [Application Structure](#application-structure)
- [Component Architecture](#component-architecture)
- [Design System](#design-system)
- [Routing Structure](#routing-structure)
- [State Management](#state-management)
- [API Integration](#api-integration)
- [Performance Optimization](#performance-optimization)
- [Mobile & Responsive Design](#mobile--responsive-design)
- [Development Workflow](#development-workflow)

---

## Technology Stack

### Core Framework
- **Frontend Framework**: Vue.js 3 with Composition API
- **Build Tool**: Vite for fast development and optimized builds
- **Language**: TypeScript for type safety
- **Styling**: CSS3 with CSS Variables + Utility Classes
- **State Management**: Pinia for reactive state management
- **Routing**: Vue Router for SPA navigation

### UI Libraries & Components
- **Charts**: Chart.js for metrics visualization
- **Icons**: Heroicons or Tabler Icons
- **Date/Time**: Day.js for lightweight date manipulation
- **HTTP Client**: Axios with interceptors
- **Notifications**: Custom notification system

### Development Tools
- **Testing**: Vitest for unit tests, Playwright for E2E
- **Linting**: ESLint + Prettier for code quality
- **Type Checking**: Vue TypeScript support
- **Hot Reload**: Vite HMR for development

---

## Application Structure

### Project Structure (Planned)
```
web/frontend/
├── public/                     # Static assets
│   ├── favicon.ico
│   ├── manifest.json          # PWA manifest
│   └── icons/                 # App icons
├── src/
│   ├── assets/                # Build-time assets
│   │   ├── images/
│   │   ├── icons/
│   │   └── styles/
│   ├── components/            # Reusable components
│   │   ├── ui/               # Base UI components
│   │   ├── forms/            # Form-specific components
│   │   ├── charts/           # Chart components
│   │   └── layout/           # Layout components
│   ├── composables/          # Vue composables (hooks)
│   │   ├── useApi.ts
│   │   ├── useAuth.ts
│   │   ├── useIncidents.ts
│   │   └── useNotifications.ts
│   ├── layouts/              # Page layouts
│   │   ├── DefaultLayout.vue
│   │   ├── AuthLayout.vue
│   │   └── MobileLayout.vue
│   ├── pages/                # Route components
│   │   ├── Dashboard.vue
│   │   ├── Incidents/
│   │   ├── Alerts/
│   │   ├── Users/
│   │   └── Settings/
│   ├── router/               # Routing configuration
│   │   └── index.ts
│   ├── stores/               # Pinia stores
│   │   ├── auth.ts
│   │   ├── incidents.ts
│   │   ├── alerts.ts
│   │   └── notifications.ts
│   ├── types/                # TypeScript type definitions
│   │   ├── api.ts
│   │   ├── incidents.ts
│   │   └── users.ts
│   ├── utils/                # Utility functions
│   │   ├── formatters.ts
│   │   ├── validators.ts
│   │   └── constants.ts
│   ├── App.vue               # Root component
│   └── main.ts               # Application entry point
├── tests/                    # Test files
│   ├── unit/
│   └── e2e/
├── index.html                # HTML entry point
├── vite.config.ts           # Vite configuration
└── package.json             # Dependencies
```

---

## Component Architecture

### Component Hierarchy

#### 1. Base Components (ui/)
**Atomic, reusable components**

```typescript
// Button Component
interface ButtonProps {
  variant: 'primary' | 'secondary' | 'danger' | 'success'
  size: 'sm' | 'md' | 'lg'
  disabled?: boolean
  loading?: boolean
  icon?: string
}

// Badge Component
interface BadgeProps {
  color: 'red' | 'yellow' | 'green' | 'blue' | 'gray'
  size: 'sm' | 'md' | 'lg'
  variant: 'solid' | 'outline' | 'soft'
}

// Modal Component
interface ModalProps {
  modelValue: boolean
  title: string
  size: 'sm' | 'md' | 'lg' | 'xl'
  persistent?: boolean
}
```

#### 2. Form Components (forms/)
**Form-specific reusable components**

```typescript
// Input Component
interface InputProps {
  modelValue: string | number
  type: 'text' | 'email' | 'password' | 'number'
  placeholder?: string
  error?: string
  disabled?: boolean
  required?: boolean
}

// Select Component
interface SelectProps {
  modelValue: any
  options: Array<{label: string, value: any}>
  placeholder?: string
  multiple?: boolean
  searchable?: boolean
}

// DatePicker Component
interface DatePickerProps {
  modelValue: Date | null
  range?: boolean
  format?: string
  minDate?: Date
  maxDate?: Date
}
```

#### 3. Business Components
**Domain-specific components**

```typescript
// IncidentCard Component
interface IncidentCardProps {
  incident: Incident
  compact?: boolean
  selectable?: boolean
  onClick?: (incident: Incident) => void
}

// AlertsList Component
interface AlertsListProps {
  alerts: Alert[]
  groupBy?: 'service' | 'severity' | 'status'
  onAlertClick?: (alert: Alert) => void
}

// StatusBadge Component
interface StatusBadgeProps {
  status: IncidentStatus
  showIcon?: boolean
  interactive?: boolean
}

// SeverityBadge Component
interface SeverityBadgeProps {
  severity: IncidentSeverity
  showIcon?: boolean
  size?: 'sm' | 'md' | 'lg'
}
```

#### 4. Layout Components
**Page structure and navigation**

```typescript
// Navbar Component
interface NavbarProps {
  user?: User
  notifications?: Notification[]
  onLogout?: () => void
}

// Sidebar Component
interface SidebarProps {
  collapsed?: boolean
  activeRoute?: string
  menuItems: MenuItem[]
}

// PageHeader Component
interface PageHeaderProps {
  title: string
  subtitle?: string
  breadcrumbs?: Breadcrumb[]
  actions?: Action[]
}
```

---

## Design System

### Color Palette
```css
:root {
  /* Primary Colors */
  --color-primary-50: #f0f9ff;
  --color-primary-500: #3b82f6;
  --color-primary-900: #1e3a8a;
  
  /* Status Colors */
  --color-success: #10b981;
  --color-warning: #f59e0b;
  --color-error: #ef4444;
  --color-info: #6366f1;
  
  /* Severity Colors */
  --color-critical: #dc2626;
  --color-high: #ea580c;
  --color-medium: #d97706;
  --color-low: #65a30d;
  
  /* Neutral Colors */
  --color-gray-50: #f9fafb;
  --color-gray-100: #f3f4f6;
  --color-gray-500: #6b7280;
  --color-gray-900: #111827;
}
```

### Typography Scale
```css
:root {
  /* Font Families */
  --font-sans: 'Inter', system-ui, sans-serif;
  --font-mono: 'JetBrains Mono', Consolas, monospace;
  
  /* Font Sizes */
  --text-xs: 0.75rem;    /* 12px */
  --text-sm: 0.875rem;   /* 14px */
  --text-base: 1rem;     /* 16px */
  --text-lg: 1.125rem;   /* 18px */
  --text-xl: 1.25rem;    /* 20px */
  --text-2xl: 1.5rem;    /* 24px */
  --text-3xl: 1.875rem;  /* 30px */
  
  /* Line Heights */
  --leading-tight: 1.25;
  --leading-normal: 1.5;
  --leading-relaxed: 1.75;
}
```

### Spacing System
```css
:root {
  --space-1: 0.25rem;   /* 4px */
  --space-2: 0.5rem;    /* 8px */
  --space-3: 0.75rem;   /* 12px */
  --space-4: 1rem;      /* 16px */
  --space-6: 1.5rem;    /* 24px */
  --space-8: 2rem;      /* 32px */
  --space-12: 3rem;     /* 48px */
  --space-16: 4rem;     /* 64px */
}
```

### Component Tokens
```css
:root {
  /* Button */
  --button-height-sm: 2rem;
  --button-height-md: 2.5rem;
  --button-height-lg: 3rem;
  --button-border-radius: 0.375rem;
  
  /* Input */
  --input-height: 2.5rem;
  --input-border-radius: 0.375rem;
  --input-border-width: 1px;
  
  /* Card */
  --card-border-radius: 0.5rem;
  --card-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  --card-padding: 1.5rem;
}
```

### Dark Mode Support
```css
[data-theme="dark"] {
  --color-bg-primary: #0f172a;
  --color-bg-secondary: #1e293b;
  --color-text-primary: #f1f5f9;
  --color-text-secondary: #cbd5e1;
  --color-border: #334155;
}
```

---

## Routing Structure

### Route Definitions
```typescript
// router/index.ts
const routes = [
  {
    path: '/',
    component: () => import('../layouts/DefaultLayout.vue'),
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('../pages/Dashboard.vue'),
        meta: { requiresAuth: true, title: 'Dashboard' }
      },
      {
        path: '/incidents',
        name: 'Incidents',
        component: () => import('../pages/Incidents/Index.vue'),
        meta: { requiresAuth: true, title: 'Incidents' }
      },
      {
        path: '/incidents/:id',
        name: 'IncidentDetail',
        component: () => import('../pages/Incidents/Detail.vue'),
        meta: { requiresAuth: true, title: 'Incident Details' }
      },
      {
        path: '/alerts',
        name: 'Alerts',
        component: () => import('../pages/Alerts/Index.vue'),
        meta: { requiresAuth: true, title: 'Alerts' }
      },
      {
        path: '/users',
        name: 'Users',
        component: () => import('../pages/Users/Index.vue'),
        meta: { requiresAuth: true, permission: 'users.read' }
      },
      {
        path: '/settings',
        name: 'Settings',
        component: () => import('../pages/Settings/Index.vue'),
        meta: { requiresAuth: true, title: 'Settings' }
      }
    ]
  },
  {
    path: '/auth',
    component: () => import('../layouts/AuthLayout.vue'),
    children: [
      {
        path: 'login',
        name: 'Login',
        component: () => import('../pages/Auth/Login.vue'),
        meta: { title: 'Login' }
      },
      {
        path: 'register',
        name: 'Register',
        component: () => import('../pages/Auth/Register.vue'),
        meta: { title: 'Register' }
      }
    ]
  }
]
```

### Navigation Guards
```typescript
// Authentication Guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/auth/login')
    return
  }
  
  if (to.meta.permission && !authStore.hasPermission(to.meta.permission)) {
    next('/unauthorized')
    return
  }
  
  next()
})
```

---

## State Management

### Store Structure
```typescript
// stores/auth.ts
export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const permissions = ref<string[]>([])
  
  const isAuthenticated = computed(() => !!token.value)
  
  const login = async (credentials: LoginCredentials) => {
    const response = await authApi.login(credentials)
    token.value = response.token
    user.value = response.user
    permissions.value = response.permissions
  }
  
  const logout = () => {
    user.value = null
    token.value = null
    permissions.value = []
  }
  
  const hasPermission = (permission: string) => {
    return permissions.value.includes(permission)
  }
  
  return {
    user,
    token,
    permissions,
    isAuthenticated,
    login,
    logout,
    hasPermission
  }
})

// stores/incidents.ts
export const useIncidentsStore = defineStore('incidents', () => {
  const incidents = ref<Incident[]>([])
  const currentIncident = ref<Incident | null>(null)
  const loading = ref(false)
  const filters = ref<IncidentFilters>({})
  
  const filteredIncidents = computed(() => {
    return incidents.value.filter(incident => {
      if (filters.value.status && incident.status !== filters.value.status) {
        return false
      }
      if (filters.value.severity && incident.severity !== filters.value.severity) {
        return false
      }
      return true
    })
  })
  
  const fetchIncidents = async () => {
    loading.value = true
    try {
      incidents.value = await incidentsApi.list(filters.value)
    } finally {
      loading.value = false
    }
  }
  
  return {
    incidents,
    currentIncident,
    loading,
    filters,
    filteredIncidents,
    fetchIncidents
  }
})
```

---

## API Integration

### API Client Structure
```typescript
// composables/useApi.ts
export const useApi = () => {
  const authStore = useAuthStore()
  
  const client = axios.create({
    baseURL: '/api',
    headers: {
      'Content-Type': 'application/json'
    }
  })
  
  // Request interceptor for auth
  client.interceptors.request.use(config => {
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  })
  
  // Response interceptor for error handling
  client.interceptors.response.use(
    response => response,
    error => {
      if (error.response?.status === 401) {
        authStore.logout()
        router.push('/auth/login')
      }
      return Promise.reject(error)
    }
  )
  
  return { client }
}

// composables/useIncidents.ts
export const useIncidents = () => {
  const { client } = useApi()
  const store = useIncidentsStore()
  
  const fetchIncidents = async (filters?: IncidentFilters) => {
    const response = await client.get('/incidents', { params: filters })
    return response.data
  }
  
  const getIncident = async (id: string) => {
    const response = await client.get(`/incidents/${id}`)
    return response.data
  }
  
  const createIncident = async (incident: CreateIncidentRequest) => {
    const response = await client.post('/incidents', incident)
    return response.data
  }
  
  const acknowledgeIncident = async (id: string, assigneeId?: string) => {
    const response = await client.put(`/incidents/${id}/acknowledge`, {
      assignee_id: assigneeId
    })
    return response.data
  }
  
  return {
    fetchIncidents,
    getIncident,
    createIncident,
    acknowledgeIncident
  }
}
```

### Real-time Updates
```typescript
// composables/useRealtime.ts
export const useRealtime = () => {
  const socket = ref<WebSocket | null>(null)
  const connected = ref(false)
  
  const connect = () => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws`
    
    socket.value = new WebSocket(wsUrl)
    
    socket.value.onopen = () => {
      connected.value = true
    }
    
    socket.value.onmessage = (event) => {
      const data = JSON.parse(event.data)
      handleRealtimeUpdate(data)
    }
    
    socket.value.onclose = () => {
      connected.value = false
      // Reconnect logic
      setTimeout(connect, 3000)
    }
  }
  
  const handleRealtimeUpdate = (data: any) => {
    const incidentsStore = useIncidentsStore()
    
    switch (data.type) {
      case 'incident_created':
        incidentsStore.incidents.unshift(data.incident)
        break
      case 'incident_updated':
        const index = incidentsStore.incidents.findIndex(i => i.id === data.incident.id)
        if (index >= 0) {
          incidentsStore.incidents[index] = data.incident
        }
        break
    }
  }
  
  return { connect, connected }
}
```

---

## Performance Optimization

### Code Splitting
```typescript
// Lazy loading routes
const routes = [
  {
    path: '/incidents',
    component: () => import('../pages/Incidents/Index.vue') // Lazy loaded
  }
]

// Dynamic imports for large components
const ChartComponent = defineAsyncComponent(() => import('./Chart.vue'))
```

### Virtual Scrolling
```vue
<!-- For large lists -->
<template>
  <VirtualList
    :items="incidents"
    :item-height="80"
    class="incident-list"
  >
    <template #default="{ item }">
      <IncidentCard :incident="item" />
    </template>
  </VirtualList>
</template>
```

### Caching Strategy
```typescript
// API response caching
const cache = new Map()

export const useCache = () => {
  const get = (key: string) => cache.get(key)
  const set = (key: string, value: any, ttl = 60000) => {
    cache.set(key, { value, expires: Date.now() + ttl })
  }
  
  const isValid = (key: string) => {
    const item = cache.get(key)
    return item && Date.now() < item.expires
  }
  
  return { get, set, isValid }
}
```

---

## Mobile & Responsive Design

### Breakpoint System
```css
:root {
  --breakpoint-sm: 640px;
  --breakpoint-md: 768px;
  --breakpoint-lg: 1024px;
  --breakpoint-xl: 1280px;
}
```

### Mobile-First Approach
```vue
<template>
  <div class="incident-card">
    <!-- Mobile layout -->
    <div class="md:hidden">
      <MobileIncidentCard :incident="incident" />
    </div>
    
    <!-- Desktop layout -->
    <div class="hidden md:block">
      <DesktopIncidentCard :incident="incident" />
    </div>
  </div>
</template>

<style>
/* Mobile-first CSS */
.incident-card {
  padding: 1rem;
}

@media (min-width: 768px) {
  .incident-card {
    padding: 1.5rem;
  }
}
</style>
```

### Touch Interactions
```vue
<template>
  <div
    class="incident-row"
    @touchstart="handleTouchStart"
    @touchend="handleTouchEnd"
  >
    <SwipeAction
      @swipe-left="acknowledgeIncident"
      @swipe-right="assignIncident"
    >
      <IncidentCard :incident="incident" />
    </SwipeAction>
  </div>
</template>
```

---

## Development Workflow

### Component Development
```bash
# Create new component
npm run generate:component ComponentName

# Run development server
npm run dev

# Run tests
npm run test
npm run test:e2e

# Build for production
npm run build
```

### Code Quality
```json
// .eslintrc.js
{
  "extends": [
    "@vue/typescript/recommended",
    "prettier"
  ],
  "rules": {
    "vue/component-name-in-template-casing": ["error", "PascalCase"],
    "vue/no-unused-vars": "error",
    "@typescript-eslint/no-unused-vars": "error"
  }
}
```

### Testing Strategy
```typescript
// Component test example
describe('IncidentCard', () => {
  it('displays incident information correctly', () => {
    const incident = {
      id: '1',
      title: 'Test Incident',
      status: 'open',
      severity: 'high'
    }
    
    const wrapper = mount(IncidentCard, {
      props: { incident }
    })
    
    expect(wrapper.text()).toContain('Test Incident')
    expect(wrapper.find('.severity-badge').classes()).toContain('severity-high')
  })
})
```

### Build Configuration
```typescript
// vite.config.ts
export default defineConfig({
  plugins: [vue()],
  build: {
    outDir: '../static',
    emptyOutDir: true,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['vue', 'vue-router', 'pinia'],
          charts: ['chart.js']
        }
      }
    }
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8080'
    }
  }
})
```

---

This UI Architecture document provides the foundation for building a scalable, maintainable, and user-friendly frontend for the incident management system. The architecture emphasizes component reusability, type safety, performance, and responsive design to deliver an excellent user experience across all devices.