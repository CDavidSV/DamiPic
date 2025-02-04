package internal

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	f *truetype.Font

	dpi               float64        = 72.0
	foreground        *image.Uniform = image.Black
	defaultBackground *image.Uniform = &image.Uniform{color.RGBA{220, 220, 220, 255}}
	spacing           float64        = 1.5

	ErrInvalidColorHex = fmt.Errorf("invalid color hex")

	contextPool = sync.Pool{
		New: func() any {
			ctx := freetype.NewContext()

			ctx.SetDPI(dpi)
			ctx.SetFont(f)

			return ctx
		},
	}
)

func init() {
	exePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Might give an error when running the tests. Use ../fonts when running the tests
	fontFilePath := filepath.Join(exePath, "fonts", "Poppins-SemiBold.ttf")

	fontBytes, err := os.ReadFile(fontFilePath)
	f, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal(err)
	}
}

func ParseHexColorString(hex string) (color.RGBA, error) {
	var rgbaColor color.RGBA

	if len(hex) != 6 {
		return rgbaColor, fmt.Errorf("invalid color string")
	}

	red, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return rgbaColor, err
	}

	green, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return rgbaColor, err
	}

	blue, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return rgbaColor, err
	}

	rgbaColor = color.RGBA{uint8(red), uint8(green), uint8(blue), 255}
	return rgbaColor, nil
}

func addLabel(img *image.RGBA, text string, fr *image.Uniform) error {
	fontSize := float64(img.Rect.Max.Y) * 0.09
	margin := float64(img.Rect.Max.X) * 0.03

	options := &truetype.Options{
		Size: fontSize,
		DPI:  dpi,
	}

	face := truetype.NewFace(f, options)

	textLines := strings.Split(strings.ReplaceAll(text, "+", " "), "\\n")

	// Get the longest line of text
	longestString := ""
	longestStringWidth := 0
	for _, line := range textLines {
		lineWidth := int(font.MeasureString(face, line).Round())

		if lineWidth > longestStringWidth {
			longestStringWidth = lineWidth
			longestString = line
		}
	}

	px := (img.Rect.Max.X / 2) - longestStringWidth/2

	// Adjusts the font size to fit the text in the image
	for px+longestStringWidth+int(margin) > img.Rect.Max.X {
		fontSize *= 0.8

		options.Size = fontSize
		face = truetype.NewFace(f, options)

		longestStringWidth = font.MeasureString(face, longestString).Round()
		px = (img.Rect.Max.X / 2) - longestStringWidth/2
	}

	// We calculate the position of the text in the image after adjusting the font size
	py := (img.Rect.Max.Y / 2) - int(fontSize*float64(len(textLines))+spacing)/2

	ctx := contextPool.Get().(*freetype.Context)
	ctx.SetFontSize(fontSize)
	ctx.SetSrc(fr)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)

	pt := freetype.Pt(px, py+int(ctx.PointToFixed(fontSize)>>6))
	textWidth := 0
	for _, line := range textLines {
		textWidth = font.MeasureString(face, line).Round()
		pt.X = fixed.I((img.Rect.Max.X / 2) - textWidth/2)
		if _, err := ctx.DrawString(line, pt); err != nil {
			return err
		}

		pt.Y += ctx.PointToFixed(fontSize * spacing)
	}

	contextPool.Put(ctx)

	return nil
}

func GeneratePlaceholderImg(width, heigth int, text string, bgColor string, textColor string) (*image.RGBA, error) {
	baseImage := image.NewRGBA(image.Rect(0, 0, width, heigth))

	bg := defaultBackground
	if bgColor != "" {
		rgbaColor, err := ParseHexColorString(bgColor)
		if err != nil {
			return nil, ErrInvalidColorHex
		}

		bg = &image.Uniform{rgbaColor}
	}

	textForeground := foreground
	if textColor != "" {
		rgbaColor, err := ParseHexColorString(textColor)
		if err != nil {
			return nil, ErrInvalidColorHex
		}

		textForeground = &image.Uniform{rgbaColor}
	}

	// Sets a color to the image
	draw.Draw(baseImage, baseImage.Bounds(), bg, image.Point{0, 0}, draw.Src)

	if err := addLabel(baseImage, text, textForeground); err != nil {
		return nil, err
	}

	return baseImage, nil
}
