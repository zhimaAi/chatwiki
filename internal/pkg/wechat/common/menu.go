// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import "github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"

type Button struct {
	Type       string    `json:"type,omitempty"`
	Name       string    `json:"name"`
	Key        string    `json:"key"`
	URL        string    `json:"url,omitempty"`
	AppID      string    `json:"appid,omitempty"`
	PagePath   string    `json:"pagepath,omitempty"`
	SubButtons []*Button `json:"sub_button"`
}

type Menu struct {
	Buttons []*Button `json:"button"`
	MenuID  int       `json:"menuid"`
}

type MatchRule struct {
	GroupID            int    `json:"group_id"`
	Sex                int    `json:"sex"`
	Country            string `json:"country"`
	Province           string `json:"province"`
	City               string `json:"city"`
	ClientPlatformType string `json:"client_platform_type"`
}

type ConditionalMenu struct {
	Buttons   []*Button  `json:"button"`
	MatchRule *MatchRule `json:"matchrule"`
	MenuID    int        `json:"menuid"`
}

type ResponseMenuGet struct {
	response.ResponseOfficialAccount

	Menus            *Menu              `json:"menu"`
	ConditionalMenus []*ConditionalMenu `json:"conditionalmenu"`
}
