-- +goose Up

CREATE INDEX ON "public"."llm_request_logs" ("create_time");

CREATE INDEX ON "public"."chat_ai_session" ("admin_user_id", "robot_id", "app_type", "create_time");