package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"strings"

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
)

func init() {
	var fontFilePath string = "./fonts/Poppins-SemiBold.ttf"

	fontBytes, err := os.ReadFile(fontFilePath)
	if err != nil {
		log.Fatal(err)
	}

	f, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal(err)
	}
}

// func parseHexColorString(color string) color.RGBA {

// }

func addLabel(img *image.RGBA, text string) error {
	fontSize := float64(img.Rect.Max.Y) * 0.04
	margin := float64(img.Rect.Max.X) * 0.05

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

	ctx := freetype.NewContext()
	ctx.SetDPI(dpi)
	ctx.SetFont(f)
	ctx.SetFontSize(fontSize)
	ctx.SetSrc(foreground)
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

	return nil
}

func generatePlaceholderImg(width, heigth int, text string) (*image.RGBA, error) {
	baseImage := image.NewRGBA(image.Rect(0, 0, width, heigth))

	// Sets a color to the image
	draw.Draw(baseImage, baseImage.Bounds(), defaultBackground, image.Point{0, 0}, draw.Src)

	if err := addLabel(baseImage, text); err != nil {
		return nil, err
	}

	return baseImage, nil
}
