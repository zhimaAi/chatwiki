// Copyright © 2016- 2025 Sesame Network Technology all right reserved

package php

import (
	"context"
	"github.com/zhimaAi/go_tools/logs"
	"os/exec"
	"sync"
	"time"

	"github.com/roadrunner-server/pool/ipc/pipe"
	"github.com/roadrunner-server/pool/payload"
	"github.com/roadrunner-server/pool/pool"
	staticPool "github.com/roadrunner-server/pool/pool/static_pool"
	"github.com/roadrunner-server/pool/process"
	"go.uber.org/zap"
)

type PhpPoolConfig struct {
	IdleTimeout          time.Duration // 自动销毁前的空闲超时时间
	Ctx                  context.Context
	RoadrunnerPoolConfig *pool.Config // RoadRunner worker pool的原始配置
	EnvVars              []string     // 传递给 PHP worker 进程的环境变量数组
	QueueSize            uint64       // 队列大小，0 表示使用默认值
}

type PhpPool struct {
	config    *PhpPoolConfig   // 配置
	log       *zap.Logger      // 日志
	mu        sync.RWMutex     // 保护 realPool 和 idleTimer
	realPool  *staticPool.Pool // php worker pool 原始配置
	idleTimer *time.Timer      // 空闲计时器
}

func NewPhpPool(cfg *PhpPoolConfig) *PhpPool {
	return &PhpPool{
		config: cfg,
		log:    zap.L().Named("php_pool"),
	}
}

// Exec
// 1. 确保池已创建（如果未创建）
// 2. 停止任何挂起的空闲销毁计时器
// 3. 执行任务
// 4. 重置空闲销毁计时器
func (lp *PhpPool) Exec(ctx context.Context, p *payload.Payload, stopCh chan struct{}) (chan *staticPool.PExec, error) {
	// 1. 停止空闲计时器
	lp.stopIdleTimer()

	// 2. 获取或创建真实的pool
	currentPool, err := lp.getOrCreatePool()
	if err != nil {
		return nil, err
	}

	// 3. 执行任务
	re, err := currentPool.Exec(ctx, p, stopCh)

	// 4. 重置空闲计时器
	lp.resetIdleTimer()

	if err != nil {
		return nil, err
	}
	return re, nil
}

// Destroy 允许外部显式销毁池
func (lp *PhpPool) Destroy(ctx context.Context) {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	lp.destroyInternal(ctx)
}

// getOrCreatePool 是一个线程安全的函数，用于获取或创建真实的pool
func (lp *PhpPool) getOrCreatePool() (*staticPool.Pool, error) {
	// 使用读锁检查池是否已存在
	lp.mu.RLock()
	if lp.realPool != nil {
		lp.mu.RUnlock()
		return lp.realPool, nil
	}
	lp.mu.RUnlock()

	// 获取写锁来创建池
	lp.mu.Lock()
	defer lp.mu.Unlock()

	// 可能在等待写锁时，其他 goroutine 已创建了池
	if lp.realPool != nil {
		return lp.realPool, nil
	}

	// 创建命令工厂函数，使用配置的环境变量
	cmdFactory := func(command []string) *exec.Cmd {
		cmd := exec.CommandContext(lp.config.Ctx, command[0], command[1:]...)
		cmd.Env = make([]string, 0, len(lp.config.EnvVars))
		cmd.Env = append(cmd.Env, lp.config.EnvVars...)
		process.IsolateProcess(cmd) // 隔离进程
		return cmd
	}

	// 使用配置的日志创建 PipeFactory
	pipeFactory := pipe.NewPipeFactory(lp.log)

	// 构建 PoolOpts
	var poolOpts []staticPool.Options
	if lp.config.QueueSize > 0 {
		poolOpts = append(poolOpts, staticPool.WithQueueSize(lp.config.QueueSize))
	}

	p, err := staticPool.NewPool(
		lp.config.Ctx,
		cmdFactory,
		pipeFactory,
		lp.config.RoadrunnerPoolConfig,
		lp.log,
		poolOpts...,
	)
	if err != nil {
		logs.Error("create pool failed: %v", err.Error())
		return nil, err
	}

	lp.realPool = p
	return lp.realPool, nil
}

// stopIdleTimer 停止当前的空闲计时器
func (lp *PhpPool) stopIdleTimer() {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.idleTimer != nil {
		lp.idleTimer.Stop()
		lp.idleTimer = nil
	}
}

// resetIdleTimer 重置空闲计时器
func (lp *PhpPool) resetIdleTimer() {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	// 重置前先停止旧的
	if lp.idleTimer != nil {
		lp.idleTimer.Stop()
	}

	if lp.config.IdleTimeout <= 0 {
		return // 如果超时时间为0或更小，则禁用空闲销毁
	}

	lp.idleTimer = time.AfterFunc(lp.config.IdleTimeout, lp.destroyIdle)
}

// destroyIdle 是 time.AfterFunc 的回调
// 它会获取锁并调用内部销毁方法
func (lp *PhpPool) destroyIdle() {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	lp.destroyInternal(lp.config.Ctx)
}

// destroyInternal 实际执行销毁操作
func (lp *PhpPool) destroyInternal(ctx context.Context) {
	if lp.realPool == nil {
		return
	}

	// 销毁池
	lp.realPool.Destroy(ctx)
	lp.realPool = nil

	// 确保计时器也被清理
	if lp.idleTimer != nil {
		lp.idleTimer.Stop()
		lp.idleTimer = nil
	}
}
