// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

const VectorDimension = 2000

const MaxRobotNum = 6

const (
	FileStatusWaitCrawl      = 5
	FileStatusCrawling       = 6
	FileStatusCrawlException = 7
	FileStatusInitial        = 0
	FileStatusException      = 3
	FileStatusPartException  = 8
	FileStatusWaitSplit      = 4
	FileStatusLearning       = 1
	FileStatusLearned        = 2
	FileStatusCancelled      = 9
	FileStatusChunking       = 10
)

const (
	QAIndexTypeQuestionAndAnswer = 1
	QAIndexTypeQuestion          = 2
)

const (
	VectorStatusInitial    = 0
	VectorStatusConverted  = 1
	VectorStatusException  = 2
	VectorStatusConverting = 3
	SplitStatusException   = 4
)

const (
	GraphStatusNotStart        = 0
	GraphStatusInitial         = 1
	GraphStatusWorking         = 4
	GraphStatusPartlyConverted = 5
	GraphStatusConverted       = 2
	GraphStatusException       = 3
)

const (
	MsgFromCustomer = 1
	MsgFromRobot    = 0
)

const (
	MsgTypeMixed = 99 //多模态输入内容
	MsgTypeText  = 1
	MsgTypeMenu  = 2
	MsgTypeImage = 3
)

const (
	FileIsTable = 1
	DocTypeQa   = 1
)

const (
	ChunkTypeNormal    = 1 //普通分段
	ChunkTypeSemantic  = 2 //语义分段
	ChunkTypeAi        = 3 //AI分段
	ChunkTypeFatherSon = 4 //父子分段
)

const (
	FatherChunkParagraphTypeFullText = 1
	FatherChunkParagraphTypeSection  = 2
)

const (
	ParagraphTypeNormal  = 1
	ParagraphTypeDocQA   = 2
	ParagraphTypeExcelQA = 3
)

const (
	SplitChunkMaxSize        = 10000
	SplitChunkMinSize        = 200
	SplitAiChunkMaxSize      = 5000
	SplitPreviewChunkMaxSize = 5000
)

const (
	VectorTypeParagraph       = 1
	VectorTypeQuestion        = 2
	VectorTypeAnswer          = 3
	VectorTypeCustom          = 4
	VectorTypeSimilarQuestion = 5
)

const (
	SearchTypeMixed    = 1
	SearchTypeVector   = 2
	SearchTypeFullText = 3
	SearchTypeGraph    = 4
)

var SeparatorsList = []map[string]any{
	{`no`: 1, `name`: `#`, `code`: []string{"\r\n# ", "\n# ", "\r# "}},
	{`no`: 2, `name`: `##`, `code`: []string{"\r\n## ", "\n## ", "\r## "}},
	{`no`: 3, `name`: `###`, `code`: []string{"\r\n### ", "\n### ", "\r### "}},
	{`no`: 4, `name`: `####`, `code`: []string{"\r\n#### ", "\n#### ", "\r#### "}},
	{`no`: 5, `name`: `#####`, `code`: []string{"\r\n##### ", "\n##### ", "\r##### "}},
	{`no`: 6, `name`: `-`, `code`: `-`},
	{`no`: 7, `name`: `space`, `code`: " "},
	{`no`: 8, `name`: `semicolon`, `code`: []string{`；`, `;`}},
	{`no`: 9, `name`: `comma`, `code`: []string{`，`, `,`}},
	{`no`: 10, `name`: `period`, `code`: []string{`。`, `.`}},
	{`no`: 11, `name`: `enter`, `code`: []string{"\r\n", "\n", "\r"}},
	{`no`: 12, `name`: `blank_line`, `code`: []string{"\r\n\r\n", "\n\n", "\r\r"}},
	{`no`: 13, `name`: `tab`, `code`: "\t"},
}

var DefaultUserRoleId int

const (
	DefaultUser   = `admin`
	DefaultPasswd = `chatwiki.com@123`
	UserTypeAdmin = 1
)

const (
	ChatTypeLibrary = 1 //仅知识库模式
	ChatTypeDirect  = 2 //直连模式
	ChatTypeMixture = 3 //混合模式
)

const (
	DocTypeLocal    = 1
	DocTypeOnline   = 2
	DocTypeCustom   = 3
	DocTypeDiy      = 4
	DocTypeOfficial = 5
)

const (
	SwitchOff = 0
	SwitchOn  = 1
)

const DefaultCustomerAvatar = `/public/user_avatar_2x.png`

const (
	DefaultCustomDomain   = `http://cloud.chatwiki.com`
	DefaultCustomH5Domain = `http://h5.wikichat.com.cn`
)

const (
	LibDocIndex = 1
	IsDraft     = 1
	IsPub       = 1
)

const (
	PartnerRightsManage = 4
	PartnerRightsEdit   = 2
)

const (
	AccessRestrictionsTypeLogin = 2
	AccessPermissionTypeLogin   = 3
)

const (
	FAQFileStatusQueuing       = 0
	FAQFileStatusAnalyzing     = 1
	FAQFileStatusExtracting    = 2
	FAQFileStatusExtracted     = 3
	FAQFileStatusExtractFailed = 4
)
const (
	FAQFileSplitStatusSuccess = 1
	FAQFileSplitStatusFailed  = 2
)
const (
	FAQChunkTypeLength       = 1
	FAQChunkTypeSeparatorsNo = 2

	FAQExtractTypeAI = 1
)

const (
	LibraryGroupTypeQA   = 0
	LibraryGroupTypeFile = 1
)

// chat_ai_library_file_doc 中的is_dir字段
const (
	IsFile = 0
	ISDir  = 1
)
