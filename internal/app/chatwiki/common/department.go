// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"errors"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetDefaultDepartmentInfo(adminUserId int) (msql.Params, error) {
	m := msql.Model(`department`, define.Postgres)
	info, err := m.Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`is_default`, cast.ToString(define.SwitchOn)).Find()
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	return info, nil
}

func GetDepartmentInfo(id int) (msql.Params, error) {
	m := msql.Model(`department`, define.Postgres)
	info, err := m.Where(`id`, cast.ToString(id)).Find()
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	return info, nil
}

func SaveDepartment(id, adminUserId int64, data msql.Datas) (int64, error) {
	if len(data) <= 0 {
		return id, errors.New(`data empty`)
	}
	data[`update_time`] = tool.Time2Int()
	m := msql.Model(`department`, define.Postgres)
	var err error
	if id > 0 {
		// 更新
		info, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Value(`id`)
		if err != nil {
			logs.Error(err.Error())
			return id, err
		}
		if len(info) == 0 {
			return id, err
		}
		_, err = m.Where(`id`, cast.ToString(id)).Update(data)
	} else {
		// 新增
		data[`create_time`] = tool.Time2Int()
		id, err = m.Insert(data, `id`)
	}
	if err != nil {
		logs.Error(err.Error())
		return id, err
	}
	return id, nil
}

// GetDepartmentLevel 获取部门层级的pid列表和层级数
func GetDepartmentLevel(adminUserId, pDid int) (int, int, error) {
	// 查询该管理员用户下的所有部门
	result, err := GetAllDepartmentList(adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return 0, 0, err
	}

	// 构建id到pid的映射
	currentId := 0
	idToPid := make(map[int]int)
	for _, v := range result {
		if cast.ToInt(v[`pid`]) == 0 {
			currentId = cast.ToInt(v[`id`])
		}
		idToPid[cast.ToInt(v["id"])] = cast.ToInt(v["pid"])
	}

	// 从指定id开始向下查找所有子部门id
	var childList []string
	maxLevel := 1 // 初始化最大层级为1
	pLevel := 0
	// 用于记录每个节点的层级
	levelMap := make(map[int]int)
	levelMap[currentId] = 1

	var findChildren func(parentId int, level int)
	findChildren = func(parentId int, level int) {
		for id, pid := range idToPid {
			if pid == parentId {
				childList = append(childList, cast.ToString(id))
				currentLevel := level + 1
				levelMap[id] = currentLevel
				if currentLevel > maxLevel {
					maxLevel = currentLevel
				}
				if id == pDid {
					pLevel = currentLevel
				}
				findChildren(id, currentLevel)
			}
		}
	}
	findChildren(currentId, 1)

	return pLevel, maxLevel, nil
}

