package particle

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yuricorredor/platformer/animation"
	"github.com/yuricorredor/platformer/assets"
	"github.com/yuricorredor/platformer/rects"
	"github.com/yuricorredor/platformer/tilemap"
	"github.com/yuricorredor/platformer/types"
)

type ProjectilesType struct {
	Particles []*Particle
}

func (projectile *ProjectilesType) Update(playerDashing float64, playerRect rects.Rect) {
	var remainingParticles []*Particle

	for _, particle := range projectile.Particles {

		particle.Frame += 1

		particle.Position.X += particle.Velocity.X
		particle.Position.Y += particle.Velocity.Y

		shouldRemoveParticle := particle.Frame > 360 ||
			tilemap.TileMap.CheckForSolid(particle.Position) ||
			(math.Abs(playerDashing) < 50 && playerRect.Colliderect(particle.Rect()))

		if !shouldRemoveParticle {
			remainingParticles = append(remainingParticles, particle)
		}
	}

	projectile.Particles = remainingParticles
}

func (projectile *ProjectilesType) Draw(screen *ebiten.Image, scollX, scollY int) {
	for _, particle := range projectile.Particles {
		particle.Draw(screen, scollX, scollY)
	}
}

func NewProjectile(position, velocity types.Vector) *Particle {
	return &Particle{
		Type:     "projectile",
		Position: position,
		Velocity: velocity,
		Frame:    0,
		Animation: animation.Animation{
			Images:        assets.Assets.Images["projectile"].Image,
			ImageDuration: 600,
			Loop:          false,
			Done:          false,
		},
	}
}

var Projectiles = &ProjectilesType{}
