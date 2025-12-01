-- +goose Up
ALTER TABLE "public"."work_flow_node"
    ADD COLUMN "loop_parent_key" varchar(32) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."work_flow_node"."loop_parent_key" IS '所属的循环节点';

ALTER TABLE "public"."work_flow_node_version"
    ADD COLUMN "loop_parent_key" varchar(32) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."work_flow_node_version"."loop_parent_key" IS '所属的循环节点';
