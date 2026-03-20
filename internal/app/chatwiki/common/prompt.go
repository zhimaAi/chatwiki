// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

type PromptItem struct {
	Priority float64 `json:"priority"`
	Subject  string  `json:"subject"`
	Describe string  `json:"describe"`
}

type StructPrompt struct {
	Role        PromptItem   `json:"role"`
	Task        PromptItem   `json:"task"`
	Constraints PromptItem   `json:"constraints"`
	Skill       PromptItem   `json:"skill"`
	Output      PromptItem   `json:"output"`
	Tone        PromptItem   `json:"tone"`
	Custom      []PromptItem `json:"custom"`
}

func GetEmptyPromptStruct(lang string) StructPrompt {
	return StructPrompt{
		Role:        PromptItem{Subject: i18n.Show(lang, `prompt_struct_role_subject`)},
		Task:        PromptItem{Subject: i18n.Show(lang, `prompt_struct_task_subject`)},
		Constraints: PromptItem{Subject: i18n.Show(lang, `prompt_struct_constraints_subject`)},
		Skill:       PromptItem{Subject: i18n.Show(lang, `prompt_struct_skill_subject`)},
		Output:      PromptItem{Subject: i18n.Show(lang, `prompt_struct_output_subject`)},
		Tone:        PromptItem{Subject: i18n.Show(lang, `prompt_struct_tone_subject`)},
		Custom:      []PromptItem{},
	}
}

func GetDefaultPromptStruct(lang string) string {
	structPrompt := GetEmptyPromptStruct(lang)
	structPrompt.Role.Describe = i18n.Show(lang, `prompt_default_role_describe`)
	structPrompt.Task.Describe = i18n.Show(lang, `prompt_default_task_describe`)
	structPrompt.Constraints.Describe = i18n.Show(lang, `prompt_default_constraints_describe`)
	structPrompt.Skill.Describe = `` // Skills are empty by default
	structPrompt.Output.Describe = i18n.Show(lang, `prompt_default_output_describe`)
	structPrompt.Tone.Describe = i18n.Show(lang, `prompt_default_tone_describe`)
	return tool.JsonEncodeNoError(structPrompt)
}

func CheckPromptConfig(lang string, promptType int, promptStruct string) (string, error) {
	structPrompt := StructPrompt{}
	err := tool.JsonDecodeUseNumber(promptStruct, &structPrompt)
	switch promptType {
	case define.PromptTypeCustom:
		//nothing to do
	case define.PromptTypeStruct:
		if err != nil {
			return ``, errors.New(i18n.Show(lang, `prompt_struct_config_error`))
		}
		for _, item := range structPrompt.Custom {
			if len(item.Describe) > 0 && len(item.Subject) == 0 {
				return ``, errors.New(i18n.Show(lang, `prompt_struct_subject_unnamed`))
			}
		}
	default:
		return ``, errors.New(i18n.Show(lang, `prompt_type_param_error`, promptType))
	}
	structPrompt = SetDdefaultFieldSubject(lang, structPrompt) // Prohibit modification of default field subjects
	if structPrompt.Custom == nil {
		structPrompt.Custom = make([]PromptItem, 0)
	}
	return tool.JsonEncodeNoError(structPrompt), nil
}

func SetDdefaultFieldSubject(lang string, structPrompt StructPrompt) StructPrompt {
	empty := GetEmptyPromptStruct(lang)
	structPrompt.Role.Subject = empty.Role.Subject
	structPrompt.Task.Subject = empty.Task.Subject
	structPrompt.Constraints.Subject = empty.Constraints.Subject
	structPrompt.Skill.Subject = empty.Skill.Subject
	structPrompt.Output.Subject = empty.Output.Subject
	structPrompt.Tone.Subject = empty.Tone.Subject
	return structPrompt
}

