-- +goose Up

CREATE TABLE "chat_ai_user_domain"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "url"         varchar(200)  NOT NULL DEFAULT '',
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_user_domain" ("admin_user_id","url");

COMMENT ON TABLE "chat_ai_user_domain" IS '自定义域名';

COMMENT ON COLUMN "chat_ai_user_domain"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_user_domain"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_user_domain"."url" IS '自定义域名';
COMMENT ON COLUMN "chat_ai_user_domain"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_user_domain"."update_time" IS '更新时间';

CREATE TABLE "chat_ai_file_info"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "file_name"         varchar(500)  NOT NULL DEFAULT '',
    "file_content"         varchar(2000)  NOT NULL DEFAULT '',
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_file_info" ("admin_user_id");

COMMENT ON TABLE "chat_ai_file_info" IS '自定义域名';

COMMENT ON COLUMN "chat_ai_file_info"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_file_info"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_file_info"."file_name" IS '文件名';
COMMENT ON COLUMN "chat_ai_file_info"."file_content" IS '文件内容';
COMMENT ON COLUMN "chat_ai_file_info"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_file_info"."update_time" IS '更新时间';
