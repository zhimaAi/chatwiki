-- +goose Up
ALTER TABLE "public"."company"
    ADD COLUMN "system_mode_selected" smallint NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."company"."system_mode_selected" IS '企业首次模式选择：0 未完成，1 已完成';

-- Preserve the mode for existing companies; new companies must select one via the column default.
UPDATE "public"."company" SET "system_mode_selected" = 1;

-- +goose Down
ALTER TABLE "public"."company" DROP COLUMN "system_mode_selected";
