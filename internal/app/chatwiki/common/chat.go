// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func BuildChatContextPair(openid string, robotId, dialogueId, curMsgId, contextPair int) []map[string]string {
	contextList := make([]map[string]string, 0)
	if len(openid) == 0 || dialogueId <= 0 || contextPair <= 0 {
		return contextList //no context required
	}
	m := msql.Model(`chat_ai_message`, define.Postgres).Where(`openid`, openid).
		Where(`robot_id`, cast.ToString(robotId)).Where(`dialogue_id`, cast.ToString(dialogueId)).
		Where(`is_valid_function_call = false`)
	if curMsgId > 0 { // Compatible with debug mode to get context
		m.Where(`id`, `<`, cast.ToString(curMsgId))
	}
	list, err := m.Order(`id desc`).Field(`id,content,is_customer`).Limit(contextPair * 4).Select()
	if err != nil {
		logs.Error(err.Error())
	}
	if len(list) == 0 {
		return contextList
	}
	//reverse
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	//foreach
	for i := 0; i < len(list)-1; i++ {
		if cast.ToInt(list[i][`is_customer`]) == define.MsgFromCustomer && cast.ToInt(list[i+1][`is_customer`]) == define.MsgFromRobot {
			contextList = append(contextList, map[string]string{`question`: list[i][`content`], `answer`: list[i+1][`content`]})
			i++ //skip answer
		}
	}
	//cut out
	if len(contextList) > contextPair {
		contextList = contextList[len(contextList)-contextPair:]
	}
	return contextList
}

func BuildOpenApiContent(params *define.ChatRequestParam, messages []adaptor.ZhimaChatCompletionMessage) []adaptor.ZhimaChatCompletionMessage {
	if params.AppType != lib_define.AppOpenApi || len(params.OpenApiContent) == 0 {
		return messages
	}
	var contents = make([]adaptor.ZhimaChatCompletionMessage, 0)
	err := tool.JsonDecode(params.OpenApiContent, &contents)
	if err != nil {
		logs.Error(err.Error())
		return messages
	}
	if len(contents) > 0 {
		messages = append(contents, messages...)
	}
	if define.IsDev {
		logs.Debug("%+v", messages)
	}
	return messages
}

func SendDefaultUnknownQuestionPrompt(params *define.ChatRequestParam, errmsg string, chanStream chan sse.Event, content *string) {
	chanStream <- sse.Event{Event: `error`, Data: `SYSERR:` + errmsg}
	code := `unknown`
	if ms := regexp.MustCompile(`ERROR\s+CODE:\s?(.*)`).FindStringSubmatch(errmsg); len(ms) > 1 {
		code = ms[1]
	}
	*content = i18n.Show(params.Lang, `gpt_error`, code)
	chanStream <- sse.Event{Event: `sending`, Data: *content}
}

func BuildLibraryChatRequestMessage(params *define.ChatRequestParam, curMsgId int64, dialogueId, sessionId int, debugLog *[]any) ([]adaptor.ZhimaChatCompletionMessage, []msql.Params, LibUseTime, error) {
	if len(params.Prompt) == 0 { //no custom is used
		prompt := params.Robot[`prompt`]
		promptStruct := params.Robot[`prompt_struct`]
		ReplaceChatVariables(params.Lang, sessionId, &prompt, &promptStruct)
		params.Prompt = BuildPromptStruct(params.Lang, cast.ToInt(params.Robot[`prompt_type`]), prompt, promptStruct)
	}
	if len(params.LibraryIds) == 0 || !CheckIds(params.LibraryIds) { //no custom is used
		params.LibraryIds = params.Robot[`library_ids`]
	}

	contextList := BuildChatContextPair(params.Openid, cast.ToInt(params.Robot[`id`]),
		dialogueId, int(curMsgId), cast.ToInt(params.Robot[`context_pair`]))

	//question optimize
	var questionopTime int64
	var optimizedQuestions []string
	if cast.ToBool(params.Robot[`enable_question_optimize`]) && len(params.LibraryIds) > 0 {
		var err error
		temp := time.Now()
		optimizedQuestions, err = GetOptimizedQuestions(params, contextList)
		questionopTime = time.Now().Sub(temp).Milliseconds()
		if err != nil {
			logs.Error(err.Error())
		}
	}

	//convert match
	list, libUseTime, err := GetMatchLibraryParagraphList(
		params.Lang,
		params.Openid,
		params.AppType,
		params.Question,
		optimizedQuestions,
		params.LibraryIds,
		cast.ToInt(params.Robot[`top_k`]),
		cast.ToFloat64(params.Robot[`similarity`]),
		cast.ToInt(params.Robot[`search_type`]),
		params.Robot,
	)
	libUseTime.QuestionOp = questionopTime
	if err != nil {
		return nil, nil, libUseTime, err
	}

	//part0:init messages
	messages := make([]adaptor.ZhimaChatCompletionMessage, 0)
	//part1:prompt
	roleType := define.PromptRoleTypeMap[cast.ToInt(params.Robot[`prompt_role_type`])]
	if cast.ToBool(params.Robot[`question_multiple_switch`]) {
		// When calling multimodal, ignore user-set prompts placed in user role, always place in system role
		roleType = define.PromptRoleTypeMap[define.PromptRoleTypeSystem]
	}
	prompt, libraryContent := FormatSystemPrompt(params.Lang, params.Prompt, list)
	if roleType == define.PromptRoleUser {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `system`, Content: libraryContent})
		*debugLog = append(*debugLog, map[string]string{`type`: `prompt`, `content`: libraryContent})
	} else {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: roleType, Content: prompt})
		*debugLog = append(*debugLog, map[string]string{`type`: `prompt`, `content`: prompt})
	}
	//part2:context_qa
	for i := range contextList {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: contextList[i][`question`]})
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `assistant`, Content: contextList[i][`answer`]})
		*debugLog = append(*debugLog, map[string]string{`type`: `context_qa`, `question`: contextList[i][`question`], `answer`: contextList[i][`answer`]})
	}
	//part3:question,prompt+question
	if roleType == define.PromptRoleUser {
		content := strings.Join([]string{params.Prompt, params.Question}, "\n\n")
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: content})
		*debugLog = append(*debugLog, map[string]string{`type`: `cur_question`, `content`: content})
	} else {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: params.Question})
		*debugLog = append(*debugLog, map[string]string{`type`: `cur_question`, `content`: params.Question})
	}
	return messages, list, libUseTime, nil
}

