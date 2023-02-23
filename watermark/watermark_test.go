// Package watermark
//
// @author: xwc1125
// @date: 2021/3/15
package watermark

import (
	"fmt"
	"testing"

	"github.com/chain5j/logger"
	"github.com/chain5j/logger/zap"
)

func init() {
	zap.InitWithConfig(&logger.LogConfig{
		Console: logger.ConsoleLogConfig{
			Level:    logger.LvlDebug,
			Modules:  "*",
			ShowPath: false,
			Format:   "",
			UseColor: true,
			Console:  true,
		},
		File: logger.FileLogConfig{},
	})
}

func TestWaterText(t *testing.T) {
	SavePath := "./logs/"
	// 水印1
	drawStr := "https://www.baidu.com"

	str := FontInfo{
		// Ttf:     "./Songti.ttc",
		Font:    nil,
		Size:    6,
		Message: drawStr,
		Position: Position{
			Position: TopRight,
			Dx:       0,
			Dy:       10,
			R:        255,
			G:        0,
			B:        0,
			A:        255,
		}}
	arr := make([]FontInfo, 0)
	arr = append(arr, str)
	// 水印2
	// str2 := FontInfo{Size: 24, Message: "努力向上，涨工资", Position: Position{
	//	Position: TopLeft,
	//	Dx:       20,
	//	Dy:       40,
	//	R:        255,
	//	G:        0,
	//	B:        0,
	//	A:        255,
	// }}
	// arr = append(arr, str2)

	// 加水印图片路径
	fileName := "../logo.png"

	w, err := New("2006/01/02")
	if err != nil {
		fmt.Println(err)
	}
	err = w.WaterText(fileName, arr, SavePath, "1")
	if err != nil {
		fmt.Println(err)
	}
	// err = w.WaterText(fileName, arr, SavePath, "")
	// if err != nil {
	//	fmt.Println(err)
	// }
}

func TestWaterImage(t *testing.T) {
	SavePath := "./logs/"
	// 水印1
	position := &Position{
		Position: TopRight,
		Dx:       20,
		Dy:       20,
		R:        255,
		G:        0,
		B:        0,
		A:        255,
	}

	// 加水印图片路径
	fileName := "./logs/test.jpeg"
	logoFileName := "../logo.png"

	w, err := New("2006/01/02")
	if err != nil {
		fmt.Println(err)
	}
	err = w.WaterImage(fileName, logoFileName, position, SavePath, "1")
	if err != nil {
		fmt.Println(err)
	}
}

func TestWaterImageAndText(t *testing.T) {
	SavePath := "./logs/"
	drawStr := "https://www.baidu.com"

	arr := make([]FontInfo, 0)
	// 水印1
	str := FontInfo{
		// Ttf: "./Songti.ttc",
		Font:    nil,
		Size:    6,
		Message: drawStr,
		Position: Position{
			Position: TopLeft,
			Dx:       20,
			Dy:       20,
			R:        255,
			G:        0,
			B:        0,
			A:        255,
		}}
	arr = append(arr, str)

	// 加水印图片路径
	fileName := "./logs/test.jpeg"

	w, err := New("2006/01/02")
	if err != nil {
		fmt.Println(err)
	}

	// 水印1
	position := &Position{
		Position: TopRight,
		Dx:       20,
		Dy:       20,
		R:        255,
		G:        0,
		B:        0,
		A:        255,
	}
	logoFileName := "../logo.png"

	MaxTextLen = 90
	imgType, resultImage, resultGif, err := w.WaterTextAndImage(fileName, arr,
		logoFileName, position)
	err = w.Save(imgType, resultImage, resultGif, SavePath, "1234")
	if err != nil {
		fmt.Println(err)
	}
}
