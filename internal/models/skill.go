package models

import "time"

// SkillMetadata 对应 SKILL.md 中的 Frontmatter
type SkillMetadata struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	Tags        []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	Author      string   `yaml:"author,omitempty" json:"author,omitempty"`
	Version     string   `yaml:"version,omitempty" json:"version,omitempty"`
}

// Skill 代表一个完整的技能对象
type Skill struct {
	SkillMetadata
	Path         string    `json:"-"` // 服务器本地绝对路径，不通过 API 暴露
	RelativePath string    `json:"path"` // 相对 data 目录的路径
	UpdatedAt    time.Time `json:"updated_at"`
}

// SkillDetail 包含 Skill 的完整信息，用于详情页接口
type SkillDetail struct {
	Skill
	ReadmeContent string `json:"readme"`    // SKILL.md 的内容
	FileTree      string `json:"file_tree"` // 文件树结构文本
}
