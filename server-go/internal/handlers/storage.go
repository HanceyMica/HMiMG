package handlers

import (
	"errors"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"hmimg-server-go/internal/config"
	"hmimg-server-go/internal/httpx"
	"hmimg-server-go/internal/middleware"
	"hmimg-server-go/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StorageHandler 存储处理器，负责相册、收藏集、图片的 CRUD 操作
type StorageHandler struct {
	DB        *gorm.DB        // 数据库连接实例
	Cfg       config.Config   // 应用配置
	DBDriver  string          // 数据库驱动类型（用于兼容不同数据库的 SQL 语法）
	UploadDir string          // 文件上传存储目录
}

// ============================================================
// 相册相关请求体结构
// ============================================================

// createAlbumRequest 创建相册请求体
type createAlbumRequest struct {
	Name        string `json:"name"`         // 相册名称
	Description string `json:"description"` // 相册描述
}

// updateAlbumRequest 更新相册请求体
type updateAlbumRequest struct {
	Name        string `json:"name"`         // 新相册名称（可选）
	Description string `json:"description"` // 新相册描述
}

// ============================================================
// 相册相关处理器
// ============================================================

// CreateAlbum 创建新相册
// 仅允许已登录用户创建
//
// POST /api/albums（需认证）
// 请求体：{"name": "...", "description": "..."}
// 成功响应：{"id": 1, "name": "...", "description": "..."}
func (h StorageHandler) CreateAlbum(c *gin.Context) {
	var req createAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 检查相册名是否已存在
	var existing models.Album
	if err := h.DB.First(&existing, "name = ?", req.Name).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Album name already exists"})
		return
	}

	// 获取当前登录用户 ID（由 RequireAuth 中间件设置）
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)

	// 创建相册记录
	album := models.Album{Name: req.Name, Description: req.Description, CreatedBy: &userID}
	if err := h.DB.Create(&album).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create album"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": album.ID, "name": album.Name, "description": album.Description})
}

// GetAlbums 获取所有相册列表（按 ID 倒序）
//
// GET /api/albums（需认证）
// 响应：相册数组
func (h StorageHandler) GetAlbums(c *gin.Context) {
	var albums []models.Album
	if err := h.DB.Order("id DESC").Find(&albums).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load albums"})
		return
	}
	out := make([]gin.H, 0, len(albums))
	for _, a := range albums {
		out = append(out, albumToJSON(a))
	}
	c.JSON(http.StatusOK, out)
}

// GetAlbum 获取单个相册详情
//
// GET /api/albums/:id（需认证）
// 成功响应：相册信息
// 失败响应：404 {"error": "Album not found"}
func (h StorageHandler) GetAlbum(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	var album models.Album
	if err := h.DB.First(&album, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}
	c.JSON(http.StatusOK, albumToJSON(album))
}

// UpdateAlbum 更新相册信息
// 仅相册创建者或管理员可以更新
//
// PUT /api/albums/:id（需认证）
// 请求体：{"name": "...", "description": "..."}
func (h StorageHandler) UpdateAlbum(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var req updateAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 查询相册
	var album models.Album
	if err := h.DB.First(&album, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	// 权限检查：仅管理员或创建者可更新
	roleAny, _ := c.Get(middleware.ContextRoleKey)
	role, _ := roleAny.(string)
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)
	if role != "admin" && (album.CreatedBy == nil || *album.CreatedBy != userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// 检查新名称是否与其他相册冲突
	if req.Name != "" && req.Name != album.Name {
		var existing models.Album
		if err := h.DB.First(&existing, "name = ?", req.Name).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Album name already exists"})
			return
		}
		album.Name = req.Name
	}
	album.Description = req.Description

	if err := h.DB.Save(&album).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update album"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Album updated", "id": id})
}

