-- +goose Up

CREATE TABLE "public"."chat_ai_model_list"
(
    "id"                    bigserial    NOT NULL primary key,
    "admin_user_id"         int4         NOT NULL DEFAULT 0,
    "model_config_id"       int4         NOT NULL DEFAULT 0,
    "model_type"            varchar(100) NOT NULL DEFAULT '',
    "use_model_name"        varchar(100) NOT NULL DEFAULT '',
    "show_model_name"       varchar(50)           default '',
    "thinking_type"         int2         NOT NULL DEFAULT 0,
    "function_call"         int2         NOT NULL DEFAULT 0,
    "input_text"            int2         NOT NULL DEFAULT 0,
    "input_voice"           int2         NOT NULL DEFAULT 0,
    "input_image"           int2         NOT NULL DEFAULT 0,
    "input_video"           int2         NOT NULL DEFAULT 0,
    "input_document"        int2         NOT NULL DEFAULT 0,
    "output_text"           int2         NOT NULL DEFAULT 0,
    "output_voice"          int2         NOT NULL DEFAULT 0,
    "output_image"          int2         NOT NULL DEFAULT 0,
    "output_video"          int2         NOT NULL DEFAULT 0,
    "vector_dimension_list" varchar(100) NOT NULL DEFAULT '',
    "create_time"           int4         NOT NULL DEFAULT 0,
    "update_time"           int4         NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "public"."chat_ai_model_list" ("model_config_id", "model_type", "use_model_name");
CREATE INDEX ON "public"."chat_ai_model_list" ("admin_user_id", "model_config_id");

COMMENT ON TABLE "public"."chat_ai_model_list" IS '可使用的模型列表';
COMMENT ON COLUMN "public"."chat_ai_model_list"."id" IS 'ID';
COMMENT ON COLUMN "public"."chat_ai_model_list"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_model_list"."model_config_id" IS '模型配置ID';
COMMENT ON COLUMN "public"."chat_ai_model_list"."model_type" IS '模型类型';
COMMENT ON COLUMN "public"."chat_ai_model_list"."use_model_name" IS '真实使用的模型(调用时的model参数)';
COMMENT ON COLUMN "public"."chat_ai_model_list"."show_model_name" IS '用于显示的模型名称(仅用于展示)';
COMMENT ON COLUMN "public"."chat_ai_model_list"."thinking_type" IS '深度思考选项:0不支持,1支持,2可选';
COMMENT ON COLUMN "public"."chat_ai_model_list"."function_call" IS '是否支持function call';
COMMENT ON COLUMN "public"."chat_ai_model_list"."input_text" IS '支持的输入类型:文本';
COMMENT ON COLUMN "public"."chat_ai_model_list"."input_voice" IS '支持的输入类型:语音';
COMMENT ON COLUMN "public"."chat_ai_model_list"."input_image" IS '支持的输入类型:图片';
COMMENT ON COLUMN "public"."chat_ai_model_list"."input_video" IS '支持的输入类型:视频';
COMMENT ON COLUMN "public"."chat_ai_model_list"."input_document" IS '支持的输入类型:文档';
COMMENT ON COLUMN "public"."chat_ai_model_list"."output_text" IS '支持的输出类型:文本';
COMMENT ON COLUMN "public"."chat_ai_model_list"."output_voice" IS '支持的输出类型:语音';
COMMENT ON COLUMN "public"."chat_ai_model_list"."output_image" IS '支持的输出类型:图片';
COMMENT ON COLUMN "public"."chat_ai_model_list"."output_video" IS '支持的输出类型:视频';
COMMENT ON COLUMN "public"."chat_ai_model_list"."vector_dimension_list" IS '向量维度列表(英文逗号分割)';
COMMENT ON COLUMN "public"."chat_ai_model_list"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."chat_ai_model_list"."update_time" IS '更新时间';
