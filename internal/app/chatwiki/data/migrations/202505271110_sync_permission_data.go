package migrations

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func init() {
	goose.AddMigrationNoTxContext(run, nil)
}

func run(_ context.Context, db *sql.DB) error {
	SyncHistoryPermissionData()
	return nil
}

func SyncHistoryPermissionData() {
	syncLibraryCreatorData()
	page := 1
	size := 200
	exec := make(chan int, 0)

	go func() {
		total := 0
		for {
			data, _, err := msql.Model(define.TableUser, define.Postgres).Where("is_deleted", define.Normal).Order("id asc").Paginate(page, size)
			if err != nil {
				logs.Error(err.Error())
				return
			}
			if len(data) <= 0 {
				exec <- total
				return
			}
			logs.Info("同步历史权限数据,  查询:%v 条", len(data))
			for _, item := range data {
				adminUserId := cast.ToInt(item[`parent_id`])
				if adminUserId == 0 {
					adminUserId = cast.ToInt(item[`id`])
				}
				if cast.ToInt(item[`parent_id`]) == 0 {
					defaultDepartment, _ := common.GetDefaultDepartmentInfo(adminUserId)
					if len(defaultDepartment) <= 0 {
						common.SaveDepartment(0, cast.ToInt64(adminUserId), msql.Datas{
							`department_name`: `默认部门`,
							`is_default`:      1,
							`admin_user_id`:   adminUserId,
						})
					}
					for _, typ := range []int{define.ObjectTypeRobot, define.ObjectTypeLibrary, define.ObjectTypeForm} {
						hasData, err := msql.Model(`permission_manage`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
							Where(`identity_id`, cast.ToString(adminUserId)).
							Where(`identity_type`, cast.ToString(define.IdentityTypeUser)).
							Where(`object_id`, cast.ToString(define.ObjectTypeAll)).
							Where(`object_type`, cast.ToString(typ)).Value(`id`)
						if err != nil {
							logs.Error(err.Error())
							continue
						}
						if len(hasData) > 0 {
							continue
						}
						common.SavePermissionManage(0, cast.ToInt64(adminUserId), msql.Datas{
							`admin_user_id`:  adminUserId,
							`identity_id`:    adminUserId,
							`object_id`:      define.ObjectTypeAll,
							`operate_rights`: define.PermissionManageRights,
							`creator`:        adminUserId,
							`identity_type`:  define.IdentityTypeUser,
							`object_type`:    typ,
						})
					}
				} else {
					for _, id := range strings.Split(item[`managed_robot_list`], `,`) {
						if cast.ToInt(id) != 0 {
							data, _ := msql.Model(`chat_ai_robot`, define.Postgres).Where(`id`, cast.ToString(id)).Value(`id`)
							if len(data) == 0 {
								continue
							}
							hasData, err := msql.Model(`permission_manage`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
								Where(`identity_id`, cast.ToString(item[`id`])).
								Where(`identity_type`, cast.ToString(define.IdentityTypeUser)).
								Where(`object_id`, cast.ToString(id)).
								Where(`object_type`, cast.ToString(define.ObjectTypeRobot)).Value(`id`)
							if err != nil {
								logs.Error(err.Error())
								continue
							}
							if len(hasData) > 0 {
								continue
							}
							common.SavePermissionManage(0, cast.ToInt64(adminUserId), msql.Datas{
								`admin_user_id`:  adminUserId,
								`identity_id`:    item[`id`],
								`object_id`:      id,
								`operate_rights`: define.PermissionManageRights,
								`creator`:        adminUserId,
								`identity_type`:  define.IdentityTypeUser,
								`object_type`:    define.ObjectTypeRobot,
							})
						}
					}
					for _, id := range strings.Split(item[`managed_library_list`], `,`) {
						if cast.ToInt(id) != 0 {
							data, _ := msql.Model(`chat_ai_library`, define.Postgres).Where(`id`, cast.ToString(id)).Value(`id`)
							if len(data) == 0 {
								continue
							}
							hasData, err := msql.Model(`permission_manage`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
								Where(`identity_id`, cast.ToString(item[`id`])).
								Where(`identity_type`, cast.ToString(define.IdentityTypeUser)).
								Where(`object_id`, cast.ToString(id)).
								Where(`object_type`, cast.ToString(define.ObjectTypeLibrary)).Value(`id`)
							if err != nil {
								logs.Error(err.Error())
								continue
							}
							if len(hasData) > 0 {
								continue
							}
							common.SavePermissionManage(0, cast.ToInt64(adminUserId), msql.Datas{
								`admin_user_id`:  adminUserId,
								`identity_id`:    item[`id`],
								`object_id`:      id,
								`operate_rights`: define.PermissionManageRights,
								`creator`:        adminUserId,
								`identity_type`:  define.IdentityTypeUser,
								`object_type`:    define.ObjectTypeLibrary,
							})
						}
					}
					for _, id := range strings.Split(item[`managed_form_list`], `,`) {
						if cast.ToInt(id) != 0 {
							data, _ := msql.Model(`form`, define.Postgres).Where(`id`, cast.ToString(id)).Value(`id`)
							if len(data) == 0 {
								continue
							}
							hasData, err := msql.Model(`permission_manage`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
								Where(`identity_id`, cast.ToString(item[`id`])).
								Where(`identity_type`, cast.ToString(define.IdentityTypeUser)).
								Where(`object_id`, cast.ToString(id)).
								Where(`object_type`, cast.ToString(define.ObjectTypeForm)).Value(`id`)
							if err != nil {
								logs.Error(err.Error())
								continue
							}
							if len(hasData) > 0 {
								continue
							}
							common.SavePermissionManage(0, cast.ToInt64(adminUserId), msql.Datas{
								`admin_user_id`:  adminUserId,
								`identity_id`:    item[`id`],
								`object_id`:      id,
								`operate_rights`: define.PermissionManageRights,
								`creator`:        adminUserId,
								`identity_type`:  define.IdentityTypeUser,
								`object_type`:    define.ObjectTypeForm,
							})
						}
					}
				}
				total++
			}
			page++
		}
	}()
	select {
	case data := <-exec:
		logs.Info("同步历史权限数据完成, 共同步 %v 条", data)
		return
	}
}

func syncLibraryCreatorData() {
	list, err := msql.Model(`chat_ai_library`, define.Postgres).Where(`type`, `in`, fmt.Sprintf(`%v,%v`, cast.ToString(define.GeneralLibraryType), cast.ToString(define.QALibraryType))).Field(`id,admin_user_id,creator`).Select()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	for _, item := range list {
		user := item[`creator`]
		if cast.ToInt(user) == 0 {
			user = item[`admin_user_id`]
		}
		adminUserId := cast.ToInt(item[`admin_user_id`])
		hasData, err := msql.Model(`permission_manage`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`identity_id`, cast.ToString(user)).
			Where(`identity_type`, cast.ToString(define.IdentityTypeUser)).
			Where(`object_id`, item[`id`]).
			Where(`object_type`, cast.ToString(define.ObjectTypeLibrary)).Value(`id`)
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		if len(hasData) > 0 {
			continue
		}
		common.SavePermissionManage(0, cast.ToInt64(adminUserId), msql.Datas{
			`admin_user_id`:  adminUserId,
			`identity_id`:    user,
			`object_id`:      item[`id`],
			`operate_rights`: define.PermissionManageRights,
			`creator`:        adminUserId,
			`identity_type`:  define.IdentityTypeUser,
			`object_type`:    define.ObjectTypeLibrary,
		})
	}
}
