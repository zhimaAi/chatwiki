// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

/************************************/

type StartNodeParams struct {
	SysGlobal []StartNodeParam `json:"sys_global"`
	DiyGlobal []StartNodeParam `json:"diy_global"`
}

type StartNodeParam struct {
	Key      string `json:"key"`
	Typ      string `json:"typ"`
	Required bool   `json:"required"`
	Desc     string `json:"desc"`
}

/************************************/

const (
	TermTypeEqual      = 1
	TermTypeNotEqual   = 2
	TermTypeContain    = 3
	TermTypeNotContain = 4
	TermTypeEmpty      = 5
	TermTypeNotEmpty   = 6
)

var TermTypes = [...]int{
	TermTypeEqual,
	TermTypeNotEqual,
	TermTypeContain,
	TermTypeNotContain,
	TermTypeEmpty,
	TermTypeNotEmpty,
}

type TermConfig struct {
	Variable string `json:"variable"`
	IsMult   bool   `json:"is_mult"`
	Type     uint   `json:"type"` //1:等于,2不等于,3包含,4不包含,5为空,6不为空
	Value    string `json:"value"`
}

func CompareEqual(single any, typ string, value string) bool {
	if single == nil {
		return false
	}
	switch typ {
	case common.TypString, common.TypArrString:
		return cast.ToString(single) == value
	case common.TypNumber, common.TypArrNumber:
		return cast.ToInt(single) == cast.ToInt(value)
	case common.TypBoole, common.TypArrBoole:
		return cast.ToBool(single) != tool.InArrayString(value, []string{`false`, `0`})
	case common.TypFloat, common.TypArrFloat:
		return cast.ToFloat64(single) == cast.ToFloat64(value)
	case common.TypObject, common.TypArrObject:
		return fmt.Sprintf(`%v`, single) == value || tool.JsonEncodeNoError(single) == value
	case common.TypParams, common.TypArrParams:
		return false //nonsupport
	}
	return false
}

func (term *TermConfig) Verify(flow *WorkFlow) bool {
	field, exist := flow.GetVariable(term.Variable)
	switch term.Type {
	case TermTypeEqual, TermTypeNotEqual:
		if term.IsMult || tool.InArrayString(field.Typ, common.TypArrays[:]) {
			return false //config error
		}
		boole := CompareEqual(field.GetVal(), field.Typ, term.Value) //equal bool
		if term.Type == TermTypeEqual {
			return boole
		} else {
			return !boole
		}
	case TermTypeEmpty, TermTypeNotEmpty:
		boole := !exist || len(field.ShowVals()) == 0 //empty bool
		if term.Type == TermTypeEmpty {
			return boole
		} else {
			return !boole
		}
	case TermTypeContain, TermTypeNotContain:
		if term.IsMult != tool.InArrayString(field.Typ, common.TypArrays[:]) {
			return false //config error
		}
		var boole bool
		if term.IsMult {
			for _, single := range field.GetVals() {
				if CompareEqual(single, field.Typ, term.Value) {
					boole = true
					break
				}
			}
		} else {
			boole = exist && strings.Contains(field.ShowVals(), term.Value)
		}
		if term.Type == TermTypeContain {
			return boole
		} else {
			return !boole
		}
	}
	return false
}

type TermNodeParams []TermNodeParam
type TermNodeParam struct {
	IsOr        bool         `json:"is_or"`
	Terms       []TermConfig `json:"terms"`
	NextNodeKey string       `json:"next_node_key"`
}

func (param *TermNodeParam) Verify(flow *WorkFlow) bool {
	for _, term := range param.Terms {
		boole := term.Verify(flow)
		if param.IsOr && boole {
			return true
		}
		if !param.IsOr && !boole {
			return false
		}
	}
	if param.IsOr {
		return false //all false
	} else {
		return true //all true
	}
}

/************************************/

type Category struct {
	Category    string `json:"category"`
	NextNodeKey string `json:"next_node_key"`
}

type CateNodeParams struct {
	LlmBaseParams
	Categorys []Category `json:"categorys"`
}

/************************************/

type CurlParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

