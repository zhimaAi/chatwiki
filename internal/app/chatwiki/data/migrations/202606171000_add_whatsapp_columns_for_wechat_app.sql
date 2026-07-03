-- +goose Up

ALTER TABLE "public"."chat_ai_wechat_app"
    ADD COLUMN IF NOT EXISTS "cams_access_key_id"  varchar(100) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS "cams_cust_space_id"  varchar(100) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_wechat_app"."cams_access_key_id" IS '阿里云ChatApp AccessKeyId(WhatsApp)';
COMMENT ON COLUMN "public"."chat_ai_wechat_app"."cams_cust_space_id" IS '阿里云ChatApp 通道ID/CustSpaceId(WhatsApp)';

-- WhatsApp 渠道按「通道(CustSpaceId)」级共用一个回调地址:同一通道下的多个业务号码
-- 派生出相同的 access_key(= 同一个 push_url),回调进来再用消息的 To 精确路由到具体号码记录。
-- 因此 access_key 不再是「一号码一条」,需要去掉建表时(202406061751)的唯一索引,改为普通索引。
DROP INDEX IF EXISTS "chat_ai_wechat_app_access_key_idx";
CREATE INDEX IF NOT EXISTS "idx_chat_ai_wechat_app_access_key" ON "public"."chat_ai_wechat_app" ("access_key");
-- 回调路由 / 保存时按通道复用 access_key 都会按 cust_space_id 过滤,加索引加速。
CREATE INDEX IF NOT EXISTS "idx_chat_ai_wechat_app_cust_space_id" ON "public"."chat_ai_wechat_app" ("cams_cust_space_id");

-- +goose Down

DROP INDEX IF EXISTS "idx_chat_ai_wechat_app_cust_space_id";
DROP INDEX IF EXISTS "idx_chat_ai_wechat_app_access_key";
CREATE UNIQUE INDEX IF NOT EXISTS "chat_ai_wechat_app_access_key_idx" ON "public"."chat_ai_wechat_app" ("access_key");

ALTER TABLE "public"."chat_ai_wechat_app"
    DROP COLUMN IF EXISTS "cams_access_key_id",
    DROP COLUMN IF EXISTS "cams_cust_space_id";
