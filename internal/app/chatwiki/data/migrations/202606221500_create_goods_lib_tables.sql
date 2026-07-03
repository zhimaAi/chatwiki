-- +goose Up

CREATE TABLE "public"."chat_ai_goods_group"
(
    "id"            bigserial   NOT NULL PRIMARY KEY,
    "admin_user_id" int4        NOT NULL DEFAULT 0,
    "parent_id"     int8        NOT NULL DEFAULT 0,
    "group_name"    varchar(15) NOT NULL DEFAULT '',
    "level"         int2        NOT NULL DEFAULT 1,
    "sort"          int4        NOT NULL DEFAULT 0,
    "create_time"   int4        NOT NULL DEFAULT 0,
    "update_time"   int4        NOT NULL DEFAULT 0
);

CREATE INDEX "idx_chat_ai_goods_group_admin_parent_sort"
    ON "public"."chat_ai_goods_group" ("admin_user_id", "parent_id", "sort", "id");

COMMENT ON TABLE "public"."chat_ai_goods_group" IS '商品库分组';
COMMENT ON COLUMN "public"."chat_ai_goods_group"."id" IS '商品分组ID';
COMMENT ON COLUMN "public"."chat_ai_goods_group"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_goods_group"."parent_id" IS '上级分组ID，0表示顶级分组';
COMMENT ON COLUMN "public"."chat_ai_goods_group"."group_name" IS '分组名称';
COMMENT ON COLUMN "public"."chat_ai_goods_group"."level" IS '分组层级，最多5级';
COMMENT ON COLUMN "public"."chat_ai_goods_group"."sort" IS '同级排序值，越小越靠前';
COMMENT ON COLUMN "public"."chat_ai_goods_group"."create_time" IS '创建时间，Unix秒级时间戳';
COMMENT ON COLUMN "public"."chat_ai_goods_group"."update_time" IS '更新时间，Unix秒级时间戳';

CREATE TABLE "public"."chat_ai_goods_library"
(
    "id"            bigserial      NOT NULL PRIMARY KEY,
    "admin_user_id" int4           NOT NULL DEFAULT 0,
    "group_id"      int8           NOT NULL DEFAULT 0,
    "goods_id"      varchar(100)   NOT NULL DEFAULT '',
    "goods_name"    varchar(100)   NOT NULL DEFAULT '',
    "category"      varchar(100)   NOT NULL DEFAULT '',
    "brand"         varchar(100)   NOT NULL DEFAULT '',
    "price"         decimal(20, 2) NOT NULL DEFAULT 0,
    "stock"         int8           NOT NULL DEFAULT 0,
    "link"          varchar(1000)  NOT NULL DEFAULT '',
    "images"        jsonb          NOT NULL DEFAULT '[]',
    "description"   text           NOT NULL DEFAULT '',
    "qa"            text           NOT NULL DEFAULT '',
    "custom_info"   text           NOT NULL DEFAULT '',
    "switch_status" int2           NOT NULL DEFAULT 1,
    "create_time"   int4           NOT NULL DEFAULT 0,
    "update_time"   int4           NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX "uniq_chat_ai_goods_library_admin_goods_id"
    ON "public"."chat_ai_goods_library" ("admin_user_id", "goods_id");
CREATE INDEX "idx_chat_ai_goods_library_admin_group"
    ON "public"."chat_ai_goods_library" ("admin_user_id", "group_id", "id");
CREATE INDEX "idx_chat_ai_goods_library_admin_switch"
    ON "public"."chat_ai_goods_library" ("admin_user_id", "switch_status", "id");
CREATE INDEX "idx_chat_ai_goods_library_goods_id_bigm"
    ON "public"."chat_ai_goods_library" USING gin ("goods_id" gin_bigm_ops);
CREATE INDEX "idx_chat_ai_goods_library_goods_name_bigm"
    ON "public"."chat_ai_goods_library" USING gin ("goods_name" gin_bigm_ops);

COMMENT ON TABLE "public"."chat_ai_goods_library" IS '商品库商品信息';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."id" IS '商品库商品ID';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."group_id" IS '商品分组ID，0表示未分组';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."goods_id" IS '客户填写的商品ID，管理员维度唯一';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."goods_name" IS '商品名称';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."category" IS '客户填写的类目文本';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."brand" IS '客户填写的品牌文本';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."price" IS '商品价格，保留2位小数';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."stock" IS '客户填写的库存文本';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."link" IS '客户填写的商品链接文本';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."images" IS '商品图片URL数组，最多5张';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."description" IS '商品描述';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."qa" IS '商品问答';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."custom_info" IS '自定义信息';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."switch_status" IS '商品启用状态：0关闭，1启用';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."create_time" IS '创建时间，Unix秒级时间戳';
COMMENT ON COLUMN "public"."chat_ai_goods_library"."update_time" IS '更新时间，Unix秒级时间戳';

-- +goose Down

DROP TABLE IF EXISTS "public"."chat_ai_goods_library";
DROP TABLE IF EXISTS "public"."chat_ai_goods_group";
