-- +goose Up

ALTER TABLE "public"."chat_ai_message"
ADD COLUMN  "reply_content_list"  jsonb ;

COMMENT ON COLUMN "chat_ai_message"."reply_content_list" IS '非ai额外回复内容列表JSON';