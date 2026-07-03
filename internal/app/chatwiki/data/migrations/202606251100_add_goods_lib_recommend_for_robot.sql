-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "goods_lib_recommend_switch"    int2           NOT NULL DEFAULT 0,
    ADD COLUMN "goods_lib_recommend_group_ids" varchar(10000) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_robot"."goods_lib_recommend_switch" IS '商品库推荐开关：0关闭，1开启';
COMMENT ON COLUMN "public"."chat_ai_robot"."goods_lib_recommend_group_ids" IS '商品库推荐范围，商品分组ID集合，空表示全部分组';

-- +goose Down

ALTER TABLE "public"."chat_ai_robot"
    DROP COLUMN IF EXISTS "goods_lib_recommend_switch",
    DROP COLUMN IF EXISTS "goods_lib_recommend_group_ids";
