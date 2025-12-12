-- +goose Up
CREATE TABLE "public"."wechat_official_account_draft_group"
(
    "id"            serial                                     NOT NULL PRIMARY KEY,
    "admin_user_id" INT4                                       NOT NULL DEFAULT 0,
    "create_time"   INT4                                       NOT NULL DEFAULT 0,
    "update_time"   INT4                                       NOT NULL DEFAULT 0,
    "group_name"    VARCHAR(30) COLLATE "pg_catalog"."default" NOT NULL
);
ALTER TABLE "public"."wechat_official_account_draft_group" OWNER TO "chatwiki";
CREATE INDEX "wechat_official_account_draft_group_admin" ON "public"."wechat_official_account_draft_group" USING btree ( "admin_user_id" "pg_catalog"."int4_ops" ASC NULLS LAST );
COMMENT ON COLUMN "public"."wechat_official_account_draft_group"."id" IS 'ID';
COMMENT ON COLUMN "public"."wechat_official_account_draft_group"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."wechat_official_account_draft_group"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."wechat_official_account_draft_group"."update_time" IS '更新时间';
COMMENT ON COLUMN "public"."wechat_official_account_draft_group"."group_name" IS '分组名';
COMMENT ON TABLE "public"."wechat_official_account_draft_group" IS '公众号草稿箱分组';