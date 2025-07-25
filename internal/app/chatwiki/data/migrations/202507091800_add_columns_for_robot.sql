-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
ADD COLUMN "prompt_role_type" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_robot"."prompt_role_type" IS '提示词角色 0:system 1:user';