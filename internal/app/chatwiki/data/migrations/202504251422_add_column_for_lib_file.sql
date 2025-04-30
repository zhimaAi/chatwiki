-- +goose Up

ALTER TABLE "chat_ai_library_file" ADD COLUMN "ocr_pdf_total" int4 NOT NULL DEFAULT 0;
ALTER TABLE "chat_ai_library_file" ADD COLUMN "ocr_pdf_index" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_library_file"."ocr_pdf_total" IS 'ocr解析pdf总页数';
COMMENT ON COLUMN "chat_ai_library_file"."ocr_pdf_index" IS 'ocr解析pdf当前页数';
