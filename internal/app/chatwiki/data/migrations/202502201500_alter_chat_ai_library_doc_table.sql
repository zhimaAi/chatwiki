-- +goose Up

ALTER TABLE "chat_ai_library_file_doc"
    ADD COLUMN "edit_user"  int4 NOT NULL DEFAULT 0,
    ADD COLUMN "seo_title"  varchar(300) NOT NULL DEFAULT '',
    ADD COLUMN "seo_desc"  varchar(1000) NOT NULL DEFAULT '',
    ADD COLUMN "seo_keywords"  varchar(1000) NOT NULL DEFAULT '';


