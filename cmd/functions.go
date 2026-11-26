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
var floorHeight float32 = float32(HEIGHT/2) + float32(BLOCK_SIZE/2)
var roofHeight float32 = float32(HEIGHT / 2)
var leftLimit float32 = float32(math.Inf(-1)) // - inf
var rightLimit float32 = float32(math.Inf(1)) // + inf

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

		HandleFloorCollision()

		if g.Plr.WorldY+float32(g.Plr.Height) >= floorHeight {
			g.Plr.WorldY = floorHeight - float32(g.Plr.Height)
			g.Plr.VelY = 0
		}

		if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyW) {
			if int(g.Plr.WorldY+float32(g.Plr.Height)) == int(floorHeight) {
				g.Plr.VelY = -15
			}
		}

		HandleSideCollisions()

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

	if flying {
		// Debugging
		if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.Plr.WorldY -= g.Camera.Speed
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
			g.Plr.WorldY += g.Camera.Speed
		}
	}

	// Recenter camera after player movement
	g.Camera.OffsetX = g.Plr.WorldX
	g.Camera.OffsetY = g.Plr.WorldY

	// Update World to Screen positions
	g.Plr.UpdateScreenPos(g.Camera.OffsetX, g.Camera.OffsetY, float32(WIDTH), float32(HEIGHT))

	// OPTIMIZE LATER!!! Needed for now (see TESTING below)
	for _, chunk := range g.Chunks {
		chunk.UpdateScreenPos(g.Camera.OffsetX, g.Camera.OffsetY)
	}

	// TESTING: Only updating screen pos of surrounding chunks (Kinda works - not really)
	// chunkOffsets := [][]int{{-HEIGHT, -WIDTH}, {-HEIGHT, 0}, {-HEIGHT, WIDTH}, {0, -WIDTH}, {0, WIDTH}, {0, 0}, {HEIGHT, -WIDTH}, {HEIGHT, 0}, {HEIGHT, WIDTH}}
	// for _, chunkOffset := range chunkOffsets {
	// 	nx := g.CurrChunkCoord.WorldX + chunkOffset[1]
	// 	ny := g.CurrChunkCoord.WorldY + chunkOffset[0]

	// 	if chunk, exists := g.Chunks[ChunkCoord{WorldX: nx, WorldY: ny}]; exists {
	// 		chunk.UpdateScreenPos(g.Camera.OffsetX, g.Camera.OffsetY)
	// 	}
	// }

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

		// FIX Later: Clean up the hovering outline left over from moving into a new chunk
		chunkOffsets := [][]int{{-HEIGHT, 0}, {0, -WIDTH}, {0, WIDTH}, {HEIGHT, 0}}
		for _, chunkOffset := range chunkOffsets {
			nx := g.CurrChunkCoord.WorldX + chunkOffset[1]
			ny := g.CurrChunkCoord.WorldY + chunkOffset[0]

			if chunk, exists := g.Chunks[ChunkCoord{WorldX: nx, WorldY: ny}]; exists {
				for i := range chunk.Blocks {
					for j := range chunk.Blocks[i] {
						block := chunk.Blocks[i][j]
						blockIsHoveredOn := block.IsHoveredOn
						if blockIsHoveredOn {
							(*block).IsHoveredOn = false
						}
					}
				}
			}
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

	// Block Breaking
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
							// Check if AT LEAST ONE of the 4 Blocks around it (Top, Bottom, Left, Right) have a NON-AIR texture
							offsets := [][]int{{-1, 0}, {0, -1}, {0, 0}, {0, 1}, {1, 0}}

							foundNonAirTexture := false
							for _, offset := range offsets {
								ni := i + offset[0]
								nj := j + offset[1]

								if ni >= 0 && ni < len(chunk.Blocks) && nj >= 0 && nj < len(chunk.Blocks[0]) {
									foundBlock := chunk.Blocks[ni][nj]
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

								block.Texture = airTexture

								// Debugging
								// fmt.Printf("Clicked Block: %v\n", block)
							}
						}
					}
				}
			}
		}
	}

	// Block Placing
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
							// Check if AT LEAST ONE of the 4 Blocks around it (Top, Bottom, Left, Right) have a NON-AIR texture
							offsets := [][]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}

							foundNonAirTexture := false
							for _, offset := range offsets {
								ni := i + offset[0]
								nj := j + offset[1]

								if ni >= 0 && ni < len(chunk.Blocks) && nj >= 0 && nj < len(chunk.Blocks[ni]) {
									foundBlock := chunk.Blocks[ni][nj]
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

								block.Texture = newTexture

								// Debugging
								// fmt.Printf("Clicked Block: %v\n", block)
							}
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
