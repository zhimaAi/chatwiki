-- +goose Up

ALTER TABLE "chat_ai_library_file_data" ADD COLUMN "isolated" bool NOT NULL DEFAULT false;

COMMENT ON COLUMN "chat_ai_library_file_data"."isolated" IS '段落被标记后,对应的文档被删除或重新分段,该段落就被孤立了';