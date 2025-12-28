package indexer

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"skills-server/internal/models"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type Indexer struct {
	dataDir string
	skills  []models.Skill
	mu      sync.RWMutex
}

func New(dataDir string) *Indexer {
	return &Indexer{
		dataDir: dataDir,
		skills:  make([]models.Skill, 0),
	}
}

// Start 启动索引器和文件监听
func (i *Indexer) Start() error {
	// 初始扫描
	if err := i.Scan(); err != nil {
		return err
	}

	// 启动监听 (简单的实现，生产环境可能需要更健壮的 Debounce)
	go i.watch()
	return nil
}

// Scan 遍历目录构建索引
func (i *Indexer) Scan() error {
	var newSkills []models.Skill

	err := filepath.WalkDir(i.dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 忽略隐藏目录
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") && d.Name() != "." {
			return filepath.SkipDir
		}

		// 检查是否是 Skill 目录 (存在 SKILL.md)
		if !d.IsDir() && strings.EqualFold(d.Name(), "SKILL.md") {
			skillDir := filepath.Dir(path)
			skill, err := i.parseSkill(skillDir, path)
			if err == nil {
				newSkills = append(newSkills, skill)
			} else {
				fmt.Printf("Error parsing skill at %s: %v\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	i.mu.Lock()
	i.skills = newSkills
	i.mu.Unlock()
	fmt.Printf("Indexed %d skills from %s\n", len(newSkills), i.dataDir)
	return nil
}

// 解析单个 Skill
func (i *Indexer) parseSkill(dirPath, mdPath string) (models.Skill, error) {
	content, err := os.ReadFile(mdPath)
	if err != nil {
		return models.Skill{}, err
	}

	// 解析 Frontmatter
	// 格式: --- \n yaml \n ---
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		return models.Skill{}, fmt.Errorf("invalid SKILL.md format: missing frontmatter")
	}

	var meta models.SkillMetadata
	if err := yaml.Unmarshal(parts[1], &meta); err != nil {
		return models.Skill{}, fmt.Errorf("yaml parse error: %w", err)
	}

	// 基础校验
	if meta.Name == "" {
		// 如果 YAML 没写名字，用文件夹名兜底
		meta.Name = filepath.Base(dirPath)
	}

	relPath, _ := filepath.Rel(i.dataDir, dirPath)
	info, _ := os.Stat(mdPath)

	return models.Skill{
		SkillMetadata: meta,
		Path:          dirPath,
		RelativePath:  relPath,
		UpdatedAt:     info.ModTime(),
	}, nil
}

// watch 监听文件变动
func (i *Indexer) watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Failed to create watcher:", err)
		return
	}
	defer watcher.Close()

	if err := watcher.Add(i.dataDir); err != nil {
		fmt.Println("Failed to watch data dir:", err)
		return
	}

	// 简单的 Debounce 机制：收到事件后等待一小段时间再扫描
	var timer *time.Timer

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// 忽略 chmod 事件
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				continue
			}

			// 如果是目录创建，需要添加到 watcher (fsnotify 不会递归监听新目录)
			if event.Op&fsnotify.Create == fsnotify.Create {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					watcher.Add(event.Name)
				}
			}

			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(1*time.Second, func() {
				fmt.Println("File change detected, rescanning...")
				i.Scan()
			})

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Watcher error:", err)
		}
	}
}

// Search 搜索接口
func (i *Indexer) Search(query string) []models.Skill {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if query == "" {
		return i.skills
	}

	query = strings.ToLower(query)
	var results []models.Skill
	for _, s := range i.skills {
		if strings.Contains(strings.ToLower(s.Name), query) ||
			strings.Contains(strings.ToLower(s.Description), query) {
			results = append(results, s)
		}
	}
	return results
}

// GetByName 获取详情
func (i *Indexer) GetByName(name string) (models.SkillDetail, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	var target models.Skill
	found := false
	for _, s := range i.skills {
		if strings.EqualFold(s.Name, name) {
			target = s
			found = true
			break
		}
	}

	if !found {
		return models.SkillDetail{}, fmt.Errorf("skill not found")
	}

	// 读取完整内容
	mdPath := filepath.Join(target.Path, "SKILL.md")
	content, err := os.ReadFile(mdPath)
	if err != nil {
		return models.SkillDetail{}, err
	}

	// 生成文件树
	tree := i.generateFileTree(target.Path)

	return models.SkillDetail{
		Skill:         target,
		ReadmeContent: string(content),
		FileTree:      tree,
	}, nil
}

func (i *Indexer) generateFileTree(root string) string {
	var builder strings.Builder
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		// 过滤逻辑
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") && d.Name() != "." {
				return filepath.SkipDir
			}
			if d.Name() == "__pycache__" || d.Name() == "node_modules" {
				return filepath.SkipDir
			}
		}

		rel, _ := filepath.Rel(root, path)
		if rel == "." {
			return nil
		}

		depth := strings.Count(rel, string(os.PathSeparator))
		indent := strings.Repeat("  ", depth)

		marker := ""
		if d.IsDir() {
			marker = "/"
		}

		builder.WriteString(fmt.Sprintf("%s- %s%s\n", indent, d.Name(), marker))
		return nil
	})
	return builder.String()
}

// GetPath 获取物理路径（用于下载）
func (i *Indexer) GetPath(name string) (string, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	for _, s := range i.skills {
		if strings.EqualFold(s.Name, name) {
			return s.Path, nil
		}
	}
	return "", fmt.Errorf("not found")
}
