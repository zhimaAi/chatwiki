-- +goose Up

CREATE TABLE "public"."llm_http_log"
(
    "id"                bigserial       NOT NULL PRIMARY KEY,
    "request_time"      timestamptz,
    "url"               varchar(8192)   NOT NULL DEFAULT '',
    "method"            varchar(64)     NOT NULL DEFAULT '',
    "request_headers"   varchar(262144) NOT NULL DEFAULT '',
    "request_body"      bytea           NOT NULL DEFAULT '',
    "response_time"     timestamptz,
    "end_time"          timestamptz,
    "duration"          int8            NOT NULL DEFAULT 0,
    "status_code"       int4            NOT NULL DEFAULT 0,
    "response_headers"  varchar(262144) NOT NULL DEFAULT '',
    "response_body"     bytea           NOT NULL DEFAULT '',
    "response_body_eof" boolean         NOT NULL DEFAULT false,
    "error"             varchar(16384)  NOT NULL DEFAULT '',
    "is_error"          boolean         NOT NULL DEFAULT false,
    "belong_day"        int4            NOT NULL DEFAULT 0
);

-- Hash indexes support long URL equality lookups without the btree index-entry size limit.
CREATE INDEX "llm_http_log_url_idx" ON "public"."llm_http_log" USING hash ("url");
CREATE INDEX "llm_http_log_method_idx" ON "public"."llm_http_log" ("method");
CREATE INDEX "llm_http_log_status_code_idx" ON "public"."llm_http_log" ("status_code");
CREATE INDEX "llm_http_log_is_error_idx" ON "public"."llm_http_log" ("is_error");
CREATE INDEX "llm_http_log_belong_day_idx" ON "public"."llm_http_log" ("belong_day");

COMMENT ON TABLE "public"."llm_http_log" IS '大模型适配器原始HTTP请求响应日志';
COMMENT ON COLUMN "public"."llm_http_log"."request_time" IS '请求开始时间，未提供时为空';
COMMENT ON COLUMN "public"."llm_http_log"."url" IS '原始请求URL，最多8192字符';
COMMENT ON COLUMN "public"."llm_http_log"."method" IS 'HTTP方法，最多64字符';
COMMENT ON COLUMN "public"."llm_http_log"."request_headers" IS '请求头JSON文本，最多262144字符，超长截断';
COMMENT ON COLUMN "public"."llm_http_log"."request_body" IS '原始请求体字节，最多保存16MiB';
COMMENT ON COLUMN "public"."llm_http_log"."response_time" IS '收到响应头的时间，未收到响应时为空';
COMMENT ON COLUMN "public"."llm_http_log"."end_time" IS '响应读取结束或关闭时间，未提供时为空';
COMMENT ON COLUMN "public"."llm_http_log"."duration" IS 'HTTP往返总耗时，单位纳秒';
COMMENT ON COLUMN "public"."llm_http_log"."status_code" IS 'HTTP状态码，未收到响应时为0';
COMMENT ON COLUMN "public"."llm_http_log"."response_headers" IS '响应头JSON文本，最多262144字符，超长截断';
COMMENT ON COLUMN "public"."llm_http_log"."response_body" IS '原始响应体字节，流式响应保留SSE，最多保存16MiB';
COMMENT ON COLUMN "public"."llm_http_log"."response_body_eof" IS '响应体是否已读取至EOF';
COMMENT ON COLUMN "public"."llm_http_log"."error" IS '传输或正文读取错误，最多16384字符';
COMMENT ON COLUMN "public"."llm_http_log"."is_error" IS '存在传输或正文读取错误，或HTTP状态码不小于400';
COMMENT ON COLUMN "public"."llm_http_log"."belong_day" IS '按请求开始时间计算的归属天，格式YYYYMMDD';

-- +goose Down
DROP TABLE "public"."llm_http_log";
