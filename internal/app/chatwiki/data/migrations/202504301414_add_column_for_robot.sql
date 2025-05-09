-- +goose Up

ALTER TABLE "chat_ai_robot" ADD COLUMN "question_guide_num" int4 NOT NULL DEFAULT 3;

COMMENT ON COLUMN "chat_ai_robot"."question_guide_num" IS '用户问题建议数量';