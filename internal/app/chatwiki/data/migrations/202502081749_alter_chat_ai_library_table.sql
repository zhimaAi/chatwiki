-- +goose Up

ALTER TABLE "chat_ai_library"
    ADD COLUMN "ai_summary"  int2 NOT NULL DEFAULT 0,
    ADD COLUMN "ai_summary_model"  varchar(100) NOT NULL DEFAULT '',
    ADD COLUMN "share_url"  varchar(200) NOT NULL DEFAULT '',
    ADD COLUMN "type" int2          NOT NULL DEFAULT 0,
    ADD COLUMN "creator"  int4 NOT NULL DEFAULT 0,
    ADD COLUMN "access_rights" int2          NOT NULL DEFAULT 0;

COMMENT ON COLUMN "chat_ai_library"."ai_summary" IS 'ai总结 0:不开启 1：开启';
COMMENT ON COLUMN "chat_ai_library"."ai_summary_model" IS 'ai总结model';
COMMENT ON COLUMN "chat_ai_library"."share_url" IS '分享链接';
COMMENT ON COLUMN "chat_ai_library"."access_rights" IS '访问权限 0：私有，1：公开 ';
COMMENT ON COLUMN "chat_ai_library"."type" IS '知识库类型 0:默认 1：对外';
COMMENT ON COLUMN "chat_ai_library"."creator" IS '创建人';

