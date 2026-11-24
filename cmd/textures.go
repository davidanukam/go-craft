package main

import (
	"image"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var BasePath string = "textures"
var BaseExt string = "png"

type TextureConfig struct {
	path string
	ext  string
}

type Textures struct {
	dirt string
}

func LoadTexture(textureName string, cfg *TextureConfig) (*ebiten.Image, image.Image, error) {
	var img *ebiten.Image
	var err error

	var sb strings.Builder

	// Build path
	if cfg == nil {
		cfg = &TextureConfig{path: BasePath, ext: BaseExt}
	} else {
		if cfg.path == "" {
			log.Fatal("No path given to load texture from (e.g assets)")
		}
		if cfg.ext == "" {
			log.Fatal("No extension given to load texture from (e.g .png)")
		}
	}

	sb.WriteString(cfg.path)
	sb.WriteString("/")
	sb.WriteString(textureName)
	sb.WriteString(".")
	sb.WriteString(cfg.ext)

	img, _, err = ebitenutil.NewImageFromFile(sb.String())
	if err != nil {
		log.Fatal(err)
	}

	return img, nil, err
}
