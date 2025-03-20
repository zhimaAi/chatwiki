-- +goose Up

CREATE TABLE "chat_ai_library_question_guide"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "library_id" int4          NOT NULL DEFAULT 0,
    "question"         varchar(1000)  NOT NULL DEFAULT '',
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_library_question_guide" ("library_id");

COMMENT ON TABLE "chat_ai_library_question_guide" IS '知识库引导表';

COMMENT ON COLUMN "chat_ai_library_question_guide"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_library_question_guide"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_library_question_guide"."library_id" IS '知识库id';
COMMENT ON COLUMN "chat_ai_library_question_guide"."question" IS '自定义域名';
COMMENT ON COLUMN "chat_ai_library_question_guide"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_library_question_guide"."update_time" IS '更新时间';
