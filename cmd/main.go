package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

	// player, err := CreatePlayer(float32(WIDTH)/2, float32(HEIGHT)/2, 32, 32, "player")
	var chunkCenterX = WIDTH / 2
	var chunkCenterY = HEIGHT / 2
	player, err := CreatePlayer(float32(chunkCenterX)-(float32(BLOCK_SIZE)/2), float32(chunkCenterY)-(float32(BLOCK_SIZE)/2), 32, 32, "player")

	if err != nil {
		log.Fatal(err)
	}

	LoadAllTextures()

	game.Plr = player
	game.Camera = &Camera{OffsetX: player.WorldX, OffsetY: player.WorldY, Speed: 5}

	GenerateChunk(0, 0, "center", "grass")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "top-left", "dirt")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "top", "stone")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "top-right", "dirt")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "left", "grass")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "right", "grass")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "bottom-left", "dirt")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "bottom", "stone")
	GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, "bottom-right", "dirt")

	CheckSurroundingChunks()
}

func main() {
	NewGame()

	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Go-craft")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// --- Helper Functions --- //

// Update game.CurrChunkCoord (pointer) to store the current chunk the player is in
func UpdateCurrChunkCoord() {
	for _, chunk := range game.Chunks {
		if game.Plr.ScreenY >= chunk.Top && game.Plr.ScreenY <= chunk.Bottom && game.Plr.ScreenX >= chunk.Left && game.Plr.ScreenX <= chunk.Right {
			// Debugging
			// fmt.Printf("UpdatedCurrChunkCoord - Chunk_WorldX: %d, Chunk_WorldY: %d, Chunk_Top: %f, Chunk_Bottom: %f, Chunk_Left: %f, Chunk_Right: %f\n", chunk.Coord.WorldX, chunk.Coord.WorldY, chunk.Top, chunk.Bottom, chunk.Left, chunk.Right)
			game.CurrChunkCoord.WorldX = chunk.Coord.WorldX
			game.CurrChunkCoord.WorldY = chunk.Coord.WorldY
			break
		}
	}
}

var TopLeftChunk, TopChunk, TopRightChunk, LeftChunk, CenterChunk, RightChunk, BottomLeftChunk, BottomChunk, BottomRightChunk *Chunk

// Check the 8 surrounding chunks
func CheckSurroundingChunks() {
	var _TopLeftChunk, _TopChunk, _TopRightChunk, _LeftChunk, _CenterChunk, _RightChunk, _BottomLeftChunk, _BottomChunk, _BottomRightChunk *Chunk

	curr := game.Chunks[*game.CurrChunkCoord]

	for _, chunk := range game.Chunks {
		// Top-Left
		if chunk.Bottom == curr.Top && chunk.Right == curr.Left {
			_TopLeftChunk = chunk
		}

		// Top
		if chunk.Bottom == curr.Top && chunk.Left == curr.Left {
			_TopChunk = chunk
		}

		// Top-Right
		if chunk.Bottom == curr.Top && chunk.Left == curr.Right {
			_TopRightChunk = chunk
		}

		// Left
		if chunk.Top == curr.Top && chunk.Right == curr.Left {
			_LeftChunk = chunk
		}

		// Center
		if chunk.Top == curr.Top && chunk.Left == curr.Left {
			_CenterChunk = chunk
		}

		// Right
		if chunk.Top == curr.Top && chunk.Left == curr.Right {
			_RightChunk = chunk
		}

		// Bottom-Left
		if chunk.Top == curr.Bottom && chunk.Right == curr.Left {
			_BottomLeftChunk = chunk
		}

		// Bottom
		if chunk.Top == curr.Bottom && chunk.Left == curr.Left {
			_BottomChunk = chunk
		}

		// Bottom-Right
		if chunk.Top == curr.Bottom && chunk.Left == curr.Right {
			_BottomRightChunk = chunk
		}
	}

	TopLeftChunk = _TopLeftChunk
	TopChunk = _TopChunk
	TopRightChunk = _TopRightChunk
	LeftChunk = _LeftChunk
	CenterChunk = _CenterChunk
	RightChunk = _RightChunk
	BottomLeftChunk = _BottomLeftChunk
	BottomChunk = _BottomChunk
	BottomRightChunk = _BottomRightChunk
}

func GenerateChunk(ChunkCoordWorldX, ChunkCoordWorldY int, side string, ty string) *Chunk {
	var newX int
	var newY int

	// index = y * Width + x
	switch side {
	case "center":
		newX = 0
		newY = 0
	case "top-left":
		newX = ChunkCoordWorldX - WIDTH
		newY = ChunkCoordWorldY - HEIGHT
	case "top":
		newX = ChunkCoordWorldX
		newY = ChunkCoordWorldY - HEIGHT
	case "top-right":
		newX = ChunkCoordWorldX + WIDTH
		newY = ChunkCoordWorldY - HEIGHT
	case "left":
		newX = ChunkCoordWorldX - WIDTH
		newY = ChunkCoordWorldY
	case "right":
		newX = ChunkCoordWorldX + WIDTH
		newY = ChunkCoordWorldY
	case "bottom-left":
		newX = ChunkCoordWorldX - WIDTH
		newY = ChunkCoordWorldY + HEIGHT
	case "bottom":
		newX = ChunkCoordWorldX
		newY = ChunkCoordWorldY + HEIGHT
	case "bottom-right":
		newX = ChunkCoordWorldX + WIDTH
		newY = ChunkCoordWorldY + HEIGHT
	}

	// fmt.Printf("newX: %f, newY: %f\n", newX, newY)
	coord := ChunkCoord{WorldX: newX, WorldY: newY}
	chunk, err := NewChunk(coord, game.Camera.OffsetX, game.Camera.OffsetY, WIDTH, HEIGHT, ty)

	if err != nil {
		log.Fatal(err)
	}

	game.Chunks[coord] = chunk

	fmt.Printf("New Chunk generated [%s] - CoordWorldX: %d, CoordWorldY: %d, Top: %f, Bottom: %f, Left: %f, Right: %f\n", side, chunk.Coord.WorldX, chunk.Coord.WorldY, chunk.Top, chunk.Bottom, chunk.Left, chunk.Right)

	return nil // chunk
}

func IsMouseInBlockBoundary(block *Block, mx, my float64) bool {
	blockScreenX := block.ScreenX
	blockScreenY := block.ScreenY
	blockWidth := block.Width
	blockHeight := block.Height

	blockRight := float64(blockScreenX) + float64(blockWidth)
	blockBottom := float64(blockScreenY) + float64(blockHeight)

	// Debugging
	// if block == chunk.Blocks[1][1] {
	// 	chunk.Blocks[1][1].IsHoveredOn = true
	// 	fmt.Printf("BWX: %f, BWY: %f, BSX: %f, BSY: %f, BR: %f, BB: %f - MX: %f, MY: %f\n", block.WorldX, block.WorldY, block.ScreenX, block.ScreenY, blockRight, blockBottom, mx, my)
	// }

	if mx > float64(blockScreenX) && mx < blockRight && my > float64(blockScreenY) && my < blockBottom {
		return true
	} else {
		return false
	}
}

// --- Game Functions --- //

func (g *Game) Update() error {

	// - Movement - //

	// Debugging
	// if ebiten.IsKeyPressed(ebiten.KeyUp) {
	// 	g.Plr.WorldY -= g.Camera.Speed
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyDown) {
	// 	g.Plr.WorldY += g.Camera.Speed
	// }

	g.Plr.VelY += g.Plr.Gravity
	g.Plr.WorldY += g.Plr.VelY

	if g.Plr.WorldY+float32(g.Plr.Height) >= float32(HEIGHT) {
		g.Plr.WorldY = float32(HEIGHT) - float32(g.Plr.Height)
		g.Plr.VelY = 0
	}

	// Debugging
	// fmt.Printf("VelY: %f\n", g.Plr.VelY)

	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyW) {
		if g.Plr.WorldY+float32(g.Plr.Height) == float32(HEIGHT) {
			g.Plr.VelY = -15
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Plr.WorldX -= g.Camera.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Plr.WorldX += g.Camera.Speed
	}

	// Recenter camera after player movement
	g.Camera.OffsetX = g.Plr.WorldX
	g.Camera.OffsetY = g.Plr.WorldY

	// Update World to Screen positions
	g.Plr.UpdateScreenPos(g.Camera.OffsetX, g.Camera.OffsetY, float32(WIDTH), float32(HEIGHT))

	for _, chunk := range g.Chunks {
		chunk.UpdateScreenPos(g.Camera.OffsetX, g.Camera.OffsetY)
	}

	// - Chunk Generation - //

	if g.Plr.ScreenY < g.Chunks[*g.CurrChunkCoord].Top || g.Plr.ScreenY > g.Chunks[*g.CurrChunkCoord].Bottom || g.Plr.ScreenX < g.Chunks[*g.CurrChunkCoord].Left || g.Plr.ScreenX > g.Chunks[*g.CurrChunkCoord].Right {
		UpdateCurrChunkCoord()

		CheckSurroundingChunks()

		if TopLeftChunk == nil {
			GenerateChunk(g.CurrChunkCoord.WorldX, g.CurrChunkCoord.WorldY, "top-left", "dirt")
		}

		if TopChunk == nil {
			GenerateChunk(g.CurrChunkCoord.WorldX, g.CurrChunkCoord.WorldY, "top", "stone")
		}

		if TopRightChunk == nil {
			GenerateChunk(g.CurrChunkCoord.WorldX, g.CurrChunkCoord.WorldY, "top-right", "dirt")
		}

		if LeftChunk == nil {
			GenerateChunk(g.CurrChunkCoord.WorldX, g.CurrChunkCoord.WorldY, "left", "stone")
		}

		if RightChunk == nil {
			GenerateChunk(g.CurrChunkCoord.WorldX, g.CurrChunkCoord.WorldY, "right", "stone")
		}

		if BottomLeftChunk == nil {
			GenerateChunk(g.CurrChunkCoord.WorldX, g.CurrChunkCoord.WorldY, "bottom-left", "dirt")
		}

		if BottomChunk == nil {
			GenerateChunk(g.CurrChunkCoord.WorldX, g.CurrChunkCoord.WorldY, "bottom", "stone")
		}

		if BottomRightChunk == nil {
			GenerateChunk(g.CurrChunkCoord.WorldX, g.CurrChunkCoord.WorldY, "bottom-right", "dirt")
		}
	}

	// - Block Breaking / Placing - //

	mx, my := ebiten.CursorPositionF()

	// Check for hovering
	chunk := g.Chunks[*g.CurrChunkCoord]
	for i := range chunk.Blocks {
		for j := range chunk.Blocks[i] {
			block := chunk.Blocks[i][j]
			if chunk != nil {
				val := IsMouseInBlockBoundary(block, mx, my)
				if val {
					(*block).IsHoveredOn = true
				} else {
					(*block).IsHoveredOn = false
				}
			}
		}
	}

	// Check for clicking

	// Breaking
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		chunk := g.Chunks[*g.CurrChunkCoord]
		for i := range chunk.Blocks {
			for j := range chunk.Blocks[i] {
				block := chunk.Blocks[i][j]
				if chunk != nil {
					val := IsMouseInBlockBoundary(block, mx, my)
					if val {

						// Check if the texture is NOT air
						if block.Texture.File != "textures/air.png" {
							airTexture, _, err := LoadTexture("air", &TextureConfig{Path: "textures", Ext: "png"})

							if err != nil {
								log.Fatal(err)
							}

							block.Texture = airTexture

							// Debugging
							fmt.Printf("Clicked Block: %v\n", block)
						}
					}
				}
			}
		}
	}

	// Placing
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton2) {
		chunk := g.Chunks[*g.CurrChunkCoord]
		for i := range chunk.Blocks {
			for j := range chunk.Blocks[i] {
				block := chunk.Blocks[i][j]
				if chunk != nil {
					val := IsMouseInBlockBoundary(block, mx, my)
					if val {

						// Check if the texture IS air
						if block.Texture.File == "textures/air.png" {
							newTexture, _, err := LoadTexture("stone", &TextureConfig{Path: "textures", Ext: "png"})

							if err != nil {
								log.Fatal(err)
							}

							block.Texture = newTexture

							// Debugging
							fmt.Printf("Clicked Block: %v\n", block)
						}
					}
				}
			}
		}
	}

	return nil
}

