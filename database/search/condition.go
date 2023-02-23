package search

import "strings"

const (
	TypeTag   = "type"   // 类型标签
	ColumnTag = "column" // 字段标签
	TableTag  = "table"  // 表标签
	OnTag     = "on"     // on标签
	JoinTag   = "join"   // 连接标签
	SkipTag   = "-"      // 连接标签

	LeftTypeTag        = "left"
	ExactTypeTag       = "exact"
	IExactTypeTag      = "iexact"
	ContainsTypeTag    = "contains"
	IContainsTypeTag   = "icontains"
	GtTypeTag          = "gt"  // 大于
	GteTypeTag         = "gte" // 大于等于
	LtTypeTag          = "lt"  // 小于
	LteTypeTag         = "lte" // 小于等于
	StartsWithTypeTag  = "startswith"
	IStartsWithTypeTag = "istartswith"
	EndWithTypeTag     = "endswith"
	IEndWithTypeTag    = "iendswith"
	IsNullTypeTag      = "isnull"
	InTypeTag          = "in"
	OrderTypeTag       = "order"
	DescTypeTag        = "desc"
	AscTypeTag         = "asc"
)

// Condition 条件
type Condition interface {
	SetWhere(k string, v []interface{})
	SetOr(k string, v []interface{})
	SetOrder(k string)
	SetJoinOn(t, on string) Condition
}

// OrmCondition Orm查询的条件
type OrmCondition struct {
	OrmWhere
	Join []*OrmJoin
}

func NewOrmCondition() *OrmCondition {
	return &OrmCondition{
		OrmWhere: OrmWhere{},
		Join:     make([]*OrmJoin, 0),
	}
}

// OrmWhere 查询条件，where、order、or
type OrmWhere struct {
	Where map[string][]interface{}
	Order []string
	Or    map[string][]interface{}
}

// OrmJoin 连接标签
type OrmJoin struct {
	Type   string
	JoinOn string
	OrmWhere
}

// SetJoinOn 设置连接
func (e *OrmJoin) SetJoinOn(t, on string) Condition {
	return nil
}

// SetWhere 设置where
func (e *OrmWhere) SetWhere(k string, v []interface{}) {
	if e.Where == nil {
		e.Where = make(map[string][]interface{})
	}
	e.Where[k] = v
}

// SetOr 设置or
func (e *OrmWhere) SetOr(k string, v []interface{}) {
	if e.Or == nil {
		e.Or = make(map[string][]interface{})
	}
	e.Or[k] = v
}

// SetOrder 设置order
func (e *OrmWhere) SetOrder(k string) {
	if e.Order == nil {
		e.Order = make([]string, 0)
	}
	e.Order = append(e.Order, k)
}

// SetJoinOn 设置join on
func (e *OrmCondition) SetJoinOn(t, on string) Condition {
	if e.Join == nil {
		e.Join = make([]*OrmJoin, 0)
	}
	join := &OrmJoin{
		Type:     t,
		JoinOn:   on,
		OrmWhere: OrmWhere{},
	}
	e.Join = append(e.Join, join)
	return join
}

type resolveSearchTag struct {
	Type   string
	Column string
	Table  string
	On     []string
	Join   string
}

// makeTag 解析search的tag标签
func makeTag(tag string) *resolveSearchTag {
	r := &resolveSearchTag{}
	tags := strings.Split(tag, ";")
	var ts []string
	for _, t := range tags {
		ts = strings.Split(t, ":")
		if len(ts) == 0 {
			continue
		}
		switch ts[0] {
		case TypeTag:
			if len(ts) > 1 {
				r.Type = ts[1]
			}
		case ColumnTag:
			if len(ts) > 1 {
				r.Column = ts[1]
			}
		case TableTag:
			if len(ts) > 1 {
				r.Table = ts[1]
			}
		case OnTag:
			if len(ts) > 1 {
				r.On = ts[1:]
			}
		case JoinTag:
			if len(ts) > 1 {
				r.Join = ts[1]
			}
		}
	}
	return r
}
