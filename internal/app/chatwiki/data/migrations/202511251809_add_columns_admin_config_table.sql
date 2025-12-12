-- +goose Up
ALTER TABLE "public"."admin_user_config"
    ADD COLUMN "comment_pull_days" INT2 NOT NULL DEFAULT 7,
    ADD COLUMN "comment_pull_limit" INT2 NOT NULL DEFAULT 10;
COMMENT ON COLUMN "public"."admin_user_config"."comment_pull_days" IS '微信公众号评论拉取时效，默认7天';
COMMENT ON COLUMN "public"."admin_user_config"."comment_pull_limit" IS '微信公众号评论拉取间隔，默认拉取结束10分钟后进行下一轮';