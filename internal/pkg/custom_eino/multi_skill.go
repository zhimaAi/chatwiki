// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package custom_eino

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk/middlewares/skill"
)

type LayeredSkillBackend struct {
	Backends []skill.Backend        // priority order: private first, public second
	Filter   func(name string) bool // returning true will filter out the skill
}

func (b *LayeredSkillBackend) List(ctx context.Context) ([]skill.FrontMatter, error) {
	seen := map[string]struct{}{}
	var out []skill.FrontMatter
	for _, backend := range b.Backends {
		items, err := backend.List(ctx)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			if b.Filter != nil && b.Filter(item.Name) {
				continue
			}
			if _, ok := seen[item.Name]; ok {
				continue
			}
			seen[item.Name] = struct{}{}
			out = append(out, item)
		}
	}
	return out, nil
}

func (b *LayeredSkillBackend) Get(ctx context.Context, name string) (skill.Skill, error) {
	if b.Filter != nil && b.Filter(name) {
		return skill.Skill{}, fmt.Errorf("skill not found: %s", name)
	}
	for _, backend := range b.Backends {
		items, err := backend.List(ctx)
		if err != nil {
			return skill.Skill{}, err
		}
		for _, item := range items {
			if item.Name == name {
				return backend.Get(ctx, name)
			}
		}
	}
	return skill.Skill{}, fmt.Errorf("skill not found: %s", name)
}
