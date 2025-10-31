-- +goose Up

COMMENT ON COLUMN "public"."llm_request_daily_stats"."type" IS '类别 1日活用户数 2日新增用户数 3总消息数 4token消耗数 5命中知识库次数 6机器人回复消息数';

