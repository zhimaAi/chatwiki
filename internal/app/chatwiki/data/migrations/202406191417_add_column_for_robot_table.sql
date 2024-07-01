-- +goose Up

ALTER TABLE "chat_ai_robot" ADD COLUMN "enable_question_guide" bool default true;
COMMENT ON COLUMN "chat_ai_robot"."enable_question_guide" IS '用户问题建议开关';
