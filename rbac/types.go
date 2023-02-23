// Package rbac
//
// @author: xwc1125
package rbac

var (
	Separator = ":"
)

// RightType 权限类型
type RightType int64

const (
	AbleAccessible = RightType(0) // 可访问
	AbleAuthorize  = RightType(1) // 可授权
)

// CasbinRule casbin规则
type CasbinRule struct {
	ID    uint   `xorm:"pk autoincr notnull" gorm:"primaryKey;autoIncrement"`
	PType string `xorm:"varchar(100) index not null default ''" gorm:"column:p_type;size:100"`
	V0    string `xorm:"varchar(100) index not null default ''" gorm:"size:100"`
	V1    string `xorm:"varchar(100) index not null default ''" gorm:"size:100"`
	V2    string `xorm:"varchar(100) index not null default ''" gorm:"size:100"`
	V3    string `xorm:"varchar(100) index not null default ''" gorm:"size:100"`
	V4    string `xorm:"varchar(100) index not null default ''" gorm:"size:100"`
	V5    string `xorm:"varchar(100) index not null default ''" gorm:"size:100"`
}

func (CasbinRule) TableName() string {
	return "sys_casbin_rule"
}

// OwnerType 持有者类型
type OwnerType interface {
	GetOwnerTypeKey() string
	GetOwnerTypeValue() string
}

type Owner interface {
	OwnerKey() string
}

// RoleDefine 角色定义
type RoleDefine struct {
	Sub string `json:"sub"` // 访问资源的用户/角色
	Dom string `json:"dom"` // 域/域租户
	Obj string `json:"obj"` // 要访问的资源/路径
	Act string `json:"act"` // 对资源访问的操作/动作
	Suf string `json:"suf"` // 附加资源
}

type Menu struct {
	MenuId   int64
	Path     string
	Action   string
	MenuType int64
}

// Permission 资源
type Permission struct {
	Id           int64       // 资源ID
	Resource     string      // 资源内容
	ResourceType string      // 资源类型
	Action       string      // 操作动作
	RightType    RightType   // 权限类型
	Extra        interface{} // 扩展内容
}
