-- +goose Up

ALTER TABLE "chat_ai_library_file" ADD COLUMN html_url varchar(500) NOT NULL DEFAULT '';
COMMENT ON COLUMN "chat_ai_library_file"."html_url" IS 'html链接';
