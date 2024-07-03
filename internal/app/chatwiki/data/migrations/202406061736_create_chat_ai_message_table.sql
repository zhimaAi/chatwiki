-- +goose Up

CREATE TABLE "chat_ai_message"
(
    "id"            serial          NOT NULL primary key,
    "admin_user_id" int4            NOT NULL DEFAULT 0,
    "robot_id"      int4            NOT NULL DEFAULT 0,
    "openid"        varchar(100)    NOT NULL DEFAULT '',
    "is_customer"   int2            NOT NULL DEFAULT 0,
    "msg_type"      int2            NOT NULL DEFAULT 0,
    "content"       varchar(1000)   NOT NULL DEFAULT '',
    "menu_json"     varchar(1000)   NOT NULL DEFAULT '',
    "quote_file"    varchar(1000)   NOT NULL DEFAULT '[]',
    "dialogue_id"   int4            NOT NULL DEFAULT 0,
    "session_id"    int4            NOT NULL DEFAULT 0,
    "create_time"   int4            NOT NULL DEFAULT 0,
    "update_time"   int4            NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_message" ("admin_user_id");
CREATE INDEX ON "chat_ai_message" ("robot_id");
CREATE INDEX ON "chat_ai_message" ("openid", "robot_id");
CREATE INDEX ON "chat_ai_message" ("dialogue_id");
CREATE INDEX ON "chat_ai_message" ("session_id");

COMMENT ON TABLE "chat_ai_message" IS '文档问答机器人-消息表';

COMMENT ON COLUMN "chat_ai_message"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_message"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_message"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "chat_ai_message"."openid" IS '客户ID';
COMMENT ON COLUMN "chat_ai_message"."is_customer" IS '消息来自:0机器人,1客户';
COMMENT ON COLUMN "chat_ai_message"."msg_type" IS '消息类型:1文本,2菜单';
COMMENT ON COLUMN "chat_ai_message"."content" IS '消息内容';
COMMENT ON COLUMN "chat_ai_message"."menu_json" IS '菜单json';
COMMENT ON COLUMN "chat_ai_message"."quote_file" IS '答案来源json数据';
COMMENT ON COLUMN "chat_ai_message"."dialogue_id" IS '对话ID';
COMMENT ON COLUMN "chat_ai_message"."session_id" IS '会话ID';
COMMENT ON COLUMN "chat_ai_message"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_message"."update_time" IS '更新时间';
