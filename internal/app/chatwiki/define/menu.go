package define

type Menu struct {
	Name     string  `json:"name"`
	UniKey   string  `json:"uni_key"`
	Path     string  `json:"path"`
	Children []*Menu `json:"children"`
}

var Menus = []Menu{
	// {
	// 	Name:   "发现",
	// 	UniKey: "Discovery",
	// 	Children: []*Menu{
	// 		{
	// 			Name:   "发现",
	// 			UniKey: "DiscoveryManage",
	// 		},
	// 	},
	// },
	{
		Name:   "机器人管理",
		UniKey: "Robot",
		Children: []*Menu{
			{
				Name:   "机器人管理",
				UniKey: "RobotManage",
				Children: []*Menu{
					{
						Name:   "创建机器人",
						UniKey: "RobotCreate",
					},
				},
			},
		},
	},
	{
		Name:   "知识库",
		UniKey: "library",
		Children: []*Menu{
			{
				Name:   "知识库",
				UniKey: "LibraryManage",
				Children: []*Menu{
					{
						Name:   "创建知识库",
						UniKey: "LibraryCreate",
					},
				},
			},
			{
				Name:   "数据库",
				UniKey: "FormManage",
				Children: []*Menu{
					{
						Name:   "创建数据库",
						UniKey: "FormCreate",
					},
				},
			},
			// {
			// 	Name:   "文档提取FAQ",
			// 	UniKey: "DocFaq",
			// 	Children: []*Menu{
			// 		{
			// 			Name:   "上传文档提取",
			// 			UniKey: "UploadDocFaq",
			// 		},
			// 	},
			// },
		},
	},
	{
		Name:   "文档",
		UniKey: "OpenLibDoc",
		Children: []*Menu{
			{
				Name:   "对外文档",
				UniKey: "OpenLibDocManage",
				Children: []*Menu{
					{
						Name:   "新建对外文档",
						UniKey: "CreateOpenLibDoc",
					},
				},
			},
		},
	},
	{
		Name:   "搜索",
		UniKey: "Search",
		Children: []*Menu{
			{
				Name:   "搜索",
				UniKey: "SearchManage",
				Children: []*Menu{
					{
						Name:   "搜索设置",
						UniKey: "SearchSets",
					},
				},
			},
		},
	},
	{
		Name:   "会话",
		UniKey: "ChatSession",
		Children: []*Menu{
			{
				Name:   "会话",
				UniKey: "ChatSessionManage",
			},
		},
	},
	{
		Name:   "系统管理",
		UniKey: "System",
		Children: []*Menu{
			{
				Name:   "模型管理",
				UniKey: "ModelManage",
			},
			{
				Name:   "Token使用",
				UniKey: "TokenManage",
			},
			{
				Name:   "团队管理",
				UniKey: "TeamManage",
			},
			{
				Name:   "自定义域名",
				UniKey: "UserDomainManage",
			},
			// {
			// 	Name:   "版本信息",
			// 	UniKey: "VersionManage",
			// },
			{
				Name:   "客户端下载",
				UniKey: "ClientSideManage",
			},
			{
				Name:   "阿里云OCR",
				UniKey: "AliyunOCRManage",
			},
			{
				Name:   "敏感词管理",
				UniKey: "SensitiveWordManage",
			},
			{
				Name:   "提示词模板库",
				UniKey: "PromptTemplateManage",
			},
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
