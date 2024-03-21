package assets

import (
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	BasePath            = "assets/data/images/"
	PlayerPath          = BasePath + "entities/player.png"
	PlayerIdlePath      = BasePath + "entities/player/idle/"
	PlayerJumpPath      = BasePath + "entities/player/jump/"
	PlayerRunPath       = BasePath + "entities/player/run/"
	PlayerSlidePath     = BasePath + "entities/player/slide/"
	PlayerWallSlidePath = BasePath + "entities/player/wall_slide/"
	DecorPath           = BasePath + "tiles/decor/"
	GrassPath           = BasePath + "tiles/grass/"
	LargeDecorPath      = BasePath + "tiles/large_decor/"
	SpawnersPath        = BasePath + "tiles/spawners/"
	StonePath           = BasePath + "tiles/stone/"
	BackgroundPath      = BasePath + "background.png"
	CloudsPath          = BasePath + "clouds/"
	LeafsPath           = BasePath + "particles/leaf/"
	ParticlePath        = BasePath + "particles/particle/"
)

var Assets = &AssetsType{
	Images: map[string]Asset{
		"player": {
			Image:                load_image(PlayerPath),
			ShouldRenderOnGame:   true,
			ShouldRenderOnEditor: true,
		},
		"player_idle": {
			Image:                load_images(PlayerIdlePath),
			ShouldRenderOnGame:   true,
			ShouldRenderOnEditor: true,
		},
		"player_jump": {
			Image:                load_images(PlayerJumpPath),
			ShouldRenderOnGame:   true,
			ShouldRenderOnEditor: true,
		},
		"player_run": {
			Image:                load_images(PlayerRunPath),
			ShouldRenderOnGame:   true,
			ShouldRenderOnEditor: true,
		},
		"player_slide": {
			Image:                load_images(PlayerSlidePath),
			ShouldRenderOnGame:   true,
			ShouldRenderOnEditor: true,
		},
		"player_wall_slide": {
			Image:                load_images(PlayerWallSlidePath),
			ShouldRenderOnGame:   true,
			ShouldRenderOnEditor: true,
		},
		"decor": {
			Image:                load_images(DecorPath),
			ShouldRenderOnGame:   true,
			ShouldRenderOnEditor: true,
		},
		"grass": {
			Image:                load_images(GrassPath),
			ShouldRenderOnGame:   true,
			ShouldRenderOnEditor: true,
		},
		"large_decor": {
			Image:                load_images(LargeDecorPath),
			ShouldRenderOnGame:   true,
			ShouldRenderOnEditor: true,
		},
		"spawners": {
			Image:                load_images(SpawnersPath),
			ShouldRenderOnGame:   false,
			ShouldRenderOnEditor: true,
		},
		"stone": {
			Image:                load_images(StonePath),
			ShouldRenderOnGame:   true,
			ShouldRenderOnEditor: true,
		},
		"background": {
			Image:                load_image(BackgroundPath),
			ShouldRenderOnGame:   true,
			ShouldRenderOnEditor: true,
		},
		"clouds": {
			Image:                load_images(CloudsPath),
			ShouldRenderOnGame:   true,
			ShouldRenderOnEditor: false,
		},
		"particle_leaf": {
			Image:                load_images(LeafsPath),
			ShouldRenderOnGame:   true,
			ShouldRenderOnEditor: false,
		},
		"particle": {
			Image:                load_images(ParticlePath),
			ShouldRenderOnGame:   true,
			ShouldRenderOnEditor: false,
		},
	},
}

type Asset struct {
	Image                []*ebiten.Image
	ShouldRenderOnGame   bool
	ShouldRenderOnEditor bool
}

type AssetsType struct {
	Images map[string]Asset
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
