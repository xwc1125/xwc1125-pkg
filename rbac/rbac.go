// Package rbac
//
// @author: xwc1125
package rbac

import (
	"fmt"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/casbin/casbin/v2/util"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"xorm.io/xorm"
)

var (
	DefaultDomain = "default"
)

type RBAC struct {
	adapter persist.Adapter
	*casbin.SyncedEnforcer
}

// NewRbacByGorm 根据gorm初始化casbin
func NewRbacByGorm(db *gorm.DB, tablePrefix string, tableName string) (*RBAC, error) {
	rbacModelOrFile := viper.GetString("rbac.model_file")
	return NewRbacByGormWithModel(rbacModelOrFile, db, tablePrefix, tableName)
}

func NewRbacByGormWithModel(rbacModelOrFile string, db *gorm.DB, tablePrefix string, tableName string) (*RBAC, error) {
	adapter, err := NewAdapterByDBUseTableName(db, tablePrefix, tableName)
	if err != nil {
		return nil, err
	}
	enforcer, err := getEnforcer(adapter, rbacModelOrFile)
	if err != nil {
		return nil, err
	}
	return &RBAC{
		adapter:        adapter,
		SyncedEnforcer: enforcer,
	}, nil
}

// NewRbacByXorm 根据xorm初始化casbin
func NewRbacByXorm(db *xorm.Engine, tablePrefix string, tableName string) (*RBAC, error) {
	return NewRbacByXormWithModel(viper.GetString("rbac.model_file"), db, tablePrefix, tableName)
}
func NewRbacByXormWithModel(rbacModelOrFile string, db *xorm.Engine, tablePrefix string, tableName string) (*RBAC, error) {
	adapter, err := xormadapter.NewAdapterByEngineWithTableName(db, tableName, tablePrefix)
	if err != nil {
		return nil, err
	}
	enforcer, err := getEnforcer(adapter, rbacModelOrFile)
	if err != nil {
		return nil, err
	}
	return &RBAC{
		adapter:        adapter,
		SyncedEnforcer: enforcer,
	}, nil
}

// 获取Enforcer
func getEnforcer(adapter persist.Adapter, rbacModelOrFile string) (*casbin.SyncedEnforcer, error) {
	e, err := casbin.NewSyncedEnforcer(rbacModelOrFile, adapter)
	if err != nil {
		return nil, err
	}
	err = e.LoadPolicy()
	if err != nil {
		return nil, err
	}

	e.SetLogger(&Logger{})
	enableLog := viper.GetBool("rbac.enable_log")
	e.EnableLog(enableLog)
	e.EnableAutoSave(true)
	e.StartAutoLoadPolicy(time.Minute)
	return e, nil
}

// ==================api=================

// AddPermissionsForOwner 给owner设置资源
func (r *RBAC) AddPermissionsForOwner(pType string, owner Owner, permissions []Permission, domain ...string) (bool, error) {
	var dom string
	if len(domain) > 0 {
		dom = domain[0]
	}
	// 先将原始权限删除
	userKey := owner.OwnerKey()

	for _, res := range permissions {
		_, err := r.GetEnforcer().AddNamedPolicy(pType, util.JoinSlice(userKey, res.Resource, res.Action, res.ResourceType, dom))
		// _, err := r.GetEnforcer().AddPermissionForUser(userKey, res.Resource, res.Action, res.ResourceType, dom)
		if err != nil {
			return false, err
		}
	}
	r.GetEnforcer().LoadPolicy()
	return true, nil
}

// DeletePermissionsForOwner 删除Owner的资源权限[p]
// policy_definition规则
func (r *RBAC) DeletePermissionsForOwner(pType string, owner Owner, params []string, domain ...string) (bool, error) {
	if len(pType) == 0 {
		pType = "p"
	}
	userKey := owner.OwnerKey()
	var vals = make([]string, 0)
	vals = append(vals, userKey)
	vals = append(vals, params...)
	vals = append(vals, domain...)
	return r.RemoveFilteredNamedPolicy(pType, 0, vals...)
}

