package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yuricorredor/platformer/assets"
	"github.com/yuricorredor/platformer/clouds"
	"github.com/yuricorredor/platformer/entities"
	"github.com/yuricorredor/platformer/tilemap"
)

var (
	gameClouds = clouds.Clouds
)

type Game struct {
	scollX       int
	scrollY      int
	screenWidth  int
	screenHeight int
}

func (g *Game) UpdateScrollPosition() {
	playerRect := entities.Player.Rect()
	g.scollX += (int(playerRect.CenterX()) - g.screenWidth/2 - g.scollX) / 15
	g.scrollY += (int(playerRect.CenterY()) - g.screenHeight/2 - g.scrollY) / 15
}

func (g *Game) Update() error {
	g.UpdateScrollPosition()
	gameClouds.Update()
	entities.Player.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(assets.Assets.Images["background"][0], nil)
	gameClouds.Draw(screen, g.scollX, g.scrollY)
	tilemap.TileMap.Draw(screen, g.scollX, g.scrollY)
	entities.Player.Draw(screen, g.scollX, g.scrollY)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.screenWidth = outsideWidth / 2
	g.screenHeight = outsideHeight / 2
	return g.screenWidth, g.screenHeight
}

func main() {
	game := &Game{}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Platformer")

	gameClouds.GenerateRandomClouds()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
