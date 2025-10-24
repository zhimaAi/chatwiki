-- +goose Up

ALTER TABLE "public"."chat_ai_library_file_doc" 
  ADD COLUMN "banner_img_url" varchar(300) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_library_file_doc"."banner_img_url" IS '首页背景图，当is_index为1时专用';