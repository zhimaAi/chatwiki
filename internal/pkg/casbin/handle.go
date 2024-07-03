// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package casbin

import (
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var (
	once    sync.Once
	Handler = new(casbinHandler)
)

type casbinHandler struct {
	syncedEnforcer *casbin.Enforcer
}

func (c *casbinHandler) Init(DB *gorm.DB) error {
	once.Do(func() {
		adapter, err := gormadapter.NewAdapterByDB(DB)
		if err != nil {
			panic(err)
		}
		confPath := "internal/pkg/casbin/rbac_model.conf"
		c.syncedEnforcer, err = casbin.NewEnforcer(confPath, adapter)
		if err != nil {
			panic(err)
		}
	})
	c.syncedEnforcer.AddFunction("matchRootFunc", MatchRootFunc)
	err := c.syncedEnforcer.LoadPolicy()
	if err != nil {
		return err
	}
	return nil
}
func MatchRootFunc(arguments ...interface{}) (interface{}, error) {
	username := arguments[0].(string)
	// check is root
	return Handler.syncedEnforcer.HasRoleForUser(username, fmt.Sprintf("%d", Root))
}

// Enforce check permission
func (c *casbinHandler) Enforce(user, uri, action string) (bool, error) {
	return c.syncedEnforcer.Enforce(user, uri, action)
}

func (c *casbinHandler) AddPolicy(roleId, uri, method string) (bool, error) {
	return c.syncedEnforcer.AddPolicy(roleId, uri, method)
}

func (c *casbinHandler) LoadPolicy() error {
	return c.syncedEnforcer.LoadPolicy()
}

func (c *casbinHandler) AddPolicies(rules [][]string) (bool, error) {
	success, err := c.syncedEnforcer.AddPolicies(rules)
	_ = c.syncedEnforcer.LoadPolicy()
	return success, err
}

func (c *casbinHandler) DeleteRole(roleId string) (bool, error) {
	return c.syncedEnforcer.DeleteRole(roleId)
}

// DeleteUserRole delete user role info o user; 1 role
func (c *casbinHandler) DeleteUserRole(pIndex int, role string) (bool, error) {
	return c.syncedEnforcer.RemoveFilteredNamedGroupingPolicy("g", pIndex, role)
}

func (c *casbinHandler) DelRoleRules(role string) (bool, error) {
	success, err := c.syncedEnforcer.RemoveNamedPolicy("p", role)
	_ = c.syncedEnforcer.LoadPolicy()
	return success, err
}

// AddUserRole add user role
func (c *casbinHandler) AddUserRole(user string, roleId string) (bool, error) {
	return c.syncedEnforcer.AddGroupingPolicy(user, roleId)
}

// UpdateUserRole add user role
func (c *casbinHandler) UpdateUserRole(user string, roleId string) (bool, error) {
	_, _ = c.DeleteUserRole(0, user)
	return c.syncedEnforcer.AddGroupingPolicy(user, roleId)
}

// AddUserRoles batch add
func (c *casbinHandler) AddUserRoles(usernames, roleIds []string) (bool, error) {
	rules := make([][]string, 0)
	for _, u := range usernames {
		for _, r := range roleIds {
			rules = append(rules, []string{u, r})
		}
	}
	return c.syncedEnforcer.AddGroupingPolicies(rules)
}

func (c *casbinHandler) GetRolesForUser(user string) ([]string, error) {
	return c.syncedEnforcer.GetRolesForUser(user)
}

func (c *casbinHandler) GetUsersForRole(role string) ([]string, error) {
	return c.syncedEnforcer.GetUsersForRole(role)
}

func (c *casbinHandler) GetPolicyForUser(role string) ([][]string, error) {
	return c.syncedEnforcer.GetFilteredPolicy(0, role)
}
