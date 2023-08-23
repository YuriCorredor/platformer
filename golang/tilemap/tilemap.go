package tilemap

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yuricorredor/platformer/assets"
	"github.com/yuricorredor/platformer/rects"
	"github.com/yuricorredor/platformer/types"
)

var (
	NeighboursOffset = []types.Vector{
		{X: -1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: -1},
		{X: 1, Y: -1}, {X: 1, Y: 0}, {X: 1, Y: 1},
		{X: 0, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: 1},
		{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1},
	}
	PhysicsTiles = []string{"grass", "stone"}
)

var TileMap = &TileMapType{
	TileSize: 16,
	Tiles:    random(),
}

type Vector = types.Vector

type Tile struct {
	Position types.Vector
	Size     int
	Variant  int
	Type     string
}

type TileMapType struct {
	TileSize     int
	Tiles        map[string]Tile
	OffGridTiles []Tile
}

func (t *TileMapType) Update() error {
	return nil
}

func (t *TileMapType) Draw(screen *ebiten.Image, scrollX, scrollY int) {
	for _, tile := range t.OffGridTiles {
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(tile.Position.X*float64(tile.Size)-float64(scrollX), tile.Position.Y*float64(tile.Size)-float64(scrollY))
		screen.DrawImage(assets.Assets.Images[tile.Type][tile.Variant], options)
	}

	screenWidth := screen.Bounds().Max.X
	screenHeight := screen.Bounds().Max.Y

	for x := scrollX / t.TileSize; x < (scrollX+screenWidth)/t.TileSize+1; x++ {
		for y := scrollY / t.TileSize; y < (scrollY+screenHeight)/t.TileSize+1; y++ {
			location := strconv.Itoa(x) + ";" + strconv.Itoa(y)
			if tile, ok := t.Tiles[location]; ok {
				options := &ebiten.DrawImageOptions{}
				options.GeoM.Translate(tile.Position.X*float64(tile.Size)-float64(scrollX), tile.Position.Y*float64(tile.Size)-float64(scrollY))
				screen.DrawImage(assets.Assets.Images[tile.Type][tile.Variant], options)
			}
		}
	}
}

func (t *TileMapType) TilesAroundPosition(position types.Vector) []Tile {
	tiles := []Tile{}
	tileLoc := types.Vector{X: position.X / float64(t.TileSize), Y: position.Y / float64(t.TileSize)}
	for _, offset := range NeighboursOffset {
		checkLocation := strconv.Itoa(int(tileLoc.X+offset.X)) + ";" + strconv.Itoa(int(tileLoc.Y+offset.Y))
		if tile, ok := t.Tiles[checkLocation]; ok {
			tiles = append(tiles, tile)
		}
	}

	return tiles
}

func (t *TileMapType) PhysicsRectsAroundPosition(position types.Vector) []rects.Rect {
	rectsList := []rects.Rect{}
	for _, tile := range t.TilesAroundPosition(position) {
		for _, physicsTile := range PhysicsTiles {
			if tile.Type == physicsTile {
				rect := rects.Rect{
					X:      tile.Position.X * float64(tile.Size),
					Y:      tile.Position.Y * float64(tile.Size),
					Width:  float64(tile.Size),
					Height: float64(tile.Size),
				}

				rectsList = append(rectsList, rect)
			}
		}
	}

	return rectsList
}

func random() map[string]Tile {
	tiles := map[string]Tile{}
	for i := 0; i < 10; i++ {
		tiles[strconv.Itoa(3+i)+";"+strconv.Itoa(10)] = Tile{Position: types.Vector{X: float64(3 + i), Y: 10}, Variant: 1, Type: "grass", Size: 16}
		tiles[strconv.Itoa(10)+";"+strconv.Itoa(5+i)] = Tile{Position: types.Vector{X: 10, Y: float64(5 + i)}, Variant: 1, Type: "stone", Size: 16}
	}

	return tiles
}
