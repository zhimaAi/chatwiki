// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_web"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/tool"
)

func GetGoodsGroupList(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	data, err := common.GetGoodsLibGroupList(common.GetLang(c), adminUserId)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, data)
}

type saveGoodsGroupReq struct {
	ID        int64  `form:"id" json:"id" binding:"gte=0"`
	ParentID  int64  `form:"parent_id" json:"parent_id" binding:"gte=0"`
	GroupName string `form:"group_name" json:"group_name" binding:"required"`
}

func SaveGoodsGroup(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	req := saveGoodsGroupReq{}
	if err := common.RequestParamsBind(&req, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if utf8.RuneCountInString(strings.TrimSpace(req.GroupName)) > define.GoodsLibGroupNameMaxLength {
		common.FmtError(c, `param_invalid`, `group_name`)
		return
	}
	id, err := common.SaveGoodsLibGroup(common.GetLang(c), adminUserId, define.GoodsLibSaveGroupParams{
		ID:        req.ID,
		ParentID:  req.ParentID,
		GroupName: req.GroupName,
	})
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, id)
}

type deleteGoodsGroupReq struct {
	ID int64 `form:"id" json:"id" binding:"required"`
}

func DeleteGoodsGroup(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	req := deleteGoodsGroupReq{}
	if err := common.RequestParamsBind(&req, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if err := common.DeleteGoodsLibGroup(common.GetLang(c), adminUserId, req.ID); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, nil)
}

type sortGoodsGroupReq struct {
	Data []define.GoodsLibGroupSortItem `form:"data" json:"data" binding:"required,min=1,dive"`
}

func SortGoodsGroup(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	req := sortGoodsGroupReq{}
	err := common.RequestParamsBind(&req, c)
	if len(req.Data) == 0 {
		data := strings.TrimSpace(c.PostForm(`data`))
		if len(data) > 0 {
			err = tool.JsonDecodeUseNumber(data, &req.Data)
		}
	}
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if err := common.SortGoodsLibGroup(common.GetLang(c), adminUserId, req.Data); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, nil)
}

func GetGoodsLibraryList(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	filter := define.GoodsLibListFilter{
		GroupID:      -1,
		SwitchStatus: -1,
		Page:         1,
		Size:         define.GoodsLibDefaultPageSize,
	}
	if err := c.ShouldBindQuery(&filter); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(filter, err, common.GetLang(c)).Error())
		return
	}
	if !common.ValidateGoodsLibraryListReq(filter) {
		common.FmtError(c, `param_invalid`, `query`)
		return
	}
	list, total, err := common.GetGoodsLibLibraryList(common.GetLang(c), adminUserId, filter)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, map[string]any{
		`list`: list, `total`: total, `page`: filter.Page, `size`: filter.Size,
	})
}

type getGoodsLibraryInfoReq struct {
	ID int64 `form:"id" json:"id" binding:"required"`
}

func GetGoodsLibraryInfo(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	req := getGoodsLibraryInfoReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	info, err := common.GetGoodsLibLibraryInfo(common.GetLang(c), adminUserId, req.ID)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, info)
}

func SaveGoodsLibrary(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	params := define.GoodsLibSaveParams{}
	if err := common.RequestParamsBind(&params, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(params, err, common.GetLang(c)).Error())
		return
	}
	params.Images = append(params.Images, c.PostFormArray(`images[]`)...)
	if len(params.Images) == 1 && strings.HasPrefix(strings.TrimSpace(params.Images[0]), `[`) {
		imagesJson := params.Images[0]
		params.Images = nil
		if err := tool.JsonDecodeUseNumber(imagesJson, &params.Images); err != nil {
			common.FmtError(c, `param_invalid`, i18n.Show(common.GetLang(c), `goods_import_header_images`))
			return
		}
	} else if len(params.Images) == 0 {
		imagesJson := strings.TrimSpace(c.PostForm(`images`))
		if strings.HasPrefix(imagesJson, `[`) {
			if err := tool.JsonDecodeUseNumber(imagesJson, &params.Images); err != nil {
				common.FmtError(c, `param_invalid`, i18n.Show(common.GetLang(c), `goods_import_header_images`))
				return
			}
		}
	}
	id, err := common.SaveGoodsLibLibrary(common.GetLang(c), adminUserId, params)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, id)
}

