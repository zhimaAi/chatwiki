-- +goose Up

ALTER TABLE "public"."chat_ai_session"
    ADD COLUMN "rel_user_id" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_session"."rel_user_id" IS '关联的用户ID';

ALTER TABLE "public"."chat_ai_receiver"
    ADD COLUMN "rel_user_id" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_receiver"."rel_user_id" IS '关联的用户ID';

update public."user" set nick_name = 'admin' where user_name = 'admin';