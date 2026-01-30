// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package work_flow

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
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
	SysGlobal   []StartNodeParam `json:"sys_global"` // deprecated field
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
	Type     uint   `json:"type"` //1:equals,2:not equals,3:contains,4:not contains,5:is empty,6:not empty
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
	ToolInfo HttpToolInfo         `json:"http_tool_info"` // HTTP tool info stored in workflow to facilitate importing triggers that add HTTP tool group
}

type HttpToolInfo struct {
	HttpToolName            string `json:"http_tool_name"`
	HttpToolNameEn          string `json:"http_tool_name_en"`
	HttpToolKey             string `json:"http_tool_key"`
	HttpToolAvatar          string `json:"http_tool_avatar"`
	HttpToolDescription     string `json:"http_tool_description"`
	HttpToolNodeKey         string `json:"http_tool_node_key"`
	HttpToolNodeName        string `json:"http_tool_node_name"`
	HttpToolNodeNameEn      string `json:"http_tool_node_name_en"`
	HttpToolNodeDescription string `json:"http_tool_node_description"`
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
	RecallNeighborSwitch    bool            `json:"recall_neighbor_switch"`
	RecallNeighborAfterNum  common.MixedInt `json:"recall_neighbor_after_num"`
	RecallNeighborBeforeNum common.MixedInt `json:"recall_neighbor_before_num"`
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
	Role           int             `json:"role"`
}

