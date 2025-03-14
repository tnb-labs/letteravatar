// Package initials generates initials-avatars.
package initials

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

func drawAvatar(bgColor, fgColor color.Color, ft *opentype.Font, size int, fontSize float64, letters []rune) (image.Image, error) {
	// Auto-calculate font size if not provided
	if fontSize <= 0 {
		fontSize = calculateFontSize(len(letters), size)
	}

	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create font face: %v", err)
	}
	defer face.Close()

	textWidth, ascent, descent, _ := calculateTextDimensions(face, string(letters))
	x := (size - textWidth) / 2
	y := size/2 + (ascent-descent)/2

	dst := newRGBA(size, size, bgColor)
	drawText(dst, face, string(letters), x, y, fgColor)

	return dst, nil
}

// calculateFontSize calculates the best font size based on image size and text length
func calculateFontSize(textLength int, imageSize int) float64 {
	// use 2/3 of the image size as the initial size estimate
	initialSize := float64(imageSize) * 2.0 / 3.0

	// adjust based on text length
	if textLength > 1 {
		initialSize = initialSize * 3.0 / float64(textLength+2)
	}

	// make sure the font size is at least 12px and does not exceed the image size
	if initialSize < 12 {
		initialSize = 12
	} else if initialSize > float64(imageSize) {
		initialSize = float64(imageSize)
	}

	return initialSize
}

// calculateTextDimensions calculates the width and height of the text
func calculateTextDimensions(face font.Face, text string) (width int, ascent int, descent int, lineGap int) {
	metrics := face.Metrics()
	ascent = metrics.Ascent.Ceil()
	descent = metrics.Descent.Ceil()
	lineGap = metrics.Height.Ceil() - ascent - descent

	// width
	var w fixed.Int26_6
	for _, r := range text {
		adv, ok := face.GlyphAdvance(r)
		if !ok {
			adv, _ = face.GlyphAdvance(' ') // if glyph not found, fall back to space
		}
		w += adv
	}

	return w.Ceil(), ascent, descent, lineGap
}

// drawText draws the text on the image
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
