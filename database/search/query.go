package search

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	// QueryTag tag标记
	QueryTag = "search"
)

const (
	// Mysql 数据库标识
	Mysql = "mysql"
	// Postgres 数据库标识
	Postgres = "postgres"
)

// ResolveSearchQuery 解析
/**
 * 	exact / iexact 等于
 * 	contains / icontains 包含
 *	gt / gte 大于 / 大于等于
 *	lt / lte 小于 / 小于等于
 *	startswith / istartswith 以…起始
 *	endswith / iendswith 以…结束
 *	in
 *	isnull
 *  order 排序		e.g. order[key]=desc     order[key]=asc
 */
func ResolveSearchQuery(driver string, q interface{}, condition Condition) {
	qType := reflect.TypeOf(q)
	qValue := reflect.ValueOf(q)
	var tag string
	var ok bool
	var t *resolveSearchTag
	for i := 0; i < qType.NumField(); i++ {
		tag, ok = "", false
		tag, ok = qType.Field(i).Tag.Lookup(QueryTag)
		if !ok {
			// 递归调用
			ResolveSearchQuery(driver, qValue.Field(i).Interface(), condition)
			continue
		}
		switch tag {
		case SkipTag:
			continue
		}
		t = makeTag(tag)
		if qValue.Field(i).IsZero() {
			continue
		}
		// 解析
		switch t.Type {
		case LeftTypeTag:
			// 左关联
			join := condition.SetJoinOn(t.Type, fmt.Sprintf(
				"left join `%s` on `%s`.`%s` = `%s`.`%s`",
				t.Join,
				t.Join,
				t.On[0],
				t.Table,
				t.On[1],
			))
			ResolveSearchQuery(driver, qValue.Field(i).Interface(), join)
		case ExactTypeTag, IExactTypeTag:
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` = ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case ContainsTypeTag, IContainsTypeTag:
			// fixme mysql不支持ilike
			if driver == Postgres && t.Type == IContainsTypeTag {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` ilike ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String() + "%"})
			} else {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` like ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String() + "%"})
			}
		case GtTypeTag:
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` > ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case GteTypeTag:
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` >= ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case LtTypeTag:
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` < ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case LteTypeTag:
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` <= ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case StartsWithTypeTag, IStartsWithTypeTag:
			if driver == Postgres && t.Type == IStartsWithTypeTag {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` ilike ?", t.Table, t.Column), []interface{}{qValue.Field(i).String() + "%"})
			} else {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` like ?", t.Table, t.Column), []interface{}{qValue.Field(i).String() + "%"})
			}
		case EndWithTypeTag, IEndWithTypeTag:
			if driver == Postgres && t.Type == IEndWithTypeTag {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` ilike ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String()})
			} else {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` like ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String()})
			}
		case InTypeTag:
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` in (?)", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case IsNullTypeTag:
			if !(qValue.Field(i).IsZero() && qValue.Field(i).IsNil()) {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` isnull", t.Table, t.Column), make([]interface{}, 0))
			}
		case OrderTypeTag:
			switch strings.ToLower(qValue.Field(i).String()) {
			case DescTypeTag, AscTypeTag:
				condition.SetOrder(fmt.Sprintf("`%s`.`%s` %s", t.Table, t.Column, qValue.Field(i).String()))
			}
		}
	}
}
