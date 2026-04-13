# HMiMG - Vuetify Frontend

This is a refactored version of the HMiMG frontend using Vue 3 and Vuetify 3.

## Recent Updates
- Uploaded images are now loaded through `/api/files/*path`.
- After a successful upload, the target album and library views refresh immediately.
- The image details page now supports renaming and previous/next navigation.

## Tech Stack
- Vue 3
- Vuetify 3 (Material Design)
- Vite
- Pinia (State Management)
- Vue Router
- Vue I18n
- Axios

## Getting Started

### Prerequisites
- Node.js 18+

### Installation
```bash
npm install
```

### Development
```bash
npm run dev
```
Runs the development server on [http://localhost:9109](http://localhost:9109).

### Build
```bash
npm run build
```
Generates production build in `dist/`.

## Configuration
1. Copy `.env.example` to `.env`.
2. Edit `.env` and set your backend API URL.

Example:
```env
VITE_API_URL=/api
```

### What `.env.example` Is For
- `.env.example` is the committed template file for frontend environment variables.
- `client-vuetify/.env` is intended for local or server-specific values and is ignored by the repository root `.gitignore`.
- When cloning or deploying the project, create `client-vuetify/.env` from `.env.example` instead of editing the example file directly.

### `VITE_API_URL` Usage
- Use `VITE_API_URL=/api` when your web server proxies `/api` to the backend on the same domain.
- Use a full URL such as `VITE_API_URL=https://your-domain.com/api` when the frontend and backend are deployed separately.
- After changing `VITE_API_URL`, run `npm run build` again because Vite injects it at build time.
