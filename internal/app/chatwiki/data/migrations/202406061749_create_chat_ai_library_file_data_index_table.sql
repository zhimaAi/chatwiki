-- +goose Up

CREATE TABLE "chat_ai_library_file_data_index"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULl DEFAULT 0,
    "library_id"    int4          NOT NULL DEFAULT 0,
    "file_id"       int4          NOT NULL DEFAULT 0,
    "data_id"       int4          NOT NULL DEFAULT 0,
    "type"          int4          NOT NULL DEFAULT 1,
    "content"       varchar(5000) NOT NULL DEFAULT '',
    "status"        int2          NOT NULL DEFAULT 0,
    "errmsg"        varchar(1000) NOT NULL DEFAULT '',
    "embedding"     vector(2000),
    "total_tokens"  int4          NOT NULL DEFAULT 0,
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_library_file_data_index" ("data_id");
CREATE INDEX ON "chat_ai_library_file_data_index" ("type");
CREATE INDEX ON "chat_ai_library_file_data_index" ("content");
CREATE INDEX ON "chat_ai_library_file_data_index" USING ivfflat ("embedding" vector_cosine_ops) WITH (lists = 100);
CREATE INDEX ON "chat_ai_library_file_data_index" using gin (to_tsvector('zhima_zh_parser', upper("content")));

COMMENT ON TABLE "chat_ai_library_file_data_index" IS '文档问答机器人-知识库数据索引';

COMMENT ON COLUMN "chat_ai_library_file_data_index"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_library_file_data_index"."id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_library_file_data_index"."library_id" IS '知识库ID';
COMMENT ON COLUMN "chat_ai_library_file_data_index"."file_id" IS '文件ID';
COMMENT ON COLUMN "chat_ai_library_file_data_index"."data_id" IS '关联知识库数据id';
COMMENT ON COLUMN "chat_ai_library_file_data_index"."type" IS '索引类型 1文档段落索引 2问题索引 3答案索引 4自定义索引';
COMMENT ON COLUMN "chat_ai_library_file_data_index"."content" IS '需要索引的数据';
COMMENT ON COLUMN "chat_ai_library_file_data_index"."status" IS '状态:0未转换,1已转换,2转换异常';
COMMENT ON COLUMN "chat_ai_library_file_data_index"."errmsg" IS '转换报错信息';
COMMENT ON COLUMN "chat_ai_library_file_data_index"."embedding" IS '转换为向量的文档';
COMMENT ON COLUMN "chat_ai_library_file_data_index"."total_tokens" IS '消耗的token';
COMMENT ON COLUMN "chat_ai_library_file_data_index"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_library_file_data_index"."update_time" IS '更新时间';


