package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yuricorredor/platformer/assets"
	"github.com/yuricorredor/platformer/rects"
	"github.com/yuricorredor/platformer/tilemap"
	"github.com/yuricorredor/platformer/types"
)

var (
	RENDER_SCALE   = 2
	MOVEMENT_SPEED = 3
	PATH           = "assets/maps/new_map.json"
)

type Editor struct {
	scrollX       int
	scrollY       int
	screenWidth   int
	screenHeight  int
	clicking      bool
	rightClicking bool
	onGrid        bool
	tileList      []string
	tileGroup     int
	tileVariant   int
	position      types.Vector
}

func (e *Editor) Update() error {
	e.HandleCursor()
	e.HandleInputs()

	return nil
}

func (e *Editor) Draw(screen *ebiten.Image) {
	tilemap.TileMap.Draw(screen, e.scrollX, e.scrollY)

	e.DrawLayout(screen)
	e.DrawCurrentTile(screen)
}

func (e *Editor) CurrentTileImage() *ebiten.Image {
	return assets.Assets.Images[e.tileList[e.tileGroup]][e.tileVariant]
}

func (e *Editor) DrawLayout(screen *ebiten.Image) {
	currentTileImage := e.CurrentTileImage()

	options := &ebiten.DrawImageOptions{}
	options.ColorScale.ScaleAlpha(0.65)
	options.GeoM.Translate(5, 5)
	screen.DrawImage(currentTileImage, options)
}

func (e *Editor) DrawCurrentTile(screen *ebiten.Image) {
	currentTileImage := e.CurrentTileImage()
	options := &ebiten.DrawImageOptions{}
	options.ColorScale.ScaleAlpha(0.65)

	if e.onGrid {
		options.GeoM.Translate(float64((int(e.position.X)+e.scrollX)/tilemap.TileMap.TileSize)*float64(tilemap.TileMap.TileSize)-float64(e.scrollX), float64((int(e.position.Y)+e.scrollY)/tilemap.TileMap.TileSize)*float64(tilemap.TileMap.TileSize)-float64(e.scrollY))
	} else {
		options.GeoM.Translate(e.position.X-float64(currentTileImage.Bounds().Dx()/2), e.position.Y-float64(currentTileImage.Bounds().Dy()/2))
	}

	screen.DrawImage(currentTileImage, options)
}

func (e *Editor) HandleCursor() {
	mouseX, mouseY := ebiten.CursorPosition()
	isCursorInsideGameScreen := mouseX >= 0 && mouseX < e.screenWidth && mouseY >= 0 && mouseY < e.screenHeight
	if isCursorInsideGameScreen {
		e.position = types.Vector{X: float64(mouseX), Y: float64(mouseY)}
	}
}

func (e *Editor) HandleInputs() {
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		e.onGrid = !e.onGrid
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		tilemap.TileMap.Save(PATH)
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		e.scrollY -= MOVEMENT_SPEED
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		e.scrollY += MOVEMENT_SPEED
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		e.scrollX -= MOVEMENT_SPEED
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		e.scrollX += MOVEMENT_SPEED
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		e.clicking = true
	} else {
		e.clicking = false
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		e.rightClicking = true
	} else {
		e.rightClicking = false
	}
	_, yoff := ebiten.Wheel()

	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		if yoff > 0 {
			e.tileVariant++
			if e.tileVariant >= len(assets.Assets.Images[e.tileList[e.tileGroup]]) {
				e.tileVariant = 0
			}
		}
		if yoff < 0 {
			e.tileVariant--
			if e.tileVariant < 0 {
				e.tileVariant = len(assets.Assets.Images[e.tileList[e.tileGroup]]) - 1
			}
		}
	} else {
		if yoff > 0 {
			e.tileGroup++
			if e.tileGroup >= len(e.tileList) {
				e.tileGroup = 0
			}
			e.tileVariant = 0
		}
		if yoff < 0 {
			e.tileGroup--
			if e.tileGroup < 0 {
				e.tileGroup = len(e.tileList) - 1
			}
			e.tileVariant = 0
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		e.RemoveTile()
		e.RemoveOffGridTile()
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !e.onGrid {
		e.AddOffgridTile()
	}

	if e.clicking && e.onGrid {
		e.AddTile()
	}
}

func (e *Editor) RemoveTile() {
	tilemap.TileMap.RemoveTile(types.Vector{X: float64((int(e.position.X) + e.scrollX) / tilemap.TileMap.TileSize), Y: float64((int(e.position.Y) + e.scrollY) / tilemap.TileMap.TileSize)})
}

func (e *Editor) RemoveOffGridTile() {
	// for tile in self.tilemap.offgrid_tiles.copy():
	//         tile_img = self.assets[tile['type']][tile['variant']]
	//         tile_rect = pygame.Rect(tile['pos'][0] - self.scroll[0], tile['pos'][1] - self.scroll[1], tile_img.get_width(), tile_img.get_height())
	//         if tile_rect.collidepoint(mouse_pos):
	//           self.tilemap.offgrid_tiles.remove(tile)
	for _, tile := range tilemap.TileMap.OffGridTiles {
		tile_image := assets.Assets.Images[tile.Type][tile.Variant]
		tileRect := rects.Rect{
			X:      tile.Position.X*float64(tilemap.TileMap.TileSize) - float64(e.scrollX),
			Y:      tile.Position.Y*float64(tilemap.TileMap.TileSize) - float64(e.scrollY),
			Width:  float64(tile_image.Bounds().Dx()),
			Height: float64(tile_image.Bounds().Dy()),
		}

		if tileRect.Contains(e.position) {
			tilemap.TileMap.RemoveOffGridTile(tile)
		}
	}
}

func (e *Editor) AddTile() {
	tile := tilemap.Tile{
		Position: types.Vector{X: float64((int(e.position.X) + e.scrollX) / tilemap.TileMap.TileSize), Y: float64((int(e.position.Y) + e.scrollY) / tilemap.TileMap.TileSize)},
		Variant:  e.tileVariant,
		Type:     e.tileList[e.tileGroup],
	}

	tilemap.TileMap.SetTile(tile)
}

func (e *Editor) AddOffgridTile() {
	tile := tilemap.Tile{
		Position: types.Vector{X: (e.position.X + float64(e.scrollX) - float64(e.CurrentTileImage().Bounds().Dx()/2)) / float64(tilemap.TileMap.TileSize), Y: (e.position.Y + float64(e.scrollY) - float64(e.CurrentTileImage().Bounds().Dy()/2)) / float64(tilemap.TileMap.TileSize)},
		Variant:  e.tileVariant,
		Type:     e.tileList[e.tileGroup],
	}

	tilemap.TileMap.SetOffGridTile(tile)
}

func (e *Editor) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	e.screenWidth = outsideWidth / RENDER_SCALE
	e.screenHeight = outsideHeight / RENDER_SCALE
	return e.screenWidth, e.screenHeight
}

func NewEditor() *Editor {
	return &Editor{
		onGrid: true,
	}
}

func main() {
	editor := NewEditor()

	tilemap.TileMap.Load(PATH)

	editor.tileList = []string{"grass", "stone", "decor", "large_decor"}

	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Platformer Editor")

	if err := ebiten.RunGame(editor); err != nil {
		panic(err)
	}
}
