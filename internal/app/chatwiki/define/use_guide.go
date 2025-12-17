// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

const UseTypeGuide = `default`

const (
	StepIsFinish    = `1`
	StepIsNotFinish = `0`
)

const (
	ProcessTypeSetModel      = `set_model`
	ProcessTypeCreateRobot   = `create_robot`
	ProcessTypeCreateLibrary = `create_library`
	ProcessTypeTestRobot     = `test_robot`
)

const (
	StepTestRobot       = `test_robot`
	StepCreateRobot     = `create_robot`
	StepRelationLibrary = `relation_library`
	StepCreateLibrary   = `create_library`
	StepImportWord      = `import_word`
	StepImportPdf       = `import_pdf`
	StepSetLlm          = `set_llm`
	StepSetEmbedding    = `set_embedding`
)

type GuideProcess struct {
	Name     string      `json:"name"`
	Key      string      `json:"type"`
	IsFinish string      `json:"is_finish"` //1 finish
	Steps    []GuideStep `json:"steps"`
}

type GuideStep struct {
	Name     string `json:"name"`
	Key      string `json:"type"`
	IsFinish string `json:"is_finish"`
}

type GuideData struct {
	ProcessList []GuideProcess `json:"process_list"`
}