// DeleteAlbum 删除相册及其所有图片
// 仅管理员或相册创建者可删除
// 删除操作在事务中执行，先删除图片文件，再删除数据库记录
//
// DELETE /api/albums/:id（需认证）
func (h StorageHandler) DeleteAlbum(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	// 查询相册
	var album models.Album
	if err := h.DB.First(&album, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	// 权限检查
	roleAny, _ := c.Get(middleware.ContextRoleKey)
	role, _ := roleAny.(string)
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)
	if role != "admin" && (album.CreatedBy == nil || *album.CreatedBy != userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// 先删除磁盘上的图片文件
	var images []models.Image
	_ = h.DB.Where("album_id = ?", id).Find(&images).Error
	for _, img := range images {
		_ = os.Remove(filepath.Join(h.UploadDir, img.Path))
	}

	// 在事务中删除数据库记录
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		// 删除相册中的所有图片记录
		if err := tx.Where("album_id = ?", id).Delete(&models.Image{}).Error; err != nil {
			return err
		}
		// 删除收藏集中引用此相册的条目
		if err := tx.Where("item_type = ? AND item_id = ?", "album", id).Delete(&models.CollectionItem{}).Error; err != nil {
			return err
		}
		// 删除相册记录
		if err := tx.Delete(&models.Album{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete album"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album deleted"})
}

// ============================================================
// 收藏集相关请求体结构
// ============================================================

// createCollectionRequest 创建收藏集请求体
type createCollectionRequest struct {
	Name        string `json:"name"`         // 收藏集名称
	Description string `json:"description"` // 收藏集描述
}

// updateCollectionRequest 更新收藏集请求体
type updateCollectionRequest struct {
	Name        string `json:"name"`         // 新收藏集名称（可选）
	Description string `json:"description"` // 新收藏集描述
}

// ============================================================
// 收藏集相关处理器
// ============================================================

// CreateCollection 创建新收藏集
//
// POST /api/collections（需认证）
// 请求体：{"name": "...", "description": "..."}
func (h StorageHandler) CreateCollection(c *gin.Context) {
	var req createCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 检查名称是否已存在
	var existing models.Collection
	if err := h.DB.First(&existing, "name = ?", req.Name).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Collection name already exists"})
		return
	}

	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)

	col := models.Collection{Name: req.Name, Description: req.Description, CreatedBy: &userID}
	if err := h.DB.Create(&col).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create collection"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": col.ID, "name": col.Name, "description": col.Description})
}

// GetCollections 获取所有收藏集列表（按 ID 倒序）
//
// GET /api/collections（需认证）
func (h StorageHandler) GetCollections(c *gin.Context) {
	var cols []models.Collection
	if err := h.DB.Order("id DESC").Find(&cols).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load collections"})
		return
	}
	out := make([]gin.H, 0, len(cols))
	for _, col := range cols {
		out = append(out, collectionToJSON(col))
	}
	c.JSON(http.StatusOK, out)
}

// GetCollection 获取收藏集详情，包括其包含的所有条目
//
// GET /api/collections/:id（需认证）
// 响应包含 children 字段，列出收藏集中的相册或子收藏集
func (h StorageHandler) GetCollection(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var col models.Collection
	if err := h.DB.First(&col, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	// 查询收藏集中的所有条目
	var items []models.CollectionItem
	if err := h.DB.Where("collection_id = ?", id).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load collection items"})
		return
	}

	// 根据条目类型批量查询详情
	children := make([]gin.H, 0)
	if len(items) > 0 {
		itemType := items[0].ItemType
		ids := make([]uint32, 0, len(items))
		for _, it := range items {
			ids = append(ids, it.ItemID)
		}

		// 统一类型的条目一起查询（避免混合类型）
		if itemType == "album" {
			var albums []models.Album
			_ = h.DB.Where("id IN ?", ids).Find(&albums).Error
			for _, a := range albums {
				j := albumToJSON(a)
				j["type"] = "album"
				children = append(children, j)
			}
		} else if itemType == "collection" {
			var cols2 []models.Collection
			_ = h.DB.Where("id IN ?", ids).Find(&cols2).Error
			for _, c2 := range cols2 {
				j := collectionToJSON(c2)
				j["type"] = "collection"
				children = append(children, j)
			}
		}
	}

	// 组装响应
	out := collectionToJSON(col)
	out["children"] = children
	c.JSON(http.StatusOK, out)
}

