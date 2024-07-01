-- +goose Up

INSERT INTO menu (id, "name", "path", parent_id, is_deleted, create_time, update_time, operate_id, operate_name,
                  uni_key)
VALUES (4, '客户端管理', '', 0, 0, 0, 0, 0, '', 'ClientSideManage');

INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
VALUES ('p', '2', 'ClientSideManage', 'GET', '', '', '');

COMMENT ON COLUMN "role"."role_type" IS '角色类型 0自定义,1:所有者,2:管理员,3:成员,4:客户端使用者';

INSERT INTO role ("name", "mark", "is_deleted", "create_time", "update_time", "operate_id", "role_type",
                  "operate_name",
                  "create_name", "parent_id")
VALUES ('客户端使用者', '', 0, 1718784806, 1718784806, 0, 4, '', '系统', 0);

INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
VALUES ('p', (SELECT "last_value" FROM role_id_seq), 'ClientSideManage', 'GET', '', '', '');
