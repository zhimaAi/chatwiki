-- +goose Up

CREATE TABLE "func_chat_robot_smart_menu"
(
    "id"                   serial      NOT NULL primary key,
    "admin_user_id"        int4        NOT NULL DEFAULT 0,
    "robot_id"             int4        NOT NULL DEFAULT 0,
    "menu_title"           varchar(255) NOT NULL DEFAULT '',
    "menu_description"     varchar(1000) NOT NULL DEFAULT '',
    "menu_content"         jsonb       NOT NULL DEFAULT '[]',
    "create_time"          int4        NOT NULL DEFAULT 0,
    "update_time"          int4        NOT NULL DEFAULT 0
);

CREATE INDEX on "func_chat_robot_smart_menu"("robot_id");
CREATE INDEX on "func_chat_robot_smart_menu"("admin_user_id");

COMMENT ON TABLE "func_chat_robot_smart_menu" IS '聊天机器人智能菜单';
COMMENT ON COLUMN "func_chat_robot_smart_menu"."id" IS '自增ID';
COMMENT ON COLUMN "func_chat_robot_smart_menu"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "func_chat_robot_smart_menu"."robot_id" IS '机器人ID，为0表示非机器人类型';
COMMENT ON COLUMN "func_chat_robot_smart_menu"."menu_title" IS '菜单标题';
COMMENT ON COLUMN "func_chat_robot_smart_menu"."menu_description" IS '菜单介绍';
COMMENT ON COLUMN "func_chat_robot_smart_menu"."menu_content" IS '菜单内容';
COMMENT ON COLUMN "func_chat_robot_smart_menu"."create_time" IS '创建时间';
COMMENT ON COLUMN "func_chat_robot_smart_menu"."update_time" IS '更新时间';
