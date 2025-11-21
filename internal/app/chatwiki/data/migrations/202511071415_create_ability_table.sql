-- +goose Up

CREATE TABLE "ability"
(
    "id"                serial NOT NULL primary key,
    "admin_user_id"     int4 NOT NULL DEFAULT 0,
    "module_type"         varchar(100) NOT NULL DEFAULT '',
    "ability_type"      varchar(100) NOT NULL DEFAULT '',
    "switch_status"     int2 NOT NULL DEFAULT 0,
    "create_time"       int4 NOT NULL DEFAULT 0,
    "update_time"       int4 NOT NULL DEFAULT 0
);

CREATE INDEX on "ability"("admin_user_id","module_type","ability_type");

COMMENT ON TABLE "ability" IS '账号探索能力表';

COMMENT ON COLUMN "ability"."id" IS '自增ID';
COMMENT ON COLUMN "ability"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "ability"."module_type" IS '模块类型：robot 机器人,work_flow 工作流';
COMMENT ON COLUMN "ability"."ability_type" IS '功能类型：keyword_reply关键词回复';
COMMENT ON COLUMN "ability"."switch_status" IS '开启状态：0关1开';
COMMENT ON COLUMN "ability"."create_time" IS '创建时间';
COMMENT ON COLUMN "ability"."update_time" IS '更新时间';
