-- +goose Up

ALTER TABLE "chat_ai_library_file_data" ADD COLUMN "graph_status" int2 NOT NULL DEFAULT 0;
ALTER TABLE "chat_ai_library_file_data" ADD COLUMN "graph_err_msg" varchar(1024) NOT NULL DEFAULT '';
COMMENT ON COLUMN "chat_ai_library_file_data"."graph_status" IS '知识图谱状态 0未开始 1构建中 2构建完成 3构建失败';
COMMENT ON COLUMN "chat_ai_library_file_data"."graph_err_msg" IS '构建错误原因';
