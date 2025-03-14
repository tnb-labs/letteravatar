# initials

[![GoDoc](https://godoc.org/github.com/weavatar/initials?status.svg)](https://godoc.org/github.com/weavatar/initials)
[![Go](https://img.shields.io/github/go-mod/go-version/weavatar/initials)](https://go.dev/)
[![Release](https://img.shields.io/github/release/weavatar/initials.svg)](https://github.com/weavatar/initials/releases)
[![Test](https://github.com/weavatar/initials/actions/workflows/test.yml/badge.svg)](https://github.com/weavatar/initials/actions)
[![Report Card](https://goreportcard.com/badge/github.com/weavatar/initials)](https://goreportcard.com/report/github.com/weavatar/initials)

Initials avatar generation for Go.

## Usage

Generate a 100x100px 'A'-initial avatar:

```go
img, err := initials.Draw(100, 'A', nil)
```

The third parameter `options *Options` can be used for customization:

```go
type Options struct {
	Font        *truetype.Font
	Palette     []color.Color
	LetterColor color.Color
	PaletteKey  string
	FontSize    int
}
```

Using a custom palette:

```go
img, err := initials.Draw(100, []rune{'A'}, &initials.Options{
	Palette: []color.Color{
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
	},
})
```

## Documentation

[https://godoc.org/github.com/weavatar/initials](https://godoc.org/github.com/weavatar/initials)

## Examples

![](example/Alice.png)
![](example/Bob.png)
![](example/Carol.png)
![](example/Dave.png)
![](example/Eve.png)
![](example/Frank.png)
![](example/Gloria.png)
![](example/Henry.png)
![](example/Isabella.png)
![](example/Жозефина.png)
![](example/Ярослав.png)

```go
package main

import (
	"image/png"
	"log"
	"os"
	"unicode/utf8"

	"github.com/weavatar/initials"
)

var names = []string{
	"Alice",
	"Bob",
	"Carol",
	"Dave",
	"Eve",
	"Frank",
	"Gloria",
	"Henry",
	"Isabella",
	"James",
	"Жозефина",
	"Ярослав",
}

func main() {
	for _, name := range names {
		firstLetter, _ := utf8.DecodeRuneInString(name)

		img, err := initials.Draw(75, []rune{firstLetter}, nil)
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(name + ".png")
		if err != nil {
			log.Fatal(err)
		}

		err = png.Encode(file, img)
		if err != nil {
			log.Fatal(err)
		}
	}
}

```

## License

The package "initials" is distributed under the terms of the MIT license.

The Roboto-Medium font is distributed under the terms of the Apache License v2.0.

See [LICENSE](LICENSE).
