package zig_simd

import "github.com/efren-garza-z/Zenith-AI/core"

// CenterOfMass representa el punto medio de los píxeles detectados
type CenterOfMass struct {
	X float64
	Y float64
}

// GetActivePixelPercentage calcula la densidad de contenido
func GetActivePixelPercentage(t *core.Tensor4D) float64 {
	if len(t.Data) == 0 {
		return 0
	}
	activeCount := 0
	for _, v := range t.Data {
		if v > 0.5 {
			activeCount++
		}
	}
	return (float64(activeCount) / float64(len(t.Data))) * 100
}

// GetCenterOfMass calcula las coordenadas (X, Y) promedio de los píxeles activos
func GetCenterOfMass(t *core.Tensor4D) CenterOfMass {
	sumX, sumY := 0.0, 0.0
	count := 0

	for y := 0; y < t.Height; y++ {
		for x := 0; x < t.Width; x++ {
			// Usamos At para obtener el valor del píxel en la posición actual
			if t.At(0, 0, y, x) > 0.5 {
				sumX += float64(x)
				sumY += float64(y)
				count++
			}
		}
	}

	if count == 0 {
		return CenterOfMass{X: -1, Y: -1} // Indica que no hay contenido
	}

	return CenterOfMass{
		X: sumX / float64(count),
		Y: sumY / float64(count),
	}
}
