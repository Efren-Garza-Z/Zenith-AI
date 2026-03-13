package zig_simd

var (
	SobelX = [][]float64{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}

	Sharpen = [][]float64{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}
)
