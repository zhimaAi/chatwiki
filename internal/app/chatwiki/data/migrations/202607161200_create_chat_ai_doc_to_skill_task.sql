-- +goose Up

CREATE TABLE "public"."chat_ai_doc_to_skill_task"
(
    "id"                serial        NOT NULL primary key,
    "admin_user_id"     int4          NOT NULL DEFAULT 0,
    "task_batch"        varchar(36)   NOT NULL DEFAULT '',
    "source_files"      text          NOT NULL DEFAULT '',
    "custom_prompt"     text          NOT NULL DEFAULT '',
    "model_config_id"   int4          NOT NULL DEFAULT 0,
    "use_model"         varchar(100)  NOT NULL DEFAULT '',
    "temperature"       float4        NOT NULL DEFAULT 1,
    "max_token"         int4          NOT NULL DEFAULT 32768,
    "status"            int2          NOT NULL DEFAULT 0,
    "skill_name"        varchar(500)  NOT NULL DEFAULT '',
    "skill_description" text          NOT NULL DEFAULT '',
    "file_name"         varchar(255)  NOT NULL DEFAULT '',
    "file_url"          varchar(1000) NOT NULL DEFAULT '',
    "file_size"         int4          NOT NULL DEFAULT 0,
    "debug_log"         text          NOT NULL DEFAULT '',
    "err_msg"           text          NOT NULL DEFAULT '',
    "start_time"        int4          NOT NULL DEFAULT 0,
    "end_time"          int4          NOT NULL DEFAULT 0,
    "create_time"       int4          NOT NULL DEFAULT 0,
    "update_time"       int4          NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "public"."chat_ai_doc_to_skill_task" ("task_batch");
CREATE INDEX ON "public"."chat_ai_doc_to_skill_task" ("admin_user_id", "status", "id");
CREATE INDEX ON "public"."chat_ai_doc_to_skill_task" ("admin_user_id", "id");

COMMENT ON TABLE "public"."chat_ai_doc_to_skill_task" IS 'Doc-to-Skill生成任务表';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."id" IS 'ID';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."task_batch" IS '任务批次UUID';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."source_files" IS '同一任务上传的源文档列表(JSON)';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."custom_prompt" IS '用户自定义提示词';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."model_config_id" IS '模型配置ID';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."use_model" IS '使用的模型(枚举值)';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."temperature" IS '模型设置-温度(0~2)';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."max_token" IS '模型设置-最大token数';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."status" IS '任务状态:0排队中,1生成中,2生成成功,3生成失败,4停止中,5已停止';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."skill_name" IS '生成的skill名称';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."skill_description" IS '生成的skill描述';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."file_name" IS '生成的skill压缩包文件名';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."file_url" IS '生成的skill压缩包文件链接';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."file_size" IS '生成的skill压缩包文件大小';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."debug_log" IS '生成日志(JSON)';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."err_msg" IS '错误信息';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."start_time" IS '开始生成时间';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."end_time" IS '结束生成时间';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."chat_ai_doc_to_skill_task"."update_time" IS '更新时间';

-- +goose Down

DROP TABLE IF EXISTS "public"."chat_ai_doc_to_skill_task";