func BuildDirectChatRequestMessage(params *define.ChatRequestParam, curMsgId int64, dialogueId, sessionId int, debugLog *[]any) ([]adaptor.ZhimaChatCompletionMessage, error) {
	if len(params.Prompt) == 0 { //no custom is used
		prompt := params.Robot[`prompt`]
		promptStruct := params.Robot[`prompt_struct`]
		ReplaceChatVariables(params.Lang, sessionId, &prompt, &promptStruct)
		params.Prompt = BuildPromptStruct(params.Lang, cast.ToInt(params.Robot[`prompt_type`]), prompt, promptStruct)
	}

	//part0:init messages
	messages := make([]adaptor.ZhimaChatCompletionMessage, 0)
	//part1:prompt
	prompt, _ := FormatSystemPrompt(params.Lang, params.Prompt, nil)
	roleType := define.PromptRoleTypeMap[cast.ToInt(params.Robot[`prompt_role_type`])]
	if cast.ToBool(params.Robot[`question_multiple_switch`]) {
		// When calling multimodal, ignore user-set prompts in user role, always place in system role
		roleType = define.PromptRoleTypeMap[define.PromptRoleTypeSystem]
	}
	if roleType != define.PromptRoleUser {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: roleType, Content: prompt})
		*debugLog = append(*debugLog, map[string]string{`type`: `prompt`, `content`: prompt})
	}
	//part2:context_qa
	contextList := BuildChatContextPair(params.Openid, cast.ToInt(params.Robot[`id`]),
		dialogueId, int(curMsgId), cast.ToInt(params.Robot[`context_pair`]))
	for i := range contextList {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: contextList[i][`question`]})
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `assistant`, Content: contextList[i][`answer`]})
		*debugLog = append(*debugLog, map[string]string{`type`: `context_qa`, `question`: contextList[i][`question`], `answer`: contextList[i][`answer`]})
	}
	//part3:cur_question
	if roleType == define.PromptRoleUser {
		content := strings.Join([]string{prompt, params.Question}, "\n\n")
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: content})
		*debugLog = append(*debugLog, map[string]string{`type`: `cur_question`, `content`: content})
	} else {
		messages = append(messages, adaptor.ZhimaChatCompletionMessage{Role: `user`, Content: params.Question})
		*debugLog = append(*debugLog, map[string]string{`type`: `cur_question`, `content`: params.Question})
	}
	return messages, nil
}

