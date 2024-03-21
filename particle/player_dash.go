package particle

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yuricorredor/platformer/animation"
	"github.com/yuricorredor/platformer/assets"
	"github.com/yuricorredor/platformer/types"
)

type DashParticlesType struct {
	Particles []*Particle
}

func (d *DashParticlesType) Update() {
	for _, particle := range d.Particles {
		kill := particle.Update()

		if kill {
			d.Particles = d.Particles[1:]
		}
	}
}

func (d *DashParticlesType) Draw(screen *ebiten.Image, scollX, scollY int) {
	for _, particle := range d.Particles {
		particle.Draw(screen, scollX, scollY)
	}
}

func CreateDashParticle(velocity types.Vector, position types.Vector) *Particle {
	return &Particle{
		Type:     "particle",
		Position: position,
		Velocity: velocity,
		Frame:    rand.Intn(7),
		Animation: animation.Animation{
			Images:        assets.Assets.Images["particle"].Image,
			ImageDuration: 6,
			Loop:          false,
			Done:          false,
		},
	}
}

var DashParticles = &DashParticlesType{}
