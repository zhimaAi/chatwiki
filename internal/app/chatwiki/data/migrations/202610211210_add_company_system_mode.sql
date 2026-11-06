-- +goose Up

ALTER TABLE "public"."company"
    ADD COLUMN "system_mode" varchar(20) NOT NULL DEFAULT 'ai_service';

COMMENT ON COLUMN "public"."company"."system_mode" IS '企业系统模式：ai_service|agent|advanced，老企业迁移为 advanced，新企业默认 ai_service';

-- Existing companies use advanced mode; new companies use the ai_service column default.
UPDATE "public"."company" SET "system_mode" = 'advanced';

-- +goose Down
ALTER TABLE "public"."company" DROP COLUMN "system_mode";
