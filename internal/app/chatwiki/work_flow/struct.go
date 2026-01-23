// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_web"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/zhimaAi/go_tools/curl"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

/************************************/

type StartNodeParams struct {
	SysGlobal   []StartNodeParam `json:"sys_global"` //废弃字段
	DiyGlobal   []StartNodeParam `json:"diy_global"`
	TriggerList []TriggerConfig  `json:"trigger_list"`
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
	Categorys     []Category `json:"categorys"`
	QuestionValue string     `json:"question_value"`
}

/************************************/

type CurlParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CurlAuthParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	AddTo string `json:"add_to"`
}

const (
	TypeNone       = 0 //none
	TypeUrlencoded = 1 //x-www-form-urlencoded
	TypeJsonBody   = 2 //application/json
)

type CurlNodeParams struct {
	Method   string               `json:"method"`
	Rawurl   string               `json:"rawurl"`
	Headers  []CurlParam          `json:"headers"`
	Params   []CurlParam          `json:"params"`
	Type     uint                 `json:"type"` //0:none,1:x-www-form-urlencoded,2:application/json
	Body     []CurlParam          `json:"body"`
	BodyRaw  string               `json:"body_raw"`
	Timeout  uint                 `json:"timeout"`
	Output   common.RecurveFields `json:"output"`
	HttpAuth []CurlAuthParam      `json:"http_auth"`
}

/************************************/

type LibsNodeParams struct {
	LibraryIds              string          `json:"library_ids"`
	SearchType              common.MixedInt `json:"search_type"`
	RrfWeight               string          `json:"rrf_weight"`
	TopK                    common.MixedInt `json:"top_k"`
	Similarity              float64         `json:"similarity"`
	RerankStatus            uint            `json:"rerank_status"`
	RerankModelConfigId     common.MixedInt `json:"rerank_model_config_id"`
	RerankUseModel          string          `json:"rerank_use_model"`
	QuestionValue           string          `json:"question_value"`
	MetaSearchSwitch        int             `json:"meta_search_switch"`
	MetaSearchType          int             `json:"meta_search_type"`
	MetaSearchConditionList string          `json:"meta_search_condition_list"`
}

/************************************/

type LlmBaseParams struct {
	ModelConfigId  common.MixedInt `json:"model_config_id"`
	UseModel       string          `json:"use_model"`
	ContextPair    common.MixedInt `json:"context_pair"`
	EnableThinking bool            `json:"enable_thinking"`
	Temperature    float32         `json:"temperature"`
	MaxToken       common.MixedInt `json:"max_token"`
	Prompt         string          `json:"prompt"`
}

