-- +goose Up

ALTER TABLE "chat_ai_robot"
    ADD COLUMN "sensitive_words_switch"   int2 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_robot"."sensitive_words_switch" IS '敏感词开关 0：关 1:开';