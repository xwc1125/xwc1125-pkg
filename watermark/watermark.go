// Package watermark
//
// @author: xwc1125
// @date: 2021/3/15
package watermark

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
	"unicode/utf8"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/xwc1125/xwc1125-pkg/utils/stringutil"
	"golang.org/x/image/font/gofont/goregular"
)

// 水印的位置
const (
	TopLeft int = iota
	TopRight
	BottomLeft
	BottomRight
	Center
)

// MaxTextLen 进行切割时的最大文字长度
var MaxTextLen = 80

// Water 水印对象
type water struct {
	pattern string // 增加按时间划分的子目录：默认没有时间划分的子目录[2006/01/02]
}

// New ...
func New(pattern string) (*water, error) {
	return &water{
		pattern: pattern,
	}, nil
}

func getTtf(ttf string) (*truetype.Font, error) {
	var (
		fontBytes []byte
		err       error
	)
	if stringutil.IsEmpty(ttf) {
		return truetype.Parse(goregular.TTF)
	} else {
		fontBytes, err = ioutil.ReadFile(ttf)
	}
	if err != nil {
		return nil, err
	}
	return freetype.ParseFont(fontBytes)
}

// WaterText 写文本
// originFileName 需要加水印的原文件
// fontInfo 文本内容
// desDir 目标文件目录
// resultFileName 结果文件名，不含路径及后缀。当空时，会自动生成
func (w *water) WaterText(originFileName string, fontInfo []FontInfo, desDir, resultFileName string) error {
	imageType, img, gif, err := w.WaterTextAndImage(originFileName, fontInfo, "", nil)
	if err != nil {
		return err
	}
	return w.Save(imageType, img, gif, desDir, resultFileName)
}

// WaterImage ...
func (w *water) WaterImage(originFileName, watermarkFileName string, position *Position, desDir, resultFileName string) error {
	imageType, img, gif, err := w.WaterTextAndImage(originFileName, nil, watermarkFileName, position)
	if err != nil {
		return err
	}
	return w.Save(imageType, img, gif, desDir, resultFileName)
}

// WaterTextAndImage ...
func (w *water) WaterTextAndImage(originFileName string, fontInfo []FontInfo, logoFileName string, logoPosition *Position) (imgType string, resultImage *image.NRGBA, resultGif *gif.GIF, err error) {
	imgFile1, _ := os.Open(originFileName)
	defer imgFile1.Close()

	// 解析图片
	_, imgType, err = image.DecodeConfig(imgFile1)
	if err != nil {
		return "", nil, nil, err
	}

	// 需要加水印的图片
	imgFile, _ := os.Open(originFileName)
	defer imgFile.Close()
	{
		if imgType == "gif" {
			resultGif, err = DecodeGif(imgFile)
			if err != nil {
				return imgType, nil, nil, err
			}
		} else {
			staticImg, err := DecodeImg(imgFile)
			if err != nil {
				return imgType, nil, nil, err
			}
			resultImage = ImageToRgba(staticImg)
		}
	}

	// 添加文字水印
	if fontInfo != nil {
		if imgType == "gif" {
			resultGif, err = WriterGifText(resultGif, fontInfo)
			if err != nil {
				return imgType, nil, nil, err
			}
		} else {
			resultImage, err = WriterText(resultImage, fontInfo)
			if err != nil {
				return imgType, nil, nil, err
			}
		}
	}

	// 添加图标logo
	// 解析水印图标
	if !stringutil.IsEmpty(logoFileName) {
		wmImgFile, _ := os.Open(logoFileName)
		defer wmImgFile.Close()
		waterImage, err := DecodeImg(wmImgFile)
		if err != nil {
			return imgType, nil, nil, err
		}
		if imgType == "gif" {
			resultGif, err = WriterGifImage(resultGif, waterImage, logoPosition)
			if err != nil {
				return imgType, nil, nil, err
			}
		} else {
			resultImage = WriterImage(resultImage, waterImage, logoPosition)
		}
	}
	return imgType, resultImage, resultGif, err
}

