-- +goose Up

ALTER TABLE "chat_ai_robot"
    ADD COLUMN "enable_common_question" bool default false,
    ADD COLUMN "common_question_list" jsonb default '[]';

ALTER TABLE "chat_ai_robot" ADD CONSTRAINT check_common_questions_is_array CHECK (jsonb_typeof(common_question_list) = 'array');

COMMENT ON COLUMN "chat_ai_robot"."enable_common_question" IS '常见问题开关';
COMMENT ON COLUMN "chat_ai_robot"."common_question_list" IS '常见问题列表';
