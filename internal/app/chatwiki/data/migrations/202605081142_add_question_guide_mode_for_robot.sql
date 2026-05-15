-- +goose Up

ALTER TABLE "chat_ai_robot" ADD COLUMN "question_guide_mode" int2 NOT NULL DEFAULT 0;
ALTER TABLE "chat_ai_robot" ADD COLUMN "question_guide_prompt" text NOT NULL DEFAULT '';
ALTER TABLE "chat_ai_robot" ADD COLUMN "question_guide_workflow_key" varchar(100) NOT NULL DEFAULT '';

COMMENT ON COLUMN "chat_ai_robot"."question_guide_mode" IS '用户问题建议生成方式: 0=默认方案, 1=自定义提示词, 2=指定工作流';
COMMENT ON COLUMN "chat_ai_robot"."question_guide_prompt" IS '用户问题建议自定义提示词';
COMMENT ON COLUMN "chat_ai_robot"."question_guide_workflow_key" IS '用户问题建议指定工作流机器人key';

-- +goose Down
