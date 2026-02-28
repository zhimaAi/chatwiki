-- +goose Up

ALTER TABLE "chat_ai_robot" 
  ALTER COLUMN "recall_neighbor_switch" SET DEFAULT true; 




