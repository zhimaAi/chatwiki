// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetUseGuide(adminUserId int, lang string) (*define.GuideData, error) {
	guideData, err := getGuideData(adminUserId, define.UseTypeGuide)
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(guideData.ProcessList) == 0 {
		guideData := define.GuideData{ProcessList: make([]define.GuideProcess, 0)}
		//model
		modelProcess, err := getModelProcess(adminUserId, lang)
		if err != nil {
			return nil, err
		}
		guideData.ProcessList = append(guideData.ProcessList, modelProcess)
		//library
		guideData.ProcessList = append(guideData.ProcessList, getCreateLibrary(lang))
		//robot
		guideData.ProcessList = append(guideData.ProcessList, getCreateRobot(lang))
		//test robot
		guideData.ProcessList = append(guideData.ProcessList, getTestRobot(lang))
		_, err = msql.Model(define.TableUseGuideProcess, define.Postgres).Insert(msql.Datas{
			`admin_user_id`:  adminUserId,
			`use_guide_type`: define.UseTypeGuide,
			`data`:           tool.JsonEncodeNoError(guideData),
			`create_time`:    time.Now().Unix(),
			`update_time`:    time.Now().Unix(),
		})
		if err != nil {
			logs.Error(err.Error())
			return nil, errors.New(i18n.Show(lang, `sys_err`))
		}
		lib_redis.DelCacheData(define.Redis, UseGuideCacheBuildHandler{
			AdminUserId:  adminUserId,
			UseGuideType: define.UseTypeGuide,
		})
		return &guideData, nil
	}
	return &guideData, nil
}

func getGuideData(adminUserId int, useGuideType string) (guideData define.GuideData, err error) {
	guideData = define.GuideData{
		ProcessList: make([]define.GuideProcess, 0),
	}
	guideInfo, err := msql.Model(define.TableUseGuideProcess, define.Postgres).
		Where(`use_guide_type`, useGuideType).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Find()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	if len(guideInfo) == 0 {
		return guideData, nil
	}
	data := guideInfo[`data`]
	err = tool.JsonDecode(data, &guideData)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	return
}

func getTestRobot(lang string) define.GuideProcess {
	process := define.GuideProcess{
		Name:     i18n.Show(lang, `use_guide_test_robot`),
		Key:      define.ProcessTypeTestRobot,
		IsFinish: define.StepIsNotFinish,
		Steps:    make([]define.GuideStep, 0),
	}
	testRobot := define.GuideStep{
		Name:     i18n.Show(lang, `use_guide_chat_test`),
		Key:      define.StepTestRobot,
		IsFinish: define.StepIsNotFinish,
	}
	process.Steps = append(process.Steps, testRobot)
	return process
}

func getCreateRobot(lang string) define.GuideProcess {
	process := define.GuideProcess{
		Name:     i18n.Show(lang, `RobotCreateName`),
		Key:      define.ProcessTypeCreateRobot,
		IsFinish: define.StepIsNotFinish,
		Steps:    make([]define.GuideStep, 0),
	}
	createRobot := define.GuideStep{
		Name:     i18n.Show(lang, `RobotCreateName`),
		Key:      define.StepCreateRobot,
		IsFinish: define.StepIsNotFinish,
	}
	relationLibrary := define.GuideStep{
		Name:     i18n.Show(lang, `use_guide_relation_library`),
		Key:      define.StepRelationLibrary,
		IsFinish: define.StepIsNotFinish,
	}
	process.Steps = append(process.Steps, createRobot, relationLibrary)
	return process
}

