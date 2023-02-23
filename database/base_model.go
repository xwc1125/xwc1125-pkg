// Package database
//
// @author: xwc1125
package database

// StatusType 状态：-1:已删除，0:被禁用，1:正常，2:未审批
type StatusType int64

const (
	StatusType_Delete = iota - 1
	StatusType_Disable
	StatusType_Ok
	StatusType_NotApproved
)

type ModelID struct {
	Id int64 `json:"id" form:"id" xorm:"pk autoincr bigint(20) notnull 'id' comment('ID')" gorm:"type:bigint(20);primaryKey;autoIncrement;not null;column:id;comment:ID"` // ID
}
type ModelTime struct {
	CreatedAt int64 `json:"createdAt" xorm:"bigint(15) notnull default(0) created comment('创建时间(毫秒)')" gorm:"type:bigint(15);not null;autoCreateTime:milli;comment:创建时间(毫秒)"` // 创建时间(毫秒)
	UpdatedAt int64 `json:"updatedAt" xorm:"bigint(15) notnull default(0) updated comment('更新时间(毫秒)')" gorm:"type:bigint(15);not null;autoUpdateTime:milli;comment:更新时间(毫秒)"` // 更新时间(毫秒)
	DeletedAt int64 `json:"-" xorm:"bigint(15) deleted comment('删除时间(毫秒)')" gorm:"type:bigint(15);index;comment:删除时间(毫秒)"`                                                    // 删除时间(毫秒)
}
type ModelStatus struct {
	Status int64  `json:"status" binding:"min=-1,max=2" xorm:"tinyint(8) notnull default(1) comment('状态：-1:已删除，0:被禁用，1:正常，2:未审批')" gorm:"type:tinyint(4);not null;default:1;content:状态：-1:已删除，0:被禁用，1:正常，2:未审批"` // 状态：状态：-1:已删除，0:被禁用，1:正常，2:未审批
	Remark string `json:"remark" xorm:"varchar(500) comment('备注')" gorm:"size:500;content:备注"`                                                                                                                   // 备注
}

type ModelBy struct {
	CreateBy int64 `json:"createBy" xorm:"bigint(20) comment('创建者')" gorm:"index;size:20;comment:创建者"` // 创建者
	UpdateBy int64 `json:"updateBy" xorm:"bigint(20) comment('更新者')" gorm:"index;size:20;comment:更新者"` // 更新者
}

// SetCreateBy 设置创建人id
func (e *ModelBy) SetCreateBy(createBy int64) {
	e.CreateBy = createBy
}

// SetUpdateBy 设置修改人id
func (e *ModelBy) SetUpdateBy(updateBy int64) {
	e.UpdateBy = updateBy
}
