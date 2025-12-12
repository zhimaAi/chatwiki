-- +goose Up


ALTER TABLE "public"."chat_ai_wechat_app" ADD COLUMN "custom_menu_status" int2 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_wechat_app"."custom_menu_status" IS '微信菜单开关：0关 1开';
