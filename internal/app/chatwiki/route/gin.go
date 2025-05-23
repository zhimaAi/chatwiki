// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package route

import (
	"net/http"

	"chatwiki/internal/app/chatwiki/business"
	"chatwiki/internal/app/chatwiki/business/manage"
	"chatwiki/internal/pkg/lib_web"
)

var Route lib_web.Route

func init() {
	//step1:initialize
	Route = make(map[string]map[string]lib_web.Action)
	Route[http.MethodGet] = make(map[string]lib_web.Action)
	Route[http.MethodPost] = make(map[string]lib_web.Action)
	Route[lib_web.NoMethod] = make(map[string]lib_web.Action)
	Route[lib_web.NoRoute] = make(map[string]lib_web.Action)
	//step2:define route
	noAuthFuns(Route[http.MethodGet], `/ping`, business.Ping)
	Route[lib_web.NoMethod][`/`] = business.NoMethod //NoMethod
	Route[lib_web.NoRoute][`/`] = business.NoRoute   //NoMethod
	noAuthFuns(Route[http.MethodGet], `/test/test`, business.Test)
	noAuthFuns(Route[http.MethodGet], `/test/domain`, business.TestDomain)
	Route[http.MethodGet][`/test/test1`] = business.Test1

	/* user API*/
	noAuthFuns(Route[http.MethodPost], `/manage/saveProfile`, manage.SaveProfile)
	Route[http.MethodPost]["/manage/refreshUserToken"] = manage.RefreshUserToken
	Route[http.MethodPost]["/manage/loginSwitch"] = manage.LoginSwitch

	/*admin API*/
	noAuthFuns(Route[http.MethodPost], `/manage/login`, manage.AdminLogin)
	Route[http.MethodGet][`/manage/checkLogin`] = manage.CheckLogin
	/* permission API*/
	noAuthFuns(Route[http.MethodGet], "/manage/getMenu", manage.GetMenu)
	noAuthFuns(Route[http.MethodGet], "/manage/checkPermission", manage.CheckPermission)
	Route[http.MethodPost]["/manage/saveMenu"] = manage.SaveMenu
	Route[http.MethodPost]["/manage/delMenu"] = manage.DelMenu
	Route[http.MethodGet]["/manage/getUserList"] = manage.GetUserList
	Route[http.MethodPost]["/manage/saveUser"] = manage.SaveUser
	Route[http.MethodPost]["/manage/resetPass"] = manage.ResetPass
	Route[http.MethodPost]["/manage/delUser"] = manage.DeleteUser
	Route[http.MethodGet]["/manage/getUser"] = manage.GetUser
	Route[http.MethodPost]["/manage/saveUserManagedDataList"] = manage.SaveUserManagedDataList
	Route[http.MethodGet]["/manage/getRoleList"] = manage.GetRoleList
	Route[http.MethodPost]["/manage/saveRole"] = manage.SaveRole
	Route[http.MethodPost]["/manage/delRole"] = manage.DelRole
	Route[http.MethodGet]["/manage/getRole"] = manage.GetRole
	/*company API*/
	noAuthFuns(Route[http.MethodGet], "/manage/getCompany", manage.GetCompany)
	Route[http.MethodPost]["/manage/saveCompany"] = manage.SaveCompany
	Route[http.MethodPost]["/manage/saveAliOcr"] = manage.SaveAliOcr
	Route[http.MethodPost][`/manage/checkAliOcr`] = manage.CheckAliOcr

	/*robot API*/
	Route[http.MethodPost][`/manage/upload`] = manage.Upload
	Route[http.MethodGet][`/manage/getRobotList`] = manage.GetRobotList
	Route[http.MethodPost][`/manage/saveRobot`] = manage.SaveRobot
	Route[http.MethodPost][`/manage/addFlowRobot`] = manage.AddFlowRobot
	Route[http.MethodPost][`/manage/editExternalConfig`] = manage.EditExternalConfig
	Route[http.MethodGet][`/manage/getRobotInfo`] = manage.GetRobotInfo
	Route[http.MethodPost][`/manage/deleteRobot`] = manage.DeleteRobot
	Route[http.MethodGet][`/manage/createPromptByAi`] = manage.CreatePromptByAi
	Route[http.MethodPost][`/manage/robotCopy`] = manage.RobotCopy
	Route[http.MethodPost][`/manage/editBaseInfo`] = manage.EditBaseInfo
	Route[http.MethodPost][`/manage/relationWorkFlow`] = manage.RelationWorkFlow

	/*library API*/
	Route[http.MethodGet][`/manage/getLibraryList`] = manage.GetLibraryList
	Route[http.MethodGet][`/manage/getLibraryInfo`] = manage.GetLibraryInfo
	Route[http.MethodPost][`/manage/createLibrary`] = manage.CreateLibrary
	Route[http.MethodPost][`/manage/deleteLibrary`] = manage.DeleteLibrary
	Route[http.MethodPost][`/manage/editLibrary`] = manage.EditLibrary
	Route[http.MethodGet][`/manage/getLibraryRobotInfo`] = manage.GetLibraryRobotInfo
	Route[http.MethodPost][`/manage/relationRobot`] = manage.RelationRobot

	/** library search API*/
	Route[http.MethodPost][`/manage/libraryAiSummary`] = manage.LibraryAiSummary
	Route[http.MethodGet][`/manage/getLibrarySearch`] = manage.GetLibrarySearch
	Route[http.MethodPost][`/manage/saveLibrarySearch`] = manage.SaveLibrarySearch

	/*open library API*/
	noAuthFuns(Route[http.MethodGet], `/manage/getCatalog`, manage.GetCatalog)
	noAuthFuns(Route[http.MethodGet], `/manage/getLibDocInfo`, manage.GetLibDocInfo)
	noAuthFuns(Route[http.MethodGet], `/manage/libDocSearch`, manage.LibDocSearch)
	noAuthFuns(Route[http.MethodGet], `/manage/questionGuideList`, manage.QuestionGuideList)
	Route[http.MethodPost][`/manage/saveLibDoc`] = manage.SaveLibDoc
	Route[http.MethodPost][`/manage/changeLibDoc`] = manage.ChangeLibDoc
	Route[http.MethodPost][`/manage/draftLibDoc`] = manage.DraftLibDoc
	Route[http.MethodPost][`/manage/uploadLibDoc`] = manage.UploadLibDoc
	Route[http.MethodGet][`/manage/exportLibDoc`] = manage.ExportLibDoc
	Route[http.MethodPost][`/manage/deleteLibDoc`] = manage.DeleteLibDoc
	Route[http.MethodPost][`/manage/saveQuestionGuide`] = manage.SaveQuestionGuide
	Route[http.MethodPost][`/manage/deleteQuestionGuide`] = manage.DeleteQuestionGuide
	Route[http.MethodPost][`/manage/saveLibDocSeo`] = manage.SaveLibDocSeo
	Route[http.MethodPost][`/manage/saveLibDocPartner`] = manage.SaveLibDocPartner
	Route[http.MethodGet][`/manage/getLibDocPartner`] = manage.GetLibDocPartner
	Route[http.MethodGet][`/manage/libDocPartnerList`] = manage.LibDocPartnerList
	Route[http.MethodPost][`/manage/deleteLibDocPartner`] = manage.DeleteLibDocPartner
	noAuthFuns(Route[http.MethodGet], `/manage/libDocHome/:key`, business.OpenHome)
	noAuthFuns(Route[http.MethodGet], `/open/doc/:key`, business.OpenDoc)
	noAuthFuns(Route[http.MethodGet], `/open/home/:key`, business.OpenHome)
	noAuthFuns(Route[http.MethodGet], `/open/search/:type/:lib_key`, business.OpenSearch)
	noAuthFuns(Route[http.MethodGet], `/open/summary/:lib_key`, business.OpenAiSummary)
	/*libFile API*/
	Route[http.MethodGet][`/manage/getLibFileList`] = manage.GetLibFileList
	Route[http.MethodPost][`/manage/addLibraryFile`] = manage.AddLibraryFile
	Route[http.MethodPost][`/manage/delLibraryFile`] = manage.DelLibraryFile
	Route[http.MethodGet][`/manage/getLibFileInfo`] = manage.GetLibFileInfo
	Route[http.MethodGet][`/manage/getLibRawFile`] = manage.GetLibRawFile
	noAuthFuns(Route[http.MethodGet], `/manage/getLibRawFileOnePage`, manage.GetLibRawFileOnePage)
	Route[http.MethodGet][`/manage/getLibFileExcelTitle`] = manage.GetLibFileExcelTitle
	Route[http.MethodPost][`/manage/readLibFileExcelTitle`] = manage.ReadLibFileExcelTitle
	Route[http.MethodPost][`/manage/restudyLibraryFile`] = manage.RestudyLibraryFile
	Route[http.MethodPost][`/manage/renewLibraryFile`] = manage.RenewLibraryFile
	Route[http.MethodPost][`/manage/editLibFile`] = manage.EditLibFile
	Route[http.MethodPost][`/manage/manualCrawl`] = manage.ManualCrawl
	Route[http.MethodPost][`/manage/constructGraph`] = manage.ConstructGraph
	Route[http.MethodGet][`/manage/getFileGraphInfo`] = manage.GetFileGraphInfo
	Route[http.MethodPost][`/manage/reconstructVector`] = manage.ReconstructVector
	Route[http.MethodPost][`/manage/reconstructCategoryVector`] = manage.ReconstructCategoryVector
	Route[http.MethodPost][`/manage/reconstructGraph`] = manage.ReconstructGraph
	Route[http.MethodPost][`/manage/cancelOcrPdf`] = manage.CancelOcrPdf
	/*paragraph API*/
	Route[http.MethodGet][`/manage/getSeparatorsList`] = manage.GetSeparatorsList
	Route[http.MethodGet][`/manage/getLibFileSplit`] = manage.GetLibFileSplit
	Route[http.MethodGet][`/manage/getLibFileSplitAiChunks`] = manage.GetLibFileSplitAiChunks
	Route[http.MethodPost][`/manage/saveLibFileSplit`] = manage.SaveLibFileSplit
	Route[http.MethodGet][`/manage/getParagraphList`] = manage.GetParagraphList
	Route[http.MethodGet][`/manage/getParagraphCount`] = manage.GetParagraphCount
	Route[http.MethodGet][`/manage/getCategoryParagraphList`] = manage.GetCategoryParagraphList
	Route[http.MethodPost][`/manage/saveCategoryParagraph`] = manage.SaveCategoryParagraph
	Route[http.MethodPost][`/manage/addParagraph`] = manage.SaveParagraph
	Route[http.MethodPost][`/manage/editParagraph`] = manage.SaveParagraph
	Route[http.MethodGet][`/manage/getSplitParagraph`] = manage.GetSplitParagraph
	Route[http.MethodPost][`/manage/saveSplitParagraph`] = manage.SaveSplitParagraph
	Route[http.MethodPost][`/manage/deleteParagraph`] = manage.DeleteParagraph
	Route[http.MethodPost][`/manage/updateParagraphCategory`] = manage.UpdateParagraphCategory
	Route[http.MethodPost][`/manage/generateSimilarQuestions`] = manage.GenerateSimilarQuestions
	/*category API*/
	Route[http.MethodGet][`/manage/getCategoryList`] = manage.GetCategoryList
	Route[http.MethodPost][`/manage/saveCategory`] = manage.SaveCategory

	/*form API*/
	Route[http.MethodGet][`/manage/getFormList`] = manage.GetFormList
	Route[http.MethodGet][`/manage/getFormInfo`] = manage.GetFormInfo
	Route[http.MethodPost][`/manage/addForm`] = manage.SaveForm
	Route[http.MethodPost][`/manage/editForm`] = manage.SaveForm
	Route[http.MethodPost][`/manage/delForm`] = manage.DelForm
	/*form field API*/
	Route[http.MethodGet][`/manage/getFormFieldList`] = manage.GetFormFieldList
	Route[http.MethodPost][`/manage/addFormField`] = manage.SaveFormField
	Route[http.MethodPost][`/manage/editFormField`] = manage.SaveFormField
	Route[http.MethodPost][`/manage/updateFormRequired`] = manage.UpdateFormRequired
	Route[http.MethodPost][`/manage/delFormField`] = manage.DelFormField
	Route[http.MethodPost][`/manage/uploadFormFile`] = manage.UploadFormFile
	Route[http.MethodPost][`/manage/getUploadFormFileProc`] = manage.GetUploadFormFileProc
	/*form entry API*/
	Route[http.MethodGet][`/manage/getFormEntryList`] = manage.GetFormEntryList
	Route[http.MethodPost][`/manage/addFormEntry`] = manage.SaveFormEntry
	Route[http.MethodPost][`/manage/editFormEntry`] = manage.SaveFormEntry
	Route[http.MethodPost][`/manage/delFormEntry`] = manage.DelFormEntry
	Route[http.MethodPost][`/manage/emptyFormEntry`] = manage.EmptyFormEntry
	Route[http.MethodGet][`/manage/exportFormEntry`] = manage.ExportFormEntry
	/*form filter API*/
	Route[http.MethodGet][`/manage/getFormFilterList`] = manage.GetFormFilterList
	Route[http.MethodGet][`/manage/getFormFilterInfo`] = manage.GetFormFilterInfo
	Route[http.MethodPost][`/manage/addFormFilter`] = manage.SaveFormFilter
	Route[http.MethodPost][`/manage/editFormFilter`] = manage.SaveFormFilter
	Route[http.MethodPost][`/manage/updateFormFilterEnabled`] = manage.UpdateFormFilterEnabled
	Route[http.MethodPost][`/manage/updateFormFilterSort`] = manage.UpdateFormFilterSort
	Route[http.MethodPost][`/manage/delFormFilter`] = manage.DelFormFilter
	/*stats*/
	Route[http.MethodGet][`/manage/stats/getActiveModels`] = manage.GetActiveModels
	Route[http.MethodGet][`/manage/stats/token`] = manage.StatToken
	Route[http.MethodGet][`/manage/stats/analyse`] = manage.StatAnalyse
	/*debug API*/
	Route[http.MethodPost][`/manage/getDialogueList`] = manage.GetDialogueList
	Route[http.MethodPost][`/manage/libraryRecallTest`] = manage.LibraryRecallTest
	noAuthFuns(Route[http.MethodGet], `/manage/getAnswerSource`, manage.GetAnswerSource)
	/*chat API*/
	noAuthFuns(Route[http.MethodGet], `/chat/getWsUrl`, business.GetWsUrl)
	noAuthFuns(Route[http.MethodGet], `/chat/isOnLine`, business.IsOnLine)
	noAuthFuns(Route[http.MethodPost], `/chat/message`, business.ChatMessage)
	noAuthFuns(Route[http.MethodPost], `/chat/message/addFeedback`, business.AddChatMessageFeedback)
	noAuthFuns(Route[http.MethodPost], `/chat/message/delFeedback`, business.DelChatMessageFeedback)
	noAuthFuns(Route[http.MethodPost], `/chat/welcome`, business.ChatWelcome)
	noAuthFuns(Route[http.MethodPost], `/chat/request`, business.ChatRequest)
	noAuthFuns(Route[http.MethodPost], `/chat/checkChatRequestPermission`, business.CheckChatRequestPermission)

	noAuthFuns(Route[http.MethodPost], `/chat/questionGuide`, business.ChatQuestionGuide)
	/*model API*/
	Route[http.MethodGet][`/manage/getModelConfigList`] = manage.GetModelConfigList
	Route[http.MethodPost][`/manage/addModelConfig`] = manage.AddModelConfig
	Route[http.MethodPost][`/manage/delModelConfig`] = manage.DelModelConfig
	Route[http.MethodPost][`/manage/editModelConfig`] = manage.EditModelConfig
	Route[http.MethodGet][`/manage/getModelConfigOption`] = manage.GetModelConfigOption
	/*Fast Command API*/
	Route[http.MethodGet][`/manage/getFastCommandList`] = manage.GetFastCommandList
	Route[http.MethodPost][`/manage/saveFastCommand`] = manage.SaveFastCommand
	Route[http.MethodPost][`/manage/sortFastCommand`] = manage.SortFastCommand
	Route[http.MethodGet][`/manage/GetFastCommandInfo`] = manage.GetFastCommandInfo
	Route[http.MethodPost][`/manage/deleteFastCommand`] = manage.DeleteFastCommand
	Route[http.MethodPost][`/manage/updateFastCommandSwitch`] = manage.UpdateFastCommandSwitch
	noAuthFuns(Route[http.MethodGet], `/chat/getFastCommandList`, business.GetFastCommandList)

	//register client side route
	RegClientSideRoute()
	/*session API*/
	Route[http.MethodGet][`/manage/getSessionChannelList`] = manage.GetSessionChannelList
	Route[http.MethodGet][`/manage/getSessionRecordList`] = manage.GetSessionRecordList
	Route[http.MethodGet][`/manage/createSessionExport`] = manage.CreateSessionExport
	/*feedback API*/
	Route[http.MethodGet][`/manage/feedback/stats`] = manage.StatMessageFeedback
	Route[http.MethodGet][`/manage/feedback/list`] = manage.GetMessageFeedbackList
	Route[http.MethodGet][`/manage/feedback/detail`] = manage.GetMessageFeedbackDetail
	/*export API*/
	Route[http.MethodGet][`/manage/getExportTaskList`] = manage.GetExportTaskList
	Route[http.MethodGet][`/manage/downloadExportFile`] = manage.DownloadExportFile
	/* diy domain API*/
	Route[http.MethodPost][`/manage/saveDiyDomain`] = manage.SaveDiyDomain
	Route[http.MethodGet][`/manage/diyDomainList`] = manage.DiyDomainList
	Route[http.MethodPost][`/manage/deleteDiyDomain`] = manage.DeleteDiyDomain
	Route[http.MethodPost][`/manage/uploadCertificate`] = manage.UploadCertificate
	Route[http.MethodPost][`/manage/uploadCheckFile`] = manage.UploadCheckFile
	/*work_flow API*/
	Route[http.MethodGet][`/manage/getNodeList`] = manage.GetNodeList
	Route[http.MethodPost][`/manage/saveNodes`] = manage.SaveNodes
	/* sensitive words API*/
	Route[http.MethodGet][`/manage/getSensitiveWordsList`] = manage.GetSensitiveWordsList
	Route[http.MethodPost][`/manage/saveSensitiveWords`] = manage.SaveSensitiveWords
	Route[http.MethodPost][`/manage/switchSensitiveWords`] = manage.SwitchSensitiveWords
	Route[http.MethodPost][`/manage/deleteSensitiveWords`] = manage.DeleteSensitiveWords
	/*receiver API*/
	Route[http.MethodGet][`/manage/getReceiverList`] = manage.GetReceiverList
	Route[http.MethodPost][`/manage/setReceiverRead`] = manage.SetReceiverRead
	/*permission manage API*/
	Route[http.MethodGet][`/manage/getDepartmentList`] = manage.GetDepartmentList
	Route[http.MethodGet][`/manage/getAllDepartment`] = manage.GetAllDepartment
	Route[http.MethodPost][`/manage/saveDepartment`] = manage.SaveDepartment
	Route[http.MethodPost][`/manage/deleteDepartment`] = manage.DeleteDepartment
	Route[http.MethodPost][`/manage/batchUpdateUserDepartment`] = manage.BatchUpdateUserDepartment
	Route[http.MethodGet][`/manage/getPermissionManageList`] = manage.GetPermissionManageList
	Route[http.MethodGet][`/manage/getPartnerManageList`] = manage.GetPartnerManageList
	Route[http.MethodPost][`/manage/batchSavePermissionManage`] = manage.BatchSavePermissionManage
	Route[http.MethodPost][`/manage/savePermissionManage`] = manage.SavePermissionManage
	Route[http.MethodPost][`/manage/deletePermissionManage`] = manage.DeletePermissionMange
	/*prompt_library API*/
	Route[http.MethodGet][`/manage/getPromptLibraryGroup`] = manage.GetPromptLibraryGroup
	Route[http.MethodPost][`/manage/savePromptLibraryGroup`] = manage.SavePromptLibraryGroup
	Route[http.MethodPost][`/manage/deletePromptLibraryGroup`] = manage.DeletePromptLibraryGroup
	Route[http.MethodGet][`/manage/getPromptLibraryItems`] = manage.GetPromptLibraryItems
	Route[http.MethodPost][`/manage/savePromptLibraryItems`] = manage.SavePromptLibraryItems
	Route[http.MethodPost][`/manage/deletePromptLibraryItems`] = manage.DeletePromptLibraryItems
	Route[http.MethodGet][`/manage/createPromptByLlm`] = manage.CreatePromptByLlm
	Route[http.MethodPost][`/manage/movePromptLibraryItems`] = manage.MovePromptLibraryItems
}

func noAuthFuns(route map[string]lib_web.Action, path string, handlerFunc lib_web.Action) map[string]lib_web.Action {
	lib_web.NoAuthRouteMap[path] = true
	route[path] = handlerFunc
	return route
}
