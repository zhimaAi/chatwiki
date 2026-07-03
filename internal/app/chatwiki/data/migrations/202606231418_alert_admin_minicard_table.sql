-- +goose Up

DROP INDEX IF EXISTS "uniq_admin_mini_card_relation_target";

CREATE UNIQUE INDEX IF NOT EXISTS "uniq_admin_mini_card_relation_target_card"
    ON "public"."admin_mini_card_relation" ("admin_user_id", "target_type", "target_id", "mini_card_id")
    WHERE "delete_time" = 0;
