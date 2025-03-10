-- +goose Up

ALTER TABLE "user"
    ADD COLUMN "client_side_login_switch" int4 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "user"."client_side_login_switch" IS '客户端登录开关:0不需要登录,1需要登录';