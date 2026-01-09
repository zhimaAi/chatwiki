-- +goose Up
CREATE TABLE "public"."robot_payment_setting"
(
    "id" serial NOT NULL PRIMARY KEY,
    "create_time"   int4 NOT NULL DEFAULT 0,
    "update_time"   int4 NOT NULL DEFAULT 0,
    "admin_user_id" int4 NOT NULL DEFAULT 0,
    "robot_id"      int4 NOT NULL DEFAULT 0,
    "try_count"     int4 NOT NULL DEFAULT 0,
    "package_type"  int2 NOT NULL DEFAULT 1,
    "count_package" jsonb NOT NULL DEFAULT '[]'::jsonb,
    "duration_package"   jsonb NOT NULL DEFAULT '[]'::jsonb,
    "contact_qrcode" varchar(300) NOT NULL DEFAULT '',
    "package_poster" varchar(300) NOT NULL DEFAULT ''
);
COMMENT ON TABLE "public"."robot_payment_setting" IS '机器人应用收费设置';
COMMENT ON COLUMN "public"."robot_payment_setting"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."robot_payment_setting"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "public"."robot_payment_setting"."try_count" IS '允许试用次数';
COMMENT ON COLUMN "public"."robot_payment_setting"."package_type" IS '套餐方案 1按次数收费 2按时长收费';
COMMENT ON COLUMN "public"."robot_payment_setting"."count_package" IS '次数套餐详情';
COMMENT ON COLUMN "public"."robot_payment_setting"."duration_package" IS '时长套餐详情';
COMMENT ON COLUMN "public"."robot_payment_setting"."contact_qrcode" IS '联系方式二维码';
COMMENT ON COLUMN "public"."robot_payment_setting"."package_poster" IS '海报地址';

CREATE TABLE "public"."robot_payment_auth_code"
(
    "id" serial NOT NULL PRIMARY KEY,
    "create_time"   int4 NOT NULL DEFAULT 0,
    "update_time"   int4 NOT NULL DEFAULT 0,
    "creator_id"    int4 NOT NULL DEFAULT 0,
    "creator_name" varchar(100) NOT NULL DEFAULT '',
    "admin_user_id" int4 NOT NULL DEFAULT 0,
    "robot_id"      int4 NOT NULL DEFAULT 0,
    "content"       varchar(100) NOT NULL DEFAULT '',
    "package_type"  int2 NOT NULL DEFAULT 0,
    "package_id"    int4 NOT NULL DEFAULT 0,
    "package_name"  varchar(100) NOT NULL DEFAULT '',
    "package_duration" int4 NOT NULL DEFAULT 0,
    "package_count"    int4 NOT NULL DEFAULT 0,
    "package_price"    decimal(20, 2) NOT NULL DEFAULT 0,
    "usage_status"  int4 NOT NULL DEFAULT 0,
    "exchanger_openid" varchar(100) NOT NULL DEFAULT '',
    "exchanger_name" varchar(100) NOT NULL DEFAULT '',
    "exchanger_avatar" varchar(1000) NOT NULL DEFAULT '',
    "exchange_time" int4 NOT NULL DEFAULT 0,
    "use_time"      int4 NOT NULL DEFAULT 0,
    "use_date"      varchar(10) NOT NULL DEFAULT CURRENT_DATE,
    "remark"        varchar(100) NOT NULL DEFAULT '',
    "create_date"   varchar(10) NOT NULL DEFAULT CURRENT_DATE,
    "exchange_date" varchar(10) NOT NULL DEFAULT CURRENT_DATE,
    "used_count"    int4 NOT NULL DEFAULT 0,
    "used_duration" int4 NOT NULL DEFAULT 0
);
COMMENT ON TABLE "public"."robot_payment_auth_code" IS '应用收费的授权码';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."creator_id" IS '创建人ID';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."creator_name" IS '创建人用户名(主要是微信公众号场景)';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."admin_user_id" IS '管理员ID';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."package_type" IS '套餐类型 1按次数收费 2按时长收费';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."package_id" IS '关联的套餐ID';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."package_duration" IS '时长';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."package_count" IS '次数';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."package_price" IS '费用';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."usage_status" IS '使用状态 1未使用 2已兑换 3已使用';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."exchanger_openid" IS '兑换者openid';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."exchanger_name" IS '兑换者名称';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."exchanger_avatar" IS '兑换者头像';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."exchange_time" IS '兑换时间';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."use_time" IS '使用时间';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."remark" IS '备注';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."create_date" IS '创建日期';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."exchange_date" IS '兑换日期';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."used_count" IS '已使用次数';
COMMENT ON COLUMN "public"."robot_payment_auth_code"."used_duration" IS '已使用天数';
CREATE INDEX on "public"."robot_payment_auth_code" USING btree ("admin_user_id", "robot_id", "create_date");
CREATE INDEX on "public"."robot_payment_auth_code" USING btree ("admin_user_id", "robot_id", "exchange_date");

CREATE TABLE "public"."robot_payment_auth_code_manager"
(
    "id" serial NOT NULL PRIMARY KEY,
    "create_time"   int4 NOT NULL DEFAULT 0,
    "update_time"   int4 NOT NULL DEFAULT 0,
    "admin_user_id" int4 NOT NULL DEFAULT 0,
    "robot_id"      int4 NOT NULL DEFAULT 0,
    "manager_openid" varchar(100) NOT NULL DEFAULT '',
    "manager_avatar" varchar(300) NOT NULL DEFAULT '',
    "manager_nickname" varchar(300) NOT NULL DEFAULT ''
);
COMMENT ON TABLE "public"."robot_payment_auth_code_manager" IS '应用收费的授权码管理员列表';
COMMENT ON COLUMN "public"."robot_payment_auth_code_manager"."admin_user_id" IS '管理员ID';
COMMENT ON COLUMN "public"."robot_payment_auth_code_manager"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "public"."robot_payment_auth_code_manager"."manager_openid" IS '管理员openid';
COMMENT ON COLUMN "public"."robot_payment_auth_code_manager"."manager_avatar" IS '管理员头像';
COMMENT ON COLUMN "public"."robot_payment_auth_code_manager"."manager_nickname" IS '管理员昵称';

CREATE TABLE "public"."robot_payment_user_try_count"
(
    "id" serial NOT NULL PRIMARY KEY,
    "create_time"   int4 NOT NULL DEFAULT 0,
    "update_time"   int4 NOT NULL DEFAULT 0,
    "admin_user_id" int4 NOT NULL DEFAULT 0,
    "robot_id"      int4 NOT NULL DEFAULT 0,
    "openid"        varchar(100) NOT NULL DEFAULT 0,
    "try_count"    int4 NOT NULL DEFAULT 0
);
COMMENT ON TABLE "public"."robot_payment_user_try_count" IS '收费应用用户试用次数';
COMMENT ON COLUMN "public"."robot_payment_user_try_count"."admin_user_id" IS '管理员ID';
COMMENT ON COLUMN "public"."robot_payment_user_try_count"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "public"."robot_payment_user_try_count"."openid" IS '用户OPENID';
COMMENT ON COLUMN "public"."robot_payment_user_try_count"."try_count" IS '试用次数';