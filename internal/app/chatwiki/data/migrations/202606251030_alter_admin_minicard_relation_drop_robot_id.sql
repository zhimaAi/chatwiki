-- +goose Up

DROP INDEX IF EXISTS "idx_admin_mini_card_relation_robot";

ALTER TABLE "public"."admin_mini_card_relation"
    DROP COLUMN IF EXISTS "robot_id";
    