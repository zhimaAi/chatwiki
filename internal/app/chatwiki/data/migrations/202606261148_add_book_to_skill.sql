-- +goose Up
CREATE TABLE "public"."chat_ai_book_to_skill_task"
(
    id                  serial        NOT NULL PRIMARY KEY,
    admin_user_id       int4          NOT NULL DEFAULT 0,
    robot_id            int4          NOT NULL DEFAULT 0,
    robot_key           varchar(50)   NOT NULL DEFAULT '',
    skill_name          varchar(60)   NOT NULL DEFAULT '',
    model_config_id     int4          NOT NULL DEFAULT 0,
    use_model           varchar(100)  NOT NULL DEFAULT '',
    status              int4          NOT NULL DEFAULT 0,
    source_files        text          NOT NULL DEFAULT '',
    skill_dir           varchar(500)  NOT NULL DEFAULT '',
    skill_file_path     varchar(500)  NOT NULL DEFAULT '',
    error_msg           varchar(2000) NOT NULL DEFAULT '',
    install_status      int4          NOT NULL DEFAULT 0,
    current_iteration   int4          NOT NULL DEFAULT 0,
    template_content    text          NOT NULL DEFAULT '',
    lang                varchar(10)   NOT NULL DEFAULT 'zh-CN',
    create_time         int4          NOT NULL DEFAULT 0,
    update_time         int4          NOT NULL DEFAULT 0
);

CREATE INDEX idx_btos_task_admin_user_id ON "public"."chat_ai_book_to_skill_task" (admin_user_id);
CREATE INDEX idx_btos_task_robot_id ON "public"."chat_ai_book_to_skill_task" (robot_id);
CREATE INDEX idx_btos_task_status ON "public"."chat_ai_book_to_skill_task" (status);

COMMENT ON TABLE "public"."chat_ai_book_to_skill_task" IS 'Book-to-Skill文档转技能任务表';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".id IS '自增主键ID';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".admin_user_id IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".robot_id IS '机器人ID';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".robot_key IS '机器人key(冗余,便于检索)';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".skill_name IS '技能名称';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".model_config_id IS '模型配置ID';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".use_model IS '使用的模型名称';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".status IS '任务状态:0=待开始,1=执行中,2=成功,3=失败,4=已停止';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".source_files IS '上传的源文件信息JSON';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".skill_dir IS '生成的技能目录路径';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".skill_file_path IS '技能文件路径';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".error_msg IS '错误信息';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".install_status IS '安装状态:0=未安装,1=已安装';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".create_time IS '创建时间(Unix时间戳)';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".update_time IS '更新时间(Unix时间戳)';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".current_iteration IS '当前迭代轮次(跨阶段累计)';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".template_content IS '自定义模板内容(为空则使用默认模板)';
COMMENT ON COLUMN "public"."chat_ai_book_to_skill_task".lang IS '任务语言(zh-CN或en-US)';
