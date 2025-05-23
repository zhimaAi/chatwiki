-- +goose Up


ALTER TABLE "user" ADD COLUMN "expire_time" int4 NOT NULL  default 0;
ALTER TABLE "user" ADD COLUMN "login_switch" int4 NOT NULL  default 1;
COMMENT ON COLUMN "user"."expire_time" IS '账号过期时间 0:永久';
COMMENT ON COLUMN "user"."login_switch" IS '登录开关 1:开,0:关';

ALTER TABLE "chat_ai_robot" ADD COLUMN "creator" int4 NOT NULL DEFAULT 0;
COMMENT ON COLUMN "chat_ai_robot"."creator" IS '创建人';

ALTER TABLE "form" ADD COLUMN "creator" int4 NOT NULL DEFAULT 0;
COMMENT ON COLUMN "form"."creator" IS '创建人';

CREATE TABLE "public"."permission_manage"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "identity_id" int4          NOT NULL DEFAULT 0,
    "object_id" int4          NOT NULL DEFAULT 0,
    "operate_rights"      int2          NOT NULL DEFAULT 0,
    "creator" int4          NOT NULL DEFAULT 0,
    "identity_type" int2 NOT NULL DEFAULT 1,
    "object_type" int2 NOT NULL DEFAULT 1,
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."permission_manage" ("admin_user_id");
CREATE INDEX ON "public"."permission_manage" ("identity_id");
CREATE INDEX ON "public"."permission_manage" ("object_id");
COMMENT ON TABLE "public"."permission_manage" IS '权限管理表';
COMMENT ON COLUMN "public"."permission_manage"."id" IS 'ID';
COMMENT ON COLUMN "public"."permission_manage"."operate_rights" IS '权限 0:无权限 1:查看 2:编辑 4:管理';
COMMENT ON COLUMN "public"."permission_manage"."identity_id" IS '身份ID';
COMMENT ON COLUMN "public"."permission_manage"."identity_type" IS '1:用户 2:部门';
COMMENT ON COLUMN "public"."permission_manage"."object_id" IS '-1为所有';
COMMENT ON COLUMN "public"."permission_manage"."object_type" IS '1:机器人,2:知识库 3:数据库';
COMMENT ON COLUMN "public"."permission_manage"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."permission_manage"."update_time" IS '更新时间';

CREATE TABLE "public"."department"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "department_name"  varchar(100)  NOT NULL DEFAULT '',
    "is_default" int2 NOT NULL DEFAULT 0,
    "pid" int4         NOT NULL DEFAULT 0,
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."department" ("admin_user_id");
CREATE INDEX ON "public"."department" ("pid");

COMMENT ON TABLE "public"."department" IS '部门表';

COMMENT ON COLUMN "public"."department"."id" IS 'ID';
COMMENT ON COLUMN "public"."department"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."department"."department_name" IS '部门名称';
COMMENT ON COLUMN "public"."department"."is_default" IS '部门类型 1:默认部门,0:普通部门';
COMMENT ON COLUMN "public"."department"."pid" IS '父级部门id';
COMMENT ON COLUMN "public"."department"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."department"."update_time" IS '更新时间';
CREATE TABLE "public"."department_member"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "department_id"  int4  NOT NULL DEFAULT 0,
    "user_id" int2 NOT NULL DEFAULT 0,
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."department_member" ("admin_user_id");
CREATE INDEX ON "public"."department_member" ("department_id");
CREATE INDEX ON "public"."department_member" ("user_id");

COMMENT ON TABLE "public"."department_member" IS '部门成员表';

COMMENT ON COLUMN "public"."department_member"."id" IS 'ID';
COMMENT ON COLUMN "public"."department_member"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."department_member"."department_id" IS '部门id';
COMMENT ON COLUMN "public"."department_member"."user_id" IS '用户id';
COMMENT ON COLUMN "public"."department_member"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."department_member"."update_time" IS '更新时间';

CREATE TABLE "public"."department_config"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "department_id" int4          NOT NULL DEFAULT 0,
    "max_level"      int2          NOT NULL DEFAULT 0,
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);