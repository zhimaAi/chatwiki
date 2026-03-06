-- +goose Up

ALTER TABLE "public"."wechat_official_account_batch_send_task"
    ADD COLUMN "task_thumb_url" TEXT NOT NULL DEFAULT '';

ALTER TABLE "public"."wechat_official_account_batch_send_task"
    ADD COLUMN "task_digest" varchar(250) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."task_thumb_url" IS '缩略图';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."task_digest" IS '摘要';
