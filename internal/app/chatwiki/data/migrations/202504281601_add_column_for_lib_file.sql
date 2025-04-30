-- +goose Up

ALTER TABLE "chat_ai_library_file" ADD COLUMN "doc_auto_renew_minute" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_library_file"."doc_auto_renew_minute" IS '更新时间，按分钟计算。比如：传10就表示：00:10，传100就表示：01:40';