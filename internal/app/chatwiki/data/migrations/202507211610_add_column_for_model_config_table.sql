-- +goose Up

ALTER TABLE "public"."chat_ai_model_config"
    ADD COLUMN "thinking_type" int2 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_model_config"."thinking_type" IS '深度思考选项:0不支持,1支持,2可选';
