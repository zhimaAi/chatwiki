-- +goose Up

create table "public"."wechat_official_article_sync_task"
(
    "id"                serial      not null primary key,
    "create_time"       int4        NOT NULL DEFAULT 0,
    "update_time"       int4        NOT NULL DEFAULT 0,
    "admin_user_id"     int4        not null default 0,
    "app_id"            varchar(50) not null default '',
    "sync_time"         varchar(50) not null default '',
    "ai_comment_switch" int4        not null default 0,
    "last_sync_time"    int4        not null default 0,
    "auto_sync_switch"  int4        not null default 0
);

CREATE INDEX on "wechat_official_article_sync_task"("admin_user_id");
CREATE INDEX on "wechat_official_article_sync_task"("app_id");

COMMENT ON TABLE "public"."wechat_official_article_sync_task" IS '公众号文章自动同步任务';
COMMENT ON COLUMN "public"."wechat_official_article_sync_task"."id" IS 'ID';
COMMENT ON COLUMN "public"."wechat_official_article_sync_task"."admin_user_id" IS '管理员ID';
COMMENT ON COLUMN "public"."wechat_official_article_sync_task"."app_id" IS '应用ID';
COMMENT ON COLUMN "public"."wechat_official_article_sync_task"."sync_time" IS '同步时间';
COMMENT ON COLUMN "public"."wechat_official_article_sync_task"."ai_comment_switch" IS 'AI评论开关，0:关闭，1:开启';
COMMENT ON COLUMN "public"."wechat_official_article_sync_task"."last_sync_time" IS '上次同步时间';
COMMENT ON COLUMN "public"."wechat_official_article_sync_task"."auto_sync_switch" IS '自动同步开关，0:关闭，1:开启';
