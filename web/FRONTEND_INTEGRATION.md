# Vue.js Frontend Integration Guide

This document provides detailed technical information about the Vue.js frontend integration with the Go backend.

## Architecture Overview

The frontend uses a modern Vue.js 3 + TypeScript stack with Vite for development and building. The architecture follows these principles:

- **Separation of Concerns**: Frontend and backend are decoupled
- **Development Efficiency**: Hot reload and fast builds with Vite
- **Type Safety**: Full TypeScript coverage for better code quality
- **Production Ready**: Optimized builds with code splitting and asset optimization

## Project Structure

```
web/frontend/
├── src/
│   ├── components/         # Reusable Vue components
│   │   ├── Navbar.vue     # Main navigation component
│   │   ├── Modal.vue      # Reusable modal dialog
│   │   └── DoughnutChart.vue  # Chart.js integration
│   ├── views/             # Page-level components (routes)
│   │   ├── Dashboard.vue  # Main dashboard with metrics
│   │   ├── Incidents.vue  # Incident management interface
│   │   └── Alerts.vue     # Alert viewing interface
│   ├── services/          # API integration layer
│   │   └── api.ts         # Axios-based API client
│   ├── types/             # TypeScript definitions
│   │   └── api.ts         # API response type definitions  
│   ├── utils/             # Utility functions
│   │   └── format.ts      # Date/time formatting helpers
│   ├── assets/            # Global styles and assets
│   │   └── main.css       # Global CSS with component styles
│   ├── router/            # Vue Router configuration
│   │   └── index.ts       # Route definitions for SPA
│   ├── App.vue            # Root component
│   └── main.ts            # Application entry point
├── dist/                  # Production build output (ignored by git)
├── package.json           # Dependencies and scripts
├── vite.config.ts         # Vite build configuration
├── tsconfig.json          # TypeScript configuration
└── .eslintrc.cjs          # ESLint configuration
```

## Build Process

The Vue.js application is built using Vite and outputs to the `web/frontend/dist/` directory. For production deployment, these files should be copied to the appropriate location for the Go backend to serve.

## Development Workflow

### Local Development Setup

1. **Start the Go backend** (required for API calls):
   ```bash
   # Option A: Using Docker Compose
   docker-compose --profile development up -d backend postgres

   # Option B: Manual Go development
   go run cmd/server/main.go
   ```

2. **Start the Vue.js development server**:
   ```bash
   cd web/frontend
   npm install
   npm run dev
   ```

3. **Access the application**:
   - Frontend: http://localhost:5173 (with hot reload)
   - Backend API: http://localhost:8080

### Development Features

- **Hot Module Replacement (HMR)**: Instant updates during development
- **API Proxy**: Automatic proxying of `/api/*` requests to Go backend
- **TypeScript Support**: Full type checking and IntelliSense
- **ESLint Integration**: Code quality and consistency checking
- **Source Maps**: Debugging support in browser dev tools

## Proxy configuration (container vs host)

The dev server proxy can be configured via environment variables. The Vite config reads `VITE_API_TARGET` (and falls back to `http://backend:8080`).

- When running the frontend inside Docker via `docker-compose` the default target `http://backend:8080` will correctly route to the backend service.

- When running the frontend on your host machine, set:

```powershell
# PowerShell example (Windows)
$env:VITE_API_TARGET = 'http://host.docker.internal:8080'
npm run dev
```

Or on macOS/Linux:

```bash
VITE_API_TARGET='http://host.docker.internal:8080' npm run dev
```

Vite loads `.env` files with `VITE_` prefixes automatically, so you can also place `VITE_API_TARGET` in `web/frontend/.env` for convenience.

## Production Build Process

### Building for Production

```bash
cd web/frontend
npm install           # Install dependencies
npm run type-check    # Verify TypeScript compilation
npm run build        # Create production build
```

### Build Output

The production build creates optimized assets in `web/frontend/dist/`:

