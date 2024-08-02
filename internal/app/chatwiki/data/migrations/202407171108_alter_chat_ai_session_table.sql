-- +goose Up

ALTER TABLE "chat_ai_session"
    ADD COLUMN "app_type" varchar(100) NOT NULL DEFAULT 'yun_h5';

COMMENT ON COLUMN "chat_ai_session"."app_type" IS '应用类型:yun_h5-WebAPP,yun_pc-嵌入网站';