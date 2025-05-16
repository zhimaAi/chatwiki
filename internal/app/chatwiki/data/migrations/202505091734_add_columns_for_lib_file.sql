-- +goose Up

ALTER TABLE "public"."chat_ai_library_file"
    ADD COLUMN "ali_ocr_job_id" varchar(64) NOT NULL DEFAULT '';

COMMENT ON COLUMN "chat_ai_library_file"."ali_ocr_job_id" IS '阿里云OCR任务id';
