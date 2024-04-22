package particle

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yuricorredor/platformer/animation"
	"github.com/yuricorredor/platformer/rects"
	"github.com/yuricorredor/platformer/types"
)

type ParticleI interface {
	Update()
	Draw(screen *ebiten.Image, scollX, scollY int)
}

type Particle struct {
	Type      string
	Angle     float64
	Position  types.Vector
	Velocity  types.Vector
	Frame     int
	Animation animation.Animation
}

func (p *Particle) Size() (int, int) {
	bounds := p.Animation.Images[0].Bounds()
	return bounds.Max.X, bounds.Max.Y
}

func (p *Particle) Rect() rects.Rect {
	width, height := p.Size()
	return rects.Rect{X: p.Position.X, Y: p.Position.Y, Width: float64(width), Height: float64(height)}
}

func (p *Particle) Update() bool {
	kill := false
	if p.Animation.Done {
		kill = true
	}

	p.Position.X += p.Velocity.X
	p.Position.Y += p.Velocity.Y
	p.Animation.Update()

	return kill
}

func (p *Particle) Draw(screen *ebiten.Image, scollX, scollY int) {
	image := p.Animation.Image()
	options := &ebiten.DrawImageOptions{}
	positionX := p.Position.X - float64(scollX) - float64(image.Bounds().Max.X/2)
	positionY := p.Position.Y - float64(scollY) - float64(image.Bounds().Max.Y/2)
	options.GeoM.Translate(positionX, positionY)
	screen.DrawImage(image, options)
}
