-- +goose Up

CREATE TABLE "public"."user_search_config"
(
    "id"            serial        NOT NULL primary key,
    "user_id" int4          NOT NULL DEFAULT 0,
    "model_config_id"               int4            NOT NULL DEFAULT 0,
    "use_model"                     varchar(100)    NOT NULL DEFAULT '',
    "rerank_status"                 int2            NOT NULL DEFAULT 0,
    "rerank_model_config_id"        int4            NOT NULL DEFAULT 0,
    "rerank_use_model"              varchar(100)    NOT NULL DEFAULT '',
    "temperature"                   float4          NOT NULL DEFAULT 0.5,
    "max_token"                     int4            NOT NULL DEFAULT 4080,
    "context_pair"                  int4            NOT NULL DEFAULT 1,
    "size"                         int4            NOT NULL DEFAULT 200,
    "similarity"                    float4          NOT NULL DEFAULT 0.9,
    "search_type"                   int2            NOT NULL DEFAULT 1,
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."user_search_config" ( "user_id");

COMMENT ON TABLE "public"."user_search_config" IS '敏感词表';

COMMENT ON COLUMN "public"."user_search_config"."id" IS 'ID';
COMMENT ON COLUMN "public"."user_search_config"."user_id" IS '用户ID';
COMMENT ON COLUMN "public". "user_search_config"."model_config_id" IS '模型配置ID';
COMMENT ON COLUMN "public". "user_search_config"."use_model" IS '使用的模型(枚举值)';
COMMENT ON COLUMN "public". "user_search_config"."rerank_status" IS 'Rerank功能开关:0关1开';
COMMENT ON COLUMN "public". "user_search_config"."rerank_model_config_id" IS 'Rerank模型配置ID';
COMMENT ON COLUMN "public". "user_search_config"."rerank_use_model" IS 'Rerank使用的模型(枚举值)';
COMMENT ON COLUMN "public". "user_search_config"."temperature" IS '模型设置-温度(0~2)';
COMMENT ON COLUMN "public". "user_search_config"."max_token" IS '模型设置-最大token数';
COMMENT ON COLUMN "public". "user_search_config"."context_pair" IS '模型设置-上下文数量';
COMMENT ON COLUMN "public". "user_search_config"."search_type" IS '检索模式:1混合,2向量,3全文';
COMMENT ON COLUMN "public"."user_search_config"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."user_search_config"."update_time" IS '更新时间';