package entities

import (
	"github.com/yuricorredor/platformer/assets"
	"github.com/yuricorredor/platformer/types"
)

var PlayerAnimations = map[string]*Animation{
	"idle": {
		Images:        assets.Assets.Images["player_idle"],
		ImageDuration: 6,
		Loop:          true,
		Done:          false,
		Offset: types.Vector{
			X: -3,
			Y: -3,
		},
	},
	"run": {
		Images:        assets.Assets.Images["player_run"],
		ImageDuration: 4,
		Loop:          true,
		Done:          false,
		Offset: types.Vector{
			X: -3,
			Y: -3,
		},
	},
	"jump": {
		Images:        assets.Assets.Images["player_jump"],
		ImageDuration: 5,
		Loop:          false,
		Done:          false,
		Offset: types.Vector{
			X: -3,
			Y: -2.5,
		},
	},
	"slide": {
		Images:        assets.Assets.Images["player_slide"],
		ImageDuration: 5,
		Loop:          false,
		Done:          false,
		Offset: types.Vector{
			X: -3,
			Y: -3,
		},
	},
	"wall_slide": {
		Images:        assets.Assets.Images["player_wall_slide"],
		ImageDuration: 5,
		Loop:          false,
		Done:          false,
		Offset: types.Vector{
			X: -3,
			Y: -3,
		},
	},
}

var Player = &PhysicsEntity{
	EntityType: "player",
	Position:   types.Vector{X: 0, Y: 0},
	Velocity:   types.Vector{X: 0, Y: 0},
	Collisions: types.Collisions{},
	Action:     "idle",
	Animations: PlayerAnimations,
}
