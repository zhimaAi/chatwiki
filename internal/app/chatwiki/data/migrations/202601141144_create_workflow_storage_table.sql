-- +goose Up
CREATE TABLE "public"."work_flow_storage_cache"
(
    "id" serial NOT NULL PRIMARY KEY,
    "admin_user_id" int4 NOT NULL DEFAULT 0,
    "robot_id"      int4 NOT NULL DEFAULT 0,
    "dialog_id"     int4 NOT NULL DEFAULT 0,
    "session_id"    int4 NOT NULL DEFAULT 0,
    "openid"        varchar(100) NOT NULL DEFAULT '',
    "storage"       text NOT NULL DEFAULT '',
    "create_time"   int4 NOT NULL DEFAULT 0,
    "update_time"   int4 NOT NULL DEFAULT 0
);
COMMENT ON TABLE "public"."work_flow_storage_cache" IS '工作流暂存缓存';
COMMENT ON COLUMN "public"."work_flow_storage_cache"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."work_flow_storage_cache"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "public"."work_flow_storage_cache"."dialog_id" IS '对话ID';
COMMENT ON COLUMN "public"."work_flow_storage_cache"."session_id" IS '会话ID';
COMMENT ON COLUMN "public"."work_flow_storage_cache"."storage" IS '暂存的数据';

CREATE INDEX ON "work_flow_storage_cache" ("dialog_id", "session_id" , "create_time");
CREATE INDEX ON "work_flow_storage_cache" ("admin_user_id");