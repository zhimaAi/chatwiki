-- +goose Up

COMMENT ON COLUMN "chat_ai_library"."type" IS '知识库类型:0普通知识库,1对外知识库,2问答知识库';