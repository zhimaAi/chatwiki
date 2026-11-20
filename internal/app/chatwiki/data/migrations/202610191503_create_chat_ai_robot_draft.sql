-- +goose Up

CREATE TABLE "public"."chat_ai_robot_draft"
(
    "id"            serial NOT NULL primary key,
    "admin_user_id" int4   NOT NULL DEFAULT 0,
    "robot_id"      int4   NOT NULL DEFAULT 0,
    "draft_data"    text   NOT NULL DEFAULT '',
    "create_time"   int4   NOT NULL DEFAULT 0,
    "update_time"   int4   NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "public"."chat_ai_robot_draft" ("robot_id");
CREATE INDEX ON "public"."chat_ai_robot_draft" ("admin_user_id");

COMMENT ON TABLE "public"."chat_ai_robot_draft" IS '机器人配置草稿表';
COMMENT ON COLUMN "public"."chat_ai_robot_draft"."id" IS 'ID';
COMMENT ON COLUMN "public"."chat_ai_robot_draft"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_robot_draft"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "public"."chat_ai_robot_draft"."draft_data" IS '草稿数据JSON:robot为机器人配置,lang_configs为多语言配置';
COMMENT ON COLUMN "public"."chat_ai_robot_draft"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."chat_ai_robot_draft"."update_time" IS '更新时间';

ALTER TABLE "public"."chat_ai_robot" ADD COLUMN IF NOT EXISTS "published_time" int4 NOT NULL DEFAULT 0;
UPDATE "public"."chat_ai_robot" SET "published_time" = "update_time" WHERE "published_time" = 0;
COMMENT ON COLUMN "public"."chat_ai_robot"."published_time" IS '最后发布时间';

-- +goose Down

ALTER TABLE "public"."chat_ai_robot" DROP COLUMN IF EXISTS "published_time";
DROP TABLE IF EXISTS "public"."chat_ai_robot_draft";
