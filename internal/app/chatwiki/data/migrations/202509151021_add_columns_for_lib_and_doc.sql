-- +goose Up

ALTER TABLE "public"."chat_ai_library" 
  ADD COLUMN "icon_template_config_id" int4 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "public"."chat_ai_library"."icon_template_config_id" IS '图标模板配置ID';


ALTER TABLE "public"."chat_ai_library_file_doc" 
  ADD COLUMN "is_dir" int2 NOT NULL DEFAULT 0,
  ADD COLUMN "quick_doc_content" text NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_library_file_doc"."is_dir" IS '是否是文件夹 0：否认 1：是';

COMMENT ON COLUMN "public"."chat_ai_library_file_doc"."quick_doc_content" IS '首页快捷文档，当is_index为1时专用';

ALTER TABLE "public"."chat_ai_library_file_doc" 
  ADD COLUMN "doc_icon" varchar(300) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_library_file_doc"."doc_icon" IS '文档icon';