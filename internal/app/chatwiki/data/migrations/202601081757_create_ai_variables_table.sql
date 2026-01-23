-- +goose Up
CREATE TABLE "public"."chat_ai_variables"
(
    "id"                serial          NOT NULL PRIMARY KEY,
    "admin_user_id"     int4            NOT NULL DEFAULT 0,
    "robot_id"          int4            NOT NULL DEFAULT 0,
    "robot_key"         varchar(100)    NOT NULL DEFAULT '',
    "variable_type"     varchar(20)     NOT NULL DEFAULT '',
    "variable_key"      varchar(20)     NOT NULL DEFAULT '',
    "variable_name"     varchar(20)     NOT NULL DEFAULT '',
    "max_input_length"  int4            NOT NULL DEFAULT 0,
    "default_value"     varchar(100)    NOT NULL DEFAULT '',
    "must_input"        int2            NOT NULL DEFAULT -1,
    "options"           text            NOT NULL DEFAULT '[]',
    "create_time"       int4            NOT NULL DEFAULT 0,
    "update_time"       int4            NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "chat_ai_variables" ("robot_key","variable_key");

COMMENT ON TABLE "public"."chat_ai_variables" IS '机器人变量';
COMMENT ON COLUMN "public"."chat_ai_variables"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_variables"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "public"."chat_ai_variables"."robot_key" IS '';
COMMENT ON COLUMN "public"."chat_ai_variables"."variable_type" IS '变量类型，input_string，input_number，select，checkbox_switch，必填';
COMMENT ON COLUMN "public"."chat_ai_variables"."variable_key" IS '字段key，英文字母和下划线组成，最长10个字符，必填';
COMMENT ON COLUMN "public"."chat_ai_variables"."variable_name" IS '字段名，最长10个字符，必填';
COMMENT ON COLUMN "public"."chat_ai_variables"."max_input_length" IS '最大输入长度，1-50之间，默认10';
COMMENT ON COLUMN "public"."chat_ai_variables"."default_value" IS '默认值';
COMMENT ON COLUMN "public"."chat_ai_variables"."must_input" IS '是否必填，1是，0否（不包含复选框）';
COMMENT ON COLUMN "public"."chat_ai_variables"."options" IS '选项配置，json';

ALTER TABLE "public"."chat_ai_session"
    ADD COLUMN "chat_prompt_variables" text NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_session"."chat_prompt_variables" IS '填写的提示词信息';
