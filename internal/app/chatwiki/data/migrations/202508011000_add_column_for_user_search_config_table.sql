-- +goose Up

ALTER TABLE "public"."user_search_config"
ADD COLUMN "summary_switch" int2 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."user_search_config"."summary_switch" IS '智能总结开关:0关,1开';