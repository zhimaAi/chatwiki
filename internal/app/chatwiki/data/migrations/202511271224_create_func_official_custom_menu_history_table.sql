-- +goose Up

CREATE TABLE "func_official_custom_menu_history"
(
    "id"                   serial      NOT NULL primary key,
    "admin_user_id"        int4        NOT NULL DEFAULT 0,
    "appid"                varchar(50) NOT NULL DEFAULT '',
    "history_menu"         jsonb NOT NULL DEFAULT '[]',
    "mew_menu"             jsonb NOT NULL DEFAULT '[]',
    "oper_user_id"         int4  NOT NULL DEFAULT 0,
    "create_time"          int4  NOT NULL DEFAULT 0,
    "update_time"          int4  NOT NULL DEFAULT 0
);

CREATE INDEX on "func_official_custom_menu_history"("appid");
CREATE INDEX on "func_official_custom_menu_history"("admin_user_id");

COMMENT ON TABLE "func_official_custom_menu_history" IS '聊天机器人智能菜单';
COMMENT ON COLUMN "func_official_custom_menu_history"."id" IS '自增ID';
COMMENT ON COLUMN "func_official_custom_menu_history"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "func_official_custom_menu_history"."appid" IS '公众号appid';
COMMENT ON COLUMN "func_official_custom_menu_history"."history_menu" IS '旧菜单';
COMMENT ON COLUMN "func_official_custom_menu_history"."mew_menu" IS '新菜单';
COMMENT ON COLUMN "func_official_custom_menu_history"."oper_user_id" IS '操作人ID';
COMMENT ON COLUMN "func_official_custom_menu_history"."create_time" IS '创建时间';
COMMENT ON COLUMN "func_official_custom_menu_history"."update_time" IS '更新时间';
