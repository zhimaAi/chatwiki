-- +goose Up

ALTER TABLE "public"."chat_ai_library_file"
    ADD COLUMN "async_split_params" text NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_library_file"."async_split_params" IS '异步任务-自动分段参数';
