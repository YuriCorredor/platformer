package entities

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yuricorredor/platformer/animation"
	"github.com/yuricorredor/platformer/assets"
	"github.com/yuricorredor/platformer/particle"
	"github.com/yuricorredor/platformer/rects"
	"github.com/yuricorredor/platformer/tilemap"
	"github.com/yuricorredor/platformer/types"
)

type EnemyEntity struct {
	EntityType string
	Position   types.Vector
	Velocity   types.Vector
	Walking    int
	Flipped    bool
	Collisions types.Collisions
	Action     string
	Animations map[string]*animation.Animation
}

func (enemy *EnemyEntity) Size() (int, int) {
	bounds := assets.Assets.Images[enemy.EntityType].Image[0].Bounds()
	return bounds.Max.X, bounds.Max.Y
}

func (enemy *EnemyEntity) Rect() rects.Rect {
	width, height := enemy.Size()
	return rects.Rect{X: enemy.Position.X, Y: enemy.Position.Y, Width: float64(width), Height: float64(height)}
}

func (enemy *EnemyEntity) SetAction(action string) {
	enemy.Action = action
}

func (enemy *EnemyEntity) Draw(screen *ebiten.Image, scrollX, scrollY int) {
	image := enemy.Animations[enemy.Action].Image()
	imageOffset := enemy.Animations[enemy.Action].Offset
	options := &ebiten.DrawImageOptions{}
	if enemy.Flipped {
		options.GeoM.Scale(-1, 1)
		options.GeoM.Translate(float64(image.Bounds().Max.X), 0)
	}
	options.GeoM.Translate(enemy.Position.X-float64(scrollX)+imageOffset.X, enemy.Position.Y-float64(scrollY)+imageOffset.Y)
	screen.DrawImage(image, options)

	options = &ebiten.DrawImageOptions{}
	gunImage := assets.Assets.Images["gun"].Image[0]
	if enemy.Flipped {
		options.GeoM.Scale(-1, 1)
		options.GeoM.Translate(enemy.Position.X-float64(scrollX)+imageOffset.X-float64(gunImage.Bounds().Max.X)+8, enemy.Position.Y-float64(scrollY)+imageOffset.Y+8)
	} else {
		options.GeoM.Translate(enemy.Position.X-float64(scrollX)+imageOffset.X+12, enemy.Position.Y-float64(scrollY)+imageOffset.Y+8)
	}

	screen.DrawImage(gunImage, options)
}

func (enemy *EnemyEntity) ResetCollisions() {
	enemy.Collisions.Top = false
	enemy.Collisions.Bottom = false
	enemy.Collisions.Left = false
	enemy.Collisions.Right = false
}

