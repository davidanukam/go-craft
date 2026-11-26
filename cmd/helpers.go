package main

import (
	"fmt"
	"log"
	"math"
	"slices"
)

// --- Helper Functions --- //

// Update game.CurrChunkCoord (pointer) to store the current chunk the player is in
func UpdateCurrChunkCoord() {
	chunkOffsets := [][]int{{-HEIGHT, -WIDTH}, {-HEIGHT, 0}, {-HEIGHT, WIDTH}, {0, -WIDTH}, {0, WIDTH}, {0, 0}, {HEIGHT, -WIDTH}, {HEIGHT, 0}, {HEIGHT, WIDTH}}
	for _, chunkOffset := range chunkOffsets {
		if chunk, exists := game.Chunks[ChunkCoord{WorldX: game.CurrChunkCoord.WorldX + chunkOffset[0], WorldY: game.CurrChunkCoord.WorldY + chunkOffset[1]}]; exists {
			if game.Plr.ScreenY >= chunk.Top && game.Plr.ScreenY <= chunk.Bottom && game.Plr.ScreenX >= chunk.Left && game.Plr.ScreenX <= chunk.Right {
				// Debugging
				// fmt.Printf("UpdatedCurrChunkCoord - Chunk_WorldX: %d, Chunk_WorldY: %d, Chunk_Top: %f, Chunk_Bottom: %f, Chunk_Left: %f, Chunk_Right: %f\n", chunk.Coord.WorldX, chunk.Coord.WorldY, chunk.Top, chunk.Bottom, chunk.Left, chunk.Right)
				game.CurrChunkCoord.WorldX = chunk.Coord.WorldX
				game.CurrChunkCoord.WorldY = chunk.Coord.WorldY
				break
			}
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
		if math.Abs(float64(block.WorldX)-float64(game.Plr.WorldX)) >= float64(BLOCK_SIZE) || math.Abs(float64(block.WorldY)-float64(game.Plr.WorldY)) >= float64(BLOCK_SIZE) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func HandleVerticalCollisions() {
	var possibleFloorLimits []float32 = []float32{}
	var possibleRoofLimits []float32 = []float32{}

	chunkOffsets := [][]int{{-HEIGHT, 0}, {0, 0}, {HEIGHT, 0}}
	for _, chunkOffset := range chunkOffsets {
		nx := game.CurrChunkCoord.WorldX + chunkOffset[1]
		ny := game.CurrChunkCoord.WorldY + chunkOffset[0]

		if chunk, exists := game.Chunks[ChunkCoord{WorldX: nx, WorldY: ny}]; exists {
			for i := range chunk.Blocks {
				for j := range chunk.Blocks[i] {
					block := chunk.Blocks[i][j]

					// Get floor Height
					if block.WorldY+block.Height > game.Plr.WorldY+float32(game.Plr.Height) {
						if math.Abs(float64(game.Plr.WorldX)-float64(block.WorldX)) < float64(BLOCK_SIZE) {
							if block.Texture.File != "textures/air.png" {
								possibleFloorLimits = append(possibleFloorLimits, block.WorldY)
							}
						}
					}

					// Get roof Height
					if block.WorldY+block.Height <= game.Plr.WorldY {
						if math.Abs(float64(game.Plr.WorldX)-float64(block.WorldX)) < float64(BLOCK_SIZE) {
							if block.Texture.File != "textures/air.png" {
								possibleRoofLimits = append(possibleRoofLimits, block.WorldY+block.Height)
							}
						}
					}
				}
			}
		}
	}

	if len(possibleFloorLimits) > 0 {
		floorHeight = slices.Min(possibleFloorLimits)
	} else {
		floorHeight = float32(math.Inf(1))
	}

	if len(possibleRoofLimits) > 0 {
		roofHeight = slices.Max(possibleRoofLimits)
	} else {
		roofHeight = float32(math.Inf(-1))
	}
}

func HandleHorizontalCollisions() {
	var possibleLeftLimits []float32 = []float32{}
	var possibleRightLimits []float32 = []float32{}

	chunkOffsets := [][]int{{0, -WIDTH}, {0, 0}, {0, WIDTH}}
	for _, chunkOffset := range chunkOffsets {
		nx := game.CurrChunkCoord.WorldX + chunkOffset[1]
		ny := game.CurrChunkCoord.WorldY + chunkOffset[0]

		if chunk, exists := game.Chunks[ChunkCoord{WorldX: nx, WorldY: ny}]; exists {
			for i := range chunk.Blocks {
				for j := range chunk.Blocks[i] {
					block := chunk.Blocks[i][j]

					if math.Abs(float64(block.WorldY)-float64(game.Plr.WorldY)) < float64(BLOCK_SIZE) && block.WorldY <= game.Plr.WorldY {
						// Get left Limits
						if block.WorldX+block.Width <= game.Plr.WorldX {
							if block.Texture.File != "textures/air.png" {
								possibleLeftLimits = append(possibleLeftLimits, block.WorldX+block.Width)
							}
						}

						// Get right Limits
						if block.WorldX >= game.Plr.WorldX+float32(game.Plr.Width) {
							if block.Texture.File != "textures/air.png" {
								possibleRightLimits = append(possibleRightLimits, block.WorldX)
							}
						}
					}
				}
			}
		}
	}

	if len(possibleLeftLimits) > 0 {
		leftLimit = slices.Max(possibleLeftLimits)
	} else {
		leftLimit = float32(math.Inf(-1))
	}

	if len(possibleRightLimits) > 0 {
		rightLimit = slices.Min(possibleRightLimits)
	} else {
		rightLimit = float32(math.Inf(1))
	}
}
