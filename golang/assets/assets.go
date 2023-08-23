package assets

import (
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	BasePath       = "assets/data/images/"
	PlayerPath     = BasePath + "entities/player.png"
	DecorPath      = BasePath + "tiles/decor/"
	GrassPath      = BasePath + "tiles/grass/"
	LargeDecorPath = BasePath + "tiles/large_decor/"
	SpawnersPath   = BasePath + "tiles/spawners/"
	StonePath      = BasePath + "tiles/stone/"
	backgroundPath = BasePath + "background.png"
	CloudsPath     = BasePath + "clouds/"
)

var Assets = &AssetsType{
	Images: map[string][]*ebiten.Image{
		"player":      load_image(PlayerPath),
		"decor":       load_images(DecorPath),
		"grass":       load_images(GrassPath),
		"large_decor": load_images(LargeDecorPath),
		"spawners":    load_images(SpawnersPath),
		"stone":       load_images(StonePath),
		"background":  load_image(backgroundPath),
		"clouds":      load_images(CloudsPath),
	},
}

type AssetsType struct {
	Images map[string][]*ebiten.Image
}

func load_image(path string) []*ebiten.Image {
	image, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		panic(err)
	}

	return []*ebiten.Image{image}
}

func load_images(path string) []*ebiten.Image {
	images := []*ebiten.Image{}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			image := load_image(path)
			images = append(images, image...)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return images
}
