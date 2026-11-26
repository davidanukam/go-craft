package main

import (
	"log"
	"math"
	"slices"
)

// --- Helper Functions --- //

// Update game.CurrChunkCoord (pointer) to store the current chunk the player is in
func UpdateCurrChunkCoord() bool {
	if game == nil || game.CurrChunkCoord == nil || game.Plr == nil {
		return false
	}

	baseX := game.CurrChunkCoord.WorldX
	baseY := game.CurrChunkCoord.WorldY
	px := game.Plr.ScreenX
	py := game.Plr.ScreenY

	for _, offset := range chunkOffsets {
		coord := ChunkCoord{
			WorldX: baseX + offset[1],
			WorldY: baseY + offset[0],
		}
		chunk, exists := game.Chunks[coord]
		if !exists || chunk == nil {
			continue
		}

		if py >= chunk.Top && py <= chunk.Bottom && px >= chunk.Left && px <= chunk.Right {
			changed := chunk.Coord.WorldX != baseX || chunk.Coord.WorldY != baseY
			game.CurrChunkCoord.WorldX = chunk.Coord.WorldX
			game.CurrChunkCoord.WorldY = chunk.Coord.WorldY
			return changed
		}
	}

	return false
}

type NeighborGrid struct {
	TopLeft, Top, TopRight          *Chunk
	Left, Center, Right             *Chunk
	BottomLeft, Bottom, BottomRight *Chunk
}

// Check the 8 surrounding chunks
func GetSurroundingChunks(base ChunkCoord, stepX, stepY int) NeighborGrid {
	get := func(dx, dy int) *Chunk {
		c, ok := game.Chunks[ChunkCoord{
			WorldX: base.WorldX + dx,
			WorldY: base.WorldY + dy,
		}]
		if !ok {
			return nil
		}
		return c
	}

	return NeighborGrid{
		TopLeft:     get(-stepX, -stepY),
		Top:         get(0, -stepY),
		TopRight:    get(stepX, -stepY),
		Left:        get(-stepX, 0),
		Center:      get(0, 0),
		Right:       get(stepX, 0),
		BottomLeft:  get(-stepX, stepY),
		Bottom:      get(0, stepY),
		BottomRight: get(stepX, stepY),
	}
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

	// fmt.Printf("New Chunk generated [%s] - CoordWorldX: %d, CoordWorldY: %d, Top: %f, Bottom: %f, Left: %f, Right: %f\n", side, chunk.Coord.WorldX, chunk.Coord.WorldY, chunk.Top, chunk.Bottom, chunk.Left, chunk.Right)

	return nil // chunk
}

type neighborSpec struct {
	dir   string
	mat   string
	chunk *Chunk
}

// Generate new chunks when/where needed
func HandleChunkGeneration() {
	if game == nil || game.CurrChunkCoord == nil || game.Plr == nil {
		return
	}

	currChunk, ok := game.Chunks[*game.CurrChunkCoord]
	if !ok || currChunk == nil {
		return
	}

	outOfBounds := game.Plr.ScreenY < currChunk.Top ||
		game.Plr.ScreenY > currChunk.Bottom ||
		game.Plr.ScreenX < currChunk.Left ||
		game.Plr.ScreenX > currChunk.Right

	if !outOfBounds {
		return
	}

	if !UpdateCurrChunkCoord() {
		return
	}

	neighbors := GetSurroundingChunks(*game.CurrChunkCoord, WIDTH, HEIGHT)

	specs := [8]neighborSpec{
		{"top-left", "dirt", neighbors.TopLeft},
		{"top", "stone", neighbors.Top},
		{"top-right", "dirt", neighbors.TopRight},
		{"left", "stone", neighbors.Left},
		{"right", "stone", neighbors.Right},
		{"bottom-left", "dirt", neighbors.BottomLeft},
		{"bottom", "stone", neighbors.Bottom},
		{"bottom-right", "dirt", neighbors.BottomRight},
	}

	chunksGen := 0
	for _, s := range specs {
		if s.chunk == nil {
			GenerateChunk(game.CurrChunkCoord.WorldX, game.CurrChunkCoord.WorldY, s.dir, s.mat)
			chunksGen++
		}
	}

	// DEBUGGING
	// fmt.Printf("Chunks generated: %d\n", chunksGen)
}

// Check if the mouse (cursor) is within a block
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

var verticalChunkOffsets = [][]int{{-HEIGHT, 0}, {0, 0}, {HEIGHT, 0}}

// Handle floor and roof collisions
func HandleVerticalCollisions() {
	var possibleFloorLimits []float32 = []float32{}
	var possibleRoofLimits []float32 = []float32{}

	for _, chunkOffset := range verticalChunkOffsets {
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

var horizontalChunkOffsets = [][]int{{0, -WIDTH}, {0, 0}, {0, WIDTH}}

// Handle left and right block collisions
func HandleHorizontalCollisions() {
	var possibleLeftLimits []float32 = []float32{}
	var possibleRightLimits []float32 = []float32{}

	for _, chunkOffset := range horizontalChunkOffsets {
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
