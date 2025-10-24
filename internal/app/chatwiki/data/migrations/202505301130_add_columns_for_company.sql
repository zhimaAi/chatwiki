-- +goose Up

ALTER TABLE "public"."company"
ADD COLUMN "top_navigate" varchar(2000) NOT NULL DEFAULT '';

COMMENT ON COLUMN "company"."top_navigate" IS '顶部导航';