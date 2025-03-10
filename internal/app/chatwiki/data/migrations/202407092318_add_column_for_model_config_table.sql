-- +goose Up

ALTER TABLE "chat_ai_model_config" ADD COLUMN "region" varchar(32) default '';

COMMENT ON COLUMN "chat_ai_model_config"."region" IS '地区';
