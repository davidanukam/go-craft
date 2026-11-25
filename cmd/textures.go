package main

import (
	"image"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var BasePath string = "textures"
var BaseExt string = "png"

type TextureConfig struct {
	Path string
	Ext  string
}

type Texture struct {
	Image *ebiten.Image
	File  string
}

var Textures map[string]*ebiten.Image = make(map[string]*ebiten.Image)

func LoadTexture(textureName string, cfg *TextureConfig) (*Texture, image.Image, error) {
	var sb strings.Builder

	// Build path
	if cfg == nil {
		cfg = &TextureConfig{Path: BasePath, Ext: BaseExt}
	} else {
		if cfg.Path == "" {
			log.Fatal("No path given to load texture from (e.g assets)")
		}
		if cfg.Ext == "" {
			log.Fatal("No extension given to load texture from (e.g .png)")
		}
	}

	sb.WriteString(cfg.Path)
	sb.WriteString("/")
	sb.WriteString(textureName)
	sb.WriteString(".")
	sb.WriteString(cfg.Ext)

	img, _, err := ebitenutil.NewImageFromFile(sb.String())
	if err != nil {
		log.Fatal(err)
	}

	return &Texture{Image: img, File: sb.String()}, nil, err
}

func LoadAllTextures() {
	textureFolder, folderErr := os.ReadDir("textures")

	if folderErr != nil {
		log.Fatal(textureFolder, folderErr)
	}

	for _, file := range textureFolder {
		img, _, imageErr := ebitenutil.NewImageFromFile("textures/" + file.Name())

		if imageErr != nil {
			log.Fatal(imageErr)
		}

		Textures[file.Name()] = img
	}
}