func (w water) Save(imgType string, resultImage *image.NRGBA, resultGif *gif.GIF, desDir, resultFileName string) error {
	var subPath string
	subPath = w.pattern
	// 创建目标目录
	dirs, err := createDir(desDir, subPath)
	if err != nil {
		return err
	}
	if stringutil.IsEmpty(resultFileName) {
		resultFileName = getRandomString(10)
	}
	newName := fmt.Sprintf("%s%s.%s", dirs, resultFileName, imgType)

	// 保存
	if resultGif != nil {
		return WriterGifToFile(resultGif, newName)
	}
	if resultImage != nil {
		return WriterToFile(imgType, resultImage, newName)
	}
	return nil
}

// DecodeImg 解析Image(png,jpeg)
func DecodeImg(imgFile *os.File) (image.Image, error) {
	img, _, err := image.Decode(imgFile)
	return img, err
}

// DecodeImg 解析gif
func DecodeGif(imgFile *os.File) (*gif.GIF, error) {
	return gif.DecodeAll(imgFile)
}

// ImageToRgba 将image.Image转换为*image.NRGBA
func ImageToRgba(img image.Image) *image.NRGBA {
	imgRgba := image.NewNRGBA(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			imgRgba.Set(x, y, img.At(x, y))
		}
	}
	return imgRgba
}

// ImageToBytes 将image.Image转换为[]bytes
func ImageToBytes(imgType string, img image.Image) ([]byte, error) {
	buf2 := new(bytes.Buffer)
	if imgType == "png" {
		if err := png.Encode(buf2, img); err != nil {
			return nil, err
		}
	} else {
		if err := jpeg.Encode(buf2, img, nil); err != nil {
			return nil, err
		}
	}
	return buf2.Bytes(), nil
}

// ImageToBytes 将image.NRGBA转换为[]bytes
func ImageNRGBAToBytes(imgType string, img *image.NRGBA) ([]byte, error) {
	buf2 := new(bytes.Buffer)
	if imgType == "png" {
		if err := png.Encode(buf2, img); err != nil {
			return nil, err
		}
	} else {
		if err := jpeg.Encode(buf2, img, nil); err != nil {
			return nil, err
		}
	}
	return buf2.Bytes(), nil
}

// WriterToFile 保存到新文件中
func WriterToFile(imgType string, img *image.NRGBA, desFileName string) error {
	newFile, err := os.Create(desFileName)
	if err != nil {
		return err
	}
	defer newFile.Close()

	if imgType == "png" {
		err = png.Encode(newFile, img)
	} else {
		err = jpeg.Encode(newFile, img, &jpeg.Options{Quality: 100})
	}
	return err
}

// WriterGifToFile 保存gif到新文件中
func WriterGifToFile(img *gif.GIF, desFileName string) error {
	newFile, err := os.Create(desFileName)
	if err != nil {
		return err
	}
	defer newFile.Close()

	return gif.EncodeAll(newFile, img)
}

// WriterText 添加文字水印函数
func WriterText(img *image.NRGBA, typeface []FontInfo) (*image.NRGBA, error) {
	var err error
	for _, t := range typeface {
		font, err := getTtf(t.Ttf)
		if err != nil {
			return nil, err
		}

		f := freetype.NewContext()
		f.SetDPI(256)
		f.SetFont(font)
		f.SetFontSize(t.Size)
		f.SetClip(img.Bounds())
		f.SetDst(img)
		f.SetSrc(image.NewUniform(color.RGBA{R: t.Position.R, G: t.Position.G, B: t.Position.B, A: t.Position.A}))

		text := t.Message
		maxTextLen := MaxTextLen
		drawArray := make([]string, 0)
		count := utf8.RuneCountInString(text) / maxTextLen
		count = count + 1
		for i := 0; i < count; i++ {
			if i == count-1 {
				drawArray = append(drawArray, text[i*maxTextLen:])
			} else {
				drawArray = append(drawArray, text[i*maxTextLen:(i+1)*maxTextLen])
			}
		}
		drawArrLen := len(drawArray)
		for i, msg := range drawArray {
			position := getTextPosition(img, f, t, msg, drawArrLen, i+1)
			pt := freetype.Pt(position.dx, position.dy)
			_, err = f.DrawString(msg, pt)
			if err != nil {
				break
			}
		}
	}
	return img, err
}

