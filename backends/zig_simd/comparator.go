package zig_simd

import (
	"github.com/efren-garza-z/zenith-ai/core"
)

// CompareSignatures realiza una comparación de solapamiento de píxeles (Índice de Jaccard simplificado).
// Requiere que ambos tensores estén binarizados (0 o 1) y centrados.
func CompareSignatures(t1, t2 *core.Tensor4D) float64 {
	// 1. Verificación de dimensiones (deben ser iguales para comparar píxel a píxel)
	if t1.Height != t2.Height || t1.Width != t2.Width {
		return 0.0
	}

	coincidencias := 0
	pixelesTotales := 0

	// 2. Recorremos los datos del tensor
	// Nota: Solo comparamos el primer canal [0] ya que es binarizado/gris
	for i := 0; i < len(t1.Data); i++ {
		p1 := t1.Data[i] > 0.5
		p2 := t2.Data[i] > 0.5

		// Si al menos uno de los dos tiene tinta en esa posición
		if p1 || p2 {
			pixelesTotales++
			// Si AMBOS tienen tinta en la misma posición exacta
			if p1 && p2 {
				coincidencias++
			}
		}
	}

	// Evitar división por cero si ambas imágenes están vacías
	if pixelesTotales == 0 {
		return 0.0
	}

	// 3. Retornamos el porcentaje de similitud
	return (float64(coincidencias) / float64(pixelesTotales)) * 100
}
