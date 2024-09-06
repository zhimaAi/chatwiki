-- +goose Up

ALTER TABLE "public"."user"
    ADD COLUMN "managed_robot_list" varchar(1000) NOT NULL DEFAULT '',
    ADD COLUMN "managed_library_list" varchar(1000) NOT NULL DEFAULT '',
    ADD COLUMN "managed_form_list" varchar(1000) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."user"."managed_robot_list" IS '管理的机器人';
COMMENT ON COLUMN "public"."user"."managed_library_list" IS '管理的知识库';
COMMENT ON COLUMN "public"."user"."managed_form_list" IS '管理的数据库';