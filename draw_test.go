package initials

import (
	"image"
	"image/color"
	"reflect"
	"testing"
)

func TestNewRGBA(t *testing.T) {
	w := 10
	h := 20
	c := color.RGBA{G: 50, B: 100, A: 200}
	img := newRGBA(w, h, c)
	if !img.Bounds().Eq(image.Rect(0, 0, w, h)) {
		t.Fatalf("bad bounds")
	}
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if !reflect.DeepEqual(img.RGBAAt(x, y), c) {
				t.Fatalf("bad color at %d, %d: %#v", x, y, img.RGBAAt(x, y))
			}
		}
	}
}

func TestDraw(t *testing.T) {
	testcases := []struct {
		size    int
		letter  rune
		options *Options
		want    image.Image
	}{
		{
			size:   50,
			letter: 'Z',
			options: &Options{
				Palette:     []color.Color{color.Black},
				LetterColor: color.White,
			},
			want: &image.RGBA{
				Rect:   image.Rect(0, 0, 50, 50),
				Stride: 50 * 4,
			},
		},
		{
			size:   40,
			letter: 'Я',
			options: &Options{
				Palette:     []color.Color{color.RGBA{R: 0x33, G: 0x66, B: 0x99, A: 0xff}},
				LetterColor: color.RGBA{R: 0xaa, G: 0xaa, B: 0xaa, A: 0xaa},
			},
			want: &image.RGBA{
				Rect:   image.Rect(0, 0, 40, 40),
				Stride: 40 * 4,
			},
		},
	}

	for i, testcase := range testcases {
		img, err := Draw(testcase.size, []rune{testcase.letter}, testcase.options)
		if err != nil {
			t.Fatalf("failed to create avatar #%d: %s", i, err)
		}
		if !img.Bounds().Eq(testcase.want.Bounds()) {
			t.Fatalf("avatar #%d has bad bounds: got %v, want %v", i, img.Bounds(), testcase.want.Bounds())
		}
	}
}

func TestPaletteKey(t *testing.T) {
	users := []string{
		"Username 1",
		"Username 2",
		"Username 3",
		"Username 4",
		"Username 5",
	}
	avatars := make(map[string]image.Image)
	for _, u := range users {
		img, err := Draw(30, []rune(u), &Options{PaletteKey: u})
		if err != nil {
			t.Fatalf("failed to create avatar for %s: %s", u, err)
		}
		avatars[u] = img
	}
	for _, u := range users {
		img, err := Draw(30, []rune(u), &Options{PaletteKey: u})
		if err != nil {
			t.Fatalf("failed to create avatar for %s: %s", u, err)
		}
		if !reflect.DeepEqual(avatars[u], img) {
			t.Fatalf("avatar mismatch for %s: %#v, %#v", u, avatars[u], img)
		}
	}
}
