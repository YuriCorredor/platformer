package types

type Vector struct {
	X float64
	Y float64
}

type Collisions struct {
	Top    bool
	Bottom bool
	Left   bool
	Right  bool
}
