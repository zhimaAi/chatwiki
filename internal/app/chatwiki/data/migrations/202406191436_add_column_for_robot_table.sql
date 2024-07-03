-- +goose Up

ALTER TABLE "chat_ai_robot" ADD COLUMN "enable_question_optimize" bool default false;
COMMENT ON COLUMN "chat_ai_robot"."enable_question_optimize" IS '问题优化开关';

