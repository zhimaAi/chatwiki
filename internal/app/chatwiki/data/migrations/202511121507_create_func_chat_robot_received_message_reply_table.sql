-- +goose Up

CREATE TABLE "func_chat_robot_received_message_reply"
(
    "id"                   serial      NOT NULL primary key,
    "admin_user_id"        int4        NOT NULL DEFAULT 0,
    "robot_id"             int4        NOT NULL DEFAULT 0,
    "rule_type"            varchar(50) NOT NULL DEFAULT '',
    "duration_type"        varchar(50) NOT NULL DEFAULT '',
    "week_duration"        jsonb       NOT NULL DEFAULT '[]',
    "start_day"            varchar(50) NOT NULL DEFAULT '',
    "end_day"              varchar(50) NOT NULL DEFAULT '',
    "start_duration"       varchar(50) NOT NULL DEFAULT '',
    "end_duration"         varchar(50) NOT NULL DEFAULT '',
    "priority_num"         int4        NOT NULL DEFAULT 0,
    "reply_interval"       int4        NOT NULL DEFAULT 0,
    "message_type"         int4        NOT NULL DEFAULT 0,
    "specify_message_type" jsonb       NOT NULL DEFAULT '[]',
    "reply_content"        jsonb       NOT NULL DEFAULT '[]',
    "reply_type"           jsonb       NOT NULL DEFAULT '[]',
    "switch_status"        int2        NOT NULL DEFAULT 0,
    "reply_num"            int4        NOT NULL DEFAULT 0,
    "create_time"          int4        NOT NULL DEFAULT 0,
    "update_time"          int4        NOT NULL DEFAULT 0
);

CREATE INDEX on "func_chat_robot_received_message_reply"("robot_id");

COMMENT ON TABLE "func_chat_robot_received_message_reply" IS '聊天机器人关键词回复规则';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."id" IS '自增ID';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."robot_id" IS '机器人ID，为0表示非机器人类型';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."rule_type" IS '规则类型：receive_reply_message_type:类型规则，receive_reply_duration:时间规则';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."duration_type" IS '时间类型：week:周，day:天，time_range:时间范围';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."week_duration" IS '周时间段';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."start_day" IS '开始指定时间';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."end_day" IS '结束指定时间';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."start_duration" IS '开始时间';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."end_duration" IS '结束时间';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."priority_num" IS '优先级';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."reply_interval" IS '触发回复间隔时间 单位秒，0表示不限制';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."message_type" IS '指定消息类型：0 全部 1 指定消息';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."specify_message_type" IS '指定的消息类型 ';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."reply_content" IS '回复内容';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."reply_type" IS '包含的回复类型：imageText 图文 text 文本 url 链接 image 图片 card 小程序';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."switch_status" IS '开启状态：0关1开';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."reply_num" IS '回复方式：0全部 1随机一条';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."create_time" IS '创建时间';
COMMENT ON COLUMN "func_chat_robot_received_message_reply"."update_time" IS '更新时间';