func WriterTexts(img *image.NRGBA, fInfo FontInfo, texts []string) (*image.NRGBA, error) {
	var err error
	font := fInfo.Font
	for _, t := range texts {
		if fInfo.Font == nil {
			font, err = getTtf(fInfo.Ttf)
			if err != nil {
				return nil, err
			}
		}

		f := freetype.NewContext()
		f.SetDPI(256)
		f.SetFont(font)
		f.SetFontSize(fInfo.Size)
		f.SetClip(img.Bounds())
		f.SetDst(img)
		f.SetSrc(image.NewUniform(color.RGBA{R: fInfo.Position.R, G: fInfo.Position.G, B: fInfo.Position.B, A: fInfo.Position.A}))

		text := t
		maxTextLen := MaxTextLen
		drawArray := make([]string, 0)
		count := utf8.RuneCountInString(text) / maxTextLen
		count = count + 1
		for i := 0; i < count; i++ {
			if i == count-1 {
				drawArray = append(drawArray, text[i*maxTextLen:])
			} else {
				drawArray = append(drawArray, text[i*maxTextLen:(i+1)*maxTextLen])
			}
		}
		drawArrLen := len(drawArray)
		for i, msg := range drawArray {
			position := getTextPosition(img, f, fInfo, msg, drawArrLen, i+1)
			pt := freetype.Pt(position.dx, position.dy)
			_, err = f.DrawString(msg, pt)
			if err != nil {
				break
			}
		}
	}
	return img, err
}

// Dxy ...
type Dxy struct {
	dx int
	dy int
}

// getTextPosition获取文本位置
// indexSec 大文本的第几个分片[从1开始]
func getTextPosition(img *image.NRGBA, f *freetype.Context, t FontInfo, msg string, totalSec, indexSec int) *Dxy {
	dxy := new(Dxy)

	// 获取字体的尺寸大小
	fixed := f.PointToFixed(t.Size)
	ceil := fixed.Ceil()
	if msg == "" {
		msg = t.Message
	}
	msgLen := (utf8.RuneCountInString(msg) / 2) * ceil
	switch int(t.Position.Position) {
	case 0:
		dxy.dx = t.Position.Dx
		dxy.dy = t.Position.Dy + ceil + (indexSec-1)*ceil
	case 1:
		dxy.dx = img.Bounds().Dx() - msgLen - t.Position.Dx
		dxy.dy = t.Position.Dy + ceil + (indexSec-1)*ceil
	case 2:
		dxy.dx = t.Position.Dx
		dxy.dy = img.Bounds().Dy() - ceil - t.Position.Dy - (totalSec-indexSec-1)*ceil
	case 3:
		dxy.dx = img.Bounds().Dx() - msgLen - t.Position.Dx
		dxy.dy = img.Bounds().Dy() - ceil - t.Position.Dy - (totalSec-indexSec-1)*ceil
	case 4:
		dxy.dx = (img.Bounds().Dx()-msgLen)/2 + t.Position.Dx + (indexSec-1)*ceil
		dxy.dy = (img.Bounds().Dy()-ceil)/2 + t.Position.Dy + (indexSec-1)*ceil
	default:
	}
	return dxy
}

