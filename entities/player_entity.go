package entities

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yuricorredor/platformer/animation"
	"github.com/yuricorredor/platformer/assets"
	"github.com/yuricorredor/platformer/particle"
	"github.com/yuricorredor/platformer/rects"
	"github.com/yuricorredor/platformer/tilemap"
	"github.com/yuricorredor/platformer/types"
)

type PlayerEntity struct {
	EntityType string
	Position   types.Vector
	Velocity   types.Vector
	Collisions types.Collisions
	Action     string
	Animations map[string]*animation.Animation
	Fliped     bool
	AirTime    int
	Jumps      int
	WallSlide  bool
	Dashing    float64
}

func (p *PlayerEntity) Draw(screen *ebiten.Image, scrollX, scrollY int) {
	if math.Abs(p.Dashing) > 50 {
		return
	}
	image := p.Animations[p.Action].Image()
	imageOffset := p.Animations[p.Action].Offset
	options := &ebiten.DrawImageOptions{}
	if p.Fliped {
		options.GeoM.Scale(-1, 1)
		options.GeoM.Translate(float64(image.Bounds().Max.X), 0)
	}
	options.GeoM.Translate(p.Position.X-float64(scrollX)+imageOffset.X, p.Position.Y-float64(scrollY)+imageOffset.Y)
	screen.DrawImage(image, options)
}

func (p *PlayerEntity) SetAction(action string) {
	p.Action = action
}

func (p *PlayerEntity) Jump() bool {
	var jumped = false

	if p.WallSlide {
		if p.Fliped && ebiten.IsKeyPressed(ebiten.KeyA) {
			p.Velocity.X = 3.5
			p.Velocity.Y = -2.5
			p.AirTime = 5
			p.Jumps = int(math.Max(0, float64(p.Jumps-1)))
			jumped = true
		} else if !p.Fliped && ebiten.IsKeyPressed(ebiten.KeyD) {
			p.Velocity.X = -3.5
			p.Velocity.Y = -2.5
			p.AirTime = 5
			p.Jumps = int(math.Max(0, float64(p.Jumps-1)))
			jumped = true
		}
	} else if p.Jumps != 0 {
		p.Velocity.Y = -3
		p.Jumps--
		p.AirTime = 5
		jumped = true
	}

	return jumped
}

func (p *PlayerEntity) Dash() {
	if p.Dashing == 0 {
		if p.Fliped {
			p.Dashing = -60
		} else {
			p.Dashing = 60
		}
	}
}

func (p *PlayerEntity) ResetCollisions() {
	p.Collisions.Top = false
	p.Collisions.Bottom = false
	p.Collisions.Left = false
	p.Collisions.Right = false
}

func (p *PlayerEntity) Size() (int, int) {
	bounds := assets.Assets.Images[p.EntityType].Image[0].Bounds()
	return bounds.Max.X, bounds.Max.Y
}

func (p *PlayerEntity) Rect() rects.Rect {
	width, height := p.Size()
	return rects.Rect{X: p.Position.X, Y: p.Position.Y, Width: float64(width), Height: float64(height)}
}