func OverDepartmentLevel(adminUserId, pid, id int) (bool, int) {
	var maxLevel = define.MaxDepartmentLevel
	// get config
	config, err := msql.Model(`department_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Value(`max_level`)
	if err != nil {
		logs.Error(err.Error())
		return true, maxLevel
	}
	if cast.ToInt(config) > 0 {
		maxLevel = cast.ToInt(config)
	}
	// 查询该管理员用户下的所有部门
	pLevel, currentLevel, err := GetDepartmentLevel(adminUserId, pid)
	if err != nil {
		logs.Error(err.Error())
		return true, maxLevel
	}
	addLevel := 0
	if pLevel >= maxLevel {
		return true, maxLevel
	}
	return currentLevel+addLevel > maxLevel, maxLevel
}

// GetDepartmentMembers 获取部门成员列表
func GetDepartmentMembers(departmentIds string) (map[int][]string, error) {
	if len(departmentIds) <= 0 {
		return nil, errors.New("invalid department id")
	}

	m := msql.Model("department_member", define.Postgres)
	list, err := m.Where("department_id", `in`, departmentIds).Field(`department_id,user_id`).Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	mapList := make(map[int][]string)
	for _, v := range list {
		mapList[cast.ToInt(v[`department_id`])] = append(mapList[cast.ToInt(v[`department_id`])], v[`user_id`])
	}
	return mapList, nil
}

// GetDepartmentMembers 获取成员部门列表
func GetUserDepartments(userId int) ([]msql.Params, error) {
	if userId <= 0 {
		return nil, errors.New("invalid department id")
	}

	list, err := GetUserDepartmentIds(cast.ToString(userId))
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	data, err := msql.Model("department", define.Postgres).Where(`id`, `in`, strings.Join(list, `,`)).Select()
	return data, nil
}

// GetDepartmentMembers 获取成员部门列表
func GetUserDepartmentIds(userIds string) ([]string, error) {
	if len(userIds) <= 0 {
		return nil, errors.New("invalid user id")
	}

	m := msql.Model("department_member", define.Postgres)
	list, err := m.Where("user_id", `in`, cast.ToString(userIds)).ColumnArr(`department_id`)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	return list, nil
}

// GetUserAllDepartmentIds 获取成员部门列表
func GetUserAllDepartmentIds(adminUserId int, userIds string) ([]string, error) {
	list, err := GetUserDepartmentIds(userIds)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	allDepartments, err := GetAllDepartmentList(adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	// 获取部门所有的上级...
	for _, v := range list {
		childList := &[]string{}
		findDepartmentParent(allDepartments, cast.ToInt(v), childList)
		list = append(list, *childList...)
	}
	return list, nil
}

func findDepartmentParent(departments []msql.Params, parentId int, data *[]string) {
	if parentId == 0 {
		return
	}
	for _, item := range departments {
		if parentId == cast.ToInt(item["id"]) {
			*data = append(*data, item[`id`])
			findDepartmentParent(departments, cast.ToInt(item["pid"]), data)
			break
		}
	}
}

func FindDepartmentChildren(departments []msql.Params, parentId int, data *[]string) {
	if parentId == 0 {
		return
	}
	for _, item := range departments {
		if parentId == cast.ToInt(item["pid"]) {
			*data = append(*data, item[`id`])
			FindDepartmentChildren(departments, cast.ToInt(item["id"]), data)
		}
	}
}

func GetAllDepartmentList(adminUserId int) ([]msql.Params, error) {
	m := msql.Model("department", define.Postgres)
	list, err := m.Where("admin_user_id", cast.ToString(adminUserId)).Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	return list, nil
}

type DepartmentInfo struct {
	Id             int               `json:"id"`
	Pid            int               `json:"pid"`
	DepartmentName string            `json:"department_name"`
	IsDefault      int               `json:"is_default"`
	CreateTime     int               `json:"create_time"`
	UpdateTime     int               `json:"update_time"`
	Children       []*DepartmentInfo `json:"children"`
	ChildrenNums   int               `json:"children_nums"`
	UserData       []msql.Params     `json:"user_data"`
}

func GetDepartmentTrees(adminUserId int) ([]*DepartmentInfo, []string, error) {
	m := msql.Model(`department`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId))
	// 获取所有部门数据
	list, err := m.Order(`id asc`).Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, nil, err
	}
	trees := make([]*DepartmentInfo, 0)
	departmentIds := make([]string, 0)
	for _, v := range list {
		departmentIds = append(departmentIds, v[`id`])
	}
	userIdMap, _ := GetDepartmentMembers(strings.Join(departmentIds, `,`))
	// 返回用户
	for _, v := range list {
		var userData = []msql.Params{}
		if userIds, ok := userIdMap[cast.ToInt(v[`id`])]; ok {
			userData, _ = msql.Model(define.TableUser, define.Postgres).Where("is_deleted", define.Normal).Where(`id`, `in`, strings.Join(userIds, `,`)).Order("id asc").Field(`id,user_name,avatar`).Select()
		}
		trees = append(trees, &DepartmentInfo{
			Id:             cast.ToInt(v[`id`]),
			Pid:            cast.ToInt(v[`pid`]),
			IsDefault:      cast.ToInt(v[`is_default`]),
			DepartmentName: cast.ToString(v[`department_name`]),
			CreateTime:     cast.ToInt(v[`create_time`]),
			UpdateTime:     cast.ToInt(v[`update_time`]),
			UserData:       userData,
		})
	}
	return ConvertListToTree(trees, 0), departmentIds, nil
}

func ConvertListToTree(list []*DepartmentInfo, parentId int) []*DepartmentInfo {
	var (
		tree []*DepartmentInfo
	)
	// 计算所有子部门数量
	for _, node := range list {
		if node.Pid == parentId {
			node.Children = ConvertListToTree(list, node.Id)
			// 计算当前节点下所有成员数
			childCount := len(node.UserData)
			for _, child := range node.Children {
				childCount += child.ChildrenNums
			}
			node.ChildrenNums = childCount
			tree = append(tree, node)
		}
	}
	return tree
}

func SaveUserDepartmentData(adminUserId int, userIds, departmentIds string) error {
	// 开启事务
	m := msql.Model("department_member", define.Postgres)
	m.Begin()

	// 删除原有部门关系
	_, err := m.Where("user_id", "in", userIds).Delete()
	if err != nil {
		m.Rollback()
		logs.Error(err.Error())
		return err
	}

	// 批量插入新的部门关系
	userIdArr := strings.Split(userIds, ",")
	for _, departmentId := range strings.Split(departmentIds, ",") {
		for _, uid := range userIdArr {
			_, err = m.Insert(msql.Datas{
				`admin_user_id`: adminUserId,
				"department_id": departmentId,
				"user_id":       cast.ToInt(uid),
				"create_time":   tool.Time2Int(),
				"update_time":   tool.Time2Int(),
			})
			if err != nil {
				m.Rollback()
				logs.Error(err.Error())
				return err
			}
		}
	}
	m.Commit()
	return nil
}
