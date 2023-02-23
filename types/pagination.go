package types

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/chain5j/chain5j-pkg/util/dateutil"
)

type PageQuery interface {
	GetNeedSearch() interface{}
	GetPageIndex() int64
	GetPageSize() int64
	GetOffset() int64
	GetStartTime() int64
	GetEndTime() int64
	GetOrder() string
}

// Pagination 分页参数
type Pagination struct {
	// 用于分页设置的参数
	PageIndex int64 `json:"pageIndex" form:"pageIndex"`    // 当前页
	PageSize  int64 `json:"pageSize" form:"pageSize"`      // 每页的数量
	offset    int64 `json:"start,omitempty" form:"offset"` // 查询起始的条数ID

	OrderName   string `json:"orderName,omitempty" form:"orderName"` // 用于指定的排序
	orderParams string

	// 时间范围
	StartTime string `json:"startTime,omitempty" form:"startTime"` // 2006-01-02
	EndTime   string `json:"endTime,omitempty" form:"endTime"`     // 2006-01-02

	startTimeMs int64
	endTimeMs   int64

	ParamMap map[string]interface{}
}

func (p *Pagination) GetPageIndex() int64 {
	if p.PageIndex <= 0 {
		p.PageIndex = 1
	}
	return p.PageIndex
}

func (p *Pagination) GetPageSize() int64 {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 500 {
		p.PageSize = 500
	}
	return p.PageSize
}

func (p *Pagination) GetOffset() int64 {
	if p.offset > 0 {
		return p.offset
	}
	p.offset = (p.GetPageIndex() - 1) * p.GetPageSize()
	if p.offset < 0 {
		p.offset = 0
	}
	return p.offset
}

// GetStartTime 获取起始时间的毫秒值
func (p *Pagination) GetStartTime() int64 {
	if p.startTimeMs > 0 {
		return p.startTimeMs
	}
	if len(p.StartTime) == 0 {
		return 0
	}
	local, err := dateutil.ParseIn(p.StartTime+" 00:00:00", dateutil.SysTimeLocation)
	if err != nil {
		return 0
	}
	p.startTimeMs = local.UnixMilli()
	return p.startTimeMs
}

// GetEndTime 获取结束时间的毫秒值
func (p *Pagination) GetEndTime() int64 {
	if p.endTimeMs > 0 {
		return p.endTimeMs
	}
	if len(p.EndTime) == 0 {
		return 0
	}
	local, err := dateutil.ParseIn(p.EndTime+" 23:59:59", dateutil.SysTimeLocation)
	if err != nil {
		return 0
	}
	p.endTimeMs = local.UnixMilli()
	return p.endTimeMs
}

// GetOrder 获取排序组合
func (p *Pagination) GetOrder() string {
	if len(p.orderParams) > 0 {
		return p.orderParams
	}
	orderName := p.OrderName
	if len(strings.TrimSpace(orderName)) == 0 {
		return ""
	}
	orderArray := strings.Split(orderName, ",")
	if len(orderArray) == 0 {
		return ""
	}
	var orderBuff bytes.Buffer
	for _, orderOne := range orderArray {
		orderOne = url.QueryEscape(orderOne)
		if strings.HasPrefix(orderOne, "-") {
			orderBuff.WriteString(orderOne[1:] + " desc")
			orderBuff.WriteString(",")
		} else if strings.HasPrefix(orderOne, "+") {
			orderBuff.WriteString(orderOne[1:] + " asc")
			orderBuff.WriteString(",")
		} else {
			orderBuff.WriteString(orderOne + " asc")
			orderBuff.WriteString(",")
		}
	}
	order := orderBuff.String()
	if len(order) > 0 {
		p.orderParams = order[:len(order)-1]
	}
	return p.orderParams
}
