// Package db_gorm
//
// @author: xwc1125
package db_gorm

import (
	"errors"

	"github.com/xwc1125/xwc1125-pkg/types"
	"github.com/xwc1125/xwc1125-pkg/types/response"
	"gorm.io/gorm"
)

func UpdateDBErr(tx *gorm.DB) types.JsonError {
	if err := QueryDBErr(tx); err != nil {
		return err
	}
	affected := tx.RowsAffected
	if affected == 0 {
		return response.NoPermission
	}
	return nil
}

func QueryDBErr(tx *gorm.DB) types.JsonError {
	err := tx.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		return response.NewStatusErr(err)
	}
	if err != nil {
		return response.NewStatusErr(err)
	}
	return nil
}
