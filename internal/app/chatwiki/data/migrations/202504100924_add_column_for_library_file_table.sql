-- +goose Up

ALTER TABLE "chat_ai_library_file" ADD COLUMN "pdf_parse_type" int2 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "chat_ai_library_file"."pdf_parse_type" IS 'pdf解析模式 1纯文本 2ocr';