func DisposeQuoteFile(adminUserId int, list []msql.Params) ([]msql.Params, string) {
	var fileSourceMap = make(map[string][]msql.Datas)
	if len(list) > 0 {
		for _, one := range list {
			var images []string
			if err := tool.JsonDecodeUseNumber(one[`images`], &images); err != nil {
				logs.Error(err.Error())
			}
			fileSourceMap[one[`file_id`]] = append(fileSourceMap[one[`file_id`]], msql.Datas{
				`admin_user_id`: adminUserId,
				`file_id`:       one[`file_id`],
				`paragraph_id`:  one[`id`],
				`word_total`:    one[`word_total`],
				`similarity`:    one[`similarity`],
				`title`:         one[`title`],
				`type`:          one[`type`],
				`content`:       one[`content`],
				`question`:      one[`question`],
				`answer`:        one[`answer`],
				`images`:        images,
				`create_time`:   tool.Time2Int(),
				`update_time`:   tool.Time2Int(),
			})
		}
	}
	quoteFile, ms := make([]msql.Params, 0), map[string]struct{}{}
	var quoteFileForSave = make([]msql.Params, len(quoteFile))
	for _, one := range list {
		if _, ok := ms[one[`file_id`]]; ok {
			continue //remove duplication
		}
		library, err := GetLibraryInfo(cast.ToInt(one[`library_id`]), cast.ToInt(one[`admin_user_id`]))
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		ms[one[`file_id`]] = struct{}{}
		quoteFile = append(quoteFile, msql.Params{
			`id`:                 one[`file_id`],
			`library_id`:         library[`id`],
			`library_name`:       library[`library_name`],
			`file_name`:          one[`file_name`],
			`answer_source_data`: tool.JsonEncodeNoError(fileSourceMap[one[`file_id`]]),
		})
		quoteFileForSave = append(quoteFileForSave, msql.Params{
			`id`:           one[`file_id`],
			`library_id`:   library[`id`],
			`library_name`: library[`library_name`],
			`file_name`:    one[`file_name`],
		})
	}
	return quoteFile, tool.JsonEncodeNoError(quoteFileForSave)
}

// CheckQaDirectReply checks if QA knowledge is hit and replies directly with the answer
func CheckQaDirectReply(list []msql.Params, robot msql.Params) (string, bool) {
	var fieldSwitch, fieldScore string
	switch cast.ToInt(robot[`chat_type`]) {
	case define.ChatTypeMixture:
		fieldSwitch, fieldScore = `mixture_qa_direct_reply_switch`, `mixture_qa_direct_reply_score`
	case define.ChatTypeLibrary:
		fieldSwitch, fieldScore = `library_qa_direct_reply_switch`, `library_qa_direct_reply_score`
	default:
		return ``, false
	}
	if len(list) > 0 && cast.ToBool(robot[fieldSwitch]) &&
		cast.ToInt(list[0][`type`]) != define.ParagraphTypeNormal &&
		len(list[0][`similarity`]) > 0 &&
		cast.ToFloat32(list[0][`similarity`]) >= cast.ToFloat32(robot[fieldScore]) {
		content := list[0][`answer`]
		if len(list[0][`images`]) > 0 { // Library paragraph contains images
			images := make([]string, 0)
			_ = tool.JsonDecodeUseNumber(list[0][`images`], &images)
			for _, image := range images {
				content += fmt.Sprintf("\n![img](%s)", image)
			}
		}
		return content, true
	}
	return ``, false
}

// GetRandomSliceReply randomly selects specified number of items from reply content list
// If requested number is greater than total list size, returns all content
// If requested number is less than or equal to 0, returns empty list
func GetRandomSliceReply(replyList []ReplyContent, num int) []ReplyContent {
	// Boundary condition check
	if len(replyList) == 0 || num <= 0 {
		return []ReplyContent{}
	}

	// If requested number is greater than total, return all
	if num >= len(replyList) {
		return replyList
	}

	// Create result slice
	result := make([]ReplyContent, 0, num)

	// Create index slice and shuffle randomly
	indexes := make([]int, len(replyList))
	for i := range indexes {
		indexes[i] = i
	}

	// Fisher-Yates shuffle algorithm to randomize indexes
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(indexes) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		indexes[i], indexes[j] = indexes[j], indexes[i]
	}

	// Take first num elements
	for i := 0; i < num; i++ {
		result = append(result, replyList[indexes[i]])
	}

	return result
}

