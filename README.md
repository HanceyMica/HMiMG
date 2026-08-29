# HMiMG - Private Image Hosting

<div align="center">
  <img src="client-vuetify/public/images/slogan_light.png" alt="HMiMG Logo" width="500" />
  <br/>
  <p>A modern, responsive, and private image management gallery.</p>
  <a href="https://github.com/HanceyMica/HMiMG/blob/master/README_zh.md">中文文档</a>
</div>

## Features

- **Nested Collections**: Organize your photos into albums and infinitely nested collections.
- **Privacy First**: Self-hosted solution with secure authentication.
- **Drag & Drop Upload**: Easy multi-file uploading.
- **Responsive Design**: Works great on desktop and mobile.
- **Dark/Light Mode**: Automatic theme switching based on system preference.
- **Admin Controls**: Manage user registration and system settings.
- **Internationalization (i18n)**: English, Simplified Chinese, and Japanese support.

## Recent Updates

- **Uploaded image access fix**: Uploaded files are now served through `/api/files/*path`, fixing image loading issues after server deployment.
- **Instant refresh after upload**: The library and album pages now refresh after uploads, and the UI jumps to the target album automatically.
- **Image details enhancements**: The image details page now supports renaming images and navigating to the previous/next image.

## Tech Stack

- **Frontend**: Vue 3 + Vuetify 3 + Vite
- **Backend**: Go + Gin + GORM
- **Database**: MySQL / PostgreSQL

## Setup (Local Development)

### Backend (Go)
1. Navigate to `server-go/`.
2. Copy `.env.example` to `.env` and set `JWT_SECRET` (database connection is optional — the installer collects it).
3. Run:
   ```bash
   go mod tidy
   go run .
   ```
   *Server runs on `http://localhost:9108`*

### Frontend (Vuetify)
1. Navigate to `client-vuetify/`.
2. Copy `.env.example` to `.env`.
3. Set `VITE_API_URL` based on your deployment:
   - Use `VITE_API_URL=/api` when frontend and backend share the same domain and `/api` is reverse-proxied to the backend
   - Use a full URL such as `VITE_API_URL=https://your-domain.com/api` for separate deployments
4. Note: `client-vuetify/.env` is a local environment file and is ignored by the repository root `.gitignore`; `client-vuetify/.env.example` is the committed template for it.
5. Run:
   ```bash
   npm install
   npm run dev
   ```
   *Client runs on `http://localhost:9109`*

## Manual Deployment (Production)

### Backend (Go)
1. Build the binary:
   ```bash
   cd server-go
   go build -o hmimg-server .
   ```
2. Run in release mode:
   ```bash
   # Linux/macOS
   export GIN_MODE=release
   ./hmimg-server
   
   # Windows (PowerShell)
   $env:GIN_MODE="release"
   .\hmimg-server.exe
   ```

### Frontend (Vuetify)
1. Build the static files:
   ```bash
   cd client-vuetify
   npm run build
   ```
2. The generated files in `dist/` should be served by a web server like Nginx or Apache.
3. If the API address changes on the server, update `client-vuetify/.env` and rebuild first because `VITE_API_URL` is embedded during the build.

## Docker Deployment

```bash
docker-compose up -d
```
- **Client**: `http://localhost:9109`
- **Server**: `http://localhost:9108`

## First-Run Installation

HMiMG uses a web installer instead of any default credentials or SQL imports:

1. Start the backend and frontend (see Setup below). Database connection may be left unconfigured in `.env` — the installer collects it.
2. Open the site; you will be redirected to `/install` automatically.
3. Follow the wizard: environment check → database (test connection, create tables; skipped if already configured via env) → custom admin account (password >= 8 chars) → site settings (title, language, max users, registration).
4. The installer locks itself after completion. Existing deployments (with users already in the DB) are detected and locked automatically on upgrade — no re-install needed.

*Note: If you use Docker Compose, the database is provisioned automatically and the wizard starts at the admin step.*

## Documentation
- [API Documentation (English)](docs/api.md)
- [API Documentation (Chinese)](docs/api_zh-cn.md)

## License
MIT
