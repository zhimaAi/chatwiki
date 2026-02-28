-- +goose Up

ALTER TABLE "public"."chat_ai_robot" 
  ADD COLUMN "recall_neighbor_topK" int4 NOT NULL DEFAULT 5;

COMMENT ON COLUMN "public"."chat_ai_robot"."recall_neighbor_switch" IS '召回相邻分段开关(废弃)';

COMMENT ON COLUMN "public"."chat_ai_robot"."recall_neighbor_before_num" IS '召回相邻分段前n条';

COMMENT ON COLUMN "public"."chat_ai_robot"."recall_neighbor_after_num" IS '召回相邻分段后n条';

COMMENT ON COLUMN "public"."chat_ai_robot"."recall_neighbor_topK" IS '召回相邻分段的topK'; 




