-- +goose Up

ALTER TABLE "chat_ai_library_file" ADD COLUMN "graph_status" int2 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_library_file"."graph_status" IS '知识图谱状态 0未开始 1构建中 2构建完成 3构建失败';
