-- +goose Up
ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "last_edit_ip" varchar(100) NOT NULL DEFAULT '',
  ADD COLUMN "last_edit_user_agent" varchar(1000) NOT NULL DEFAULT '',
    ADD COLUMN "draft_save_type" varchar(20) NOT NULL DEFAULT '',
ADD COLUMN "draft_save_time" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_robot"."last_edit_ip" IS '上次编辑人ip';

COMMENT ON COLUMN "public"."chat_ai_robot"."last_edit_user_agent" IS '上次编辑人user_agent';

COMMENT ON COLUMN "public"."chat_ai_robot"."draft_save_type" IS '草稿保存类型，automatic：自动，handle：手动';

COMMENT ON COLUMN "public"."chat_ai_robot"."draft_save_time" IS '草稿保存时间';

ALTER TABLE "public"."work_flow_version"
    ADD COLUMN "last_edit_ip" varchar(100) NOT NULL DEFAULT '',
  ADD COLUMN "last_edit_user_agent" varchar(1000) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."work_flow_version"."last_edit_ip" IS '上次编辑人ip';

COMMENT ON COLUMN "public"."work_flow_version"."last_edit_user_agent" IS '上次编辑人user_agent';

CREATE TABLE "public"."admin_user_config"
(
    "id"            serial NOT NULL PRIMARY KEY,
    "admin_user_id" int4   NOT NULL DEFAULT 0,
    "create_time"   int4   NOT NULL DEFAULT 0,
    "update_time"   int4   NOT NULL DEFAULT 0,
    "draft_exptime" int4   NOT NULL DEFAULT 10
)
;

ALTER TABLE "public"."admin_user_config" OWNER TO "chatwiki";

CREATE INDEX "admin_user_config_idx" ON "public"."admin_user_config" USING btree (
    "admin_user_id" "pg_catalog"."int4_ops" ASC NULLS LAST
    );

COMMENT
ON COLUMN "public"."admin_user_config"."id" IS 'ID';

COMMENT
ON COLUMN "public"."admin_user_config"."admin_user_id" IS '管理员ID';

COMMENT
ON COLUMN "public"."admin_user_config"."create_time" IS '创建时间';

COMMENT
ON COLUMN "public"."admin_user_config"."update_time" IS '更新时间';

COMMENT
ON COLUMN "public"."admin_user_config"."draft_exptime" IS '草稿箱编辑锁过期时间';

COMMENT
ON TABLE "public"."admin_user_config" IS '系统配置';