// UpdateCollection 更新收藏集信息
// 仅管理员或创建者可更新
//
// PUT /api/collections/:id（需认证）
func (h StorageHandler) UpdateCollection(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var req updateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var col models.Collection
	if err := h.DB.First(&col, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	// 权限检查
	roleAny, _ := c.Get(middleware.ContextRoleKey)
	role, _ := roleAny.(string)
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)
	if role != "admin" && (col.CreatedBy == nil || *col.CreatedBy != userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// 检查新名称是否冲突
	if req.Name != "" && req.Name != col.Name {
		var existing models.Collection
		if err := h.DB.First(&existing, "name = ?", req.Name).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Collection name already exists"})
			return
		}
		col.Name = req.Name
	}
	col.Description = req.Description

	if err := h.DB.Save(&col).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update collection"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Collection updated", "id": id})
}

// DeleteCollection 删除收藏集
// 仅管理员或创建者可删除
// 会同时删除收藏集中的所有条目记录
//
// DELETE /api/collections/:id（需认证）
func (h StorageHandler) DeleteCollection(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var col models.Collection
	if err := h.DB.First(&col, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	// 权限检查
	roleAny, _ := c.Get(middleware.ContextRoleKey)
	role, _ := roleAny.(string)
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)
	if role != "admin" && (col.CreatedBy == nil || *col.CreatedBy != userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// 在事务中删除
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		// 删除引用此收藏集的条目（如其他收藏集中引用了此收藏集）
		if err := tx.Where("item_type = ? AND item_id = ?", "collection", id).Delete(&models.CollectionItem{}).Error; err != nil {
			return err
		}
		// 删除此收藏集中的所有条目
		if err := tx.Where("collection_id = ?", id).Delete(&models.CollectionItem{}).Error; err != nil {
			return err
		}
		// 删除收藏集记录
		if err := tx.Delete(&models.Collection{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete collection"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Collection deleted"})
}

// ============================================================
// 收藏集条目操作
// ============================================================

// addToCollectionRequest 添加到收藏集请求体
type addToCollectionRequest struct {
	CollectionID uint32 `json:"collectionId"` // 目标收藏集 ID
	ItemType     string `json:"itemType"`     // 条目类型：album 或 collection
	ItemName     string `json:"itemName"`     // 条目名称（通过名称查找 ID）
}

// AddToCollection 将相册或收藏集添加到另一个收藏集
// 同一收藏集内的条目必须类型一致（不能同时包含相册和收藏集）
//
// POST /api/collections/add（需认证）
// 请求体：{"collectionId": 1, "itemType": "album", "itemName": "我的相册"}
func (h StorageHandler) AddToCollection(c *gin.Context) {
	var req addToCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.CollectionID == 0 || req.ItemName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if req.ItemType != "album" && req.ItemType != "collection" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item type"})
		return
	}

	// 验证目标收藏集存在
	var collection models.Collection
	if err := h.DB.First(&collection, "id = ?", req.CollectionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	// 根据类型查找条目 ID
	var itemID uint32
	if req.ItemType == "album" {
		var album models.Album
		if err := h.DB.First(&album, "name = ?", req.ItemName).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		itemID = album.ID
	} else {
		var col models.Collection
		if err := h.DB.First(&col, "name = ?", req.ItemName).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Target collection not found"})
			return
		}
		// 不允许将自己添加到自己
		if col.ID == req.CollectionID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add collection to itself"})
			return
		}
		itemID = col.ID
	}

	// 检查收藏集是否已包含混合类型条目
	var existingOne models.CollectionItem
	if err := h.DB.First(&existingOne, "collection_id = ?", req.CollectionID).Error; err == nil {
		if existingOne.ItemType != req.ItemType {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Collection cannot contain mixed types"})
			return
		}
	}

	// 检查条目是否已在收藏集中
	var exists models.CollectionItem
	if err := h.DB.First(&exists, "collection_id = ? AND item_type = ? AND item_id = ?", req.CollectionID, req.ItemType, itemID).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item already in collection"})
		return
	}

	// 创建收藏集条目
	toCreate := models.CollectionItem{CollectionID: req.CollectionID, ItemType: req.ItemType, ItemID: itemID}
	if err := h.DB.Create(&toCreate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Added successfully"})
}

