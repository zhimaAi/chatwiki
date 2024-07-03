-- +goose Up

CREATE TABLE "chat_ai_answer_source"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "message_id"    int4          NOT NULL DEFAULT 0,
    "file_id"       int4          NOT NULL DEFAULT 0,
    "paragraph_id"  int4          NOT NULL DEFAULT 0,
    "word_total"    int4          NOT NULL DEFAULT 0,
    "similarity"    varchar(20)   NOT NULL DEFAULT '',
    "content"       varchar(5000) NOT NULL DEFAULT '',
    "title"         varchar(100)  NOT NULL DEFAULT '',
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_answer_source" ("admin_user_id");
CREATE INDEX ON "chat_ai_answer_source" ("message_id", "file_id");

COMMENT ON TABLE "chat_ai_answer_source" IS '文档问答机器人-答案来源';

COMMENT ON COLUMN "chat_ai_answer_source"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_answer_source"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_answer_source"."message_id" IS '聊天消息ID';
COMMENT ON COLUMN "chat_ai_answer_source"."file_id" IS '文件ID';
COMMENT ON COLUMN "chat_ai_answer_source"."paragraph_id" IS '分段ID';
COMMENT ON COLUMN "chat_ai_answer_source"."word_total" IS '单词数量';
COMMENT ON COLUMN "chat_ai_answer_source"."similarity" IS '相似度';
COMMENT ON COLUMN "chat_ai_answer_source"."content" IS '文档段落';
COMMENT ON COLUMN "chat_ai_answer_source"."title" IS '分段标题';
COMMENT ON COLUMN "chat_ai_answer_source"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_answer_source"."update_time" IS '更新时间';
