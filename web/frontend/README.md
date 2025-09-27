# Incident Management Frontend

This is the Vue.js + TypeScript frontend for the Incident Management System.

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Type-Check, Compile and Minify for Production

```sh
npm run build
```

### Run Unit Tests with Vitest

```sh
npm run test:unit
```

### Lint with ESLint

```sh
npm run lint
```

### Format with Prettier

```sh
npm run format
```

## Project Structure

```
src/
├── assets/         # Static assets (CSS, images)
├── components/     # Reusable Vue components
├── composables/    # Vue 3 Composition API composables
├── router/         # Vue Router configuration
├── services/       # API service layer
├── stores/         # Pinia state management stores
├── types/          # TypeScript type definitions
├── utils/          # Utility functions
└── views/          # Page components
```

## Build Configuration

The build is configured to output to `../static/` (the web/static directory) so that the Go backend can serve the built Vue.js application.

## Development Features

- **Vue 3** with Composition API
- **TypeScript** for type safety
- **Vue Router** for SPA routing
- **Pinia** for state management
- **Axios** for HTTP requests
- **Chart.js** for data visualization (to be integrated)
- **ESLint + Prettier** for code quality
- **Vite** for fast development and building

## API Integration

The frontend communicates with the Go backend through API endpoints:
- `/api/incidents` - Incident management
- `/api/alerts` - Alert management  
- `/api/metrics` - Dashboard metrics

All API calls are proxied to `http://localhost:8080` during development.