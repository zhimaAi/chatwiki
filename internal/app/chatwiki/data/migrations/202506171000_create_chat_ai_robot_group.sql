-- +goose Up

CREATE TABLE "public"."chat_ai_robot_group" (
    "id" serial NOT NULL primary key,
    "admin_user_id" int4 NOT NULL DEFAULT 0,
    "group_name" varchar(50) NOT NULL DEFAULT '',
    "create_time" int4 NOT NULL DEFAULT 0,
    "update_time" int4 NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."chat_ai_robot_group" ("admin_user_id");

COMMENT ON TABLE "public"."chat_ai_robot_group" IS '机器人分组表';

COMMENT ON COLUMN "public"."chat_ai_robot_group"."id" IS 'ID';

COMMENT ON COLUMN "public"."chat_ai_robot_group"."admin_user_id" IS '管理员用户ID';

COMMENT ON COLUMN "public"."chat_ai_robot_group"."group_name" IS '群名称';

COMMENT ON COLUMN "public"."chat_ai_robot_group"."create_time" IS '创建时间';

COMMENT ON COLUMN "public"."chat_ai_robot_group"."update_time" IS '更新时间';

ALTER TABLE "public"."chat_ai_robot"
ADD COLUMN "group_id" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_robot"."group_id" IS '关联分组ID';