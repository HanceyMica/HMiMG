package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"hmimg-server-go/internal/config"
	"hmimg-server-go/internal/dbstate"
	"hmimg-server-go/internal/handlers"
	"hmimg-server-go/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewRouter 构建全部路由
// /api/install/* 始终挂载（安装完成后接口内部返回 404）
// 其余 /api/* 经 requireInstalled 中间件：未安装时统一 503
func NewRouter(cfg config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	r.MaxMultipartMemory = 32 << 20

	installHandler := handlers.InstallHandler{Cfg: cfg}
	install := r.Group("/api/install")
	{
		install.GET("/status", installHandler.Status)
		install.POST("/database", installHandler.Database)
		install.POST("/admin", installHandler.Admin)
		install.POST("/site", installHandler.Site)
	}

	api := r.Group("/api")
	api.Use(requireInstalled())
	{
		authHandler := handlers.AuthHandler{Cfg: cfg}
		settingsHandler := handlers.SettingsHandler{}
		storageHandler := handlers.StorageHandler{Cfg: cfg, DBDriver: cfg.DBDriver, UploadDir: cfg.UploadDir}

		api.POST("/login", authHandler.Login)
		api.POST("/register", authHandler.Register)

		api.GET("/settings/public", settingsHandler.GetPublic)
		api.GET("/files/*path", storageHandler.GetUploadedFile)

		protected := api.Group("")
		protected.Use(middleware.RequireAuth(cfg))
		{
			protected.PUT("/admin/update", authHandler.UpdateProfile)
			protected.GET("/admin/users", middleware.RequireAdmin(), authHandler.ListUsers)
			protected.PUT("/admin/users/:id/role", middleware.RequireAdmin(), authHandler.UpdateUserRole)
			protected.DELETE("/admin/users/:id", middleware.RequireAdmin(), authHandler.DeleteUser)

			protected.GET("/settings", middleware.RequireAdmin(), settingsHandler.GetAll)
			protected.PUT("/settings", middleware.RequireAdmin(), settingsHandler.Update)

			protected.POST("/albums", storageHandler.CreateAlbum)
			protected.GET("/albums", storageHandler.GetAlbums)
			protected.GET("/albums/:id", storageHandler.GetAlbum)
			protected.PUT("/albums/:id", storageHandler.UpdateAlbum)
			protected.DELETE("/albums/:id", storageHandler.DeleteAlbum)

			protected.POST("/collections", storageHandler.CreateCollection)
			protected.GET("/collections", storageHandler.GetCollections)
			protected.GET("/collections/:id", storageHandler.GetCollection)
			protected.GET("/collections/:id/random", storageHandler.GetRandomFromCollection)
			protected.PUT("/collections/:id", storageHandler.UpdateCollection)
			protected.DELETE("/collections/:id", storageHandler.DeleteCollection)
			protected.POST("/collections/add", storageHandler.AddToCollection)

			protected.POST("/upload", storageHandler.UploadImages)
			protected.GET("/images", storageHandler.GetImages)
			protected.GET("/images/:id", storageHandler.GetImage)
			protected.PUT("/images/:id", storageHandler.UpdateImage)
			protected.DELETE("/images/:id", storageHandler.DeleteImage)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// API 路径未命中路由：保持 JSON 404
		if strings.HasPrefix(path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		// 前后端不分离：托管前端静态文件（FRONTEND_DIR 配置时启用）
		if cfg.FrontendDir != "" {
			if serveStatic(c, cfg.FrontendDir, path) {
				return
			}
		}

		if path == "/" || path == "" {
			c.String(http.StatusOK, "HMiMG API Server")
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	})

	return r
}

// serveStatic 从前端目录提供静态文件，未命中时回退 index.html（SPA 路由）
// 返回 false 表示目录不可用（按无前端处理）
func serveStatic(c *gin.Context, dir, urlPath string) bool {
	clean := filepath.Clean("/" + urlPath)
	full := filepath.Join(dir, clean)

	// 防路径穿越：解析后必须仍位于前端目录内（目录本身视为根路径）
	if full != dir && !strings.HasPrefix(full, dir+string(os.PathSeparator)) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return true
	}

	if info, err := os.Stat(full); err == nil && !info.IsDir() {
		c.File(full)
		return true
	}

	index := filepath.Join(dir, "index.html")
	if _, err := os.Stat(index); err != nil {
		return false
	}
	c.File(index)
	return true
}

// requireInstalled 未安装时拦截全部业务接口
// 安装向导完成（dbstate.SetInstalled）后同进程立即放行，无需重启
func requireInstalled() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !dbstate.Installed() {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "Not installed"})
			return
		}
		c.Next()
	}
}
