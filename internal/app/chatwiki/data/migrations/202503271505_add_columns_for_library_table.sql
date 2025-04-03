-- +goose Up

ALTER TABLE "chat_ai_library" ADD COLUMN "chunk_type" int2 NOT NULL DEFAULT 1;
ALTER TABLE "chat_ai_library" ADD COLUMN "normal_chunk_default_separators_no" varchar(100) NOT NULL DEFAULT '11,12';
ALTER TABLE "chat_ai_library" ADD COLUMN "normal_chunk_default_chunk_size" int4 NOT NULL DEFAULT 512;
ALTER TABLE "chat_ai_library" ADD COLUMN "normal_chunk_default_chunk_overlap" int4 NOT NULL DEFAULT 50;
ALTER TABLE "chat_ai_library" ADD COLUMN "semantic_chunk_default_chunk_size" int4 NOT NULL DEFAULT 512;
ALTER TABLE "chat_ai_library" ADD COLUMN "semantic_chunk_default_chunk_overlap" int4 NOT NULL DEFAULT 50;
ALTER TABLE "chat_ai_library" ADD COLUMN "semantic_chunk_default_threshold" int2 NOT NULL DEFAULT 90;

COMMENT ON COLUMN "chat_ai_library"."chunk_type" IS '分段方式 1普通分段 2语义分段';
COMMENT ON COLUMN "chat_ai_library"."normal_chunk_default_separators_no" IS '普通分段模式下分隔符序号集';
COMMENT ON COLUMN "chat_ai_library"."normal_chunk_default_chunk_size" IS '普通分段模式下分段最大长度';
COMMENT ON COLUMN "chat_ai_library"."normal_chunk_default_chunk_overlap" IS '普通分段模式下分段重叠长度';
COMMENT ON COLUMN "chat_ai_library"."semantic_chunk_default_chunk_size" IS '语义模式下分段最大长度';
COMMENT ON COLUMN "chat_ai_library"."semantic_chunk_default_chunk_overlap" IS '语义分段模式下分段重叠长度';
COMMENT ON COLUMN "chat_ai_library"."semantic_chunk_default_threshold" IS '语义分段模式下默认断点阈值';
