-- +goose Up

CREATE TABLE "chat_ai_wechat_app"
(
    "id"            serial       NOT NULL primary key,
    "admin_user_id" int4         NOT NULL DEFAULT 0,
    "robot_id"      int4         NOT NULL DEFAULT 0,
    "app_name"      varchar(100) NOT NULL DEFAULT '',
    "app_id"        varchar(100) NOT NULL DEFAULT '',
    "app_secret"    varchar(500) NOT NULL DEFAULT '',
    "app_avatar"    varchar(500) NOT NULL DEFAULT '',
    "access_key"    varchar(100) NOT NULL DEFAULT '',
    "app_type"      varchar(100) NOT NULL DEFAULT '',
    "set_type"      int2         NOT NULL DEFAULT 1,
    "refresh_token" varchar(500) NOT NULL DEFAULT '',
    "robot_key"     varchar(15)  NOT NULL DEFAULT '',
    "create_time"   int4         NOT NULL DEFAULT 0,
    "update_time"   int4         NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_wechat_app" ("admin_user_id");
CREATE INDEX ON "chat_ai_wechat_app" ("robot_id");
CREATE UNIQUE INDEX ON "chat_ai_wechat_app" ("app_id");
CREATE UNIQUE INDEX ON "chat_ai_wechat_app" ("access_key");

COMMENT ON TABLE "chat_ai_wechat_app" IS '文档问答机器人-微信应用表';

COMMENT ON COLUMN "chat_ai_wechat_app"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_wechat_app"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_wechat_app"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "chat_ai_wechat_app"."app_name" IS '应用名称';
COMMENT ON COLUMN "chat_ai_wechat_app"."app_id" IS '应用ID';
COMMENT ON COLUMN "chat_ai_wechat_app"."app_secret" IS '应用secret';
COMMENT ON COLUMN "chat_ai_wechat_app"."app_avatar" IS '应用头像';
COMMENT ON COLUMN "chat_ai_wechat_app"."access_key" IS '对外的key';
COMMENT ON COLUMN "chat_ai_wechat_app"."app_type" IS '应用类型:official_account公众号,mini_program小程序,wechat_kefu微信客服';
COMMENT ON COLUMN "chat_ai_wechat_app"."set_type" IS '接入类型:1密码接入,2授权接入';
COMMENT ON COLUMN "chat_ai_wechat_app"."refresh_token" IS '授权接入的刷新令牌';
COMMENT ON COLUMN "chat_ai_wechat_app"."robot_key" IS '机器人key';
COMMENT ON COLUMN "chat_ai_wechat_app"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_wechat_app"."update_time" IS '更新时间';
