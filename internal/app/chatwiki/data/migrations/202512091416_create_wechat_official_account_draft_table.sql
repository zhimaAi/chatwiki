-- +goose Up
CREATE TABLE "public"."wechat_official_account_draft"
(
    "id"                serial                                      NOT NULL PRIMARY KEY,
    "admin_user_id"     INT4                                        NOT NULL DEFAULT 0,
    "create_time"       INT4                                        NOT NULL DEFAULT 0,
    "update_time"       INT4                                        NOT NULL DEFAULT 0,
    "app_id"            VARCHAR(30) COLLATE "pg_catalog"."default"  NOT NULL,
    "article_type"      VARCHAR(20) COLLATE "pg_catalog"."default"  NOT NULL,
    "thumb_url"         TEXT COLLATE "pg_catalog"."default"         NOT NULL,
    "title"             VARCHAR(255) COLLATE "pg_catalog"."default" NOT NULL,
    "draft_create_time" INT4                                        NOT NULL,
    "draft_update_time" INT4                                        NOT NULL,
    "media_id"          VARCHAR(191) COLLATE "pg_catalog"."default" NOT NULL,
    "digest"            VARCHAR(255) COLLATE "pg_catalog"."default" NOT NULL,
    "group_id"          INT4                                        NOT NULL DEFAULT 0
);
ALTER TABLE "public"."wechat_official_account_draft" OWNER TO "chatwiki";
CREATE INDEX "wechat_official_account_draft_admin" ON "public"."wechat_official_account_draft" USING btree ( "admin_user_id" "pg_catalog"."int4_ops" ASC NULLS LAST, "app_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST );
CREATE INDEX "wechat_official_account_draft_group_index" ON "public"."wechat_official_account_draft" USING btree ( "group_id" "pg_catalog"."int4_ops" ASC NULLS LAST );
COMMENT ON COLUMN "public"."wechat_official_account_draft"."id" IS 'ID';
COMMENT ON COLUMN "public"."wechat_official_account_draft"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."wechat_official_account_draft"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."wechat_official_account_draft"."update_time" IS '更新时间';
COMMENT ON COLUMN "public"."wechat_official_account_draft"."app_id" IS '公众号ID';
COMMENT ON COLUMN "public"."wechat_official_account_draft"."article_type" IS '文章类型 newspic：图片消息，news：图文消息';
COMMENT ON COLUMN "public"."wechat_official_account_draft"."thumb_url" IS '封面预览图';
COMMENT ON COLUMN "public"."wechat_official_account_draft"."title" IS '文章标题';
COMMENT ON COLUMN "public"."wechat_official_account_draft"."draft_create_time" IS '文章创建时间';
COMMENT ON COLUMN "public"."wechat_official_account_draft"."draft_update_time" IS '文章更新时间';
COMMENT ON COLUMN "public"."wechat_official_account_draft"."media_id" IS '图文消息ID';
COMMENT ON COLUMN "public"."wechat_official_account_draft"."digest" IS '摘要';
COMMENT ON COLUMN "public"."wechat_official_account_draft"."group_id" IS '公众号草稿';