# UI Architecture - Incident Management System

This document outlines the planned user interface architecture, component structure, and design system for the incident management system frontend.

## Table of Contents
- [Technology Stack](#technology-stack)
- [Application Structure](#application-structure)
- [Component Architecture](#component-architecture)
- [Design System](#design-system)
- [Layout Design](#layout-design)
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
web/
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

## Layout Design

### Desktop Layout Architecture

#### Overall Application Layout (1440px+ desktop)
```
┌─────────────────────────────────────────────────────────────────────────┐
│ Header Navigation (Fixed)                                    Height: 64px │
├─────────────────────────────────────────────────────────────────────────┤
│ Sidebar │                Main Content Area                               │
│ 280px   │                                                               │
│ (Fixed) │  ┌─────────────────────────────────────────────────────────┐  │
│         │  │ Page Header                               Height: 80px   │  │
│         │  ├─────────────────────────────────────────────────────────┤  │
│         │  │                                                         │  │
│         │  │ Content Area (Scrollable)                              │  │
│         │  │                                                         │  │
│         │  │                                                         │  │
│         │  └─────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Header Navigation Component
**Dimensions & Behavior:**
- **Height**: 64px (fixed)
- **Width**: 100vw (full viewport)
- **Position**: Fixed top, z-index: 1000
- **Background**: White/Dark theme adaptive
- **Shadow**: 0 1px 3px rgba(0, 0, 0, 0.1)

**Internal Layout (Left to Right):**
```
┌─────┬─────────────────┬───────────────────────────┬─────────┬─────────┐
│Logo │ Breadcrumbs     │        Search Bar         │Alerts   │Profile  │
│120px│ Flex-1          │         320px             │ 40px    │ 200px   │
└─────┴─────────────────┴───────────────────────────┴─────────┴─────────┘
```

**Component Behaviors:**
- **Logo**: Clickable, navigates to dashboard
- **Breadcrumbs**: Auto-generated based on current route, max 4 levels
- **Search Bar**: Global search with autocomplete, 300ms debounce
- **Alerts**: Badge counter, dropdown on click showing recent notifications
- **Profile**: Dropdown menu with user info, settings, logout

#### Sidebar Navigation Component
**Dimensions & Behavior:**
- **Width**: 280px (expanded), 72px (collapsed)
- **Height**: calc(100vh - 64px)
- **Position**: Fixed left, below header
- **Transition**: 0.3s ease-in-out for collapse/expand

**Menu Structure:**
```
┌─────────────────────────────┐
│ Dashboard         [Icon]    │ 48px height
├─────────────────────────────┤
│ Incidents (12)    [Icon]    │ 48px height
├─────────────────────────────┤
│ Alerts (5)        [Icon]    │ 48px height
├─────────────────────────────┤
│ Reports           [Icon]    │ 48px height
├─────────────────────────────┤
│ Users             [Icon]    │ 48px height (if permission)
├─────────────────────────────┤
│ Settings          [Icon]    │ 48px height
└─────────────────────────────┘
```

**Collapsed State (72px wide):**
- Show only icons (24px size)
- Tooltip on hover showing menu text
- Active state indicator (3px left border)

#### Main Content Area
**Dimensions & Behavior:**
- **Width**: calc(100vw - 280px) when sidebar expanded
- **Width**: calc(100vw - 72px) when sidebar collapsed
- **Margin-left**: Adjusts based on sidebar state
- **Padding**: 24px on all sides
- **Background**: #f8fafc (light theme), #0f172a (dark theme)

### Page-Specific Layouts

#### Dashboard Layout
```
┌─────────────────────────────────────────────────────────────────┐
│ Page Header (Welcome back, John)                    Height: 80px │
├─────────────────────────────────────────────────────────────────┤
│ Stats Cards Row                                     Height: 120px│
│ ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐     │
│ │Open        │ │High Sev    │ │Avg Response│ │Active Users│     │
│ │Incidents   │ │Incidents   │ │Time        │ │           │     │
│ │    24      │ │     8      │ │   4.2min   │ │    12     │     │
│ └────────────┘ └────────────┘ └────────────┘ └────────────┘     │
├─────────────────────────────────────────────────────────────────┤
│ Charts Section                                      Height: 400px│
│ ┌─────────────────────────────────┐ ┌─────────────────────────┐   │
│ │ Incident Trends (Line Chart)    │ │ Severity Distribution   │   │
│ │          60% width              │ │      40% width          │   │
│ │                                 │ │    (Doughnut Chart)     │   │
│ └─────────────────────────────────┘ └─────────────────────────┘   │
├─────────────────────────────────────────────────────────────────┤
│ Recent Incidents Table                              Min-height: 400px│
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

**Stats Cards Behavior:**
- **Size**: 280px width × 120px height
- **Gap**: 24px between cards
- **Hover**: Subtle lift effect (translateY: -2px, shadow increase)
- **Click**: Navigate to detailed view
- **Animation**: Count-up animation on load (1.5s duration)

#### Incidents List Layout
```
┌─────────────────────────────────────────────────────────────────┐
│ Page Header + Actions                               Height: 80px │
│ ┌─────────────────────────────┐ ┌─────────────────────────────┐ │
│ │ Incidents (156)             │ │ [Create] [Export] [Filter]  │ │
│ └─────────────────────────────┘ └─────────────────────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│ Filters Bar (Collapsible)                          Height: 56px │
│ [Status ▼] [Severity ▼] [Assignee ▼] [Date Range] [Clear All]   │
├─────────────────────────────────────────────────────────────────┤
│ Incidents Table/Cards                               Flex: 1     │
│                                                                 │
│ ┌─────────────────────────────────────────────────────────────┐ │
│ │ [□] ID-001  High    Server Down      John   2h ago    Open │ │
│ ├─────────────────────────────────────────────────────────────┤ │
│ │ [□] ID-002  Medium  DB Connection    Jane   4h ago    Open │ │
│ ├─────────────────────────────────────────────────────────────┤ │
│ │ ... more rows ...                                           │ │
│ └─────────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│ Pagination                                          Height: 60px │
│                              [1] [2] [3] ... [10] Next         │
└─────────────────────────────────────────────────────────────────┘
```

**Table Row Behavior:**
- **Height**: 64px per row
- **Hover**: Background color change (#f1f5f9)
- **Selection**: Checkbox with bulk actions
- **Click**: Navigate to incident detail
- **Status Badge**: Color-coded, 8px height indicator
- **Severity Badge**: Icon + text, hover shows description

#### Incident Detail Layout
```
┌─────────────────────────────────────────────────────────────────┐
│ Header with Actions                                 Height: 100px│
│ ┌─────────────────────┐ ┌─────────────────────────────────────┐ │
│ │← Back | ID-001      │ │[Acknowledge] [Assign] [Close] [...] │ │
│ │High Severity        │ │                                     │ │
│ └─────────────────────┘ └─────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│ Main Content (2-Column)                                         │
│ ┌─────────────────────────────────────┐ ┌─────────────────────┐ │
│ │ Left Panel (Details)                │ │ Right Panel (Info)  │ │
│ │              70% width              │ │      30% width      │ │
│ │                                     │ │                     │ │
│ │ ┌─────────────────────────────────┐ │ │ ┌─────────────────┐ │ │
│ │ │ Description                     │ │ │ │ Status: Open    │ │ │
│ │ │ Timeline                        │ │ │ │ Created: 2h ago │ │ │
│ │ │ Related Alerts                  │ │ │ │ Assignee: John  │ │ │
│ │ │ Comments                        │ │ │ │ Severity: High  │ │ │
│ │ └─────────────────────────────────┘ │ │ └─────────────────┘ │ │
│ │                                     │ │ ┌─────────────────┐ │ │
│ │                                     │ │ │ Actions History │ │ │
│ │                                     │ │ │ Notifications   │ │ │
│ │                                     │ │ └─────────────────┘ │ │
│ └─────────────────────────────────────┘ └─────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

**Component Behaviors:**
- **Timeline**: Vertical timeline with timestamps, auto-scroll to latest
- **Comments**: Real-time updates, markdown support
- **Status Updates**: Animated transitions between states
- **Action Buttons**: Loading states, confirmation modals for destructive actions


#### Login Page Layout (Two-column, right-shifted)

```
┌───────────────────────────────────────────────────────────────────────────┐
│                                                                           │
│  Left Gap (Illustration / empty)   │         Right Content Area           │
│  35% width (reserved)              │  65% width (vertical centered)       │
│                                    │                                      │
│                                    │     ┌─────────────────────────────┐  │
│                                    │     │           [Logo]            │  │
│                                    │     │     Incident Management     │  │
│                                    │     │         System              │  │
│                                    │     └─────────────────────────────┘  │
│                                    │                                      │
│                                    │     ┌─────────────────────────────┐  │
│                                    │     │           Sign In           │  │
│                                    │     │  Email                      │  │
│                                    │     │  [input field]              │  │
│                                    │     │  Password                   │  │
│                                    │     │  [input field] [eye icon]   │  │
│                                    │     │  [□] Remember me            │  │
│                                    │     │  [Sign In] (full width)     │  │
│                                    │     └─────────────────────────────┘  │
│                                    │                                      │
└───────────────────────────────────────────────────────────────────────────┘
```

Design notes:
- Left column: reserved space (35% of viewport width) for brand illustration, product tagline or left intentionally blank. Keeps visual balance and creates a right-shifted form.
- Right column: content area vertically centered; the login card sits with max-width 420px and aligns visually to the right column center.

Layout Specifications:
- Page container: 100vw × 100vh, display: grid with two columns (35% / 65%).
- Login card: max-width: 420px; padding: 24px; border-radius: 12px; box-shadow: var(--card-shadow);
- Left gap: min-width 280px on wide screens; can host illustration or remain empty.
- Form fields: 44px height, 16px horizontal padding, 8px border-radius.
- Primary button: full width of card, 48px height.

Behavior & UX:
- The right-shifted layout guides user's focus to the right while keeping branding visible.
- Logo in card is optional; clicking it navigates home.
- Email field auto-focus; client-side validation with inline error messages.
- Password toggle (eye icon) and optional strength hint (if enabled) below the field.
- Remember me persisted via secure cookie/localStorage per policy.
- Sign-in button shows spinner and disabled state while authenticating.
- Inline error area below the form for server errors (slides in, aria-live="assertive").

Responsive rules:
- Desktop (>=1024px): Two-column grid (35%/65%). Left column visible and provides visual balance.
- Tablet (768px–1023px): Left column reduces to 30%; right column contains the card centered. If viewport narrow, left column can collapse to 20%.
- Mobile (<768px): Single column stack; left column content collapses above the card (or hidden) and the card becomes full width with 16px page padding.

Accessibility & Security:
- Focus order: email -> password -> remember -> sign-in -> secondary links.
- ARIA: form fields have aria-label/aria-describedby for errors.
- Rate limiting and captcha appear as previously specified (server-side).

#### Register Page Layout (Two-column, right-shifted)

```
┌───────────────────────────────────────────────────────────────────────────┐
│                                                                           │
│  Left Gap (Illustration / empty)   │         Right Content Area           │
│  35% width (reserved)              │  65% width (vertical centered)       │
│                                    │                                      │
│                                    │     ┌─────────────────────────────┐  │
│                                    │     │           [Logo]            │  │
│                                    │     │     Incident Management     │  │
│                                    │     │         System              │  │
│                                    │     └─────────────────────────────┘  │
│                                    │                                      │
│                                    │     ┌─────────────────────────────┐  │
│                                    │     │        Create Account       │  │
│                                    │     │  Full Name                  │  │
│                                    │     │  [input field]              │  │
│                                    │     │  Email                      │  │
│                                    │     │  [input field]              │  │
│                                    │     │  Password                   │  │
│                                    │     │  [input field] [eye icon]   │  │
│                                    │     │  Strength bar & hint        │  │
│                                    │     │  Confirm Password           │  │
│                                    │     │  [input field]              │  │
│                                    │     │  [□] I agree to Terms       │  │
│                                    │     │  [Create Account] (full)    │  │
│                                    │     └─────────────────────────────┘  │
│                                    │                                      │
└───────────────────────────────────────────────────────────────────────────┘
```

Design notes:
- Layout mirrors login: left reserved column (35%), right column contains the registration card (max-width 480px).
- Use visual hierarchy: labels, compact spacing, and clear affordances for password strength and validation.

Layout Specifications:
- Grid: two columns (35% / 65%) on wide screens.
- Register card: max-width 480px; padding: 28px; border-radius: 12px; subtle shadow.
- Form fields: 44px height; consistent spacing (12–16px gap) between inputs.

Behavior & UX:
- Real-time validation for name/email/password; password strength meter updates as user types.
- Confirm password verifies equality with main password field and shows immediate feedback.
- Terms checkbox required; Create Account button disabled until validations pass.
- Server-side errors displayed in inline banner with aria-live and focus management.

Responsive rules:
- Desktop (>=1024px): Two-column layout (left gap visible).
- Tablet (768px–1023px): Left column reduces; card remains centered in right column.
- Mobile (<768px): Single column; left gap hidden and card uses full width with 16px padding.

Security & UX:
- Captcha shown after repeated failed attempts or suspicious activity.
- Rate limiting enforced server-side; client shows explanatory messages.
- Email verification flow triggered after successful sign-up; token and expiry rules handled by backend.

### Component Specifications

#### Auth Layout (Shared)
**Overall Structure:**
```
┌─────────────────────────────────────────────────────────────────┐
│                    Auth Content Area                            │
│                    (Centered, Full Height)                      │
│                                                                 │
│                    ┌─────────────────────────────┐              │
│                    │       Page Content          │              │
│                    │       (Login/Register)      │              │
│                    └─────────────────────────────┘              │
│                                                                 │
│                    Footer Links                                 │
│                    [Privacy] [Terms] [Support]                   │
└─────────────────────────────────────────────────────────────────┘
```

**Layout Specifications:**
- **Background**: Full-screen gradient or branded background image
- **Content Container**: Flexbox centered both horizontally and vertically
- **Max Width**: 400px for login, 480px for register
- **Footer**: Fixed bottom, small text links, neutral color
- **Branding**: Consistent logo placement across all auth pages

**Responsive Considerations:**
- **Mobile**: Full width with 16px padding, footer becomes inline
- **Tablet**: Centered with background visible on sides
- **Desktop**: Centered with optional background pattern/image

### Component Specifications

#### Button Component Sizes
```css
/* Size Specifications */
.btn-sm {
  height: 32px;
  padding: 0 12px;
  font-size: 14px;
  border-radius: 6px;
}

.btn-md {
  height: 40px;
  padding: 0 16px;
  font-size: 16px;
  border-radius: 8px;
}

.btn-lg {
  height: 48px;
  padding: 0 24px;
  font-size: 18px;
  border-radius: 8px;
}
```

**Button Behaviors:**
- **Hover**: 150ms ease transition, brightness(1.1)
- **Active**: Scale(0.98) for 100ms
- **Loading**: Spinner replaces text, maintains width
- **Disabled**: opacity: 0.5, pointer-events: none

#### Modal Component Sizes
```css
/* Modal Size Specifications */
.modal-sm {
  width: 400px;
  max-width: 90vw;
}

.modal-md {
  width: 600px;
  max-width: 90vw;
}

.modal-lg {
  width: 800px;
  max-width: 95vw;
}

.modal-xl {
  width: 1200px;
  max-width: 95vw;
}
```

**Modal Behaviors:**
- **Open Animation**: fadeIn 200ms + slideUp 20px
- **Close Animation**: fadeOut 150ms + slideDown 20px
- **Backdrop Click**: Close modal (unless persistent)
- **Escape Key**: Close modal
- **Focus Management**: Trap focus within modal

#### Card Component
```css
.card {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  background: white;
  transition: all 0.2s ease;
}

.card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-1px);
}

.card-padding {
  padding: 24px;
}
```

#### Form Input Specifications
```css
.input-field {
  height: 44px;
  padding: 0 16px;
  border-radius: 8px;
  border: 1px solid #d1d5db;
  font-size: 16px;
  transition: all 0.2s ease;
}

.input-field:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.input-error {
  border-color: #ef4444;
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
}
```

### Responsive Breakpoint Behaviors

#### Tablet Layout (768px - 1023px)
- **Sidebar**: Overlay mode, hidden by default
- **Header**: Hamburger menu appears (left side)
- **Content**: Full width with 16px padding
- **Cards**: 2-column grid instead of 4-column
- **Tables**: Horizontal scroll with sticky first column

#### Mobile Layout (< 768px)
- **Navigation**: Bottom tab bar (60px height)
- **Header**: Simplified, search becomes modal
- **Content**: Single column, 12px padding
- **Cards**: Stack vertically, compact mode
- **Forms**: Full-width inputs, larger touch targets (48px min)
- **Modals**: Full-screen on small devices

### Performance Behaviors

#### Loading States
- **Skeleton Loading**: Grey placeholder shapes matching content
- **Progressive Loading**: Show basic layout first, then populate
- **Lazy Loading**: Images and components load when in viewport
- **Infinite Scroll**: Load more items as user scrolls (incidents list)

#### Animation Guidelines
- **Micro-interactions**: 150-300ms duration
- **Page transitions**: 200-400ms duration
- **Loading spinners**: 1s rotation cycle
- **Hover effects**: 150ms ease-out
- **Focus indicators**: Immediate (0ms)

### Accessibility Specifications

#### Focus Management
- **Focus Visible**: 2px solid #3b82f6 outline
- **Tab Order**: Logical flow, skip links provided
- **Screen Readers**: ARIA labels, live regions for updates

#### Color Contrast
- **Text**: Minimum 4.5:1 contrast ratio
- **Interactive Elements**: Minimum 3:1 contrast ratio
- **Status Indicators**: Not relying solely on color

#### Keyboard Navigation
- **Arrow Keys**: Navigate within components
- **Enter/Space**: Activate buttons and links
- **Escape**: Close modals and dropdowns
- **Tab/Shift+Tab**: Navigate between interactive elements

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