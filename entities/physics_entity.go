package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yuricorredor/platformer/rects"
)

type PhysicsEntity interface {
	Update() error
	Draw(screen *ebiten.Image, scrollX, scrollY int)
	SetAction(action string)
	Size() (int, int)
	Rect() rects.Rect
}
