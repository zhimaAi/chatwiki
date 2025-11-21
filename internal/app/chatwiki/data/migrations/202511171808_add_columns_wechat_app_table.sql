-- +goose Up

ALTER TABLE "public"."chat_ai_wechat_app"
    ADD COLUMN "encrypt_key" varchar(100) NOT NULL DEFAULT '',
    ADD COLUMN "verification_token" varchar(100) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_wechat_app"."encrypt_key" IS '飞书消息解密key';
COMMENT ON COLUMN "public"."chat_ai_wechat_app"."verification_token" IS '飞书刷新token';