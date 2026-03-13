package zig_simd

/*
#cgo LDFLAGS: -L${SRCDIR}/../../ -lconv
#include <stdlib.h>

void zig_convolve_fast(double* input, double* output, int h, int w, double* k, int k_size);
void zig_gaussian_blur(double* input, double* output, int h, int w, double* k, int k_size);
void zig_threshold(double* input, double* output, int h, int w, double threshold);
*/
import "C"
import (
	"unsafe"

	"github.com/efren-garza-z/zenith-ai/core"
)

// ConvolveWithZig usa el motor de Zig optimizado con instrucciones SIMD
func ConvolveWithZig(input *core.Tensor4D, kernel [][]float64) *core.Tensor4D {
	kSize := len(kernel)
	output := core.NewTensor4D(input.Batch, input.Channels, input.Height, input.Width)

	// Aplanamos el kernel (3x3 -> 9 elementos)
	flatKernel := make([]float64, kSize*kSize)
	for i := range kernel {
		copy(flatKernel[i*kSize:], kernel[i])
	}

	// Procesamos cada canal de la imagen
	pixelsPerChannel := input.Height * input.Width
	for b := 0; b < input.Batch; b++ {
		for c := 0; c < input.Channels; c++ {
			offset := (b * input.Channels * pixelsPerChannel) + (c * pixelsPerChannel)

			C.zig_convolve_fast(
				(*C.double)(unsafe.Pointer(&input.Data[offset])),
				(*C.double)(unsafe.Pointer(&output.Data[offset])),
				C.int(input.Height),
				C.int(input.Width),
				(*C.double)(unsafe.Pointer(&flatKernel[0])),
				C.int(kSize),
			)
		}
	}
	return output
}

func GaussianBlur(input *core.Tensor4D) *core.Tensor4D {
	// Kernel Gaussiano 5x5 estándar
	kernel := [][]float64{
		{1 / 256.0, 4 / 256.0, 6 / 256.0, 4 / 256.0, 1 / 256.0},
		{4 / 256.0, 16 / 256.0, 24 / 256.0, 16 / 256.0, 4 / 256.0},
		{6 / 256.0, 24 / 256.0, 36 / 256.0, 24 / 256.0, 6 / 256.0},
		{4 / 256.0, 16 / 256.0, 24 / 256.0, 16 / 256.0, 4 / 256.0},
		{1 / 256.0, 4 / 256.0, 6 / 256.0, 4 / 256.0, 1 / 256.0},
	}

	// Aquí reutilizamos la lógica de pasar el tensor que ya tenemos
	// pero llamando a C.zig_gaussian_blur
	return applyZigFilter(input, kernel, "blur")
}

// applyZigFilter es nuestra función maestra que maneja la lógica de punteros y offsets
func applyZigFilter(input *core.Tensor4D, kernel [][]float64, filterType string) *core.Tensor4D {
	kSize := len(kernel)
	output := core.NewTensor4D(input.Batch, input.Channels, input.Height, input.Width)

	// 1. Aplanamos el kernel
	flatKernel := make([]float64, kSize*kSize)
	for i := range kernel {
		copy(flatKernel[i*kSize:], kernel[i])
	}

	// 2. Procesamos cada canal
	pixelsPerChannel := input.Height * input.Width
	for b := 0; b < input.Batch; b++ {
		for c := 0; c < input.Channels; c++ {
			offset := (b * input.Channels * pixelsPerChannel) + (c * pixelsPerChannel)

			// 3. Decidimos qué función de C llamar
			inputPtr := (*C.double)(unsafe.Pointer(&input.Data[offset]))
			outputPtr := (*C.double)(unsafe.Pointer(&output.Data[offset]))
			kernelPtr := (*C.double)(unsafe.Pointer(&flatKernel[0]))

			if filterType == "blur" {
				C.zig_gaussian_blur(inputPtr, outputPtr, C.int(input.Height), C.int(input.Width), kernelPtr, C.int(kSize))
			} else {
				C.zig_convolve_fast(inputPtr, outputPtr, C.int(input.Height), C.int(input.Width), kernelPtr, C.int(kSize))
			}
		}
	}
	return output
}

// Threshold convierte la imagen a blanco y negro puro (Binarización)
// El valor 'limit' suele estar entre 0.1 y 0.5 para detección de bordes
func Threshold(input *core.Tensor4D, limit float64) *core.Tensor4D {
	output := core.NewTensor4D(input.Batch, input.Channels, input.Height, input.Width)
	pixelsPerChannel := input.Height * input.Width

	for b := 0; b < input.Batch; b++ {
		for c := 0; c < input.Channels; c++ {
			offset := (b * input.Channels * pixelsPerChannel) + (c * pixelsPerChannel)

			C.zig_threshold(
				(*C.double)(unsafe.Pointer(&input.Data[offset])),
				(*C.double)(unsafe.Pointer(&output.Data[offset])),
				C.int(input.Height),
				C.int(input.Width),
				C.double(limit),
			)
		}
	}
	return output
}
