# 将 HMiMG 打包为 Docker Image（并用 docker compose 运行）

本项目由两部分组成：
- `server/`：Koa 后端 API（默认端口 9108）
- `client/`：Next.js 前端（默认端口 9109）

以下步骤会：
- 构建两个镜像：`hmimg-server` 与 `hmimg-client`
- 通过 `docker-compose.yml` 一键启动：MySQL + 后端 + 前端

---

## 1. 准备 Docker 文件（已提供）

仓库中已包含：
- `server/Dockerfile`
- `client/Dockerfile`
- `docker-compose.yml`

---

## 2. 准备后端配置文件（必须）

后端依赖 `server/config/config.js`（该文件应包含数据库账号、JWT 密钥等敏感信息，默认应被 Git 忽略）。

在项目根目录执行：

```bash
cp server/config/config.example.js server/config/config.js
```

编辑 `server/config/config.js`，重点修改：
- `jwtSecret`：生产环境必须使用强随机字符串
- `database.connection.host/user/password/database`

在 `docker-compose.yml` 中，后端会把该配置文件以只读方式挂载到容器：

- `./server/config/config.js:/app/config/config.js:ro`

---

## 3. 启动（docker compose）

在项目根目录执行：

```bash
docker compose up -d --build
```

启动后：
- 前端：`http://localhost:9109`
- 后端：`http://localhost:9108`（API 前缀：`/api`）

后端启动时会自动：
- 执行数据库迁移
- 若数据库中不存在 `admin` 用户，会自动创建默认管理员

默认管理员：
- 用户名：`admin`
- 密码：`admin`

建议首次登录后立即修改密码。

---

## 4. 配置前端 API 地址（NEXT_PUBLIC_API_URL）

前端通过 `NEXT_PUBLIC_API_URL` 指定后端 API 地址。

本仓库提供的 `docker-compose.yml` 默认构建参数为：

- `http://localhost:9108/api`

如果你部署在服务器上并且使用域名/反向代理，建议改为：

- `https://你的域名/api`

修改方式：

1. 打开 `docker-compose.yml`
2. 找到 `client.build.args.NEXT_PUBLIC_API_URL`
3. 修改后重新构建：

```bash
docker compose up -d --build
```

---

## 5. 数据持久化（强烈建议）

`docker-compose.yml` 已默认启用数据卷：

- `hmimg_mysql`：MySQL 数据目录持久化
- `hmimg_uploads`：后端上传目录持久化（容器内 `/app/uploads`）

这意味着即使容器重建：
- 数据库数据不会丢
- 已上传图片不会丢

---

## 6. 常见问题

### 6.1 前端打不开 / 前端请求 API 失败

- 检查 `NEXT_PUBLIC_API_URL` 是否正确
- 检查 `server` 容器是否正常启动：`docker compose ps`
- 检查后端日志：`docker compose logs -f server`

### 6.2 注册按钮显示/隐藏不符合预期

登录页会请求：
- `GET /api/settings/public`

返回字段 `allow_registration` 会控制是否显示注册入口。

### 6.3 上传失败

当前后端限制：
- 单次最多 20 个文件
- 仅允许 `jpeg/png/gif/webp`

同时检查：
- 上传目录卷是否可写

---

## 7. 生产部署建议（可选）

生产环境推荐：
- 使用反向代理（Nginx/Caddy/Traefik）统一暴露一个域名
- `/` 代理到前端（9109）
- `/api` 代理到后端（9108）
- 开启 HTTPS
- 设置强 `JWT_SECRET`

