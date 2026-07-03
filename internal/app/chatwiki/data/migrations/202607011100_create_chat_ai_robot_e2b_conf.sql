-- +goose Up

CREATE TABLE "public"."chat_ai_robot_e2b_conf"
(
    "id"              serial        NOT NULL primary key,
    "admin_user_id"   int4          NOT NULL DEFAULT 0,
    "robot_key"       varchar(15)   NOT NULL DEFAULT '',
    "switch_status"   int2          NOT NULL DEFAULT 0,
    "api_key"         varchar(500)  NOT NULL DEFAULT '',
    "api_base_url"    varchar(1000) NOT NULL DEFAULT '',
    "sandbox_domain"  varchar(1000) NOT NULL DEFAULT '',
    "template"        varchar(500)  NOT NULL DEFAULT '',
    "timeout"         int4          NOT NULL DEFAULT 0,
    "command_timeout" int4          NOT NULL DEFAULT 0,
    "command_user"    varchar(100)  NOT NULL DEFAULT '',
    "create_time"     int4          NOT NULL DEFAULT 0,
    "update_time"     int4          NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "public"."chat_ai_robot_e2b_conf" ("robot_key");
CREATE INDEX ON "public"."chat_ai_robot_e2b_conf" ("admin_user_id");

COMMENT ON TABLE "public"."chat_ai_robot_e2b_conf" IS '机器人E2B配置表';
COMMENT ON COLUMN "public"."chat_ai_robot_e2b_conf"."id" IS 'ID';
COMMENT ON COLUMN "public"."chat_ai_robot_e2b_conf"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_robot_e2b_conf"."robot_key" IS '机器人key';
COMMENT ON COLUMN "public"."chat_ai_robot_e2b_conf"."switch_status" IS 'E2B开关:0关,1开';
COMMENT ON COLUMN "public"."chat_ai_robot_e2b_conf"."api_key" IS 'E2B APIKey';
COMMENT ON COLUMN "public"."chat_ai_robot_e2b_conf"."api_base_url" IS 'E2B APIBaseURL';
COMMENT ON COLUMN "public"."chat_ai_robot_e2b_conf"."sandbox_domain" IS 'E2B SandboxDomain';
COMMENT ON COLUMN "public"."chat_ai_robot_e2b_conf"."template" IS 'E2B Template';
COMMENT ON COLUMN "public"."chat_ai_robot_e2b_conf"."timeout" IS 'E2B Timeout';
COMMENT ON COLUMN "public"."chat_ai_robot_e2b_conf"."command_timeout" IS 'E2B COMMAND_TIMEOUT';
COMMENT ON COLUMN "public"."chat_ai_robot_e2b_conf"."command_user" IS 'E2B COMMAND_USER';
COMMENT ON COLUMN "public"."chat_ai_robot_e2b_conf"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."chat_ai_robot_e2b_conf"."update_time" IS '更新时间';

-- +goose Down

DROP TABLE IF EXISTS "public"."chat_ai_robot_e2b_conf";
