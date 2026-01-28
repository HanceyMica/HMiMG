package server

import (
	"net/http"
	"os"
	"path/filepath"

	"hmimg-server-go/internal/config"
	"hmimg-server-go/internal/handlers"
	"hmimg-server-go/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, cfg config.Config) *gin.Engine {
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

	api := r.Group("/api")
	{
		authHandler := handlers.AuthHandler{DB: db, Cfg: cfg}
		settingsHandler := handlers.SettingsHandler{DB: db}
		storageHandler := handlers.StorageHandler{DB: db, Cfg: cfg, DBDriver: cfg.DBDriver, UploadDir: cfg.UploadDir}

		api.POST("/login", authHandler.Login)
		api.POST("/register", authHandler.Register)

		api.GET("/settings/public", settingsHandler.GetPublic)

		protected := api.Group("")
		protected.Use(middleware.RequireAuth(cfg))
		{
			protected.PUT("/admin/update", authHandler.UpdateProfile)

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
		}
	}

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// If path is root or empty, show welcome message (like Node implementation)
		if path == "/" || path == "" {
			c.String(http.StatusOK, "HMiMG API Server")
			return
		}

		fullPath := filepath.Join(cfg.UploadDir, path)
		if fi, err := os.Stat(fullPath); err == nil && !fi.IsDir() {
			c.File(fullPath)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	})

	return r
}
