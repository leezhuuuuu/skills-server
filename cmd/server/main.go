package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"skills-server/internal/config"
	"skills-server/internal/handlers"
	"skills-server/internal/indexer"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed all:web_dist
var webDist embed.FS

func main() {
	cfg := config.Load()

	// 初始化索引
	idx := indexer.New(cfg.DataDir)
	if err := idx.Start(); err != nil {
		fmt.Printf("Warning: Failed to scan data dir: %v. Make sure %s exists.\n", err, cfg.DataDir)
	}

	h := handlers.New(idx)

	r := gin.Default()

	// CORS
	r.Use(cors.Default())

	// API Routes
	api := r.Group("/api/v1")
	{
		api.GET("/skills", h.ListSkills)
		api.GET("/skills/:name", h.GetSkill)
		api.GET("/download/:name", h.DownloadSkill)
	}

	// LLM Routes
	r.GET("/skill.md", h.GetGuideMD)
	r.GET("/skill/:name_with_ext", func(c *gin.Context) {
		// 这里的路由稍微有点 tricky，为了兼容 /skill/name.md 和 前端的 /skill/name
		name := c.Param("name_with_ext")
		if strings.HasSuffix(name, ".md") {
			h.GetSkillMD(c)
			return
		}
		// 如果不是 .md，交给前端路由处理（Fallthrough）
		c.Next()
	})

	// Frontend Static Files
	// 提取嵌入文件系统的子目录
	staticFS, _ := fs.Sub(webDist, "web_dist")
	r.StaticFS("/assets", http.FS(staticFS)) // 这主要服务 assets 目录

	// SPA 路由处理 (所有找不到的路由都返回 index.html)
	r.NoRoute(func(c *gin.Context) {
		// 避免 API 404 返回 HTML
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		c.FileFromFS("index.html", http.FS(staticFS))
	})

	fmt.Printf("Skills Server running on :%s, serving %s\n", cfg.Port, cfg.DataDir)
	r.Run(":" + cfg.Port)
}