// BuildKeywordReplyMessage builds keyword reply message
func BuildKeywordReplyMessage(params *define.ChatRequestParam) ([]ReplyContent, bool, error) {
	//part0:init messages
	var replyList []ReplyContent
	// Check if keyword should skip AI reply
	var keywordSkipAI = false

	// Check switch
	robotId := cast.ToInt(params.Robot[`id`])
	adminUserId := cast.ToInt(params.Robot[`admin_user_id`])
	// Keyword reply
	robotAbilityConfig := GetRobotAbilityConfigByAbilityType(adminUserId, robotId, RobotAbilityAutoReply)
	if len(robotAbilityConfig) == 0 {
		// Keyword reply not enabled
		return replyList, false, nil
	}

	// Get all keyword cache
	robotKeywordReplyList, err := GetRobotKeywordReplyListByRobotId(robotId)
	if err != nil {
		return replyList, false, err
	}

	// Whether keyword reply skips AI
	keywordSkipAI = cast.ToInt(robotAbilityConfig[`ai_reply_status`]) != define.SwitchOn

	// Question check
	question := GetFirstQuestionByInput(params.Question) // Special handling for multimodal input

	// Loop check and construct message
	for _, robotKeywordReply := range robotKeywordReplyList {
		// Check keyword
		if robotKeywordReply.SwitchStatus != define.SwitchOn {
			continue
		}
		keywordFlag := false
		// Exact match FullKeyword
		for _, keyword := range robotKeywordReply.FullKeyword {
			if question == keyword {
				// Match success, construct message
				keywordFlag = true
				break
			}
		}

		// Contains match HalfKeyword
		for _, keyword := range robotKeywordReply.HalfKeyword {
			if strings.Contains(question, keyword) {
				// Match success, construct message
				keywordFlag = true
				break
			}
		}

		if keywordFlag {
			// Match success, determine reply type
			if robotKeywordReply.ReplyNum == 0 {
				replyList = append(replyList, robotKeywordReply.ReplyContent...)
			} else {
				// Randomly select ReplyNum items
				// Randomly slice ReplyNum items from array
				selectReplyList := GetRandomSliceReply(robotKeywordReply.ReplyContent, robotKeywordReply.ReplyNum)
				if len(selectReplyList) > 0 {
					replyList = append(replyList, selectReplyList...)
				}
			}
		}
	}
	// Loop replyList to mark source
	replyList = FormatReplyListToDb(replyList, RobotAbilityKeywordReply)
	// Check whether to continue AI
	if len(replyList) == 0 {
		// No keyword matched, continue AI
		return replyList, false, nil
	}
	// Return message
	return replyList, keywordSkipAI, nil
}

