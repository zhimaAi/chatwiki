-- +goose Up

ALTER TABLE "chat_ai_library_file" ADD COLUMN "chunk_type" int2 NOT NULL DEFAULT 0;
ALTER TABLE "chat_ai_library_file" ADD COLUMN "semantic_chunk_size" int2 NOT NULL DEFAULT 512;
ALTER TABLE "chat_ai_library_file" ADD COLUMN "semantic_chunk_overlap" int4 NOT NULL DEFAULT 50;
ALTER TABLE "chat_ai_library_file" ADD COLUMN "semantic_chunk_threshold" int4 NOT NULL DEFAULT 90;

COMMENT ON COLUMN "chat_ai_library_file"."chunk_type" IS '分段方式 1普通分段 2语义分段';
COMMENT ON COLUMN "chat_ai_library_file"."semantic_chunk_size" IS '语义模式下分段最大长度';
COMMENT ON COLUMN "chat_ai_library_file"."semantic_chunk_overlap" IS '语义模式下分段重叠长度';
COMMENT ON COLUMN "chat_ai_library_file"."semantic_chunk_threshold" IS '语义模式下分段断点阈值';