const (
	TypeNone       = 0 //none
	TypeUrlencoded = 1 //x-www-form-urlencoded
	TypeJsonBody   = 2 //application/json
)

type CurlNodeParams struct {
	Method  string               `json:"method"`
	Rawurl  string               `json:"rawurl"`
	Headers []CurlParam          `json:"headers"`
	Params  []CurlParam          `json:"params"`
	Type    uint                 `json:"type"` //0:none,1:x-www-form-urlencoded,2:application/json
	Body    []CurlParam          `json:"body"`
	BodyRaw string               `json:"body_raw"`
	Timeout uint                 `json:"timeout"`
	Output  common.RecurveFields `json:"output"`
}

/************************************/

type LibsNodeParams struct {
	LibraryIds          string          `json:"library_ids"`
	SearchType          common.MixedInt `json:"search_type"`
	TopK                common.MixedInt `json:"top_k"`
	Similarity          float64         `json:"similarity"`
	RerankStatus        uint            `json:"rerank_status"`
	RerankModelConfigId common.MixedInt `json:"rerank_model_config_id"`
	RerankUseModel      string          `json:"rerank_use_model"`
}

/************************************/

type LlmBaseParams struct {
	ModelConfigId common.MixedInt `json:"model_config_id"`
	UseModel      string          `json:"use_model"`
	ContextPair   common.MixedInt `json:"context_pair"`
	Temperature   float32         `json:"temperature"`
	MaxToken      common.MixedInt `json:"max_token"`
	Prompt        string          `json:"prompt"`
}

