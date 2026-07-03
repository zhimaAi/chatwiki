-- +goose Up

CREATE TABLE IF NOT EXISTS "public"."admin_mini_card"
(
    "id"            serial8      NOT NULL,
    "admin_user_id" int8         NOT NULL DEFAULT 0,
    "title"         varchar(255) NOT NULL DEFAULT '',
    "appid"         varchar(100) NOT NULL DEFAULT '',
    "page_path"     varchar(500) NOT NULL DEFAULT '',
    "thumb_url"     varchar(500) NOT NULL DEFAULT '',
    "create_time"   int8         NOT NULL DEFAULT 0,
    "update_time"   int8         NOT NULL DEFAULT 0,
    "delete_time"   int8         NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE INDEX IF NOT EXISTS "idx_admin_mini_card_admin_delete_id"
    ON "public"."admin_mini_card" ("admin_user_id", "delete_time", "id");
CREATE INDEX IF NOT EXISTS "idx_admin_mini_card_admin_appid_delete"
    ON "public"."admin_mini_card" ("admin_user_id", "appid", "delete_time");

COMMENT ON TABLE "public"."admin_mini_card" IS '管理员小程序卡片';
COMMENT ON COLUMN "public"."admin_mini_card"."id" IS '自增ID';
COMMENT ON COLUMN "public"."admin_mini_card"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."admin_mini_card"."title" IS '小程序卡片标题';
COMMENT ON COLUMN "public"."admin_mini_card"."appid" IS '小程序appid';
COMMENT ON COLUMN "public"."admin_mini_card"."page_path" IS '小程序页面路径';
COMMENT ON COLUMN "public"."admin_mini_card"."thumb_url" IS '小程序卡片封面';
COMMENT ON COLUMN "public"."admin_mini_card"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."admin_mini_card"."update_time" IS '更新时间';
COMMENT ON COLUMN "public"."admin_mini_card"."delete_time" IS '删除时间';

CREATE TABLE IF NOT EXISTS "public"."admin_mini_card_relation"
(
    "id"            serial8     NOT NULL,
    "admin_user_id" int8        NOT NULL DEFAULT 0,
    "mini_card_id"  int8        NOT NULL DEFAULT 0,
    "robot_id"      int8        NOT NULL DEFAULT 0,
    "library_id"    int8        NOT NULL DEFAULT 0,
    "target_type"   varchar(32) NOT NULL DEFAULT '',
    "target_id"     int8        NOT NULL DEFAULT 0,
    "create_time"   int8        NOT NULL DEFAULT 0,
    "update_time"   int8        NOT NULL DEFAULT 0,
    "delete_time"   int8        NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS "uniq_admin_mini_card_relation_target"
    ON "public"."admin_mini_card_relation" ("admin_user_id", "target_type", "target_id")
    WHERE "delete_time" = 0;
CREATE INDEX IF NOT EXISTS "idx_admin_mini_card_relation_card"
    ON "public"."admin_mini_card_relation" ("admin_user_id", "mini_card_id", "delete_time");
CREATE INDEX IF NOT EXISTS "idx_admin_mini_card_relation_library"
    ON "public"."admin_mini_card_relation" ("admin_user_id", "library_id", "target_type", "delete_time");

COMMENT ON TABLE "public"."admin_mini_card_relation" IS '管理员小程序卡片关联关系';
COMMENT ON COLUMN "public"."admin_mini_card_relation"."id" IS '自增ID';
COMMENT ON COLUMN "public"."admin_mini_card_relation"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."admin_mini_card_relation"."mini_card_id" IS '小程序卡片ID';
COMMENT ON COLUMN "public"."admin_mini_card_relation"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "public"."admin_mini_card_relation"."library_id" IS '知识库ID';
COMMENT ON COLUMN "public"."admin_mini_card_relation"."target_type" IS '关联对象类型:library_qa问答知识库,library_paragraph文档段落,robot_prompt大模型自定义提示词';
COMMENT ON COLUMN "public"."admin_mini_card_relation"."target_id" IS '关联对象ID';
COMMENT ON COLUMN "public"."admin_mini_card_relation"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."admin_mini_card_relation"."update_time" IS '更新时间';
COMMENT ON COLUMN "public"."admin_mini_card_relation"."delete_time" IS '删除时间';
