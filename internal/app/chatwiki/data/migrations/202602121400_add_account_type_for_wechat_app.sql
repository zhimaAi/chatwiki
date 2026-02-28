-- +goose Up

ALTER TABLE "public"."chat_ai_wechat_app"
    ADD COLUMN "account_type" int4 NOT NULL DEFAULT -1;
COMMENT ON COLUMN "public"."chat_ai_wechat_app"."account_type" IS '账号类型:默认-1表示未知,1订阅号,2服务号,3小程序';
