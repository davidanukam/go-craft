package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// --- Game Functions --- //

var flying bool = false

var floorHeight float32 = float32(math.Inf(1)) // +inf
var roofHeight float32 = float32(math.Inf(-1)) // -inf
var leftLimit float32 = float32(math.Inf(-1))  // -inf
var rightLimit float32 = float32(math.Inf(1))  // +inf

var hoveredBlock *Block = nil
var hoveredBlockChunk *Chunk = nil

var chunkOffsets = [9][2]int{
	{-HEIGHT, -WIDTH}, {-HEIGHT, 0}, {-HEIGHT, WIDTH},
	{0, -WIDTH}, {0, 0}, {0, WIDTH},
	{HEIGHT, -WIDTH}, {HEIGHT, 0}, {HEIGHT, WIDTH},
}

func (g *Game) Update() error {

	// - Movement - //

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		flying = !flying
		floorHeight = float32(HEIGHT/2) + float32(BLOCK_SIZE/2)
		roofHeight = float32(HEIGHT / 2)
	}

	// DEBUGGING (Using center of screen as floorHeight)
	// floorHeight := float32(HEIGHT/2) + float32(BLOCK_SIZE/2)

	if !flying {
		g.Plr.VelY += g.Plr.Gravity
		g.Plr.WorldY += g.Plr.VelY

		if g.Plr.VelY < 0 {
			if g.Plr.WorldY <= roofHeight {
				g.Plr.WorldY = roofHeight
				g.Plr.VelY = 0
			}
		}

		HandleVerticalCollisions()

		if g.Plr.WorldY+float32(g.Plr.Height) >= floorHeight {
			g.Plr.WorldY = floorHeight - float32(g.Plr.Height)
			g.Plr.VelY = 0
		}

		if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyW) {
			if int(g.Plr.WorldY+float32(g.Plr.Height)) == int(floorHeight) {
				if g.Plr.WorldY > roofHeight {
					g.Plr.VelY = -15
				}
			}
		}

		HandleHorizontalCollisions()

		if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			if g.Plr.WorldX > leftLimit {
				g.Plr.WorldX -= g.Camera.Speed
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			if g.Plr.WorldX+float32(g.Plr.Width) < rightLimit {
				g.Plr.WorldX += g.Camera.Speed
			}
		}

	}

	// Used for Debugging
	if flying {
		if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.Plr.WorldY -= g.Camera.Speed
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
			g.Plr.WorldY += g.Camera.Speed
		}

		if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			g.Plr.WorldX -= g.Camera.Speed
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			g.Plr.WorldX += g.Camera.Speed
		}
	}

	// Recenter camera after player movement
	g.Camera.OffsetX = g.Plr.WorldX
	g.Camera.OffsetY = g.Plr.WorldY

	// Update Player World to Screen positions
	g.Plr.UpdateScreenPos(g.Camera.OffsetX, g.Camera.OffsetY, float32(WIDTH), float32(HEIGHT))

	// Updating Chunk World to Screen positions
	for _, offset := range chunkOffsets {
		coord := ChunkCoord{
			WorldX: g.CurrChunkCoord.WorldX + offset[1],
			WorldY: g.CurrChunkCoord.WorldY + offset[0],
		}
		if chunk, exists := g.Chunks[coord]; exists {
			chunk.UpdateScreenPos(g.Camera.OffsetX, g.Camera.OffsetY)
		}
	}

	// - Chunk Generation - //

	HandleChunkGeneration()

	// - Block Interaction  - //

	mx, my := ebiten.CursorPositionF()

	// Check for hovering
	if chunk, exists := g.Chunks[*g.CurrChunkCoord]; exists {
		for i := range chunk.Blocks {
			for j := range chunk.Blocks[i] {
				block := chunk.Blocks[i][j]
				if block != nil {
					val := IsMouseInBlockBoundary(block, mx, my)
					if val {
						(*block).IsHoveredOn = true
						hoveredBlock = block
						hoveredBlockChunk = chunk
					} else {
						(*block).IsHoveredOn = false
					}
				}
			}
		}
	}

	// Block Breaking
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		if hoveredBlock != nil && hoveredBlockChunk != nil {
			for i := range hoveredBlockChunk.Blocks {
				for j := range hoveredBlockChunk.Blocks[i] {
					// Check if the texture is NOT air
					if hoveredBlock.Texture.File != "textures/air.png" {
						// Check if AT LEAST ONE of the 4 Blocks around it (Top, Bottom, Left, Right) have a NON-AIR texture
						offsets := [][]int{{-1, 0}, {0, -1}, {0, 0}, {0, 1}, {1, 0}}

						foundNonAirTexture := false
						for _, offset := range offsets {
							ni := i + offset[0]
							nj := j + offset[1]

							if ni >= 0 && ni < len(hoveredBlockChunk.Blocks) && nj >= 0 && nj < len(hoveredBlockChunk.Blocks[0]) {
								foundBlock := hoveredBlockChunk.Blocks[ni][nj]
								if foundBlock.Texture.File != "textures/air.png" {
									foundNonAirTexture = true
									break
									// DEBUGGING
									// fmt.Printf("Found Block - %v\n", foundBlock.Texture.File)
								}
							}
						}

						if foundNonAirTexture {
							airTexture, _, err := LoadTexture("air", &TextureConfig{Path: "textures", Ext: "png"})

							if err != nil {
								log.Fatal(err)
							}

							hoveredBlock.Texture = airTexture

							// Debugging
							// fmt.Printf("Clicked Block: %v\n", hoveredBlock)
						}
					}
				}
			}
		}
	}

	// Block Placing
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton2) {
		if hoveredBlock != nil && hoveredBlockChunk != nil {
			for i := range hoveredBlockChunk.Blocks {
				for j := range hoveredBlockChunk.Blocks[i] {
					// Check if the texture IS air
					if hoveredBlock.Texture.File == "textures/air.png" {
						// Check if AT LEAST ONE of the 4 Blocks around it (Top, Bottom, Left, Right) have a NON-AIR texture
						offsets := [][]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}

						foundNonAirTexture := false
						for _, offset := range offsets {
							ni := i + offset[0]
							nj := j + offset[1]

							if ni >= 0 && ni < len(hoveredBlockChunk.Blocks) && nj >= 0 && nj < len(hoveredBlockChunk.Blocks[ni]) {
								foundBlock := hoveredBlockChunk.Blocks[ni][nj]
								if foundBlock.Texture.File != "textures/air.png" {
									foundNonAirTexture = true
									break
									// DEBUGGING
									// fmt.Printf("Found Block - %v\n", foundBlock.Texture.File)
								}

							}
						}

						if foundNonAirTexture {
							newTexture, _, err := LoadTexture("cobblestone", &TextureConfig{Path: "textures", Ext: "png"})

							if err != nil {
								log.Fatal(err)
							}

							hoveredBlock.Texture = newTexture

							// Debugging
							// fmt.Printf("Clicked Block: %v\n", hoveredBlock)
						}
					}
				}
			}
		}
	}

	return nil
}

