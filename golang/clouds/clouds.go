package clouds

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yuricorredor/platformer/assets"
	"github.com/yuricorredor/platformer/types"
)

var Clouds = &CloudsType{
	CloudImages: assets.Assets.Images["clouds"],
	Count:       16,
}

type Cloud struct {
	Position types.Vector
	Image    *ebiten.Image
	Speed    float64
	Depth    float64
}

func (c *Cloud) Update() {
	c.Position.X += float64(c.Speed)
}

func (c *Cloud) Draw(screen *ebiten.Image, scrollX, scrollY int) {
	options := &ebiten.DrawImageOptions{}
	renderX := int(c.Position.X - float64(scrollX)*float64(c.Depth))
	renderY := int(c.Position.Y - float64(scrollY)*float64(c.Depth))
	screenWidth := screen.Bounds().Max.X
	screenHeight := screen.Bounds().Max.Y
	imageWidth := c.Image.Bounds().Max.X
	imageHeight := c.Image.Bounds().Max.Y
	options.GeoM.Translate(float64(renderX%(screenWidth+imageWidth)-imageWidth), float64(renderY%(screenHeight+imageHeight)-imageHeight))
	screen.DrawImage(c.Image, options)
}

type CloudsType struct {
	CloudImages []*ebiten.Image
	Clouds      []Cloud
	Count       int
}

func (c *CloudsType) Update() {
	for i := 0; i < c.Count; i++ {
		c.Clouds[i].Update()
	}
}

func (c *CloudsType) Draw(screen *ebiten.Image, scrollX, scrollY int) {
	for i := 0; i < c.Count; i++ {
		c.Clouds[i].Draw(screen, scrollX, scrollY)
	}
}

func (c *CloudsType) GenerateRandomClouds() {
	c.Clouds = randomClouds(c)
}

func randomClouds(Clouds *CloudsType) []Cloud {
	clouds := make([]Cloud, Clouds.Count)

	for i := 0; i < Clouds.Count; i++ {
		clouds[i] = Cloud{
			Position: types.Vector{
				X: rand.Float64() * 99999,
				Y: rand.Float64() * 99999,
			},
			Image: Clouds.CloudImages[rand.Intn(len(Clouds.CloudImages))],
			Speed: rand.Float64()*0.05 + 0.05,
			Depth: rand.Float64()*0.6 + 0.2,
		}
	}

	return clouds
}
