-- +goose Up

CREATE TABLE "llm_request_logs"
(
    "id"                serial NOT NULL primary key,
    "admin_user_id"     int4 NOT NULL DEFAULT 0,
    "openid"            varchar(100) NOT NULL DEFAULT '',
    "corp"              varchar(32) NOT NULL DEFAULT '',
    "model"             varchar(32) NOT NULL DEFAULT '',
    "app_type"          varchar(32) NOT NULL DEFAULT '',
    "type"              varchar(16) NOT NULL DEFAULT '',
    "prompt_token"      int4 NOT NULL DEFAULT 0,
    "completion_token"  int4 NOT NULL DEFAULT 0,
    "source_robot_id"   int4 NOT NULL DEFAULT 0,
    "source_robot"      json NOT NULL DEFAULT '{}',
    "source_library"    json NOT NULL DEFAULT '{}',
    "source_file"       json NOT NULL DEFAULT '{}',
    "request_detail"    json NOT NULL DEFAULT '{}',
    "response_detail"   json NOT NULL DEFAULT '{}',
    "create_time"       int4 NOT NULL DEFAULT 0,
    "update_time"       int4 NOT NULL DEFAULT 0
);

CREATE INDEX on "llm_request_logs"("admin_user_id", "openid", "corp", "model", "app_type", "type", "source_robot_id");

COMMENT ON TABLE "llm_request_logs" IS '大模型请求日志';

COMMENT ON COLUMN "llm_request_logs"."id" IS '自增ID';
COMMENT ON COLUMN "llm_request_logs"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "llm_request_logs"."openid" IS '客户ID';
COMMENT ON COLUMN "llm_request_logs"."corp" IS '模型的企业';
COMMENT ON COLUMN "llm_request_logs"."model" IS '模型名';
COMMENT ON COLUMN "llm_request_logs"."app_type" IS '应用类型';
COMMENT ON COLUMN "llm_request_logs"."type" IS '类型。llm、embedding';
COMMENT ON COLUMN "llm_request_logs"."prompt_token" IS '输入token量';
COMMENT ON COLUMN "llm_request_logs"."completion_token" IS '输出token量';
COMMENT ON COLUMN "llm_request_logs"."source_robot_id" IS '来源机器人id';
COMMENT ON COLUMN "llm_request_logs"."source_robot" IS '来源机器人';
COMMENT ON COLUMN "llm_request_logs"."source_library" IS '来源知识库';
COMMENT ON COLUMN "llm_request_logs"."source_file" IS '来源文档';
COMMENT ON COLUMN "llm_request_logs"."request_detail" IS '请求的一些细节数据';
COMMENT ON COLUMN "llm_request_logs"."response_detail" IS '响应的一些细节数据';
COMMENT ON COLUMN "llm_request_logs"."create_time" IS '创建时间';
COMMENT ON COLUMN "llm_request_logs"."update_time" IS '更新时间';
