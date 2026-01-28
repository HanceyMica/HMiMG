# HMiMG Backend API Documentation

Base URL: `/api`

All endpoints except Login/Register/Public Settings require the `Authorization` header with a valid JWT token.
Header format: `Authorization: Bearer <token>`

Common error response:
```json
{ "error": "..." }
```

## Authentication

### Login
- **URL**: `/login`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "username": "admin",
    "password": "password"
  }
  ```
- **Response**:
  ```json
  {
    "token": "jwt_token_string",
    "user": { "id": 1, "username": "admin", "role": "admin" }
  }
  ```

### Register
- **URL**: `/register`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "username": "user",
    "password": "password",
    "email": "user@example.com",
    "phone": "1234567890"
  }
  ```
- **Response**:
  ```json
  { "message": "Registered successfully" }
  ```

### Update Admin Profile
- **URL**: `/admin/update`
- **Method**: `PUT`
- **Headers**: `Authorization: Bearer <token>`
- **Body** (optional fields):
  ```json
  {
    "username": "newadmin",
    "email": "new@example.com",
    "phone": "+1234567890",
    "oldPassword": "currentPassword",
    "password": "newPassword"
  }
  ```
- **Response**:
  ```json
  { "message": "Profile updated successfully", "passwordChanged": true }
  ```

---

## Public Settings (No Auth)

### Get Public Settings
- **URL**: `/settings/public`
- **Method**: `GET`
- **Response**:
  ```json
  {
    "allow_registration": true,
    "website_title": "HMiMG",
    "default_language": "zh"
  }
  ```

## Albums

### Get All Albums
- **URL**: `/albums`
- **Method**: `GET`
- **Response**: Array of album objects.

### Get Single Album
- **URL**: `/albums/:id`
- **Method**: `GET`
- **Response**: Album object.

### Create Album
- **URL**: `/albums`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "name": "My Trip",
    "description": "Photos from 2024"
  }
  ```
- **Response**:
  ```json
  { "id": 1, "name": "My Trip", "description": "..." }
  ```

### Update Album
- **URL**: `/albums/:id`
- **Method**: `PUT`
- **Body**:
  ```json
  {
    "name": "Updated Name",
    "description": "Updated Description"
  }
  ```
- **Response**:
  ```json
  { "message": "Album updated", "id": "1" }
  ```

### Delete Album
- **URL**: `/albums/:id`
- **Method**: `DELETE`
- **Response**:
  ```json
  { "message": "Album deleted" }
  ```

---

## Collections

### Get All Collections
- **URL**: `/collections`
- **Method**: `GET`
- **Response**: Array of collection objects.

### Get Single Collection
- **URL**: `/collections/:id`
- **Method**: `GET`
- **Response**: Collection object with `children` array (albums or sub-collections).

### Create Collection
- **URL**: `/collections`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "name": "My Portfolio",
    "description": "Best works"
  }
  ```
- **Response**:
  ```json
  { "id": 1, "name": "My Portfolio", ... }
  ```

### Update Collection
- **URL**: `/collections/:id`
- **Method**: `PUT`
- **Body**:
  ```json
  { "name": "New Name", "description": "..." }
  ```
- **Response**:
  ```json
  { "message": "Collection updated", "id": "1" }
  ```

### Delete Collection
- **URL**: `/collections/:id`
- **Method**: `DELETE`
- **Response**:
  ```json
  { "message": "Collection deleted" }
  ```

### Add Item to Collection
- **URL**: `/collections/add`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "collectionId": 1,
    "itemType": "album",  // or "collection"
    "itemName": "My Trip" // Name of the album/collection to add
  }
  ```
- **Response**:
  ```json
  { "message": "Added successfully" }
  ```

### Get Random Image from Collection
- **URL**: `/collections/:id/random`
- **Method**: `GET`
- **Query Params**: `?type=json` (default) or `?type=redirect`
  - `type=json`: Returns the image object metadata.
  - `type=redirect`: Redirects (302) directly to the image file URL.
- **Notes**:
  - This endpoint requires authentication. Using `type=redirect` directly in a plain `<img>` tag is not possible unless you can attach auth headers.
- **Response** (type=json):
  ```json
  {
    "id": 5,
    "filename": "1738021356789-123456.jpg",
    "original_name": "photo.jpg",
    ...
  }
  ```

---

## Images

### Upload Images
- **URL**: `/upload`
- **Method**: `POST`
- **Content-Type**: `multipart/form-data`
- **Body**:
  - `images`: File(s) (binary)
  - `albumId`: ID of the target album
- **Notes**:
  - Max 20 files per request.
  - Allowed MIME types: `image/jpeg`, `image/png`, `image/gif`, `image/webp`.
- **Response**:
  ```json
  { "ids": [1, 2], "count": 2 }
  ```

### Get Images
- **URL**: `/images`
- **Method**: `GET`
- **Query Params**: `?albumId=1` (optional filter)
- **Response**: Array of image objects.

### Get Image Details
- **URL**: `/images/:id`
- **Method**: `GET`
- **Response**: Image object with metadata.

---

## Settings (Admin Only)

### Get Settings
- **URL**: `/settings`
- **Method**: `GET`
- **Response**:
  ```json
  { "allow_registration": "true", "max_users": "100", "website_title": "HMiMG" }
  ```

### Update Settings
- **URL**: `/settings`
- **Method**: `PUT`
- **Body**:
  ```json
  { "allow_registration": "false", "max_users": "50", "website_title": "HMiMG" }
  ```
- **Response**:
  ```json
  { "message": "Settings updated" }
  ```
