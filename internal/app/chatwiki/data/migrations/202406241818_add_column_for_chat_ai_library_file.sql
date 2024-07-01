-- +goose Up

ALTER TABLE "chat_ai_library_file"
    ADD COLUMN "doc_auto_renew_frequency" int2 NOT NULL DEFAULT 1,
    ADD COLUMN "doc_last_renew_time" int4 NOT NULL DEFAULT 0,
    ADD COLUMN "remark" varchar(255) NOT NULL DEFAULT '';
COMMENT ON COLUMN "chat_ai_library_file"."doc_auto_renew_frequency" IS '在线文档自动更新频率 1不自动更新 2每天 3每3天 4每7天 5每30天';
COMMENT ON COLUMN "chat_ai_library_file"."doc_last_renew_time" IS '在线文档最后一次更新时间';
COMMENT ON COLUMN "chat_ai_library_file"."remark" IS '备注';