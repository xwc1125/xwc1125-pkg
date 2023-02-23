// Package resourcetree
//
// @author: xwc1125
package resourcetree

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

type Org struct {
	Id   int64
	Pid  int64
	Name string
}

func (o *Org) IsChild(id int64) (b bool) {
	return o.Pid == id
}

func (o *Org) GetId() (i int64) {
	return o.Id
}

func TestTree(t *testing.T) {
	orgList := []TreeItem{
		&Org{Id: 12, Pid: 0, Name: "1st Layer"},
		&Org{Id: 131, Pid: 12, Name: "2nd Layer(A)"},
		&Org{Id: 132, Pid: 12, Name: "2nd Layer(B)"},
		&Org{Id: 7634, Pid: 131, Name: "3rd Layer"},
	}
	orgData, _ := json.Marshal(orgList)
	fmt.Printf("原始数据：%s\n", orgData)

	nodeTree := &Tree{
		Data:     orgList[0], // 根主节点
		Children: []*Tree{},
	}
	MakeTree(orgList, nodeTree, nil, nil)
	data, _ := json.Marshal(nodeTree)
	fmt.Printf("生成树：%s", data)
}

type Group2 struct {
	Id        int64  `json:"id"`
	Alias     string `json:"alias"`
	FullAlias string `json:"fullAlias"`

	ParentId        int64  `json:"parentId"`
	ParentAlias     string `json:"parentAlias"`
	ParentFullAlias string `json:"parentFullAlias"`
}

func (g *Group2) IsChild(id int64) (b bool) {
	return g.ParentId == id
}
func (g *Group2) GetId() (i int64) {
	return g.Id
}
func TestGroup2(t *testing.T) {
	list := []TreeItem{
		&Group2{1, "FPF引擎", "", 0, "", ""},
		&Group2{2, "IDP引擎", "", 0, "", ""},
		&Group2{3, "桂妃山", "", 1, "", ""},
		&Group2{4, "南沙", "", 1, "", ""},
		&Group2{5, "大边", "", 1, "", ""},
		&Group2{6, "一场", "", 3, "", ""},
		&Group2{9, "一线", "", 6, "", ""},
		&Group2{7, "二场", "", 3, "", ""},
		&Group2{10, "一线", "", 7, "", ""},
		&Group2{8, "洗消中心", "", 5, "", ""},
	}

	nodeTree := &Tree{
		Data:     list[0], // 根主节点
		Children: []*Tree{},
	}
	MakeTree(list, nodeTree, func(ti, tj *Tree) bool {
		data1 := ti.Data.(*Group2)
		data2 := tj.Data.(*Group2)
		return data1.Alias < data2.Alias
	}, func(parent, current *Tree) {
		parrentG := parent.Data.(*Group2)
		v := current.Data.(*Group2)

		parentAlias := parrentG.FullAlias
		if len(parentAlias) == 0 {
			parentAlias = parrentG.Alias
		}
		v.FullAlias = parentAlias + "/" + v.Alias
		v.ParentFullAlias = parentAlias
		fullList := strings.Split(v.FullAlias, "/")
		if len(fullList) >= 2 {
			v.ParentAlias = fullList[len(fullList)-2]
		} else {
			v.ParentAlias = parentAlias
		}
	})
	bytes, _ := json.MarshalIndent(nodeTree, "", "    ")
	fmt.Printf("%s\n", bytes)

	fmt.Println("OK")
}
