-- +goose Up

ALTER TABLE "public"."chat_ai_library_file"
ADD COLUMN "delete_time" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_library_file"."delete_time" IS '删除时间';

ALTER TABLE "public"."chat_ai_library_file_data"
ADD COLUMN "delete_time" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_library_file_data"."delete_time" IS '删除时间';

ALTER TABLE "public"."chat_ai_library_file_data_index"
ADD COLUMN "delete_time" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_library_file_data_index"."delete_time" IS '删除时间';