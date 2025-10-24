-- +goose Up

ALTER TABLE "chat_ai_session"
    ADD COLUMN "app_id" varchar(100) NOT NULL DEFAULT '';

COMMENT ON COLUMN "chat_ai_session"."app_id" IS '应用ID';

COMMENT ON TABLE "chat_ai_wechat_app" IS '文档问答机器人-微信应用表';