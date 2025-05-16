-- +goose Up

ALTER TABLE "public"."company"
    ADD COLUMN "ali_ocr_key" varchar(64) NOT NULL DEFAULT '',
    ADD COLUMN "ali_ocr_secret" varchar(64) NOT NULL DEFAULT '',
    ADD COLUMN "ali_ocr_switch" int2 NOT NULL DEFAULT 2;

COMMENT ON COLUMN "company"."ali_ocr_key" IS '阿里云OCR配置AccessKeyID';
COMMENT ON COLUMN "company"."ali_ocr_secret" IS '阿里云OCR配置AccessKeySecret';
COMMENT ON COLUMN "company"."ali_ocr_switch" IS '阿里云OCR配置开关 1开启 2关闭';