// GetRandomFromCollection 从收藏集中随机获取一张图片
// 支持递归获取（收藏集可以包含子收藏集）
// 返回类型可以是 JSON 或重定向到图片文件
//
// GET /api/collections/:id/random（需认证）
// 查询参数：type=json（默认）或 type=redirect
func (h StorageHandler) GetRandomFromCollection(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	returnType := c.DefaultQuery("type", "json")

	var col models.Collection
	if err := h.DB.First(&col, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	// 递归获取收藏集及其子收藏集中的所有相册 ID
	albumIDs, err := h.getAlbumIDsFromCollection(id, map[uint32]bool{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get albums"})
		return
	}
	if len(albumIDs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No albums in this collection"})
		return
	}

	// 从这些相册中随机选一张图片
	var image models.Image
	q := h.DB.Where("album_id IN ?", albumIDs)
	dialect := strings.ToLower(h.DBDriver)
	// 不同数据库的随机函数不同
	if dialect == "postgres" || dialect == "pg" || dialect == "postgresql" {
		q = q.Order("RANDOM()")
	} else {
		q = q.Order("RAND()")
	}
	if err := q.First(&image).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No images found in this collection"})
		return
	}

	// 根据返回类型决定响应方式
	if returnType == "redirect" {
		base := httpx.Origin(c.Request)
		c.Redirect(http.StatusFound, base+"/api/files/"+url.PathEscape(image.Path))
		return
	}
	c.JSON(http.StatusOK, imageToJSON(image))
}

// getAlbumIDsFromCollection 递归获取收藏集中所有相册 ID
// 用于 GetRandomFromCollection 的随机图片选择
//
// 参数：
//   - collectionID: 要查询的收藏集 ID
//   - visited: 已访问的收藏集 ID 集合（防止循环引用）
//
// 返回值：
//   - []uint32: 所有相册 ID 列表
//   - error: 查询失败时的错误
func (h StorageHandler) getAlbumIDsFromCollection(collectionID uint32, visited map[uint32]bool) ([]uint32, error) {
	// 检测循环引用
	if visited[collectionID] {
		return nil, nil
	}
	visited[collectionID] = true

	var items []models.CollectionItem
	if err := h.DB.Where("collection_id = ?", collectionID).Find(&items).Error; err != nil {
		return nil, err
	}

	var albumIDs []uint32
	for _, it := range items {
		if it.ItemType == "album" {
			// 直接是相册，添加到结果集
			albumIDs = append(albumIDs, it.ItemID)
		} else if it.ItemType == "collection" {
			// 是子收藏集，递归获取
			childAlbumIDs, err := h.getAlbumIDsFromCollection(it.ItemID, visited)
			if err != nil {
				return nil, err
			}
			albumIDs = append(albumIDs, childAlbumIDs...)
		}
	}
	return albumIDs, nil
}

// ============================================================
// 图片上传和获取
// ============================================================

// UploadImages 上传图片到指定相册
// 支持批量上传（最多 20 张），仅允许图片格式（JPEG、PNG、GIF、WebP）
//
// POST /api/upload（需认证）
// 表单字段：albumId（相册ID）, images（文件数组）
// 成功响应：{"ids": [1,2,3], "count": 3}
func (h StorageHandler) UploadImages(c *gin.Context) {
	albumIDStr := c.PostForm("albumId")
	if albumIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Album ID required"})
		return
	}
	albumID64, err := strconv.ParseUint(albumIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Album ID required"})
		return
	}
	albumID := uint32(albumID64)

	// 验证相册存在
	var album models.Album
	if err := h.DB.First(&album, "id = ?", albumID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	// 解析 multipart 表单
	form, err := c.MultipartForm()
	if err != nil || form == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}
	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}
	if len(files) > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Too many files"})
		return
	}

	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)

	// 逐个处理上传的文件
	insertedIDs := make([]uint32, 0, len(files))
	var firstFilename string
	for i, file := range files {
		// 验证文件类型
		if ok := isAllowedMime(file); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed."})
			return
		}

		// 生成唯一文件名
		filename := uniqueFilename(file.Filename)
		if i == 0 {
			firstFilename = filename
		}

		// 保存文件到磁盘
		dst := filepath.Join(h.UploadDir, filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// 创建图片记录
		img := models.Image{
			Filename:     filename,
			OriginalName: file.Filename,
			Path:         filename,
			Size:         file.Size,
			Mimetype:     file.Header.Get("Content-Type"),
			AlbumID:      albumID,
			UploadedBy:   &userID,
		}
		if img.Mimetype == "" {
			img.Mimetype = "application/octet-stream"
		}
		if err := h.DB.Create(&img).Error; err != nil {
			// 数据库写入失败，删除已保存的文件
			_ = os.Remove(dst)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record image"})
			return
		}
		insertedIDs = append(insertedIDs, img.ID)
	}

	// 如果相册没有封面，设置第一张上传的图片为封面
	if album.CoverImage == nil && firstFilename != "" {
		cover := firstFilename
		album.CoverImage = &cover
		_ = h.DB.Save(&album).Error
	}

	c.JSON(http.StatusOK, gin.H{"ids": insertedIDs, "count": len(files)})
}