// BuildSubscribeReply builds subscribe reply message
func BuildSubscribeReply(params *define.ChatRequestParam, subscribeScene string) ([]ReplyContent, error) {
	//part0:init messages
	var replyList []ReplyContent
	// Check switch
	adminUserId := cast.ToInt(params.AppInfo[`admin_user_id`])
	appid := cast.ToString(params.AppInfo[`app_id`]) // Get app_id from params.AppInfo
	// Keyword reply
	useAbility := CheckUseAbilityByAbilityType(adminUserId, RobotAbilitySubscribeReply)
	if !useAbility {
		// Keyword reply not enabled
		return replyList, nil
	}

	// 1. Get today's weekday as int value
	// time.Weekday: Sunday=0, Monday=1, Tuesday=2, Wednesday=3, Thursday=4, Friday=5, Saturday=6
	weekday := cast.ToInt(time.Now().Weekday())

	var subscribeReplyList []RobotSubscribeReply
	var err error
	needCheck := true

	if needCheck && len(replyList) == 0 {
		// Time check
		subscribeReplyList, err = GetRobotSubscribeReplyListByAppid(adminUserId, appid, define.RuleTypeSubscribeDuration)
		if err != nil {

			return replyList, err
		}

		// Loop check and construct message
		for _, subscribeReply := range subscribeReplyList {
			// Check keyword
			if subscribeReply.SwitchStatus != define.SwitchOn {
				continue
			}
			if subscribeReply.RuleType != define.RuleTypeSubscribeDuration {
				// Not time rule
				continue
			}
			// Time rule and enabled
			checkFlag := false
			switch subscribeReply.DurationType {
			case DurationTypeWeek:
				// Week
				for _, day := range subscribeReply.WeekDuration {
					if day == weekday {
						checkFlag = true
						break
					}
				}
				break
			case DurationTypeDay:
				// Every day
				checkFlag = true
				break
			case DurationTypeTimeRange:
				// Time range
				checkFlag = IsTodayInDateRange(subscribeReply.StartDay, subscribeReply.EndDay)
				break
			default:
				// Default
				break
			}
			// Check whether to continue
			if !checkFlag {
				continue
			}
			// Compare time, skip if not in range
			if !NowInHHmmRangeSimple(subscribeReply.StartDuration, subscribeReply.EndDuration) {
				continue
			}
			// Match message then skip message check
			needCheck = false
			// Check time interval
			if subscribeReply.ReplyInterval > 0 {
				// Check time interval
				var lastTime int
				lastTime, err = GetSubscribeReplyLastTime(adminUserId, subscribeReply.ID, params.Openid)
				if err != nil {
					continue
				}
				nextTime := lastTime + subscribeReply.ReplyInterval
				if nextTime > tool.Time2Int() {
					// Not satisfied within time interval
					break
				}
				// Set time interval
				err = SetSubscribeReplyLastTime(adminUserId, subscribeReply.ID, tool.Time2Int(), params.Openid)
				if err != nil {
					return nil, err
				}
			}

			// Match success, construct message
			if subscribeReply.ReplyNum == 0 {
				replyList = append(replyList, subscribeReply.ReplyContent...)
			} else {
				// Randomly select ReplyNum items
				// Randomly slice ReplyNum items from array
				selectReplyList := GetRandomSliceReply(subscribeReply.ReplyContent, subscribeReply.ReplyNum)
				if len(selectReplyList) > 0 {
					replyList = append(replyList, selectReplyList...)
				}
			}
			break
		}
	}

	if needCheck && len(replyList) == 0 {
		// Get all source checks
		subscribeReplyList, err = GetRobotSubscribeReplyListByAppid(adminUserId, appid, define.RuleTypeSubscribeSource)
		if err != nil {

			return replyList, err
		}
		// Source check
		for _, subscribeReply := range subscribeReplyList {
			// Check keyword
			if subscribeReply.SwitchStatus != define.SwitchOn {
				continue
			}
			if subscribeReply.RuleType != define.RuleTypeSubscribeSource {
				// Not source type
				continue
			}
			// Check source
			if !tool.InArrayString(subscribeScene, subscribeReply.SubscribeSource) {
				// Not in specified source
				continue
			}
			// Match message then skip message check
			needCheck = false
			// Match success, construct message
			if subscribeReply.ReplyNum == 0 {
				replyList = append(replyList, subscribeReply.ReplyContent...)
			} else {
				// Randomly select ReplyNum items
				// Randomly slice ReplyNum items from array
				selectReplyList := GetRandomSliceReply(subscribeReply.ReplyContent, subscribeReply.ReplyNum)
				if len(selectReplyList) > 0 {
					replyList = append(replyList, selectReplyList...)
				}
			}
			break
		}
	}

	if needCheck && len(replyList) == 0 {
		// Check by default follow enable
		subscribeReplyList, err = GetRobotSubscribeReplyListByAppid(adminUserId, appid, define.RuleTypeSubscribeDefault)
		if err != nil {

			return replyList, err
		}
		for _, subscribeReply := range subscribeReplyList {
			if subscribeReply.SwitchStatus != define.SwitchOn {
				continue
			}
			if subscribeReply.RuleType != define.RuleTypeSubscribeDefault {
				// Not default type
				continue
			}

			// Match success, construct message
			if subscribeReply.ReplyNum == 0 {
				replyList = append(replyList, subscribeReply.ReplyContent...)
			} else {
				// Randomly select ReplyNum items
				// Randomly slice ReplyNum items from array
				selectReplyList := GetRandomSliceReply(subscribeReply.ReplyContent, subscribeReply.ReplyNum)
				if len(selectReplyList) > 0 {
					replyList = append(replyList, selectReplyList...)
				}
			}
		}
	}
	// Loop replyList to mark source
	replyList = FormatReplyListToDb(replyList, RobotAbilitySubscribeReply)
	// Check whether to continue AI
	if len(replyList) == 0 {
		// No keyword matched, continue AI
		return replyList, nil
	}
	// Return message
	return replyList, nil
}

