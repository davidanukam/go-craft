package main

import (
	"log"
)

type ChunkCoord struct {
	WorldX int
	WorldY int
}

type Chunk struct {
	Coord  ChunkCoord
	Blocks [][]*Block
	Width  int
	Height int
	Top    float32
	Bottom float32
	Left   float32
	Right  float32
}

func NewChunk(coord ChunkCoord, camX, camY float32, width, height int, ty string) (*Chunk, error) {
	var chunk Chunk

	chunk.Coord = coord

	chunk.Top = float32(chunk.Coord.WorldY) - camY + float32(chunk.Height/2)
	chunk.Bottom = float32(chunk.Coord.WorldY) + float32(chunk.Height) - camY + float32(chunk.Height/2)
	chunk.Left = float32(chunk.Coord.WorldX) - camX + float32(chunk.Width/2)
	chunk.Right = float32(chunk.Coord.WorldX) + float32(chunk.Width) - camX + float32(chunk.Width/2)

	chunk.Width = width
	chunk.Height = height

	chunk.Blocks = make([][]*Block, int(width*height))

	var block *Block
	var err error

	var blockSize int = 32

	for i := range height / blockSize {
		for j := range width / blockSize {
			block, err = NewBlock(float32(chunk.Coord.WorldX+j*32), float32(chunk.Coord.WorldY+i*32), 0, 0, 32, 32, ty) // "grass"

			if err != nil {
				log.Fatal(err)
			}

			chunk.Blocks[i] = append(chunk.Blocks[i], block)
		}
	}

	for i := range chunk.Blocks {
		for j := range chunk.Blocks[i] {
			chunk.Blocks[i][j].UpdateScreenPos(camX, camY, chunk.Width, chunk.Height)
		}
	}

	return &chunk, nil
}

func (c *Chunk) UpdateScreenPos(camX, camY float32) {
	c.Top = float32(c.Coord.WorldY) - camY + float32(c.Height/2)
	c.Bottom = float32(c.Coord.WorldY) + float32(c.Height) - camY + float32(c.Height/2)
	c.Left = float32(c.Coord.WorldX) - camX + float32(c.Width/2)
	c.Right = float32(c.Coord.WorldX) + float32(c.Width) - camX + float32(c.Width/2)

	for i := range c.Blocks {
		for j := range c.Blocks[i] {
			c.Blocks[i][j].UpdateScreenPos(camX, camY, c.Width, c.Height)
		}
	}
}
