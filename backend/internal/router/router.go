package router

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/kuonji/blog/internal/controller"
	"github.com/kuonji/blog/internal/middleware"
)

func Setup(mode string) *gin.Engine {
	gin.SetMode(mode)
	r := gin.New()

	// 中间件
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// 静态文件服务（上传的图片）
	// 使用绝对路径确保无论从哪里启动都能找到 uploads 目录
	uploadDir := "uploads"
	if exe, err := os.Executable(); err == nil {
		absDir := filepath.Join(filepath.Dir(exe), uploadDir)
		if _, err := os.Stat(absDir); err == nil {
			uploadDir = absDir
		}
	}
	r.Static("/uploads", uploadDir)

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 公开接口
		public := api.Group("")
		{
			// 文章
			public.GET("/articles", controller.GetArticles)
			public.GET("/articles/:id", controller.GetArticle)

			// 分类
			public.GET("/categories", controller.GetCategories)
			public.GET("/categories/:id", controller.GetCategory)

			// 标签
			public.GET("/tags", controller.GetTags)
			public.GET("/tags/:id", controller.GetTag)

			// 系列
			public.GET("/series", controller.GetSeries)
			public.GET("/series/:id", controller.GetSeriesDetail)

			// 评论
			public.GET("/articles/:id/comments", controller.GetComments)
			public.POST("/articles/:id/comments", controller.CreateComment)

			// 搜索
			public.GET("/search", controller.Search)
		}

		// 管理接口（需要认证）
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth())
		{
			// 图片上传
			admin.POST("/upload", controller.UploadImage)
			// 文章管理
			admin.POST("/articles", controller.CreateArticle)
			admin.PUT("/articles/:id", controller.UpdateArticle)
			admin.DELETE("/articles/:id", controller.DeleteArticle)
			admin.PUT("/articles/:id/status", controller.UpdateArticleStatus)

			// 分类管理
			admin.POST("/categories", controller.CreateCategory)
			admin.PUT("/categories/:id", controller.UpdateCategory)
			admin.DELETE("/categories/:id", controller.DeleteCategory)

			// 标签管理
			admin.POST("/tags", controller.CreateTag)
			admin.PUT("/tags/:id", controller.UpdateTag)
			admin.DELETE("/tags/:id", controller.DeleteTag)

			// 系列管理
			admin.POST("/series", controller.CreateSeries)
			admin.PUT("/series/:id", controller.UpdateSeries)
			admin.DELETE("/series/:id", controller.DeleteSeries)

			// 评论管理
			admin.PUT("/comments/:id/status", controller.UpdateCommentStatus)
			admin.DELETE("/comments/:id", controller.DeleteComment)

			// 统计
			admin.GET("/stats", controller.GetStats)
		}

		// 认证
		auth := api.Group("/auth")
		{
			auth.POST("/login", controller.Login)
		}
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
