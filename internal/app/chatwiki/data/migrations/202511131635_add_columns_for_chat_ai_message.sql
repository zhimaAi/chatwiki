-- +goose Up

ALTER TABLE "public"."chat_ai_message"
ADD COLUMN  "received_message_type" varchar(100) NOT NULL DEFAULT '',
ADD COLUMN  "received_message"  jsonb ,
ADD COLUMN  "media_id_to_oss_url" varchar(255) NOT NULL DEFAULT '',
ADD COLUMN  "thumb_media_id_to_oss_url" varchar(255) NOT NULL DEFAULT '';

COMMENT ON COLUMN "chat_ai_message"."received_message_type" IS '收到消息类型';
COMMENT ON COLUMN "chat_ai_message"."received_message" IS '收到消息内容';
COMMENT ON COLUMN "chat_ai_message"."media_id_to_oss_url" IS '媒体文件ID转换后的OSS地址';
COMMENT ON COLUMN "chat_ai_message"."thumb_media_id_to_oss_url" IS '缩略图媒体文件ID转换后的OSS地址';