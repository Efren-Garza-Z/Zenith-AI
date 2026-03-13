package zig_simd

import "Zenith-AI/core"

// GetActivePixelPercentage calcula qué porcentaje de la imagen es "tinta" (píxeles en 1)
// Funciona mejor después de aplicar Threshold.
func GetActivePixelPercentage(t *core.Tensor4D) float64 {
	if len(t.Data) == 0 {
		return 0
	}

	activeCount := 0
	for _, v := range t.Data {
		// Como está binarizado, los píxeles activos son 1.0
		if v > 0.5 {
			activeCount++
		}
	}

	// Retorna el porcentaje (0.0 a 100.0)
	return (float64(activeCount) / float64(len(t.Data))) * 100
}