func (params *LlmBaseParams) Verify(adminUserId int) error {
	if params.ModelConfigId <= 0 || len(params.UseModel) == 0 {
		return errors.New(`请选择使用的LLM模型`)
	}
	//check model_config_id and use_model
	if ok := common.CheckModelIsValid(adminUserId, params.ModelConfigId.Int(), params.UseModel, common.Llm); !ok {
		return errors.New(`使用的LLM模型选择错误`)
	}
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
	QuestionValue string `json:"question_value"`
	LibsNodeKey   string `json:"libs_node_key"`
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

const StaffAll = 1 //转接类型:1自动分配,2指定客服,3指定客服组
const StaffIds = 2
const StaffGroup = 3

type ManualNodeParams struct {
	SwitchType    common.MixedInt `json:"switch_type"`
	SwitchStaff   string          `json:"switch_staff"`
	SwitchContent string          `json:"switch_content"`
}

/************************************/

type QuestionOptimizeNodeParams struct {
	LlmBaseParams
	QuestionValue string `json:"question_value"`
}

func (params *QuestionOptimizeNodeParams) Verify(adminUserId int) error {
	if len(params.QuestionValue) == 0 {
		return errors.New(`用户问题不能为空`)
	}
	return params.LlmBaseParams.Verify(adminUserId)
}

/************************************/

type ParamsExtractorNodeParams struct {
	LlmBaseParams
	QuestionValue string               `json:"question_value"`
	Output        common.RecurveFields `json:"output"`
}

func (params *ParamsExtractorNodeParams) Verify(adminUserId int) error {
	if len(params.QuestionValue) == 0 {
		return errors.New(`用户问题不能为空`)
	}
	if err := params.LlmBaseParams.Verify(adminUserId); err != nil {
		return err
	}
	//输出字段校验
	return params.Output.Verify()
}

/************************************/

type FormFieldTyp struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type FormFieldValue struct {
	FormFieldTyp
	Value string `json:"value"`
}

type FormInsertNodeParams struct {
	FormId common.MixedInt  `json:"form_id"`
	Datas  []FormFieldValue `json:"datas"`
}

type FormDeleteNodeParams struct {
	FormId common.MixedInt              `json:"form_id"`
	Typ    common.MixedInt              `json:"typ"`
	Where  []define.FormFilterCondition `json:"where"`
}

type FormUpdateNodeParams struct {
	FormId common.MixedInt              `json:"form_id"`
	Typ    common.MixedInt              `json:"typ"`
	Where  []define.FormFilterCondition `json:"where"`
	Datas  []FormFieldValue             `json:"datas"`
}

type FormFieldOrder struct {
	FormFieldTyp
	IsAsc bool `json:"is_asc"`
}

type FormSelectNodeParams struct {
	FormId common.MixedInt              `json:"form_id"`
	Typ    common.MixedInt              `json:"typ"`
	Where  []define.FormFilterCondition `json:"where"`
	Fields []FormFieldTyp               `json:"fields"`
	Order  []FormFieldOrder             `json:"order"`
	Limit  common.MixedInt              `json:"limit"`
}

/************************************/

type CodeRunParams struct {
	Field    string `json:"field"`
	Variable string `json:"variable"`
}

type CodeRunNodeParams struct {
	Params    []CodeRunParams      `json:"params"`
	MainFunc  string               `json:"main_func"`
	Timeout   uint                 `json:"timeout"`
	Output    common.RecurveFields `json:"output"`
	Exception string               `json:"exception"`
}

type McpNodeParams struct {
	ProviderId uint           `json:"provider_id"`
	ToolName   string         `json:"tool_name"`
	Arguments  map[string]any `json:"arguments"`
	Output     string         `json:"output"`
}

type LoopNodeParams struct {
	LoopType           string             `json:"loop_type"`           //array : loop from array ,number : loop from number.Get the loop count
	LoopArrays         []common.LoopField `json:"loop_arrays"`         //Set when LoopType equals 'array'
	LoopNumber         common.MixedInt    `json:"loop_number"`         //Set when LoopType equals 'number'
	IntermediateParams []common.LoopField `json:"intermediate_params"` //intermediate_params
	Output             []common.LoopField `json:"output"`              //
}

type BatchNodeParams struct {
	BatchArrays  []common.LoopField `json:"batch_arrays"`
	ChanNumber   common.MixedInt    `json:"chan_number"`
	MaxRunNumber common.MixedInt    `json:"max_run_number"`
	Output       []common.LoopField `json:"output"`
}

type PluginNodeParams struct {
	Name   string               `json:"name"`
	Type   string               `json:"type"`
	Params map[string]any       `json:"params"`
	Output common.RecurveFields `json:"output_obj"`
}

type ImageGenerationParams struct {
	UseModel            string               `json:"use_model"`
	ModelConfigId       string               `json:"model_config_id"`
	Size                string               `json:"size"`
	ImageNum            string               `json:"image_num"`
	Prompt              string               `json:"prompt"`
	InputImages         []string             `json:"input_images"`
	ImageWatermark      string               `json:"image_watermark"`       //是否添加水印，1添加，0不添加
	ImageOptimizePrompt string               `json:"image_optimize_prompt"` //是否开启优化提示词，1开启，0不开启
	Output              common.RecurveFields `json:"output"`
}

type JsonEncodeParams struct {
	InputVariable string               `json:"input_variable"`
	Outputs       common.RecurveFields `json:"output"`
}

type JsonDecodeParams struct {
	InputVariable string               `json:"input_variable"`
	Outputs       common.RecurveFields `json:"output"`
}

type TextToAudioNodeParams struct {
	ModelId       int    `json:"model_id"`
	ModelConfigId int    `json:"model_config_id"`
	UseModel      string `json:"use_model"`
	VoiceType     string `json:"voice_type"` //希望查询音色类型，可用选项: system, voice_cloning, voice_generation, all
	Arguments     struct {
		Text          string `json:"text"`
		ModelId       int    `json:"model_id"`
		UseModel      string `json:"use_model"`
		ModelConfigId int    `json:"model_config_id"`
		VoiceSetting  struct {
			VoiceId           string  `json:"voice_id"`           // 音色编号
			VoiceName         string  `json:"voice_name"`         // 冗余
			Speed             float32 `json:"speed"`              // 语速
			Vol               int     `json:"vol"`                // 音量
			Pitch             int     `json:"pitch"`              // 语调
			Emotion           string  `json:"emotion"`            // 情绪
			TextNormalization bool    `json:"text_normalization"` // 中英文规范化
		} `json:"voice_setting"` // 声音设置
		AudioSetting struct {
			SampleRate int    `json:"sample_rate"` // 采样率
			Bitrate    int    `json:"bitrate"`     // 比特率
			Format     string `json:"format"`      // 格式
			Channel    int    `json:"channel"`     // 声道数
			ForceCbr   bool   `json:"force_cbr"`   // 恒定比特率
		} `json:"audio_setting"` // 音频设置
		VoiceModify struct {
			Pitch        int    `json:"pitch"`         // 音高调整
			Intensity    int    `json:"intensity"`     // 强度调整（力量感/柔和）
			Timbre       int    `json:"timbre"`        // 音色调整
			SoundEffects string `json:"sound_effects"` // 音效设置
		} `json:"voice_modify"` // 音调调整
		LanguageBoost  string `json:"language_boost"` // 小语种识别能力
		OutputFormat   string `json:"output_format"`
		SubtitleEnable bool   `json:"subtitle_enable"`
		AigcWatermark  any    `json:"aigc_watermark"`
		TagMap         any    `json:"tag_map"`
	} `json:"arguments"`
	Output common.RecurveFields `json:"output"`
}

type VoiceCloneNodeParams struct {
	ModelConfigId int `json:"model_config_id"`
	Arguments     struct {
		ModelConfigId any    `json:"model_config_id"`
		ModelId       any    `json:"model_id"`
		FileUrl       string `json:"file_url"` //待复刻音频的url，要上传后转成file_id传给minimax
		VoiceId       string `json:"voice_id"` //克隆音色的 voice_id
		ClonePrompt   struct {
			PromptAudioUrl string `json:"prompt_audio_url"` //示例音频的url，要上传后转成file_id传给minimax
			PromptText     string `json:"prompt_text"`      //示例音频的对应文本，需确保和音频内容一致，句末需有标点符号做结尾
		} `json:"clone_prompt"` //音色复刻示例音频，提供本参数将有助于增强语音合成的音色相似度和稳定性
		Text                    string `json:"text"`                      //复刻试听参数，限制 1000 字符以内
		LanguageBoost           string `json:"language_boost"`            //是否增强对指定的小语种和方言的识别能力
		Model                   string `json:"model"`                     //复刻试听参数
		NeedNoiseReduction      bool   `json:"need_noise_reduction"`      //开启降噪
		NeedVolumeNormalization bool   `json:"need_volume_normalization"` //是否开启音量归一化
		AigcWatermark           any    `json:"aigc_watermark"`            //是否在合成试听音频的末尾添加音频节奏标识，默认值为 false
		TagMap                  any    `json:"tag_map"`
	} `json:"arguments"`
	Output common.RecurveFields `json:"output"`
}

type LibraryImportParams struct {
	ImportType                string               `json:"import_type"`                  //content 导入内容，url 导入url
	LibraryId                 string               `json:"library_id"`                   //知识库id 不能为空
	LibraryGroupId            string               `json:"library_group_id"`             //知识库分组id 0表示未分组
	NormalUrl                 string               `json:"normal_url"`                   //普通知识库：文档url
	NormalTitle               string               `json:"normal_title"`                 //普通知识库：文档标题
	NormalContent             string               `json:"normal_content"`               //普通知识库：文档内容
	NormalUrlRepeatOp         string               `json:"normal_url_repeat_op"`         //普通知识库：url重复时的操作 import 依然导入，not import 不导入，update 更新内容
	QaQuestion                string               `json:"qa_question"`                  //问答知识库：分段问题
	QaAnswer                  string               `json:"qa_answer"`                    //问答知识库：分段答案
	QaImagesVariable          string               `json:"qa_images_variable"`           //问答知识库：答案附图 array<string>
	QaSimilarQuestionVariable string               `json:"qa_similar_question_variable"` //问答知识库：相似问法 array<string>
	QaRepeatOp                string               `json:"qa_repeat_op"`                 //问答知识库：问题重复时的操作 import 依然导入，not import 不导入，update 更新内容
	Outputs                   common.RecurveFields `json:"outputs"`                      //输出固定一个msg string
}

type WorkflowNodeParams struct {
	RobotId   int `json:"robot_id"`
	RobotInfo any `json:"robot_info"` // 前端使用的临时数据
	Params    []struct {
		StartNodeParam
		Variable string `json:"variable"` //对应的全部变量
		Tags     any    `json:"tags"`     // 前端使用的临时数据
	} `json:"params"`
	Output    common.RecurveFields `json:"output"`
	Exception string               `json:"exception"`
}

type Message struct {
	Type    string `json:"type"`    //text image voice
	Content string `json:"content"` //content
}

type FinishNodeParams struct {
	OutType  string               `json:"out_type"` //variable返回消息和变量，message返回消息
	Messages []Message            `json:"messages"` //具体的消息
	Outputs  common.RecurveFields `json:"outputs"`  //返回的变量
}

/************************************/

type ImmediatelyReplyNodeParams struct {
	Content string `json:"content"`
}

/************************************/

var LoopAllowNodeTypes = []int{
	NodeTypeRemark,
	NodeTypeTerm,
	NodeTypeCate,
	NodeTypeCurl,
	NodeTypeLibs,
	NodeTypeLlm,
	NodeTypeFinish,
	NodeTypeAssign,
	NodeTypeReply,
	NodeTypeQuestionOptimize,
	NodeTypeParamsExtractor,
	NodeTypeFormInsert,
	NodeTypeFormDelete,
	NodeTypeFormUpdate,
	NodeTypeFormSelect,
	NodeTypeCodeRun,
	NodeTypeMcp,
	NodeTypeLoopEnd,
	NodeTypeLoopStart,
	NodeTypePlugin,
	NodeTypeImageGeneration,
	NodeTypeJsonEncode,
	NodeTypeJsonDecode,
	NodeTypeTextToAudio,
	NodeTypeVoiceClone,
	NodeTypeLibraryImport,
	NodeTypeWorkflow,
	NodeTypeImmediatelyReply,
}

var BatchAllowNodeTypes = []int{
	NodeTypeRemark,
	NodeTypeTerm,
	NodeTypeCate,
	NodeTypeCurl,
	NodeTypeLibs,
	NodeTypeLlm,
	NodeTypeFinish,
	NodeTypeAssign,
	NodeTypeReply,
	NodeTypeQuestionOptimize,
	NodeTypeParamsExtractor,
	NodeTypeFormInsert,
	NodeTypeFormDelete,
	NodeTypeFormUpdate,
	NodeTypeFormSelect,
	NodeTypeCodeRun,
	NodeTypeMcp,
	NodeTypeBatchStart,
	NodeTypePlugin,
	NodeTypeImageGeneration,
	NodeTypeJsonEncode,
	NodeTypeJsonDecode,
	NodeTypeTextToAudio,
	NodeTypeVoiceClone,
	NodeTypeLibraryImport,
	NodeTypeWorkflow,
	NodeTypeImmediatelyReply,
}

type ExportLibrary struct {
	Libs struct {
		LibraryIds string `json:"library_ids"`
	} `json:"libs"`
}
type ExportBaseFormInfo struct {
	FormDescription string `json:"form_description"`
	FormId          any    `json:"form_id"`
	FormName        string `json:"form_name"`
}
type ExportFormDeleteInfo struct {
	FormDelete ExportBaseFormInfo `json:"form_delete"`
}
type ExportFormInsertInfo struct {
	FormInsert ExportBaseFormInfo `json:"form_insert"`
}
type ExportFormUpdateInfo struct {
	FormUpdate ExportBaseFormInfo `json:"form_update"`
}
type ExportFormSelectInfo struct {
	FormSelect ExportBaseFormInfo `json:"form_select"`
}

type NodeInfo struct {
	DataRaw string `json:"dataRaw"`
}

/************************************/

type NodeParams struct {
	Start            StartNodeParams            `json:"start"`
	Term             TermNodeParams             `json:"term"`
	Cate             CateNodeParams             `json:"cate"`
	Curl             CurlNodeParams             `json:"curl"`
	Libs             LibsNodeParams             `json:"libs"`
	Llm              LlmNodeParams              `json:"llm"`
	Assign           AssignNodeParams           `json:"assign"`
	Reply            ReplyNodeParams            `json:"reply"`
	Manual           ManualNodeParams           `json:"manual"`
	QuestionOptimize QuestionOptimizeNodeParams `json:"question_optimize"`
	ParamsExtractor  ParamsExtractorNodeParams  `json:"params_extractor"`
	FormInsert       FormInsertNodeParams       `json:"form_insert"`
	FormDelete       FormDeleteNodeParams       `json:"form_delete"`
	FormUpdate       FormUpdateNodeParams       `json:"form_update"`
	FormSelect       FormSelectNodeParams       `json:"form_select"`
	CodeRun          CodeRunNodeParams          `json:"code_run"`
	Mcp              McpNodeParams              `json:"mcp"`
	Loop             LoopNodeParams             `json:"loop"`
	Plugin           PluginNodeParams           `json:"plugin"`
	Batch            BatchNodeParams            `json:"batch"`
	Finish           FinishNodeParams           `json:"finish"`
	ImageGeneration  ImageGenerationParams      `json:"image_generation"`
	JsonEncode       JsonEncodeParams           `json:"json_encode"`
	JsonDecode       JsonDecodeParams           `json:"json_decode"`
	TextToAudio      TextToAudioNodeParams      `json:"text_to_audio"`
	VoiceClone       VoiceCloneNodeParams       `json:"voice_clone"`
	LibraryImport    LibraryImportParams        `json:"library_import"`
	Workflow         WorkflowNodeParams         `json:"workflow"`
	ImmediatelyReply ImmediatelyReplyNodeParams `json:"immediately_reply"`
}

func FillDiyGlobalBlanks(output TriggerOutputParam, start *StartNodeParams) {
	if len(output.Variable) == 0 {
		return //未配置变量映射
	}
	key, found := strings.CutPrefix(output.Variable, `global.`)
	if !found {
		return //映射错误的,不管
	}
	for _, param := range start.DiyGlobal {
		if param.Key == key {
			return //已存在的变量直接跳过
		}
	}
	start.DiyGlobal = append(start.DiyGlobal, output.StartNodeParam)
}

func DisposeNodeParams(nodeType int, nodeParams string) NodeParams {
	params := NodeParams{}
	_ = tool.JsonDecodeUseNumber(nodeParams, &params)
	params.Start.SysGlobal = make([]StartNodeParam, 0) //废弃字段
	if params.Start.DiyGlobal == nil {
		params.Start.DiyGlobal = make([]StartNodeParam, 0)
	}
	if params.Start.TriggerList == nil {
		params.Start.TriggerList = make([]TriggerConfig, 0)
	}
	if nodeType == NodeTypeStart && len(params.Start.TriggerList) == 0 { //默认值处理
		chatTrigger := GetTriggerChatConfig()
		params.Start.TriggerList = []TriggerConfig{chatTrigger}
		for _, output := range chatTrigger.Outputs {
			params.Start.DiyGlobal = append(params.Start.DiyGlobal, output.StartNodeParam)
		}
	}
	if nodeType == NodeTypeStart { //开始节点触发器输出变量旧数据兼容
		for i, trigger := range params.Start.TriggerList {
			outputs, exist := GetTriggerOutputsByType(trigger.TriggerType)
			if !exist {
				continue
			}
			if trigger.TriggerType == TriggerTypeOfficial {
				switch trigger.TriggerOfficialConfig.MsgType {
				case define.TriggerOfficialMessage:
					outputs = GetMessage()
				case define.TriggerOfficialQrCodeScan:
					outputs = GetQrcodeScan()
				case define.TriggerOfficialSubscribeUnScribe:
					outputs = GetSubscribeUnsubscribe()
				case define.TriggerOfficialMenuClick:
					outputs = GetMenuClick()
				}
			}
			//采集旧的变量映射数据
			variableMap := make(map[string]string)
			for _, output := range trigger.Outputs {
				variableMap[output.Key] = output.Variable
			}
			//将变量映替换到新的配置上
			for idx, output := range outputs {
				if variable, ok := variableMap[output.Key]; ok {
					outputs[idx].Variable = variable
				}
			}
			params.Start.TriggerList[i].Outputs = outputs
			//补充开始节点自定义全局变量
			for _, output := range outputs {
				FillDiyGlobalBlanks(output, &params.Start)
			}
		}
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
	if params.ParamsExtractor.Output == nil {
		params.ParamsExtractor.Output = make(common.RecurveFields, 0)
	}
	if params.FormInsert.Datas == nil {
		params.FormInsert.Datas = make([]FormFieldValue, 0)
	}
	if params.FormDelete.Where == nil {
		params.FormDelete.Where = make([]define.FormFilterCondition, 0)
	}
	if params.FormUpdate.Where == nil {
		params.FormUpdate.Where = make([]define.FormFilterCondition, 0)
	}
	if params.FormUpdate.Datas == nil {
		params.FormUpdate.Datas = make([]FormFieldValue, 0)
	}
	if params.FormSelect.Where == nil {
		params.FormSelect.Where = make([]define.FormFilterCondition, 0)
	}
	if params.FormSelect.Fields == nil {
		params.FormSelect.Fields = make([]FormFieldTyp, 0)
	}
	if params.FormSelect.Order == nil {
		params.FormSelect.Order = make([]FormFieldOrder, 0)
	}
	if params.CodeRun.Params == nil {
		params.CodeRun.Params = make([]CodeRunParams, 0)
	}
	if params.CodeRun.Output == nil {
		params.CodeRun.Output = make(common.RecurveFields, 0)
	}
	if params.Mcp.Arguments == nil {
		params.Mcp.Arguments = make(map[string]any)
	}
	if params.Plugin.Output == nil {
		params.Plugin.Output = make(common.RecurveFields, 0)
	}

	return params
}

type WorkFlowNode struct {
	NodeType      common.MixedInt `json:"node_type"`
	NodeName      string          `json:"node_name"`
	NodeKey       string          `json:"node_key"`
	NodeParams    NodeParams      `json:"node_params"`
	NodeInfoJson  map[string]any  `json:"node_info_json"`
	NextNodeKey   string          `json:"next_node_key"`
	LoopParentKey string          `json:"loop_parent_key"`
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
	case NodeTypeLlm, NodeTypeReply, NodeTypeImmediatelyReply:
		variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, `special.llm_reply_content`))
	case NodeTypeQuestionOptimize:
		variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, `special.question_optimize_reply_content`))
	case NodeTypeParamsExtractor:
		for variable := range common.SimplifyFields(node.NodeParams.ParamsExtractor.Output) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //上一个节点,兼容旧数据
				variables = append(variables, variable)
			}
		}
	case NodeTypeFormSelect:
		for _, variable := range []string{`output_list`, `row_num`} {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //上一个节点,兼容旧数据
				variables = append(variables, variable)
			}
		}
	case NodeTypeCodeRun:
		for variable := range common.SimplifyFields(node.NodeParams.CodeRun.Output) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //上一个节点,兼容旧数据
				variables = append(variables, variable)
			}
		}
	case NodeTypeMcp:
		variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, `special.mcp_reply_content`))
	case NodeTypeLoop:
		loopOutput := node.NodeParams.Loop.Output
		formatOutput := make(common.RecurveFields, 0)
		for _, loopField := range loopOutput {
			formatRecurveField := common.RecurveField{}
			err := tool.JsonDecode(tool.JsonEncodeNoError(loopField), &formatRecurveField)
			if err != nil {
				logs.Error(`node type loop field decode ` + err.Error())
			}
			formatOutput = append(formatOutput, formatRecurveField)
		}
		for variable := range common.SimplifyFields(formatOutput) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //上一个节点,兼容旧数据
				variables = append(variables, variable)
			}
		}
	case NodeTypeBatch:
		batchOutput := node.NodeParams.Batch.Output
		formatOutput := make(common.RecurveFields, 0)
		for _, batchField := range batchOutput {
			formatRecurveField := common.RecurveField{}
			err := tool.JsonDecode(tool.JsonEncodeNoError(batchField), &formatRecurveField)
			if err != nil {
				logs.Error(`node type loop field decode ` + err.Error())
			}
			formatOutput = append(formatOutput, formatRecurveField)
		}
		for variable := range common.SimplifyFields(formatOutput) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //上一个节点,兼容旧数据
				variables = append(variables, variable)
			}
		}
	case NodeTypePlugin:
		for variable := range common.SimplifyFields(node.NodeParams.Plugin.Output) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //上一个节点,兼容旧数据
				variables = append(variables, variable)
			}
		}
	case NodeTypeImageGeneration:
		for variable := range common.SimplifyFields(node.NodeParams.ImageGeneration.Output) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //上一个节点,兼容旧数据
				variables = append(variables, variable)
			}
		}
	case NodeTypeJsonEncode:
		for variable := range common.SimplifyFields(node.NodeParams.JsonEncode.Outputs) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //上一个节点,兼容旧数据
				variables = append(variables, variable)
			}
		}
	case NodeTypeJsonDecode:
		simpleFields := make(common.SimpleFields)
		node.NodeParams.JsonDecode.Outputs.SimplifyFieldsDeep(&simpleFields, ``)
		for variable := range simpleFields {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //上一个节点,兼容旧数据
				variables = append(variables, variable)
			}
		}
	case NodeTypeTextToAudio:
		for variable := range common.SimplifyFields(node.NodeParams.TextToAudio.Output) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //上一个节点,兼容旧数据
				variables = append(variables, variable)
			}
		}
	case NodeTypeVoiceClone:
		for variable := range common.SimplifyFields(node.NodeParams.VoiceClone.Output) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //上一个节点,兼容旧数据
				variables = append(variables, variable)
			}
		}
	case NodeTypeLibraryImport:
		for variable := range common.SimplifyFields(node.NodeParams.LibraryImport.Outputs) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //上一个节点,兼容旧数据
				variables = append(variables, variable)
			}
		}
	case NodeTypeWorkflow:
		for variable := range common.SimplifyFields(node.NodeParams.Workflow.Output) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //上一个节点,兼容旧数据
				variables = append(variables, variable)
			}
		}
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

