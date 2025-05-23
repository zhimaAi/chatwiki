-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "optimize_question_model_config_id" int4 NOT NULL DEFAULT 0,
    ADD COLUMN "optimize_question_use_model" varchar(32) NOT NULL DEFAULT '',
    ADD COLUMN "optimize_question_dialogue_background" varchar(1024) NOT NULL DEFAULT '';

COMMENT ON COLUMN "chat_ai_robot"."optimize_question_model_config_id" IS '问题优化模型ID';
COMMENT ON COLUMN "chat_ai_robot"."optimize_question_use_model" IS '问题优化模型名称';
COMMENT ON COLUMN "chat_ai_robot"."optimize_question_dialogue_background" IS '问题优化对话背景';
