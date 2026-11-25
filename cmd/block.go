package main

import (
	"log"
)

type Block struct {
	WorldX      float32
	WorldY      float32
	ScreenX     float32
	ScreenY     float32
	Width       float32
	Height      float32
	Texture     *Texture
	IsHoveredOn bool
}

func NewBlock(worldX, worldY, screenX, screenY, width, height float32, textureName string) (*Block, error) {
	// FIX: Block textures are chosen based on the worldY for now
	if worldY < float32(736)/2 {
		textureName = "air"
	}

	texture, _, err := LoadTexture(textureName, &TextureConfig{Path: "textures", Ext: "png"})

	if err != nil {
		log.Fatal(err)
	}

	// DEBUGGING
	// var texture *ebiten.Image
	// texture = Textures[textureName]

	return &Block{WorldX: worldX, WorldY: worldY, ScreenX: screenX, ScreenY: screenY, Width: width, Height: height, Texture: texture, IsHoveredOn: false}, nil
}

func (b *Block) UpdateScreenPos(camX, camY float32, Width, Height int) {
	b.ScreenX = float32(b.WorldX) - camX + float32(Width/2)
	b.ScreenY = float32(b.WorldY) - camY + float32(Height/2)
}
