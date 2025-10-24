// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package chatwiki

import (
	"chatwiki/internal/app/chatwiki/business"
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/data/migrations"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/initialize"
	"chatwiki/internal/pkg/casbin"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_web"
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pressly/goose/v3"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/mq"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func Run() {
	//initialize
	initialize.Initialize()
	//postgres table
	PostgresTable()
	//producer handle
	define.ProducerHandle = mq.NewProducerHandle().SetWorkNum(3).SetHostAndPort(define.Config.Nsqd[`host`], cast.ToUint(define.Config.Nsqd[`port`]))
	//consumer handle
	define.ConsumerHandle = mq.NewConsumerHandle().SetHostAndPort(define.Config.NsqLookup[`host`], cast.ToUint(define.Config.NsqLookup[`port`]))
	//web start
	go lib_web.WebRun(define.WebService)
	//pprof api
	go func() {
		err := http.ListenAndServe(":55557", nil)
		if err != nil {
			logs.Error(err.Error())
		}
	}()
	//consumer start
	StartConsumer()

	//cron tasks
	StartCronTasks()
}

func Stop() {
	define.ConsumerHandle.Stop()
	lib_web.Shutdown(define.WebService)
	define.ProducerHandle.Stop()
}

func StartConsumer() {
	common.RunTask(define.ConvertHtmlTopic, define.ConvertHtmlChannel, 1, business.ConvertHtml)
	common.RunTask(define.ConvertVectorTopic, define.ConvertVectorChannel, 2, business.ConvertVector)
	common.RunTask(define.ConvertGraphTopic, define.ConvertGraphChannel, 10, business.ConvertGraph)
	common.RunTask(define.CrawlArticleTopic, define.CrawlArticleChannel, 2, business.CrawlArticle)
	common.RunTask(lib_define.PushMessage, lib_define.PushChannel, 10, business.AppPush)
	common.RunTask(lib_define.PushEvent, lib_define.PushChannel, 5, business.AppPush)
	common.RunTask(define.ExportTaskTopic, define.ExportTaskChannel, 5, business.ExportTask)
	common.RunTask(define.ExtractFaqFilesTopic, define.ExtractFaqFilesChannel, 5, business.ExtractFaqFiles)
}

func StartCronTasks() {
	c := cron.New()
	_, _ = c.AddFunc("@every 1m", func() { logs.Debug("cron test") })
	_, _ = c.AddFunc("@every 1m", func() { business.RenewCrawl() })
	_, _ = c.AddFunc("@every 1h", func() { business.DeleteFormEntry() })
	_, _ = c.AddFunc("@every 1h", func() { business.DeleteExportFile() })
	_, _ = c.AddFunc("@every 1h", func() { business.DeleteConvertHtml() })
	_, _ = c.AddFunc("@every 1h", func() { business.DeleteClientSide() })
	_, _ = c.AddFunc("@every 1h", func() { business.DeleteDownloadFile() })
	_, _ = c.AddFunc("@every 15s", common.DeleteReceiver)
	_, _ = c.AddFunc("@every 5s", func() { business.CheckAliOcrRequest() })
	_, _ = c.AddFunc("0 0 * * *", func() { business.UpdateLibraryFileData() })
	c.Start()
	logs.Debug("cron start")
}

//go:embed data/migrations/*.sql
var embedMigrations embed.FS

func PostgresTable() {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		define.Config.Postgres["host"], define.Config.Postgres["port"],
		define.Config.Postgres["user"], define.Config.Postgres["password"],
		define.Config.Postgres["dbname"], define.Config.Postgres["sslmode"])

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "data/migrations", goose.WithAllowMissing()); err != nil {
		panic(err)
	}

	InitDefaultRole()
	userId, insert, err := CreateDefaultUser()
	if err != nil {
		logs.Error(err.Error())
	}
	if insert {
		CreateDefaultBaaiModel(userId)
	}
	common.SetNeo4jStatus(cast.ToInt(userId), cast.ToBool(define.Config.Neo4j["enabled"]))
	InitRoleRootPermissions()
	InitRoleUserPermissions()
}

func CreateDefaultUser() (int64, bool, error) {
	m := msql.Model(define.TableUser, define.Postgres)
	user, err := m.Where(`"user_name"`, define.DefaultUser).Find()
	if err != nil {
		panic(err)
	}
	if len(user) > 0 {
		return cast.ToInt64(user["id"]), false, nil
	}

	salt := tool.Random(20)
	id, err := msql.Model(define.TableUser, define.Postgres).Insert(msql.Datas{
		`user_name`:   define.DefaultUser,
		`salt`:        salt,
		`password`:    tool.MD5(define.DefaultPasswd + salt),
		`user_type`:   define.UserTypeAdmin,
		`user_roles`:  define.DefaultUserRoleId,
		`create_time`: tool.Time2Int(),
		`update_time`: tool.Time2Int(),
	}, "id")
	if err != nil {
		panic(err)
	}
	migrations.SyncHistoryPermissionData()
	return id, true, nil
}
func CreateDefaultRole(userId int64) {
	var defaultRole = []string{define.DefaultRoleRoot, define.DefaultRoleAdmin, define.DefaultRoleUser}
	for k, role := range defaultRole {
		_, err := msql.Model(define.TableRole, define.Postgres).Insert(msql.Datas{
			`name`:        role,
			"role_type":   k + 1,
			`create_name`: "系统",
			`create_time`: tool.Time2Int(),
			`update_time`: tool.Time2Int(),
		})
		if err != nil {
			logs.Error(`role create err:%s`, err.Error())
		}
		if userId <= 0 || role != define.DefaultRoleRoot {
			continue
		}
	}
}

