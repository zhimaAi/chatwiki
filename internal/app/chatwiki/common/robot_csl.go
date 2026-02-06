// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_crypto"
	"chatwiki/internal/pkg/lib_web"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

const RobotCslKey = `A8R0pvUT7fuCLm4idsBV1IkSD6EZg9HO`

type RobotCsl struct {
	StartTime time.Time           `json:"start_time"`
	EndTime   time.Time           `json:"end_time"`
	Robot     msql.Params         `json:"robot"`               // Robot configuration information
	FileName  string              `json:"file_name"`           // Exported file name
	Category  msql.Params         `json:"category"`            // File segment refinement configuration
	Librarys  map[int]*LibraryCsl `json:"librarys,omitempty"`  // Knowledge bases involved in the robot
	Forms     map[int]*FormCsl    `json:"forms,omitempty"`     // Data tables involved in the robot
	Nodes     []msql.Params       `json:"nodes"`               // Workflow draft node data
	Workflows []*RobotCsl         `json:"workflows,omitempty"` // Workflows associated with the robot
}

func NewRobotCsl() *RobotCsl {
	robotCsl := RobotCsl{
		StartTime: time.Now(),
		Robot:     make(msql.Params),
		FileName:  fmt.Sprintf(`robot_%s.csl`, tool.Date(`YmdHis`)),
		Category:  make(msql.Params),
		Librarys:  make(map[int]*LibraryCsl),
		Forms:     make(map[int]*FormCsl),
		Nodes:     make([]msql.Params, 0),
		Workflows: make([]*RobotCsl, 0),
	}
	return &robotCsl
}

func (robotCsl *RobotCsl) Output() (content string, err error) {
	robotCsl.EndTime = time.Now()
	origData, err := json.Marshal(robotCsl)
	if err != nil {
		return
	}
	crypted, err := lib_crypto.AesEncrypt(origData, []byte(RobotCslKey))
	if err != nil {
		return
	}
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func ParseRobotCsl(content string) (robotCsl *RobotCsl, err error) {
	crypted, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return
	}
	decrypt, err := lib_crypto.AesDecrypt(crypted, []byte(RobotCslKey))
	if err != nil {
		return
	}
	robotCsl = NewRobotCsl()
	err = json.Unmarshal(decrypt, &robotCsl)
	return
}

type LibraryCsl struct {
	Library       msql.Params   `json:"library"`        // Knowledge base configuration information
	QuestionGuide []msql.Params `json:"question_guide"` // External document guide
	LibGroups     []msql.Params `json:"lib_groups"`     // Knowledge base grouping information
	LibFiles      []*LibFileCsl `json:"lib_files"`      // Knowledge base file information
	FileDocs      []msql.Params `json:"file_docs"`      // External document file list
}

