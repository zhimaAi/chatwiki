-- +goose Up

ALTER TABLE "public"."chat_ai_library_file_data_index"
    ADD COLUMN "embedding2000" vector(2000);

COMMENT ON COLUMN "public"."chat_ai_library_file_data_index"."embedding2000" IS '固定2000维度向量的文档';

CREATE INDEX ON "public"."chat_ai_library_file_data_index" USING ivfflat ("embedding2000" vector_cosine_ops) WITH (lists = 100);