func InitDefaultRole() {
	rolesMap := map[int]string{
		define.RoleTypeRoot:  define.DefaultRoleRoot,
		define.RoleTypeAdmin: define.DefaultRoleAdmin,
		define.RoleTypeUser:  define.DefaultRoleUser,
	}
	for roleType, roleName := range rolesMap {
		row, err := msql.Model(define.TableRole, define.Postgres).Where(`role_type`, cast.ToString(roleType)).Find()
		if err != nil {
			panic(err.Error())
		}
		if len(row) > 0 {
			if roleType == define.RoleTypeRoot {
				define.DefaultUserRoleId = cast.ToInt(row[`id`])
			}
			continue
		}
		id, err := msql.Model(define.TableRole, define.Postgres).Insert(msql.Datas{
			`name`:        roleName,
			`role_type`:   roleType,
			`create_name`: "系统",
			`create_time`: tool.Time2Int(),
			`update_time`: tool.Time2Int(),
		}, `id`)
		if roleType == define.RoleTypeRoot {
			define.DefaultUserRoleId = cast.ToInt(id)
		}
	}
}

func CreateDefaultBaaiModel(userId int64) {
	modelInfo, ok := common.GetModelInfoByDefine(common.ModelBaai)
	if !ok {
		logs.Error(`modelInfo not found`)
		return
	}
	_, err := msql.Model("chat_ai_model_config", define.Postgres).Insert(msql.Datas{
		`admin_user_id`:   userId,
		`model_define`:    common.ModelBaai,
		`model_types`:     strings.Join(modelInfo.SupportedType, `,`),
		`api_endpoint`:    "http://host.docker.internal:50001",
		`deployment_name`: "",
		`create_time`:     tool.Time2Int(),
		`update_time`:     tool.Time2Int(),
	})
	if err != nil {
		logs.Error("baai model create err:%s", err.Error())
	}
}

func InitRoleRootPermissions() {
	roles, err := msql.Model(define.TableRole, define.Postgres).
		Where(`role_type`, `in`, fmt.Sprintf(`%d,%d`, define.RoleTypeRoot, define.RoleTypeAdmin)).
		Field(`id,name,role_type`).
		Select()
	if err != nil {
		panic(err.Error())
	}

	for _, role := range roles {

		// reset role name
		if cast.ToInt(role[`role_type`]) == define.RoleTypeRoot && role[`name`] != define.DefaultRoleRoot {
			_, err = msql.Model(define.TableRole, define.Postgres).Where(`id`, role[`id`]).Update(msql.Datas{`name`: define.DefaultRoleRoot})
			if err != nil {
				panic(err.Error())
			}
		}
		if cast.ToInt(role[`role_type`]) == define.RoleTypeAdmin && role[`name`] != define.DefaultRoleAdmin {
			_, err = msql.Model(define.TableRole, define.Postgres).Where(`id`, role[`id`]).Update(msql.Datas{`name`: define.DefaultRoleAdmin})
			if err != nil {
				panic(err.Error())
			}
		}

		// reset role permissions
		_, err = casbin.Handler.DelRoleRules(role[`id`])
		for _, item := range define.GetAllUniKeyList() {
			_, err := casbin.Handler.AddPolicies([][]string{{role[`id`], item, "GET"}})
			if err != nil {
				panic(err.Error())
			}
		}
	}
}

func InitRoleUserPermissions() {
	roleInfo, err := msql.Model(define.TableRole, define.Postgres).Where(`role_type`, cast.ToString(define.RoleTypeUser)).Field(`id,name`).Find()
	if err != nil {
		panic(err.Error())
	}
	if len(roleInfo) == 0 {
		panic(`no user role`)
	}
	if roleInfo[`name`] != define.DefaultRoleUser {
		_, err = msql.Model(define.TableRole, define.Postgres).Where(`id`, roleInfo[`id`]).Update(msql.Datas{`name`: define.DefaultRoleUser})
		if err != nil {
			panic(err.Error())
		}
	}
	rolePermissions := make([][]string, 0)
	_, err = casbin.Handler.DelRoleRules(roleInfo[`id`])
	for _, item := range define.UserUniKeyList {
		rolePermissions = append(rolePermissions, []string{roleInfo[`id`], item, "GET"})
		_, err := casbin.Handler.AddPolicies([][]string{{roleInfo[`id`], item, "GET"}})
		if err != nil {
			panic(err.Error())
		}
	}
}
