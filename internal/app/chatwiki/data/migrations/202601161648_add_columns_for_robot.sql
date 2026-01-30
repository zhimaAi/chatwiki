-- +goose Up

ALTER TABLE "chat_ai_robot" ADD COLUMN "recall_neighbor_switch" bool NOT NULL DEFAULT false;
COMMENT ON COLUMN "chat_ai_robot"."recall_neighbor_switch" IS '召回邻居开关';

ALTER TABLE "chat_ai_robot" ADD COLUMN "recall_neighbor_before_num" INT4 NOT NULL DEFAULT 1;
COMMENT ON COLUMN "chat_ai_robot"."recall_neighbor_before_num" IS '召回邻居前n条';

ALTER TABLE "chat_ai_robot" ADD COLUMN "recall_neighbor_after_num" INT4 NOT NULL DEFAULT 1;
COMMENT ON COLUMN "chat_ai_robot"."recall_neighbor_after_num" IS '召回邻居后n条';




