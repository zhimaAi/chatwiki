// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

const (
	TableGoodsLibGroup   = `chat_ai_goods_group`
	TableGoodsLibLibrary = `chat_ai_goods_library`
)

const (
	GoodsLibMaxGroupLevel       = 5
	GoodsLibGroupNameMaxLength  = 15
	GoodsLibBaseInfoMaxLength   = 100
	GoodsLibLinkMaxLength       = 1000
	GoodsLibDetailMaxLength     = 1000
	GoodsLibImageLimitSize      = 10 * 1024 * 1024
	GoodsLibImportFileLimitSize = 100 * 1024 * 1024
	GoodsLibDefaultPageSize     = 20
)

const (
	GoodsLibSwitchOff = 0
	GoodsLibSwitchOn  = 1
)

const (
	GoodsLibRecommendSearchTypeDetail = `detail`
	GoodsLibRecommendSearchTypeBasic  = `basic`
)

var GoodsLibImportAllowExt = []string{`xlsx`, `csv`}

type GoodsLibGroup struct {
	ID              int64            `json:"id"`
	ParentID        int64            `json:"parent_id"`
	GroupName       string           `json:"group_name"`
	Level           int              `json:"level"`
	Sort            int              `json:"sort"`
	GoodsCount      int              `json:"goods_count"`
	TotalGoodsCount int              `json:"total_goods_count"`
	Children        []*GoodsLibGroup `json:"children"`
}

type GoodsLibSaveGroupParams struct {
	ID        int64
	ParentID  int64
	GroupName string
}

type GoodsLibGroupSortItem struct {
	ID   int64 `form:"id" json:"id" binding:"required"`
	Sort int   `form:"sort" json:"sort"`
}

type GoodsLibListFilter struct {
	GroupID      int64  `form:"group_id" json:"group_id"`
	GroupIDs     string `form:"group_ids" json:"group_ids"`
	Keyword      string `form:"keyword" json:"keyword"`
	SwitchStatus int    `form:"switch_status" json:"switch_status"`
	Page         int    `form:"page" json:"page"`
	Size         int    `form:"size" json:"size"`
}

type GoodsLibSaveParams struct {
	ID           int64    `form:"id" json:"id" binding:"gte=0"`
	GroupID      int64    `form:"group_id" json:"group_id" binding:"gte=0"`
	GoodsID      string   `form:"goods_id" json:"goods_id" binding:"required"`
	GoodsName    string   `form:"goods_name" json:"goods_name" binding:"required"`
	Category     string   `form:"category" json:"category"`
	Brand        string   `form:"brand" json:"brand"`
	Price        float64  `form:"price" json:"price" binding:"gte=0"`
	Stock        int64    `form:"stock" json:"stock"`
	Link         string   `form:"link" json:"link"`
	Images       []string `form:"images" json:"images"`
	Description  string   `form:"description" json:"description"`
	QA           string   `form:"qa" json:"qa"`
	CustomInfo   string   `form:"custom_info" json:"custom_info"`
	SwitchStatus *int     `form:"switch_status" json:"switch_status"`
}

type GoodsLibImportHeader struct {
	Field string
	Name  string
}

type GoodsLibImportError struct {
	Row int `json:"row"`
	GoodsLibSaveParams
	Message string `json:"message"`
}

type GoodsLibImportResult struct {
	TotalCount   int                    `json:"total_count"`
	CreatedCount int                    `json:"created_count"`
	UpdatedCount int                    `json:"updated_count"`
	FailedCount  int                    `json:"failed_count"`
	Errors       []GoodsLibImportError  `json:"errors"`
	Headers      []GoodsLibImportHeader `json:"headers"`
}
