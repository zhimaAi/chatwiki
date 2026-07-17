// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package route

import (
	"chatwiki/internal/app/chatwiki/business/manage"
	"net/http"
)

func RegBookToSkillRoute() {
	Route[http.MethodPost][`/manage/bookToSkill/createTask`] = manage.CreateBookToSkillTask
	Route[http.MethodGet][`/manage/bookToSkill/taskList`] = manage.GetBookToSkillTaskList
	Route[http.MethodGet][`/manage/bookToSkill/taskProgress`] = manage.GetBookToSkillTaskProgress
	Route[http.MethodPost][`/manage/bookToSkill/stopTask`] = manage.StopBookToSkillTask
	Route[http.MethodPost][`/manage/bookToSkill/retryTask`] = manage.RetryBookToSkillTask
	Route[http.MethodGet][`/manage/bookToSkill/taskLog`] = manage.GetBookToSkillTaskLog
	Route[http.MethodPost][`/manage/bookToSkill/installSkill`] = manage.InstallBookToSkill
	Route[http.MethodGet][`/manage/bookToSkill/downloadSkill`] = manage.DownloadBookToSkill
}
