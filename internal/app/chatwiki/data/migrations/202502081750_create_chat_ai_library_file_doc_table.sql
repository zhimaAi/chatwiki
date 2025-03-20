-- +goose Up

CREATE TABLE "chat_ai_library_file_doc"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "library_id"    int4          NOT NULL DEFAULT 0,
    "file_id"       int4          NOT NULL DEFAULT 0,
    "doc_key"         varchar(50)  NOT NULL DEFAULT '',
    "title"         varchar(100)  NOT NULL DEFAULT '',
    "pid"           int4          NOT NULL DEFAULT 0,
    "sort"   int8          NOT NULL DEFAULT 0,
    "summary"       varchar(5000) NOT NULL DEFAULT '',
    "summary_embedding"     vector(2000),
    "content"       text          NOT NULL DEFAULT '',
    "is_draft"      int2          NOT NULL DEFAULT 0,
    "is_index"      int2          NOT NULL DEFAULT 0,
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_library_file_doc" ("admin_user_id");
CREATE INDEX ON "chat_ai_library_file_doc" ("library_id");
CREATE INDEX ON "chat_ai_library_file_doc" ("pid");
CREATE INDEX ON "chat_ai_library_file_doc" ("file_id");
CREATE INDEX ON "chat_ai_library_file_doc" ("doc_key");

COMMENT ON TABLE "chat_ai_library_file_doc" IS '文档问答机器人-对外知识库文档';

COMMENT ON COLUMN "chat_ai_library_file_doc"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_library_file_doc"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_library_file_doc"."library_id" IS '知识库ID';
COMMENT ON COLUMN "chat_ai_library_file_doc"."file_id" IS '文件id';
COMMENT ON COLUMN "chat_ai_library_file_doc"."pid" IS '父文档id';
COMMENT ON COLUMN "chat_ai_library_file_doc"."sort" IS '排序';
COMMENT ON COLUMN "chat_ai_library_file_doc"."title" IS '标题';
COMMENT ON COLUMN "chat_ai_library_file_doc"."is_draft" IS '草稿 0:否 1：是';
COMMENT ON COLUMN "chat_ai_library_file_doc"."is_index" IS '是否首页 0:否 1：是';
COMMENT ON COLUMN "chat_ai_library_file_doc"."summary" IS 'ai总结';
COMMENT ON COLUMN "chat_ai_library_file_doc"."summary_embedding" IS 'ai总结嵌入';
COMMENT ON COLUMN "chat_ai_library_file_doc"."content" IS '文档内容';
COMMENT ON COLUMN "chat_ai_library_file_doc"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_library_file_doc"."update_time" IS '更新时间';
