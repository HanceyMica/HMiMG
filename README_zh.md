# HMiMG - 私有图片管理图库

<div align="center">
  <img src="client-vuetify/public/images/slogan_light.png" alt="HMiMG Logo" width="500" />
  <br/>
  <p>一个现代、响应式且注重隐私的图片管理系统。</p>
  <a href="https://github.com/HanceyMica/HMiMG/blob/master/README.md">English Doc</a>
</div>

## 功能特性

- **嵌套合集**: 将照片整理到相册和无限层级的嵌套合集中。
- **隐私优先**: 自托管解决方案，配备安全的身份验证。
- **拖拽上传**: 支持多文件批量上传。
- **响应式设计**: 完美适配桌面端和移动端。
- **深色/浅色模式**: 根据系统偏好自动切换主题。
- **后台管理**: 管理用户注册和系统设置。
- **多语言 (i18n)**: 支持英文、简体中文、日语。

## 技术栈

- **前端**: Vue 3 + Vuetify 3 + Vite
- **后端**: Go + Gin + GORM
- **数据库**: MySQL / PostgreSQL

## 安装部署 (本地开发)

### 后端 (Go)
1. 进入 `server-go/` 目录。
2. 将 `.env.example` 复制为 `.env`。
3. 在 `.env` 中配置您的数据库连接。
4. 运行:
   ```bash
   go mod tidy
   go run .
   ```
   *服务器运行在 `http://localhost:9108`*

### 前端 (Vuetify)
1. 进入 `client-vuetify/` 目录。
2. 运行:
   ```bash
   npm install
   npm run dev
   ```
   *客户端运行在 `http://localhost:9109`*

## 手动生产部署

### 后端 (Go)
1. 编译二进制文件:
   ```bash
   cd server-go
   go build -o hmimg-server .
   ```
2. 在 Release 模式下运行:
   ```bash
   # Linux/macOS
   export GIN_MODE=release
   ./hmimg-server
   
   # Windows (PowerShell)
   $env:GIN_MODE="release"
   .\hmimg-server.exe
   ```

### 前端 (Vuetify)
1. 构建静态文件:
   ```bash
   cd client-vuetify
   npm run build
   ```
2. 将 `dist/` 目录下生成的静态文件部署到 Nginx 或 Apache 等 Web 服务器上。

## Docker 部署

```bash
docker-compose up -d
```
- **前端**: `http://localhost:9109`
- **后端**: `http://localhost:9108`

## 数据库初始化

HMiMG 会自动处理数据库设置：

1. **手动创建数据库**: 在您的 MySQL 或 PostgreSQL 服务器中手动创建一个空的数据库（例如 `hmimg_db`）。
2. **自动迁移**: 在配置好你的 `.env` 文件后，当 Go 后端启动时，它会自动创建所有必要的表结构。
3. **预设数据**: 首次运行时，系统会自动填充默认设置并创建默认管理员账户。

*注意：如果您使用 Docker Compose 部署，数据库将自动创建并配置。*

## 默认管理员
- 用户名: `admin`
- 密码: `admin` (请登录后立即修改！)

## 文档
- [API 文档 (英文)](docs/api.md)
- [API 文档 (中文)](docs/api_zh-cn.md)

## 许可证
MIT
