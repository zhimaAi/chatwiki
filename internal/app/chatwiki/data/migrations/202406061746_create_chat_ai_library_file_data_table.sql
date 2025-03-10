-- +goose Up

CREATE TABLE "chat_ai_library_file_data"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "library_id"    int4          NOT NULL DEFAULT 0,
    "file_id"       int4          NOT NULL DEFAULT 0,
    "number"        int4          NOT NULL DEFAULT 0,
    "page_num"      int4          NOT NULL DEFAULT 0,
    "type"          int4          NOT NULL DEFAULT 0,
    "title"         varchar(100)  NOT NULL DEFAULT '',
    "content"       varchar(5000) NOT NULL DEFAULT '',
    "question"      varchar(5000) NOT NULL DEFAULT '',
    "answer"        varchar(5000) NOT NULL DEFAULT '',
    "word_total"    int4          NOT NULL DEFAULT 0,
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_library_file_data" ("admin_user_id");
CREATE INDEX ON "chat_ai_library_file_data" ("library_id");
CREATE INDEX ON "chat_ai_library_file_data" ("file_id");

COMMENT ON TABLE "chat_ai_library_file_data" IS '文档问答机器人-知识库数据';

COMMENT ON COLUMN "chat_ai_library_file_data"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_library_file_data"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_library_file_data"."library_id" IS '知识库ID';
COMMENT ON COLUMN "chat_ai_library_file_data"."file_id" IS '文件ID';
COMMENT ON COLUMN "chat_ai_library_file_data"."number" IS '文档编号';
COMMENT ON COLUMN "chat_ai_library_file_data"."page_num" IS '文档页码';
COMMENT ON COLUMN "chat_ai_library_file_data"."type" IS '数据类型 1普通段落 2文档问答 3excel问答';
COMMENT ON COLUMN "chat_ai_library_file_data"."title" IS '分段标题';
COMMENT ON COLUMN "chat_ai_library_file_data"."word_total" IS '单词数量';
COMMENT ON COLUMN "chat_ai_library_file_data"."content" IS '文档段落';
COMMENT ON COLUMN "chat_ai_library_file_data"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_library_file_data"."update_time" IS '更新时间';
