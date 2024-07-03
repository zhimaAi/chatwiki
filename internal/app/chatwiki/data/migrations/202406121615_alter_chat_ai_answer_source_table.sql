-- +goose Up

ALTER TABLE "chat_ai_answer_source"
    ADD COLUMN "type"     int4          NOT NULL DEFAULT 1,
    ADD COLUMN "question" varchar(5000) NOT NULL DEFAULT '',
    ADD COLUMN "answer"   varchar(5000) NOT NULL DEFAULT '';

COMMENT ON COLUMN "chat_ai_answer_source"."type" IS '数据类型 1普通段落 2文档问答 3excel问答';