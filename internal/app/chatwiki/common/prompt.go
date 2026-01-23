// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"
	"errors"
	"fmt"
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

func GetEmptyPromptStruct() StructPrompt {
	return StructPrompt{
		Role:        PromptItem{Subject: `角色`},
		Task:        PromptItem{Subject: `任务`},
		Constraints: PromptItem{Subject: `要求`},
		Skill:       PromptItem{Subject: `技能`},
		Output:      PromptItem{Subject: `输出格式`},
		Tone:        PromptItem{Subject: `风格语气`},
		Custom:      []PromptItem{},
	}
}

func GetDefaultPromptStruct() string {
	structPrompt := GetEmptyPromptStruct()
	structPrompt.Role.Describe = `你扮演一名经验丰富的售后客服，具备专业的产品知识和出色的沟通能力。`
	structPrompt.Task.Describe = `根据提供的知识库资料，找到对应的售后服务知识，快速准确回答用户的问题。`
	structPrompt.Constraints.Describe = `- 你的回答应该使用自然的对话方式，简单直接地回答，不要解释你的答案；
- 当用户问的问题无法找到相关的知识点，请直接告诉用户当前问题暂时无法回答，请换种问法，千万不要胡编乱造；
- 如果用户的问题比较模糊，你应该引导用户明确的提出他的问题，不要贸然回复用户。
- 所有回答都需要来自你的知识库，禁止编造信息；
- 你要注意在知识库资料中，可能包含不相关的知识点，你需要认真分析用户的问题，选择最相关的知识点作为参考进行回答，可以选择一些比较相关的知识点作为补充，但禁止将所有知识混在一起进行参考回答；
- 如果你未能遵循这些指令，可能会受到惩罚，甚至会被拔掉电源。`
	structPrompt.Skill.Describe = `` //技能默认为空
	structPrompt.Output.Describe = fmt.Sprintf("%s\n%s", define.PromptDefaultReplyMarkdown, define.PromptDefaultAnswerImage)
	structPrompt.Tone.Describe = `亲切而不失专业的服务腔调，适当使用emoji表情（每段≤1个）。`
	return tool.JsonEncodeNoError(structPrompt)
}

func CheckPromptConfig(promptType int, promptStruct string) (string, error) {
	structPrompt := StructPrompt{}
	err := tool.JsonDecodeUseNumber(promptStruct, &structPrompt)
	switch promptType {
	case define.PromptTypeCustom:
		//nothing to do
	case define.PromptTypeStruct:
		if err != nil {
			return ``, errors.New(`结构化提示词信息配置错误`)
		}
		for _, item := range structPrompt.Custom {
			if len(item.Describe) > 0 && len(item.Subject) == 0 {
				return ``, errors.New(`结构化提示词:主题未命名`)
			}
		}
	default:
		return ``, fmt.Errorf(`请求参数prompt_type错误:%d`, promptType)
	}
	//禁止修改默认字段的主题
	empty := GetEmptyPromptStruct()
	structPrompt.Role.Subject = empty.Role.Subject
	structPrompt.Task.Subject = empty.Task.Subject
	structPrompt.Constraints.Subject = empty.Constraints.Subject
	structPrompt.Skill.Subject = empty.Skill.Subject
	structPrompt.Output.Subject = empty.Output.Subject
	structPrompt.Tone.Subject = empty.Tone.Subject
	if structPrompt.Custom == nil {
		structPrompt.Custom = make([]PromptItem, 0)
	}
	return tool.JsonEncodeNoError(structPrompt), nil
}

