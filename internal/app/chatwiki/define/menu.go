// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

type Menu struct {
	Name     string  `json:"name"`
	UniKey   string  `json:"uni_key"`
	Path     string  `json:"path"`
	Children []*Menu `json:"children"`
}

var Menus = []Menu{
	// {
	// 	Name:   "[[ZM--DiscoveryName--ZM]]",
	// 	UniKey: "Discovery",
	// 	Children: []*Menu{
	// 		{
	// 			Name:   "[[ZM--DiscoveryName--ZM]]",
	// 			UniKey: "DiscoveryManage",
	// 		},
	// 	},
	// },
	{
		Name:   "[[ZM--AbilityName--ZM]]",
		UniKey: "Ability",
		Children: []*Menu{
			{
				Name:   "[[ZM--AbilityName--ZM]]",
				UniKey: "AbilityCenter",
			},
		},
	},
	{
		Name:   "[[ZM--RobotManageName--ZM]]",
		UniKey: "Robot",
		Children: []*Menu{
			{
				Name:   "[[ZM--RobotManageName--ZM]]",
				UniKey: "RobotManage",
				Children: []*Menu{
					{
						Name:   "[[ZM--RobotCreateName--ZM]]",
						UniKey: "RobotCreate",
					},
				},
			},
		},
	},
	{
		Name:   "[[ZM--KnowledgeBaseName--ZM]]",
		UniKey: "library",
		Children: []*Menu{
			{
				Name:   "[[ZM--KnowledgeBaseName--ZM]]",
				UniKey: "LibraryManage",
				Children: []*Menu{
					{
						Name:   "[[ZM--LibraryCreateName--ZM]]",
						UniKey: "LibraryCreate",
					},
				},
			},
			{
				Name:   "[[ZM--DatabaseName--ZM]]",
				UniKey: "FormManage",
				Children: []*Menu{
					{
						Name:   "[[ZM--FormCreateName--ZM]]",
						UniKey: "FormCreate",
					},
				},
			},
			{
				Name:   "[[ZM--DocFaqName--ZM]]",
				UniKey: "DocFaq",
				Children: []*Menu{
					{
						Name:   "[[ZM--UploadDocFaqName--ZM]]",
						UniKey: "UploadDocFaq",
					},
				},
			},
		},
	},
	{
		Name:   "[[ZM--DocumentName--ZM]]",
		UniKey: "OpenLibDoc",
		Children: []*Menu{
			{
				Name:   "[[ZM--OpenLibDocName--ZM]]",
				UniKey: "OpenLibDocManage",
				Children: []*Menu{
					{
						Name:   "[[ZM--CreateOpenLibDocName--ZM]]",
						UniKey: "CreateOpenLibDoc",
					},
				},
			},
		},
	},
	{
		Name:   "[[ZM--SearchName--ZM]]",
		UniKey: "Search",
		Children: []*Menu{
			{
				Name:   "[[ZM--SearchName--ZM]]",
				UniKey: "SearchManage",
				Children: []*Menu{
					{
						Name:   "[[ZM--SearchSetsName--ZM]]",
						UniKey: "SearchSets",
					},
				},
			},
		},
	},
	{
		Name:   "[[ZM--ChatSessionName--ZM]]",
		UniKey: "ChatSession",
		Children: []*Menu{
			{
				Name:   "[[ZM--ChatSessionName--ZM]]",
				UniKey: "ChatSessionManage",
			},
		},
	},
	{
		Name:   "[[ZM--SystemManageName--ZM]]",
		UniKey: "System",
		Children: []*Menu{
			{
				Name:   "[[ZM--ModelManageName--ZM]]",
				UniKey: "ModelManage",
			},
			{
				Name:   "[[ZM--TokenManageName--ZM]]",
				UniKey: "TokenManage",
			},
			{
				Name:   "[[ZM--TeamManageName--ZM]]",
				UniKey: "TeamManage",
			},
			{
				Name:   "[[ZM--UserDomainManageName--ZM]]",
				UniKey: "UserDomainManage",
			},
			// {
			// 	Name:   "[[ZM--VersionManageName--ZM]]",
			// 	UniKey: "VersionManage",
			// },
			{
				Name:   "[[ZM--ClientSideManageName--ZM]]",
				UniKey: "ClientSideManage",
			},
			{
				Name:   "[[ZM--AliyunOCRManageName--ZM]]",
				UniKey: "AliyunOCRManage",
			},
			{
				Name:   "[[ZM--SensitiveWordManageName--ZM]]",
				UniKey: "SensitiveWordManage",
			},
			{
				Name:   "[[ZM--PromptTemplateManageName--ZM]]",
				UniKey: "PromptTemplateManage",
			},
			//{
			//	Name:   "[[ZM--LogInfoName--ZM]]",
			//	UniKey: "Logs",
			//},
			//{
			//	Name:   "[[ZM--OfficialAccountManageName--ZM]]",
			//	UniKey: "OfficialAccountManage",
			//},
		},
	},
}

var AllUniKeyList = []string{
	"RobotCreate",
	"RobotManage",
	"LibraryCreate",
	"LibraryManage",
	"FormCreate",
	"FormManage",
	"ModelManage",
	"TokenManage",
	"TeamManage",
	"AccountManage",
	"CompanyManage",
	"ClientSideManage",
	"LibrarySearch",
}

var UserUniKeyList = []string{
	"RobotManage",
	"LibraryManage",
	"FormManage",
	"AccountManage",
	"ClientSideManage",
}

var MustUniKeyList = []string{
	"RobotManage",
}
var ContainsUniKeyList = []map[string]string{
	{
		"RobotCreate": "RobotManage",
	},
	{
		"LibraryCreate": "LibraryManage",
	},
	{
		"FormCreate": "FormManage",
	},
}

func GetAllUniKeyList() []string {
	var uniKeys []string
	for _, item := range Menus {
		uniKeys = append(uniKeys, RescurseUniKeyList(item.Children)...)
	}
	return uniKeys
}

func RescurseUniKeyList(menu []*Menu) (uniKeys []string) {
	for _, item := range menu {
		if item.UniKey != "" {
			uniKeys = append(uniKeys, item.UniKey)
		}
		if item.Children != nil {
			uniKeys = append(uniKeys, RescurseUniKeyList(item.Children)...)
		}
	}
	return uniKeys
}
