// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package route

import (
	"chatwiki/internal/app/chatwiki/business"
	"chatwiki/internal/app/chatwiki/business/manage"
	"chatwiki/internal/pkg/lib_web"
	"net/http"
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
	Route[http.MethodGet][`/test/test1`] = business.Test1

	/* user API*/
	noAuthFuns(Route[http.MethodPost], `/manage/saveProfile`, manage.SaveProfile)
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
	Route[http.MethodGet]["/manage/getRoleList"] = manage.GetRoleList
	Route[http.MethodPost]["/manage/saveRole"] = manage.SaveRole
	Route[http.MethodPost]["/manage/delRole"] = manage.DelRole
	Route[http.MethodGet]["/manage/getRole"] = manage.GetRole
	/*company API*/
	noAuthFuns(Route[http.MethodGet], "/manage/getCompany", manage.GetCompany)
	Route[http.MethodPost]["/manage/saveCompany"] = manage.SaveCompany

	/*robot API*/
	Route[http.MethodPost][`/manage/upload`] = manage.Upload
	Route[http.MethodGet][`/manage/getRobotList`] = manage.GetRobotList
	Route[http.MethodPost][`/manage/saveRobot`] = manage.SaveRobot
	Route[http.MethodPost][`/manage/editExternalConfig`] = manage.EditExternalConfig
	Route[http.MethodGet][`/manage/getRobotInfo`] = manage.GetRobotInfo
	Route[http.MethodPost][`/manage/deleteRobot`] = manage.DeleteRobot
	/*library API*/
	Route[http.MethodGet][`/manage/getLibraryList`] = manage.GetLibraryList
	Route[http.MethodGet][`/manage/getLibraryInfo`] = manage.GetLibraryInfo
	Route[http.MethodPost][`/manage/createLibrary`] = manage.CreateLibrary
	Route[http.MethodPost][`/manage/deleteLibrary`] = manage.DeleteLibrary
	Route[http.MethodPost][`/manage/editLibrary`] = manage.EditLibrary
	/*libFile API*/
	Route[http.MethodGet][`/manage/getLibFileList`] = manage.GetLibFileList
	Route[http.MethodPost][`/manage/addLibraryFile`] = manage.AddLibraryFile
	Route[http.MethodPost][`/manage/delLibraryFile`] = manage.DelLibraryFile
	Route[http.MethodGet][`/manage/getLibFileInfo`] = manage.GetLibFileInfo
	Route[http.MethodGet][`/manage/getLibFileExcelTitle`] = manage.GetLibFileExcelTitle
	Route[http.MethodPost][`/manage/renewLibraryFile`] = manage.RenewLibraryFile
	Route[http.MethodPost][`/manage/editLibFile`] = manage.EditLibFile
	/*paragraph API*/
	Route[http.MethodGet][`/manage/getSeparatorsList`] = manage.GetSeparatorsList
	Route[http.MethodGet][`/manage/getLibFileSplit`] = manage.GetLibFileSplit
	Route[http.MethodPost][`/manage/saveLibFileSplit`] = manage.SaveLibFileSplit
	Route[http.MethodGet][`/manage/getParagraphList`] = manage.GetParagraphList
	Route[http.MethodPost][`/manage/addParagraph`] = manage.SaveParagraph
	Route[http.MethodPost][`/manage/editParagraph`] = manage.SaveParagraph
	Route[http.MethodPost][`/manage/deleteParagraph`] = manage.DeleteParagraph
	/*debug API*/
	Route[http.MethodPost][`/manage/getDialogueList`] = manage.GetDialogueList
	Route[http.MethodPost][`/manage/libraryRecallTest`] = manage.LibraryRecallTest
	Route[http.MethodGet][`/manage/getAnswerSource`] = manage.GetAnswerSource
	/*chat API*/
	noAuthFuns(Route[http.MethodGet], `/chat/getWsUrl`, business.GetWsUrl)
	noAuthFuns(Route[http.MethodGet], `/chat/isOnLine`, business.IsOnLine)
	noAuthFuns(Route[http.MethodPost], `/chat/message`, business.ChatMessage)
	noAuthFuns(Route[http.MethodPost], `/chat/welcome`, business.ChatWelcome)
	noAuthFuns(Route[http.MethodPost], `/chat/request`, business.ChatRequest)
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
	noAuthFuns(Route[http.MethodGet], `/chat/getFastCommandList`, business.GetFastCommandList)
	//register client side route
	RegClientSideRoute()
}

func noAuthFuns(route map[string]lib_web.Action, path string, handlerFunc lib_web.Action) map[string]lib_web.Action {
	lib_web.NoAuthRouteMap[path] = true
	route[path] = handlerFunc
	return route
}