func getCreateLibrary(lang string) define.GuideProcess {
	process := define.GuideProcess{
		Name:     i18n.Show(lang, `LibraryCreateName`),
		Key:      define.ProcessTypeCreateLibrary,
		IsFinish: define.StepIsNotFinish,
		Steps:    make([]define.GuideStep, 0),
	}
	createLibrary := define.GuideStep{
		Name:     i18n.Show(lang, `LibraryCreateName`),
		Key:      define.StepCreateLibrary,
		IsFinish: define.StepIsNotFinish,
	}
	importWord := define.GuideStep{
		Name:     i18n.Show(lang, `use_guide_import_word`),
		Key:      define.StepImportWord,
		IsFinish: define.StepIsNotFinish,
	}
	importPdf := define.GuideStep{
		Name:     i18n.Show(lang, `use_guide_import_pdf`),
		Key:      define.StepImportPdf,
		IsFinish: define.StepIsNotFinish,
	}
	process.Steps = append(process.Steps, createLibrary, importWord, importPdf)
	return process
}

func getModelProcess(adminUserId int, lang string) (define.GuideProcess, error) {
	modelProcess := define.GuideProcess{
		Name:     i18n.Show(lang, `use_guide_model_check`),
		Key:      define.ProcessTypeSetModel,
		IsFinish: define.StepIsNotFinish,
		Steps:    make([]define.GuideStep, 0),
	}
	setLlm, setText, err := checkSetModel(adminUserId)
	if err != nil {
		return modelProcess, err
	}
	llmSet := define.GuideStep{
		Name:     i18n.Show(lang, `use_guide_config_llm`),
		Key:      define.StepSetLlm,
		IsFinish: define.StepIsNotFinish,
	}
	if setLlm {
		llmSet.IsFinish = define.StepIsFinish
	}
	textSet := define.GuideStep{
		Name:     i18n.Show(lang, `use_guide_config_embedding`),
		Key:      define.StepSetEmbedding,
		IsFinish: define.StepIsNotFinish,
	}
	if setText {
		textSet.IsFinish = define.StepIsFinish
	}
	if setText && setLlm {
		modelProcess.IsFinish = define.StepIsFinish
	}
	modelProcess.Steps = append(modelProcess.Steps, llmSet, textSet)
	return modelProcess, nil
}

func checkSetModel(adminUserId int) (setLlm bool, setText bool, _ error) {
	m := msql.Model(`chat_ai_model_list`, define.Postgres)
	if id, err := m.Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`model_type`, Llm).Value(`id`); err == nil && cast.ToUint(id) > 0 {
		setLlm = true
	}
	if id, err := m.Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`model_type`, TextEmbedding).Value(`id`); err == nil && cast.ToUint(id) > 0 {
		setText = true
	}
	return
}

func SetStepFinish(adminUserId int, stepType string) error {
	guideData := define.GuideData{}
	err := lib_redis.GetCacheWithBuild(define.Redis, UseGuideCacheBuildHandler{
		AdminUserId:  adminUserId,
		UseGuideType: define.UseTypeGuide,
	}, &guideData, time.Hour*12)
	if err != nil {
		return err
	}
	for key, process := range guideData.ProcessList {
		if process.IsFinish == define.StepIsFinish {
			continue
		}
		processIsFinishNumber := 0
		for key2, step := range process.Steps {
			if step.Key == stepType {
				guideData.ProcessList[key].Steps[key2].IsFinish = define.StepIsFinish
			}
			if guideData.ProcessList[key].Steps[key2].IsFinish == define.StepIsFinish {
				processIsFinishNumber++
			}
		}
		if processIsFinishNumber == len(process.Steps) {
			guideData.ProcessList[key].IsFinish = define.StepIsFinish
		}
	}
	_, err = msql.Model(define.TableUseGuideProcess, define.Postgres).
		Where(`use_guide_type`, define.UseTypeGuide).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Update(msql.Datas{`data`: tool.JsonEncodeNoError(guideData)})
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	lib_redis.DelCacheData(define.Redis, UseGuideCacheBuildHandler{
		AdminUserId:  adminUserId,
		UseGuideType: define.UseTypeGuide,
	})
	return nil
}

type UseGuideCacheBuildHandler struct {
	UseGuideType string
	AdminUserId  int
}

func (h UseGuideCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.use.guide.%d.%s`, h.AdminUserId, h.UseGuideType)
}
func (h UseGuideCacheBuildHandler) GetCacheData() (any, error) {
	return getGuideData(h.AdminUserId, h.UseGuideType)
}
