-- +goose Up

CREATE TABLE "chatclaw_token_log"
(
    "id"           bigserial    NOT NULL primary key,
    "user_id"      int4         NOT NULL DEFAULT 0,
    "token"        text         NOT NULL DEFAULT '',
    "token_hash"   varchar(64)  NOT NULL DEFAULT '',
    "os_type"      varchar(64)  NOT NULL DEFAULT '',
    "os_version"   varchar(128) NOT NULL DEFAULT '',
    "client_ip"    varchar(64)  NOT NULL DEFAULT '',
    "expired_at"   int8         NOT NULL DEFAULT 0,
    "status"       int2         NOT NULL DEFAULT 1,
    "revoke_time"  int8         NOT NULL DEFAULT 0,
    "revoke_by"    int4         NOT NULL DEFAULT 0,
    "revoke_reason" varchar(128) NOT NULL DEFAULT '',
    "create_time"  int4         NOT NULL DEFAULT 0
);

CREATE INDEX ON "chatclaw_token_log" ("user_id");
CREATE INDEX ON "chatclaw_token_log" ("create_time");
CREATE INDEX ON "chatclaw_token_log" ("token_hash");
CREATE INDEX ON "chatclaw_token_log" ("user_id", "status", "create_time");

COMMENT ON TABLE "chatclaw_token_log" IS 'ChatClaw 客户端 Token 签发日志';

COMMENT ON COLUMN "chatclaw_token_log"."id"           IS '自增ID';
COMMENT ON COLUMN "chatclaw_token_log"."user_id"      IS '用户ID';
COMMENT ON COLUMN "chatclaw_token_log"."token"        IS '签发的 JWT Token';
COMMENT ON COLUMN "chatclaw_token_log"."token_hash"   IS 'token 的 sha256';
COMMENT ON COLUMN "chatclaw_token_log"."os_type"      IS '客户端操作系统类型（如 Windows / macOS / Linux）';
COMMENT ON COLUMN "chatclaw_token_log"."os_version"   IS '客户端操作系统版本号';
COMMENT ON COLUMN "chatclaw_token_log"."client_ip"    IS '客户端 IP 地址';
COMMENT ON COLUMN "chatclaw_token_log"."expired_at"   IS 'Token 过期时间戳（Unix 秒）';
COMMENT ON COLUMN "chatclaw_token_log"."status"       IS '状态:1生效中,2已下线';
COMMENT ON COLUMN "chatclaw_token_log"."revoke_time"  IS '下线时间戳';
COMMENT ON COLUMN "chatclaw_token_log"."revoke_by"    IS '下线操作人ID';
COMMENT ON COLUMN "chatclaw_token_log"."revoke_reason" IS '下线原因';
COMMENT ON COLUMN "chatclaw_token_log"."create_time"  IS '记录创建时间戳（Unix 秒）';
