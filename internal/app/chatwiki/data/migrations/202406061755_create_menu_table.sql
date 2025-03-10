-- +goose Up

CREATE TABLE "menu"
(
    "id" serial NOT NULL primary key,
    "name" varchar(100) NOT NULL DEFAULT '' ,
    "uni_key" varchar(100) DEFAULT '',
    "path" varchar(500) NOT NULL DEFAULT '',
    "parent_id" int4 NOT NULL DEFAULT 0,
    "is_deleted" int2 NOT NULL DEFAULT 0,
    "create_time" int4 NOT NULL DEFAULT 0,
    "update_time" int4 NOT NULL DEFAULT 0,
    "operate_id" int4 NOT NULL DEFAULT 0,
    "operate_name" varchar(128) NOT NULL DEFAULT ''
);

COMMENT ON COLUMN "menu"."id" IS 'ID';
COMMENT ON COLUMN "menu"."name" IS '名称';
COMMENT ON COLUMN "menu"."uni_key" IS '唯一标识';
COMMENT ON COLUMN "menu"."is_deleted" IS '是否删除 1:删除';
COMMENT ON COLUMN "menu"."create_time" IS '创建时间';
COMMENT ON COLUMN "menu"."update_time" IS '更新时间';
COMMENT ON COLUMN "menu"."operate_id" IS '操作人id';
COMMENT ON COLUMN "menu"."operate_name" IS '操作人';

INSERT INTO menu
(id, "name", "path", parent_id, is_deleted, create_time, update_time, operate_id, operate_name, uni_key)
VALUES(1, '机器人管理', '/manage/getRobotList,/manage/getRobotInfo', 0, 0, 0, 0, 0, '', 'RobotManage');

INSERT INTO menu
(id, "name", "path", parent_id, is_deleted, create_time, update_time, operate_id, operate_name, uni_key)
VALUES(2, '知识库管理', '/manage/getLibraryList,/manage/getLibraryInfo', 0, 0, 0, 0, 0, '', 'LibraryManage');

INSERT INTO menu
(id, "name", "path", parent_id, is_deleted, create_time, update_time, operate_id, operate_name, uni_key)
VALUES(3, '系统设置', '', 0, 0, 0, 0, 0, '', 'SystemManage');
