-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "rrf_weight" varchar(100) NOT NULL DEFAULT '';
COMMENT ON COLUMN "public"."chat_ai_robot"."rrf_weight" IS 'RRF算法权重配置 eg:{"vector":50,"search":30,"graph":20}';

ALTER TABLE "public"."user_search_config"
    ADD COLUMN "rrf_weight" varchar(100) NOT NULL DEFAULT '';
COMMENT ON COLUMN "public"."user_search_config"."rrf_weight" IS 'RRF算法权重配置 eg:{"vector":50,"search":30,"graph":20}';
