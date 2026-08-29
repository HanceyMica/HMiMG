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
	"hmimg-server-go/internal/dbstate"
	"hmimg-server-go/internal/httpx"
	"hmimg-server-go/internal/middleware"
	"hmimg-server-go/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StorageHandler 相册/合集/图片存储处理器
type StorageHandler struct {
	Cfg       config.Config
	DBDriver  string
	UploadDir string
}

// db 获取当前数据库连接（来自 dbstate）
func (h StorageHandler) db() *gorm.DB {
	return dbstate.DB()
}

type createAlbumRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h StorageHandler) CreateAlbum(c *gin.Context) {
	var req createAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var existing models.Album
	if err := h.db().First(&existing, "name = ?", req.Name).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Album name already exists"})
		return
	}
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)
	album := models.Album{Name: req.Name, Description: req.Description, CreatedBy: &userID}
	if err := h.db().Create(&album).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create album"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": album.ID, "name": album.Name, "description": album.Description})
}

func (h StorageHandler) GetAlbums(c *gin.Context) {
	var albums []models.Album
	if err := h.db().Order("id DESC").Find(&albums).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load albums"})
		return
	}
	out := make([]gin.H, 0, len(albums))
	for _, a := range albums {
		out = append(out, albumToJSON(a))
	}
	c.JSON(http.StatusOK, out)
}

func (h StorageHandler) GetAlbum(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	var album models.Album
	if err := h.db().First(&album, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}
	c.JSON(http.StatusOK, albumToJSON(album))
}

type updateAlbumRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

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
	var album models.Album
	if err := h.db().First(&album, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}
	roleAny, _ := c.Get(middleware.ContextRoleKey)
	role, _ := roleAny.(string)
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)
	if role != "admin" && (album.CreatedBy == nil || *album.CreatedBy != userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	if req.Name != "" && req.Name != album.Name {
		var existing models.Album
		if err := h.db().First(&existing, "name = ?", req.Name).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Album name already exists"})
			return
		}
		album.Name = req.Name
	}
	album.Description = req.Description
	if err := h.db().Save(&album).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update album"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Album updated", "id": id})
}

func (h StorageHandler) DeleteAlbum(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	var album models.Album
	if err := h.db().First(&album, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}
	roleAny, _ := c.Get(middleware.ContextRoleKey)
	role, _ := roleAny.(string)
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)
	if role != "admin" && (album.CreatedBy == nil || *album.CreatedBy != userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var images []models.Image
	_ = h.db().Where("album_id = ?", id).Find(&images).Error
	for _, img := range images {
		_ = os.Remove(filepath.Join(h.UploadDir, img.Path))
	}

	if err := h.db().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("album_id = ?", id).Delete(&models.Image{}).Error; err != nil {
			return err
		}
		if err := tx.Where("item_type = ? AND item_id = ?", "album", id).Delete(&models.CollectionItem{}).Error; err != nil {
			return err
		}
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

type createCollectionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h StorageHandler) CreateCollection(c *gin.Context) {
	var req createCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var existing models.Collection
	if err := h.db().First(&existing, "name = ?", req.Name).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Collection name already exists"})
		return
	}
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)
	col := models.Collection{Name: req.Name, Description: req.Description, CreatedBy: &userID}
	if err := h.db().Create(&col).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create collection"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": col.ID, "name": col.Name, "description": col.Description})
}

func (h StorageHandler) GetCollections(c *gin.Context) {
	var cols []models.Collection
	if err := h.db().Order("id DESC").Find(&cols).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load collections"})
		return
	}
	out := make([]gin.H, 0, len(cols))
	for _, col := range cols {
		out = append(out, collectionToJSON(col))
	}
	c.JSON(http.StatusOK, out)
}

