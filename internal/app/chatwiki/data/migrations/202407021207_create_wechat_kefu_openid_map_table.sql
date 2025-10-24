-- +goose Up

CREATE TABLE "wechat_kefu_openid_map"
(
    "id"              serial       NOT NULL primary key,
    "admin_user_id"   int4         NOT NULL DEFAULT 0,
    "openid"          varchar(100) NOT NULL DEFAULT '',
    "app_id"          varchar(100) NOT NULL DEFAULT '',
    "external_userid" varchar(100) NOT NULL DEFAULT '',
    "open_kfid"       varchar(100) NOT NULL DEFAULT '',
    "create_time"     int4         NOT NULL DEFAULT 0,
    "update_time"     int4         NOT NULL DEFAULT 0
);

CREATE INDEX ON "wechat_kefu_openid_map" ("admin_user_id");
CREATE UNIQUE INDEX ON "wechat_kefu_openid_map" ("openid");

COMMENT ON TABLE "wechat_kefu_openid_map" IS '微客客服-自定义openid映射关系表';

COMMENT ON COLUMN "wechat_kefu_openid_map"."id" IS 'ID';
COMMENT ON COLUMN "wechat_kefu_openid_map"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "wechat_kefu_openid_map"."openid" IS '客户ID';
COMMENT ON COLUMN "wechat_kefu_openid_map".app_id IS '应用ID';
COMMENT ON COLUMN "wechat_kefu_openid_map".external_userid IS '客户UserID';
COMMENT ON COLUMN "wechat_kefu_openid_map".open_kfid IS '客服账号ID';
COMMENT ON COLUMN "wechat_kefu_openid_map"."create_time" IS '创建时间';
COMMENT ON COLUMN "wechat_kefu_openid_map"."update_time" IS '更新时间';