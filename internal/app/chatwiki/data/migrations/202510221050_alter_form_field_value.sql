-- +goose Up

ALTER TABLE "public"."form_field_value"
    ALTER COLUMN "string_content" TYPE text;