func getImagePosition(img *image.NRGBA, watermark image.Image, position *Position) *Dxy {
	dxy := new(Dxy)
	// 左上角为(0,0),往下x+，往右y+
	switch int(position.Position) {
	case 0:
		dxy.dx = position.Dx + position.Dx
		dxy.dy = position.Dy + position.Dy
	case 1:
		dxy.dx = img.Bounds().Dx() - watermark.Bounds().Dx() - position.Dx
		dxy.dy = position.Dy + position.Dy
	case 2:
		dxy.dx = position.Dx + position.Dx
		dxy.dy = img.Bounds().Dy() - watermark.Bounds().Dy() - position.Dy
	case 3:
		dxy.dx = img.Bounds().Dx() - watermark.Bounds().Dx() - position.Dx
		dxy.dy = img.Bounds().Dy() - watermark.Bounds().Dy() - position.Dy
	case 4:
		dxy.dx = (img.Bounds().Dx()-watermark.Bounds().Dx())/2 + position.Dx
		dxy.dy = (img.Bounds().Dy()-watermark.Bounds().Dy())/2 + position.Dy
	default:
	}
	return dxy
}

// WriterGifText 添加文字水印函数
func WriterGifText(gifImg2 *gif.GIF, typeface []FontInfo) (*gif.GIF, error) {
	var err error

	gifs := make([]*image.Paletted, 0)
	x0 := 0
	y0 := 0
	yuan := 0
	for k, gifImg := range gifImg2.Image {
		img := image.NewNRGBA(gifImg.Bounds())
		if k == 0 {
			x0 = img.Bounds().Dx()
			y0 = img.Bounds().Dy()
		}

		if k == 0 && gifImg2.Image[k+1].Bounds().Dx() > x0 && gifImg2.Image[k+1].Bounds().Dy() > y0 {
			yuan = 1
			break
		}
		if x0 == img.Bounds().Dx() && y0 == img.Bounds().Dy() {
			for y := 0; y < img.Bounds().Dy(); y++ {
				for x := 0; x < img.Bounds().Dx(); x++ {
					img.Set(x, y, gifImg.At(x, y))
				}
			}
			// todo 使用callback
			img, err = WriterText(img, typeface) // 添加文字水印
			if err != nil {
				break
			}
			// 定义一个新的图片调色板img.Bounds()：使用原图的颜色域，gifimg.Palette：使用原图的调色板
			p1 := image.NewPaletted(gifImg.Bounds(), gifImg.Palette)
			// 把绘制过文字的图片添加到新的图片调色板上
			draw.Draw(p1, gifImg.Bounds(), img, image.ZP, draw.Src)
			// 把添加过文字的新调色板放入调色板slice
			gifs = append(gifs, p1)
		} else {
			gifs = append(gifs, gifImg)
		}
	}
	if yuan == 1 {
		return nil, errors.New("gif: image block is out of bounds")
	} else {
		if err != nil {
			return nil, err
		}

		g1 := &gif.GIF{
			Image:     gifs,
			Delay:     gifImg2.Delay,
			LoopCount: gifImg2.LoopCount,
		}
		return g1, nil
	}
}

// WriterImage 往图片上加上图片水印
func WriterImage(img *image.NRGBA, watermark image.Image, position *Position) *image.NRGBA {
	log().Trace("start witer logo watermark")
	dxy := getImagePosition(img, watermark, position)
	// 把水印写到右下角，并向0坐标各偏移10个像素
	offset := image.Pt(
		dxy.dx,
		dxy.dy,
	)
	b := img.Bounds()
	m := image.NewNRGBA(b)

	draw.Draw(m, b, img, image.ZP, draw.Src)
	draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)
	log().Trace("end witer logo watermark")
	return m
}

