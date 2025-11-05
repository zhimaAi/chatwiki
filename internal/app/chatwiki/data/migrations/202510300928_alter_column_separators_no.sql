-- +goose Up

ALTER TABLE "public"."chat_ai_library_file"
    ALTER COLUMN "separators_no" TYPE varchar(500),
    ALTER COLUMN "father_chunk_separators_no" TYPE varchar(500),
    ALTER COLUMN "son_chunk_separators_no" TYPE varchar(500);

ALTER TABLE "public"."chat_ai_library"
    ALTER COLUMN "normal_chunk_default_separators_no" TYPE varchar(500),
    ALTER COLUMN "father_chunk_separators_no" TYPE varchar(500),
    ALTER COLUMN "son_chunk_separators_no" TYPE varchar(500);

ALTER TABLE "public"."chat_ai_faq_files"
    ALTER COLUMN "separators_no" TYPE varchar(500);