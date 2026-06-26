-- +goose Up

ALTER TABLE "public"."chat_ai_message"
ALTER COLUMN "media_id_to_oss_url" TYPE varchar(1000),
ALTER COLUMN "thumb_media_id_to_oss_url" TYPE varchar(1000);

-- +goose Down

ALTER TABLE "public"."chat_ai_message"
ALTER COLUMN "media_id_to_oss_url" TYPE varchar(255),
ALTER COLUMN "thumb_media_id_to_oss_url" TYPE varchar(255);
