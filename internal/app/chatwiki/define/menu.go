package define

type Menu struct {
	Name     string  `json:"name"`
	UniKey   string  `json:"uni_key"`
	Path     string  `json:"path"`
	Children []*Menu `json:"children"`
}

var Menus = []Menu{
	{
		Name:   "机器人管理",
		UniKey: "Robot",
		Children: []*Menu{
			{
				Name:   "创建机器人",
				UniKey: "RobotCreate",
			},
			{
				Name:   "机器人管理",
				UniKey: "RobotManage",
			},
		},
	},
	{
		Name:   "知识库管理",
		UniKey: "library",
		Children: []*Menu{
			{
				Name:   "创建知识库",
				UniKey: "LibraryCreate",
			},
			{
				Name:   "知识库管理",
				UniKey: "LibraryManage",
			},
		},
	},
	{
		Name:   "数据库",
		UniKey: "Form",
		Children: []*Menu{
			{
				Name:   "创建数据库",
				UniKey: "FormCreate",
			},
			{
				Name:   "数据库管理",
				UniKey: "FormManage",
			},
		},
	},
	{
		Name:   "搜索",
		UniKey: "Search",
		Children: []*Menu{
			{
				Name:   "知识库搜索",
				UniKey: "LibrarySearch",
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
				Name:   "账号设置",
				UniKey: "AccountManage",
			},
			{
				Name:   "企业设置",
				UniKey: "CompanyManage",
			},
			{
				Name:   "客户端下载",
				UniKey: "ClientSideManage",
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
