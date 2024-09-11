-- +goose Up

DELETE FROM "public"."role" WHERE role_type = 4;
UPDATE "public"."user" SET user_roles = (SELECT id FROM "public"."role" WHERE role_type = 1) WHERE user_name = 'admin';
