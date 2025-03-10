-- +goose Up

ALTER TABLE "chat_ai_robot"
    ADD COLUMN "prompt_type"   int2 NOT NULL DEFAULT 0,
    ADD COLUMN "prompt_struct" text NOT NULL DEFAULT '';

COMMENT ON COLUMN "chat_ai_robot"."prompt_type" IS '提示词类型:0自定义提示词,1结构化提示词';
COMMENT ON COLUMN "chat_ai_robot"."prompt" IS '自定义提示词';
COMMENT ON COLUMN "chat_ai_robot"."prompt_struct" IS '结构化提示词配置json';

UPDATE chat_ai_robot
SET prompt=concat(prompt, E'\n\n- 请使用markdown格式回答问题。')
WHERE show_type = 1;

UPDATE chat_ai_robot
SET prompt=concat(prompt,
                  E'\n\n- 当你选择的知识点中包含图片、链接数据时，你需要在你的答案对应位置输出这些数据，不要改写或忽略这些数据。');