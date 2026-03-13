-- +goose Up

ALTER TABLE "public"."user"
    ADD COLUMN "regist_from" varchar(30) NOT NULL DEFAULT 'cloud';

COMMENT ON COLUMN "public"."user"."regist_from" IS '注册来源';
