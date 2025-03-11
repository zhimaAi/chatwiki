-- +goose Up

CREATE TABLE "chat_ai_chat_monitor"
(
    "id"               bigserial    NOT NULL primary key,
    "admin_user_id"    int4         NOT NULL DEFAULT 0,
    "robot_id"         int4         NOT NULL DEFAULT 0,
    "robot_name"       varchar(100) NOT NULL DEFAULT '',
    "application_type" int4         NOT NULL DEFAULT 0,
    "openid"           varchar(100) NOT NULL DEFAULT '',
    "content"          text         NOT NULL DEFAULT '',
    "start_time"       int8         NOT NULL DEFAULT 0,
    "end_time"         int8         NOT NULL DEFAULT 0,
    "all_use_time"     int8         NOT NULL DEFAULT 0,
    "is_error"         int2         NOT NULL DEFAULT 0,
    "error_msg"        varchar(100) NOT NULL DEFAULT '',
    "question_op"      int8         NOT NULL DEFAULT 0,
    "recall_time"      int8         NOT NULL DEFAULT 0,
    "rerank_time"      int8         NOT NULL DEFAULT 0,
    "request_time"     int8         NOT NULL DEFAULT 0,
    "llm_call_time"    int8         NOT NULL DEFAULT 0,
    "debug_log"        text         NOT NULL DEFAULT '[]',
    "node_logs"        text         NOT NULL DEFAULT '[]',
    "create_time"      int4         NOT NULL DEFAULT 0,
    "update_time"      int4         NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_chat_monitor" ("admin_user_id");
CREATE INDEX ON "chat_ai_chat_monitor" ("robot_id");
CREATE INDEX ON "chat_ai_chat_monitor" ("openid");
CREATE INDEX ON "chat_ai_chat_monitor" ("application_type");
CREATE INDEX ON "chat_ai_chat_monitor" ("is_error");

COMMENT ON COLUMN "chat_ai_chat_monitor"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_chat_monitor"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_chat_monitor"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "chat_ai_chat_monitor"."robot_name" IS '机器人名称';
COMMENT ON COLUMN "chat_ai_chat_monitor"."application_type" IS '应用类型:0聊天机器人,1工作流';
COMMENT ON COLUMN "chat_ai_chat_monitor"."openid" IS '客户openid';
COMMENT ON COLUMN "chat_ai_chat_monitor"."content" IS '消息内容';
COMMENT ON COLUMN "chat_ai_chat_monitor"."start_time" IS '开始时间(毫秒)';
COMMENT ON COLUMN "chat_ai_chat_monitor"."end_time" IS '结束(毫秒)';
COMMENT ON COLUMN "chat_ai_chat_monitor"."all_use_time" IS '总耗时(毫秒)';
COMMENT ON COLUMN "chat_ai_chat_monitor"."is_error" IS '是否出错中断';
COMMENT ON COLUMN "chat_ai_chat_monitor"."error_msg" IS '聊天报错信息';
COMMENT ON COLUMN "chat_ai_chat_monitor"."question_op" IS '问题优化(毫秒)';
COMMENT ON COLUMN "chat_ai_chat_monitor"."recall_time" IS '知识库召回(毫秒)';
COMMENT ON COLUMN "chat_ai_chat_monitor"."rerank_time" IS 'rerank重排(毫秒)';
COMMENT ON COLUMN "chat_ai_chat_monitor"."request_time" IS 'llm请求时间(毫秒)';
COMMENT ON COLUMN "chat_ai_chat_monitor"."llm_call_time" IS 'llm调用时间(毫秒)';
COMMENT ON COLUMN "chat_ai_chat_monitor"."debug_log" IS 'Prompt日志json';
COMMENT ON COLUMN "chat_ai_chat_monitor"."node_logs" IS '节点耗时日志json';
COMMENT ON COLUMN "chat_ai_chat_monitor"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_chat_monitor"."update_time" IS '更新时间';

COMMENT ON TABLE "chat_ai_chat_monitor" IS '聊天响应监控';