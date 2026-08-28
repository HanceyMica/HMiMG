# HMiMG 后端 API 文档

基础 URL: `/api`

除登录/注册/公开设置外，所有接口都需要在 Header 中携带有效的 JWT Token。
Header 格式: `Authorization: Bearer <token>`

通用错误响应：
```json
{ "error": "..." }
```

## 认证 (Authentication)

### 登录 (Login)
- **URL**: `/login`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "username": "admin",
    "password": "password"
  }
  ```
- **响应**:
  ```json
  {
    "token": "jwt_token_string",
    "user": { "id": 1, "username": "admin", "role": "admin" }
  }
  ```

### 注册 (Register)
- **URL**: `/register`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "username": "user",
    "password": "password",
    "email": "user@example.com",
    "phone": "1234567890"
  }
  ```
- **响应**:
  ```json
  { "message": "Registered successfully" }
  ```

### 更新管理员资料 (Update Admin Profile)
- **URL**: `/admin/update`
- **方法**: `PUT`
- **Header**: `Authorization: Bearer <token>`
- **请求体** (可选字段):
  ```json
  {
    "username": "newadmin",
    "email": "new@example.com",
    "phone": "+1234567890",
    "oldPassword": "currentPassword",
    "password": "newPassword"
  }
  ```
- **响应**:
  ```json
  { "message": "Profile updated successfully", "passwordChanged": true }
  ```

---

## 公开设置 (Public Settings) - 无需登录

### 获取公开设置
- **URL**: `/settings/public`
- **方法**: `GET`
- **响应**:
  ```json
  {
    "allow_registration": true,
    "website_title": "HMiMG",
    "default_language": "zh"
  }
  ```

## 相册 (Albums)

### 获取所有相册
- **URL**: `/albums`
- **方法**: `GET`
- **响应**: 相册对象数组。

### 获取单个相册
- **URL**: `/albums/:id`
- **方法**: `GET`
- **响应**: 相册对象。

### 创建相册
- **URL**: `/albums`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "name": "My Trip",
    "description": "Photos from 2024"
  }
  ```
- **响应**:
  ```json
  { "id": 1, "name": "My Trip", "description": "..." }
  ```

### 更新相册
- **URL**: `/albums/:id`
- **方法**: `PUT`
- **请求体**:
  ```json
  {
    "name": "Updated Name",
    "description": "Updated Description"
  }
  ```
- **响应**:
  ```json
  { "message": "Album updated", "id": "1" }
  ```

### 删除相册
- **URL**: `/albums/:id`
- **方法**: `DELETE`
- **响应**:
  ```json
  { "message": "Album deleted" }
  ```

---

## 合集 (Collections)

### 获取所有合集
- **URL**: `/collections`
- **方法**: `GET`
- **响应**: 合集对象数组。

### 获取单个合集
- **URL**: `/collections/:id`
- **方法**: `GET`
- **响应**: 合集对象，包含 `children` 数组（子相册或子合集）。

### 创建合集
- **URL**: `/collections`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "name": "My Portfolio",
    "description": "Best works"
  }
  ```
- **响应**:
  ```json
  { "id": 1, "name": "My Portfolio", ... }
  ```

### 更新合集
- **URL**: `/collections/:id`
- **方法**: `PUT`
- **请求体**:
  ```json
  { "name": "New Name", "description": "..." }
  ```
- **响应**:
  ```json
  { "message": "Collection updated", "id": "1" }
  ```

### 删除合集
- **URL**: `/collections/:id`
- **方法**: `DELETE`
- **响应**:
  ```json
  { "message": "Collection deleted" }
  ```

