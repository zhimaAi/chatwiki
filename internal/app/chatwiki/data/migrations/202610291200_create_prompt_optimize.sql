-- +goose Up

CREATE TABLE "public"."chat_ai_prompt_optimize_record"
(
    "id"                bigserial     NOT NULL PRIMARY KEY,
    "admin_user_id"     int4          NOT NULL,
    "robot_id"          int4          NOT NULL,
    "request_id"        varchar(64)   NOT NULL,
    "demand"            text          NOT NULL,
    "scene"             varchar(16)   NOT NULL CHECK ("scene" IN ('config', 'conversation')),
    "mode"              varchar(16)   NOT NULL CHECK ("mode" IN ('optimize', 'rewrite')),
    "model_config_id"   int4          NOT NULL,
    "use_model"         varchar(255)  NOT NULL,
    "base_version_id"   int8          NOT NULL DEFAULT 0,
    "message_id"        int8          NOT NULL DEFAULT 0,
    "input_context"     text          NOT NULL DEFAULT '',
    "summary"           text          NOT NULL DEFAULT '',
    "version_id"        int8          NOT NULL DEFAULT 0,
    "status"            varchar(16)   NOT NULL DEFAULT 'pending' CHECK ("status" IN ('pending', 'success', 'failed', 'cancelled')),
    "error"             text          NOT NULL DEFAULT '',
    "create_time"       int8          NOT NULL,
    "update_time"       int8          NOT NULL,
    UNIQUE ("admin_user_id", "robot_id", "request_id")
);

CREATE INDEX "prompt_optimize_record_history" ON "public"."chat_ai_prompt_optimize_record" ("admin_user_id", "robot_id", "id" DESC);
CREATE INDEX "prompt_optimize_record_context" ON "public"."chat_ai_prompt_optimize_record" ("admin_user_id", "robot_id", "create_time" DESC, "id" DESC) WHERE "status" = 'success';

COMMENT ON TABLE "public"."chat_ai_prompt_optimize_record" IS 'AI优化提示词-对话轮次记录';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."id" IS '主键ID';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."admin_user_id" IS '管理员ID（租户隔离）';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."robot_id" IS '机器人ID，支持聊天机器人和Agent';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."request_id" IS '请求唯一标识，用于防止重复生成';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."demand" IS '本轮用户提出的提示词调整要求';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."scene" IS '入口场景：config配置入口，conversation对话入口';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."mode" IS '生成方式：optimize优化，rewrite重写';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."model_config_id" IS '本轮使用的模型配置ID';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."use_model" IS '本轮使用的模型名称';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."base_version_id" IS '基准提示词版本ID，0表示使用机器人当前已保存的提示词';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."message_id" IS '原业务AI回复消息ID，0表示未引用业务问答';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."input_context" IS '本轮模型输入上下文JSON，包含需求、基准提示词、能力说明及可选业务问答';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."summary" IS '本轮生成结果摘要';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."version_id" IS '成功生成的提示词版本ID，0表示尚未生成正式版本';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."status" IS '轮次状态：pending处理中，success成功，failed失败，cancelled已取消';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."error" IS '失败或取消原因';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."create_time" IS '创建时间（秒级时间戳）';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_record"."update_time" IS '更新时间（秒级时间戳）';

CREATE TABLE "public"."chat_ai_prompt_optimize_version"
(
    "id"                bigserial     NOT NULL PRIMARY KEY,
    "admin_user_id"     int4          NOT NULL,
    "robot_id"          int4          NOT NULL,
    "record_id"         int8          NOT NULL UNIQUE REFERENCES "public"."chat_ai_prompt_optimize_record" ("id"),
    "version"           int4          NOT NULL CHECK ("version" > 0),
    "prompt"            text          NOT NULL,
    "create_time"       int8          NOT NULL,
    UNIQUE ("admin_user_id", "robot_id", "version")
);

COMMENT ON TABLE "public"."chat_ai_prompt_optimize_version" IS 'AI优化提示词-生成版本（不记录采用状态）';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_version"."id" IS '主键ID';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_version"."admin_user_id" IS '管理员ID（租户隔离）';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_version"."robot_id" IS '所属机器人ID';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_version"."record_id" IS '关联的优化轮次记录ID，一个成功轮次对应一个版本';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_version"."version" IS '机器人内递增的版本号，从1开始';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_version"."prompt" IS '生成的完整提示词正文';
COMMENT ON COLUMN "public"."chat_ai_prompt_optimize_version"."create_time" IS '创建时间（秒级时间戳）';

-- +goose Down

DROP TABLE IF EXISTS "public"."chat_ai_prompt_optimize_version";
DROP TABLE IF EXISTS "public"."chat_ai_prompt_optimize_record";
