package entities

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yuricorredor/platformer/types"
)

type Animation struct {
	Images        []*ebiten.Image
	ImageDuration int
	Loop          bool
	Done          bool
	Frame         int
	Offset        types.Vector
}

func (a *Animation) Update() {
	if a.Loop {
		a.Frame = (a.Frame + 1) % (len(a.Images) * a.ImageDuration)
	} else {
		a.Frame = int(math.Min(float64(a.Frame+1), float64((len(a.Images)-1)*a.ImageDuration)))
		if a.Frame >= a.ImageDuration*(len(a.Images)-1) {
			a.Done = true
		}
	}
}

func (a *Animation) Image() *ebiten.Image {
	return a.Images[a.Frame/a.ImageDuration]
}