func BuildReceivedMessageReply(params *define.ChatRequestParam, messageType string) ([]ReplyContent, error) {
	//part0:init messages
	var replyList []ReplyContent
	// Check switch
	robotId := cast.ToInt(params.Robot[`id`])
	adminUserId := cast.ToInt(params.Robot[`admin_user_id`])
	// Keyword reply
	robotAbilityConfig := GetRobotAbilityConfigByAbilityType(adminUserId, robotId, RobotAbilityAutoReply)
	if len(robotAbilityConfig) == 0 {
		// Keyword reply not enabled
		return replyList, nil
	}

	// 1. Get today's weekday as int value
	// time.Weekday: Sunday=0, Monday=1, Tuesday=2, Wednesday=3, Thursday=4, Friday=5, Saturday=6
	weekday := cast.ToInt(time.Now().Weekday())
	// Get all keyword cache
	receivedMessageRuleList, err := GetRobotReceivedMessageReplyListByRobotId(robotId, RuleTypeDuration)
	if err != nil {

		return replyList, err
	}

	messageTypeCheck := true
	// Loop check and construct message
	for _, receivedMessageRule := range receivedMessageRuleList {
		// Check keyword
		if receivedMessageRule.SwitchStatus != define.SwitchOn {
			continue
		}
		if receivedMessageRule.RuleType != RuleTypeDuration {
			// Not time rule
			continue
		}
		// Time rule and enabled
		checkFlag := false
		switch receivedMessageRule.DurationType {
		case DurationTypeWeek:
			// Week
			for _, day := range receivedMessageRule.WeekDuration {
				if day == weekday {
					checkFlag = true
					break
				}
			}
			break
		case DurationTypeDay:
			// Every day
			checkFlag = true
			break
		case DurationTypeTimeRange:
			// Time range
			checkFlag = IsTodayInDateRange(receivedMessageRule.StartDay, receivedMessageRule.EndDay)
			break
		default:
			// Default
			break
		}
		// Check whether to continue
		if !checkFlag {
			continue
		}
		// Compare time, skip if not in range
		if !NowInHHmmRangeSimple(receivedMessageRule.StartDuration, receivedMessageRule.EndDuration) {
			continue
		}
		// Match message then skip message check
		messageTypeCheck = false
		// Check time interval
		if receivedMessageRule.ReplyInterval > 0 {
			// Check time interval
			var lastTime int
			lastTime, err = GetReceivedMessageReplyLastTime(robotId, receivedMessageRule.ID, params.Openid)
			if err != nil {
				continue
			}
			nextTime := lastTime + receivedMessageRule.ReplyInterval
			if nextTime > tool.Time2Int() {
				// Not satisfied within time interval
				break
			}
			// Set time interval
			err = SetReceivedMessageReplyLastTime(robotId, receivedMessageRule.ID, tool.Time2Int(), params.Openid)
			if err != nil {
				return nil, err
			}
		}

		// Match success, construct message
		if receivedMessageRule.ReplyNum == 0 {
			replyList = append(replyList, receivedMessageRule.ReplyContent...)
		} else {
			// Randomly select ReplyNum items
			// Randomly slice ReplyNum items from array
			selectReplyList := GetRandomSliceReply(receivedMessageRule.ReplyContent, receivedMessageRule.ReplyNum)
			if len(selectReplyList) > 0 {
				replyList = append(replyList, selectReplyList...)
			}
		}
		break
	}

	if messageTypeCheck && len(replyList) == 0 {
		// Check by message type
		receivedMessageRuleList, err = GetRobotReceivedMessageReplyListByRobotId(robotId, RuleTypeMessageType)
		if err != nil {

			return replyList, err
		}
		for _, receivedMessageRule := range receivedMessageRuleList {
			if receivedMessageRule.SwitchStatus != define.SwitchOn {
				continue
			}
			if receivedMessageRule.RuleType != RuleTypeMessageType {
				// Not specified message type
				continue
			}
			checkFlag := false
			switch receivedMessageRule.MessageType {
			case MessageTypeAll:
				checkFlag = true
				break
			case MessageTypeSpecify:
				// Specified message type
				for _, msgType := range receivedMessageRule.SpecifyMessageType {
					if messageType == msgType {
						checkFlag = true
						break
					}
				}
				break
			default:
				// Default
				break
			}
			// Check whether to continue
			if !checkFlag {
				continue
			}
			// Check time interval
			if receivedMessageRule.ReplyInterval > 0 {
				// Check time interval
				var lastTime int
				lastTime, err = GetReceivedMessageReplyLastTime(robotId, receivedMessageRule.ID, params.Openid)
				if err != nil {
					continue
				}
				nextTime := lastTime + receivedMessageRule.ReplyInterval
				if nextTime > tool.Time2Int() {
					// Not satisfied within time interval
					break
				}
				// Set time interval
				err = SetReceivedMessageReplyLastTime(robotId, receivedMessageRule.ID, tool.Time2Int(), params.Openid)
				if err != nil {
					return nil, err
				}
			}
			// Match success, construct message
			if receivedMessageRule.ReplyNum == 0 {
				replyList = append(replyList, receivedMessageRule.ReplyContent...)
			} else {
				// Randomly select ReplyNum items
				// Randomly slice ReplyNum items from array
				selectReplyList := GetRandomSliceReply(receivedMessageRule.ReplyContent, receivedMessageRule.ReplyNum)
				if len(selectReplyList) > 0 {
					replyList = append(replyList, selectReplyList...)
				}
			}

			break
		}
	}

	// Loop replyList to mark source
	replyList = FormatReplyListToDb(replyList, RobotAbilityReceivedMessageReply)
	// Check whether to continue AI
	if len(replyList) == 0 {
		// No keyword matched, continue AI
		return replyList, nil
	}
	// Return message
	return replyList, nil
}