func (params *LlmBaseParams) Verify(adminUserId int, lang string) error {
	if params.ModelConfigId <= 0 || len(params.UseModel) == 0 {
		return errors.New(i18n.Show(lang, `llm_model_please_select`))
	}
	//check model_config_id and use_model
	if ok := common.CheckModelIsValid(adminUserId, params.ModelConfigId.Int(), params.UseModel, common.Llm); !ok {
		return errors.New(i18n.Show(lang, `llm_model_selection_error`))
	}
	if params.ContextPair < 0 || params.ContextPair > 50 {
		return errors.New(i18n.Show(lang, `llm_context_pair_range_error`))
	}
	if params.Temperature < 0 || params.Temperature > 2 {
		return errors.New(i18n.Show(lang, `llm_temperature_range_error`))
	}
	if params.MaxToken < 0 {
		return errors.New(i18n.Show(lang, `llm_max_token_error`))
	}
	//if len(params.Prompt) == 0 {
	//	return errors.New(i18n.Show(lang, `llm_prompt_cannot_be_empty`))
	//}

	if !tool.InArrayInt(params.Role, []int{0, 1, 2}) {
		return errors.New(i18n.Show(lang, `llm_role_param_error`))
	}

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

const StaffAll = 1 //transfer type: 1 automatic allocation, 2 designated customer service, 3 designated customer service group
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

func (params *QuestionOptimizeNodeParams) Verify(adminUserId int, lang string) error {
	if len(params.QuestionValue) == 0 {
		return errors.New(i18n.Show(lang, `question_value_cannot_be_empty`))
	}
	return params.LlmBaseParams.Verify(adminUserId, lang)
}

/************************************/

type ParamsExtractorNodeParams struct {
	LlmBaseParams
	QuestionValue string               `json:"question_value"`
	Output        common.RecurveFields `json:"output"`
}

func (params *ParamsExtractorNodeParams) Verify(adminUserId int, lang string) error {
	if len(params.QuestionValue) == 0 {
		return errors.New(i18n.Show(lang, `question_value_cannot_be_empty`))
	}
	if err := params.LlmBaseParams.Verify(adminUserId, lang); err != nil {
		return err
	}
	//output field validation
	return params.Output.Verify(lang)
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
	Name         string               `json:"name"`
	Type         string               `json:"type"`
	Params       map[string]any       `json:"params"`
	Output       common.RecurveFields `json:"output_obj"`
	NoAuthFilter bool                 `json:"no_auth_filter,omitempty"`
}

type ImageGenerationParams struct {
	UseModel            string               `json:"use_model"`
	ModelConfigId       string               `json:"model_config_id"`
	Size                string               `json:"size"`
	ImageNum            string               `json:"image_num"`
	Prompt              string               `json:"prompt"`
	InputImages         []string             `json:"input_images"`
	ImageWatermark      string               `json:"image_watermark"`       //whether to add watermark, 1 add, 0 not add
	ImageOptimizePrompt string               `json:"image_optimize_prompt"` //whether to enable optimized prompt words, 1 enable, 0 disable
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
	VoiceType     string `json:"voice_type"` //desired voice type, options: system, voice_cloning, voice_generation, all
	Arguments     struct {
		Text          string `json:"text"`
		ModelId       int    `json:"model_id"`
		UseModel      string `json:"use_model"`
		ModelConfigId int    `json:"model_config_id"`
		VoiceSetting  struct {
			VoiceId           string  `json:"voice_id"`           // voice ID
			VoiceName         string  `json:"voice_name"`         // redundant
			Speed             float32 `json:"speed"`              // speed
			Vol               int     `json:"vol"`                // volume
			Pitch             int     `json:"pitch"`              // pitch
			Emotion           string  `json:"emotion"`            // emotion
			TextNormalization bool    `json:"text_normalization"` // text normalization
		} `json:"voice_setting"` // voice setting
		AudioSetting struct {
			SampleRate int    `json:"sample_rate"` // sample rate
			Bitrate    int    `json:"bitrate"`     // bitrate
			Format     string `json:"format"`      // format
			Channel    int    `json:"channel"`     // channel
			ForceCbr   bool   `json:"force_cbr"`   // constant bitrate
		} `json:"audio_setting"` // audio setting
		VoiceModify struct {
			Pitch        int    `json:"pitch"`         // pitch adjustment
			Intensity    int    `json:"intensity"`     // intensity adjustment (powerful/soft)
			Timbre       int    `json:"timbre"`        // timbre adjustment
			SoundEffects string `json:"sound_effects"` // sound effect setting
		} `json:"voice_modify"` // pitch adjustment
		LanguageBoost  string `json:"language_boost"` // minor language recognition capability
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
		FileUrl       string `json:"file_url"` //url of the audio to replicate, should be uploaded and converted to file_id to send to minimax
		VoiceId       string `json:"voice_id"` //voice_id of cloned voice
		ClonePrompt   struct {
			PromptAudioUrl string `json:"prompt_audio_url"` //url of the example audio, should be uploaded and converted to file_id to send to minimax
			PromptText     string `json:"prompt_text"`      //corresponding text of the example audio, should ensure consistency with audio content, punctuation required at the end of sentences
		} `json:"clone_prompt"` //voice cloning example audio, providing this parameter will help enhance voice similarity and stability for speech synthesis
		Text                    string `json:"text"`                      //replication listening parameter, limited to 1000 characters
		LanguageBoost           string `json:"language_boost"`            //whether to enhance recognition capability for specified minor languages and dialects
		Model                   string `json:"model"`                     //replication listening parameter
		NeedNoiseReduction      bool   `json:"need_noise_reduction"`      //enable noise reduction
		NeedVolumeNormalization bool   `json:"need_volume_normalization"` //enable volume normalization
		AigcWatermark           any    `json:"aigc_watermark"`            //whether to add audio rhythm marker at the end of synthesized listening audio, default is false
		TagMap                  any    `json:"tag_map"`
	} `json:"arguments"`
	Output common.RecurveFields `json:"output"`
}

type LibraryImportParams struct {
	ImportType                string               `json:"import_type"`                  //content import content, url import url
	LibraryId                 string               `json:"library_id"`                   //knowledge base id cannot be empty
	LibraryGroupId            string               `json:"library_group_id"`             //knowledge base group id 0 means ungrouped
	NormalUrl                 string               `json:"normal_url"`                   //normal knowledge base: document url
	NormalTitle               string               `json:"normal_title"`                 //normal knowledge base: document title
	NormalContent             string               `json:"normal_content"`               //normal knowledge base: document content
	NormalUrlRepeatOp         string               `json:"normal_url_repeat_op"`         //normal knowledge base: operation when url repeats import still import, not import do not import, update update content
	QaQuestion                string               `json:"qa_question"`                  //Q&A knowledge base: segmented question
	QaAnswer                  string               `json:"qa_answer"`                    //Q&A knowledge base: segmented answer
	QaImagesVariable          string               `json:"qa_images_variable"`           //Q&A knowledge base: answer images array<string>
	QaSimilarQuestionVariable string               `json:"qa_similar_question_variable"` //Q&A knowledge base: similar questions array<string>
	QaRepeatOp                string               `json:"qa_repeat_op"`                 //Q&A knowledge base: operation when question repeats import still import, not import do not import, update update content
	Outputs                   common.RecurveFields `json:"outputs"`                      //output fixed msg string
}

type WorkflowNodeParams struct {
	RobotId   int `json:"robot_id"`
	RobotInfo any `json:"robot_info"` // temporary data used by frontend
	Params    []struct {
		StartNodeParam
		Variable string `json:"variable"` //corresponding variables
		Tags     any    `json:"tags"`     // temporary data used by frontend
	} `json:"params"`
	Output    common.RecurveFields `json:"output"`
	Exception string               `json:"exception"`
}

type Message struct {
	Type    string `json:"type"`    //text image voice
	Content string `json:"content"` //content
}

type FinishNodeParams struct {
	OutType  string               `json:"out_type"` //variable returns message and variable, message returns message
	Messages []Message            `json:"messages"` //specific messages
	Outputs  common.RecurveFields `json:"outputs"`  //returned variables
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
	NodeTypeHttpTool,
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
	NodeTypeQuestion,
}

var BatchAllowNodeTypes = []int{
	NodeTypeRemark,
	NodeTypeTerm,
	NodeTypeCate,
	NodeTypeCurl,
	NodeTypeHttpTool,
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
	NodeTypeQuestion,
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
	Question         QuestionParams             `json:"question"`
}

func FillDiyGlobalBlanks(output TriggerOutputParam, start *StartNodeParams) {
	if len(output.Variable) == 0 {
		return //variable mapping not configured
	}
	key, found := strings.CutPrefix(output.Variable, `global.`)
	if !found {
		return //incorrect mapping, ignore
	}
	for _, param := range start.DiyGlobal {
		if param.Key == key {
			return //skip existing variables
		}
	}
	start.DiyGlobal = append(start.DiyGlobal, output.StartNodeParam)
}

func DisposeNodeParams(nodeType int, nodeParams, lang string) NodeParams {
	params := NodeParams{}
	_ = tool.JsonDecodeUseNumber(nodeParams, &params)
	params.Start.SysGlobal = make([]StartNodeParam, 0) //deprecated field
	if params.Start.DiyGlobal == nil {
		params.Start.DiyGlobal = make([]StartNodeParam, 0)
	}
	if params.Start.TriggerList == nil {
		params.Start.TriggerList = make([]TriggerConfig, 0)
	}
	if nodeType == NodeTypeStart && len(params.Start.TriggerList) == 0 { //default value handling
		chatTrigger := GetTriggerChatConfig(lang)
		params.Start.TriggerList = []TriggerConfig{chatTrigger}
		for _, output := range chatTrigger.Outputs {
			params.Start.DiyGlobal = append(params.Start.DiyGlobal, output.StartNodeParam)
		}
	}
	if nodeType == NodeTypeStart { //start node trigger output variable old data compatibility
		for i, trigger := range params.Start.TriggerList {
			outputs, exist := GetTriggerOutputsByType(trigger.TriggerType, lang)
			if !exist {
				continue
			}
			if trigger.TriggerType == TriggerTypeOfficial {
				switch trigger.TriggerOfficialConfig.MsgType {
				case define.TriggerOfficialMessage:
					outputs = GetMessage(lang)
				case define.TriggerOfficialQrCodeScan:
					outputs = GetQrcodeScan(lang)
				case define.TriggerOfficialSubscribeUnScribe:
					outputs = GetSubscribeUnsubscribe(lang)
				case define.TriggerOfficialMenuClick:
					outputs = GetMenuClick(lang)
				}
			}
			//collect old variable mapping data
			variableMap := make(map[string]string)
			for _, output := range trigger.Outputs {
				variableMap[output.Key] = output.Variable
			}
			//replace variable mapping to new configuration
			for idx, output := range outputs {
				if variable, ok := variableMap[output.Key]; ok {
					outputs[idx].Variable = variable
				}
			}
			params.Start.TriggerList[i].Outputs = outputs
			//supplement start node custom global variables
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
	case NodeTypeCurl, NodeTypeHttpTool:
		for variable := range common.SimplifyFields(node.NodeParams.Curl.Output) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //previous node, compatible with old data
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
			if len(last) > 0 && last[0] { //previous node, compatible with old data
				variables = append(variables, variable)
			}
		}
	case NodeTypeFormSelect:
		for _, variable := range []string{`output_list`, `row_num`} {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //previous node, compatible with old data
				variables = append(variables, variable)
			}
		}
	case NodeTypeCodeRun:
		for variable := range common.SimplifyFields(node.NodeParams.CodeRun.Output) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //previous node, compatible with old data
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
			if len(last) > 0 && last[0] { //previous node, compatible with old data
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
			if len(last) > 0 && last[0] { //previous node, compatible with old data
				variables = append(variables, variable)
			}
		}
	case NodeTypePlugin:
		for variable := range common.SimplifyFields(node.NodeParams.Plugin.Output) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //previous node, compatible with old data
				variables = append(variables, variable)
			}
		}
	case NodeTypeImageGeneration:
		for variable := range common.SimplifyFields(node.NodeParams.ImageGeneration.Output) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //previous node, compatible with old data
				variables = append(variables, variable)
			}
		}
	case NodeTypeJsonEncode:
		for variable := range common.SimplifyFields(node.NodeParams.JsonEncode.Outputs) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //previous node, compatible with old data
				variables = append(variables, variable)
			}
		}
	case NodeTypeJsonDecode:
		simpleFields := make(common.SimpleFields)
		node.NodeParams.JsonDecode.Outputs.SimplifyFieldsDeep(&simpleFields, ``)
		for variable := range simpleFields {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //previous node, compatible with old data
				variables = append(variables, variable)
			}
		}
	case NodeTypeTextToAudio:
		for variable := range common.SimplifyFields(node.NodeParams.TextToAudio.Output) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //previous node, compatible with old data
				variables = append(variables, variable)
			}
		}
	case NodeTypeVoiceClone:
		for variable := range common.SimplifyFields(node.NodeParams.VoiceClone.Output) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //previous node, compatible with old data
				variables = append(variables, variable)
			}
		}
	case NodeTypeLibraryImport:
		for variable := range common.SimplifyFields(node.NodeParams.LibraryImport.Outputs) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //previous node, compatible with old data
				variables = append(variables, variable)
			}
		}
	case NodeTypeWorkflow:
		for variable := range common.SimplifyFields(node.NodeParams.Workflow.Output) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] { //previous node, compatible with old data
				variables = append(variables, variable)
			}
		}
	case NodeTypeQuestion:
		for variable := range common.SimplifyFields(node.NodeParams.Question.Outputs) {
			variables = append(variables, fmt.Sprintf(`%s.%s`, node.NodeKey, variable))
			if len(last) > 0 && last[0] {
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
	//system global variables
	variables := SysGlobalVariables()
	//upper level node variables
	for _, node := range (*fn)[nodeKey] {
		variables = append(variables, node.GetVariables(true)...)
	}
	//recursively get upper level variables
	nodes := make([]*WorkFlowNode, 0)
	fn.recurveSetFrom(nodeKey, &nodes)
	for _, node := range nodes {
		variables = append(variables, node.GetVariables()...)
	}
	//deduplicate
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

func VerifyWorkFlowNodes(nodeList []WorkFlowNode, adminUserId int, lang string) (startNodeKey, modelConfigIds, libraryIds string, questionMultipleSwitch bool, err error) {
	startNodeCount, finishNodeCount := 0, 0
	fromNodes := make(FromNodes)
	for i, node := range nodeList {
		if err = node.Verify(adminUserId, lang); err != nil {
			return
		}
		if node.NodeType <= NodeTypeEdges {
			continue
		}
		if node.NodeType == NodeTypeStart {
			startNodeKey = node.NodeKey
			for _, trigger := range node.NodeParams.Start.TriggerList {
				if trigger.TriggerType == TriggerTypeChat && trigger.TriggerSwitch && trigger.TriggerChatConfig.QuestionMultipleSwitch {
					questionMultipleSwitch = true //Conversation Trigger - Multimodal Input - Enabled
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
		if node.NodeType == NodeTypeQuestion && node.NodeParams.Question.AnswerType == define.QuestionAnswerTypeMenu {
			if len(node.NodeParams.Question.ReplyContentList) > 0 {
				for _, reply := range node.NodeParams.Question.ReplyContentList {
					if len(reply.SmartMenu.MenuContent) > 0 {
						for _, menu := range reply.SmartMenu.MenuContent {
							fromNodes.AddRelation(&nodeList[i], menu.NextNodeKey)
							boolExist := false
							for _, nodeTemp := range nodeList {
								if nodeTemp.NodeKey == menu.NextNodeKey {
									boolExist = true
									break
								}
							}
							if !boolExist {
								err = errors.New(i18n.Show(lang, `node_next_node_not_exist`, node.NodeName))
								return
							}
						}
					}
				}
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
		err = errors.New(i18n.Show(lang, `workflow_only_one_start_node`))
		return
	}
	if finishNodeCount == 0 {
		err = errors.New(i18n.Show(lang, `workflow_need_finish_node`))
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
			err = errors.New(i18n.Show(lang, `workflow_has_isolated_node`, node.NodeName))
			return
		}
		//Verify that the selected variable must exist
		err = verifyNode(adminUserId, node, fromNodes, nodeList, lang)
		if err != nil {
			return
		}
	}
	var libraryArr []string
	//Collect the set of used model IDs
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

func verifyNode(adminUserId int, node WorkFlowNode, fromNodes FromNodes, nodeList []WorkFlowNode, lang string) (err error) {
	variables := fromNodes.GetVariableList(node.NodeKey)
	switch node.NodeType {
	case NodeTypeTerm:
		for _, param := range node.NodeParams.Term {
			for _, term := range param.Terms {
				if !tool.InArrayString(term.Variable, variables) {
					err = errors.New(i18n.Show(lang, `workflow_node_variable_not_exist`, node.NodeName, term.Variable))
					return
				}
			}
		}
	case NodeTypeCate:
		if len(node.NodeParams.Cate.QuestionValue) > 0 && !tool.InArrayString(node.NodeParams.Cate.QuestionValue, variables) {
			err = errors.New(i18n.Show(lang, `workflow_node_question_variable_not_exist`, node.NodeName, node.NodeParams.Cate.QuestionValue))
			return
		}
	case NodeTypeCurl, NodeTypeHttpTool:
		for _, param := range node.NodeParams.Curl.Headers {
			if variable, ok := CheckVariablePlaceholder(param.Value, variables); !ok {
				err = errors.New(i18n.Show(lang, `workflow_node_headers_variable_not_exist`, node.NodeName, variable))
				return
			}
		}
		for _, param := range node.NodeParams.Curl.Params {
			if variable, ok := CheckVariablePlaceholder(param.Value, variables); !ok {
				err = errors.New(i18n.Show(lang, `workflow_node_params_variable_not_exist`, node.NodeName, variable))
				return
			}
		}
		switch node.NodeParams.Curl.Type {
		case TypeUrlencoded:
			for _, param := range node.NodeParams.Curl.Body {
				if variable, ok := CheckVariablePlaceholder(param.Value, variables); !ok {
					err = errors.New(i18n.Show(lang, `workflow_node_body_variable_not_exist`, node.NodeName, variable))
					return
				}
			}
		case TypeJsonBody:
			if variable, ok := CheckVariablePlaceholder(node.NodeParams.Curl.BodyRaw, variables); !ok {
				err = errors.New(i18n.Show(lang, `workflow_node_json_body_variable_not_exist`, node.NodeName, variable))
				return
			}
		}
	case NodeTypeLibs:
		if len(node.NodeParams.Libs.QuestionValue) > 0 && !tool.InArrayString(node.NodeParams.Libs.QuestionValue, variables) {
			err = errors.New(i18n.Show(lang, `workflow_node_question_variable_not_exist`, node.NodeName, node.NodeParams.Libs.QuestionValue))
			return
		}
	case NodeTypeLlm:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.Llm.Prompt, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_prompt_variable_not_exist`, node.NodeName, variable))
			return
		}
		if len(node.NodeParams.Llm.LibsNodeKey) > 0 {
			variable := fmt.Sprintf(`%s.%s`, node.NodeParams.Llm.LibsNodeKey, `special.lib_paragraph_list`)
			if !tool.InArrayString(variable, variables) {
				err = errors.New(i18n.Show(lang, `workflow_node_knowledge_base_ref_not_from_parent_node`, node.NodeName))
				return
			}
		}
		if len(node.NodeParams.Llm.QuestionValue) > 0 && !tool.InArrayString(node.NodeParams.Llm.QuestionValue, variables) {
			err = errors.New(i18n.Show(lang, `workflow_node_question_variable_not_exist`, node.NodeName, node.NodeParams.Llm.QuestionValue))
			return
		}
	case NodeTypeAssign:
		for i, param := range node.NodeParams.Assign {
			if !tool.InArrayString(param.Variable, variables) { //Custom variable does not exist
				err = errors.New(i18n.Show(lang, `workflow_node_custom_global_variable_not_exist`, node.NodeName, param.Variable))
				return
			}
			if variable, ok := CheckVariablePlaceholder(param.Value, variables); !ok {
				err = errors.New(i18n.Show(lang, `workflow_node_line_variable_not_exist`, node.NodeName, i+1, variable))
				return
			}
		}
	case NodeTypeReply:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.Reply.Content, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_insert_variable_not_exist`, node.NodeName, variable))
			return
		}
	case NodeTypeQuestionOptimize:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.QuestionOptimize.Prompt, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_conversation_context_variable_not_exist`, node.NodeName, variable))
			return
		}
		if !tool.InArrayString(node.NodeParams.QuestionOptimize.QuestionValue, variables) {
			err = errors.New(i18n.Show(lang, `workflow_node_question_variable_not_exist`, node.NodeName, node.NodeParams.QuestionOptimize.QuestionValue))
			return
		}
	case NodeTypeParamsExtractor:
		if !tool.InArrayString(node.NodeParams.ParamsExtractor.QuestionValue, variables) {
			err = errors.New(i18n.Show(lang, `workflow_node_question_variable_not_exist`, node.NodeName, node.NodeParams.ParamsExtractor.QuestionValue))
			return
		}
	case NodeTypeFormInsert:
		for _, field := range node.NodeParams.FormInsert.Datas {
			if variable, ok := CheckVariablePlaceholder(field.Value, variables); !ok {
				err = errors.New(i18n.Show(lang, `workflow_node_insert_variable_not_exist`, node.NodeName, variable))
				return
			}
		}
	case NodeTypeFormDelete:
		for _, field := range node.NodeParams.FormDelete.Where {
			if variable, ok := CheckVariablePlaceholder(field.RuleValue1, variables); !ok {
				err = errors.New(i18n.Show(lang, `workflow_node_insert_variable_not_exist`, node.NodeName, variable))
				return
			}
			if variable, ok := CheckVariablePlaceholder(field.RuleValue2, variables); !ok {
				err = errors.New(i18n.Show(lang, `workflow_node_insert_variable_not_exist`, node.NodeName, variable))
				return
			}
		}
	case NodeTypeFormUpdate:
		for _, field := range node.NodeParams.FormUpdate.Where {
			if variable, ok := CheckVariablePlaceholder(field.RuleValue1, variables); !ok {
				err = errors.New(i18n.Show(lang, `workflow_node_insert_variable_not_exist`, node.NodeName, variable))
				return
			}
			if variable, ok := CheckVariablePlaceholder(field.RuleValue2, variables); !ok {
				err = errors.New(i18n.Show(lang, `workflow_node_insert_variable_not_exist`, node.NodeName, variable))
				return
			}
		}
		for _, field := range node.NodeParams.FormUpdate.Datas {
			if variable, ok := CheckVariablePlaceholder(field.Value, variables); !ok {
				err = errors.New(i18n.Show(lang, `workflow_node_insert_variable_not_exist`, node.NodeName, variable))
				return
			}
		}
	case NodeTypeFormSelect:
		for _, field := range node.NodeParams.FormSelect.Where {
			if variable, ok := CheckVariablePlaceholder(field.RuleValue1, variables); !ok {
				err = errors.New(i18n.Show(lang, `workflow_node_insert_variable_not_exist`, node.NodeName, variable))
				return
			}
			if variable, ok := CheckVariablePlaceholder(field.RuleValue2, variables); !ok {
				err = errors.New(i18n.Show(lang, `workflow_node_insert_variable_not_exist`, node.NodeName, variable))
				return
			}
		}
	case NodeTypeCodeRun:
		for _, param := range node.NodeParams.CodeRun.Params {
			if !tool.InArrayString(param.Variable, variables) {
				err = errors.New(i18n.Show(lang, `workflow_node_variable_not_exist`, node.NodeName, param.Variable))
				return
			}
		}
	case NodeTypeMcp:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.Mcp.ToolName, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_mcp_tool_variable_not_exist`, node.NodeName, variable))
			return
		}
	case NodeTypePlugin:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.Plugin.Name, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_variable_not_exist`, node.NodeName, variable))
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.Plugin.Type, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_variable_not_exist`, node.NodeName, variable))
			return
		}
	case NodeTypeTextToAudio:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.TextToAudio.Arguments.Text, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.TextToAudio.Arguments.VoiceSetting.VoiceId, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
			return
		}
	case NodeTypeVoiceClone:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.VoiceClone.Arguments.FileUrl, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.VoiceClone.Arguments.VoiceId, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.VoiceClone.Arguments.ClonePrompt.PromptAudioUrl, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.VoiceClone.Arguments.ClonePrompt.PromptText, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.VoiceClone.Arguments.Text, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
			return
		}
	case NodeTypeWorkflow:
		for _, param := range node.NodeParams.Workflow.Params {
			if variable, ok := CheckVariablePlaceholder(param.Variable, variables); !ok {
				err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
				return
			}
		}
	case NodeTypeLoop:
		//verify loop arrays
		if node.NodeParams.Loop.LoopType == common.LoopTypeArray {
			if len(node.NodeParams.Loop.LoopArrays) == 0 {
				err = errors.New(i18n.Show(lang, `workflow_node_not_add_loop_array`, node.NodeName))
				return
			}
			//Verify reference to previous node
			bFind := VerityLoopParams(node.NodeParams.Loop.LoopArrays, nodeList)
			if !bFind {
				err = errors.New(i18n.Show(lang, `workflow_node_loop_array_variable_not_exist`, node.NodeName))
				return
			}
		} else if node.NodeParams.Loop.LoopNumber.Int() <= 0 {
			err = errors.New(i18n.Show(lang, `workflow_node_loop_number_bigger_than_zero`, node.NodeName))
			return
		}
		//Verify if intermediate variable initial value exists
		bFind := VerityLoopParams(node.NodeParams.Loop.IntermediateParams, nodeList)
		if !bFind {
			err = errors.New(i18n.Show(lang, `workflow_node_loop_intermediate_not_exist`, node.NodeName))
			return
		}
		//child node
		childNodes := make([]WorkFlowNode, 0)
		for _, vNode := range nodeList {
			if vNode.LoopParentKey == node.NodeKey {
				childNodes = append(childNodes, vNode)
			}
		}
		//Verify if output comes from intermediate variable
		bFind = VerityLoopParams(node.NodeParams.Loop.Output, []WorkFlowNode{node})
		if !bFind {
			err = errors.New(i18n.Show(lang, `workflow_node_loop_output_not_exist`, node.NodeName))
			return
		}
		//Verify child nodes
		_, _, _, err = VerityLoopWorkflowNodes(adminUserId, node, childNodes, LoopAllowNodeTypes, i18n.Show(lang, `loop_node`), lang)
		if err != nil {
			return
		}
	case NodeTypeBatch:
		if len(node.NodeParams.Batch.BatchArrays) == 0 {
			err = errors.New(i18n.Show(lang, `workflow_node_batch_array_not_added`, node.NodeName))
			return
		}
		//Verify if execution array initial value exists
		bFind := VerityLoopParams(node.NodeParams.Batch.BatchArrays, nodeList)
		if !bFind {
			err = errors.New(i18n.Show(lang, `workflow_node_batch_array_variable_not_exist`, node.NodeName))
			return
		}
		childNodes := make([]WorkFlowNode, 0)
		for _, vNode := range nodeList {
			if vNode.LoopParentKey == node.NodeKey {
				childNodes = append(childNodes, vNode)
			}
		}
		//Verify if output value exists
		bFind = VerityLoopParams(node.NodeParams.Batch.Output, childNodes)
		if !bFind {
			err = errors.New(i18n.Show(lang, `workflow_node_batch_output_variable_not_exist`, node.NodeName))
			return
		}
		//Verify child nodes
		_, _, _, err = VerityLoopWorkflowNodes(adminUserId, node, childNodes, BatchAllowNodeTypes, i18n.Show(lang, `batch_processing`), lang)
		if err != nil {
			return
		}
	case NodeTypeFinish:
		if node.NodeParams.Finish.OutType == define.FinishNodeOutTypeMessage {
			for _, field := range node.NodeParams.Finish.Messages {
				if variable, ok := CheckVariablePlaceholder(field.Content, variables); !ok {
					err = errors.New(i18n.Show(lang, `workflow_node_finish_message_variable_not_exist`, node.NodeName, variable))
					return
				}
			}
		}

	case NodeTypeImageGeneration:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.ImageGeneration.Prompt, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_image_generation_prompt_variable_not_exist`, node.NodeName, variable))
			return
		}
		for _, imageUrl := range node.NodeParams.ImageGeneration.InputImages {
			if variable, ok := CheckVariablePlaceholder(imageUrl, variables); !ok {
				err = errors.New(i18n.Show(lang, `workflow_node_image_generation_input_image_variable_not_exist`, node.NodeName, variable))
				return
			}
		}
	case NodeTypeJsonEncode:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.JsonEncode.InputVariable, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
			return
		}
	case NodeTypeJsonDecode:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.JsonDecode.InputVariable, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
			return
		}
	case NodeTypeLibraryImport:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.LibraryImport.NormalTitle, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.LibraryImport.NormalUrl, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.LibraryImport.NormalContent, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.LibraryImport.QaQuestion, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.LibraryImport.QaAnswer, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
			return
		}
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.LibraryImport.QaSimilarQuestionVariable, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_content_variable_not_exist`, node.NodeName, variable))
		}
	case NodeTypeImmediatelyReply:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.ImmediatelyReply.Content, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_insert_variable_not_exist`, node.NodeName, variable))
			return
		}
	case NodeTypeQuestion:
		if variable, ok := CheckVariablePlaceholder(node.NodeParams.Question.AnswerText, variables); !ok {
			err = errors.New(i18n.Show(lang, `workflow_node_variable_not_exist`, node.NodeName, variable))
			return
		}
		if len(node.NodeParams.Question.ReplyContentList) > 0 {
			for _, menu := range node.NodeParams.Question.ReplyContentList {
				if len(menu.SmartMenu.MenuContent) > 0 {
					for _, menuContent := range menu.SmartMenu.MenuContent {
						if variable, ok := CheckVariablePlaceholder(menuContent.Content, variables); !ok {
							err = errors.New(i18n.Show(lang, `workflow_node_variable_not_exist`, node.NodeName, variable))
							return
						}
					}
				}
			}
		}
	}
	return nil
}

func (node *WorkFlowNode) Verify(adminUserId int, lang string) error {
	if !tool.InArrayInt(node.NodeType.Int(), NodeTypes[:]) {
		return errors.New(i18n.Show(lang, `node_type_param_error`, node.NodeType))
	}
	if len(node.NodeName) == 0 {
		return errors.New(i18n.Show(lang, `node_name_cannot_be_empty`, node.NodeKey))
	}
	if len(node.NodeKey) == 0 || !common.IsMd5Str(node.NodeKey) {
		return errors.New(i18n.Show(lang, `node_key_param_format_error`, node.NodeName))
	}
	if len(node.NextNodeKey) > 0 && !common.IsMd5Str(node.NextNodeKey) {
		return errors.New(i18n.Show(lang, `node_next_node_key_param_format_error`, node.NodeName))
	}
	if len(node.NextNodeKey) == 0 && node.LoopParentKey == `` && !tool.InArrayInt(node.NodeType.Int(), []int{NodeTypeRemark, NodeTypeEdges, NodeTypeFinish, NodeTypeManual, NodeTypeLoopEnd, NodeTypeQuestion}) {
		return errors.New(i18n.Show(lang, `node_not_specify_next_node`, node.NodeName))
	}
	var err error
	switch node.NodeType {
	case NodeTypeStart:
		err = node.NodeParams.Start.Verify(lang)
	case NodeTypeTerm:
		err = node.NodeParams.Term.Verify(lang)
	case NodeTypeCate:
		err = node.NodeParams.Cate.Verify(adminUserId, lang)
	case NodeTypeCurl, NodeTypeHttpTool:
		err = node.NodeParams.Curl.Verify(lang)
	case NodeTypeLibs:
		err = node.NodeParams.Libs.Verify(adminUserId, lang)
	case NodeTypeLlm:
		err = node.NodeParams.Llm.Verify(adminUserId, lang)
	case NodeTypeAssign:
		err = node.NodeParams.Assign.Verify(node, lang)
	case NodeTypeReply:
		err = node.NodeParams.Reply.Verify(lang)
	case NodeTypeManual:
		err = node.NodeParams.Manual.Verify(adminUserId, lang)
	case NodeTypeQuestionOptimize:
		err = node.NodeParams.QuestionOptimize.Verify(adminUserId, lang)
	case NodeTypeParamsExtractor:
		err = node.NodeParams.ParamsExtractor.Verify(adminUserId, lang)
	case NodeTypeFormInsert:
		err = node.NodeParams.FormInsert.Verify(adminUserId, lang)
	case NodeTypeFormDelete:
		err = node.NodeParams.FormDelete.Verify(adminUserId, lang)
	case NodeTypeFormUpdate:
		err = node.NodeParams.FormUpdate.Verify(adminUserId, lang)
	case NodeTypeFormSelect:
		err = node.NodeParams.FormSelect.Verify(adminUserId, lang)
	case NodeTypeCodeRun:
		err = node.NodeParams.CodeRun.Verify(lang)
	case NodeTypeMcp:
		err = node.NodeParams.Mcp.Verify(adminUserId, lang)
	case NodeTypeLoop:
		err = node.NodeParams.Loop.Verify(node.NodeName, lang)
	case NodeTypePlugin:
		err = node.NodeParams.Plugin.Verify(adminUserId, lang)
	case NodeTypeBatch:
		err = node.NodeParams.Batch.Verify(node.NodeName, lang)
	case NodeTypeFinish:
		err = node.NodeParams.Finish.Verify(node.NodeName, lang)
	case NodeTypeImageGeneration:
		err = node.NodeParams.ImageGeneration.Verify(node.NodeName, lang)
	case NodeTypeJsonEncode:
		err = node.NodeParams.JsonEncode.Verify(node.NodeName, lang)
	case NodeTypeJsonDecode:
		err = node.NodeParams.JsonDecode.Verify(node.NodeName, lang)
	case NodeTypeTextToAudio:
		err = node.NodeParams.TextToAudio.Verify(lang)
	case NodeTypeVoiceClone:
		err = node.NodeParams.VoiceClone.Verify(adminUserId, lang)
	case NodeTypeLibraryImport:
		err = node.NodeParams.LibraryImport.Verify(adminUserId, lang)
	case NodeTypeWorkflow:
		err = node.NodeParams.Workflow.Verify(adminUserId, lang)
	case NodeTypeImmediatelyReply:
		err = node.NodeParams.ImmediatelyReply.Verify(lang)
	case NodeTypeQuestion:
		err = node.NodeParams.Question.Verify(node.NodeName, lang)
	}

	if err != nil {
		return errors.New(i18n.Show(lang, `node_error_with_detail`, node.NodeName, err.Error()))
	}
	return nil
}

func (params *StartNodeParams) Verify(lang string) error {
	maps := map[string]struct{}{}
	for _, item := range params.DiyGlobal {
		if !common.IsVariableName(item.Key) {
			return errors.New(i18n.Show(lang, `custom_global_variable_name_format_error`, item.Key))
		}
		if tool.InArrayString(fmt.Sprintf(`global.%s`, item.Key), SysGlobalVariables()) {
			return errors.New(i18n.Show(lang, `custom_global_variable_conflict_with_system`, item.Key))
		}
		if !tool.InArrayString(item.Typ, []string{common.TypString, common.TypNumber, common.TypArrString, common.TypArrObject, common.TypBoole, common.TypArrNumber}) {
			return errors.New(i18n.Show(lang, `custom_global_variable_type_not_supported`, item.Key))
		}
		if _, ok := maps[item.Key]; ok {
			return errors.New(i18n.Show(lang, `custom_global_variable_duplicate_definition`, item.Key))
		}
		maps[item.Key] = struct{}{}
	}
	//Trigger parameter validation
	if len(params.TriggerList) == 0 {
		return errors.New(i18n.Show(lang, `workflow_at_least_one_trigger`))
	}
	for i, trigger := range params.TriggerList {
		for _, output := range trigger.Outputs {
			if len(output.Variable) == 0 {
				continue //Trigger output supports not configuring variable mapping
			}
			key, _ := strings.CutPrefix(output.Variable, `global.`)
			if _, ok := maps[key]; !ok {
				return errors.New(i18n.Show(lang, `trigger_output_variable_mapping_error`, i+1, trigger.TriggerName, output.Key, output.Desc))
			}
		}
	}
	return nil
}

func (params *TermNodeParams) Verify(lang string) error {
	if params == nil || len(*params) == 0 {
		return errors.New(i18n.Show(lang, `config_params_cannot_be_empty`))
	}
	for i, item := range *params {
		if len(item.Terms) == 0 {
			return errors.New(i18n.Show(lang, `branch_config_empty`, i+1))
		}
		for j, term := range item.Terms {
			if len(term.Variable) == 0 {
				return errors.New(i18n.Show(lang, `branch_condition_variable_select`, i+1, j+1))
			}
			if !common.IsVariableNames(term.Variable) {
				return errors.New(i18n.Show(lang, `branch_condition_variable_format_error`, i+1, j+1))
			}
			if term.IsMult { // Array types do not support equal and not equal
				if !tool.InArrayInt(int(term.Type), TermTypes[2:]) {
					return errors.New(i18n.Show(lang, `branch_condition_match_error`, i+1, j+1))
				}
			} else {
				if !tool.InArrayInt(int(term.Type), TermTypes[:]) {
					return errors.New(i18n.Show(lang, `branch_condition_match_error`, i+1, j+1))
				}
			}
			if !tool.InArrayInt(int(term.Type), []int{TermTypeEmpty, TermTypeNotEmpty}) && len(term.Value) == 0 {
				return errors.New(i18n.Show(lang, `branch_condition_input_match_value`, i+1, j+1))
			}
		}
		if len(item.NextNodeKey) == 0 || !common.IsMd5Str(item.NextNodeKey) {
			return errors.New(i18n.Show(lang, `branch_next_node_not_specified`, i+1))
		}
	}
	return nil
}

func (params *CateNodeParams) Verify(adminUserId int, lang string) error {
	if err := params.LlmBaseParams.Verify(adminUserId, lang); err != nil {
		return err
	}
	if len(params.Categorys) == 0 {
		return errors.New(i18n.Show(lang, `cate_category_list_empty`))
	}
	for i, category := range params.Categorys {
		if len(category.Category) == 0 {
			return errors.New(i18n.Show(lang, `cate_category_name_empty`, i+1))
		}
		if len(category.NextNodeKey) == 0 || !common.IsMd5Str(category.NextNodeKey) {
			return errors.New(i18n.Show(lang, `cate_next_node_param_error`, i+1))
		}
	}
	return nil
}

func (params *CurlNodeParams) Verify(lang string) error {
	if !tool.InArrayString(params.Method, []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}) {
		return errors.New(i18n.Show(lang, `request_method_param_error`))
	}
	if _, err := url.Parse(params.Rawurl); err != nil || len(params.Rawurl) == 0 {
		return errors.New(i18n.Show(lang, `request_link_empty_or_error`))
	}
	for _, header := range params.Headers {
		if len(header.Key) == 0 || len(header.Value) == 0 {
			return errors.New(i18n.Show(lang, `header_key_value_cannot_be_empty`))
		}
		if header.Key == `Content-Type` {
			return errors.New(i18n.Show(lang, `header_content_type_not_allowed`))
		}
	}
	for _, param := range params.Params {
		if len(param.Key) == 0 || len(param.Value) == 0 {
			return errors.New(i18n.Show(lang, `params_key_value_cannot_be_empty`))
		}
	}
	if params.Method != http.MethodGet {
		switch params.Type {
		case TypeNone:
		case TypeUrlencoded:
			for _, param := range params.Body {
				if len(param.Key) == 0 || len(param.Value) == 0 {
					return errors.New(i18n.Show(lang, `body_key_value_cannot_be_empty`))
				}
			}
		case TypeJsonBody:
			if len(params.BodyRaw) == 0 {
				return errors.New(i18n.Show(lang, `json_body_cannot_be_empty`))
			}
			var temp any
			if err := tool.JsonDecodeUseNumber(params.BodyRaw, &temp); err != nil {
				return errors.New(i18n.Show(lang, `body_not_json_string`))
			}
		default:
			return errors.New(i18n.Show(lang, `body_param_type_error`))
		}
	}
	if params.Timeout > 60 {
		return errors.New(i18n.Show(lang, `request_timeout_max_value`))
	}
	//Output field validation
	return params.Output.Verify(lang)
}

func (params *LibsNodeParams) Verify(adminUserId int, lang string) error {
	if len(params.LibraryIds) == 0 || !common.CheckIds(params.LibraryIds) {
		return errors.New(i18n.Show(lang, `related_library_empty_or_param_error`))
	}
	if len(params.QuestionValue) == 0 {
		return errors.New(i18n.Show(lang, `question_value_cannot_be_empty`))
	}
	for _, libraryId := range strings.Split(params.LibraryIds, `,`) {
		info, err := common.GetLibraryInfo(cast.ToInt(libraryId), adminUserId)
		if err != nil {
			logs.Error(err.Error())
			return err
		}
		if len(info) == 0 {
			return errors.New(i18n.Show(lang, `related_library_not_exist`, libraryId))
		}
	}
	if !tool.InArrayInt(params.SearchType.Int(), []int{define.SearchTypeMixed, define.SearchTypeVector, define.SearchTypeFullText}) {
		return errors.New(i18n.Show(lang, `library_search_mode_param_error`))
	}
	if err := common.CheckRrfWeight(params.RrfWeight, lang); err != nil {
		return err
	}
	if params.TopK <= 0 || params.TopK > 500 {
		return errors.New(i18n.Show(lang, `library_search_topk_range`))
	}
	if params.Similarity < 0 || params.Similarity > 1 {
		return errors.New(i18n.Show(lang, `library_search_similarity_range`))
	}
	if params.RerankStatus > 0 || params.RerankModelConfigId != 0 || len(params.RerankUseModel) > 0 {
		if params.RerankModelConfigId <= 0 || len(params.RerankUseModel) == 0 {
			return errors.New(i18n.Show(lang, `rerank_model_please_select`))
		}
		if ok := common.CheckModelIsValid(adminUserId, params.RerankModelConfigId.Int(), params.RerankUseModel, common.Rerank); !ok {
			return errors.New(i18n.Show(lang, `rerank_model_selection_error`))
		}
	}
	return nil
}

func (params *LlmNodeParams) Verify(adminUserId int, lang string) error {
	if err := params.LlmBaseParams.Verify(adminUserId, lang); err != nil {
		return err
	}
	if len(params.Prompt) == 0 {
		return errors.New(i18n.Show(lang, `llm_prompt_cannot_be_empty`))
	}
	if len(params.QuestionValue) == 0 {
		return errors.New(i18n.Show(lang, `question_value_cannot_be_empty`))
	}
	if len(params.LibsNodeKey) > 0 && !common.IsMd5Str(params.LibsNodeKey) {
		return errors.New(i18n.Show(lang, `knowledge_base_ref_node_param_format_error`))
	}
	return nil
}

func (params *AssignNodeParams) Verify(node *WorkFlowNode, lang string) error {
	if params == nil || len(*params) == 0 {
		return errors.New(i18n.Show(lang, `config_params_cannot_be_empty`))
	}
	for i, param := range *params {
		if len(param.Variable) == 0 {
			return errors.New(i18n.Show(lang, `line_select_variable`, i+1))
		}
		if node.LoopParentKey == `` && (!strings.HasPrefix(param.Variable, `global.`) || !common.IsVariableNames(param.Variable)) {
			return errors.New(i18n.Show(lang, `line_variable_format_error`, i+1))
		}
		if tool.InArrayString(param.Variable, SysGlobalVariables()) {
			return errors.New(i18n.Show(lang, `line_global_variable_assign_error`, i+1))
		}
	}
	return nil
}

func (params *ReplyNodeParams) Verify(lang string) error {
	if len(params.Content) == 0 {
		return errors.New(i18n.Show(lang, `message_content_cannot_be_empty`))
	}
	return nil
}

func (params *ManualNodeParams) Verify(adminUserId int, lang string) error {
	return errors.New(i18n.Show(lang, `cloud_version_manual_node_only`))
}

func checkFormId(adminUserId, formId int, lang string) error {
	if formId <= 0 {
		return errors.New(i18n.Show(lang, `workflow_node_db_table_not_selected`))
	}
	form, err := msql.Model(`form`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(formId)).Where(`delete_time`, `0`).Field(`id`).Find()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	if len(form) == 0 {
		return errors.New(i18n.Show(lang, `workflow_node_db_table_info_not_exist`))
	}
	return nil
}

func checkFormDatas(adminUserId, formId int, datas []FormFieldValue, lang string) error {
	if len(datas) == 0 {
		return errors.New(i18n.Show(lang, `field_list_cannot_be_empty`))
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
			return errors.New(i18n.Show(lang, `line_field_name_param_cannot_be_empty`, i+1))
		}
		field, ok := fields[data.Name]
		if !ok {
			return errors.New(i18n.Show(lang, `line_field_not_exist_in_table`, i+1, data.Name))
		}
		if len(data.Type) == 0 || data.Type != field[`type`] {
			return errors.New(i18n.Show(lang, `line_field_type_not_consistent`, i+1, data.Name, field[`description`]))
		}
		if len(data.Value) == 0 && cast.ToBool(field[`required`]) {
			return errors.New(i18n.Show(lang, `line_field_required_cannot_be_empty`, i+1, data.Name, field[`description`]))
		}
		if _, ok := maps[data.Name]; ok {
			return errors.New(i18n.Show(lang, `line_field_duplicate_in_list`, i+1, data.Name, field[`description`]))
		}
		maps[data.Name] = struct{}{}
	}
	return nil
}

func checkFormWhere(adminUserId, formId int, where []define.FormFilterCondition, lang string) error {
	if len(where) == 0 {
		return errors.New(i18n.Show(lang, `condition_list_cannot_be_empty`))
	}
	fields, err := msql.Model(`form_field`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`form_id`, cast.ToString(formId)).ColumnMap(`name,type,description`, `id`)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	fields[`0`] = msql.Params{`name`: `id`, `type`: `integer`, `description`: `ID`} //Append an ID for compatibility processing
	for i, condition := range where {
		if condition.FormFieldId < 0 { //Note: here it can be equal to 0
			return errors.New(i18n.Show(lang, `line_field_select_param_invalid`, i+1))
		}
		field, ok := fields[cast.ToString(condition.FormFieldId)]
		if !ok {
			return errors.New(i18n.Show(lang, `line_field_not_exist_in_table`, i+1))
		}
		if err = condition.Check(field[`type`], true); err != nil {
			return errors.New(i18n.Show(lang, `line_field_validation_error`, i+1, field[`name`], field[`description`], err.Error()))
		}
	}
	return nil
}

func checkFormFields(adminUserId, formId int, Fields []FormFieldTyp, lang string) error {
	if len(Fields) == 0 {
		return errors.New(i18n.Show(lang, `field_list_cannot_be_empty`))
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
			return errors.New(i18n.Show(lang, `line_field_name_param_cannot_be_empty`, i+1))
		}
		field, ok := fields[data.Name]
		if !ok {
			return errors.New(i18n.Show(lang, `line_field_type_not_consistent`, i+1, data.Name, field[`description`]))
		}
		if len(data.Type) == 0 || data.Type != field[`type`] {
			return errors.New(i18n.Show(lang, `line_field_type_not_consistent`, i+1, data.Name, field[`description`]))
		}
		if _, ok := maps[data.Name]; ok {
			return errors.New(i18n.Show(lang, `line_field_duplicate_in_list`, i+1, data.Name, field[`description`]))
		}
		maps[data.Name] = struct{}{}
	}
	return nil
}

func (params *FormInsertNodeParams) Verify(adminUserId int, lang string) error {
	if err := checkFormId(adminUserId, params.FormId.Int(), lang); err != nil {
		return err
	}
	if err := checkFormDatas(adminUserId, params.FormId.Int(), params.Datas, lang); err != nil {
		return err
	}
	return nil
}

func (params *FormDeleteNodeParams) Verify(adminUserId int, lang string) error {
	if err := checkFormId(adminUserId, params.FormId.Int(), lang); err != nil {
		return err
	}
	if !tool.InArrayInt(params.Typ.Int(), []int{1, 2}) {
		return errors.New(i18n.Show(lang, `condition_relationship_param_error`))
	}
	if err := checkFormWhere(adminUserId, params.FormId.Int(), params.Where, lang); err != nil {
		return err
	}
	return nil
}

func (params *FormUpdateNodeParams) Verify(adminUserId int, lang string) error {
	if err := checkFormId(adminUserId, params.FormId.Int(), lang); err != nil {
		return err
	}
	if !tool.InArrayInt(params.Typ.Int(), []int{1, 2}) {
		return errors.New(i18n.Show(lang, `condition_relationship_param_error`))
	}
	if err := checkFormWhere(adminUserId, params.FormId.Int(), params.Where, lang); err != nil {
		return err
	}
	if err := checkFormDatas(adminUserId, params.FormId.Int(), params.Datas, lang); err != nil {
		return err
	}
	return nil
}

func (params *FormSelectNodeParams) Verify(adminUserId int, lang string) error {
	if err := checkFormId(adminUserId, params.FormId.Int(), lang); err != nil {
		return err
	}
	if !tool.InArrayInt(params.Typ.Int(), []int{1, 2}) {
		return errors.New(i18n.Show(lang, `condition_relationship_param_error`))
	}
	if err := checkFormWhere(adminUserId, params.FormId.Int(), params.Where, lang); err != nil {
		return err
	}
	if err := checkFormFields(adminUserId, params.FormId.Int(), params.Fields, lang); err != nil {
		return err
	}
	for _, order := range params.Order {
		if !tool.InArrayString(order.Name, []string{`id`, `create_time`, `update_time`}) {
			return errors.New(i18n.Show(lang, `sort_operation_not_support`, order.Name))
		}
	}
	if params.Limit <= 0 || params.Limit > 1000 {
		return errors.New(i18n.Show(lang, `query_quantity_range`))
	}
	return nil
}

func (params *CodeRunNodeParams) Verify(lang string) error {
	maps := map[string]struct{}{}
	for idx, param := range params.Params {
		if !common.IsVariableName(param.Field) {
			return errors.New(i18n.Show(lang, `custom_input_param_key_format_error`, param.Field))
		}
		if len(param.Variable) == 0 {
			return errors.New(i18n.Show(lang, `line_custom_input_param_select_variable`, idx+1))
		}
		if !common.IsVariableNames(param.Variable) {
			return errors.New(i18n.Show(lang, `line_custom_input_param_variable_format_error`, idx+1))
		}
		if _, ok := maps[param.Field]; ok {
			return errors.New(i18n.Show(lang, `custom_input_param_key_duplicate_definition`, param.Field))
		}
		maps[param.Field] = struct{}{}
	}
	ok, err := regexp.MatchString(`function\s+main\s*\(.*\)\s*\{`, params.MainFunc)
	if err != nil || !ok {
		return errors.New(i18n.Show(lang, `javascript_code_missing_main_function`))
	}
	if params.Timeout < 1 || params.Timeout > 60 {
		return errors.New(i18n.Show(lang, `code_run_timeout_range`))
	}
	if err = params.Output.Verify(lang); err != nil {
		return err
	}
	if len(params.Exception) == 0 || !common.IsMd5Str(params.Exception) {
		return errors.New(i18n.Show(lang, `exception_handling_next_node_not_specified`))
	}
	return nil
}

func (params *McpNodeParams) Verify(adminUserId int, lang string) error {
	info, err := msql.Model(`mcp_provider`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(params.ProviderId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	if len(info) == 0 {
		return errors.New(i18n.Show(lang, `please_select_mcp_tool`))
	}
	if cast.ToInt(info[`has_auth`]) != 1 {
		return errors.New(i18n.Show(lang, `please_authorize_mcp_tool_first`))
	}
	var tools []mcp.Tool
	err = json.Unmarshal([]byte(info[`tools`]), &tools)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	if len(tools) == 0 {
		return errors.New(i18n.Show(lang, `no_available_tool_found`))
	}

	var mcpTool *mcp.Tool
	for _, t := range tools {
		if t.Name == params.ToolName {
			mcpTool = &t
			break
		}
	}
	if mcpTool == nil {
		return errors.New(i18n.Show(lang, `mcp_tool_not_matched`))
	}

	if err := common.ValidateMcpToolArguments(*mcpTool, params.Arguments); err != nil {
		return err
	}

	return nil
}

func (params *LoopNodeParams) Verify(nodeName string, lang string) error {
	if !tool.InArray(params.LoopType, []string{common.LoopTypeArray, common.LoopTypeNumber}) {
		return errors.New(i18n.Show(lang, `loop_type_error`, nodeName))
	}
	if params.LoopType == common.LoopTypeArray {
		if len(params.LoopArrays) == 0 {
			return errors.New(i18n.Show(lang, `loop_array_cannot_be_empty`, nodeName))
		}
	} else {
		if params.LoopNumber <= 0 {
			return errors.New(i18n.Show(lang, `loop_number_must_greater_than_zero`, nodeName))
		}
	}
	return nil
}

// VerityLoopParams validates output or loop parameters
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

func VerityLoopWorkflowNodes(adminUserId int, loopNode WorkFlowNode, nodeList []WorkFlowNode, allowNodeTypes []int, nodeTypeDesc string, lang string) (startNodeKey, modelConfigIds, libraryIds string, err error) {
	startNodeCount, finishNodeCount := 0, 0
	fromNodes := make(FromNodes)
	for i, node := range nodeList {
		if !tool.InArrayInt(node.NodeType.Int(), allowNodeTypes) {
			err = errors.New(i18n.Show(lang, `workflow_sub_node_type_error`, nodeTypeDesc, node.NodeType))
			return
		}
		//node base verify
		if err = node.Verify(adminUserId, lang); err != nil {
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
		err = errors.New(i18n.Show(lang, `workflow_subtype_only_one_entry_node`, nodeTypeDesc))
		return
	}
	if finishNodeCount == 0 {
		// err = errors.New(`Workflow must have a termination loop or exit node`)
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
			err = errors.New(i18n.Show(lang, `workflow_subtype_has_isolated_node`, nodeTypeDesc, node.NodeName))
			return
		}
		//skip verification temporarily - loop nodes executed separately will lack inputs from all previous nodes, making code difficult to handle
		//err = verifyNode(adminUserId, node, fromNodes, make([]WorkFlowNode, 0))
		//if err != nil {
		//	return
		//}
	}
	var libraryArr []string
	//collect model IDs being used
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
func (param *PluginNodeParams) Verify(adminUserId int, lang string) error {
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

		// Validate schema
		schema, ok := resp.Data.(map[string]any)
		if !ok {
			return errors.New(i18n.Show(lang, `invalid_schema_format`))
		}

		for key, raw := range schema {
			rule, _ := raw.(map[string]any)
			req, _ := rule["required"].(bool)
			typ, _ := rule["type"].(string)
			val, exists := param.Params[key]
			if req && !exists {
				return errors.New(i18n.Show(lang, `missing_required_field`, key))
			}
			if !exists {
				continue
			}

			switch typ {
			case "string":
				if _, ok := val.(string); !ok {
					return errors.New(i18n.Show(lang, `field_must_be_string`, key))
				}
			case "number":
				if _, ok := val.(float64); !ok {
					return errors.New(i18n.Show(lang, `field_must_be_number`, key))
				}
			case "boolean":
				if _, ok := val.(bool); !ok {
					return errors.New(i18n.Show(lang, `field_must_be_boolean`, key))
				}
			default:
				return errors.New(i18n.Show(lang, `unknown_type_for_field`, key, typ))
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
		//replace placeholders before plugin validation
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

func (params *BatchNodeParams) Verify(nodeName string, lang string) error {
	if len(params.BatchArrays) == 0 {
		return errors.New(i18n.Show(lang, `batch_execution_array_cannot_be_empty`, nodeName))
	}
	if params.ChanNumber.Int() < 1 || params.ChanNumber.Int() > 10 {
		return errors.New(i18n.Show(lang, `batch_concurrent_execution_number_error`, nodeName))
	}
	if params.MaxRunNumber.Int() < 1 || params.MaxRunNumber.Int() > 500 {
		return errors.New(i18n.Show(lang, `batch_max_run_number_error`, nodeName))
	}
	return nil
}

func (params *FinishNodeParams) Verify(nodeName string, lang string) error {
	if len(params.OutType) > 0 && !tool.InArray(params.OutType, []string{define.FinishNodeOutTypeMessage, define.FinishNodeOutTypeVariable}) {
		return errors.New(i18n.Show(lang, `output_type_error`, nodeName))
	}
	return nil
}

func (params *ImageGenerationParams) Verify(nodeName string, lang string) error {
	if cast.ToInt(params.ImageNum) < 0 || cast.ToInt(params.ImageNum) > 15 {
		return errors.New(i18n.Show(lang, `image_number_error`, nodeName))
	}
	if !tool.InArray(params.Size, define.ImageSizes) {
		return errors.New(i18n.Show(lang, `image_size_error`, nodeName))
	}
	return nil
}

func (params *JsonEncodeParams) Verify(nodeName string, lang string) error {
	if len(params.InputVariable) == 0 {
		return errors.New(i18n.Show(lang, `json_encode_input_missing`, nodeName))
	}
	return nil
}

func (params *JsonDecodeParams) Verify(nodeName string, lang string) error {
	if len(params.InputVariable) == 0 {
		return errors.New(i18n.Show(lang, `json_decode_input_missing`, nodeName))
	}
	return nil
}

func (params *TextToAudioNodeParams) Verify(lang string) error {
	if params.Arguments.ModelId <= 0 {
		return errors.New(i18n.Show(lang, `please_select_model_config`))
	}

	validVoiceTypes := []string{"system", "voice_cloning", "voice_generation", "all"}
	if len(params.VoiceType) > 0 && !tool.InArrayString(params.VoiceType, validVoiceTypes) {
		return errors.New(i18n.Show(lang, `voice_type_param_error`))
	}

	// Validate Text, max 10000 characters
	if len(params.Arguments.Text) == 0 {
		return errors.New(i18n.Show(lang, `text_content_cannot_be_empty`))
	}
	if len(params.Arguments.Text) > 10000 {
		return errors.New(i18n.Show(lang, `text_content_length_exceed_limit`))
	}

	// Validate VoiceSetting
	voiceSetting := params.Arguments.VoiceSetting
	if len(voiceSetting.VoiceId) == 0 {
		return errors.New(i18n.Show(lang, `voice_setting_voice_id_cannot_be_empty`))
	}

	// Validate Speed - recommended range [0.5, 2.0]
	if voiceSetting.Speed < 0 || voiceSetting.Speed > 100 {
		return errors.New(i18n.Show(lang, `voice_setting_speed_range_error`))
	}

	// Validate Vol - volume range
	if voiceSetting.Vol < 0 || voiceSetting.Vol > 100 {
		return errors.New(i18n.Show(lang, `voice_setting_vol_range_error`))
	}

	// Validate Pitch - pitch range
	if voiceSetting.Pitch < -100 || voiceSetting.Pitch > 100 {
		return errors.New(i18n.Show(lang, `voice_setting_pitch_range_error`))
	}

	// Validate AudioSetting - audio settings
	audioSetting := params.Arguments.AudioSetting
	// Validate SampleRate - sample rate
	if audioSetting.SampleRate > 0 {
		validSampleRates := []int{8000, 16000, 22050, 24000, 32000, 44100}
		if !tool.InArrayInt(audioSetting.SampleRate, validSampleRates) {
			return errors.New(i18n.Show(lang, `audio_setting_sample_rate_error`))
		}
	}

	// Validate Bitrate - bit rate
	if audioSetting.Bitrate > 0 {
		validBitrates := []int{32000, 64000, 128000, 256000}
		if !tool.InArrayInt(audioSetting.Bitrate, validBitrates) {
			return errors.New(i18n.Show(lang, `audio_setting_bitrate_error`))
		}
	}

	// Validate Format - audio format
	if len(audioSetting.Format) > 0 {
		validFormats := []string{"mp3", "pcm", "flac", "wav"}
		if !tool.InArrayString(audioSetting.Format, validFormats) {
			return errors.New(i18n.Show(lang, `audio_setting_format_error`))
		}
	}

	// Validate Channel - channel count
	if audioSetting.Channel > 0 {
		if audioSetting.Channel != 1 && audioSetting.Channel != 2 {
			return errors.New(i18n.Show(lang, `audio_setting_channel_error`))
		}
	}

	// Validate LanguageBoost
	if len(params.Arguments.LanguageBoost) > 0 {
		validLanguageBoosts := []string{"auto", "Chinese", "English", "Japanese", "Korean", "French", "German", "Spanish", "Russian", "Arabic"}
		if !tool.InArrayString(params.Arguments.LanguageBoost, validLanguageBoosts) {
			return errors.New(i18n.Show(lang, `language_boost_param_error`))
		}
	}

	return nil
}

func (param *VoiceCloneNodeParams) Verify(adminUserId int, lang string) error {
	if len(param.Arguments.FileUrl) == 0 {
		return errors.New(i18n.Show(lang, `file_url_cannot_be_empty`))
	}
	// Validate VoiceId - MiniMax requirement: length 8-256, starts with letter, allows alphanumeric chars, -, _, end cannot be - or _
	if len(param.Arguments.VoiceId) == 0 {
		return errors.New(i18n.Show(lang, `voice_id_cannot_be_empty`))
	}

	// Validate ClonePrompt - optional parameter, but if provided must be complete
	if len(param.Arguments.ClonePrompt.PromptAudioUrl) > 0 || len(param.Arguments.ClonePrompt.PromptText) > 0 {
		// If example audio URL is provided, validate format
		if len(param.Arguments.ClonePrompt.PromptAudioUrl) == 0 {
			return errors.New(i18n.Show(lang, `prompt_audio_url_cannot_be_empty`))
		}
		if _, err := url.Parse(param.Arguments.ClonePrompt.PromptAudioUrl); err != nil {
			return errors.New(i18n.Show(lang, `clone_prompt_audio_url_format_error`))
		}

		// If example audio is provided, validate text
		if len(param.Arguments.ClonePrompt.PromptText) == 0 {
			return errors.New(i18n.Show(lang, `prompt_text_cannot_be_empty`))
		}
	}

	return nil
}

func (param *LibraryImportParams) Verify(adminUserId int, lang string) error {
	if cast.ToInt(param.LibraryId) == 0 {
		return errors.New(i18n.Show(lang, `please_select_library`))
	}
	libraryInfo, err := common.GetLibraryInfo(cast.ToInt(param.LibraryId), adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return errors.New(i18n.Show(lang, `get_library_detail_failed`))
	}
	if len(libraryInfo) == 0 {
		return errors.New(i18n.Show(lang, `library_not_exist`))
	}
	if cast.ToInt(param.LibraryGroupId) > 0 {
		libraryGroupInfo, sqlErr := msql.Model(`chat_ai_library_group`, define.Postgres).
			Where(`library_id`, cast.ToString(param.LibraryId)).
			Where(`id`, param.LibraryGroupId).Find()
		if sqlErr != nil {
			logs.Error(sqlErr.Error())
			return errors.New(i18n.Show(lang, `get_library_group_detail_failed`))
		}
		if len(libraryGroupInfo) == 0 {
			return errors.New(i18n.Show(lang, `library_group_not_exist`))
		}
	}

	if cast.ToInt(libraryInfo[`type`]) == define.GeneralLibraryType {
		if param.ImportType == define.LibraryImportContent {
			if param.NormalTitle == `` || param.NormalContent == `` {
				return errors.New(i18n.Show(lang, `please_fill_content_and_title`))
			}
		} else if param.ImportType == define.LibraryImportUrl {
			if param.NormalUrl == `` {
				return errors.New(i18n.Show(lang, `please_fill_import_url`))
			}
		}
	} else if cast.ToInt(libraryInfo[`type`]) == define.QALibraryType {
		if param.QaQuestion == `` || param.QaAnswer == `` {
			return errors.New(i18n.Show(lang, `please_fill_question_and_answer`))
		}
	}
	return nil
}

func (param *WorkflowNodeParams) Verify(adminUserId int, lang string) error {
	if param.RobotId <= 0 {
		return errors.New(i18n.Show(lang, `robot_id_cannot_be_empty`))
	}
	info, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(param.RobotId)).
		Find()
	if err != nil {
		return err
	}
	if len(info) == 0 || cast.ToInt(info[`application_type`]) != define.ApplicationTypeFlow {
		return errors.New(i18n.Show(lang, `selected_workflow_invalid`))
	}

	maps := map[string]struct{}{}
	for _, item := range param.Params {
		if !common.IsVariableName(item.Key) {
			return errors.New(i18n.Show(lang, `custom_global_variable_name_format_error`, item.Key))
		}
		if tool.InArrayString(fmt.Sprintf(`global.%s`, item.Key), SysGlobalVariables()) {
			return errors.New(i18n.Show(lang, `custom_global_variable_conflict_with_system`, item.Key))
		}
		if !tool.InArrayString(item.Typ, []string{common.TypString, common.TypNumber, common.TypArrString, common.TypArrObject}) {
			return errors.New(i18n.Show(lang, `custom_global_variable_type_not_supported`, item.Key))
		}
		if item.Required && len(item.Variable) == 0 {
			return errors.New(i18n.Show(lang, `required_param_missing`, item.Key))
		}
		if _, ok := maps[item.Key]; ok {
			return errors.New(i18n.Show(lang, `custom_global_variable_duplicate_definition`, item.Key))
		}
		maps[item.Key] = struct{}{}
	}

	return nil
}

func (params *ImmediatelyReplyNodeParams) Verify(lang string) error {
	if len(params.Content) == 0 {
		return errors.New(i18n.Show(lang, `message_content_cannot_be_empty`))
	}
	return nil
}

type QuestionMenu struct {
	MenuLabel string `json:"menu_label"`
}

type ReplyContent struct {
	ReplyType string    `json:"reply_type"`
	Type      string    `json:"type" form:"type"`
	SmartMenu SmartMenu `json:"smart_menu,omitempty" form:"smart_menu"` // Smart menu, passed when outputting
}

type SmartMenu struct {
	MenuDescription string             `json:"menu_description"`
	MenuContent     []SmartMenuContent `json:"menu_content"`
}

type SmartMenuContent struct {
	MenuType    string `json:"menu_type"`     // menu type: 0 normal, 1 keyword click menu
	SerialNo    string `json:"serial_no"`     // serial number
	Content     string `json:"content"`       // content; if keyword click menu, this is the keyword
	NextNodeKey string `json:"next_node_key"` // next node
}

type QuestionParams struct {
	AnswerType       string               `json:"answer_type"`
	AnswerText       string               `json:"answer_text"`
	ReplyContentList []ReplyContent       `json:"reply_content_list"`
	Outputs          common.RecurveFields `json:"outputs"`
}

func (params *QuestionParams) Verify(nodeName string, lang string) error {
	if !tool.InArray(params.AnswerType, []string{define.QuestionAnswerTypeText, define.QuestionAnswerTypeMenu}) {
		return errors.New(i18n.Show(lang, `node_answer_type_param_error`, nodeName))
	}
	if params.AnswerType == define.QuestionAnswerTypeMenu {
		if len(params.ReplyContentList) == 0 {
			return errors.New(i18n.Show(lang, `smart_menu_param_error`, nodeName))
		}
		replyContent := params.ReplyContentList[0]
		if replyContent.ReplyType != common.ReplyTypeSmartMenu {
			return errors.New(i18n.Show(lang, `smart_menu_item_type_error`, nodeName))
		}
		if len(replyContent.SmartMenu.MenuContent) == 0 {
			return errors.New(i18n.Show(lang, `smart_menu_at_least_one`, nodeName))
		}
	} else if params.AnswerType == define.QuestionAnswerTypeText {
		if len(params.AnswerText) == 0 {
			return errors.New(i18n.Show(lang, `question_content_param_error`, nodeName))
		}
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
