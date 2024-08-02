-- +goose Up
CREATE TABLE "chat_ai_robot_apikey"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "robot_key" varchar(15)          NOT NULL DEFAULT '',
    "key"        varchar(128)  NOT NULL DEFAULT '',
    "status" int2          NOT NULL DEFAULT 1,
    "expire_time"   int4          NOT NULL DEFAULT 0,
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_robot_apikey" ("robot_key");
CREATE INDEX ON "chat_ai_robot_apikey" ("key");


COMMENT ON TABLE "chat_ai_robot_apikey" IS 'apikey表';

COMMENT ON COLUMN "chat_ai_robot_apikey"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_robot_apikey"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_robot_apikey"."robot_key" IS '机器人key';
COMMENT ON COLUMN "chat_ai_robot_apikey"."key" IS 'apikey';
COMMENT ON COLUMN "chat_ai_robot_apikey".status IS '状态 1:启用';
COMMENT ON COLUMN "chat_ai_robot_apikey"."expire_time" IS '过期时间 0 为永久';
COMMENT ON COLUMN "chat_ai_robot_apikey"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_robot_apikey"."update_time" IS '更新时间';
