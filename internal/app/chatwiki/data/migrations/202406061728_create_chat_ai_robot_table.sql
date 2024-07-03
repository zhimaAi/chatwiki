-- +goose Up

CREATE TABLE "chat_ai_robot"
(
    "id"                            serial          NOT NULL primary key,
    "admin_user_id"                 int4            NOT NULL DEFAULT 0,
    "robot_key"                     varchar(15)     NOT NULL DEFAULT '',
    "robot_name"                    varchar(100)    NOT NULL DEFAULT '',
    "robot_intro"                   varchar(1000)   NOT NULL DEFAULT '',
    "robot_avatar"                  varchar(500)    NOT NULL DEFAULT '',
    "prompt"                        varchar(1000)   NOT NULL DEFAULT '',
    "library_ids"                   varchar(1000)   NOT NULL DEFAULT '',
    "welcomes"                      varchar(1000)   NOT NULL DEFAULT '',
    "model_config_id"               int4            NOT NULL DEFAULT 0,
    "use_model"                     varchar(100)    NOT NULL DEFAULT '',
    "rerank_status"                 int2            NOT NULL DEFAULT 0,
    "rerank_model_config_id"        int4            NOT NULL DEFAULT 0,
    "rerank_use_model"              varchar(100)    NOT NULL DEFAULT '',
    "temperature"                   float4          NOT NULL DEFAULT 0.5,
    "max_token"                     int4            NOT NULL DEFAULT 4080,
    "context_pair"                  int4            NOT NULL DEFAULT 1,
    "top_k"                         int4            NOT NULL DEFAULT 5,
    "similarity"                    float4          NOT NULL DEFAULT 0.9,
    "search_type"                   int2            NOT NULL DEFAULT 1,
    "external_config_h5"            varchar(2000)   NOT NULL DEFAULT '',
    "external_config_pc"            varchar(2000)   NOT NULL DEFAULT '',
    "unknown_question_prompt"       json,
    "is_direct"                     bool            NOT NULL DEFAULT false,
    "create_time"                   int4            NOT NULL DEFAULT 0,
    "update_time"                   int4            NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_robot" ("admin_user_id");
CREATE INDEX ON "chat_ai_robot" ("model_config_id");
CREATE INDEX ON "chat_ai_robot" ("rerank_model_config_id");
CREATE UNIQUE INDEX ON "chat_ai_robot" ("robot_key");

COMMENT ON TABLE "chat_ai_robot" IS '文档问答机器人-机器人';

COMMENT ON COLUMN "chat_ai_robot"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_robot"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_robot"."robot_key" IS '对外的key';
COMMENT ON COLUMN "chat_ai_robot"."robot_name" IS '机器人名称';
COMMENT ON COLUMN "chat_ai_robot"."robot_intro" IS '机器人简介';
COMMENT ON COLUMN "chat_ai_robot"."robot_avatar" IS '机器人头像';
COMMENT ON COLUMN "chat_ai_robot"."prompt" IS '提示语';
COMMENT ON COLUMN "chat_ai_robot"."library_ids" IS '知识库ID集合';
COMMENT ON COLUMN "chat_ai_robot"."welcomes" IS '欢迎语{content,[question]}';
COMMENT ON COLUMN "chat_ai_robot"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_robot"."update_time" IS '更新时间';
COMMENT ON COLUMN "chat_ai_robot"."model_config_id" IS '模型配置ID';
COMMENT ON COLUMN "chat_ai_robot"."use_model" IS '使用的模型(枚举值)';
COMMENT ON COLUMN "chat_ai_robot"."rerank_status" IS 'Rerank功能开关:0关1开';
COMMENT ON COLUMN "chat_ai_robot"."rerank_model_config_id" IS 'Rerank模型配置ID';
COMMENT ON COLUMN "chat_ai_robot"."rerank_use_model" IS 'Rerank使用的模型(枚举值)';
COMMENT ON COLUMN "chat_ai_robot"."temperature" IS '模型设置-温度(0~2)';
COMMENT ON COLUMN "chat_ai_robot"."max_token" IS '模型设置-最大token数';
COMMENT ON COLUMN "chat_ai_robot"."context_pair" IS '模型设置-上下文数量';
COMMENT ON COLUMN "chat_ai_robot"."search_type" IS '检索模式:1混合,2向量,3全文';
COMMENT ON COLUMN "chat_ai_robot"."external_config_h5" IS 'h5-对外服务配置json';
COMMENT ON COLUMN "chat_ai_robot"."external_config_pc" IS 'pc-对外服务配置json';
COMMENT ON COLUMN "chat_ai_robot"."unknown_question_prompt" IS '未知问题提示语和问题';
COMMENT ON COLUMN "chat_ai_robot"."is_direct" IS '是否直连llm';