// IsTodayInDateRange concise but robust version (recommended for general scenarios)
func IsTodayInDateRange(start, end string) bool {
	today := time.Now()
	sd, _ := time.Parse("2006-01-02", start)
	ed, _ := time.Parse("2006-01-02", end)
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, sd.Location())
	return !today.Before(sd) && !today.After(ed)
}

func NowInHHmmRangeSimple(start, end string) bool {
	now := time.Now()
	loc := now.Location()

	// Parse HH:mm to temporary time (date will be ignored)
	startTime, _ := time.Parse("15:04", start)
	endTime, _ := time.Parse("15:04", end)

	// Construct today's start and end time (only take hour and minute)
	startT := time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), 0, 0, loc)
	endT := time.Date(now.Year(), now.Month(), now.Day(), endTime.Hour(), endTime.Minute(), 0, 0, loc)

	return !now.Before(startT) && !now.After(endT)
}

func OnlyReceivedMessageReply(params *define.ChatRequestParam) (msql.Params, error) {
	monitor := NewMonitor(params)
	message, err := OnlyReceivedMessageReplyHandle(params, monitor)
	if len(message) > 0 {
		monitor.Save(err)
	}
	return message, err
}

func OnlyReceivedMessageReplyHandle(params *define.ChatRequestParam, monitor *Monitor) (msql.Params, error) {
	var err error
	dialogueId := params.DialogueId
	sessionId, err := GetSessionId(params, dialogueId)
	customer, err := GetCustomerInfo(params.Openid, params.AdminUserId)
	//msgType := getMsgTypeByReceivedMessageType(params.ReceivedMessageType)
	msgType := define.MsgTypeText
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	var receivedMessageJson string
	if len(params.ReceivedMessage) > 0 {
		// Reply content
		receivedMessageJson = tool.JsonEncodeNoError(params.ReceivedMessage)
	}

	// Display content
	showContent, isName := lib_define.MsgTypeNameMap[params.ReceivedMessageType]
	if !isName {
		showContent = i18n.Show(params.Lang, `msg_type_unknown`)
	}
	showContent = i18n.Show(params.Lang, `received_message_type`, showContent)
	// Display image message
	if params.ReceivedMessageType == lib_define.MsgTypeImage && params.MediaIdToOssUrl != `` {
		msgType = define.MsgTypeImage
		showContent = params.MediaIdToOssUrl
	}

	message := msql.Datas{
		`admin_user_id`:             params.AdminUserId,
		`robot_id`:                  params.Robot[`id`],
		`openid`:                    params.Openid,
		`dialogue_id`:               dialogueId,
		`session_id`:                sessionId,
		`is_customer`:               define.MsgFromCustomer,
		`msg_type`:                  msgType,
		`content`:                   showContent,
		`received_message_type`:     params.ReceivedMessageType,
		`received_message`:          receivedMessageJson,
		`media_id_to_oss_url`:       params.MediaIdToOssUrl,
		`thumb_media_id_to_oss_url`: params.ThumbMediaIdToOssUrl,
		`menu_json`:                 ``,
		`quote_file`:                `[]`,
		`create_time`:               tool.Time2Int(),
		`update_time`:               tool.Time2Int(),
	}
	if len(customer) > 0 {
		message[`nickname`] = customer[`nickname`]
		message[`name`] = customer[`name`]
		message[`avatar`] = customer[`avatar`]
	}
	lastChat := msql.Datas{
		`last_chat_time`:    message[`create_time`],
		`last_chat_message`: MbSubstr(cast.ToString(message[`content`]), 0, 1000),
	}
	id, err := msql.Model(`chat_ai_message`, define.Postgres).Insert(message, `id`)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	UpLastChat(dialogueId, sessionId, lastChat, define.MsgFromCustomer)
	//websocket notify
	ReceiverChangeNotify(params.AdminUserId, `c_message`, ToStringMap(message, `id`, id))

	debugLog := make([]any, 0) //debug log
	defer func() {
		monitor.DebugLog = debugLog // Record monitoring data
	}()

	var receivedMessageReplyList []ReplyContent
	// Handle received message reply
	receivedMessageReplyList, _ = BuildReceivedMessageReply(params, params.ReceivedMessageType)

	if len(receivedMessageReplyList) == 0 {
		logs.Error(`received message reply list is empty`)
		return msql.Params{}, nil
	}
	var (
		content, menuJson, reasoningContent string
		requestTime                         int64
		chatResp                            = adaptor.ZhimaChatCompletionResponse{}
		llmStartTime                        = time.Now()
	)

	// Record monitoring data
	monitor.LlmCallTime = time.Now().Sub(llmStartTime).Milliseconds()
	monitor.RequestTime, monitor.Error = requestTime, err

	if *params.IsClose { //client break
		return nil, errors.New(`client break`)
	}

	quoteFile, _ := make([]msql.Params, 0), map[string]struct{}{}
	var quoteFileForSave = make([]msql.Params, len(quoteFile))
	quoteFileJson, _ := tool.JsonEncode(quoteFileForSave)

	message = msql.Datas{
		`admin_user_id`:          params.AdminUserId,
		`robot_id`:               params.Robot[`id`],
		`openid`:                 params.Openid,
		`dialogue_id`:            dialogueId,
		`session_id`:             sessionId,
		`is_customer`:            define.MsgFromRobot,
		`request_time`:           requestTime,
		`recall_time`:            monitor.LibUseTime.RecallTime,
		`msg_type`:               define.MsgTypeText,
		`content`:                content,
		`reasoning_content`:      reasoningContent,
		`is_valid_function_call`: chatResp.IsValidFunctionCall,
		`menu_json`:              menuJson,
		`quote_file`:             quoteFileJson,
		`create_time`:            tool.Time2Int(),
		`update_time`:            tool.Time2Int(),
	}
	if len(params.Robot) > 0 {
		message[`nickname`] = `` //none
		message[`name`] = params.Robot[`robot_name`]
		message[`avatar`] = params.Robot[`robot_avatar`]
	}
	if len(receivedMessageReplyList) > 0 {
		// Reply content
		receivedMessageReplyListJson := tool.JsonEncodeNoError(receivedMessageReplyList)
		message[`reply_content_list`] = receivedMessageReplyListJson
	}

	lastChat = msql.Datas{
		`last_chat_time`:    message[`create_time`],
		`last_chat_message`: MbSubstr(cast.ToString(message[`content`]), 0, 1000),
	}
	id, err = msql.Model(`chat_ai_message`, define.Postgres).Insert(message, `id`)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	UpLastChat(dialogueId, sessionId, lastChat, define.MsgFromRobot)
	//websocket notify
	ReceiverChangeNotify(params.AdminUserId, `ai_message`, ToStringMap(message, `id`, id))

	message["prompt_tokens"] = chatResp.PromptToken
	message["completion_tokens"] = chatResp.CompletionToken
	message["use_model"] = params.Robot["use_model"]
	return ToStringMap(message, `id`, id), nil
}