// WriterGifImage 添加图片水印函数
func WriterGifImage(gifImg2 *gif.GIF, waterImage image.Image, position *Position) (*gif.GIF, error) {
	var err error
	gifs := make([]*image.Paletted, 0)
	x0 := 0
	y0 := 0
	yuan := 0
	for k, gifImg := range gifImg2.Image {
		img := image.NewNRGBA(gifImg.Bounds())
		if k == 0 {
			x0 = img.Bounds().Dx()
			y0 = img.Bounds().Dy()
		}

		if k == 0 && gifImg2.Image[k+1].Bounds().Dx() > x0 && gifImg2.Image[k+1].Bounds().Dy() > y0 {
			yuan = 1
			break
		}
		if x0 == img.Bounds().Dx() && y0 == img.Bounds().Dy() {
			for y := 0; y < img.Bounds().Dy(); y++ {
				for x := 0; x < img.Bounds().Dx(); x++ {
					img.Set(x, y, gifImg.At(x, y))
				}
			}
			// todo 使用callback
			img = WriterImage(img, waterImage, position) // 添加image水印
			// 定义一个新的图片调色板img.Bounds()：使用原图的颜色域，gifimg.Palette：使用原图的调色板
			p1 := image.NewPaletted(gifImg.Bounds(), gifImg.Palette)
			// 把绘制过文字的图片添加到新的图片调色板上
			draw.Draw(p1, gifImg.Bounds(), img, image.ZP, draw.Src)
			// 把添加过文字的新调色板放入调色板slice
			gifs = append(gifs, p1)
		} else {
			gifs = append(gifs, gifImg)
		}
	}
	if yuan == 1 {
		return nil, errors.New("gif: image block is out of bounds")
	} else {
		if err != nil {
			return nil, err
		}

		g1 := &gif.GIF{
			Image:     gifs,
			Delay:     gifImg2.Delay,
			LoopCount: gifImg2.LoopCount,
		}
		return g1, nil
	}
}

// gifFontWater gif图片水印
func gifFontWater(file string, typeface []FontInfo) (g *gif.GIF, err error) {
	imgFile, _ := os.Open(file)
	defer imgFile.Close()

	gifImg2, _ := DecodeGif(imgFile)
	return WriterGifText(gifImg2, typeface)
}

// staticFontWater png,jpeg图片水印
func staticFontWater(file string, typeface []FontInfo) (*image.NRGBA, error) {
	// 需要加水印的图片
	imgFile, _ := os.Open(file)
	defer imgFile.Close()

	staticImg, err := DecodeImg(imgFile)
	if err != nil {
		return nil, err
	}
	img := ImageToRgba(staticImg)
	return WriterText(img, typeface) // 添加文字水印
}

// FontInfo 定义添加的文字信息
type FontInfo struct {
	Ttf      string         `json:"ttf" mapstructure:"ttf"`           // 文字字体
	Font     *truetype.Font `json:"font" mapstructure:"font"`         // 文字字体
	Size     float64        `json:"size" mapstructure:"size"`         // 文字大小
	Message  string         `json:"message" mapstructure:"message"`   // 文字内容
	Position Position       `json:"position" mapstructure:"position"` // 文字存放位置
}

// Position ...
type Position struct {
	Position int   `json:"position" mapstructure:"position"` // 文字存放位置
	Dx       int   `json:"dx" mapstructure:"dx"`             // 文字x轴留白距离
	Dy       int   `json:"dy" mapstructure:"dy"`             // 文字y轴留白距离
	R        uint8 `json:"r" mapstructure:"r"`               // 文字颜色值RGBA中的R值
	G        uint8 `json:"g" mapstructure:"g"`               // 文字颜色值RGBA中的G值
	B        uint8 `json:"b" mapstructure:"b"`               // 文字颜色值RGBA中的B值
	A        uint8 `json:"a" mapstructure:"a"`               // 文字颜色值RGBA中的A值
}

// getRandomString 生成图片名字
func getRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	bytesLen := len(bytes)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(bytesLen)])
	}
	return string(result)
}

// createDir 检查并生成存放图片的目录
func createDir(SavePath, subPath string) (string, error) {
	var dirs string
	if subPath == "" {
		dirs = fmt.Sprintf("%s/", SavePath)
	} else {
		format := time.Now().Format(subPath)
		dirs = fmt.Sprintf("%s/%s/", SavePath, format)
	}

	_, err := os.Stat(dirs)
	if err != nil {
		err = os.MkdirAll(dirs, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return dirs, nil
}
