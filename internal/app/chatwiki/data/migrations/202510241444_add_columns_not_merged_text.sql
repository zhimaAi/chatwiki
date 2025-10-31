-- +goose Up

ALTER TABLE "public"."chat_ai_library_file"
    ADD COLUMN "not_merged_text" bool default false;
COMMENT ON COLUMN "public"."chat_ai_library_file"."not_merged_text" IS '自动合并较小分段开关:0开1关';

ALTER TABLE "public"."chat_ai_library"
    ADD COLUMN "normal_chunk_default_not_merged_text" bool default false;
COMMENT ON COLUMN "public"."chat_ai_library"."normal_chunk_default_not_merged_text" IS '普通分段模式下自动合并较小分段开关:0开1关';