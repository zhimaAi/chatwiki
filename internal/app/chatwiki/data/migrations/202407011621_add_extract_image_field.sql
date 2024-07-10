-- +goose Up

ALTER TABLE "chat_ai_library_file" ADD COLUMN "enable_extract_image" bool default false;
COMMENT ON COLUMN "chat_ai_library_file"."enable_extract_image" IS '是否提取图片';

ALTER TABLE "chat_ai_library_file_data" ADD COLUMN "images" json default '[]';
ALTER TABLE "chat_ai_library_file_data" ADD CONSTRAINT check_images_array CHECK (json_typeof(images) = 'array');
COMMENT ON COLUMN "chat_ai_library_file_data"."images" IS '图片列表';
