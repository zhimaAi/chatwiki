-- +goose Up

ALTER TABLE "chat_ai_library"
    ALTER COLUMN "graph_model_config_id" TYPE int4;

ALTER TABLE "public"."sensitive_words_relation"
    ALTER COLUMN "robot_id" TYPE int4;

ALTER TABLE "chat_ai_library_file"
    ALTER COLUMN "semantic_chunk_model_config_id" TYPE int4;

ALTER TABLE "chat_ai_library_file_data"
    ALTER COLUMN "category_id" TYPE int4;

ALTER TABLE "public"."department_member"
    ALTER COLUMN "user_id" TYPE int4;

ALTER TABLE "public"."wechat_official_account_comment_rule"
    ALTER COLUMN "model_config_id" TYPE int4;

ALTER TABLE "public"."wechat_official_comment_list"
    ALTER COLUMN "comment_rule_id" TYPE int4,
    ALTER COLUMN "user_comment_id" TYPE int4,
    ALTER COLUMN "draft_id" TYPE int4;

ALTER TABLE "public"."wechat_official_account_batch_send_task"
    ALTER COLUMN "res_total" TYPE int4,
    ALTER COLUMN "res_succ" TYPE int4,
    ALTER COLUMN "res_err" TYPE int4,
    ALTER COLUMN "common_total" TYPE int4,
    ALTER COLUMN "comment_rule_id" TYPE int4,
    ALTER COLUMN "max_comment_id" TYPE int4;
