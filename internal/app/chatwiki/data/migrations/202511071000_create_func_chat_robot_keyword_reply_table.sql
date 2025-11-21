-- +goose Up

CREATE TABLE "func_chat_robot_keyword_reply"
(
    "id"                serial NOT NULL primary key,
    "admin_user_id"     int4 NOT NULL DEFAULT 0,
    "robot_id"          int4 NOT NULL DEFAULT 0,
    "name"              varchar(255) NOT NULL DEFAULT '',
    "full_keyword"      jsonb   NOT NULL DEFAULT '[]',
    "half_keyword"      jsonb   NOT NULL DEFAULT '[]',
    "reply_content"     jsonb   NOT NULL DEFAULT '[]',
    "reply_type"        jsonb   NOT NULL DEFAULT '[]',
    "switch_status"     int2 NOT NULL DEFAULT 0,
    "reply_num"        int4 NOT NULL DEFAULT 0,
    "create_time"       int4 NOT NULL DEFAULT 0,
    "update_time"       int4 NOT NULL DEFAULT 0
);

CREATE INDEX on "func_chat_robot_keyword_reply"("robot_id");

COMMENT ON TABLE "func_chat_robot_keyword_reply" IS '聊天机器人关键词回复规则';

COMMENT ON COLUMN "func_chat_robot_keyword_reply"."id" IS '自增ID';
COMMENT ON COLUMN "func_chat_robot_keyword_reply"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "func_chat_robot_keyword_reply"."robot_id" IS '机器人ID，为0表示非机器人类型';
COMMENT ON COLUMN "func_chat_robot_keyword_reply"."name" IS '规则名称';
COMMENT ON COLUMN "func_chat_robot_keyword_reply"."full_keyword" IS '精准匹配关键词';
COMMENT ON COLUMN "func_chat_robot_keyword_reply"."half_keyword" IS '模糊匹配关键词';
COMMENT ON COLUMN "func_chat_robot_keyword_reply"."reply_content" IS '回复内容';
COMMENT ON COLUMN "func_chat_robot_keyword_reply"."reply_type" IS '包含的回复类型：imageText 图文 text 文本 url 链接 image 图片 card 小程序';
COMMENT ON COLUMN "func_chat_robot_keyword_reply"."switch_status" IS '开启状态：0关1开';
COMMENT ON COLUMN "func_chat_robot_keyword_reply"."reply_num" IS '回复方式：0全部 1随机一条';
COMMENT ON COLUMN "func_chat_robot_keyword_reply"."create_time" IS '创建时间';
COMMENT ON COLUMN "func_chat_robot_keyword_reply"."update_time" IS '更新时间';
