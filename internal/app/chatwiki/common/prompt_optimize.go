// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"context"
	"errors"
	"sort"
	"strings"
	"time"
)

// Pick complete rounds without truncating or summarizing their contents.
func SelectPromptOptimizeHistory(rounds []define.PromptOptimizeRound, now time.Time) []define.PromptOptimizeRound {
	cutoff := now.Add(-define.PromptOptimizeHistoryWindow).Unix()
	selected := make([]define.PromptOptimizeRound, 0, len(rounds))
	for _, r := range rounds {
		if r.Status == define.PromptOptimizeStatusSuccess && r.CreatedAt >= cutoff {
			selected = append(selected, r)
		}
	}
	sort.Slice(selected, func(i, j int) bool {
		if selected[i].CreatedAt == selected[j].CreatedAt {
			return selected[i].ID > selected[j].ID
		}
		return selected[i].CreatedAt > selected[j].CreatedAt
	})
	if len(selected) > define.PromptOptimizeHistoryRounds {
		selected = selected[:define.PromptOptimizeHistoryRounds]
	}
	for i, j := 0, len(selected)-1; i < j; i, j = i+1, j-1 {
		selected[i], selected[j] = selected[j], selected[i]
	}
	return selected
}

// PromptOptimizeSystemPrompt is the only mode-dependent part of optimization.
// History selection and request assembly must remain identical for all modes.
func PromptOptimizeSystemPrompt(lang, scene, mode string) (string, error) {
	var key string
	switch scene + "/" + mode {
	case define.PromptOptimizeSceneConfig + "/" + define.PromptOptimizeModeOptimize:
		key = "prompt_optimize_system_config_optimize"
	case define.PromptOptimizeSceneConfig + "/" + define.PromptOptimizeModeRewrite:
		key = "prompt_optimize_system_config_rewrite"
	case define.PromptOptimizeSceneConversation + "/" + define.PromptOptimizeModeOptimize:
		key = "prompt_optimize_system_conversation_optimize"
	case define.PromptOptimizeSceneConversation + "/" + define.PromptOptimizeModeRewrite:
		key = "prompt_optimize_system_conversation_rewrite"
	default:
		return "", errors.New(i18n.Show(lang, "param_invalid", "scene/mode"))
	}
	return i18n.Show(lang, "prompt_optimize_system_base") + i18n.Show(lang, key) +
		"\n\n" + i18n.Show(lang, "prompt_optimize_output_format"), nil
}

func promptOptimizeHistoryReply(round define.PromptOptimizeRound) string {
	if round.Reply == "" {
		return round.Prompt // Older generated versions did not store an explanation.
	}
	if round.Prompt == "" {
		return round.Reply
	}
	return define.PromptOptimizePromptOpenTag + "\n" + round.Prompt + "\n" + define.PromptOptimizePromptCloseTag +
		"\n\n" + define.PromptOptimizeReplyOpenTag + "\n" + round.Reply + "\n" + define.PromptOptimizeReplyCloseTag
}

// Provider details are logged by the service, never exposed to clients.
func PromptOptimizeError(lang string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return errors.New(i18n.Show(lang, "prompt_optimize_timeout"))
	}
	if errors.Is(err, context.Canceled) {
		return errors.New(i18n.Show(lang, "prompt_optimize_cancelled"))
	}
	text := strings.ToLower(err.Error())
	for _, marker := range []string{"context_length", "context window", "maximum context", "too many tokens", "input too long"} {
		if strings.Contains(text, marker) {
			return errors.New(i18n.Show(lang, "prompt_optimize_context_exceeded"))
		}
	}
	return errors.New(i18n.Show(lang, "prompt_optimize_generation_failed"))
}
