-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "work_flow_ids" varchar(1000) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."chat_ai_robot"."work_flow_ids" IS '关联的工作流ID集';
