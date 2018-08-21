package utils

import (
	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
	"log"
)

// 初始化 pdf合同模版的基础参数
type XJDPdfConfig struct {
	PageWidth     float64 // 页宽
	PageHeight    float64 // 页高
	Side          float64 // 边
	Top           float64 // 头部
	Bottom        float64 // 底部
	TableLineSize float64 // 表格尺寸大小
	LineHeight    float64 // pdf文字行间距
}

const FONT_NAME = "XJD"                    //普通字体
const FONT_PATH = "./font/SIMYOU.TTF"      //微软字体（默认字体）
const FONT_PATH_BOLD = "./font/msyhbd.ttf" //加粗字体
const FONT_NAME_BOLD = "XJDB"              //加粗字体
const DEFAULT_FONT_SIZE = 13               //中文字，高和宽等于这个值，英文为一半

type XJDPdf struct {
	gopdf.GoPdf
	Config       *XJDPdfConfig
	CurrFontName string
}

var (
	defaultBr          float64
	defaultLineHight   float64
	defaultLineFontNum int
)

// 生成一份合同，对象
func NewXJDPdf(config *XJDPdfConfig) *XJDPdf {
	if config == nil {
		config = &XJDPdfConfig{PageWidth: 595.28, PageHeight: 841.89, Side: 50, Top: 60, Bottom: 60, LineHeight: 1, TableLineSize: 0.5} //A4纸
	}
	defaultLineHight = float64(float64(DEFAULT_FONT_SIZE) + float64(DEFAULT_FONT_SIZE)/5)
	defaultBr = float64(float64(DEFAULT_FONT_SIZE) + float64(DEFAULT_FONT_SIZE)/2)
	num := (config.PageWidth - config.Side*2) / DEFAULT_FONT_SIZE
	defaultLineFontNum = int(num)
	pdf := &XJDPdf{}
	pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: config.PageWidth, H: config.PageHeight}})
	pdf.SetTopMargin(config.Top)
	pdf.SetLeftMargin(config.Side)
	pdf.SetLineWidth(config.Side)
	pdf.Config = config
	//设置默认字体
	var err error
	err = pdf.AddTTFFont(FONT_NAME, FONT_PATH)
	HandleError("load font error", err)
	err = pdf.AddTTFFont(FONT_NAME_BOLD, FONT_PATH_BOLD)
	HandleError("load boldfont error", err)
	err = pdf.SetFont(FONT_NAME, "", DEFAULT_FONT_SIZE)
	HandleError("set font error", err)
	pdf.CurrFontName = FONT_NAME
	//设置线
	pdf.SetLineWidth(1)
	return pdf
}

// pdf 加一个页面
func (pdf *XJDPdf) ZAddPage() *XJDPdf {
	pdf.AddPage()
	return pdf
}

// 4象限，向右偏移
func (pdf *XJDPdf) AddX(size float64) *XJDPdf {
	x := pdf.GetX()
	pdf.SetX(x + size)
	return pdf
}

// 4象限，向下偏移
func (pdf *XJDPdf) AddY(size float64) *XJDPdf {
	y := pdf.GetY()
	pdf.SetY(y + size)
	return pdf
}

// 普通的文字写入
func (pdf *XJDPdf) Write(text string) *XJDPdf {
	pdf.Cell(nil, text)
	return pdf
}

func (pdf *XJDPdf) WriteWithLineFont(text string, font int) *XJDPdf {
	x, y := pdf.GetX(), pdf.GetY()

	pdf.Cell(nil, text)
	wide, _ := pdf.MeasureTextWidth(text)
	lineHight := float64(float64(font) + float64(font)/5)
	pdf.Line(x, y+lineHight, x+wide, y+lineHight)

	pdf.SetX(x + wide)
	return pdf
}

// 在文字下面 加入一个下划线
func (pdf *XJDPdf) WriteWithLine(text string) *XJDPdf {
	x, y := pdf.GetX(), pdf.GetY()
	// log.Println("x,y:", x, y)
	pdf.Cell(nil, text)
	wide, _ := pdf.MeasureTextWidth(text)
	pdf.Line(x, y+defaultLineHight, x+wide, y+defaultLineHight)
	// log.Println("wide:", wide)
	// log.Println(x, y+defaultLineHight, x+wide, y+defaultLineHight)
	//划线后x轴值不是线的末尾，设置回去
	pdf.SetX(x + wide)
	return pdf
}

//固定长度中间写划线
func (pdf *XJDPdf) WriteInCenterWithLineFixedlength(text string, length int, font int) *XJDPdf {
	lentext := LenOfSee(text)
	//去除逻辑错误
	if length < 1 {
		length = 1
	}
	if font < 0 {
		font = 10
	}
	if lentext > length {
		text = text[:length]
	} else {
		//前后缀空格
		prefix := GetBlank((length - lentext) / 2)
		suffix := GetBlank((length - lentext) - (length-lentext)/2)
		text = prefix + text + suffix
	}
	return pdf.WriteWithLineFont(text, font)
}

