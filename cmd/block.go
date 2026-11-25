package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Block struct {
	WorldX      float32
	WorldY      float32
	ScreenX     float32
	ScreenY     float32
	Width       float32
	Height      float32
	Texture     *ebiten.Image
	IsHoveredOn bool
}

func NewBlock(worldX, worldY, screenX, screenY, w, h float32, textureName string) (*Block, error) {
	var texture *ebiten.Image
	var err error

	texture, _, err = LoadTexture(textureName, &TextureConfig{path: "textures", ext: "png"})

	if err != nil {
		log.Fatal(err)
	}

	return &Block{WorldX: worldX, WorldY: worldY, ScreenX: screenX, ScreenY: screenY, Width: w, Height: h, Texture: texture, IsHoveredOn: false}, nil
}

func (b *Block) UpdateScreenPos(camX, camY float32, Width, Height int) {
	b.ScreenX = float32(b.WorldX) - camX + float32(Width/2)
	b.ScreenY = float32(b.WorldY) - camY + float32(Height/2)
}