func BuildPromptStruct(lang string, promptType int, prompt, promptStruct string) string {
	switch promptType {
	case define.PromptTypeStruct:
		sp := StructPrompt{}
		if err := tool.JsonDecodeUseNumber(promptStruct, &sp); err != nil {
			logs.Error(`promptStruct:%s,err:%v`, promptStruct, err)
		}
		sp = SetDdefaultFieldSubject(lang, sp) // Prohibit modification of default field subjects
		mds := make([]string, 0)
		if len(sp.Role.Describe) > 0 {
			mds = append(mds, fmt.Sprintf("## %s\n%s", sp.Role.Subject, sp.Role.Describe))
		}
		if len(sp.Task.Describe) > 0 {
			mds = append(mds, fmt.Sprintf("## %s\n%s", sp.Task.Subject, sp.Task.Describe))
		}
		if len(sp.Constraints.Describe) > 0 {
			mds = append(mds, fmt.Sprintf("## %s\n%s", sp.Constraints.Subject, sp.Constraints.Describe))
		}
		if len(sp.Skill.Describe) > 0 {
			mds = append(mds, fmt.Sprintf("## %s\n%s", sp.Skill.Subject, sp.Skill.Describe))
		}
		if len(sp.Output.Describe) > 0 {
			mds = append(mds, fmt.Sprintf("## %s\n%s", sp.Output.Subject, sp.Output.Describe))
		}
		if len(sp.Tone.Describe) > 0 {
			mds = append(mds, fmt.Sprintf("## %s\n%s", sp.Tone.Subject, sp.Tone.Describe))
		}
		for _, item := range sp.Custom {
			if len(item.Subject) > 0 && len(item.Describe) > 0 {
				mds = append(mds, fmt.Sprintf("## %s\n%s", item.Subject, item.Describe))
			}
		}
		return strings.Join(mds, "\n")
	default:
		return prompt
	}
}

func FormatSystemPrompt(lang string, prompt string, list []msql.Params) (string, string) {
	output := fmt.Sprintf("# %s\n%s", i18n.Show(lang, `prompt_system`), prompt)
	knowledges := make([]string, 0)
	for idx, one := range list {
		var images []string
		if err := tool.JsonDecode(one[`images`], &images); err != nil {
			logs.Error(err.Error())
		}
		var imgs string
		for _, image := range images {
			ext := strings.ToLower(strings.TrimLeft(filepath.Ext(image), `.`))
			if tool.InArrayString(ext, define.VideoAllowExt) {
				imgs += fmt.Sprintf("\n![video](%s)", image)
			} else {
				if !IsUrl(image) { // if it is not a url add the domain name
					image = define.Config.WebService[`image_domain`] + image
				}
				imgs += fmt.Sprintf("\n![image](%s)", image)
			}
		}
		if cast.ToInt(one[`type`]) == define.ParagraphTypeNormal {
			knowledges = append(knowledges, fmt.Sprintf("## %s\n%s%s", i18n.Show(lang, `prompt_library_section`, idx+1), one[`content`], imgs))
		} else {
			var similarQuestions []string
			if err := tool.JsonDecode(one[`similar_questions`], &similarQuestions); err != nil {
				logs.Error(err.Error())
			}
			var similar string
			if len(similarQuestions) > 0 {
				similar = fmt.Sprintf("\n%s：%s", i18n.Show(lang, `prompt_similar_questions`), strings.Join(similarQuestions, `/`))
			}
			knowledges = append(knowledges, fmt.Sprintf("## %s\n%s:%s%s\n%s:%s%s", i18n.Show(lang, `prompt_library_section`, idx+1),
				i18n.Show(lang, `prompt_question`), one[`question`], similar, i18n.Show(lang, `prompt_answer`), one[`answer`], imgs))
		}
	}
	var libraryOutput string
	if len(knowledges) > 0 {
		output += fmt.Sprintf("\n# %s\n%s", i18n.Show(lang, `prompt_library`), strings.Join(knowledges, "\n"))
		libraryOutput = fmt.Sprintf("# %s\n%s", i18n.Show(lang, `prompt_library`), strings.Join(knowledges, "\n"))
	}
	return UnifyLineBreak(output), UnifyLineBreak(libraryOutput) // Unify line break processing
}

func UnifyLineBreak(content string) string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	content = strings.ReplaceAll(content, "\n", "\r\n")
	return content
}

func CreatePromptByAi(lang string, demand string, adminUserId, modelConfigId int, useModel string) (string, error) {
	messages := []adaptor.ZhimaChatCompletionMessage{
		{Role: `system`, Content: define.PromptDefaultCreatePrompt},
		{Role: `user`, Content: demand},
	}
	chatResp, _, err := RequestChat(lang, adminUserId, cast.ToString(adminUserId), nil, lib_define.AppYunPc,
		modelConfigId, useModel, messages, nil, 0.5, 2000)
	if err != nil {
		logs.Error(err.Error())
		return ``, err
	}
	chatResp.Result, _ = strings.CutPrefix(chatResp.Result, "```json")
	chatResp.Result, _ = strings.CutSuffix(chatResp.Result, "```")
	promptStruct, err := CheckPromptConfig(lang, define.PromptTypeStruct, chatResp.Result)
	if err != nil {
		return ``, fmt.Errorf(`%s`, chatResp.Result)
	}
	return promptStruct, nil
}

