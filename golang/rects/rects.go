package rects

type Rect struct {
	X, Y, Width, Height float64
}

func (r *Rect) Coliderect(other Rect) bool {
	return r.Left() < other.Right() && r.Right() > other.Left() && r.Top() < other.Bottom() && r.Bottom() > other.Top()
}

func (r *Rect) Top() float64 {
	return r.Y
}

func (r *Rect) Bottom() float64 {
	return r.Y + r.Height
}

func (r *Rect) Left() float64 {
	return r.X
}

func (r *Rect) Right() float64 {
	return r.X + r.Width
}

func (r *Rect) SetLeft(left float64) {
	r.X = left
}

func (r *Rect) SetRight(right float64) {
	r.X = right - r.Width
}

func (r *Rect) SetTop(top float64) {
	r.Y = top
}

func (r *Rect) SetBottom(bottom float64) {
	r.Y = bottom - r.Height
}

func (r *Rect) CenterX() float64 {
	return r.X + r.Width/2
}

func (r *Rect) CenterY() float64 {
	return r.Y + r.Height/2
}