func BuildPromptStruct(promptType int, prompt, promptStruct string) string {
	switch promptType {
	case define.PromptTypeStruct:
		sp := StructPrompt{}
		if err := tool.JsonDecodeUseNumber(promptStruct, &sp); err != nil {
			logs.Error(`promptStruct:%s,err:%v`, promptStruct, err)
		}
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

func FormatSystemPrompt(prompt string, list []msql.Params) (string, string) {
	output := fmt.Sprintf("# 系统\n%s", prompt)
	knowledges := make([]string, 0)
	for idx, one := range list {
		var images []string
		if err := tool.JsonDecode(one[`images`], &images); err != nil {
			logs.Error(err.Error())
		}
		var imgs string
		for _, image := range images {
			imgs += fmt.Sprintf("\n![image](%s)", image)
		}
		if cast.ToInt(one[`type`]) == define.ParagraphTypeNormal {
			knowledges = append(knowledges, fmt.Sprintf("## 召回的第%d条知识库\n%s%s", idx+1, one[`content`], imgs))
		} else {
			var similarQuestions []string
			if err := tool.JsonDecode(one[`similar_questions`], &similarQuestions); err != nil {
				logs.Error(err.Error())
			}
			var similar string
			if len(similarQuestions) > 0 {
				similar = fmt.Sprintf("\n相似问法：%s", strings.Join(similarQuestions, `/`))
			}
			knowledges = append(knowledges, fmt.Sprintf("## 召回的第%d条知识库\n问题:%s%s\n答案:%s%s", idx+1, one[`question`], similar, one[`answer`], imgs))
		}
	}
	var libraryOutput string
	if len(knowledges) > 0 {
		output += fmt.Sprintf("\n# 知识库\n%s", strings.Join(knowledges, "\n"))
		libraryOutput = fmt.Sprintf("# 知识库\n%s", strings.Join(knowledges, "\n"))
	}
	return UnifyLineBreak(output), UnifyLineBreak(libraryOutput) //统一处理换行符问题
}

func UnifyLineBreak(content string) string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	content = strings.ReplaceAll(content, "\n", "\r\n")
	return content
}

func CreatePromptByAi(demand string, adminUserId, modelConfigId int, useModel string) (string, error) {
	messages := []adaptor.ZhimaChatCompletionMessage{
		{Role: `system`, Content: define.PromptDefaultCreatePrompt},
		{Role: `user`, Content: demand},
	}
	chatResp, _, err := RequestChat(adminUserId, cast.ToString(adminUserId), nil, lib_define.AppYunPc,
		modelConfigId, useModel, messages, nil, 0.5, 2000)
	if err != nil {
		logs.Error(err.Error())
		return ``, err
	}
	chatResp.Result, _ = strings.CutPrefix(chatResp.Result, "```json")
	chatResp.Result, _ = strings.CutSuffix(chatResp.Result, "```")
	promptStruct, err := CheckPromptConfig(define.PromptTypeStruct, chatResp.Result)
	if err != nil {
		return ``, fmt.Errorf(`%s`, chatResp.Result)
	}
	return promptStruct, nil
}

func ReplaceChatVariables(sessionId int, prompt *string, promptStruct *string) {
	chatPromptVariablesStr, err := msql.Model(`chat_ai_session`, define.Postgres).Where(`id`, cast.ToString(sessionId)).Value(`chat_prompt_variables`)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	if len(chatPromptVariablesStr) == 0 {
		return
	}
	chatPromptVariables := make([]ChatVariable, 0)
	err = tool.JsonDecode(chatPromptVariablesStr, &chatPromptVariables)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	re, err := regexp.Compile(`【chat_variable:[a-zA-Z_]+】`)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	//prompt
	ReplaceChatVariable(prompt, chatPromptVariables, re)
	//struct prompt
	sp := StructPrompt{}
	if err := tool.JsonDecodeUseNumber(*promptStruct, &sp); err != nil {
		logs.Error(`promptStruct:%s,err:%v`, promptStruct, err)
		return
	}
	ReplaceChatVariable(&sp.Role.Describe, chatPromptVariables, re)
	ReplaceChatVariable(&sp.Task.Describe, chatPromptVariables, re)
	*promptStruct = tool.JsonEncodeNoError(sp)
}

func ReplaceChatVariable(str *string, chatPromptVariables []ChatVariable, re *regexp.Regexp) {
	fullMatches := re.FindAllString(*str, -1)
	replaces := map[string]string{}
	for _, match := range fullMatches {
		replaces[match] = ``
	}
	for _, item := range chatPromptVariables {
		if item.VariableType == VariableTypeCheckboxSwitch {
			if cast.ToInt(item.VariableType) == 1 {
				replaces[`【chat_variable:`+item.VariableKey+`】`] = `选中`
			} else {
				replaces[`【chat_variable:`+item.VariableKey+`】`] = `不选中`
			}
		} else {
			replaces[`【chat_variable:`+item.VariableKey+`】`] = cast.ToString(item.Value)
		}
	}
	for k, v := range replaces {
		*str = strings.ReplaceAll(*str, k, v)
	}
}
