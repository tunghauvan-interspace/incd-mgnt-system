# Vue.js Frontend Integration

This document explains how the Vue.js frontend integrates with the Go backend.

## Build Process

The Vue.js application is built using Vite and outputs to the `web/static/` directory, which the Go backend serves directly.

### Development

```bash
cd web/frontend
npm run dev
```

The development server runs on http://localhost:5173 with API proxy to the Go backend on port 8080.

### Production Build

```bash
cd web/frontend  
npm run build
```

This generates optimized static assets in `web/static/`:
- `index.html` - Single page application entry point
- `css/` - Stylesheet bundles  
- `js/` - JavaScript bundles

## Go Backend Integration

The Go backend should serve the Vue.js SPA by:

1. **Serving static assets** from `web/static/` for CSS/JS files
2. **Serving `index.html`** for all frontend routes to enable client-side routing
3. **Keeping API routes** (`/api/*`) unchanged for the frontend to consume

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
    http.ServeFile(w, r, "web/static/index.html")
}

// In your router setup:
// Static files
http.Handle("/css/", http.StripPrefix("/", http.FileServer(http.Dir("web/static/"))))
http.Handle("/js/", http.StripPrefix("/", http.FileServer(http.Dir("web/static/"))))

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

All API data structures match the existing Go backend implementation.