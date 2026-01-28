# HMiMG - Private Image Hosting

<div align="center">
  <img src="client/public/images/slogan_light.png" alt="HMiMG Logo" width="500" />
  <br/>
  <p>A modern, responsive, and private image management gallery.</p>
</div>

## Features

- **Nested Collections**: Organize your photos into albums and infinitely nested collections.
- **Privacy First**: Self-hosted solution with secure authentication.
- **Drag & Drop Upload**: Easy multi-file uploading.
- **Responsive Design**: Works great on desktop and mobile.
- **Dark/Light Mode**: Automatic theme switching based on system preference.
- **Admin Controls**: Manage user registration and system settings.
- **Internationalization (i18n)**: English, Simplified Chinese, and Japanese support with per-user preference persisted via cookie.

## Setup

1. **Database Setup**
   - Create a database named `hmimg_db` in your MySQL or PostgreSQL server.

2. **Configuration**
   - Navigate to `server/config/`.
   - Copy `config.example.js` to `config.js`.
   - Edit `server/config/config.js` with your database credentials and JWT secret.
   ```bash
   cp server/config/config.example.js server/config/config.js
   ```
   - Environment variables:
     - `PORT` (default: `9108`)
     - `JWT_SECRET` (required in production, use a strong random string)
     - `NEXT_PUBLIC_API_URL` (client API base, default `http://localhost:9108/api`)

3. **Installation**
   - Run `npm run install-all` in the root directory (or `npm install` in `server` and `client` separately).

4. **Running the App**
   - Run `npm run dev` in the root directory to start both server and client.
   - **Server** runs on `http://localhost:9108`
   - **Client** runs on `http://localhost:9109`

## Default Admin
- Username: `admin`
- Password: `admin` (Please change this immediately after login!)
- Email: `admin@yourdomain.com` (sample only; actual seeded value may vary, change after first login)

## Documentation
- [API Documentation (English)](docs/api.md)
- [API Documentation (Chinese)](docs/api_zh-cn.md)

## License
MIT
