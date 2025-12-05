-- +goose Up

CREATE INDEX on "func_chat_robot_subscribe_reply"("admin_user_id");
CREATE INDEX on "func_chat_robot_subscribe_reply"("appid");

ALTER TABLE "func_chat_robot_subscribe_reply" DROP COLUMN IF EXISTS "message_type";
ALTER TABLE "func_chat_robot_subscribe_reply" DROP COLUMN IF EXISTS "specify_message_type";

DROP INDEX IF EXISTS "func_chat_robot_subscribe_reply_robot_id_idx";