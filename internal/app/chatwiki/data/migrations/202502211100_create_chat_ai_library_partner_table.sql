-- +goose Up

CREATE TABLE "chat_ai_library_partner"
(
    "id"            serial        NOT NULL primary key,
    "user_id" int4          NOT NULL DEFAULT 0,
    "library_id" int4          NOT NULL DEFAULT 0,
    "operate_rights"      int2          NOT NULL DEFAULT 0,
    "creator" int4          NOT NULL DEFAULT 0,
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_library_partner" ("user_id");
CREATE INDEX ON "chat_ai_library_partner" ("library_id");

COMMENT ON TABLE "chat_ai_library_partner" IS '文档问答机器人-知识库协作者';

COMMENT ON COLUMN "chat_ai_library_partner"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_library_partner"."user_id" IS '用户ID';
COMMENT ON COLUMN "chat_ai_library_partner"."library_id" IS '知识库ID';
COMMENT ON COLUMN "chat_ai_library_partner"."operate_rights" IS '权限 0:无权限 2：编辑 4：管理';
COMMENT ON COLUMN "chat_ai_library_partner"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_library_partner"."update_time" IS '更新时间';

