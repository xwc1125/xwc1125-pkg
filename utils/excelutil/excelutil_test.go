// Package excelutil
//
// @author: xwc1125
package excelutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/tealeg/xlsx"
)

type Person struct {
	Name    string
	Gender  string
	Age     string
	Tel     string
	Marry   string
	Address string
}

func TestXlsxRow_GenerateRow(t *testing.T) {
	peo := []Person{
		{
			Name:    "1",
			Gender:  "男",
			Age:     "18",
			Tel:     "18888888888",
			Marry:   "已婚",
			Address: "中国",
		},
		{
			Name:    "2",
			Gender:  "男",
			Age:     "18",
			Tel:     "18888888888",
			Marry:   "已婚",
			Address: "中国",
		},
	}
	GeneratePeopleExcel(peo)
}

func GeneratePeopleExcel(peo []Person) (error, bool) {
	t1 := time.Now()
	defer func() {
		fmt.Println(time.Since(t1))
	}()
	t := make([]string, 0)
	t = append(t, "姓名")
	t = append(t, "性别")
	t = append(t, "年龄")
	t = append(t, "电话")
	t = append(t, "婚配")
	t = append(t, "现居地")
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("sheet")
	if err != nil {
		return err, false
	}
	titleRow := sheet.AddRow()
	xlsRow := NewRow(titleRow, t)
	err = xlsRow.SetRowTitle()
	if err != nil {
		return err, false
	}
	for _, v := range peo {
		currentRow := sheet.AddRow()
		tmp := make([]string, 0)
		tmp = append(tmp, v.Name)
		tmp = append(tmp, v.Gender)
		tmp = append(tmp, v.Age)
		tmp = append(tmp, v.Tel)
		tmp = append(tmp, v.Marry)
		tmp = append(tmp, v.Address)

		xlsRow := NewRow(currentRow, tmp)
		err := xlsRow.GenerateRow()
		if err != nil {
			return err, false
		}
	}
	err = file.Save("./Excel/人员信息.xlsx")
	if err != nil {
		return err, false
	}
	return nil, true
}
