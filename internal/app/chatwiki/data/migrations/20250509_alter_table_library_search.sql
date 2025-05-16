-- +goose Up

ALTER TABLE "user_search_config" ADD COLUMN "prompt_type" int4 NOT NULL DEFAULT 0;
ALTER TABLE "user_search_config" ADD COLUMN "prompt" varchar(500) NOT NULL DEFAULT '';

COMMENT ON COLUMN "user_search_config"."prompt" IS '提示词';