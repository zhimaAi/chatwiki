-- +goose Up

ALTER TABLE "chat_ai_message"
    ADD COLUMN "reasoning_content"   text          NOT NULL DEFAULT '';

COMMENT ON COLUMN "chat_ai_message"."reasoning_content" IS '推理过程';