```
dist/
├── index.html              # SPA entry point
├── css/
│   ├── index-[hash].css   # Main application styles
│   ├── Incidents-[hash].css  # Component-specific styles
│   └── ...                # Other component styles
└── js/
    ├── index-[hash].js    # Main application bundle
    ├── Incidents-[hash].js   # Route-based code splitting
    └── ...                # Other component chunks
```

**Key Features:**
- **Code Splitting**: Automatic route-based and component-based splitting
- **Asset Hashing**: Cache-busting with content-based hashes
- **Minification**: JavaScript and CSS minification for smaller bundles
- **Tree Shaking**: Unused code elimination
- **Source Maps**: Available for production debugging (optional)

### Build Configuration

The Vite configuration (`vite.config.ts`) includes:

```typescript
export default defineConfig({
  // Development server configuration
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',  // Go backend
        changeOrigin: true,
      }
    }
  },
  
  // Production build configuration  
  build: {
    outDir: 'dist',           // Build output directory
    emptyOutDir: true,        // Clean directory before build
    rollupOptions: {
      output: {
        // Organize assets by type with content hashing
        entryFileNames: 'js/[name]-[hash].js',
        chunkFileNames: 'js/[name]-[hash].js',
        assetFileNames: (info) => {
          // CSS files go to css/ directory
          if (/\.css$/i.test(info.name || '')) {
            return 'css/[name]-[hash][extname]'
          }
          return '[name]-[hash][extname]'
        }
      }
    }
  }
})
```

**Note**: Build artifacts are excluded from git to maintain clean version control.

## Go Backend Integration

The Go backend should serve the Vue.js SPA by:

1. **Serving static assets** from the appropriate directory for CSS/JS files  
2. **Serving `index.html`** for all frontend routes to enable client-side routing
3. **Keeping API routes** (`/api/*`) unchanged for the frontend to consume

**Deployment**: Copy the contents of `web/frontend/dist/` to your static file serving directory.

### Suggested Go Handler Changes

```go
// Serve Vue.js frontend for all non-API routes
func (h *Handler) handleSPA(w http.ResponseWriter, r *http.Request) {
    // Check if it's an API route
    if strings.HasPrefix(r.URL.Path, "/api/") {
        http.NotFound(w, r)
        return
    }
    
    // Serve index.html for all frontend routes
    http.ServeFile(w, r, "path/to/frontend/index.html")
}

// In your router setup:
// Static files (adjust path as needed)
http.Handle("/css/", http.StripPrefix("/", http.FileServer(http.Dir("path/to/static/"))))
http.Handle("/js/", http.StripPrefix("/", http.FileServer(http.Dir("path/to/static/"))))

// API routes (existing)
http.HandleFunc("/api/incidents", h.handleIncidents)
http.HandleFunc("/api/alerts", h.handleAlerts)
http.HandleFunc("/api/metrics", h.handleGetMetrics)

// SPA fallback (catch-all for frontend routes)
http.HandleFunc("/", h.handleSPA)
```

## Features Implemented

### ✅ Completed
- Vue 3 + TypeScript setup with Vite
- Vue Router for SPA navigation
- Pinia for state management
- Axios for API integration
- Chart.js for dashboard visualizations
- Modal components for incident/alert details
- Responsive design
- TypeScript type definitions for API
- Production build optimization

### 🚀 Ready for Integration
- Dashboard with metrics and charts
- Incidents list with acknowledge/resolve actions
- Alerts list with detail modals  
- Modern responsive UI
- Error handling and loading states

## API Compatibility

The frontend expects these API endpoints (unchanged from original):
- `GET /api/incidents` - List incidents
- `PUT /api/incidents/:id/acknowledge` - Acknowledge incident
- `PUT /api/incidents/:id/resolve` - Resolve incident
- `GET /api/alerts` - List alerts
- `GET /api/metrics` - Get dashboard metrics