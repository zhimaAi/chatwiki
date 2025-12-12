-- +goose Up

CREATE TABLE "trigger_config"
(
    "id"                   serial       NOT NULL primary key,
    "admin_user_id"        int4         NOT NULL DEFAULT 0,
    "switch_status"        int2         NOT NULL DEFAULT 1,
    "name"                 varchar(255) NOT NULL DEFAULT '',
    "trigger_type"         varchar(255) NOT NULL DEFAULT '',
    "from_type"            varchar(255) NOT NULL DEFAULT '',
    "intro"                varchar(500) NOT NULL DEFAULT '',
    "author"               varchar(255) NOT NULL DEFAULT '',
    "create_time"          int4         NOT NULL DEFAULT 0,
    "update_time"          int4         NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX on "trigger_config"("admin_user_id" , "trigger_type");

COMMENT ON TABLE "trigger_config" IS '触发器';
COMMENT ON COLUMN "trigger_config"."id" IS '自增ID';
COMMENT ON COLUMN "trigger_config"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "trigger_config"."switch_status" IS '开关，1开启，0关闭';
COMMENT ON COLUMN "trigger_config"."name" IS '触发器名称';
COMMENT ON COLUMN "trigger_config"."trigger_type" IS '触发器类型,1会话触发器 2测试触发器  3定时触发器';
COMMENT ON COLUMN "trigger_config"."from_type" IS '触发器来源，inherent内置';
COMMENT ON COLUMN "trigger_config"."intro" IS '介绍';
COMMENT ON COLUMN "trigger_config"."author" IS '作者';
COMMENT ON COLUMN "trigger_config"."create_time" IS '创建时间';
COMMENT ON COLUMN "trigger_config"."update_time" IS '更新时间';

-- work_flow_trigger
CREATE TABLE "work_flow_trigger"
(
    "id"                   serial       NOT NULL primary key,
    "admin_user_id"        int4         NOT NULL DEFAULT 0,
    "robot_id"             int4         NOT NULL DEFAULT 0,
    "trigger_type"         varchar(255) NOT NULL DEFAULT '',
    "trigger_json"         json         NOT NULL DEFAULT '{}',
    "cron_entry_id"        int4         NOT NULL DEFAULT 0,
    "is_finish"            int2         NOT NULL DEFAULT 0,
    "last_msg"             varchar(1000) NOT NULL DEFAULT '',
    "create_time"          int4         NOT NULL DEFAULT 0,
    "update_time"          int4         NOT NULL DEFAULT 0
);

CREATE INDEX on "work_flow_trigger"("trigger_type" , "cron_entry_id" , "is_finish");
CREATE INDEX on "work_flow_trigger"("admin_user_id","robot_id");

COMMENT ON TABLE "work_flow_trigger" IS '工作流触发器';
COMMENT ON COLUMN "work_flow_trigger"."id" IS '自增ID';
COMMENT ON COLUMN "work_flow_trigger"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "work_flow_trigger"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "work_flow_trigger"."trigger_type" IS '触发器类型，3定时触发器，4公众号触发器';
COMMENT ON COLUMN "work_flow_trigger"."trigger_json" IS '触发器配置';
COMMENT ON COLUMN "work_flow_trigger"."cron_entry_id" IS 'linux crontab 任务ID';
COMMENT ON COLUMN "work_flow_trigger"."is_finish" IS '针对一次性任务，1已完成，0待处理';
COMMENT ON COLUMN "work_flow_trigger"."last_msg" IS '最后一次结果';
COMMENT ON COLUMN "work_flow_trigger"."create_time" IS '创建时间';
COMMENT ON COLUMN "work_flow_trigger"."update_time" IS '更新时间';
