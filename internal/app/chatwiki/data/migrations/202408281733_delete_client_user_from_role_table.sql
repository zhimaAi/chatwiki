-- +goose Up

DELETE FROM "public"."role" WHERE id = 1;
UPDATE "public"."user" SET user_roles = 2 WHERE id = 1;