var numberOfBlocksDrawn int

func (g *Game) Draw(screen *ebiten.Image) {
	// sky
	vector.FillRect(screen, 0, 0, float32(WIDTH), float32(HEIGHT), color.NRGBA{149, 203, 231, 255}, true)

	// DEBUGGING
	// vector.StrokeLine(screen, 0, roofHeight, float32(WIDTH), roofHeight, 5, color.White, true)

	// vector.StrokeLine(screen, leftLimit, 0, leftLimit, float32(HEIGHT), 5, color.NRGBA{255, 0, 0, 255}, true)

	// vector.StrokeLine(screen, rightLimit, 0, rightLimit, float32(HEIGHT), 5, color.NRGBA{0, 0, 255, 255}, true)

	// keys := make([]int, 0, len(g.Chunks))
	// for k := range g.Chunks {
	// 	keys = append(keys, k)
	// }
	// sort.Ints(keys)

	// Only check the current chunk and the 8 arround it when drawing
	numberOfBlocksDrawn = 0
	chunkOffsets := [][]int{{-HEIGHT, -WIDTH}, {-HEIGHT, 0}, {-HEIGHT, WIDTH}, {0, -WIDTH}, {0, WIDTH}, {0, 0}, {HEIGHT, -WIDTH}, {HEIGHT, 0}, {HEIGHT, WIDTH}}
	for _, chunkOffset := range chunkOffsets {
		nx := g.CurrChunkCoord.WorldX + chunkOffset[1]
		ny := g.CurrChunkCoord.WorldY + chunkOffset[0]

		if chunk, exists := g.Chunks[ChunkCoord{WorldX: nx, WorldY: ny}]; exists {

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
						if blockTexture != nil && blockTexture.File != "textures/air.png" {
							block_op := &ebiten.DrawImageOptions{}
							block_op.GeoM.Translate(float64(blockScreenX), float64(blockScreenY))

							screen.DrawImage(blockTexture.Image, block_op)

							numberOfBlocksDrawn += 1

							// DUBUGGING
							// vector.StrokeLine(screen, float32(blockScreenX), 0, float32(blockScreenX), float32(HEIGHT), 5, color.NRGBA{255, 0, 0, 255}, true)
							// vector.StrokeLine(screen, 0, float32(blockScreenY), float32(WIDTH), float32(blockScreenY), 5, color.NRGBA{0, 0, 255, 255}, true)
						}
						// Draw hover outline if hovering on block
						if blockIsHoveredOn {
							vector.StrokeRect(screen, blockScreenX+4, blockScreenY+4, blockWidth-8, blockHeight-8, 4, color.Black, true)
						}
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

	fps := ebiten.ActualFPS()
	tps := ebiten.ActualTPS()
	o := fmt.Sprintf("NumOfBlockRendering: %d, PlrWX: %f, PlrWY: %f\nCameraOX: %f, CameraOY: %f, Chunks: %d\nCurrChunkCoordWorldX: %d, CurrChunkCoordWorldY: %d, FPS: %f, TPS: %f\n", numberOfBlocksDrawn, g.Plr.WorldX, g.Plr.WorldY, g.Camera.OffsetX, g.Camera.OffsetY, len(g.Chunks), (*g).CurrChunkCoord.WorldX, (*g).CurrChunkCoord.WorldY, fps, tps)

	ebitenutil.DebugPrint(screen, o)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}