func (enemy *EnemyEntity) Update() error {
	enemy.Animations[enemy.Action].Update()

	var movement = types.Vector{X: 0, Y: 0}

	if enemy.Walking != 0 {
		enemyRect := enemy.Rect()
		positionToCheck := types.Vector{X: enemyRect.CenterX(), Y: enemy.Position.Y + 20}
		if enemy.Flipped {
			positionToCheck.X -= 7
		} else {
			positionToCheck.X += 7
		}

		if tilemap.TileMap.CheckForSolid(positionToCheck) {
			if enemy.Collisions.Right || enemy.Collisions.Left {
				enemy.Flipped = !enemy.Flipped
			} else {
				if enemy.Flipped {
					movement.X -= 0.5
				} else {
					movement.X = 0.5
				}
			}
		} else {
			enemy.Flipped = !enemy.Flipped
		}

		if enemy.Flipped {
			movement.X -= 0.5
		} else {
			movement.X = 0.5
		}

		enemy.Walking = int(math.Max(0, float64(enemy.Walking-1)))

		if enemy.Walking == 0 {
			distanceEnemyPlayer := &types.Vector{
				X: Player.Position.X - enemy.Position.X,
				Y: Player.Position.Y - enemy.Position.Y,
			}

			if math.Abs(distanceEnemyPlayer.Y) < 16 {
				projectilePosition := types.Vector{
					X: enemy.Position.X,
					Y: enemy.Position.Y + 8,
				}
				projectileVelocity := types.Vector{
					X: 0,
					Y: 0,
				}
				if enemy.Flipped && distanceEnemyPlayer.X < 0 {
					projectileVelocity.X = -1.5
					particle.Projectiles.Particles = append(particle.Projectiles.Particles, particle.NewProjectile(projectilePosition, projectileVelocity))
					for i := 0; i < 4; i++ {
						angle := rand.Float64() * math.Pi * 2
						particle.SparksParticles.Particles = append(particle.SparksParticles.Particles, particle.NewSpark(angle, projectilePosition, types.Vector{X: 1, Y: 1}))
					}
				} else if !enemy.Flipped && distanceEnemyPlayer.X > 0 {
					projectileVelocity.X = 1.5
					particle.Projectiles.Particles = append(particle.Projectiles.Particles, particle.NewProjectile(projectilePosition, projectileVelocity))
					for i := 0; i < 4; i++ {
						angle := rand.Float64() * math.Pi
						particle.SparksParticles.Particles = append(particle.SparksParticles.Particles, particle.NewSpark(angle, projectilePosition, types.Vector{X: 1, Y: 1}))
					}
				}
			}
		}

	} else if rand.Intn(100) == 1 {
		enemy.Walking = rand.Intn(120)
	}

	enemy.ResetCollisions()

	frameMovement := types.Vector{X: movement.X + enemy.Velocity.X, Y: movement.Y + enemy.Velocity.Y}

	enemy.Position.X += frameMovement.X
	entityRect := enemy.Rect()
	rectsList := tilemap.TileMap.PhysicsRectsAroundPosition(enemy.Position)

	for _, rect := range rectsList {
		if entityRect.Colliderect(rect) {
			if frameMovement.X > 0 {
				entityRect.SetRight(rect.Left())
				enemy.Collisions.Right = true
			}
			if frameMovement.X < 0 {
				entityRect.SetLeft(rect.Right())
				enemy.Collisions.Left = true
			}
			enemy.Position.X = entityRect.X
		}
	}

	enemy.Position.Y += frameMovement.Y
	entityRect = enemy.Rect()
	rectsList = tilemap.TileMap.PhysicsRectsAroundPosition(enemy.Position)
	for _, rect := range rectsList {
		if entityRect.Colliderect(rect) {
			if frameMovement.Y > 0 {
				entityRect.SetBottom(rect.Top())
				enemy.Collisions.Bottom = true
			}
			if frameMovement.Y < 0 {
				entityRect.SetTop(rect.Bottom())
				enemy.Collisions.Top = true
			}
			enemy.Position.Y = entityRect.Y
		}
	}

	if movement.X > 0 {
		enemy.Flipped = false
	}
	if movement.X < 0 {
		enemy.Flipped = true
	}

	enemy.Velocity.Y = math.Min(3, enemy.Velocity.Y+0.1)

	if enemy.Collisions.Bottom || enemy.Collisions.Top {
		enemy.Velocity.Y = 0
	}

	if movement.X != 0 {
		enemy.SetAction("run")
	} else {
		enemy.SetAction("idle")
	}

	return nil
}

func CreateEnemy(position types.Vector) *EnemyEntity {
	return &EnemyEntity{
		EntityType: "enemy",
		Position:   position,
		Velocity:   types.Vector{X: 0, Y: 0},
		Walking:    0,
		Flipped:    false,
		Collisions: types.Collisions{},
		Action:     "idle",
		Animations: EnemyAnimations,
	}
}

var EnemyAnimations = map[string]*animation.Animation{
	"idle": {
		Images:        assets.Assets.Images["enemy_idle"].Image,
		ImageDuration: 8,
		Loop:          true,
		Done:          false,
		Offset: types.Vector{
			X: -3,
			Y: -3,
		},
	},
	"run": {
		Images:        assets.Assets.Images["enemy_run"].Image,
		ImageDuration: 4,
		Loop:          true,
		Done:          false,
		Offset: types.Vector{
			X: -3,
			Y: -3,
		},
	},
}