func (h StorageHandler) GetCollection(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	var col models.Collection
	if err := h.db().First(&col, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	var items []models.CollectionItem
	if err := h.db().Where("collection_id = ?", id).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load collection items"})
		return
	}

	children := make([]gin.H, 0)
	if len(items) > 0 {
		itemType := items[0].ItemType
		ids := make([]uint32, 0, len(items))
		for _, it := range items {
			ids = append(ids, it.ItemID)
		}
		if itemType == "album" {
			var albums []models.Album
			_ = h.db().Where("id IN ?", ids).Find(&albums).Error
			for _, a := range albums {
				j := albumToJSON(a)
				j["type"] = "album"
				children = append(children, j)
			}
		} else if itemType == "collection" {
			var cols2 []models.Collection
			_ = h.db().Where("id IN ?", ids).Find(&cols2).Error
			for _, c2 := range cols2 {
				j := collectionToJSON(c2)
				j["type"] = "collection"
				children = append(children, j)
			}
		}
	}

	out := collectionToJSON(col)
	out["children"] = children
	c.JSON(http.StatusOK, out)
}

type updateCollectionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type updateImageRequest struct {
	OriginalName    string `json:"original_name"`
	OriginalNameAlt string `json:"originalName"`
}

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
	if err := h.db().First(&col, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}
	roleAny, _ := c.Get(middleware.ContextRoleKey)
	role, _ := roleAny.(string)
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)
	if role != "admin" && (col.CreatedBy == nil || *col.CreatedBy != userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	if req.Name != "" && req.Name != col.Name {
		var existing models.Collection
		if err := h.db().First(&existing, "name = ?", req.Name).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Collection name already exists"})
			return
		}
		col.Name = req.Name
	}
	col.Description = req.Description
	if err := h.db().Save(&col).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update collection"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Collection updated", "id": id})
}

func (h StorageHandler) DeleteCollection(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	var col models.Collection
	if err := h.db().First(&col, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}
	roleAny, _ := c.Get(middleware.ContextRoleKey)
	role, _ := roleAny.(string)
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)
	if role != "admin" && (col.CreatedBy == nil || *col.CreatedBy != userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	if err := h.db().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("item_type = ? AND item_id = ?", "collection", id).Delete(&models.CollectionItem{}).Error; err != nil {
			return err
		}
		if err := tx.Where("collection_id = ?", id).Delete(&models.CollectionItem{}).Error; err != nil {
			return err
		}
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

type addToCollectionRequest struct {
	CollectionID uint32 `json:"collectionId"`
	ItemType     string `json:"itemType"`
	ItemID       uint32 `json:"itemId"`
	ItemName     string `json:"itemName"`
}

func (h StorageHandler) AddToCollection(c *gin.Context) {
	var req addToCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.CollectionID == 0 || (req.ItemID == 0 && req.ItemName == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if req.ItemType != "album" && req.ItemType != "collection" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item type"})
		return
	}

	roleAny, _ := c.Get(middleware.ContextRoleKey)
	role, _ := roleAny.(string)
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)
	owns := func(owner *uint32) bool {
		return role == "admin" || (owner != nil && *owner == userID)
	}

	var collection models.Collection
	if err := h.db().First(&collection, "id = ?", req.CollectionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}
	if !owns(collection.CreatedBy) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var itemID uint32
	if req.ItemType == "album" {
		var album models.Album
		var itemErr error
		if req.ItemID != 0 {
			itemErr = h.db().First(&album, "id = ?", req.ItemID).Error
		} else {
			itemErr = h.db().First(&album, "name = ?", req.ItemName).Error
		}
		if itemErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		if !owns(album.CreatedBy) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		itemID = album.ID
	} else {
		var col models.Collection
		var itemErr error
		if req.ItemID != 0 {
			itemErr = h.db().First(&col, "id = ?", req.ItemID).Error
		} else {
			itemErr = h.db().First(&col, "name = ?", req.ItemName).Error
		}
		if itemErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Target collection not found"})
			return
		}
		if col.ID == req.CollectionID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add collection to itself"})
			return
		}
		if !owns(col.CreatedBy) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		itemID = col.ID
	}

	var existingOne models.CollectionItem
	if err := h.db().First(&existingOne, "collection_id = ?", req.CollectionID).Error; err == nil {
		if existingOne.ItemType != req.ItemType {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Collection cannot contain mixed types"})
			return
		}
	}

	var exists models.CollectionItem
	if err := h.db().First(&exists, "collection_id = ? AND item_type = ? AND item_id = ?", req.CollectionID, req.ItemType, itemID).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item already in collection"})
		return
	}

	toCreate := models.CollectionItem{CollectionID: req.CollectionID, ItemType: req.ItemType, ItemID: itemID}
	if err := h.db().Create(&toCreate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Added successfully"})
}

