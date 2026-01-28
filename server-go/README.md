# server-go (Go backend)

该目录是对现有 Node.js(Koa) 后端的 Go 版本重构，实现同一套 `/api` 接口，供前端复用。

## 端口与前缀

- 默认监听端口：`9108`（可通过 `PORT` 覆盖）
- API 前缀：`/api`
- 上传文件静态访问：`http(s)://<host>/<filename>`

## 环境变量

后端通过环境变量读取配置：

- `PORT`（默认：`9108`）
- `JWT_SECRET`（生产环境必须强随机）
- `DB_DRIVER`：`mysql` / `postgres`
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `UPLOAD_DIR`（默认：`uploads`）

程序启动时会自动读取当前目录下的 `.env` 文件（如果存在），并写入到环境变量中（仅在对应环境变量未设置时生效）。示例模板见 `.env.example`。

## 本地运行

在 `server-go/` 目录执行：

```bash
go mod tidy
go run .
```

启动后：
- `GET /` → `HMiMG API Server`
- `POST /api/login`、`POST /api/register` 等接口可用

## 数据库初始化

启动时会自动：

- AutoMigrate 创建/更新表结构（表名与 Node 版本一致：`hmimg_*`）
- 写入默认 Settings（`max_users`/`allow_registration`/`website_title`）
- 若不存在 `admin` 用户则创建默认管理员（用户名/密码 `admin/admin`）

