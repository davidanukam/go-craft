package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Block struct {
	WorldX  float32
	WorldY  float32
	ScreenX float32
	ScreenY float32
	W       float32
	H       float32
	Texture *ebiten.Image
}

func NewBlock(worldX, worldY, screenX, screenY, w, h float32, textureName string) (*Block, error) {
	var image *ebiten.Image
	var err error

	image, _, err = LoadTexture(textureName, &TextureConfig{path: "textures", ext: "png"})

	if err != nil {
		log.Fatal(err)
	}

	return &Block{worldX, worldY, screenX, screenY, w, h, image}, nil
}

func (b *Block) Update(camX, camY float32, Width, Height int) {
	b.ScreenX = float32(b.WorldX) - camX + float32(Width/2)
	b.ScreenY = float32(b.WorldY) - camY + float32(Height/2)
}
