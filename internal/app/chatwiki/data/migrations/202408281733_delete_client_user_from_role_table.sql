-- +goose Up

DELETE FROM "public"."role" WHERE role_type = 4;