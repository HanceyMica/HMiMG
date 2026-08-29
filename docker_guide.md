# 将 HMiMG 打包为 Docker Image（并用 docker compose 运行）

本项目由两部分组成：
- `server-go/`：Go 后端 API（默认端口 9108）
- `client-vuetify/`：Vue 3 + Vuetify 前端（默认端口 9109）

以下步骤会：
- 构建两个镜像：`hmimg-server` 与 `hmimg-client`
- 通过 `docker-compose.yml` 一键启动：MySQL + 后端 + 前端

---

## 1. 准备 Docker 文件（已提供）

仓库中已包含：
- `server-go/Dockerfile`
- `client-vuetify/Dockerfile`
- `docker-compose.yml`

---

## 2. 准备后端环境变量（必须）

后端依赖环境变量配置。在项目根目录或 `server-go/` 目录中：

```bash
cp server-go/.env.example server-go/.env
```

编辑 `server-go/.env`，重点修改：
- `JWT_SECRET`：生产环境必须使用强随机字符串
- `DB_HOST/USER/PASSWORD/NAME`：数据库连接信息

在 `docker-compose.yml` 中，后端会使用这些环境变量进行配置。

---

## 3. 启动（docker compose）

在项目根目录执行：

```bash
docker compose up -d --build
```

启动后：
- 前端：`http://localhost:9109`
- 后端：`http://localhost:9108`（API 前缀：`/api`）

首次启动后访问前端会自动进入安装向导（`/install`）：
1. 环境检查
2. 数据库配置（compose 已通过环境变量注入数据库连接，向导会自动跳过此步）
3. 设置管理员账号（自定义用户名/密码，密码至少 8 位）
4. 站点设置（标题/默认语言/最大用户数/是否开放注册）

安装完成后向导自动锁定，需重新登录即可使用。

---

## 4. 配置前端 API 地址（VITE_API_URL）

前端通过 `VITE_API_URL` 指定后端 API 地址。

本仓库提供的 `docker-compose.yml` 默认构建参数为：
- `http://localhost:9108/api`

如果你部署在服务器上并且使用域名/反向代理，建议改为：
- `https://你的域名/api`

修改方式：
1. 打开 `docker-compose.yml`
2. 找到 `client.build.args.VITE_API_URL`
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
- 检查 `VITE_API_URL` 是否正确
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
同时检查上传目录卷是否可写。

---

## 7. 生产部署建议（可选）

生产环境推荐：
- 使用反向代理（Nginx/Caddy/Traefik）统一暴露一个域名
- `/` 代理到前端（9109）
- `/api` 代理到后端（9108）
- 开启 HTTPS
- 设置强 `JWT_SECRET`
