package tilemap

import (
	"encoding/json"
	"os"
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
}

type Vector = types.Vector

type Tile struct {
	Position types.Vector
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

func (t *TileMapType) Draw(screen *ebiten.Image, scrollX, scrollY int, renderContext string) {
	for _, tile := range t.OffGridTiles {
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(tile.Position.X*float64(t.TileSize)-float64(scrollX), tile.Position.Y*float64(t.TileSize)-float64(scrollY))

		var shouldRender bool
		if renderContext == "game" {
			shouldRender = assets.Assets.Images[tile.Type].ShouldRenderOnGame
		} else if renderContext == "editor" {
			shouldRender = assets.Assets.Images[tile.Type].ShouldRenderOnEditor
		}

		if shouldRender {
			screen.DrawImage(assets.Assets.Images[tile.Type].Image[tile.Variant], options)
		}
	}

	screenWidth := screen.Bounds().Max.X
	screenHeight := screen.Bounds().Max.Y

	for x := scrollX / t.TileSize; x < (scrollX+screenWidth)/t.TileSize+1; x++ {
		for y := scrollY / t.TileSize; y < (scrollY+screenHeight)/t.TileSize+1; y++ {
			location := strconv.Itoa(x) + ";" + strconv.Itoa(y)
			if tile, ok := t.Tiles[location]; ok {
				options := &ebiten.DrawImageOptions{}
				options.GeoM.Translate(tile.Position.X*float64(t.TileSize)-float64(scrollX), tile.Position.Y*float64(t.TileSize)-float64(scrollY))

				var shouldRender bool
				if renderContext == "game" {
					shouldRender = assets.Assets.Images[tile.Type].ShouldRenderOnGame
				} else if renderContext == "editor" {
					shouldRender = assets.Assets.Images[tile.Type].ShouldRenderOnEditor
				}

				if shouldRender {
					screen.DrawImage(assets.Assets.Images[tile.Type].Image[tile.Variant], options)
				}
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
					X:      tile.Position.X * float64(t.TileSize),
					Y:      tile.Position.Y * float64(t.TileSize),
					Width:  float64(t.TileSize),
					Height: float64(t.TileSize),
				}

				rectsList = append(rectsList, rect)
			}
		}
	}

	return rectsList
}

func (t *TileMapType) CheckForSolid(position types.Vector) bool {
	tileLoc := types.Vector{X: position.X / float64(t.TileSize), Y: position.Y / float64(t.TileSize)}
	location := strconv.Itoa(int(tileLoc.X)) + ";" + strconv.Itoa(int(tileLoc.Y))
	if tile, ok := t.Tiles[location]; ok {
		for _, physicsTile := range PhysicsTiles {
			if tile.Type == physicsTile {
				return true
			}
		}
	}

	return false
}

func (t *TileMapType) Extract(pairs []types.Pair, keep bool) []Tile {
	matches := []Tile{}

	for _, pair := range pairs {
		for _, tile := range t.OffGridTiles {
			if tile.Type == pair.AssetType && tile.Variant == pair.AssetVariant {
				matches = append(matches, tile)
				if !keep {
					t.RemoveOffGridTile(tile)
				}
			}
		}

		for _, tile := range t.Tiles {
			if tile.Type == pair.AssetType && tile.Variant == pair.AssetVariant {
				toAppend := tile
				toAppend.Position.X *= float64(t.TileSize)
				toAppend.Position.Y *= float64(t.TileSize)
				matches = append(matches, toAppend)

				if !keep {
					t.RemoveTile(types.Vector{X: tile.Position.X, Y: tile.Position.Y})
				}
			}
		}
	}

	return matches
}

func (t *TileMapType) SetTile(tile Tile) {
	t.Tiles[strconv.Itoa(int(tile.Position.X))+";"+strconv.Itoa(int(tile.Position.Y))] = tile
}

func (t *TileMapType) RemoveTile(position types.Vector) {
	delete(t.Tiles, strconv.Itoa(int(position.X))+";"+strconv.Itoa(int(position.Y)))
}

func (t *TileMapType) SetOffGridTile(tile Tile) {
	t.OffGridTiles = append(t.OffGridTiles, tile)
}

func (t *TileMapType) RemoveOffGridTile(tile Tile) {
	for i, offGridTile := range t.OffGridTiles {
		if offGridTile.Position == tile.Position {
			t.OffGridTiles = append(t.OffGridTiles[:i], t.OffGridTiles[i+1:]...)
			break
		}
	}
}

func (t *TileMapType) ToJSONString() string {
	jsonString, err := json.MarshalIndent(t, "", "")
	if err != nil {
		panic(err)
	}

	return string(jsonString)
}

func (t *TileMapType) Save(path string) error {
	// Create file if not exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Create(path)
	}

	// Open file
	f, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}

	// Write to file
	_, err = f.WriteString(t.ToJSONString())
	if err != nil {
		return err
	}

	// Close file
	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func (t *TileMapType) Load(path string) error {
	// Open file
	f, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	// Read file
	var tileMap TileMapType
	err = json.NewDecoder(f).Decode(&tileMap)
	if err != nil {
		return err
	}

	// Close file
	err = f.Close()
	if err != nil {
		return err
	}

	// Set tilemap
	TileMap = &tileMap

	return nil
}
