// Package excelutil
//
// @author: xwc1125
package excelutil

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/tealeg/xlsx"
)

type XlsxRow struct {
	Row  *xlsx.Row
	Data []string
}

func NewRow(row *xlsx.Row, data []string) *XlsxRow {
	return &XlsxRow{
		Row:  row,
		Data: data,
	}
}

func (row *XlsxRow) SetRowTitle() error {
	return generateRow(row.Row, row.Data)
}

func (row *XlsxRow) GenerateRow() error {
	return generateRow(row.Row, row.Data)
}

func generateRow(row *xlsx.Row, rowStr []string) error {
	if rowStr == nil {
		return errors.New("no data to generate xlsx!")
	}
	for _, v := range rowStr {
		cell := row.AddCell()
		cell.SetString(v)
	}
	return nil
}

func Write(w http.ResponseWriter, req *http.Request, fileName string) {
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s", "file.xls")) // 文件名
	w.Header().Set("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	// 最主要的一句
	http.ServeFile(w, req, fileName)
}
