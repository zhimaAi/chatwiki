-- +goose Up

ALTER TABLE "public"."chat_ai_model_config"
    ADD COLUMN "config_name" varchar(50) default '';

COMMENT ON COLUMN "public"."chat_ai_model_config"."config_name" IS '配置名称(配置的别名,仅用于展示)';

COMMENT ON COLUMN "public"."chat_ai_model_config"."model_types" IS '模型类型(英文逗号拼接)--禁止编辑';
COMMENT ON COLUMN "public"."chat_ai_model_config"."deployment_name" IS '部署名称(非必填)--已废弃,请勿使用';
COMMENT ON COLUMN "public"."chat_ai_model_config"."show_model_name" IS '用于显示的模型名称--已废弃,请勿使用';
COMMENT ON COLUMN "public"."chat_ai_model_config"."thinking_type" IS '深度思考选项:0不支持,1支持,2可选--已废弃,请勿使用';