func (params *LlmBaseParams) Verify(adminUserId int) error {
	if params.ModelConfigId <= 0 || len(params.UseModel) == 0 {
		return errors.New(`请选择使用的LLM模型`)
	}
	//check model_config_id and use_model
	config, err := common.GetModelConfigInfo(params.ModelConfigId.Int(), adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	if len(config) == 0 || !tool.InArrayString(common.Llm, strings.Split(config[`model_types`], `,`)) {
		return errors.New(`使用的LLM服务商选择错误`)
	}
	modelInfo, _ := common.GetModelInfoByDefine(config[`model_define`])
	if !tool.InArrayString(params.UseModel, modelInfo.LlmModelList) && !common.IsMultiConfModel(config[`model_define`]) {
		return errors.New(`使用的LLM模型选择错误`)
	}
	////check function call
	//if modelInfo.SupportedFunctionCallList == nil {
	//	return errors.New(`LLM模型不支持FunctionCall`)
	//}
	//if modelInfo.CheckFancCallRequest != nil {
	//	if err = modelInfo.CheckFancCallRequest(modelInfo, config, params.UseModel); err != nil {
	//		return errors.New(`LLM模型不支持FunctionCall`)
	//	}
	//} else if !tool.InArrayString(params.UseModel, modelInfo.SupportedFunctionCallList) {
	//	return errors.New(`LLM模型不支持FunctionCall`)
	//}
	if params.ContextPair < 0 || params.ContextPair > 50 {
		return errors.New(`上下文数量范围0~50`)
	}
	if params.Temperature < 0 || params.Temperature > 2 {
		return errors.New(`LLM模型温度取值范围0~2`)
	}
	if params.MaxToken < 0 {
		return errors.New(`LLM模型最大token取值错误`)
	}
	//if len(params.Prompt) == 0 {
	//	return errors.New(`提示词内容不能为空`)
	//}
	return nil
}

/************************************/

type LlmNodeParams struct {
	LlmBaseParams
}

/************************************/

type AssignNodeParams []AssignNodeParam
type AssignNodeParam struct {
	Variable string `json:"variable"`
	Value    string `json:"value"`
}

/************************************/

type ReplyNodeParams struct {
	Content string `json:"content"`
}

/************************************/

type ManualNodeParams struct {
}

/************************************/

type NodeParams struct {
	Start  StartNodeParams  `json:"start"`
	Term   TermNodeParams   `json:"term"`
	Cate   CateNodeParams   `json:"cate"`
	Curl   CurlNodeParams   `json:"curl"`
	Libs   LibsNodeParams   `json:"libs"`
	Llm    LlmNodeParams    `json:"llm"`
	Assign AssignNodeParams `json:"assign"`
	Reply  ReplyNodeParams  `json:"reply"`
	Manual ManualNodeParams `json:"manual"`
}

func DisposeNodeParams(nodeType int, nodeParams string) NodeParams {
	params := NodeParams{}
	_ = tool.JsonDecodeUseNumber(nodeParams, &params)
	if nodeType == NodeTypeStart {
		params.Start.SysGlobal = []StartNodeParam{
			{Key: `question`, Typ: common.TypString, Required: true, Desc: `用户消息`},
			{Key: `openid`, Typ: common.TypString, Required: true, Desc: `用户openid`},
		}
	} else {
		params.Start.SysGlobal = make([]StartNodeParam, 0)
	}
	if params.Start.DiyGlobal == nil {
		params.Start.DiyGlobal = make([]StartNodeParam, 0)
	}
	if params.Term == nil {
		params.Term = make(TermNodeParams, 0)
	}
	if params.Cate.Categorys == nil {
		params.Cate.Categorys = make([]Category, 0)
	}
	if params.Curl.Headers == nil {
		params.Curl.Headers = make([]CurlParam, 0)
	}
	if params.Curl.Params == nil {
		params.Curl.Params = make([]CurlParam, 0)
	}
	if params.Curl.Body == nil {
		params.Curl.Body = make([]CurlParam, 0)
	}
	if params.Curl.Output == nil {
		params.Curl.Output = make(common.RecurveFields, 0)
	}
	if params.Assign == nil {
		params.Assign = make([]AssignNodeParam, 0)
	}
	return params
}

type WorkFlowNode struct {
	NodeType     common.MixedInt `json:"node_type"`
	NodeName     string          `json:"node_name"`
	NodeKey      string          `json:"node_key"`
	NodeParams   NodeParams      `json:"node_params"`
	NodeInfoJson map[string]any  `json:"node_info_json"`
	NextNodeKey  string          `json:"next_node_key"`
}

func (node *WorkFlowNode) GetVariables(last ...bool) []string {
	variables := make([]string, 0)
	switch node.NodeType {
	case NodeTypeStart:
		for _, param := range node.NodeParams.Start.DiyGlobal {
			variables = append(variables, fmt.Sprintf(`global.%s`, param.Key))
		}
	case NodeTypeCurl:
		for variable := range common.SimplifyFields(node.NodeParams.Curl.Output) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //上一个节点,兼容旧数据
				variables = append(variables, variable)
			}
		}
	case NodeTypeLibs:
		variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, `special.lib_paragraph_list`))
	case NodeTypeLlm, NodeTypeReply:
		variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, `special.llm_reply_content`))
	}
	return variables
}

type FromNodes map[string][]*WorkFlowNode

func (fn *FromNodes) AddRelation(node *WorkFlowNode, nextNodeKey string) {
	if _, ok := (*fn)[nextNodeKey]; !ok {
		(*fn)[nextNodeKey] = make([]*WorkFlowNode, 0)
	}
	(*fn)[nextNodeKey] = append((*fn)[nextNodeKey], node)
}

func (fn *FromNodes) recurveSetFrom(nodeKey string, nodes *[]*WorkFlowNode) {
	for _, node := range (*fn)[nodeKey] {
		var exist bool
		for i := range *nodes {
			if node.NodeKey == (*nodes)[i].NodeKey {
				exist = true
				break
			}
		}
		if exist {
			continue
		}
		*nodes = append(*nodes, node)
		fn.recurveSetFrom(node.NodeKey, nodes)
	}
}

func (fn *FromNodes) GetVariableList(nodeKey string) []string {
	//系统全局变量
	variables := SysGlobalVariables()
	//上一级节点变量
	for _, node := range (*fn)[nodeKey] {
		variables = append(variables, node.GetVariables(true)...)
	}
	//递归上上级变量
	nodes := make([]*WorkFlowNode, 0)
	fn.recurveSetFrom(nodeKey, &nodes)
	for _, node := range nodes {
		variables = append(variables, node.GetVariables()...)
	}
	//去重
	newVs := make([]string, 0)
	maps := map[string]struct{}{}
	for _, variable := range variables {
		if _, ok := maps[variable]; ok {
			continue
		}
		maps[variable] = struct{}{}
		newVs = append(newVs, variable)
	}
	return newVs
}

