package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yuricorredor/platformer/assets"
	"github.com/yuricorredor/platformer/clouds"
	"github.com/yuricorredor/platformer/entities"
	"github.com/yuricorredor/platformer/particle"
	"github.com/yuricorredor/platformer/tilemap"
)

var (
	gameClouds = &clouds.CloudsType{
		CloudImages: assets.Assets.Images["clouds"],
		Count:       16,
	}
	leafs = &particle.LeafsType{}
)

type Game struct {
	scollX       int
	scrollY      int
	screenWidth  int
	screenHeight int
}

func (g *Game) Update() error {
	g.updateScrollPosition()
	gameClouds.Update()
	entities.Player.Update()
	leafs.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(assets.Assets.Images["background"][0], nil)

	gameClouds.Draw(screen, g.scollX, g.scrollY)
	tilemap.TileMap.Draw(screen, g.scollX, g.scrollY)
	entities.Player.Draw(screen, g.scollX, g.scrollY)
	leafs.Draw(screen, g.scollX, g.scrollY)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.screenWidth = outsideWidth / 2
	g.screenHeight = outsideHeight / 2
	return g.screenWidth, g.screenHeight
}

func (g *Game) updateScrollPosition() {
	playerRect := entities.Player.Rect()
	g.scollX += (int(playerRect.CenterX()) - g.screenWidth/2 - g.scollX) / 15
	g.scrollY += (int(playerRect.CenterY()) - g.screenHeight/2 - g.scrollY) / 15
}

func (g *Game) loadMap(mapName string) {
	tilemap.TileMap.Load("assets/data/maps/" + mapName + ".json")
}

func main() {
	game := &Game{}
	game.loadMap("0")

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Platformer")

	gameClouds.GenerateRandomClouds()
	leafs = particle.NewLeafs()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}