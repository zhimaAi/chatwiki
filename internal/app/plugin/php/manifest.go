// Copyright © 2016- 2025 Sesame Network Technology all right reserved

package php

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Resource struct {
	Memory     uint64 `json:"memory"`
	Permission any    `json:"permission"`
}

// PluginManifest 插件清单结构体
type PluginManifest struct {
	Name       string   `json:"name"`
	Title      string   `json:"title"`
	Version    string   `json:"version"`
	Type       string   `json:"type"`
	Compatible string   `json:"compatible"`
	Resource   Resource `json:"resource"`
	HasLoaded  bool     `json:"has_loaded"` // 是否已加载
}

// GetPluginManifest 获取插件清单信息
func GetPluginManifest(name string) (*PluginManifest, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	pluginDir := filepath.Join(workDir, "php", "plugins", name)
	manifestPath := filepath.Join(pluginDir, "manifest.json")

	// 读取 manifest.json 文件
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest for plugin %s: %w", name, err)
	}

	var manifest PluginManifest
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to parse manifest for plugin %s: %w", name, err)
	}

	return &manifest, nil
}