func CheckVariablePlaceholder(content string, variables []string) (string, bool) {
	for _, item := range regexp.MustCompile(`【(([a-f0-9]{32}\.)?[a-zA-Z_][a-zA-Z0-9_\-.]*)】`).FindAllStringSubmatch(content, -1) {
		if len(item) > 1 && !tool.InArrayString(item[1], variables) {
			return item[1], false
		}
	}
	return ``, true
}

func VerifyWorkFlowNodes(nodeList []WorkFlowNode, adminUserId int) (startNodeKey, modelConfigIds, libraryIds string, err error) {
	startNodeCount, finishNodeCount := 0, 0
	fromNodes := make(FromNodes)
	for i, node := range nodeList {
		if err = node.Verify(adminUserId); err != nil {
			return
		}
		if node.NodeType <= NodeTypeEdges {
			continue
		}
		if node.NodeType == NodeTypeStart {
			startNodeKey = node.NodeKey
			startNodeCount++
		}
		if !tool.InArrayInt(node.NodeType.Int(), []int{NodeTypeFinish, NodeTypeManual}) {
			fromNodes.AddRelation(&nodeList[i], node.NextNodeKey)
		}
		if node.NodeType == NodeTypeTerm {
			for _, param := range node.NodeParams.Term {
				fromNodes.AddRelation(&nodeList[i], param.NextNodeKey)
			}
		}
		if node.NodeType == NodeTypeCate {
			for _, category := range node.NodeParams.Cate.Categorys {
				fromNodes.AddRelation(&nodeList[i], category.NextNodeKey)
			}
		}
		if tool.InArrayInt(node.NodeType.Int(), []int{NodeTypeFinish, NodeTypeManual}) {
			finishNodeCount++
		}
	}
	if startNodeCount != 1 {
		err = errors.New(`工作流有且仅有一个开始节点`)
		return
	}
	if finishNodeCount == 0 {
		err = errors.New(`工作流必须存在一个结束节点`)
		return
	}
	for _, node := range nodeList {
		if tool.InArrayInt(node.NodeType.Int(), []int{NodeTypeRemark, NodeTypeEdges, NodeTypeStart}) {
			continue
		}
		if _, ok := fromNodes[node.NodeKey]; !ok {
			err = errors.New(`工作流存在游离节点:` + node.NodeName)
			return
		}
		//校验选择的变量必须存在
		variables := fromNodes.GetVariableList(node.NodeKey)
		switch node.NodeType {
		case NodeTypeTerm:
			for _, param := range node.NodeParams.Term {
				for _, term := range param.Terms {
					if !tool.InArrayString(term.Variable, variables) {
						err = errors.New(node.NodeName + `节点选择的变量不存在:` + term.Variable)
						return
					}
				}
			}
		case NodeTypeCurl:
			for _, param := range node.NodeParams.Curl.Headers {
				if variable, ok := CheckVariablePlaceholder(param.Value, variables); !ok {
					err = errors.New(node.NodeName + `节点Headers变量不存在:` + variable)
					return
				}
			}
			for _, param := range node.NodeParams.Curl.Params {
				if variable, ok := CheckVariablePlaceholder(param.Value, variables); !ok {
					err = errors.New(node.NodeName + `节点Params变量不存在:` + variable)
					return
				}
			}
			switch node.NodeParams.Curl.Type {
			case TypeUrlencoded:
				for _, param := range node.NodeParams.Curl.Body {
					if variable, ok := CheckVariablePlaceholder(param.Value, variables); !ok {
						err = errors.New(node.NodeName + `节点Body变量不存在:` + variable)
						return
					}
				}
			case TypeJsonBody:
				if variable, ok := CheckVariablePlaceholder(node.NodeParams.Curl.BodyRaw, variables); !ok {
					err = errors.New(node.NodeName + `节点JsonBody变量不存在:` + variable)
					return
				}
			}
		//case NodeTypeCate:
		//	if variable, ok := CheckVariablePlaceholder(node.NodeParams.Cate.Prompt, variables); !ok {
		//		err = errors.New(node.NodeName + `节点提示词变量不存在:` + variable)
		//		return
		//	}
		case NodeTypeLlm:
			if variable, ok := CheckVariablePlaceholder(node.NodeParams.Llm.Prompt, variables); !ok {
				err = errors.New(node.NodeName + `节点提示词变量不存在:` + variable)
				return
			}
		case NodeTypeAssign:
			for i, param := range node.NodeParams.Assign {
				if !tool.InArrayString(param.Variable, variables) { //自定义变量不存在
					err = errors.New(node.NodeName + `自定义全局变量不存在:` + param.Variable)
					return
				}
				if variable, ok := CheckVariablePlaceholder(param.Value, variables); !ok {
					err = errors.New(node.NodeName + fmt.Sprintf(`第%d行:变量不存在:`, i+1) + variable)
					return
				}
			}
		case NodeTypeReply:
			if variable, ok := CheckVariablePlaceholder(node.NodeParams.Reply.Content, variables); !ok {
				err = errors.New(node.NodeName + `节点插入的变量不存在:` + variable)
				return
			}
		}
	}
	var libraryArr []string
	//采集使用的模型id集合
	for _, node := range nodeList {
		var modelConfigId int
		switch node.NodeType {
		case NodeTypeCate:
			modelConfigId = node.NodeParams.Cate.ModelConfigId.Int()
		case NodeTypeLibs:
			modelConfigId = node.NodeParams.Libs.RerankModelConfigId.Int()
			if node.NodeParams.Libs.LibraryIds != "" {
				libraryArr = append(libraryArr, strings.Split(node.NodeParams.Libs.LibraryIds, `,`)...)
			}
		case NodeTypeLlm:
			modelConfigId = node.NodeParams.Llm.ModelConfigId.Int()
		}
		if modelConfigId > 0 {
			if len(modelConfigIds) > 0 {
				modelConfigIds += `,`
			}
			modelConfigIds += cast.ToString(modelConfigId)
		}
	}
	libraryIds = strings.Join(libraryArr, `,`)
	return
}

