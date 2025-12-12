-- +goose Up

CREATE TABLE "func_official_custom_menu"
(
    "id"                   serial      NOT NULL primary key,
    "admin_user_id"        int4        NOT NULL DEFAULT 0,
    "appid"                varchar(50) NOT NULL DEFAULT '',
    "seq_id"               int4  NOT NULL DEFAULT 0,
    "menu_name"            varchar(255) NOT NULL DEFAULT '',
    "menu_level"           int4  NOT NULL DEFAULT 0,
    "parent_menu_id"       int4  NOT NULL DEFAULT 0,
    "choose_act_item"      int4  NOT NULL DEFAULT 0,
    "act_params"           jsonb NOT NULL DEFAULT '[]',
    "oper_user_id"         int4  NOT NULL DEFAULT 0,
    "template_id"          int4  NOT NULL DEFAULT 0,
    "batch_id"             int4  NOT NULL DEFAULT 0,
    "create_time"          int4  NOT NULL DEFAULT 0,
    "update_time"          int4  NOT NULL DEFAULT 0
);

CREATE INDEX on "func_official_custom_menu"("appid");
CREATE INDEX on "func_official_custom_menu"("admin_user_id");

COMMENT ON TABLE "func_official_custom_menu" IS '聊天机器人智能菜单';
COMMENT ON COLUMN "func_official_custom_menu"."id" IS '自增ID';
COMMENT ON COLUMN "func_official_custom_menu"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "func_official_custom_menu"."appid" IS '公众号appid';
COMMENT ON COLUMN "func_official_custom_menu"."seq_id" IS '排序id';
COMMENT ON COLUMN "func_official_custom_menu"."menu_name" IS '菜单名称';
COMMENT ON COLUMN "func_official_custom_menu"."menu_level" IS '菜单层级 1根菜单 2二级菜单';
COMMENT ON COLUMN "func_official_custom_menu"."parent_menu_id" IS '上级节点id';
COMMENT ON COLUMN "func_official_custom_menu"."choose_act_item" IS '默认选中栏目 0有子节点时无此功能 1发送消息 2跳转网页 3跳转小程序 4人工客服 5推送事件';
COMMENT ON COLUMN "func_official_custom_menu"."act_params" IS '配置json串';
COMMENT ON COLUMN "func_official_custom_menu"."oper_user_id" IS '操作人ID';
COMMENT ON COLUMN "func_official_custom_menu"."template_id" IS '模版id';
COMMENT ON COLUMN "func_official_custom_menu"."batch_id" IS '批量id';
COMMENT ON COLUMN "func_official_custom_menu"."create_time" IS '创建时间';
COMMENT ON COLUMN "func_official_custom_menu"."update_time" IS '更新时间';
