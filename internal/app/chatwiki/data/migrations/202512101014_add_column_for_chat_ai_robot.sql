-- +goose Up


ALTER TABLE "public"."chat_ai_robot" ADD COLUMN "en_name" varchar(50) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_robot"."en_name" IS '英文名称';

CREATE UNIQUE INDEX ON "public"."chat_ai_robot"
    ("admin_user_id", "en_name") WHERE "en_name" != '';