func (node *WorkFlowNode) Verify(adminUserId int) error {
	if !tool.InArrayInt(node.NodeType.Int(), NodeTypes[:]) {
		return errors.New(`节点类型参数错误:` + cast.ToString(node.NodeType))
	}
	if len(node.NodeName) == 0 {
		return errors.New(`节点名称不能为空:` + node.NodeKey)
	}
	if len(node.NodeKey) == 0 || !common.IsMd5Str(node.NodeKey) {
		return errors.New(`节点NodeKey参数为空或格式错误:` + node.NodeName)
	}
	if len(node.NextNodeKey) > 0 && !common.IsMd5Str(node.NextNodeKey) {
		return errors.New(`节点NextNodeKey参数格式错误:` + node.NodeName)
	}
	if len(node.NextNodeKey) == 0 && !tool.InArrayInt(node.NodeType.Int(), []int{NodeTypeRemark, NodeTypeEdges, NodeTypeFinish, NodeTypeManual}) {
		return errors.New(`节点没有指定下一个节点:` + node.NodeName)
	}
	var err error
	switch node.NodeType {
	case NodeTypeStart:
		err = node.NodeParams.Start.Verify()
	case NodeTypeTerm:
		err = node.NodeParams.Term.Verify()
	case NodeTypeCate:
		err = node.NodeParams.Cate.Verify(adminUserId)
	case NodeTypeCurl:
		err = node.NodeParams.Curl.Verify()
	case NodeTypeLibs:
		err = node.NodeParams.Libs.Verify(adminUserId)
	case NodeTypeLlm:
		err = node.NodeParams.Llm.Verify(adminUserId)
	case NodeTypeAssign:
		err = node.NodeParams.Assign.Verify()
	case NodeTypeReply:
		err = node.NodeParams.Reply.Verify()
	case NodeTypeManual:
		err = node.NodeParams.Manual.Verify(adminUserId)
	}
	if err != nil {
		return errors.New(node.NodeName + `节点:` + err.Error())
	}
	return nil
}

