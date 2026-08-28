package server

import (
	"net/http"

	"hmimg-server-go/internal/config"
	"hmimg-server-go/internal/handlers"
	"hmimg-server-go/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewRouter 创建并配置 Gin 路由引擎
// 设置所有 API 路由、中间件和 CORS 配置
//
// 参数：
//   - db: GORM 数据库连接实例
//   - cfg: 应用配置对象
//
// 返回值：
//   - *gin.Engine: 配置好的 Gin 路由引擎
func NewRouter(db *gorm.DB, cfg config.Config) *gin.Engine {
	// 创建默认的 Gin 引擎（包含 Logger 和 Recovery 中间件）
	r := gin.New()
	r.Use(gin.Recovery())

	// 配置 CORS（跨域资源共享）
	// 允许所有来源、所有常用 HTTP 方法和指定的头部
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // 允许所有来源（开发环境适用，生产环境建议限制）
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // 不允许携带凭证（Cookie 等）
	}))

	// 设置文件上传大小限制为 32MB
	r.MaxMultipartMemory = 32 << 20

	// ============================================================
	// API 路由组 /api
	// ============================================================
	api := r.Group("/api")
	{
		// 实例化处理器（注入数据库连接和配置）
		authHandler := handlers.AuthHandler{DB: db, Cfg: cfg}
		settingsHandler := handlers.SettingsHandler{DB: db}
		storageHandler := handlers.StorageHandler{DB: db, Cfg: cfg, DBDriver: cfg.DBDriver, UploadDir: cfg.UploadDir}

		// --------------------------------------------------------
		// 公开接口（无需认证）
		// --------------------------------------------------------

		// POST /api/login - 用户登录
		// 请求体：{"username": "...", "password": "..."}
		// 响应：{"token": "...", "user": {"id": ..., "username": "...", "role": "..."}}
		api.POST("/login", authHandler.Login)

		// POST /api/register - 用户注册
		// 请求体：{"username": "...", "password": "...", "email": "...", "phone": "..."}
		// 响应：{"message": "Registered successfully"}
		// 注意：是否允许注册由系统设置 allow_registration 控制
		api.POST("/register", authHandler.Register)

		// GET /api/settings/public - 获取公开设置（无需认证）
		// 响应：{"allow_registration": false, "website_title": "...", "default_language": "..."}
		api.GET("/settings/public", settingsHandler.GetPublic)

		// GET /api/files/*path - 获取上传的文件（静态文件服务）
		// 用于直接访问用户上传的图片等文件
		api.GET("/files/*path", storageHandler.GetUploadedFile)

		// --------------------------------------------------------
		// 受保护接口（需要登录认证）
		// --------------------------------------------------------
		protected := api.Group("")
		protected.Use(middleware.RequireAuth(cfg)) // 应用认证中间件
		{
			// PUT /api/admin/update - 更新当前用户个人信息
			// 可更新：username, email, phone, password
			protected.PUT("/admin/update", authHandler.UpdateProfile)
			protected.GET("/admin/users", middleware.RequireAdmin(), authHandler.ListUsers)
			protected.PUT("/admin/users/:id/role", middleware.RequireAdmin(), authHandler.UpdateUserRole)
			protected.DELETE("/admin/users/:id", middleware.RequireAdmin(), authHandler.DeleteUser)

			// --------------------------------------------------------
			// 管理员专属接口（需认证 + 管理员角色）
			// --------------------------------------------------------

			// GET /api/settings - 获取所有设置项（需管理员权限）
			protected.GET("/settings", middleware.RequireAdmin(), settingsHandler.GetAll)

			// PUT /api/settings - 更新系统设置（需管理员权限）
			// 可更新：max_users, allow_registration, website_title, default_language
			protected.PUT("/settings", middleware.RequireAdmin(), settingsHandler.Update)

			// --------------------------------------------------------
			// 相册管理接口（需认证）
			// --------------------------------------------------------

			// POST /api/albums - 创建相册
			protected.POST("/albums", storageHandler.CreateAlbum)

			// GET /api/albums - 获取所有相册（按 ID 倒序）
			protected.GET("/albums", storageHandler.GetAlbums)

			// GET /api/albums/:id - 获取单个相册详情
			protected.GET("/albums/:id", storageHandler.GetAlbum)

			// PUT /api/albums/:id - 更新相册（仅创建者或管理员）
			protected.PUT("/albums/:id", storageHandler.UpdateAlbum)

			// DELETE /api/albums/:id - 删除相册及其所有图片（仅创建者或管理员）
			protected.DELETE("/albums/:id", storageHandler.DeleteAlbum)

			// --------------------------------------------------------
			// 收藏集管理接口（需认证）
			// --------------------------------------------------------

			// POST /api/collections - 创建收藏集
			protected.POST("/collections", storageHandler.CreateCollection)

			// GET /api/collections - 获取所有收藏集（按 ID 倒序）
			protected.GET("/collections", storageHandler.GetCollections)

			// GET /api/collections/:id - 获取收藏集详情（包括包含的条目）
			protected.GET("/collections/:id", storageHandler.GetCollection)

			// GET /api/collections/:id/random - 从收藏集中随机获取一张图片
			// 查询参数：type=json（默认）或 type=redirect
			protected.GET("/collections/:id/random", storageHandler.GetRandomFromCollection)

			// PUT /api/collections/:id - 更新收藏集（仅创建者或管理员）
			protected.PUT("/collections/:id", storageHandler.UpdateCollection)

			// DELETE /api/collections/:id - 删除收藏集（仅创建者或管理员）
			protected.DELETE("/collections/:id", storageHandler.DeleteCollection)

			// POST /api/collections/add - 添加相册或收藏集到收藏集
			// 请求体：{"collectionId": 1, "itemType": "album", "itemName": "相册名"}
			protected.POST("/collections/add", storageHandler.AddToCollection)

			// --------------------------------------------------------
			// 图片管理接口（需认证）
			// --------------------------------------------------------

			// POST /api/upload - 上传图片（支持批量，最多 20 张）
			// 表单字段：albumId（相册ID）, images（文件数组）
			// 注意：需要特殊处理，Content-Type 为 multipart/form-data
			protected.POST("/upload", storageHandler.UploadImages)

			// GET /api/images - 获取图片列表
			// 查询参数：albumId（可选，筛选指定相册）
			protected.GET("/images", storageHandler.GetImages)

			// GET /api/images/:id - 获取单张图片详情
			protected.GET("/images/:id", storageHandler.GetImage)

			// PUT /api/images/:id - 更新图片信息（仅上传者或管理员）
			// 目前仅支持更新原始文件名
			protected.PUT("/images/:id", storageHandler.UpdateImage)
			protected.DELETE("/images/:id", storageHandler.DeleteImage)
		}
	}

	// 处理根路径请求
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/" || path == "" {
			// 根路径返回 API 标识字符串
			c.String(http.StatusOK, "HMiMG API Server")
			return
		}
		// 其他未匹配路径返回 404
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	})

	return r
}