func (r *RBAC) DeleteResource(pType string, perm Permission, params []string, domain ...string) (bool, error) {
	if len(pType) == 0 {
		pType = "p"
	}
	var vals = make([]string, 0)
	vals = append(vals, "")
	vals = append(vals, perm.Resource)
	vals = append(vals, perm.Action)
	vals = append(vals, params...)
	vals = append(vals, domain...)
	return r.RemoveFilteredNamedPolicy(pType, 0, vals...)
}

// HasPermission 检查用户是否有权限【middleware中调用】
func (r *RBAC) HasPermission(owner Owner, path string, method string, extra []string, domain ...string) (bool, error) {
	userKey := owner.OwnerKey()
	// 用户权限判断
	{
		b, err := r.HasPermission2(userKey, path, method, extra, domain...)
		if err != nil {
			log().Error(err.Error())
		}
		if b {
			return b, nil
		}
	}
	// 取出用户的所有角色
	{
		roles, err := r.GetImplicitRolesForUser(userKey, domain...)
		if err != nil {
			return false, err
		}
		// 角色权限判断
		for _, role := range roles {
			b, err := r.HasPermission2(role, path, method, extra, domain...)
			if err != nil {
				log().Error(err.Error())
			}
			if b {
				return b, nil
			}
		}
	}
	return false, nil
}

func (r *RBAC) HasPermission2(subject string, resource string, action string, extra []string, domain ...string) (bool, error) {
	var rvals = make([]interface{}, 0)
	rvals = append(rvals, subject)
	rvals = append(rvals, resource)
	rvals = append(rvals, action)
	for _, s := range extra {
		rvals = append(rvals, s)
	}
	for _, s := range domain {
		rvals = append(rvals, s)
	}
	return r.GetEnforcer().Enforce(rvals...)
}

// 资源绑定

// AddProvidersForOwner 给Owner添加提供者
func (r *RBAC) AddProvidersForOwner(pType string, owner Owner, providers []Owner, domain ...string) (bool, error) {
	if len(pType) == 0 {
		pType = "g"
	}
	var rules [][]string
	userKey := owner.OwnerKey()

	for _, provider := range providers {
		rule := []string{userKey, provider.OwnerKey()}
		rule = append(rule, domain...)
		rules = append(rules, rule)
	}
	return r.AddNamedGroupingPolicies(pType, rules)
}

// GetProvidersForOwner 根据Owner获取提供者
func (r *RBAC) GetProvidersForOwner(pType string, owner Owner, domain ...string) ([]string, error) {
	userKey := owner.OwnerKey()
	if len(pType) == 0 {
		pType = "g"
	}
	roleManager := r.GetModel()["g"][pType].RM
	if roleManager == nil {
		return nil, fmt.Errorf("role manager emtpy.pType=%s", pType)
	}
	return roleManager.GetRoles(userKey, domain...)
}

// GetOwnersForProvider 根据rid获取用户集合
func (r *RBAC) GetOwnersForProvider(pType string, provider Owner, domain ...string) ([]string, error) {
	roleKey := provider.OwnerKey()
	if len(pType) == 0 {
		pType = "g"
	}
	roleManager := r.GetModel()["g"][pType].RM
	if roleManager == nil {
		return nil, fmt.Errorf("role manager emtpy.pType=%s", pType)
	}
	return roleManager.GetUsers(roleKey, domain...)
}

// HasProviderForOwner 判断Owner是否拥有角色
func (r *RBAC) HasProviderForOwner(pType string, owner Owner, provider Owner, domain ...string) (bool, error) {
	if len(pType) == 0 {
		pType = "g"
	}
	roles, err := r.GetProvidersForOwner(pType, owner, domain...)
	if err != nil {
		return false, err
	}
	hasRole := false
	rKey := provider.OwnerKey()
	for _, r := range roles {
		if strings.EqualFold(rKey, r) {
			hasRole = true
			break
		}
	}

	return hasRole, nil
}

// DeleteProviderForOwner 删除Owner的某个角色
func (r *RBAC) DeleteProviderForOwner(pType string, owner Owner, provider Owner, domain ...string) (bool, error) {
	if len(pType) == 0 {
		pType = "g"
	}
	args := []string{owner.OwnerKey(), provider.OwnerKey()}
	args = append(args, domain...)
	return r.RemoveNamedGroupingPolicy(pType, args)
}