// GetUploadedFile 获取已上传的文件（静态文件服务）
// 路径参数为相对于上传目录的文件路径
//
// GET /api/files/*path
func (h StorageHandler) GetUploadedFile(c *gin.Context) {
	relPath, err := sanitizeUploadPath(c.Param("path"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file path"})
		return
	}

	fullPath := filepath.Join(h.UploadDir, relPath)
	fi, statErr := os.Stat(fullPath)
	if statErr != nil || fi.IsDir() {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(fullPath)
}

// GetImages 获取图片列表（可按相册筛选）
//
// GET /api/images（需认证）
// 查询参数：albumId（可选，筛选指定相册）
// 按 ID 倒序返回
func (h StorageHandler) GetImages(c *gin.Context) {
	albumIDStr := c.Query("albumId")
	q := h.DB.Model(&models.Image{}).Order("id DESC")
	if albumIDStr != "" {
		albumID64, err := strconv.ParseUint(albumIDStr, 10, 32)
		if err == nil {
			q = q.Where("album_id = ?", uint32(albumID64))
		}
	}
	var images []models.Image
	if err := q.Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load images"})
		return
	}
	out := make([]gin.H, 0, len(images))
	for _, img := range images {
		out = append(out, imageToJSON(img))
	}
	c.JSON(http.StatusOK, out)
}

// GetImage 获取单张图片详情
//
// GET /api/images/:id（需认证）
// 响应包含图片所属相册的名称
func (h StorageHandler) GetImage(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var img models.Image
	if err := h.DB.Order("id DESC").First(&img, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// 查询所属相册名称
	var album models.Album
	albumErr := h.DB.First(&album, "id = ?", img.AlbumID).Error
	out := imageToJSON(img)
	if albumErr != nil || album.Name == "" {
		out["album_name"] = nil
	} else {
		out["album_name"] = album.Name
	}
	c.JSON(http.StatusOK, out)
}

// updateImageRequest 更新图片请求体
type updateImageRequest struct {
	OriginalName    string `json:"original_name"`    // 新的原始文件名（下划线命名）
	OriginalNameAlt string `json:"originalName"`    // 新的原始文件名（驼峰命名，兼容旧客户端）
}

// UpdateImage 更新图片信息（目前仅支持更新原始文件名）
// 仅管理员或上传者可更新
//
// PUT /api/images/:id（需认证）
// 请求体：{"original_name": "新文件名.jpg"}
func (h StorageHandler) UpdateImage(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var req updateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 优先使用下划线命名，其次驼峰命名
	name := strings.TrimSpace(req.OriginalName)
	if name == "" {
		name = strings.TrimSpace(req.OriginalNameAlt)
	}
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image name is required"})
		return
	}

	// 查询图片
	var img models.Image
	if err := h.DB.First(&img, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// 权限检查：仅管理员或上传者可更新
	roleAny, _ := c.Get(middleware.ContextRoleKey)
	role, _ := roleAny.(string)
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)
	if role != "admin" && (img.UploadedBy == nil || *img.UploadedBy != userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	img.OriginalName = name
	if err := h.DB.Save(&img).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update image"})
		return
	}

	c.JSON(http.StatusOK, imageToJSON(img))
}

// ============================================================
// 工具函数
// ============================================================

// parseUintParam 解析 URL 路径参数中的无符号整数
// 用于解析 :id 等路径参数
//
// 参数：
//   - c: Gin Context
//   - key: 参数名（如 "id"）
//
// 返回值：
//   - uint32: 解析出的整数
//   - bool: 解析是否成功
func parseUintParam(c *gin.Context, key string) (uint32, bool) {
	idStr := c.Param(key)
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return 0, false
	}
	return uint32(id64), true
}

