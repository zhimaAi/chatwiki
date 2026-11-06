// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import (
	"maps"
	"sync"

	"github.com/zhimaAi/go_tools/msql"
)

// ManualTransferState belongs to one request and is shared by its nested calls.
// Only a successful transfer publishes a result; the first success wins.
type ManualTransferState struct {
	transferMu  sync.Mutex
	mu          sync.RWMutex
	content     string
	transferred bool
	done        chan struct{}
	robot       msql.Params // original channel robot, preserved across nested robot changes
	source      ManualTransferSource
	sourceSet   bool
}

// ManualTransferSource identifies the customer-facing entry, independently of
// the robot and conversation currently executing a nested call.
type ManualTransferSource struct {
	Robot          msql.Params
	AdminUserId    int
	AppType        string
	AppInfo        msql.Params
	WechatappAppid string
	Openid         string
	SessionId      int
}

func NewManualTransferState(robot msql.Params) *ManualTransferState {
	return &ManualTransferState{done: make(chan struct{}), robot: maps.Clone(robot)}
}

// SetSource is called after the entry has resolved its session. Nested callers
// share the state but cannot replace the entry with their own session.
func (s *ManualTransferState) SetSource(params *ChatRequestParam, sessionID int) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.sourceSet {
		return
	}
	s.source = ManualTransferSource{
		Robot: maps.Clone(s.robot), AdminUserId: params.AdminUserId,
		AppType: params.AppType, AppInfo: maps.Clone(params.AppInfo),
		WechatappAppid: params.WechatappAppid, Openid: params.Openid, SessionId: sessionID,
	}
	s.sourceSet = true
}

func (s *ManualTransferState) Source() (ManualTransferSource, bool) {
	if s == nil {
		return ManualTransferSource{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	source := s.source
	source.Robot, source.AppInfo = maps.Clone(source.Robot), maps.Clone(source.AppInfo)
	return source, s.sourceSet
}

func (s *ManualTransferState) SourceRobot() msql.Params {
	if s == nil {
		return nil
	}
	return s.robot
}

func (s *ManualTransferState) Result() (string, bool) {
	if s == nil {
		return "", false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.content, s.transferred
}

func (s *ManualTransferState) Done() <-chan struct{} {
	if s == nil {
		return nil
	}
	return s.done
}

// Transfer serializes manual requests so a sibling cannot transfer again after
// the first success. A failed attempt leaves later attempts available.
func (s *ManualTransferState) Transfer(content string, run func() error) error {
	s.transferMu.Lock()
	defer s.transferMu.Unlock()
	if _, ok := s.Result(); ok {
		return nil
	}
	if err := run(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.content, s.transferred = content, true
	close(s.done)
	return nil
}
