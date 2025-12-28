package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"skills-server/internal/indexer"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Index *indexer.Indexer
}

func New(idx *indexer.Indexer) *Handler {
	return &Handler{Index: idx}
}

// API: List Skills (Search)
func (h *Handler) ListSkills(c *gin.Context) {
	query := c.Query("q")
	skills := h.Index.Search(query)
	c.JSON(http.StatusOK, gin.H{"skills": skills})
}

// API: Get Detail
func (h *Handler) GetSkill(c *gin.Context) {
	name := c.Param("name")
	detail, err := h.Index.GetByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}
	c.JSON(http.StatusOK, detail)
}

// API: Download ZIP
func (h *Handler) DownloadSkill(c *gin.Context) {
	name := c.Param("name")
	path, err := h.Index.GetPath(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", name))
	c.Header("Content-Type", "application/zip")

	// 实时压缩
	zw := zip.NewWriter(c.Writer)
	defer zw.Close()

	// 遍历目录并写入 ZIP
	// 注意：baseDir 是为了在 zip 中保持文件夹结构
	// 例如 pdf-toolkit/SKILL.md 而不是直接 SKILL.md
	baseDir := filepath.Base(path)

	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 过滤
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || info.Name() == "__pycache__") {
			return filepath.SkipDir
		}

		// 创建 zip header
		relPath, err := filepath.Rel(path, filePath)
		if err != nil {
			return err
		}

		// 根目录本身不写入
		if relPath == "." {
			return nil
		}

		// ZIP 内部路径
		zipPath := filepath.Join(baseDir, relPath)
		// 统一转为 forward slash (zip 标准)
		zipPath = strings.ReplaceAll(zipPath, "\\", "/")

		if info.IsDir() {
			zipPath += "/"
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = zipPath
		header.Method = zip.Deflate

		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	if err != nil {
		// 这里如果出错比较尴尬，因为 header 可能已经发了
		fmt.Println("Zip error:", err)
	}
}

// LLM: Global Guide
func (h *Handler) GetGuideMD(c *gin.Context) {
	md := `# Skills Registry

Welcome to the Skills Registry.

## Installation
To use skills from this registry, please install the MCP client:

` + "```bash\ncurl -LsSf https://astral.sh/uv/install.sh | sh\nuv tool install skills-mcp\n```" + `

## Available Skills
`
	skills := h.Index.Search("")
	for _, s := range skills {
		md += fmt.Sprintf("- **%s**: %s\n", s.Name, s.Description)
	}

	c.String(http.StatusOK, md)
}

// LLM: Skill Context
func (h *Handler) GetSkillMD(c *gin.Context) {
	// 尝试获取 name (兼容 /api/v1/skills/:name)
	name := c.Param("name")

	// 如果为空，尝试获取 name_with_ext (兼容 /skill/:name_with_ext)
	if name == "" {
		name = c.Param("name_with_ext")
	}

	// 简单处理：去掉后缀 .md
	name = strings.TrimSuffix(name, ".md")

	detail, err := h.Index.GetByName(name)
	if err != nil {
		c.String(http.StatusNotFound, "Skill not found")
		return
	}

	md := fmt.Sprintf("# %s\n\n> %s\n\n## Metadata\n- Version: %s\n- Author: %s\n\n## File Structure\n```\n%s\n```\n\n## Documentation\n\n%s",
		detail.Name, detail.Description, detail.Version, detail.Author, detail.FileTree, detail.ReadmeContent)

	c.String(http.StatusOK, md)
}
