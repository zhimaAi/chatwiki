-- +goose Up

CREATE TABLE "model_backup_config"
(
    "id"              serial       NOT NULL primary key,
    "admin_user_id"   int4         NOT NULL DEFAULT 0,
    "model_config_id" int4         NOT NULL DEFAULT 0,
    "use_model"       varchar(100) NOT NULL DEFAULT '',
    "create_time"     int4         NOT NULL DEFAULT 0,
    "update_time"     int4         NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "model_backup_config" ("admin_user_id");

COMMENT ON TABLE  "model_backup_config" IS '账号维度备用模型配置(仅LLM)';
COMMENT ON COLUMN "model_backup_config"."id"              IS '自增ID';
COMMENT ON COLUMN "model_backup_config"."admin_user_id"   IS '管理员用户ID';
COMMENT ON COLUMN "model_backup_config"."model_config_id" IS '备用模型供应商配置id(chat_ai_model_config.id)';
COMMENT ON COLUMN "model_backup_config"."use_model"       IS '备用具体LLM模型名(chat_ai_model_list.use_model_name)';
COMMENT ON COLUMN "model_backup_config"."create_time"     IS '创建时间戳(Unix秒)';
COMMENT ON COLUMN "model_backup_config"."update_time"     IS '更新时间戳(Unix秒)';

-- +goose Down

DROP TABLE IF EXISTS "model_backup_config";
