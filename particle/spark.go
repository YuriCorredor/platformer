package particle

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yuricorredor/platformer/animation"
	"github.com/yuricorredor/platformer/assets"
	"github.com/yuricorredor/platformer/types"
)

type SparksType struct {
	Particles []*Particle
}

// for spark in self.sparks.copy():
// kill = spark.update()
// spark.render(self.display, offset=render_scroll)
// if kill:
// 		self.sparks.remove(spark)

// class Spark:
//     def __init__(self, pos, angle, speed):
//         self.pos = list(pos)
//         self.angle = angle
//         self.speed = speed

//     def update(self):
//         self.pos[0] += math.cos(self.angle) * self.speed
//         self.pos[1] += math.sin(self.angle) * self.speed

//         self.speed = max(0, self.speed - 0.1)
//         return not self.speed

//     def render(self, surf, offset=(0, 0)):
//         render_points = [
//             (self.pos[0] + math.cos(self.angle) * self.speed * 3 - offset[0], self.pos[1] + math.sin(self.angle) * self.speed * 3 - offset[1]),
//             (self.pos[0] + math.cos(self.angle + math.pi * 0.5) * self.speed * 0.5 - offset[0], self.pos[1] + math.sin(self.angle + math.pi * 0.5) * self.speed * 0.5 - offset[1]),
//             (self.pos[0] + math.cos(self.angle + math.pi) * self.speed * 3 - offset[0], self.pos[1] + math.sin(self.angle + math.pi) * self.speed * 3 - offset[1]),
//             (self.pos[0] + math.cos(self.angle - math.pi * 0.5) * self.speed * 0.5 - offset[0], self.pos[1] + math.sin(self.angle - math.pi * 0.5) * self.speed * 0.5 - offset[1]),
//         ]

//         pygame.draw.polygon(surf, (255, 255, 255), render_points)

func (spark *SparksType) Update() {
	remainingSparks := []*Particle{}

	for _, particle := range spark.Particles {
		particle.Position.X += math.Cos(particle.Angle) * particle.Velocity.X
		particle.Position.Y += math.Sin(particle.Angle) * particle.Velocity.Y

		particle.Velocity.X = math.Max(0, particle.Velocity.X-0.1)
		particle.Velocity.Y = math.Max(0, particle.Velocity.Y-0.1)

		if particle.Velocity.X == 0 && particle.Velocity.Y == 0 {
			continue
		}

		remainingSparks = append(remainingSparks, particle)
	}

	spark.Particles = remainingSparks
}

func (spark *SparksType) Draw(screen *ebiten.Image, scollX, scollY int) {
	whiteImage := ebiten.NewImage(3, 3)
	op := &ebiten.DrawTrianglesOptions{}
	whiteImage.Fill(color.White)
	whiteSubImage := whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

	for _, particle := range spark.Particles {
		screen.DrawTriangles([]ebiten.Vertex{
			{
				DstX: float32(particle.Position.X + math.Cos(particle.Angle)*particle.Velocity.X*3 - float64(scollX)),
				DstY: float32(particle.Position.Y + math.Sin(particle.Angle)*particle.Velocity.Y*3 - float64(scollY)),
			},
			{
				DstX: float32(particle.Position.X + math.Cos(particle.Angle+math.Pi*0.5)*particle.Velocity.X*0.5 - float64(scollX)),
				DstY: float32(particle.Position.Y + math.Sin(particle.Angle+math.Pi*0.5)*particle.Velocity.Y*0.5 - float64(scollY)),
			},
			{
				DstX: float32(particle.Position.X + math.Cos(particle.Angle+math.Pi)*particle.Velocity.X*3 - float64(scollX)),
				DstY: float32(particle.Position.Y + math.Sin(particle.Angle+math.Pi)*particle.Velocity.Y*3 - float64(scollY)),
			},
			{
				DstX: float32(particle.Position.X + math.Cos(particle.Angle-math.Pi*0.5)*particle.Velocity.X*0.5 - float64(scollX)),
				DstY: float32(particle.Position.Y + math.Sin(particle.Angle-math.Pi*0.5)*particle.Velocity.Y*0.5 - float64(scollY)),
			},
		}, []uint16{0, 1, 2, 0, 2, 3}, whiteSubImage, op)
	}
}

func NewSpark(angle float64, position, velocity types.Vector) *Particle {
	return &Particle{
		Type:     "spark",
		Angle:    angle,
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

var SparksParticles = &SparksType{}
