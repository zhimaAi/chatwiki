-- +goose Up
CREATE TABLE "fast_command"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "robot_id" int4          NOT NULL DEFAULT 0,
    "title"        varchar(20)  NOT NULL DEFAULT '',
    "typ" int2          NOT NULL DEFAULT 1,
    "sort" int2 NOT NULL DEFAULT 0,
    "content"       varchar(500) NOT NULL DEFAULT '',
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);
COMMENT ON TABLE "fast_command" IS '快捷指令表';

COMMENT ON COLUMN "fast_command"."id" IS 'ID';
COMMENT ON COLUMN "fast_command"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "fast_command"."robot_id" IS '机器人id';
COMMENT ON COLUMN "fast_command"."title" IS '标题';
COMMENT ON COLUMN "fast_command".typ IS '指令类型 1:输入文本,2:跳转网页';
COMMENT ON COLUMN "fast_command".sort IS '排序';
COMMENT ON COLUMN "fast_command"."content" IS '内容';
COMMENT ON COLUMN "fast_command"."create_time" IS '创建时间';
COMMENT ON COLUMN "fast_command"."update_time" IS '更新时间';