func (p *PlayerEntity) Update() error {
	p.Animations[p.Action].Update()

	p.ResetCollisions()
	var movement = types.Vector{X: 0, Y: 0}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		movement.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		movement.X += 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.Jump()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyShift) {
		p.Dash()
	}

	frameMovement := types.Vector{X: movement.X + p.Velocity.X, Y: movement.Y + p.Velocity.Y}

	p.Position.X += frameMovement.X
	rectsList := tilemap.TileMap.PhysicsRectsAroundPosition(p.Position)
	entityRect := p.Rect()
	for _, rect := range rectsList {
		if entityRect.Coliderect(rect) {
			if frameMovement.X > 0 {
				entityRect.SetRight(rect.Left())
				p.Collisions.Right = true
			}
			if frameMovement.X < 0 {
				entityRect.SetLeft(rect.Right())
				p.Collisions.Left = true
			}
			p.Position.X = entityRect.X
		}
	}

	p.Position.Y += frameMovement.Y
	rectsList = tilemap.TileMap.PhysicsRectsAroundPosition(p.Position)
	entityRect = p.Rect()
	for _, rect := range rectsList {
		if entityRect.Coliderect(rect) {
			if frameMovement.Y > 0 {
				entityRect.SetBottom(rect.Top())
				p.Collisions.Bottom = true
			}
			if frameMovement.Y < 0 {
				entityRect.SetTop(rect.Bottom())
				p.Collisions.Top = true
			}
			p.Position.Y = entityRect.Y
		}
	}

	if p.Collisions.Bottom {
		p.Jumps = 1
		p.AirTime = 0
	} else {
		p.AirTime++
	}

	if (p.Collisions.Left || p.Collisions.Right) && p.AirTime > 4 {
		p.WallSlide = true
		p.Velocity.Y = math.Min(float64(p.Velocity.Y), 0.5)
		if p.Collisions.Right {
			p.Fliped = false
		} else {
			p.Fliped = true
		}
	} else {
		p.WallSlide = false
	}

	if !p.WallSlide {
		if p.AirTime > 4 {
			p.SetAction("jump")
		} else if movement.X != 0 {
			p.SetAction("run")
		} else {
			p.SetAction("idle")
		}
	} else {
		p.SetAction("wall_slide")
	}

	rect := p.Rect()
	position := types.Vector{
		X: rect.CenterX(),
		Y: rect.CenterY(),
	}

	if math.Abs(p.Dashing) == 60 || math.Abs(p.Dashing) == 50 {
		for i := 0; i < 20; i++ {
			angle := rand.Float64() * math.Pi * 2
			speed := rand.Float64()*0.5 + 0.5
			velocity := types.Vector{
				X: math.Cos(angle) * speed,
				Y: math.Sin(angle) * speed,
			}
			particle.DashParticles.Particles =
				append(particle.DashParticles.Particles, particle.CreateDashParticle(velocity, position))
		}
	}

	if math.Abs(p.Dashing) > 50 {
		p.Velocity.X = math.Abs(p.Dashing) / p.Dashing * 8
		if math.Abs(p.Dashing) == 51 {
			p.Velocity.X *= 0.1
		}

		velocity := types.Vector{
			X: math.Abs(p.Dashing) / p.Dashing * rand.Float64() * 3,
			Y: 0,
		}
		particle.DashParticles.Particles =
			append(particle.DashParticles.Particles, particle.CreateDashParticle(velocity, position))
	}

	if p.Dashing > 0 {
		p.Dashing = math.Max(0, p.Dashing-1)
	}
	if p.Dashing < 0 {
		p.Dashing = math.Min(0, p.Dashing+1)
	}

	if movement.X > 0 {
		p.Fliped = false
	}
	if movement.X < 0 {
		p.Fliped = true
	}

	if p.Velocity.X > 0 {
		p.Velocity.X = math.Max(p.Velocity.X-0.1, 0)
	} else if p.Velocity.X < 0 {
		p.Velocity.X = math.Min(p.Velocity.X+0.1, 0)
	}

	p.Velocity.Y = math.Min(3, p.Velocity.Y+0.1)

	if p.Collisions.Bottom || p.Collisions.Top {
		p.Velocity.Y = 0
	}

	return nil
}

var PlayerAnimations = map[string]*animation.Animation{
	"idle": {
		Images:        assets.Assets.Images["player_idle"].Image,
		ImageDuration: 6,
		Loop:          true,
		Done:          false,
		Offset: types.Vector{
			X: -3,
			Y: -3,
		},
	},
	"run": {
		Images:        assets.Assets.Images["player_run"].Image,
		ImageDuration: 4,
		Loop:          true,
		Done:          false,
		Offset: types.Vector{
			X: -3,
			Y: -3,
		},
	},
	"jump": {
		Images:        assets.Assets.Images["player_jump"].Image,
		ImageDuration: 5,
		Loop:          false,
		Done:          false,
		Offset: types.Vector{
			X: -3,
			Y: -2.5,
		},
	},
	"slide": {
		Images:        assets.Assets.Images["player_slide"].Image,
		ImageDuration: 5,
		Loop:          false,
		Done:          false,
		Offset: types.Vector{
			X: -3,
			Y: -3,
		},
	},
	"wall_slide": {
		Images:        assets.Assets.Images["player_wall_slide"].Image,
		ImageDuration: 5,
		Loop:          false,
		Done:          false,
		Offset: types.Vector{
			X: -3,
			Y: -3,
		},
	},
}

var Player = &PlayerEntity{
	EntityType: "player",
	Position:   types.Vector{X: 0, Y: 0},
	Velocity:   types.Vector{X: 0, Y: 0},
	Collisions: types.Collisions{},
	Action:     "idle",
	Animations: PlayerAnimations,
	Jumps:      1,
}
