-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "normal_library_id" integer NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_robot"."normal_library_id" IS '新版客服工作区普通知识库';

UPDATE "public"."chat_ai_robot" AS robot
SET "normal_library_id" = library."id"
FROM "public"."chat_ai_library" AS library
WHERE robot."default_library_id" = library."id"
    AND library."type" = 0
    AND robot."admin_user_id" = library."admin_user_id";

-- +goose Down

ALTER TABLE "public"."chat_ai_robot"
    DROP COLUMN "normal_library_id";
