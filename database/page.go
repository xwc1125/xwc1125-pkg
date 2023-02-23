package database

// PageWhereOrder 分页条件
type PageWhereOrder struct {
	Order string
	Where string
	Value []interface{}
}

// OrderByCol 排序信息
type OrderByCol struct {
	Column string // 排序字段
	Asc    bool   // 是否正序
}

// PageResult 分页返回数据
type PageResult struct {
	Page    *Paging     `json:"page"`    // 分页信息
	Results interface{} `json:"results"` // 数据
}

// CursorResult Cursor分页返回数据
type CursorResult struct {
	Results interface{} `json:"results"` // 数据
	Cursor  string      `json:"cursor"`  // 下一页
}

// Paging 分页请求数据
type Paging struct {
	Page  int `json:"page"`  // 页码
	Limit int `json:"limit"` // 每页条数
	Total int `json:"total"` // 总数据条数
}

func (p *Paging) Offset() int {
	offset := 0
	if p.Page > 0 {
		offset = (p.Page - 1) * p.Limit
	}
	return offset
}

func (p *Paging) TotalPage() int {
	if p.Total == 0 || p.Limit == 0 {
		return 0
	}
	totalPage := p.Total / p.Limit
	if p.Total%p.Limit > 0 {
		totalPage++
	}
	return totalPage
}
