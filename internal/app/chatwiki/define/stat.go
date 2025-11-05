// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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
