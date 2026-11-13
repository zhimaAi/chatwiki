-- +goose Up

ALTER TABLE "public"."chat_ai_message"
    ADD COLUMN "customer_message_id" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_message"."customer_message_id" IS 'AI回复对应的客户提问消息ID，0表示旧记录或未关联提问';

-- +goose Down

ALTER TABLE "public"."chat_ai_message"
    DROP COLUMN "customer_message_id";
