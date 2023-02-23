// Package resourcetree
//
// @author: xwc1125
package resourcetree

import "sort"

type (
	// TreeItem 数据接口类型
	TreeItem interface {
		// IsChild 判断是否是指定id的子节点数据
		IsChild(id int64) bool
		// GetId 获取数据id
		GetId() int64
	}

	// Tree 层级结构树
	Tree struct {
		Data     TreeItem `json:"data"`
		Children []*Tree  `json:"children"`
	}
)

// MakeTree 生成层级结构树
// 最初的参数node必须是根主节点
func MakeTree(vList []TreeItem, node *Tree, sortFunc func(ti, tj *Tree) bool, callback func(parent, current *Tree)) (remainList []TreeItem) {
	return makeTree(vList, node, sortFunc, callback)
}

func makeTree(vList []TreeItem, node *Tree, sortFunc func(ti, tj *Tree) bool, callback func(parent, current *Tree)) (remainList []TreeItem) {
	// 返回数据集合中某个节点下的直属子节点集合
	var (
		childs []*Tree
	)
	childs, vList = findChild(node, vList, callback)
	for _, child := range childs {
		// 将直属子节点加入树的当前节点下
		node.Children = append(node.Children, child)
		if sortFunc != nil {
			// 用 alias 排序，alias相等的元素保持原始顺序
			sort.SliceStable(node.Children, func(i, j int) bool {
				return sortFunc(node.Children[i], node.Children[j])
			})
		}
		// 判断该子节点下是否还有下属，如果有，继续递归
		if hasChild(child, vList) {
			vList = makeTree(vList, child, sortFunc, callback)
		}
	}
	if sortFunc != nil {
		// 用 alias 排序，alias相等的元素保持原始顺序
		sort.SliceStable(childs, func(i, j int) bool {
			return sortFunc(childs[i], childs[j])
		})
	}
	return vList
}

// 返回数据集合中某个节点下的直属子节点集合
func findChild(node *Tree, vList []TreeItem, callback func(parent *Tree, current *Tree)) (ret []*Tree, remainList []TreeItem) {
	var delIds = make([]int, 0, len(vList))
	for i, v := range vList {
		if v.IsChild(node.Data.GetId()) {
			current := &Tree{
				Data:     v,
				Children: []*Tree{},
			}
			if callback != nil {
				callback(node, current)
			}
			ret = append(ret, current)
			delIds = append(delIds, i)
		}
	}
	for i := len(delIds) - 1; i >= 0; i-- {
		i2 := len(vList)
		i3 := delIds[i]
		if i3 == (i2 - 1) {
			// 最后一个
			vList = vList[:i3]
		} else {
			vList = append(vList[:i3], vList[i3+1:]...)

		}
	}
	return ret, vList
}

// 判断该子节点下是否还有下属
func hasChild(node *Tree, vList []TreeItem) (has bool) {
	for _, v := range vList {
		if v.IsChild(node.Data.GetId()) {
			has = true
			break
		}
	}
	return has
}
