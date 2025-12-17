// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

const (
	TokenAppTypeRobot    string = `chatwiki_robot`
	TokenAppTypeWorkflow string = `workflow`
	TokenAppTypeOther    string = `other`
)

func GetTokenAppTypes() []string {
	return []string{
		TokenAppTypeRobot,
		TokenAppTypeWorkflow,
		TokenAppTypeOther,
	}
}
