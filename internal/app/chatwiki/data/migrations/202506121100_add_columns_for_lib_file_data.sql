-- +goose Up

ALTER TABLE "public"."chat_ai_library_file_data"
ADD COLUMN "yesterday_hits" int4 NOT NULL DEFAULT 0,
ADD COLUMN "today_hits" int4 NOT NULL DEFAULT 0,
ADD COLUMN "total_hits" int8 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_library_file_data"."yesterday_hits" IS '昨日访问量';

COMMENT ON COLUMN "chat_ai_library_file_data"."today_hits" IS '今日访问量';

COMMENT ON COLUMN "chat_ai_library_file_data"."total_hits" IS '总访问量';