func BuildLibraryCsl(lang string, libraryId, adminUserId int, importData bool) (libraryCsl *LibraryCsl, err error) {
	libraryCsl = &LibraryCsl{
		Library:       make(msql.Params),
		QuestionGuide: make([]msql.Params, 0),
		LibGroups:     make([]msql.Params, 0),
		LibFiles:      make([]*LibFileCsl, 0),
		FileDocs:      make([]msql.Params, 0),
	}
	if libraryId <= 0 {
		err = errors.New(i18n.Show(lang, `library_id_param_error`))
		return
	}
	// Knowledge base configuration information
	library, sqlErr := GetLibraryInfo(libraryId, adminUserId)
	if sqlErr != nil {
		err = sqlErr
		return
	}
	if len(library) == 0 {
		err = errors.New(i18n.Show(lang, `library_info_not_exist`))
		return
	}
	libraryCsl.Library = library

	if !importData {
		return
	}

	// External document guide
	questionGuide, sqlErr := msql.Model(`chat_ai_library_question_guide `, define.Postgres).
		Where(`library_id`, cast.ToString(libraryId)).Order(`id`).Select()
	if sqlErr != nil {
		err = sqlErr
		return
	}
	libraryCsl.QuestionGuide = questionGuide
	// Knowledge base grouping information
	libGroups, sqlErr := msql.Model(`chat_ai_library_group`, define.Postgres).
		Where(`library_id`, cast.ToString(libraryId)).Order(`id`).Select()
	if sqlErr != nil {
		err = sqlErr
		return
	}
	libraryCsl.LibGroups = libGroups
	// Knowledge base file information
	libFiles, sqlErr := msql.Model(`chat_ai_library_file`, define.Postgres).
		Where(`library_id`, cast.ToString(libraryId)).Order(`id`).Select()
	if sqlErr != nil {
		err = sqlErr
		return
	}
	for _, libFile := range libFiles {
		libFileCsl, sqlErr := BuildLibFileCsl(libFile)
		if sqlErr == nil {
			libraryCsl.LibFiles = append(libraryCsl.LibFiles, libFileCsl)
		} else {
			logs.Error(sqlErr.Error())
		}
	}
	// External document file list
	fileDocs, sqlErr := msql.Model(`chat_ai_library_file_doc `, define.Postgres).Where(`delete_time`, `0`).
		Where(`library_id`, cast.ToString(libraryId)).Order(`sort desc`).Order(`id`).Select()
	if sqlErr != nil {
		err = sqlErr
		return
	}
	libraryCsl.FileDocs = fileDocs
	return
}

