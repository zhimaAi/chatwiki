// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package lib_define

const (
	SmartMenuTypeNormal = "0" //普通菜单
	SmartMenuTypeKey    = "1" //关键词菜单
)

// SmartMenu 智能菜单
type SmartMenu struct {
	ID              int                `json:"id"`
	AdminUserID     int                `json:"admin_user_id"`
	RobotID         int                `json:"robot_id"`
	MenuTitle       string             `json:"menu_title"`
	MenuDescription string             `json:"menu_description"`
	MenuContent     []SmartMenuContent `json:"menu_content"`
	CreateTime      int                `json:"create_time"`
	UpdateTime      int                `json:"update_time"`
}

// SmartMenuContent 智能菜单内容
type SmartMenuContent struct {
	ID       string `json:"id" form:"id"`
	MenuType string `json:"menu_type" form:"menu_type"` //菜单类型 0：普通 1：点击关键词菜单
	SerialNo string `json:"serial_no" form:"serial_no"` //序号
	Content  string `json:"content" form:"content"`     //内容  如果是点击关键词菜单，则是关键词
	RuleID   string `json:"rule_id" form:"rule_id"`     // 关键词规则id
}
