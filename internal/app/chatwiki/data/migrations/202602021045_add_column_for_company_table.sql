-- +goose Up

ALTER TABLE "public"."company"
    ADD COLUMN "cookie_tip_positions" varchar(250) NOT NULL DEFAULT 'webapp,login';

ALTER TABLE "public"."company"
    ADD COLUMN "cookie_tip_ip_locations" varchar(250) NOT NULL DEFAULT 'overseas';

COMMENT ON COLUMN "public"."company"."cookie_tip_positions" IS 'cookie隐私提示位置：home webapp login，多个英文逗号分割';
COMMENT ON COLUMN "public"."company"."cookie_tip_ip_locations" IS 'cookie隐私提示IP范围，domestic国内，overseas海外';