func CheckVariablePlaceholderExist(content string) bool {
	for _, item := range regexp.MustCompile(`【(([a-f0-9]{32}\.)?[a-zA-Z_][a-zA-Z0-9_\-.]*)】`).FindAllStringSubmatch(content, -1) {
		if len(item) > 1 {
			return true
		}
	}
	return false
}

func GetVariablePlaceholders(content string) []string {
	variables := make([]string, 0)
	for _, item := range regexp.MustCompile(`【(([a-f0-9]{32}\.)?[a-zA-Z_][a-zA-Z0-9_\-.]*)】`).FindAllStringSubmatch(content, -1) {
		if len(item) > 1 && !tool.InArrayString(item[1], variables) {
			variables = append(variables, item[1])
		}
	}
	return variables
}

func GetFirstVariable(content string) string {
	variables := GetVariablePlaceholders(content)
	variable := ``
	for _, v := range variables {
		if len(v) > 0 {
			variable = v
			break
		}
	}
	return variable
}

func RemoveVariablePlaceholders(content string) string {
	return regexp.MustCompile(`【(([a-f0-9]{32}\.)?[a-zA-Z_][a-zA-Z0-9_\-.]*)】`).ReplaceAllString(content, "")
}

func VerifyWorkFlowNodes(nodeList []WorkFlowNode, adminUserId int) (startNodeKey, modelConfigIds, libraryIds string, questionMultipleSwitch bool, err error) {
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
			for _, trigger := range node.NodeParams.Start.TriggerList {
				if trigger.TriggerType == TriggerTypeChat && trigger.TriggerSwitch && trigger.TriggerChatConfig.QuestionMultipleSwitch {
					questionMultipleSwitch = true //会话触发器-多模态输入-开启
				}
			}
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
		if node.NodeType == NodeTypeCodeRun {
			fromNodes.AddRelation(&nodeList[i], node.NodeParams.CodeRun.Exception)
		}
		if node.NodeType == NodeTypeWorkflow {
			fromNodes.AddRelation(&nodeList[i], node.NodeParams.Workflow.Exception)
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
		if node.LoopParentKey != `` {
			continue
		}
		if _, ok := fromNodes[node.NodeKey]; !ok {
			err = errors.New(`工作流存在游离节点:` + node.NodeName)
			return
		}
		//校验选择的变量必须存在
		err = verifyNode(adminUserId, node, fromNodes, nodeList)
		if err != nil {
			return
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
		case NodeTypeQuestionOptimize:
			modelConfigId = node.NodeParams.QuestionOptimize.ModelConfigId.Int()
		case NodeTypeParamsExtractor:
			modelConfigId = node.NodeParams.ParamsExtractor.ModelConfigId.Int()
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

func verifyNode(adminUserId int, node WorkFlowNode, fromNodes FromNodes, nodeList []WorkFlowNode) (err error) {
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
	case NodeTypeCate:
		if len(node.NodeParams.Cate.QuestionValue) > 0 && !tool.InArrayString(node.NodeParams.Cate.QuestionValue, variables) {
			err = errors.New(node.NodeName + `节点问题变量不存在:` + node.NodeParams.Cate.QuestionValue)
			return
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
	case NodeTypeLibs:
		if len(node.NodeParams.Libs.QuestionValue) > 0 && !tool.InArrayString(node.NodeParams.Libs.QuestionValue, variables) {
			err = errors.New(node.NodeName + `节点问题变量不存在:` + node.NodeParams.Libs.QuestionValue)
			return
		}
	case NodeTypeLlm:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.Llm.Prompt, variables); !ok {
			err = errors.New(node.NodeName + `节点提示词变量不存在:` + variable)
			return
		}
		if len(node.NodeParams.Llm.LibsNodeKey) > 0 {
			variable := fmt.Sprintf(`%s.%s`, node.NodeParams.Llm.LibsNodeKey, `special.lib_paragraph_list`)
			if !tool.InArrayString(variable, variables) {
				err = errors.New(node.NodeName + `节点的知识库引用选择的不是上级检索知识库节点`)
				return
			}
		}
		if len(node.NodeParams.Llm.QuestionValue) > 0 && !tool.InArrayString(node.NodeParams.Llm.QuestionValue, variables) {
			err = errors.New(node.NodeName + `节点问题变量不存在:` + node.NodeParams.Llm.QuestionValue)
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
	case NodeTypeQuestionOptimize:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.QuestionOptimize.Prompt, variables); !ok {
			err = errors.New(node.NodeName + `节点对话背景变量不存在:` + variable)
			return
		}
		if !tool.InArrayString(node.NodeParams.QuestionOptimize.QuestionValue, variables) {
			err = errors.New(node.NodeName + `节点问题变量不存在:` + node.NodeParams.QuestionOptimize.QuestionValue)
			return
		}
	case NodeTypeParamsExtractor:
		if !tool.InArrayString(node.NodeParams.ParamsExtractor.QuestionValue, variables) {
			err = errors.New(node.NodeName + `节点问题变量不存在:` + node.NodeParams.ParamsExtractor.QuestionValue)
			return
		}
	case NodeTypeFormInsert:
		for _, field := range node.NodeParams.FormInsert.Datas {
			if variable, ok := CheckVariablePlaceholder(field.Value, variables); !ok {
				err = errors.New(node.NodeName + `节点插入的变量不存在:` + variable)
				return
			}
		}
	case NodeTypeFormDelete:
		for _, field := range node.NodeParams.FormDelete.Where {
			if variable, ok := CheckVariablePlaceholder(field.RuleValue1, variables); !ok {
				err = errors.New(node.NodeName + `节点插入的变量不存在:` + variable)
				return
			}
			if variable, ok := CheckVariablePlaceholder(field.RuleValue2, variables); !ok {
				err = errors.New(node.NodeName + `节点插入的变量不存在:` + variable)
				return
			}
		}
	case NodeTypeFormUpdate:
		for _, field := range node.NodeParams.FormUpdate.Where {
			if variable, ok := CheckVariablePlaceholder(field.RuleValue1, variables); !ok {
				err = errors.New(node.NodeName + `节点插入的变量不存在:` + variable)
				return
			}
			if variable, ok := CheckVariablePlaceholder(field.RuleValue2, variables); !ok {
				err = errors.New(node.NodeName + `节点插入的变量不存在:` + variable)
				return
			}
		}
		for _, field := range node.NodeParams.FormUpdate.Datas {
			if variable, ok := CheckVariablePlaceholder(field.Value, variables); !ok {
				err = errors.New(node.NodeName + `节点插入的变量不存在:` + variable)
				return
			}
		}
	case NodeTypeFormSelect:
		for _, field := range node.NodeParams.FormSelect.Where {
			if variable, ok := CheckVariablePlaceholder(field.RuleValue1, variables); !ok {
				err = errors.New(node.NodeName + `节点插入的变量不存在:` + variable)
				return
			}
			if variable, ok := CheckVariablePlaceholder(field.RuleValue2, variables); !ok {
				err = errors.New(node.NodeName + `节点插入的变量不存在:` + variable)
				return
			}
		}
	case NodeTypeCodeRun:
		for _, param := range node.NodeParams.CodeRun.Params {
			if !tool.InArrayString(param.Variable, variables) {
				err = errors.New(node.NodeName + `节点选择的变量不存在:` + param.Variable)
				return
			}
		}
	case NodeTypeMcp:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.Mcp.ToolName, variables); !ok {
			err = errors.New(node.NodeName + `节点MCP工具变量不存在:` + variable)
			return
		}
	case NodeTypePlugin:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.Plugin.Name, variables); !ok {
			err = errors.New(node.NodeName + `节点选择的变量不存在:` + variable)
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.Plugin.Type, variables); !ok {
			err = errors.New(node.NodeName + `节点选择的变量不存在:` + variable)
			return
		}
	case NodeTypeTextToAudio:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.TextToAudio.Arguments.Text, variables); !ok {
			err = errors.New(node.NodeName + `节点内容变量不存在:` + variable)
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.TextToAudio.Arguments.VoiceSetting.VoiceId, variables); !ok {
			err = errors.New(node.NodeName + `节点内容变量不存在:` + variable)
			return
		}
	case NodeTypeVoiceClone:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.VoiceClone.Arguments.FileUrl, variables); !ok {
			err = errors.New(node.NodeName + `节点内容变量不存在:` + variable)
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.VoiceClone.Arguments.VoiceId, variables); !ok {
			err = errors.New(node.NodeName + `节点内容变量不存在:` + variable)
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.VoiceClone.Arguments.ClonePrompt.PromptAudioUrl, variables); !ok {
			err = errors.New(node.NodeName + `节点内容变量不存在:` + variable)
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.VoiceClone.Arguments.ClonePrompt.PromptText, variables); !ok {
			err = errors.New(node.NodeName + `节点内容变量不存在:` + variable)
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.VoiceClone.Arguments.Text, variables); !ok {
			err = errors.New(node.NodeName + `节点内容变量不存在:` + variable)
			return
		}
	case NodeTypeWorkflow:
		for _, param := range node.NodeParams.Workflow.Params {
			if variable, ok := CheckVariablePlaceholder(param.Variable, variables); !ok {
				err = errors.New(node.NodeName + `节点内容变量不存在:` + variable)
				return
			}
		}
	case NodeTypeLoop:
		//verify loop arrays
		if node.NodeParams.Loop.LoopType == common.LoopTypeArray {
			if len(node.NodeParams.Loop.LoopArrays) == 0 {
				err = errors.New(fmt.Sprintf(`%s未添加循环数组`, node.NodeName))
				return
			}
			//验证对前置节点的引用
			bFind := VerityLoopParams(node.NodeParams.Loop.LoopArrays, nodeList)
			if !bFind {
				err = errors.New(fmt.Sprintf(`%s未找到循环数组引用的节点输出<array>参数`, node.NodeName))
				return
			}
		} else if node.NodeParams.Loop.LoopNumber.Int() <= 0 {
			err = errors.New(fmt.Sprintf(`%s循环次数必须大于0`, node.NodeName))
			return
		}
		//验证中间变量初始值是否存在
		bFind := VerityLoopParams(node.NodeParams.Loop.IntermediateParams, nodeList)
		if !bFind {
			err = errors.New(fmt.Sprintf(`%s未找到中间变量的引用`, node.NodeName))
			return
		}
		//child node
		childNodes := make([]WorkFlowNode, 0)
		for _, vNode := range nodeList {
			if vNode.LoopParentKey == node.NodeKey {
				childNodes = append(childNodes, vNode)
			}
		}
		//验证输出是否来自于中间变量
		bFind = VerityLoopParams(node.NodeParams.Loop.Output, []WorkFlowNode{node})
		if !bFind {
			err = errors.New(fmt.Sprintf(`%s输出节点中未找到引用的中间变量参数`, node.NodeName))
			return
		}
		//验证子节点
		_, _, _, err = VerityLoopWorkflowNodes(adminUserId, node, childNodes, LoopAllowNodeTypes, `循环节点`)
		if err != nil {
			return
		}
	case NodeTypeBatch:
		if len(node.NodeParams.Batch.BatchArrays) == 0 {
			err = errors.New(fmt.Sprintf(`%s未添加执行数组`, node.NodeName))
			return
		}
		//验证执行数组的初始值是否存在
		bFind := VerityLoopParams(node.NodeParams.Batch.BatchArrays, nodeList)
		if !bFind {
			err = errors.New(fmt.Sprintf(`未找到【%s】执行数组引用的变量`, node.NodeName))
			return
		}
		childNodes := make([]WorkFlowNode, 0)
		for _, vNode := range nodeList {
			if vNode.LoopParentKey == node.NodeKey {
				childNodes = append(childNodes, vNode)
			}
		}
		//验证输出值是否存在
		bFind = VerityLoopParams(node.NodeParams.Batch.Output, childNodes)
		if !bFind {
			err = errors.New(fmt.Sprintf(`未找到【%s】的输出中引用的子节点输出`, node.NodeName))
			return
		}
		//验证子节点
		_, _, _, err = VerityLoopWorkflowNodes(adminUserId, node, childNodes, BatchAllowNodeTypes, `批处理`)
		if err != nil {
			return
		}
	case NodeTypeFinish:
		if node.NodeParams.Finish.OutType == define.FinishNodeOutTypeMessage {
			for _, field := range node.NodeParams.Finish.Messages {
				if variable, ok := CheckVariablePlaceholder(field.Content, variables); !ok {
					err = errors.New(node.NodeName + fmt.Sprintf(`变量(%s)不存在:`, variable))
					return
				}
			}
		}

	case NodeTypeImageGeneration:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.ImageGeneration.Prompt, variables); !ok {
			err = errors.New(node.NodeName + `节点提示词选择的变量不存在:` + variable)
			return
		}
		for _, imageUrl := range node.NodeParams.ImageGeneration.InputImages {
			if variable, ok := CheckVariablePlaceholder(imageUrl, variables); !ok {
				err = errors.New(node.NodeName + `节点输入图片选择的变量不存在:` + variable)
				return
			}
		}
	case NodeTypeJsonEncode:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.JsonEncode.InputVariable, variables); !ok {
			err = errors.New(node.NodeName + `节点选择的输入变量不存在:` + variable)
			return
		}
	case NodeTypeJsonDecode:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.JsonDecode.InputVariable, variables); !ok {
			err = errors.New(node.NodeName + `节点选择的输入变量不存在:` + variable)
			return
		}
	case NodeTypeLibraryImport:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.LibraryImport.NormalTitle, variables); !ok {
			err = errors.New(node.NodeName + `节点选择的变量不存在:` + variable)
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.LibraryImport.NormalUrl, variables); !ok {
			err = errors.New(node.NodeName + `节点选择的变量不存在:` + variable)
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.LibraryImport.NormalContent, variables); !ok {
			err = errors.New(node.NodeName + `节点选择的变量不存在:` + variable)
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.LibraryImport.QaQuestion, variables); !ok {
			err = errors.New(node.NodeName + `节点选择的变量不存在:` + variable)
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.LibraryImport.QaAnswer, variables); !ok {
			err = errors.New(node.NodeName + `节点选择的变量不存在:` + variable)
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.LibraryImport.QaSimilarQuestionVariable, variables); !ok {
			err = errors.New(node.NodeName + `节点选择的变量不存在:` + variable)
		}
	case NodeTypeImmediatelyReply:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.ImmediatelyReply.Content, variables); !ok {
			err = errors.New(node.NodeName + `节点插入的变量不存在:` + variable)
			return
		}
	}
	return nil
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
	if len(node.NextNodeKey) == 0 && node.LoopParentKey == `` && !tool.InArrayInt(node.NodeType.Int(), []int{NodeTypeRemark, NodeTypeEdges, NodeTypeFinish, NodeTypeManual, NodeTypeLoopEnd}) {
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
		err = node.NodeParams.Assign.Verify(node)
	case NodeTypeReply:
		err = node.NodeParams.Reply.Verify()
	case NodeTypeManual:
		err = node.NodeParams.Manual.Verify(adminUserId)
	case NodeTypeQuestionOptimize:
		err = node.NodeParams.QuestionOptimize.Verify(adminUserId)
	case NodeTypeParamsExtractor:
		err = node.NodeParams.ParamsExtractor.Verify(adminUserId)
	case NodeTypeFormInsert:
		err = node.NodeParams.FormInsert.Verify(adminUserId)
	case NodeTypeFormDelete:
		err = node.NodeParams.FormDelete.Verify(adminUserId)
	case NodeTypeFormUpdate:
		err = node.NodeParams.FormUpdate.Verify(adminUserId)
	case NodeTypeFormSelect:
		err = node.NodeParams.FormSelect.Verify(adminUserId)
	case NodeTypeCodeRun:
		err = node.NodeParams.CodeRun.Verify()
	case NodeTypeMcp:
		err = node.NodeParams.Mcp.Verify(adminUserId)
	case NodeTypeLoop:
		err = node.NodeParams.Loop.Verify(node.NodeName)
	case NodeTypePlugin:
		err = node.NodeParams.Plugin.Verify(adminUserId)
	case NodeTypeBatch:
		err = node.NodeParams.Batch.Verify(node.NodeName)
	case NodeTypeFinish:
		err = node.NodeParams.Finish.Verify(node.NodeName)
	case NodeTypeImageGeneration:
		err = node.NodeParams.ImageGeneration.Verify(node.NodeName)
	case NodeTypeJsonEncode:
		err = node.NodeParams.JsonEncode.Verify(node.NodeName)
	case NodeTypeJsonDecode:
		err = node.NodeParams.JsonDecode.Verify(node.NodeName)
	case NodeTypeTextToAudio:
		err = node.NodeParams.TextToAudio.Verify()
	case NodeTypeVoiceClone:
		err = node.NodeParams.VoiceClone.Verify(adminUserId)
	case NodeTypeLibraryImport:
		err = node.NodeParams.LibraryImport.Verify(adminUserId)
	case NodeTypeWorkflow:
		err = node.NodeParams.Workflow.Verify(adminUserId)
	case NodeTypeImmediatelyReply:
		err = node.NodeParams.ImmediatelyReply.Verify()
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
		if !tool.InArrayString(item.Typ, []string{common.TypString, common.TypNumber, common.TypArrString, common.TypArrObject, common.TypBoole, common.TypArrNumber}) {
			return errors.New(fmt.Sprintf(`自定义全局变量类型不支持:%s`, item.Key))
		}
		if _, ok := maps[item.Key]; ok {
			return errors.New(fmt.Sprintf(`自定义全局变量名重复定义:%s`, item.Key))
		}
		maps[item.Key] = struct{}{}
	}
	//触发器参数校验
	if len(params.TriggerList) == 0 {
		return errors.New(`工作流至少添加一个触发器`)
	}
	for i, trigger := range params.TriggerList {
		for _, output := range trigger.Outputs {
			if len(output.Variable) == 0 {
				continue //触发器输出支持不配置变量映射
			}
			key, _ := strings.CutPrefix(output.Variable, `global.`)
			if _, ok := maps[key]; !ok {
				return errors.New(fmt.Sprintf(`第%d个触发器(%s)的输出变量%s(%s)映射配置错误`, i+1, trigger.TriggerName, output.Key, output.Desc))
			}
		}
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
	if len(params.QuestionValue) == 0 {
		return errors.New(`用户问题不能为空`)
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
	if err := common.CheckRrfWeight(params.RrfWeight, define.LangZhCn); err != nil {
		return err
	}
	if params.TopK <= 0 || params.TopK > 500 {
		return errors.New(`知识库检索TopK范围1~500`)
	}
	if params.Similarity < 0 || params.Similarity > 1 {
		return errors.New(`知识库检索相似度阈值0~1`)
	}
	if params.RerankStatus > 0 || params.RerankModelConfigId != 0 || len(params.RerankUseModel) > 0 {
		if params.RerankModelConfigId <= 0 || len(params.RerankUseModel) == 0 {
			return errors.New(`请选择使用的Rerank模型`)
		}
		if ok := common.CheckModelIsValid(adminUserId, params.RerankModelConfigId.Int(), params.RerankUseModel, common.Rerank); !ok {
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
	if len(params.QuestionValue) == 0 {
		return errors.New(`用户问题不能为空`)
	}
	if len(params.LibsNodeKey) > 0 && !common.IsMd5Str(params.LibsNodeKey) {
		return errors.New(`知识库引用节点参数格式错误`)
	}
	return nil
}

func (params *AssignNodeParams) Verify(node *WorkFlowNode) error {
	if params == nil || len(*params) == 0 {
		return errors.New(`配置参数不能为空`)
	}
	for i, param := range *params {
		if len(param.Variable) == 0 {
			return errors.New(fmt.Sprintf(`第%d行:请选择变量`, i+1))
		}
		if node.LoopParentKey == `` && (!strings.HasPrefix(param.Variable, `global.`) || !common.IsVariableNames(param.Variable)) {
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

func checkFormId(adminUserId, formId int) error {
	if formId <= 0 {
		return errors.New(`未选择操作的数据表`)
	}
	form, err := msql.Model(`form`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(formId)).Where(`delete_time`, `0`).Field(`id`).Find()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	if len(form) == 0 {
		return errors.New(`数据表信息不存在`)
	}
	return nil
}

func checkFormDatas(adminUserId, formId int, datas []FormFieldValue) error {
	if len(datas) == 0 {
		return errors.New(`字段列表不能为空`)
	}
	fields, err := msql.Model(`form_field`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`form_id`, cast.ToString(formId)).ColumnMap(`type,required,description`, `name`)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	maps := map[string]struct{}{}
	for i, data := range datas {
		if len(data.Name) == 0 {
			return errors.New(fmt.Sprintf(`第%d行:字段名参数不能为空`, i+1))
		}
		field, ok := fields[data.Name]
		if !ok {
			return errors.New(fmt.Sprintf(`第%d行:字段名%s不存在于数据表`, i+1, data.Name))
		}
		if len(data.Type) == 0 || data.Type != field[`type`] {
			return errors.New(fmt.Sprintf(`第%d行:字段名%s(%s)类型与数据表不一致`, i+1, data.Name, field[`description`]))
		}
		if len(data.Value) == 0 && cast.ToBool(field[`required`]) {
			return errors.New(fmt.Sprintf(`第%d行:字段名%s(%s)为必要字段,不能为空`, i+1, data.Name, field[`description`]))
		}
		if _, ok := maps[data.Name]; ok {
			return errors.New(fmt.Sprintf(`第%d行:字段名%s(%s)重复出现在字段列表`, i+1, data.Name, field[`description`]))
		}
		maps[data.Name] = struct{}{}
	}
	return nil
}

func checkFormWhere(adminUserId, formId int, where []define.FormFilterCondition) error {
	if len(where) == 0 {
		return errors.New(`条件列表不能为空`)
	}
	fields, err := msql.Model(`form_field`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`form_id`, cast.ToString(formId)).ColumnMap(`name,type,description`, `id`)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	fields[`0`] = msql.Params{`name`: `id`, `type`: `integer`, `description`: `ID`} //追加一个ID,用于兼容处理
	for i, condition := range where {
		if condition.FormFieldId < 0 { //特比注意,这里可以等于0
			return errors.New(fmt.Sprintf(`第%d行:选择字段参数非法`, i+1))
		}
		field, ok := fields[cast.ToString(condition.FormFieldId)]
		if !ok {
			return errors.New(fmt.Sprintf(`第%d行:选择的字段不存在于数据表`, i+1))
		}
		if err = condition.Check(field[`type`], true); err != nil {
			return errors.New(fmt.Sprintf(`第%d行:字段名%s(%s)校验错误:%s`, i+1, field[`name`], field[`description`], err.Error()))
		}
	}
	return nil
}

func checkFormFields(adminUserId, formId int, Fields []FormFieldTyp) error {
	if len(Fields) == 0 {
		return errors.New(`字段列表不能为空`)
	}
	fields, err := msql.Model(`form_field`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`form_id`, cast.ToString(formId)).ColumnMap(`type,description`, `name`)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	maps := map[string]struct{}{}
	for i, data := range Fields {
		if len(data.Name) == 0 {
			return errors.New(fmt.Sprintf(`第%d行:字段名参数不能为空`, i+1))
		}
		field, ok := fields[data.Name]
		if !ok {
			return errors.New(fmt.Sprintf(`第%d行:字段名%s不存在于数据表`, i+1, data.Name))
		}
		if len(data.Type) == 0 || data.Type != field[`type`] {
			return errors.New(fmt.Sprintf(`第%d行:字段名%s(%s)类型与数据表不一致`, i+1, data.Name, field[`description`]))
		}
		if _, ok := maps[data.Name]; ok {
			return errors.New(fmt.Sprintf(`第%d行:字段名%s(%s)重复出现在字段列表`, i+1, data.Name, field[`description`]))
		}
		maps[data.Name] = struct{}{}
	}
	return nil
}

func (params *FormInsertNodeParams) Verify(adminUserId int) error {
	if err := checkFormId(adminUserId, params.FormId.Int()); err != nil {
		return err
	}
	if err := checkFormDatas(adminUserId, params.FormId.Int(), params.Datas); err != nil {
		return err
	}
	return nil
}

func (params *FormDeleteNodeParams) Verify(adminUserId int) error {
	if err := checkFormId(adminUserId, params.FormId.Int()); err != nil {
		return err
	}
	if !tool.InArrayInt(params.Typ.Int(), []int{1, 2}) {
		return errors.New(`条件之间关系参数错误`)
	}
	if err := checkFormWhere(adminUserId, params.FormId.Int(), params.Where); err != nil {
		return err
	}
	return nil
}

func (params *FormUpdateNodeParams) Verify(adminUserId int) error {
	if err := checkFormId(adminUserId, params.FormId.Int()); err != nil {
		return err
	}
	if !tool.InArrayInt(params.Typ.Int(), []int{1, 2}) {
		return errors.New(`条件之间关系参数错误`)
	}
	if err := checkFormWhere(adminUserId, params.FormId.Int(), params.Where); err != nil {
		return err
	}
	if err := checkFormDatas(adminUserId, params.FormId.Int(), params.Datas); err != nil {
		return err
	}
	return nil
}

func (params *FormSelectNodeParams) Verify(adminUserId int) error {
	if err := checkFormId(adminUserId, params.FormId.Int()); err != nil {
		return err
	}
	if !tool.InArrayInt(params.Typ.Int(), []int{1, 2}) {
		return errors.New(`条件之间关系参数错误`)
	}
	if err := checkFormWhere(adminUserId, params.FormId.Int(), params.Where); err != nil {
		return err
	}
	if err := checkFormFields(adminUserId, params.FormId.Int(), params.Fields); err != nil {
		return err
	}
	for _, order := range params.Order {
		if !tool.InArrayString(order.Name, []string{`id`, `create_time`, `update_time`}) {
			return fmt.Errorf(`不支持%s用于排序操作`, order.Name)
		}
	}
	if params.Limit <= 0 || params.Limit > 1000 {
		return errors.New(`查询数量范围:1~1000`)
	}
	return nil
}

func (params *CodeRunNodeParams) Verify() error {
	maps := map[string]struct{}{}
	for idx, param := range params.Params {
		if !common.IsVariableName(param.Field) {
			return errors.New(fmt.Sprintf(`自定义输入参数KEY格式错误:%s`, param.Field))
		}
		if len(param.Variable) == 0 {
			return errors.New(fmt.Sprintf(`第%d个自定义输入参数:请选择变量`, idx+1))
		}
		if !common.IsVariableNames(param.Variable) {
			return errors.New(fmt.Sprintf(`第%d个自定义输入参数:变量格式错误`, idx+1))
		}
		if _, ok := maps[param.Field]; ok {
			return errors.New(fmt.Sprintf(`自定义输入参数KEY重复定义:%s`, param.Field))
		}
		maps[param.Field] = struct{}{}
	}
	ok, err := regexp.MatchString(`function\s+main\s*\(.*\)\s*\{`, params.MainFunc)
	if err != nil || !ok {
		return errors.New(`JavaScript代码缺少main函数`)
	}
	if params.Timeout < 1 || params.Timeout > 60 {
		return errors.New(`代码运行超时时间范围1~60秒`)
	}
	if err = params.Output.Verify(); err != nil {
		return err
	}
	if len(params.Exception) == 0 || !common.IsMd5Str(params.Exception) {
		return errors.New(`异常处理:下一个节点未指定或格式错误`)
	}
	return nil
}

func (params *McpNodeParams) Verify(adminUserId int) error {
	info, err := msql.Model(`mcp_provider`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(params.ProviderId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	if len(info) == 0 {
		return errors.New(`请选择mcp工具`)
	}
	if cast.ToInt(info[`has_auth`]) != 1 {
		return errors.New(`请先授权mcp工具`)
	}
	var tools []mcp.Tool
	err = json.Unmarshal([]byte(info[`tools`]), &tools)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	if len(tools) == 0 {
		return errors.New(`未找到可用工具`)
	}

	var mcpTool *mcp.Tool
	for _, t := range tools {
		if t.Name == params.ToolName {
			mcpTool = &t
			break
		}
	}
	if mcpTool == nil {
		return errors.New(`匹配到工具`)
	}

	if err := common.ValidateMcpToolArguments(*mcpTool, params.Arguments); err != nil {
		return err
	}

	return nil
}

func (params *LoopNodeParams) Verify(nodeName string) error {
	if !tool.InArray(params.LoopType, []string{common.LoopTypeArray, common.LoopTypeNumber}) {
		return errors.New(nodeName + `循环类型错误`)
	}
	if params.LoopType == common.LoopTypeArray {
		if len(params.LoopArrays) == 0 {
			return errors.New(nodeName + `循环数组不能为空`)
		}
	} else {
		if params.LoopNumber <= 0 {
			return errors.New(nodeName + `循环数字必须大于0`)
		}
	}
	return nil
}

// VerityLoopParams 验证输出或者循环参数
func VerityLoopParams(fields []common.LoopField, nodeList []WorkFlowNode) bool {
	for _, field := range fields {
		if field.Value == `` {
			continue
		}
		if field.Value != `` && CheckVariablePlaceholderExist(field.Value) && FindNodeByUseKey(nodeList, field.Value) == nil {
			return false
		}
	}
	return true
}

func VerityLoopWorkflowNodes(adminUserId int, loopNode WorkFlowNode, nodeList []WorkFlowNode, allowNodeTypes []int, nodeTypeDesc string) (startNodeKey, modelConfigIds, libraryIds string, err error) {
	startNodeCount, finishNodeCount := 0, 0
	fromNodes := make(FromNodes)
	for i, node := range nodeList {
		if !tool.InArrayInt(node.NodeType.Int(), allowNodeTypes) {
			err = errors.New(fmt.Sprintf(nodeTypeDesc+`的子节点类型错误 %d`, node.NodeType))
			return
		}
		//node base verify
		if err = node.Verify(adminUserId); err != nil {
			return
		}
		//start node
		if tool.InArray(node.NodeType.Int(), []int{NodeTypeLoopStart, NodeTypeBatchStart}) {
			startNodeKey = node.NodeKey
			startNodeCount++
		}
		if node.NodeType.Int() != NodeTypeLoopEnd {
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
		if node.NodeType == NodeTypeCodeRun {
			fromNodes.AddRelation(&nodeList[i], node.NodeParams.CodeRun.Exception)
		}
		if node.NodeType == NodeTypeWorkflow {
			fromNodes.AddRelation(&nodeList[i], node.NodeParams.Workflow.Exception)
		}
		if node.NextNodeKey == loopNode.NodeKey {
			finishNodeCount++
		}
	}
	if startNodeCount != 1 {
		err = errors.New(nodeTypeDesc + `仅能有一个入口节点`)
		return
	}
	if finishNodeCount == 0 {
		//err = errors.New(`工作流必须存在一个终止循环或出口节点`)
		//return
	}
	for _, node := range nodeList {
		//remark node continue
		skipNodeTypes := []int{
			NodeTypeRemark,
			NodeTypeLoopStart,
			NodeTypeBatchStart,
		}
		if tool.InArrayInt(node.NodeType.Int(), skipNodeTypes) {
			continue
		}
		if _, ok := fromNodes[node.NodeKey]; !ok {
			err = errors.New(nodeTypeDesc + `中存在游离节点:` + node.NodeName)
			return
		}
		//暂时不进行验证 循环节点单独执行会缺少前面所有节点的输入，代码难处理
		//err = verifyNode(adminUserId, node, fromNodes, make([]WorkFlowNode, 0))
		//if err != nil {
		//	return
		//}
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
		case NodeTypeQuestionOptimize:
			modelConfigId = node.NodeParams.QuestionOptimize.ModelConfigId.Int()
		case NodeTypeParamsExtractor:
			modelConfigId = node.NodeParams.ParamsExtractor.ModelConfigId.Int()
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
func (param *PluginNodeParams) Verify(adminUserId int) error {
	u := define.Config.Plugin[`endpoint`] + "/manage/plugin/local-plugins/run"

	if param.Type == `notice` {
		resp := &lib_web.Response{}
		request := curl.Post(u).Header(`admin_user_id`, cast.ToString(adminUserId))
		request.Param("name", param.Name)
		request.Param("action", "default/get-schema")
		request.Param("params", "{}")
		err := request.ToJSON(resp)
		if err != nil {
			return err
		}
		if resp.Res != 0 {
			return errors.New(resp.Msg)
		}

		// 验证 schema
		schema, ok := resp.Data.(map[string]any)
		if !ok {
			return errors.New("invalid schema format")
		}

		for key, raw := range schema {
			rule, _ := raw.(map[string]any)
			req, _ := rule["required"].(bool)
			typ, _ := rule["type"].(string)
			val, exists := param.Params[key]
			if req && !exists {
				return fmt.Errorf("missing required field: %s", key)
			}
			if !exists {
				continue
			}

			switch typ {
			case "string":
				if _, ok := val.(string); !ok {
					return fmt.Errorf("field %s must be string", key)
				}
			case "number":
				if _, ok := val.(float64); !ok {
					return fmt.Errorf("field %s must be number", key)
				}
			case "boolean":
				if _, ok := val.(bool); !ok {
					return fmt.Errorf("field %s must be boolean", key)
				}
			default:
				return fmt.Errorf("unknown type for field %s: %s", key, typ)
			}
		}
	} else if param.Type == `extension` {
		resp := &lib_web.Response{}
		request := curl.Post(u).Header(`admin_user_id`, cast.ToString(adminUserId))
		request.Param("name", param.Name)
		request.Param("action", "default/check-schema")
		params, err := json.Marshal(param.Params)
		if err != nil {
			return err
		}
		//插件验证前替换占位符
		paramsStr := regexp.MustCompile(`【([a-f0-9]{32}\.)?[a-zA-Z_][a-zA-Z0-9_\-.]*】`).ReplaceAllString(string(params), `0`)
		//logs.Debug(paramsStr)
		err = request.Param("params", paramsStr).ToJSON(resp)
		if err != nil {
			return err
		}
		if resp.Res != 0 {
			logs.Debug("resp %v", resp.Res)
			return errors.New(resp.Msg)
		}
	}

	return nil
}

func (params *BatchNodeParams) Verify(nodeName string) error {
	if len(params.BatchArrays) == 0 {
		return errors.New(fmt.Sprintf(`【%s】执行数组不能为空`, nodeName))
	}
	if params.ChanNumber.Int() < 1 || params.ChanNumber.Int() > 10 {
		return errors.New(fmt.Sprintf(`【%s】执行并发数错误(1-10)`, nodeName))
	}
	if params.MaxRunNumber.Int() < 1 || params.MaxRunNumber.Int() > 500 {
		return errors.New(fmt.Sprintf(`【%s】执行最大运行数错误(1-500)`, nodeName))
	}
	return nil
}

func (params *FinishNodeParams) Verify(nodeName string) error {
	if len(params.OutType) > 0 && !tool.InArray(params.OutType, []string{define.FinishNodeOutTypeMessage, define.FinishNodeOutTypeVariable}) {
		return errors.New(fmt.Sprintf(`【%s】输出类型错误`, nodeName))
	}
	return nil
}

func (params *ImageGenerationParams) Verify(nodeName string) error {
	if cast.ToInt(params.ImageNum) < 0 || cast.ToInt(params.ImageNum) > 15 {
		return errors.New(fmt.Sprintf(`【%s】图片数量错误(0-15)`, nodeName))
	}
	if !tool.InArray(params.Size, define.ImageSizes) {
		return errors.New(fmt.Sprintf(`【%s】图片尺寸错误`, nodeName))
	}
	return nil
}

func (params *JsonEncodeParams) Verify(nodeName string) error {
	if len(params.InputVariable) == 0 {
		return errors.New(fmt.Sprintf(`【%s】缺少输入`, nodeName))
	}
	return nil
}

func (params *JsonDecodeParams) Verify(nodeName string) error {
	if len(params.InputVariable) == 0 {
		return errors.New(fmt.Sprintf(`【%s】缺少输入`, nodeName))
	}
	return nil
}

func (params *TextToAudioNodeParams) Verify() error {
	if params.Arguments.ModelId <= 0 {
		return errors.New(`请选择模型配置`)
	}

	validVoiceTypes := []string{"system", "voice_cloning", "voice_generation", "all"}
	if len(params.VoiceType) > 0 && !tool.InArrayString(params.VoiceType, validVoiceTypes) {
		return errors.New(`voiceType参数错误，可选值: system, voice_cloning, voice_generation, all`)
	}

	// 验证Text，最大10000字符
	if len(params.Arguments.Text) == 0 {
		return errors.New(`text内容不能为空`)
	}
	if len(params.Arguments.Text) > 10000 {
		return errors.New(`text内容长度不能超过10000字符`)
	}

	// 验证VoiceSetting
	voiceSetting := params.Arguments.VoiceSetting
	if len(voiceSetting.VoiceId) == 0 {
		return errors.New(`voice_setting.voice_id不能为空`)
	}

	// 验证Speed - 建议范围[0.5, 2.0]
	if voiceSetting.Speed < 0 || voiceSetting.Speed > 100 {
		return errors.New(`voice_setting.speed取值范围建议0.5~2.0`)
	}

	// 验证Vol - 音量范围
	if voiceSetting.Vol < 0 || voiceSetting.Vol > 100 {
		return errors.New(`voice_setting.vol取值范围0~100`)
	}

	// 验证Pitch - 音调范围
	if voiceSetting.Pitch < -100 || voiceSetting.Pitch > 100 {
		return errors.New(`voice_setting.pitch取值范围-100~100`)
	}

	// 验证AudioSetting - 音频设置
	audioSetting := params.Arguments.AudioSetting
	// 验证SampleRate - 采样率
	if audioSetting.SampleRate > 0 {
		validSampleRates := []int{8000, 16000, 22050, 24000, 32000, 44100}
		if !tool.InArrayInt(audioSetting.SampleRate, validSampleRates) {
			return errors.New(`audio_setting.sample_rate必须是以下值之一: 8000, 16000, 22050, 24000, 32000, 44100`)
		}
	}

	// 验证Bitrate - 比特率
	if audioSetting.Bitrate > 0 {
		validBitrates := []int{32000, 64000, 128000, 256000}
		if !tool.InArrayInt(audioSetting.Bitrate, validBitrates) {
			return errors.New(`audio_setting.bitrate必须是以下值之一: 32000, 64000, 128000, 256000`)
		}
	}

	// 验证Format - 音频格式
	if len(audioSetting.Format) > 0 {
		validFormats := []string{"mp3", "pcm", "flac", "wav"}
		if !tool.InArrayString(audioSetting.Format, validFormats) {
			return errors.New(`audio_setting.format必须是以下值之一: mp3, pcm, flac, wav`)
		}
	}

	// 验证Channel - 声道数
	if audioSetting.Channel > 0 {
		if audioSetting.Channel != 1 && audioSetting.Channel != 2 {
			return errors.New(`audio_setting.channel必须是1(单声道)或2(双声道)`)
		}
	}

	// 验证LanguageBoost
	if len(params.Arguments.LanguageBoost) > 0 {
		validLanguageBoosts := []string{"auto", "Chinese", "English", "Japanese", "Korean", "French", "German", "Spanish", "Russian", "Arabic"}
		if !tool.InArrayString(params.Arguments.LanguageBoost, validLanguageBoosts) {
			return errors.New(`language_boost参数错误`)
		}
	}

	return nil
}

func (param *VoiceCloneNodeParams) Verify(adminUserId int) error {
	if len(param.Arguments.FileUrl) == 0 {
		return errors.New(`file_url不能为空`)
	}
	// 验证VoiceId - MiniMax要求：长度8-256，首字符为字母，允许数字字母-_，末位不可为-_
	if len(param.Arguments.VoiceId) == 0 {
		return errors.New(`voice_id不能为空`)
	}

	// 验证ClonePrompt - 可选参数，但如果提供则需要完整
	if len(param.Arguments.ClonePrompt.PromptAudioUrl) > 0 || len(param.Arguments.ClonePrompt.PromptText) > 0 {
		// 如果提供了示例音频URL，验证格式
		if len(param.Arguments.ClonePrompt.PromptAudioUrl) == 0 {
			return errors.New(`提供了prompt_text则prompt_audio_url不能为空`)
		}
		if _, err := url.Parse(param.Arguments.ClonePrompt.PromptAudioUrl); err != nil {
			return errors.New(`clone_prompt.prompt_audio_url格式错误`)
		}

		// 如果提供了示例音频，验证文本
		if len(param.Arguments.ClonePrompt.PromptText) == 0 {
			return errors.New(`提供了prompt_audio_url则prompt_text不能为空`)
		}
	}

	return nil
}

func (param *LibraryImportParams) Verify(adminUserId int) error {
	if cast.ToInt(param.LibraryId) == 0 {
		return errors.New(`请选择导入的知识库`)
	}
	libraryInfo, err := common.GetLibraryInfo(cast.ToInt(param.LibraryId), adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return errors.New(`获取知识库明细失败`)
	}
	if len(libraryInfo) == 0 {
		return errors.New(`知识库不存在`)
	}
	if cast.ToInt(param.LibraryGroupId) > 0 {
		libraryGroupInfo, sqlErr := msql.Model(`chat_ai_library_group`, define.Postgres).
			Where(`library_id`, cast.ToString(param.LibraryId)).
			Where(`id`, param.LibraryGroupId).Find()
		if sqlErr != nil {
			logs.Error(sqlErr.Error())
			return errors.New(`获取知识库分组明细失败`)
		}
		if len(libraryGroupInfo) == 0 {
			return errors.New(`知识库分组不存在`)
		}
	}
	if cast.ToInt(libraryInfo[`type`]) == define.GeneralLibraryType {
		if param.ImportType == define.LibraryImportContent {
			if param.NormalTitle == `` || param.NormalContent == `` {
				return errors.New(`请填写导入内容和标题`)
			}
		} else if param.ImportType == define.LibraryImportUrl {
			if param.NormalUrl == `` {
				return errors.New(`请填写导入的URL`)
			}
		}
	} else if cast.ToInt(libraryInfo[`type`]) == define.QALibraryType {
		if param.QaQuestion == `` || param.QaAnswer == `` {
			return errors.New(`请填写导入的问题和答案`)
		}
	}
	return nil
}

func (param *WorkflowNodeParams) Verify(adminUserId int) error {
	if param.RobotId <= 0 {
		return errors.New(`robot_id不能为空`)
	}
	info, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(param.RobotId)).
		Find()
	if err != nil {
		return err
	}
	if len(info) == 0 || cast.ToInt(info[`application_type`]) != define.ApplicationTypeFlow {
		return errors.New("选择的工作流不合法")
	}

	maps := map[string]struct{}{}
	for _, item := range param.Params {
		if !common.IsVariableName(item.Key) {
			return errors.New(fmt.Sprintf(`自定义全局变量名格式错误:%s`, item.Key))
		}
		if tool.InArrayString(fmt.Sprintf(`global.%s`, item.Key), SysGlobalVariables()) {
			return errors.New(fmt.Sprintf(`自定义全局变量与系统变量同名:%s`, item.Key))
		}
		if !tool.InArrayString(item.Typ, []string{common.TypString, common.TypNumber, common.TypArrString, common.TypArrObject}) {
			return errors.New(fmt.Sprintf(`自定义全局变量类型不支持:%s`, item.Key))
		}
		if item.Required && len(item.Variable) == 0 {
			return errors.New(fmt.Sprintf(`缺少必填参数:%s`, item.Key))
		}
		if _, ok := maps[item.Key]; ok {
			return errors.New(fmt.Sprintf(`自定义全局变量名重复定义:%s`, item.Key))
		}
		maps[item.Key] = struct{}{}
	}

	return nil
}

func (params *ImmediatelyReplyNodeParams) Verify() error {
	if len(params.Content) == 0 {
		return errors.New(`消息内容不能为空`)
	}
	return nil
}


func GetVariablePlaceholders2(content string) []string {
	variables := make([]string, 0)
	for _, item := range regexp.MustCompile(`【?(([a-f0-9]{32}\.)?[a-zA-Z_][a-zA-Z0-9_\-.]*)】?`).FindAllStringSubmatch(content, -1) {
		if len(item) > 1 && !tool.InArrayString(item[1], variables) {
			if !strings.Contains(item[1], `.`) {
				continue
			}
			variables = append(variables, item[1])
		}
	}
	return variables
}

func ExtractVariables(data any, result *[]string) {
	if result == nil {
		return
	}
	if data == nil {
		return
	}
	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.String:
		str := val.String()
		matches := GetVariablePlaceholders2(str)
		*result = append(*result, matches...)
	case reflect.Map:
		for _, key := range val.MapKeys() {
			value := val.MapIndex(key)
			ExtractVariables(value.Interface(), result)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i).Interface()
			ExtractVariables(elem, result)
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if !field.CanInterface() {
				continue
			}
			ExtractVariables(field.Interface(), result)
		}
	default:
		return
	}
}