// SubscribeReplyHandle handles subscribe reply
func SubscribeReplyHandle(params *define.ChatRequestParam, subscribeScene string) (msql.Params, error) {
	// Display content
	var subscribeReplyList []ReplyContent
	appid := cast.ToString(params.AppInfo[`app_id`])
	// Handle received message reply
	subscribeReplyList, _ = BuildSubscribeReply(params, subscribeScene)

	if len(subscribeReplyList) == 0 {
		logs.Error(`subscribe reply list is empty`)
		return msql.Params{}, nil
	}
	var (
		content, menuJson, reasoningContent string
		requestTime                         int64
		chatResp                            = adaptor.ZhimaChatCompletionResponse{}
	)

	quoteFile, _ := make([]msql.Params, 0), map[string]struct{}{}
	var quoteFileForSave = make([]msql.Params, len(quoteFile))
	quoteFileJson, _ := tool.JsonEncode(quoteFileForSave)

	message := msql.Datas{
		`admin_user_id`:          params.AdminUserId,
		`robot_id`:               0,
		`appid`:                  appid,
		`openid`:                 params.Openid,
		`dialogue_id`:            0,
		`session_id`:             0,
		`is_customer`:            define.MsgFromRobot,
		`request_time`:           requestTime,
		`recall_time`:            tool.Time2Int(),
		`msg_type`:               define.MsgTypeText,
		`content`:                content,
		`reasoning_content`:      reasoningContent,
		`is_valid_function_call`: chatResp.IsValidFunctionCall,
		`menu_json`:              menuJson,
		`quote_file`:             quoteFileJson,
		`create_time`:            tool.Time2Int(),
		`update_time`:            tool.Time2Int(),
	}
	if len(params.Robot) > 0 {
		message[`nickname`] = `` //none
		message[`name`] = params.Robot[`robot_name`]
		message[`avatar`] = params.Robot[`robot_avatar`]
	}
	if len(subscribeReplyList) > 0 {
		// Reply content
		subscribeReplyListJson := tool.JsonEncodeNoError(subscribeReplyList)
		message[`reply_content_list`] = subscribeReplyListJson
	}
	message["prompt_tokens"] = chatResp.PromptToken
	message["completion_tokens"] = chatResp.CompletionToken
	message["use_model"] = params.Robot["use_model"]
	return ToStringMap(message, `id`, 0), nil
}
