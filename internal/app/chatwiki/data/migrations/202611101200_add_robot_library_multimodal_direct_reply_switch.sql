-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "library_multimodal_direct_reply_switch" int2 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "public"."chat_ai_robot"."library_multimodal_direct_reply_switch" IS '仅知识库-多模态内容直接回复未知问题开关：1开启 2关闭';

-- +goose Down

ALTER TABLE "public"."chat_ai_robot"
    DROP COLUMN "library_multimodal_direct_reply_switch";
