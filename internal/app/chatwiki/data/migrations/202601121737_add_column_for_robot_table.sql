-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "is_top" INT2 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_robot"."is_top" IS '是否置顶：0默认 1置顶';

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "sort_num" INT4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_robot"."sort_num" IS '排序数字，数字越大越靠前';
