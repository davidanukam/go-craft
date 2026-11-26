package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var WIDTH, HEIGHT int = 736, 736 // 352, 352
var BLOCK_SIZE int = 32

type Game struct {
	CurrChunkCoord *ChunkCoord
	Chunks         map[ChunkCoord]*Chunk
	Plr            *Player
	Camera         *Camera
}

var game *Game

// Crate a new Game object
func NewGame() {
	game = &Game{CurrChunkCoord: &ChunkCoord{0, 0}, Chunks: make(map[ChunkCoord]*Chunk, 0), Plr: nil, Camera: nil}

	var chunkCenterX = WIDTH / 2
	var chunkCenterY = HEIGHT / 2
	player, err := CreatePlayer(float32(chunkCenterX)-(float32(BLOCK_SIZE)/2), float32(chunkCenterY)-(float32(BLOCK_SIZE)/2), BLOCK_SIZE, BLOCK_SIZE, "player")

	if err != nil {
		log.Fatal(err)
	}

	LoadAllTextures()

	game.Plr = player
	game.Camera = &Camera{OffsetX: player.WorldX, OffsetY: player.WorldY, Speed: 4}

	GenerateChunk(0, 0, "center", "grass")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "top-left", "dirt")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "top", "stone")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "top-right", "dirt")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "left", "grass")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "right", "grass")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "bottom-left", "dirt")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "bottom", "stone")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "bottom-right", "dirt")
}

func main() {
	NewGame()

	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Go-craft")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