func RequestChatWiki(lang string, path, method, token string, params map[string]string) (lib_web.Response, error) {
	link := fmt.Sprintf(`http://127.0.0.1:%s%s`, define.Config.WebService[`port`], path)
	request := curl.NewRequest(link, method).Header(`token`, token).Header(`lang`, lang)
	for key, item := range params {
		request.Param(key, item)
	}
	resp, err := request.Response()
	if err != nil {
		return lib_web.Response{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return lib_web.Response{}, errors.New(fmt.Sprintf(`SYSTEM ERROR:%d`, resp.StatusCode))
	}
	code := lib_web.Response{}
	if err = request.ToJSON(&code); err != nil {
		return lib_web.Response{}, err
	}
	if code.Res != lib_web.CommonSuccess {
		return code, errors.New(code.Msg)
	}
	return code, nil
}

func (libraryCsl *LibraryCsl) Import(lang string, adminUserId, userId int, cslIdMaps *CslIdMaps, models *DefaultModelParams, token string) error {
	// Knowledge base configuration information
	oldLibraryId := cast.ToInt(libraryCsl.Library[`id`])
	libraryData := make(map[string]string)
	for key, val := range libraryCsl.Library {
		switch key {
		case `id`, `admin_user_id`, `share_url`, `creator`, `group_id`:
		case `model_config_id`:
			if cast.ToInt(val) > 0 {
				libraryData[key] = cast.ToString(models.VectorModelConfigId)
			}
		case `use_model`:
			if len(val) > 0 {
				libraryData[key] = models.VectorUseModel
			}
		case `summary_model_config_id`, `graph_model_config_id`, `ai_chunk_model_config_id`:
			if cast.ToInt(val) > 0 {
				libraryData[key] = cast.ToString(models.LlmModelConfigId)
			}
		case `ai_summary_model`, `graph_use_model`, `ai_chunk_model`:
			if len(val) > 0 {
				libraryData[key] = models.LlmUseModel
			}
		case `avatar`:
			libraryData[`avatar_from_template`] = val
		default:
			libraryData[key] = val
		}
	}

	base_library_name := libraryData[`library_name`]
	library_name_arr := strings.Split(libraryData[`library_name`], "_")
	if cast.ToInt(library_name_arr[len(library_name_arr)-1]) > 0 { // If the last character is a digit
		base_library_name = strings.Join(library_name_arr[:len(library_name_arr)-1], "_")
	}

	libraryLen, _ := msql.Model(`chat_ai_library`, define.Postgres).
		Where(`library_name`, `like`, base_library_name).Count()
	if libraryLen > 0 {
		base_library_name = base_library_name + "_" + cast.ToString(libraryLen)
	}
	libraryData[`library_name`] = base_library_name

	libraryData[`is_default`] = cast.ToString(define.NotDefault)
	code, err := RequestChatWiki(lang, `/manage/createLibrary`, http.MethodPost, token, libraryData)
	if err != nil {
		return err
	}
	newLibraryId := cast.ToInt(cast.ToStringMap(code.Data)[`id`])
	cslIdMaps.Librarys[oldLibraryId] = newLibraryId
	// External document guide
	for _, questionGuide := range libraryCsl.QuestionGuide {
		questionGuideData := msql.Datas{`admin_user_id`: adminUserId, `library_id`: newLibraryId}
		for key, val := range questionGuide {
			if !tool.InArrayString(key, []string{`id`, `admin_user_id`, `library_id`}) {
				questionGuideData[key] = val
			}
		}
		_, err = msql.Model(`chat_ai_library_question_guide`, define.Postgres).Insert(questionGuideData, `id`)
		if err != nil {
			logs.Error(err.Error())
		}
	}
	// Knowledge base grouping information
	for _, libGroup := range libraryCsl.LibGroups {
		oldLibGroupId := cast.ToInt(libGroup[`id`])
		libGroupData := msql.Datas{`admin_user_id`: adminUserId, `library_id`: newLibraryId}
		for key, val := range libGroup {
			if !tool.InArrayString(key, []string{`id`, `admin_user_id`, `library_id`}) {
				libGroupData[key] = val
			}
		}
		newLibGroupId, err := msql.Model(`chat_ai_library_group`, define.Postgres).Insert(libGroupData, `id`)
		if err != nil {
			logs.Error(err.Error())
		}
		cslIdMaps.LibGroups[oldLibGroupId] = int(newLibGroupId)
	}
	// Knowledge base file information
	for _, LibFile := range libraryCsl.LibFiles {
		err = LibFile.Import(lang, adminUserId, newLibraryId, cslIdMaps, models, token)
		if err != nil {
			logs.Error(err.Error())
		}
	}
	// External document file list
	tempMaps := make(map[int]int) // New file id => old pid
	for _, fileDoc := range libraryCsl.FileDocs {
		oldFileDocId := cast.ToInt(fileDoc[`id`])
		oldFileDocPId := cast.ToInt(fileDoc[`pid`])
		newLibFileId := cslIdMaps.LibFiles[cast.ToInt(fileDoc[`file_id`])]
		newFileDocId, err := SaveLibDoc(adminUserId, userId, newLibraryId, 0, newLibFileId, 0,
			cast.ToInt(fileDoc[`is_index`]), cast.ToInt(fileDoc[`is_draft`]), fileDoc[`title`], nil, fileDoc[`content`], cast.ToInt(fileDoc[`is_dir`]), fileDoc[`doc_icon`])
		if err != nil {
			logs.Error(err.Error())
		}
		if oldFileDocPId > 0 { // If there is a parent, record it first
			tempMaps[newFileDocId] = oldFileDocPId
		}
		cslIdMaps.FileDocs[oldFileDocId] = newFileDocId
	}
	for newFileDocId, oldFileDocPId := range tempMaps {
		if newFileDocPId := cslIdMaps.FileDocs[oldFileDocPId]; newFileDocPId > 0 { // Update pid relationship
			if _, err = ChangeLibDoc(newFileDocId, msql.Datas{`pid`: newFileDocPId}); err != nil {
				logs.Error(err.Error())
			}
		}
	}
	return nil
}

type LibFileCsl struct {
	LibFile  msql.Params   `json:"lib_file"`  // Knowledge base file information
	FileData []msql.Params `json:"file_data"` // File segment list
}

func BuildLibFileCsl(libFile msql.Params) (libFileCsl *LibFileCsl, err error) {
	libFileCsl = &LibFileCsl{
		LibFile:  libFile,
		FileData: make([]msql.Params, 0),
	}
	// File segment list
	fileData, sqlErr := msql.Model(`chat_ai_library_file_data `, define.Postgres).
		Where(`file_id`, libFile[`id`]).Order(`page_num,father_chunk_paragraph_number,number`).Order(`id`).Select()
	if sqlErr != nil {
		err = sqlErr
		return
	}
	libFileCsl.FileData = fileData
	return
}

func (libFileCsl *LibFileCsl) Import(lang string, adminUserId, libraryId int, cslIdMaps *CslIdMaps, models *DefaultModelParams, token string) error {
	// Knowledge base file information
	oldLibFileId := cast.ToInt(libFileCsl.LibFile[`id`])
	libFileData := msql.Datas{`admin_user_id`: adminUserId, `library_id`: libraryId}
	for key, val := range libFileCsl.LibFile {
		switch key {
		case `id`, `admin_user_id`, `library_id`, `ai_chunk_task_id`, `ali_ocr_job_id`, `group_id`:
		case `status`:
			libFileData[key] = define.FileStatusLearned // Set as completed
		case `semantic_chunk_model_config_id`:
			if cast.ToInt(val) > 0 {
				libFileData[key] = models.VectorModelConfigId
			}
		case `semantic_chunk_use_model`:
			if len(val) > 0 {
				libFileData[key] = models.VectorUseModel
			}
		case `ai_chunk_model_config_id`:
			if cast.ToInt(val) > 0 {
				libFileData[key] = models.LlmModelConfigId
			}
		case `ai_chunk_model`:
			if len(val) > 0 {
				libFileData[key] = models.LlmUseModel
			}
		default:
			libFileData[key] = val
		}
	}
	newLibFileId, err := msql.Model(`chat_ai_library_file`, define.Postgres).Insert(libFileData, `id`)
	if err != nil {
		return err
	}
	cslIdMaps.LibFiles[oldLibFileId] = int(newLibFileId)
	// File segment list
	for _, fileData := range libFileCsl.FileData {
		params := map[string]string{`library_id`: cast.ToString(libraryId), `file_id`: cast.ToString(newLibFileId)}
		for key, val := range fileData {
			switch key {
			case `id`, `admin_user_id`, `library_id`, `file_id`:
			case `images`:
				params[`images_json`] = val
			case `category_id`:
				params[key] = cast.ToString(cslIdMaps.Category[cast.ToInt(val)])
			case `group_id`:
				params[key] = cast.ToString(cslIdMaps.LibGroups[cast.ToInt(val)])
			default:
				params[key] = val
			}
		}
		_, err = RequestChatWiki(lang, `/manage/addParagraph`, http.MethodPost, token, params)
		if err != nil {
			logs.Error(err.Error())
		}
	}
	return nil
}

type FormCsl struct {
	Form                msql.Params   `json:"form"`                  // Data table configuration information
	FormEntry           []msql.Params `json:"form_entry"`            // Data table row data
	FormField           []msql.Params `json:"form_field"`            // Data table field information
	FormFieldValue      []msql.Params `json:"form_field_value"`      // Data table key-value pairs
	FormFilter          []msql.Params `json:"form_filter"`           // Data table classification configuration
	FormFilterCondition []msql.Params `json:"form_filter_condition"` // Data table classification conditions
}

func BuildFormCsl(lang string, formId, adminUserId int, importData bool) (formCsl *FormCsl, err error) {
	formCsl = &FormCsl{
		Form:                make(msql.Params),
		FormEntry:           make([]msql.Params, 0),
		FormField:           make([]msql.Params, 0),
		FormFieldValue:      make([]msql.Params, 0),
		FormFilter:          make([]msql.Params, 0),
		FormFilterCondition: make([]msql.Params, 0),
	}
	if formId <= 0 {
		err = errors.New(i18n.Show(lang, `form_id_param_error`))
		return
	}
	// Data table configuration information
	form, sqlErr := msql.Model(`form`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`id`, cast.ToString(formId)).Find()
	if sqlErr != nil {
		err = sqlErr
		return
	}
	if len(form) == 0 {
		err = errors.New(i18n.Show(lang, `form_info_not_exist`))
		return
	}
	formCsl.Form = form
	// Data table row data
	if importData {
		formEntry, sqlErr := msql.Model(`form_entry`, define.Postgres).
			Where(`form_id`, cast.ToString(formId)).Where(`delete_time`, `0`).Order(`id`).Select()
		if sqlErr != nil {
			err = sqlErr
			return
		}
		formCsl.FormEntry = formEntry
	}
	// Data table field information
	formField, sqlErr := msql.Model(`form_field`, define.Postgres).
		Where(`form_id`, cast.ToString(formId)).Order(`id`).Select()
	if sqlErr != nil {
		err = sqlErr
		return
	}
	formCsl.FormField = formField
	// Data table key-value pairs
	if len(formCsl.FormEntry) > 0 {
		formEntryIds := make([]string, 0)
		for _, item := range formCsl.FormEntry {
			formEntryIds = append(formEntryIds, item[`id`])
		}
		formFieldValue, sqlErr := msql.Model(`form_field_value`, define.Postgres).
			Where(`form_entry_id`, `in`, strings.Join(formEntryIds, `,`)).Order(`id`).Select()
		if sqlErr != nil {
			err = sqlErr
			return
		}
		formCsl.FormFieldValue = formFieldValue
	}
	// Data table classification configuration
	formFilter, sqlErr := msql.Model(`form_filter`, define.Postgres).
		Where(`form_id`, cast.ToString(formId)).Order(`id`).Select()
	if sqlErr != nil {
		err = sqlErr
		return
	}
	formCsl.FormFilter = formFilter
	// Data table classification conditions
	if len(formFilter) > 0 {
		formFilterIds := make([]string, 0)
		for _, item := range formFilter {
			formFilterIds = append(formFilterIds, cast.ToString(item[`id`]))
		}
		formFilterCondition, sqlErr := msql.Model(`form_filter_condition`, define.Postgres).
			Where(`form_filter_id`, `in`, strings.Join(formFilterIds, `,`)).Order(`id`).Select()
		if sqlErr != nil {
			err = sqlErr
			return
		}
		formCsl.FormFilterCondition = formFilterCondition
	}
	return
}

func (formCsl *FormCsl) Import(adminUserId int, cslIdMaps *CslIdMaps) error {
	// Data table configuration information
	oldFormId := cast.ToInt(formCsl.Form[`id`])
	formData := msql.Datas{`admin_user_id`: adminUserId}
	for key, val := range formCsl.Form {
		if !tool.InArrayString(key, []string{`id`, `admin_user_id`}) {
			formData[key] = val
		}
	}

	tableName := strings.Split(cast.ToString(formData["name"]), "_")
	realTableName := ""
	tableSuffix := cast.ToInt(tableName[len(tableName)-1])
	if tableSuffix > 0 { // If there is a suffix
		realTableName = strings.Join(tableName[:len(tableName)-1], "_")
	} else {
		realTableName = strings.Join(tableName, "_")
	}

	// Validate table with same name
	dbNameCount, err := msql.Model(`form`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`name`, `like`, realTableName).Count(`id`)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	if dbNameCount > 0 {
		formData[`name`] = realTableName + "_" + cast.ToString(dbNameCount+1)
	}

	newFormId, err := msql.Model(`form`, define.Postgres).Insert(formData, `id`)
	if err != nil {
		return err
	}
	cslIdMaps.Forms[oldFormId] = int(newFormId)
	// New table name
	cslIdMaps.FormsName[int(newFormId)] = cast.ToString(formData[`name`])
	// Data table field information
	for _, formField := range formCsl.FormField {
		oldFormFieldId := cast.ToInt(formField[`id`])
		formFieldData := msql.Datas{`admin_user_id`: adminUserId, `form_id`: newFormId}
		for key, val := range formField {
			if !tool.InArrayString(key, []string{`id`, `admin_user_id`, `form_id`}) {
				formFieldData[key] = val
			}
		}
		newFormFieldId, err := msql.Model(`form_field`, define.Postgres).Insert(formFieldData, `id`)
		if err != nil {
			return err
		}
		cslIdMaps.FormFields[oldFormFieldId] = int(newFormFieldId)
	}

	// Data table row data
	for _, formEntry := range formCsl.FormEntry {
		oldFormEntryId := cast.ToInt(formEntry[`id`])
		formEntryData := msql.Datas{`admin_user_id`: adminUserId, `form_id`: newFormId}
		for key, val := range formEntry {
			if !tool.InArrayString(key, []string{`id`, `admin_user_id`, `form_id`}) {
				formEntryData[key] = val
			}
		}
		newFormEntryId, err := msql.Model(`form_entry`, define.Postgres).Insert(formEntryData, `id`)
		if err != nil {
			return err
		}
		cslIdMaps.FormEntrys[oldFormEntryId] = int(newFormEntryId)
	}

	// Data table key-value pairs
	for _, formFieldValue := range formCsl.FormFieldValue {
		formEntryId := cslIdMaps.FormEntrys[cast.ToInt(formFieldValue[`form_entry_id`])]
		formFieldId := cslIdMaps.FormFields[cast.ToInt(formFieldValue[`form_field_id`])]
		formFieldValueData := msql.Datas{`admin_user_id`: adminUserId, `form_entry_id`: formEntryId, `form_field_id`: formFieldId}
		for key, val := range formFieldValue {
			if !tool.InArrayString(key, []string{`id`, `admin_user_id`, `form_entry_id`, `form_field_id`}) {
				formFieldValueData[key] = val
			}
		}
		_, err = msql.Model(`form_field_value`, define.Postgres).Insert(formFieldValueData, `id`)
		if err != nil {
			return err
		}
	}
	// Data table classification configuration
	for _, formFilter := range formCsl.FormFilter {
		oldFormFilterId := cast.ToInt(formFilter[`id`])
		formFilterData := msql.Datas{`form_id`: newFormId}
		for key, val := range formFilter {
			if !tool.InArrayString(key, []string{`id`, `form_id`}) {
				formFilterData[key] = val
			}
		}
		newFormFilterId, err := msql.Model(`form_filter`, define.Postgres).Insert(formFilterData, `id`)
		if err != nil {
			return err
		}
		cslIdMaps.FormFilters[oldFormFilterId] = int(newFormFilterId)
	}
	// Data table classification conditions
	for _, formFilterCondition := range formCsl.FormFilterCondition {
		formFilterId := cslIdMaps.FormFilters[cast.ToInt(formFilterCondition[`form_filter_id`])]
		formFieldId := cslIdMaps.FormFields[cast.ToInt(formFilterCondition[`form_field_id`])]
		formFilterConditionData := msql.Datas{`form_filter_id`: formFilterId, `form_field_id`: formFieldId}
		for key, val := range formFilterCondition {
			if !tool.InArrayString(key, []string{`id`, `form_filter_id`, `form_field_id`}) {
				formFilterConditionData[key] = val
			}
		}
		_, err = msql.Model(`form_filter_condition`, define.Postgres).Insert(formFilterConditionData, `id`)
		if err != nil {
			return err
		}
	}
	return nil
}

type CslIdMaps struct {
	Librarys    map[int]int    // Knowledge base ID (old=>new)
	LibFiles    map[int]int    // Knowledge base file ID (old=>new)
	Category    map[int]int    // Segment classification (refined) ID (old=>new)
	LibGroups   map[int]int    // Knowledge base segment grouping ID (old=>new)
	FileDocs    map[int]int    // External document file ID (old=>new)
	Forms       map[int]int    // Data table ID (old=>new)
	FormEntrys  map[int]int    // Data table row data ID (old=>new)
	FormFields  map[int]int    // Data table field ID (old=>new)
	FormFilters map[int]int    // Data table classification ID (old=>new)
	Workflows   map[int]int    // Workflow ID (old=>new)
	FormsName   map[int]string // Data table name (old=>new)
}

func NewCslIdMaps() *CslIdMaps {
	cslIdMaps := CslIdMaps{
		Librarys:    make(map[int]int),
		LibFiles:    make(map[int]int),
		Category:    make(map[int]int),
		LibGroups:   make(map[int]int),
		FileDocs:    make(map[int]int),
		Forms:       make(map[int]int),
		FormEntrys:  make(map[int]int),
		FormFields:  make(map[int]int),
		FormFilters: make(map[int]int),
		Workflows:   make(map[int]int),
	}
	return &cslIdMaps
}

type DefaultModelParams struct {
	LlmModelConfigId    int    `json:"llm_model_config_id"`
	LlmUseModel         string `json:"llm_use_model"`
	VectorModelConfigId int    `json:"vector_model_config_id"`
	VectorUseModel      string `json:"vector_use_model"`
	RerankModelConfigId int    `json:"rerank_model_config_id"`
	RerankUseModel      string `json:"rerank_use_model"`
}

func GetDefaultModelParams(lang string, adminUserId int) (params *DefaultModelParams, err error) {
	params = &DefaultModelParams{}
	configs, sqlErr := msql.Model(`chat_ai_model_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Order(`id desc`).Select()
	if sqlErr != nil {
		err = sqlErr
		return
	}
	sort.Slice(configs, func(i, j int) bool {
		return tool.InArrayString(configs[i][`model_define`], []string{ModelChatWiki}) // Priority model to select
	})
	for _, config := range configs {
		modelInfo, ok := GetModelInfoByConfig(lang, adminUserId, cast.ToInt(config[`id`]))
		if !ok {
			continue
		}
		if params.LlmModelConfigId == 0 && tool.InArrayString(Llm, strings.Split(config[`model_types`], `,`)) {
			if models := modelInfo.GetLlmModelList(); len(models) > 0 {
				params.LlmModelConfigId = cast.ToInt(config[`id`])
				params.LlmUseModel = models[0]
				// Prefer to use models that support function calls
				if models = modelInfo.GetFunctionCallModels(); len(models) > 0 {
					params.LlmUseModel = models[0]
				}
			}
		}
		if params.VectorModelConfigId == 0 && tool.InArrayString(TextEmbedding, strings.Split(config[`model_types`], `,`)) {
			if models := modelInfo.GetVectorModelList(); len(models) > 0 {
				params.VectorModelConfigId = cast.ToInt(config[`id`])
				params.VectorUseModel = models[0]
			}
		}
		if params.RerankModelConfigId == 0 && tool.InArrayString(Rerank, strings.Split(config[`model_types`], `,`)) {
			if models := modelInfo.GetRerankModelList(); len(models) > 0 {
				params.RerankModelConfigId = cast.ToInt(config[`id`])
				params.RerankUseModel = models[0]
			}
		}
	}
	if params.LlmModelConfigId == 0 {
		return nil, errors.New(i18n.Show(lang, `llm_model_provider_not_configured`))
	}
	if params.VectorModelConfigId == 0 {
		return nil, errors.New(i18n.Show(lang, `embedding_model_provider_not_configured`))
	}
	return
}
