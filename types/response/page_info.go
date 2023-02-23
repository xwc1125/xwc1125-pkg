package response

type PageInfo struct {
	PageIndex int64 `json:"pageIndex"` // 当前页
	PageSize  int64 `json:"pageSize"`  // 每页的数量

	Total int64       `json:"total"`           // 总条数
	List  interface{} `json:"list"`            // 数据集
	Pages int64       `json:"pages,omitempty"` // 总页数

	Filter   interface{} `json:"filter,omitempty"`   // 如果为空置则忽略字段
	PrePage  int64       `json:"prePage,omitempty"`  // 前一页
	NextPage int64       `json:"nextPage,omitempty"` // 前一页
}

type PageFilter struct {
	Key      string `json:"key"`                // 主键
	Label    string `json:"label"`              // 便签
	Show     int    `json:"show"`               // 0:不显示给用户看，1:默认显示，2:默认不显示
	Sortable bool   `json:"sortable,omitempty"` // 是否支持排序
}
