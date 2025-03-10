-- +goose Up

CREATE TABLE "form_filter_condition"
(
    "id"                serial NOT NULL PRIMARY KEY,
    "form_filter_id"    int4 NOT NULL DEFAULT 0,
    "form_field_id"     int4 NOT NULL DEFAULT 0,
    "rule"              varchar(100) NOT NULL DEFAULT '',
    "rule_value1"       varchar(100) NOT NULL DEFAULT '',
    "rule_value2"       varchar(100) NOT NULL DEFAULT '',
    "create_time"       int4 NOT NULL DEFAULT 0,
    "update_time"       int4 NOT NULL DEFAULT 0
);

CREATE INDEX ON "form_filter_condition" ("form_filter_id", "form_field_id");

COMMENT ON TABLE "form_filter_condition" IS '表单过滤器条件';

COMMENT ON COLUMN "form_filter_condition"."id" IS 'ID';
COMMENT ON COLUMN "form_filter_condition"."form_field_id" IS '过滤器id';
COMMENT ON COLUMN "form_filter_condition"."rule" IS '条件规则';
COMMENT ON COLUMN "form_filter_condition"."rule_value1" IS '规则对应的数值';
COMMENT ON COLUMN "form_filter_condition"."rule_value2" IS '规则对应的数值，有些规则需要多个数值，这里是第二个，如不需要则该字段为空';
COMMENT ON COLUMN "form_filter_condition"."create_time" IS '创建时间';
COMMENT ON COLUMN "form_filter_condition"."update_time" IS '更新时间';
