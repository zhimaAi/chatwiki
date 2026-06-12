-- +goose Up

COMMENT ON COLUMN "chat_ai_robot"."application_type" IS '应用类型:0聊天机器人,1工作流,2clawbot机器人';
COMMENT ON COLUMN "chat_ai_chat_monitor"."application_type" IS '应用类型:0聊天机器人,1工作流,2clawbot机器人';

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "search_knowledge_close" int2 NOT NULL DEFAULT 0,
    ADD COLUMN "query_local_docs_close" int2 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_robot"."search_knowledge_close" IS '禁用查询知识库:0开,1关';
COMMENT ON COLUMN "public"."chat_ai_robot"."query_local_docs_close" IS '禁用查询本地文档:0开,1关';
