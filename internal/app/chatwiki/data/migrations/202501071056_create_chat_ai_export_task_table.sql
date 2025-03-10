-- +goose Up

CREATE TABLE "chat_ai_export_task"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "robot_id"      int4          NOT NULL DEFAULT 0,
    "file_name"     varchar(100)  NOT NULL DEFAULT '',
    "source"        int4          NOT NULL DEFAULT 0,
    "params"        text          NOT NULL DEFAULT '',
    "status"        int4          NOT NULL DEFAULT 0,
    "file_url"      varchar(500)  NOT NULL DEFAULT '',
    "err_msg"       varchar(1000) NOT NULL DEFAULT '',
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_export_task" ("admin_user_id");
CREATE INDEX ON "chat_ai_export_task" ("robot_id");

COMMENT ON TABLE "chat_ai_export_task" IS '文档问答机器人-导出记录';

COMMENT ON COLUMN "chat_ai_export_task"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_export_task"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_export_task"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "chat_ai_export_task"."file_name" IS '文件名称';
COMMENT ON COLUMN "chat_ai_export_task"."source" IS '来源类型';
COMMENT ON COLUMN "chat_ai_export_task"."params" IS '导出任务需要的参数';
COMMENT ON COLUMN "chat_ai_export_task"."status" IS '状态:0等待导出,1导出中,2导出成功,3导出失败';
COMMENT ON COLUMN "chat_ai_export_task"."file_url" IS '文件链接';
COMMENT ON COLUMN "chat_ai_export_task"."err_msg" IS '错误原因';
COMMENT ON COLUMN "chat_ai_export_task"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_export_task"."update_time" IS '更新时间';