func (params *StartNodeParams) Verify() error {
	maps := map[string]struct{}{}
	for _, item := range params.DiyGlobal {
		if !common.IsVariableName(item.Key) {
			return errors.New(fmt.Sprintf(`自定义全局变量名格式错误:%s`, item.Key))
		}
		if tool.InArrayString(fmt.Sprintf(`global.%s`, item.Key), SysGlobalVariables()) {
			return errors.New(fmt.Sprintf(`自定义全局变量与系统变量同名:%s`, item.Key))
		}
		if !tool.InArrayString(item.Typ, []string{common.TypString, common.TypNumber, common.TypArrString}) {
			return errors.New(fmt.Sprintf(`自定义全局变量类型不支持:%s`, item.Key))
		}
		if _, ok := maps[item.Key]; ok {
			return errors.New(fmt.Sprintf(`自定义全局变量名重复定义:%s`, item.Key))
		}
		maps[item.Key] = struct{}{}
	}
	return nil
}

func (params *TermNodeParams) Verify() error {
	if params == nil || len(*params) == 0 {
		return errors.New(`配置参数不能为空`)
	}
	for i, item := range *params {
		if len(item.Terms) == 0 {
			return errors.New(fmt.Sprintf(`第%d分支配置为空`, i+1))
		}
		for j, term := range item.Terms {
			if len(term.Variable) == 0 {
				return errors.New(fmt.Sprintf(`第%d分支的第%d条件:请选择变量`, i+1, j+1))
			}
			if !common.IsVariableNames(term.Variable) {
				return errors.New(fmt.Sprintf(`第%d分支的第%d条件:变量格式错误`, i+1, j+1))
			}
			if term.IsMult { //数组类型的不支持 等于和不等于
				if !tool.InArrayInt(int(term.Type), TermTypes[2:]) {
					return errors.New(fmt.Sprintf(`第%d分支的第%d条件:匹配条件错误`, i+1, j+1))
				}
			} else {
				if !tool.InArrayInt(int(term.Type), TermTypes[:]) {
					return errors.New(fmt.Sprintf(`第%d分支的第%d条件:匹配条件错误`, i+1, j+1))
				}
			}
			if !tool.InArrayInt(int(term.Type), []int{TermTypeEmpty, TermTypeNotEmpty}) && len(term.Value) == 0 {
				return errors.New(fmt.Sprintf(`第%d分支的第%d条件:请输入匹配值`, i+1, j+1))
			}
		}
		if len(item.NextNodeKey) == 0 || !common.IsMd5Str(item.NextNodeKey) {
			return errors.New(fmt.Sprintf(`第%d分支:下一个节点未指定或格式错误`, i+1))
		}
	}
	return nil
}

func (params *CateNodeParams) Verify(adminUserId int) error {
	if err := params.LlmBaseParams.Verify(adminUserId); err != nil {
		return err
	}
	if len(params.Categorys) == 0 {
		return errors.New(`分类列表不能为空`)
	}
	for i, category := range params.Categorys {
		if len(category.Category) == 0 {
			return errors.New(fmt.Sprintf(`第%d个分类:分类名称为空`, i+1))
		}
		if len(category.NextNodeKey) == 0 || !common.IsMd5Str(category.NextNodeKey) {
			return errors.New(fmt.Sprintf(`第%d个分类:下一个节点未指定或格式错误`, i+1))
		}
	}
	return nil
}

