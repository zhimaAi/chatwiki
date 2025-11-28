-- +goose Up

ALTER TABLE "public"."chat_ai_wechat_app" ADD COLUMN "sort" int4 NOT NULL DEFAULT '0';

COMMENT ON COLUMN "public"."chat_ai_wechat_app"."sort" IS '排序,越小越靠前';
