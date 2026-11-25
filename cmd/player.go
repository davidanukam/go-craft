package main

import (
	"log"
)

type Player struct {
	WorldX  float32
	WorldY  float32
	ScreenX float32
	ScreenY float32
	Width   int
	Height  int
	Texture *Texture
	VelX    float32
	VelY    float32
	Gravity float32
}

func CreatePlayer(worldX, worldY float32, width, height int, textureName string) (*Player, error) {
	texture, _, err := LoadTexture(textureName, &TextureConfig{Path: "assets", Ext: "png"})

	if err != nil {
		log.Fatal(err)
	}

	return &Player{WorldX: worldX, WorldY: worldY, ScreenX: worldX, ScreenY: worldY, Width: width, Height: height, Texture: texture, VelX: 0, VelY: 0, Gravity: 1}, nil
}

func (p *Player) UpdateScreenPos(camX, camY, Width, Height float32) error {
	p.ScreenX = p.WorldX - camX + Width/2
	p.ScreenY = p.WorldY - camY + Height/2

	return nil
}