func (h StorageHandler) GetRandomFromCollection(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	returnType := c.DefaultQuery("type", "json")

	var col models.Collection
	if err := h.db().First(&col, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	albumIDs, err := h.getAlbumIDsFromCollection(id, map[uint32]bool{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get albums"})
		return
	}
	if len(albumIDs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No albums in this collection"})
		return
	}

	var image models.Image
	q := h.db().Where("album_id IN ?", albumIDs)
	dialect := strings.ToLower(h.DBDriver)
	if dialect == "postgres" || dialect == "pg" || dialect == "postgresql" {
		q = q.Order("RANDOM()")
	} else {
		q = q.Order("RAND()")
	}
	if err := q.First(&image).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No images found in this collection"})
		return
	}

	if returnType == "redirect" {
		base := httpx.Origin(c.Request, h.Cfg.TrustProxy)
		c.Redirect(http.StatusFound, base+"/api/files/"+url.PathEscape(image.Path))
		return
	}
	c.JSON(http.StatusOK, imageToJSON(image))
}

func (h StorageHandler) getAlbumIDsFromCollection(collectionID uint32, visited map[uint32]bool) ([]uint32, error) {
	if visited[collectionID] {
		return nil, nil
	}
	visited[collectionID] = true

	var items []models.CollectionItem
	if err := h.db().Where("collection_id = ?", collectionID).Find(&items).Error; err != nil {
		return nil, err
	}
	var albumIDs []uint32
	for _, it := range items {
		if it.ItemType == "album" {
			albumIDs = append(albumIDs, it.ItemID)
			continue
		}
		if it.ItemType == "collection" {
			childAlbumIDs, err := h.getAlbumIDsFromCollection(it.ItemID, visited)
			if err != nil {
				return nil, err
			}
			albumIDs = append(albumIDs, childAlbumIDs...)
		}
	}
	return albumIDs, nil
}

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

	var album models.Album
	if err := h.db().First(&album, "id = ?", albumID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

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

	insertedIDs := make([]uint32, 0, len(files))
	failedFiles := make([]gin.H, 0)
	var firstFilename string
	for _, file := range files {
		if ok := isAllowedMime(file); !ok {
			failedFiles = append(failedFiles, gin.H{"filename": file.Filename, "error": "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed."})
			continue
		}
		filename := uniqueFilename(h.UploadDir, file.Filename)
		dst := filepath.Join(h.UploadDir, filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			failedFiles = append(failedFiles, gin.H{"filename": file.Filename, "error": "Failed to save file"})
			continue
		}
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
		if err := h.db().Create(&img).Error; err != nil {
			_ = os.Remove(dst)
			failedFiles = append(failedFiles, gin.H{"filename": file.Filename, "error": "Failed to record image"})
			continue
		}
		if firstFilename == "" {
			firstFilename = filename
		}
		insertedIDs = append(insertedIDs, img.ID)
	}

	if len(insertedIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded successfully", "failures": failedFiles})
		return
	}

	if album.CoverImage == nil && firstFilename != "" {
		cover := firstFilename
		album.CoverImage = &cover
		_ = h.db().Save(&album).Error
	}

	c.JSON(http.StatusOK, gin.H{"ids": insertedIDs, "count": len(insertedIDs), "failures": failedFiles})
}

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

func (h StorageHandler) GetImages(c *gin.Context) {
	albumIDStr := c.Query("albumId")
	q := h.db().Model(&models.Image{}).Order("id DESC")
	if albumIDStr != "" {
		albumID64, err := strconv.ParseUint(albumIDStr, 10, 32)
		if err == nil {
			q = q.Where("album_id = ?", uint32(albumID64))
		}
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load images"})
		return
	}

	page, pageSize := paginationParams(c)
	var images []models.Image
	if err := q.Limit(pageSize).Offset((page - 1) * pageSize).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load images"})
		return
	}
	out := make([]gin.H, 0, len(images))
	for _, img := range images {
		out = append(out, imageToJSON(img))
	}
	c.JSON(http.StatusOK, gin.H{"items": out, "total": total, "page": page, "pageSize": pageSize})
}

