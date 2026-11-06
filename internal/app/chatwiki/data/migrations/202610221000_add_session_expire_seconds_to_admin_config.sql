-- +goose Up
ALTER TABLE "public"."admin_user_config"
    ADD COLUMN IF NOT EXISTS "session_expire_seconds" INT4 NOT NULL DEFAULT 86400;

COMMENT ON COLUMN "public"."admin_user_config"."session_expire_seconds" IS '机器人会话无新消息自动结束时长，单位秒，默认24小时';

-- +goose Down

ALTER TABLE "public"."admin_user_config"
    DROP COLUMN IF EXISTS "session_expire_seconds";
