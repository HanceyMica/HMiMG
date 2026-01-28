# HMiMG - 私有图片管理图库

<div align="center">
  <img src="client/public/images/slogan_light.png" alt="HMiMG Logo" width="500" />
  <br/>
  <p>一个现代、响应式且注重隐私的图片管理系统。</p>
</div>

## 功能特性

- **嵌套合集**: 将照片整理到相册和无限层级的嵌套合集中。
- **隐私优先**: 自托管解决方案，配备安全的身份验证。
- **拖拽上传**: 支持多文件批量上传。
- **响应式设计**: 完美适配桌面端和移动端。
- **深色/浅色模式**: 根据系统偏好自动切换主题。
- **后台管理**: 管理用户注册和系统设置。
- **多语言 (i18n)**: 支持英文、简体中文、日语，语言偏好会通过 Cookie 记忆。

## 安装部署

1. **数据库设置**
   - 在您的 MySQL 或 PostgreSQL 服务器中创建一个名为 `hmimg_db` 的数据库。

2. **配置**
   - 进入 `server/config/` 目录。
   - 将 `config.example.js` 复制为 `config.js`。
   - 编辑 `server/config/config.js`，填入您的数据库凭据和 JWT 密钥。
   ```bash
   cp server/config/config.example.js server/config/config.js
   ```
   - 环境变量（可选）：
     - `PORT`（默认：`3001`）
     - `JWT_SECRET`（生产环境必填，请使用强随机字符串）
     - `NEXT_PUBLIC_API_URL`（前端 API 地址，默认：`http://localhost:3001/api`）

3. **安装依赖**
   - 在根目录下运行 `npm run install-all`（或者分别在 `server` 和 `client` 目录下运行 `npm install`）。

4. **启动应用**
   - 在根目录下运行 `npm run dev` 以同时启动服务器和客户端。
   - **服务器 (后端)** 运行在 `http://localhost:3001`
   - **客户端 (前端)** 运行在 `http://localhost:3000`

## 默认管理员
- 用户名: `admin`
- 密码: `admin` (请登录后立即修改！)
- 邮箱: `admin@yourdomain.com`（示例；首次启动自动创建的默认值可能不同，请登录后修改）

## 文档
- [API 文档 (英文)](docs/api.md)
- [API 文档 (中文)](docs/api_zh-cn.md)

## 许可证
MIT
