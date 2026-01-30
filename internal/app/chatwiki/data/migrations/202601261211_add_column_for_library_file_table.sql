-- +goose Up

ALTER TABLE "public"."chat_ai_library_file" ADD COLUMN "thumb_path" varchar(500) NOT NULL DEFAULT '';


COMMENT ON COLUMN "public"."chat_ai_library_file"."thumb_path" IS '文件缩略图url';