func paginationParams(c *gin.Context) (int, int) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	if err != nil || pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}
	return page, pageSize
}

func (h StorageHandler) GetImage(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	var img models.Image
	if err := h.db().Order("id DESC").First(&img, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}
	var album models.Album
	albumErr := h.db().First(&album, "id = ?", img.AlbumID).Error
	out := imageToJSON(img)
	if albumErr != nil || album.Name == "" {
		out["album_name"] = nil
	} else {
		out["album_name"] = album.Name
	}
	c.JSON(http.StatusOK, out)
}

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

	name := strings.TrimSpace(req.OriginalName)
	if name == "" {
		name = strings.TrimSpace(req.OriginalNameAlt)
	}
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image name is required"})
		return
	}

	var img models.Image
	if err := h.db().First(&img, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	roleAny, _ := c.Get(middleware.ContextRoleKey)
	role, _ := roleAny.(string)
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)
	if role != "admin" && (img.UploadedBy == nil || *img.UploadedBy != userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	img.OriginalName = name
	if err := h.db().Save(&img).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update image"})
		return
	}

	c.JSON(http.StatusOK, imageToJSON(img))
}

func (h StorageHandler) DeleteImage(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var img models.Image
	if err := h.db().First(&img, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	roleAny, _ := c.Get(middleware.ContextRoleKey)
	role, _ := roleAny.(string)
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)
	if role != "admin" && (img.UploadedBy == nil || *img.UploadedBy != userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.db().Delete(&models.Image{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image"})
		return
	}
	_ = os.Remove(filepath.Join(h.UploadDir, img.Path))

	var album models.Album
	if err := h.db().First(&album, "id = ?", img.AlbumID).Error; err == nil {
		if album.CoverImage != nil && *album.CoverImage == img.Filename {
			var next models.Image
			if err := h.db().Where("album_id = ?", img.AlbumID).Order("id ASC").First(&next).Error; err == nil {
				cover := next.Filename
				album.CoverImage = &cover
			} else {
				album.CoverImage = nil
			}
			_ = h.db().Save(&album).Error
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted"})
}

func parseUintParam(c *gin.Context, key string) (uint32, bool) {
	idStr := c.Param(key)
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return 0, false
	}
	return uint32(id64), true
}

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
	if filepath.IsAbs(cleaned) || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", errors.New("invalid path")
	}

	return cleaned, nil
}

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

func uniqueFilename(uploadDir, original string) string {
	ext := filepath.Ext(original)
	if ext == "" {
		ext = ".bin"
	}
	for {
		filename := strconv.FormatInt(time.Now().UnixMilli(), 10) + "-" + strconv.Itoa(rand.Intn(1_000_000_000)) + ext
		if _, err := os.Stat(filepath.Join(uploadDir, filename)); os.IsNotExist(err) {
			return filename
		}
	}
}

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