type deleteGoodsLibraryReq struct {
	ID int64 `form:"id" json:"id" binding:"required"`
}

func DeleteGoodsLibrary(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	req := deleteGoodsLibraryReq{}
	if err := common.RequestParamsBind(&req, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if err := common.DeleteGoodsLibLibrary(common.GetLang(c), adminUserId, req.ID); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, nil)
}

type updateGoodsLibrarySwitchReq struct {
	ID           int64 `form:"id" json:"id" binding:"required"`
	SwitchStatus *int  `form:"switch_status" json:"switch_status" binding:"required"`
}

func UpdateGoodsLibrarySwitch(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	req := updateGoodsLibrarySwitchReq{}
	if err := common.RequestParamsBind(&req, c); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if *req.SwitchStatus != define.GoodsLibSwitchOff && *req.SwitchStatus != define.GoodsLibSwitchOn {
		common.FmtError(c, `param_invalid`, `switch_status`)
		return
	}
	if err := common.UpdateGoodsLibLibrarySwitch(common.GetLang(c), adminUserId, req.ID, *req.SwitchStatus); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, nil)
}

func UploadGoodsLibraryImage(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	fileHeader, err := c.FormFile(`file`)
	if err != nil || fileHeader == nil || fileHeader.Size <= 0 ||
		fileHeader.Size > define.GoodsLibImageLimitSize {
		common.FmtError(c, `param_invalid`, `file`)
		return
	}
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fileHeader.Filename), `.`))
	if !tool.InArrayString(ext, define.ImageAllowExt) {
		common.FmtError(c, `param_invalid`, `file`)
		return
	}
	uploadInfo, err := common.SaveGoodsLibImage(common.GetLang(c), adminUserId, fileHeader)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, uploadInfo)
}

func ImportGoodsLibrary(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	groupId := cast.ToInt64(c.PostForm(`group_id`))
	fileHeader, err := c.FormFile(`file`)
	if err != nil || fileHeader == nil || fileHeader.Size <= 0 ||
		fileHeader.Size > define.GoodsLibImportFileLimitSize {
		common.FmtError(c, `param_invalid`, `file`)
		return
	}
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fileHeader.Filename), `.`))
	if !tool.InArrayString(ext, define.GoodsLibImportAllowExt) {
		common.FmtError(c, `param_invalid`, `file`)
		return
	}
	result, err := common.ImportGoodsLibLibrary(common.GetLang(c), adminUserId, groupId, fileHeader)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	common.FmtOk(c, result)
}

func ExportGoodsLibrary(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	filter := define.GoodsLibListFilter{
		GroupID:      -1,
		SwitchStatus: -1,
	}
	if err := c.ShouldBindQuery(&filter); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(filter, err, common.GetLang(c)).Error())
		return
	}
	if !common.ValidateGoodsLibraryFilter(filter) {
		common.FmtError(c, `param_invalid`, `query`)
		return
	}
	filePath, fileName, err := common.ExportGoodsLibLibrary(common.GetLang(c), adminUserId, filter, false)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	sendGoodsLibFile(c, filePath, fileName)
}

func DownloadGoodsLibraryImportTemplate(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}
	filePath, fileName, err := common.ExportGoodsLibLibrary(common.GetLang(c), adminUserId, define.GoodsLibListFilter{}, true)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	sendGoodsLibFile(c, filePath, fileName)
}

func sendGoodsLibFile(c *gin.Context, filePath, fileName string) {
	c.FileAttachment(filePath, fileName)
	go func(path string) {
		time.Sleep(time.Minute)
		_ = os.Remove(path)
	}(filePath)
}
