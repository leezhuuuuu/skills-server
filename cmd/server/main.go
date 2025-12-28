package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path"
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
	// Avoid redirect loops when serving SPA paths.
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	// CORS
	r.Use(cors.Default())

	// Frontend Static Files
	// 提取嵌入文件系统的子目录
	staticFS, err := fs.Sub(webDist, "web_dist")
	if err != nil {
		log.Fatalf("Failed to load embedded frontend: %v", err)
	}
	indexHTML, err := fs.ReadFile(staticFS, "index.html")
	if err != nil {
		log.Fatalf("Failed to read embedded index.html: %v", err)
	}

	serveIndex := func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	}

	r.GET("/", func(c *gin.Context) {
		serveIndex(c)
	})

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
		// 如果不是 .md，交给前端路由处理（SPA）
		serveIndex(c)
	})

	// Serve embedded assets with explicit Content-Type.
	serveAsset := func(c *gin.Context) {
		p := strings.TrimPrefix(c.Param("filepath"), "/")
		if p == "" {
			c.Status(http.StatusNotFound)
			return
		}
		clean := path.Clean("/" + p)
		if strings.HasPrefix(clean, "/..") {
			c.Status(http.StatusBadRequest)
			return
		}
		filePath := path.Join("assets", strings.TrimPrefix(clean, "/"))
		data, err := fs.ReadFile(staticFS, filePath)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		contentType := mime.TypeByExtension(path.Ext(filePath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		c.Data(http.StatusOK, contentType, data)
	}
	r.GET("/assets/*filepath", serveAsset)
	r.HEAD("/assets/*filepath", serveAsset)

	// SPA 路由处理 (所有找不到的路由都返回 index.html)
	r.NoRoute(func(c *gin.Context) {
		// 避免 API 404 返回 HTML
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		serveIndex(c)
	})

	fmt.Printf("Skills Server running on :%s, serving %s\n", cfg.Port, cfg.DataDir)
	r.Run(":" + cfg.Port)
}
