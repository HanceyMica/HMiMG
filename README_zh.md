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

## 最近更新

- **上传图片访问修复**: 已上传图片统一通过后端 `/api/files/*path` 读取，修复服务器部署后旧图片无法访问的问题。
- **上传后即时刷新**: 上传成功后会通知图库页和相册页刷新数据，并自动跳转到目标相册。
- **图片详情增强**: 图片详情页已支持修改图片名、查看上一张、下一张。

## 技术栈

- **前端**: Vue 3 + Vuetify 3 + Vite
- **后端**: Go + Gin + GORM
- **数据库**: MySQL / PostgreSQL

## 安装部署 (本地开发)

### 后端 (Go)
1. 进入 `server-go/` 目录。
2. 将 `.env.example` 复制为 `.env`，设置 `JWT_SECRET`（数据库连接可留空，安装向导会收集）。
3. 运行:
   ```bash
   go mod tidy
   go run .
   ```
   *服务器运行在 `http://localhost:9108`*

### 前端 (Vuetify)
1. 进入 `client-vuetify/` 目录。
2. 将 `.env.example` 复制为 `.env`。
3. 根据部署方式修改 `VITE_API_URL`：
   - 若前后端同域部署，并由反向代理将 `/api` 转发到后端，推荐直接使用 `VITE_API_URL=/api`
   - 若前后端分离部署，可填写完整地址，如 `VITE_API_URL=https://your-domain.com/api`
4. 注意：`client-vuetify/.env` 属于本地环境文件，会被仓库根目录 `.gitignore` 排除，不会提交到 Git；`client-vuetify/.env.example` 用于提供可提交的配置模板。
5. 运行:
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
3. 若服务器 API 地址发生变化，请先修改 `client-vuetify/.env` 后重新构建；`Vite` 会在构建时把 `VITE_API_URL` 写入产物。

## Docker 部署

```bash
docker-compose up -d
```
- **前端**: `http://localhost:9109`
- **后端**: `http://localhost:9108`

## 首次安装

HMiMG 使用 Web 安装向导，不再依赖默认账号或 SQL 导入：

1. 启动前后端服务（见下文）。`.env` 中数据库连接可留空——安装向导会收集并写回。
2. 打开站点，将自动跳转 `/install`。
3. 按向导操作：环境检查 → 数据库（测试连接、创建数据表；若已通过环境变量配置则自动跳过）→ 自定义管理员账号（密码至少 8 位）→ 站点设置（标题/默认语言/最大用户数/开放注册）。
4. 安装完成后向导自动锁定。旧版本部署（库中已有用户）升级时自动识别并补锁，无需重装。

*注意：Docker Compose 部署时数据库自动就绪，向导从管理员步骤开始。*

## 文档
- [API 文档 (英文)](docs/api.md)
- [API 文档 (中文)](docs/api_zh-cn.md)

## 许可证
MIT