// sanitizeUploadPath 清理和验证上传文件路径
// 防止路径遍历攻击（../）和绝对路径访问
//
// 参数：
//   - raw: 原始路径字符串
//
// 返回值：
//   - string: 清理后的相对路径
//   - error: 路径无效时的错误
func sanitizeUploadPath(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	trimmed = strings.TrimPrefix(trimmed, "/")
	if trimmed == "" {
		return "", errors.New("empty path")
	}

	cleaned := filepath.Clean(trimmed)
	if cleaned == "." || cleaned == "" {
		return "", errors.New("empty path")
	}
	// 禁止绝对路径、.. 和以 .. 开头的路径
	if filepath.IsAbs(cleaned) || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", errors.New("invalid path")
	}

	return cleaned, nil
}

// albumToJSON 将 Album 模型转换为 JSON 响应格式
func albumToJSON(a models.Album) gin.H {
	var cover interface{} = nil
	if a.CoverImage != nil {
		cover = *a.CoverImage
	}
	var createdBy interface{} = nil
	if a.CreatedBy != nil {
		createdBy = *a.CreatedBy
	}
	return gin.H{
		"id":          a.ID,
		"name":        a.Name,
		"description": a.Description,
		"created_by":  createdBy,
		"cover_image": cover,
		"created_at":  a.CreatedAt,
		"updated_at":  a.UpdatedAt,
	}
}

// collectionToJSON 将 Collection 模型转换为 JSON 响应格式
func collectionToJSON(col models.Collection) gin.H {
	var createdBy interface{} = nil
	if col.CreatedBy != nil {
		createdBy = *col.CreatedBy
	}
	return gin.H{
		"id":          col.ID,
		"name":        col.Name,
		"description": col.Description,
		"created_by":  createdBy,
		"created_at":  col.CreatedAt,
		"updated_at":  col.UpdatedAt,
	}
}

// imageToJSON 将 Image 模型转换为 JSON 响应格式
func imageToJSON(img models.Image) gin.H {
	var uploadedBy interface{} = nil
	if img.UploadedBy != nil {
		uploadedBy = *img.UploadedBy
	}
	return gin.H{
		"id":            img.ID,
		"filename":      img.Filename,
		"original_name": img.OriginalName,
		"path":          img.Path,
		"size":          img.Size,
		"mimetype":      img.Mimetype,
		"album_id":      img.AlbumID,
		"uploaded_by":   uploadedBy,
		"created_at":    img.CreatedAt,
		"updated_at":    img.UpdatedAt,
	}
}

// uniqueFilename 生成唯一的文件名
// 格式：时间戳-随机数.扩展名
// 时间戳使用毫秒级 Unix 时间，确保唯一性
func uniqueFilename(original string) string {
	ext := filepath.Ext(original)
	if ext == "" {
		ext = ".bin"
	}
	// 使用当前时间戳和随机数生成唯一名称
	rand.Seed(time.Now().UnixNano())
	return strconv.FormatInt(time.Now().UnixMilli(), 10) + "-" + strconv.Itoa(rand.Intn(1_000_000_000)) + ext
}

// isAllowedMime 检查文件是否为允许的图片类型
// 检查方式：1. Content-Type 头 2. 文件内容魔数（更可靠）
//
// 参数：
//   - file: 上传的文件头信息
//
// 返回值：
//   - bool: 是否为允许的类型
func isAllowedMime(file *multipart.FileHeader) bool {
	ct := file.Header.Get("Content-Type")
	allowed := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if allowed[ct] {
		return true
	}
	// 通过文件内容检测 MIME 类型（更可靠，可防止伪造 Content-Type）
	f, err := file.Open()
	if err != nil {
		return false
	}
	defer f.Close()
	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	detected := http.DetectContentType(buf[:n])
	return allowed[detected]
}
