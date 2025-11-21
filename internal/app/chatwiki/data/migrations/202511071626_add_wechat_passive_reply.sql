-- +goose Up

ALTER TABLE "public"."chat_ai_wechat_app"
    ADD COLUMN "account_customer_type" int4 NOT NULL DEFAULT -1;
COMMENT ON COLUMN "public"."chat_ai_wechat_app"."account_customer_type" IS '微信应用认证类型:默认-1,表示未知,认为已认证';

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "wechat_not_verify_hand_get_reply" varchar(100) NOT NULL DEFAULT E'正在思考中，请稍后点击下方蓝字\r\n获取回复👇👇👇',
    ADD COLUMN "wechat_not_verify_hand_get_word"  varchar(100) NOT NULL DEFAULT '👉👉点我获取回复👈👈',
    ADD COLUMN "wechat_not_verify_hand_get_next"  varchar(100) NOT NULL DEFAULT '内容较多，点此查看下文';
COMMENT ON COLUMN "public"."chat_ai_robot"."wechat_not_verify_hand_get_reply" IS '微信应用未认证:回复消息提示语';
COMMENT ON COLUMN "public"."chat_ai_robot"."wechat_not_verify_hand_get_word" IS '微信应用未认证:获取回复蓝字文案';
COMMENT ON COLUMN "public"."chat_ai_robot"."wechat_not_verify_hand_get_next" IS '微信应用未认证:获取下文蓝字文案';