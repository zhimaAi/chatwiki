-- +goose Up


ALTER TABLE "public"."chat_ai_model_list" ADD COLUMN "image_generation" varchar(500) NOT NULL DEFAULT '{}';

COMMENT ON COLUMN "public"."chat_ai_model_list"."image_generation" IS '图片模型配置参数';