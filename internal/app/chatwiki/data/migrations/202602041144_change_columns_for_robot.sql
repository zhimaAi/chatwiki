-- +goose Up

ALTER TABLE "public"."chat_ai_robot" RENAME COLUMN "recall_neighbor_topK" TO "recall_neighbor_top_k";