// DeleteOwner 删除用户
func (r *RBAC) DeleteOwner(pType string, owner Owner, domain ...string) (bool, error) {
	if len(pType) == 0 {
		pType = "g"
	}
	var vals = make([]string, 0)
	vals = append(vals, owner.OwnerKey())
	vals = append(vals, domain...)
	return r.RemoveFilteredNamedGroupingPolicy(pType, 0, vals...)
}

// DeleteAllProvidersForOwner 删除Owner的所有角色
func (r *RBAC) DeleteAllProvidersForOwner(pType string, owner Owner, domain ...string) (bool, error) {
	if len(pType) == 0 {
		pType = "g"
	}
	var vals = make([]string, 0)
	vals = append(vals, owner.OwnerKey())
	vals = append(vals, domain...)
	return r.RemoveFilteredNamedGroupingPolicy(pType, 1, vals...)
}

// DeleteOwnerLike 删除带前缀为providerPrefix的用户
func (r *RBAC) DeleteOwnerLike(pType string, owner Owner, providerPrefix string, domain ...string) (bool, error) {
	if len(pType) == 0 {
		pType = "g"
	}
	provides, err := r.GetProvidersForOwner(pType, owner, domain...)
	if err != nil {
		return false, err
	}
	var roles = make([][]string, 0)
	for _, provide := range provides {
		if strings.HasPrefix(provide, providerPrefix) {
			var rule = make([]string, 0)
			rule = append(rule, owner.OwnerKey())
			rule = append(rule, provide)
			rule = append(rule, domain...)
			roles = append(roles, rule)
		}
	}

	return r.RemoveGroupingPolicies(roles)
	// return r.RemoveFilteredNamedGroupingPolicy(pType, 0, vals...)
}

// DeleteProvidersForOwnerLike 删除带前缀为providerPrefix的provider
func (r *RBAC) DeleteProvidersForOwnerLike(pType string, owner Owner, providerPrefix string, domain ...string) (bool, error) {
	if len(pType) == 0 {
		pType = "g"
	}
	provides, err := r.GetProvidersForOwner(pType, owner, domain...)
	if err != nil {
		return false, err
	}
	var roles = make([][]string, 0)
	for _, provide := range provides {
		if strings.HasPrefix(provide, providerPrefix) {
			var rule = make([]string, 0)
			rule = append(rule, "")
			rule = append(rule, owner.OwnerKey()) // provider位于第二位置
			rule = append(rule, provide)
			rule = append(rule, domain...)
			roles = append(roles, rule)
		}
	}

	return r.RemoveGroupingPolicies(roles)
}

// DeletePolicyForOwner 删除用户对应的权限[g]
// role_definition规则
func (r *RBAC) DeletePolicyForOwner(pType string, owner Owner, domain ...string) (bool, error) {
	if len(pType) == 0 {
		pType = "g"
	}
	userKey := owner.OwnerKey()
	var vals = make([]string, 0)
	vals = append(vals, userKey)
	vals = append(vals, domain...)
	return r.RemoveFilteredNamedGroupingPolicy(pType, 0, vals...)
}

func (r *RBAC) GetEnforcer() *casbin.SyncedEnforcer {
	return r.SyncedEnforcer
}

// ======================others======================

// GetAllResourcesByOwner 通过uid获取用户的所有资源[用于解析权限]
func (r *RBAC) GetAllResourcesByOwner(owner Owner, domain ...string) map[string]interface{} {
	userKey := owner.OwnerKey()
	allRes := make(map[string]interface{})

	myRes := r.GetEnforcer().GetPermissionsForUser(userKey)
	log().Info("解析权限", "permissions", myRes)

	// 获取用户的隐形角色
	implicitRoles, err := r.GetEnforcer().GetImplicitRolesForUser(userKey, domain...)
	if err != nil {
		return nil
	}
	for _, v := range implicitRoles {
		// 查询用户隐形角色的资源权限
		subRes := r.GetEnforcer().GetPermissionsForUser(v)
		log().Info("-------------------------------------------------")
		log().Info(fmt.Sprintf("subRes[%s], len(res)=> %d", v, len(subRes)))
		log().Info("subRes[%s], res=> %s", v, subRes)
		log().Info("-------------------------------------------------")
		allRes[v] = subRes
	}

	allRes["myRes"] = myRes
	return allRes
}
