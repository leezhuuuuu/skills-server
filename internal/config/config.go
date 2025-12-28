package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Port    string
	DataDir string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dataDir := os.Getenv("SKILLS_DATA_DIR")
	if dataDir == "" {
		// 默认为当前工作目录下的 data
		cwd, _ := os.Getwd()
		dataDir = filepath.Join(cwd, "data")
	}

	return &Config{
		Port:    port,
		DataDir: dataDir,
	}
}
