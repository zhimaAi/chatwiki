// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

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
	MultiNode  bool     `json:"multiNode"` // 是否支持多节点运行，默认为单节点运行，即不支持多节点运行，如果为 true 则支持多节点运行，即支持多节点运行，此时需要配置
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