// getSessionChatVariables loads all prompt variables (including system & workflow variables)
// and returns them with the unified placeholder regexp used across the chat flow.
func getSessionChatVariables(sessionId int) ([]ChatVariable, *regexp.Regexp, error) {
	if sessionId <= 0 {
		return nil, nil, nil
	}

	chatPromptVariablesStr, err := msql.Model(`chat_ai_session`, define.Postgres).
		Where(`id`, cast.ToString(sessionId)).
		Value(`chat_prompt_variables`)
	if err != nil {
		logs.Error(err.Error())
		return nil, nil, err
	}
	if len(chatPromptVariablesStr) == 0 {
		return nil, nil, nil
	}

	chatPromptVariables := make([]ChatVariable, 0)
	if err := tool.JsonDecode(chatPromptVariablesStr, &chatPromptVariables); err != nil {
		logs.Error(err.Error())
		return nil, nil, err
	}

	re, err := regexp.Compile(`【chat_variable:[a-zA-Z_]+】`)
	if err != nil {
		logs.Error(err.Error())
		return nil, nil, err
	}

	return chatPromptVariables, re, nil
}

func ReplaceChatVariables(lang string, sessionId int, prompt *string, promptStruct *string) {
	chatPromptVariables, re, err := getSessionChatVariables(sessionId)
	if err != nil || len(chatPromptVariables) == 0 || re == nil {
		return
	}
	//prompt
	ReplaceChatVariable(lang, prompt, chatPromptVariables, re)
	//struct prompt
	sp := StructPrompt{}
	if err := tool.JsonDecodeUseNumber(*promptStruct, &sp); err != nil {
		logs.Error(`promptStruct:%s,err:%v`, promptStruct, err)
		return
	}
	ReplaceChatVariable(lang, &sp.Role.Describe, chatPromptVariables, re)
	ReplaceChatVariable(lang, &sp.Task.Describe, chatPromptVariables, re)
	ReplaceChatVariable(lang, &sp.Constraints.Describe, chatPromptVariables, re)
	ReplaceChatVariable(lang, &sp.Skill.Describe, chatPromptVariables, re)
	ReplaceChatVariable(lang, &sp.Output.Describe, chatPromptVariables, re)
	ReplaceChatVariable(lang, &sp.Tone.Describe, chatPromptVariables, re)
	if len(sp.Custom) > 0 {
		for _, item := range sp.Custom {
			ReplaceChatVariable(lang, &item.Describe, chatPromptVariables, re)
		}
	}
	*promptStruct = tool.JsonEncodeNoError(sp)
}

// ReplaceMetaSearchChatVariables replaces chat variable placeholders in robot metadata filter config.
// It only affects string-type conditions and keeps the original configuration when:
// - metadata search is disabled, or
// - there is no active chat session / variables.
func ReplaceMetaSearchChatVariables(lang string, sessionId int, robot *msql.Params) {
	if cast.ToInt((*robot)[`meta_search_switch`]) != define.MetaSearchSwitchOn {
		return
	}

	raw := strings.TrimSpace((*robot)[`meta_search_condition_list`])
	if raw == "" || raw == "{}" || raw == "null" {
		return
	}

	chatPromptVariables, re, err := getSessionChatVariables(sessionId)
	if err != nil || len(chatPromptVariables) == 0 || re == nil {
		return
	}

	conds := make([]MetaSearchCondition, 0)
	if err := tool.JsonDecode(raw, &conds); err != nil {
		logs.Error(err.Error())
		return
	}
	if len(conds) == 0 {
		return
	}

	for i := range conds {
		// Support all types (string, number, time) as they may contain placeholders
		v := conds[i].Value
		ReplaceChatVariable(lang, &v, chatPromptVariables, re)
		conds[i].Value = v
	}

	replacedResult := tool.JsonEncodeNoError(conds)
	(*robot)[`meta_search_condition_list`] = replacedResult
}

func ReplaceChatVariable(lang string, str *string, chatPromptVariables []ChatVariable, re *regexp.Regexp) {
	fullMatches := re.FindAllString(*str, -1)
	replaces := map[string]string{}
	for _, match := range fullMatches {
		replaces[match] = ``
	}
	for _, item := range chatPromptVariables {
		if item.VariableType == VariableTypeCheckboxSwitch {
			if cast.ToInt(item.Value) == 1 {
				replaces[`【chat_variable:`+item.VariableKey+`】`] = i18n.Show(lang, `chat_variable_selected`)
			} else {
				replaces[`【chat_variable:`+item.VariableKey+`】`] = i18n.Show(lang, `chat_variable_unselected`)
			}
		} else {
			replaces[`【chat_variable:`+item.VariableKey+`】`] = cast.ToString(item.Value)
		}
	}
	for k, v := range replaces {
		*str = strings.ReplaceAll(*str, k, v)
	}
}
