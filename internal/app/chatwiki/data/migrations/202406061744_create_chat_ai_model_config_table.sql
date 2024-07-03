-- +goose Up

CREATE TABLE "chat_ai_model_config"
(
    "id"              serial       NOT NULL primary key,
    "admin_user_id"   int4         NOT NULL DEFAULT 0,
    "model_define"    varchar(50)  NOT NULL DEFAULT '',
    "model_types"     varchar(100) NOT NULL DEFAULT '',
    "deployment_name" varchar(500) NOT NULL DEFAULT '',
    "api_endpoint"    varchar(500) NOT NULL DEFAULT '',
    "api_key"         varchar(500) NOT NULL DEFAULT '',
    "secret_key"      varchar(500) NOT NULL DEFAULT '',
    "api_version"     varchar(500) NOT NULL DEFAULT '',
    "app_id"          int4         NOT NULL DEFAULT 0,
    "create_time"     int4         NOT NULL DEFAULT 0,
    "update_time"     int4         NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_model_config" ("admin_user_id", "model_define");

COMMENT ON TABLE "chat_ai_model_config" IS '文档问答机器人-模型配置';

COMMENT ON COLUMN "chat_ai_model_config"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_model_config"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_model_config"."model_define" IS '模型定义的枚举值';
COMMENT ON COLUMN "chat_ai_model_config"."model_types" IS '模型类型(英文逗号拼接)';
COMMENT ON COLUMN "chat_ai_model_config"."deployment_name" IS '部署名称(非必填)';
COMMENT ON COLUMN "chat_ai_model_config"."api_endpoint" IS 'API域名(非必填)';
COMMENT ON COLUMN "chat_ai_model_config"."api_key" IS 'API的key(必填项)';
COMMENT ON COLUMN "chat_ai_model_config"."secret_key" IS 'API的secret(非必填)';
COMMENT ON COLUMN "chat_ai_model_config"."api_version" IS 'api-version(非必填)';
COMMENT ON COLUMN "chat_ai_model_config"."app_id" IS '应用id';
COMMENT ON COLUMN "chat_ai_model_config"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_model_config"."update_time" IS '更新时间';