### 向合集添加项目
- **URL**: `/collections/add`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "collectionId": 1,
    "itemType": "album",  // 或 "collection"
    "itemId": 2            // 要添加的相册或合集 ID（推荐）
  }
  ```
- **说明**：
  - 通过 `itemId`（推荐）或旧版 `itemName` 标识被添加项。
  - 需要对目标合集和被添加项拥有所有权（或为管理员）。
- **响应**:
  ```json
  { "message": "Added successfully" }
  ```

### 从合集中获取随机图片
- **URL**: `/collections/:id/random`
- **方法**: `GET`
- **查询参数**: `?type=json` (默认) 或 `?type=redirect`
  - `type=json`: 返回图片对象的元数据 JSON。
  - `type=redirect`: 直接重定向 (302) 到图片文件的 URL。
- **注意**：
  - 该接口需要鉴权。由于 `<img>` 标签默认无法携带 `Authorization` Header，因此不建议直接使用 `type=redirect` 作为公开外链。
- **响应** (type=json):
  ```json
  {
    "id": 5,
    "filename": "1738021356789-123456.jpg",
    "original_name": "photo.jpg",
    ...
  }
  ```

---

## 图片 (Images)

### 上传图片
- **URL**: `/upload`
- **方法**: `POST`
- **Content-Type**: `multipart/form-data`
- **请求体**:
  - `images`: 文件 (二进制)
  - `albumId`: 目标相册 ID
- **说明**：
  - 单次请求最多上传 20 个文件。
  - 允许的 MIME 类型：`image/jpeg`、`image/png`、`image/gif`、`image/webp`。
  - 文件逐个处理，失败文件会记录在 `failures` 中，其余文件正常上传。
- **响应**:
  ```json
  { "ids": [1, 2], "count": 2, "failures": [] }
  ```

### 获取图片列表
- **URL**: `/images`
- **方法**: `GET`
- **查询参数**:
  - `?albumId=1` (可选过滤)
  - `?page=1&pageSize=50` (分页；`pageSize` 最大 200)
- **响应**:
  ```json
  { "items": [ /* 图片对象 */ ], "total": 120, "page": 1, "pageSize": 50 }
  ```

### 获取图片详情
- **URL**: `/images/:id`
- **方法**: `GET`
- **响应**: 图片对象及元数据。

### 更新图片
- **URL**: `/images/:id`
- **方法**: `PUT`
- **请求体**:
  ```json
  { "original_name": "新名称" }
  ```
- **说明**：需要为上传者本人或管理员。
- **响应**：更新后的图片对象。

### 删除图片
- **URL**: `/images/:id`
- **方法**: `DELETE`
- **说明**：需要为上传者本人或管理员。删除文件及数据库记录；若为相册封面则自动更新封面。
- **响应**:
  ```json
  { "message": "Image deleted" }
  ```

---

## 用户 (Users) - 仅限管理员

### 获取用户列表
- **URL**: `/admin/users`
- **方法**: `GET`
- **响应**: 用户对象数组（`id`、`username`、`email`、`phone`、`role`、`created_at`）。

### 更新用户角色
- **URL**: `/admin/users/:id/role`
- **方法**: `PUT`
- **请求体**:
  ```json
  { "role": "admin" }  // 或 "user"
  ```
- **说明**：不能修改自己的角色；不能降级最后一名管理员。
- **响应**:
  ```json
  { "message": "Role updated" }
  ```

### 删除用户
- **URL**: `/admin/users/:id`
- **方法**: `DELETE`
- **说明**：不能删除自己的账号；不能删除最后一名管理员。
- **响应**:
  ```json
  { "message": "User deleted" }
  ```

---

## 设置 (Settings) - 仅限管理员

### 获取设置
- **URL**: `/settings`
- **方法**: `GET`
- **响应**:
  ```json
  { "allow_registration": "true", "max_users": "100", "website_title": "HMiMG" }
  ```

### 更新设置
- **URL**: `/settings`
- **方法**: `PUT`
- **请求体**:
  ```json
  { "allow_registration": "false", "max_users": "50", "website_title": "HMiMG" }
  ```
- **响应**:
  ```json
  { "message": "Settings updated" }
  ```