var num int

func (g *Game) Draw(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, float32(WIDTH), float32(HEIGHT), color.NRGBA{173, 216, 230, 255}, true)

	// keys := make([]int, 0, len(g.Chunks))
	// for k := range g.Chunks {
	// 	keys = append(keys, k)
	// }

	// sort.Ints(keys)

	num = 0
	for k := range g.Chunks {
		chunk := g.Chunks[k]

		// Chunk-Level Culling
		if chunk.Left+float32(chunk.Width) < 0 || chunk.Left > float32(WIDTH) ||
			chunk.Top+float32(chunk.Height) < 0 || chunk.Top > float32(HEIGHT) {
			continue
		}

		// Block-Level Culling
		for i := range chunk.Blocks {
			for j := range chunk.Blocks[i] {
				block := chunk.Blocks[i][j]

				blockScreenX := block.ScreenX
				blockScreenY := block.ScreenY
				blockWidth := block.Width
				blockHeight := block.Height
				blockTexture := block.Texture
				blockIsHoveredOn := block.IsHoveredOn

				if blockScreenX+blockWidth >= 0 && blockScreenX <= float32(WIDTH) && blockScreenY+blockHeight >= 0 && blockScreenY <= float32(chunk.Height) {
					if blockTexture != nil {
						block_op := &ebiten.DrawImageOptions{}
						block_op.GeoM.Translate(float64(blockScreenX), float64(blockScreenY))

						screen.DrawImage(blockTexture.Image, block_op)

						if blockIsHoveredOn {
							vector.StrokeRect(screen, blockScreenX+4, blockScreenY+4, blockWidth-8, blockHeight-8, 4, color.Black, true)
						}

						// DUBUGGING
						// vector.StrokeLine(screen, float32(blockScreenX), 0, float32(blockScreenX), float32(HEIGHT), 5, color.NRGBA{255, 0, 0, 255}, true)
						// vector.StrokeLine(screen, 0, float32(blockScreenY), float32(WIDTH), float32(blockScreenY), 5, color.NRGBA{0, 0, 255, 255}, true)
						num += 1
					}
				}
			}
		}
	}

	// DEBUG
	// vector.StrokeLine(screen, float32(g.Chunks[*g.CurrChunkCoord].Coord.WorldX), 0, float32(g.Chunks[*g.CurrChunkCoord].Coord.WorldX), float32(HEIGHT), 5, color.NRGBA{255, 0, 0, 255}, true)

	// vector.StrokeLine(screen, g.Chunks[*g.CurrChunkCoord].Left, 0, g.Chunks[*g.CurrChunkCoord].Left, float32(HEIGHT), 5, color.NRGBA{0, 255, 0, 255}, true)

	plr_op := &ebiten.DrawImageOptions{}
	plr_op.GeoM.Translate(float64(g.Plr.ScreenX), float64(g.Plr.ScreenY))

	if g.Plr.Texture != nil {
		screen.DrawImage(g.Plr.Texture.Image, plr_op)
	}

	o := fmt.Sprintf("NumOfBlockRendering: %d, PlrWX: %f, PlrWY: %f\nCameraOX: %f, CameraOY: %f, Chunks: %d, CurrChunkCoordWorldX: %d, CurrChunkCoordWorldY: %d", num, g.Plr.WorldX, g.Plr.WorldY, g.Camera.OffsetX, g.Camera.OffsetY, len(g.Chunks), (*g).CurrChunkCoord.WorldX, (*g).CurrChunkCoord.WorldY)
	ebitenutil.DebugPrint(screen, o)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}
