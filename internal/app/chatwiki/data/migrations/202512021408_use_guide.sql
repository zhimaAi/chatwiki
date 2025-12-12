-- +goose Up
-- library
ALTER TABLE "public"."chat_ai_library" ADD COLUMN "is_default" int2 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "public"."chat_ai_library"."is_default" IS '是否为默认创建，1是，2不是';

-- robot
ALTER TABLE "public"."chat_ai_robot" ADD COLUMN "is_default" int2 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "public"."chat_ai_robot"."is_default" IS '是否为默认创建，1是，2不是';


-- use guide record
CREATE TABLE "public"."use_guide_process" (
  "id" serial NOT NULL primary key,
  "admin_user_id" int4 NOT NULL DEFAULT 0,
  "use_guide_type" varchar(50) NOT NULL DEFAULT '',
  "data" jsonb NOT NULL DEFAULT '{}',
  "create_time" int4 NOT NULL DEFAULT 0,
  "update_time" int4 NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."use_guide_process" ("admin_user_id" , "use_guide_type");

COMMENT ON TABLE "public"."use_guide_process" IS '使用指引';
COMMENT ON COLUMN "public"."use_guide_process"."admin_user_id" IS '管理员ID';
COMMENT ON COLUMN "public"."use_guide_process"."use_guide_type" IS '使用指引类型';
COMMENT ON COLUMN "public"."use_guide_process"."data" IS '配置内容';