// 自动换行，默认长度换行
func (pdf *XJDPdf) WriteMassage(text string) *XJDPdf {
	length := Strlen(text)
	wideFlag := pdf.Config.PageWidth - pdf.Config.Side*2
	var i int
	var offset int
	var flag bool
	for i = 0; i <= length-defaultLineFontNum; i += defaultLineFontNum {
		offset = 0
		flag = true
		for flag {
			if i+defaultLineFontNum+offset > length {
				flag = false
			}
			strline := Substr(text, i, defaultLineFontNum+offset)
			wide, _ := pdf.MeasureTextWidth(strline)
			if wide < wideFlag {
				offset++
			} else {
				flag = false
			}
		}
		pdf.Write(Substr(text, i, defaultLineFontNum+offset))
		pdf.DoDefaultBR()
		i += offset
	}
	if i < length {
		pdf.Write(Substr(text, i, length-i))
	}
	return pdf
}

// 设置不同的字体
func (pdf *XJDPdf) SetFontFamily(family string) *XJDPdf {
	pdf.CurrFontName = family
	return pdf
}

// 普通字体 设置文字(居中)
func (pdf *XJDPdf) WriteInCenter(text string, fontSize int) *XJDPdf {
	//1.设置字体
	pdf.SetFontSize(fontSize)
	wide, _ := pdf.MeasureTextWidth(text)
	//2.居中定位,写
	// offset := (PAGEWIDE - pdf.GetX() - wide) / 2
	offset := (pdf.Config.PageWidth - wide) / 2
	pdf.SetX(offset)
	pdf.Cell(nil, text)
	pdf.DoBR(GetBR(fontSize))
	if fontSize != DEFAULT_FONT_SIZE {
		//3.字体设置回默认的
		pdf.SetDefaultFontSize()
	}
	return pdf
}

//右对齐
func (pdf *XJDPdf) WriteInRight(text string) *XJDPdf {

	wide, _ := pdf.MeasureTextWidth(text)
	offset := pdf.Config.PageWidth - pdf.Config.Side - wide
	pdf.SetX(offset)
	pdf.Cell(nil, text)
	pdf.DoBR(GetBR(DEFAULT_FONT_SIZE))
	return pdf
}

func (pdf *XJDPdf) WriteAnyPlace(percent float64, text string, fontSize int) *XJDPdf {
	if fontSize == DEFAULT_FONT_SIZE {
		pdf.SetX(pdf.Config.PageWidth * percent)
		pdf.Cell(nil, text)
	} else {
		pdf.SetFontSize(fontSize)

		pdf.SetX(pdf.Config.PageWidth * percent)
		pdf.Cell(nil, text)

		pdf.SetDefaultFontSize()
	}
	return pdf
}

//自动分段(自定宽度，前缀长度，行间距，所有内容)
func (pdf *XJDPdf) WritePassageAnyWidth(wideFlag int, preLength float64, lineSpace int, text string) *XJDPdf {
	r := []rune(text)
	length := len(r)
	if length < wideFlag {
		return pdf.AddX(preLength).Write(text).DoBR(GetBR(lineSpace))
	}
	var i int

	for length > wideFlag {
		length -= wideFlag
		pdf.AddX(preLength).Write(string(r[i : i+wideFlag])).DoBR(GetBR(lineSpace))
		i += wideFlag
	}
	return pdf.AddX(preLength).Write(string(r[i:])).DoBR(GetBR(lineSpace))
}

// 设置文字大小
func (pdf *XJDPdf) SetFontSize(fontSize int) *XJDPdf {
	err := pdf.SetFont(pdf.CurrFontName, "", fontSize)
	HandleError("set font error", err)
	return pdf
}

// 多大字体回车
func (pdf *XJDPdf) DoBR(height float64) *XJDPdf {
	if pdf.GetY() > pdf.Config.PageHeight-pdf.Config.Bottom-DEFAULT_FONT_SIZE {
		pdf.AddPage()
	}
	pdf.Br(height)
	return pdf
}

// 回车
func GetBR(fontSize int) float64 {
	return float64(float64(fontSize) + float64(fontSize)/2)
}

// 默认回车
func (pdf *XJDPdf) DoDefaultBR() *XJDPdf {
	if pdf.GetY() > pdf.Config.PageHeight-pdf.Config.Bottom-DEFAULT_FONT_SIZE {
		pdf.AddPage()
	}
	pdf.Br(defaultBr)
	return pdf
}

// 字体大小
func (pdf *XJDPdf) SetDefaultFontSize() *XJDPdf {
	err := pdf.SetFont(FONT_NAME, "", DEFAULT_FONT_SIZE)
	HandleError("set font error", err)
	return pdf
}

// 日志打印输出
func HandleError(prefix string, err error) bool {
	if err != nil {
		log.Fatalln(err.Error())
		panic(errors.New(prefix))
		return true
	} else {
		return false
	}
}

// 输出这个pdf,保存为文件
func (pdf *XJDPdf) Out(filename string) *XJDPdf {
	pdf.WritePdf(filename)
	return pdf
}
