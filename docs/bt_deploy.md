# HMiMG 宝塔面板部署教程

适用：宝塔 Linux 面板（Nginx）。架构：前端静态站 + Go 后端常驻进程 + MySQL，`/api` 反向代理到后端。

```
浏览器 ──> 宝塔 Nginx(80/443)
             ├── /       → 前端静态文件 (dist)
             └── /api/*  → 反代 127.0.0.1:9108 (hmimg-server)
```

---

## 1. 宝塔环境准备

软件商店安装：

| 组件 | 用途 |
|------|------|
| Nginx | 站点 + 反向代理 |
| MySQL 5.7 / 8.0 | 数据库 |
| **Supervisor 管理器** | 守护 Go 后端进程（关键） |

防火墙（宝塔安全组）：只放行 `80`、`443`、宝塔面板端口。**9108 不要放行**，只走本机反代。

## 2. 创建数据库

宝塔 → 数据库 → 添加数据库：

- 库名：`hmimg_db`
- 用户名：`hmimg`
- 密码：自动生成并记下
- 字符集：`utf8mb4`

> 安装向导只建表，不建库。库必须先存在。

## 3. 部署后端

### 3.1 本地编译（Linux amd64）

开发机（Windows PowerShell）：

```powershell
cd server-go
$env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"
go build -o hmimg-server .
```

macOS/Linux 直接 `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o hmimg-server .`

### 3.2 准备 .env

```ini
PORT=9108
JWT_SECRET=换成长随机串_至少32位
TRUST_PROXY=true

# 数据库可留空，安装向导会收集并写回
# DB_DRIVER=mysql
# DB_HOST=127.0.0.1
# DB_PORT=3306
# DB_USER=hmimg
# DB_PASSWORD=你的库密码
# DB_NAME=hmimg_db
```

> `TRUST_PROXY=true`：已过 Nginx 反代，让后端取真实客户端 IP。

### 3.3 上传

宝塔文件管理，上传到 `/www/hmimg/server/`：

```
/www/hmimg/server/
├── hmimg-server   # 二进制，权限 755
└── .env
```

`uploads/` 目录无需手建，后端自动创建。

### 3.4 Supervisor 守护

软件商店 → Supervisor 管理器 → 添加守护进程：

- 名称：`hmimg-server`
- 启动用户：`root`
- **运行目录**：`/www/hmimg/server`（必须！`.env`、`uploads` 按工作目录相对寻址）
- 启动命令：`/www/hmimg/server/hmimg-server`

启动后验证：

```bash
curl http://127.0.0.1:9108/api/install/status
```

返回 JSON 即成功。

## 4. 部署前端

### 4.1 本地构建

```bash
cd client-vuetify
# .env 里：VITE_API_URL=/api   （同域反代模式）
bun run build   # 或 npm run build
```

产物在 `dist/`。

### 4.2 宝塔建站

宝塔 → 网站 → 添加站点：

- 域名：`img.example.com`（你的域名）
- 根目录：`/www/wwwroot/hmimg`
- PHP版本：**纯静态**

上传 `dist/` 内全部内容到 `/www/wwwroot/hmimg`。

### 4.3 伪静态（SPA 路由支持）

站点设置 → 伪静态：

```nginx
location / {
    try_files $uri $uri/ /index.html;
}
```

### 4.4 反向代理 /api

站点设置 → 反向代理 → 添加：

- 代理名称：`api`
- 目标URL：`http://127.0.0.1:9108`

或直接在配置文件手写（效果相同，可控性更强）：

```nginx
location /api {
    proxy_pass http://127.0.0.1:9108;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;

    # 上传放宽（后端单次最多 20 张图，Nginx 默认 1m 会 413）
    client_max_body_size 200m;

    # 大文件上传超时放宽
    proxy_read_timeout 300s;
    proxy_send_timeout 300s;
}
```

## 5. HTTPS

站点设置 → SSL → Let's Encrypt 申请证书 → 开启「强制 HTTPS」。

> 安装向导与登录 Cookie 在 HTTPS 下最安全，建议上线即配。

## 6. 首次安装

浏览器打开 `https://img.example.com`，自动跳 `/install` 向导：

1. **环境检查**：确认三项通过
2. **数据库**：主机 `127.0.0.1`、端口 `3306`、用户/密码/库名用第 2 步创建的；先「测试连接」再「保存并创建数据表」
3. **管理员**：自定义用户名 + 密码（≥8 位）
4. **站点设置**：标题 / 默认语言 / 最大用户数 / 是否开放注册 → 完成

向导自动锁定。Docker/老库升级场景自动识别，无需重装。

## 7. 常见问题

| 症状 | 原因 | 处理 |
|------|------|------|
| 图片上传 413 | Nginx body 限制 | `client_max_body_size 200m` |
| 大图上传超时 | 反代超时 | `proxy_read_timeout 300s` |
| 接口 502 | 后端进程挂了 | Supervisor 里查日志，确认工作目录是 `/www/hmimg/server` |
| 前端 404（刷新页面） | 伪静态没配 | 第 4.3 步 `try_files` |
| 图片能传不能显示 | 反代漏了 `/api` | 检查 location /api 是否生效 |
| 安装向导连不上库 | 库未建 / 账号错 | 宝塔数据库页核对；MySQL 8 注意认证插件 |
| 后端启动读不到 .env | 工作目录不对 | Supervisor 运行目录必须 `/www/hmimg/server` |

## 8. 更新版本

1. 本地重新编译后端二进制 + 前端 `bun run build`
2. Supervisor 停止 hmimg-server → 覆盖二进制 → 启动
3. 覆盖 `/www/wwwroot/hmimg` 前端文件

数据库无需动：已有部署（库中有用户）启动时自动识别为已安装。
