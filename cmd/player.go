package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	WorldX  float32
	WorldY  float32
	ScreenX float32
	ScreenY float32
	Width   int
	Height  int
	Texture *ebiten.Image
	VelX    float32
	VelY    float32
	Gravity float32
}

func CreatePlayer(worldX, worldY float32, w, h int, textureName string) (*Player, error) {
	var texture *ebiten.Image
	var err error

	texture, _, err = LoadTexture(textureName, &TextureConfig{path: "assets", ext: "png"})

	if err != nil {
		log.Fatal(err)
	}

	return &Player{WorldX: worldX, WorldY: worldY, ScreenX: worldX, ScreenY: worldY, Width: w, Height: h, Texture: texture, VelX: 0, VelY: 0, Gravity: 1}, nil
}

func (p *Player) Update(camX, camY, Width, Height float32) error {
	p.ScreenX = p.WorldX - camX + Width/2
	p.ScreenY = p.WorldY - camY + Height/2

	return nil
}
