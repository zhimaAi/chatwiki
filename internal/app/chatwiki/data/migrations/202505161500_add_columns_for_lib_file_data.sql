-- +goose Up

ALTER TABLE "public"."chat_ai_library_file_data"
ADD COLUMN "split_status" int2 NOT NULL DEFAULT 0,
ADD COLUMN "split_err_msg" varchar(1024) NOT NULL DEFAULT '';

COMMENT ON COLUMN "chat_ai_library_file_data"."split_status" IS '0:正常,4:失败';