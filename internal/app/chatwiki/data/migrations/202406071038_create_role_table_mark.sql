-- +goose Up
ALTER TABLE public."role" ALTER COLUMN mark SET DEFAULT '';

UPDATE public."role" SET mark='系统角色，拥有所有权限，不可修改和删除' WHERE id=1;
UPDATE public."role" SET mark='团队管理员，拥有等同于管理员的权限，可以编辑，不可删除' WHERE id=2;
UPDATE public."role" SET mark='团队成员' WHERE id=3;