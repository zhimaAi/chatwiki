-- +goose Up

CREATE TABLE "llm_model_error_logs"
(
    "id"               serial       NOT NULL primary key,
    "admin_user_id"    int4         NOT NULL DEFAULT 0,
    "date"             date         NOT NULL,
    "model_config_id"  int4         NOT NULL DEFAULT 0,
    "corp"             varchar(32)  NOT NULL DEFAULT '',
    "model"            varchar(100) NOT NULL DEFAULT '',
    "robot_id"         int4         NOT NULL DEFAULT 0,
    "robot_name"       varchar(100) NOT NULL DEFAULT '',
    "source_robot"     json         NOT NULL DEFAULT '{}',
    "application_type" int2         NOT NULL DEFAULT 0,
    "error_msg"        text         NOT NULL DEFAULT '',
    "create_time"      int4         NOT NULL DEFAULT 0
);

CREATE INDEX ON "llm_model_error_logs" ("admin_user_id", "create_time");
CREATE INDEX ON "llm_model_error_logs" ("admin_user_id", "date");

COMMENT ON TABLE  "llm_model_error_logs" IS '模型调用异常明细日志(每次报错一行)';
COMMENT ON COLUMN "llm_model_error_logs"."id"               IS '自增ID';
COMMENT ON COLUMN "llm_model_error_logs"."admin_user_id"    IS '管理员用户ID';
COMMENT ON COLUMN "llm_model_error_logs"."date"             IS '日期(写入时落,前端按日分组用)';
COMMENT ON COLUMN "llm_model_error_logs"."model_config_id"  IS '供应商配置id';
COMMENT ON COLUMN "llm_model_error_logs"."corp"             IS '模型企业(供应商),如 DeepSeek';
COMMENT ON COLUMN "llm_model_error_logs"."model"            IS '具体模型名(use_model),如 deepseek-reasoner';
COMMENT ON COLUMN "llm_model_error_logs"."robot_id"         IS '应用id,0=系统/其他';
COMMENT ON COLUMN "llm_model_error_logs"."robot_name"       IS '应用名快照(写入时定格,改名后保留历史名)';
COMMENT ON COLUMN "llm_model_error_logs"."source_robot"     IS '机器人快照';
COMMENT ON COLUMN "llm_model_error_logs"."application_type" IS '快照:0机器人/1工作流(供前端区分跳转目标)';
COMMENT ON COLUMN "llm_model_error_logs"."error_msg"        IS '详细报错文本';
COMMENT ON COLUMN "llm_model_error_logs"."create_time"      IS 'Unix时间戳';

-- +goose Down

DROP TABLE IF EXISTS "llm_model_error_logs";
