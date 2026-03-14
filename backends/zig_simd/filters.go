package zig_simd

import "github.com/efren-garza-z/zenith-ai/core"

// Dilate engrosa los píxeles activos. Esto da "tolerancia" al movimiento.
func Dilate(t *core.Tensor4D) *core.Tensor4D {
	newTensor := core.NewTensor4D(t.Batch, t.Channels, t.Height, t.Width)

	for y := 1; y < t.Height-1; y++ {
		for x := 1; x < t.Width-1; x++ {
			// Si el píxel actual o alguno de sus vecinos está activo...
			if t.At(0, 0, y, x) > 0.5 ||
				t.At(0, 0, y-1, x) > 0.5 || t.At(0, 0, y+1, x) > 0.5 ||
				t.At(0, 0, y, x-1) > 0.5 || t.At(0, 0, y, x+1) > 0.5 {
				newTensor.Set(0, 0, y, x, 1.0)
				newTensor.Set(0, 1, y, x, 1.0)
				newTensor.Set(0, 2, y, x, 1.0)
			}
		}
	}
	return newTensor
}
