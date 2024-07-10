-- +goose Up

ALTER TABLE "chat_ai_answer_source" ADD COLUMN "images" json default '[]';
COMMENT ON COLUMN "chat_ai_answer_source"."images" IS '关联的图片';

