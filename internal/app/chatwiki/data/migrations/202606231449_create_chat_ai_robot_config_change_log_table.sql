-- +goose Up

CREATE TABLE "public"."chat_ai_robot_config_change_log"
(
    "id"               bigserial    NOT NULL primary key,
    "admin_user_id"    int4         NOT NULL DEFAULT 0,
    "robot_id"         int4         NOT NULL DEFAULT 0,
    "robot_key"        varchar(64)  NOT NULL DEFAULT '',
    "application_type" int2         NOT NULL DEFAULT 0,
    "module"           varchar(50)  NOT NULL DEFAULT '',
    "section"          varchar(50)  NOT NULL DEFAULT '',
    "oper_user_id"     int4         NOT NULL DEFAULT 0,
    "oper_user_name"   varchar(100) NOT NULL DEFAULT '',
    "before_content"   jsonb        NOT NULL DEFAULT '{}',
    "after_content"    jsonb        NOT NULL DEFAULT '{}',
    "change_detail"    jsonb        NOT NULL DEFAULT '[]',
    "create_time"      int4         NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."chat_ai_robot_config_change_log" ("admin_user_id", "create_time");
CREATE INDEX ON "public"."chat_ai_robot_config_change_log" ("robot_id");

COMMENT ON TABLE  "public"."chat_ai_robot_config_change_log" IS '机器人/工作流配置变更日志';
COMMENT ON COLUMN "public"."chat_ai_robot_config_change_log"."id"               IS '自增ID';
COMMENT ON COLUMN "public"."chat_ai_robot_config_change_log"."admin_user_id"    IS '企业账号(管理员用户ID)';
COMMENT ON COLUMN "public"."chat_ai_robot_config_change_log"."robot_id"         IS '机器人/工作流ID';
COMMENT ON COLUMN "public"."chat_ai_robot_config_change_log"."robot_key"        IS '机器人/工作流key(冗余,便于检索)';
COMMENT ON COLUMN "public"."chat_ai_robot_config_change_log"."application_type" IS '应用类型 0=聊天 1=工作流 2=Claw';
COMMENT ON COLUMN "public"."chat_ai_robot_config_change_log"."module"           IS '修改模块 base_info/relation_library/relation_workflow/robot_setting/unknown_summary/external_config/lang_config/clean_cache/chat_variable';
COMMENT ON COLUMN "public"."chat_ai_robot_config_change_log"."section"          IS '修改所属界面板块(按部署语言,冗余,便于检索/肉眼分组)';
COMMENT ON COLUMN "public"."chat_ai_robot_config_change_log"."oper_user_id"     IS '操作人(当前登录用户)ID';
COMMENT ON COLUMN "public"."chat_ai_robot_config_change_log"."oper_user_name"   IS '操作人用户名(冗余)';
COMMENT ON COLUMN "public"."chat_ai_robot_config_change_log"."before_content"   IS '变化字段的修改前内容(JSON,仅含本次变化字段,原始值)';
COMMENT ON COLUMN "public"."chat_ai_robot_config_change_log"."after_content"    IS '变化字段的修改后内容(JSON,仅含本次变化字段,原始值)';
COMMENT ON COLUMN "public"."chat_ai_robot_config_change_log"."change_detail"    IS '字段级变更明细(JSON数组,每项含section/field/label/before/after,原始值未翻译)';
COMMENT ON COLUMN "public"."chat_ai_robot_config_change_log"."create_time"      IS '创建时间(Unix时间戳)';

-- +goose Down

DROP TABLE IF EXISTS "public"."chat_ai_robot_config_change_log";
