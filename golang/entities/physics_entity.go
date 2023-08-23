package entities

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yuricorredor/platformer/assets"
	"github.com/yuricorredor/platformer/rects"
	"github.com/yuricorredor/platformer/tilemap"
	"github.com/yuricorredor/platformer/types"
)

var Player = &PhysicsEntity{
	EntityType: "player",
	Position:   types.Vector{X: 0, Y: 0},
	Velocity:   types.Vector{X: 0, Y: 0},
	Collisions: types.Collisions{},
}

type PhysicsEntity struct {
	EntityType string
	Position   types.Vector
	Velocity   types.Vector
	Collisions types.Collisions
}

func (p *PhysicsEntity) Update() error {
	if p.EntityType == "player" {
		p.HandlePlayerMovement()
	}

	return nil
}

func (p *PhysicsEntity) Draw(screen *ebiten.Image, scrollX, scrollY int) {
	if p.EntityType == "player" {
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(p.Position.X-float64(scrollX), p.Position.Y-float64(scrollY))
		screen.DrawImage(assets.Assets.Images[p.EntityType][0], options)
	}
}

func (p *PhysicsEntity) Jump() {
	p.Velocity.Y = -3
}

func (p *PhysicsEntity) ResetCollisions() {
	p.Collisions.Top = false
	p.Collisions.Bottom = false
	p.Collisions.Left = false
	p.Collisions.Right = false
}

func (p *PhysicsEntity) Size() (int, int) {
	bounds := assets.Assets.Images[p.EntityType][0].Bounds()
	return bounds.Max.X, bounds.Max.Y
}

func (p *PhysicsEntity) Rect() rects.Rect {
	width, height := p.Size()
	return rects.Rect{X: p.Position.X, Y: p.Position.Y, Width: float64(width), Height: float64(height)}
}

func (p *PhysicsEntity) HandlePlayerMovement() {
	p.ResetCollisions()
	var movement = types.Vector{X: 0, Y: 0}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		movement.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		movement.X += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.Jump()
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

	p.Velocity.Y = math.Min(float64(3), float64(p.Velocity.Y)+0.1)

	if p.Collisions.Bottom || p.Collisions.Top {
		p.Velocity.Y = 0
	}
}
