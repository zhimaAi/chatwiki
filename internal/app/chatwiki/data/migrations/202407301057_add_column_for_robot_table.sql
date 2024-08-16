-- +goose Up

ALTER TABLE "chat_ai_robot" ADD COLUMN "form_ids" varchar(1000) NOT NULL DEFAULT '';

COMMENT ON COLUMN "chat_ai_robot"."form_ids" IS '数据表ID集合';
