// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package php

import (
	"chatwiki/internal/pkg/lib_web"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/roadrunner-server/goridge/v3/pkg/frame"
	"github.com/roadrunner-server/pool/payload"
	"github.com/roadrunner-server/pool/pool"
)

const LambdaMode = "lambda"
const ConsumerMode = "consumer"
const WebMode = "web"

type PhpPoolWrap struct {
	StopCh       chan struct{} // 停止信号
	LambdaPool   *PhpPool
	ConsumerPool *PhpPool // todo
	WebPool      *PhpPool // todo
}

type MultiPool struct {
	PhpWorkerPath   string
	Mu              sync.Mutex
	PluginWorkerMap map[string]*PhpPoolWrap
}

func NewPhpPluginPool(phpPath string) *MultiPool {
	return &MultiPool{
		PhpWorkerPath:   phpPath,
		PluginWorkerMap: make(map[string]*PhpPoolWrap),
	}
}

// LoadPhpPlugin 加载php插件
// currentVersion 当前主版本号，如: "2.0.4"
func (p *MultiPool) LoadPhpPlugin(name string, currentVersion string, config []string) error {
	manifest, err := GetPluginManifest(name)
	if err != nil {
		return err
	}

	// 检查版本兼容性
	if err := checkVersionCompatibility(manifest.Compatible, currentVersion); err != nil {
		return err
	}

	p.Mu.Lock()
	defer p.Mu.Unlock()

	config = append(config, "RR_MODE="+LambdaMode)
	lambdaPool := NewPhpPool(&PhpPoolConfig{
		IdleTimeout: 600 * time.Second,
		Ctx:         context.Background(),
		RoadrunnerPoolConfig: &pool.Config{
			Debug:      false,
			Command:    []string{"php", p.PhpWorkerPath},
			NumWorkers: 1,
			MaxJobs:    10,
			Supervisor: &pool.SupervisorConfig{
				MaxWorkerMemory: manifest.Resource.Memory,
			},
			DynamicAllocatorOpts: &pool.DynamicAllocationOpts{
				MaxWorkers: 3,
			},
		},
		EnvVars:   config,
		QueueSize: 100,
	})
	p.PluginWorkerMap[name] = &PhpPoolWrap{
		LambdaPool: lambdaPool,
		StopCh:     make(chan struct{}),
	}

	// todo consumer pool

	// todo web pool

	return nil
}

// UnloadPhpPlugin 卸载插件
func (p *MultiPool) UnloadPhpPlugin(name string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if workers, ok := p.PluginWorkerMap[name]; ok {
		close(workers.StopCh)
		delete(p.PluginWorkerMap, name)
	}
}

func (p *MultiPool) ExecLambdaPhpPlugin(name string, title string, action string, params map[string]any) (*lib_web.Response, error) {
	p.Mu.Lock()
	workers, ok := p.PluginWorkerMap[name]
	p.Mu.Unlock()

	if !ok {
		return nil, fmt.Errorf("插件%s未启用", title)
	}

	// 构建请求参数
	req := map[string]any{
		"plugin": name,
		"action": action,
		"params": params,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal invoke params: %w", err)
	}

	// 创建 payload
	pld := &payload.Payload{
		Codec:   frame.CodecProto,
		Context: nil,
		Body:    body,
	}

	// 执行任务
	ctx := context.Background()
	resultCh, err := workers.LambdaPool.Exec(ctx, pld, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute php plugin: %w", err)
	}

	var r *payload.Payload
	select {
	case pl := <-resultCh:
		if pl.Error() != nil {
			return nil, fmt.Errorf("failed to execute php plugin: %w", pl.Error())
		}
		if pl.Payload().Flags&frame.STREAM != 0 {
			return nil, fmt.Errorf("streaming is not supported")
		}
		r = pl.Payload()
	default:
		return nil, fmt.Errorf("empty Response")
	}

	var resp lib_web.Response
	err = json.Unmarshal(r.Body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Response: %w", err)
	}
	return &resp, nil
}

// checkVersionCompatibility 检查版本兼容性
// compatibleStr 格式如: ">=2.0.4,<3.1.0" 或 ">=2.0.4"
// currentVersion 当前主版本号，如: "2.0.4"
func checkVersionCompatibility(compatibleStr, currentVersion string) error {
	if compatibleStr == "" {
		return nil
	}

	// 解析当前版本
	currentVer, err := semver.NewVersion(currentVersion)
	if err != nil {
		return fmt.Errorf("版本信息读取失败: %v", err.Error())
	}

	// 解析兼容性约束（支持多个约束，用逗号分隔）
	constraints := strings.Split(compatibleStr, ",")
	var constraintList []*semver.Constraints

	for _, constraintStr := range constraints {
		constraintStr = strings.TrimSpace(constraintStr)
		if constraintStr == "" {
			continue
		}

		constraint, err := semver.NewConstraint(constraintStr)
		if err != nil {
			return fmt.Errorf("版本信息设置有误 %s: %w", constraintStr, err)
		}
		constraintList = append(constraintList, constraint)
	}

	// 检查所有约束是否都满足
	for _, constraint := range constraintList {
		if !constraint.Check(currentVer) {
			return fmt.Errorf("版本号不兼容,当前主系统版本: %s 该插件仅能匹配 %s", currentVersion, compatibleStr)
		}
	}

	return nil
}
