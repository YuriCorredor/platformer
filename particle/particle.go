package particle

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yuricorredor/platformer/assets"
	"github.com/yuricorredor/platformer/entities"
	"github.com/yuricorredor/platformer/rects"
	"github.com/yuricorredor/platformer/tilemap"
	"github.com/yuricorredor/platformer/types"
)

type Particle struct {
	Type      string
	Position  types.Vector
	Velocity  types.Vector
	Frame     int
	Animation entities.Animation
}

type LeafsType struct {
	Particles []*Particle
	Spawners  []rects.Rect
}

func (l *LeafsType) Update() {
	for _, rect := range l.Spawners {
		if rand.Float64()*40000 < rect.Width*rect.Height {
			position := types.Vector{
				X: float64(rect.X) + rand.Float64()*float64(rect.Width),
				Y: float64(rect.Y) + rand.Float64()*float64(rect.Height),
			}

			l.Particles = append(l.Particles, newLeaf(position))
		}
	}
}

func (l *LeafsType) Draw(screen *ebiten.Image, scollX, scollY int) {
	for _, particle := range l.Particles {
		kill := particle.Update()
		particle.Draw(screen, scollX, scollY)
		particle.Position.X += math.Sin(float64(particle.Animation.Frame)*0.035) * 0.3
		if kill {
			l.Particles = l.Particles[1:]
		}
	}
}

func NewLeafs() *LeafsType {
	leafs := &LeafsType{
		Particles: []*Particle{},
		Spawners:  []rects.Rect{},
	}

	trees := tilemap.TileMap.Extract(types.Pair{
		AssetType:    "large_decor",
		AssetVariant: 2,
	}, true)

	for _, tree := range trees {
		treeRect := rects.Rect{
			X:      tree.Position.X*float64(tilemap.TileMap.TileSize) + 4,
			Y:      tree.Position.Y*float64(tilemap.TileMap.TileSize) + 4,
			Width:  24,
			Height: 12,
		}

		leafs.Spawners = append(leafs.Spawners, treeRect)
	}

	return leafs
}

func newLeaf(position types.Vector) *Particle {
	return &Particle{
		Type:     "leaf",
		Position: position,
		Velocity: types.Vector{
			X: -0.1,
			Y: 0.3,
		},
		Frame: rand.Intn(len(assets.Assets.Images["particle_leaf"]) - 1),
		Animation: entities.Animation{
			Images:        assets.Assets.Images["particle_leaf"],
			ImageDuration: 20,
			Loop:          false,
			Done:          false,
		},
	}
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
	positionX := p.Position.X - float64(scollX) - float64(image.Bounds().Max.X)
	positionY := p.Position.Y - float64(scollY) - float64(image.Bounds().Max.Y)
	options.GeoM.Translate(positionX, positionY)
	screen.DrawImage(image, options)
}
