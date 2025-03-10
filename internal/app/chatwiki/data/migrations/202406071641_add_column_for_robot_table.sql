-- +goose Up

ALTER TABLE "chat_ai_robot"
    ADD COLUMN "chat_type" int2 default 3,
    ADD COLUMN "library_qa_direct_reply_switch" bool default false,
    ADD COLUMN "library_qa_direct_reply_score" float4 default 0.9,
    ADD COLUMN "mixture_qa_direct_reply_switch" bool default false,
    ADD COLUMN "mixture_qa_direct_reply_score" float4 default 0.9;
ALTER TABLE "chat_ai_robot" DROP COLUMN IF EXISTS "is_direct";

COMMENT ON COLUMN "chat_ai_robot"."library_qa_direct_reply_switch" IS '仅知识库-QA文档直接回复答案开关';
COMMENT ON COLUMN "chat_ai_robot"."library_qa_direct_reply_score" IS '仅知识库-QA文档直接回复答案的相似度';
COMMENT ON COLUMN "chat_ai_robot"."mixture_qa_direct_reply_switch" IS '混合模式-QA文档直接回复答案开关';
COMMENT ON COLUMN "chat_ai_robot"."mixture_qa_direct_reply_score" IS '混合模式-QA文档直接回复答案的相似度';
COMMENT ON COLUMN "chat_ai_robot"."chat_type" IS '聊天模式 1仅知识库 2直连 3混合';