func (params *CurlNodeParams) Verify() error {
	if !tool.InArrayString(params.Method, []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}) {
		return errors.New(`请求方式参数错误`)
	}
	if _, err := url.Parse(params.Rawurl); err != nil || len(params.Rawurl) == 0 {
		return errors.New(`请求链接为空或错误`)
	}
	for _, header := range params.Headers {
		if len(header.Key) == 0 || len(header.Value) == 0 {
			return errors.New(`请求头的键值对不能为空`)
		}
		if header.Key == `Content-Type` {
			return errors.New(`请求头Content-Type不允许被设置`)
		}
	}
	for _, param := range params.Params {
		if len(param.Key) == 0 || len(param.Value) == 0 {
			return errors.New(`params的键值对不能为空`)
		}
	}
	if params.Method != http.MethodGet {
		switch params.Type {
		case TypeNone:
		case TypeUrlencoded:
			for _, param := range params.Body {
				if len(param.Key) == 0 || len(param.Value) == 0 {
					return errors.New(`body的键值对不能为空`)
				}
			}
		case TypeJsonBody:
			if len(params.BodyRaw) == 0 {
				return errors.New(`JSONBody不能为空`)
			}
			var temp any
			if err := tool.JsonDecodeUseNumber(params.BodyRaw, &temp); err != nil {
				return errors.New(`Body不是一个JSON字符串`)
			}
		default:
			return errors.New(`body参数类型选择错误`)
		}
	}
	if params.Timeout > 60 {
		return errors.New(`请求超时时间最大值60秒`)
	}
	//输出字段校验
	return params.Output.Verify()
}

func (params *LibsNodeParams) Verify(adminUserId int) error {
	if len(params.LibraryIds) == 0 || !common.CheckIds(params.LibraryIds) {
		return errors.New(`关联知识库为空或参数错误`)
	}
	for _, libraryId := range strings.Split(params.LibraryIds, `,`) {
		info, err := common.GetLibraryInfo(cast.ToInt(libraryId), adminUserId)
		if err != nil {
			logs.Error(err.Error())
			return err
		}
		if len(info) == 0 {
			return errors.New(`关联知识库不存在ID:` + libraryId)
		}
	}
	if !tool.InArrayInt(params.SearchType.Int(), []int{define.SearchTypeMixed, define.SearchTypeVector, define.SearchTypeFullText}) {
		return errors.New(`知识库检索模式参数错误`)
	}
	if params.TopK <= 0 || params.TopK > 10 {
		return errors.New(`知识库检索TopK范围1~10`)
	}
	if params.Similarity < 0 || params.Similarity > 1 {
		return errors.New(`知识库检索相似度阈值0~1`)
	}
	if params.RerankStatus > 0 || params.RerankModelConfigId != 0 || len(params.RerankUseModel) > 0 {
		if params.RerankModelConfigId <= 0 || len(params.RerankUseModel) == 0 {
			return errors.New(`请选择使用的Rerank模型`)
		}
		config, err := common.GetModelConfigInfo(params.RerankModelConfigId.Int(), adminUserId)
		if err != nil {
			logs.Error(err.Error())
			return err
		}
		if len(config) == 0 || !tool.InArrayString(common.Rerank, strings.Split(config[`model_types`], `,`)) {
			return errors.New(`使用的Rerank服务商选择错误`)
		}
		modelInfo, _ := common.GetModelInfoByDefine(config[`model_define`])
		if !tool.InArrayString(params.RerankUseModel, modelInfo.RerankModelList) {
			return errors.New(`使用的Rerank模型选择错误`)
		}
	}
	return nil
}

func (params *LlmNodeParams) Verify(adminUserId int) error {
	if err := params.LlmBaseParams.Verify(adminUserId); err != nil {
		return err
	}
	if len(params.Prompt) == 0 {
		return errors.New(`提示词内容不能为空`)
	}
	return nil
}

func (params *AssignNodeParams) Verify() error {
	if params == nil || len(*params) == 0 {
		return errors.New(`配置参数不能为空`)
	}
	for i, param := range *params {
		if len(param.Variable) == 0 {
			return errors.New(fmt.Sprintf(`第%d行:请选择变量`, i+1))
		}
		if !strings.HasPrefix(param.Variable, `global.`) || !common.IsVariableNames(param.Variable) {
			return errors.New(fmt.Sprintf(`第%d行:变量格式错误`, i+1))
		}
		if tool.InArrayString(param.Variable, SysGlobalVariables()) {
			return errors.New(fmt.Sprintf(`第%d行:系统全局变量禁止被赋值`, i+1))
		}
	}
	return nil
}

func (params *ReplyNodeParams) Verify() error {
	if len(params.Content) == 0 {
		return errors.New(`消息内容不能为空`)
	}
	return nil
}

func (params *ManualNodeParams) Verify(adminUserId int) error {
	return errors.New(`仅云版支持转人工节点`)
}
