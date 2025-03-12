// Package letteravatar generates letter-avatars.
package letteravatar

import (
	"fmt"
	"hash/fnv"
	"image"
	"image/color"
	"image/draw"
	"math/rand/v2"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Options are letter-avatar parameters.
type Options struct {
	Font        *opentype.Font
	Palette     []color.Color
	LetterColor color.Color
	FontSize    int

	// PaletteKey is used to pick the background color from the Palette.
	// Using the same PaletteKey leads to the same background color being picked.
	// If PaletteKey is empty (default) the background color is picked randomly.
	PaletteKey string
}

var defaultLetterColor = color.RGBA{R: 0xf0, G: 0xf0, B: 0xf0, A: 0xf0}

// Draw generates a new letter-avatar image of the given size using the given letter
// with the given options. Default parameters are used if a nil *Options is passed.
func Draw(size int, letters []rune, options *Options) (image.Image, error) {
	if options == nil {
		options = &Options{}
	}
	if options.Font == nil {
		options.Font = defaultFont
	}
	if options.Palette == nil {
		options.Palette = defaultPalette
	}
	if options.LetterColor == nil {
		options.LetterColor = defaultLetterColor
	}

	var bgColor color.Color = color.RGBA{A: 0xff}
	if len(options.Palette) > 0 {
		if len(options.PaletteKey) > 0 {
			bgColor = options.Palette[randomIndex(len(options.Palette), options.PaletteKey)]
		} else {
			bgColor = options.Palette[rand.IntN(len(options.Palette))]
		}
	}

	return drawAvatar(bgColor, options.LetterColor, options.Font, size, float64(options.FontSize), letters)
}

func drawAvatar(bgColor, fgColor color.Color, f *opentype.Font, size int, fontSize float64, letters []rune) (image.Image, error) {
	dst := newRGBA(size, size, bgColor)

	src, err := drawString(bgColor, fgColor, f, size, fontSize, letters)
	if err != nil {
		return nil, err
	}

	r := src.Bounds().Add(dst.Bounds().Size().Div(2)).Sub(src.Bounds().Size().Div(2))
	draw.Draw(dst, r, src, src.Bounds().Min, draw.Src)

	return dst, nil
}

func drawString(bgColor, fgColor color.Color, f *opentype.Font, size int, fontSize float64, letters []rune) (image.Image, error) {
	// 自动计算字体大小
	if fontSize <= 0 {
		fontSize = calculateFontSize(len(letters), size)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create font face: %v", err)
	}
	defer face.Close()

	// 计算文本尺寸以居中显示
	textWidth, ascent, descent, _ := calculateTextDimensions(face, string(letters))

	x := (size - textWidth) / 2
	y := size/2 + (ascent-descent)/2

	dst := newRGBA(x, y, bgColor)
	drawText(dst, face, string(letters), x, y, fgColor)

	return dst, nil
}

// calculateFontSize 根据图像大小和文本长度确定最佳字体大小
func calculateFontSize(textLength int, imageSize int) float64 {
	// 以图像大小的2/3作为初始大小估计
	initialSize := float64(imageSize) * 2.0 / 3.0

	// 根据文本长度调整
	if textLength > 1 {
		initialSize = initialSize * 3.0 / float64(textLength+2)
	}

	// 确保字体大小至少为12px且不超过图像大小
	if initialSize < 12 {
		initialSize = 12
	} else if initialSize > float64(imageSize) {
		initialSize = float64(imageSize)
	}

	return initialSize
}

// calculateTextDimensions 计算文本的宽度和高度，返回宽度、上升部分高度、下降部分高度和行间距
func calculateTextDimensions(face font.Face, text string) (width int, ascent int, descent int, lineGap int) {
	metrics := face.Metrics()
	ascent = metrics.Ascent.Ceil()
	descent = metrics.Descent.Ceil()
	lineGap = metrics.Height.Ceil() - ascent - descent

	// 计算文本宽度
	var w fixed.Int26_6
	for _, r := range text {
		adv, ok := face.GlyphAdvance(r)
		if !ok {
			adv, _ = face.GlyphAdvance(' ') // 如果字符不支持，使用空格的宽度
		}
		w += adv
	}

	return w.Ceil(), ascent, descent, lineGap
}

// drawText 将文本绘制到图像上
func drawText(img *image.RGBA, face font.Face, text string, x, y int, textColor color.Color) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textColor),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)},
	}
	d.DrawString(text)
}

func newRGBA(w, h int, c color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), image.NewUniform(c), image.Point{}, draw.Src)
	return img
}

func randomIndex(n int, key string) int {
	h := fnv.New64a()
	if _, err := h.Write([]byte(key)); err != nil {
		return 0
	}
	rd := rand.New(rand.NewPCG(h.Sum64(), (h.Sum64()>>1)|1))
	return rd.IntN(n)
}
