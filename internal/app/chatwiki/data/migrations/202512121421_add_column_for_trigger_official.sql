-- +goose Up
ALTER TABLE "public"."work_flow_trigger" ADD COLUMN "find_key" varchar(255) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."work_flow_trigger"."find_key" IS '查找的key 通过缓存进行查找';

CREATE INDEX ON "work_flow_trigger" ("find_key");