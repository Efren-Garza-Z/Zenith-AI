package core

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

// Tensor4D representa un volumen de datos [Batch, Channels, Height, Width]
type Tensor4D struct {
	Data     []float64
	Batch    int
	Channels int
	Height   int
	Width    int
	// Strides para navegación rápida en el array plano
	sB, sC, sH int
}

// NewTensor4D inicializa la estructura y calcula los saltos de memoria (strides)
func NewTensor4D(b, c, h, w int) *Tensor4D {
	sH := w
	sC := h * w
	sB := c * h * w
	return &Tensor4D{
		Data:     make([]float64, b*c*h*w),
		Batch:    b,
		Channels: c,
		Height:   h,
		Width:    w,
		sH:       sH,
		sC:       sC,
		sB:       sB,
	}
}

// At devuelve el valor en la coordenada específica
func (t *Tensor4D) At(b, c, h, w int) float64 {
	return t.Data[b*t.sB+c*t.sC+h*t.sH+w]
}

// Set asigna un valor en la coordenada específica
func (t *Tensor4D) Set(b, c, h, w int, val float64) {
	t.Data[b*t.sB+c*t.sC+h*t.sH+w] = val
}

// saveTensorToImage es una utilidad para visualizar lo que la red neuronal "ve"
func SaveTensorToImage(filename string, t *Tensor4D) {
	img := image.NewRGBA(image.Rect(0, 0, t.Width, t.Height))

	for y := 0; y < t.Height; y++ {
		for x := 0; x < t.Width; x++ {
			// Tomamos los valores de los 3 canales (RGB)
			// Usamos Abs() porque los gradientes de Sobel pueden ser negativos
			r := uint8(clamp(t.At(0, 0, y, x) * 255))
			g := uint8(clamp(t.At(0, 1, y, x) * 255))
			b := uint8(clamp(t.At(0, 2, y, x) * 255))

			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	f, _ := os.Create(filename)
	defer f.Close()
	png.Encode(f, img)
}

// clamp asegura que el valor esté entre 0 y 255
func clamp(v float64) float64 {
	if v < 0 {
		return -v
	} // Para Sobel, el valor absoluto nos da el borde
	if v > 255 {
		return 255
	}
	return v
}
