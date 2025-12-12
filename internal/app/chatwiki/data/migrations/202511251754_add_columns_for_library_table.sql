-- +goose Up

ALTER TABLE "public"."chat_ai_library" ADD COLUMN "official_app_id" varchar(50) NOT NULL DEFAULT '';
ALTER TABLE "public"."chat_ai_library" ADD COLUMN "sync_official_history_type" int4 NOT NULL DEFAULT 3;
ALTER TABLE "public"."chat_ai_library" ADD COLUMN "enable_cron_sync_official_content" boolean NOT NULL DEFAULT false;
ALTER TABLE "public"."chat_ai_library" ADD COLUMN "sync_official_content_status" int2 NOT NULL DEFAULT 1;
ALTER TABLE "public"."chat_ai_library" ADD COLUMN "sync_official_content_last_err_msg" varchar(1000) NOT NULL DEFAULT '';


COMMENT ON COLUMN "public"."chat_ai_library"."type" IS '知识库类型:0普通知识库,1对外知识库,2问答知识库,3公众号知识库';
COMMENT ON COLUMN "public"."chat_ai_library"."official_app_id" IS '公众号知识库关联的app_id';
COMMENT ON COLUMN "public"."chat_ai_library"."sync_official_history_type" IS '公众号知识库获取内容的时间段 1半年内 2一年内 3三年以内 10全部';
COMMENT ON COLUMN "public"."chat_ai_library"."enable_cron_sync_official_content" IS '公众号知识库每天定时获取最新发布内容开关';
COMMENT ON COLUMN "public"."chat_ai_library"."sync_official_content_status" IS '同步公众号知识库状态 1未同步 2同步中 3同步失败';
COMMENT ON COLUMN "public"."chat_ai_library"."sync_official_content_last_err_msg" IS '同步公众号知识库时上一次同步失败的原因';

ALTER TABLE "public"."chat_ai_library_file" ADD COLUMN "official_article_id" varchar(100) NOT NULL DEFAULT '';
ALTER TABLE "public"."chat_ai_library_file" ADD COLUMN "official_article_update_time" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_library_file"."official_article_id" IS '公众号文章id';
COMMENT ON COLUMN "public"."chat_ai_library_file"."official_article_update_time" IS '公众号文章更新时间';

