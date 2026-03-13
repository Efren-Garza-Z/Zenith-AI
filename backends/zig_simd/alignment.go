package zig_simd

import (
	"github.com/efren-garza-z/zenith-ai/core"
)

// CenterSignature desplaza los píxeles para que el Centro de Masa esté en el centro exacto del tensor
func CenterSignature(t *core.Tensor4D) *core.Tensor4D {
	com := GetCenterOfMass(t)

	// Si no hay píxeles (imagen vacía), devolvemos el original
	if com.X == -1 {
		return t
	}

	// Calculamos el centro teórico del lienzo
	targetX := float64(t.Width) / 2
	targetY := float64(t.Height) / 2

	// Calculamos cuánto debemos mover cada píxel (Offset)
	offsetX := int(targetX - com.X)
	offsetY := int(targetY - com.Y)

	// Creamos un nuevo tensor limpio para la firma centrada
	newTensor := core.NewTensor4D(t.Batch, t.Channels, t.Height, t.Width)

	// Recorremos el tensor original y movemos los píxeles al nuevo
	for y := 0; y < t.Height; y++ {
		for x := 0; x < t.Width; x++ {
			// Si el píxel está activo (tinta)
			if t.At(0, 0, y, x) > 0.5 {
				newX := x + offsetX
				newY := y + offsetY

				// Verificamos que el nuevo punto esté dentro de los límites
				if newX >= 0 && newX < t.Width && newY >= 0 && newY < t.Height {
					// Seteamos el píxel en los 3 canales (RGB) para mantenerlo blanco
					newTensor.Set(0, 0, newY, newX, 1.0)
					newTensor.Set(0, 1, newY, newX, 1.0)
					newTensor.Set(0, 2, newY, newX, 1.0)
				}
			}
		}
	}

	return newTensor
}
