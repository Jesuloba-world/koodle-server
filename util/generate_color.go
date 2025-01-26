package util

import (
	"math/rand/v2"

	"github.com/lucasb-eyer/go-colorful"
)

func GenerateColor() string {
	// Generate cool color using HSL
	// H: 120-300 (green to purple range)
	// S: 0.4-0.8 (40-80% saturation)
	// L: 0.4-0.7 (40-70% lightness)
	h := 120.0 + rand.Float64()*180.0 // Range 120-300
	s := 0.4 + rand.Float64()*0.4     // Range 0.4-0.8
	l := 0.4 + rand.Float64()*0.3     // Range 0.4-0.7

	c := colorful.Hsl(h, s, l)
	return c.Hex()
}
